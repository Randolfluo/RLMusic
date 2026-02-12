package g

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper" // 配置文件读取包
)

// 全局统计变量
var (
	ApiCallCount int64 // API调用次数
)

type Config struct {
	Server struct {
		Mode          string // debug | release
		Port          string
		DbType        string // mongodb | sqlite
		DbAutoMigrate bool   // 是否自动迁移数据库表结构
		DbLogMode     string // silent | error | warn | info
	}
	Log struct {
		Level     string // debug | info | warn | error
		Prefix    string
		Format    string // text | json
		Directory string
	}
	JWT struct { //JWT鉴权
		Secret string
		Expire int64
		Issuer string
	}
	SQLite struct {
		Dsn string
	}
	Redis struct {
		DB       int    // 指定 Redis 数据库
		Addr     string // 服务器地址:端口
		Password string // 密码
	}
	Session struct {
		Name   string
		Salt   string
		MaxAge int
	}
	BasicPath struct {
		FilePath string
		FileName string
	}
	Email struct {
		Host       string // 服务器地址, 例如 smtp.qq.com 前往要发邮件的邮箱查看其 smtp 协议
		Port       int    // 前往要发邮件的邮箱查看其 smtp 协议端口, 大多为 465
		From       string // 发件人 要发邮件的邮箱
		IsSSL      bool   // 是否开启 SSL
		Secret     string // 密钥, 不是邮箱登录密码, 是开启 smtp 服务后获取的一串验证码
		Nickname   string // 发件人昵称, 通常为自己的邮箱名
		CaptchaLen int    //数字验证码长度
	}
	Captcha struct {
		SendEmail  bool // 是否通过邮箱发送验证码
		ExpireTime int  // 过期时间
	}
	CaptchaSetting struct {
		Height   int
		Width    int
		Length   int
		MaxSkew  float64
		DotCount int
	}
	QwenTTS struct {
		ApiKey     string  `mapstructure:"api_key"`
		Model      string  `mapstructure:"model"`
		Voice      string  `mapstructure:"voice"`
		Format     string  `mapstructure:"format"`      // mp3, wav, pcm, flac
		SampleRate int     `mapstructure:"sample_rate"` // 22050, 44100, 48000
		Volume     int     `mapstructure:"volume"`      // 0-100
		Rate       float64 `mapstructure:"rate"`        // 0.5-2.0
		Pitch      float64 `mapstructure:"pitch"`       // 0.5-2.0
	}
	SiliconFlow struct {
		ApiKey string `mapstructure:"api_key"`
		Model  string `mapstructure:"model"`
	}
}

var Conf *Config
var StartTime time.Time

// 从指定路径读取配置文件
func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()                                   // 允许使用环境变量
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // SERVER_APPMODE => SERVER.APPMODE

	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}

	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败: " + err.Error())
	}

	log.Println("配置文件内容加载成功: ", path)
	return Conf
}

// 获取全局配置对象
func GetConfig() *Config {
	if Conf == nil {
		log.Panic("配置文件未初始化")
		return nil
	}
	return Conf
}

// 数据库类型
func (*Config) DbType() string {
	if Conf.Server.DbType == "" {
		log.Panic("未设置数据库类型")
	}
	return Conf.Server.DbType
}

// 数据库连接字符串
func (*Config) DbDSN() string {
	switch Conf.Server.DbType {
	// case "mysql":
	// 	conf := Conf.Mysql
	// 	return fmt.Sprintf(
	// 		"%s:%s@tcp(%s:%s)/%s?%s",
	// 		conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Config,
	// 	)
	case "sqlite":
		return filepath.Join(Conf.BasicPath.FilePath, Conf.BasicPath.FileName, Conf.SQLite.Dsn)
	default:
		log.Panic("未设置数据库连接字符串")
		return ""
	}
}
