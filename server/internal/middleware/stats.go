package middleware

import (
	g "server/internal/global"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// StatsMiddleware 统计API调用次数
func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除获取系统状态的接口，避免轮询造成统计数据虚高
		if c.Request.URL.Path != "/api/system/status" {
			atomic.AddInt64(&g.ApiCallCount, 1)
		}
		c.Next()
	}
}
