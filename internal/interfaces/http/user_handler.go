package http

import (
	"golang_dev_docker/internal/domain/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 用戶處理器
type UserHandler struct {
	userRepo user.UserRepository
}

// NewUserHandler 創建用戶處理器
func NewUserHandler(userRepo user.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetUser 獲取用戶資訊
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	u, err := h.userRepo.FindByID(user.UserID(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// Hello 簡單的問候
func (h *UserHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Hello, World!",
	})
}
