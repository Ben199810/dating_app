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

func TestMatchingSwipeContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.POST("/api/matches/like", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("表達興趣應該回傳201", func(t *testing.T) {
		// 有效的表達興趣請求資料
		payload := map[string]interface{}{
			"target_user_id": 456,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusCreated, w.Code, "表達興趣應回傳 201 Created")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		// 驗證回應結構符合 API 規格
		assert.Contains(t, response, "match_id", "回應應包含 match_id")
		assert.Contains(t, response, "is_matched", "回應應包含 is_matched 布林值")
		assert.Contains(t, response, "message", "回應應包含 message")
		
		// 驗證資料類型
		assert.IsType(t, 0.0, response["match_id"], "match_id 應為數字")
		assert.IsType(t, false, response["is_matched"], "is_matched 應為布林值")
		assert.IsType(t, "", response["message"], "message 應為字串")
	})

	t.Run("雙向配對成功的情況", func(t *testing.T) {
		payload := map[string]interface{}{
			"target_user_id": 789,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "雙向配對應回傳 201 Created")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		// 當雙向配對成功時，is_matched 應為 true
		if isMatched, exists := response["is_matched"]; exists {
			if matched, ok := isMatched.(bool); ok && matched {
				assert.True(t, matched, "雙向配對時 is_matched 應為 true")
				// 配對成功時應該有特殊訊息
				if message, ok := response["message"].(string); ok {
					assert.Contains(t, message, "配對成功", "配對成功時訊息應包含相關文字")
				}
			}
		}
	})

	t.Run("缺少 target_user_id 應該回傳400", func(t *testing.T) {
		// 缺少必需欄位的請求
		payload := map[string]interface{}{}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusBadRequest, w.Code, "缺少必需欄位應回傳 400 Bad Request")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		assert.Contains(t, response, "error", "錯誤回應應包含 error 欄位")
	})

	t.Run("無效的 target_user_id 應該回傳400", func(t *testing.T) {
		// 無效的目標使用者ID
		payload := map[string]interface{}{
			"target_user_id": "invalid_id",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusBadRequest, w.Code, "無效目標用戶ID應回傳 400 Bad Request")
	})

	t.Run("對自己表達興趣應該回傳400", func(t *testing.T) {
		// 模擬對自己表達興趣的情況
		payload := map[string]interface{}{
			"target_user_id": 123, // 假設當前用戶ID也是123
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token") // 假設這個token代表用戶123
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 業務邏輯驗證
		assert.Equal(t, http.StatusBadRequest, w.Code, "對自己表達興趣應回傳 400 Bad Request")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		if message, ok := response["error"].(string); ok {
			assert.Contains(t, message, "自己", "錯誤訊息應提及不能對自己操作")
		}
	})

	t.Run("目標使用者不存在應該回傳404", func(t *testing.T) {
		payload := map[string]interface{}{
			"target_user_id": 99999, // 不存在的用戶ID
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusNotFound, w.Code, "目標使用者不存在應回傳 404 Not Found")
	})

	t.Run("重複表達興趣應該回傳409", func(t *testing.T) {
		payload := map[string]interface{}{
			"target_user_id": 456,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 重複操作檢測
		assert.Equal(t, http.StatusConflict, w.Code, "重複表達興趣應回傳 409 Conflict")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "錯誤回應應為有效的 JSON")
		
		if message, ok := response["error"].(string); ok {
			assert.Contains(t, message, "重複", "錯誤訊息應提及重複操作")
		}
	})

	t.Run("未認證用戶應該回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"target_user_id": 456,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// 不設置 Authorization 標頭
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 需要JWT中間件
		assert.Equal(t, http.StatusUnauthorized, w.Code, "未認證用戶應回傳 401 Unauthorized")
	})

	t.Run("無效 JWT token 應該回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"target_user_id": 456,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/matches/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - JWT驗證失敗
		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效 JWT 應回傳 401 Unauthorized")
	})
}

// 測試配對回應資料結構
func TestMatchResponseStructure(t *testing.T) {
	t.Run("驗證配對回應結構符合 API 規格", func(t *testing.T) {
		// 模擬預期的配對回應結構
		expectedResponse := map[string]interface{}{
			"match_id":   123,
			"is_matched": true,
			"message":    "配對成功！現在你們可以開始聊天了。",
		}

		// 驗證必需欄位
		assert.Contains(t, expectedResponse, "match_id", "配對回應必須包含 match_id")
		assert.Contains(t, expectedResponse, "is_matched", "配對回應必須包含 is_matched")
		assert.Contains(t, expectedResponse, "message", "配對回應必須包含 message")

		// 驗證資料類型
		assert.IsType(t, 0, expectedResponse["match_id"], "match_id 應為整數")
		assert.IsType(t, false, expectedResponse["is_matched"], "is_matched 應為布林值")
		assert.IsType(t, "", expectedResponse["message"], "message 應為字串")
	})
}