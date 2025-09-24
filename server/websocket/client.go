package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 寫入等待時間
	writeWait = 10 * time.Second

	// Pong 等待時間
	pongWait = 60 * time.Second

	// Ping 間隔 (必須小於 pongWait)
	pingPeriod = (pongWait * 9) / 10

	// 最大訊息大小
	maxMessageSize = 512
)

// readPump 處理從客戶端讀取訊息
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var message Message
		if err := c.Conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 讀取錯誤: %v", err)
			}
			break
		}

		// 設置發送者 ID
		message.UserID = c.UserID

		// 處理不同類型的訊息
		c.handleMessage(message)
	}
}

// writePump 處理向客戶端寫入訊息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 關閉了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket 寫入錯誤: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 處理收到的訊息
func (c *Client) handleMessage(message Message) {
	switch message.Type {
	case "chat":
		// 轉發聊天訊息給指定用戶
		if message.TargetID != 0 {
			c.Hub.SendToUser(message.TargetID, message)
		}

	case "typing":
		// 轉發打字狀態給指定用戶
		if message.TargetID != 0 {
			c.Hub.SendToUser(message.TargetID, message)
		}

	case "ping":
		// 回應 ping
		c.Send <- Message{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": time.Now().Unix(),
			},
		}

	default:
		log.Printf("未知的訊息類型: %s, 來自用戶: %d", message.Type, c.UserID)
	}
}

// SendMessage 發送訊息給此客戶端
func (c *Client) SendMessage(message Message) {
	select {
	case c.Send <- message:
	default:
		// 客戶端緩衝區滿，關閉連接
		close(c.Send)
	}
}
