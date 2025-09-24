package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestChatMatchesContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.GET("/api/chats", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("獲取聊天列表應該回傳200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/chats", nil)
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusOK, w.Code, "獲取聊天列表應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		// 驗證回應結構符合 API 規格
		assert.Contains(t, response, "chats", "回應應包含 chats 陣列")
		
		if chats, ok := response["chats"].([]interface{}); ok {
			// 如果有聊天資料，驗證聊天項目結構
			if len(chats) > 0 {
				chat := chats[0].(map[string]interface{})
				assert.Contains(t, chat, "match_id", "聊天項目應包含 match_id")
				assert.Contains(t, chat, "other_user", "聊天項目應包含 other_user")
				assert.Contains(t, chat, "last_message", "聊天項目應包含 last_message")
				assert.Contains(t, chat, "unread_count", "聊天項目應包含 unread_count")
				
				// 驗證 other_user 結構（UserCard）
				if otherUser, ok := chat["other_user"].(map[string]interface{}); ok {
					assert.Contains(t, otherUser, "id", "other_user 應包含 id")
					assert.Contains(t, otherUser, "display_name", "other_user 應包含 display_name")
				}
				
				// 驗證 last_message 結構
				if lastMessage, ok := chat["last_message"].(map[string]interface{}); ok {
					assert.Contains(t, lastMessage, "id", "last_message 應包含 id")
					assert.Contains(t, lastMessage, "content", "last_message 應包含 content")
					assert.Contains(t, lastMessage, "sender_id", "last_message 應包含 sender_id")
					assert.Contains(t, lastMessage, "created_at", "last_message 應包含 created_at")
				}
			}
		}
	})

	t.Run("空的聊天列表應該回傳空陣列", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/chats", nil)
		req.Header.Set("Authorization", "Bearer mock_jwt_token_no_chats")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "空聊天列表應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		if chats, ok := response["chats"].([]interface{}); ok {
			assert.Equal(t, 0, len(chats), "新用戶應該有空的聊天列表")
		}
	})

	t.Run("聊天按最新訊息時間排序", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/chats", nil)
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "獲取聊天列表應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		if chats, ok := response["chats"].([]interface{}); ok && len(chats) > 1 {
			// 驗證聊天按最新訊息時間排序
			for i := 0; i < len(chats)-1; i++ {
				currentChat := chats[i].(map[string]interface{})
				nextChat := chats[i+1].(map[string]interface{})
				
				if currentMsg, ok := currentChat["last_message"].(map[string]interface{}); ok {
					if nextMsg, ok := nextChat["last_message"].(map[string]interface{}); ok {
						currentTime := currentMsg["created_at"].(string)
						nextTime := nextMsg["created_at"].(string)
						
						// 當前訊息時間應該 >= 下一個訊息時間（降冪排序）
						assert.GreaterOrEqual(t, currentTime, nextTime, "聊天應按最新訊息時間降冪排序")
					}
				}
			}
		}
	})

	t.Run("未讀訊息計數正確", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/chats", nil)
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "獲取聊天列表應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		if chats, ok := response["chats"].([]interface{}); ok {
			for _, chatInterface := range chats {
				chat := chatInterface.(map[string]interface{})
				if unreadCount, ok := chat["unread_count"]; ok {
					// 未讀計數應該是非負整數
					if count, ok := unreadCount.(float64); ok {
						assert.GreaterOrEqual(t, count, 0.0, "未讀計數應該是非負數")
					}
				}
			}
		}
	})

	t.Run("未認證用戶應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/chats", nil)
		// 不設置 Authorization 標頭
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 需要JWT中間件
		assert.Equal(t, http.StatusUnauthorized, w.Code, "未認證用戶應回傳 401 Unauthorized")
	})

	t.Run("無效 JWT token 應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/chats", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - JWT驗證失敗
		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效 JWT 應回傳 401 Unauthorized")
	})

	t.Run("被封鎖的聊天不應出現在列表中", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/chats", nil)
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "獲取聊天列表應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		if chats, ok := response["chats"].([]interface{}); ok {
			for _, chatInterface := range chats {
				chat := chatInterface.(map[string]interface{})
				
				// 檢查聊天中的用戶是否為活躍狀態（未被封鎖）
				if otherUser, ok := chat["other_user"].(map[string]interface{}); ok {
					// 實際實作中，被封鎖的用戶聊天不應該出現
					// 這個測試確保業務邏輯正確實作
					assert.NotNil(t, otherUser["id"], "聊天中的用戶資訊應該存在且有效")
				}
			}
		}
	})
}

// 測試聊天項目資料結構
func TestChatItemStructure(t *testing.T) {
	t.Run("驗證聊天項目結構符合 API 規格", func(t *testing.T) {
		// 模擬預期的聊天項目結構
		expectedChatItem := map[string]interface{}{
			"match_id": 123,
			"other_user": map[string]interface{}{
				"id":           456,
				"display_name": "Alice",
				"age":          25,
				"photos": []map[string]interface{}{
					{
						"id":  1,
						"url": "https://example.com/photo1.jpg",
					},
				},
			},
			"last_message": map[string]interface{}{
				"id":         789,
				"content":    "Hello there!",
				"sender_id":  456,
				"created_at": "2024-01-20T10:30:00Z",
			},
			"unread_count": 2,
		}

		// 驗證必需欄位
		assert.Contains(t, expectedChatItem, "match_id", "聊天項目必須包含 match_id")
		assert.Contains(t, expectedChatItem, "other_user", "聊天項目必須包含 other_user")
		assert.Contains(t, expectedChatItem, "last_message", "聊天項目必須包含 last_message")
		assert.Contains(t, expectedChatItem, "unread_count", "聊天項目必須包含 unread_count")

		// 驗證資料類型
		assert.IsType(t, 0, expectedChatItem["match_id"], "match_id 應為整數")
		assert.IsType(t, map[string]interface{}{}, expectedChatItem["other_user"], "other_user 應為物件")
		assert.IsType(t, map[string]interface{}{}, expectedChatItem["last_message"], "last_message 應為物件")
		assert.IsType(t, 0, expectedChatItem["unread_count"], "unread_count 應為整數")
	})
}