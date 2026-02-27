package main

import (
	"flag"
	"log"
	"net"
	"net/url"
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
	defaultPort := "12345"
	serverAddr := strings.TrimSpace(conf.Server.Port)
	if serverAddr == "" {
		serverAddr = ":" + defaultPort
	}
	if strings.Contains(serverAddr, "://") {
		if parsed, err := url.Parse(serverAddr); err == nil && parsed.Host != "" {
			serverAddr = parsed.Host
		}
	}
	serverAddr = strings.TrimSpace(serverAddr)
	if serverAddr == "localhost" || serverAddr == "127.0.0.1" || serverAddr == "0.0.0.0" {
		serverAddr = serverAddr + ":" + defaultPort
	}
	isPortOnly := true
	for _, r := range serverAddr {
		if r < '0' || r > '9' {
			isPortOnly = false
			break
		}
	}
	if isPortOnly {
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
		// 获取本机局域网IP
		ips := []string{}
		if interfaces, err := net.Interfaces(); err == nil {
			for _, i := range interfaces {
				if addrs, err := i.Addrs(); err == nil {
					for _, addr := range addrs {
						var ip net.IP
						switch v := addr.(type) {
						case *net.IPNet:
							ip = v.IP
						case *net.IPAddr:
							ip = v.IP
						}
						if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
							ips = append(ips, ip.String())
						}
					}
				}
			}
		}

		port := strings.TrimPrefix(serverAddr, "0.0.0.0:")
		log.Printf("Serving HTTP on:\n")
		log.Printf("  - Local:   http://localhost:%s/\n", port)
		for _, ip := range ips {
			log.Printf("  - Network: http://%s:%s/\n", ip, port)
		}
	} else {
		log.Printf("Serving HTTP on http://%s/ ... \n", serverAddr)
	}
	r.Run(serverAddr)
}
