# 靜態檔案結構說明

## 資料夾結構

```
static/
├── html/           # HTML 頁面檔案
│   ├── index.html          # 首頁
│   └── chat_websocket.html # WebSocket 聊天室頁面
├── css/            # CSS 樣式檔案
├── js/             # JavaScript 檔案
└── assets/         # 靜態資源 (圖片、字型等)
```

## 路由對應

| 路由路徑 | 檔案位置 | 說明 |
|---------|---------|------|
| `/` | `static/html/index.html` | 首頁 |
| `/chat` | `static/html/chat_websocket.html` | 聊天室（新路由） |
| `/chat_websocket.html` | `static/html/chat_websocket.html` | 聊天室（向後兼容） |
| `/static/*` | `static/*` | 靜態檔案服務 |

## 使用方式

### 直接訪問頁面

- 首頁：<http://localhost:8080/>
- 聊天室：<http://localhost:8080/chat>

### 訪問靜態資源

- CSS 檔案：<http://localhost:8080/static/css/style.css>
- JavaScript 檔案：<http://localhost:8080/static/js/app.js>
- 圖片資源：<http://localhost:8080/static/assets/logo.png>

## 開發建議

1. **HTML 檔案**：放在 `static/html/` 目錄
2. **CSS 檔案**：放在 `static/css/` 目錄
3. **JavaScript 檔案**：放在 `static/js/` 目錄
4. **圖片、字型等資源**：放在 `static/assets/` 目錄

## 未來擴展

這個結構可以輕鬆支援：

- 多個頁面的網站
- 獨立的 CSS 和 JS 檔案
- 靜態資源的版本控制
- 前端構建工具的整合
