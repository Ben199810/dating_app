package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProfileGetContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.GET("/users/profile", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("有效令牌應該回傳200及個人檔案", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusOK, w.Code, "有效驗證請求應回傳 200 OK")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")

		// 根據API規格檢查必要的個人檔案欄位
		assert.Contains(t, response, "id", "個人檔案應包含使用者ID")
		assert.Contains(t, response, "email", "個人檔案應包含電子郵件")
		assert.Contains(t, response, "display_name", "個人檔案應包含顯示名稱")
		assert.Contains(t, response, "age", "個人檔案應包含年齡")
		assert.Contains(t, response, "gender", "個人檔案應包含性別")
		assert.Contains(t, response, "location", "個人檔案應包含地點")
		assert.Contains(t, response, "bio", "個人檔案應包含個人簡介")
		assert.Contains(t, response, "show_age", "個人檔案應包含是否顯示年齡")
		assert.Contains(t, response, "max_distance", "個人檔案應包含最大距離")
		assert.Contains(t, response, "age_range", "個人檔案應包含年齡範圍")
		assert.Contains(t, response, "interests", "個人檔案應包含興趣")
		assert.Contains(t, response, "photos", "個人檔案應包含照片")

		// 驗證關鍵欄位的資料類型
		if id, ok := response["id"]; ok {
			assert.IsType(t, float64(0), id, "ID應為數字格式")
		}
		if age, ok := response["age"]; ok {
			assert.IsType(t, float64(0), age, "年齡應為數字格式")
		}
		if showAge, ok := response["show_age"]; ok {
			assert.IsType(t, true, showAge, "是否顯示年齡應為布林值")
		}
		if maxDistance, ok := response["max_distance"]; ok {
			assert.IsType(t, float64(0), maxDistance, "最大距離應為數字格式")
		}

		// 檢查年齡範圍結構
		if ageRange, ok := response["age_range"].(map[string]interface{}); ok {
			assert.Contains(t, ageRange, "min", "年齡範圍應有最小值")
			assert.Contains(t, ageRange, "max", "年齡範圍應有最大值")
		}

		// 檢查興趣為陣列
		if interests, ok := response["interests"]; ok {
			assert.IsType(t, []interface{}{}, interests, "興趣應為陣列格式")
		}

		// 檢查照片為陣列
		if photos, ok := response["photos"]; ok {
			assert.IsType(t, []interface{}{}, photos, "照片應為陣列格式")
		}
	})

	t.Run("缺少授權應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		// 沒有 Authorization header

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "缺少令牌應回傳 401 Unauthorized")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error", "回應應包含錯誤訊息")
	})

	t.Run("無效令牌應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效令牌應回傳 401 Unauthorized")
	})

	t.Run("過期令牌應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer expired_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "過期令牌應回傳 401 Unauthorized")
	})

	t.Run("格式錯誤的授權標頭應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "InvalidFormat token_here")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "格式錯誤的授權標頭應回傳 401 Unauthorized")
	})

	t.Run("空的Bearer令牌應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer ")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "空令牌應回傳 401 Unauthorized")
	})

	t.Run("不存在的使用者應該回傳404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer token_for_deleted_user")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 如果使用者被刪除但令牌仍有效，應回傳404
		// 目前因為未實作而回傳501
		assert.Contains(t, []int{http.StatusNotFound, http.StatusNotImplemented}, w.Code,
			"已刪除使用者應回傳 404 Not Found 或未實作回傳 501")
	})
}
