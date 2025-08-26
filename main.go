package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	RegisterRoutes(r)
	r.Run(":8080") // 啟動伺服器在 8080 port
}
