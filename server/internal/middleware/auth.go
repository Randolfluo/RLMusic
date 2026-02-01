// TODO
// JWTAuth	JWT认证中间件
package middleware

import (
	g "server/internal/global"
	"server/internal/handle"
	"server/internal/model"
	"server/internal/utils/jwt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		db := c.MustGet(g.CtxDB).(*gorm.DB)

		authorization := c.Request.Header.Get("Authorization") //获取token
		if authorization == "" {
			handle.ReturnError(c, g.ErrTokenNotExist, nil)
			return
		}

		// token 的正确格式: `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handle.ReturnError(c, g.ErrTokenType, nil)
			return
		}

		claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
		if err != nil {
			handle.ReturnError(c, g.ErrTokenWrong, err)
			return
		}

		// 判断 token 已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handle.ReturnError(c, g.ErrTokenRuntime, nil)
			return
		}

		user, err := model.GetUserAuthInfoById(db, claims.UserId)
		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
			return
		}

		// gin context
		c.Set(g.CtxUserAuth, user)
	}
}
