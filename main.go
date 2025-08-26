package main

import (
	"github.com/gin-gonic/gin"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Hobbies []string `json:"hobbies"`
}

func main() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"text": "Hello, World!",
		})
	})
	r.GET("/user", func(c *gin.Context) {
		user := User{
			ID:      1,
			Name:    "Ben",
			Email:   "ben@example.com",
			Hobbies: []string{"coding", "music", "travel"},
		}
		c.JSON(200, user)
	})
	r.Run(":8080") // 啟動伺服器在 8080 port
}
