# 資料模型設計：18+ 交友聊天應用程式

*Phase 1 輸出 - 基於功能規格提取的實體模型*

## 核心實體總覽

基於功能規格分析，系統包含 9 個關鍵實體：

1. **User** - 用戶基本資料
2. **UserProfile** - 用戶詳細檔案  
3. **Match** - 配對記錄
4. **ChatMessage** - 聊天訊息
5. **Report** - 檢舉記錄
6. **Block** - 封鎖記錄
7. **Interest** - 興趣標籤
8. **Photo** - 用戶照片
9. **AgeVerification** - 年齡驗證

---

## 1. User（用戶）

**目的**: 儲存用戶基本帳戶資訊和認證資料

### 欄位

```go
type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Email        string    `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"`
    BirthDate    time.Time `gorm:"not null" json:"birth_date"`
    IsVerified   bool      `gorm:"default:false" json:"is_verified"`
    IsActive     bool      `gorm:"default:true" json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### 驗證規則

- **Email**: 必填，唯一，有效 email 格式
- **PasswordHash**: 必填，bcrypt 雜湊，最少 60 字元
- **BirthDate**: 必填，用戶必須年滿 18 歲
- **IsVerified**: 預設 false，通過年齡驗證後設為 true
- **IsActive**: 預設 true，帳戶停用時設為 false

### 關聯

- **一對一**: UserProfile
- **一對一**: AgeVerification  
- **一對多**: Photo (用戶照片)
- **多對多**: Match (透過 User1ID, User2ID)

---

## 2. UserProfile（用戶檔案）

**目的**: 儲存用戶展示資訊和配對相關屬性

### 欄位

```go
type UserProfile struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"`
    DisplayName string    `gorm:"not null;size:50" json:"display_name"`
    Bio         string    `gorm:"size:500" json:"bio"`
    Gender      Gender    `gorm:"not null" json:"gender"`
    ShowAge     bool      `gorm:"default:true" json:"show_age"`
    LocationLat *float64  `json:"location_lat"`
    LocationLng *float64  `json:"location_lng"`
    MaxDistance int       `gorm:"default:50" json:"max_distance"` // km
    AgeRangeMin int       `gorm:"default:18" json:"age_range_min"`
    AgeRangeMax int       `gorm:"default:99" json:"age_range_max"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // 關聯
    User User `gorm:"constraint:OnDelete:CASCADE" json:"user"`
}

type Gender string
const (
    GenderMale   Gender = "male"
    GenderFemale Gender = "female"
    GenderOther  Gender = "other"
)
```

### 驗證規則

- **UserID**: 必填，外鍵關聯 User.ID
- **DisplayName**: 必填，1-50 字元，用於顯示
- **Bio**: 選填，最多 500 字元
- **Gender**: 必填，限制為 male/female/other  
- **LocationLat/Lng**: 選填，地理座標 (-90~90, -180~180)
- **MaxDistance**: 預設 50km，配對範圍限制
- **AgeRange**: Min >= 18, Max <= 99, Min <= Max

### 關聯

- **屬於**: User
- **多對多**: Interest (透過 user_interests 中間表)

---

## 3. Match（配對）

**目的**: 管理用戶間的雙向配對狀態

### 欄位

```go
type Match struct {
    ID        uint        `gorm:"primaryKey" json:"id"`
    User1ID   uint        `gorm:"not null;index" json:"user1_id"`
    User2ID   uint        `gorm:"not null;index" json:"user2_id"`
    Status    MatchStatus `gorm:"not null;default:'pending'" json:"status"`
    CreatedAt time.Time   `json:"created_at"`
    MatchedAt *time.Time  `json:"matched_at"`
    
    // 關聯
    User1 User `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE" json:"user1"`
    User2 User `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE" json:"user2"`
}

type MatchStatus string
const (
    StatusPending  MatchStatus = "pending"  // 等待對方回應
    StatusMatched  MatchStatus = "matched"  // 雙方配對成功
    StatusDeclined MatchStatus = "declined" // 對方拒絕配對
)
```

### 驗證規則

- **User1ID, User2ID**: 必填，且 User1ID < User2ID（防重複）
- **Status**: 預設 pending，只能是三個定義值之一
- **(User1ID, User2ID)**: 複合唯一索引
- **MatchedAt**: Status = matched 時必填

### 狀態轉換

```text
pending → matched   (User2 接受配對)
pending → declined  (User2 拒絕配對)
```

### 關聯

- **屬於**: User (User1, User2)
- **一對多**: ChatMessage

---

## 4. ChatMessage（聊天訊息）

**目的**: 儲存配對用戶間的聊天記錄

### 欄位

```go
type ChatMessage struct {
    ID          uint            `gorm:"primaryKey" json:"id"`
    MatchID     uint            `gorm:"not null;index" json:"match_id"`
    SenderID    uint            `gorm:"not null;index" json:"sender_id"`
    Content     string          `gorm:"not null;size:1000" json:"content"`
    MessageType MessageType     `gorm:"default:'text'" json:"message_type"`
    IsRead      bool            `gorm:"default:false" json:"is_read"`
    SentAt      time.Time       `gorm:"not null" json:"sent_at"`
    
    // 關聯
    Match  Match `gorm:"constraint:OnDelete:CASCADE" json:"match"`
    Sender User  `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sender"`
}

