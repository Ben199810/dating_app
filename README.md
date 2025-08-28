# 聊天應用程式 - DDD 事件驅動架構

這是一個使用 Go 語言開發的即時聊天應用程式，採用領域驅動設計（DDD）和事件驅動架構模式。

## 架構概覽

```
cmd/
└── server/           # 應用程式入口點
    └── main.go

internal/
├── domain/          # 領域層
│   ├── chat/        # 聊天領域
│   │   ├── message.go
│   │   └── room.go
│   └── user/        # 用戶領域
│       └── user.go
├── application/     # 應用層
│   ├── event/       # 事件定義
│   │   ├── event.go
│   │   └── chat_events.go
│   └── service/     # 應用服務
│       ├── chat_service.go
│       └── errors.go
├── infrastructure/ # 基礎設施層
│   ├── repository/  # 倉儲實現
│   │   ├── memory_repository.go
│   │   └── user_repository.go
│   ├── websocket/   # WebSocket 實現
│   │   └── manager.go
│   └── event/       # 事件總線實現
│       └── memory_event_bus.go
└── interfaces/      # 介面層
    └── http/        # HTTP 處理器
        ├── chat_handler.go
        └── user_handler.go

web/                 # 靜態檔案
├── chat.html
└── chat_websocket.html
```

## 核心概念

### 領域驅動設計 (DDD)

- **領域層 (Domain Layer)**: 包含業務實體、值對象和領域服務
- **應用層 (Application Layer)**: 協調領域對象執行業務用例
- **基礎設施層 (Infrastructure Layer)**: 提供技術實現細節
- **介面層 (Interface Layer)**: 處理外部請求和回應

### 事件驅動架構

- **事件總線**: 解耦組件間的通訊
- **事件處理器**: 響應領域事件
- **非同步處理**: 提高系統響應性和擴展性

## 功能特性

- ✅ 即時聊天訊息
- ✅ 用戶加入/離開聊天室
- ✅ WebSocket 連線管理
- ✅ 事件驅動的訊息廣播
- ✅ 記憶體儲存（開發環境）
- ✅ Docker 支援

## 快速開始

### 本地開發

1. 確保已安裝 Go 1.22+
2. 克隆專案並進入目錄
3. 安裝依賴：

   ```bash
   go mod tidy
   ```

4. 運行應用程式：

   ```bash
   go run cmd/server/main.go
   ```

5. 訪問 <http://localhost:8080/chat_websocket.html>

### Docker 部署

1. 構建並運行：

   ```bash
   docker-compose up --build
   ```

2. 訪問 <http://localhost:8080/chat_websocket.html>

## API 端點

### REST API

- `GET /api/v1/hello` - 問候端點
- `GET /api/v1/users/:id` - 獲取用戶資訊
- `POST /api/v1/messages` - 發送訊息
- `GET /api/v1/messages?room_id=xxx` - 獲取聊天室訊息
- `POST /api/v1/rooms/join` - 加入聊天室
- `POST /api/v1/rooms/leave` - 離開聊天室

### WebSocket

- `GET /ws?user_id=xxx&room_id=xxx` - WebSocket 連線

## 事件類型

- `message.sent` - 訊息發送事件
- `user.joined` - 用戶加入事件
- `user.left` - 用戶離開事件

## 擴展建議

### 生產環境改進

1. **資料庫集成**: 使用 PostgreSQL 替代記憶體儲存
2. **身份認證**: 實現 JWT 認證機制
3. **日誌系統**: 集成結構化日誌
4. **監控指標**: 添加 Prometheus 指標
5. **配置管理**: 使用環境變數和配置檔案
6. **安全性**: 實現適當的 CORS 和來源檢查

### 功能擴展

1. **私人訊息**: 實現一對一聊天
2. **檔案上傳**: 支援圖片和檔案分享
3. **訊息歷史**: 實現訊息分頁和搜尋
4. **用戶狀態**: 顯示在線/離線狀態
5. **聊天室管理**: 創建/刪除聊天室功能

## 開發注意事項

- 遵循 DDD 原則，保持領域邏輯的純淨性
- 使用事件來處理跨聚合的通訊
- 確保適當的錯誤處理和日誌記錄
- 編寫單元測試和集成測試
- 定期重構以保持代碼品質

## 貢獻指南

1. Fork 專案
2. 創建功能分支
3. 遵循現有的代碼風格
4. 編寫測試
5. 提交 Pull Request

## 授權

MIT License
