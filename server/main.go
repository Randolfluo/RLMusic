package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 实例
	r := gin.Default()

	// 配置 CORS 中间件（如果 vite 代理配置正确，生产环境可能需要，开发环境有代理可跳过，但加上更稳妥）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 示例：首页每日推荐接口
	// 前端请求：/api/recommend/songs
	// Vite 代理重写后转发至：/recommend/songs
	r.GET("/recommend/songs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"dailySongs": []gin.H{
					{
						"id":   1,
						"name": "Golang 后端测试歌曲",
						"al": gin.H{
							"picUrl": "https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg",
						},
						"ar": []gin.H{
							{"name": "Golang"},
						},
					},
				},
			},
		})
	})

	// 启动服务，监听 3000 端口
	// 请确保 .env 文件中 VITE_MUSIC_API = "http://localhost:3000"
	r.Run(":3000")
}
