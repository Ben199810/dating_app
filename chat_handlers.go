package main

import (
	"net/http"
	"sync"
	"github.com/gin-gonic/gin"
)

type Message struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

var (
	messages []Message
	mu       sync.Mutex
)

func getMessagesHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, messages)
}

func postMessageHandler(c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	messages = append(messages, msg)
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
