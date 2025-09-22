package repository

import "golang_dev_docker/domain/entity"

// UserRepository 定義用戶相關的儲存庫會具備哪些操作
type UserRepository interface {
	Create(user *entity.UserInformation) error
	GetByEmail(email string) (*entity.UserInformation, error)
	GetByID(id int) (*entity.UserInformation, error)
	Update(user *entity.UserInformation) error
	Delete(id int) error

	// 非敏感查詢
	GetUserProfile(id int) (*entity.UserInformation, error)
	UpdateUserProfile(user *entity.UserInformation) error
	
	// 新增的交友軟體功能
	GetUsersByLocation(lat, lng float64, radiusKm int, limit int) ([]*entity.UserInformation, error)
	GetUsersByAgeRange(minAge, maxAge int, limit int) ([]*entity.UserInformation, error)
	GetUsersByGender(gender entity.Gender, limit int) ([]*entity.UserInformation, error)
	UpdateLastActiveTime(userID int) error
	SearchUsers(filters map[string]interface{}, limit, offset int) ([]*entity.UserInformation, error)
}
