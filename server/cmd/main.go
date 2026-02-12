package main

import (
	"flag"
	"log"
	server "server/internal"
	g "server/internal/global"
	"server/internal/middleware"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	g.StartTime = time.Now()
	configPath := flag.String("c", "config.yml", "配置文件路径")
	flag.Parse()
	conf := g.ReadConfig(*configPath)

	_ = server.InitLogger(conf)     // 初始化日志
	db := server.InitDatabase(conf) // 初始化数据库
	// rdb := server.InitRedis(conf)   // 初始化缓存

	gin.SetMode(conf.Server.Mode)
	r := gin.New()
	r.SetTrustedProxies([]string{"*"})

	r.Use(gin.Logger(), gin.Recovery())

	//r.Use(middleware.CORS())
	r.Use(middleware.WithGormDB(db))
	// r.Use(middleware.WithRedisDB(rdb))
	//r.Use(middleware.WithCookieStore(conf.Session.Name, conf.Session.Salt))

	// 静态资源: 封面图
	r.Static("/covers", "./data/covers")

	server.RegisterHandlers(r, nil)

	//运行服务
	serverAddr := conf.Server.Port
	if serverAddr[0] == ':' || strings.HasPrefix(serverAddr, "0.0.0.0:") {
		log.Printf("Serving HTTP on (http://localhost:%s/) ... \n", strings.Split(serverAddr, ":")[1])
	} else {
		log.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	}
	r.Run(serverAddr)
}
