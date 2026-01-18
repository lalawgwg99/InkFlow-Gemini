# 常見問題 (FAQ)

本文件收集了使用者最常遇到的問題和解決方案。

## 目錄

- [安裝問題](#安裝問題)
- [設定問題](#設定問題)
- [轉換問題](#轉換問題)
- [圖片問題](#圖片問題)
- [微信 API 問題](#微信-api-問題)
- [進階問題](#進階問題)

---

## 安裝問題

### Q1: 提示 "command not found: md2wechat"

**原因**：二進位檔案不在 PATH 中

**解決方案 A**：加入到 PATH

```bash
# 暫時加入
export PATH=$PATH:/usr/local/bin

# 永久加入（加入到 ~/.bashrc）
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

**解決方案 B**：使用完整路徑

```bash
/usr/local/bin/md2wechat --help
```

---

### Q2: Go 安裝失敗 "go: module not found"

**原因**：網路問題或模組路徑錯誤

**解決方案**：

```bash
# 1. 設定 Go 代理（中國大陸使用者）
export GOPROXY=https://goproxy.cn,direct

# 2. 清理快取
go clean -modcache

# 3. 重新安裝
go install github.com/geekjourneyx/md2wechat-skill/cmd/md2wechat@latest
```

---

### Q3: macOS 提示 "無法開啟，因為無法驗證開發者"

**原因**：macOS 安全限制

**解決方案**：

```bash
# 允許任何來源的應用程式（系統偏好設定 > 安全性與隱私權）
# 或使用命令列
sudo xattr -cr /Applications/md2wechat
```

---

## 設定問題

### Q4: 提示 "WECHAT_APPID is required"

**原因**：未設定微信公眾號憑證

**解決方案 A**：使用設定檔

```bash
md2wechat config init
# 編輯產生的 md2wechat.yaml，填入：
# wechat:
#   appid: "your_appid"
#   secret: "your_secret"
```

**解決方案 B**：使用環境變數

```bash
export WECHAT_APPID="wx1234567890abcdef"
export WECHAT_SECRET="your_secret_here"
```

**如何取得 AppID 和 Secret**：

1. 前往 **[微信開發者平台](https://developers.weixin.qq.com/platform)**

2. 登入後，選擇你的公眾號（如果沒有，請先註冊）

3. 點擊左側選單 **「設定與開發」** → **「基本設定」**

4. 在「開發者ID」區域可以看到：
   - **開發者ID(AppID)**：直接複製即可
   - **開發者密碼(AppSecret)**：點擊「重設」按鈕取得

   > **警告**：AppSecret 非常重要，請妥善保管，不要洩露給他人！

5. 複製這兩個值到設定檔或環境變數中

---

### Q5: 設定檔不生效

**原因**：設定檔位置或格式錯誤

**解決方案**：

```bash
# 1. 檢查設定檔位置
md2wechat config show

# 2. 驗證設定檔格式
cat md2wechat.yaml

# 3. 重新初始化設定
md2wechat config init
```

**支援的設定檔位置**：

- `./md2wechat.yaml`（目前目錄，優先順序最高）
- `~/.md2wechat.yaml`
- `~/.config/md2wechat/config.yaml`

---

### Q6: API 模式提示需要 API Key

**原因**：API 模式需要 [md2wechat.cn](https://md2wechat.cn) 的 API Key

**解決方案 A**：取得 API Key

1. 前往 [md2wechat.cn](https://md2wechat.cn)
2. 註冊帳號並取得 API Key
3. 設定：

```bash
export MD2WECHAT_API_KEY="your_key"
```

**解決方案 B**：使用 AI 模式（不需要 md2wechat API Key）

```bash
md2wechat convert article.md --mode ai --theme autumn-warm
```

---

## 轉換問題

### Q7: 轉換結果為空或亂碼

**可能原因 1**：Markdown 檔案編碼問題

**解決方案**：

```bash
# 檢查檔案編碼
file article.md

# 轉換為 UTF-8
iconv -f GBK -t UTF-8 article.md > article-utf8.md
```

**可能原因 2**：Markdown 格式錯誤

**解決方案**：使用 Markdown 檢查工具

```bash
# 安裝 markdownlint
npm install -g markdownlint-cli

# 檢查檔案
markdownlint article.md
```

---

### Q8: AI 模式轉換失敗

**原因**：AI API Key 未設定或無效

**解決方案**：

```bash
# 1. 設定 Claude API Key
export IMAGE_API_KEY="your_claude_api_key"

# 2. 驗證
md2wechat config validate

# 3. 重試
md2wechat convert article.md --mode ai
```

---

### Q9: HTML 樣式在微信中顯示不正常

**原因**：微信編輯器會過濾部分 CSS

**解決方案**：

1. **使用 API 模式**（更穩定）

```bash
md2wechat convert article.md --mode api
```

1. **檢查是否使用了內聯樣式**

微信只支援內聯 style 屬性，不支援 `<style>` 標籤。

1. **嘗試簡化 Markdown**

```markdown
<!-- 避免複雜巢狀 -->
## 簡單標題

這是段落內容。

> 這是引用
```

---

## 圖片問題

### Q10: 圖片上傳失敗 "upload material failed"

**可能原因 1**：圖片格式不支援

**解決方案**：

```bash
# 支援的格式：jpg, png, gif, bmp, webp
# 轉換圖片格式
convert input.tiff output.jpg
```

**可能原因 2**：圖片太大

**解決方案**：

```bash
# 程式會自動壓縮，但可以先手動壓縮
# 設定壓縮參數
# md2wechat.yaml:
image:
  compress: true
  max_width: 1920
  max_size_mb: 5
```

**可能原因 3**：微信 API 頻率限制

**解決方案**：等待幾分鐘後重試

```bash
# 分批上傳圖片
md2wechat upload_image image1.jpg
sleep 5
md2wechat upload_image image2.jpg
```

---

### Q11: AI 產生圖片失敗

**原因**：圖片產生 API Key 未設定或額度不足

**解決方案**：

```bash
# 1. 設定 API Key
export IMAGE_API_KEY="your_openai_or_claude_key"

# 2. 驗證
md2wechat generate_image "test prompt"

# 3. 檢查 API 額度
# 登入對應的 API 提供商查看剩餘額度
```

---

### Q12: 圖片連結未被替換

**原因**：未使用 `--upload` 參數

**解決方案**：

```bash
# 必須使用 --upload 參數
md2wechat convert article.md --upload -o output.html
```

---

## 微信 API 問題

### Q13: 第一次呼叫 API 提示 "IP 不在白名單中"

**現象**：第一次呼叫微信 API 時，回傳錯誤：

```
ip xxx.xxx.xxx.xxx not in whitelist
```

**原因**：微信為了安全，要求伺服器 IP 必須在白名單中才能呼叫 API。

**解決方案**：

1. **取得你的伺服器 IP 位址**

```bash
# 查看你的公網 IP
curl ifconfig.me
# 或
curl ip.sb
# 或
curl ipinfo.io/ip
```

1. **將 IP 加入微信白名單**

   - 前往 [微信開發者平台](https://developers.weixin.qq.com/platform)
   - 選擇你的公眾號
   - 點擊 **「設定與開發」** → **「基本設定」**
   - 找到 **「IP白名單」** 區域
   - 點擊「設定」
   - 輸入你的伺服器 IP 位址（多個 IP 用換行分隔）
   - 點擊「確定」儲存

2. **等待生效並重試**

```bash
# 白名單設定通常幾分鐘內生效
# 等待 5 分鐘後重試
sleep 300
md2wechat convert article.md --upload --draft
```

> **注意**：
>
> - 如果你使用本機電腦測試，需要加入你本機網路的公網 IP
> - 如果使用雲端伺服器，加入雲端伺服器的公網 IP
> - 如果使用 GitHub Actions 等動態 IP 環境，建議使用固定 IP 的伺服器

---

### Q14: 提示 "access_token expired"

**原因**：微信 access_token 過期（通常 2 小時）

**解決方案**：程式會自動重新整理，如果持續失敗：

```bash
# 1. 檢查 AppID 和 Secret 是否正確
md2wechat config show --show-secret

# 2. 重新設定
md2wechat config init
```

---

### Q15: 草稿建立失敗 "create draft failed"

**可能原因 1**：公眾號權限不足

**解決方案**：

- 確保公眾號已認證
- 確保有素材管理權限
- 檢查是否超過草稿數量限制

**可能原因 2**：內容包含敏感詞

**解決方案**：

- 檢查文章內容
- 嘗試簡化內容後重試

```bash
# 先儲存為 JSON，檢查內容
md2wechat convert article.md --save-draft draft.json
cat draft.json
```

---

### Q16: API 呼叫頻率限制

**現象**：提示 "api freq limit"

**解決方案**：

```bash
# 方案 1：等待後重試
sleep 60

# 方案 2：分批處理
for file in articles/*.md; do
  md2wechat convert "$file" --draft
  sleep 5
done
```

---

## 進階問題

### Q17: 如何在 CI/CD 中使用？

**解決方案**：使用環境變數或 Secrets

```yaml
# .github/workflows/publish.yml
name: Publish to WeChat
on: [push]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install md2wechat
        run: go install github.com/geekjourneyx/md2wechat-skill/cmd/md2wechat@latest
      - name: Convert and publish
        env:
          WECHAT_APPID: ${{ secrets.WECHAT_APPID }}
          WECHAT_SECRET: ${{ secrets.WECHAT_SECRET }}
        run: |
          md2wechat convert article.md --upload --draft
```

---

### Q18: 如何自訂主題？

**解決方案**：使用 custom-prompt

```bash
md2wechat convert article.md --mode ai --custom-prompt "
請使用以下配色：
- 主色：#e53e3e（紅色）
- 副色：#3182ce（藍色）
- 背景：#f7fafc（淺灰）
- 字體：16px，行高 1.8

請確保：
1. 所有 CSS 使用內聯 style
2. 圖片使用佔位符 <!-- IMG:index -->
3. 不使用外部樣式表
"
```

---

### Q19: 如何批次轉換多個檔案？

**解決方案**：使用 Shell 腳本

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

  echo "Waiting 5 seconds..."
  sleep 5
done
```

---

### Q19: 如何除錯問題？

**解決方案**：啟用詳細記錄

```bash
# 方案 1：查看指令輸出
md2wechat convert article.md --preview 2>&1 | tee debug.log

# 方案 2：逐步測試
md2wechat config validate       # 1. 驗證設定
md2wechat upload_image test.jpg  # 2. 測試圖片上傳
md2wechat convert test.md --preview  # 3. 測試轉換
```

---

### Q20: 如何取得協助？

**解決方案**：

1. **查看指令說明**

```bash
md2wechat --help
md2wechat convert --help
```

1. **查看文件**

- [安裝指南](INSTALL.md)
- [設定指南](CONFIG.md)
- [使用教學](USAGE.md)

1. **提交 Issue**

前往 [GitHub Issues](https://github.com/geekjourneyx/md2wechat-skill/issues)

---

## 仍然無法解決？

請提供以下資訊：

1. **版本資訊**

   ```bash
   md2wechat --version
   go version
   ```

2. **設定資訊**

   ```bash
   md2wechat config show
   ```

3. **錯誤資訊**

   ```bash
   md2wechat convert article.md 2>&1
   ```

4. **系統資訊**

   ```bash
   uname -a  # Linux/macOS
   # 或
   systeminfo  # Windows
   ```

將以上資訊提交到 [GitHub Issues](https://github.com/geekjourneyx/md2wechat-skill/issues)，我們會儘快回覆。
