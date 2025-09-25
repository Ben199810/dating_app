package redis

import (
	"fmt"
	"time"
)

// Config Redis 配置結構
type Config struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	MaxRetries   int           `yaml:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DefaultConfig 回傳預設的 Redis 配置
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// GetAddr 獲取 Redis 地址字串
func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Validate 驗證配置的有效性
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("Redis host 不能為空")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("Redis port 必須在 1-65535 範圍內")
	}
	if c.DB < 0 {
		return fmt.Errorf("Redis DB 不能為負數")
	}
	if c.PoolSize <= 0 {
		c.PoolSize = 10 // 設定預設值
	}
	if c.MinIdleConns < 0 {
		c.MinIdleConns = 0
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = 3
	}
	if c.DialTimeout <= 0 {
		c.DialTimeout = 5 * time.Second
	}
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 3 * time.Second
	}
	if c.WriteTimeout <= 0 {
		c.WriteTimeout = 3 * time.Second
	}
	return nil
}

// ConnectionStringOptions Redis 連接字串選項
type ConnectionStringOptions struct {
	MaxRetries      int
	RetryDelay      time.Duration
	PoolSize        int
	MinIdleConns    int
	MaxConnAge      time.Duration
	PoolTimeout     time.Duration
	IdleTimeout     time.Duration
	IdleCheckFreq   time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ConnMaxLifetime time.Duration
}

// DefaultConnectionOptions 回傳預設連接選項
func DefaultConnectionOptions() *ConnectionStringOptions {
	return &ConnectionStringOptions{
		MaxRetries:      3,
		RetryDelay:      100 * time.Millisecond,
		PoolSize:        10,
		MinIdleConns:    5,
		MaxConnAge:      time.Hour,
		PoolTimeout:     4 * time.Second,
		IdleTimeout:     5 * time.Minute,
		IdleCheckFreq:   time.Minute,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		ConnMaxLifetime: time.Hour,
	}
}
