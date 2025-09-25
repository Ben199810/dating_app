# Go 聊天室應用程式的 Copilot 指導說明

## 架構概述

這是一個採用 **領域驅動設計 (Domain-Driven Design, DDD)** 的聊天室應用程式，使用 Go、Docker、MySQL 和 WebSocket 建構。專案遵循 4 層架構：

- **領域層** (`domain/`): 業務邏輯、實體和儲存庫介面
- **基礎設施層** (`infrastructure/`): 資料庫實作、外部適配器
- **展示層** (`server/`): HTTP 處理器、WebSocket、中介軟體
- **元件層** (`component/`): 共用工具如驗證器

## 關鍵開發模式

### 服務初始化模式

服務在 `main.go` 中使用建構子注入初始化：

```go
userRepo, authRepo := mysql.NewUserRepository(db)
userService := service.NewUserService(userRepo, authRepo)
handler.SetUserService(userService) // 全域處理器依賴注入
```

### 儲存庫模式

- 介面定義位於 `domain/repository/`
- MySQL 實作位於 `infrastructure/mysql/`
- 永遠回傳領域實體，不是資料庫模型

### 配置管理

各環境專用的 YAML 配置檔案位於 `config/`:

- `development.yaml` (預設)、`production.yaml`、`test.yaml`
- 設定 `APP_ENV` 環境變數來切換環境
- 透過 `config.LoadConfig("")` 存取

## 關鍵開發工作流程

### Docker 開發環境設定

```bash
# 完整環境包含資料庫
docker compose -f build/docker-compose.yaml up -d

# 可用服務:
# - 應用程式: localhost:8080
# - MySQL: localhost:3306  
# - phpMyAdmin: localhost:8081
```

### 本機開發

```bash
# 預設使用 development.yaml
go run .

# 覆寫環境
APP_ENV=production go run .

# 熱重載（需要安裝 air）
air
```

## 專案特定慣例

### 實體設計

- 列舉型別作為自定義類型：`type Gender string` 搭配常數
- 實體上的領域方法：`StringArray.Contains()`、`User.IsActive()`
- 純粹的領域邏輯，不依賴基礎設施

### 處理器模式

- 全域服務注入：`handler.SetUserService(userService)`
- RESTful API 結構：`/api/user/register`、`/api/auth/login`
- 為 HTML 前端提供靜態檔案服務

### 資料庫整合

- 使用 GORM 搭配 MySQL 驅動程式
- 實體映射位於 `infrastructure/mysql/entity_mapper.go`
- docker-compose 中的資料庫健康檢查

## WebSocket 實作

- 背景訊息處理：`go handler.HandleMessages()`
- 前端透過 JavaScript 連接到 WebSocket 端點
- 即時聊天功能與使用者認證整合

## 測試策略

- 透過 `test.yaml` 進行環境隔離
- 儲存庫介面支援輕鬆的模擬測試
- 執行測試：`go test ./...` 或 `go test -cover ./...`

## 檔案組織規則

- **靜態前端**：全部放在 `static/` (html/、css/、js/、assets/)
- **API 處理器**：在 `server/handler/` 中按功能分組
- **業務邏輯**：絕不放在處理器中，永遠放在 `domain/service/`
- **資料庫程式碼**：只能放在 `infrastructure/` 實作中

## 常見錯誤避免

- 不要將業務邏輯放在處理器中 - 使用領域服務
- 不要從處理器直接存取資料庫 - 使用儲存庫
- 不要忘記為不同環境設定 `APP_ENV`
- 永遠使用依賴注入，不使用全域變數（除了處理器設定）
