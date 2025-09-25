package repository

import (
	"context"

	"golang_dev_docker/domain/entity"
)

// UserRepository 用戶數據儲存庫介面
// 提供用戶基本資料的持久化操作，包括註冊、認證、查詢等功能
type UserRepository interface {
	// Create 創建新用戶
	// 用於用戶註冊流程，保存基本帳戶資訊
	Create(ctx context.Context, user *entity.User) error

	// GetByID 根據 ID 獲取用戶
	// 用於身份驗證和用戶資訊查詢
	GetByID(ctx context.Context, id uint) (*entity.User, error)

	// GetByEmail 根據 Email 獲取用戶
	// 用於登入驗證和重複註冊檢查
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update 更新用戶資訊
	// 用於修改用戶基本資料（如驗證狀態、啟用狀態等）
	Update(ctx context.Context, user *entity.User) error

	// Delete 軟刪除用戶（設為非啟用狀態）
	// 用於用戶停用帳戶
	Delete(ctx context.Context, id uint) error

	// SetVerified 設定用戶驗證狀態
	// 用於年齡驗證完成後更新狀態
	SetVerified(ctx context.Context, id uint, verified bool) error

	// SetActive 設定用戶啟用狀態
	// 用於帳戶啟用/停用管理
	SetActive(ctx context.Context, id uint, active bool) error
}

// UserProfileRepository 用戶檔案數據儲存庫介面
// 提供用戶展示資訊的持久化操作，包括個人檔案、偏好設定等功能
type UserProfileRepository interface {
	// Create 創建用戶檔案
	// 用於註冊完成後建立展示檔案
	Create(ctx context.Context, profile *entity.UserProfile) error

	// GetByUserID 根據用戶 ID 獲取檔案
	// 用於檔案展示和編輯
	GetByUserID(ctx context.Context, userID uint) (*entity.UserProfile, error)

	// Update 更新用戶檔案
	// 用於檔案資訊修改
	Update(ctx context.Context, profile *entity.UserProfile) error

	// Delete 刪除用戶檔案
	// 用於帳戶刪除時清理相關資料
	Delete(ctx context.Context, userID uint) error

	// UpdateLocation 更新用戶位置資訊
	// 用於地理位置配對功能
	UpdateLocation(ctx context.Context, userID uint, lat, lng *float64) error

	// UpdateMatchingPreferences 更新配對偏好設定
	// 用於修改配對範圍、年齡偏好等設定
	UpdateMatchingPreferences(ctx context.Context, userID uint, maxDistance, ageMin, ageMax int) error
}

// PhotoRepository 用戶照片數據儲存庫介面
// 提供用戶照片的持久化操作，包括上傳、排序、刪除等功能
type PhotoRepository interface {
	// Create 添加用戶照片
	// 用於照片上傳功能
	Create(ctx context.Context, photo *entity.Photo) error

	// GetByUserID 獲取用戶所有照片
	// 用於檔案展示和照片管理
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Photo, error)

	// GetByID 根據 ID 獲取照片
	// 用於照片操作和權限驗證
	GetByID(ctx context.Context, id uint) (*entity.Photo, error)

	// Update 更新照片資訊
	// 用於修改照片描述、排序等
	Update(ctx context.Context, photo *entity.Photo) error

	// Delete 刪除照片
	// 用於照片移除功能
	Delete(ctx context.Context, id uint) error

	// SetPrimary 設定主要照片
	// 用於變更用戶主要展示照片
	SetPrimary(ctx context.Context, userID, photoID uint) error

	// UpdateOrder 更新照片排序
	// 用於照片順序調整功能
	UpdateOrder(ctx context.Context, userID uint, photoOrders []struct {
		PhotoID uint
		Order   int
	}) error
}

// InterestRepository 興趣標籤數據儲存庫介面
// 提供興趣標籤的持久化操作，支援配對演算法
type InterestRepository interface {
	// GetAll 獲取所有可用興趣標籤
	// 用於註冊和檔案編輯時的選項展示
	GetAll(ctx context.Context) ([]*entity.Interest, error)

	// GetByUserID 獲取用戶的興趣標籤
	// 用於檔案展示和配對演算法
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Interest, error)

	// SetUserInterests 設定用戶興趣標籤
	// 用於更新用戶興趣偏好
	SetUserInterests(ctx context.Context, userID uint, interestIDs []uint) error

	// Create 創建新興趣標籤
	// 用於系統管理功能
	Create(ctx context.Context, interest *entity.Interest) error

	// Update 更新興趣標籤
	// 用於系統管理功能
	Update(ctx context.Context, interest *entity.Interest) error

	// Delete 刪除興趣標籤
	// 用於系統管理功能
	Delete(ctx context.Context, id uint) error
}

// AgeVerificationRepository 年齡驗證數據儲存庫介面
// 提供年齡驗證記錄的持久化操作，確保18+合規要求
type AgeVerificationRepository interface {
	// Create 創建年齡驗證記錄
	// 用於用戶註冊時的年齡驗證流程
	Create(ctx context.Context, verification *entity.AgeVerification) error

	// GetByUserID 根據用戶 ID 獲取驗證記錄
	// 用於檢查用戶驗證狀態
	GetByUserID(ctx context.Context, userID uint) (*entity.AgeVerification, error)

	// Update 更新驗證記錄
	// 用於驗證狀態變更和審核結果更新
	Update(ctx context.Context, verification *entity.AgeVerification) error

	// GetPendingVerifications 獲取待審核的驗證記錄
	// 用於管理員審核功能
	GetPendingVerifications(ctx context.Context, limit int) ([]*entity.AgeVerification, error)

	// SetVerificationStatus 設定驗證狀態
	// 用於審核通過或拒絕操作
	SetVerificationStatus(ctx context.Context, userID uint, status entity.VerificationStatus, reviewerID *uint, notes string) error
}
