package entity

import (
	"errors"
	"strings"
	"time"
)

// InterestCategory 興趣類別枚舉
type InterestCategory string

const (
	InterestCategoryHobbies    InterestCategory = "hobbies"    // 愛好
	InterestCategorySports     InterestCategory = "sports"     // 運動
	InterestCategoryMusic      InterestCategory = "music"      // 音樂
	InterestCategoryMovies     InterestCategory = "movies"     // 電影
	InterestCategoryFood       InterestCategory = "food"       // 美食
	InterestCategoryTravel     InterestCategory = "travel"     // 旅行
	InterestCategoryReading    InterestCategory = "reading"    // 閱讀
	InterestCategoryTechnology InterestCategory = "technology" // 科技
	InterestCategoryArt        InterestCategory = "art"        // 藝術
	InterestCategoryFitness    InterestCategory = "fitness"    // 健身
	InterestCategoryGaming     InterestCategory = "gaming"     // 遊戲
	InterestCategoryNature     InterestCategory = "nature"     // 自然
	InterestCategoryOther      InterestCategory = "other"      // 其他
)

// IsValid 檢查興趣類別是否有效
func (ic InterestCategory) IsValid() bool {
	validCategories := []InterestCategory{
		InterestCategoryHobbies,
		InterestCategorySports,
		InterestCategoryMusic,
		InterestCategoryMovies,
		InterestCategoryFood,
		InterestCategoryTravel,
		InterestCategoryReading,
		InterestCategoryTechnology,
		InterestCategoryArt,
		InterestCategoryFitness,
		InterestCategoryGaming,
		InterestCategoryNature,
		InterestCategoryOther,
	}

	for _, category := range validCategories {
		if ic == category {
			return true
		}
	}
	return false
}

// GetDisplayName 獲取興趣類別的顯示名稱
func (ic InterestCategory) GetDisplayName() string {
	switch ic {
	case InterestCategoryHobbies:
		return "愛好"
	case InterestCategorySports:
		return "運動"
	case InterestCategoryMusic:
		return "音樂"
	case InterestCategoryMovies:
		return "電影"
	case InterestCategoryFood:
		return "美食"
	case InterestCategoryTravel:
		return "旅行"
	case InterestCategoryReading:
		return "閱讀"
	case InterestCategoryTechnology:
		return "科技"
	case InterestCategoryArt:
		return "藝術"
	case InterestCategoryFitness:
		return "健身"
	case InterestCategoryGaming:
		return "遊戲"
	case InterestCategoryNature:
		return "自然"
	case InterestCategoryOther:
		return "其他"
	default:
		return string(ic)
	}
}

// Interest 興趣標籤實體
type Interest struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Name        string           `gorm:"not null;uniqueIndex;size:100" json:"name"`
	Category    InterestCategory `gorm:"not null" json:"category"`
	Description *string          `gorm:"size:300" json:"description,omitempty"`
	IsActive    bool             `gorm:"default:true" json:"is_active"`
	UsageCount  int              `gorm:"default:0" json:"usage_count"` // 使用此興趣的用戶數量

	// 時間戳記
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 關聯 - 多對多關聯將在UserProfile實體完成後添加
	// Users []UserProfile `gorm:"many2many:user_interests;" json:"users,omitempty"`
}

// Validate 驗證興趣標籤資料
func (i *Interest) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("name 是必填欄位")
	}

	if len(i.Name) > 100 {
		return errors.New("name 不能超過 100 字元")
	}

	if !i.Category.IsValid() {
		return errors.New("category 必須是有效的興趣類別")
	}

	// 檢查描述長度（如果有提供）
	if i.Description != nil && len(*i.Description) > 300 {
		return errors.New("description 不能超過 300 字元")
	}

	// 檢查使用數量不能為負數
	if i.UsageCount < 0 {
		return errors.New("usage_count 不能為負數")
	}

	return nil
}

// IsActiveInterest 檢查興趣是否啟用
func (i *Interest) IsActiveInterest() bool {
	return i.IsActive
}

// Activate 啟用興趣標籤
func (i *Interest) Activate() {
	i.IsActive = true
	i.UpdatedAt = time.Now()
}

// Deactivate 停用興趣標籤
func (i *Interest) Deactivate() {
	i.IsActive = false
	i.UpdatedAt = time.Now()
}

// IncrementUsage 增加使用次數
func (i *Interest) IncrementUsage() {
	i.UsageCount++
	i.UpdatedAt = time.Now()
}

// DecrementUsage 減少使用次數
func (i *Interest) DecrementUsage() {
	if i.UsageCount > 0 {
		i.UsageCount--
		i.UpdatedAt = time.Now()
	}
}

// SetDescription 設定描述
func (i *Interest) SetDescription(description string) {
	cleanDescription := strings.TrimSpace(description)
	if cleanDescription == "" {
		i.Description = nil
	} else {
		i.Description = &cleanDescription
	}
	i.UpdatedAt = time.Now()
}

// HasDescription 檢查是否有描述
func (i *Interest) HasDescription() bool {
	return i.Description != nil && strings.TrimSpace(*i.Description) != ""
}

// GetDescription 獲取描述（如果有）
func (i *Interest) GetDescription() string {
	if i.Description != nil {
		return *i.Description
	}
	return ""
}

// GetCategoryDisplayName 獲取類別的顯示名稱
func (i *Interest) GetCategoryDisplayName() string {
	return i.Category.GetDisplayName()
}

// IsPopular 檢查是否為熱門興趣（使用數量 >= 10）
func (i *Interest) IsPopular() bool {
	return i.UsageCount >= 10
}

// GetPopularityLevel 獲取熱門程度等級
func (i *Interest) GetPopularityLevel() string {
	switch {
	case i.UsageCount >= 100:
		return "非常熱門"
	case i.UsageCount >= 50:
		return "很熱門"
	case i.UsageCount >= 10:
		return "熱門"
	case i.UsageCount >= 1:
		return "普通"
	default:
		return "新興"
	}
}
