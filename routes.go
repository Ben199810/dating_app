package main

import (
	"golang_dev_docker/server/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 設定靜態檔案服務
	r.Static("/static", "./static")

	// 主要頁面路由
	r.StaticFile("/", "./static/html/index.html")
	r.StaticFile("/chat", "./static/html/chat_websocket.html")
	r.StaticFile("/register", "./static/html/register.html")
	r.StaticFile("/login", "./static/html/login.html")
	r.StaticFile("/profile", "./static/html/profile.html")

	// API 路由群組
	api := r.Group("/api")
	{
		// 基本功能
		api.GET("/status", handler.HealthCheckHandler)        // 健康檢查端點
		api.POST("/user/register", handler.CreateUserHandler) // 新增用戶註冊端點

		// 身份驗證
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.LoginHandler) // 用戶登入端點
		}

		// 用戶個人資料
		user := api.Group("/user")
		{
			user.GET("/profile", handler.GetUserProfileHandler)    // 獲取用戶個人資料
			user.PUT("/profile", handler.UpdateUserProfileHandler) // 更新用戶個人資料
		}

		// 用戶資料管理 (需要實作 handler 初始化)
		// api.PUT("/users/:id/basic-info", handler.UpdateBasicInfo)
		// api.PUT("/users/:id/location", handler.UpdateLocation)
		// api.POST("/users/:id/photos", handler.AddPhoto)
		// api.GET("/users/:id/nearby", handler.FindNearbyUsers)
		// api.GET("/users/:id/search", handler.SearchUsers)
		// api.POST("/users/:id/profile", handler.CreateProfile)
	}

	// WebSocket 路由
	r.GET("/ws", handler.WSHandler)
}
