package websocket

import (
	"context"
	"encoding/json"
	"golang_dev_docker/internal/application/event"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Manager WebSocket 連線管理器
type Manager struct {
	clients    map[*websocket.Conn]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Client WebSocket 客戶端
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	userID string
	roomID string
}

// NewManager 創建 WebSocket 管理器
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*websocket.Conn]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 運行 WebSocket 管理器
func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.conn] = client
			m.mu.Unlock()

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client.conn]; ok {
				delete(m.clients, client.conn)
				close(client.send)
			}
			m.mu.Unlock()

		case message := <-m.broadcast:
			m.mu.RLock()
			for conn, client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, conn)
				}
			}
			m.mu.RUnlock()
		}
	}
}

// HandleWebSocket 處理 WebSocket 連線
func (m *Manager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 開發環境，生產環境需要更嚴格的檢查
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 從查詢參數獲取用戶ID和房間ID
	userID := r.URL.Query().Get("user_id")
	roomID := r.URL.Query().Get("room_id")

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		roomID: roomID,
	}

	m.register <- client

	go client.writePump()
	go client.readPump(m)
}

// readPump 讀取客戶端訊息
func (c *Client) readPump(manager *Manager) {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		manager.broadcast <- message
	}
}

// writePump 發送訊息給客戶端
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// WebSocketEventHandler WebSocket 事件處理器
type WebSocketEventHandler struct {
	manager *Manager
}

// NewWebSocketEventHandler 創建 WebSocket 事件處理器
func NewWebSocketEventHandler(manager *Manager) *WebSocketEventHandler {
	return &WebSocketEventHandler{
		manager: manager,
	}
}

// Handle 處理事件
func (h *WebSocketEventHandler) Handle(ctx context.Context, evt event.Event) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	h.manager.broadcast <- data
	return nil
}

// CanHandle 檢查是否能處理該事件類型
func (h *WebSocketEventHandler) CanHandle(eventType string) bool {
	return eventType == "message.sent" ||
		eventType == "user.joined" ||
		eventType == "user.left"
}
