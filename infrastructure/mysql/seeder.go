package mysql

import (
	"log"
	"time"

	"golang_dev_docker/domain/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seeder 種子資料管理器
type Seeder struct {
	db *gorm.DB
}

// NewSeeder 創建種子資料管理器
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// SeedAll 植入所有種子資料
func (s *Seeder) SeedAll() error {
	log.Println("開始植入種子資料...")

	// 按依賴順序植入資料
	if err := s.seedInterests(); err != nil {
		return err
	}

	if err := s.seedDemoUsers(); err != nil {
		return err
	}

	log.Println("種子資料植入完成")
	return nil
}

// seedInterests 植入興趣資料
func (s *Seeder) seedInterests() error {
	log.Println("植入興趣資料...")

	interests := []entity.Interest{
		{Name: "音樂", Category: "娛樂", IsActive: true},
		{Name: "電影", Category: "娛樂", IsActive: true},
		{Name: "閱讀", Category: "文化", IsActive: true},
		{Name: "旅行", Category: "生活", IsActive: true},
		{Name: "運動", Category: "健康", IsActive: true},
		{Name: "健身", Category: "健康", IsActive: true},
		{Name: "瑜伽", Category: "健康", IsActive: true},
		{Name: "烹飪", Category: "生活", IsActive: true},
		{Name: "攝影", Category: "藝術", IsActive: true},
		{Name: "繪畫", Category: "藝術", IsActive: true},
		{Name: "舞蹈", Category: "藝術", IsActive: true},
		{Name: "游泳", Category: "運動", IsActive: true},
		{Name: "跑步", Category: "運動", IsActive: true},
		{Name: "登山", Category: "運動", IsActive: true},
		{Name: "騎車", Category: "運動", IsActive: true},
		{Name: "咖啡", Category: "飲食", IsActive: true},
		{Name: "紅酒", Category: "飲食", IsActive: true},
		{Name: "茶道", Category: "文化", IsActive: true},
		{Name: "動漫", Category: "娛樂", IsActive: true},
		{Name: "遊戲", Category: "娛樂", IsActive: true},
		{Name: "寵物", Category: "生活", IsActive: true},
		{Name: "園藝", Category: "生活", IsActive: true},
	}

	for _, interest := range interests {
		// 檢查是否已存在
		var existing entity.Interest
		result := s.db.Where("name = ?", interest.Name).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := s.db.Create(&interest).Error; err != nil {
				return err
			}
			log.Printf("興趣 '%s' 已植入", interest.Name)
		}
	}

	log.Println("興趣資料植入完成")
	return nil
}

// seedDemoUsers 植入示範用戶
func (s *Seeder) seedDemoUsers() error {
	log.Println("植入示範用戶...")

	// 檢查是否已有用戶
	var userCount int64
	s.db.Model(&entity.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("已存在用戶資料，跳過示範用戶植入")
		return nil
	}

	// 創建密碼哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("demo123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	demoUsers := []struct {
		User    entity.User
		Profile entity.UserProfile
	}{
		{
			User: entity.User{
				Email:        "alice@example.com",
				PasswordHash: string(passwordHash),
				BirthDate:    time.Date(1995, 6, 15, 0, 0, 0, 0, time.UTC),
				IsVerified:   true,
				IsActive:     true,
			},
			Profile: entity.UserProfile{
				DisplayName: "Alice",
				Bio:         "喜歡旅行和攝影的女生，期待遇到有趣的靈魂~",
				Gender:      entity.GenderFemale,
				ShowAge:     true,
				LocationLat: ptr(25.0330),
				LocationLng: ptr(121.5654),
				MaxDistance: 30,
				AgeRangeMin: 25,
				AgeRangeMax: 35,
			},
		},
		{
			User: entity.User{
				Email:        "bob@example.com",
				PasswordHash: string(passwordHash),
				BirthDate:    time.Date(1992, 3, 22, 0, 0, 0, 0, time.UTC),
				IsVerified:   true,
				IsActive:     true,
			},
			Profile: entity.UserProfile{
				DisplayName: "Bob",
				Bio:         "熱愛運動和音樂，正在尋找生活中的另一半",
				Gender:      entity.GenderMale,
				ShowAge:     true,
				LocationLat: ptr(25.0478),
				LocationLng: ptr(121.5318),
				MaxDistance: 25,
				AgeRangeMin: 22,
				AgeRangeMax: 30,
			},
		},
		{
			User: entity.User{
				Email:        "charlie@example.com",
				PasswordHash: string(passwordHash),
				BirthDate:    time.Date(1998, 9, 8, 0, 0, 0, 0, time.UTC),
				IsVerified:   true,
				IsActive:     true,
			},
			Profile: entity.UserProfile{
				DisplayName: "Charlie",
				Bio:         "咖啡愛好者，週末喜歡看電影和閱讀",
				Gender:      entity.GenderOther,
				ShowAge:     true,
				LocationLat: ptr(25.0419),
				LocationLng: ptr(121.5430),
				MaxDistance: 40,
				AgeRangeMin: 20,
				AgeRangeMax: 35,
			},
		},
	}

	for i, userData := range demoUsers {
		// 創建用戶
		if err := s.db.Create(&userData.User).Error; err != nil {
			return err
		}

		// 設定個人檔案的用戶ID
		userData.Profile.UserID = userData.User.ID
		if err := s.db.Create(&userData.Profile).Error; err != nil {
			return err
		}

		log.Printf("示範用戶 %d (%s) 已植入", i+1, userData.User.Email)
	}

	log.Println("示範用戶植入完成")
	return nil
}

// ClearAll 清除所有資料（保留表結構）
func (s *Seeder) ClearAll() error {
	log.Println("警告: 正在清除所有資料...")

	// 按相反順序刪除（避免外鍵約束問題）
	tables := []string{
		"blocks", "reports", "chat_messages", "matches",
		"age_verifications", "user_interests", "photos",
		"user_profiles", "users", "interests",
	}

	for _, table := range tables {
		if err := s.db.Exec("DELETE FROM " + table).Error; err != nil {
			log.Printf("警告: 清除表 %s 失敗: %v", table, err)
		} else {
			log.Printf("表 %s 資料已清除", table)
		}
	}

	log.Println("所有資料已清除")
	return nil
}

// ptr 輔助函數：創建指標
func ptr[T any](v T) *T {
	return &v
}
