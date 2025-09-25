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

func TestLoginContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.POST("/api/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("有效登入應該回傳200", func(t *testing.T) {
		// 有效的登入請求資料
		payload := map[string]interface{}{
			"email":    "user@example.com",
			"password": "SecurePassword123",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusOK, w.Code, "有效登入應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		assert.Contains(t, response, "access_token", "回應應包含存取令牌")
		assert.Contains(t, response, "refresh_token", "回應應包含刷新令牌")
		assert.Contains(t, response, "user", "回應應包含使用者資訊")
		
		// 驗證令牌為字串格式
		assert.IsType(t, "", response["access_token"], "存取令牌應為字串格式")
		assert.IsType(t, "", response["refresh_token"], "刷新令牌應為字串格式")
		
		// 驗證使用者物件存在且包含基本屬性
		user, ok := response["user"].(map[string]interface{})
		assert.True(t, ok, "使用者資訊應為物件")
		if ok {
			assert.Contains(t, user, "id", "使用者應有ID")
			assert.Contains(t, user, "email", "使用者應有電子郵件")
			assert.Contains(t, user, "display_name", "使用者應有顯示名稱")
		}
	})

	t.Run("無效憑證應該回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "user@example.com",
			"password": "WrongPassword",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效憑證應回傳 401 Unauthorized")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error", "回應應包含錯誤訊息")
	})

	t.Run("不存在的使用者應該回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "AnyPassword123",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "不存在使用者應回傳 401 Unauthorized")
	})

	t.Run("缺少電子郵件應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"password": "SecurePassword123",
			// 缺少 email
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "缺少電子郵件應回傳 400 Bad Request")
	})

	t.Run("缺少密碼應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"email": "user@example.com",
			// 缺少 password
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "缺少密碼應回傳 400 Bad Request")
	})

	t.Run("無效電子郵件格式應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "invalid-email-format",
			"password": "SecurePassword123",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "無效電子郵件格式應回傳 400 Bad Request")
	})

	t.Run("空的請求內容應該回傳400", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "空的請求內容應回傳 400 Bad Request")
	})

	t.Run("無效JSON應該回傳400", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("{invalid-json}")))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "無效JSON應回傳 400 Bad Request")
	})
}