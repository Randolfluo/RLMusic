package handle

import (
	g "server/internal/global"
	"server/internal/utils/ai"

	"github.com/gin-gonic/gin"
)

type AIHandler struct{}

// ChatRequest 对话请求参数
type ChatRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// QwenChat 硅基流动 Qwen 对话接口
func (h *AIHandler) QwenChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	reply, err := ai.ChatWithQwen(req.Prompt)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	ReturnSuccess(c, gin.H{
		"reply": reply,
	})
}

// TTSRequest 语音合成请求参数
type TTSRequest struct {
	Text       string  `json:"text" binding:"required"`
	Voice      string  `json:"voice"`
	Format     string  `json:"format"`
	SampleRate int     `json:"sample_rate"`
	Volume     int     `json:"volume"`
	Rate       float64 `json:"rate"`
	Pitch      float64 `json:"pitch"`
}

// QwenTTS 阿里百炼语音合成接口
func (h *AIHandler) QwenTTS(c *gin.Context) {
	var req TTSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	params := &ai.Params{
		Voice:      req.Voice,
		Format:     req.Format,
		SampleRate: req.SampleRate,
		Volume:     req.Volume,
		Rate:       req.Rate,
		Pitch:      req.Pitch,
	}

	filePath, err := ai.SynthesizeAudio(req.Text, params)
	if err != nil {
		ReturnError(c, g.Err, err) // 使用通用的服务器错误
		return
	}

	ReturnSuccess(c, gin.H{
		"path": filePath,
	})
}
