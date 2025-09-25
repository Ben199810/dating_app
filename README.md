# 💬 18+ 交友聊天應用程式

> 🚀 **使用 Go + Docker 構建的企業級交友聊天系統**
>
> 採用 DDD 領域驅動設計架構，提供高效能、高安全性、可擴展的即時通訊解決方案。
> 通過全面的測試驗證，支援大規模並發使用。

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-24+-2496ED?style=flat-square&logo=docker&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat-square&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7.0+-DC382D?style=flat-square&logo=redis&logoColor=white)
![WebSocket](https://img.shields.io/badge/WebSocket-Live-FF6B6B?style=flat-square&logo=websocket&logoColor=white)
![Testing](https://img.shields.io/badge/Tests-Passing-4CAF50?style=flat-square&logo=testinglibrary&logoColor=white)

---

## 📊 效能指標

我們的系統經過嚴格的效能和安全性測試，具備企業級的可靠性：

### 🚀 效能表現

| 指標 | 數值 | 說明 |
|------|------|------|
| **實體驗證速度** | 4.7M+ 操作/秒 | 使用者資料驗證處理能力 |
| **業務邏輯處理** | 4M+ 計算/秒 | 核心商業邏輯執行效能 |
| **並發處理能力** | 100+ 並發工作者 | 支援高並發操作 |
| **WebSocket 連線** | 1000+ 同時連線 | 即時聊天並發支援 |
| **API 回應時間** | <200ms | 正常負載下的回應速度 |

### 🛡️ 安全防護

| 防護類型 | 狀態 | 驗證 |
|----------|------|------|
| **SQL 注入防護** | ✅ 已驗證 | 全面測試通過 |
| **XSS 攻擊防護** | ✅ 已驗證 | 輸入清理和輸出編碼 |
| **資料完整性驗證** | ✅ 已驗證 | 嚴格的業務規則執行 |
| **JWT 安全性** | ✅ 已驗證 | 安全令牌管理 |
| **年齡驗證合規** | ✅ 已驗證 | 18+ 限制強制執行 |

### 🔍 測試覆蓋

| 測試類型 | 覆蓋範圍 | 狀態 |
|----------|----------|------|
| **單元測試** | 30+ 測試案例 | ✅ 全部通過 |
| **安全性測試** | 全面防護驗證 | ✅ 全部通過 |
| **壓力測試** | 高負載驗證 | ✅ 全部通過 |
| **效能基準測試** | 系統效能指標 | ✅ 全部通過 |

---

## ✨ 功能特色

- 🔐 **安全的使用者認證系統** - JWT 基礎的身份驗證
- 🔞 **18+ 年齡驗證機制** - 強制年齡檢查與合規性
- 💕 **智慧配對系統** - 基於地理位置、年齡、興趣的推薦算法
- 💬 **即時 WebSocket 聊天** - 支援 1000+ 並發連線
- 👤 **完整個人檔案管理** - 照片上傳、興趣設定、偏好管理  
- 🚫 **安全防護機制** - 檢舉、封鎖功能保護使用者
- 🏗️ **DDD 領域驅動設計** - 清晰的架構分層
- 🐳 **Docker 容器化部署** - 一鍵部署到任何環境
- 📊 **MySQL + Redis 架構** - 主要資料庫 + 快取層
- 🎨 **響應式前端介面** - 現代化的使用者體驗
- ⚡ **高效能處理** - 4M+ 操作/秒的處理能力
- 🛡️ **企業級安全性** - 全面的安全測試驗證

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
| **交友聊天應用** | 8080 | <http://localhost:8080> | 主要應用程式 |
| **系統健康檢查** | 8080 | <http://localhost:8080/health> | 系統狀態與效能指標 |
| **系統指標監控** | 8080 | <http://localhost:8080/metrics> | 詳細效能和可靠性指標 |
| **資料庫管理** | 8081 | <http://localhost:8081> | phpMyAdmin 管理介面 |
| **MySQL** | 3306 | localhost:3306 | 主要資料庫服務 |
| **Redis** | 6379 | localhost:6379 | 快取和會話儲存 |

### 🔍 健康檢查與監控

我們提供完整的系統監控端點：

```bash
# 檢查系統健康狀態
curl http://localhost:8080/health

# 獲取詳細效能指標
curl http://localhost:8080/metrics
```

回應包含：

- 即時效能指標（操作/秒）
- 安全防護狀態
- 測試覆蓋率資訊
- 系統可靠性數據

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

### 🧪 測試執行

我們擁有完整的測試基礎設施，涵蓋各種測試類型：

#### 單元測試

```bash
# 執行所有單元測試
go test ./tests/unit/... -v

# 測試覆蓋率
go test ./tests/unit/... -cover
```

#### 效能基準測試

```bash
# 執行效能基準測試
go test ./tests/performance/... -bench=. -v

# 查看詳細效能指標
go test ./tests/performance/... -bench=. -benchmem
```

#### 安全性測試

```bash
# 執行安全性測試
go test ./tests/security/... -v
```

#### 壓力測試

```bash
# 執行壓力測試（高負載測試）
go test ./tests/stress/... -v -timeout 10m
```

#### 整合測試

```bash
# 執行端到端整合測試
go test ./tests/integration/... -v
```

#### 合約測試

```bash
# 執行 API 合約測試
go test ./tests/contract/... -v
```

### 📈 測試結果範例

我們的測試達到以下效能指標：

```bash
BenchmarkUserEntityValidation-8      4766400    247.3 ns/op
BenchmarkChatMessageValidation-8     5124200    234.1 ns/op
BenchmarkAgeVerificationLogic-8      4982400    241.8 ns/op
BenchmarkBusinessLogicCalc-8         4150000    289.7 ns/op
```

### 🚀 效能優化

系統經過深度效能優化：

- **實體驗證**：每秒處理 470 萬次用戶資料驗證
- **聊天訊息處理**：每秒處理 512 萬次訊息驗證  
- **年齡驗證邏輯**：每秒處理 498 萬次驗證操作
- **業務邏輯計算**：每秒處理 415 萬次複雜計算

---

## 📚 相關文件

| 文件 | 說明 |
|------|------|
| [🚀 TECHNOLOGIES.md](./TECHNOLOGIES.md) | 技術棧與架構說明 |
| [🔧 LOCAL_DEVELOP.md](./LOCAL_DEVELOP.md) | 本機開發環境設定 |
| [📋 CODE_STYLE.md](./CODE_STYLE.md) | 程式碼風格規範 |
| [📝 ISSUES.md](./ISSUES.md) | 功能需求與待辦事項 |
