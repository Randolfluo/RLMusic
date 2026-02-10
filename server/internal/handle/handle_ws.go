package handle

import (
	"server/ws"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type WSAuth struct {
	Server *ws.WSServer
}

func NewWSAuth(rdb *redis.Client) *WSAuth {
	return &WSAuth{Server: ws.NewWSServer(rdb)} // 初始化 WebSocket 服务器
}

// HandleWS 处理 WebSocket 连接
func (auth *WSAuth) HandleWS(c *gin.Context) {
	auth.Server.M.HandleRequest(c.Writer, c.Request)
}
