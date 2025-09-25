package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 升級器配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 在生產環境中應該檢查來源
		return true
	},
}

// Client 代表一個 WebSocket 客戶端連接
type Client struct {
	ID       string          // 客戶端唯一標識符
	UserID   uint            // 用戶ID
	Conn     *websocket.Conn // WebSocket 連接
	Send     chan []byte     // 發送緩衝通道
	Manager  *Manager        // 管理器引用
	LastPong time.Time       // 最後心跳時間
}

// Message 代表 WebSocket 訊息
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Manager 管理所有 WebSocket 連接
type Manager struct {
	// 已註冊的客戶端
	clients map[*Client]bool

	// 用戶ID到客戶端的映射
	userClients map[uint]*Client

	// 聊天室到客戶端的映射
	chatRooms map[uint][]*Client

	// 註冊請求來自客戶端
	register chan *Client

	// 註銷請求來自客戶端
	unregister chan *Client

	// 廣播訊息到所有客戶端
	broadcast chan []byte

	// 發送訊息到特定用戶
	sendToUser chan *UserMessage

	// 發送訊息到聊天室
	sendToChat chan *ChatMessage

	// 上下文用於優雅關閉
	ctx    context.Context
	cancel context.CancelFunc

	// 讀寫鎖
	mu sync.RWMutex
}

// UserMessage 用戶訊息結構
type UserMessage struct {
	UserID  uint   `json:"user_id"`
	Message []byte `json:"message"`
}

// ChatMessage 聊天室訊息結構
type ChatMessage struct {
	ChatID  uint   `json:"chat_id"`
	Message []byte `json:"message"`
}

// NewManager 建立新的 WebSocket 管理器
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		clients:     make(map[*Client]bool),
		userClients: make(map[uint]*Client),
		chatRooms:   make(map[uint][]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		sendToUser:  make(chan *UserMessage),
		sendToChat:  make(chan *ChatMessage),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Run 啟動 WebSocket 管理器
func (m *Manager) Run() {
	log.Println("WebSocket 管理器已啟動")

	for {
		select {
		case client := <-m.register:
			m.handleRegister(client)

		case client := <-m.unregister:
			m.handleUnregister(client)

		case message := <-m.broadcast:
			m.handleBroadcast(message)

		case userMsg := <-m.sendToUser:
			m.handleSendToUser(userMsg)

		case chatMsg := <-m.sendToChat:
			m.handleSendToChat(chatMsg)

		case <-m.ctx.Done():
			log.Println("WebSocket 管理器正在關閉...")
			m.closeAllConnections()
			return
		}
	}
}

// HandleWebSocket 處理 WebSocket 連接請求
func (m *Manager) HandleWebSocket(c *gin.Context) {
	// 從查詢參數或標頭獲取用戶ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要認證"})
		return
	}

	// 升級 HTTP 連接到 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升級失敗: %v", err)
		return
	}

	// 建立客戶端
	client := &Client{
		ID:       generateClientID(),
		UserID:   userID.(uint),
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Manager:  m,
		LastPong: time.Now(),
	}

	// 註冊客戶端
	m.register <- client

	// 啟動客戶端的讀寫協程
	go client.writePump()
	go client.readPump()
}

// handleRegister 處理客戶端註冊
func (m *Manager) handleRegister(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果用戶已有連接，關閉舊連接
	if oldClient, exists := m.userClients[client.UserID]; exists {
		m.disconnectClient(oldClient)
	}

	// 註冊新客戶端
	m.clients[client] = true
	m.userClients[client.UserID] = client

	log.Printf("用戶 %d 已連接 WebSocket (客戶端ID: %s)", client.UserID, client.ID)

	// 發送連接成功確認
	welcomeMsg := Message{
		Type: "connected",
		Data: map[string]interface{}{
			"client_id": client.ID,
			"user_id":   client.UserID,
		},
	}

	if data, err := json.Marshal(welcomeMsg); err == nil {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
		}
	}

	// 通知其他相關用戶該用戶已上線
	m.broadcastUserOnlineStatus(client.UserID, true)
}

