package handle

import (
	"log/slog"
	"os"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"

	"github.com/gin-gonic/gin"
)

type FileAuth struct{}

// InitFolder 在 basicPath 下创建文件夹
func (*FileAuth) InitFolder(c *gin.Context) {
	// 检查是否已初始化 (持久化检查)
	db := GetDB(c)
	const KeyBaseFolderCreated = "is_base_folder_created"
	if val, _ := model.GetSystemSetting(db, KeyBaseFolderCreated); val == "true" {
		ReturnSuccess(c, "基础文件夹已初始化(已跳过创建)")
		return
	}

	conf := g.GetConfig().BasicPath
	if conf.FilePath == "" || conf.FileName == "" {
		slog.Error("BasicPath not configured")
		ReturnError(c, g.Err, "服务器配置错误: BasicPath 未配置")
		return
	}

	// 拼接完整路径
	fullPath := filepath.Join(conf.FilePath, conf.FileName)

	//检查文件夹是否存在
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		// 如果文件夹存在，也视为初始化成功，并写入数据库状态
		_ = model.SetSystemSetting(db, KeyBaseFolderCreated, "true")
		ReturnSuccess(c, "文件夹已存在")
		return
	}

	// 创建文件夹
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		slog.Error("Failed to create folder", "path", fullPath, "error", err)
		ReturnError(c, g.Err, "创建文件夹失败: "+err.Error())
		return
	}

	// 记录初始化状态
	if err := model.SetSystemSetting(db, KeyBaseFolderCreated, "true"); err != nil {
		slog.Error("Failed to save system setting", "error", err)
		// 不阻断流程，仅记录日志
	}

	ReturnSuccess(c, "文件夹创建成功")
}

type InitUserFolderReq struct {
	Username string `json:"username" binding:"required"`
}

// InitUserFolder 初始化用户目录
func (*FileAuth) InitUserFolder(c *gin.Context) {
	var req InitUserFolderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	user, err := model.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		ReturnError(c, g.ErrUserNotExist, err)
		return
	}

	if user.IsCreateFile {
		ReturnSuccess(c, "该用户目录已初始化")
		return
	}

	// 检查基础文件夹是否已初始化
	const KeyBaseFolderCreated = "is_base_folder_created"
	if val, _ := model.GetSystemSetting(db, KeyBaseFolderCreated); val != "true" {
		ReturnError(c, g.Err, "基础文件夹尚未初始化，请先初始化基础文件夹")
		return
	}

	conf := g.GetConfig().BasicPath
	if conf.FilePath == "" || conf.FileName == "" {
		slog.Error("BasicPath not configured")
		ReturnError(c, g.Err, "服务器配置错误: BasicPath 未配置")
		return
	}

	// 用户主目录: 基础路径/基础文件夹名/用户名
	// 确保在基础文件夹(localmusicplayerFoleder)下创建
	userPath := filepath.Join(conf.FilePath, conf.FileName, req.Username)

	// 需要创建的子目录
	// "" 代表用户主目录本身
	subDirs := []string{"", "setting", "public", "private"}

	// 创建主目录和子目录
	for _, dir := range subDirs {
		fullPath := filepath.Join(userPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			slog.Error("Failed to create folder", "path", fullPath, "error", err)
			ReturnError(c, g.Err, "创建目录失败: "+err.Error())
			return
		}
	}

	// 更新用户状态
	if err := model.UpdateUserCreateFileStatus(db, req.Username); err != nil {
		slog.Error("Failed to update user status", "error", err)
		ReturnError(c, g.ErrDbOp, "更新用户状态失败")
		return
	}

	ReturnSuccess(c, "用户目录初始化成功")
}
