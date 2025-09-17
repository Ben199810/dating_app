package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"golang_dev_docker/domain/service"
	"golang_dev_docker/infrastructure/mysql"
	"golang_dev_docker/server/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 初始化資料庫連線
	db := initDB()
	defer db.Close()

	// 初始化 Repository 和 Service
	userRepo := mysql.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// 設定 Handler 的依賴
	handler.SetUserService(userService)

	// 啟動 WebSocket 處理
	go handler.HandleMessages()

	// 設定路由
	r := gin.Default()
	RegisterRoutes(r)

	// 啟動伺服器
	r.Run(":8080")
}

func initDB() *sql.DB {
	// 從環境變數獲取資料庫連線資訊
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "chat_user")
	dbPassword := getEnv("DB_PASSWORD", "chat_password")
	dbName := getEnv("DB_NAME", "chat_app")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("無法連接到資料庫:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("無法ping資料庫:", err)
	}

	log.Println("資料庫連線成功")
	return db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
