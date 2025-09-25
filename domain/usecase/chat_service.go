package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

// ChatService 聊天業務邏輯服務
// 負責聊天訊息處理、歷史管理、WebSocket連接管理等核心業務邏輯
type ChatService struct {
	chatRepo      repository.ChatRepository
	chatListRepo  repository.ChatListRepository
	websocketRepo repository.WebSocketRepository
	matchRepo     repository.MatchRepository
	userRepo      repository.UserRepository
}

// NewChatService 創建新的聊天服務實例
func NewChatService(
	chatRepo repository.ChatRepository,
	chatListRepo repository.ChatListRepository,
	websocketRepo repository.WebSocketRepository,
	matchRepo repository.MatchRepository,
	userRepo repository.UserRepository,
) *ChatService {
	return &ChatService{
		chatRepo:      chatRepo,
		chatListRepo:  chatListRepo,
		websocketRepo: websocketRepo,
		matchRepo:     matchRepo,
		userRepo:      userRepo,
	}
}

// SendMessageRequest 發送訊息請求
type SendMessageRequest struct {
	SenderID   uint               `json:"sender_id" validate:"required"`
	ReceiverID uint               `json:"receiver_id" validate:"required"`
	MatchID    uint               `json:"match_id" validate:"required"`
	Type       entity.MessageType `json:"type" validate:"required"`
	Content    string             `json:"content" validate:"required"`
	FileName   *string            `json:"file_name,omitempty"`
	FileSize   *int64             `json:"file_size,omitempty"`
	FilePath   *string            `json:"file_path,omitempty"`
}

// MessageResponse 訊息回應
type MessageResponse struct {
	Message *entity.ChatMessage `json:"message"`
	Success bool                `json:"success"`
	Error   string              `json:"error,omitempty"`
}

// ChatHistoryRequest 聊天歷史請求
type ChatHistoryRequest struct {
	UserID  uint       `json:"user_id" validate:"required"`
	MatchID uint       `json:"match_id" validate:"required"`
	Limit   int        `json:"limit"`
	Before  *time.Time `json:"before,omitempty"`
}

// ChatHistoryResponse 聊天歷史回應
type ChatHistoryResponse struct {
	Messages   []*entity.ChatMessage `json:"messages"`
	TotalCount int                   `json:"total_count"`
	HasMore    bool                  `json:"has_more"`
}

// SendMessage 發送聊天訊息
// 驗證配對關係、處理訊息發送、更新聊天活躍狀態
func (s *ChatService) SendMessage(ctx context.Context, req *SendMessageRequest) (*MessageResponse, error) {
	// 驗證請求資料
	if err := s.validateSendMessageRequest(req); err != nil {
		return &MessageResponse{
			Success: false,
			Error:   fmt.Sprintf("訊息驗證失敗: %v", err),
		}, nil
	}

	// 驗證配對關係是否存在且有效
	match, err := s.matchRepo.GetMatchByID(ctx, req.MatchID)
	if err != nil {
		return &MessageResponse{
			Success: false,
			Error:   "配對關係不存在",
		}, nil
	}

	// 驗證用戶是否為配對的一方
	if !s.isUserInMatch(match, req.SenderID) {
		return &MessageResponse{
			Success: false,
			Error:   "無權限在此配對中發送訊息",
		}, nil
	}

	// 驗證配對狀態（只有已配對成功才能聊天）
	if match.Status != entity.MatchStatusMatched {
		return &MessageResponse{
			Success: false,
			Error:   "配對尚未成功，無法發送訊息",
		}, nil
	}

	// 驗證發送者和接收者
	if !s.isValidChatPair(match, req.SenderID, req.ReceiverID) {
		return &MessageResponse{
			Success: false,
			Error:   "發送者和接收者必須是配對的雙方",
		}, nil
	}

	// 創建訊息
	message := &entity.ChatMessage{
		MatchID:    req.MatchID,
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Type:       req.Type,
		Content:    req.Content,
		Status:     entity.MessageStatusSent,
	}

	// 處理檔案訊息的額外欄位
	if req.Type == entity.MessageTypeImage || req.Type == entity.MessageTypeFile {
		message.FileName = req.FileName
		message.FileSize = req.FileSize
		message.FilePath = req.FilePath
	}

	// 保存訊息
	if err := s.chatRepo.CreateMessage(ctx, message); err != nil {
		return &MessageResponse{
			Success: false,
			Error:   fmt.Sprintf("保存訊息失敗: %v", err),
		}, nil
	}

	// 更新聊天活躍狀態
	if err := s.chatListRepo.UpdateChatActivity(ctx, req.MatchID, time.Now()); err != nil {
		// 記錄錯誤但不影響訊息發送
		// 可以考慮使用日誌系統記錄
	}

	return &MessageResponse{
		Message: message,
		Success: true,
	}, nil
}

