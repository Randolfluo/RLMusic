package server

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/encrypt"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ringBufferHandler 自定义 slog Handler，拦截日志写入环形缓冲区
type ringBufferHandler struct {
	inner slog.Handler
}

func (h *ringBufferHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *ringBufferHandler) Handle(ctx context.Context, r slog.Record) error {
	msg := r.Message
	r.Attrs(func(a slog.Attr) bool {
		msg += " " + a.Key + "=" + a.Value.String()
		return true
	})
	g.LogBuffer.Append(g.LogEntry{
		Time:    r.Time.Format(time.DateTime),
		Level:   r.Level.String(),
		Message: msg,
	})
	return h.inner.Handle(ctx, r)
}

func (h *ringBufferHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ringBufferHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *ringBufferHandler) WithGroup(name string) slog.Handler {
	return &ringBufferHandler{inner: h.inner.WithGroup(name)}
}

// stdLogBridge 将标准库 log 输出桥接到 slog
type stdLogBridge struct{}

func (stdLogBridge) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		slog.Info(msg)
	}
	return len(p), nil
}

// ginLogBridge 双写：原始输出保留 ANSI 颜色，环形缓冲区捕获纯文本
type ginLogBridge struct {
	original io.Writer
}

func (b ginLogBridge) Write(p []byte) (n int, err error) {
	// 写入原始目标（终端保留 ANSI 彩色）
	b.original.Write(p)
	// 直接写入环形缓冲区（绕过 slog，避免 stdout 重复输出）
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		g.LogBuffer.Append(g.LogEntry{
			Time:    time.Now().Format(time.DateTime),
			Level:   "INFO",
			Message: msg,
		})
	}
	return len(p), nil
}

// 根据配置文件初始化 slog 日志
func InitLogger(conf *g.Config) *slog.Logger {
	g.LogBuffer = g.NewLogRingBuffer(5000)

	level := getLogLevel(conf.Log.Level)

	option := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}

	handler := getLogHandler(conf, option)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// 将标准库 log 输出重定向到 slog，使 log.Printf/Println 也进入环形缓冲区
	log.SetOutput(&stdLogBridge{})
	// 将 Gin debug 模式的路由列表等输出也重定向到 slog
	gin.DefaultWriter = &ginLogBridge{original: gin.DefaultWriter}

	return logger
}

// 获取日志级别
func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// 获取日志处理器
func getLogHandler(conf *g.Config, option *slog.HandlerOptions) slog.Handler {
	var writer io.Writer = os.Stdout

	if conf.Log.Directory != "" {
		if err := os.MkdirAll(conf.Log.Directory, os.ModePerm); err == nil {
			logFileName := filepath.Join(conf.Log.Directory, "server-"+time.Now().Format(time.DateOnly)+".log")
			logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				writer = io.MultiWriter(os.Stdout, logFile)
			}
		}
	}

	var inner slog.Handler
	if conf.Log.Format == "json" {
		inner = slog.NewJSONHandler(writer, option)
	} else {
		inner = slog.NewTextHandler(writer, option)
	}

	return &ringBufferHandler{inner: inner}
}

// 根据配置文件初始化数据库
func InitDatabase(conf *g.Config) *gorm.DB {
	dbtype := conf.DbType()
	dsn := conf.DbDSN()

	level := getDBLogLevel(conf.Server.DbLogMode)
	config := getGormConfig(level)

	if dbtype == "sqlite" {
		if err := os.MkdirAll(filepath.Dir(dsn), os.ModePerm); err != nil {
			log.Panic("创建数据库目录失败:", err)
		}
	}

	db, err := openDatabase(dbtype, dsn, config)
	if err != nil {
		log.Panic("数据库连接失败:", err)
	}

	log.Printf("数据库连接成功 类型:%s DSN:%s", dbtype, dsn)

	if conf.Server.DbAutoMigrate {
		if err := model.MakeMigrate(db); err != nil {
			log.Fatal("数据库迁移失败:", err)
		}
		log.Println("数据库自动迁移成功")
		initAdminUser(db)
	}
	return db
}

// 获取数据库日志级别
func getDBLogLevel(mode string) logger.LogLevel {
	switch mode {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	default:
		return logger.Error
	}
}

// 获取 Gorm 配置
func getGormConfig(level logger.LogLevel) *gorm.Config {
	return &gorm.Config{
		Logger:                                   logger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

// 打开数据库连接
func openDatabase(dbtype string, dsn string, config *gorm.Config) (*gorm.DB, error) {
	switch dbtype {
	// case "mysql":
	// 	return gorm.Open(mysql.Open(dsn), config)
	case "sqlite":
		return gorm.Open(sqlite.Open(dsn), config)
	default:
		log.Panic("不支持的数据库类型:", dbtype)
		return nil, nil
	}
}

func initAdminUser(db *gorm.DB) {
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		log.Printf("检查管理员账号失败: %v", err)
		return
	}

	if count == 0 {
		hash, err := encrypt.BcryptHash("123456")
		if err != nil {
			log.Printf("管理员密码加密失败: %v", err)
			return
		}
		admin := model.User{
			Username:  "admin",
			Password:  hash,
			UserGroup: "admin",
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("管理员账号创建失败: %v", err)
		} else {
			log.Println("管理员账号创建成功: admin/123456")
		}
	}
}

// 根据配置文件初始化 Redis
func InitRedis(conf *g.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis 连接失败:", err)
	}

	log.Printf("Redis 连接成功 地址:%s DB:%d", conf.Redis.Addr, conf.Redis.DB)
	return rdb
}
