package main

import (
	"github.com/gin-gonic/gin"
	"golang_dev_docker/server/handler"
)

func RegisterRoutes(r *gin.Engine) {
	// 設定靜態檔案服務
	r.Static("/static", "./static")

	// 主要頁面路由
	r.StaticFile("/", "./static/html/index.html")
	r.StaticFile("/chat", "./static/html/chat_websocket.html")

	// API 路由
	r.GET("/hello", handler.HelloHandler)
	r.GET("/user", handler.UserHandler)
	r.GET("/ws", handler.WSHandler)
}
