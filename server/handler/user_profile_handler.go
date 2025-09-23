package handler

import (
	"net/http"
	"strconv"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/service"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	userProfileService *service.UserProfileService
}

func NewUserProfileHandler(userProfileService *service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		userProfileService: userProfileService,
	}
}

// UpdateBasicInfoRequest 更新基本資訊請求
type UpdateBasicInfoRequest struct {
	Age       *int           `json:"age,omitempty"`
	Gender    *entity.Gender `json:"gender,omitempty"`
	Bio       *string        `json:"bio,omitempty"`
	Interests []string       `json:"interests,omitempty"`
}

// UpdateLocationRequest 更新位置請求
type UpdateLocationRequest struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	City      *string  `json:"city,omitempty"`
	Country   *string  `json:"country,omitempty"`
}

// AddPhotoRequest 新增照片請求
type AddPhotoRequest struct {
	PhotoURL  string  `json:"photo_url" binding:"required"`
	Caption   *string `json:"caption,omitempty"`
	IsPrimary bool    `json:"is_primary"`
}

// UpdateBasicInfo 更新用戶基本資訊
func (h *UserProfileHandler) UpdateBasicInfo(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶ID"})
		return
	}

	var req UpdateBasicInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請求格式錯誤", "details": err.Error()})
		return
	}

	// 分別調用不同的服務方法
	if req.Age != nil || req.Gender != nil {
		err = h.userProfileService.UpdateBasicInfo(userID, req.Age, req.Gender)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if req.Bio != nil || req.Interests != nil {
		err = h.userProfileService.UpdateProfileInfo(userID, req.Bio, req.Interests)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "基本資訊更新成功"})
}

// UpdateLocation 更新用戶位置
func (h *UserProfileHandler) UpdateLocation(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶ID"})
		return
	}

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請求格式錯誤", "details": err.Error()})
		return
	}

	err = h.userProfileService.UpdateUserLocation(userID, req.Latitude, req.Longitude, req.City, req.Country)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "位置資訊更新成功"})
}

// AddPhoto 新增用戶照片
func (h *UserProfileHandler) AddPhoto(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶ID"})
		return
	}

	var req AddPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請求格式錯誤", "details": err.Error()})
		return
	}

	err = h.userProfileService.AddUserPhoto(userID, req.PhotoURL, req.Caption, req.IsPrimary)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "照片新增成功"})
}

// FindNearbyUsers 尋找附近用戶
func (h *UserProfileHandler) FindNearbyUsers(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶ID"})
		return
	}

	// 取得查詢參數
	radiusStr := c.DefaultQuery("radius", "10") // 預設10公里
	limitStr := c.DefaultQuery("limit", "20")   // 預設20個結果

	radius, err := strconv.Atoi(radiusStr)
	if err != nil || radius <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的距離範圍"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的結果數量限制"})
		return
	}

	users, err := h.userProfileService.GetNearbyUsers(userID, float64(radius), limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 過濾敏感資訊（如密碼）
	var safeUsers []gin.H
	for _, user := range users {
		safeUser := gin.H{
			"id":             user.ID,
			"username":       user.Username,
			"age":            user.Age,
			"gender":         user.Gender,
			"is_verified":    user.IsVerified,
			"profile_views":  user.ProfileViews,
			"last_active_at": user.LastActiveAt,
			// 注意：bio, interests, city, country 現在在 UserProfile 表中
			// 如果需要這些資訊，需要另外查詢 UserProfile
		}
		safeUsers = append(safeUsers, safeUser)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": safeUsers,
		"count": len(safeUsers),
	})
}

// SearchUsers 搜尋相容用戶
func (h *UserProfileHandler) SearchUsers(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶ID"})
		return
	}

	// 取得查詢參數
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的結果數量限制"})
		return
	}

	users, err := h.userProfileService.SearchCompatibleUsers(userID, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 過濾敏感資訊
	var safeUsers []gin.H
	for _, user := range users {
		safeUser := gin.H{
			"id":             user.ID,
			"username":       user.Username,
			"age":            user.Age,
			"gender":         user.Gender,
			"is_verified":    user.IsVerified,
			"profile_views":  user.ProfileViews,
			"last_active_at": user.LastActiveAt,
			// 注意：bio, interests, city, country 現在在 UserProfile 表中
		}
		safeUsers = append(safeUsers, safeUser)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": safeUsers,
		"count": len(safeUsers),
	})
}

// CreateProfile 創建詳細個人資料
func (h *UserProfileHandler) CreateProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶ID"})
		return
	}

	var profile entity.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請求格式錯誤", "details": err.Error()})
		return
	}

	err = h.userProfileService.CreateUserProfile(userID, &profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "個人資料創建成功"})
}
