# CSS 架構說明

## 檔案結構

```text
static/css/
├── common.css    # 通用樣式（全域重置、按鈕、毛玻璃效果等）
├── index.css     # 首頁專用樣式
└── chat.css      # 聊天室專用樣式
```

## 使用方式

### 在 HTML 中引用

**首頁 (index.html):**

```html
<!-- 只需引用 index.css，會自動載入 common.css -->
<link rel="stylesheet" href="/static/css/index.css">
```

**聊天室 (chat_websocket.html):**

```html
<!-- 只需引用 chat.css，會自動載入 common.css -->
<link rel="stylesheet" href="/static/css/chat.css">
```

### CSS 檔案結構

每個頁面專用的 CSS 檔案都會透過 `@import` 自動引入通用樣式：

```css
/* 在 index.css 和 chat.css 的開頭 */
@import url('/static/css/common.css');
```

⚠️ **重要**:

- HTML 中只需引用頁面專用的 CSS 檔案
- `@import` 語句必須在 CSS 檔案的最開頭
- 通用樣式會自動載入，無需手動引用

## 通用樣式類別 (common.css)

### 按鈕樣式

```html
<!-- 主要按鈕 -->
<button class="btn btn-primary">主要按鈕</button>

<!-- 次要按鈕 -->
<button class="btn btn-secondary">次要按鈕</button>
```

### 毛玻璃效果

```html
<!-- 毛玻璃卡片（包含懸停效果） -->
<div class="glass-card">
    <p>內容</p>
</div>

<!-- 純毛玻璃效果（無懸停） -->
<div class="glass-effect">
    <p>內容</p>
</div>
```

### 狀態指示器

```html
<span class="status-indicator"></span>服務運行中
```

### 響應式工具類

```html
<!-- 在行動裝置上隱藏 -->
<div class="mobile-hidden">桌面版內容</div>

<!-- 在行動裝置上垂直排列 -->
<div class="mobile-stack">
    <div>項目1</div>
    <div>項目2</div>
</div>
```

## 整合後的變更

### 已移除的重複樣式

- ✅ 全域重置樣式 (`*` 選擇器)
- ✅ body 基本樣式
- ✅ 按鈕樣式 (`.btn`, `.btn-primary`, `.btn-secondary`)
- ✅ 狀態指示器動畫 (`.status-indicator`, `@keyframes pulse`)

### HTML 更新

- ✅ `index.html`: 添加 `glass-effect` 和 `glass-card` 類別
- ✅ `chat_websocket.html`: 更新 CSS 引用順序

## 最佳實踐

1. **CSS 模組化**: 使用 `@import` 在 CSS 檔案中引入依賴，而不是在 HTML 中多次引用
2. **類別使用**: 優先使用通用類別，避免重複樣式
3. **毛玻璃效果**:
   - 使用 `glass-card` 需要懸停動畫的卡片
   - 使用 `glass-effect` 純背景效果
4. **按鈕設計**:
   - 主要操作使用 `btn btn-primary`
   - 次要操作使用 `btn btn-secondary`

## 檔案大小優化

整合前:

- common.css: ~3KB
- index.css: ~8KB (包含重複樣式)

整合後:

- common.css: ~3.5KB (包含更多通用樣式)
- index.css: ~4KB (移除重複樣式)

**總減少**: ~3.5KB (約 30% 減少)
