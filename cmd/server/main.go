package main

import (
	"golang_dev_docker/internal/application/service"
	"golang_dev_docker/internal/domain/chat"
	"golang_dev_docker/internal/domain/user"
	eventBus "golang_dev_docker/internal/infrastructure/event"
	"golang_dev_docker/internal/infrastructure/repository"
	"golang_dev_docker/internal/infrastructure/websocket"
	"golang_dev_docker/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化倉儲
	userRepo := repository.NewInMemoryUserRepository()
	messageRepo := repository.NewInMemoryMessageRepository()
	roomRepo := repository.NewInMemoryChatRoomRepository()

	// 初始化事件總線
	eventBusInstance := eventBus.NewInMemoryEventBus()

	// 初始化 WebSocket 管理器
	wsManager := websocket.NewManager()
	wsEventHandler := websocket.NewWebSocketEventHandler(wsManager)

	// 註冊事件處理器
	eventBusInstance.Subscribe("message.sent", wsEventHandler)
	eventBusInstance.Subscribe("user.joined", wsEventHandler)
	eventBusInstance.Subscribe("user.left", wsEventHandler)

	// 初始化應用服務
	chatService := service.NewChatService(messageRepo, roomRepo, userRepo, eventBusInstance)

	// 初始化處理器
	chatHandler := http.NewChatHandler(chatService)
	userHandler := http.NewUserHandler(userRepo)

	// 創建測試數據
	createTestData(userRepo, roomRepo)

	// 啟動 WebSocket 管理器
	go wsManager.Run()

	// 設置路由
	r := gin.Default()
	setupRoutes(r, chatHandler, userHandler, wsManager)

	// 啟動伺服器
	r.Run(":8080")
}

// createTestData 創建測試數據
func createTestData(userRepo user.UserRepository, roomRepo chat.ChatRoomRepository) {
	// 創建測試用戶
	testUser := user.NewUser("user_1", "Ben", "ben@example.com", []string{"coding", "music", "travel"})
	userRepo.Save(testUser)

	// 創建測試聊天室
	testRoom := chat.NewChatRoom("general", "General chat room")
	testRoom.AddMember(chat.UserID("user_1"))
	roomRepo.Save(testRoom)
}

// setupRoutes 設置路由
func setupRoutes(r *gin.Engine, chatHandler *http.ChatHandler, userHandler *http.UserHandler, wsManager *websocket.Manager) {
	// API 路由
	api := r.Group("/api/v1")
	{
		// 聊天相關
		api.POST("/messages", chatHandler.SendMessage)
		api.GET("/messages", chatHandler.GetMessages)
		api.POST("/rooms/join", chatHandler.JoinRoom)
		api.POST("/rooms/leave", chatHandler.LeaveRoom)

		// 用戶相關
		api.GET("/users/:id", userHandler.GetUser)
		api.GET("/hello", userHandler.Hello)
	}

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		wsManager.HandleWebSocket(c.Writer, c.Request)
	})

	// 靜態檔案
	r.StaticFile("/chat.html", "./web/chat.html")
	r.StaticFile("/chat_websocket.html", "./web/chat_websocket.html")
}
