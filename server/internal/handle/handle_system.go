package handle

import (
	g "server/internal/global"
	"server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type SystemAuth struct{}

type UpdateConfigReq struct {
	FilePath string `json:"filepath"`
}

type DurationStatsVO struct {
	UserTotalDuration   int64 `json:"user_total_duration"`
	SystemTotalDuration int64 `json:"system_total_duration"`
}

// GetSettings 获取所有系统设置
func (*SystemAuth) GetSettings(c *gin.Context) {
	db := GetDB(c)
	info, err := model.GetSystemInfoStruct(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, info)
}

// GetDurationStats 获取用户和系统总时长
func (*SystemAuth) GetDurationStats(c *gin.Context) {
	db := GetDB(c)
	// 获取当前登录用户
	authUser := GetCurrentUser(c)
	if authUser == nil {
		ReturnError(c, g.ErrTokenEmpty, nil)
		return
	}

	user, err := model.GetUserAuthInfoByName(db, authUser.Username)
	if err != nil {
		ReturnError(c, g.ErrUserNotExist, err)
		return
	}

	info, err := model.GetSystemInfoStruct(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 容错处理
	if info == nil {
		info = &model.SystemInfoStruct{}
	}

	ReturnSuccess(c, DurationStatsVO{
		UserTotalDuration:   user.TotalDuration,
		SystemTotalDuration: info.TotalDuration,
	})
}

// UpdateConfig 更新系统配置
func (*SystemAuth) UpdateConfig(c *gin.Context) {
	// 1. Check Auth (Admin only)
	userVal, exists := c.Get(g.CtxUserAuth)
	if !exists {
		ReturnError(c, g.ErrTokenRuntime, nil)
		return
	}
	user := userVal.(*model.User)
	if user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, nil)
		return
	}

	// 2. Bind JSON
	var req UpdateConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if req.FilePath == "" {
		ReturnError(c, g.ErrRequest, "File path cannot be empty")
		return
	}

	// 3. Update Viper Config
	v := viper.New()
	v.SetConfigFile("config.yml")
	if err := v.ReadInConfig(); err != nil {
		ReturnError(c, g.Err, "Read config failed: "+err.Error())
		return
	}

	v.Set("BasicPath.FilePath", req.FilePath)
	if err := v.WriteConfig(); err != nil {
		ReturnError(c, g.Err, "Write config failed: "+err.Error())
		return
	}

	// 4. Update Memory Config
	g.Conf.BasicPath.FilePath = req.FilePath

	ReturnSuccess(c, nil)
}