// handleUnregister 處理客戶端註銷
func (m *Manager) handleUnregister(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[client]; ok {
		m.disconnectClient(client)
		log.Printf("用戶 %d 已斷開 WebSocket 連接", client.UserID)

		// 通知其他相關用戶該用戶已離線
		m.broadcastUserOnlineStatus(client.UserID, false)
	}
}

// disconnectClient 斷開客戶端連接
func (m *Manager) disconnectClient(client *Client) {
	// 從所有映射中移除客戶端
	delete(m.clients, client)
	delete(m.userClients, client.UserID)

	// 從所有聊天室中移除客戶端
	for chatID, clients := range m.chatRooms {
		for i, c := range clients {
			if c == client {
				m.chatRooms[chatID] = append(clients[:i], clients[i+1:]...)
				if len(m.chatRooms[chatID]) == 0 {
					delete(m.chatRooms, chatID)
				}
				break
			}
		}
	}

	// 關閉發送通道和連接
	close(client.Send)
	client.Conn.Close()
}

// handleBroadcast 處理廣播訊息
func (m *Manager) handleBroadcast(message []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		select {
		case client.Send <- message:
		default:
			// 發送失敗，關閉連接
			go func(c *Client) {
				m.unregister <- c
			}(client)
		}
	}
}

// handleSendToUser 處理發送訊息給特定用戶
func (m *Manager) handleSendToUser(userMsg *UserMessage) {
	m.mu.RLock()
	client, exists := m.userClients[userMsg.UserID]
	m.mu.RUnlock()

	if !exists {
		log.Printf("用戶 %d 未在線，無法發送訊息", userMsg.UserID)
		return
	}

	select {
	case client.Send <- userMsg.Message:
	default:
		// 發送失敗，關閉連接
		go func() {
			m.unregister <- client
		}()
	}
}

// handleSendToChat 處理發送訊息到聊天室
func (m *Manager) handleSendToChat(chatMsg *ChatMessage) {
	m.mu.RLock()
	clients, exists := m.chatRooms[chatMsg.ChatID]
	m.mu.RUnlock()

	if !exists {
		return
	}

	for _, client := range clients {
		select {
		case client.Send <- chatMsg.Message:
		default:
			// 發送失敗，關閉連接
			go func(c *Client) {
				m.unregister <- c
			}(client)
		}
	}
}

// broadcastUserOnlineStatus 廣播用戶在線狀態
func (m *Manager) broadcastUserOnlineStatus(userID uint, isOnline bool) {
	statusMsg := Message{
		Type: "user_online_status",
		Data: map[string]interface{}{
			"user_id":   userID,
			"is_online": isOnline,
		},
	}

	if data, err := json.Marshal(statusMsg); err == nil {
		// 這裡應該只發送給相關的用戶（如配對的用戶）
		// 暫時廣播給所有用戶
		go func() {
			m.broadcast <- data
		}()
	}
}

// SendToUser 發送訊息給特定用戶
func (m *Manager) SendToUser(userID uint, messageType string, data interface{}) {
	msg := Message{
		Type: messageType,
		Data: data,
	}

	if msgData, err := json.Marshal(msg); err == nil {
		userMsg := &UserMessage{
			UserID:  userID,
			Message: msgData,
		}

		select {
		case m.sendToUser <- userMsg:
		default:
			log.Printf("發送給用戶 %d 的訊息隊列已滿", userID)
		}
	}
}

// SendToChat 發送訊息到聊天室
func (m *Manager) SendToChat(chatID uint, messageType string, data interface{}) {
	msg := Message{
		Type: messageType,
		Data: data,
	}

	if msgData, err := json.Marshal(msg); err == nil {
		chatMsg := &ChatMessage{
			ChatID:  chatID,
			Message: msgData,
		}

		select {
		case m.sendToChat <- chatMsg:
		default:
			log.Printf("發送到聊天室 %d 的訊息隊列已滿", chatID)
		}
	}
}

