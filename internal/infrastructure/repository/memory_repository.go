package repository

import (
	"errors"
	"golang_dev_docker/internal/domain/chat"
	"sync"
)

var (
	ErrMessageNotFound = errors.New("message not found")
	ErrRoomNotFound    = errors.New("room not found")
	ErrUserNotFound    = errors.New("user not found")
)

// InMemoryMessageRepository 記憶體訊息倉儲實現
type InMemoryMessageRepository struct {
	messages map[chat.MessageID]*chat.Message
	mu       sync.RWMutex
}

// NewInMemoryMessageRepository 創建記憶體訊息倉儲
func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{
		messages: make(map[chat.MessageID]*chat.Message),
	}
}

// Save 儲存訊息
func (r *InMemoryMessageRepository) Save(message *chat.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[message.ID] = message
	return nil
}

// FindByRoomID 根據聊天室ID查找訊息
func (r *InMemoryMessageRepository) FindByRoomID(roomID chat.RoomID) ([]*chat.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var messages []*chat.Message
	for _, message := range r.messages {
		if message.RoomID == roomID {
			messages = append(messages, message)
		}
	}
	return messages, nil
}

// FindByID 根據ID查找訊息
func (r *InMemoryMessageRepository) FindByID(id chat.MessageID) (*chat.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if message, exists := r.messages[id]; exists {
		return message, nil
	}
	return nil, ErrMessageNotFound
}

// InMemoryChatRoomRepository 記憶體聊天室倉儲實現
type InMemoryChatRoomRepository struct {
	rooms map[chat.RoomID]*chat.ChatRoom
	mu    sync.RWMutex
}

// NewInMemoryChatRoomRepository 創建記憶體聊天室倉儲
func NewInMemoryChatRoomRepository() *InMemoryChatRoomRepository {
	return &InMemoryChatRoomRepository{
		rooms: make(map[chat.RoomID]*chat.ChatRoom),
	}
}

// Save 儲存聊天室
func (r *InMemoryChatRoomRepository) Save(room *chat.ChatRoom) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms[room.ID] = room
	return nil
}

// FindByID 根據ID查找聊天室
func (r *InMemoryChatRoomRepository) FindByID(id chat.RoomID) (*chat.ChatRoom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if room, exists := r.rooms[id]; exists {
		return room, nil
	}
	return nil, ErrRoomNotFound
}

// FindAll 查找所有聊天室
func (r *InMemoryChatRoomRepository) FindAll() ([]*chat.ChatRoom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var rooms []*chat.ChatRoom
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}
	return rooms, nil
}
