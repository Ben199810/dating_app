package mysql

import (
	"context"
	"fmt"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"

	"gorm.io/gorm"
)

// MySQLUserRepository MySQL 用戶儲存庫實作
type MySQLUserRepository struct {
	db *gorm.DB
}

// NewUserRepository 創建新的 MySQL 用戶儲存庫
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &MySQLUserRepository{db: db}
}

// Create 創建新用戶
func (r *MySQLUserRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("創建用戶失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 獲取用戶
func (r *MySQLUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用戶不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢用戶失敗: %w", err)
	}
	return &user, nil
}

// GetByEmail 根據 Email 獲取用戶
func (r *MySQLUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用戶不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢用戶失敗: %w", err)
	}
	return &user, nil
}

// Update 更新用戶資訊
func (r *MySQLUserRepository) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("更新用戶失敗: %w", err)
	}
	return nil
}

// Delete 軟刪除用戶（設為非啟用狀態）
func (r *MySQLUserRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return fmt.Errorf("刪除用戶失敗: %w", err)
	}
	return nil
}

// SetVerified 設定用戶驗證狀態
func (r *MySQLUserRepository) SetVerified(ctx context.Context, id uint, verified bool) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("is_verified", verified).Error; err != nil {
		return fmt.Errorf("設定驗證狀態失敗: %w", err)
	}
	return nil
}

// SetActive 設定用戶啟用狀態
func (r *MySQLUserRepository) SetActive(ctx context.Context, id uint, active bool) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("is_active", active).Error; err != nil {
		return fmt.Errorf("設定啟用狀態失敗: %w", err)
	}
	return nil
}

// MySQLUserProfileRepository MySQL 用戶檔案儲存庫實作
type MySQLUserProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository 創建新的 MySQL 用戶檔案儲存庫
func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &MySQLUserProfileRepository{db: db}
}

// Create 創建用戶檔案
func (r *MySQLUserProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("創建用戶檔案失敗: %w", err)
	}
	return nil
}

// GetByUserID 根據用戶 ID 獲取檔案
func (r *MySQLUserProfileRepository) GetByUserID(ctx context.Context, userID uint) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用戶檔案不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢用戶檔案失敗: %w", err)
	}
	return &profile, nil
}

// Update 更新用戶檔案
func (r *MySQLUserProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return fmt.Errorf("更新用戶檔案失敗: %w", err)
	}
	return nil
}

// Delete 刪除用戶檔案
func (r *MySQLUserProfileRepository) Delete(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.UserProfile{}).Error; err != nil {
		return fmt.Errorf("刪除用戶檔案失敗: %w", err)
	}
	return nil
}

