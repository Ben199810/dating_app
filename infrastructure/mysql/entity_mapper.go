package mysql

import (
	"fmt"
	"golang_dev_docker/domain/entity"

	"gorm.io/gorm"
)

// AutoMigrateEntities 自動遷移所有實體到資料庫
func AutoMigrateEntities(db *gorm.DB) error {
	entities := []interface{}{
		// 用戶相關實體
		&entity.User{},
		&entity.UserProfile{},
		&entity.Photo{},
		&entity.Interest{},
		&entity.AgeVerification{},

		// 配對相關實體
		&entity.Match{},

		// 聊天相關實體
		&entity.ChatMessage{},

		// 檢舉和封鎖實體
		&entity.Report{},
		&entity.Block{},
	}

	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			return err
		}
	}

	// 創建關聯表
	if err := createAssociationTables(db); err != nil {
		return err
	}

	// 創建索引
	if err := createIndexes(db); err != nil {
		return err
	}

	return nil
}

// createAssociationTables 創建關聯表
func createAssociationTables(db *gorm.DB) error {
	// 用戶興趣關聯表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_interests (
			user_id INT UNSIGNED NOT NULL,
			interest_id INT UNSIGNED NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, interest_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (interest_id) REFERENCES interests(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		return err
	}

	// WebSocket 連接表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS websocket_connections (
			id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			user_id INT UNSIGNED NOT NULL,
			connection_id VARCHAR(255) NOT NULL UNIQUE,
			connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user_connections (user_id),
			INDEX idx_connection_id (connection_id)
		)
	`).Error; err != nil {
		return err
	}

	// 內容審核日誌表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS moderation_logs (
			id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			content_type VARCHAR(50) NOT NULL,
			content_id INT UNSIGNED NOT NULL,
			user_id INT UNSIGNED NOT NULL,
			moderator_id INT UNSIGNED NULL,
			action VARCHAR(50) NOT NULL,
			reason TEXT,
			is_automatic BOOLEAN DEFAULT FALSE,
			confidence DECIMAL(3,2) NULL,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (moderator_id) REFERENCES users(id) ON DELETE SET NULL,
			INDEX idx_content (content_type, content_id),
			INDEX idx_user_moderation (user_id),
			INDEX idx_moderator (moderator_id),
			INDEX idx_action (action),
			INDEX idx_created_at (created_at)
		)
	`).Error; err != nil {
		return err
	}

	return nil
}

// createIndexes 創建重要索引以提升查詢性能
func createIndexes(db *gorm.DB) error {
	indexes := []string{
		// 用戶表索引
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_active_verified ON users(is_active, is_verified)",
		"CREATE INDEX IF NOT EXISTS idx_users_birth_date ON users(birth_date)",

		// 用戶檔案表索引
		"CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_profiles_gender ON user_profiles(gender)",
		"CREATE INDEX IF NOT EXISTS idx_user_profiles_location ON user_profiles(location_lat, location_lng)",
		"CREATE INDEX IF NOT EXISTS idx_user_profiles_age_range ON user_profiles(age_range_min, age_range_max)",

		// 照片表索引
		"CREATE INDEX IF NOT EXISTS idx_photos_user_id ON photos(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_photos_primary ON photos(user_id, is_primary)",
		"CREATE INDEX IF NOT EXISTS idx_photos_order ON photos(user_id, display_order)",

		// 配對表索引
		"CREATE INDEX IF NOT EXISTS idx_matches_users ON matches(user1_id, user2_id)",
		"CREATE INDEX IF NOT EXISTS idx_matches_user1 ON matches(user1_id)",
		"CREATE INDEX IF NOT EXISTS idx_matches_user2 ON matches(user2_id)",
		"CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status)",
		"CREATE INDEX IF NOT EXISTS idx_matches_created_at ON matches(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_matches_matched_at ON matches(matched_at)",

		// 聊天訊息表索引
		"CREATE INDEX IF NOT EXISTS idx_chat_messages_match_id ON chat_messages(match_id)",
		"CREATE INDEX IF NOT EXISTS idx_chat_messages_sender ON chat_messages(sender_id)",
		"CREATE INDEX IF NOT EXISTS idx_chat_messages_receiver ON chat_messages(receiver_id)",
		"CREATE INDEX IF NOT EXISTS idx_chat_messages_status ON chat_messages(status)",
		"CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_chat_messages_unread ON chat_messages(receiver_id, status, is_deleted)",

		// 檢舉表索引
		"CREATE INDEX IF NOT EXISTS idx_reports_reporter ON reports(reporter_id)",
		"CREATE INDEX IF NOT EXISTS idx_reports_reported_user ON reports(reported_user_id)",
		"CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status)",
		"CREATE INDEX IF NOT EXISTS idx_reports_category ON reports(category)",
		"CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at)",

		// 封鎖表索引
		"CREATE INDEX IF NOT EXISTS idx_blocks_blocker ON blocks(blocker_id)",
		"CREATE INDEX IF NOT EXISTS idx_blocks_blocked ON blocks(blocked_id)",
		"CREATE INDEX IF NOT EXISTS idx_blocks_relationship ON blocks(blocker_id, blocked_id)",

		// 年齡驗證表索引
		"CREATE INDEX IF NOT EXISTS idx_age_verification_user_id ON age_verifications(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_age_verification_status ON age_verifications(status)",
		"CREATE INDEX IF NOT EXISTS idx_age_verification_reviewer ON age_verifications(reviewer_id)",

		// 興趣表索引
		"CREATE INDEX IF NOT EXISTS idx_interests_name ON interests(name)",
		"CREATE INDEX IF NOT EXISTS idx_interests_category ON interests(category)",

		// 用戶興趣關聯表索引
		"CREATE INDEX IF NOT EXISTS idx_user_interests_user ON user_interests(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_interests_interest ON user_interests(interest_id)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			// 記錄錯誤但繼續執行，因為索引可能已存在
			continue
		}
	}

	return nil
}

