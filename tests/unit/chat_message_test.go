package unit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang_dev_docker/domain/entity"
)

func TestMessageType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		msgType  entity.MessageType
		expected bool
	}{
		{
			name:     "Valid text message type",
			msgType:  entity.MessageTypeText,
			expected: true,
		},
		{
			name:     "Valid image message type",
			msgType:  entity.MessageTypeImage,
			expected: true,
		},
		{
			name:     "Valid file message type",
			msgType:  entity.MessageTypeFile,
			expected: true,
		},
		{
			name:     "Valid system message type",
			msgType:  entity.MessageTypeSystem,
			expected: true,
		},
		{
			name:     "Invalid message type",
			msgType:  "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.msgType.IsValid())
		})
	}
}

func TestMessageStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   entity.MessageStatus
		expected bool
	}{
		{
			name:     "Valid sent status",
			status:   entity.MessageStatusSent,
			expected: true,
		},
		{
			name:     "Valid delivered status",
			status:   entity.MessageStatusDelivered,
			expected: true,
		},
		{
			name:     "Valid read status",
			status:   entity.MessageStatusRead,
			expected: true,
		},
		{
			name:     "Invalid status",
			status:   "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestChatMessage_Validate(t *testing.T) {
	tests := []struct {
		name        string
		message     entity.ChatMessage
		expectError bool
		errorText   string
	}{
		{
			name: "Valid text message",
			message: entity.ChatMessage{
				MatchID:    1,
				SenderID:   1,
				ReceiverID: 2,
				Type:       entity.MessageTypeText,
				Content:    "Hello world",
				Status:     entity.MessageStatusSent,
			},
			expectError: false,
		},
		{
			name: "Missing match ID",
			message: entity.ChatMessage{
				SenderID:   1,
				ReceiverID: 2,
				Type:       entity.MessageTypeText,
				Content:    "Hello world",
				Status:     entity.MessageStatusSent,
			},
			expectError: true,
			errorText:   "match_id 是必填欄位",
		},
		{
			name: "Same sender and receiver",
			message: entity.ChatMessage{
				MatchID:    1,
				SenderID:   1,
				ReceiverID: 1,
				Type:       entity.MessageTypeText,
				Content:    "Hello world",
				Status:     entity.MessageStatusSent,
			},
			expectError: true,
			errorText:   "發送者和接收者不能是同一人",
		},
		{
			name: "Empty content",
			message: entity.ChatMessage{
				MatchID:    1,
				SenderID:   1,
				ReceiverID: 2,
				Type:       entity.MessageTypeText,
				Content:    "",
				Status:     entity.MessageStatusSent,
			},
			expectError: true,
			errorText:   "content 是必填欄位",
		},
		{
			name: "Content too long",
			message: entity.ChatMessage{
				MatchID:    1,
				SenderID:   1,
				ReceiverID: 2,
				Type:       entity.MessageTypeText,
				Content:    strings.Repeat("A", 1001),
				Status:     entity.MessageStatusSent,
			},
			expectError: true,
			errorText:   "content 不能超過 1000 字元",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChatMessage_MessageTypeCheckers(t *testing.T) {
	tests := []struct {
		name        string
		msgType     entity.MessageType
		isText      bool
		isFile      bool
		isSystem    bool
	}{
		{
			name:     "Text message",
			msgType:  entity.MessageTypeText,
			isText:   true,
			isFile:   false,
			isSystem: false,
		},
		{
			name:     "Image message",
			msgType:  entity.MessageTypeImage,
			isText:   false,
			isFile:   true,
			isSystem: false,
		},
		{
			name:     "File message",
			msgType:  entity.MessageTypeFile,
			isText:   false,
			isFile:   true,
			isSystem: false,
		},
		{
			name:     "System message",
			msgType:  entity.MessageTypeSystem,
			isText:   false,
			isFile:   false,
			isSystem: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &entity.ChatMessage{
				MatchID:    1,
				SenderID:   1,
				ReceiverID: 2,
				Type:       tt.msgType,
				Content:    "test content",
				Status:     entity.MessageStatusSent,
			}

			assert.Equal(t, tt.isText, message.IsTextMessage())
			assert.Equal(t, tt.isFile, message.IsFileMessage())
			assert.Equal(t, tt.isSystem, message.IsSystemMessage())
		})
	}
}

func TestChatMessage_StatusMethods(t *testing.T) {
	message := &entity.ChatMessage{
		MatchID:    1,
		SenderID:   1,
		ReceiverID: 2,
		Type:       entity.MessageTypeText,
		Content:    "Hello world",
		Status:     entity.MessageStatusSent,
	}

	// 初始狀態
	assert.False(t, message.IsRead())
	assert.False(t, message.IsDelivered())

	// 標記為已送達
	message.MarkAsDelivered()
	assert.Equal(t, entity.MessageStatusDelivered, message.Status)
	assert.True(t, message.IsDelivered())
	assert.False(t, message.IsRead())
	assert.NotNil(t, message.DeliveredAt)

	// 標記為已讀
	message.MarkAsRead()
	assert.Equal(t, entity.MessageStatusRead, message.Status)
	assert.True(t, message.IsRead())
	assert.True(t, message.IsDelivered())
	assert.NotNil(t, message.ReadAt)
}

func TestChatMessage_UserMethods(t *testing.T) {
	message := &entity.ChatMessage{
		MatchID:    1,
		SenderID:   1,
		ReceiverID: 2,
		Type:       entity.MessageTypeText,
		Content:    "Hello world",
		Status:     entity.MessageStatusSent,
	}

	assert.True(t, message.IsSentBy(1))
	assert.False(t, message.IsSentBy(2))
	assert.True(t, message.IsReceivedBy(2))
	assert.False(t, message.IsReceivedBy(1))
}

func TestChatMessage_FileInfo(t *testing.T) {
	// 文字訊息沒有檔案資訊
	textMessage := &entity.ChatMessage{
		Type:    entity.MessageTypeText,
		Content: "Hello world",
	}

	fileName, fileSize, filePath, hasFile := textMessage.GetFileInfo()
	assert.False(t, hasFile)
	assert.Empty(t, fileName)
	assert.Zero(t, fileSize)
	assert.Empty(t, filePath)

	// 檔案訊息
	fileMessage := &entity.ChatMessage{
		Type:    entity.MessageTypeFile,
		Content: "Sent a file",
	}

	fileMessage.SetFileInfo("document.pdf", 1024, "/uploads/document.pdf")

	fileName, fileSize, filePath, hasFile = fileMessage.GetFileInfo()
	assert.True(t, hasFile)
	assert.Equal(t, "document.pdf", fileName)
	assert.Equal(t, int64(1024), fileSize)
	assert.Equal(t, "/uploads/document.pdf", filePath)
}

func TestChatMessage_FileValidation(t *testing.T) {
	tests := []struct {
		name        string
		message     entity.ChatMessage
		expectError bool
		errorText   string
	}{
		{
			name: "Valid file message",
			message: func() entity.ChatMessage {
				fileName := "test.jpg"
				fileSize := int64(1024)
				filePath := "/uploads/test.jpg"
				return entity.ChatMessage{
					MatchID:    1,
					SenderID:   1,
					ReceiverID: 2,
					Type:       entity.MessageTypeImage,
					Content:    "Image sent",
					Status:     entity.MessageStatusSent,
					FileName:   &fileName,
					FileSize:   &fileSize,
					FilePath:   &filePath,
				}
			}(),
			expectError: false,
		},
		{
			name: "File message missing filename",
			message: func() entity.ChatMessage {
				fileSize := int64(1024)
				filePath := "/uploads/test.jpg"
				return entity.ChatMessage{
					MatchID:    1,
					SenderID:   1,
					ReceiverID: 2,
					Type:       entity.MessageTypeImage,
					Content:    "Image sent",
					Status:     entity.MessageStatusSent,
					FileSize:   &fileSize,
					FilePath:   &filePath,
				}
			}(),
			expectError: true,
			errorText:   "檔案訊息必須提供檔案名稱",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}