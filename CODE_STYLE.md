# 📝 程式碼風格規範

> 🎯 **統一的程式碼風格是團隊協作的基石**
>
> 本文件定義了專案的程式碼風格標準，確保程式碼的可讀性和一致性。

![Go Style](https://img.shields.io/badge/Style-Go%20Standard-00ADD8?style=flat-square&logo=go&logoColor=white)
![Code Quality](https://img.shields.io/badge/Quality-High-4CAF50?style=flat-square&logo=codeclimate&logoColor=white)

---

## 📏 命名規範

### 🔤 變數命名

| 類型 | 規範 | 範例 | 說明 |
|------|------|------|------|
| **私有變數** | 小駝峰式 (camelCase) | `userName`, `userAge` | 包內部使用 |
| **公開變數** | 大駝峰式 (PascalCase) | `UserName`, `UserAge` | 對外公開 |
| **常數** | 全大寫 + 底線 | `MAX_RETRY_COUNT` | 不可變值 |
| **全域變數** | 前綴 + 大駝峰 | `GlobalUserCount` | 避免使用 |

```go
// ✅ 正確示例
var userName string           // 私有變數
var UserID int               // 公開變數
const MAX_CONNECTIONS = 100   // 常數

// ❌ 錯誤示例
var UserName string          // 私有變數不應大寫
var user_id int             // 不使用底線
```

### 🔧 函式命名

| 類型 | 規範 | 範例 | 說明 |
|------|------|------|------|
| **私有函式** | 小駝峰式 + 動詞開頭 | `getUser()`, `validateEmail()` | 包內部使用 |
| **公開函式** | 大駝峰式 + 動詞開頭 | `GetUser()`, `CreateUser()` | 對外公開 |
| **建構函式** | `New` + 類型名稱 | `NewUserService()` | 物件建立 |
| **測試函式** | `Test` + 函式名 | `TestGetUser()` | 單元測試 |

```go
// ✅ 正確示例
func getUserByID(id int) *User {...}      // 私有函式
func CreateUser(req *Request) *User {...} // 公開函式
func NewUserService() *UserService {...}  // 建構函式

// ❌ 錯誤示例
func GetUserByID(id int) *User {...}      // 私有函式不應大寫
func create_user(req *Request) *User {...} // 不使用底線
```

### 🏗️ 類型命名

| 類型 | 規範 | 範例 | 說明 |
|------|------|------|------|
| **結構體** | 大駝峰式 + 名詞 | `User`, `UserProfile` | 資料結構 |
| **介面** | 大駝峰式 + `er`後綴 | `UserRepository`, `Validator` | 行為定義 |
| **錯誤** | `Err` + 描述 | `ErrUserNotFound` | 錯誤類型 |

```go
// ✅ 正確示例
type User struct {...}
type UserRepository interface {...}
var ErrUserNotFound = errors.New("user not found")

// ❌ 錯誤示例
type user struct {...}              // 結構體應大寫
type UserRepositoryInterface {...}  // 避免 Interface 後綴
```

---

## 📁 檔案與套件規範

### 📂 檔案命名

| 類型 | 規範 | 範例 | 說明 |
|------|------|------|------|
| **Go 檔案** | 小寫 + 底線 | `user_service.go` | 功能描述 |
| **測試檔案** | 檔名 + `_test` | `user_service_test.go` | 測試檔案 |
| **範例檔案** | 檔名 + `_example` | `user_service_example.go` | 範例程式 |

### 📦 套件命名

```go
// ✅ 正確示例
package user          // 簡潔明確
package repository    // 功能導向
package handler       // 職責清楚

// ❌ 錯誤示例
package userService   // 避免駝峰
package user_repo     // 避免底線
package utils         // 過於泛用
```

---

## 🎨 程式碼格式化

### 🔧 自動格式化工具

```bash
# 格式化程式碼
go fmt ./...

# 整理 import
go mod tidy
goimports -w .

# 靜態分析
go vet ./...
golint ./...
```

### 📝 註解規範

```go
// ✅ 公開函式必須有註解
// GetUser 根據 ID 獲取使用者資訊
// 如果使用者不存在，回傳 ErrUserNotFound 錯誤
func GetUser(id int) (*User, error) {
    // 私有函式內的重要邏輯也需要註解
    if id <= 0 {
        return nil, ErrInvalidUserID
    }
    // ...
}

// ✅ 結構體註解
// User 代表系統中的使用者實體
// 包含使用者的基本資訊和狀態
type User struct {
    ID       int    `json:"id"`        // 使用者唯一識別碼
    Username string `json:"username"`  // 使用者名稱
    Email    string `json:"email"`     // 電子郵件地址
}
```

---

## 🏗️ 架構規範

### 📂 DDD 層次劃分

```text
📦 專案結構
├── 🎯 domain/          # 領域層 - 業務邏輯核心
│   ├── entity/        # 實體 - 業務物件
│   ├── repository/    # 儲存庫介面 - 資料存取抽象
│   └── service/       # 領域服務 - 複雜業務邏輯
├── 🔧 infrastructure/ # 基礎設施層 - 技術實作
│   └── mysql/        # 資料庫實作
├── 🌐 server/         # 展示層 - 對外介面
│   └── handler/      # HTTP 處理器
└── 🧩 component/     # 共用元件
    └── validator/    # 驗證器
```

### 🔄 依賴規則

```go
// ✅ 正確：基礎設施層依賴領域層
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) GetUser(id int) (*entity.User, error) {
    // 實作細節...
}

// ❌ 錯誤：領域層不應依賴基礎設施層
func (u *User) SaveToDB(db *sql.DB) error {
    // 領域實體不應知道資料庫細節
}
```

---

## ✅ 程式碼品質檢查

### 🧪 測試規範

```go
// ✅ 測試函式命名
func TestUserService_CreateUser(t *testing.T) {
    // 測試邏輯...
}

func TestUserService_CreateUser_WithInvalidEmail(t *testing.T) {
    // 特定情況測試...
}

// ✅ 表格驅動測試
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        expected bool
    }{
        {"valid email", "user@example.com", true},
        {"invalid email", "invalid-email", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ValidateEmail(tt.email)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 📊 品質指標

| 指標 | 目標值 | 工具 | 說明 |
|------|--------|------|------|
| **測試覆蓋率** | ≥ 80% | `go test -cover` | 程式碼測試覆蓋度 |
| **循環複雜度** | ≤ 10 | `gocyclo` | 函式複雜度控制 |
| **程式碼重複** | ≤ 5% | `dupl` | 避免程式碼重複 |
| **技術債務** | A級 | `SonarQube` | 程式碼品質評估 |

---

## 🔧 開發工具配置

### VS Code 設定

```json
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.vetOnSave": "package",
    "go.testTimeout": "30s",
    "editor.formatOnSave": true
}
```

### Git Hooks

```bash
#!/bin/sh
# pre-commit hook
go fmt ./...
go vet ./...
go test ./...
```

---

## 📚 參考資源

| 資源 | 連結 | 說明 |
|------|------|------|
| **Go 官方風格** | [Effective Go](https://golang.org/doc/effective_go.html) | 官方程式碼風格指南 |
| **Google 風格** | [Go Style Guide](https://google.github.io/styleguide/go/) | Google Go 風格指南 |
| **Uber 風格** | [Uber Go Style Guide](https://github.com/uber-go/guide) | Uber Go 最佳實踐 |
