package service

import (
	"context"
	"golang_dev_docker/internal/application/event"
	"golang_dev_docker/internal/domain/chat"
	"golang_dev_docker/internal/domain/user"
)

// ChatService 聊天應用服務
type ChatService struct {
	messageRepo chat.MessageRepository
	roomRepo    chat.ChatRoomRepository
	userRepo    user.UserRepository
	eventBus    event.EventBus
}

// NewChatService 創建聊天服務
func NewChatService(
	messageRepo chat.MessageRepository,
	roomRepo chat.ChatRoomRepository,
	userRepo user.UserRepository,
	eventBus event.EventBus,
) *ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		eventBus:    eventBus,
	}
}

// SendMessage 發送訊息
func (s *ChatService) SendMessage(ctx context.Context, userID string, roomID string, content string) (*chat.Message, error) {
	// 驗證用戶存在
	_, err := s.userRepo.FindByID(user.UserID(userID))
	if err != nil {
		return nil, err
	}

	// 驗證聊天室存在
	room, err := s.roomRepo.FindByID(chat.RoomID(roomID))
	if err != nil {
		return nil, err
	}

	// 檢查用戶是否為聊天室成員
	if !room.IsMember(chat.UserID(userID)) {
		return nil, ErrUserNotInRoom
	}

	// 創建訊息
	message := chat.NewMessage(chat.UserID(userID), content, chat.RoomID(roomID))

	// 儲存訊息
	if err := s.messageRepo.Save(message); err != nil {
		return nil, err
	}

	// 發布事件
	messageSentEvent := event.NewMessageSentEvent(userID, roomID, content)
	if err := s.eventBus.Publish(ctx, messageSentEvent); err != nil {
		// 記錄錯誤但不阻止流程
		// TODO: 實現日誌系統
	}

	return message, nil
}

// GetMessages 獲取聊天室訊息
func (s *ChatService) GetMessages(ctx context.Context, roomID string) ([]*chat.Message, error) {
	return s.messageRepo.FindByRoomID(chat.RoomID(roomID))
}

// JoinRoom 用戶加入聊天室
func (s *ChatService) JoinRoom(ctx context.Context, userID string, roomID string) error {
	// 驗證用戶存在
	_, err := s.userRepo.FindByID(user.UserID(userID))
	if err != nil {
		return err
	}

	// 獲取聊天室
	room, err := s.roomRepo.FindByID(chat.RoomID(roomID))
	if err != nil {
		return err
	}

	// 添加用戶到聊天室
	room.AddMember(chat.UserID(userID))

	// 儲存聊天室
	if err := s.roomRepo.Save(room); err != nil {
		return err
	}

	// 發布事件
	userJoinedEvent := event.NewUserJoinedEvent(userID, roomID)
	if err := s.eventBus.Publish(ctx, userJoinedEvent); err != nil {
		// 記錄錯誤但不阻止流程
	}

	return nil
}

// LeaveRoom 用戶離開聊天室
func (s *ChatService) LeaveRoom(ctx context.Context, userID string, roomID string) error {
	// 獲取聊天室
	room, err := s.roomRepo.FindByID(chat.RoomID(roomID))
	if err != nil {
		return err
	}

	// 從聊天室移除用戶
	room.RemoveMember(chat.UserID(userID))

	// 儲存聊天室
	if err := s.roomRepo.Save(room); err != nil {
		return err
	}

	// 發布事件
	userLeftEvent := event.NewUserLeftEvent(userID, roomID)
	if err := s.eventBus.Publish(ctx, userLeftEvent); err != nil {
		// 記錄錯誤但不阻止流程
	}

	return nil
}
