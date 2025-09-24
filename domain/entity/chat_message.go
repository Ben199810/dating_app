package entity

import (
	"errors"
	"strings"
	"time"
)

// MessageType 訊息類型枚舉
type MessageType string

const (
	MessageTypeText   MessageType = "text"   // 文字訊息
	MessageTypeImage  MessageType = "image"  // 圖片訊息
	MessageTypeFile   MessageType = "file"   // 檔案訊息
	MessageTypeSystem MessageType = "system" // 系統訊息
)

// IsValid 檢查訊息類型是否有效
func (mt MessageType) IsValid() bool {
	return mt == MessageTypeText || mt == MessageTypeImage ||
		mt == MessageTypeFile || mt == MessageTypeSystem
}

// MessageStatus 訊息狀態枚舉
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"      // 已發送
	MessageStatusDelivered MessageStatus = "delivered" // 已送達
	MessageStatusRead      MessageStatus = "read"      // 已讀
)

// IsValid 檢查訊息狀態是否有效
func (ms MessageStatus) IsValid() bool {
	return ms == MessageStatusSent || ms == MessageStatusDelivered || ms == MessageStatusRead
}

// ChatMessage 聊天訊息實體
type ChatMessage struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	MatchID    uint          `gorm:"not null;index" json:"match_id"`
	SenderID   uint          `gorm:"not null;index" json:"sender_id"`
	ReceiverID uint          `gorm:"not null;index" json:"receiver_id"`
	Type       MessageType   `gorm:"not null;default:'text'" json:"type"`
	Content    string        `gorm:"not null;size:1000" json:"content"`
	Status     MessageStatus `gorm:"not null;default:'sent'" json:"status"`

	// 檔案相關（用於圖片或檔案訊息）
	FileName *string `gorm:"size:255" json:"file_name,omitempty"`
	FileSize *int64  `json:"file_size,omitempty"`
	FilePath *string `gorm:"size:500" json:"file_path,omitempty"`

	// 時間戳記
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	ReadAt      *time.Time `json:"read_at"`

	// 關聯 - 將在其他實體完成後添加
	// Match    Match `gorm:"constraint:OnDelete:CASCADE" json:"match"`
	// Sender   User  `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sender"`
	// Receiver User  `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE" json:"receiver"`
}

// Validate 驗證聊天訊息資料
func (cm *ChatMessage) Validate() error {
	if cm.MatchID == 0 {
		return errors.New("match_id 是必填欄位")
	}

	if cm.SenderID == 0 {
		return errors.New("sender_id 是必填欄位")
	}

	if cm.ReceiverID == 0 {
		return errors.New("receiver_id 是必填欄位")
	}

	if cm.SenderID == cm.ReceiverID {
		return errors.New("發送者和接收者不能是同一人")
	}

	if !cm.Type.IsValid() {
		return errors.New("type 必須是 text、image、file 或 system")
	}

	if strings.TrimSpace(cm.Content) == "" {
		return errors.New("content 是必填欄位")
	}

	if len(cm.Content) > 1000 {
		return errors.New("content 不能超過 1000 字元")
	}

	if !cm.Status.IsValid() {
		return errors.New("status 必須是 sent、delivered 或 read")
	}

	// 檔案訊息的額外驗證
	if cm.Type == MessageTypeImage || cm.Type == MessageTypeFile {
		if cm.FileName == nil || strings.TrimSpace(*cm.FileName) == "" {
			return errors.New("檔案訊息必須提供檔案名稱")
		}

		if cm.FileSize == nil || *cm.FileSize <= 0 {
			return errors.New("檔案訊息必須提供有效的檔案大小")
		}

		if cm.FilePath == nil || strings.TrimSpace(*cm.FilePath) == "" {
			return errors.New("檔案訊息必須提供檔案路徑")
		}
	}

	return nil
}

// IsTextMessage 檢查是否為文字訊息
func (cm *ChatMessage) IsTextMessage() bool {
	return cm.Type == MessageTypeText
}

// IsFileMessage 檢查是否為檔案訊息
func (cm *ChatMessage) IsFileMessage() bool {
	return cm.Type == MessageTypeImage || cm.Type == MessageTypeFile
}

// IsSystemMessage 檢查是否為系統訊息
func (cm *ChatMessage) IsSystemMessage() bool {
	return cm.Type == MessageTypeSystem
}

// MarkAsDelivered 標記訊息為已送達
func (cm *ChatMessage) MarkAsDelivered() {
	if cm.Status == MessageStatusSent {
		cm.Status = MessageStatusDelivered
		now := time.Now()
		cm.DeliveredAt = &now
		cm.UpdatedAt = now
	}
}

// MarkAsRead 標記訊息為已讀
func (cm *ChatMessage) MarkAsRead() {
	cm.Status = MessageStatusRead
	now := time.Now()
	cm.ReadAt = &now
	cm.UpdatedAt = now

	// 如果還沒標記為已送達，也一併標記
	if cm.DeliveredAt == nil {
		cm.DeliveredAt = &now
	}
}

// IsRead 檢查訊息是否已讀
func (cm *ChatMessage) IsRead() bool {
	return cm.Status == MessageStatusRead
}

// IsDelivered 檢查訊息是否已送達
func (cm *ChatMessage) IsDelivered() bool {
	return cm.Status == MessageStatusDelivered || cm.Status == MessageStatusRead
}

// IsSentBy 檢查訊息是否由指定用戶發送
func (cm *ChatMessage) IsSentBy(userID uint) bool {
	return cm.SenderID == userID
}

// IsReceivedBy 檢查訊息是否由指定用戶接收
func (cm *ChatMessage) IsReceivedBy(userID uint) bool {
	return cm.ReceiverID == userID
}

// GetFileInfo 獲取檔案資訊（如果有）
func (cm *ChatMessage) GetFileInfo() (fileName string, fileSize int64, filePath string, hasFile bool) {
	if !cm.IsFileMessage() {
		return "", 0, "", false
	}

	if cm.FileName != nil && cm.FileSize != nil && cm.FilePath != nil {
		return *cm.FileName, *cm.FileSize, *cm.FilePath, true
	}

	return "", 0, "", false
}

// SetFileInfo 設定檔案資訊
func (cm *ChatMessage) SetFileInfo(fileName string, fileSize int64, filePath string) {
	cm.FileName = &fileName
	cm.FileSize = &fileSize
	cm.FilePath = &filePath
	cm.UpdatedAt = time.Now()
}
