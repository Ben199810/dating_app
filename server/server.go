package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang_dev_docker/domain/usecase"
	"golang_dev_docker/infrastructure/mysql"
	"golang_dev_docker/infrastructure/redis"
	"golang_dev_docker/server/middleware"
	"golang_dev_docker/server/websocket"

	"github.com/gin-gonic/gin"
)

// WebSocketNotifierAdapter WebSocket 管理器適配器
// 實作 ChatService 所需的 WebSocketNotifier 介面
type WebSocketNotifierAdapter struct {
	manager *websocket.Manager
}

// SendToUser 發送訊息給特定用戶
func (w *WebSocketNotifierAdapter) SendToUser(userID uint, message interface{}) error {
	if w.manager == nil {
		return fmt.Errorf("WebSocket 管理器未初始化")
	}

	// 將訊息發送給指定用戶的所有連接
	// 提取訊息類型，如果沒有則使用預設類型
	messageType := "notification"
	if msgMap, ok := message.(map[string]interface{}); ok {
		if msgType, exists := msgMap["type"]; exists {
			if typeStr, ok := msgType.(string); ok {
				messageType = typeStr
			}
		}
	}

	w.manager.SendToUser(userID, messageType, message)
	return nil
}

// BroadcastToUsers 廣播訊息給多個用戶
func (w *WebSocketNotifierAdapter) BroadcastToUsers(userIDs []uint, message interface{}) error {
	if w.manager == nil {
		return fmt.Errorf("WebSocket 管理器未初始化")
	}

	// 提取訊息類型
	messageType := "notification"
	if msgMap, ok := message.(map[string]interface{}); ok {
		if msgType, exists := msgMap["type"]; exists {
			if typeStr, ok := msgType.(string); ok {
				messageType = typeStr
			}
		}
	}

	// 依次發送給每個用戶
	for _, userID := range userIDs {
		w.manager.SendToUser(userID, messageType, message)
	}
	return nil
}

// ServerConfig 伺服器配置
type ServerConfig struct {
	Port                    int           `yaml:"port"`
	Mode                    string        `yaml:"mode"` // gin 模式: debug, release, test
	ReadTimeout             time.Duration `yaml:"read_timeout"`
	WriteTimeout            time.Duration `yaml:"write_timeout"`
	IdleTimeout             time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes          int           `yaml:"max_header_bytes"`
	GracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout"`
	EnablePprof             bool          `yaml:"enable_pprof"`
	EnableMetrics           bool          `yaml:"enable_metrics"`
	StaticPath              string        `yaml:"static_path"`
	UploadPath              string        `yaml:"upload_path"`
	JWTSecret               string        `yaml:"jwt_secret"`
}

// DefaultServerConfig 預設伺服器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:                    8080,
		Mode:                    gin.DebugMode,
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            10 * time.Second,
		IdleTimeout:             60 * time.Second,
		MaxHeaderBytes:          1 << 20, // 1MB
		GracefulShutdownTimeout: 30 * time.Second,
		EnablePprof:             false,
		EnableMetrics:           false,
		StaticPath:              "./static",
		UploadPath:              "./uploads",
		JWTSecret:               "your-secret-key",
	}
}

// Server 伺服器結構
type Server struct {
	config       *ServerConfig
	engine       *gin.Engine
	httpServer   *http.Server
	dbManager    *mysql.DatabaseManager
	redisClient  *redis.RedisClient
	cacheService *redis.CacheService
	wsManager    *websocket.Manager
	chatHandler  *websocket.ChatHandler

	// 業務服務
	userService     *usecase.UserService
	matchingService *usecase.MatchingService
	chatService     *usecase.ChatService

	// 中間件
	jwtAuth        *middleware.JWTAuthMiddleware
	wsAuth         *middleware.WebSocketAuthMiddleware
	corsMiddleware *middleware.CORSMiddleware
	rateLimiters   map[string]*middleware.RateLimitMiddleware
}

