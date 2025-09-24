package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS 配置跨域資源分享
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 允許的來源 (開發環境較寬鬆，生產環境應該更嚴格)
		allowedOrigins := map[string]bool{
			"http://localhost:8080": true,
			"http://127.0.0.1:8080": true,
			"http://localhost:3000": true, // 前端開發伺服器
			"http://127.0.0.1:3000": true,
		}

		// 檢查來源是否被允許
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 小時

		// 處理預檢請求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// SecurityHeaders 設置安全相關的 HTTP 標頭
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 防止 XSS 攻擊
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// 強制 HTTPS (生產環境)
		if gin.Mode() == gin.ReleaseMode {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Content Security Policy
		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: https:; "+
				"font-src 'self'; "+
				"connect-src 'self' ws: wss:; "+
				"frame-ancestors 'none'")

		// 防止資訊洩漏
		c.Header("Server", "ChatApp/1.0")

		c.Next()
	})
}

// RequestID 為每個請求添加唯一 ID
func RequestID() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 如果沒有提供 Request ID，產生一個簡單的
			requestID = generateSimpleID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)

		c.Next()
	})
}

// generateSimpleID 產生簡單的請求 ID
func generateSimpleID() string {
	// 使用時間戳作為簡單的 ID (生產環境應該使用更好的方法)
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
