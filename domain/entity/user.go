package entity

type UserInformation struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // 不序列化密碼
}
