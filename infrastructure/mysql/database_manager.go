package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Database        string        `yaml:"database"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Charset         string        `yaml:"charset"`
	Collation       string        `yaml:"collation"`
	Timezone        string        `yaml:"timezone"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	LogLevel        string        `yaml:"log_level"`
	SlowThreshold   time.Duration `yaml:"slow_threshold"`
}

// DefaultDatabaseConfig 預設資料庫配置
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "localhost",
		Port:            3306,
		Database:        "dating_app",
		Username:        "root",
		Password:        "password",
		Charset:         "utf8mb4",
		Collation:       "utf8mb4_unicode_ci",
		Timezone:        "Asia/Taipei",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
		LogLevel:        "info",
		SlowThreshold:   200 * time.Millisecond,
	}
}

// DatabaseManager 資料庫管理器
type DatabaseManager struct {
	db     *gorm.DB
	sqlDB  *sql.DB
	config *DatabaseConfig
}

// NewDatabaseManager 建立新的資料庫管理器
func NewDatabaseManager(config *DatabaseConfig) (*DatabaseManager, error) {
	dsn := buildDSN(config)
	
	// 配置 GORM logger
	gormLogger := getGormLogger(config.LogLevel, config.SlowThreshold)
	
	// 連接資料庫
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("連接資料庫失敗: %w", err)
	}

	// 獲取底層的 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("獲取 SQL DB 失敗: %w", err)
	}

	// 配置連接池
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// 測試連接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("資料庫連接測試失敗: %w", err)
	}

	manager := &DatabaseManager{
		db:     db,
		sqlDB:  sqlDB,
		config: config,
	}

	log.Printf("資料庫連接成功: %s:%d/%s", config.Host, config.Port, config.Database)
	
	return manager, nil
}

// buildDSN 建立資料庫連接字符串
func buildDSN(config *DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=true&loc=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
		config.Collation,
		config.Timezone,
	)
}

