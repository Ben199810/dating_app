package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/usecase"
)

// ChatHandler 聊天處理器
type ChatHandler struct {
	chatService     *usecase.ChatService
	matchingService *usecase.MatchingService
}

// 全域聊天處理器實例
var chatHandler *ChatHandler

// SetChatService 設置聊天處理器的服務依賴
func SetChatService(chatService *usecase.ChatService, matchingService *usecase.MatchingService) {
	chatHandler = &ChatHandler{
		chatService:     chatService,
		matchingService: matchingService,
	}
}

// SendMessageRequest 發送訊息請求結構
type SendMessageRequest struct {
	MatchID  uint    `json:"match_id" binding:"required"`
	Content  string  `json:"content" binding:"required"`
	Type     string  `json:"type"` // "text", "image", "file"
	FileName *string `json:"file_name,omitempty"`
	FileSize *int64  `json:"file_size,omitempty"`
	FilePath *string `json:"file_path,omitempty"`
}

// GetChatMatchesHandler 獲取用戶的配對聊天列表
// GET /chat/matches
func GetChatMatchesHandler(c *gin.Context) {
	if chatHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "聊天服務未初始化",
		})
		return
	}

	// 從 JWT token 中獲取用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶ID格式錯誤",
		})
		return
	}

	// 獲取用戶的配對成功列表（用於聊天）
	matchList, err := chatHandler.matchingService.GetUserMatches(c.Request.Context(), userIDUint, entity.MatchStatusMatched)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "獲取配對列表失敗",
			"message": err.Error(),
		})
		return
	}

	// 獲取聊天列表（包含最後訊息等資訊）
	chatList, err := chatHandler.chatService.GetActiveChatList(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "獲取聊天列表失敗",
			"message": err.Error(),
		})
		return
	}

	// 構建回應資料 - 合併配對和聊天資訊
	var chatMatches []map[string]interface{}

	// 建立聊天資訊的映射表（按 match_id）
	chatMap := make(map[uint]interface{})
	for _, chat := range chatList {
		chatMap[chat.MatchID] = map[string]interface{}{
			"last_message":      chat.LastMessage,
			"last_message_time": chat.LastMessageTime,
			"unread_count":      chat.UnreadCount,
		}
	}

	// 組合配對與聊天資訊
	for _, match := range matchList.Matches {
		chatInfo := map[string]interface{}{
			"match_id":          match.ID,
			"matched_at":        match.MatchedAt,
			"status":            match.Status,
			"last_message":      nil,
			"last_message_time": nil,
			"unread_count":      0,
		}

		// 如果存在聊天記錄，合併資訊
		if chat, exists := chatMap[match.ID]; exists {
			if chatData, ok := chat.(map[string]interface{}); ok {
				chatInfo["last_message"] = chatData["last_message"]
				chatInfo["last_message_time"] = chatData["last_message_time"]
				chatInfo["unread_count"] = chatData["unread_count"]
			}
		}

		// 確定對方用戶資訊
		var otherUserID uint
		if match.User1ID == userIDUint {
			otherUserID = match.User2ID
		} else {
			otherUserID = match.User1ID
		}
		chatInfo["other_user_id"] = otherUserID

		chatMatches = append(chatMatches, chatInfo)
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"chat_matches": chatMatches,
		"total_count":  len(chatMatches),
	})
}

