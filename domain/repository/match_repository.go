package repository

import (
	"context"

	"golang_dev_docker/domain/entity"
)

// MatchRepository 配對數據儲存庫介面
// 提供配對系統的持久化操作，包括滑動、配對查詢、推薦演算法等功能
type MatchRepository interface {
	// CreateSwipe 記錄滑動動作
	// 用於記錄用戶對其他用戶的 like/pass 操作
	CreateSwipe(ctx context.Context, match *entity.Match) error

	// GetMatch 獲取配對記錄
	// 用於檢查兩個用戶之間的配對狀態
	GetMatch(ctx context.Context, user1ID, user2ID uint) (*entity.Match, error)

	// GetMatchByID 根據 ID 獲取配對記錄
	// 用於配對詳情查詢和聊天權限驗證
	GetMatchByID(ctx context.Context, id uint) (*entity.Match, error)

	// UpdateMatchStatus 更新配對狀態
	// 用於雙向配對成功時狀態更新
	UpdateMatchStatus(ctx context.Context, matchID uint, status entity.MatchStatus) error

	// ProcessSwipe 處理滑動動作並檢查是否配對成功
	// 包含業務邏輯：檢查對方是否已 like，如是則創建雙向配對
	ProcessSwipe(ctx context.Context, userID, targetUserID uint, action entity.SwipeAction) (*entity.Match, bool, error)

	// GetUserMatches 獲取用戶的所有配對記錄
	// 用於聊天列表和配對歷史展示
	GetUserMatches(ctx context.Context, userID uint, status entity.MatchStatus) ([]*entity.Match, error)

	// GetMatchedUsers 獲取用戶配對成功的用戶列表
	// 用於聊天對象列表展示
	GetMatchedUsers(ctx context.Context, userID uint) ([]*entity.User, error)

	// HasUserSwiped 檢查用戶是否已經滑動過目標用戶
	// 用於避免重複滑動和推薦去重
	HasUserSwiped(ctx context.Context, userID, targetUserID uint) (bool, error)

	// Delete 刪除配對記錄
	// 用於取消配對功能
	Delete(ctx context.Context, id uint) error
}

// MatchingAlgorithmRepository 配對演算法數據儲存庫介面
// 提供配對推薦系統的數據查詢功能，支援地理位置、興趣、年齡等條件篩選
type MatchingAlgorithmRepository interface {
	// GetPotentialMatches 獲取潛在配對對象
	// 根據地理位置、年齡範圍、興趣匹配等條件推薦用戶
	GetPotentialMatches(ctx context.Context, userID uint, params PotentialMatchParams) ([]*entity.User, error)

	// GetUsersNearby 獲取附近的用戶
	// 基於地理位置的用戶推薦
	GetUsersNearby(ctx context.Context, userID uint, lat, lng float64, maxDistanceKm int, limit int) ([]*entity.User, error)

	// GetUsersByAgeRange 根據年齡範圍獲取用戶
	// 用於年齡偏好篩選
	GetUsersByAgeRange(ctx context.Context, userID uint, minAge, maxAge int, limit int) ([]*entity.User, error)

	// GetUsersByCommonInterests 根據共同興趣獲取用戶
	// 用於興趣匹配推薦
	GetUsersByCommonInterests(ctx context.Context, userID uint, limit int) ([]*entity.User, error)

	// CalculateCompatibilityScore 計算用戶相容性分數
	// 用於排序推薦結果
	CalculateCompatibilityScore(ctx context.Context, user1ID, user2ID uint) (float64, error)

	// GetMatchingStats 獲取配對統計數據
	// 用於系統監控和推薦演算法優化
	GetMatchingStats(ctx context.Context, userID uint) (*MatchingStats, error)
}

// PotentialMatchParams 潛在配對查詢參數
type PotentialMatchParams struct {
	// 地理位置篩選
	Latitude    *float64 // 用戶緯度
	Longitude   *float64 // 用戶經度
	MaxDistance *int     // 最大距離（公里）

	// 年齡篩選
	MinAge *int // 最小年齡
	MaxAge *int // 最大年齡

	// 性別偏好
	PreferredGender *entity.Gender // 偏好性別

	// 興趣篩選
	RequireCommonInterests bool // 是否要求共同興趣
	MinCommonInterests     *int // 最少共同興趣數量

	// 分頁參數
	Limit  int // 返回數量限制
	Offset int // 分頁偏移

	// 排除條件
	ExcludeSwipedUsers  bool // 排除已滑動過的用戶
	ExcludeBlockedUsers bool // 排除已封鎖的用戶
}

// MatchingStats 配對統計資料
type MatchingStats struct {
	TotalSwipes    int     // 總滑動次數
	LikesGiven     int     // 給出的 like 數
	LikesReceived  int     // 收到的 like 數
	TotalMatches   int     // 總配對數
	ActiveMatches  int     // 活躍配對數（有聊天記錄）
	MatchRate      float64 // 配對成功率（配對數/滑動數）
	PopularityRate float64 // 受歡迎度（收到like/被滑動）
}
