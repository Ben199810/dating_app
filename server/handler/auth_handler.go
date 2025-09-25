package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"golang_dev_docker/domain/usecase"
)

// AuthHandler 認證處理器
type AuthHandler struct {
	userService *usecase.UserService
	jwtSecret   string
}

// 全域認證處理器實例
var authHandler *AuthHandler

// SetUserServiceForAuth 設置認證處理器的用戶服務依賴
func SetUserServiceForAuth(userService *usecase.UserService, jwtSecret string) {
	authHandler = &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

// RegisterRequest 註冊請求結構
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	BirthDate   string `json:"birth_date" binding:"required"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=50"`
	Gender      string `json:"gender" binding:"required"`
	Biography   string `json:"biography"`
}

// LoginRequest 登入請求結構
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterHandler 處理用戶註冊請求
// POST /api/auth/register
func RegisterHandler(c *gin.Context) {
	if authHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "認證服務未初始化",
		})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求資料格式錯誤",
			"details": err.Error(),
		})
		return
	}

	// 解析出生日期
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "出生日期格式錯誤，請使用 YYYY-MM-DD 格式",
		})
		return
	}

	// 構建服務層請求
	serviceReq := &usecase.RegisterRequest{
		Email:       req.Email,
		Password:    req.Password,
		BirthDate:   birthDate,
		DisplayName: req.DisplayName,
		Gender:      req.Gender,
		Biography:   req.Biography,
	}

	// 調用用戶服務註冊
	userResponse, err := authHandler.userService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		// 根據錯誤類型返回適當的 HTTP 狀態碼
		if err.Error() == "用戶必須年滿18歲才能註冊" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "年齡驗證失敗",
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "此 Email 已被註冊" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Email 已存在",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "註冊失敗",
			"message": err.Error(),
		})
		return
	}

	// 成功回應
	c.JSON(http.StatusCreated, gin.H{
		"message": "註冊成功",
		"user_id": userResponse.ID,
		"user": gin.H{
			"id":          userResponse.ID,
			"email":       userResponse.Email,
			"is_verified": userResponse.IsVerified,
			"age":         userResponse.Age,
		},
	})
}

// LoginHandler 處理用戶登入請求
// POST /api/auth/login
func LoginHandler(c *gin.Context) {
	if authHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "認證服務未初始化",
		})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求資料格式錯誤",
			"details": err.Error(),
		})
		return
	}

	// 構建服務層請求
	serviceReq := &usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// 調用用戶服務登入
	userResponse, err := authHandler.userService.Login(c.Request.Context(), serviceReq)
	if err != nil {
		// 登入錯誤統一返回 401
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "登入失敗",
			"message": err.Error(),
		})
		return
	}

	// 生成 JWT token
	token, err := generateJWT(userResponse.ID, authHandler.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token 生成失敗",
		})
		return
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"token":   token,
		"user": gin.H{
			"id":          userResponse.ID,
			"email":       userResponse.Email,
			"is_verified": userResponse.IsVerified,
			"age":         userResponse.Age,
		},
	})
}

// generateJWT 生成 JWT token
func generateJWT(userID uint, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24小時過期
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
