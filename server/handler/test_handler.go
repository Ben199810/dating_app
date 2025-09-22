package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

type TestHandler struct {
	userRepo repository.UserRepository
}

func NewTestHandler(userRepo repository.UserRepository) *TestHandler {
	return &TestHandler{
		userRepo: userRepo,
	}
}

// TestUsersByGender 測試按性別搜尋用戶
func (h *TestHandler) TestUsersByGender(c *gin.Context) {
	gender := c.Param("gender")
	limitStr := c.DefaultQuery("limit", "10")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 limit 參數"})
		return
	}

	// 轉換為 Gender 枚舉
	var genderEnum string
	switch gender {
	case "male", "female", "other":
		genderEnum = gender
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的性別參數"})
		return
	}

	users, err := h.userRepo.GetUsersByGender(entity.Gender(genderEnum), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 過濾敏感資訊
	var safeUsers []gin.H
	for _, user := range users {
		safeUser := gin.H{
			"id":       user.ID,
			"username": user.Username,
			"age":      user.Age,
			"gender":   user.Gender,
			"bio":      user.Bio,
			"city":     user.City,
			"country":  user.Country,
		}
		safeUsers = append(safeUsers, safeUser)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": safeUsers,
		"count": len(safeUsers),
	})
}

// TestUsersById 測試根據 ID 獲取用戶
func (h *TestHandler) TestUsersById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶 ID"})
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用戶不存在"})
		return
	}

	// 過濾敏感資訊
	safeUser := gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"age":           user.Age,
		"gender":        user.Gender,
		"bio":           user.Bio,
		"interests":     user.Interests,
		"location_lat":  user.LocationLat,
		"location_lng":  user.LocationLng,
		"city":          user.City,
		"country":       user.Country,
		"is_verified":   user.IsVerified,
		"status":        user.Status,
		"profile_views": user.ProfileViews,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}

	c.JSON(http.StatusOK, safeUser)
}