package websocket

import (
	"encoding/json"
	"log"
	"time"
)

// ChatHandler 處理聊天相關的 WebSocket 訊息
type ChatHandler struct {
	manager *Manager
}

// NewChatHandler 建立新的聊天處理器
func NewChatHandler(manager *Manager) *ChatHandler {
	return &ChatHandler{
		manager: manager,
	}
}

// ChatMessageData 聊天訊息數據結構
type ChatMessageData struct {
	ChatID    uint      `json:"chat_id"`
	SenderID  uint      `json:"sender_id"`
	Content   string    `json:"content"`
	MessageID uint      `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
}

// TypingStatusData 輸入狀態數據結構
type TypingStatusData struct {
	ChatID   uint `json:"chat_id"`
	SenderID uint `json:"sender_id"`
	IsTyping bool `json:"is_typing"`
}

// MessageReadData 訊息已讀狀態數據結構
type MessageReadData struct {
	ChatID    uint `json:"chat_id"`
	MessageID uint `json:"message_id"`
	ReaderID  uint `json:"reader_id"`
}

// OnlineStatusData 在線狀態數據結構
type OnlineStatusData struct {
	UserID   uint      `json:"user_id"`
	IsOnline bool      `json:"is_online"`
	LastSeen time.Time `json:"last_seen"`
}

// BroadcastNewMessage 廣播新聊天訊息
func (h *ChatHandler) BroadcastNewMessage(chatID, senderID, messageID uint, content string) {
	messageData := ChatMessageData{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		MessageID: messageID,
		Timestamp: time.Now(),
	}

	// 發送到聊天室的所有參與者
	h.manager.SendToChat(chatID, "new_message", messageData)

	log.Printf("廣播新訊息到聊天室 %d: 來自用戶 %d", chatID, senderID)
}

// BroadcastTypingStatus 廣播輸入狀態
func (h *ChatHandler) BroadcastTypingStatus(chatID, senderID uint, isTyping bool) {
	typingData := TypingStatusData{
		ChatID:   chatID,
		SenderID: senderID,
		IsTyping: isTyping,
	}

	// 發送到聊天室的其他參與者（不包括發送者）
	h.broadcastToChatExceptSender(chatID, senderID, "typing_status", typingData)

	log.Printf("廣播輸入狀態到聊天室 %d: 用戶 %d %s", chatID, senderID,
		map[bool]string{true: "開始輸入", false: "停止輸入"}[isTyping])
}

// BroadcastMessageRead 廣播訊息已讀狀態
func (h *ChatHandler) BroadcastMessageRead(chatID, messageID, readerID uint) {
	readData := MessageReadData{
		ChatID:    chatID,
		MessageID: messageID,
		ReaderID:  readerID,
	}

	// 發送到聊天室的所有參與者
	h.manager.SendToChat(chatID, "message_read", readData)

	log.Printf("廣播訊息已讀狀態到聊天室 %d: 用戶 %d 已讀訊息 %d", chatID, readerID, messageID)
}

// BroadcastUserOnlineStatus 廣播用戶在線狀態給相關聊天室
func (h *ChatHandler) BroadcastUserOnlineStatus(userID uint, isOnline bool, chatIDs []uint) {
	statusData := OnlineStatusData{
		UserID:   userID,
		IsOnline: isOnline,
		LastSeen: time.Now(),
	}

	// 發送到用戶參與的所有聊天室
	for _, chatID := range chatIDs {
		h.broadcastToChatExceptSender(chatID, userID, "user_online_status", statusData)
	}

	log.Printf("廣播用戶在線狀態: 用戶 %d %s", userID,
		map[bool]string{true: "上線", false: "離線"}[isOnline])
}

// BroadcastMatchNotification 廣播配對通知
func (h *ChatHandler) BroadcastMatchNotification(userID1, userID2 uint, chatID uint) {
	timestamp := time.Now()

	// 分別發送給兩個用戶
	h.manager.SendToUser(userID1, "new_match", map[string]interface{}{
		"match_user_id": userID2,
		"chat_id":       chatID,
		"timestamp":     timestamp,
	})

	h.manager.SendToUser(userID2, "new_match", map[string]interface{}{
		"match_user_id": userID1,
		"chat_id":       chatID,
		"timestamp":     timestamp,
	})

	log.Printf("廣播配對通知: 用戶 %d 和 %d 配對成功，聊天室 %d", userID1, userID2, chatID)
}

// BroadcastLikeNotification 廣播喜歡通知
func (h *ChatHandler) BroadcastLikeNotification(fromUserID, toUserID uint) {
	likeData := map[string]interface{}{
		"from_user_id": fromUserID,
		"timestamp":    time.Now(),
	}

	h.manager.SendToUser(toUserID, "received_like", likeData)

	log.Printf("廣播喜歡通知: 用戶 %d 喜歡了用戶 %d", fromUserID, toUserID)
}

// BroadcastSuperLikeNotification 廣播超級喜歡通知
func (h *ChatHandler) BroadcastSuperLikeNotification(fromUserID, toUserID uint) {
	superLikeData := map[string]interface{}{
		"from_user_id": fromUserID,
		"timestamp":    time.Now(),
	}

	h.manager.SendToUser(toUserID, "received_super_like", superLikeData)

	log.Printf("廣播超級喜歡通知: 用戶 %d 超級喜歡了用戶 %d", fromUserID, toUserID)
}

// BroadcastPhotoLikeNotification 廣播照片喜歡通知
func (h *ChatHandler) BroadcastPhotoLikeNotification(fromUserID, toUserID, photoID uint) {
	photoLikeData := map[string]interface{}{
		"from_user_id": fromUserID,
		"photo_id":     photoID,
		"timestamp":    time.Now(),
	}

	h.manager.SendToUser(toUserID, "photo_liked", photoLikeData)

	log.Printf("廣播照片喜歡通知: 用戶 %d 喜歡了用戶 %d 的照片 %d", fromUserID, toUserID, photoID)
}

// BroadcastMessageDelivered 廣播訊息已送達狀態
func (h *ChatHandler) BroadcastMessageDelivered(chatID, messageID, senderID uint) {
	deliveredData := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"timestamp":  time.Now(),
	}

	// 只發送給訊息發送者
	h.manager.SendToUser(senderID, "message_delivered", deliveredData)

	log.Printf("廣播訊息送達狀態: 聊天室 %d 的訊息 %d 已送達", chatID, messageID)
}

// broadcastToChatExceptSender 向聊天室廣播訊息但排除發送者
func (h *ChatHandler) broadcastToChatExceptSender(chatID, excludeUserID uint, messageType string, data interface{}) {
	// 獲取聊天室的所有客戶端
	h.manager.mu.RLock()
	clients := h.manager.chatRooms[chatID]
	h.manager.mu.RUnlock()

	if len(clients) == 0 {
		return
	}

	// 準備訊息
	msg := Message{
		Type: messageType,
		Data: data,
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("序列化訊息失敗: %v", err)
		return
	}

	// 發送給聊天室的所有客戶端，但排除指定用戶
	for _, client := range clients {
		if client.UserID != excludeUserID {
			select {
			case client.Send <- msgData:
			default:
				// 發送失敗，可能連接已關閉
				log.Printf("向客戶端 %s 發送訊息失敗", client.ID)
			}
		}
	}
}

// HandleChatMessage 處理聊天訊息相關的 WebSocket 訊息
func (h *ChatHandler) HandleChatMessage(client *Client, msgType string, data interface{}) {
	switch msgType {
	case "send_message":
		h.handleSendMessage(client, data)
	case "typing_start":
		h.handleTypingStart(client, data)
	case "typing_stop":
		h.handleTypingStop(client, data)
	case "message_read":
		h.handleMessageRead(client, data)
	case "join_chat":
		h.handleJoinChat(client, data)
	case "leave_chat":
		h.handleLeaveChat(client, data)
	default:
		log.Printf("未處理的聊天訊息類型: %s", msgType)
	}
}

// handleSendMessage 處理發送訊息
func (h *ChatHandler) handleSendMessage(client *Client, data interface{}) {
	msgData, ok := data.(map[string]interface{})
	if !ok {
		log.Printf("發送訊息數據格式錯誤")
		return
	}

	chatID, ok := msgData["chat_id"].(float64)
	if !ok {
		log.Printf("聊天室ID格式錯誤")
		return
	}

	content, ok := msgData["content"].(string)
	if !ok {
		log.Printf("訊息內容格式錯誤")
		return
	}

	// TODO: 這裡應該調用聊天服務保存訊息到資料庫
	// messageID := chatService.SaveMessage(uint(chatID), client.UserID, content)

	// 暫時使用時間戳作為訊息ID
	messageID := uint(time.Now().UnixNano())

	// 廣播新訊息
	h.BroadcastNewMessage(uint(chatID), client.UserID, messageID, content)
}

// handleTypingStart 處理開始輸入
func (h *ChatHandler) handleTypingStart(client *Client, data interface{}) {
	msgData, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	chatID, ok := msgData["chat_id"].(float64)
	if !ok {
		return
	}

	h.BroadcastTypingStatus(uint(chatID), client.UserID, true)
}

// handleTypingStop 處理停止輸入
func (h *ChatHandler) handleTypingStop(client *Client, data interface{}) {
	msgData, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	chatID, ok := msgData["chat_id"].(float64)
	if !ok {
		return
	}

	h.BroadcastTypingStatus(uint(chatID), client.UserID, false)
}

// handleMessageRead 處理訊息已讀
func (h *ChatHandler) handleMessageRead(client *Client, data interface{}) {
	msgData, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	chatID, ok := msgData["chat_id"].(float64)
	if !ok {
		return
	}

	messageID, ok := msgData["message_id"].(float64)
	if !ok {
		return
	}

	// TODO: 更新資料庫中的已讀狀態
	// chatService.MarkAsRead(uint(chatID), uint(messageID), client.UserID)

	h.BroadcastMessageRead(uint(chatID), uint(messageID), client.UserID)
}

// handleJoinChat 處理加入聊天室
func (h *ChatHandler) handleJoinChat(client *Client, data interface{}) {
	msgData, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	chatID, ok := msgData["chat_id"].(float64)
	if !ok {
		return
	}

	h.manager.JoinChat(client.UserID, uint(chatID))

	// 發送確認訊息
	confirmData := map[string]interface{}{
		"chat_id": uint(chatID),
		"status":  "joined",
	}

	if msgBytes, err := json.Marshal(Message{Type: "chat_joined", Data: confirmData}); err == nil {
		select {
		case client.Send <- msgBytes:
		default:
		}
	}
}

// handleLeaveChat 處理離開聊天室
func (h *ChatHandler) handleLeaveChat(client *Client, data interface{}) {
	msgData, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	chatID, ok := msgData["chat_id"].(float64)
	if !ok {
		return
	}

	h.manager.LeaveChat(client.UserID, uint(chatID))

	// 發送確認訊息
	confirmData := map[string]interface{}{
		"chat_id": uint(chatID),
		"status":  "left",
	}

	if msgBytes, err := json.Marshal(Message{Type: "chat_left", Data: confirmData}); err == nil {
		select {
		case client.Send <- msgBytes:
		default:
		}
	}
}

// GetChatParticipants 獲取聊天室參與者
func (h *ChatHandler) GetChatParticipants(chatID uint) []uint {
	h.manager.mu.RLock()
	defer h.manager.mu.RUnlock()

	clients := h.manager.chatRooms[chatID]
	participants := make([]uint, 0, len(clients))

	for _, client := range clients {
		participants = append(participants, client.UserID)
	}

	return participants
}

// GetOnlineUsersInChat 獲取聊天室中的在線用戶
func (h *ChatHandler) GetOnlineUsersInChat(chatID uint) []uint {
	return h.GetChatParticipants(chatID) // 在聊天室中的用戶都是在線的
}

// IsUserInChat 檢查用戶是否在聊天室中
func (h *ChatHandler) IsUserInChat(userID, chatID uint) bool {
	h.manager.mu.RLock()
	defer h.manager.mu.RUnlock()

	clients := h.manager.chatRooms[chatID]
	for _, client := range clients {
		if client.UserID == userID {
			return true
		}
	}

	return false
}
