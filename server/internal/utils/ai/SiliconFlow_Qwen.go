package ai

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	g "server/internal/global"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/time/rate"
)

// 硅基流动速率限制 (较默认配额提高5倍)
const (
	RPM = 5000   // 每分钟请求数
	TPM = 250000 // 每分钟 Token 数
)

var limiter *rate.Limiter

func init() {
	limiter = rate.NewLimiter(rate.Limit(RPM)/60, 5)
}

// ChatWithQwen 调用硅基流动 Qwen 模型进行对话
func ChatWithQwen(prompt string) (string, error) {
	conf := g.GetConfig()

	// 1. 获取 API Key
	apiKey := conf.SiliconFlow.ApiKey
	if apiKey == "" {
		apiKey = os.Getenv("SILICON_FLOW_API_KEY")
	}
	if apiKey == "" {
		return "", fmt.Errorf("SiliconFlow API Key not found in config or environment")
	}

	// 2. 配置 OpenAI Client
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.siliconflow.cn/v1"
	config.HTTPClient = &http.Client{Timeout: 60 * time.Second}

	client := openai.NewClientWithConfig(config)

	// 3. 确定模型
	modelName := conf.SiliconFlow.Model
	if modelName == "" {
		modelName = "Pro/deepseek-ai/DeepSeek-V3.2"
	}

	// 4. 等待限流器放行
	ctx := context.Background()
	if err := limiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limiter wait failed: %v", err)
	}

	log.Printf("[SiliconFlow] Calling model: %s", modelName)

	// 5. 发起请求
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: modelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 512,
		},
	)

	if err != nil {
		return "", fmt.Errorf("api call failed: %v", err)
	}

	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content

		usage := resp.Usage
		log.Printf("[SiliconFlow] Token usage: Prompt=%d, Completion=%d, Total=%d",
			usage.PromptTokens, usage.CompletionTokens, usage.TotalTokens)

		return content, nil
	}

	return "", fmt.Errorf("no response content")
}