type MessageType string
const (
    MessageText  MessageType = "text"
    MessageImage MessageType = "image"
    MessageEmoji MessageType = "emoji"
)
```

### 驗證規則

- **MatchID**: 必填，必須是 Status = matched 的配對
- **SenderID**: 必填，必須是該配對的參與者之一
- **Content**: 必填，1-1000 字元
- **MessageType**: 預設 text，限制為定義值
- **SentAt**: 必填，訊息發送時間

### 索引優化

- **(MatchID, SentAt)**: 複合索引，支援聊天記錄查詢
- **SenderID**: 索引，支援用戶訊息查詢

---

## 5. Report（檢舉）

**目的**: 記錄用戶檢舉行為和處理狀態

### 欄位

```go
type Report struct {
    ID           uint         `gorm:"primaryKey" json:"id"`
    ReporterID   uint         `gorm:"not null;index" json:"reporter_id"`
    ReportedID   uint         `gorm:"not null;index" json:"reported_id"`
    ReportType   ReportType   `gorm:"not null" json:"report_type"`
    Reason       string       `gorm:"size:500" json:"reason"`
    Status       ReportStatus `gorm:"default:'pending'" json:"status"`
    CreatedAt    time.Time    `json:"created_at"`
    ProcessedAt  *time.Time   `json:"processed_at"`
    AdminNote    string       `gorm:"size:500" json:"admin_note"`
    
    // 關聯
    Reporter User `gorm:"foreignKey:ReporterID;constraint:OnDelete:CASCADE" json:"reporter"`
    Reported User `gorm:"foreignKey:ReportedID;constraint:OnDelete:CASCADE" json:"reported"`
}

type ReportType string
const (
    ReportSpam        ReportType = "spam"
    ReportHarassment  ReportType = "harassment"
    ReportFakeProfile ReportType = "fake_profile"
    ReportInappropriate ReportType = "inappropriate"
    ReportOther       ReportType = "other"
)

type ReportStatus string
const (
    ReportPending  ReportStatus = "pending"
    ReportReviewed ReportStatus = "reviewed"
    ReportResolved ReportStatus = "resolved"
)
```

### 驗證規則

- **ReporterID ≠ ReportedID**: 不能檢舉自己
- **ReportType**: 必填，限制為定義值
- **Reason**: 選填，最多 500 字元詳細說明
- **(ReporterID, ReportedID, ReportType)**: 複合唯一索引（防重複檢舉）

---

## 6. Block（封鎖）

**目的**: 管理用戶間的封鎖關係

### 欄位

```go
type Block struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    BlockerID  uint      `gorm:"not null;index" json:"blocker_id"`
    BlockedID  uint      `gorm:"not null;index" json:"blocked_id"`
    CreatedAt  time.Time `json:"created_at"`
    
    // 關聯
    Blocker User `gorm:"foreignKey:BlockerID;constraint:OnDelete:CASCADE" json:"blocker"`
    Blocked User `gorm:"foreignKey:BlockedID;constraint:OnDelete:CASCADE" json:"blocked"`
}
```

### 驗證規則

- **BlockerID ≠ BlockedID**: 不能封鎖自己
- **(BlockerID, BlockedID)**: 複合唯一索引
- **影響**: 封鎖後不會在配對中出現，無法收發訊息

---

## 7. Interest（興趣）

**目的**: 預定義的興趣標籤，用於配對推薦

### 欄位

```go
type Interest struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Name        string `gorm:"uniqueIndex;not null;size:50" json:"name"`
    Category    string `gorm:"not null;size:50" json:"category"`
    IsActive    bool   `gorm:"default:true" json:"is_active"`
    DisplayOrder int   `gorm:"default:0" json:"display_order"`
}

