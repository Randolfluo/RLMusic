package handle

import (
	"fmt"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/jwt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"net"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

// GetLocalIPs 获取局域网IP地址列表
func (*SystemAuth) GetLocalIPs(c *gin.Context) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		ReturnError(c, g.Err, "获取网络接口失败: "+err.Error())
		return
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 过滤掉回环地址和 IPv6 地址 (目前主要关注 IPv4 局域网访问)
			// 如果需要 IPv6，可以去掉 To4() 的判断
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	// 获取前端传来的端口，如果没有则使用后端配置的端口
	targetPort := c.Query("port")
	if targetPort == "" {
		targetPort = g.GetConfig().Server.Port
	}

	var urls []string
	for _, ip := range ips {
		// 构造 URL: http://IP:Port/
		urls = append(urls, fmt.Sprintf("http://%s:%s/", ip, targetPort))
	}

	ReturnSuccess(c, gin.H{
		"ips":  ips,
		"port": targetPort,
		"urls": urls,
	})
}

// ExportDatabaseToExcel 导出数据库到Excel
func (*SystemAuth) ExportDatabaseToExcel(c *gin.Context) {
	// 1. Check Admin Permission
	user := GetCurrentUser(c)
	if user == nil || user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "权限不足，仅管理员可执行")
		return
	}

	db := GetDB(c)
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			// log error
		}
	}()

	// 2. Export Users
	var users []model.User
	if err := db.Find(&users).Error; err == nil {
		sheetName := "Users"
		f.NewSheet(sheetName)
		headers := []string{"ID", "Username", "Nickname", "Email", "UserGroup", "TotalDuration", "ListeningDuration", "CreatedAt"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, u := range users {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), u.ID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), u.Username)
			f.SetCellValue(sheetName, "D"+strconv.Itoa(row), u.Email)
			f.SetCellValue(sheetName, "E"+strconv.Itoa(row), u.UserGroup)
			f.SetCellValue(sheetName, "F"+strconv.Itoa(row), u.TotalDuration)
			f.SetCellValue(sheetName, "G"+strconv.Itoa(row), u.ListeningDuration)
			if u.CreatedAt != nil {
				f.SetCellValue(sheetName, "H"+strconv.Itoa(row), u.CreatedAt.Format(time.DateTime))
			}
		}
	}

	// 3. Export Songs
	var songs []model.Song
	if err := db.Find(&songs).Error; err == nil {
		sheetName := "Songs"
		f.NewSheet(sheetName)
		headers := []string{"ID", "Title", "Artist", "Album", "Duration", "FilePath", "PlayCount", "Year", "Format"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, s := range songs {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), s.ID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), s.Title)
			f.SetCellValue(sheetName, "C"+strconv.Itoa(row), s.ArtistName)
			f.SetCellValue(sheetName, "D"+strconv.Itoa(row), s.AlbumName)
			f.SetCellValue(sheetName, "E"+strconv.Itoa(row), s.Duration)
			f.SetCellValue(sheetName, "F"+strconv.Itoa(row), s.FilePath)
			f.SetCellValue(sheetName, "G"+strconv.Itoa(row), s.PlayCount)
			f.SetCellValue(sheetName, "H"+strconv.Itoa(row), s.Year)
			f.SetCellValue(sheetName, "I"+strconv.Itoa(row), s.Format)
		}
	}

	// 4. Export Playlists
	var playlists []model.Playlist
	if err := db.Find(&playlists).Error; err == nil {
		sheetName := "Playlists"
		f.NewSheet(sheetName)
		headers := []string{"ID", "Title", "OwnerID", "IsPublic", "PlayCount", "TotalSongs"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, p := range playlists {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), p.ID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), p.Title)
			f.SetCellValue(sheetName, "C"+strconv.Itoa(row), p.OwnerID)
			f.SetCellValue(sheetName, "D"+strconv.Itoa(row), p.IsPublic)
			f.SetCellValue(sheetName, "E"+strconv.Itoa(row), p.PlayCount)
			f.SetCellValue(sheetName, "F"+strconv.Itoa(row), p.TotalSongs)
		}
	}

	// 5. Export Albums
	var albums []model.Album
	if err := db.Preload("Artist").Find(&albums).Error; err == nil {
		sheetName := "Albums"
		f.NewSheet(sheetName)
		headers := []string{"ID", "Title", "Artist", "ReleaseDate"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, a := range albums {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), a.ID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), a.Title)
			f.SetCellValue(sheetName, "C"+strconv.Itoa(row), a.Artist.Name)
			if a.ReleaseDate != nil {
				f.SetCellValue(sheetName, "D"+strconv.Itoa(row), a.ReleaseDate.Format(time.DateOnly))
			}
		}
	}

	// 6. Export Artists
	var artists []model.Artist
	if err := db.Find(&artists).Error; err == nil {
		sheetName := "Artists"
		f.NewSheet(sheetName)
		headers := []string{"ID", "Name", "Description"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, a := range artists {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), a.ID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), a.Name)
			f.SetCellValue(sheetName, "C"+strconv.Itoa(row), a.Description)
		}
	}

	// Delete default Sheet1 if unused
	f.DeleteSheet("Sheet1")

	// Save to disk (same directory as data.db)
	// data.db path: g.Conf.SQLite.Dsn (relative to server/data/)
	// But let's check config to be sure where data.db is

	// Assuming data.db is in filepath.Join(g.GetConfig().BasicPath.FilePath, g.GetConfig().BasicPath.FileName, "data")
	// Or simply use the configured SQLite DSN path if it's absolute, but usually it's relative.
	// Let's use the standard data directory:
	saveDir := filepath.Join(g.GetConfig().BasicPath.FilePath, g.GetConfig().BasicPath.FileName, "data")
	fileName := fmt.Sprintf("database_export_%s.xlsx", time.Now().Format("20060102150405"))
	savePath := filepath.Join(saveDir, fileName)

	if err := f.SaveAs(savePath); err != nil {
		ReturnError(c, g.Err, "保存Excel文件失败: "+err.Error())
		return
	}

	ReturnSuccess(c, gin.H{
		"message":   "导出成功",
		"file_path": savePath,
		"file_name": fileName,
	})
}

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
