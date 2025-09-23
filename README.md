# 💬 即時聊天室應用程式

> 🚀 **使用 Go + Docker 構建的現代化即時聊天室系統**
>
> 採用 DDD 領域驅動設計架構，提供高效能、可擴展的即時通訊解決方案。

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-24+-2496ED?style=flat-square&logo=docker&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat-square&logo=mysql&logoColor=white)
![WebSocket](https://img.shields.io/badge/WebSocket-Live-FF6B6B?style=flat-square&logo=websocket&logoColor=white)

---

## ✨ 功能特色

- 🔐 **使用者註冊與登入系統**
- 💬 **即時 WebSocket 聊天**
- 👤 **個人檔案管理**
- 🏗️ **DDD 架構設計**
- 🐳 **Docker 容器化部署**
- 📊 **MySQL 資料持久化**
- 🎨 **響應式前端介面**

---

## 🚀 快速開始

### 📋 系統需求

| 軟體 | 版本 | 說明 |
|------|------|------|
| **Docker** | 24+ | 容器化平台 |
| **Docker Compose** | 2.0+ | 服務編排工具 |
| **Go** | 1.23+ | *(可選，僅開發需要)* |

### 🏃‍♂️ 一鍵啟動

```bash
# 1. 克隆專案
git clone <repository-url>
cd golang_dev_docker

# 2. 啟動所有服務
docker compose -f build/docker-compose.yaml up -d

# 3. 訪問應用程式
open http://localhost:8080
```

### 🌐 服務端點

| 服務 | 端口 | 網址 | 說明 |
|------|------|------|------|
| **聊天室應用** | 8080 | <http://localhost:8080> | 主要應用程式 |
| **資料庫管理** | 8081 | <http://localhost:8081> | phpMyAdmin |
| **MySQL** | 3306 | localhost:3306 | 資料庫服務 |

---

## ⚙️ 配置管理

本專案使用 YAML 檔案管理不同環境的配置設定，支援開發、測試、生產環境的無縫切換。

### 📁 配置檔案

| 檔案 | 環境 | 用途 |
|------|------|------|
| `config/development.yaml` | 開發 | 本機開發設定 |
| `config/production.yaml` | 生產 | 正式環境設定 |
| `config/test.yaml` | 測試 | 單元測試設定 |

### 🔧 環境變數設定

```bash
# 設定環境變數
export APP_ENV=development  # 載入 development.yaml
export APP_ENV=production   # 載入 production.yaml
export APP_ENV=test        # 載入 test.yaml

# 如果未設定，預設載入 development.yaml
```

### 📝 配置結構

```yaml
# 資料庫設定
database:
  host: localhost
  port: 3306
  user: chat_user
  password: chat_password
  dbname: chat_app
  charset: utf8mb4
  parseTime: true
  loc: Local

# 伺服器設定
server:
  port: 8080
  mode: debug  # gin 模式: debug, release, test

# 日誌設定
logging:
  level: info
  format: json
```

---

## 🛠️ 開發模式

### 本機開發

```bash
# 使用預設環境 (development)
go run .

# 指定特定環境
APP_ENV=production go run .

# 開發模式熱重載 (需要 air)
go install github.com/cosmtrek/air@latest
air
```

### 測試執行

```bash
# 執行所有測試
go test ./...

# 執行特定模組測試
go test ./domain/service/...

# 測試覆蓋率
go test -cover ./...
```

---

## 📚 相關文件

| 文件 | 說明 |
|------|------|
| [🚀 TECHNOLOGIES.md](./TECHNOLOGIES.md) | 技術棧與架構說明 |
| [🔧 LOCAL_DEVELOP.md](./LOCAL_DEVELOP.md) | 本機開發環境設定 |
| [📋 CODE_STYLE.md](./CODE_STYLE.md) | 程式碼風格規範 |
| [📝 ISSUES.md](./ISSUES.md) | 功能需求與待辦事項 |
