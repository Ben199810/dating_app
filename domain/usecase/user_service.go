package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

// UserService 用戶業務邏輯服務
// 負責用戶註冊、認證、個人檔案管理等核心業務邏輯
type UserService struct {
	userRepo            repository.UserRepository
	userProfileRepo     repository.UserProfileRepository
	photoRepo           repository.PhotoRepository
	interestRepo        repository.InterestRepository
	ageVerificationRepo repository.AgeVerificationRepository
}

// NewUserService 創建新的用戶服務實例
func NewUserService(
	userRepo repository.UserRepository,
	userProfileRepo repository.UserProfileRepository,
	photoRepo repository.PhotoRepository,
	interestRepo repository.InterestRepository,
	ageVerificationRepo repository.AgeVerificationRepository,
) *UserService {
	return &UserService{
		userRepo:            userRepo,
		userProfileRepo:     userProfileRepo,
		photoRepo:           photoRepo,
		interestRepo:        interestRepo,
		ageVerificationRepo: ageVerificationRepo,
	}
}

// RegisterRequest 用戶註冊請求
type RegisterRequest struct {
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,min=8"`
	BirthDate   time.Time `json:"birth_date" validate:"required"`
	DisplayName string    `json:"display_name" validate:"required,min=2,max=50"`
	Gender      string    `json:"gender" validate:"required"`
	Biography   string    `json:"biography" validate:"max=500"`
}

// LoginRequest 用戶登入請求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse 用戶回應資料
type UserResponse struct {
	ID         uint                `json:"id"`
	Email      string              `json:"email"`
	IsVerified bool                `json:"is_verified"`
	IsActive   bool                `json:"is_active"`
	CreatedAt  time.Time           `json:"created_at"`
	Age        int                 `json:"age"`
	Profile    *entity.UserProfile `json:"profile,omitempty"`
}

// Register 用戶註冊
// 處理用戶註冊流程，包括年齡驗證、密碼加密、創建基本檔案
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// 驗證請求資料
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, fmt.Errorf("註冊資料驗證失敗: %w", err)
	}

	// 檢查年齡（必須18+）
	user := &entity.User{
		Email:     strings.ToLower(req.Email),
		BirthDate: req.BirthDate,
	}
	if !user.IsAdult() {
		return nil, errors.New("用戶必須年滿18歲才能註冊")
	}

	// 檢查 Email 是否已存在
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("此 Email 已被註冊")
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密碼加密失敗: %w", err)
	}
	user.PasswordHash = string(hashedPassword)
	user.IsVerified = false // 需要年齡驗證
	user.IsActive = true

	// 創建用戶
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("創建用戶失敗: %w", err)
	}

	// 創建基本用戶檔案
	profile := &entity.UserProfile{
		UserID:      user.ID,
		DisplayName: req.DisplayName,
		Bio:         req.Biography,
		Gender:      entity.Gender(req.Gender),
	}

	if err := s.userProfileRepo.Create(ctx, profile); err != nil {
		// 如果檔案創建失敗，回滾用戶創建（簡化處理）
		s.userRepo.Delete(ctx, user.ID)
		return nil, fmt.Errorf("創建用戶檔案失敗: %w", err)
	}

	// 創建年齡驗證記錄
	ageVerification := &entity.AgeVerification{
		UserID:            user.ID,
		Method:            entity.VerificationMethodOther, // 預設，等待用戶選擇
		DocumentNumber:    "",                             // 等待用戶填寫
		DocumentImagePath: "",                             // 等待用戶上傳
		Status:            entity.VerificationStatusPending,
	}

	if err := s.ageVerificationRepo.Create(ctx, ageVerification); err != nil {
		// 記錄錯誤但不阻塞註冊流程
		// 可以稍後通過其他方式創建驗證記錄
	}

	return &UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Age:        user.GetAge(),
		Profile:    profile,
	}, nil
}

// Login 用戶登入
// 處理用戶登入驗證，返回用戶資訊
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*UserResponse, error) {
	// 驗證請求資料
	if err := s.validateLoginRequest(req); err != nil {
		return nil, fmt.Errorf("登入資料驗證失敗: %w", err)
	}

	// 根據 Email 查找用戶
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, errors.New("Email 或密碼錯誤")
	}

	// 檢查用戶狀態
	if !user.IsActive {
		return nil, errors.New("帳戶已被停用")
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("Email 或密碼錯誤")
	}

	// 獲取用戶檔案
	profile, err := s.userProfileRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		// 檔案不存在不影響登入，設為 nil
		profile = nil
	}

	return &UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Age:        user.GetAge(),
		Profile:    profile,
	}, nil
}