// GetChatHistory 獲取聊天歷史
// 支援分頁載入聊天記錄
func (s *ChatService) GetChatHistory(ctx context.Context, req *ChatHistoryRequest) (*ChatHistoryResponse, error) {
	// 驗證用戶權限
	match, err := s.matchRepo.GetMatchByID(ctx, req.MatchID)
	if err != nil {
		return nil, errors.New("配對不存在")
	}

	if !s.isUserInMatch(match, req.UserID) {
		return nil, errors.New("無權限查看此聊天記錄")
	}

	// 設定預設分頁大小
	limit := req.Limit
	if limit <= 0 {
		limit = 20 // 預設載入20條訊息
	}
	if limit > 100 {
		limit = 100 // 最大限制100條
	}

	// 獲取聊天記錄
	messages, err := s.chatRepo.GetMessagesByMatch(ctx, req.MatchID, limit+1, req.Before)
	if err != nil {
		return nil, fmt.Errorf("獲取聊天記錄失敗: %w", err)
	}

	// 判斷是否還有更多訊息
	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit] // 移除多餘的一條訊息
	}

	return &ChatHistoryResponse{
		Messages:   messages,
		TotalCount: len(messages),
		HasMore:    hasMore,
	}, nil
}

// MarkMessagesAsRead 標記訊息為已讀
// 用於用戶進入聊天室時標記未讀訊息
func (s *ChatService) MarkMessagesAsRead(ctx context.Context, matchID, userID uint) error {
	// 驗證用戶權限
	match, err := s.matchRepo.GetMatchByID(ctx, matchID)
	if err != nil {
		return errors.New("配對不存在")
	}

	if !s.isUserInMatch(match, userID) {
		return errors.New("無權限操作此聊天")
	}

	return s.chatRepo.MarkMessagesAsRead(ctx, matchID, userID)
}

// GetChatList 獲取用戶聊天列表
// 返回用戶所有聊天對話的預覽資訊
func (s *ChatService) GetChatList(ctx context.Context, userID uint) ([]*repository.ChatListItem, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用戶不存在")
	}

	if !user.IsActive {
		return nil, errors.New("用戶未啟用")
	}

	return s.chatListRepo.GetChatList(ctx, userID)
}

// GetActiveChatList 獲取活躍聊天列表
// 只返回有訊息交換的聊天對話
func (s *ChatService) GetActiveChatList(ctx context.Context, userID uint) ([]*repository.ChatListItem, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用戶不存在")
	}

	if !user.IsActive {
		return nil, errors.New("用戶未啟用")
	}

	return s.chatListRepo.GetActiveChatList(ctx, userID)
}

// GetUnreadCount 獲取未讀訊息數量
// 用於顯示聊天列表中的未讀數字
func (s *ChatService) GetUnreadCount(ctx context.Context, matchID, userID uint) (int64, error) {
	// 驗證用戶權限
	match, err := s.matchRepo.GetMatchByID(ctx, matchID)
	if err != nil {
		return 0, errors.New("配對不存在")
	}

	if !s.isUserInMatch(match, userID) {
		return 0, errors.New("無權限查看此聊天")
	}

	return s.chatRepo.GetUnreadCount(ctx, matchID, userID)
}

// DeleteMessage 刪除訊息
// 軟刪除訊息（標記為已刪除，不實際刪除）
func (s *ChatService) DeleteMessage(ctx context.Context, messageID, userID uint) error {
	// 獲取訊息
	message, err := s.chatRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return errors.New("訊息不存在")
	}

	// 驗證權限（只有發送者可以刪除）
	if message.SenderID != userID {
		return errors.New("只能刪除自己發送的訊息")
	}

	return s.chatRepo.DeleteMessage(ctx, messageID)
}

