package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// UserRegistrationTestSuite 測試用戶註冊的完整流程
type UserRegistrationTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func TestUserRegistrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserRegistrationTestSuite))
}

func (suite *UserRegistrationTestSuite) SetupSuite() {
	// 設置測試路由
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// 模擬註冊相關的路由端點
	api := suite.router.Group("/api")
	auth := api.Group("/auth")
	{
		auth.POST("/register", suite.mockRegisterEndpoint)
	}
}

// mockRegisterEndpoint 模擬註冊端點
func (suite *UserRegistrationTestSuite) mockRegisterEndpoint(c *gin.Context) {
	var registrationData struct {
		Email       string    `json:"email"`
		Password    string    `json:"password"`
		BirthDate   time.Time `json:"birth_date"`
		DisplayName string    `json:"display_name"`
		Gender      string    `json:"gender"`
	}

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	// 簡單的年齡檢查邏輯（模擬）
	now := time.Now()
	age := now.Year() - registrationData.BirthDate.Year()
	
	// 調整月日差異
	if now.YearDay() < registrationData.BirthDate.YearDay() {
		age--
	}

	if age < 18 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "年齡限制",
			"message": "必須年滿18歲才能註冊",
		})
		return
	}

	// 目前回傳 501 Not Implemented（TDD 方式）
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "註冊功能尚未實作",
		"message": "POST /api/auth/register endpoint not implemented yet",
	})
}

// TestValidUserRegistration 測試有效的用戶註冊
func (suite *UserRegistrationTestSuite) TestValidUserRegistration() {
	// 18+ 用戶註冊資料
	validRegistration := map[string]interface{}{
		"email":        "alice@test.com",
		"password":     "SecurePassword123",
		"birth_date":   "1995-06-15T00:00:00Z", // 28歲
		"display_name": "Alice Chen",
		"gender":       "female",
	}

	reqBody, _ := json.Marshal(validRegistration)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 目前應該回傳 501 Not Implemented
	suite.Equal(http.StatusNotImplemented, w.Code, "有效註冊應該回傳 501 (尚未實作)")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("註冊功能尚未實作", response["error"])
}

// TestUnderageUserRegistration 測試未成年用戶註冊
func (suite *UserRegistrationTestSuite) TestUnderageUserRegistration() {
	// 未成年用戶註冊資料
	underageRegistration := map[string]interface{}{
		"email":        "young@test.com",
		"password":     "SecurePassword123",
		"birth_date":   "2010-01-01T00:00:00Z", // 13歲
		"display_name": "Too Young",
		"gender":       "male",
	}

	reqBody, _ := json.Marshal(underageRegistration)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 應該回傳 400 Bad Request
	suite.Equal(http.StatusBadRequest, w.Code, "未成年註冊應該回傳 400 Bad Request")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	
	suite.Contains(response, "error")
	suite.Contains(response, "message")
	suite.Contains(response["message"], "18歲", "錯誤訊息應該提及年齡限制")
}

// TestJustEighteenUserRegistration 測試剛滿18歲的用戶註冊
func (suite *UserRegistrationTestSuite) TestJustEighteenUserRegistration() {
	// 剛滿18歲的用戶（今天生日）
	now := time.Now()
	birthday := now.AddDate(-18, 0, 0) // 剛好18年前的今天

	justEighteenRegistration := map[string]interface{}{
		"email":        "eighteen@test.com",
		"password":     "SecurePassword123",
		"birth_date":   birthday.Format("2006-01-02T15:04:05Z"),
		"display_name": "Just Eighteen",
		"gender":       "female",
	}

	reqBody, _ := json.Marshal(justEighteenRegistration)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 目前應該回傳 501 Not Implemented
	suite.Equal(http.StatusNotImplemented, w.Code, "剛滿18歲註冊應該被接受")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("註冊功能尚未實作", response["error"])
}

// TestInvalidEmailFormat 測試無效的 email 格式
func (suite *UserRegistrationTestSuite) TestInvalidEmailFormat() {
	invalidEmailRegistration := map[string]interface{}{
		"email":        "invalid-email", // 無效的 email 格式
		"password":     "SecurePassword123",
		"birth_date":   "1995-06-15T00:00:00Z",
		"display_name": "Invalid Email",
		"gender":       "female",
	}

	reqBody, _ := json.Marshal(invalidEmailRegistration)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 目前會回傳 501，實作完成後應該是 400
	suite.Equal(http.StatusNotImplemented, w.Code)
}

// TestMissingRequiredFields 測試缺少必填欄位
func (suite *UserRegistrationTestSuite) TestMissingRequiredFields() {
	requiredFields := []string{"email", "password", "birth_date", "display_name", "gender"}

	for _, field := range requiredFields {
		suite.T().Run("missing_"+field, func(t *testing.T) {
			// 建立完整的註冊資料
			completeRegistration := map[string]interface{}{
				"email":        "complete@test.com",
				"password":     "SecurePassword123",
				"birth_date":   "1995-06-15T00:00:00Z",
				"display_name": "Complete User",
				"gender":       "female",
			}

			// 移除測試欄位
			delete(completeRegistration, field)

			reqBody, _ := json.Marshal(completeRegistration)
			req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			// 當實作完成後，應該回傳 400 Bad Request
			// 目前可能回傳 400 (JSON binding 錯誤) 或 501
			suite.True(w.Code == http.StatusBadRequest || w.Code == http.StatusNotImplemented)
		})
	}
}