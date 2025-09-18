package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang_dev_docker/config"
	"golang_dev_docker/domain/service"
	"golang_dev_docker/infrastructure/mysql"
	"golang_dev_docker/server/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 載入配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatal("無法載入配置:", err)
	}

	// 初始化資料庫連線
	db := initDB(cfg.Database)
	defer db.Close()

	// 初始化 Repository 和 Service
	userRepo := mysql.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// 設定 Handler 的依賴
	handler.SetUserService(userService)

	// 啟動 WebSocket 處理
	go handler.HandleMessages()

	// 設定 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 設定路由
	r := gin.Default()
	RegisterRoutes(r)

	// 啟動伺服器
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

func initDB(dbConfig config.DatabaseConfig) *sql.DB {
	// 使用配置建構 DSN
	dsn := dbConfig.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("無法連接到資料庫:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("無法ping資料庫:", err)
	}

	log.Printf("資料庫連線成功 - Host: %s, DB: %s", dbConfig.Host, dbConfig.DBName)
	return db
}