// GetChatStats 獲取用戶聊天統計
// 提供用戶聊天行為的數據分析
func (s *ChatService) GetChatStats(ctx context.Context, userID uint) (*repository.ChatStats, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用戶不存在")
	}

	if !user.IsActive {
		return nil, errors.New("用戶未啟用")
	}

	return s.chatListRepo.GetChatStats(ctx, userID)
}

// WebSocket 連接管理相關方法

// HandleUserConnect 處理用戶WebSocket連接
// 記錄用戶上線狀態
func (s *ChatService) HandleUserConnect(ctx context.Context, userID uint, connectionID string) error {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("用戶不存在")
	}

	if !user.IsActive {
		return errors.New("用戶未啟用")
	}

	// 記錄連接
	return s.websocketRepo.StoreConnection(ctx, userID, connectionID, time.Now())
}

// HandleUserDisconnect 處理用戶WebSocket斷線
// 清理連接記錄並更新最後上線時間
func (s *ChatService) HandleUserDisconnect(ctx context.Context, connectionID string) error {
	// 移除連接記錄
	return s.websocketRepo.RemoveConnection(ctx, connectionID)
}

// IsUserOnline 檢查用戶是否在線
// 用於聊天列表顯示在線狀態
func (s *ChatService) IsUserOnline(ctx context.Context, userID uint) (bool, error) {
	return s.websocketRepo.IsUserOnline(ctx, userID)
}

// GetUserConnections 獲取用戶的所有連接
// 用於多設備訊息推送
func (s *ChatService) GetUserConnections(ctx context.Context, userID uint) ([]string, error) {
	return s.websocketRepo.GetUserConnections(ctx, userID)
}

// UpdateLastSeen 更新用戶最後上線時間
// 用於顯示「最後上線於...」
func (s *ChatService) UpdateLastSeen(ctx context.Context, userID uint) error {
	return s.websocketRepo.UpdateLastSeen(ctx, userID, time.Now())
}

// GetLastSeen 獲取用戶最後上線時間
func (s *ChatService) GetLastSeen(ctx context.Context, userID uint) (*time.Time, error) {
	return s.websocketRepo.GetLastSeen(ctx, userID)
}

// 私有輔助方法

// validateSendMessageRequest 驗證發送訊息請求
func (s *ChatService) validateSendMessageRequest(req *SendMessageRequest) error {
	if req.SenderID == 0 {
		return errors.New("發送者ID不能為空")
	}

	if req.ReceiverID == 0 {
		return errors.New("接收者ID不能為空")
	}

	if req.MatchID == 0 {
		return errors.New("配對ID不能為空")
	}

	if req.SenderID == req.ReceiverID {
		return errors.New("不能發送訊息給自己")
	}

	if !req.Type.IsValid() {
		return errors.New("無效的訊息類型")
	}

	if strings.TrimSpace(req.Content) == "" {
		return errors.New("訊息內容不能為空")
	}

	if len(req.Content) > 1000 {
		return errors.New("訊息內容不能超過1000個字符")
	}

	// 檔案訊息的額外驗證
	if req.Type == entity.MessageTypeImage || req.Type == entity.MessageTypeFile {
		if req.FileName == nil || strings.TrimSpace(*req.FileName) == "" {
			return errors.New("檔案訊息必須提供檔案名稱")
		}

		if req.FilePath == nil || strings.TrimSpace(*req.FilePath) == "" {
			return errors.New("檔案訊息必須提供檔案路徑")
		}
	}

	return nil
}

// isUserInMatch 檢查用戶是否為配對的一方
func (s *ChatService) isUserInMatch(match *entity.Match, userID uint) bool {
	return match.User1ID == userID || match.User2ID == userID
}

// isValidChatPair 檢查發送者和接收者是否為有效的聊天對象
func (s *ChatService) isValidChatPair(match *entity.Match, senderID, receiverID uint) bool {
	return (match.User1ID == senderID && match.User2ID == receiverID) ||
		(match.User1ID == receiverID && match.User2ID == senderID)
}
