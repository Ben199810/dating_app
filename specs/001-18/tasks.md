# Tasks: 18+ 交友聊天網頁應用程式

**Input**: Design documents from `/specs/001-18/`  
**Prerequisites**: plan.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓, quickstart.md ✓

## 技術規格摘要

- **語言**: Go 1.23+ (後端), HTML5/CSS3/ES6+ JavaScript (前端)
- **框架**: Gin (Web框架), GORM (ORM), WebSocket
- **資料庫**: MySQL 8.0+ (主要), Redis 7.0+ (快取)
- **架構**: DDD四層架構 (domain/, infrastructure/, server/, component/)
- **核心實體**: 9個 (User, UserProfile, Match, ChatMessage, Report, Block, Interest, Photo, AgeVerification)
- **API端點**: 認證、用戶管理、配對、聊天、安全功能
- **特殊需求**: 18+年齡驗證, WebSocket即時通訊, 1000並發支援

## Phase 3.1: 專案設定與基礎架構

- [x] T001 更新 go.mod 依賴項目 (WebSocket, Redis client, GORM 外鍵支援)
- [x] T002 [P] 建立 Redis 連接配置與初始化 (infrastructure/redis/)
- [x] T003 [P] 設定 WebSocket 連接管理器結構 (server/websocket/)
- [x] T004 [P] 配置 CORS 和安全標頭中介軟體 (server/middleware/)
- [x] T005 更新 Docker Compose 加入 Redis 服務 (build/docker-compose.yaml)

## Phase 3.2: 測試優先開發 (TDD) ⚠️ 必須在實作前完成

### ⚠️ CRITICAL: 這些測試必須先寫好並且失敗，才能開始任何實作

### 合約測試 (API 端點)

- [ ] T006 [P] POST /api/auth/register 合約測試 (tests/contract/auth_register_test.go)
- [ ] T007 [P] POST /api/auth/login 合約測試 (tests/contract/auth_login_test.go)
- [ ] T008 [P] GET /api/users/profile 合約測試 (tests/contract/user_profile_get_test.go)
- [ ] T009 [P] PUT /api/users/profile 合約測試 (tests/contract/user_profile_put_test.go)
- [ ] T010 [P] GET /api/matches/discover 合約測試 (tests/contract/match_discover_test.go)
- [ ] T011 [P] POST /api/matches/like 合約測試 (tests/contract/match_like_test.go)
- [ ] T012 [P] GET /api/chats 合約測試 (tests/contract/chat_list_test.go)
- [ ] T013 [P] POST /api/chats/{match_id}/messages 合約測試 (tests/contract/chat_message_test.go)
- [ ] T014 [P] POST /api/reports 合約測試 (tests/contract/report_test.go)
- [ ] T015 [P] POST /api/blocks 合約測試 (tests/contract/block_test.go)

### 整合測試 (使用者故事)

- [ ] T016 [P] 18+年齡驗證整合測試 (tests/integration/age_verification_test.go)
- [ ] T017 [P] 用戶註冊到登入流程測試 (tests/integration/auth_flow_test.go)
- [ ] T018 [P] 雙向配對機制整合測試 (tests/integration/matching_flow_test.go)
- [ ] T019 [P] 即時聊天 WebSocket 測試 (tests/integration/websocket_chat_test.go)
- [ ] T020 [P] 檢舉封鎖功能測試 (tests/integration/safety_features_test.go)

## Phase 3.3: 核心實體模型 (測試失敗後才實作)

### 資料模型建立

- [ ] T021 [P] User 實體模型 (domain/entity/user.go)
- [ ] T022 [P] UserProfile 實體模型 (domain/entity/user_profile.go)
- [ ] T023 [P] Match 實體模型 (domain/entity/match.go)
- [ ] T024 [P] ChatMessage 實體模型 (domain/entity/chat_message.go)
- [ ] T025 [P] Report 實體模型 (domain/entity/report.go)
- [ ] T026 [P] Block 實體模型 (domain/entity/block.go)
- [ ] T027 [P] Interest 實體模型 (domain/entity/interest.go)
- [ ] T028 [P] Photo 實體模型 (domain/entity/photo.go)
- [ ] T029 [P] AgeVerification 實體模型 (domain/entity/age_verification.go)

### Repository 介面定義

- [ ] T030 [P] UserRepository 介面 (domain/repository/user_repository.go)
- [ ] T031 [P] MatchRepository 介面 (domain/repository/match_repository.go)
- [ ] T032 [P] ChatRepository 介面 (domain/repository/chat_repository.go)
- [ ] T033 [P] ReportRepository 介面 (domain/repository/report_repository.go)
- [ ] T034 [P] InterestRepository 介面 (domain/repository/interest_repository.go)

## Phase 3.4: 服務層實作

### 業務邏輯服務