// GetProfile 獲取用戶完整檔案資訊
func (s *UserService) GetProfile(ctx context.Context, userID uint) (*UserResponse, error) {
	// 獲取基本用戶資訊
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("用戶不存在: %w", err)
	}

	// 獲取用戶檔案
	profile, err := s.userProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("獲取用戶檔案失敗: %w", err)
	}

	return &UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Age:        user.GetAge(),
		Profile:    profile,
	}, nil
}

// UpdateProfile 更新用戶檔案
type UpdateProfileRequest struct {
	DisplayName      *string        `json:"display_name,omitempty"`
	Biography        *string        `json:"biography,omitempty"`
	LocationLat      *float64       `json:"location_lat,omitempty"`
	LocationLng      *float64       `json:"location_lng,omitempty"`
	MaxDistance      *int           `json:"max_distance,omitempty"`
	AgeRangeMin      *int           `json:"age_range_min,omitempty"`
	AgeRangeMax      *int           `json:"age_range_max,omitempty"`
	InterestedGender *entity.Gender `json:"interested_gender,omitempty"`
	InterestIDs      []uint         `json:"interest_ids,omitempty"`
}

// UpdateProfile 更新用戶檔案
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*UserResponse, error) {
	// 獲取現有檔案
	profile, err := s.userProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("獲取用戶檔案失敗: %w", err)
	}

	// 更新檔案資料
	if req.DisplayName != nil {
		if strings.TrimSpace(*req.DisplayName) == "" {
			return nil, errors.New("顯示名稱不能為空")
		}
		profile.DisplayName = *req.DisplayName
	}

	if req.Biography != nil {
		profile.Bio = *req.Biography
	}

	// Note: InterestedGender not found in UserProfile entity - removing this feature
	// if req.InterestedGender != nil {
	//	profile.InterestedGender = *req.InterestedGender
	// }

	// 更新檔案
	if err := s.userProfileRepo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("更新檔案失敗: %w", err)
	}

	// 更新位置資訊
	if req.LocationLat != nil && req.LocationLng != nil {
		if err := s.userProfileRepo.UpdateLocation(ctx, userID, req.LocationLat, req.LocationLng); err != nil {
			return nil, fmt.Errorf("更新位置失敗: %w", err)
		}
	}

	// 更新配對偏好
	if req.MaxDistance != nil || req.AgeRangeMin != nil || req.AgeRangeMax != nil {
		maxDistance := profile.MaxDistance
		ageMin := profile.AgeRangeMin
		ageMax := profile.AgeRangeMax

		if req.MaxDistance != nil {
			maxDistance = *req.MaxDistance
		}
		if req.AgeRangeMin != nil {
			ageMin = *req.AgeRangeMin
		}
		if req.AgeRangeMax != nil {
			ageMax = *req.AgeRangeMax
		}

		if err := s.userProfileRepo.UpdateMatchingPreferences(ctx, userID, maxDistance, ageMin, ageMax); err != nil {
			return nil, fmt.Errorf("更新配對偏好失敗: %w", err)
		}
	}

	// 更新興趣標籤
	if req.InterestIDs != nil {
		if err := s.interestRepo.SetUserInterests(ctx, userID, req.InterestIDs); err != nil {
			return nil, fmt.Errorf("更新興趣標籤失敗: %w", err)
		}
	}

	// 返回更新後的檔案
	return s.GetProfile(ctx, userID)
}

// GetUserPhotos 獲取用戶照片
func (s *UserService) GetUserPhotos(ctx context.Context, userID uint) ([]*entity.Photo, error) {
	return s.photoRepo.GetByUserID(ctx, userID)
}

// AddPhoto 添加用戶照片
func (s *UserService) AddPhoto(ctx context.Context, userID uint, imageURL, description string) (*entity.Photo, error) {
	// 檢查用戶照片數量限制（例如最多6張）
	photos, err := s.photoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("獲取用戶照片失敗: %w", err)
	}

	if len(photos) >= 6 {
		return nil, errors.New("照片數量已達上限(6張)")
	}

	// 創建照片記錄
	photo := &entity.Photo{
		UserID:       userID,
		Type:         entity.PhotoTypeProfile,
		FileName:     description,               // Using description as filename for now
		FilePath:     imageURL,                  // Using imageURL as file path
		FileSize:     0,                         // Would need to be calculated
		MimeType:     "image/jpeg",              // Default, should be determined from actual file
		DisplayOrder: len(photos) + 1,           // 新照片排在最後
		Status:       entity.PhotoStatusPending, // 需要審核
	}

	// 如果是第一張照片，設為主要照片
	if len(photos) == 0 {
		photo.IsMain = true
	}

	if err := s.photoRepo.Create(ctx, photo); err != nil {
		return nil, fmt.Errorf("添加照片失敗: %w", err)
	}

	return photo, nil
}

