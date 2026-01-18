// Package main provides the md2wechat CLI tool
package main

import (
	"fmt"
	"os"
	"strings"

	"bufio"
	"context"
	"strconv"
	
	"github.com/geekjourneyx/md2wechat-skill/internal/llm"
	"github.com/geekjourneyx/md2wechat-skill/internal/writer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// writeCmd 写作命令
var writeCmd = &cobra.Command{
	Use:   "write [input]",
	Short: "Writer Style Assistant - Write with creator styles",
	Long: `Assisted writing with customizable creator styles.

Default style: Dan Koe (profound, sharp, grounded)

Examples:
  # Interactive mode
  md2wechat write

  # Write from idea
  md2wechat write --style dan-koe

  # Refine existing content
  md2wechat write --style dan-koe --input-type fragment article.md

  # Generate with cover
  md2wechat write --style dan-koe --cover`,
	Args:  cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := runWrite(cmd, args); err != nil {
			responseError(err)
		}
	},
}

// write 命令参数
var (
	writeStyle       string
	writeInputType   string
	writeArticleType string
	writeLength      string
	writeTitle       string
	writeOutput      string
	writeCover       bool
	writeCoverOnly   bool
	writeListStyles  bool
	writeStyleDetail bool
	writeIdea        bool // 新增：靈感模式
)

func init() {
	// 添加 flags
	writeCmd.Flags().StringVar(&writeStyle, "style", "dan-koe", "Writer style")
	writeCmd.Flags().StringVar(&writeInputType, "input-type", "idea", "Input type: idea/fragment/outline/title")
	writeCmd.Flags().StringVar(&writeArticleType, "article-type", "essay", "Article type: essay/commentary/story/tutorial/review")
	writeCmd.Flags().StringVar(&writeLength, "length", "medium", "Article length: short/medium/long")
	writeCmd.Flags().StringVar(&writeTitle, "title", "", "Article title")
	writeCmd.Flags().StringVarP(&writeOutput, "output", "o", "", "Output file")
	writeCmd.Flags().BoolVar(&writeCover, "cover", false, "Generate matching cover")
	writeCmd.Flags().BoolVar(&writeCoverOnly, "cover-only", false, "Generate cover only")
	writeCmd.Flags().BoolVar(&writeListStyles, "list", false, "List all available styles")
	writeCmd.Flags().BoolVar(&writeStyleDetail, "detail", false, "Show detailed style info")
	writeCmd.Flags().BoolVar(&writeIdea, "idea", false, "Generate writing ideas (AI Brainstorming)")
}

// runWrite 执行写作命令
func runWrite(cmd *cobra.Command, args []string) error {
	// 处理列出风格
	if writeListStyles {
		return runListStyles()
	}

	// 處理靈感生成模式
	if writeIdea {
		return runIdeaGenerator()
	}

	// 获取输入内容
	input := ""
	if len(args) > 0 {
		// 从文件读取
		content, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("读取文件: %w", err)
		}
		input = string(content)

		// 如果没有明确指定输入类型，默认为 fragment
		if writeInputType == "idea" {
			writeInputType = "fragment"
		}
	}

	// 如果没有输入，进入交互模式
	if input == "" {
		return runInteractiveWrite()
	}

	// 执行写作
	return executeWrite(input)
}

// runListStyles 列出所有风格
func runListStyles() error {
	asst := writer.NewAssistant()
	result := asst.ListStyles()

	if !result.Success {
		return fmt.Errorf("%s", result.Error)
	}

	if writeStyleDetail {
		// 详细模式
		for _, style := range result.Styles {
			fmt.Println(writer.FormatStyleSummary(style))
			fmt.Println("---")
		}
	} else {
		// 简洁模式
		fmt.Println(writer.FormatStyleList(result.Styles))
	}

	return nil
}

