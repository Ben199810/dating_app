package repository

import (
	"golang_dev_docker/domain/entity"
)

// UserRepository 用戶儲存庫接口
type UserRepository interface {
	// ExistsByUsername 檢查用戶名是否已存在
	ExistsByUsername(username string) (bool, error)

	// ExistsByEmail 檢查電子郵件是否已存在
	ExistsByEmail(email string) (bool, error)

	// Save 保存用戶
	Save(user *entity.User) error

	// FindByID 根據ID查找用戶
	FindByID(id string) (*entity.User, error)

	// FindByUsername 根據用戶名查找用戶
	FindByUsername(username string) (*entity.User, error)

	// FindByEmail 根據電子郵件查找用戶
	FindByEmail(email string) (*entity.User, error)
}
