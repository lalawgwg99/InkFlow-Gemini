package main

import (
	"encoding/json"
	"fmt"
	"os"

	"strings"
	
	"github.com/geekjourneyx/md2wechat-skill/internal/converter"
	"github.com/geekjourneyx/md2wechat-skill/internal/draft"
	"github.com/geekjourneyx/md2wechat-skill/internal/image"
	"github.com/geekjourneyx/md2wechat-skill/internal/wechat"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// convertCmd convert 命令
var convertCmd = &cobra.Command{
	Use:   "convert <markdown_file>",
	Short: "將 Markdown 轉換為社群貼文格式 (HTML/Thread)",
	Long: `SocialContent-AI: 將 Markdown 文章轉換為適合 IG、FB、X (Twitter) 的格式。

支援的轉換模式 (Mode):
  - api:    使用 md2wechat.cn API (穩定的公眾號 HTML 轉換)
  - ai:     使用 Gemini/Claude AI 生成內容與 HTML (推薦)
  - thread: 生成 X (Twitter) 貼文串結構 (自動斷句編號)
  - card:   生成 IG/FB 圖片卡片 HTML (可用於截圖發文)

支援的主題 (Theme):
  API 模式: default, bytedance, apple, sports, chinese, cyber
  AI 模式: autumn-warm (秋日), spring-fresh (春日), ocean-calm (深海), custom
  社群模式: magazine-dark (黑金雜誌), minimalist-gray (極簡灰調), tech-neon (賽博霓虹)`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := runConvert(cmd, args); err != nil {
			responseError(err)
		}
	},
}

// convert 命令参数
var (
	convertMode         string
	convertTheme        string
	convertAPIKey       string
	convertFontSize     string
	convertCustomPrompt string
	convertOutput       string
	convertPreview      bool
	convertUpload       bool
	convertDraft        bool
	convertSaveDraft    string
	convertCoverImage   string // 封面图片路径
)

func init() {
	// 添加 flags
	// 添加 flags
	convertCmd.Flags().StringVar(&convertMode, "mode", "api", "轉換模式: api, ai, thread, card")
	convertCmd.Flags().StringVar(&convertTheme, "theme", "default", "主題名稱")
	convertCmd.Flags().StringVar(&convertAPIKey, "api-key", "", "md2wechat.cn API Key (僅 API 模式需要)")
	convertCmd.Flags().StringVar(&convertFontSize, "font-size", "medium", "字體大小: small/medium/large (僅 API 模式)")
	convertCmd.Flags().StringVar(&convertCustomPrompt, "custom-prompt", "", "自定義 AI 提示詞 (僅 AI 模式)")
	convertCmd.Flags().StringVarP(&convertOutput, "output", "o", "", "輸出檔案路徑")
	convertCmd.Flags().BoolVar(&convertPreview, "preview", false, "僅預覽，不進行圖片上傳")
	convertCmd.Flags().BoolVar(&convertUpload, "upload", false, "上傳圖片並替換網址 (微信模式)")
	convertCmd.Flags().BoolVar(&convertDraft, "draft", false, "轉換後直接建立微信草稿")
	convertCmd.Flags().StringVar(&convertSaveDraft, "save-draft", "", "保存草稿 JSON 到文件")
	convertCmd.Flags().StringVar(&convertCoverImage, "cover", "", "草稿封面圖片路徑 (使用 --draft 時必填)")
}

