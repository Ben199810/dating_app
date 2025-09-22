package handler

import (
	"golang_dev_docker/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var userService *service.UserService
var authService *service.AuthService
var userProfileService *service.UserProfileService

func SetUserService(us *service.UserService) {
	userService = us
}

func SetAuthService(as *service.AuthService) {
	authService = as
}

func SetUserProfileService(ups *service.UserProfileService) {
	userProfileService = ups
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

// LoginHandler 處理用戶登入請求
func LoginHandler(c *gin.Context) {
	if authService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "認證服務未初始化"})
		return
	}

	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式: " + err.Error()})
		return
	}

	user, err := authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"user":    user,
	})
}

// GetUserProfileHandler 獲取用戶個人資料
func GetUserProfileHandler(c *gin.Context) {
	if userProfileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用戶資料服務未初始化"})
		return
	}

	// 從查詢參數或 JWT Token 中獲取用戶 ID
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少用戶 ID"})
		return
	}

	profile, err := userProfileService.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功獲取用戶資料",
		"profile": profile,
	})
}

// UpdateUserProfileHandler 更新用戶個人資料
func UpdateUserProfileHandler(c *gin.Context) {
	if userProfileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用戶資料服務未初始化"})
		return
	}

	var req service.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式: " + err.Error()})
		return
	}

	// 從查詢參數或 JWT Token 中獲取用戶 ID
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少用戶 ID"})
		return
	}

	err := userProfileService.UpdateUserProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用戶資料更新成功",
	})
}
