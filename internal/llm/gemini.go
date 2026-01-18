package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// GeminiClient Google Gemini 客戶端實現
type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
	log    *zap.Logger
}

// NewGeminiClient 創建新的 Gemini 客戶端
func NewGeminiClient(ctx context.Context, apiKey string, modelName string, log *zap.Logger) (*GeminiClient, error) {
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is empty")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	if modelName == "" {
		modelName = os.Getenv("GEMINI_MODEL")
	}
	if modelName == "" {
		modelName = "gemini-3-pro" // 2026 預設最強旗艦模型
	}

	model := client.GenerativeModel(modelName)
	
	// 設定生成參數 (可選，根據需求調整)
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	
	// 安全設定：放寬限制以避免過度過濾創意內容
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}

	return &GeminiClient{
		client: client,
		model:  model,
		log:    log,
	}, nil
}

// GenerateContent 生成文本內容
func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	if c.log != nil {
		c.log.Info("Sending request to Gemini", zap.String("model", c.model.Name()))
	}
	
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Gemini generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("Gemini returned empty response")
	}

	// 提取文本內容
	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			result += string(txt)
		}
	}

	return result, nil
}

// GenerateImage 生成圖片 (預留接口)
func (c *GeminiClient) GenerateImage(ctx context.Context, prompt string) (string, error) {
	return "", fmt.Errorf("Gemini image generation not implemented yet")
}

// Close 關閉客戶端
func (c *GeminiClient) Close() error {
	return c.client.Close()
}
