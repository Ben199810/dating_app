package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang_dev_docker/domain/usecase"
)

// UserHandler 用戶處理器
type UserHandler struct {
	userService *usecase.UserService
}

// 全域用戶處理器實例
var userHandler *UserHandler

// SetUserService 設置用戶處理器的服務依賴
func SetUserService(userService *usecase.UserService) {
	userHandler = &UserHandler{
		userService: userService,
	}
}

// UpdateProfileRequest 更新檔案請求結構
type UpdateProfileRequest struct {
	DisplayName *string  `json:"display_name,omitempty"`
	Biography   *string  `json:"biography,omitempty"`
	LocationLat *float64 `json:"location_lat,omitempty"`
	LocationLng *float64 `json:"location_lng,omitempty"`
	MaxDistance *int     `json:"max_distance,omitempty"`
	AgeRangeMin *int     `json:"age_range_min,omitempty"`
	AgeRangeMax *int     `json:"age_range_max,omitempty"`
	InterestIDs []uint   `json:"interest_ids,omitempty"`
}

// PhotoUploadRequest 照片上傳請求結構
type PhotoUploadRequest struct {
	ImageURL    string `json:"image_url" binding:"required"`
	Description string `json:"description"`
}

// GetProfileHandler 獲取用戶檔案
// GET /users/profile
func GetProfileHandler(c *gin.Context) {
	if userHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶服務未初始化",
		})
		return
	}

	// 從 JWT token 中獲取用戶 ID (假設已經通過認證中間件)
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

	// 獲取用戶檔案
	userResponse, err := userHandler.userService.GetProfile(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "獲取檔案失敗",
			"message": err.Error(),
		})
		return
	}

	// 獲取用戶照片
	photos, err := userHandler.userService.GetUserPhotos(c.Request.Context(), userIDUint)
	if err != nil {
		// 照片獲取失敗不阻塞主要檔案返回
		photos = nil
	}

	// 獲取用戶興趣
	interests, err := userHandler.userService.GetUserInterests(c.Request.Context(), userIDUint)
	if err != nil {
		// 興趣獲取失敗不阻塞主要檔案返回
		interests = nil
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":          userResponse.ID,
			"email":       userResponse.Email,
			"is_verified": userResponse.IsVerified,
			"age":         userResponse.Age,
			"profile":     userResponse.Profile,
		},
		"photos":    photos,
		"interests": interests,
	})
}

// UpdateProfileHandler 更新用戶檔案
// PUT /users/profile
func UpdateProfileHandler(c *gin.Context) {
	if userHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶服務未初始化",
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

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求資料格式錯誤",
			"details": err.Error(),
		})
		return
	}

	// 構建服務層請求
	serviceReq := &usecase.UpdateProfileRequest{
		DisplayName: req.DisplayName,
		Biography:   req.Biography,
		LocationLat: req.LocationLat,
		LocationLng: req.LocationLng,
		MaxDistance: req.MaxDistance,
		AgeRangeMin: req.AgeRangeMin,
		AgeRangeMax: req.AgeRangeMax,
		InterestIDs: req.InterestIDs,
	}

	// 調用用戶服務更新檔案
	userResponse, err := userHandler.userService.UpdateProfile(c.Request.Context(), userIDUint, serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "更新檔案失敗",
			"message": err.Error(),
		})
		return
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"message": "檔案更新成功",
		"user": gin.H{
			"id":          userResponse.ID,
			"email":       userResponse.Email,
			"is_verified": userResponse.IsVerified,
			"age":         userResponse.Age,
			"profile":     userResponse.Profile,
		},
	})
}

// UploadPhotoHandler 上傳用戶照片
// POST /users/photos
func UploadPhotoHandler(c *gin.Context) {
	if userHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶服務未初始化",
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

	var req PhotoUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求資料格式錯誤",
			"details": err.Error(),
		})
		return
	}

	// 調用用戶服務添加照片
	photo, err := userHandler.userService.AddPhoto(c.Request.Context(), userIDUint, req.ImageURL, req.Description)
	if err != nil {
		if err.Error() == "照片數量已達上限(6張)" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "照片上傳失敗",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "照片上傳失敗",
			"message": err.Error(),
		})
		return
	}

	// 成功回應
	c.JSON(http.StatusCreated, gin.H{
		"message": "照片上傳成功",
		"photo": gin.H{
			"id":        photo.ID,
			"image_url": photo.FilePath,
			"is_main":   photo.IsMain,
			"status":    photo.Status,
		},
	})
}

// GetUserPhotosByIDHandler 根據用戶ID獲取照片 (公開API，用於配對顯示)
// GET /users/:id/photos
func GetUserPhotosByIDHandler(c *gin.Context) {
	if userHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶服務未初始化",
		})
		return
	}

	// 獲取URL參數中的用戶ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用戶ID格式錯誤",
		})
		return
	}

	// 獲取用戶照片
	photos, err := userHandler.userService.GetUserPhotos(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "獲取照片失敗",
			"message": err.Error(),
		})
		return
	}

	// 只返回已通過審核的照片
	var approvedPhotos []map[string]interface{}
	for _, photo := range photos {
		if photo.Status == "approved" { // 假設這是已通過的狀態
			approvedPhotos = append(approvedPhotos, map[string]interface{}{
				"id":        photo.ID,
				"image_url": photo.FilePath,
				"is_main":   photo.IsMain,
				"order":     photo.DisplayOrder,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"photos": approvedPhotos,
	})
}

// DeletePhotoHandler 刪除用戶照片
// DELETE /users/photos/:photoId
func DeletePhotoHandler(c *gin.Context) {
	if userHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用戶服務未初始化",
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

	// 獲取URL參數中的照片ID
	photoIDStr := c.Param("photoId")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "照片ID格式錯誤",
		})
		return
	}

	// 調用用戶服務刪除照片
	err = userHandler.userService.DeletePhoto(c.Request.Context(), userIDUint, uint(photoID))
	if err != nil {
		if err.Error() == "無權限操作此照片" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "無權限操作",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "刪除失敗",
			"message": err.Error(),
		})
		return
	}

	// 成功回應
	c.JSON(http.StatusOK, gin.H{
		"message": "照片刪除成功",
	})
}
