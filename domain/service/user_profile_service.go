package service

import (
	"errors"
	"golang_dev_docker/component/validator"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"strconv"
)

// UpdateUserProfileRequest 更新用戶個人資料請求
type UpdateUserProfileRequest struct {
	Username string          `json:"username"`
	Age      *int            `json:"age"`
	Gender   *entity.Gender  `json:"gender"`
	Country  string          `json:"country"`
	City     string          `json:"city"`
	Bio      *string         `json:"bio"`
}

type UserProfileService struct {
	userRepo           repository.UserRepository
	userProfileRepo    repository.UserProfileRepository
	userPhotoRepo      repository.UserPhotoRepository
	userPreferenceRepo repository.UserPreferenceRepository
}

func NewUserProfileService(
	userRepo repository.UserRepository,
	profileRepo repository.UserProfileRepository,
	photoRepo repository.UserPhotoRepository,
	preferenceRepo repository.UserPreferenceRepository,
) *UserProfileService {
	return &UserProfileService{
		userRepo:           userRepo,
		userProfileRepo:    profileRepo,
		userPhotoRepo:      photoRepo,
		userPreferenceRepo: preferenceRepo,
	}
}

// UpdateUserBasicInfo 更新用戶基本資訊
func (s *UserProfileService) UpdateUserBasicInfo(userID int, age *int, gender *entity.Gender, bio *string, interests []string) error {
	// 驗證輸入
	if err := validator.ValidateAge(age); err != nil {
		return err
	}

	genderStr := ""
	if gender != nil {
		genderStr = string(*gender)
	}
	if err := validator.ValidateGender(&genderStr); err != nil {
		return err
	}

	if err := validator.ValidateBio(bio); err != nil {
		return err
	}

	if err := validator.ValidateInterests(interests); err != nil {
		return err
	}

	// 取得用戶現有資訊
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用戶不存在")
	}

	// 更新資訊
	user.Age = age
	user.Gender = gender
	user.Bio = bio
	user.Interests = entity.StringArray(interests)

	return s.userRepo.Update(user)
}

// UpdateUserLocation 更新用戶位置
func (s *UserProfileService) UpdateUserLocation(userID int, lat, lng *float64, city, country *string) error {
	// 驗證位置
	if err := validator.ValidateLocation(lat, lng); err != nil {
		return err
	}

	// 取得用戶現有資訊
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用戶不存在")
	}

	// 更新位置資訊
	user.LocationLat = lat
	user.LocationLng = lng
	user.City = city
	user.Country = country

	return s.userRepo.Update(user)
}

// AddUserPhoto 新增用戶照片
func (s *UserProfileService) AddUserPhoto(userID int, photoURL string, caption *string, isPrimary bool) error {
	// 驗證照片URL
	if err := validator.ValidatePhotoURL(photoURL); err != nil {
		return err
	}

	// 如果設為主要照片，需要將其他照片設為非主要
	if isPrimary {
		// 取得現有主要照片
		primaryPhoto, err := s.userPhotoRepo.GetPrimaryPhoto(userID)
		if err != nil {
			return err
		}

		// 如果有現有主要照片，將其設為非主要
		if primaryPhoto != nil {
			primaryPhoto.IsPrimary = false
			if err := s.userPhotoRepo.UpdatePhoto(primaryPhoto); err != nil {
				return err
			}
		}
	}

	// 創建新照片
	photo := &entity.UserPhoto{
		UserID:    userID,
		PhotoURL:  photoURL,
		IsPrimary: isPrimary,
		Caption:   caption,
		Status:    entity.PhotoStatusPending,
		Order:     0, // 可以根據現有照片數量設定
	}

	return s.userPhotoRepo.CreatePhoto(photo)
}

// CreateUserProfile 創建詳細用戶資料
func (s *UserProfileService) CreateUserProfile(userID int, profile *entity.UserProfile) error {
	// 驗證身高體重
	if err := validator.ValidateHeight(profile.Height); err != nil {
		return err
	}
	if err := validator.ValidateWeight(profile.Weight); err != nil {
		return err
	}

	profile.UserID = userID
	return s.userProfileRepo.CreateProfile(profile)
}

// FindNearbyUsers 尋找附近用戶
func (s *UserProfileService) FindNearbyUsers(userID int, radiusKm int, limit int) ([]*entity.UserInformation, error) {
	// 取得當前用戶位置
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用戶不存在")
	}

	if user.LocationLat == nil || user.LocationLng == nil {
		return nil, errors.New("用戶未設定位置")
	}

	// 更新用戶最後活躍時間
	s.userRepo.UpdateLastActiveTime(userID)

	// 搜尋附近用戶
	return s.userRepo.GetUsersByLocation(*user.LocationLat, *user.LocationLng, radiusKm, limit)
}

// SearchCompatibleUsers 搜尋相容的用戶
func (s *UserProfileService) SearchCompatibleUsers(userID int, limit int) ([]*entity.UserInformation, error) {
	// 取得用戶偏好設定
	preference, err := s.userPreferenceRepo.GetPreferenceByUserID(userID)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]interface{})

	if preference != nil {
		if preference.PreferredGender != nil {
			filters["gender"] = *preference.PreferredGender
		}
		if preference.AgeMin != nil {
			filters["age_min"] = *preference.AgeMin
		}
		if preference.AgeMax != nil {
			filters["age_max"] = *preference.AgeMax
		}
	}

	// 排除自己
	// 這裡可以添加更複雜的過濾邏輯

	return s.userRepo.SearchUsers(filters, limit, 0)
}

// GetUserProfile 獲取用戶個人資料
func (s *UserProfileService) GetUserProfile(userIDStr string) (*entity.UserInformation, error) {
	// 轉換用戶 ID
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, errors.New("無效的用戶 ID")
	}

	// 獲取用戶基本資訊
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用戶不存在")
	}

	// 清空密碼欄位
	user.Password = ""

	return user, nil
}

// UpdateUserProfile 更新用戶個人資料
func (s *UserProfileService) UpdateUserProfile(userIDStr string, req *UpdateUserProfileRequest) error {
	// 轉換用戶 ID
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return errors.New("無效的用戶 ID")
	}

	// 驗證用戶是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用戶不存在")
	}

	// 更新用戶基本資訊
	if req.Username != "" {
		if err := validator.ValidateUsername(req.Username); err != nil {
			return err
		}
		user.Username = req.Username
	}

	if req.Country != "" {
		user.Country = &req.Country
	}

	if req.City != "" {
		user.City = &req.City
	}

	// 更新額外的個人資料欄位
	if req.Age != nil || req.Gender != nil || req.Bio != nil {
		if err := s.UpdateUserBasicInfo(userID, req.Age, req.Gender, req.Bio, nil); err != nil {
			return err
		}
	}

	// 更新用戶基本資訊
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}