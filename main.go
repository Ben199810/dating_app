package main

import (
	"log"
	"os"
	"time"

	"golang_dev_docker/config"
	"golang_dev_docker/infrastructure/mysql"
	"golang_dev_docker/infrastructure/redis"
	"golang_dev_docker/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 載入配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("載入配置失敗: %v", err)
	}

	// 調試：輸出資料庫配置
	log.Printf("資料庫配置: Host=%s, Port=%d, User=%s, DBName=%s", 
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.DBName)
	log.Printf("DSN: %s", cfg.Database.GetDSN())

	// 建立伺服器實例
	serverConfig := &server.ServerConfig{
		Port:                    cfg.Server.Port,
		Mode:                    cfg.Server.Mode,
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            10 * time.Second,
		IdleTimeout:             60 * time.Second,
		MaxHeaderBytes:          1 << 20, // 1MB
		GracefulShutdownTimeout: 30 * time.Second,
		StaticPath:              "./static",
		UploadPath:              "./uploads",
		JWTSecret:               "dating-app-secret-key", // TODO: 從環境變數或配置讀取
	}

	srv := server.NewServer(serverConfig)

	// 初始化資料庫
	dbConfig := &mysql.DatabaseConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Database:        cfg.Database.DBName,
		Username:        cfg.Database.User,
		Password:        cfg.Database.Password,
		Charset:         "utf8mb4",
		Collation:       "utf8mb4_unicode_ci",
		Timezone:        "UTC",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
		LogLevel:        "info",
		SlowThreshold:   200 * time.Millisecond,
	}

	if err := srv.InitializeDatabase(dbConfig); err != nil {
		log.Fatalf("初始化資料庫失敗: %v", err)
	}

	// 初始化 Redis 配置
	redisConfig := &redis.RedisConfig{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}

	if err := srv.InitializeRedis(redisConfig); err != nil {
		log.Printf("Redis 初始化失敗，繼續無快取模式: %v", err)
	}

	// 初始化 WebSocket
	srv.InitializeWebSocket()

	// 初始化中間件
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	srv.InitializeMiddleware(env)

	// 設置路由
	srv.SetupRoutes()

	// 啟動伺服器（支援優雅關閉）
	log.Printf("啟動 %s 環境的伺服器...", env)
	if err := srv.StartWithGracefulShutdown(); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}