- [ ] T035 [P] AuthService - 註冊登入邏輯 (domain/service/auth_service.go)
- [ ] T036 [P] UserProfileService - 檔案管理 (domain/service/user_profile_service.go)
- [ ] T037 [P] MatchService - 配對邏輯 (domain/service/match_service.go)
- [ ] T038 [P] ChatService - 聊天邏輯 (domain/service/chat_service.go)
- [ ] T039 [P] SafetyService - 檢舉封鎖邏輯 (domain/service/safety_service.go)
- [ ] T040 年齡驗證邏輯與工具函式 (component/validator/age_validator.go)

## Phase 3.5: 基礎設施層實作

### MySQL Repository 實作

- [ ] T041 [P] UserRepository MySQL 實作 (infrastructure/mysql/user_repository.go)
- [ ] T042 [P] MatchRepository MySQL 實作 (infrastructure/mysql/match_repository.go)
- [ ] T043 [P] ChatRepository MySQL 實作 (infrastructure/mysql/chat_repository.go)
- [ ] T044 [P] ReportRepository MySQL 實作 (infrastructure/mysql/report_repository.go)
- [ ] T045 [P] InterestRepository MySQL 實作 (infrastructure/mysql/interest_repository.go)

### Redis 快取實作

- [ ] T046 [P] User 快取層 (infrastructure/redis/user_cache.go)
- [ ] T047 [P] Match 快取層 (infrastructure/redis/match_cache.go)
- [ ] T048 [P] Online 狀態管理 (infrastructure/redis/online_status.go)

### 資料庫遷移

- [ ] T049 建立資料庫遷移腳本 (build/migrations/001_create_tables.sql)
- [ ] T050 興趣標籤初始資料 (build/migrations/002_seed_interests.sql)

## Phase 3.6: API 端點實作

### 認證相關端點

- [ ] T051 POST /api/auth/register 處理器 (server/handler/auth_handler.go)
- [ ] T052 POST /api/auth/login 處理器 (server/handler/auth_handler.go)
- [ ] T053 POST /api/auth/logout 處理器 (server/handler/auth_handler.go)

### 用戶管理端點

- [ ] T054 GET /api/users/profile 處理器 (server/handler/user_handler.go)
- [ ] T055 PUT /api/users/profile 處理器 (server/handler/user_handler.go)
- [ ] T056 POST /api/users/photos 處理器 (server/handler/user_handler.go)

### 配對系統端點

- [ ] T057 GET /api/matches/discover 處理器 (server/handler/match_handler.go)
- [ ] T058 POST /api/matches/like 處理器 (server/handler/match_handler.go)
- [ ] T059 GET /api/matches/pending 處理器 (server/handler/match_handler.go)
- [ ] T060 POST /api/matches/{id}/respond 處理器 (server/handler/match_handler.go)

### 聊天系統端點

- [ ] T061 GET /api/chats 處理器 (server/handler/chat_handler.go)
- [ ] T062 GET /api/chats/{id}/messages 處理器 (server/handler/chat_handler.go)
- [ ] T063 POST /api/chats/{id}/messages 處理器 (server/handler/chat_handler.go)

### 安全功能端點

- [ ] T064 POST /api/reports 處理器 (server/handler/safety_handler.go)
- [ ] T065 POST /api/blocks 處理器 (server/handler/safety_handler.go)
- [ ] T066 GET /api/blocks 處理器 (server/handler/safety_handler.go)

## Phase 3.7: WebSocket 即時通訊

- [ ] T067 WebSocket 連接管理器 (server/websocket/manager.go)
- [ ] T068 WebSocket 客戶端處理 (server/websocket/client.go)
- [ ] T069 聊天訊息即時廣播 (server/websocket/chat_handler.go)
- [ ] T070 配對成功通知 (server/websocket/notification_handler.go)
- [ ] T071 用戶在線狀態管理 (server/websocket/status_handler.go)

## Phase 3.8: 前端頁面整合

- [ ] T072 [P] 更新註冊頁面支援年齡驗證 (static/html/register.html, static/js/register.js)
- [ ] T073 [P] 更新個人檔案頁面 (static/html/profile.html, static/js/profile.js)
- [ ] T074 [P] 建立配對頁面 (static/html/matching.html, static/js/matching.js)
- [ ] T075 [P] 更新聊天頁面支援 WebSocket (static/html/chat_websocket.html, static/js/chat.js)
- [ ] T076 [P] 建立安全功能頁面 (static/html/safety.html, static/js/safety.js)

## Phase 3.9: 中介軟體與安全

- [ ] T077 JWT 認證中介軟體 (server/middleware/auth_middleware.go)
- [ ] T078 請求限流中介軟體 (server/middleware/rate_limit_middleware.go)
- [ ] T079 日誌記錄中介軟體 (server/middleware/logging_middleware.go)
- [ ] T080 錯誤處理中介軟體 (server/middleware/error_middleware.go)

