package usecase

import (
	"context"
	"errors"
	"fmt"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

// MatchingService 配對業務邏輯服務
// 負責配對演算法、滑動處理、推薦系統等核心業務邏輯
type MatchingService struct {
	matchRepo     repository.MatchRepository
	algorithmRepo repository.MatchingAlgorithmRepository
	userRepo      repository.UserRepository
	profileRepo   repository.UserProfileRepository
}

// NewMatchingService 創建新的配對服務實例
func NewMatchingService(
	matchRepo repository.MatchRepository,
	algorithmRepo repository.MatchingAlgorithmRepository,
	userRepo repository.UserRepository,
	profileRepo repository.UserProfileRepository,
) *MatchingService {
	return &MatchingService{
		matchRepo:     matchRepo,
		algorithmRepo: algorithmRepo,
		userRepo:      userRepo,
		profileRepo:   profileRepo,
	}
}

// SwipeRequest 滑動請求
type SwipeRequest struct {
	UserID       uint               `json:"user_id" validate:"required"`
	TargetUserID uint               `json:"target_user_id" validate:"required"`
	Action       entity.SwipeAction `json:"action" validate:"required"`
}

// SwipeResponse 滑動回應
type SwipeResponse struct {
	Success bool          `json:"success"`
	IsMatch bool          `json:"is_match"`
	Match   *entity.Match `json:"match,omitempty"`
	Message string        `json:"message"`
}

// PotentialMatchRequest 潛在配對請求
type PotentialMatchRequest struct {
	UserID                 uint           `json:"user_id" validate:"required"`
	Limit                  int            `json:"limit"`
	MaxDistance            *int           `json:"max_distance,omitempty"`
	MinAge                 *int           `json:"min_age,omitempty"`
	MaxAge                 *int           `json:"max_age,omitempty"`
	PreferredGender        *entity.Gender `json:"preferred_gender,omitempty"`
	RequireCommonInterests bool           `json:"require_common_interests"`
	MinCommonInterests     *int           `json:"min_common_interests,omitempty"`
}

// MatchListResponse 配對列表回應
type MatchListResponse struct {
	Matches    []*entity.Match `json:"matches"`
	TotalCount int             `json:"total_count"`
}

// ProcessSwipe 處理用戶滑動動作
// 記錄滑動並檢查是否形成雙向配對
func (s *MatchingService) ProcessSwipe(ctx context.Context, req *SwipeRequest) (*SwipeResponse, error) {
	// 驗證請求
	if err := s.validateSwipeRequest(req); err != nil {
		return nil, fmt.Errorf("滑動請求驗證失敗: %w", err)
	}

	// 檢查是否嘗試對自己滑動
	if req.UserID == req.TargetUserID {
		return &SwipeResponse{
			Success: false,
			IsMatch: false,
			Message: "不能對自己進行滑動操作",
		}, nil
	}

	// 檢查目標用戶是否存在且啟用
	targetUser, err := s.userRepo.GetByID(ctx, req.TargetUserID)
	if err != nil {
		return &SwipeResponse{
			Success: false,
			IsMatch: false,
			Message: "目標用戶不存在",
		}, nil
	}

	if !targetUser.IsActive || !targetUser.IsVerified {
		return &SwipeResponse{
			Success: false,
			IsMatch: false,
			Message: "目標用戶未啟用或未驗證",
		}, nil
	}

	// 檢查是否已經滑動過
	hasSwipped, err := s.matchRepo.HasUserSwiped(ctx, req.UserID, req.TargetUserID)
	if err != nil {
		return nil, fmt.Errorf("檢查滑動歷史失敗: %w", err)
	}

	if hasSwipped {
		return &SwipeResponse{
			Success: false,
			IsMatch: false,
			Message: "已經對該用戶進行過滑動操作",
		}, nil
	}

	// 處理滑動動作
	match, isMatch, err := s.matchRepo.ProcessSwipe(ctx, req.UserID, req.TargetUserID, req.Action)
	if err != nil {
		return nil, fmt.Errorf("處理滑動失敗: %w", err)
	}

	response := &SwipeResponse{
		Success: true,
		IsMatch: isMatch,
		Match:   match,
	}

	if isMatch {
		response.Message = "恭喜！你們配對成功了"
	} else if req.Action == entity.SwipeActionLike {
		response.Message = "已送出喜歡"
	} else {
		response.Message = "已跳過"
	}

	return response, nil
}

// GetPotentialMatches 獲取潛在配對對象
// 根據用戶偏好和配對演算法推薦合適的配對對象
func (s *MatchingService) GetPotentialMatches(ctx context.Context, req *PotentialMatchRequest) ([]*entity.User, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("用戶不存在: %w", err)
	}

	if !user.IsActive || !user.IsVerified {
		return nil, errors.New("用戶未啟用或未驗證")
	}

	// 獲取用戶檔案（用於推薦參數）
	profile, err := s.profileRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("獲取用戶檔案失敗: %w", err)
	}

	// 建立查詢參數
	params := s.buildMatchingParams(profile, req)

	// 獲取潛在配對
	potentialMatches, err := s.algorithmRepo.GetPotentialMatches(ctx, req.UserID, params)
	if err != nil {
		return nil, fmt.Errorf("獲取潛在配對失敗: %w", err)
	}

	// 根據相容性分數排序（已在 repository 層處理）
	return potentialMatches, nil
}

