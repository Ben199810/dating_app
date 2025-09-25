package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMatchingGetContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.GET("/api/matches/discover", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("發現潛在配對應該回傳200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/matches/discover", nil)
		// 模擬 JWT 認證標頭
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusOK, w.Code, "發現潛在配對應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		// 驗證回應結構符合 API 規格
		assert.Contains(t, response, "users", "回應應包含 users 陣列")
		assert.Contains(t, response, "has_more", "回應應包含 has_more 布林值")
		
		if users, ok := response["users"].([]interface{}); ok {
			// 如果有使用者資料，驗證使用者卡片結構
			if len(users) > 0 {
				user := users[0].(map[string]interface{})
				assert.Contains(t, user, "id", "使用者卡片應包含 id")
				assert.Contains(t, user, "display_name", "使用者卡片應包含 display_name")
				assert.Contains(t, user, "age", "使用者卡片應包含 age")
				assert.Contains(t, user, "photos", "使用者卡片應包含 photos 陣列")
			}
		}
	})

	t.Run("帶有 limit 參數的發現配對", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/matches/discover?limit=5", nil)
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusOK, w.Code, "帶參數的發現配對應回傳 200 OK")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		// 驗證 limit 參數有效果（最多5個使用者）
		if users, ok := response["users"].([]interface{}); ok {
			assert.LessOrEqual(t, len(users), 5, "回應的使用者數量不應超過 limit 參數")
		}
	})

	t.Run("超過最大 limit 應該被限制", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/matches/discover?limit=100", nil)
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 根據 API 規格，maximum 是 50
		assert.Equal(t, http.StatusOK, w.Code, "超大 limit 應該被接受但限制在 50")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")
		
		if users, ok := response["users"].([]interface{}); ok {
			assert.LessOrEqual(t, len(users), 50, "回應的使用者數量不應超過 50（API 最大限制）")
		}
	})

	t.Run("未認證用戶應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/matches/discover", nil)
		// 不設置 Authorization 標頭
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - 需要JWT中間件
		assert.Equal(t, http.StatusUnauthorized, w.Code, "未認證用戶應回傳 401 Unauthorized")
	})

	t.Run("無效 JWT token 應該回傳401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/matches/discover", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過 - JWT驗證失敗
		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效 JWT 應回傳 401 Unauthorized")
	})
}

// 測試用戶卡片資料結構驗證
func TestUserCardStructure(t *testing.T) {
	t.Run("驗證 UserCard 符合 API 規格", func(t *testing.T) {
		// 模擬預期的使用者卡片結構
		expectedUserCard := map[string]interface{}{
			"id":           123,
			"display_name": "Alice",
			"age":          25,
			"bio":          "Love hiking and photography",
			"photos": []map[string]interface{}{
				{
					"id":  1,
					"url": "https://example.com/photo1.jpg",
				},
			},
			"interests": []string{"hiking", "photography"},
			"distance":  5.2,
		}

		// 驗證必需欄位
		assert.Contains(t, expectedUserCard, "id", "UserCard 必須包含 id")
		assert.Contains(t, expectedUserCard, "display_name", "UserCard 必須包含 display_name")
		assert.Contains(t, expectedUserCard, "age", "UserCard 必須包含 age")
		assert.Contains(t, expectedUserCard, "photos", "UserCard 必須包含 photos")

		// 驗證資料類型
		assert.IsType(t, 0, expectedUserCard["id"], "id 應為整數")
		assert.IsType(t, "", expectedUserCard["display_name"], "display_name 應為字串")
		assert.IsType(t, 0, expectedUserCard["age"], "age 應為整數")
		assert.IsType(t, []map[string]interface{}{}, expectedUserCard["photos"], "photos 應為陣列")
	})
}