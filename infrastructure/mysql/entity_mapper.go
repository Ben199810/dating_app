package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"golang_dev_docker/domain/entity"
)

// StringArrayMapper 處理 StringArray 與資料庫 JSON 欄位之間的轉換
type StringArrayMapper struct{}

// ToDatabase 將 StringArray 轉換為資料庫可儲存的格式
func (m *StringArrayMapper) ToDatabase(sa entity.StringArray) (driver.Value, error) {
	if len(sa) == 0 {
		return "[]", nil
	}
	return json.Marshal([]string(sa))
}

// FromDatabase 將資料庫中的 JSON 資料轉換為 StringArray
func (m *StringArrayMapper) FromDatabase(value interface{}) (entity.StringArray, error) {
	if value == nil {
		return entity.StringArray{}, nil
	}

	var result []string

	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, &result); err != nil {
			return nil, err
		}
	case string:
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return nil, err
		}
	default:
		return entity.StringArray{}, nil
	}

	return entity.StringArray(result), nil
}

// UserEntityMapper 處理 UserInformation Entity 與資料庫記錄之間的轉換
type UserEntityMapper struct {
	stringArrayMapper *StringArrayMapper
}

func NewUserEntityMapper() *UserEntityMapper {
	return &UserEntityMapper{
		stringArrayMapper: &StringArrayMapper{},
	}
}

// UserDBModel 資料庫模型，用於與資料庫互動
type UserDBModel struct {
	ID           int     `db:"id"`
	Username     string  `db:"username"`
	Email        string  `db:"email"`
	Password     string  `db:"password"`
	Age          *int    `db:"age"`
	Gender       *string `db:"gender"`
	IsVerified   bool    `db:"is_verified"`
	Status       string  `db:"status"`
	LastActiveAt *string `db:"last_active_at"` // 使用 string 處理時間
	ProfileViews int     `db:"profile_views"`
	CreatedAt    string  `db:"created_at"` // 使用 string 處理時間
	UpdatedAt    string  `db:"updated_at"` // 使用 string 處理時間
}

// ToEntity 將資料庫模型轉換為領域實體
func (m *UserEntityMapper) ToEntity(dbModel *UserDBModel) (*entity.UserInformation, error) {
	// 轉換性別
	var gender *entity.Gender
	if dbModel.Gender != nil {
		g := entity.Gender(*dbModel.Gender)
		gender = &g
	}

	// 轉換狀態
	status := entity.UserStatus(dbModel.Status)

	// 這裡可以添加時間轉換邏輯，暫時使用基本轉換
	user := &entity.UserInformation{
		ID:           dbModel.ID,
		Username:     dbModel.Username,
		Email:        dbModel.Email,
		Password:     dbModel.Password,
		Age:          dbModel.Age,
		Gender:       gender,
		IsVerified:   dbModel.IsVerified,
		Status:       status,
		ProfileViews: dbModel.ProfileViews,
		// 時間欄位需要額外的轉換邏輯
	}

	return user, nil
}

// ToDBModel 將領域實體轉換為資料庫模型
func (m *UserEntityMapper) ToDBModel(user *entity.UserInformation) (*UserDBModel, error) {
	// 轉換性別
	var gender *string
	if user.Gender != nil {
		g := string(*user.Gender)
		gender = &g
	}

	// 轉換狀態
	status := string(user.Status)

	dbModel := &UserDBModel{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		Age:          user.Age,
		Gender:       gender,
		IsVerified:   user.IsVerified,
		Status:       status,
		ProfileViews: user.ProfileViews,
		// 時間欄位需要額外的轉換邏輯
	}

	return dbModel, nil
}

// PrepareForInsert 為插入操作準備資料庫模型的值
func (m *UserEntityMapper) PrepareForInsert(user *entity.UserInformation) ([]interface{}, error) {
	dbModel, err := m.ToDBModel(user)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		dbModel.Username,
		dbModel.Email,
		dbModel.Password,
		dbModel.Age,
		dbModel.Gender,
		dbModel.IsVerified,
		dbModel.Status,
		dbModel.ProfileViews,
	}, nil
}

// UserProfileEntityMapper 處理 UserProfile Entity 與資料庫記錄之間的轉換
type UserProfileEntityMapper struct {
	stringArrayMapper *StringArrayMapper
}

func NewUserProfileEntityMapper() *UserProfileEntityMapper {
	return &UserProfileEntityMapper{
		stringArrayMapper: &StringArrayMapper{},
	}
}

