package http

import (
	"context"
	"golang_dev_docker/internal/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChatHandler HTTP 聊天處理器
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler 創建聊天處理器
func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// SendMessageRequest 發送訊息請求
type SendMessageRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	RoomID  string `json:"room_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// SendMessage 發送訊息
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.chatService.SendMessage(context.Background(), req.UserID, req.RoomID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// GetMessages 獲取訊息
func (h *ChatHandler) GetMessages(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	messages, err := h.chatService.GetMessages(context.Background(), roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// JoinRoomRequest 加入聊天室請求
type JoinRoomRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoomID string `json:"room_id" binding:"required"`
}

// JoinRoom 加入聊天室
func (h *ChatHandler) JoinRoom(c *gin.Context) {
	var req JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.chatService.JoinRoom(context.Background(), req.UserID, req.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// LeaveRoom 離開聊天室
func (h *ChatHandler) LeaveRoom(c *gin.Context) {
	var req JoinRoomRequest // 使用相同的結構
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.chatService.LeaveRoom(context.Background(), req.UserID, req.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
