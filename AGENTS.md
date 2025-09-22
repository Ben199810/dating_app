# 本機電腦啟動開發環境

參考 `DEVELOP.md` 文件中的說明來在本機電腦上啟動 Docker 開發環境。

# 專案編碼風格與規範

## 變數命名規範

- private 變數: 使用小駝峰式命名法 (e.g. `myVariable`)
- public 變數: 使用大駝峰式命名法 (e.g. `MyVariable`)

## 函式命名規範

- 動詞開頭，例如 `GetUser`, `SetUser`
- 私有函式使用小駝峰式命名法 (e.g. `getUser`)
- 公有函式使用大駝峰式命名法 (e.g. `GetUser`)

# 檔案結構

```
.
├── build
│   ├── docker-compose.yaml # 本機環境開發的 docker-compose 配置
│   ├── Dockerfile          # 應用程式的 Dockerfile
│   └── init.sql            # 資料庫初始化 SQL 腳本(本機環境開發使用)
├── config
│   ├── config.go        # 配置加載邏輯
│   ├── development.yaml # 本機開發環境配置
│   ├── production.yaml  # 生產環境配置
│   └── test.yaml        # 測試環境配置
├── domain               # 領域模型
│   ├── entity           # 實體
│   ├── repository       # 儲存庫介面
│   └── service          # 服務介面 (業務邏輯)
├── infrastructure       # 基礎設施層
│   └── mysql            # MySQL 相關
├── server               # 伺服器
│   └── handler          # HTTP 處理器
├── static               # 靜態資源
│   ├── css              # CSS 檔案
│   │   ├── common.css      # 共用樣式
│   │   ├── index.css       # 首頁樣式
│   │   ├── login.css       # 登入頁樣式
│   │   ├── register.css    # 註冊頁樣式
│   │   └── chat.css        # 聊天頁樣式
│   ├── html             # HTML 檔案
│   │   ├── chat_websocket.html # 聊天頁面
│   │   ├── index.html         # 首頁
│   │   └── register.html      # 註冊頁面
│   └── js               # JavaScript 檔案
├── .env.example        # 環境變數範例檔案
├── go.mod              # Go 模組檔案
├── go.sum              # Go 模組校驗檔案
├── main.go             # 應用程式入口
├── router.go           # 路由定義
└── README.md           # 專案說明文件
```

# 更多等待實作功能

參考 `ISSUES.md` 文件中的待辦事項列表。
