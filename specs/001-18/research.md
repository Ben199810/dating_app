# 研究分析：18+ 交友聊天應用程式

**Phase 0 研究輸出**：解決技術決策與實作細節

## 研究範圍

基於技術背景分析，以下項目需要深入研究：

### 1. WebSocket 即時通訊架構

### 2. MySQL + Redis 資料層設計

### 3. 18+ 年齡驗證機制

### 4. 雙向配對演算法

### 5. 併發用戶處理模式

### 6. 前後端整合策略

---

## WebSocket 即時通訊架構

**研究問題**: 如何在 Go 中實現可支援 1000 並發用戶的 WebSocket 架構？

### 決策：Gorilla WebSocket + 連接池管理

**選擇理由**：

- Gorilla WebSocket 是 Go 生態最成熟的 WebSocket 庫
- 原生支持 Go 的併發模型（goroutine）
- 提供完整的連接管理和訊息處理機制
- 社群支持度高，文件完善

**架構設計**：

```go
type WebSocketManager struct {
    clients    map[string]*websocket.Conn  // userId -> connection
    broadcast  chan []byte                 // 廣播頻道
    register   chan *Client               // 註冊新連接
    unregister chan *Client               // 移除連接
}

type Client struct {
    conn   *websocket.Conn
    userID string
    send   chan []byte
}
```

**併發處理模式**：

- 每個連接一個讀取 goroutine
- 每個連接一個寫入 goroutine  
- 一個中央訊息分發 goroutine
- 連接池管理 goroutine

**替代方案考慮**：

1. **nhooyr.io/websocket**: 更現代的 API，但生態較小
2. **原生 net/http WebSocket**: 功能有限，需要更多手動實作
3. **Socket.IO Go 實作**: 複雜度高，不符合輕量化需求

**效能預估**：

- 支援 1000 併發連接：✅ 可達成
- 延遲 < 100ms：✅ 可達成
- 記憶體使用：每連接約 4KB，1000 連接約 4MB

---

## MySQL + Redis 資料層設計

**研究問題**: 如何優化資料庫架構以支援高併發的配對查詢和聊天記錄？

### 決策：MySQL 主儲存 + Redis 快取 + 讀寫分離

**MySQL 架構**：

```sql
-- 核心表結構
users (用戶基本資料)
├── id, email, password_hash
├── birth_date (年齡驗證關鍵)
└── is_verified, is_active

user_profiles (用戶詳細資料)
├── user_id, display_name, bio
├── interests, photos
└── location_lat, location_lng

matches (配對記錄)
├── user1_id, user2_id
├── status (pending/matched/declined)
└── created_at, matched_at

chat_messages (聊天訊息)
├── match_id, sender_id
├── content, message_type
└── sent_at, is_read
```

**Redis 快取策略**：

```bash
user:profile:{user_id}     # 用戶資料快取 (TTL: 1小時)
match:pending:{user_id}    # 待處理配對 (TTL: 24小時)  
chat:online:{user_id}      # 在線狀態 (TTL: 5分鐘)
location:nearby:{lat}:{lng} # 附近用戶快取 (TTL: 30分鐘)
```

**查詢優化**：

1. **配對查詢**: 地理位置索引 + 年齡範圍索引
2. **聊天載入**: match_id + sent_at 複合索引
3. **用戶搜尋**: 興趣標籤 JSON 索引（MySQL 8.0+）

**替代方案考慮**：

1. **PostgreSQL**: JSONB 支援更好，但團隊 MySQL 經驗豐富
2. **MongoDB**: 適合非關聯資料，但事務支援有限
3. **純 Redis**: 效能最佳，但資料持久性風險

---

## 18+ 年齡驗證機制

**研究問題**: 如何實現符合法規的年齡驗證，平衡使用者體驗和合規要求？

### 決策：生日驗證 + 身份文件上傳（分階段實作）

#### 階段 1: 基礎驗證（MVP）

- 用戶註冊時輸入生日
- 系統計算年齡，拒絕未滿 18 歲用戶
- Email 驗證確保帳戶真實性

```go
func ValidateAge(birthDate time.Time) error {
    age := calculateAge(birthDate, time.Now())
    if age < 18 {
        return ErrUnderAge
    }
    return nil
}
```

#### 階段 2: 強化驗證（生產環境）

- 要求上傳身份證明文件
- 第三方身份驗證服務整合
- 定期重新驗證機制

**法規合規要求**：

1. **資料保護**: 生日資料加密存儲
2. **審核記錄**: 驗證過程完整日誌
3. **隱私權**: 最小必要資料原則

**替代方案考慮**：

1. **第三方 KYC 服務**: 成本高但最可靠
2. **信用卡驗證**: 間接驗證但非直接證明
3. **社群媒體整合**: 便利但可靠性低

---

## 雙向配對演算法

**研究問題**: 如何實現高效的雙向配對系統，確保只有雙方都同意才能建立聊天？

### 決策：狀態機 + 非同步處理模式

**配對狀態流程**：