// JoinChat 將用戶加入聊天室
func (m *Manager) JoinChat(userID, chatID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.userClients[userID]
	if !exists {
		return
	}

	// 檢查是否已在聊天室中
	for _, c := range m.chatRooms[chatID] {
		if c == client {
			return // 已在聊天室中
		}
	}

	// 加入聊天室
	m.chatRooms[chatID] = append(m.chatRooms[chatID], client)
	log.Printf("用戶 %d 加入聊天室 %d", userID, chatID)
}

// LeaveChat 將用戶從聊天室移除
func (m *Manager) LeaveChat(userID, chatID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.userClients[userID]
	if !exists {
		return
	}

	clients := m.chatRooms[chatID]
	for i, c := range clients {
		if c == client {
			m.chatRooms[chatID] = append(clients[:i], clients[i+1:]...)
			if len(m.chatRooms[chatID]) == 0 {
				delete(m.chatRooms, chatID)
			}
			log.Printf("用戶 %d 離開聊天室 %d", userID, chatID)
			break
		}
	}
}

// IsUserOnline 檢查用戶是否在線
func (m *Manager) IsUserOnline(userID uint) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.userClients[userID]
	return exists
}

// GetOnlineUsers 獲取在線用戶列表
func (m *Manager) GetOnlineUsers() []uint {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]uint, 0, len(m.userClients))
	for userID := range m.userClients {
		users = append(users, userID)
	}

	return users
}

// Shutdown 優雅關閉管理器
func (m *Manager) Shutdown() {
	log.Println("正在關閉 WebSocket 管理器...")
	m.cancel()
}

// closeAllConnections 關閉所有連接
func (m *Manager) closeAllConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for client := range m.clients {
		close(client.Send)
		client.Conn.Close()
	}

	// 清空所有映射
	m.clients = make(map[*Client]bool)
	m.userClients = make(map[uint]*Client)
	m.chatRooms = make(map[uint][]*Client)
}

// 生成客戶端ID
func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

// readPump 處理從客戶端讀取訊息
func (c *Client) readPump() {
	defer func() {
		c.Manager.unregister <- c
		c.Conn.Close()
	}()

	// 設定讀取限制和超時
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.LastPong = time.Now()
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 讀取錯誤: %v", err)
			}
			break
		}

		// 處理收到的訊息
		c.handleMessage(message)
	}
}

// writePump 處理向客戶端寫入訊息
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket 寫入錯誤: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 處理客戶端發送的訊息
func (c *Client) handleMessage(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("解析訊息失敗: %v", err)
		return
	}

	switch msg.Type {
	case "ping":
		// 回應 pong
		pongMsg := Message{Type: "pong", Data: nil}
		if data, err := json.Marshal(pongMsg); err == nil {
			select {
			case c.Send <- data:
			default:
			}
		}

	case "join_chat":
		// 加入聊天室
		if chatData, ok := msg.Data.(map[string]interface{}); ok {
			if chatID, ok := chatData["chat_id"].(float64); ok {
				c.Manager.JoinChat(c.UserID, uint(chatID))
			}
		}

	case "leave_chat":
		// 離開聊天室
		if chatData, ok := msg.Data.(map[string]interface{}); ok {
			if chatID, ok := chatData["chat_id"].(float64); ok {
				c.Manager.LeaveChat(c.UserID, uint(chatID))
			}
		}

	case "typing_start", "typing_stop", "send_message", "message_read":
		// 這些訊息類型已移動到 chat_handler.go 中處理
		log.Printf("收到用戶 %d 的 %s 訊息，應通過 ChatHandler 處理", c.UserID, msg.Type)

	default:
		log.Printf("未知的訊息類型: %s", msg.Type)
	}
}
