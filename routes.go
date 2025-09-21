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

	// API 路由群組
	api := r.Group("/api")
	{
		api.GET("/status", handler.HealthCheckHandler)        // 健康檢查端點
		api.POST("/user/register", handler.CreateUserHandler) // 新增用戶註冊端點
	}

	// WebSocket 路由
	r.GET("/ws", handler.WSHandler)
}
