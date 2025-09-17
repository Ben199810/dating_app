package repository

import "golang_dev_docker/domain/entity"

// 定義註冊使用者的實作介面
type RegisterUser interface {
	// 先查詢是否有 user 資料
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	// 新增 user 資料到資料庫
	RegisterUser(user *entity.User) error
}
