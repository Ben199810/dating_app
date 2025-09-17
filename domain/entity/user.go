package entity

import "time"

type UserInformation struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // 不序列化密碼
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