// NewServer 建立新的伺服器實例
func NewServer(config *ServerConfig) *Server {
	// 設置 Gin 模式
	gin.SetMode(config.Mode)

	server := &Server{
		config:       config,
		engine:       gin.New(),
		rateLimiters: make(map[string]*middleware.RateLimitMiddleware),
	}

	return server
}

// InitializeDatabase 初始化資料庫連接
func (s *Server) InitializeDatabase(dbConfig *mysql.DatabaseConfig) error {
	var err error
	s.dbManager, err = mysql.NewDatabaseManager(dbConfig)
	if err != nil {
		return fmt.Errorf("初始化資料庫失敗: %w", err)
	}

	// 執行資料庫遷移和種子資料植入
	if err := s.dbManager.InitializeDatabase(false); err != nil {
		log.Printf("警告：資料庫初始化失敗: %v", err)
		// 不返回錯誤，允許應用程式繼續運行
	}

	log.Println("資料庫連接初始化成功")
	return nil
}

// InitializeRedis 初始化 Redis 連接
func (s *Server) InitializeRedis(redisConfig *redis.RedisConfig) error {
	var err error
	s.redisClient, err = redis.NewRedisClient(redisConfig)
	if err != nil {
		return fmt.Errorf("初始化 Redis 失敗: %w", err)
	}

	// 初始化快取服務
	s.cacheService = redis.NewCacheService(s.redisClient, "dating_app")

	log.Println("Redis 連接初始化成功")
	return nil
}

// InitializeWebSocket 初始化 WebSocket 管理器
func (s *Server) InitializeWebSocket() {
	s.wsManager = websocket.NewManager()
	s.chatHandler = websocket.NewChatHandler(s.wsManager)

	// 啟動 WebSocket 管理器
	go s.wsManager.Run()

	log.Println("WebSocket 管理器初始化成功")
}

// InitializeServices 初始化業務服務
func (s *Server) InitializeServices() error {
	if s.dbManager == nil {
		return fmt.Errorf("資料庫管理器未初始化")
	}

	// 獲取資料庫連接
	db := s.dbManager.GetDB()
	if db == nil {
		return fmt.Errorf("無法獲取資料庫連接")
	}

	// 創建儲存庫實例
	userRepo := mysql.NewUserRepository(db)
	userProfileRepo := mysql.NewUserProfileRepository(db)
	matchRepo := mysql.NewMatchRepository(db)
	chatRepo := mysql.NewChatRepository(db)
	chatListRepo := mysql.NewChatListRepository(db)
	websocketRepo := mysql.NewWebSocketRepository(db)

	// 創建 Redis 快取服務（如果可用）
	var matchingCache *redis.MatchingCacheService

	if s.redisClient != nil {
		// 初始化會話快取服務但暫時不存儲引用（將在後續整合到認證中間件）
		_ = redis.NewSessionCacheService(s.redisClient)

		// 初始化配對快取服務
		matchingCache = redis.NewMatchingCacheService(s.redisClient)

		log.Println("Redis 快取服務初始化成功")
	}

	// 初始化用戶服務
	// TODO: 需要添加缺少的儲存庫（PhotoRepository, InterestRepository, AgeVerificationRepository）
	s.userService = usecase.NewUserService(
		userRepo,
		userProfileRepo,
		nil, // photoRepo - 需要實作
		nil, // interestRepo - 需要實作
		nil, // ageVerificationRepo - 需要實作
	)

	// 初始化配對服務
	// TODO: 需要添加 MatchingAlgorithmRepository
	s.matchingService = usecase.NewMatchingService(
		matchRepo,
		nil, // algorithmRepo - 需要實作
		userRepo,
		userProfileRepo,
	)

	// 設定配對快取（如果可用）
	if matchingCache != nil {
		s.matchingService.SetCache(matchingCache)
		log.Println("配對服務快取整合完成")
	}

	// 初始化聊天服務
	s.chatService = usecase.NewChatService(
		chatRepo,
		chatListRepo,
		websocketRepo,
		matchRepo,
		userRepo,
	)

	// 整合 WebSocket 通知器（需要在 WebSocket 管理器初始化後設定）
	if s.wsManager != nil {
		wsNotifier := &WebSocketNotifierAdapter{manager: s.wsManager}
		s.chatService.SetWebSocketNotifier(wsNotifier)
		log.Println("聊天服務 WebSocket 通知整合完成")
	}

	log.Println("業務服務初始化成功")
	return nil
}

