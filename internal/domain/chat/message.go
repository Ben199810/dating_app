package chat

import (
	"time"

	"github.com/google/uuid"
)

// Message 代表聊天訊息的領域實體
type Message struct {
	ID        MessageID `json:"id"`
	UserID    UserID    `json:"user_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	RoomID    RoomID    `json:"room_id"`
}

// MessageID 訊息唯一識別符
type MessageID string

// UserID 用戶唯一識別符
type UserID string

// RoomID 聊天室唯一識別符
type RoomID string

// NewMessage 創建新訊息
func NewMessage(userID UserID, content string, roomID RoomID) *Message {
	return &Message{
		ID:        MessageID(uuid.New().String()),
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
		RoomID:    roomID,
	}
}

// IsValid 驗證訊息是否有效
func (m *Message) IsValid() bool {
	return m.UserID != "" && m.Content != "" && m.RoomID != ""
}

// MessageRepository 訊息倉儲介面
type MessageRepository interface {
	Save(message *Message) error
	FindByRoomID(roomID RoomID) ([]*Message, error)
	FindByID(id MessageID) (*Message, error)
}
