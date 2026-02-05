package handle

import (
	g "server/internal/global"
	"server/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SearchSuggest 搜索建议 (quick search for dropdown)
func SearchSuggest(c *gin.Context) {
	keyword := c.Query("keywords")
	if keyword == "" {
		ReturnError(c, g.ErrRequest, "关键词不能为空")
		return
	}

	result, err := model.SearchSuggest(GetDB(c), keyword)
	if err != nil {
		ReturnError(c, g.Err, err.Error())
		return
	}

	ReturnSuccess(c, result)
}

// SearchDetail 搜索详情 (full search page)
func SearchDetail(c *gin.Context) {
	keyword := c.Query("keywords")
	typeStr := c.Query("type") // currently only support "1" for songs
	pageStr := c.DefaultQuery("offset", "1")
	limitStr := c.DefaultQuery("limit", "30")

	if keyword == "" {
		ReturnError(c, g.ErrRequest, "关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// In Netease API, type 1 is song, 10 is album, 100 is artist, 1000 is playlist
	// For now, let's just implement song search
	if typeStr == "1" || typeStr == "" {
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
	} else {
		ReturnError(c, g.ErrRequest, "暂不支持该类型搜索")
	}
}
