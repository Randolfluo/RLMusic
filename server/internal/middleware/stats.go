package middleware

import (
	g "server/internal/global"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// StatsMiddleware 统计API调用次数
func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		atomic.AddInt64(&g.ApiCallCount, 1)
		c.Next()
	}
}
