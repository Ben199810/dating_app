package repository

import "golang_dev_docker/domain/entity"

// AuthRepository 定義認證相關的儲存庫操作
type AuthRepository interface {
	// 根據 email 獲取用戶用於登入驗證
	GetUserByEmail(email string) (*entity.UserInformation, error)

	// 根據 username 獲取用戶用於登入驗證
	GetUserByUsername(username string) (*entity.UserInformation, error)

	// 更新最後登入時間
	UpdateLastLoginTime(userID int) error

	// 檢查用戶是否存在（用於註冊時檢查）
	UserExists(email, username string) (bool, error)
}
