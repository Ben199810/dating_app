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
	ID            int             `db:"id"`
	Username      string          `db:"username"`
	Email         string          `db:"email"`
	Password      string          `db:"password"`
	Age           *int            `db:"age"`
	Gender        *string         `db:"gender"`
	Bio           *string         `db:"bio"`
	Interests     interface{}     `db:"interests"`     // 原始 JSON 資料
	LocationLat   *float64        `db:"location_lat"`
	LocationLng   *float64        `db:"location_lng"`
	City          *string         `db:"city"`
	Country       *string         `db:"country"`
	IsVerified    bool            `db:"is_verified"`
	Status        string          `db:"status"`
	LastActiveAt  *string         `db:"last_active_at"` // 使用 string 處理時間
	ProfileViews  int             `db:"profile_views"`
	CreatedAt     string          `db:"created_at"`      // 使用 string 處理時間
	UpdatedAt     string          `db:"updated_at"`      // 使用 string 處理時間
}

// ToEntity 將資料庫模型轉換為領域實體
func (m *UserEntityMapper) ToEntity(dbModel *UserDBModel) (*entity.UserInformation, error) {
	// 轉換 StringArray
	interests, err := m.stringArrayMapper.FromDatabase(dbModel.Interests)
	if err != nil {
		return nil, err
	}
	
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
		ID:            dbModel.ID,
		Username:      dbModel.Username,
		Email:         dbModel.Email,
		Password:      dbModel.Password,
		Age:           dbModel.Age,
		Gender:        gender,
		Bio:           dbModel.Bio,
		Interests:     interests,
		LocationLat:   dbModel.LocationLat,
		LocationLng:   dbModel.LocationLng,
		City:          dbModel.City,
		Country:       dbModel.Country,
		IsVerified:    dbModel.IsVerified,
		Status:        status,
		ProfileViews:  dbModel.ProfileViews,
		// 時間欄位需要額外的轉換邏輯
	}
	
	return user, nil
}

// ToDBModel 將領域實體轉換為資料庫模型
func (m *UserEntityMapper) ToDBModel(user *entity.UserInformation) (*UserDBModel, error) {
	// 轉換 StringArray
	interests, err := m.stringArrayMapper.ToDatabase(user.Interests)
	if err != nil {
		return nil, err
	}
	
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
		Bio:          user.Bio,
		Interests:    interests,
		LocationLat:  user.LocationLat,
		LocationLng:  user.LocationLng,
		City:         user.City,
		Country:      user.Country,
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
		dbModel.Bio,
		dbModel.Interests,
		dbModel.LocationLat,
		dbModel.LocationLng,
		dbModel.City,
		dbModel.Country,
		dbModel.IsVerified,
		dbModel.Status,
		dbModel.ProfileViews,
	}, nil
}