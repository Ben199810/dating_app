# Tasks: 18+ 交友聊天網頁應用程式

**Input**: 設計文件來自 `/specs/001-18/`

**Prerequisites**: plan.md (必需), research.md, data-model.md, contracts/api-spec.yaml, quickstart.md

## 執行流程 (主要)

```text
1. 從功能目錄載入 plan.md
   → 如果找不到：錯誤 "找不到實作計劃"
   → 提取：技術棧、函式庫、結構
2. 載入可選設計文件：
   → data-model.md: 提取實體 → 模型任務
   → contracts/: 每個檔案 → 合約測試任務
   → research.md: 提取決策 → 設置任務
3. 按類別生成任務：
   → 設置：專案初始化、依賴項、代碼檢查
   → 測試：合約測試、整合測試
   → 核心：模型、服務、端點
   → 整合：資料庫、中間件、日誌
   → 優化：單元測試、效能、文件
4. 應用任務規則：
   → 不同檔案 = 標記 [P] 平行執行
   → 相同檔案 = 順序執行 (無 [P])
   → 測試優先於實作 (TDD)
5. 順序編號任務 (T001, T002...)
6. 生成依賴關係圖
7. 創建平行執行範例
8. 驗證任務完整性：
   → 所有合約都有測試？
   → 所有實體都有模型？
   → 所有端點都已實作？
9. 返回：成功 (任務準備執行)
```

## 格式：`[ID] [P?] 描述`

- **[P]**：可以平行執行（不同檔案，無依賴）
- 在描述中包含確切的檔案路徑

## 路徑約定

基於 plan.md 分析，這是一個 Go Web 應用程式，使用現有的 DDD 架構：

- **領域層**：`domain/entity/`, `domain/repository/`, `domain/usecase/`
- **基礎設施層**：`infrastructure/mysql/`, `infrastructure/redis/`
- **展示層**：`server/handler/`, `server/websocket/`, `server/middleware/`
- **前端**：`static/html/`, `static/css/`, `static/js/`
- **測試**：`tests/contract/`, `tests/integration/`

## Phase 3.1: 設置

- [x] T001 更新專案依賴 (WebSocket、Redis client、JWT)
- [x] T002 [P] 配置資料庫遷移腳本 (build/init.sql)
- [x] T003 [P] 更新 Docker compose 配置支援 Redis
- [x] T004 [P] 配置年齡驗證環境變數

## Phase 3.2: 測試優先 (TDD) ⚠️ 必須在 3.3 之前完成

> **重要**：這些測試必須撰寫並且必須失敗才能進行任何實作

### 合約測試

- [x] T005 [P] 合約測試 POST /api/auth/register 在 tests/contract/auth_register_test.go
- [x] T006 [P] 合約測試 POST /api/auth/login 在 tests/contract/auth_login_test.go
- [x] T007 [P] 合約測試 GET /users/profile 在 tests/contract/profile_get_test.go
- [x] T008 [P] 合約測試 PUT /users/profile 在 tests/contract/profile_update_test.go
- [x] T009 [P] 合約測試 POST /users/photos 在 tests/contract/photos_upload_test.go
- [x] T010 [P] 合約測試 GET /matching/potential 在 tests/contract/matching_get_test.go
- [x] T011 [P] 合約測試 POST /matching/swipe 在 tests/contract/matching_swipe_test.go
- [x] T012 [P] 合約測試 GET /chat/matches 在 tests/contract/chat_matches_test.go
- [x] T013 [P] 合約測試 POST /chat/messages 在 tests/contract/chat_send_test.go

### 整合測試

- [x] T014 [P] 整合測試用戶註冊流程（18+驗證）在 tests/integration/user_registration_test.go
- [x] T015 [P] 整合測試配對演算法在 tests/integration/matching_algorithm_test.go
- [x] T016 [P] 整合測試 WebSocket 聊天在 tests/integration/websocket_chat_test.go
- [x] T017 [P] 整合測試檢舉系統在 tests/integration/report_system_test.go

## Phase 3.3: 核心實作（只有在測試失敗後）

### 領域實體模型

- [x] T018 [P] User 實體在 domain/entity/user.go
- [x] T019 [P] UserProfile 實體在 domain/entity/user_profile.go
- [x] T020 [P] Match 實體在 domain/entity/match.go
- [x] T021 [P] ChatMessage 實體在 domain/entity/chat_message.go
- [x] T022 [P] Report 實體在 domain/entity/report.go
- [x] T023 [P] Block 實體在 domain/entity/block.go
- [x] T024 [P] Interest 實體在 domain/entity/interest.go
- [x] T025 [P] Photo 實體在 domain/entity/photo.go
- [x] T026 [P] AgeVerification 實體在 domain/entity/age_verification.go