// DropAllTables 刪除所有表（用於測試或重置）
func DropAllTables(db *gorm.DB) error {
	tables := []string{
		"moderation_logs",
		"websocket_connections",
		"user_interests",
		"blocks",
		"reports",
		"chat_messages",
		"matches",
		"age_verifications",
		"photos",
		"interests",
		"user_profiles",
		"users",
	}

	// 禁用外鍵檢查
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// 刪除所有表
	for _, table := range tables {
		if err := db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
			return err
		}
	}

	// 重新啟用外鍵檢查
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return err
	}

	return nil
}

// SeedData 插入種子資料
func SeedData(db *gorm.DB) error {
	// 輔助函數：將字符串轉換為字符串指針
	stringPtr := func(s string) *string {
		return &s
	}

	// 創建基本興趣標籤
	interests := []*entity.Interest{
		{Name: "音樂", Category: entity.InterestCategoryMusic, Description: stringPtr("喜歡聽音樂、演奏樂器")},
		{Name: "電影", Category: entity.InterestCategoryMovies, Description: stringPtr("電影愛好者")},
		{Name: "旅行", Category: entity.InterestCategoryTravel, Description: stringPtr("喜歡旅遊、探索新地方")},
		{Name: "運動", Category: entity.InterestCategorySports, Description: stringPtr("各種運動愛好")},
		{Name: "健身", Category: entity.InterestCategoryFitness, Description: stringPtr("健身、保持身材")},
		{Name: "閱讀", Category: entity.InterestCategoryReading, Description: stringPtr("讀書、學習新知識")},
		{Name: "美食", Category: entity.InterestCategoryFood, Description: stringPtr("美食愛好者、烹飪")},
		{Name: "攝影", Category: entity.InterestCategoryArt, Description: stringPtr("攝影、拍照")},
		{Name: "繪畫", Category: entity.InterestCategoryArt, Description: stringPtr("繪畫、藝術創作")},
		{Name: "寵物", Category: entity.InterestCategoryHobbies, Description: stringPtr("喜歡動物、養寵物")},
		{Name: "遊戲", Category: entity.InterestCategoryGaming, Description: stringPtr("電子遊戲、桌遊")},
		{Name: "舞蹈", Category: entity.InterestCategoryArt, Description: stringPtr("跳舞、各種舞蹈")},
		{Name: "瑜伽", Category: entity.InterestCategoryFitness, Description: stringPtr("瑜伽、冥想")},
		{Name: "咖啡", Category: entity.InterestCategoryHobbies, Description: stringPtr("咖啡愛好者")},
		{Name: "購物", Category: entity.InterestCategoryHobbies, Description: stringPtr("購物、時尚")},
		{Name: "科技", Category: entity.InterestCategoryTechnology, Description: stringPtr("科技產品、程式設計")},
		{Name: "自然", Category: entity.InterestCategoryNature, Description: stringPtr("喜歡大自然、戶外活動")},
		{Name: "車", Category: entity.InterestCategoryHobbies, Description: stringPtr("汽車、機車愛好")},
		{Name: "派對", Category: entity.InterestCategoryHobbies, Description: stringPtr("聚會、派對")},
		{Name: "冒險", Category: entity.InterestCategoryHobbies, Description: stringPtr("冒險活動、極限運動")},
	}

	for _, interest := range interests {
		// 檢查是否已存在
		var existingCount int64
		db.Model(&entity.Interest{}).Where("name = ?", interest.Name).Count(&existingCount)
		if existingCount == 0 {
			if err := db.Create(interest).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// GetDatabaseTables 獲取資料庫中所有表名
func GetDatabaseTables(db *gorm.DB) ([]string, error) {
	var tables []string
	if err := db.Raw("SHOW TABLES").Pluck("Tables_in_"+getDatabaseName(db), &tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

// getDatabaseName 獲取資料庫名稱
func getDatabaseName(db *gorm.DB) string {
	var dbName string
	db.Raw("SELECT DATABASE()").Scan(&dbName)
	return dbName
}

// ValidateDatabase 驗證資料庫結構完整性
func ValidateDatabase(db *gorm.DB) error {
	requiredTables := []string{
		"users",
		"user_profiles",
		"photos",
		"interests",
		"age_verifications",
		"matches",
		"chat_messages",
		"reports",
		"blocks",
		"user_interests",
		"websocket_connections",
		"moderation_logs",
	}

	existingTables, err := GetDatabaseTables(db)
	if err != nil {
		return err
	}

	tableExists := make(map[string]bool)
	for _, table := range existingTables {
		tableExists[table] = true
	}

	var missingTables []string
	for _, requiredTable := range requiredTables {
		if !tableExists[requiredTable] {
			missingTables = append(missingTables, requiredTable)
		}
	}

	if len(missingTables) > 0 {
		return fmt.Errorf("缺少必要的資料庫表: %v", missingTables)
	}

	return nil
}
