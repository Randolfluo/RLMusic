package handle

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/encrypt"
	"server/internal/utils/jwt"
	"strings"
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
	Email     string     `json:"email"`
	Avatar    string     `json:"avatar"`
	Token     string     `json:"token"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	// 校验密码长度 (解密后)
	if len(req.Password) < 4 || len(req.Password) > 20 {
		ReturnError(c, g.ErrRequest, errors.New("密码长度必须在4-20位之间"))
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
	user, err := model.CreateUser(db, req.Username, hashedPassword, req.Email)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 创建默认歌单 "我喜欢的音乐" (私有)
	if _, err := model.CreatePlaylist(db, user.ID, "我喜欢的音乐", "因为热爱，所以收藏", false); err != nil {
		slog.Warn("Failed to create default playlist for user", "username", req.Username, "error", err)
	}

	// 自动创建用户目录
	// 忽略错误，因为这不应该阻止注册流程 (可能系统文件夹还没初始化，或者管理员还没配置路径)
	// 用户可以在后续通过扫描等操作前被提示需要管理员初始化
	if err := CreateUserFolder(db, req.Username); err != nil {
		slog.Warn("Failed to auto-create user folder during registration", "user", req.Username, "error", err)
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
		Email:     user.Email,
		Avatar:    user.Avatar,
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

type UserInfoVO struct {
	ID            int        `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Avatar        string     `json:"avatar"`
	LastLogin     *time.Time `json:"last_login"`
	TotalSongs    int64      `json:"total_songs"`
	TotalAlbums   int64      `json:"total_albums"`
	TotalArtists  int64      `json:"total_artists"`
	TotalDuration int64      `json:"total_duration"`
	FavoriteSong  string     `json:"favorite_song"`
}

// GetUserInfo 获取当前登录用户的详细信息
func (*UserAuth) GetUserInfo(c *gin.Context) {
	// 从Context中获取当前用户（由中间件设置）
	currentUser := GetCurrentUser(c)
	if currentUser == nil {
		ReturnError(c, g.ErrTokenEmpty, nil)
		return
	}

	// 重新从数据库查一遍，确保信息是最新的
	db := GetDB(c)
	user, err := model.GetUserAuthInfoByName(db, currentUser.Username)
	if err != nil {
		ReturnError(c, g.ErrUserNotExist, err)
		return
	}

	sysInfo, _ := model.GetSystemInfoStruct(db)
	if sysInfo == nil {
		sysInfo = &model.SystemInfoStruct{}
	}

	ReturnSuccess(c, UserInfoVO{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Avatar:        user.Avatar,
		LastLogin:     user.LastLogin,
		TotalSongs:    sysInfo.TotalSongs,
		TotalAlbums:   sysInfo.TotalAlbums,
		TotalArtists:  sysInfo.TotalArtists,
		TotalDuration: user.ListeningDuration, // 保持为用户的听歌时长
		FavoriteSong:  user.FavoriteSong,
	})
}

// UploadAvatar 上传用户头像
func (*UserAuth) UploadAvatar(c *gin.Context) {
	currentUser := GetCurrentUser(c)
	if currentUser == nil {
		ReturnError(c, g.ErrTokenEmpty, nil)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowExts[ext] {
		ReturnError(c, g.ErrRequest, errors.New("file type not allowed"))
		return
	}

	// 准备存储目录
	conf := g.GetConfig().BasicPath
	baseDir := filepath.Join(conf.FilePath, conf.FileName, "avatar")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 生成文件名: user_id_timestamp.ext
	filename := fmt.Sprintf("user_%d_%d%s", currentUser.ID, time.Now().Unix(), ext)
	savePath := filepath.Join(baseDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 相对 web 路径 (注意: 这里返回的路径需要和下面注册的路由匹配)
	// 假设路由注册为 /file/avatar/:filename (无需认证读取)
	// 这里我们暂定为 /api/file/avatar/:filename
	webPath := "/api/file/avatar/" + filename

	db := GetDB(c)
	if err := model.UpdateUserAvatar(db, currentUser.ID, webPath); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, webPath)
}

// GetAvatar 获取用户头像
func (*UserAuth) GetAvatar(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.Status(http.StatusNotFound)
		return
	}

	// 安全检查
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		c.Status(http.StatusForbidden)
		return
	}

	conf := g.GetConfig().BasicPath
	filePath := filepath.Join(conf.FilePath, conf.FileName, "avatar", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(filePath)
}
