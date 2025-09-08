package main

import (
	"golang_dev_docker/server/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	go handler.HandleMessages()
	r := gin.Default()
	RegisterRoutes(r)
	r.Run(":8080") // 啟動伺服器在 8080 port
}