// GetUserMatches 獲取用戶的配對列表
// 返回用戶所有配對成功的記錄
func (s *MatchingService) GetUserMatches(ctx context.Context, userID uint, status entity.MatchStatus) (*MatchListResponse, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("用戶不存在: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("用戶未啟用")
	}

	// 獲取配對記錄
	matches, err := s.matchRepo.GetUserMatches(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("獲取配對記錄失敗: %w", err)
	}

	return &MatchListResponse{
		Matches:    matches,
		TotalCount: len(matches),
	}, nil
}

// GetMatchedUsers 獲取配對成功的用戶列表
// 用於聊天對象列表展示
func (s *MatchingService) GetMatchedUsers(ctx context.Context, userID uint) ([]*entity.User, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("用戶不存在: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("用戶未啟用")
	}

	return s.matchRepo.GetMatchedUsers(ctx, userID)
}

// UnmatchUser 取消配對
// 允許用戶取消與某個用戶的配對關係
func (s *MatchingService) UnmatchUser(ctx context.Context, userID, targetUserID uint) error {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("用戶不存在: %w", err)
	}

	if !user.IsActive {
		return errors.New("用戶未啟用")
	}

	// 查找配對記錄
	match, err := s.matchRepo.GetMatch(ctx, userID, targetUserID)
	if err != nil {
		return fmt.Errorf("配對記錄不存在: %w", err)
	}

	// 檢查配對狀態
	if match.Status != entity.MatchStatusMatched {
		return errors.New("只能取消已配對成功的關係")
	}

	// 更新配對狀態為未配對
	return s.matchRepo.UpdateMatchStatus(ctx, match.ID, entity.MatchStatusUnmatched)
}

// GetMatchingStats 獲取用戶配對統計
// 提供用戶配對數據的分析統計
func (s *MatchingService) GetMatchingStats(ctx context.Context, userID uint) (*repository.MatchingStats, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("用戶不存在: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("用戶未啟用")
	}

	return s.algorithmRepo.GetMatchingStats(ctx, userID)
}