// InitializeMiddleware 初始化中間件
func (s *Server) InitializeMiddleware(env string) {
	// JWT 認證中間件
	s.jwtAuth = middleware.NewJWTAuthMiddleware(s.config.JWTSecret)

	// WebSocket 認證中間件
	s.wsAuth = middleware.NewWebSocketAuthMiddleware(s.config.JWTSecret)

	// CORS 中間件
	s.corsMiddleware = middleware.CreateCORSMiddleware(env)

	// 速率限制中間件
	s.rateLimiters["api"] = middleware.CreateAPIRateLimiter()
	s.rateLimiters["login"] = middleware.CreateLoginRateLimiter()
	s.rateLimiters["message"] = middleware.CreateMessageRateLimiter()
	s.rateLimiters["swipe"] = middleware.CreateSwipeRateLimiter()
	s.rateLimiters["photo"] = middleware.CreatePhotoUploadRateLimiter()
	s.rateLimiters["register"] = middleware.CreateRegistrationRateLimiter()

	log.Println("中間件初始化成功")
}

// SetupRoutes 設置路由
func (s *Server) SetupRoutes() {
	// 全域中間件
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())
	s.engine.Use(s.corsMiddleware.Handler())

	// 靜態檔案服務
	s.engine.Static("/static", s.config.StaticPath)
	s.engine.Static("/uploads", s.config.UploadPath)

	// 首頁路由 - 重導向到 static/index.html
	s.engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// 健康檢查端點
	s.engine.GET("/health", s.healthCheckHandler)
	s.engine.GET("/api/health", s.apiHealthCheckHandler)

	// WebSocket 端點
	wsGroup := s.engine.Group("/ws")
	wsGroup.Use(s.wsAuth.CreateAuthMiddlewareChain()...)
	wsGroup.GET("/chat", s.wsManager.HandleWebSocket)

	// API 路由組
	apiGroup := s.engine.Group("/api")
	apiGroup.Use(s.rateLimiters["api"].Handler())

	// 認證路由（不需要 JWT 認證）
	authGroup := apiGroup.Group("/auth")
	authGroup.Use(s.rateLimiters["login"].Handler())
	{
		// TODO: 添加認證相關的路由
		// authGroup.POST("/register", registerHandler)
		// authGroup.POST("/login", loginHandler)
		// authGroup.POST("/logout", logoutHandler)
	}

	// 需要認證的路由
	protectedGroup := apiGroup.Group("")
	protectedGroup.Use(s.jwtAuth.AuthMiddleware())
	{
		// 用戶相關路由
		userGroup := protectedGroup.Group("/user")
		_ = userGroup // TODO: 添加用戶相關的路由
		{
			// TODO: 添加用戶相關的路由
			// userGroup.GET("/profile", getUserProfileHandler)
			// userGroup.PUT("/profile", updateUserProfileHandler)
		}

		// 配對相關路由
		matchGroup := protectedGroup.Group("/match")
		matchGroup.Use(s.rateLimiters["swipe"].Handler())
		_ = matchGroup // TODO: 添加配對相關的路由
		{
			// TODO: 添加配對相關的路由
			// matchGroup.GET("/candidates", getMatchCandidatesHandler)
			// matchGroup.POST("/swipe", swipeHandler)
		}

		// 聊天相關路由
		chatGroup := protectedGroup.Group("/chat")
		chatGroup.Use(s.rateLimiters["message"].Handler())
		_ = chatGroup // TODO: 添加聊天相關的路由
		{
			// TODO: 添加聊天相關的路由
			// chatGroup.GET("/list", getChatListHandler)
			// chatGroup.GET("/:chat_id/messages", getChatMessagesHandler)
			// chatGroup.POST("/:chat_id/messages", sendMessageHandler)
		}

		// 照片上傳路由
		photoGroup := protectedGroup.Group("/photos")
		photoGroup.Use(s.rateLimiters["photo"].Handler())
		_ = photoGroup // TODO: 添加照片相關的路由
		{
			// TODO: 添加照片相關的路由
			// photoGroup.POST("/upload", uploadPhotoHandler)
			// photoGroup.DELETE("/:photo_id", deletePhotoHandler)
		}
	}

	log.Println("路由設置完成")
}

