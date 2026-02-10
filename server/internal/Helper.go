package server

import (
	"context"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/encrypt"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 根据配置文件初始化 slog 日志
func InitLogger(conf *g.Config) *slog.Logger {
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

	handler := getLogHandler(conf.Log.Format, option)
	logger := slog.New(handler)
	slog.SetDefault(logger)
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
func getLogHandler(format string, option *slog.HandlerOptions) slog.Handler {
	if format == "json" {
		return slog.NewJSONHandler(os.Stdout, option)
	}
	return slog.NewTextHandler(os.Stdout, option)
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
