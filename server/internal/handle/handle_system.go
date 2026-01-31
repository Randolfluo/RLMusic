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

// GetSettings 获取所有系统设置
func (*SystemAuth) GetSettings(c *gin.Context) {
	db := GetDB(c)
	settings, err := model.GetAllSystemSettings(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, settings)
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
