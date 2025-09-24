package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// MatchingAlgorithmTestSuite 測試配對演算法的完整功能
type MatchingAlgorithmTestSuite struct {
	suite.Suite
	router     *gin.Engine
	authTokens map[string]string // user_id -> jwt_token
	testUsers  map[string]TestUser
}

// TestUser 測試用戶資料結構
type TestUser struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Gender      string    `json:"gender"`
	Age         int       `json:"age"`
	BirthDate   time.Time `json:"birth_date"`
	Interests   []string  `json:"interests"`
	Location    string    `json:"location"`
}

func TestMatchingAlgorithmTestSuite(t *testing.T) {
	suite.Run(t, new(MatchingAlgorithmTestSuite))
}

func (suite *MatchingAlgorithmTestSuite) SetupSuite() {
	// 設置測試路由
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// 模擬配對相關的路由端點
	api := suite.router.Group("/api")
	matches := api.Group("/matches")
	{
		// 模擬 JWT 認證中間件
		matches.Use(suite.mockJWTAuth())
		
		// 這些端點目前會回傳 501 Not Implemented（遵循 TDD）
		matches.GET("/discover", suite.mockGetMatchingCandidates)
		matches.POST("/like", suite.mockCreateMatch)
	}

	suite.authTokens = make(map[string]string)
	suite.testUsers = make(map[string]TestUser)
	
	// 初始化測試用戶資料
	suite.setupTestUsers()
}

func (suite *MatchingAlgorithmTestSuite) setupTestUsers() {
	now := time.Now()

	users := map[string]TestUser{
		"alice": {
			ID:          1,
			Email:       "alice@test.com",
			DisplayName: "Alice Chen",
			Gender:      "female",
			Age:         25,
			BirthDate:   now.AddDate(-25, 0, 0),
			Interests:   []string{"travel", "music", "cooking"},
			Location:    "台北市",
		},
		"bob": {
			ID:          2,
			Email:       "bob@test.com",
			DisplayName: "Bob Wang",
			Gender:      "male", 
			Age:         30,
			BirthDate:   now.AddDate(-30, 0, 0),
			Interests:   []string{"travel", "sports", "movies"},
			Location:    "台北市",
		},
		"carol": {
			ID:          3,
			Email:       "carol@test.com",
			DisplayName: "Carol Lin",
			Gender:      "female",
			Age:         22,
			BirthDate:   now.AddDate(-22, 0, 0),
			Interests:   []string{"art", "books", "coffee"},
			Location:    "新北市",
		},
		"david": {
			ID:          4,
			Email:       "david@test.com",
			DisplayName: "David Liu",
			Gender:      "male",
			Age:         28,
			BirthDate:   now.AddDate(-28, 0, 0),
			Interests:   []string{"travel", "music", "photography"},
			Location:    "台中市",
		},
	}

	suite.testUsers = users
	
	// 設置模擬 JWT tokens
	for identifier := range users {
		suite.authTokens[identifier] = fmt.Sprintf("mock_jwt_token_for_%s", identifier)
	}
}

// mockJWTAuth 模擬 JWT 認證中間件
func (suite *MatchingAlgorithmTestSuite) mockJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Authorization header 中提取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供認證令牌"})
			c.Abort()
			return
		}
		
		// 簡單驗證 token 格式
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的令牌格式"})
			c.Abort()
			return
		}
		
		token := authHeader[7:]
		
		// 根據 token 找到對應的用戶
		var currentUser TestUser
		var found bool
		for identifier, expectedToken := range suite.authTokens {
			if token == expectedToken {
				currentUser = suite.testUsers[identifier]
				found = true
				break
			}
		}
		
		if !found {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證令牌"})
			c.Abort()
			return
		}
		
		// 將用戶資訊存入上下文
		c.Set("user_id", currentUser.ID)
		c.Set("current_user", currentUser)
		c.Next()
	}
}