// runInteractiveWrite 交互式写作模式
func runInteractiveWrite() error {
	fmt.Println("📝 Writer Style Assistant")
	fmt.Println()

	// 显示可用风格
	asst := writer.NewAssistant()
	styles := asst.GetAvailableStyles()

	fmt.Printf("可用风格 (%d 个):\n", len(styles))
	for _, styleName := range styles {
		style, _ := asst.GetStyleInfo(styleName)
		fmt.Printf("  - %s (%s)\n", style.Name, style.EnglishName)
	}
	fmt.Println()

	// 選單
	fmt.Println("請選擇模式：")
	fmt.Println("1. 🧠 靈感產生器 (AI 幫我想要寫什麼)")
	fmt.Println("2. ✍️  自由寫作 (我有主題了)")
	fmt.Print("\n請選擇 [1-2] (默認 2): ")
	modeInput := readLine()
	
	if modeInput == "1" {
		writeIdea = true
		return runIdeaGenerator()
	}

	// 获取输入
	fmt.Print("請選擇風格 [默認: dan-koe]: ")
	styleInput := readLine()
	if styleInput == "" {
		styleInput = "dan-koe"
	}

	fmt.Print("请输入你的观点或内容 (Ctrl+D 结束):\n")
	input := readMultiline()
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("输入不能为空")
	}

	// 构建请求
	req := &writer.WriteRequest{
		Input:     input,
		InputType: writer.GetInputTypeFromString(writeInputType),
		StyleName: styleInput,
		Length:    writer.GetLengthFromString(writeLength),
	}

	// 执行写作
	result := asst.Write(req)

	if result.IsAIRequest {
		// AI 模式：返回提示词
		output := map[string]interface{}{
			"success": true,
			"mode":    "ai",
			"action":  "ai_write_request",
			"style":   result.Style.Name,
			"prompt":  result.Prompt,
		}

		if writeCover {
			coverGen := writer.NewCoverGenerator(asst.GetStyleManager())
			coverResult, _ := coverGen.GeneratePrompt(&writer.GenerateCoverRequest{
				StyleName:      styleInput,
				ArticleContent: input,
			})
			if coverResult.Success {
				output["cover_prompt"] = coverResult.Prompt
			}
		}

		printJSON(output)
		return nil
	}

	if !result.Success {
		return fmt.Errorf("%s", result.Error)
	}

	// 输出结果
	if writeOutput != "" {
		if err := os.WriteFile(writeOutput, []byte(result.Article), 0644); err != nil {
			return fmt.Errorf("保存文件: %w", err)
		}
		log.Info("article saved", zap.String("file", writeOutput))
	} else {
		fmt.Println("\n=== 生成文章 ===")
		fmt.Println(result.Article)
		fmt.Println("\n=== 金句 ===")
		for i, quote := range result.Quotes {
			fmt.Printf("%d. %s\n", i+1, quote)
		}
	}

	return nil
}

// runIdeaGenerator 執行靈感生成
func runIdeaGenerator() error {
	log.Info("啟動靈感產生器...")
	
	// 初始化 Gemini Client
	client, err := llm.NewGeminiClient()
	if err != nil {
		return fmt.Errorf("初始化 AI 失敗 (請檢查 GEMINI_API_KEY): %w", err)
	}
	defer client.Close()

	// 選擇風格
	asst := writer.NewAssistant()
	fmt.Println("\n🧠 靈感產生器 - 根據風格提供選題建議")
	fmt.Print("請選擇風格 (例如 dan-koe, taiwan-ecommerce) [默認: dan-koe]: ")
	styleName := readLine()
	if styleName == "" {
		styleName = "dan-koe"
	}

	style, err := asst.GetStyleInfo(styleName)
	if err != nil {
		return fmt.Errorf("風格不存在: %w", err)
	}

	fmt.Printf("\n正在召喚 %s 風格的靈感 muse...\n", style.Name)

	// 構建提示詞
	prompt := fmt.Sprintf(`
你是一位專業的內容策略專家，熟悉 "%s" 的寫作風格。
該風格的核心描述為："%s"。
寫作 DNA：
%s

請為我生成 5 個「高病毒傳播潛力」的寫作主題/標題。
這些主題必須非常符合該風格的受眾（例如 Dan Koe 針對創作者/自律，台灣電商針對消費者痛點）。

輸出格式要求：
直接輸出 5 行，每行一個主題，不要有編號或其他廢話。
例如：
如何在新的一年徹底擺脫拖延
為什麼原本的努力方向都錯了
(略)
`, style.Name, style.Description, strings.Join(style.CoreBeliefs, "\n"))

	// 調用 AI
	ideasText, err := client.GenerateContent(context.Background(), prompt)
	if err != nil {
		return fmt.Errorf("生成靈感失敗: %w", err)
	}

	// 解析與顯示
	ideas := strings.Split(strings.TrimSpace(ideasText), "\n")
	var validIdeas []string
	for _, idea := range ideas {
		idea = strings.TrimSpace(idea)
		// 移除可能的編號 (1. , - )
		idea = strings.TrimLeft(idea, "1234567890.- ") 
		if idea != "" {
			validIdeas = append(validIdeas, idea)
		}
	}

	if len(validIdeas) == 0 {
		return fmt.Errorf("AI 沒有產生有效的靈感，請重試")
	}

	fmt.Println("\n👇 以下是為您生成的靈感主題：")
	for i, idea := range validIdeas {
		fmt.Printf("[%d] %s\n", i+1, idea)
	}
	fmt.Println("[0] 不滿意，重新生成")
	fmt.Println("[q] 退出")

	fmt.Print("\n請選擇想要撰寫的主題編號: ")
	choice := readLine()

	if choice == "q" || choice == "Q" {
		return nil
	}
	
	if choice == "0" {
		return runIdeaGenerator() // 遞迴重試
	}

	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(validIdeas) {
		fmt.Println("無效的選擇")
		return nil
	}

	selectedTopic := validIdeas[idx-1]
	fmt.Printf("\n✅ 已選擇主題: %s\n", selectedTopic)
	fmt.Println("正在以此主題生成文章...\n")

	// 直接進入文章生成流程 (這次是真的調用 AI)
	return executeWriteWithAI(client, selectedTopic, styleName, writer.GetLengthFromString(writeLength))
}

