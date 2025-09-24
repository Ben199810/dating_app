package entity

import (
	"errors"
	"strings"
	"time"
)

// Gender 性別枚舉
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// IsValid 檢查性別是否有效
func (g Gender) IsValid() bool {
	return g == GenderMale || g == GenderFemale || g == GenderOther
}

// UserProfile 用戶檔案實體
type UserProfile struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	DisplayName string    `gorm:"not null;size:50" json:"display_name"`
	Bio         string    `gorm:"size:500" json:"bio"`
	Gender      Gender    `gorm:"not null" json:"gender"`
	ShowAge     bool      `gorm:"default:true" json:"show_age"`
	LocationLat *float64  `json:"location_lat"`
	LocationLng *float64  `json:"location_lng"`
	MaxDistance int       `gorm:"default:50" json:"max_distance"` // km
	AgeRangeMin int       `gorm:"default:18" json:"age_range_min"`
	AgeRangeMax int       `gorm:"default:99" json:"age_range_max"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 關聯 - 將在User實體完成後添加
	// User User `gorm:"constraint:OnDelete:CASCADE" json:"user"`
}

// Validate 驗證用戶檔案資料
func (up *UserProfile) Validate() error {
	// 檢查必填欄位
	if up.UserID == 0 {
		return errors.New("user_id 是必填欄位")
	}

	if strings.TrimSpace(up.DisplayName) == "" {
		return errors.New("display_name 是必填欄位")
	}

	// 檢查顯示名稱長度
	if len(up.DisplayName) > 50 {
		return errors.New("display_name 不能超過 50 字元")
	}

	// 檢查個人簡介長度
	if len(up.Bio) > 500 {
		return errors.New("bio 不能超過 500 字元")
	}

	// 檢查性別有效性
	if !up.Gender.IsValid() {
		return errors.New("gender 必須是 male、female 或 other")
	}

	// 檢查地理座標範圍
	if up.LocationLat != nil && (*up.LocationLat < -90 || *up.LocationLat > 90) {
		return errors.New("location_lat 必須在 -90 到 90 之間")
	}

	if up.LocationLng != nil && (*up.LocationLng < -180 || *up.LocationLng > 180) {
		return errors.New("location_lng 必須在 -180 到 180 之間")
	}

	// 檢查距離範圍
	if up.MaxDistance < 1 || up.MaxDistance > 1000 {
		return errors.New("max_distance 必須在 1 到 1000 公里之間")
	}

	// 檢查年齡範圍
	if up.AgeRangeMin < 18 || up.AgeRangeMin > 99 {
		return errors.New("age_range_min 必須在 18 到 99 之間")
	}

	if up.AgeRangeMax < 18 || up.AgeRangeMax > 99 {
		return errors.New("age_range_max 必須在 18 到 99 之間")
	}

	if up.AgeRangeMin > up.AgeRangeMax {
		return errors.New("age_range_min 不能大於 age_range_max")
	}

	return nil
}

// HasLocation 檢查是否設定了位置資訊
func (up *UserProfile) HasLocation() bool {
	return up.LocationLat != nil && up.LocationLng != nil
}

// SetLocation 設定位置座標
func (up *UserProfile) SetLocation(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("緯度必須在 -90 到 90 之間")
	}

	if lng < -180 || lng > 180 {
		return errors.New("經度必須在 -180 到 180 之間")
	}

	up.LocationLat = &lat
	up.LocationLng = &lng
	up.UpdatedAt = time.Now()

	return nil
}

// ClearLocation 清除位置資訊
func (up *UserProfile) ClearLocation() {
	up.LocationLat = nil
	up.LocationLng = nil
	up.UpdatedAt = time.Now()
}

// UpdateAgeRange 更新年齡範圍偏好
func (up *UserProfile) UpdateAgeRange(min, max int) error {
	if min < 18 || min > 99 {
		return errors.New("最小年齡必須在 18 到 99 之間")
	}

	if max < 18 || max > 99 {
		return errors.New("最大年齡必須在 18 到 99 之間")
	}

	if min > max {
		return errors.New("最小年齡不能大於最大年齡")
	}

	up.AgeRangeMin = min
	up.AgeRangeMax = max
	up.UpdatedAt = time.Now()

	return nil
}

// UpdateMaxDistance 更新最大配對距離
func (up *UserProfile) UpdateMaxDistance(distance int) error {
	if distance < 1 || distance > 1000 {
		return errors.New("配對距離必須在 1 到 1000 公里之間")
	}

	up.MaxDistance = distance
	up.UpdatedAt = time.Now()

	return nil
}
