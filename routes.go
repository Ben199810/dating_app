package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/hello", helloHandler)
	r.GET("/user", userHandler)
	r.GET("/messages", getMessagesHandler)
	r.POST("/messages", postMessageHandler)
	r.StaticFile("/chat.html", "./chat.html")
	r.GET("/ws", wsHandler)
}
