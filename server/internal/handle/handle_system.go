package handle

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	g "server/internal/global"
	"server/internal/model"
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

	// 7. Export Histories
	var histories []model.History
	if err := db.Find(&histories).Error; err == nil {
		sheetName := "Histories"
		f.NewSheet(sheetName)
		headers := []string{"ID", "UserID", "SongID", "CreatedAt"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, h := range histories {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), h.ID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), h.UserID)
			f.SetCellValue(sheetName, "C"+strconv.Itoa(row), h.SongID)
			f.SetCellValue(sheetName, "D"+strconv.Itoa(row), h.CreatedAt.Format(time.DateTime))
		}
	}

	// 8. Export PlaylistSongs (Join Table)
	type PlaylistSong struct {
		PlaylistID int
		SongID     int
	}
	var playlistSongs []PlaylistSong
	if err := db.Table("playlist_songs").Find(&playlistSongs).Error; err == nil {
		sheetName := "PlaylistSongs"
		f.NewSheet(sheetName)
		headers := []string{"PlaylistID", "SongID"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, ps := range playlistSongs {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), ps.PlaylistID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), ps.SongID)
		}
	}

	// 9. Export SongArtists (Join Table)
	type SongArtist struct {
		SongID   int
		ArtistID int
	}
	var songArtists []SongArtist
	if err := db.Table("song_artists").Find(&songArtists).Error; err == nil {
		sheetName := "SongArtists"
		f.NewSheet(sheetName)
		headers := []string{"SongID", "ArtistID"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, h)
		}
		for i, sa := range songArtists {
			row := i + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), sa.SongID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), sa.ArtistID)
		}
	}

	// Delete default Sheet1 if unused
	f.DeleteSheet("Sheet1")

	// Set Headers for Download
	fileName := fmt.Sprintf("database_export_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("File-Name", fileName)
	c.Header("Access-Control-Expose-Headers", "File-Name")

	// Write to Response
	if _, err := f.WriteTo(c.Writer); err != nil {
		// Log error (can't return JSON error here as headers already sent)
		fmt.Printf("Export failed: %v\n", err)
	}
}

type SystemAuth struct{}

type UpdateConfigReq struct {
	FilePath string `json:"filepath"`
}

type StatsVO struct {
	SongCount     int64 `json:"song_count"`
	SongVolume    int64 `json:"song_volume"`
	AlbumCount    int64 `json:"album_count"`
	ArtistCount   int64 `json:"artist_count"`
	MusicDuration int64 `json:"music_duration"`

	PlaylistCount int64 `json:"playlist_count"`
	UserCount     int64 `json:"user_count"`
	SystemUptime  int64 `json:"system_uptime"`
}

type SystemStatusVO struct {
	CpuUsage     float64 `json:"cpu_usage"`
	MemUsage     float64 `json:"mem_usage"`
	ApiCallCount int64   `json:"api_call_count"`
	SystemUptime int64   `json:"system_uptime"`
	GoRoutines   int     `json:"go_routines"`
	TotalVolume  int64   `json:"total_volume"`
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

	ReturnSuccess(c, StatsVO{
		SongCount:     info.TotalSongs,
		SongVolume:    info.TotalVolume,
		AlbumCount:    info.TotalAlbums,
		ArtistCount:   info.TotalArtists,
		MusicDuration: info.TotalDuration,
		PlaylistCount: totalPlaylists,
		UserCount:     totalUsers,
		SystemUptime:  systemRunTime,
	})
}

