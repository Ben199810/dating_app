package handler

import (
	"github.com/gin-gonic/gin"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Hobbies []string `json:"hobbies"`
}

func UserHandler(c *gin.Context) {
	user := User{
		ID:      1,
		Name:    "Ben",
		Email:   "ben@example.com",
		Hobbies: []string{"coding", "music", "travel"},
	}
	c.JSON(200, user)
}
