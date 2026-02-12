package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	g "server/internal/global"

	"github.com/google/uuid"
)

const (
	// QwenTTSApiUrl 是阿里百炼语音合成 API 地址
	QwenTTSApiUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation"
)

// QwenTTSRequest 阿里百炼语音合成请求体
// 参考: https://help.aliyun.com/zh/model-studio/qwen-tts-api
type QwenTTSRequest struct {
	Model      string `json:"model"`
	Input      Input  `json:"input"`
	Parameters Params `json:"parameters,omitempty"`
}

type Input struct {
	Text string `json:"text"`
}

type Params struct {
	Voice      string  `json:"voice"`
	Format     string  `json:"format,omitempty"`
	SampleRate int     `json:"sample_rate,omitempty"`
	Volume     int     `json:"volume,omitempty"`
	Rate       float64 `json:"rate,omitempty"`
	Pitch      float64 `json:"pitch,omitempty"`
}

// QwenTTSResponse 阿里百炼语音合成响应体
type QwenTTSResponse struct {
	RequestId string `json:"request_id"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Output    struct {
		Audio struct {
			Url string `json:"url"`
		} `json:"audio"`
		FinishReason string `json:"finish_reason"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// SynthesizeAudio 执行语音合成并保存文件
// 返回保存的文件路径和错误信息
func SynthesizeAudio(text string, params *Params) (string, error) {
	conf := g.GetConfig()

	// 获取 API Key
	apiKey := conf.QwenTTS.ApiKey
	if apiKey == "" {
		apiKey = os.Getenv("QWEN_TTS_API_KEY")
	}
	if apiKey == "" {
		log.Println("[QwenTTS] API Key 未找到")
		return "", fmt.Errorf("Qwen TTS API Key not found in config or environment")
	}
	log.Println("[QwenTTS] 开始语音合成, 使用模型:", conf.QwenTTS.Model)

	// 准备请求参数
	model := conf.QwenTTS.Model
	if model == "" {
		model = "qwen3-tts-flash" // 默认模型
	}

	// 合并参数
	req := QwenTTSRequest{
		Model: model,
		Input: Input{Text: text},
		Parameters: Params{
			Voice:      conf.QwenTTS.Voice,
			Format:     conf.QwenTTS.Format,
			SampleRate: conf.QwenTTS.SampleRate,
			Volume:     conf.QwenTTS.Volume,
			Rate:       conf.QwenTTS.Rate,
			Pitch:      conf.QwenTTS.Pitch,
		},
	}

	// 如果传入了自定义参数，覆盖默认配置
	if params != nil {
		if params.Voice != "" {
			req.Parameters.Voice = params.Voice
		}
		if params.Format != "" {
			req.Parameters.Format = params.Format
		}
		if params.SampleRate != 0 {
			req.Parameters.SampleRate = params.SampleRate
		}
		if params.Volume != 0 {
			req.Parameters.Volume = params.Volume
		}
		if params.Rate != 0 {
			req.Parameters.Rate = params.Rate
		}
		if params.Pitch != 0 {
			req.Parameters.Pitch = params.Pitch
		}
	}

	// 默认 Format 为 flac
	if req.Parameters.Format == "" {
		req.Parameters.Format = "flac"
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", QwenTTSApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	// 非流式输出，不需要 X-DashScope-SSE: enable

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应
	var ttsResp QwenTTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&ttsResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if ttsResp.Code != "" {
		return "", fmt.Errorf("api error: %s - %s", ttsResp.Code, ttsResp.Message)
	}

	if ttsResp.Output.Audio.Url == "" {
		return "", fmt.Errorf("no audio url in response")
	}

	// 下载音频文件
	audioUrl := ttsResp.Output.Audio.Url
	audioResp, err := http.Get(audioUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %v", err)
	}
	defer audioResp.Body.Close()

	if audioResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download audio, status: %d", audioResp.StatusCode)
	}

	// 确定保存路径
	// 对应存储路径为config.yml定义的BasicPath下的boke（没有该文件夹则自动创建）
	basePath := filepath.Join(conf.BasicPath.FilePath, conf.BasicPath.FileName, "boke")
	log.Printf("[QwenTTS] 准备保存音频文件到目录: %s", basePath)
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Printf("[QwenTTS] 创建目录失败: %v", err)
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 生成文件名
	fileName := fmt.Sprintf("%s.%s", uuid.New().String(), req.Parameters.Format)
	filePath := filepath.Join(basePath, fileName)

	// 保存文件
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("[QwenTTS] 创建文件失败: %v", err)
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, audioResp.Body)
	if err != nil {
		log.Printf("[QwenTTS] 保存音频内容失败: %v", err)
		return "", fmt.Errorf("failed to save audio content: %v", err)
	}
	log.Printf("[QwenTTS] 音频文件保存成功: %s, 大小: %d 字节", filePath, written)

	return filePath, nil
}