// UserProfileDBModel 用戶檔案資料庫模型
type UserProfileDBModel struct {
	ID              int         `db:"id"`
	UserID          int         `db:"user_id"`
	Bio             *string     `db:"bio"`
	Interests       interface{} `db:"interests"`
	LocationLat     *float64    `db:"location_lat"`
	LocationLng     *float64    `db:"location_lng"`
	City            *string     `db:"city"`
	Country         *string     `db:"country"`
	Height          *int        `db:"height"`
	Weight          *int        `db:"weight"`
	Education       *string     `db:"education"`
	Occupation      *string     `db:"occupation"`
	Company         *string     `db:"company"`
	Relationship    *string     `db:"relationship"`
	LookingFor      interface{} `db:"looking_for"`
	Languages       interface{} `db:"languages"`
	Hobbies         interface{} `db:"hobbies"`
	Lifestyle       interface{} `db:"lifestyle"`
	PetPreference   *string     `db:"pet_preference"`
	DrinkingHabit   *string     `db:"drinking_habit"`
	SmokingHabit    *string     `db:"smoking_habit"`
	ExerciseHabit   *string     `db:"exercise_habit"`
	SocialMediaLink *string     `db:"social_media_link"`
	PersonalityType *string     `db:"personality_type"`
	Zodiac          *string     `db:"zodiac"`
	Religion        *string     `db:"religion"`
	CreatedAt       string      `db:"created_at"`
	UpdatedAt       string      `db:"updated_at"`
}

// ToEntity 將資料庫模型轉換為 UserProfile 實體
func (m *UserProfileEntityMapper) ToEntity(dbModel *UserProfileDBModel) (*entity.UserProfile, error) {
	// 轉換各種 StringArray
	interests, err := m.stringArrayMapper.FromDatabase(dbModel.Interests)
	if err != nil {
		return nil, err
	}

	lookingFor, err := m.stringArrayMapper.FromDatabase(dbModel.LookingFor)
	if err != nil {
		return nil, err
	}

	languages, err := m.stringArrayMapper.FromDatabase(dbModel.Languages)
	if err != nil {
		return nil, err
	}

	hobbies, err := m.stringArrayMapper.FromDatabase(dbModel.Hobbies)
	if err != nil {
		return nil, err
	}

	lifestyle, err := m.stringArrayMapper.FromDatabase(dbModel.Lifestyle)
	if err != nil {
		return nil, err
	}

	profile := &entity.UserProfile{
		ID:              dbModel.ID,
		UserID:          dbModel.UserID,
		Bio:             dbModel.Bio,
		Interests:       interests,
		LocationLat:     dbModel.LocationLat,
		LocationLng:     dbModel.LocationLng,
		City:            dbModel.City,
		Country:         dbModel.Country,
		Height:          dbModel.Height,
		Weight:          dbModel.Weight,
		Education:       dbModel.Education,
		Occupation:      dbModel.Occupation,
		Company:         dbModel.Company,
		Relationship:    dbModel.Relationship,
		LookingFor:      lookingFor,
		Languages:       languages,
		Hobbies:         hobbies,
		Lifestyle:       lifestyle,
		PetPreference:   dbModel.PetPreference,
		DrinkingHabit:   dbModel.DrinkingHabit,
		SmokingHabit:    dbModel.SmokingHabit,
		ExerciseHabit:   dbModel.ExerciseHabit,
		SocialMediaLink: dbModel.SocialMediaLink,
		PersonalityType: dbModel.PersonalityType,
		Zodiac:          dbModel.Zodiac,
		Religion:        dbModel.Religion,
		// 時間欄位需要額外的轉換邏輯
	}

	return profile, nil
}

// ToDBModel 將 UserProfile 實體轉換為資料庫模型
func (m *UserProfileEntityMapper) ToDBModel(profile *entity.UserProfile) (*UserProfileDBModel, error) {
	// 轉換各種 StringArray
	interests, err := m.stringArrayMapper.ToDatabase(profile.Interests)
	if err != nil {
		return nil, err
	}

	lookingFor, err := m.stringArrayMapper.ToDatabase(profile.LookingFor)
	if err != nil {
		return nil, err
	}

	languages, err := m.stringArrayMapper.ToDatabase(profile.Languages)
	if err != nil {
		return nil, err
	}

	hobbies, err := m.stringArrayMapper.ToDatabase(profile.Hobbies)
	if err != nil {
		return nil, err
	}

	lifestyle, err := m.stringArrayMapper.ToDatabase(profile.Lifestyle)
	if err != nil {
		return nil, err
	}

	dbModel := &UserProfileDBModel{
		ID:              profile.ID,
		UserID:          profile.UserID,
		Bio:             profile.Bio,
		Interests:       interests,
		LocationLat:     profile.LocationLat,
		LocationLng:     profile.LocationLng,
		City:            profile.City,
		Country:         profile.Country,
		Height:          profile.Height,
		Weight:          profile.Weight,
		Education:       profile.Education,
		Occupation:      profile.Occupation,
		Company:         profile.Company,
		Relationship:    profile.Relationship,
		LookingFor:      lookingFor,
		Languages:       languages,
		Hobbies:         hobbies,
		Lifestyle:       lifestyle,
		PetPreference:   profile.PetPreference,
		DrinkingHabit:   profile.DrinkingHabit,
		SmokingHabit:    profile.SmokingHabit,
		ExerciseHabit:   profile.ExerciseHabit,
		SocialMediaLink: profile.SocialMediaLink,
		PersonalityType: profile.PersonalityType,
		Zodiac:          profile.Zodiac,
		Religion:        profile.Religion,
		// 時間欄位需要額外的轉換邏輯
	}

	return dbModel, nil
}
