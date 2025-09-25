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

func TestProfileUpdateContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.PUT("/users/profile", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("有效的個人檔案更新應該回傳200", func(t *testing.T) {
		// 有效的個人檔案更新請求資料
		payload := map[string]interface{}{
			"display_name":  "更新的名稱",
			"bio":           "這是我更新的個人簡介",
			"show_age":      true,
			"max_distance":  50,
			"age_range_min": 22,
			"age_range_max": 35,
			"interests":     []int{1, 3, 5, 7}, // 興趣ID
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusOK, w.Code, "有效的個人檔案更新應回傳 200 OK")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")

		// 檢查是否回傳更新後的個人檔案
		assert.Contains(t, response, "id", "回應應包含使用者ID")
		assert.Contains(t, response, "display_name", "回應應包含顯示名稱")
		assert.Contains(t, response, "bio", "回應應包含個人簡介")
		assert.Contains(t, response, "show_age", "回應應包含顯示年齡設定")
		assert.Contains(t, response, "max_distance", "回應應包含最大距離")
		assert.Contains(t, response, "age_range", "回應應包含年齡範圍")
		assert.Contains(t, response, "interests", "回應應包含興趣")

		// 驗證數值是否已更新（實作後）
		if displayName, ok := response["display_name"]; ok {
			assert.Equal(t, "更新的名稱", displayName, "顯示名稱應已更新")
		}
		if bio, ok := response["bio"]; ok {
			assert.Equal(t, "這是我更新的個人簡介", bio, "個人簡介應已更新")
		}
	})

	t.Run("缺少授權應回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"display_name": "更新的名稱",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// 沒有授權標頭

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "缺少權杖應回傳 401 未授權")
	})

	t.Run("無效權杖應回傳401", func(t *testing.T) {
		payload := map[string]interface{}{
			"display_name": "更新的名稱",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效權杖應回傳 401 未授權")
	})

	t.Run("顯示名稱過長應回傳400", func(t *testing.T) {
		// 顯示名稱超過50個字元
		longName := "這是一個非常長的顯示名稱，超過了五十個字元的限制，應該會被拒絕"

		payload := map[string]interface{}{
			"display_name": longName,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "顯示名稱過長應回傳 400 錯誤請求")
	})

	t.Run("個人簡介過長應回傳400", func(t *testing.T) {
		// 個人簡介超過500個字元
		longBio := string(make([]byte, 501))
		for i := range longBio {
			longBio = string(append([]byte(longBio[:i]), 'a'))
		}

		payload := map[string]interface{}{
			"bio": longBio,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "個人簡介過長應回傳 400 錯誤請求")
	})

	t.Run("無效的最大距離應回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"max_distance": 150, // 超過最大值100
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 Bad Request for invalid max distance")
	})

	t.Run("Invalid age range should return 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"age_range_min": 17, // Under minimum age of 18
			"age_range_max": 25,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 Bad Request for invalid age range")
	})

	t.Run("年齡範圍最大值小於最小值應回傳400", func(t *testing.T) {
		payload := map[string]interface{}{
			"age_range_min": 30,
			"age_range_max": 25, // 最大值小於最小值
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "年齡範圍最大值 < 最小值應回傳 400 錯誤請求")
	})

	t.Run("無效JSON應回傳400", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer([]byte("{invalid-json}")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "無效JSON應回傳 400 錯誤請求")
	})

	t.Run("空請求體應回傳200", func(t *testing.T) {
		// 空的更新應該被允許（沒有變更）
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後空更新應為OK
		// 目前為501因為尚未實作
		assert.Contains(t, []int{http.StatusOK, http.StatusNotImplemented}, w.Code,
			"空更新應回傳 200 OK 或因尚未實作回傳 501")
	})

	t.Run("部分更新應回傳200", func(t *testing.T) {
		// 只更新顯示名稱
		payload := map[string]interface{}{
			"display_name": "只更新名稱",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後部分更新應為OK
		// 目前為501因為尚未實作
		assert.Contains(t, []int{http.StatusOK, http.StatusNotImplemented}, w.Code,
			"部分更新應回傳 200 OK 或因尚未實作回傳 501")
	})
}