// CalculateCompatibilityScore 計算兩個用戶的相容性分數
// 用於推薦演算法和配對品質評估
func (s *MatchingService) CalculateCompatibilityScore(ctx context.Context, user1ID, user2ID uint) (float64, error) {
	// 驗證用戶存在
	if _, err := s.userRepo.GetByID(ctx, user1ID); err != nil {
		return 0, fmt.Errorf("用戶1不存在: %w", err)
	}

	if _, err := s.userRepo.GetByID(ctx, user2ID); err != nil {
		return 0, fmt.Errorf("用戶2不存在: %w", err)
	}

	return s.algorithmRepo.CalculateCompatibilityScore(ctx, user1ID, user2ID)
}

// GetUsersNearby 獲取附近用戶
// 基於地理位置推薦附近的用戶
func (s *MatchingService) GetUsersNearby(ctx context.Context, userID uint, maxDistance int, limit int) ([]*entity.User, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("用戶不存在: %w", err)
	}

	if !user.IsActive || !user.IsVerified {
		return nil, errors.New("用戶未啟用或未驗證")
	}

	// 獲取用戶檔案位置資訊
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("獲取用戶檔案失敗: %w", err)
	}

	if profile.LocationLat == nil || profile.LocationLng == nil {
		return nil, errors.New("用戶尚未設定位置資訊")
	}

	return s.algorithmRepo.GetUsersNearby(ctx, userID, *profile.LocationLat, *profile.LocationLng, maxDistance, limit)
}

// GetUsersByCommonInterests 根據共同興趣推薦用戶
// 基於興趣匹配推薦有共同愛好的用戶
func (s *MatchingService) GetUsersByCommonInterests(ctx context.Context, userID uint, limit int) ([]*entity.User, error) {
	// 驗證用戶存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("用戶不存在: %w", err)
	}

	if !user.IsActive || !user.IsVerified {
		return nil, errors.New("用戶未啟用或未驗證")
	}

	return s.algorithmRepo.GetUsersByCommonInterests(ctx, userID, limit)
}

// 私有輔助方法

// validateSwipeRequest 驗證滑動請求
func (s *MatchingService) validateSwipeRequest(req *SwipeRequest) error {
	if req.UserID == 0 {
		return errors.New("用戶ID不能為空")
	}

	if req.TargetUserID == 0 {
		return errors.New("目標用戶ID不能為空")
	}

	if !req.Action.IsValid() {
		return errors.New("無效的滑動動作")
	}

	return nil
}

// buildMatchingParams 建立配對查詢參數
func (s *MatchingService) buildMatchingParams(profile *entity.UserProfile, req *PotentialMatchRequest) repository.PotentialMatchParams {
	params := repository.PotentialMatchParams{
		Limit:                  req.Limit,
		RequireCommonInterests: req.RequireCommonInterests,
		ExcludeSwipedUsers:     true, // 默認排除已滑動的用戶
		ExcludeBlockedUsers:    true, // 默認排除被封鎖的用戶
	}

	// 設定預設值
	if params.Limit <= 0 {
		params.Limit = 10 // 默認返回10個推薦
	}
	if params.Limit > 50 {
		params.Limit = 50 // 最大限制50個
	}

	// 地理位置篩選
	if profile.LocationLat != nil && profile.LocationLng != nil {
		params.Latitude = profile.LocationLat
		params.Longitude = profile.LocationLng

		// 使用請求中的距離或用戶檔案中的偏好距離
		if req.MaxDistance != nil {
			params.MaxDistance = req.MaxDistance
		} else {
			params.MaxDistance = &profile.MaxDistance
		}
	}

	// 年齡篩選
	if req.MinAge != nil {
		params.MinAge = req.MinAge
	} else {
		params.MinAge = &profile.AgeRangeMin
	}

	if req.MaxAge != nil {
		params.MaxAge = req.MaxAge
	} else {
		params.MaxAge = &profile.AgeRangeMax
	}

	// 性別偏好
	if req.PreferredGender != nil {
		params.PreferredGender = req.PreferredGender
	}

	// 共同興趣要求
	if req.MinCommonInterests != nil {
		params.MinCommonInterests = req.MinCommonInterests
	}

	return params
}
