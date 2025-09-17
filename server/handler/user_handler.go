package handler

import (
	"golang_dev_docker/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Hobbies []string `json:"hobbies"`
}

var userService *service.UserService

func SetUserService(us *service.UserService) {
	userService = us
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

func CreateUserHandler(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用戶服務未初始化"})
		return
	}

	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式: " + err.Error()})
		return
	}

	user, err := userService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用戶建立成功",
		"user":    user,
	})
}