// SetPrimaryPhoto 設定主要照片
func (s *UserService) SetPrimaryPhoto(ctx context.Context, userID, photoID uint) error {
	// 驗證照片所有權
	photo, err := s.photoRepo.GetByID(ctx, photoID)
	if err != nil {
		return fmt.Errorf("照片不存在: %w", err)
	}

	if photo.UserID != userID {
		return errors.New("無權限操作此照片")
	}

	return s.photoRepo.SetPrimary(ctx, userID, photoID)
}

// DeletePhoto 刪除照片
func (s *UserService) DeletePhoto(ctx context.Context, userID, photoID uint) error {
	// 驗證照片所有權
	photo, err := s.photoRepo.GetByID(ctx, photoID)
	if err != nil {
		return fmt.Errorf("照片不存在: %w", err)
	}

	if photo.UserID != userID {
		return errors.New("無權限操作此照片")
	}

	return s.photoRepo.Delete(ctx, photoID)
}

// GetAvailableInterests 獲取所有可用的興趣標籤
func (s *UserService) GetAvailableInterests(ctx context.Context) ([]*entity.Interest, error) {
	return s.interestRepo.GetAll(ctx)
}

// GetUserInterests 獲取用戶的興趣標籤
func (s *UserService) GetUserInterests(ctx context.Context, userID uint) ([]*entity.Interest, error) {
	return s.interestRepo.GetByUserID(ctx, userID)
}

// SubmitAgeVerification 提交年齡驗證
func (s *UserService) SubmitAgeVerification(ctx context.Context, userID uint, method entity.VerificationMethod, documentNumber string, documentImagePath string) error {
	// 檢查是否已有驗證記錄
	verification, err := s.ageVerificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		// 沒有記錄，創建新的
		verification = &entity.AgeVerification{
			UserID:            userID,
			Method:            method,
			DocumentNumber:    documentNumber,
			DocumentImagePath: documentImagePath,
			Status:            entity.VerificationStatusPending,
		}
		return s.ageVerificationRepo.Create(ctx, verification)
	}

	// 更新現有記錄
	verification.Method = method
	verification.DocumentNumber = documentNumber
	verification.DocumentImagePath = documentImagePath
	verification.Status = entity.VerificationStatusPending

	return s.ageVerificationRepo.Update(ctx, verification)
}

// GetAgeVerificationStatus 獲取年齡驗證狀態
func (s *UserService) GetAgeVerificationStatus(ctx context.Context, userID uint) (*entity.AgeVerification, error) {
	return s.ageVerificationRepo.GetByUserID(ctx, userID)
}

// 私有輔助方法

// validateRegisterRequest 驗證註冊請求
func (s *UserService) validateRegisterRequest(req *RegisterRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("Email 不能為空")
	}

	if !strings.Contains(req.Email, "@") {
		return errors.New("Email 格式不正確")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("密碼不能為空")
	}

	if len(req.Password) < 8 {
		return errors.New("密碼長度至少8個字符")
	}

	if strings.TrimSpace(req.DisplayName) == "" {
		return errors.New("顯示名稱不能為空")
	}

	if len(req.DisplayName) < 2 || len(req.DisplayName) > 50 {
		return errors.New("顯示名稱長度必須在2-50個字符之間")
	}

	if strings.TrimSpace(req.Gender) == "" {
		return errors.New("性別不能為空")
	}

	// 驗證性別值
	gender := entity.Gender(req.Gender)
	if !gender.IsValid() {
		return errors.New("性別值不正確")
	}

	if req.BirthDate.IsZero() {
		return errors.New("出生日期不能為空")
	}

	if len(req.Biography) > 500 {
		return errors.New("個人簡介不能超過500個字符")
	}

	return nil
}

// validateLoginRequest 驗證登入請求
func (s *UserService) validateLoginRequest(req *LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("Email 不能為空")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("密碼不能為空")
	}

	return nil
}
