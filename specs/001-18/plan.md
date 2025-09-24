
# Implementation Plan: 18+ 交友聊天網頁應用程式

**Branch**: `001-18` | **Date**: 2025-01-24 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-18/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
18+ 交友聊天網頁應用程式：使用者必須註冊並驗證年滿18歲才能使用，提供滑動配對（雙向同意機制）和一對一即時聊天功能，包含檢舉與封鎖機制。技術棧採用 Go 後端、HTML/CSS/JavaScript 前端、MySQL 主要儲存、Redis 快取、WebSocket 即時通訊。

## Technical Context
**Language/Version**: Go 1.23+, HTML5/CSS3/ES6+ JavaScript  
**Primary Dependencies**: Gin (Web框架), GORM (ORM), WebSocket, MySQL Driver, Redis Client  
**Storage**: MySQL 8.0+ (主要資料庫), Redis 7.0+ (快取與會話)  
**Testing**: Go標準測試工具 `go test`, 前端單元測試  
**Target Platform**: Web瀏覽器 (現代瀏覽器支援)  
**Project Type**: web (frontend + backend)  
**Performance Goals**: 支援1000並發用戶, WebSocket即時通訊<100ms延遲  
**Constraints**: 18+年齡驗證強制要求, GDPR個資保護合規, 即時通訊可靠性  
**Scale/Scope**: 初期10k用戶規模, 擴展至100k用戶, 24/7高可用性

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**基於憲章 v1.1.0 的檢查項目：**

### I. 領域驅動設計（DDD）架構

✅ **PASS**: 設計將採用四層架構（領域層/基礎設施層/展示層/元件層）

- 使用者、配對、聊天等業務邏輯封裝在 domain/ 層
- 儲存庫介面在 domain/repository/，實作在 infrastructure/mysql/
- 業務邏輯不放在 HTTP 處理器中

### II. 事件驅動開發

✅ **PASS**: WebSocket 事件驅動架構適合即時聊天需求

- WebSocket 連接管理透過事件處理器
- 配對成功、新訊息等狀態變更透過事件傳播
- 使用背景 goroutine 處理異步事件

### III. 四層架構分離（不可協商）

✅ **PASS**: 嚴格遵循架構分離

- domain/: 使用者實體、配對邏輯、聊天業務規則
- infrastructure/: MySQL/Redis 實作
- server/: HTTP API、WebSocket 處理器
- component/: 共用驗證器、工具函式

### IV. 依賴注入模式

✅ **PASS**: 服務初始化使用建構子注入

- 在 main.go 中建構服務依賴
- 處理器透過全域設定接收服務
- 測試使用介面模擬

### V. 繁體中文優先

✅ **PASS**: 文件與回應訊息使用繁體中文

- API 錯誤訊息使用繁體中文
- 使用者介面文字繁體中文
- 程式碼變數名稱使用英文

### VI. 程式碼風格一致性

✅ **PASS**: 遵循 Go 標準命名規範

- 結構體大駝峰式：User, Match, ChatMessage
- 介面 er 後綴：UserRepository, ChatService
- 錯誤 Err 前綴：ErrUserNotFound, ErrInvalidAge

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure]
```

**Structure Decision**: [DEFAULT to Option 1 unless Technical Context indicates web/mobile app]

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

✅ **完成**: research.md 已生成，所有技術決策已確定

## Phase 1: Design & Contracts

*Prerequisites: research.md complete ✅*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh copilot`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

✅ **完成**: 所有 Phase 1 文件已生成

- data-model.md: 9個核心實體設計完成
- contracts/api-spec.yaml: OpenAPI 3.0 規格完成
- quickstart.md: 開發者快速上手指南完成

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:

基於 Phase 1 設計文件生成具體實作任務：

1. **測試任務** (優先執行)
   - 每個 API 端點 → 合約測試任務 [P]
   - 每個實體 → 單元測試任務 [P] 
   - 每個使用者故事 → 整合測試任務

2. **實作任務** (依賴順序)
   - 資料模型建立：User, Match, ChatMessage 等 9 個實體
   - 儲存庫實作：UserRepository, MatchRepository 等
   - 服務層實作：AuthService, MatchService, ChatService
   - API 處理器：認證、配對、聊天端點
   - WebSocket 處理器：即時通訊功能

3. **整合任務**
   - 資料庫遷移腳本
   - Docker 環境設定
   - 前端頁面更新

**排序策略**:

- **TDD 順序**: 測試先於實作，確保品質
- **依賴順序**: 模型 → 儲存庫 → 服務 → 處理器
- **標記 [P]**: 可平行執行的獨立任務

**預估輸出**: 30-35 個已排序的具體任務

**重要提醒**: 此階段由 /tasks 指令執行，NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:

- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:

- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [ ] Complexity deviations documented

## Post-Design Constitution Check

**重新驗證憲章 v1.1.0 的所有要求**：

### I. 領域驅動設計（DDD）架構

✅ **PASS**: data-model.md 確認實體設計遵循 DDD 原則

- 9個核心實體：User, UserProfile, Match, ChatMessage, Report, Block, Interest, Photo, AgeVerification
- 實體封裝業務邏輯：Match 狀態轉換, User 年齡驗證
- 清晰的聚合邊界：User-UserProfile, Match-ChatMessage

### II. 事件驅動開發

✅ **PASS**: API 規格和 WebSocket 設計支援事件驅動

- WebSocket 事件：chat_message, user_online, match_notification
- 異步處理：配對成功觸發通知事件
- 背景處理：訊息分發、在線狀態管理

### III. 四層架構分離（不可協商）

✅ **PASS**: 設計文件嚴格遵循四層分離

- Domain 層：實體模型、業務規則、儲存庫介面
- Infrastructure 層：MySQL/Redis 實作、外部服務
- Server 層：REST API、WebSocket 處理器
- Component 層：驗證器、工具函式

### IV. 依賴注入模式

✅ **PASS**: API 設計支援依賴注入模式

- 服務層依賴儲存庫介面，不依賴具體實作
- 處理器依賴服務層，透過建構子注入
- 測試可使用 Mock 介面

### V. 繁體中文優先

✅ **PASS**: 所有設計文件使用繁體中文

- API 錯誤訊息繁體中文：「必須年滿18歲才能註冊」
- 文件說明繁體中文：實體描述、驗證規則
- 程式碼維持英文：User, Match, ChatMessage

### VI. 程式碼風格一致性

✅ **PASS**: 資料模型遵循 Go 命名規範

- 實體名稱：User, UserProfile, Match, ChatMessage
- 列舉類型：MatchStatus, MessageType, Gender
- 常數命名：StatusPending, MessageText, GenderMale

---
*Based on Constitution v1.1.0 - See `/memory/constitution.md`*