// mockGetMatchingCandidates 模擬配對候選人發現端點
func (suite *MatchingAlgorithmTestSuite) mockGetMatchingCandidates(c *gin.Context) {
	// 目前回傳 501 Not Implemented（TDD 方式）
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "配對候選人發現功能尚未實作",
		"message": "GET /api/matches/discover endpoint not implemented yet",
	})
	
	// 當實作完成時，應該回傳類似這樣的結構：
	/*
	currentUser, _ := c.Get("current_user")
	user := currentUser.(TestUser)
	
	candidates := suite.findMatchingCandidates(user)
	
	c.JSON(http.StatusOK, gin.H{
		"candidates": candidates,
		"total": len(candidates),
		"page": 1,
		"limit": 10,
	})
	*/
}

// mockCreateMatch 模擬創建配對端點
func (suite *MatchingAlgorithmTestSuite) mockCreateMatch(c *gin.Context) {
	// 目前回傳 501 Not Implemented（TDD 方式）
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "配對功能尚未實作",
		"message": "POST /api/matches/like endpoint not implemented yet",
	})
	
	// 當實作完成時，應該處理點讚邏輯：
	/*
	var request struct {
		TargetUserID uint   `json:"target_user_id"`
		Action       string `json:"action"` // "like" or "pass"
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}
	
	currentUser, _ := c.Get("current_user")
	user := currentUser.(TestUser)
	
	result := suite.processMatchAction(user, request.TargetUserID, request.Action)
	c.JSON(http.StatusOK, result)
	*/
}

