package entity

import (
	"errors"
	"time"
)

// User 用戶基本資料實體
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	BirthDate    time.Time `gorm:"not null" json:"birth_date"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 關聯 - 將在其他實體建立後添加
	// Profile        *UserProfile     `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	// Photos         []Photo          `gorm:"foreignKey:UserID" json:"photos,omitempty"`
	// AgeVerification *AgeVerification `gorm:"foreignKey:UserID" json:"age_verification,omitempty"`
}

// IsAdult 檢查用戶是否年滿 18 歲
func (u *User) IsAdult() bool {
	now := time.Now()
	age := now.Year() - u.BirthDate.Year()

	// 調整年齡計算（如果還沒到生日）
	if now.YearDay() < u.BirthDate.YearDay() {
		age--
	}

	return age >= 18
}

// GetAge 計算用戶當前年齡
func (u *User) GetAge() int {
	now := time.Now()
	age := now.Year() - u.BirthDate.Year()

	// 調整年齡計算（如果還沒到生日）
	if now.YearDay() < u.BirthDate.YearDay() {
		age--
	}

	return age
}

// IsEligible 檢查用戶是否符合使用應用程式的條件
func (u *User) IsEligible() bool {
	return u.IsActive && u.IsVerified && u.IsAdult()
}

// Validate 驗證用戶資料的完整性
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email 是必填欄位")
	}

	if u.PasswordHash == "" {
		return errors.New("password 是必填欄位")
	}

	if u.BirthDate.IsZero() {
		return errors.New("birth_date 是必填欄位")
	}

	// 檢查年齡限制
	if !u.IsAdult() {
		return errors.New("用戶必須年滿 18 歲")
	}

	return nil
}

// Deactivate 停用用戶帳戶
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate 啟用用戶帳戶
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// MarkAsVerified 標記用戶已通過年齡驗證
func (u *User) MarkAsVerified() {
	u.IsVerified = true
	u.UpdatedAt = time.Now()
}
