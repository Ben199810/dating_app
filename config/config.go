package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 代表整個應用程式的配置結構
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	Redis    RedisConfig    `yaml:"redis"`
}

// DatabaseConfig 代表資料庫配置
type DatabaseConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parseTime"`
	Loc       string `yaml:"loc"`
}

// ServerConfig 代表伺服器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// LoggingConfig 代表日誌配置
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// RedisConfig 代表 Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// GetDSN 建構資料庫連線字串
func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		db.User, db.Password, db.Host, db.Port, db.DBName,
		db.Charset, db.ParseTime, db.Loc)
}

// LoadConfig 載入指定環境的配置檔案
func LoadConfig(env string) (*Config, error) {
	if env == "" {
		env = getEnv("APP_ENV", "development")
	}

	// 建構配置檔案路徑
	configPath := filepath.Join("config", fmt.Sprintf("%s.yaml", env))

	// 檢查檔案是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置檔案不存在: %s", configPath)
	}

	// 讀取檔案內容
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("無法讀取配置檔案 %s: %w", configPath, err)
	}

	// 解析 YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("無法解析配置檔案 %s: %w", configPath, err)
	}

	return &config, nil
}

// LoadConfigFromPath 從指定路徑載入配置檔案
func LoadConfigFromPath(configPath string) (*Config, error) {
	// 檢查檔案是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置檔案不存在: %s", configPath)
	}

	// 讀取檔案內容
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("無法讀取配置檔案 %s: %w", configPath, err)
	}

	// 解析 YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("無法解析配置檔案 %s: %w", configPath, err)
	}

	return &config, nil
}

// getEnv 取得環境變數，如果不存在則回傳預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
