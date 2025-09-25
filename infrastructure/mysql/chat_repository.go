package mysql

import (
	"context"
	"fmt"
	"time"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"

	"gorm.io/gorm"
)

// MySQLChatRepository MySQL 聊天儲存庫實作
type MySQLChatRepository struct {
	db *gorm.DB
}

// NewChatRepository 創建新的 MySQL 聊天儲存庫
func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &MySQLChatRepository{db: db}
}

// CreateMessage 創建新訊息
func (r *MySQLChatRepository) CreateMessage(ctx context.Context, message *entity.ChatMessage) error {
	if err := r.db.WithContext(ctx).Create(message).Error; err != nil {
		return fmt.Errorf("創建聊天訊息失敗: %w", err)
	}
	return nil
}

// GetMessageByID 根據 ID 獲取訊息
func (r *MySQLChatRepository) GetMessageByID(ctx context.Context, id uint) (*entity.ChatMessage, error) {
	var message entity.ChatMessage
	if err := r.db.WithContext(ctx).First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("訊息不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢訊息失敗: %w", err)
	}
	return &message, nil
}

// GetMessages 獲取聊天訊息列表
func (r *MySQLChatRepository) GetMessages(ctx context.Context, params repository.MessageQueryParams) ([]*entity.ChatMessage, error) {
	query := r.db.WithContext(ctx).Model(&entity.ChatMessage{})

	// 基本篩選
	if params.MatchID != nil {
		query = query.Where("match_id = ?", *params.MatchID)
	}

	if params.SenderID != nil {
		query = query.Where("sender_id = ?", *params.SenderID)
	}

	if params.ReceiverID != nil {
		query = query.Where("receiver_id = ?", *params.ReceiverID)
	}

	// 訊息類型篩選
	if params.MessageType != nil {
		query = query.Where("message_type = ?", *params.MessageType)
	}

	// 狀態篩選
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	// 時間範圍篩選
	if params.StartTime != nil {
		query = query.Where("created_at >= ?", *params.StartTime)
	}

	if params.EndTime != nil {
		query = query.Where("created_at <= ?", *params.EndTime)
	}

	if params.Before != nil {
		query = query.Where("created_at < ?", *params.Before)
	}

	if params.After != nil {
		query = query.Where("created_at > ?", *params.After)
	}

	// 排除已刪除訊息（除非特別要求）
	if !params.IncludeDeleted {
		query = query.Where("is_deleted = ?", false)
	}

	// 排序
	orderBy := "created_at"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}

	if params.OrderDesc {
		orderBy += " DESC"
	}
	query = query.Order(orderBy)

	// 分頁
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	var messages []*entity.ChatMessage
	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("獲取聊天訊息失敗: %w", err)
	}

	return messages, nil
}

// GetMessagesByMatch 根據配對 ID 獲取訊息
func (r *MySQLChatRepository) GetMessagesByMatch(ctx context.Context, matchID uint, limit int, before *time.Time) ([]*entity.ChatMessage, error) {
	query := r.db.WithContext(ctx).
		Where("match_id = ? AND is_deleted = ?", matchID, false).
		Order("created_at DESC")

	if before != nil {
		query = query.Where("created_at < ?", *before)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	var messages []*entity.ChatMessage
	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("獲取配對聊天訊息失敗: %w", err)
	}

	return messages, nil
}

// UpdateMessage 更新訊息內容
func (r *MySQLChatRepository) UpdateMessage(ctx context.Context, message *entity.ChatMessage) error {
	if err := r.db.WithContext(ctx).Save(message).Error; err != nil {
		return fmt.Errorf("更新聊天訊息失敗: %w", err)
	}
	return nil
}

// UpdateMessageStatus 更新訊息狀態
func (r *MySQLChatRepository) UpdateMessageStatus(ctx context.Context, messageID uint, status entity.MessageStatus) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status": status,
	}

	// 根據狀態更新對應的時間戳
	switch status {
	case entity.MessageStatusDelivered:
		updates["delivered_at"] = now
	case entity.MessageStatusRead:
		updates["read_at"] = now
	}

	if err := r.db.WithContext(ctx).Model(&entity.ChatMessage{}).Where("id = ?", messageID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新訊息狀態失敗: %w", err)
	}
	return nil
}

