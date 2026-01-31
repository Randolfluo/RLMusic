package handle

import (
	"errors"
	"log/slog"
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/encrypt"
	"server/internal/utils/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserAuth struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginVO struct {
	ID        int        `json:"id"`
	LastLogin *time.Time `json:"last_login"`
	Username  string     `json:"username"`
	Token     string     `json:"token"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	Email    string `json:"email" binding:"required"`
	// EmailCaptcha string `json:"emailCaptcha" binding:"required"`
}

type RegisterVO struct{}

type CaptchaVO struct {
	Id   string `json:"id"`
	B64s string `json:"b64s"`
}

type GetEmailCaptchaReq struct {
	Email string `json:"email" binding:"required,min=5,max=254"`
}

type GetEmailCaptchaVO struct{}

// Register 用户注册
func (*UserAuth) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	db := GetDB(c)

	// 检查用户是否已存在
	if _, err := model.GetUserAuthInfoByName(db, req.Username); err == nil {
		ReturnError(c, g.ErrUserExist, nil)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 检查是否有由于软删除占用用户名的记录，如果有则重命名该旧记录，释放用户名
	if err := model.RenameSoftDeletedUser(db, req.Username); err != nil {
		slog.Error("Failed to rename old deleted user", "username", req.Username, "error", err)
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 密码加密
	hashedPassword, err := encrypt.BcryptHash(req.Password)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 创建用户
	if _, err := model.CreateUser(db, req.Username, hashedPassword, req.Email); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, RegisterVO{})
}

// Login 用户登录
func (*UserAuth) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 获取用户信息
	user, err := model.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 验证密码
	if !encrypt.BcryptCheck(req.Password, user.Password) {
		ReturnError(c, g.ErrPassword, nil)
		return
	}

	// 生成JWT令牌
	conf := g.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), user.ID)
	if err != nil {
		ReturnError(c, g.ErrTokenCreate, err)
		return
	}

	// 更新登录信息
	if err := model.UpdateUserLoginInfo(db, user.ID); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	slog.Info("用户登录成功", "username", user.Username)
	ReturnSuccess(c, LoginVO{
		ID:        user.ID,
		Username:  user.Username,
		LastLogin: user.LastLogin,
		Token:     token,
	})
}

// Logout 用户退出登录
func (*UserAuth) Logout(c *gin.Context) {
	// 目前使用JWT无状态认证，服务端无需特殊处理，客户端丢弃Token即可
	// 后续如果引入Redis，可在此处将Token加入黑名单
	ReturnSuccess(c, nil)
}

// DeleteUser 用户注销
func (*UserAuth) DeleteUser(c *gin.Context) {
	db := GetDB(c)
	user := GetCurrentUser(c)

	if err := model.DeleteUser(db, user.ID); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}
