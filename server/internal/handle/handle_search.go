package handle

import (
	g "server/internal/global"
	"server/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getSearchParams(c *gin.Context) (string, int, int) {
	keyword := c.Query("keywords")
	pageStr := c.DefaultQuery("offset", "1")
	limitStr := c.DefaultQuery("limit", "30")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	return keyword, page, limit
}

// SearchSong 搜索歌曲
func SearchSong(c *gin.Context) {
	keyword, page, limit := getSearchParams(c)
	if keyword == "" {
		ReturnError(c, g.ErrRequest, "关键词不能为空")
		return
	}

	songs, total, err := model.SearchSongs(GetDB(c), keyword, page, limit)
	if err != nil {
		ReturnError(c, g.Err, err.Error())
		return
	}

	ReturnSuccess(c, gin.H{
		"result": gin.H{
			"songs":     songs,
			"songCount": total,
		},
		"code": 200,
	})
}

// SearchArtist 搜索歌手
func SearchArtist(c *gin.Context) {
	keyword, page, limit := getSearchParams(c)
	if keyword == "" {
		ReturnError(c, g.ErrRequest, "关键词不能为空")
		return
	}

	artists, total, err := model.SearchArtists(GetDB(c), keyword, page, limit)
	if err != nil {
		ReturnError(c, g.Err, err.Error())
		return
	}

	ReturnSuccess(c, gin.H{
		"result": gin.H{
			"artists":     artists,
			"artistCount": total,
		},
		"code": 200,
	})
}

// SearchAlbum 搜索专辑
func SearchAlbum(c *gin.Context) {
	keyword, page, limit := getSearchParams(c)
	if keyword == "" {
		ReturnError(c, g.ErrRequest, "关键词不能为空")
		return
	}

	albums, total, err := model.SearchAlbums(GetDB(c), keyword, page, limit)
	if err != nil {
		ReturnError(c, g.Err, err.Error())
		return
	}

	ReturnSuccess(c, gin.H{
		"result": gin.H{
			"albums":     albums,
			"albumCount": total,
		},
		"code": 200,
	})
}

// SearchPlaylist 搜索歌单
func SearchPlaylist(c *gin.Context) {
	keyword, page, limit := getSearchParams(c)
	if keyword == "" {
		ReturnError(c, g.ErrRequest, "关键词不能为空")
		return
	}

	playlists, total, err := model.SearchPlaylists(GetDB(c), keyword, page, limit)
	if err != nil {
		ReturnError(c, g.Err, err.Error())
		return
	}

	ReturnSuccess(c, gin.H{
		"result": gin.H{
			"playlists":     playlists,
			"playlistCount": total,
		},
		"code": 200,
	})
}
