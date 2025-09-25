package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// WebSocketAuthMiddleware WebSocket 認證中間件
type WebSocketAuthMiddleware struct {
	jwtSecret []byte
}

// NewWebSocketAuthMiddleware 建立新的 WebSocket 認證中間件
func NewWebSocketAuthMiddleware(jwtSecret string) *WebSocketAuthMiddleware {
	return &WebSocketAuthMiddleware{
		jwtSecret: []byte(jwtSecret),
	}
}

// JWTClaims JWT 聲明結構
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// AuthenticateWebSocket WebSocket 連接認證
func (m *WebSocketAuthMiddleware) AuthenticateWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 從查詢參數獲取 token
		token := c.Query("token")

		// 2. 如果查詢參數中沒有，則從 Authorization 標頭獲取
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				// 移除 "Bearer " 前綴
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
				} else {
					token = authHeader
				}
			}
		}

		// 3. 如果還是沒有 token，嘗試從 Cookie 獲取
		if token == "" {
			if cookie, err := c.Cookie("auth_token"); err == nil {
				token = cookie
			}
		}

		// 4. 檢查 token 是否存在
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少認證 token",
				"code":  "MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		// 5. 驗證和解析 token
		claims, err := m.validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("無效的 token: %v", err),
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// 6. 檢查 token 是否過期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token 已過期",
				"code":  "TOKEN_EXPIRED",
			})
			c.Abort()
			return
		}

		// 7. 將用戶資訊存儲到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("jwt_claims", claims)

		// 8. 記錄認證成功
		c.Header("X-Authenticated-User", fmt.Sprintf("%d", claims.UserID))

		c.Next()
	}
}

// validateToken 驗證 JWT token
func (m *WebSocketAuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 驗證簽名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的簽名方法: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析 token 失敗: %w", err)
	}

	// 檢查 token 是否有效
	if !token.Valid {
		return nil, fmt.Errorf("無效的 token")
	}

	// 獲取 claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("無法獲取 token claims")
	}

	return claims, nil
}

// RateLimitMiddleware WebSocket 連接速率限制中間件
func (m *WebSocketAuthMiddleware) RateLimitMiddleware(connectionsPerMinute int) gin.HandlerFunc {
	// 簡單的記憶體存儲，生產環境建議使用 Redis
	connectionCounts := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// 清理過期的記錄（超過1分鐘）
		if records, exists := connectionCounts[clientIP]; exists {
			var validRecords []time.Time
			for _, timestamp := range records {
				if now.Sub(timestamp) < time.Minute {
					validRecords = append(validRecords, timestamp)
				}
			}
			connectionCounts[clientIP] = validRecords
		}

		// 檢查連接次數
		if len(connectionCounts[clientIP]) >= connectionsPerMinute {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "連接過於頻繁，請稍後再試",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		// 記錄新的連接
		connectionCounts[clientIP] = append(connectionCounts[clientIP], now)

		c.Next()
	}
}

// WebSocketCORSMiddleware WebSocket CORS 中間件
func (m *WebSocketAuthMiddleware) WebSocketCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 在生產環境中應該檢查允許的來源
		allowedOrigins := []string{
			"http://localhost:8080",
			"https://your-domain.com",
		}

		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		if !isAllowed && origin != "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "不允許的來源",
				"code":  "CORS_NOT_ALLOWED",
			})
			c.Abort()
			return
		}

		// 設置 CORS 標頭
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// 處理 OPTIONS 預檢請求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// UserStatusMiddleware 檢查用戶狀態中間件
func (m *WebSocketAuthMiddleware) UserStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用戶 ID 不存在",
				"code":  "USER_ID_MISSING",
			})
			c.Abort()
			return
		}

		// TODO: 這裡應該檢查用戶是否被封鎖、停用等
		// userService.IsUserActive(userID.(uint))
		// userService.IsUserBlocked(userID.(uint))
		
		// 設置用戶最後活動時間
		c.Set("last_activity", time.Now())
		
		// 暫時記錄用戶狀態檢查
		fmt.Printf("用戶狀態檢查: 用戶 %v\n", userID)
		
		c.Next()
	}
}// WebSocketSecurityMiddleware 安全中間件
func (m *WebSocketAuthMiddleware) WebSocketSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 設置安全標頭
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// 檢查請求是否為 WebSocket 升級請求
		if c.Request.Header.Get("Upgrade") != "websocket" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "非 WebSocket 請求",
				"code":  "NOT_WEBSOCKET_REQUEST",
			})
			c.Abort()
			return
		}

		// 檢查連接標頭
		if c.Request.Header.Get("Connection") != "Upgrade" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "無效的連接標頭",
				"code":  "INVALID_CONNECTION_HEADER",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoggingMiddleware WebSocket 連接記錄中間件
func (m *WebSocketAuthMiddleware) LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// WebSocketMetricsMiddleware 指標收集中間件
func (m *WebSocketAuthMiddleware) WebSocketMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 記錄指標
		duration := time.Since(start)
		userID, exists := c.Get("user_id")

		// TODO: 發送指標到監控系統
		// metrics.RecordWebSocketConnection(userID, duration, c.Writer.Status())

		if exists {
			fmt.Printf("WebSocket 連接指標 - 用戶: %v, 耗時: %v, 狀態: %d\n",
				userID, duration, c.Writer.Status())
		} else {
			fmt.Printf("WebSocket 連接指標 - 匿名用戶, 耗時: %v, 狀態: %d\n",
				duration, c.Writer.Status())
		}
	}
}

// CreateAuthMiddlewareChain 建立完整的認證中間件鏈
func (m *WebSocketAuthMiddleware) CreateAuthMiddlewareChain() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.LoggingMiddleware(),           // 記錄
		m.WebSocketMetricsMiddleware(),  // 指標
		m.WebSocketCORSMiddleware(),     // CORS
		m.WebSocketSecurityMiddleware(), // 安全
		m.RateLimitMiddleware(10),       // 速率限制（每分鐘10次連接）
		m.AuthenticateWebSocket(),       // JWT 認證
		m.UserStatusMiddleware(),        // 用戶狀態檢查
	}
}

// ValidateUserPermission 驗證用戶權限
func (m *WebSocketAuthMiddleware) ValidateUserPermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用戶未認證",
				"code":  "USER_NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}

		// TODO: 實作權限檢查邏輯
		// hasPermission := permissionService.UserHasPermission(userID.(uint), requiredPermission)
		// if !hasPermission {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "error": "權限不足",
		//         "code":  "INSUFFICIENT_PERMISSION",
		//     })
		//     c.Abort()
		//     return
		// }

		fmt.Printf("用戶 %v 權限檢查: %s\n", userID, requiredPermission)

		c.Next()
	}
}

// ExtractUserFromToken 從 token 中提取用戶資訊（工具函數）
func (m *WebSocketAuthMiddleware) ExtractUserFromToken(tokenString string) (*JWTClaims, error) {
	return m.validateToken(tokenString)
}

// GenerateToken 生成 JWT token（工具函數）
func (m *WebSocketAuthMiddleware) GenerateToken(userID uint, username, email string, duration time.Duration) (string, error) {
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
