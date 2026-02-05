package handle

import (
	g "server/internal/global"
	"server/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddHistory 记录播放历史
func (*SongAuth) AddHistory(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	var req struct {
		SongID int `json:"song_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	db := GetDB(c)
	if err := model.AddHistory(db, user.ID, req.SongID); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// GetHistory 获取播放历史
func (*SongAuth) GetHistory(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	db := GetDB(c)
	histories, total, err := model.GetUserHistory(db, user.ID, page, limit)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 转换格式适配前端
	// 这里我们直接返回历史记录列表，前端可能需要适配一下字段
	// 或者我们在此处转换为前端惯用的 Song 列表格式
	var songList []gin.H
	for _, h := range histories {
		s := h.Song
		artistName := "Unknown"
		if s.ArtistID != nil {
			artistName = s.ArtistName
		}
		albumName := "Unknown"
		if s.AlbumID != nil {
			albumName = s.AlbumName
		}
		coverUrl := ""
		if s.CoverID != nil && s.Cover.ID != 0 {
			coverUrl = "/covers/" + s.Cover.Path
		}

		songList = append(songList, gin.H{
			"id":          s.ID,
			"title":       s.Title,
			"artist":      artistName, // Compatible with some frontends
			"artist_name": artistName,
			"album":       albumName, // Compatible with some frontends
			"album_title": albumName,
			"duration":    s.Duration,
			"cover_url":   coverUrl,
			"played_at":   h.CreatedAt,
		})
	}

	ReturnSuccess(c, gin.H{
		"list":  songList,
		"total": total,
	})
}

// ClearHistory 清空播放历史
func (*SongAuth) ClearHistory(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	db := GetDB(c)
	if err := model.ClearUserHistory(db, user.ID); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}
