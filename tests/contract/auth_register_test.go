package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.POST("/api/auth/register", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("有效註冊應該回傳201", func(t *testing.T) {
		// 有效的註冊請求資料
		payload := map[string]interface{}{
			"email":        "user@example.com",
			"password":     "SecurePassword123",
			"birth_date":   "1995-06-15",
			"display_name": "John",
			"gender":       "male",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusCreated, w.Code, "有效註冊應回傳 201 Created")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		assert.Contains(t, response, "message", "回應應包含訊息")
		assert.Contains(t, response, "user_id", "回應應包含使用者ID")
		assert.Equal(t, "註冊成功", response["message"])
	})

	t.Run("無效電子郵件應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":        "invalid-email",
			"password":     "SecurePassword123",
			"birth_date":   "1995-06-15",
			"display_name": "John",
			"gender":       "male",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "無效電子郵件應回傳 400 Bad Request")
	})

	t.Run("未成年使用者應該回傳400", func(t *testing.T) {
		// 計算18歲以下的日期
		underageDate := time.Now().AddDate(-17, 0, 0).Format("2006-01-02")
		
		payload := map[string]interface{}{
			"email":        "underage@example.com",
			"password":     "SecurePassword123",
			"birth_date":   underageDate,
			"display_name": "Minor",
			"gender":       "male",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "未成年使用者應回傳 400 Bad Request")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error", "回應應包含錯誤訊息")
	})

	t.Run("缺少必填欄位應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "user@example.com",
			"password": "SecurePassword123",
			// 缺少 birth_date, display_name, gender
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "缺少必填欄位應回傳 400 Bad Request")
	})

	t.Run("弱密碼應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":        "user@example.com",
			"password":     "123", // 密碼太弱
			"birth_date":   "1995-06-15",
			"display_name": "John",
			"gender":       "male",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "弱密碼應回傳 400 Bad Request")
	})

	t.Run("無效性別應該回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":        "user@example.com",
			"password":     "SecurePassword123",
			"birth_date":   "1995-06-15",
			"display_name": "John",
			"gender":       "invalid", // 無效的性別
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "無效性別應回傳 400 Bad Request")
	})
}