```text
1. User A 對 User B 表示興趣 → status: "pending"
2. User B 查看 User A 的邀請
3a. User B 同意 → status: "matched" → 建立聊天室
3b. User B 拒絕 → status: "declined" → 結束流程
```

**資料結構設計**：

```go
type Match struct {
    ID        uint      `gorm:"primaryKey"`
    User1ID   uint      `gorm:"not null;index"`
    User2ID   uint      `gorm:"not null;index"`
    Status    MatchStatus `gorm:"not null;default:'pending'"`
    CreatedAt time.Time
    MatchedAt *time.Time
}

type MatchStatus string
const (
    StatusPending  MatchStatus = "pending"
    StatusMatched  MatchStatus = "matched" 
    StatusDeclined MatchStatus = "declined"
)
```

**防重複機制**：

- (user1_id, user2_id) 唯一索引，較小 ID 在前
- 避免 A→B 和 B→A 同時存在

**效能優化**：

- Redis 快取用戶的 pending 配對列表
- 批次處理配對通知，減少資料庫查詢

**替代方案考慮**：

1. **評分演算法**: 依相容度排序，但計算複雜
2. **推薦引擎**: ML 驅動配對，但初期資料不足
3. **地理位置優先**: 簡單但可能限制配對機會

---

## 併發用戶處理模式

**研究問題**: 如何確保 1000 併發用戶的系統穩定性和效能？

### 決策：連接池 + 限流 + 優雅降級

**併發控制策略**：

```go
type ConcurrencyManager struct {
    maxConnections int
    activeUsers    sync.Map        // 在線用戶追蹤
    rateLimiter    *rate.Limiter   // 請求限流
    circuitBreaker *breaker.Breaker // 熔斷器
}
```

**效能優化措施**：

1. **連接復用**: 資料庫連接池，避免頻繁建立連接
2. **請求限流**: 每用戶每分鐘最多 60 請求
3. **熔斷機制**: 資料庫故障時自動降級
4. **資源監控**: CPU/記憶體/連接數即時監控

**擴容策略**：

- **垂直擴容**: 增加伺服器規格（初期方案）
- **水平擴容**: Load Balancer + 多實例（成長後）
- **讀寫分離**: MySQL Master-Slave 架構

**壓力測試計畫**：

```bash
# 模擬 1000 併發 WebSocket 連接
wstest -c 1000 -d 30s ws://localhost:8080/ws

# 模擬配對請求壓力測試
ab -n 10000 -c 100 http://localhost:8080/api/match
```

**替代方案考慮**：

1. **微服務拆分**: 聊天、配對、用戶各自獨立，但複雜度增加
2. **訊息佇列**: RabbitMQ/Kafka 處理異步任務，但引入額外依賴
3. **CDN 加速**: 靜態資源 CDN，但成本考量

---

## 前後端整合策略

**研究問題**: 如何設計 RESTful API 和 WebSocket 協定，確保前後端順暢整合？

### 決策：RESTful API + WebSocket 雙協定

**API 設計原則**：

```typescript
// RESTful API 負責 CRUD 操作
POST   /api/auth/register    // 註冊
POST   /api/auth/login       // 登入
GET    /api/users/profile    // 取得個人資料
POST   /api/matches/like     // 表示興趣
GET    /api/matches/pending  // 待處理配對

// WebSocket 負責即時通訊
{
  "type": "message_send",
  "data": {
    "match_id": 123,
    "content": "Hello!"
  }
}
```

**前端架構**：

- **Vanilla JavaScript**: 避免框架複雜性，直接 DOM 操作
- **模組化設計**: 每頁面獨立 JS 檔案
- **WebSocket 客戶端**: 自動重連 + 心跳機制

**錯誤處理策略**：

```javascript
// 統一錯誤處理
class ApiClient {
    async request(url, options) {
        try {
            const response = await fetch(url, options);
            if (!response.ok) {
                throw new ApiError(response.status, await response.text());
            }
            return await response.json();
        } catch (error) {
            this.handleError(error);
            throw error;
        }
    }
}
```

**安全性考量**：

1. **CSRF 防護**: SameSite cookie + CSRF token
2. **XSS 防護**: 所有用戶輸入 HTML 編碼
3. **JWT 認證**: Access token + Refresh token 機制

**替代方案考慮**：

1. **GraphQL**: 減少請求次數，但學習曲線陡峭
2. **gRPC-Web**: 效能優異，但瀏覽器支援有限
3. **Server-Sent Events**: 單向通訊替代，但功能受限

---

## 總結建議

### 優先實作順序

1. **Phase 1**: 基礎用戶註冊登入 + 年齡驗證
2. **Phase 2**: 用戶資料管理 + 基礎配對功能  
3. **Phase 3**: WebSocket 即時聊天實作
4. **Phase 4**: 效能優化 + 擴容準備

### 風險評估

- **高風險**: WebSocket 併發處理，需要充分測試
- **中風險**: 資料庫效能優化，可透過索引解決
- **低風險**: 前端整合，技術相對成熟

### 技術決策摘要

所有研究項目的技術選型都已確定，沒有遺留的 "NEEDS CLARIFICATION" 項目。架構設計符合憲章要求，可以進入 Phase 1 設計階段。