### 儲存庫介面

- [x] T027 [P] UserRepository 介面在 domain/repository/user_repository.go
- [x] T028 [P] MatchRepository 介面在 domain/repository/match_repository.go
- [x] T029 [P] ChatRepository 介面在 domain/repository/chat_repository.go
- [x] T030 [P] ReportRepository 介面在 domain/repository/report_repository.go

### MySQL 實作

- [ ] T031 [P] MySQL User 實作在 infrastructure/mysql/user_repository.go
- [ ] T032 [P] MySQL Match 實作在 infrastructure/mysql/match_repository.go
- [ ] T033 [P] MySQL Chat 實作在 infrastructure/mysql/chat_repository.go
- [ ] T034 [P] MySQL Report 實作在 infrastructure/mysql/report_repository.go
- [ ] T035 [P] 實體映射器更新在 infrastructure/mysql/entity_mapper.go

### 領域服務

- [ ] T036 [P] UserService（註冊、認證）在 domain/usecase/user_service.go
- [ ] T037 [P] MatchingService（配對演算法）在 domain/usecase/matching_service.go
- [ ] T038 [P] ChatService（訊息處理）在 domain/usecase/chat_service.go
- [ ] T039 [P] ReportService（檢舉處理）在 domain/usecase/report_service.go

### API 端點實作

- [ ] T040 POST /api/auth/register 處理器在 server/handler/auth_handler.go
- [ ] T041 POST /api/auth/login 處理器在 server/handler/auth_handler.go
- [ ] T042 GET /users/profile 處理器在 server/handler/user_handler.go
- [ ] T043 PUT /users/profile 處理器在 server/handler/user_handler.go
- [ ] T044 POST /users/photos 處理器在 server/handler/user_handler.go
- [ ] T045 GET /matching/potential 處理器在 server/handler/matching_handler.go
- [ ] T046 POST /matching/swipe 處理器在 server/handler/matching_handler.go
- [ ] T047 GET /chat/matches 處理器在 server/handler/chat_handler.go
- [ ] T048 POST /chat/messages 處理器在 server/handler/chat_handler.go

### 前端頁面

- [ ] T049 [P] 註冊頁面更新（年齡驗證）在 static/html/register.html
- [ ] T050 [P] 個人檔案頁面在 static/html/profile.html
- [ ] T051 [P] 配對頁面在 static/html/matching.html
- [ ] T052 [P] 聊天室頁面更新在 static/html/chat_websocket.html

### 前端 JavaScript

- [ ] T053 [P] 註冊表單驗證在 static/js/register.js
- [ ] T054 [P] 個人檔案管理在 static/js/profile.js
- [ ] T055 [P] 配對介面在 static/js/matching.js
- [ ] T056 聊天 WebSocket 更新在 static/js/chat.js

## Phase 3.4: 整合

### WebSocket 實作

- [ ] T057 WebSocket 管理器在 server/websocket/manager.go
- [ ] T058 [P] 聊天訊息廣播在 server/websocket/chat_handler.go
- [ ] T059 [P] WebSocket 中間件在 server/middleware/websocket_auth.go

### 認證與授權

- [ ] T060 JWT 中間件在 server/middleware/jwt_auth.go
- [ ] T061 [P] 年齡驗證中間件在 server/middleware/age_verification.go
- [ ] T062 [P] API 速率限制在 server/middleware/rate_limit.go

### Redis 快取

- [ ] T063 [P] Redis 配置在 infrastructure/redis/config.go
- [ ] T064 [P] 會話快取在 infrastructure/redis/session_cache.go
- [ ] T065 [P] 配對快取在 infrastructure/redis/matching_cache.go

### 資料庫整合

- [ ] T066 連接 UserService 到資料庫
- [ ] T067 連接 MatchingService 到資料庫和快取
- [ ] T068 連接 ChatService 到資料庫和 WebSocket
- [ ] T069 [P] 資料庫遷移和種子資料

## Phase 3.5: 優化

### 單元測試

- [ ] T070 [P] User 實體單元測試在 tests/unit/user_entity_test.go
- [ ] T071 [P] UserService 單元測試在 tests/unit/user_service_test.go
- [ ] T072 [P] MatchingService 單元測試在 tests/unit/matching_service_test.go
- [ ] T073 [P] 年齡驗證單元測試在 tests/unit/age_validation_test.go

### 效能與安全

