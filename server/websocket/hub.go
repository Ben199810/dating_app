package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message 表示 WebSocket 訊息
type Message struct {
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
	UserID   uint        `json:"user_id,omitempty"`
	TargetID uint        `json:"target_id,omitempty"`
}

// Client 表示 WebSocket 客戶端
type Client struct {
	ID     string          `json:"id"`
	UserID uint            `json:"user_id"`
	Conn   *websocket.Conn `json:"-"`
	Send   chan Message    `json:"-"`
	Hub    *Hub            `json:"-"`
}

// Hub 管理所有 WebSocket 客戶端
type Hub struct {
	// 已註冊的客戶端
	clients map[*Client]bool

	// 廣播通道
	broadcast chan Message

	// 客戶端註冊請求
	register chan *Client

	// 客戶端取消註冊請求
	unregister chan *Client

	// 用戶 ID 到客戶端的映射
	userClients map[uint]*Client

	// 互斥鎖保護併發存取
	mutex sync.RWMutex
}

// WebSocket 升級器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 在生產環境中，應該檢查來源
		return true
	},
}

// NewHub 建立新的 WebSocket Hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		userClients: make(map[uint]*Client),
	}
}

// Run 啟動 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 註冊新客戶端
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true
	h.userClients[client.UserID] = client

	log.Printf("WebSocket 客戶端已連接: UserID=%d, ClientID=%s", client.UserID, client.ID)

	// 發送連接成功訊息
	client.Send <- Message{
		Type: "connection",
		Data: map[string]interface{}{
			"status":  "connected",
			"message": "WebSocket 連接成功",
		},
	}
}

// unregisterClient 取消註冊客戶端
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		delete(h.userClients, client.UserID)
		close(client.Send)

		log.Printf("WebSocket 客戶端已斷線: UserID=%d, ClientID=%s", client.UserID, client.ID)
	}
}

// broadcastMessage 廣播訊息
func (h *Hub) broadcastMessage(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// 如果指定了目標用戶，只發送給該用戶
	if message.TargetID != 0 {
		if client, ok := h.userClients[message.TargetID]; ok {
			select {
			case client.Send <- message:
			default:
				// 客戶端緩衝區滿，關閉連接
				h.unregisterClient(client)
			}
		}
		return
	}

	// 廣播給所有客戶端
	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			// 客戶端緩衝區滿，關閉連接
			h.unregisterClient(client)
		}
	}
}

// SendToUser 發送訊息給特定用戶
func (h *Hub) SendToUser(userID uint, message Message) {
	message.TargetID = userID
	h.broadcast <- message
}

// BroadcastToAll 廣播訊息給所有用戶
func (h *Hub) BroadcastToAll(message Message) {
	h.broadcast <- message
}

// GetOnlineUsers 取得線上用戶 ID 列表
func (h *Hub) GetOnlineUsers() []uint {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	users := make([]uint, 0, len(h.userClients))
	for userID := range h.userClients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline 檢查用戶是否線上
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, online := h.userClients[userID]
	return online
}

// HandleWebSocket 處理 WebSocket 連接
func (h *Hub) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升級失敗: %v", err)
		return
	}

	// 從查詢參數或 JWT token 取得用戶 ID
	// 這裡暫時使用查詢參數，實際應該驗證 JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		log.Println("WebSocket 連接缺少 user_id 參數")
		conn.Close()
		return
	}

	var userID uint
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		log.Printf("無效的 user_id 參數: %s", userIDStr)
		conn.Close()
		return
	}

	// 產生客戶端 ID
	clientID := fmt.Sprintf("client_%d_%d", userID, time.Now().Unix())

	client := &Client{
		ID:     clientID,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan Message, 256),
		Hub:    h,
	}

	// 註冊客戶端
	h.register <- client

	// 啟動讀寫協程
	go client.writePump()
	go client.readPump()
}