// SendChatMessageHandler 發送聊天訊息
// POST /chat/messages
func SendChatMessageHandler(c *gin.Context) {
	if chatHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "聊天服務未初始化",
		})
		return
	}

	// 從 JWT token 中獲取用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶ID格式錯誤",
		})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求資料格式錯誤",
			"details": err.Error(),
		})
		return
	}

	// 驗證訊息類型
	var msgType entity.MessageType
	switch req.Type {
	case "image":
		msgType = entity.MessageTypeImage
	case "file":
		msgType = entity.MessageTypeFile
	case "", "text":
		msgType = entity.MessageTypeText
	default:
		msgType = entity.MessageTypeText
	}

	// 首先需要確定接收者ID（配對中的另一方）
	matchList, err := chatHandler.matchingService.GetUserMatches(c.Request.Context(), userIDUint, entity.MatchStatusMatched)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取配對資訊失敗",
		})
		return
	}

	// 找到對應的配對記錄
	var targetMatch *entity.Match
	for _, match := range matchList.Matches {
		if match.ID == req.MatchID {
			targetMatch = match
			break
		}
	}

	if targetMatch == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "配對不存在或無權限",
		})
		return
	}

	// 確定接收者ID
	var receiverID uint
	if targetMatch.User1ID == userIDUint {
		receiverID = targetMatch.User2ID
	} else {
		receiverID = targetMatch.User1ID
	}

	// 構建服務層請求
	serviceReq := &usecase.SendMessageRequest{
		SenderID:   userIDUint,
		ReceiverID: receiverID,
		MatchID:    req.MatchID,
		Type:       msgType,
		Content:    req.Content,
		FileName:   req.FileName,
		FileSize:   req.FileSize,
		FilePath:   req.FilePath,
	}

	// 調用聊天服務發送訊息
	response, err := chatHandler.chatService.SendMessage(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "發送訊息失敗",
			"message": err.Error(),
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "發送失敗",
			"message": response.Error,
		})
		return
	}

	// 成功回應
	c.JSON(http.StatusCreated, gin.H{
		"message": "訊息發送成功",
		"data": gin.H{
			"id":         response.Message.ID,
			"content":    response.Message.Content,
			"type":       response.Message.Type,
			"created_at": response.Message.CreatedAt,
			"status":     response.Message.Status,
		},
	})
}

// GetChatHistoryHandler 獲取聊天歷史
// GET /chat/matches/:matchId/messages
func GetChatHistoryHandler(c *gin.Context) {
	if chatHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "聊天服務未初始化",
		})
		return
	}

	// 從 JWT token 中獲取用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶ID格式錯誤",
		})
		return
	}

	// 獲取URL參數中的配對ID
	matchIDStr := c.Param("matchId")
	matchID, err := strconv.ParseUint(matchIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "配對ID格式錯誤",
		})
		return
	}

	// 獲取查詢參數
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	// 解析 before 參數（用於分頁）
	var before *time.Time
	if beforeStr := c.Query("before"); beforeStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, beforeStr); err == nil {
			before = &parsedTime
		}
	}

	// 構建服務層請求
	historyReq := &usecase.ChatHistoryRequest{
		UserID:  userIDUint,
		MatchID: uint(matchID),
		Limit:   limit,
		Before:  before,
	}

	// 調用聊天服務獲取歷史
	history, err := chatHandler.chatService.GetChatHistory(c.Request.Context(), historyReq)
	if err != nil {
		if err.Error() == "配對不存在" || err.Error() == "無權限查看此聊天記錄" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "無權限",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "獲取聊天歷史失敗",
			"message": err.Error(),
		})
		return
	}

	// 構建回應資料
	var messages []map[string]interface{}
	for _, msg := range history.Messages {
		messages = append(messages, map[string]interface{}{
			"id":          msg.ID,
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
			"type":        msg.Type,
			"content":     msg.Content,
			"file_name":   msg.FileName,
			"file_path":   msg.FilePath,
			"status":      msg.Status,
			"created_at":  msg.CreatedAt,
		})
	}

	// 標記訊息為已讀
	if err := chatHandler.chatService.MarkMessagesAsRead(c.Request.Context(), uint(matchID), userIDUint); err != nil {
		// 記錄錯誤但不影響歷史返回
		// 可以考慮使用日誌系統記錄
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"messages":    messages,
		"total_count": history.TotalCount,
		"has_more":    history.HasMore,
		"limit":       limit,
	})
}