// MarkMessagesAsRead 批量標記訊息為已讀
func (r *MySQLChatRepository) MarkMessagesAsRead(ctx context.Context, matchID, userID uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":  entity.MessageStatusRead,
		"read_at": now,
	}

	if err := r.db.WithContext(ctx).Model(&entity.ChatMessage{}).
		Where("match_id = ? AND receiver_id = ? AND status != ?", matchID, userID, entity.MessageStatusRead).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("批量標記訊息為已讀失敗: %w", err)
	}
	return nil
}

// DeleteMessage 軟刪除訊息
func (r *MySQLChatRepository) DeleteMessage(ctx context.Context, messageID uint) error {
	if err := r.db.WithContext(ctx).Model(&entity.ChatMessage{}).Where("id = ?", messageID).Update("is_deleted", true).Error; err != nil {
		return fmt.Errorf("刪除訊息失敗: %w", err)
	}
	return nil
}

// GetUnreadCount 獲取未讀訊息數量
func (r *MySQLChatRepository) GetUnreadCount(ctx context.Context, matchID, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.ChatMessage{}).
		Where("match_id = ? AND receiver_id = ? AND status != ? AND is_deleted = ?", matchID, userID, entity.MessageStatusRead, false).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("獲取未讀訊息數量失敗: %w", err)
	}
	return count, nil
}

// GetLastMessage 獲取最後一條訊息
func (r *MySQLChatRepository) GetLastMessage(ctx context.Context, matchID uint) (*entity.ChatMessage, error) {
	var message entity.ChatMessage
	if err := r.db.WithContext(ctx).
		Where("match_id = ? AND is_deleted = ?", matchID, false).
		Order("created_at DESC").
		First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 沒有訊息不算錯誤
		}
		return nil, fmt.Errorf("獲取最後訊息失敗: %w", err)
	}
	return &message, nil
}

// MySQLChatListRepository MySQL 聊天列表儲存庫實作
type MySQLChatListRepository struct {
	db *gorm.DB
}

// NewChatListRepository 創建新的 MySQL 聊天列表儲存庫
func NewChatListRepository(db *gorm.DB) repository.ChatListRepository {
	return &MySQLChatListRepository{db: db}
}

