package handle

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	g "server/internal/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileAuth struct{}

// InitFolder 在 basicPath 下创建文件夹
func (*FileAuth) InitFolder(c *gin.Context) {
	conf := g.GetConfig().BasicPath
	if conf.FilePath == "" {
		slog.Error("BasicPath not configured")
		ReturnError(c, g.Err, "服务器配置错误: BasicPath 未配置")
		return
	}

	// 拼接完整路径
	fullPath := filepath.Join(conf.FilePath, conf.FileName)

	//检查文件夹是否存在
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		// 如果文件夹存在，也视为初始化成功
		ReturnSuccess(c, "文件夹已存在")
		return
	}

	// 创建文件夹
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		slog.Error("Failed to create folder", "path", fullPath, "error", err)
		ReturnError(c, g.Err, "创建文件夹失败: "+err.Error())
		return
	}

	ReturnSuccess(c, "文件夹创建成功")
}

// CreateUserFolder 内部调用：创建用户目录
func CreateUserFolder(db *gorm.DB, username string) error {
	conf := g.GetConfig().BasicPath
	if conf.FilePath == "" {
		return errors.New("服务器配置错误: BasicPath 未配置")
	}

	// 用户主目录: 基础路径/基础文件夹名/用户名
	userPath := filepath.Join(conf.FilePath, conf.FileName, username)

	// 创建用户主目录
	if err := os.MkdirAll(userPath, 0755); err != nil {
		slog.Error("Failed to create user folder", "path", userPath, "error", err)
		return err
	}

	return nil
}
