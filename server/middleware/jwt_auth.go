package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware JWT 認證中間件
type JWTAuthMiddleware struct {
	jwtSecret []byte
}

// NewJWTAuthMiddleware 建立新的 JWT 認證中間件
func NewJWTAuthMiddleware(jwtSecret string) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		jwtSecret: []byte(jwtSecret),
	}
}

// AuthMiddleware JWT 認證中間件
func (m *JWTAuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從標頭獲取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少 Authorization 標頭",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		// 檢查 Bearer 格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "無效的 Authorization 格式，需要 'Bearer token'",
				"code":  "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 驗證 token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("無效的 token: %v", err),
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// 檢查 token 是否過期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token 已過期",
				"code":  "TOKEN_EXPIRED",
			})
			c.Abort()
			return
		}

		// 將用戶資訊存儲到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可選的認證中間件（用戶可以未登入訪問）
func (m *JWTAuthMiddleware) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 沒有 token，繼續執行但不設置用戶資訊
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			// 格式錯誤，但不阻止請求
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.validateToken(tokenString)
		if err != nil {
			// token 無效，但不阻止請求
			c.Next()
			return
		}

		// token 有效，設置用戶資訊
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// RefreshTokenMiddleware 刷新 token 中間件
func (m *JWTAuthMiddleware) RefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查 token 是否即將過期（剩餘時間少於30分鐘）
		if claims, exists := c.Get("jwt_claims"); exists {
			jwtClaims := claims.(*JWTClaims)

			if jwtClaims.ExpiresAt != nil {
				timeUntilExpiry := time.Until(jwtClaims.ExpiresAt.Time)

				if timeUntilExpiry < 30*time.Minute {
					// 生成新的 token
					newToken, err := m.GenerateToken(
						jwtClaims.UserID,
						jwtClaims.Username,
						jwtClaims.Email,
						24*time.Hour, // 24小時有效期
					)

					if err == nil {
						c.Header("X-New-Token", newToken)
					}
				}
			}
		}

		c.Next()
	}
}

// AdminAuthMiddleware 管理員認證中間件
func (m *JWTAuthMiddleware) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先執行基本認證
		m.AuthMiddleware()(c)

		if c.IsAborted() {
			return
		}

		// 檢查用戶是否為管理員
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用戶資訊不存在",
				"code":  "USER_INFO_MISSING",
			})
			c.Abort()
			return
		}

		// TODO: 檢查用戶是否為管理員
		// isAdmin := userService.IsAdmin(userID.(uint))
		// if !isAdmin {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "error": "需要管理員權限",
		//         "code":  "ADMIN_REQUIRED",
		//     })
		//     c.Abort()
		//     return
		// }

		fmt.Printf("管理員認證通過: 用戶 %v\n", userID)
		c.Next()
	}
}

// AgeVerificationMiddleware 年齡驗證中間件
func (m *JWTAuthMiddleware) AgeVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "需要登入",
				"code":  "LOGIN_REQUIRED",
			})
			c.Abort()
			return
		}

		// TODO: 檢查用戶年齡驗證狀態
		// isAgeVerified := userService.IsAgeVerified(userID.(uint))
		// if !isAgeVerified {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "error": "需要通過年齡驗證",
		//         "code":  "AGE_VERIFICATION_REQUIRED",
		//     })
		//     c.Abort()
		//     return
		// }

		fmt.Printf("年齡驗證通過: 用戶 %v\n", userID)
		c.Next()
	}
}

// UserStatusMiddleware 用戶狀態檢查中間件
func (m *JWTAuthMiddleware) UserStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// TODO: 檢查用戶狀態
		// user := userService.GetUser(userID.(uint))
		// if user.Status == "banned" {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "error": "帳戶已被封禁",
		//         "code":  "ACCOUNT_BANNED",
		//     })
		//     c.Abort()
		//     return
		// }
		//
		// if user.Status == "suspended" {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "error": "帳戶已被暫停",
		//         "code":  "ACCOUNT_SUSPENDED",
		//     })
		//     c.Abort()
		//     return
		// }

		// 更新最後活動時間
		c.Set("last_activity", time.Now())

		fmt.Printf("用戶狀態檢查通過: 用戶 %v\n", userID)
		c.Next()
	}
}

// validateToken 驗證 JWT token
func (m *JWTAuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的簽名方法: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析 token 失敗: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("無效的 token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("無法獲取 token claims")
	}

	return claims, nil
}

// GenerateToken 生成 JWT token
func (m *JWTAuthMiddleware) GenerateToken(userID uint, username, email string, duration time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dating-app",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.jwtSecret)
}

// ExtractTokenFromRequest 從請求中提取 token
func (m *JWTAuthMiddleware) ExtractTokenFromRequest(c *gin.Context) (string, error) {
	// 1. 從 Authorization 標頭提取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	// 2. 從查詢參數提取
	if token := c.Query("token"); token != "" {
		return token, nil
	}

	// 3. 從 Cookie 提取
	if cookie, err := c.Cookie("auth_token"); err == nil {
		return cookie, nil
	}

	return "", fmt.Errorf("未找到 token")
}

// SetTokenCookie 設置 token cookie
func (m *JWTAuthMiddleware) SetTokenCookie(c *gin.Context, token string, maxAge int) {
	c.SetCookie(
		"auth_token", // cookie 名稱
		token,        // cookie 值
		maxAge,       // 最大存活時間（秒）
		"/",          // 路徑
		"",           // 域名
		false,        // 安全（生產環境應設為 true）
		true,         // HTTP Only
	)
}

// ClearTokenCookie 清除 token cookie
func (m *JWTAuthMiddleware) ClearTokenCookie(c *gin.Context) {
	c.SetCookie(
		"auth_token",
		"",
		-1, // 立即過期
		"/",
		"",
		false,
		true,
	)
}

// GetUserFromContext 從上下文獲取用戶資訊
func GetUserFromContext(c *gin.Context) (uint, string, string, bool) {
	userID, exists1 := c.Get("user_id")
	username, exists2 := c.Get("username")
	email, exists3 := c.Get("email")

	if !exists1 || !exists2 || !exists3 {
		return 0, "", "", false
	}

	return userID.(uint), username.(string), email.(string), true
}

// RequireAuth 檢查是否已認證的工具函數
func RequireAuth(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}