- [ ] T074 [P] 配對查詢效能測試（<200ms）
- [ ] T075 [P] WebSocket 併發測試（1000 用戶）
- [ ] T076 [P] SQL 注入防護驗證
- [ ] T077 [P] XSS 防護測試

### 文件與部署

- [ ] T078 [P] 更新 API 文件
- [ ] T079 [P] 更新 README.md
- [ ] T080 [P] Docker 生產配置
- [ ] T081 執行手動測試（quickstart.md 場景）

## 依賴關係

### 階段依賴

- 設置任務 (T001-T004) → 測試任務 (T005-T017)
- 測試任務 (T005-T017) → 實作任務 (T018-T056)
- 實體模型 (T018-T026) → 儲存庫介面 (T027-T030)
- 儲存庫介面 (T027-T030) → MySQL 實作 (T031-T035)
- 實體和儲存庫 → 領域服務 (T036-T039)
- 領域服務 → API 處理器 (T040-T048)
- API 處理器 → 前端整合 (T049-T056)
- 核心實作 → 整合 (T057-T069)
- 整合 → 優化 (T070-T081)

### 具體依賴

- T018-T026 (實體) 阻塞 T031-T035 (MySQL 實作)
- T036 (UserService) 阻塞 T040-T044 (用戶端點)
- T037 (MatchingService) 阻塞 T045-T046 (配對端點)
- T038 (ChatService) 阻塞 T047-T048, T057-T058 (聊天功能)
- T057 (WebSocket 管理器) 阻塞 T058-T059 (WebSocket 功能)
- T060 (JWT 中間件) 阻塞大部分 API 端點
- T066-T069 (資料庫整合) 阻塞 T081 (手動測試)

## 平行執行範例

### 階段 3.2：合約測試（並行）

```bash
# 同時啟動 T005-T013：
Task: "合約測試 POST /api/auth/register 在 tests/contract/auth_register_test.go"
Task: "合約測試 POST /api/auth/login 在 tests/contract/auth_login_test.go"
Task: "合約測試 GET /users/profile 在 tests/contract/profile_get_test.go"
Task: "合約測試 PUT /users/profile 在 tests/contract/profile_update_test.go"
Task: "合約測試 POST /users/photos 在 tests/contract/photos_upload_test.go"
Task: "合約測試 GET /matching/potential 在 tests/contract/matching_get_test.go"
Task: "合約測試 POST /matching/swipe 在 tests/contract/matching_swipe_test.go"
Task: "合約測試 GET /chat/matches 在 tests/contract/chat_matches_test.go"
Task: "合約測試 POST /chat/messages 在 tests/contract/chat_send_test.go"
```

### 階段 3.3：實體模型（並行）

```bash
# 同時啟動 T018-T026：
Task: "User 實體在 domain/entity/user.go"
Task: "UserProfile 實體在 domain/entity/user_profile.go"  
Task: "Match 實體在 domain/entity/match.go"
Task: "ChatMessage 實體在 domain/entity/chat_message.go"
Task: "Report 實體在 domain/entity/report.go"
Task: "Block 實體在 domain/entity/block.go"
Task: "Interest 實體在 domain/entity/interest.go"
Task: "Photo 實體在 domain/entity/photo.go"
Task: "AgeVerification 實體在 domain/entity/age_verification.go"
```

### 階段 3.4：中間件（並行）

```bash
# 同時啟動 T059, T061, T062：
Task: "WebSocket 中間件在 server/middleware/websocket_auth.go"
Task: "年齡驗證中間件在 server/middleware/age_verification.go"
Task: "API 速率限制在 server/middleware/rate_limit.go"
```

## 注意事項

- [P] 任務 = 不同檔案，無依賴關係
- 在實作前驗證測試失敗
- 每個任務後提交代碼
- 避免：模糊任務、相同檔案衝突

## 任務生成規則

- 每個合約端點 → 一個合約測試任務 (標記 [P])
- 每個實體在 data-model → 一個模型創建任務 (標記 [P])
- 每個端點 → 一個實作任務（如果共享檔案則不平行）
- 每個用戶故事 → 一個整合測試 (標記 [P])
- 不同檔案 = 可以平行 [P]
- 相同檔案 = 順序執行（無 [P]）

## 驗證檢查清單

- ✅ 所有 API 端點都有合約測試
- ✅ 所有實體都有模型任務
- ✅ 所有服務都有單元測試
- ✅ 18+ 年齡驗證機制包含
- ✅ WebSocket 即時聊天功能
- ✅ 雙向配對機制
- ✅ 檢舉與封鎖系統
- ✅ MySQL + Redis 資料層
- ✅ JWT 認證授權
- ✅ Docker 化部署支援
