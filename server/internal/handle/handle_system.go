package handle

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	g "server/internal/global"
	"server/internal/model"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"net"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
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
	user := GetCurrentUser(c)
	if user == nil || user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "权限不足，仅管理员可执行")
		return
	}

	db := GetDB(c)
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("close excel failed: %v\n", err)
		}
	}()

	tables, err := db.Migrator().GetTables()
	if err != nil {
		ReturnError(c, g.ErrDbOp, "读取数据库表失败: "+err.Error())
		return
	}
	sort.Strings(tables)

	usedSheetNames := map[string]struct{}{}
	createdSheets := 0

	for _, tableName := range tables {
		sheetName := makeExcelSheetName(tableName, usedSheetNames)
		if createdSheets == 0 {
			if err := f.SetSheetName("Sheet1", sheetName); err != nil {
				fmt.Printf("set first sheet name failed: %v\n", err)
				continue
			}
		} else {
			if _, err := f.NewSheet(sheetName); err != nil {
				fmt.Printf("create sheet failed, table=%s err=%v\n", tableName, err)
				continue
			}
		}

		if err := writeTableToSheet(db, f, tableName, sheetName); err != nil {
			fmt.Printf("write table failed, table=%s err=%v\n", tableName, err)
			continue
		}

		createdSheets++
	}

	if createdSheets == 0 {
		f.SetCellValue("Sheet1", "A1", "No tables found")
	}

	fileName := fmt.Sprintf("database_export_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("File-Name", fileName)
	c.Header("Access-Control-Expose-Headers", "File-Name")

	if _, err := f.WriteTo(c.Writer); err != nil {
		fmt.Printf("Export failed: %v\n", err)
	}
}

func writeTableToSheet(db *gorm.DB, f *excelize.File, tableName, sheetName string) error {
	rows, err := db.Table(tableName).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for i, col := range columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, col)
	}

	rowNum := 2
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return err
		}

		for i, val := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, rowNum)
			f.SetCellValue(sheetName, cell, normalizeExcelCellValue(val))
		}
		rowNum++
	}

	return rows.Err()
}

func normalizeExcelCellValue(value any) any {
	switch v := value.(type) {
	case nil:
		return ""
	case []byte:
		return string(v)
	case time.Time:
		return v.Format(time.DateTime)
	default:
		return v
	}
}

func makeExcelSheetName(raw string, used map[string]struct{}) string {
	clean := strings.NewReplacer("\\", "_", "/", "_", "*", "_", "?", "_", "[", "_", "]", "_", ":", "_").Replace(raw)
	clean = strings.TrimSpace(clean)
	if clean == "" {
		clean = "Sheet"
	}

	trimTo := func(s string, max int) string {
		rs := []rune(s)
		if len(rs) <= max {
			return s
		}
		return string(rs[:max])
	}

	candidate := trimTo(clean, 31)
	if _, exists := used[candidate]; !exists {
		used[candidate] = struct{}{}
		return candidate
	}

	index := 2
	for {
		suffix := fmt.Sprintf("_%d", index)
		base := trimTo(clean, 31-len([]rune(suffix)))
		candidate = base + suffix
		if _, exists := used[candidate]; !exists {
			used[candidate] = struct{}{}
			return candidate
		}
		index++
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
