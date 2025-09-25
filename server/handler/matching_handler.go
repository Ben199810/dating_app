package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/usecase"
)

// MatchingHandler 配對處理器
type MatchingHandler struct {
	matchingService *usecase.MatchingService
}

// 全域配對處理器實例
var matchingHandler *MatchingHandler

// SetMatchingService 設置配對處理器的服務依賴
func SetMatchingService(matchingService *usecase.MatchingService) {
	matchingHandler = &MatchingHandler{
		matchingService: matchingService,
	}
}

// SwipeRequest 滑動請求結構
type SwipeRequest struct {
	TargetUserID uint   `json:"target_user_id" binding:"required"`
	Action       string `json:"action" binding:"required"` // "like" or "pass"
}

// GetPotentialMatchesHandler 獲取潛在配對對象
// GET /matching/potential
func GetPotentialMatchesHandler(c *gin.Context) {
	if matchingHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "配對服務未初始化",
		})
		return
	}

	// 從 JWT token 中獲取用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶ID格式錯誤",
		})
		return
	}

	// 獲取查詢參數
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10 // 預設限制10個
	}

	// 構建請求參數
	req := &usecase.PotentialMatchRequest{
		UserID: userIDUint,
		Limit:  limit,
	}

	// 解析可選參數
	if maxDistanceStr := c.Query("max_distance"); maxDistanceStr != "" {
		if maxDistance, err := strconv.Atoi(maxDistanceStr); err == nil {
			req.MaxDistance = &maxDistance
		}
	}

	if minAgeStr := c.Query("min_age"); minAgeStr != "" {
		if minAge, err := strconv.Atoi(minAgeStr); err == nil {
			req.MinAge = &minAge
		}
	}

	if maxAgeStr := c.Query("max_age"); maxAgeStr != "" {
		if maxAge, err := strconv.Atoi(maxAgeStr); err == nil {
			req.MaxAge = &maxAge
		}
	}

	if genderStr := c.Query("preferred_gender"); genderStr != "" {
		gender := entity.Gender(genderStr)
		if gender.IsValid() {
			req.PreferredGender = &gender
		}
	}

	if commonInterestsStr := c.Query("require_common_interests"); commonInterestsStr == "true" {
		req.RequireCommonInterests = true
		if minCommonStr := c.Query("min_common_interests"); minCommonStr != "" {
			if minCommon, err := strconv.Atoi(minCommonStr); err == nil {
				req.MinCommonInterests = &minCommon
			}
		}
	}

	// 調用配對服務獲取潛在配對
	potentialUsers, err := matchingHandler.matchingService.GetPotentialMatches(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "用戶未啟用或未驗證" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "帳戶狀態異常",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "獲取潛在配對失敗",
			"message": err.Error(),
		})
		return
	}

	// 構建回應資料
	var users []map[string]interface{}
	for _, user := range potentialUsers {
		users = append(users, map[string]interface{}{
			"id":    user.ID,
			"age":   user.GetAge(),
			"email": user.Email, // 可能需要隱藏或部分隱藏
		})
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"potential_matches": users,
		"total_count":       len(users),
		"limit":             limit,
	})
}

// SwipeHandler 處理滑動配對
// POST /matching/swipe
func SwipeHandler(c *gin.Context) {
	if matchingHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "配對服務未初始化",
		})
		return
	}

	// 從 JWT token 中獲取用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶ID格式錯誤",
		})
		return
	}

	var req SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求資料格式錯誤",
			"details": err.Error(),
		})
		return
	}

	// 驗證動作參數
	var action entity.SwipeAction
	switch req.Action {
	case "like":
		action = entity.SwipeActionLike
	case "pass":
		action = entity.SwipeActionPass
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "無效的滑動動作",
			"message": "動作必須是 'like' 或 'pass'",
		})
		return
	}

	// 構建服務層請求
	serviceReq := &usecase.SwipeRequest{
		UserID:       userIDUint,
		TargetUserID: req.TargetUserID,
		Action:       action,
	}

	// 調用配對服務處理滑動
	swipeResponse, err := matchingHandler.matchingService.ProcessSwipe(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "滑動處理失敗",
			"message": err.Error(),
		})
		return
	}

	// 根據服務回應構建 HTTP 回應
	if !swipeResponse.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": swipeResponse.Message,
		})
		return
	}

	response := gin.H{
		"success":  true,
		"is_match": swipeResponse.IsMatch,
		"message":  swipeResponse.Message,
	}

	// 如果配對成功，返回配對資訊
	if swipeResponse.IsMatch && swipeResponse.Match != nil {
		response["match"] = gin.H{
			"id":         swipeResponse.Match.ID,
			"matched_at": swipeResponse.Match.MatchedAt,
			"status":     swipeResponse.Match.Status,
		}

		// 設置 HTTP 狀態為 201 Created 表示新配對產生
		c.JSON(http.StatusCreated, response)
		return
	}

	// 成功回應
	c.JSON(http.StatusOK, response)
}

// GetUserMatchesHandler 獲取用戶配對列表
// GET /matching/matches
func GetUserMatchesHandler(c *gin.Context) {
	if matchingHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "配對服務未初始化",
		})
		return
	}

	// 從 JWT token 中獲取用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶ID格式錯誤",
		})
		return
	}

	// 獲取狀態過濾參數
	statusStr := c.DefaultQuery("status", "matched")
	var status entity.MatchStatus
	switch statusStr {
	case "matched":
		status = entity.MatchStatusMatched
	case "pending":
		status = entity.MatchStatusPending
	case "unmatched":
		status = entity.MatchStatusUnmatched
	default:
		status = entity.MatchStatusMatched // 預設為配對成功狀態
	}

	// 調用配對服務獲取配對列表
	matchList, err := matchingHandler.matchingService.GetUserMatches(c.Request.Context(), userIDUint, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "獲取配對列表失敗",
			"message": err.Error(),
		})
		return
	}

	// 構建回應資料
	var matches []map[string]interface{}
	for _, match := range matchList.Matches {
		matches = append(matches, map[string]interface{}{
			"id":         match.ID,
			"user1_id":   match.User1ID,
			"user2_id":   match.User2ID,
			"status":     match.Status,
			"matched_at": match.MatchedAt,
		})
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"matches":     matches,
		"total_count": matchList.TotalCount,
		"status":      statusStr,
	})
}
