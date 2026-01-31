// 注入数据库连接到 gin.Context 的中间件
package middleware

import (
	g "server/internal/global"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// WithRedisDB 将 redis.Client 注入到 gin.Context
// handler 中通过 c.MustGet(g.CtxRedis).(*redis.Client) 来使用
func WithRedisDB(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(g.CtxRedis, rdb)
		ctx.Next()
	}
}

// WithGormDB 将 gorm.DB 注入到 gin.Context
// handler 中通过 c.MustGet(g.CtxDB).(*gorm.DB) 来使用
func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(g.CtxDB, db)
		ctx.Next()
	}
}
