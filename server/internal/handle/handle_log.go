package handle

import (
	"encoding/json"
	"fmt"
	g "server/internal/global"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetLogs 获取历史日志（admin 权限）
func (*SystemAuth) GetLogs(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil || user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "权限不足，仅管理员可执行")
		return
	}

	level := strings.ToUpper(c.Query("level"))
	search := strings.ToLower(c.Query("search"))
	limit := 200
	if l := c.Query("limit"); l != "" {
		if n, err := parseInt(l); err == nil && n > 0 && n <= 5000 {
			limit = n
		}
	}

	all := g.LogBuffer.GetAll()

	var filtered []g.LogEntry
	for i := len(all) - 1; i >= 0 && len(filtered) < limit; i-- {
		entry := all[i]
		if level != "" && !strings.EqualFold(entry.Level, level) {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(entry.Message), search) {
			continue
		}
		filtered = append(filtered, entry)
	}

	// Reverse back to chronological order
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	ReturnSuccess(c, gin.H{
		"logs":  filtered,
		"total": len(all),
	})
}

// StreamLogs SSE 实时推送日志（admin 权限）
func (*SystemAuth) StreamLogs(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil || user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "权限不足，仅管理员可执行")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	subID, ch := g.LogBuffer.Subscribe()
	defer g.LogBuffer.Unsubscribe(subID)

	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		ReturnError(c, g.Err, "SSE not supported")
		return
	}

	// 立即发送 SSE 注释以触发 HTTP 响应头发送，避免 fetch 挂起
	fmt.Fprintf(c.Writer, ": ok\n\n")
	flusher.Flush()

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(entry)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid number: %s", s)
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