// UpdateLocation 更新用戶位置資訊
func (r *MySQLUserProfileRepository) UpdateLocation(ctx context.Context, userID uint, lat, lng *float64) error {
	updates := map[string]interface{}{
		"location_lat": lat,
		"location_lng": lng,
	}

	if err := r.db.WithContext(ctx).Model(&entity.UserProfile{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新位置資訊失敗: %w", err)
	}
	return nil
}

// UpdateMatchingPreferences 更新配對偏好設定
func (r *MySQLUserProfileRepository) UpdateMatchingPreferences(ctx context.Context, userID uint, maxDistance, ageMin, ageMax int) error {
	updates := map[string]interface{}{
		"max_distance":  maxDistance,
		"age_range_min": ageMin,
		"age_range_max": ageMax,
	}

	if err := r.db.WithContext(ctx).Model(&entity.UserProfile{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新配對偏好設定失敗: %w", err)
	}
	return nil
}

// MySQLPhotoRepository MySQL 照片儲存庫實作
type MySQLPhotoRepository struct {
	db *gorm.DB
}

// NewPhotoRepository 創建新的 MySQL 照片儲存庫
func NewPhotoRepository(db *gorm.DB) repository.PhotoRepository {
	return &MySQLPhotoRepository{db: db}
}

// Create 添加用戶照片
func (r *MySQLPhotoRepository) Create(ctx context.Context, photo *entity.Photo) error {
	if err := r.db.WithContext(ctx).Create(photo).Error; err != nil {
		return fmt.Errorf("添加照片失敗: %w", err)
	}
	return nil
}

// GetByUserID 獲取用戶所有照片
func (r *MySQLPhotoRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Photo, error) {
	var photos []*entity.Photo
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("display_order").Find(&photos).Error; err != nil {
		return nil, fmt.Errorf("獲取用戶照片失敗: %w", err)
	}
	return photos, nil
}

// GetByID 根據 ID 獲取照片
func (r *MySQLPhotoRepository) GetByID(ctx context.Context, id uint) (*entity.Photo, error) {
	var photo entity.Photo
	if err := r.db.WithContext(ctx).First(&photo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("照片不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢照片失敗: %w", err)
	}
	return &photo, nil
}

// Update 更新照片資訊
func (r *MySQLPhotoRepository) Update(ctx context.Context, photo *entity.Photo) error {
	if err := r.db.WithContext(ctx).Save(photo).Error; err != nil {
		return fmt.Errorf("更新照片失敗: %w", err)
	}
	return nil
}

// Delete 刪除照片
func (r *MySQLPhotoRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Photo{}, id).Error; err != nil {
		return fmt.Errorf("刪除照片失敗: %w", err)
	}
	return nil
}

// SetPrimary 設定主要照片
func (r *MySQLPhotoRepository) SetPrimary(ctx context.Context, userID, photoID uint) error {
	// 開始事務
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先將該用戶的所有照片設為非主要
		if err := tx.Model(&entity.Photo{}).Where("user_id = ?", userID).Update("is_primary", false).Error; err != nil {
			return fmt.Errorf("重置主要照片失敗: %w", err)
		}

		// 設定指定照片為主要
		if err := tx.Model(&entity.Photo{}).Where("id = ? AND user_id = ?", photoID, userID).Update("is_primary", true).Error; err != nil {
			return fmt.Errorf("設定主要照片失敗: %w", err)
		}

		return nil
	})
}

// UpdateOrder 更新照片排序
func (r *MySQLPhotoRepository) UpdateOrder(ctx context.Context, userID uint, photoOrders []struct {
	PhotoID uint
	Order   int
}) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, order := range photoOrders {
			if err := tx.Model(&entity.Photo{}).
				Where("id = ? AND user_id = ?", order.PhotoID, userID).
				Update("display_order", order.Order).Error; err != nil {
				return fmt.Errorf("更新照片排序失敗 (照片ID: %d): %w", order.PhotoID, err)
			}
		}
		return nil
	})
}

// MySQLInterestRepository MySQL 興趣標籤儲存庫實作
type MySQLInterestRepository struct {
	db *gorm.DB
}

// NewInterestRepository 創建新的 MySQL 興趣標籤儲存庫
func NewInterestRepository(db *gorm.DB) repository.InterestRepository {
	return &MySQLInterestRepository{db: db}
}

// GetAll 獲取所有可用興趣標籤
func (r *MySQLInterestRepository) GetAll(ctx context.Context) ([]*entity.Interest, error) {
	var interests []*entity.Interest
	if err := r.db.WithContext(ctx).Find(&interests).Error; err != nil {
		return nil, fmt.Errorf("獲取興趣標籤失敗: %w", err)
	}
	return interests, nil
}

// GetByUserID 獲取用戶的興趣標籤
func (r *MySQLInterestRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Interest, error) {
	var interests []*entity.Interest
	// 假設有一個 user_interests 關聯表
	if err := r.db.WithContext(ctx).
		Table("interests").
		Select("interests.*").
		Joins("INNER JOIN user_interests ON interests.id = user_interests.interest_id").
		Where("user_interests.user_id = ?", userID).
		Find(&interests).Error; err != nil {
		return nil, fmt.Errorf("獲取用戶興趣標籤失敗: %w", err)
	}
	return interests, nil
}

