package entity

import (
	"time"
)

// Gender 性別枚舉
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// UserStatus 用戶狀態枚舉
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

// StringArray 純粹的領域類型，用於處理字串陣列
// 不包含任何基礎設施相關的邏輯
type StringArray []string

// IsEmpty 檢查陣列是否為空
func (sa StringArray) IsEmpty() bool {
	return len(sa) == 0
}

// Contains 檢查是否包含特定元素
func (sa StringArray) Contains(item string) bool {
	for _, s := range sa {
		if s == item {
			return true
		}
	}
	return false
}

// Add 添加元素（返回新的陣列，保持不可變性）
func (sa StringArray) Add(item string) StringArray {
	if sa.Contains(item) {
		return sa
	}
	newArray := make([]string, len(sa)+1)
	copy(newArray, sa)
	newArray[len(sa)] = item
	return StringArray(newArray)
}

// Remove 移除元素（返回新的陣列，保持不可變性）
func (sa StringArray) Remove(item string) StringArray {
	var result []string
	for _, s := range sa {
		if s != item {
			result = append(result, s)
		}
	}
	return StringArray(result)
}

// UserInformation 完整的用戶資訊結構體
type UserInformation struct {
	ID            int         `json:"id"`
	Username      string      `json:"username"`
	Email         string      `json:"email"`
	Password      string      `json:"-"` // 不序列化密碼
	Age           *int        `json:"age,omitempty"`
	Gender        *Gender     `json:"gender,omitempty"`
	Bio           *string     `json:"bio,omitempty"`
	Interests     StringArray `json:"interests,omitempty"`
	LocationLat   *float64    `json:"location_lat,omitempty"`
	LocationLng   *float64    `json:"location_lng,omitempty"`
	City          *string     `json:"city,omitempty"`
	Country       *string     `json:"country,omitempty"`
	IsVerified    bool        `json:"is_verified"`
	Status        UserStatus  `json:"status"`
	LastActiveAt  *time.Time  `json:"last_active_at,omitempty"`
	ProfileViews  int         `json:"profile_views"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// UserProfile 詳細的用戶個人檔案資訊
type UserProfile struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	Height          *int        `json:"height,omitempty"`        // 身高 (cm)
	Weight          *int        `json:"weight,omitempty"`        // 體重 (kg)
	Education       *string     `json:"education,omitempty"`  // 教育背景
	Occupation      *string     `json:"occupation,omitempty"` // 職業
	Company         *string     `json:"company,omitempty"`      // 公司
	Relationship    *string     `json:"relationship,omitempty"` // 感情狀態
	LookingFor      StringArray `json:"looking_for,omitempty"`    // 尋找什麼關係
	Languages       StringArray `json:"languages,omitempty"`       // 語言能力
	Hobbies         StringArray `json:"hobbies,omitempty"`           // 興趣愛好
	Lifestyle       StringArray `json:"lifestyle,omitempty"`       // 生活方式
	PetPreference   *string     `json:"pet_preference,omitempty"` // 寵物偏好
	DrinkingHabit   *string     `json:"drinking_habit,omitempty"` // 飲酒習慣
	SmokingHabit    *string     `json:"smoking_habit,omitempty"`   // 吸菸習慣
	ExerciseHabit   *string     `json:"exercise_habit,omitempty"` // 運動習慣
	SocialMediaLink *string     `json:"social_media_link,omitempty"` // 社群媒體連結
	PersonalityType *string     `json:"personality_type,omitempty"`   // 人格類型
	Zodiac          *string     `json:"zodiac,omitempty"`                      // 星座
	Religion        *string     `json:"religion,omitempty"`                  // 宗教信仰
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// PhotoStatus 照片狀態枚舉
type PhotoStatus string

const (
	PhotoStatusPending  PhotoStatus = "pending"  // 等待審核
	PhotoStatusApproved PhotoStatus = "approved" // 已審核通過
	PhotoStatusRejected PhotoStatus = "rejected" // 審核不通過
)

// UserPhoto 用戶照片管理
type UserPhoto struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	PhotoURL    string      `json:"photo_url"`
	ThumbnailURL *string    `json:"thumbnail_url,omitempty"`
	IsPrimary   bool        `json:"is_primary"`
	Order       int         `json:"order"`        // 照片排序
	Status      PhotoStatus `json:"status"`
	Caption     *string     `json:"caption,omitempty"` // 照片說明
	IsVerified  bool        `json:"is_verified"`   // 是否為認證照片
	UploadedAt  time.Time   `json:"uploaded_at"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// UserPreference 用戶配對偏好設定
type UserPreference struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`
	PreferredGender *Gender    `json:"preferred_gender,omitempty"`
	AgeMin         *int        `json:"age_min,omitempty"`
	AgeMax         *int        `json:"age_max,omitempty"`
	DistanceMax    *int        `json:"distance_max,omitempty"` // 最大距離 (km)
	HeightMin      *int        `json:"height_min,omitempty"`
	HeightMax      *int        `json:"height_max,omitempty"`
	Education      StringArray `json:"education,omitempty"`
	Interests      StringArray `json:"interests,omitempty"`
	Lifestyle      StringArray `json:"lifestyle,omitempty"`
	ShowMe         bool        `json:"show_me"`               // 是否顯示我的資料給別人
	ShowDistance   bool        `json:"show_distance"`   // 是否顯示距離
	ShowAge        bool        `json:"show_age"`             // 是否顯示年齡
	ShowLastActive bool        `json:"show_last_active"` // 是否顯示最後上線時間
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
