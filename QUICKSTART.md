# 新手入門指南（5分鐘上手）

> **你不需要懂程式設計！** 按照下面的步驟操作即可。

---

## 🚀 Claude Code 使用者（最簡單）

如果你使用 **Claude Code**，只需兩步：

### 第一步：安裝外掛

在 Claude Code 中執行：

```bash
/plugin marketplace add geekjourneyx/md2wechat-skill
/plugin install md2wechat@geekjourneyx-md2wechat-skill
```

### 第二步：開始使用

直接和 Claude 對話：

```
請用秋日暖光主題將 article.md 轉換為微信公眾號格式
```

```
幫我把這篇技術文章轉換後上傳到微信草稿箱
```

完成！🎉

---

*如果你不使用 Claude Code，請繼續閱讀下面的內容。*

---

## 第一步：安裝軟體

### 選擇你的系統，點擊下載

| 你的系統 | 下載連結 | 安裝位置 |
|----------|----------|----------|
| Windows 10/11 | [下載 .exe](https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-windows-amd64.exe) | 任意資料夾或 `C:\Windows\System32\` |
| Mac (Intel晶片) | [下載](https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-darwin-amd64) | `/usr/local/bin/` 或 `~/.local/bin/` |
| Mac (M1/M2晶片) | [下載](https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-darwin-arm64) | `/usr/local/bin/` 或 `~/.local/bin/` |
| Linux | [下載](https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-linux-amd64) | `/usr/local/bin/` 或 `~/.local/bin/` |

---

### 安裝步驟（圖文說明）

#### Windows 使用者

1. 下載 `md2wechat-windows-amd64.exe`
2. 可以重新命名為 `md2wechat.exe`（方便輸入）
3. **方法 A（推薦）**：直接放到你想放的資料夾，用時開啟 CMD 切換到那個資料夾
4. **方法 B（全域可用）**：複製到 `C:\Windows\System32\`
5. 測試：開啟「命令提示字元」或「PowerShell」，輸入 `md2wechat --help`

#### Mac / Linux 使用者

**方法一：一鍵安裝（最簡單）**

```bash
# 複製這條指令，貼到終端機，按 Enter
curl -fsSL https://raw.githubusercontent.com/geekjourneyx/md2wechat-skill/main/scripts/install.sh | bash
```

**方法二：手動安裝**

```bash
# 1. 下載
curl -Lo md2wechat https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-linux-amd64

# 2. 新增執行權限
chmod +x md2wechat

# 3. 移動到系統目錄
sudo mv md2wechat /usr/local/bin/

# 4. 測試
md2wechat --help
```

**方法三：使用者目錄安裝（無需 sudo）**

```bash
# 1. 建立使用者 bin 目錄
mkdir -p ~/.local/bin

# 2. 下載到使用者目錄
curl -Lo ~/.local/bin/md2wechat https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-linux-amd64

# 3. 新增執行權限
chmod +x ~/.local/bin/md2wechat

# 4. 加入到 PATH（只需一次）
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc   # 如果你用 bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc    # 如果你用 zsh
source ~/.bashrc   # 或 source ~/.zshrc

# 5. 測試
md2wechat --help
```

---

### 驗證安裝

輸入以下指令，如果看到說明資訊，表示安裝成功：

```bash
md2wechat --help
```

---

## 第二步：設定微信（只需1次）

### 2.1 取得微信公眾號密碼

1. 用瀏覽器開啟：<https://mp.weixin.qq.com>
2. 登入你的公眾號
3. 點擊左上角「**設定與開發**」→「**基本設定**」
4. 複製這兩個資訊：
   - **開發者ID(AppID)**：類似 `wx1234567890abcdef`
   - **開發者密碼(AppSecret)**：點擊「重設」取得

### 2.2 產生設定檔

開啟「**終端機**」（Mac/Linux）或「**命令提示字元**」（Windows）：

```bash
# 輸入這個指令，按 Enter
md2wechat config init
```

這會建立一個 `md2wechat.yaml` 檔案，用記事本開啟它。

### 2.3 填寫你的資訊

用記事本開啟 `md2wechat.yaml`，修改這兩行：

```yaml
wechat:
  appid: "wx1234567890abcdef"    # ← 貼上你的 AppID
  secret: "your_secret_here"      # ← 貼上你的 Secret
```

儲存檔案，完成！

---

## 第三步：開始使用

### 3.1 準備你的文章

你的文章用 Markdown 格式寫，儲存為 `我的文章.md`

**什麼是 Markdown？**

- 一種簡單的寫作格式
- 用 `#` 表示標題
- 用 `![圖片](地址)` 插入圖片
- [Markdown 教學](https://www.markdownguide.org/zh-cn/basic-syntax/)

### 3.2 轉換文章

```bash
# 預覽效果（先看看怎麼樣）
md2wechat convert 我的文章.md --preview

# 滿意後，直接發送到微信草稿箱
md2wechat convert 我的文章.md --draft
```

### 3.3 在微信中查看

1. 登入微信公眾號後台
2. 點擊「**新的創作**」→「**草稿箱**」
3. 你的文章已經在那裡了！
4. 編輯後發表即可

---

## 常用指令一覽

| 你想做什麼 | 輸入這個指令 |
|-------------|---------------|
| 預覽文章 | `md2wechat convert 文章.md --preview` |
| 發送到草稿箱 | `md2wechat convert 文章.md --draft` |
| 使用精美主題 | `md2wechat convert 文章.md --mode ai --theme autumn-warm` |
| 查看設定 | `md2wechat config show` |
| 檢查設定是否正確 | `md2wechat config validate` |

---

## 精美主題推薦

| 指令 | 效果 |
|------|------|
| `--theme autumn-warm` | 🟠 秋日暖光（溫暖療癒） |
| `--theme spring-fresh` | 🟢 春日清新（生機盎然） |
| `--theme ocean-calm` | 🔵 深海靜謐（理性專業） |

**用法範例**：

```bash
md2wechat convert 我的文章.md --mode ai --theme autumn-warm --draft
```

---

## 遇到問題？

### 問題 1：提示「指令不存在」

**Windows**：把下載的 `md2wechat.exe` 放到 `C:\Windows\System32\` 資料夾

**Mac/Linux**：

```bash
# 給檔案執行權限
chmod +x md2wechat

# 移動到系統目錄
sudo mv md2wechat /usr/local/bin/
```

### 問題 2：提示「WECHAT_APPID is required」

表示你還沒設定，回到「第二步：設定微信」

### 問題 3：圖片沒有上傳

需要加 `--upload` 參數：

```bash
md2wechat convert 文章.md --upload --draft
```

---

## 完整範例

假設你有一篇文章叫 `產品發布.md`：

```bash
# 第一步：預覽效果
md2wechat convert 產品發布.md --mode ai --theme autumn-warm --preview

# 第二步：滿意後，上傳圖片並發送到草稿箱
md2wechat convert 產品發布.md --mode ai --theme autumn-warm --upload --draft
```

就這麼簡單！

---

## 下一步

- 查看 [使用教學](docs/USAGE.md) 了解更多功能
- 查看 [常見問題](docs/FAQ.md) 解決更多問題