// GetSystemStatus 获取系统实时状态 (CPU, Mem, API Count)
func (*SystemAuth) GetSystemStatus(c *gin.Context) {
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

	// 系统运行时间
	systemRunTime := int64(time.Since(g.StartTime).Seconds())

	// Go Routines
	goRoutines := runtime.NumGoroutine()

	// Total Song Volume (from SystemInfo)
	var totalVolume int64 = 0
	if info, err := model.GetSystemInfoStruct(GetDB(c)); err == nil && info != nil {
		totalVolume = info.TotalVolume
	}

	ReturnSuccess(c, SystemStatusVO{
		CpuUsage:     cpuUsage,
		MemUsage:     memUsage,
		ApiCallCount: apiCount,
		SystemUptime: systemRunTime,
		GoRoutines:   goRoutines,
		TotalVolume:  totalVolume,
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

// ResetSystem 重置系统数据 (Admin only)
func (*SystemAuth) ResetSystem(c *gin.Context) {
	// 1. Check Auth (Admin only)
	user := GetCurrentUser(c)
	if user == nil || user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "权限不足，仅管理员可执行")
		return
	}

	// 2. Define data directory
	// 假设 data 目录在运行根目录下，或者根据配置获取
	// 这里默认使用 ./data
	dataDir := "data"
	absPath, err := filepath.Abs(dataDir)
	if err != nil {
		ReturnError(c, g.Err, "无法解析数据目录路径")
		return
	}

	// 3. Delete all contents in data directory except database file
	// 遍历目录内容
	files, err := filepath.Glob(filepath.Join(absPath, "*"))
	if err != nil {
		ReturnError(c, g.Err, "无法读取数据目录: "+err.Error())
		return
	}

	deletedCount := 0
	var errorMsgs []string

	for _, f := range files {
		baseName := filepath.Base(f)
		// 跳过数据库文件 (假设扩展名为 .db 或者名字为 data.db)
		if strings.HasSuffix(baseName, ".db") || baseName == "data.db" || baseName == "music.db" {
			continue
		}

		// 强制删除
		if err := os.RemoveAll(f); err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("无法删除 %s: %v", baseName, err))
		} else {
			deletedCount++
		}
	}

	// 4. Clean up Database (Preserve Admin)
	db := GetDB(c)
	tx := db.Begin()
	if tx.Error != nil {
		ReturnError(c, g.ErrDbOp, "无法开启事务")
		return
	}

	// 禁用外键约束 (SQLite)
	tx.Exec("PRAGMA foreign_keys = OFF")

	// 清空数据表
	// 关联表
	tx.Exec("DELETE FROM playlist_songs")
	tx.Exec("DELETE FROM song_artists")
	tx.Exec("DELETE FROM user_subscribed_playlists")
	tx.Exec("DELETE FROM podcast_episodes")

	// 主实体表
	tx.Exec("DELETE FROM histories")
	tx.Exec("DELETE FROM songs")
	tx.Exec("DELETE FROM playlists")
	tx.Exec("DELETE FROM albums")
	tx.Exec("DELETE FROM artists")
	tx.Exec("DELETE FROM covers")
	tx.Exec("DELETE FROM podcasts")

	// 重置 SQLite 自增计数器 (可选，但推荐)
	tx.Exec("DELETE FROM sqlite_sequence WHERE name IN ('songs', 'playlists', 'albums', 'artists', 'covers', 'podcasts', 'histories')")

	// 删除非管理员用户
	if err := tx.Where("user_group != ?", "admin").Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		ReturnError(c, g.ErrDbOp, "清理用户失败: "+err.Error())
		return
	}

	// 重新启用外键约束
	tx.Exec("PRAGMA foreign_keys = ON")

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ReturnError(c, g.ErrDbOp, "数据库提交失败: "+err.Error())
		return
	}

	msg := fmt.Sprintf("文件清理: 成功删除 %d 个项目 (保留数据库)。数据库清理: 已重置所有数据，仅保留管理员账户。", deletedCount)
	if len(errorMsgs) > 0 {
		msg += fmt.Sprintf(" 文件删除警告: %s", strings.Join(errorMsgs, "; "))
		ReturnSuccess(c, gin.H{
			"message": msg,
			"warning": true,
		})
	} else {
		ReturnSuccess(c, gin.H{
			"message": msg,
		})
	}
}
