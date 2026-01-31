package handle

import (
	"log/slog"
	"net/http"
	g "server/internal/global"
	"server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Response 通用响应结构体
type Response[T any] struct {
	Code    int    `json:"code"`    // 业务状态码,0表示成功,其他表示失败
	Message string `json:"message"` // 响应消息
	Data    T      `json:"data"`    // 响应数据
}

// ReturnHttpResponse 返回HTTP响应
func ReturnHttpResponse(c *gin.Context, httpCode, code int, msg string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// ReturnResponse 返回业务响应
func ReturnResponse(c *gin.Context, r g.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

// ReturnSuccess 返回成功响应
func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, g.Success, data)
}

// ReturnError 返回错误响应
// 所有可预料的错误(业务错误+系统错误)在业务层面处理,返回HTTP 200
// 不可预料的错误会触发panic,由gin中间件捕获并返回HTTP 500
func ReturnError(c *gin.Context, r g.Result, data any) {
	slog.Info("[ReturnError] " + r.Msg())

	errMsg := r.Msg()
	if data != nil {
		switch v := data.(type) {
		case error:
			errMsg = v.Error()
		case string:
			errMsg = v
		}
		slog.Error(errMsg)
	}

	c.AbortWithStatusJSON(http.StatusOK, Response[any]{
		Code:    r.Code(),
		Message: r.Msg(),
		Data:    errMsg,
	})
}

// GetDB 获取数据库连接
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CtxDB).(*gorm.DB)
}

// GetRDB 获取Redis连接
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(g.CtxRedis).(*redis.Client)
}

// GetCurrentUser 获取当前登录用户
func GetCurrentUser(c *gin.Context) *model.User {
	return c.MustGet(g.CtxUserAuth).(*model.User)
}