// executeWriteWithAI 使用 AI 直接生成文章 (整合 Gemini)
func executeWriteWithAI(client llm.LLMClient, topic, styleName string, length writer.Length) error {
	asst := writer.NewAssistant()
	
	// 1. 構建 Prompt (使用 Assistant 既有邏輯)
	req := &writer.WriteRequest{
		Input:     topic,
		InputType: writer.InputTypeIdea,
		StyleName: styleName,
		Length:    length,
	}
	
	// 這一步只會拿到 Prompt，因為 IsAIRequest 原本是設計給外部的
	result := asst.Write(req) 
	if !result.IsAIRequest {
		return fmt.Errorf("無法構建 AI 提示詞")
	}
	
	prompt := result.Prompt
	log.Info("正在生成文章內容...", zap.String("字數", string(length)))

	// 2. 調用 Gemini
	ctx := context.Background()
	content, err := client.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("AI 生成失敗: %w", err)
	}

	// 3. 輸出結果
	fmt.Println("\n=== ✨ 生成結果 ✨ ===\n")
	fmt.Println(content)
	fmt.Println("\n=====================\n")

	// 4. 保存文件 (Optional)
	if writeOutput != "" {
		if err := os.WriteFile(writeOutput, []byte(content), 0644); err != nil {
			log.Error("保存失敗", zap.Error(err))
		} else {
			log.Info("文章已保存", zap.String("檔案", writeOutput))
		}
	} else {
		// 詢問是否保存
		fmt.Print("是否保存此文章? (y/n) [y]: ")
		save := readLine()
		if save == "" || strings.ToLower(save) == "y" {
			// 自動產生檔名
			filename := "output_" + strings.ReplaceAll(topic, " ", "_") + ".md"
			// 移除特殊字符
			filename = strings.ReplaceAll(filename, "/", "_")
			filename = strings.ReplaceAll(filename, "\\", "_")
			filename = strings.ReplaceAll(filename, "?", "")
			
			if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
				log.Error("保存失敗", zap.Error(err))
			} else {
				fmt.Printf("✅ 文章已保存至: %s\n", filename)
			}
		}
	}

	return nil
}

// executeWrite 执行写作
func executeWrite(input string) error {
	asst := writer.NewAssistant()

	req := &writer.WriteRequest{
		Input:     input,
		InputType: writer.GetInputTypeFromString(writeInputType),
		StyleName: writer.ParseStyleInput(writeStyle),
		Title:     writeTitle,
		Length:    writer.GetLengthFromString(writeLength),
	}

	result := asst.Write(req)

	if result.IsAIRequest {
		// AI 模式：返回提示词
		output := map[string]interface{}{
			"success": true,
			"mode":    "ai",
			"action":  "ai_write_request",
			"style":   result.Style.Name,
			"prompt":  result.Prompt,
		}

		if writeCover || writeCoverOnly {
			coverGen := writer.NewCoverGenerator(asst.GetStyleManager())
			coverResult, err := coverGen.GeneratePrompt(&writer.GenerateCoverRequest{
				StyleName:      req.StyleName,
				ArticleTitle:   req.Title,
				ArticleContent: input,
			})
			if err == nil && coverResult.Success {
				output["cover_prompt"] = coverResult.Prompt
				output["cover_explanation"] = coverResult.Explanation
			}
		}

		printJSON(output)
		return nil
	}

	if !result.Success {
		return fmt.Errorf("%s", result.Error)
	}

	// 只生成封面
	if writeCoverOnly {
		return generateCover(asst, req)
	}

	// 输出文章
	if writeOutput != "" {
		if err := os.WriteFile(writeOutput, []byte(result.Article), 0644); err != nil {
			return fmt.Errorf("保存文件: %w", err)
		}
		log.Info("article saved", zap.String("file", writeOutput))
	} else {
		fmt.Println("\n=== 生成文章 ===")
		fmt.Println(result.Article)
		fmt.Println("\n=== 金句 ===")
		for i, quote := range result.Quotes {
			fmt.Printf("%d. %s\n", i+1, quote)
		}
	}

	// 如果需要封面
	if writeCover {
		return generateCover(asst, req)
	}

	return nil
}

// generateCover 生成封面
func generateCover(asst *writer.Assistant, req *writer.WriteRequest) error {
	coverGen := writer.NewCoverGenerator(asst.GetStyleManager())

	coverReq := &writer.GenerateCoverRequest{
		StyleName:      req.StyleName,
		ArticleTitle:   req.Title,
		ArticleContent: req.Input,
	}

	result, err := coverGen.GeneratePrompt(coverReq)
	if err != nil {
		return fmt.Errorf("生成封面提示词: %w", err)
	}

	fmt.Println("\n=== 封面提示词 ===")
	fmt.Println(result.Prompt)

	if result.Explanation != "" {
		fmt.Println("\n---")
		fmt.Println("📖 隐喻说明:", result.Explanation)
	}

	return nil
}

// readLine 读取一行输入
func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// readMultiline 读取多行输入
func readMultiline() string {
	var lines []string
	for {
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
