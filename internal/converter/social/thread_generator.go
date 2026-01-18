package social

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ThreadGenerator 負責將長文轉換為 X (Twitter) 貼文串
type ThreadGenerator struct {
	MaxCharsPerTweet int
}

// NewThreadGenerator 建立新的貼文串生成器
func NewThreadGenerator() *ThreadGenerator {
	return &ThreadGenerator{
		MaxCharsPerTweet: 280,
	}
}

// GenerateThread 將文章內容轉換為貼文串切片
func (g *ThreadGenerator) GenerateThread(content string) []string {
	// 1. 預處理：統一換行符，去除首尾空白
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.TrimSpace(content)

	// 2. 按段落分割
	paragraphs := strings.Split(content, "\n\n")

	var tweets []string
	var currentTweet strings.Builder
	
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// 檢查加入當前段落是否會超長 (保留一些緩衝給編號和間距)
		// 編號格式預估: " (1/5)" 佔約 6-8 字符，我們預留 15 字符
		if utf8.RuneCountInString(currentTweet.String()) + utf8.RuneCountInString(p) + 2 > (g.MaxCharsPerTweet - 15) {
			// 如果已經有內容，先送出當前 tweet
			if currentTweet.Len() > 0 {
				tweets = append(tweets, currentTweet.String())
				currentTweet.Reset()
			}
			
			// 如果單一段落本身就超長，需要強制分割 (較少見，但需處理)
			if utf8.RuneCountInString(p) > (g.MaxCharsPerTweet - 15) {
				// 簡單暴力分割，後續可優化為按句號分割
				runes := []rune(p)
				for len(runes) > 0 {
					end := g.MaxCharsPerTweet - 15
					if end > len(runes) {
						end = len(runes)
					}
					tweets = append(tweets, string(runes[:end]))
					runes = runes[end:]
				}
				continue
			}
		}

		if currentTweet.Len() > 0 {
			currentTweet.WriteString("\n\n")
		}
		currentTweet.WriteString(p)
	}

	// 加入最後一段
	if currentTweet.Len() > 0 {
		tweets = append(tweets, currentTweet.String())
	}

	// 3. 加入編號與裝飾
	formattedTweets := make([]string, len(tweets))
	total := len(tweets)
	
	for i, t := range tweets {
		// 第一則貼文加上 Thread 標記
		prefix := ""
		if i == 0 && total > 1 {
			prefix = "🧵 "
		}

		// 編號逻辑 (1/N)
		suffix := fmt.Sprintf(" (%d/%d)", i+1, total)
		
		formattedTweets[i] = fmt.Sprintf("%s%s%s", prefix, t, suffix)
	}

	return formattedTweets
}
