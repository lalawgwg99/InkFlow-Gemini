# 使用教學

本文件詳細說明 md2wechat 的各種使用方式。

## 目錄

- [Claude Code 整合](#claude-code-整合)
- [基礎使用](#基礎使用)
- [轉換模式](#轉換模式)
- [圖片處理](#圖片處理)
- [主題訂製](#主題訂製)
- [草稿管理](#草稿管理)
- [完整範例](#完整範例)

---

## Claude Code 整合

### 安裝（最簡單）

在 Claude Code 中執行：

```bash
/plugin marketplace add geekjourneyx/md2wechat-skill
/plugin install md2wechat@geekjourneyx-md2wechat-skill
```

### 使用方式

安裝後，直接與 Claude 對話即可：

```
請用秋日暖光主題將 article.md 轉換為微信公眾號格式，並上傳到草稿箱
```

```
幫我把這篇技術文章用深海靜謐主題轉換，預覽效果給我看
```

Claude 會自動：

1. 呼叫 md2wechat 轉換 Markdown
2. 套用你選擇的主題
3. 上傳圖片到微信
4. 建立草稿或顯示預覽

---

## 基礎使用

### 最簡單的例子

```bash
# 預覽轉換結果（不上傳圖片）
md2wechat convert article.md --preview
```

### 常用指令組合

```bash
# 1. 預覽模式 - 快速查看效果
md2wechat convert article.md --preview

# 2. 儲存到檔案
md2wechat convert article.md -o output.html

# 3. 上傳圖片並輸出 HTML
md2wechat convert article.md --upload -o output.html

# 4. 完整流程 - 上傳圖片 + 建立草稿
md2wechat convert article.md --upload --draft
```

---

## 轉換模式

### API 模式（推薦新手）

使用 md2wechat.cn API 進行轉換，穩定可靠。

```bash
md2wechat convert article.md --mode api --api-key "your_key"
```

**特點**：

- 轉換速度快
- 結果穩定一致
- 需要註冊 API Key

**可用主題**：

- `default` - 預設主題
- `bytedance` - 字節跳動風格
- `apple` - Apple 極簡風格
- `sports` - 運動活力風格
- `chinese` - 中國傳統文化風格
- `cyber` - 賽博龐克風格

### AI 模式（適合訂製）

使用 AI 產生 HTML，更加靈活。

```bash
md2wechat convert article.md --mode ai --theme autumn-warm
```

**特點**：

- 高度可訂製
- 主題更精美
- 需要 AI API Key

**可用主題**：

- `autumn-warm` - 秋日暖光
- `spring-fresh` - 春日清新
- `ocean-calm` - 深海靜謐
- `custom` - 自訂

### 模式對比

| 特性 | API 模式 | AI 模式 |
|------|---------|---------|
| 速度 | 快 | 較慢 |
| 穩定性 | 高 | 中 |
| 主題選擇 | 基礎 | 豐富 |
| 成本 | 需要 API Key | 需要 AI Key |
| 適用場景 | 日常使用 | 追求美觀 |

---

## 圖片處理

### 圖片語法

在 Markdown 中使用標準圖片語法：

```markdown
<!-- 本機圖片：會上傳到微信 -->
![圖片描述](./images/photo.jpg)

<!-- 線上圖片：會先下載再上傳 -->
![圖片描述](https://example.com/image.jpg)

<!-- AI 產生圖片：會呼叫 API 產生 -->
![圖片描述](__generate:A cute orange cat__)
```

### 自動上傳

```bash
# 自動上傳所有圖片
md2wechat convert article.md --upload

# 上傳並替換 HTML 中的圖片連結
md2wechat convert article.md --upload -o output.html
```

### 手動上傳單個圖片

```bash
# 上傳本機圖片
md2wechat upload_image ./photo.jpg

# 下載並上傳線上圖片
md2wechat download_and_upload https://example.com/image.jpg
```

### AI 產生圖片

```bash
# 產生圖片並上傳
md2wechat generate_image "A beautiful sunset over mountains"
```

輸出範例：

```json
{
  "success": true,
  "prompt": "A beautiful sunset over mountains",
  "media_id": "12345***6789",
  "wechat_url": "http://mmbiz.qpic.cn/..."
}
```

### 圖片壓縮

程式會自動壓縮超過限制的圖片：

- 寬度超過 1920px → 等比縮放到 1920px
- 大小超過 5MB → 壓縮品質
- 格式轉換 → PNG → JPEG（可選）

設定壓縮參數：

```yaml
# md2wechat.yaml
image:
  compress: true
  max_width: 1920      # 最大寬度
  max_size_mb: 5       # 最大大小（MB）
```

---

## 主題訂製

### 使用內建主題

```bash
# 秋日暖光
md2wechat convert article.md --mode ai --theme autumn-warm

# 春日清新
md2wechat convert article.md --mode ai --theme spring-fresh

# 深海靜謐
md2wechat convert article.md --mode ai --theme ocean-calm
```

### 主題預覽

| 主題 | 色調 | 風格 |
|------|------|------|
| autumn-warm | 橙色 | 溫暖療癒 |
| spring-fresh | 綠色 | 生機盎然 |
| ocean-calm | 藍色 | 理性專業 |

### 自訂提示詞

```bash
md2wechat convert article.md --mode ai --custom-prompt "
請使用藍色配色方案，建立專業的技術部落格風格。
標題使用深藍色 #1a365d，正文使用 #2d3748。
"
```

### 設定預設主題

在設定檔中設定：

```yaml
api:
  default_theme: "autumn-warm"  # 設定預設主題
```

---

## 草稿管理

### 建立微信草稿

```bash
# 直接建立草稿
md2wechat convert article.md --draft

# 先上傳圖片再建立草稿
md2wechat convert article.md --upload --draft
```

### 儲存草稿 JSON

```bash
# 儲存草稿到檔案（不提交到微信）
md2wechat convert article.md --save-draft draft.json

# 查看草稿檔案
cat draft.json
```

草稿 JSON 格式：

```json
{
  "articles": [
    {
      "title": "文章標題",
      "content": "<section>...</section>",
      "digest": "文章摘要..."
    }
  ]
}
```

### 從 JSON 建立草稿

```bash
md2wechat create_draft draft.json
```

---

## 完整範例

### 範例 1：新手入門

```bash
# 1. 首次使用，初始化設定
md2wechat config init
# 編輯 md2wechat.yaml，填入微信 AppID 和 Secret

# 2. 驗證設定
md2wechat config validate

# 3. 預覽轉換
md2wechat convert my-article.md --preview

# 4. 建立草稿
md2wechat convert my-article.md --draft
```

### 範例 2：使用精美主題

```bash
# 1. 使用 AI 模式 + 秋日暖光主題
md2wechat convert my-article.md \
  --mode ai \
  --theme autumn-warm \
  --preview

# 2. 滿意後，上傳圖片並建立草稿
md2wechat convert my-article.md \
  --mode ai \
  --theme autumn-warm \
  --upload \
  --draft
```

### 範例 3：批次處理

```bash
#!/bin/bash
# batch-convert.sh

for file in articles/*.md; do
  echo "Converting $file..."
  md2wechat convert "$file" \
    --mode ai \
    --theme autumn-warm \
    --upload \
    --draft
done
```

### 範例 4：CI/CD 整合

```bash
#!/bin/bash
# .github/workflows/publish.yml

# 設定環境變數
export WECHAT_APPID="${{ secrets.WECHAT_APPID }}"
export WECHAT_SECRET="${{ secrets.WECHAT_SECRET }}"

# 轉換並建立草稿
md2wechat convert article.md \
  --upload \
  --draft \
  --save-draft /outputs/draft.json
```

---

## 進階技巧

### 組合使用模式

```bash
# 使用 API 模式轉換，但用 AI 模式的主題提示詞
md2wechat convert article.md \
  --mode api \
  --custom-prompt "參考 autumn-warm 主題的配色"
```

### 僅處理圖片

```bash
# 提取所有圖片連結
md2wechat convert article.md --preview | grep IMG

# 上傳所有圖片並儲存 URL
md2wechat convert article.md --upload -o temp.html
```

### 除錯模式

```bash
# 查看詳細記錄
md2wechat convert article.md --preview 2>&1 | tee debug.log
```

---

## 故障排除

### 問題：轉換結果為空

**原因**：Markdown 內容為空或格式錯誤

**解決**：

```bash
# 檢查檔案內容
cat article.md

# 檢查檔案編碼
file article.md
```

### 問題：圖片未替換

**原因**：未使用 `--upload` 參數

**解決**：

```bash
md2wechat convert article.md --upload -o output.html
```

### 問題：草稿建立失敗

**原因**：微信 API 權限不足或呼叫頻率限制

**解決**：

```bash
# 檢查設定
md2wechat config validate

# 先儲存 JSON，手動上傳
md2wechat convert article.md --save-draft draft.json
```

---

## 下一步

- 查看 [FAQ](FAQ.md) 了解常見問題
- 查看 [範例檔案](../examples/) 了解更多用法
