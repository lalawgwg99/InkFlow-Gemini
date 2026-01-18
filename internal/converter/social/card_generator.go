package social

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// CardGenerator 負責將 Markdown 轉換為適合 IG/FB 的卡片式 HTML
type CardGenerator struct {
	Width  int
	Height int
}

// NewCardGenerator 建立新的卡片生成器 (預設直式 1080x1350)
func NewCardGenerator() *CardGenerator {
	return &CardGenerator{
		Width:  1080,
		Height: 1350,
	}
}

// CardSlide 代表一張卡片的數據
type CardSlide struct {
	ID      int
	Content template.HTML // 允許 HTML 標籤
	Type    string        // "cover", "body", "quote", "end"
	Theme   string        // 主題 class 名稱
}

// GenerateHTML 將文章轉換為卡片式 HTML 頁面
func (g *CardGenerator) GenerateHTML(markdown string, themeName string) (string, error) {
	// 簡單的 Markdown 解析與分頁邏輯
	// 這裡先做個簡單的分頁：按 H2 標題或分隔線分頁
	
	slides := g.parseMarkdownToSlides(markdown, themeName)
	
	// 渲染 HTML 模板
	tmpl, err := template.New("carousel").Parse(carouselTemplate)
	if err != nil {
		return "", fmt.Errorf("compile template failed: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Slides": slides,
		"Theme":  themeName,
		"Width":  g.Width,
		"Height": g.Height,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template failed: %w", err)
	}

	return buf.String(), nil
}

func (g *CardGenerator) parseMarkdownToSlides(markdown, themeName string) []CardSlide {
	var slides []CardSlide
	
	// 預處理
	lines := strings.Split(markdown, "\n")
	var currentContent strings.Builder
	slideCount := 1

	// 輔助函數：添加一張卡片
	addSlide := func(contentType string) {
		htmlContent := simpleMarkdownToHTML(currentContent.String())
		slides = append(slides, CardSlide{
			ID:      slideCount,
			Content: template.HTML(htmlContent),
			Type:    contentType,
			Theme:   themeName,
		})
		slideCount++
		currentContent.Reset()
	}

	for _, line := range lines {
		// 遇到 H1 或 H2 視為新卡片的開始 (除了第一張)
		if (strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ")) && currentContent.Len() > 0 {
			addSlide("body")
		}
		
		// 遇到分隔線 --- 強制分頁
		if strings.TrimSpace(line) == "---" {
			if currentContent.Len() > 0 {
				addSlide("body")
			}
			continue
		}

		currentContent.WriteString(line + "\n")
	}

	// 加入最後一張
	if currentContent.Len() > 0 {
		addSlide("end")
	}

	return slides
}

// 簡單的 Markdown 轉 HTML (暫時替代完整 Parser，後續可接 goldmark)
func simpleMarkdownToHTML(md string) string {
	html := md
	html = strings.ReplaceAll(html, "\n", "<br>")
	// 簡單處理標題，去除 # 號，改為加粗或大字體 class
	if strings.Contains(html, "# ") {
		// 實際專案建議用正規 Markdown Parser，這裡僅為示範分頁邏輯
	}
	return html
}

// carouselTemplate 是生成卡片預覽的 HTML 模板
// 包含 CSS Scroll Snap，方便在瀏覽器預覽並截圖
const carouselTemplate = `
<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>IG Carousel Preview</title>
    <style>
        :root {
            --card-width: {{.Width}}px;
            --card-height: {{.Height}}px;
            --bg-color: #1a1a1a;
            --text-color: #ffffff;
            --accent-color: #d4af37;
        }

        /* Reset */
        * { box-sizing: border-box; margin: 0; padding: 0; }

        body {
            background-color: #333;
            font-family: 'Source Han Serif', 'Noto Serif TC', serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            padding: 20px;
        }

        .carousel-container {
            display: flex;
            gap: 20px;
            overflow-x: auto;
            scroll-snap-type: x mandatory;
            padding-bottom: 20px;
            max-width: 100%;
        }

        .slide {
            flex: 0 0 auto;
            width: var(--card-width);
            height: var(--card-height);
            background-color: var(--bg-color);
            color: var(--text-color);
            position: relative;
            scroll-snap-align: center;
            overflow: hidden;
            border-radius: 8px; /* 預覽時圓角，輸出時可去掉 */
            box-shadow: 0 10px 30px rgba(0,0,0,0.5);
            /* 縮放以適應螢幕預覽 */
            zoom: 0.5; 
        }

        /* 內容佈局 */
        .slide-content {
            padding: 80px;
            height: 100%;
            display: flex;
            flex-direction: column;
            justify-content: center;
        }

        .slide-footer {
            position: absolute;
            bottom: 40px;
            left: 0;
            width: 100%;
            text-align: center;
            font-size: 24px;
            color: var(--accent-color);
            opacity: 0.8;
        }

        /* Theme: magazine-dark */
        .theme-magazine-dark {
            background-color: #1a1a1a;
            color: #f0f0f0;
            background-image: url("data:image/svg+xml,%3Csvg width='100' height='100' viewBox='0 0 100 100' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M11 18c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm48 25c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm-43-7c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm63 31c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM34 90c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm56-76c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM12 86c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm28-65c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm23-11c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-6 60c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm29 22c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 2.24 5 5 5zM32 63c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 2.24 5 5 5zm57-13c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 2.24 5 5 5zm-9-21c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM60 91c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM35 41c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM12 60c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2z' fill='%23d4af37' fill-opacity='0.05' fill-rule='evenodd'/%3E%3C/svg%3E");
            /* 微動態：背景緩慢移動 */
            animation: bgMove 60s linear infinite;
        }
        @keyframes bgMove {
            0% { background-position: 0 0; }
            100% { background-position: 100% 100%; }
        }

        .theme-magazine-dark h1 {
            color: #d4af37;
            font-size: 64px;
            margin-bottom: 60px;
            line-height: 1.2;
            font-family: 'Times New Roman', serif;
            letter-spacing: -1px;
            /* 金屬光澤文字 */
            background: linear-gradient(45deg, #d4af37, #f0e68c, #d4af37);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            text-shadow: 0 2px 10px rgba(212, 175, 55, 0.3);
        }

        .theme-magazine-dark h2 {
            color: #f0f0f0;
            font-size: 48px;
            margin-top: 40px;
            margin-bottom: 30px;
            border-left: 6px solid #d4af37;
            padding-left: 20px;
        }

        .theme-magazine-dark p {
            font-size: 32px;
            line-height: 1.8;
            margin-bottom: 30px;
            opacity: 0.9;
            color: #e0e0e0;
        }

        .theme-magazine-dark strong {
            color: #d4af37;
            font-weight: 700;
        }

        /* 引用樣式 */
        .theme-magazine-dark blockquote {
            border-top: 2px solid #333;
            border-bottom: 2px solid #333;
            padding: 40px 0;
            margin: 40px 0;
            font-style: italic;
            font-size: 36px;
            color: #a0a0a0;
            text-align: center;
        }

        /* Theme: minimalist-gray */
        .theme-minimalist-gray {
            --bg-color: #f8f8f8;
            --text-color: #333333;
            --accent-color: #666666;
            background-color: #f8f8f8;
            background-image: radial-gradient(#e0e0e0 1px, transparent 1px);
            background-size: 40px 40px;
        }
        
        .theme-minimalist-gray .slide-content {
            padding: 100px; /* 更大的留白 */
        }

        .theme-minimalist-gray h1 {
            font-weight: 300;
            color: #000;
            letter-spacing: 2px;
            margin-bottom: 80px;
            font-size: 56px;
        }

        .theme-minimalist-gray p {
            font-weight: 300;
            color: #444;
            font-size: 30px;
            line-height: 2.2;
        }
        
        .theme-minimalist-gray strong {
             background: linear-gradient(0deg, #e0e0e0 0%, #e0e0e0 40%, transparent 40%);
             font-weight: 500;
             color: #000;
        }

        /* Theme: tech-neon */
        .theme-tech-neon {
            background-color: #050510;
            color: #ffffff;
            font-family: 'JetBrains Mono', 'Courier New', monospace;
            /* 網格背景 */
            background-image: linear-gradient(#00f3ff 1px, transparent 1px), linear-gradient(90deg, #00f3ff 1px, transparent 1px);
            background-size: 50px 50px;
        }

        .theme-tech-neon .slide-content {
            background: rgba(5, 5, 16, 0.9); /* 半透明遮罩讓文字清晰 */
            margin: 20px;
            border: 1px solid #333;
            box-shadow: 0 0 20px rgba(0, 243, 255, 0.15);
        }

        .theme-tech-neon h1 {
            color: #00f3ff;
            font-size: 50px;
            text-transform: uppercase;
            text-shadow: 0 0 15px rgba(0, 243, 255, 0.8);
            border-left: 10px solid #bc13fe;
            padding-left: 30px;
        }

        .theme-tech-neon p {
            font-size: 28px;
            color: #e0e0e0;
            line-height: 1.6;
        }

        .theme-tech-neon strong {
            color: #bc13fe;
            text-shadow: 0 0 5px rgba(188, 19, 254, 0.6);
        }

        /* Hover Interactive Effects (Web Preview Only) */
        .slide:hover {
            transform: scale(1.02);
            transition: transform 0.3s ease;
            z-index: 10;
        }
        
        .theme-magazine-dark:hover h1 {
            background-position: 100% 0;
            transition: background-position 0.5s;
        }
        
        .theme-minimalist-gray:hover {
            background-size: 45px 45px; /* Grid expands */
            transition: background-size 0.5s ease;
        }

        .theme-tech-neon:hover {
            box-shadow: 0 0 50px rgba(0, 243, 255, 0.3);
            transition: box-shadow 0.3s ease;
        }
    </style>
</head>
<body>
    <div class="carousel-container">
        {{range .Slides}}
        <div class="slide theme-{{.Theme}}">
            <div class="slide-content">
                {{.Content}}
            </div>
            <div class="slide-footer">
                {{.ID}} / {{len $.Slides}}
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>
`
