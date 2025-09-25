package repository

import (
	"context"
	"time"

	"golang_dev_docker/domain/entity"
)

// ChatRepository 聊天數據儲存庫介面
// 提供聊天系統的持久化操作，包括訊息發送、歷史查詢、狀態管理等功能
type ChatRepository interface {
	// CreateMessage 創建新訊息
	// 用於發送聊天訊息
	CreateMessage(ctx context.Context, message *entity.ChatMessage) error

	// GetMessageByID 根據 ID 獲取訊息
	// 用於訊息查詢和權限驗證
	GetMessageByID(ctx context.Context, id uint) (*entity.ChatMessage, error)

	// GetMessages 獲取聊天訊息列表
	// 用於聊天歷史展示，支援分頁和時間篩選
	GetMessages(ctx context.Context, params MessageQueryParams) ([]*entity.ChatMessage, error)

	// GetMessagesByMatch 根據配對 ID 獲取訊息
	// 用於特定配對的聊天歷史展示
	GetMessagesByMatch(ctx context.Context, matchID uint, limit int, before *time.Time) ([]*entity.ChatMessage, error)

	// UpdateMessage 更新訊息內容
	// 用於訊息編輯功能（如果支援）
	UpdateMessage(ctx context.Context, message *entity.ChatMessage) error

	// UpdateMessageStatus 更新訊息狀態
	// 用於已讀標記和送達確認
	UpdateMessageStatus(ctx context.Context, messageID uint, status entity.MessageStatus) error

	// MarkMessagesAsRead 批量標記訊息為已讀
	// 用於進入聊天室時標記所有未讀訊息
	MarkMessagesAsRead(ctx context.Context, matchID, userID uint) error

	// DeleteMessage 軟刪除訊息
	// 用於訊息刪除功能
	DeleteMessage(ctx context.Context, messageID uint) error

	// GetUnreadCount 獲取未讀訊息數量
	// 用於聊天列表未讀數字顯示
	GetUnreadCount(ctx context.Context, matchID, userID uint) (int64, error)

	// GetLastMessage 獲取最後一條訊息
	// 用於聊天列表最後訊息預覽
	GetLastMessage(ctx context.Context, matchID uint) (*entity.ChatMessage, error)
}

// ChatListRepository 聊天列表數據儲存庫介面
// 提供聊天列表相關的查詢功能，優化聊天室列表展示效能
type ChatListRepository interface {
	// GetChatList 獲取用戶聊天列表
	// 用於主聊天介面列表展示
	GetChatList(ctx context.Context, userID uint) ([]*ChatListItem, error)

	// GetActiveChatList 獲取有訊息的活躍聊天列表
	// 用於篩選有聊天記錄的配對
	GetActiveChatList(ctx context.Context, userID uint) ([]*ChatListItem, error)

	// UpdateChatActivity 更新聊天活躍時間
	// 用於聊天列表排序
	UpdateChatActivity(ctx context.Context, matchID uint, lastMessageTime time.Time) error

	// GetChatStats 獲取聊天統計數據
	// 用於用戶聊天行為分析
	GetChatStats(ctx context.Context, userID uint) (*ChatStats, error)
}

// WebSocketRepository WebSocket 連接數據儲存庫介面
// 提供 WebSocket 連接管理的數據持久化功能，支援即時通訊
type WebSocketRepository interface {
	// StoreConnection 儲存 WebSocket 連接資訊
	// 用於記錄用戶在線狀態
	StoreConnection(ctx context.Context, userID uint, connectionID string, connectedAt time.Time) error

	// RemoveConnection 移除 WebSocket 連接資訊
	// 用於用戶下線時清理連接記錄
	RemoveConnection(ctx context.Context, connectionID string) error

	// GetUserConnections 獲取用戶的所有連接
	// 用於多設備消息推送
	GetUserConnections(ctx context.Context, userID uint) ([]string, error)

	// IsUserOnline 檢查用戶是否在線
	// 用於顯示在線狀態
	IsUserOnline(ctx context.Context, userID uint) (bool, error)

	// GetOnlineUsers 獲取在線用戶列表
	// 用於統計和推薦功能
	GetOnlineUsers(ctx context.Context) ([]uint, error)

	// UpdateLastSeen 更新用戶最後上線時間
	// 用於離線狀態顯示
	UpdateLastSeen(ctx context.Context, userID uint, lastSeenAt time.Time) error

	// GetLastSeen 獲取用戶最後上線時間
	// 用於顯示「最後上線於...」
	GetLastSeen(ctx context.Context, userID uint) (*time.Time, error)
}

// MessageQueryParams 訊息查詢參數
type MessageQueryParams struct {
	// 基本篩選
	MatchID    *uint // 配對 ID
	SenderID   *uint // 發送者 ID
	ReceiverID *uint // 接收者 ID

	// 訊息類型篩選
	MessageType *entity.MessageType // 訊息類型

	// 狀態篩選
	Status *entity.MessageStatus // 訊息狀態

	// 時間範圍篩選
	StartTime *time.Time // 開始時間
	EndTime   *time.Time // 結束時間
	Before    *time.Time // 在此時間之前（用於分頁）
	After     *time.Time // 在此時間之後

	// 分頁參數
	Limit  int // 返回數量限制
	Offset int // 分頁偏移

	// 排序選項
	OrderBy   string // 排序欄位（created_at, updated_at）
	OrderDesc bool   // 是否降序排列

	// 其他選項
	IncludeDeleted bool // 是否包含已刪除訊息
}

// ChatListItem 聊天列表項目
type ChatListItem struct {
	MatchID          uint                `json:"match_id"`
	OtherUser        *entity.User        `json:"other_user"`         // 對話用戶資訊
	OtherUserProfile *entity.UserProfile `json:"other_user_profile"` // 對話用戶檔案
	LastMessage      *entity.ChatMessage `json:"last_message"`       // 最後一條訊息
	UnreadCount      int64               `json:"unread_count"`       // 未讀數量
	LastMessageTime  time.Time           `json:"last_message_time"`  // 最後訊息時間
	IsOnline         bool                `json:"is_online"`          // 對方是否在線
	LastSeen         *time.Time          `json:"last_seen"`          // 對方最後上線時間
}

// ChatStats 聊天統計資料
type ChatStats struct {
	TotalChats          int            `json:"total_chats"`           // 總聊天數
	ActiveChats         int            `json:"active_chats"`          // 活躍聊天數（近期有訊息）
	MessagesSent        int            `json:"messages_sent"`         // 發送訊息數
	MessagesReceived    int            `json:"messages_received"`     // 接收訊息數
	AverageResponseTime float64        `json:"average_response_time"` // 平均回覆時間（分鐘）
	LongestChat         int            `json:"longest_chat"`          // 最長聊天訊息數
	ChatActivity        map[string]int `json:"chat_activity"`         // 每日聊天活躍度
}
