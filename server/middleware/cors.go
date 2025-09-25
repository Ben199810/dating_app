package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins     []string      // 允許的來源
	AllowMethods     []string      // 允許的 HTTP 方法
	AllowHeaders     []string      // 允許的標頭
	ExposeHeaders    []string      // 暴露的標頭
	AllowCredentials bool          // 是否允許憑證
	MaxAge           time.Duration // 預檢請求快取時間
	AllowWildcard    bool          // 是否允許通配符來源
	AllowBrowserExt  bool          // 是否允許瀏覽器擴展
	AllowWebSockets  bool          // 是否允許 WebSocket
	AllowFiles       bool          // 是否允許 file:// 協議
}

// DefaultCORSConfig 預設 CORS 配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-Id",
			"X-Response-Time",
		},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    false,
		AllowBrowserExt:  false,
		AllowWebSockets:  true,
		AllowFiles:       false,
	}
}

// ProductionCORSConfig 生產環境 CORS 配置
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-Id",
			"X-New-Token",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
		AllowWildcard:    false,
		AllowBrowserExt:  false,
		AllowWebSockets:  true,
		AllowFiles:       false,
	}
}

// DevelopmentCORSConfig 開發環境 CORS 配置
func DevelopmentCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"*",
		},
		ExposeHeaders: []string{
			"*",
		},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
		AllowWildcard:    true,
		AllowBrowserExt:  true,
		AllowWebSockets:  true,
		AllowFiles:       true,
	}
}

// CORSMiddleware CORS 中間件
type CORSMiddleware struct {
	config CORSConfig
}

// NewCORSMiddleware 建立新的 CORS 中間件
func NewCORSMiddleware(config CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		config: config,
	}
}

// Handler 返回 CORS 中間件處理函數
func (c *CORSMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")

		// 檢查來源是否被允許
		if !c.isOriginAllowed(origin) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "來源不被允許",
				"code":  "CORS_ORIGIN_NOT_ALLOWED",
			})
			return
		}

		// 設置 CORS 標頭
		c.setCORSHeaders(ctx, origin)

		// 處理預檢請求
		if ctx.Request.Method == http.MethodOptions {
			c.handlePreflightRequest(ctx)
			return
		}

		ctx.Next()
	}
}

// isOriginAllowed 檢查來源是否被允許
func (c *CORSMiddleware) isOriginAllowed(origin string) bool {
	// 如果沒有指定來源，允許
	if origin == "" {
		return true
	}

	// 檢查特殊協議
	if !c.config.AllowBrowserExt && (strings.HasPrefix(origin, "moz-extension://") || strings.HasPrefix(origin, "chrome-extension://")) {
		return false
	}

	if !c.config.AllowFiles && strings.HasPrefix(origin, "file://") {
		return false
	}

	// 檢查允許的來源列表
	for _, allowedOrigin := range c.config.AllowOrigins {
		if c.matchOrigin(allowedOrigin, origin) {
			return true
		}
	}

	return false
}

// matchOrigin 匹配來源
func (c *CORSMiddleware) matchOrigin(pattern, origin string) bool {
	if pattern == "*" {
		return c.config.AllowWildcard
	}

	if pattern == origin {
		return true
	}

	// 支援通配符匹配（如 *.example.com）
	if c.config.AllowWildcard && strings.Contains(pattern, "*") {
		return c.matchWildcard(pattern, origin)
	}

	return false
}

// matchWildcard 通配符匹配
func (c *CORSMiddleware) matchWildcard(pattern, origin string) bool {
	// 簡單的通配符匹配實現
	if strings.HasPrefix(pattern, "*.") {
		domain := strings.TrimPrefix(pattern, "*.")
		return strings.HasSuffix(origin, domain)
	}

	return false
}

// setCORSHeaders 設置 CORS 標頭
func (c *CORSMiddleware) setCORSHeaders(ctx *gin.Context, origin string) {
	// Access-Control-Allow-Origin
	if origin != "" && c.isOriginAllowed(origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	} else if len(c.config.AllowOrigins) == 1 && c.config.AllowOrigins[0] == "*" && c.config.AllowWildcard {
		ctx.Header("Access-Control-Allow-Origin", "*")
	}

	// Access-Control-Allow-Credentials
	if c.config.AllowCredentials {
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}

	// Access-Control-Expose-Headers
	if len(c.config.ExposeHeaders) > 0 {
		ctx.Header("Access-Control-Expose-Headers", strings.Join(c.config.ExposeHeaders, ", "))
	}

	// Vary header
	ctx.Header("Vary", "Origin")
}