// SetUserInterests 設定用戶興趣標籤
func (r *MySQLInterestRepository) SetUserInterests(ctx context.Context, userID uint, interestIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先刪除現有關聯
		if err := tx.Exec("DELETE FROM user_interests WHERE user_id = ?", userID).Error; err != nil {
			return fmt.Errorf("刪除現有興趣關聯失敗: %w", err)
		}

		// 新增新的關聯
		for _, interestID := range interestIDs {
			if err := tx.Exec("INSERT INTO user_interests (user_id, interest_id) VALUES (?, ?)", userID, interestID).Error; err != nil {
				return fmt.Errorf("新增興趣關聯失敗 (興趣ID: %d): %w", interestID, err)
			}
		}

		return nil
	})
}

// Create 創建新興趣標籤
func (r *MySQLInterestRepository) Create(ctx context.Context, interest *entity.Interest) error {
	if err := r.db.WithContext(ctx).Create(interest).Error; err != nil {
		return fmt.Errorf("創建興趣標籤失敗: %w", err)
	}
	return nil
}

// Update 更新興趣標籤
func (r *MySQLInterestRepository) Update(ctx context.Context, interest *entity.Interest) error {
	if err := r.db.WithContext(ctx).Save(interest).Error; err != nil {
		return fmt.Errorf("更新興趣標籤失敗: %w", err)
	}
	return nil
}

// Delete 刪除興趣標籤
func (r *MySQLInterestRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Interest{}, id).Error; err != nil {
		return fmt.Errorf("刪除興趣標籤失敗: %w", err)
	}
	return nil
}

// MySQLAgeVerificationRepository MySQL 年齡驗證儲存庫實作
type MySQLAgeVerificationRepository struct {
	db *gorm.DB
}

// NewAgeVerificationRepository 創建新的 MySQL 年齡驗證儲存庫
func NewAgeVerificationRepository(db *gorm.DB) repository.AgeVerificationRepository {
	return &MySQLAgeVerificationRepository{db: db}
}

// Create 創建年齡驗證記錄
func (r *MySQLAgeVerificationRepository) Create(ctx context.Context, verification *entity.AgeVerification) error {
	if err := r.db.WithContext(ctx).Create(verification).Error; err != nil {
		return fmt.Errorf("創建年齡驗證記錄失敗: %w", err)
	}
	return nil
}

// GetByUserID 根據用戶 ID 獲取驗證記錄
func (r *MySQLAgeVerificationRepository) GetByUserID(ctx context.Context, userID uint) (*entity.AgeVerification, error) {
	var verification entity.AgeVerification
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("年齡驗證記錄不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢年齡驗證記錄失敗: %w", err)
	}
	return &verification, nil
}

// Update 更新驗證記錄
func (r *MySQLAgeVerificationRepository) Update(ctx context.Context, verification *entity.AgeVerification) error {
	if err := r.db.WithContext(ctx).Save(verification).Error; err != nil {
		return fmt.Errorf("更新年齡驗證記錄失敗: %w", err)
	}
	return nil
}

// GetPendingVerifications 獲取待審核的驗證記錄
func (r *MySQLAgeVerificationRepository) GetPendingVerifications(ctx context.Context, limit int) ([]*entity.AgeVerification, error) {
	var verifications []*entity.AgeVerification
	query := r.db.WithContext(ctx).Where("status = ?", entity.VerificationStatusPending).Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&verifications).Error; err != nil {
		return nil, fmt.Errorf("獲取待審核驗證記錄失敗: %w", err)
	}
	return verifications, nil
}

// SetVerificationStatus 設定驗證狀態
func (r *MySQLAgeVerificationRepository) SetVerificationStatus(ctx context.Context, userID uint, status entity.VerificationStatus, reviewerID *uint, notes string) error {
	updates := map[string]interface{}{
		"status":       status,
		"review_notes": notes,
	}

	if reviewerID != nil {
		updates["reviewer_id"] = *reviewerID
	}

	if err := r.db.WithContext(ctx).Model(&entity.AgeVerification{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("設定驗證狀態失敗: %w", err)
	}
	return nil
}