// 中間表
type UserInterest struct {
    UserID     uint `gorm:"primaryKey" json:"user_id"`
    InterestID uint `gorm:"primaryKey" json:"interest_id"`
}
```

### 預設資料

```text
運動: 健身, 瑜伽, 游泳, 籃球, 足球, 網球
娛樂: 電影, 音樂, 閱讀, 遊戲, 攝影, 旅遊
生活: 美食, 烹飪, 咖啡, 購物, 寵物, 園藝
文化: 藝術, 博物館, 演唱會, 戲劇, 語言學習
戶外: 登山, 露營, 釣魚, 海邊, 騎車, 慢跑
```

---

## 8. Photo（用戶照片）

**目的**: 管理用戶個人檔案照片

### 欄位

```go
type Photo struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    UserID      uint   `gorm:"not null;index" json:"user_id"`
    URL         string `gorm:"not null;size:500" json:"url"`
    IsPrimary   bool   `gorm:"default:false" json:"is_primary"`
    DisplayOrder int   `gorm:"default:0" json:"display_order"`
    CreatedAt   time.Time `json:"created_at"`
    
    // 關聯
    User User `gorm:"constraint:OnDelete:CASCADE" json:"user"`
}
```

### 驗證規則

- **UserID**: 每個用戶最多 6 張照片
- **URL**: 必填，有效的圖片 URL
- **IsPrimary**: 每個用戶只能有一張主要照片
- **DisplayOrder**: 照片排序，0-5

---

## 9. AgeVerification（年齡驗證）

**目的**: 強化年齡驗證記錄

### 欄位

```go
type AgeVerification struct {
    ID             uint                  `gorm:"primaryKey" json:"id"`
    UserID         uint                  `gorm:"uniqueIndex;not null" json:"user_id"`
    Method         VerificationMethod    `gorm:"not null" json:"method"`
    DocumentURL    string                `gorm:"size:500" json:"document_url"`
    Status         VerificationStatus    `gorm:"default:'pending'" json:"status"`
    VerifiedAt     *time.Time            `json:"verified_at"`
    RejectionReason string               `gorm:"size:200" json:"rejection_reason"`
    CreatedAt      time.Time             `json:"created_at"`
    
    // 關聯
    User User `gorm:"constraint:OnDelete:CASCADE" json:"user"`
}

type VerificationMethod string
const (
    MethodBirthDate VerificationMethod = "birth_date"
    MethodIDCard    VerificationMethod = "id_card"
    MethodPassport  VerificationMethod = "passport"
    MethodLicense   VerificationMethod = "license"
)

type VerificationStatus string
const (
    VerificationPending  VerificationStatus = "pending"
    VerificationApproved VerificationStatus = "approved"
    VerificationRejected VerificationStatus = "rejected"
)
```

---

## 資料庫關聯圖

```text
User (1:1) ←→ UserProfile
User (1:1) ←→ AgeVerification  
User (1:n) ←→ Photo
User (n:m) ←→ Interest (透過 UserInterest)

User (n:m) ←→ Match (透過 User1ID, User2ID)
Match (1:n) ←→ ChatMessage

User (1:n) ←→ Report (作為 Reporter)
User (1:n) ←→ Report (作為 Reported)
User (1:n) ←→ Block (作為 Blocker)  
User (1:n) ←→ Block (作為 Blocked)
```

## 效能索引策略

### 核心查詢索引

1. **配對查詢**: `(gender, age_range_min, age_range_max, location_lat, location_lng)`
2. **聊天載入**: `(match_id, sent_at DESC)`
3. **用戶搜尋**: `(is_active, is_verified)`
4. **封鎖檢查**: `(blocker_id, blocked_id)`

### Redis 快取策略

```text
用戶資料: user:profile:{user_id} (TTL: 1小時)
配對列表: match:pending:{user_id} (TTL: 24小時)
在線狀態: user:online:{user_id} (TTL: 5分鐘)
附近用戶: location:nearby:{lat}:{lng} (TTL: 30分鐘)
```

## 安全性考量

### 資料保護

- **PasswordHash**: bcrypt 加密，never 回傳給前端
- **BirthDate**: 僅用於年齡計算，前端僅顯示年齡
- **LocationLat/Lng**: 精確度控制，只顯示大概區域

### 隱私控制

- **ShowAge**: 用戶可選擇是否顯示年齡
- **Block**: 封鎖用戶不會出現在任何推薦中
- **Report**: 檢舉記錄完整保存，支援追蹤

這個資料模型設計完全符合功能規格要求，支援所有核心功能：註冊登入、年齡驗證、個人檔案、配對系統、即時聊天、檢舉封鎖。所有實體間的關聯關係明確，索引策略支援高效查詢。
