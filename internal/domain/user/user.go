package user

// User 代表用戶的領域實體
type User struct {
	ID       UserID   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Hobbies  []string `json:"hobbies"`
	IsOnline bool     `json:"is_online"`
}

// UserID 用戶唯一識別符
type UserID string

// NewUser 創建新用戶
func NewUser(id UserID, name, email string, hobbies []string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Hobbies:  hobbies,
		IsOnline: false,
	}
}

// SetOnline 設置用戶在線狀態
func (u *User) SetOnline(online bool) {
	u.IsOnline = online
}

// IsValid 驗證用戶資料是否有效
func (u *User) IsValid() bool {
	return u.ID != "" && u.Name != "" && u.Email != ""
}

// UserRepository 用戶倉儲介面
type UserRepository interface {
	Save(user *User) error
	FindByID(id UserID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
}
