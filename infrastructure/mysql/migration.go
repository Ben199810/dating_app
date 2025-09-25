package mysql

import (
	"fmt"
	"log"

	"golang_dev_docker/domain/entity"

	"gorm.io/gorm"
)

// Migration 資料庫遷移管理器
type Migration struct {
	db *gorm.DB
}

// NewMigration 創建遷移管理器
func NewMigration(db *gorm.DB) *Migration {
	return &Migration{db: db}
}

// AutoMigrate 自動遷移所有表結構
func (m *Migration) AutoMigrate() error {
	log.Println("開始執行資料庫遷移...")

	// 定義需要遷移的模型
	models := []interface{}{
		&entity.User{},
		&entity.UserProfile{},
		&entity.Photo{},
		&entity.Interest{},
		&entity.AgeVerification{},
		&entity.Match{},
		&entity.ChatMessage{},
		&entity.Report{},
		&entity.Block{},
	}

	// 執行自動遷移
	for _, model := range models {
		if err := m.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("遷移模型 %T 失敗: %w", model, err)
		}
		log.Printf("模型 %T 遷移成功", model)
	}

	// 創建必要的索引
	if err := m.createIndexes(); err != nil {
		return fmt.Errorf("創建索引失敗: %w", err)
	}

	log.Println("資料庫遷移完成")
	return nil
}

// createIndexes 創建必要的索引
func (m *Migration) createIndexes() error {
	log.Println("創建資料庫索引...")

	indexes := []struct {
		table string
		sql   string
		desc  string
	}{
		{
			table: "users",
			sql:   "CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
			desc:  "用戶信箱索引",
		},
		{
			table: "users",
			sql:   "CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)",
			desc:  "用戶啟用狀態索引",
		},
		{
			table: "user_profiles",
			sql:   "CREATE INDEX IF NOT EXISTS idx_profiles_location ON user_profiles(latitude, longitude)",
			desc:  "用戶位置索引",
		},
		{
			table: "user_profiles",
			sql:   "CREATE INDEX IF NOT EXISTS idx_profiles_age_prefs ON user_profiles(min_age_preference, max_age_preference)",
			desc:  "年齡偏好索引",
		},
		{
			table: "matches",
			sql:   "CREATE INDEX IF NOT EXISTS idx_matches_users ON matches(user1_id, user2_id)",
			desc:  "配對用戶索引",
		},
		{
			table: "matches",
			sql:   "CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status)",
			desc:  "配對狀態索引",
		},
		{
			table: "chat_messages",
			sql:   "CREATE INDEX IF NOT EXISTS idx_messages_match ON chat_messages(match_id, created_at)",
			desc:  "聊天訊息索引",
		},
		{
			table: "chat_messages",
			sql:   "CREATE INDEX IF NOT EXISTS idx_messages_receiver ON chat_messages(receiver_id, created_at)",
			desc:  "訊息接收者索引",
		},
		{
			table: "photos",
			sql:   "CREATE INDEX IF NOT EXISTS idx_photos_user ON photos(user_id, display_order)",
			desc:  "用戶照片索引",
		},
	}

	for _, idx := range indexes {
		if err := m.db.Exec(idx.sql).Error; err != nil {
			log.Printf("警告: 創建索引失敗 (%s): %v", idx.desc, err)
			// 不中斷遷移，只是記錄警告
		} else {
			log.Printf("索引創建成功: %s", idx.desc)
		}
	}

	log.Println("索引創建完成")
	return nil
}

// DropAllTables 刪除所有表（用於重置資料庫）
func (m *Migration) DropAllTables() error {
	log.Println("警告: 正在刪除所有資料表...")

	// 禁用外鍵約束檢查
	if err := m.db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("禁用外鍵檢查失敗: %w", err)
	}

	// 獲取所有表名
	tables := []string{
		"blocks", "reports", "chat_messages", "matches",
		"age_verifications", "user_interests", "interests",
		"photos", "user_profiles", "users",
	}

	// 刪除表
	for _, table := range tables {
		if err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			log.Printf("警告: 刪除表 %s 失敗: %v", table, err)
		} else {
			log.Printf("表 %s 已刪除", table)
		}
	}

	// 啟用外鍵約束檢查
	if err := m.db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("啟用外鍵檢查失敗: %w", err)
	}

	log.Println("所有表已刪除")
	return nil
}

// CheckTablesExist 檢查表是否存在
func (m *Migration) CheckTablesExist() (map[string]bool, error) {
	tables := map[string]bool{
		"users":             false,
		"user_profiles":     false,
		"photos":            false,
		"interests":         false,
		"user_interests":    false,
		"age_verifications": false,
		"matches":           false,
		"chat_messages":     false,
		"reports":           false,
		"blocks":            false,
	}

	for table := range tables {
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s'", table)
		if err := m.db.Raw(query).Scan(&count).Error; err != nil {
			return nil, fmt.Errorf("檢查表 %s 失敗: %w", table, err)
		}
		tables[table] = count > 0
	}

	return tables, nil
}
