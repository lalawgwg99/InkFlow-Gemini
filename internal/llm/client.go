package llm

import (
	"context"
)

// LLMClient 定義通用 LLM 客戶端接口
type LLMClient interface {
	// GenerateContent 生成文本內容
	GenerateContent(ctx context.Context, prompt string) (string, error)
	
	// GenerateImage 生成圖片 (如果支持)
	// 返回圖片 URL 或 Base64
	GenerateImage(ctx context.Context, prompt string) (string, error)
	
	// Close 關閉客戶端連接
	Close() error
}
