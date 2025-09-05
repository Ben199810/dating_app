package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 設定靜態檔案服務
	r.Static("/static", "./static")

	// 主要頁面路由
	r.StaticFile("/", "./static/html/index.html")
	r.StaticFile("/chat", "./static/html/chat_websocket.html")
	r.StaticFile("/chat_websocket.html", "./static/html/chat_websocket.html") // 保持向後兼容

	// API 路由
	r.GET("/hello", helloHandler)
	r.GET("/user", userHandler)
	r.GET("/ws", wsHandler)
}