// runConvert 执行转换
func runConvert(cmd *cobra.Command, args []string) error {
	markdownFile := args[0]

	log.Info("開始轉換...",
		zap.String("檔案", markdownFile),
		zap.String("模式", convertMode),
		zap.String("主題", convertTheme))

	// 读取 Markdown 文件
	markdown, err := os.ReadFile(markdownFile)
	if err != nil {
		return fmt.Errorf("read markdown file: %w", err)
	}

	// 创建转换器
	conv := converter.NewConverter(cfg, log)

	// 构建转换请求
	req := &converter.ConvertRequest{
		Markdown:     string(markdown),
		Mode:         converter.ConvertMode(convertMode),
		Theme:        convertTheme,
		APIKey:       convertAPIKey,
		FontSize:     convertFontSize,
		CustomPrompt: convertCustomPrompt,
	}

	// 执行转换
	result := conv.Convert(req)

	if !result.Success {
		return fmt.Errorf("conversion failed: %s", result.Error)
	}

	log.Info("轉換成功！",
		zap.String("模式", string(result.Mode)),
		zap.String("主題", result.Theme),
		zap.Int("圖片數量", len(result.Images)))

	// 根据模式处理结果
	if convertMode == "ai" && converter.IsAIRequest(result) {
		// AI 模式需要外部处理
		return handleAIResult(result, markdownFile)
	}

	// 處理 Social 模式結果
	if convertMode == "thread" || convertMode == "card" {
		return handleSocialResult(result, markdownFile)
	}

	// 處理圖片
	if convertUpload || convertDraft {
		if err := processImages(result); err != nil {
			log.Warn("圖片處理部分失敗", zap.Error(err))
		}
	}

	// 输出结果
	if convertSaveDraft != "" {
		if err := saveDraft(result); err != nil {
			return fmt.Errorf("save draft: %w", err)
		}
	}

	if convertDraft {
		if err := createWeChatDraft(result, convertCoverImage); err != nil {
			return fmt.Errorf("create draft: %w", err)
		}
	}

	// 输出 HTML
	outputHTML(result.HTML, convertOutput, convertPreview)

	return nil
}

// handleAIResult 处理 AI 模式结果
func handleAIResult(result *converter.ConvertResult, markdownFile string) error {
	prompt, images, ok := converter.GetAIRequestInfo(result)
	if !ok {
		return fmt.Errorf("invalid AI request result")
	}

	log.Info("AI mode request prepared",
		zap.Int("image_count", len(images)),
		zap.Int("prompt_length", len(prompt)))

	// 输出 AI 请求信息
	response := map[string]any{
		"success":       true,
		"mode":          "ai",
		"action":        "ai_request",
		"markdown_file": markdownFile,
		"prompt":        prompt,
		"images":        images,
	}

	printJSON(response)

	if convertOutput != "" {
		// 同时保存原始 markdown 到输出文件，方便用户使用
		if err := os.WriteFile(convertOutput, []byte(prompt), 0644); err != nil {
			log.Warn("failed to save prompt", zap.Error(err))
		}
	}

	return nil
}

// handleSocialResult 處理 Social 模式結果
func handleSocialResult(result *converter.ConvertResult, markdownFile string) error {
	if result.Mode == converter.ModeThread {
		// 輸出 Thread 結構
		fmt.Printf("\n=== X (Twitter) 貼文串預覽 (%d 則貼文) ===\n\n", len(result.ThreadTweets))
		for _, t := range result.ThreadTweets {
			fmt.Println(t)
			fmt.Println("\n---")
		}
		
		if convertOutput != "" {
			// 保存為 JSON 或純文本
			content := strings.Join(result.ThreadTweets, "\n\n---\n\n")
			if err := os.WriteFile(convertOutput, []byte(content), 0644); err != nil {
				log.Error("儲存貼文串失敗", zap.Error(err))
			} else {
				log.Info("貼文串已儲存", zap.String("檔案", convertOutput))
			}
		}
		return nil
	}

	if result.Mode == converter.ModeCard {
		// 輸出 HTML
		outputHTML(result.CardHTML, convertOutput, convertPreview)
		return nil
	}

	return nil
}

// processImages 处理图片上传
func processImages(result *converter.ConvertResult) error {
	if len(result.Images) == 0 {
		log.Info("沒有圖片需要處理")
		return nil
	}

	processor := image.NewProcessor(cfg, log)

	for i, imgRef := range result.Images {
		log.Info("processing image",
			zap.Int("index", i),
			zap.String("type", string(imgRef.Type)),
			zap.String("original", imgRef.Original))

		var uploadResult *image.UploadResult
		var err error

		switch imgRef.Type {
		case converter.ImageTypeLocal:
			uploadResult, err = processor.UploadLocalImage(imgRef.Original)
		case converter.ImageTypeOnline:
			uploadResult, err = processor.DownloadAndUpload(imgRef.Original)
		case converter.ImageTypeAI:
			// AI 生成的图片需要先调用生成 API
			genResult, genErr := processor.GenerateAndUpload(imgRef.AIPrompt)
			if genErr != nil {
				err = genErr
			} else {
				uploadResult = &image.UploadResult{
					MediaID:   genResult.MediaID,
					WechatURL: genResult.WechatURL,
				}
			}
		}

		if err != nil {
			log.Warn("image upload failed",
				zap.Int("index", i),
				zap.Error(err))
			continue
		}

		// 更新图片 URL
		result.Images[i].WechatURL = uploadResult.WechatURL

		log.Info("image uploaded",
			zap.Int("index", i),
			zap.String("media_id", maskMediaID(uploadResult.MediaID)),
			zap.String("wechat_url", uploadResult.WechatURL))
	}

	// 替换 HTML 中的图片占位符
	result.HTML = converter.ReplaceImagePlaceholders(result.HTML, result.Images)

	return nil
}

