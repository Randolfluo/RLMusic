package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	g "server/internal/global"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/time/rate"
)

// 定义速率限制常量 (根据用户提供的 Qwen/Qwen3-8B 限制)
const (
	// RPM (Requests Per Minute): 1,000
	RateLimitRPM = 1000
	// TPM (Tokens Per Minute): 50,000
	RateLimitTPM = 50000
)

var limiter *rate.Limiter

func init() {
	// 初始化速率限制器
	// 将 RPM 转换为每秒允许的请求数 (Limit)
	limit := rate.Limit(float64(RateLimitRPM) / 60.0)
	// burst 设为 1，表示不允许突发超过限制
	limiter = rate.NewLimiter(limit, 1)
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
	
	client := openai.NewClientWithConfig(config)

	// 3. 确定模型
	modelName := conf.SiliconFlow.Model
	if modelName == "" {
		modelName = "Qwen/Qwen3-8B"
	}

	// 4. 等待限流
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
			MaxTokens: 512, // 简单限制
		},
	)

	if err != nil {
		return "", fmt.Errorf("api call failed: %v", err)
	}

	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		
		// 记录 Token 使用情况
		usage := resp.Usage
		log.Printf("[SiliconFlow] Token usage: Prompt=%d, Completion=%d, Total=%d", 
			usage.PromptTokens, usage.CompletionTokens, usage.TotalTokens)
		
		return content, nil
	}

	return "", fmt.Errorf("no response content")
}
