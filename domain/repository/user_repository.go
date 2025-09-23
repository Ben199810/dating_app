package repository

import "golang_dev_docker/domain/entity"

// UserRepository 定義用戶相關的儲存庫會具備哪些操作
type UserRepository interface {
	CreateUser(user *entity.UserInformation) error
	GetUserByEmail(email string) (*entity.UserInformation, error)
	GetUserByID(id int) (*entity.UserInformation, error)
	UpdateUser(user *entity.UserInformation) error
	DeleteUser(id int) error
	
	// 新增的交友軟體功能
	GetUsersByLocation(lat, lng float64, radiusKm int, limit int) ([]*entity.UserInformation, error)
	GetUsersByAgeRange(minAge, maxAge int, limit int) ([]*entity.UserInformation, error)
	GetUsersByGender(gender entity.Gender, limit int) ([]*entity.UserInformation, error)
	UpdateLastActiveTime(userID int) error
	SearchUsers(filters map[string]interface{}, limit, offset int) ([]*entity.UserInformation, error)
}
