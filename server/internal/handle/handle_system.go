package handle

import (
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/jwt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/viper"
)

type SystemAuth struct{}

type UpdateConfigReq struct {
	FilePath string `json:"filepath"`
}

type StatsVO struct {
	SongCount     int64 `json:"song_count"`
	AlbumCount    int64 `json:"album_count"`
	ArtistCount   int64 `json:"artist_count"`
	MusicDuration int64 `json:"music_duration"`

	PlaylistCount int64 `json:"playlist_count"`
	UserCount     int64 `json:"user_count"`
	SystemUptime  int64 `json:"system_uptime"`

	UserListeningDuration int64 `json:"user_listening_duration"`
	UserScannedDuration   int64 `json:"user_scanned_duration"`

	CpuUsage     float64 `json:"cpu_usage"`
	MemUsage     float64 `json:"mem_usage"`
	ApiCallCount int64   `json:"api_call_count"`
}

// GetStats 获取系统统计信息（包含原有 Settings 和 Duration）
func (*SystemAuth) GetStats(c *gin.Context) {
	db := GetDB(c)

	// 1. 获取系统基础信息
	info, err := model.GetSystemInfoStruct(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	if info == nil {
		info = &model.SystemInfoStruct{}
	}

	// 补充统计: 歌单数量
	var totalPlaylists int64
	db.Model(&model.Playlist{}).Count(&totalPlaylists)

	// 补充统计: 用户数量
	var totalUsers int64
	db.Model(&model.User{}).Where("is_delete = ?", false).Count(&totalUsers)

	// 补充统计: 系统运行时间
	systemRunTime := int64(time.Since(g.StartTime).Seconds())

	// 系统负载 (CPU & Mem)
	var cpuUsage float64 = 0
	percent, err := cpu.Percent(0, false)
	if err == nil && len(percent) > 0 {
		cpuUsage = percent[0]
	}
	var memUsage float64 = 0
	vMem, err := mem.VirtualMemory()
	if err == nil {
		memUsage = vMem.UsedPercent
	}

	// API 调用次数
	apiCount := atomic.LoadInt64(&g.ApiCallCount)

	// 2. 获取当前用户时长信息 (可选登录)
	var listeningDuration int64 = 0
	var totalDuration int64 = 0
	authorization := c.GetHeader("Authorization")
	if authorization != "" {
		parts := strings.Split(authorization, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
			if err == nil && time.Now().Unix() <= claims.ExpiresAt.Unix() {
				if user, err := model.GetUserAuthInfoById(db, claims.UserId); err == nil {
					listeningDuration = user.ListeningDuration
					totalDuration = user.TotalDuration
				}
			}
		}
	}

	ReturnSuccess(c, StatsVO{
		SongCount:             info.TotalSongs,
		AlbumCount:            info.TotalAlbums,
		ArtistCount:           info.TotalArtists,
		MusicDuration:         info.TotalDuration,
		PlaylistCount:         totalPlaylists,
		UserCount:             totalUsers,
		SystemUptime:          systemRunTime,
		UserListeningDuration: listeningDuration,
		UserScannedDuration:   totalDuration,
		CpuUsage:              cpuUsage,
		MemUsage:              memUsage,
		ApiCallCount:          apiCount,
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
