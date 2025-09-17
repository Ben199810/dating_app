package repository

import "golang_dev_docker/domain/entity"

// UserRepository 定義用戶資料庫操作的介面
type UserRepository interface {
	Create(user *entity.UserInformation) error
	GetByEmail(email string) (*entity.UserInformation, error)
	GetByID(id int) (*entity.UserInformation, error)
	Update(user *entity.UserInformation) error
	Delete(id int) error
}
