package event

import (
	"time"

	"github.com/google/uuid"
)

// MessageSentEvent 訊息發送事件
type MessageSentEvent struct {
	BaseEvent
	UserID  string `json:"user_id"`
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

// NewMessageSentEvent 創建訊息發送事件
func NewMessageSentEvent(userID, roomID, content string) *MessageSentEvent {
	return &MessageSentEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			Type:        "message.sent",
			Timestamp:   time.Now(),
			AggregateID: roomID,
		},
		UserID:  userID,
		RoomID:  roomID,
		Content: content,
	}
}

// UserJoinedEvent 用戶加入事件
type UserJoinedEvent struct {
	BaseEvent
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

// NewUserJoinedEvent 創建用戶加入事件
func NewUserJoinedEvent(userID, roomID string) *UserJoinedEvent {
	return &UserJoinedEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			Type:        "user.joined",
			Timestamp:   time.Now(),
			AggregateID: roomID,
		},
		UserID: userID,
		RoomID: roomID,
	}
}

// UserLeftEvent 用戶離開事件
type UserLeftEvent struct {
	BaseEvent
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

// NewUserLeftEvent 創建用戶離開事件
func NewUserLeftEvent(userID, roomID string) *UserLeftEvent {
	return &UserLeftEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			Type:        "user.left",
			Timestamp:   time.Now(),
			AggregateID: roomID,
		},
		UserID: userID,
		RoomID: roomID,
	}
}