// TestDiscoverCandidates 測試配對候選人發現功能
func (suite *MatchingAlgorithmTestSuite) TestDiscoverCandidates() {
	// 以 Alice 的身分請求配對候選人
	req := suite.createAuthenticatedRequest("GET", "/api/matches/discover", nil, "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 目前應該回傳 501 (尚未實作)
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	
	// 驗證錯誤回應結構
	suite.Contains(response, "error")
	suite.Contains(response, "message")
	suite.Equal("配對候選人發現功能尚未實作", response["error"])
	
	// 當實作完成後，應該驗證以下結構：
	// - candidates: 候選人陣列
	// - 每個候選人包含: id, display_name, age, photos, interests, distance
	// - 不包含自己
	// - 按相似度排序
}

// TestDiscoverWithFilters 測試帶篩選條件的候選人發現
func (suite *MatchingAlgorithmTestSuite) TestDiscoverWithFilters() {
	// 測試年齡範圍篩選
	req := suite.createAuthenticatedRequest("GET", "/api/matches/discover?min_age=25&max_age=35", nil, "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	// 驗證基本回應結構
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response, "error")
	
	// 當實作完成時，應驗證：
	// - 只返回符合年齡範圍的候選人
	// - 距離篩選邏輯
	// - 興趣篩選邏輯
}

// TestSwipeLikeFlow 測試完整的點讚配對流程
func (suite *MatchingAlgorithmTestSuite) TestSwipeLikeFlow() {
	// Alice 對 Bob 點讚
	likeData := map[string]interface{}{
		"target_user_id": suite.testUsers["bob"].ID,
		"action":         "like",
	}
	
	reqBody, _ := json.Marshal(likeData)
	req := suite.createAuthenticatedRequest("POST", "/api/matches/like", bytes.NewBuffer(reqBody), "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response, "error")
	suite.Equal("配對功能尚未實作", response["error"])
	
	// 當實作完成時，應驗證：
	// 1. 第一次點讚回傳 pending 狀態
	// 2. Bob 回讚後雙向配對成功
	// 3. 配對記錄正確存入資料庫
}

// TestSwipePassFlow 測試拒絕流程
func (suite *MatchingAlgorithmTestSuite) TestSwipePassFlow() {
	// Alice 拒絕 Carol
	passData := map[string]interface{}{
		"target_user_id": suite.testUsers["carol"].ID,
		"action":         "pass",
	}
	
	reqBody, _ := json.Marshal(passData)
	req := suite.createAuthenticatedRequest("POST", "/api/matches/like", bytes.NewBuffer(reqBody), "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	// 當實作完成時，應驗證：
	// - Carol 不再出現在 Alice 的候選清單中
	// - 拒絕記錄正確存入
}

// TestDuplicateLikePrevention 測試防止重複點讚
func (suite *MatchingAlgorithmTestSuite) TestDuplicateLikePrevention() {
	likeData := map[string]interface{}{
		"target_user_id": suite.testUsers["david"].ID,
		"action":         "like",
	}
	
	reqBody, _ := json.Marshal(likeData)
	
	// 第一次點讚
	req := suite.createAuthenticatedRequest("POST", "/api/matches/like", bytes.NewBuffer(reqBody), "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	// 第二次點讚同一個人（當實作完成時應該回傳錯誤）
	reqBody, _ = json.Marshal(likeData)
	req = suite.createAuthenticatedRequest("POST", "/api/matches/like", bytes.NewBuffer(reqBody), "alice")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	// 當實作完成時，第二次點讚應該回傳 400 錯誤
	// 並包含 "已經對此用戶表達過意願" 的錯誤訊息
}

// TestSelfLikePrevention 測試防止自己點讚自己  
func (suite *MatchingAlgorithmTestSuite) TestSelfLikePrevention() {
	likeData := map[string]interface{}{
		"target_user_id": suite.testUsers["alice"].ID, // 自己點讚自己
		"action":         "like",
	}
	
	reqBody, _ := json.Marshal(likeData)
	req := suite.createAuthenticatedRequest("POST", "/api/matches/like", bytes.NewBuffer(reqBody), "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	// 當實作完成時，應該回傳 400 錯誤
	// 並包含 "無法對自己進行操作" 的錯誤訊息
}

// TestBlockedUserFiltering 測試封鎖用戶過濾
func (suite *MatchingAlgorithmTestSuite) TestBlockedUserFiltering() {
	// 模擬 Alice 已封鎖 Bob 的情況
	// 當實作完成時，需要檢查資料庫中的 blocks 表
	
	req := suite.createAuthenticatedRequest("GET", "/api/matches/discover", nil, "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	// 當實作完成時，Bob 不應該出現在候選清單中
}

// TestMatchingWithInterests 測試基於興趣的配對推薦
func (suite *MatchingAlgorithmTestSuite) TestMatchingWithInterests() {
	// Alice 的興趣: ["travel", "music", "cooking"]  
	// David 的興趣: ["travel", "music", "photography"] - 2個共同興趣
	// Bob 的興趣: ["travel", "sports", "movies"] - 1個共同興趣
	
	req := suite.createAuthenticatedRequest("GET", "/api/matches/discover", nil, "alice")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotImplemented, w.Code)
	
	// 當實作完成時，應驗證：
	// - David（2個共同興趣）排在 Bob（1個共同興趣）前面
	// - 興趣相似度計算邏輯正確
	// - 候選人按相似度排序
}

// TestUnauthorizedAccess 測試未認證訪問
func (suite *MatchingAlgorithmTestSuite) TestUnauthorizedAccess() {
	// 不帶 Authorization header
	req, _ := http.NewRequest("GET", "/api/matches/discover", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("未提供認證令牌", response["error"])
}

// TestInvalidToken 測試無效令牌
func (suite *MatchingAlgorithmTestSuite) TestInvalidToken() {
	// 使用無效的 token
	req, _ := http.NewRequest("GET", "/api/matches/discover", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("無效的認證令牌", response["error"])
}

// createAuthenticatedRequest 建立帶認證的測試請求
func (suite *MatchingAlgorithmTestSuite) createAuthenticatedRequest(method, url string, body *bytes.Buffer, userIdentifier string) *http.Request {
	var req *http.Request
	var err error
	
	if body != nil {
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	
	// 設置模擬 JWT token
	token, exists := suite.authTokens[userIdentifier]
	suite.Require().True(exists, "User %s not found in auth tokens", userIdentifier)
	req.Header.Set("Authorization", "Bearer "+token)
	
	return req
}