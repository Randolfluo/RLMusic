package main

import (
	"flag"
	"log"
	"path/filepath"
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

	r.Use(middleware.CORS())
	r.Use(middleware.WithGormDB(db))
	// r.Use(middleware.WithRedisDB(rdb))
	//r.Use(middleware.WithCookieStore(conf.Session.Name, conf.Session.Salt))

	// 静态资源: 封面图
	coverPath := filepath.Join(conf.BasicPath.FilePath, conf.BasicPath.FileName, "data", "cover")
	r.Static("/covers", coverPath)

	// 静态资源: 播客开场白 (QwenTTS 生成的音频)
	podcastPath := filepath.Join(conf.BasicPath.FilePath, conf.BasicPath.FileName, "data", "Podcast")
	r.Static("/podcast", podcastPath)

	server.RegisterHandlers(r, nil)

	//运行服务
	serverAddr := strings.TrimSpace(conf.Server.Port)
	if serverAddr == "" {
		serverAddr = ":12345"
	}
	if !strings.Contains(serverAddr, ":") {
		serverAddr = ":" + serverAddr
	}
	if strings.HasPrefix(serverAddr, "localhost:") {
		serverAddr = "0.0.0.0:" + strings.TrimPrefix(serverAddr, "localhost:")
	} else if strings.HasPrefix(serverAddr, "127.0.0.1:") {
		serverAddr = "0.0.0.0:" + strings.TrimPrefix(serverAddr, "127.0.0.1:")
	} else if strings.HasPrefix(serverAddr, ":") {
		serverAddr = "0.0.0.0" + serverAddr
	}
	if strings.HasPrefix(serverAddr, "0.0.0.0:") {
		log.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	} else {
		log.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	}
	r.Run(serverAddr)
}