// healthCheckHandler 基本健康檢查
func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "dating-app",
	})
}

// apiHealthCheckHandler API 健康檢查（包含依賴服務）
func (s *Server) apiHealthCheckHandler(c *gin.Context) {
	health := gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "dating-app-api",
		"version":   "1.0.0",
		"checks":    gin.H{},
	}

	// 檢查資料庫連接
	if s.dbManager != nil {
		if err := s.dbManager.Health(); err != nil {
			health["status"] = "degraded"
			health["checks"].(gin.H)["database"] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			health["checks"].(gin.H)["database"] = gin.H{
				"status": "healthy",
				"info":   s.dbManager.GetConnectionInfo(),
			}
		}
	}

	// 檢查 Redis 連接
	if s.cacheService != nil {
		if err := s.cacheService.Health(); err != nil {
			health["status"] = "degraded"
			health["checks"].(gin.H)["redis"] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			health["checks"].(gin.H)["redis"] = gin.H{
				"status": "healthy",
			}
		}
	}

	// 檢查 WebSocket 管理器
	if s.wsManager != nil {
		onlineUsers := s.wsManager.GetOnlineUsers()
		health["checks"].(gin.H)["websocket"] = gin.H{
			"status":       "healthy",
			"online_users": len(onlineUsers),
		}
	}

	c.JSON(http.StatusOK, health)
}

// Start 啟動伺服器
func (s *Server) Start() error {
	// 創建 HTTP 伺服器
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.config.Port),
		Handler:        s.engine,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}

	log.Printf("伺服器啟動在端口 %d", s.config.Port)

	// 啟動伺服器
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("伺服器啟動失敗: %w", err)
	}

	return nil
}

// StartWithGracefulShutdown 啟動伺服器並支援優雅關閉
func (s *Server) StartWithGracefulShutdown() error {
	// 啟動伺服器
	go func() {
		if err := s.Start(); err != nil {
			log.Printf("伺服器啟動錯誤: %v", err)
		}
	}()

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在關閉伺服器...")

	// 優雅關閉
	return s.Shutdown()
}

// Shutdown 優雅關閉伺服器
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.GracefulShutdownTimeout)
	defer cancel()

	// 關閉 HTTP 伺服器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("伺服器關閉錯誤: %v", err)
	}

	// 關閉 WebSocket 管理器
	if s.wsManager != nil {
		s.wsManager.Shutdown()
	}

	// 關閉資料庫連接
	if s.dbManager != nil {
		if err := s.dbManager.Close(); err != nil {
			log.Printf("資料庫關閉錯誤: %v", err)
		}
	}

	// 關閉 Redis 連接
	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			log.Printf("Redis 關閉錯誤: %v", err)
		}
	}

	log.Println("伺服器已優雅關閉")
	return nil
}

// GetEngine 獲取 Gin 引擎（用於測試）
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// GetDBManager 獲取資料庫管理器
func (s *Server) GetDBManager() *mysql.DatabaseManager {
	return s.dbManager
}

// GetCacheService 獲取快取服務
func (s *Server) GetCacheService() *redis.CacheService {
	return s.cacheService
}

// GetWebSocketManager 獲取 WebSocket 管理器
func (s *Server) GetWebSocketManager() *websocket.Manager {
	return s.wsManager
}

// GetChatHandler 獲取聊天處理器
func (s *Server) GetChatHandler() *websocket.ChatHandler {
	return s.chatHandler
}
