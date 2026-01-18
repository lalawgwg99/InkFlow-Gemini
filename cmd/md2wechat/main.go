package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"github.com/geekjourneyx/md2wechat-skill/internal/draft"
	"github.com/geekjourneyx/md2wechat-skill/internal/image"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	cfg *config.Config
	log *zap.Logger
)

// initConfig 初始化配置（延迟加载，允许 help 命令无需配置）
func initConfig() error {
	if cfg != nil && log != nil {
		return nil
	}

	var err error
	cfg, err = config.Load()
	if err != nil {
		return err
	}

	log, err = zap.NewProduction()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "md2wechat",
		Short: "社群內容生成神器 (SocialContent-AI)",
		Long: `SocialContent-AI (原 md2wechat): 專為台灣社群經營者設計的內容生成工具。
將 Markdown 文章一鍵轉換為 IG/FB 爆款卡片、X 貼文串，或微信公眾號格式。

環境變數 (Environment Variables):
  GEMINI_API_KEY                 Google Gemini API Key (推薦使用)
  WECHAT_APPID                   微信公眾號 AppID (僅微信模式需要)
  WECHAT_SECRET                  微信 API Secret (僅微信模式需要)
  IMAGE_API_KEY                  圖片生成 API Key (如 OpenAI DALL-E)
  COMPRESS_IMAGES                自動壓縮過大圖片 (預設: true)

範例:
  md2wechat convert post.md --mode card --theme magazine-dark  (製作 IG 卡片)
  md2wechat convert thought.md --mode thread                   (製作 X 貼文串)
  md2wechat upload_image ./photo.jpg                           (上傳素材)
  md2wechat generate_image "一隻在寫程式的貓"                  (AI 生圖)`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// upload_image command
	var uploadImageCmd = &cobra.Command{
		Use:   "upload_image <file_path>",
		Short: "上傳本地圖片到素材庫 (微信模式)",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			processor := image.NewProcessor(cfg, log)
			result, err := processor.UploadLocalImage(filePath)
			if err != nil {
				responseError(err)
				return
			}
			responseSuccess(result)
		},
	}
	rootCmd.AddCommand(uploadImageCmd)

	// download_and_upload command
	var downloadAndUploadCmd = &cobra.Command{
		Use:   "download_and_upload <url>",
		Short: "下載網路圖片並上傳",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			processor := image.NewProcessor(cfg, log)
			result, err := processor.DownloadAndUpload(url)
			if err != nil {
				responseError(err)
				return
			}
			responseSuccess(result)
		},
	}
	rootCmd.AddCommand(downloadAndUploadCmd)

	// generate_image command
	var generateImageCmdSize string
	var generateImageCmd = &cobra.Command{
		Use:   "generate_image <prompt>",
		Short: "使用 AI 生成圖片並上傳",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			prompt := args[0]
			processor := image.NewProcessor(cfg, log)

			// 如果指定了尺寸，临时覆盖配置
			if generateImageCmdSize != "" {
				result, err := processor.GenerateAndUploadWithSize(prompt, generateImageCmdSize)
				if err != nil {
					responseError(err)
					return
				}
				responseSuccess(result)
				return
			}

			result, err := processor.GenerateAndUpload(prompt)
			if err != nil {
				responseError(err)
				return
			}
			responseSuccess(result)
		},
	}
	generateImageCmd.Flags().StringVar(&generateImageCmdSize, "size", "", "圖片尺寸 (例如: 2560x1440)")
	generateImageCmd.Flags().StringVar(&generateImageCmdSize, "s", "", "圖片尺寸 (簡寫)")
	rootCmd.AddCommand(generateImageCmd)

	// create_draft command
	var createDraftCmd = &cobra.Command{
		Use:   "create_draft <json_file>",
		Short: "從 JSON 建立草稿 (進階)",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			jsonFile := args[0]
			svc := draft.NewService(cfg, log)
			result, err := svc.CreateDraftFromFile(jsonFile)
			if err != nil {
				responseError(err)
				return
			}
			responseSuccess(result)
		},
	}
	rootCmd.AddCommand(createDraftCmd)

	// convert command
	rootCmd.AddCommand(convertCmd)

	// config command
	rootCmd.AddCommand(configCmd)

	// write command
	rootCmd.AddCommand(writeCmd)

	// test-draft command
	rootCmd.AddCommand(testHTMLCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		responseError(err)
		os.Exit(1)
	}
}

func responseSuccess(data any) {
	response := map[string]any{
		"success": true,
		"data":    data,
	}
	printJSON(response)
}

func responseError(err error) {
	response := map[string]any{
		"success": false,
		"error":   err.Error(),
	}
	printJSON(response)
	os.Exit(1)
}

func printJSON(v any) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "JSON encode error: %v\n", err)
		os.Exit(1)
	}
}

// maskMediaID 遮蔽 media_id 用于日志
func maskMediaID(id string) string {
	if len(id) < 8 {
		return "***"
	}
	return id[:4] + "***" + id[len(id)-4:]
}