// GetChatList 獲取用戶聊天列表
func (r *MySQLChatListRepository) GetChatList(ctx context.Context, userID uint) ([]*repository.ChatListItem, error) {
	var items []*repository.ChatListItem

	// 複雜查詢，需要獲取配對、對方用戶、最後訊息、未讀數等信息
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT 
			m.id as match_id,
			CASE 
				WHEN m.user1_id = ? THEN m.user2_id 
				ELSE m.user1_id 
			END as other_user_id,
			(
				SELECT content 
				FROM chat_messages cm 
				WHERE cm.match_id = m.id AND cm.is_deleted = false 
				ORDER BY cm.created_at DESC 
				LIMIT 1
			) as last_message_content,
			(
				SELECT cm.created_at 
				FROM chat_messages cm 
				WHERE cm.match_id = m.id AND cm.is_deleted = false 
				ORDER BY cm.created_at DESC 
				LIMIT 1
			) as last_message_time,
			(
				SELECT COUNT(*) 
				FROM chat_messages cm 
				WHERE cm.match_id = m.id 
				AND cm.receiver_id = ? 
				AND cm.status != 'read' 
				AND cm.is_deleted = false
			) as unread_count
		FROM matches m 
		WHERE (m.user1_id = ? OR m.user2_id = ?) 
		AND m.status = 'matched'
		ORDER BY last_message_time DESC
	`, userID, userID, userID, userID).Rows()

	if err != nil {
		return nil, fmt.Errorf("獲取聊天列表失敗: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var matchID, otherUserID uint
		var lastMessageContent *string
		var lastMessageTime *time.Time
		var unreadCount int64

		if err := rows.Scan(&matchID, &otherUserID, &lastMessageContent, &lastMessageTime, &unreadCount); err != nil {
			return nil, fmt.Errorf("掃描聊天列表行失敗: %w", err)
		}

		// 獲取對方用戶資訊
		var otherUser entity.User
		if err := r.db.WithContext(ctx).First(&otherUser, otherUserID).Error; err != nil {
			continue // 跳過無法獲取用戶資訊的記錄
		}

		// 獲取對方用戶檔案
		var otherUserProfile entity.UserProfile
		r.db.WithContext(ctx).Where("user_id = ?", otherUserID).First(&otherUserProfile)

		item := &repository.ChatListItem{
			MatchID:          matchID,
			OtherUser:        &otherUser,
			OtherUserProfile: &otherUserProfile,
			UnreadCount:      unreadCount,
		}

		if lastMessageTime != nil {
			item.LastMessageTime = *lastMessageTime
		}

		// 如果有最後訊息，創建訊息對象
		if lastMessageContent != nil {
			item.LastMessage = &entity.ChatMessage{
				Content: *lastMessageContent,
			}
		}

		items = append(items, item)
	}

	return items, nil
}

// GetActiveChatList 獲取有訊息的活躍聊天列表
func (r *MySQLChatListRepository) GetActiveChatList(ctx context.Context, userID uint) ([]*repository.ChatListItem, error) {
	// 與 GetChatList 類似，但只返回有訊息的聊天
	items, err := r.GetChatList(ctx, userID)
	if err != nil {
		return nil, err
	}

	var activeItems []*repository.ChatListItem
	for _, item := range items {
		if item.LastMessage != nil {
			activeItems = append(activeItems, item)
		}
	}

	return activeItems, nil
}

// UpdateChatActivity 更新聊天活躍時間
func (r *MySQLChatListRepository) UpdateChatActivity(ctx context.Context, matchID uint, lastMessageTime time.Time) error {
	// 這個方法可能需要一個專門的活躍度表，或者通過最後訊息時間來實現
	// 在這個簡單實作中，我們可能不需要額外的操作，因為最後訊息時間已經在訊息表中
	return nil
}

// GetChatStats 獲取聊天統計數據
func (r *MySQLChatListRepository) GetChatStats(ctx context.Context, userID uint) (*repository.ChatStats, error) {
	stats := &repository.ChatStats{
		ChatActivity: make(map[string]int),
	}

	// 總聊天數
	var totalChats int64
	if err := r.db.WithContext(ctx).Model(&entity.Match{}).
		Where("(user1_id = ? OR user2_id = ?) AND status = ?", userID, userID, entity.MatchStatusMatched).
		Count(&totalChats).Error; err != nil {
		return nil, fmt.Errorf("獲取總聊天數失敗: %w", err)
	}
	stats.TotalChats = int(totalChats)

	// 活躍聊天數（近期有訊息）
	var activeChats int64
	if err := r.db.WithContext(ctx).
		Table("matches m").
		Joins("INNER JOIN chat_messages cm ON m.id = cm.match_id").
		Where("(m.user1_id = ? OR m.user2_id = ?) AND m.status = ? AND cm.created_at > ?",
			userID, userID, entity.MatchStatusMatched, time.Now().AddDate(0, 0, -7)).
		Count(&activeChats).Error; err != nil {
		return nil, fmt.Errorf("獲取活躍聊天數失敗: %w", err)
	}
	stats.ActiveChats = int(activeChats)

	// 發送訊息數
	var messagesSent int64
	if err := r.db.WithContext(ctx).Model(&entity.ChatMessage{}).
		Where("sender_id = ? AND is_deleted = ?", userID, false).
		Count(&messagesSent).Error; err != nil {
		return nil, fmt.Errorf("獲取發送訊息數失敗: %w", err)
	}
	stats.MessagesSent = int(messagesSent)

	// 接收訊息數
	var messagesReceived int64
	if err := r.db.WithContext(ctx).Model(&entity.ChatMessage{}).
		Where("receiver_id = ? AND is_deleted = ?", userID, false).
		Count(&messagesReceived).Error; err != nil {
		return nil, fmt.Errorf("獲取接收訊息數失敗: %w", err)
	}
	stats.MessagesReceived = int(messagesReceived)

	// 最長聊天（訊息數最多的聊天）
	var longestChat int64
	r.db.WithContext(ctx).
		Table("chat_messages").
		Select("COUNT(*) as message_count").
		Joins("INNER JOIN matches ON chat_messages.match_id = matches.id").
		Where("(matches.user1_id = ? OR matches.user2_id = ?) AND chat_messages.is_deleted = ?", userID, userID, false).
		Group("chat_messages.match_id").
		Order("message_count DESC").
		Limit(1).
		Scan(&longestChat)

	stats.LongestChat = int(longestChat)

	return stats, nil
}

// MySQLWebSocketRepository MySQL WebSocket 儲存庫實作
type MySQLWebSocketRepository struct {
	db *gorm.DB
}

// NewWebSocketRepository 創建新的 MySQL WebSocket 儲存庫
func NewWebSocketRepository(db *gorm.DB) repository.WebSocketRepository {
	return &MySQLWebSocketRepository{db: db}
}

// StoreConnection 儲存 WebSocket 連接資訊
func (r *MySQLWebSocketRepository) StoreConnection(ctx context.Context, userID uint, connectionID string, connectedAt time.Time) error {
	// 這需要一個連接表，先創建一個簡單的結構
	connection := map[string]interface{}{
		"user_id":       userID,
		"connection_id": connectionID,
		"connected_at":  connectedAt,
	}

	if err := r.db.WithContext(ctx).Table("websocket_connections").Create(connection).Error; err != nil {
		return fmt.Errorf("儲存WebSocket連接失敗: %w", err)
	}
	return nil
}

// RemoveConnection 移除 WebSocket 連接資訊
func (r *MySQLWebSocketRepository) RemoveConnection(ctx context.Context, connectionID string) error {
	if err := r.db.WithContext(ctx).Table("websocket_connections").Where("connection_id = ?", connectionID).Delete(nil).Error; err != nil {
		return fmt.Errorf("移除WebSocket連接失敗: %w", err)
	}
	return nil
}

// GetUserConnections 獲取用戶的所有連接
func (r *MySQLWebSocketRepository) GetUserConnections(ctx context.Context, userID uint) ([]string, error) {
	var connectionIDs []string
	if err := r.db.WithContext(ctx).Table("websocket_connections").
		Where("user_id = ?", userID).
		Pluck("connection_id", &connectionIDs).Error; err != nil {
		return nil, fmt.Errorf("獲取用戶連接失敗: %w", err)
	}
	return connectionIDs, nil
}

// IsUserOnline 檢查用戶是否在線
func (r *MySQLWebSocketRepository) IsUserOnline(ctx context.Context, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Table("websocket_connections").
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("檢查用戶在線狀態失敗: %w", err)
	}
	return count > 0, nil
}

// GetOnlineUsers 獲取在線用戶列表
func (r *MySQLWebSocketRepository) GetOnlineUsers(ctx context.Context) ([]uint, error) {
	var userIDs []uint
	if err := r.db.WithContext(ctx).Table("websocket_connections").
		Distinct("user_id").
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, fmt.Errorf("獲取在線用戶列表失敗: %w", err)
	}
	return userIDs, nil
}

// UpdateLastSeen 更新用戶最後上線時間
func (r *MySQLWebSocketRepository) UpdateLastSeen(ctx context.Context, userID uint, lastSeenAt time.Time) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", userID).
		Update("last_seen_at", lastSeenAt).Error; err != nil {
		return fmt.Errorf("更新最後上線時間失敗: %w", err)
	}
	return nil
}

// GetLastSeen 獲取用戶最後上線時間
func (r *MySQLWebSocketRepository) GetLastSeen(ctx context.Context, userID uint) (*time.Time, error) {
	var lastSeen *time.Time
	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", userID).
		Pluck("last_seen_at", &lastSeen).Error; err != nil {
		return nil, fmt.Errorf("獲取最後上線時間失敗: %w", err)
	}
	return lastSeen, nil
}