## Phase 3.10: 系統整合

- [ ] T081 更新 main.go 整合所有服務與中介軟體
- [ ] T082 更新 routes.go 添加所有新端點
- [ ] T083 配置檔案更新 (config/development.yaml, config/production.yaml)
- [ ] T084 Docker 環境測試與調整

## Phase 3.11: 效能優化與完善

### 單元測試

- [ ] T085 [P] User 實體單元測試 (tests/unit/entity/user_test.go)
- [ ] T086 [P] Match 邏輯單元測試 (tests/unit/service/match_service_test.go)
- [ ] T087 [P] Age 驗證單元測試 (tests/unit/validator/age_validator_test.go)
- [ ] T088 [P] WebSocket 管理器單元測試 (tests/unit/websocket/manager_test.go)

### 效能與文件

- [ ] T089 [P] API 效能測試 (<100ms 回應時間)
- [ ] T090 [P] WebSocket 併發測試 (1000 連接)
- [ ] T091 [P] 更新 README.md 使用說明
- [ ] T092 [P] 更新 .github/copilot-instructions.md
- [ ] T093 執行 quickstart.md 完整測試流程

## 依賴關係

**Phase 執行順序**:

1. 設定 (T001-T005) → 測試 (T006-T020) → 模型 (T021-T034) → 服務 (T035-T040)
2. 基礎設施 (T041-T050) → API (T051-T066) → WebSocket (T067-T071)
3. 前端 (T072-T076) → 中介軟體 (T077-T080) → 整合 (T081-T084)
4. 完善 (T085-T093)

**關鍵依賴**:

- T006-T020 (所有測試) 必須在 T021+ (實作) 之前完成
- T021-T029 (實體) 完成後才能開始 T030-T034 (Repository介面)
- T030-T040 (介面與服務) 完成後才能開始 T041-T048 (實作)
- T049-T050 (資料庫) 完成後才能執行整合測試
- T081-T082 (整合) 完成後才能執行 T093 (完整測試)

## 平行執行範例

### 階段 1: 合約測試 (可平行)

```text
Task: "POST /api/auth/register 合約測試 in tests/contract/auth_register_test.go"
Task: "POST /api/auth/login 合約測試 in tests/contract/auth_login_test.go"  
Task: "GET /api/users/profile 合約測試 in tests/contract/user_profile_get_test.go"
Task: "PUT /api/users/profile 合約測試 in tests/contract/user_profile_put_test.go"
```

### 階段 2: 實體建立 (可平行)

```text
Task: "User 實體模型 in domain/entity/user.go"
Task: "UserProfile 實體模型 in domain/entity/user_profile.go"
Task: "Match 實體模型 in domain/entity/match.go"
Task: "ChatMessage 實體模型 in domain/entity/chat_message.go"
```

### 階段 3: Repository 實作 (可平行)

```text
Task: "UserRepository MySQL 實作 in infrastructure/mysql/user_repository.go"
Task: "MatchRepository MySQL 實作 in infrastructure/mysql/match_repository.go"
Task: "ChatRepository MySQL 實作 in infrastructure/mysql/chat_repository.go"
```

## 驗證檢查清單

**任務生成規則驗證** ✅:

- [x] 所有合約端點都有對應測試 (T006-T015)
- [x] 所有實體都有模型任務 (T021-T029)
- [x] 所有測試都在實作之前 (T006-T020 → T021+)
- [x] 平行任務確實獨立 (不同檔案，無依賴)
- [x] 每個任務都指定確切檔案路徑
- [x] 無任務與其他 [P] 任務修改同一檔案

**功能完整性驗證** ✅:

- [x] 18+ 年齡驗證功能 (T016, T040, T051)
- [x] 雙向配對機制 (T018, T023, T037, T057-T060)
- [x] WebSocket 即時聊天 (T019, T067-T071)
- [x] 檢舉封鎖安全功能 (T020, T039, T064-T066)
- [x] 9個核心實體全覆蓋 (T021-T029)
- [x] 完整 API 端點實作 (T051-T066)
- [x] 前端整合更新 (T072-T076)

## 註記

- **[P]** = 可平行執行的任務（不同檔案，無依賴關係）
- 每個任務完成後進行 git commit
- TDD 原則：確保測試失敗後再開始實作
- 遵循 DDD 四層架構：domain/ → infrastructure/ → server/ → component/
- 所有錯誤訊息使用繁體中文（憲章要求）
- Go 命名規範：大駝峰實體、介面 er 後綴、錯誤 Err 前綴

**預估完成時間**: 25-30 工作天（假設每日完成 3-4 任務）  
**技術複雜度**: 中高（WebSocket 併發處理、配對演算法）  
**風險項目**: T067-T071 (WebSocket)、T089-T090 (效能測試)