// getGormLogger 獲取 GORM logger
func getGormLogger(level string, slowThreshold time.Duration) logger.Interface {
	var logLevel logger.LogLevel

	switch level {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	return logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             slowThreshold,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

// GetDB 獲取 GORM 資料庫實例
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

// GetSQLDB 獲取原始 SQL 資料庫實例
func (dm *DatabaseManager) GetSQLDB() *sql.DB {
	return dm.sqlDB
}

// Close 關閉資料庫連接
func (dm *DatabaseManager) Close() error {
	return dm.sqlDB.Close()
}

// Ping 測試資料庫連接
func (dm *DatabaseManager) Ping() error {
	return dm.sqlDB.Ping()
}

// Stats 獲取連接池統計
func (dm *DatabaseManager) Stats() sql.DBStats {
	return dm.sqlDB.Stats()
}

// 事務管理

// Transaction 執行事務
func (dm *DatabaseManager) Transaction(fn func(*gorm.DB) error) error {
	return dm.db.Transaction(fn)
}

// BeginTransaction 開始事務
func (dm *DatabaseManager) BeginTransaction() *gorm.DB {
	return dm.db.Begin()
}

// 資料庫操作工具方法

// CreateTables 建立表格（用於測試）
func (dm *DatabaseManager) CreateTables(models ...interface{}) error {
	return dm.db.AutoMigrate(models...)
}

// DropTables 刪除表格（用於測試）
func (dm *DatabaseManager) DropTables(models ...interface{}) error {
	for _, model := range models {
		if err := dm.db.Migrator().DropTable(model); err != nil {
			return err
		}
	}
	return nil
}

// TableExists 檢查表格是否存在
func (dm *DatabaseManager) TableExists(tableName string) bool {
	return dm.db.Migrator().HasTable(tableName)
}

// ExecuteSQL 執行原始 SQL
func (dm *DatabaseManager) ExecuteSQL(sql string, values ...interface{}) error {
	return dm.db.Exec(sql, values...).Error
}

// QuerySQL 查詢原始 SQL
func (dm *DatabaseManager) QuerySQL(dest interface{}, sql string, values ...interface{}) error {
	return dm.db.Raw(sql, values...).Scan(dest).Error
}

// 健康檢查

// Health 檢查資料庫健康狀態
func (dm *DatabaseManager) Health() error {
	if err := dm.Ping(); err != nil {
		return fmt.Errorf("資料庫 ping 失敗: %w", err)
	}

	// 檢查連接池狀態
	stats := dm.Stats()
	if stats.OpenConnections >= dm.config.MaxOpenConns {
		return fmt.Errorf("連接池已滿: %d/%d", stats.OpenConnections, dm.config.MaxOpenConns)
	}

	return nil
}

// GetConnectionInfo 獲取連接資訊
func (dm *DatabaseManager) GetConnectionInfo() map[string]interface{} {
	stats := dm.Stats()
	return map[string]interface{}{
		"max_open_connections":     dm.config.MaxOpenConns,
		"max_idle_connections":     dm.config.MaxIdleConns,
		"open_connections":         stats.OpenConnections,
		"idle_connections":         stats.Idle,
		"in_use_connections":       stats.InUse,
		"wait_count":               stats.WaitCount,
		"wait_duration":            stats.WaitDuration,
		"max_idle_closed":          stats.MaxIdleClosed,
		"max_idle_time_closed":     stats.MaxIdleTimeClosed,
		"max_lifetime_closed":      stats.MaxLifetimeClosed,
		"connection_max_lifetime":  dm.config.ConnMaxLifetime,
		"connection_max_idle_time": dm.config.ConnMaxIdleTime,
		"database":                 dm.config.Database,
		"host":                     dm.config.Host,
		"port":                     dm.config.Port,
	}
}

// 資料庫維護操作

// OptimizeTable 優化表格
func (dm *DatabaseManager) OptimizeTable(tableName string) error {
	sql := fmt.Sprintf("OPTIMIZE TABLE %s", tableName)
	return dm.ExecuteSQL(sql)
}

// AnalyzeTable 分析表格
func (dm *DatabaseManager) AnalyzeTable(tableName string) error {
	sql := fmt.Sprintf("ANALYZE TABLE %s", tableName)
	return dm.ExecuteSQL(sql)
}

// CheckTable 檢查表格
func (dm *DatabaseManager) CheckTable(tableName string) error {
	sql := fmt.Sprintf("CHECK TABLE %s", tableName)
	return dm.ExecuteSQL(sql)
}

// RepairTable 修復表格
func (dm *DatabaseManager) RepairTable(tableName string) error {
	sql := fmt.Sprintf("REPAIR TABLE %s", tableName)
	return dm.ExecuteSQL(sql)
}

// GetTableSize 獲取表格大小
func (dm *DatabaseManager) GetTableSize(tableName string) (map[string]interface{}, error) {
	var result struct {
		TableName  string `json:"table_name"`
		DataLength int64  `json:"data_length"`
		IndexLength int64  `json:"index_length"`
		TableRows  int64  `json:"table_rows"`
	}

	sql := `
		SELECT 
			table_name,
			data_length,
			index_length,
			table_rows
		FROM information_schema.tables 
		WHERE table_schema = ? AND table_name = ?
	`

	err := dm.db.Raw(sql, dm.config.Database, tableName).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"table_name":   result.TableName,
		"data_length":  result.DataLength,
		"index_length": result.IndexLength,
		"table_rows":   result.TableRows,
		"total_size":   result.DataLength + result.IndexLength,
	}, nil
}

// GetDatabaseSize 獲取資料庫大小
func (dm *DatabaseManager) GetDatabaseSize() (map[string]interface{}, error) {
	var result struct {
		DatabaseName string `json:"database_name"`
		TotalSize    int64  `json:"total_size"`
		TableCount   int64  `json:"table_count"`
	}

	sql := `
		SELECT 
			table_schema as database_name,
			ROUND(SUM(data_length + index_length)) as total_size,
			COUNT(*) as table_count
		FROM information_schema.tables 
		WHERE table_schema = ?
		GROUP BY table_schema
	`

	err := dm.db.Raw(sql, dm.config.Database).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"database_name": result.DatabaseName,
		"total_size":    result.TotalSize,
		"table_count":   result.TableCount,
	}, nil
}

// 備份和恢復（基本實現）

// CreateBackup 創建備份（簡單的結構備份）
func (dm *DatabaseManager) CreateBackup(outputFile string) error {
	// 這裡應該實現真正的備份邏輯
	// 可以使用 mysqldump 或其他備份工具
	log.Printf("備份功能需要實現: %s", outputFile)
	return nil
}

// RestoreBackup 恢復備份
func (dm *DatabaseManager) RestoreBackup(inputFile string) error {
	// 這裡應該實現真正的恢復邏輯
	log.Printf("恢復功能需要實現: %s", inputFile)
	return nil
}

// 連接池監控

// StartConnectionPoolMonitor 啟動連接池監控
func (dm *DatabaseManager) StartConnectionPoolMonitor(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			stats := dm.Stats()
			log.Printf("連接池狀態 - 打開: %d, 使用中: %d, 閒置: %d, 等待: %d",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
		}
	}()
}

// 上下文支援

// WithContext 創建帶上下文的資料庫實例
func (dm *DatabaseManager) WithContext(ctx context.Context) *gorm.DB {
	return dm.db.WithContext(ctx)
}

// 配置資訊

// GetConfig 獲取資料庫配置
func (dm *DatabaseManager) GetConfig() *DatabaseConfig {
	return dm.config
}