// saveDraft 保存草稿 JSON 到文件
func saveDraft(result *converter.ConvertResult) error {
	articles := []draft.Article{
		{
			Title:   "Draft Article", // TODO: 从 markdown 提取标题
			Content: result.HTML,
		},
	}

	draftData := map[string]any{
		"articles": articles,
	}

	jsonData, err := json.MarshalIndent(draftData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal draft: %w", err)
	}

	if err := os.WriteFile(convertSaveDraft, jsonData, 0644); err != nil {
		return fmt.Errorf("write draft file: %w", err)
	}

	log.Info("draft saved", zap.String("file", convertSaveDraft))
	return nil
}

// createWeChatDraft 创建微信草稿
func createWeChatDraft(result *converter.ConvertResult, coverImagePath string) error {
	svc := draft.NewService(cfg, log)

	// 检查封面图片（微信要求必须有封面图）
	if coverImagePath == "" {
		return &DraftError{
			Message: "创建草稿需要封面图片",
			Hint:    "请使用 --cover 参数指定封面图片路径，例如: --cover /path/to/cover.jpg\n" +
				"或者先上传封面图片到微信素材库: md2wechat upload_image /path/to/cover.jpg",
		}
	}

	// 上传封面图片到微信素材库
	log.Info("uploading cover image", zap.String("path", coverImagePath))
	coverMediaID, err := uploadCoverImage(coverImagePath)
	if err != nil {
		return fmt.Errorf("上传封面图片失败: %w", err)
	}
	log.Info("cover image uploaded", zap.String("media_id", maskMediaID(coverMediaID)))

	// 提取标题（TODO: 从 markdown frontmatter 或第一个标题获取）
	title := "Article Title"

	draftResult, err := svc.CreateDraft([]draft.Article{
		{
			Title:          title,
			Content:        result.HTML,
			Digest:         draft.GenerateDigestFromContent(result.HTML, 120),
			ThumbMediaID:   coverMediaID,
			ShowCoverPic:   1, // 显示封面
		},
	})

	if err != nil {
		return fmt.Errorf("create draft: %w", err)
	}

	log.Info("draft created",
		zap.String("media_id", maskMediaID(draftResult.MediaID)),
		zap.String("draft_url", draftResult.DraftURL))

	return nil
}

// uploadCoverImage 上传封面图片到微信素材库
func uploadCoverImage(imagePath string) (string, error) {
	svc := wechat.NewService(cfg, log)
	result, err := svc.UploadMaterial(imagePath)
	if err != nil {
		return "", err
	}
	return result.MediaID, nil
}

// DraftError 草稿错误
type DraftError struct {
	Message string
	Hint    string
}

func (e *DraftError) Error() string {
	msg := fmt.Sprintf("草稿错误: %s", e.Message)
	if e.Hint != "" {
		msg += fmt.Sprintf("\n💡 提示:\n   %s", e.Hint)
	}
	return msg
}

// outputHTML 输出 HTML
func outputHTML(html, outputPath string, preview bool) {
	if preview || outputPath == "" {
		// 预览模式或未指定输出，输出到标准输出
		fmt.Println("\n=== HTML Output ===")
		fmt.Println(html)
		fmt.Println("\n=== End ===")
	}

	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
			log.Error("failed to write output file", zap.Error(err))
		} else {
			log.Info("html saved", zap.String("file", outputPath))
		}
	}
}