// handlePreflightRequest 處理預檢請求
func (c *CORSMiddleware) handlePreflightRequest(ctx *gin.Context) {
	// Access-Control-Request-Method
	requestMethod := ctx.Request.Header.Get("Access-Control-Request-Method")
	if requestMethod != "" && !c.isMethodAllowed(requestMethod) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "HTTP 方法不被允許",
			"code":  "CORS_METHOD_NOT_ALLOWED",
		})
		return
	}

	// Access-Control-Request-Headers
	requestHeaders := ctx.Request.Header.Get("Access-Control-Request-Headers")
	if requestHeaders != "" && !c.areHeadersAllowed(requestHeaders) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "請求標頭不被允許",
			"code":  "CORS_HEADERS_NOT_ALLOWED",
		})
		return
	}

	// 設置預檢回應標頭
	if len(c.config.AllowMethods) > 0 {
		ctx.Header("Access-Control-Allow-Methods", strings.Join(c.config.AllowMethods, ", "))
	}

	if len(c.config.AllowHeaders) > 0 {
		ctx.Header("Access-Control-Allow-Headers", strings.Join(c.config.AllowHeaders, ", "))
	}

	if c.config.MaxAge > 0 {
		ctx.Header("Access-Control-Max-Age", strconv.Itoa(int(c.config.MaxAge.Seconds())))
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

// isMethodAllowed 檢查 HTTP 方法是否被允許
func (c *CORSMiddleware) isMethodAllowed(method string) bool {
	for _, allowedMethod := range c.config.AllowMethods {
		if allowedMethod == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed 檢查請求標頭是否被允許
func (c *CORSMiddleware) areHeadersAllowed(headers string) bool {
	// 如果允許所有標頭
	for _, allowedHeader := range c.config.AllowHeaders {
		if allowedHeader == "*" {
			return true
		}
	}

	// 解析請求標頭
	requestHeaders := strings.Split(headers, ",")
	for _, header := range requestHeaders {
		header = strings.TrimSpace(strings.ToLower(header))
		if !c.isHeaderAllowed(header) {
			return false
		}
	}

	return true
}

// isHeaderAllowed 檢查單個標頭是否被允許
func (c *CORSMiddleware) isHeaderAllowed(header string) bool {
	header = strings.ToLower(header)

	for _, allowedHeader := range c.config.AllowHeaders {
		if strings.ToLower(allowedHeader) == header || allowedHeader == "*" {
			return true
		}
	}

	// 始終允許的簡單標頭
	simpleHeaders := []string{
		"cache-control",
		"content-language",
		"content-type",
		"expires",
		"last-modified",
		"pragma",
	}

	for _, simpleHeader := range simpleHeaders {
		if header == simpleHeader {
			return true
		}
	}

	return false
}

// WebSocketCORSMiddleware WebSocket 專用的 CORS 中間件
func (c *CORSMiddleware) WebSocketCORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")

		// 檢查來源
		if !c.isOriginAllowed(origin) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "WebSocket 來源不被允許",
				"code":  "WEBSOCKET_CORS_NOT_ALLOWED",
			})
			return
		}

		// 檢查是否為 WebSocket 升級請求
		if !c.config.AllowWebSockets {
			if ctx.Request.Header.Get("Upgrade") == "websocket" {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "WebSocket 連接不被允許",
					"code":  "WEBSOCKET_NOT_ALLOWED",
				})
				return
			}
		}

		// 設置 CORS 標頭
		c.setCORSHeaders(ctx, origin)

		ctx.Next()
	}
}

// CreateCORSMiddleware 建立適合環境的 CORS 中間件
func CreateCORSMiddleware(env string) *CORSMiddleware {
	var config CORSConfig

	switch env {
	case "production":
		config = ProductionCORSConfig([]string{
			"https://your-domain.com",
			"https://www.your-domain.com",
		})
	case "development":
		config = DevelopmentCORSConfig()
	case "test":
		config = DefaultCORSConfig()
	default:
		config = DefaultCORSConfig()
	}

	return NewCORSMiddleware(config)
}
