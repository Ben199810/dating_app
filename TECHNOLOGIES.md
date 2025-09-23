# 🚀 技術棧與架構

> 採用現代化的 **DDD 領域驅動設計 (Domain-Driven Design)** 架構，構建可擴展且易維護的即時聊天室應用程式。

---

## 🎨 前端技術棧

### 📋 核心技術

| 技術 | 版本 | 用途 | 文件 |
|------|------|------|------|
| **HTML5** | Latest | 頁面結構與語意化標記 | [MDN Docs](https://developer.mozilla.org/docs/Web/HTML) |
| **CSS3** | Latest | 樣式設計與響應式佈局 | [MDN Docs](https://developer.mozilla.org/docs/Web/CSS) |
| **Vanilla JavaScript** | ES6+ | 動態互動與 WebSocket 通訊 | [MDN Docs](https://developer.mozilla.org/docs/Web/JavaScript) |

### 📁 前端資料夾結構

```text
static/
├── 📂 assets/          # 靜態資源
│   ├── img/           # 圖片素材
│   ├── fonts/         # 字型檔案
│   └── icons/         # 圖示資源
├── 🎨 css/            # 樣式表
│   ├── common.css     # 全域通用樣式
│   ├── index.css      # 首頁專用樣式
│   ├── chat.css       # 聊天室樣式
│   ├── login.css      # 登入頁面樣式
│   ├── register.css   # 註冊頁面樣式
│   └── profile.css    # 個人檔案樣式
├── 📄 html/           # HTML 頁面
│   ├── index.html     # 首頁
│   ├── chat_websocket.html  # 聊天室
│   ├── login.html     # 登入頁面
│   ├── register.html  # 註冊頁面
│   └── profile.html   # 個人檔案
└── ⚡ js/             # JavaScript 腳本
    ├── index.js       # 首頁邏輯
    └── chat.js        # 聊天室邏輯
```

### 🎯 前端開發指南

| 檔案類型 | 建議位置 | 命名規範 | 說明 |
|----------|----------|----------|------|
| 🏗️ **HTML** | `static/html/` | `kebab-case.html` | 語意化標記，響應式設計 |
| 🎨 **CSS** | `static/css/` | `kebab-case.css` | 模組化樣式，BEM 方法論 |
| ⚡ **JavaScript** | `static/js/` | `camelCase.js` | ES6+ 語法，模組化開發 |
| 🖼️ **Assets** | `static/assets/` | `descriptive-name` | 優化圖片，適當格式 |

---

## ⚙️ 後端技術棧

### 🛠️ 核心技術

| 技術 | 版本 | 角色 | 特色 |
|------|------|------|------|
| **Go** | 1.23+ | 主要開發語言 | 高效能、強型別、並發處理 |
| **Gin** | v1.10+ | Web 框架 | 快速、輕量級 HTTP 路由 |
| **GORM** | v1.25+ | ORM 工具 | 優雅的資料庫抽象層 |
| **WebSocket** | - | 即時通訊 | 雙向即時資料傳輸 |
| **MySQL** | 8.0+ | 主資料庫 | 關聯式資料儲存 |
| **Redis** | 7.0+ | 快取系統 | 高效能記憶體資料庫 |
| **Docker** | 24+ | 容器化 | 環境一致性與部署便利 |

### 🏗️ DDD 架構層次

```text
📦 golang_dev_docker/
├── 🎯 domain/                    # 🟦 領域層 (Domain Layer)
│   ├── entity/                  #    業務實體定義
│   │   └── user.go             #    使用者實體
│   ├── repository/             #    儲存庫介面
│   │   ├── auth_repository.go   #    認證儲存庫
│   │   ├── user_repository.go   #    使用者儲存庫
│   │   └── user_profile_repository.go  # 使用者檔案儲存庫
│   └── service/                #    業務邏輯服務
│       ├── auth_user_login.go   #    登入業務邏輯
│       ├── create_new_user.go   #    註冊業務邏輯
│       └── user_profile_service.go     # 檔案管理邏輯
├── 🔧 infrastructure/           # 🟩 基礎設施層 (Infrastructure)
│   └── mysql/                  #    MySQL 實作
│       ├── entity_mapper.go     #    實體映射器
│       ├── user_repository.go   #    使用者儲存庫實作
│       └── user_profile_repository.go  # 檔案儲存庫實作
├── 🌐 server/                  # 🟨 展示層 (Presentation Layer)
│   ├── handler/                #    HTTP 處理器
│   │   ├── health_check_handler.go    # 健康檢查
│   │   ├── user_handler.go            # 使用者相關 API
│   │   ├── user_profile_handler.go    # 檔案相關 API
│   │   ├── test_handler.go            # 測試端點
│   │   └── websocket_handler.go       # WebSocket 處理
│   ├── middleware/             #    中介軟體
│   └── swagger/               #    API 文件
├── 🧩 component/              # 🟪 共用元件
│   └── validator/             #    資料驗證器
│       └── validator.go       #    輸入驗證邏輯
├── ⚙️ config/                # 配置管理
│   ├── config.go             #    配置讀取邏輯
│   ├── development.yaml      #    開發環境配置
│   ├── production.yaml       #    正式環境配置
│   └── test.yaml            #    測試環境配置
├── 🐳 build/                 # Docker 相關
│   ├── docker-compose.yaml   #    服務編排
│   ├── Dockerfile           #    應用程式映像
│   └── init.sql            #    資料庫初始化
├── 📄 static/               # 前端靜態檔案
├── 🛣️ routes.go            #    路由配置
└── 🚀 main.go              #    應用程式入口
```

### 🎯 架構設計原則

| 層次 | 職責 | 依賴方向 | 核心概念 |
|------|------|----------|----------|
| **🎯 Domain** | 業務邏輯與規則 | 不依賴其他層 | 實體、值物件、聚合根 |
| **🔧 Infrastructure** | 技術實作細節 | 依賴 Domain | 資料庫、外部服務 |
| **🌐 Presentation** | 使用者介面 | 依賴 Domain | HTTP、WebSocket、JSON |
| **🧩 Component** | 共用功能 | 被其他層使用 | 驗證、工具、輔助 |

---

## 🗄️ 資料庫設計

### � 核心資料表

| 資料表 | 用途 | 主要欄位 | 索引 |
|--------|------|----------|------|
| **users** | 使用者基本資訊 | id, username, email, password | username, email, status |
| **user_profiles** | 使用者詳細檔案 | user_id, bio, interests, location | user_id, location |

### 🔐 安全性設計

- ✅ **密碼加密**：bcrypt 雜湊演算法
- ✅ **JWT 認證**：無狀態身份驗證
- ✅ **輸入驗證**：防止 SQL 注入與 XSS
- ✅ **CORS 設定**：跨域請求安全控制

---

## 🔗 學習資源與文件

### 📚 官方文件

| 技術 | 文件連結 | 重點主題 |
|------|----------|----------|
| **Go** | [golang.org](https://golang.org/doc/) | 語法、並發、模組 |
| **Gin** | [gin-gonic.com](https://gin-gonic.com/docs/) | 路由、中介軟體、綁定 |
| **GORM** | [gorm.io](https://gorm.io/docs/) | 模型、關聯、遷移 |
| **MySQL** | [dev.mysql.com](https://dev.mysql.com/doc/) | 查詢優化、索引設計 |
| **Redis** | [redis.io](https://redis.io/documentation) | 資料結構、持久化 |
| **WebSocket** | [MDN WebSocket](https://developer.mozilla.org/docs/Web/API/WebSockets_API) | 連接管理、訊息處理 |

### 🎓 進階學習

- 📖 **DDD 領域驅動設計**：[Domain-Driven Design 參考](https://domainlanguage.com/ddd/reference/)
- 🏗️ **Clean Architecture**：Robert C. Martin 的架構設計
- 🔄 **微服務架構**：分散式系統設計模式
- 🧪 **測試驅動開發**：單元測試與整合測試

---

## 🚀 開發環境與工具

### 🛠️ 開發工具鏈

| 工具 | 版本 | 用途 |
|------|------|------|
| **Docker** | 24+ | 容器化開發 |
| **Go** | 1.23+ | 程式開發 |
| **MySQL** | 8.0+ | 資料庫 |
| **Redis** | 7.0+ | 快取 |

### 🎯 品質保證

- ✅ **程式碼格式化**：`gofmt`, `goimports`
- ✅ **靜態分析**：`golint`, `go vet`
- ✅ **單元測試**：`go test`
- ✅ **API 文件**：Swagger/OpenAPI
- ✅ **容器化**：Docker & Docker Compose
