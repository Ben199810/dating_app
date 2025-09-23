# 技術

使用 DDD 事件驅動設計 (Domain-Driven Design) 架構來組織程式碼。

## 前端架構

- HTML
- CSS
- JavaScript

### 資料結構

```
.
├── static
│   ├── assets
│   ├── css
│   │   └── common.css # CSS 通用樣式
│   ├── html
│   └── js
```

### 開發建議

1. **HTML 檔案**：放在 `static/html/` 目錄
2. **CSS 檔案**：放在 `static/css/` 目錄
3. **JavaScript 檔案**：放在 `static/js/` 目錄
4. **圖片、字型等資源**：放在 `static/assets/` 目錄

## 後端架構

- golang
- gin
- websocket

### 資料結構

```
.
├── domain             # 領域層
│   ├── entity        # 實體定義
│   ├── repository    # 儲存庫介面
│   └── service       # 業務邏輯(聚合根)
├── infrastructure     # 基礎設施層
│   └── mysql         # MySQL 實作
├── server             # 伺服器層
│   └── handler       # HTTP 處理器
├── component          # 共用元件
│   └── validator     # 資料驗證
├── router.go          # 路由設定
└── main.go           # 入口程式
```
