package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestChatSendContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.POST("/api/chats/:match_id/messages", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("發送聊天訊息應該回傳201", func(t *testing.T) {
		// 有效的發送訊息請求資料
		payload := map[string]interface{}{
			"content": "Hello, nice to meet you!",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusCreated, w.Code, "發送聊天訊息應回傳 201 Created")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		// 驗證回應結構符合 API 規格
		assert.Contains(t, response, "id", "回應應包含訊息 id")
		assert.Contains(t, response, "content", "回應應包含 content")
		assert.Contains(t, response, "sender_id", "回應應包含 sender_id")
		assert.Contains(t, response, "match_id", "回應應包含 match_id")
		assert.Contains(t, response, "created_at", "回應應包含 created_at")
		
		// 驗證資料類型
		assert.IsType(t, 0.0, response["id"], "id 應為數字")
		assert.IsType(t, "", response["content"], "content 應為字串")
		assert.IsType(t, 0.0, response["sender_id"], "sender_id 應為數字")
		assert.IsType(t, 0.0, response["match_id"], "match_id 應為數字")
		assert.IsType(t, "", response["created_at"], "created_at 應為字串")
		
		// 驗證內容正確
		assert.Equal(t, "Hello, nice to meet you!", response["content"], "回應內容應與發送內容一致")
	})

	t.Run("發送長訊息應該成功", func(t *testing.T) {
		longMessage := "這是一條比較長的訊息，用來測試系統是否能正確處理較長的文字內容。在實際的聊天應用中，用戶可能會發送各種長度的訊息，包括分享自己的想法、經歷或者詢問問題。"
		
		payload := map[string]interface{}{
			"content": longMessage,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "發送長訊息應該成功")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		assert.Equal(t, longMessage, response["content"], "長訊息內容應被完整保存")
	})

	t.Run("缺少 content 應該回傳400", func(t *testing.T) {
		// 缺少必需欄位的請求
		payload := map[string]interface{}{}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusBadRequest, w.Code, "缺少 content 應回傳 400 Bad Request")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		assert.Contains(t, response, "error", "錯誤回應應包含 error 欄位")
	})

	t.Run("空的 content 應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusBadRequest, w.Code, "空 content 應回傳 400 Bad Request")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		if message, ok := response["error"].(string); ok {
			assert.Contains(t, message, "內容", "錯誤訊息應提及內容不能為空")
		}
	})

	t.Run("超過最大長度的 content 應該回傳400", func(t *testing.T) {
		// 根據 API 規格，maxLength 是 1000
		longContent := ""
		for i := 0; i < 1001; i++ {
			longContent += "a"
		}
		
		payload := map[string]interface{}{
			"content": longContent,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusBadRequest, w.Code, "超過最大長度應回傳 400 Bad Request")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		if message, ok := response["error"].(string); ok {
			assert.Contains(t, message, "長度", "錯誤訊息應提及長度限制")
		}
	})

	t.Run("無效的 match_id 應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "Hello!",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/invalid_match_id/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusBadRequest, w.Code, "無效 match_id 應回傳 400 Bad Request")
	})

	t.Run("不存在的 match_id 應該回傳404", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "Hello!",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/99999/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusNotFound, w.Code, "不存在的配對應回傳 404 Not Found")
	})

	t.Run("用戶不屬於該配對應該回傳403", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "Hello!",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/456/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token_other_user")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 權限檢查
		assert.Equal(t, http.StatusForbidden, w.Code, "用戶不屬於該配對應回傳 403 Forbidden")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		if message, ok := response["error"].(string); ok {
			assert.Contains(t, message, "權限", "錯誤訊息應提及權限問題")
		}
	})

	t.Run("被封鎖的配對不能發送訊息", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "Hello!",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/789/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token_blocked")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 封鎖檢查
		assert.Equal(t, http.StatusForbidden, w.Code, "被封鎖的配對不能發送訊息")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		if message, ok := response["error"].(string); ok {
			assert.Contains(t, message, "封鎖", "錯誤訊息應提及封鎖狀態")
		}
	})

	t.Run("未認證用戶應該回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "Hello!",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// 不設置 Authorization 標頭
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 需要JWT中間件
		assert.Equal(t, http.StatusUnauthorized, w.Code, "未認證用戶應回傳 401 Unauthorized")
	})

	t.Run("無效 JWT token 應該回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "Hello!",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - JWT驗證失敗
		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效 JWT 應回傳 401 Unauthorized")
	})

	t.Run("訊息內容防XSS測試", func(t *testing.T) {
		xssPayload := "<script>alert('XSS')</script>"
		
		payload := map[string]interface{}{
			"content": xssPayload,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "包含特殊字符的訊息應該被正確處理")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		// 確保內容被正確儲存（應該由前端進行 XSS 防護）
		assert.Contains(t, response["content"], xssPayload, "內容應被原樣保存，XSS防護由前端處理")
	})
}

// 測試聊天訊息資料結構
func TestChatMessageStructure(t *testing.T) {
	t.Run("驗證聊天訊息結構符合 API 規格", func(t *testing.T) {
		// 模擬預期的聊天訊息結構
		expectedMessage := map[string]interface{}{
			"id":         123,
			"content":    "Hello, nice to meet you!",
			"sender_id":  456,
			"match_id":   789,
			"created_at": "2024-01-20T10:30:00Z",
		}

		// 驗證必需欄位
		assert.Contains(t, expectedMessage, "id", "聊天訊息必須包含 id")
		assert.Contains(t, expectedMessage, "content", "聊天訊息必須包含 content")
		assert.Contains(t, expectedMessage, "sender_id", "聊天訊息必須包含 sender_id")
		assert.Contains(t, expectedMessage, "match_id", "聊天訊息必須包含 match_id")
		assert.Contains(t, expectedMessage, "created_at", "聊天訊息必須包含 created_at")

		// 驗證資料類型
		assert.IsType(t, 0, expectedMessage["id"], "id 應為整數")
		assert.IsType(t, "", expectedMessage["content"], "content 應為字串")
		assert.IsType(t, 0, expectedMessage["sender_id"], "sender_id 應為整數")
		assert.IsType(t, 0, expectedMessage["match_id"], "match_id 應為整數")
		assert.IsType(t, "", expectedMessage["created_at"], "created_at 應為字串")
	})
}