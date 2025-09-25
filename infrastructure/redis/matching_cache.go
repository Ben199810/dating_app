package redis

import (
	"fmt"
	"time"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

// MatchingCacheService 配對快取服務
// 負責快取配對相關資料以提升系統性能
type MatchingCacheService struct {
	client *RedisClient
	ttl    time.Duration // 快取過期時間
}

// NewMatchingCacheService 創建配對快取服務實例
func NewMatchingCacheService(client *RedisClient) *MatchingCacheService {
	return &MatchingCacheService{
		client: client,
		ttl:    time.Hour * 24, // 預設24小時過期
	}
}

// SetTTL 設定快取過期時間
func (m *MatchingCacheService) SetTTL(ttl time.Duration) {
	m.ttl = ttl
}

// CachePotentialMatches 快取潛在配對候選人列表
func (m *MatchingCacheService) CachePotentialMatches(userID uint, matches []*entity.User) error {
	key := m.potentialMatchesKey(userID)
	return m.client.SetJSON(key, matches, m.ttl)
}

// GetPotentialMatches 獲取快取的潛在配對候選人
func (m *MatchingCacheService) GetPotentialMatches(userID uint) ([]*entity.User, error) {
	key := m.potentialMatchesKey(userID)
	var matches []*entity.User

	err := m.client.GetJSON(key, &matches)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// InvalidatePotentialMatches 清除潛在配對快取
func (m *MatchingCacheService) InvalidatePotentialMatches(userID uint) error {
	key := m.potentialMatchesKey(userID)
	_, err := m.client.Delete(key)
	return err
}

// CacheUserMatches 快取用戶配對列表
func (m *MatchingCacheService) CacheUserMatches(userID uint, status entity.MatchStatus, matches []*entity.Match) error {
	key := m.userMatchesKey(userID, status)
	return m.client.SetJSON(key, matches, m.ttl)
}

// GetUserMatches 獲取快取的用戶配對列表
func (m *MatchingCacheService) GetUserMatches(userID uint, status entity.MatchStatus) ([]*entity.Match, error) {
	key := m.userMatchesKey(userID, status)
	var matches []*entity.Match

	err := m.client.GetJSON(key, &matches)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// InvalidateUserMatches 清除用戶配對快取
func (m *MatchingCacheService) InvalidateUserMatches(userID uint, status entity.MatchStatus) error {
	key := m.userMatchesKey(userID, status)
	_, err := m.client.Delete(key)
	return err
}

// InvalidateAllUserMatches 清除用戶所有配對狀態的快取
func (m *MatchingCacheService) InvalidateAllUserMatches(userID uint) error {
	pattern := fmt.Sprintf("matching:matches:%d:*", userID)
	keys, err := m.client.Keys(pattern)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		_, err = m.client.Delete(keys...)
		return err
	}

	return nil
}

// CacheMatchingStats 快取配對統計資料
func (m *MatchingCacheService) CacheMatchingStats(userID uint, stats *repository.MatchingStats) error {
	key := m.matchingStatsKey(userID)
	return m.client.SetJSON(key, stats, time.Hour*6) // 統計資料6小時過期
}

// GetMatchingStats 獲取快取的配對統計資料
func (m *MatchingCacheService) GetMatchingStats(userID uint) (*repository.MatchingStats, error) {
	key := m.matchingStatsKey(userID)
	var stats repository.MatchingStats

	err := m.client.GetJSON(key, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// InvalidateMatchingStats 清除配對統計快取
func (m *MatchingCacheService) InvalidateMatchingStats(userID uint) error {
	key := m.matchingStatsKey(userID)
	_, err := m.client.Delete(key)
	return err
}

// CacheCompatibilityScore 快取相容性分數
func (m *MatchingCacheService) CacheCompatibilityScore(user1ID, user2ID uint, score float64) error {
	key := m.compatibilityScoreKey(user1ID, user2ID)
	return m.client.SetJSON(key, score, time.Hour*12) // 相容性分數12小時過期
}

// GetCompatibilityScore 獲取快取的相容性分數
func (m *MatchingCacheService) GetCompatibilityScore(user1ID, user2ID uint) (float64, error) {
	key := m.compatibilityScoreKey(user1ID, user2ID)
	var score float64

	err := m.client.GetJSON(key, &score)
	if err != nil {
		return 0, err
	}

	return score, nil
}

// CacheNearbyUsers 快取附近用戶列表
func (m *MatchingCacheService) CacheNearbyUsers(userID uint, maxDistance int, users []*entity.User) error {
	key := m.nearbyUsersKey(userID, maxDistance)
	return m.client.SetJSON(key, users, time.Minute*30) // 附近用戶30分鐘過期
}

// GetNearbyUsers 獲取快取的附近用戶列表
func (m *MatchingCacheService) GetNearbyUsers(userID uint, maxDistance int) ([]*entity.User, error) {
	key := m.nearbyUsersKey(userID, maxDistance)
	var users []*entity.User

	err := m.client.GetJSON(key, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// CacheCommonInterestUsers 快取共同興趣用戶列表
func (m *MatchingCacheService) CacheCommonInterestUsers(userID uint, users []*entity.User) error {
	key := m.commonInterestUsersKey(userID)
	return m.client.SetJSON(key, users, time.Hour*2) // 共同興趣用戶2小時過期
}

// GetCommonInterestUsers 獲取快取的共同興趣用戶列表
func (m *MatchingCacheService) GetCommonInterestUsers(userID uint) ([]*entity.User, error) {
	key := m.commonInterestUsersKey(userID)
	var users []*entity.User

	err := m.client.GetJSON(key, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// InvalidateUserCache 清除特定用戶的所有配對快取
func (m *MatchingCacheService) InvalidateUserCache(userID uint) error {
	// 清除潛在配對
	if err := m.InvalidatePotentialMatches(userID); err != nil {
		fmt.Printf("清除潛在配對快取失敗: %v\n", err)
	}

	// 清除所有配對狀態快取
	if err := m.InvalidateAllUserMatches(userID); err != nil {
		fmt.Printf("清除用戶配對快取失敗: %v\n", err)
	}

	// 清除配對統計
	if err := m.InvalidateMatchingStats(userID); err != nil {
		fmt.Printf("清除配對統計快取失敗: %v\n", err)
	}

	// 清除附近用戶快取（清除該用戶相關的所有距離範圍）
	pattern := fmt.Sprintf("matching:nearby:%d:*", userID)
	keys, err := m.client.Keys(pattern)
	if err == nil && len(keys) > 0 {
		if _, err := m.client.Delete(keys...); err != nil {
			fmt.Printf("清除附近用戶快取失敗: %v\n", err)
		}
	}

	// 清除共同興趣用戶快取
	key := m.commonInterestUsersKey(userID)
	if _, err := m.client.Delete(key); err != nil {
		fmt.Printf("清除共同興趣用戶快取失敗: %v\n", err)
	}

	return nil
}

// GetCacheStats 獲取配對快取統計資訊
func (m *MatchingCacheService) GetCacheStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 統計各類快取數量
	patterns := map[string]string{
		"potential_matches":     "matching:potential:*",
		"user_matches":          "matching:matches:*",
		"matching_stats":        "matching:stats:*",
		"compatibility_scores":  "matching:compatibility:*",
		"nearby_users":          "matching:nearby:*",
		"common_interest_users": "matching:interests:*",
	}

	for category, pattern := range patterns {
		keys, err := m.client.Keys(pattern)
		if err != nil {
			stats[category] = 0
		} else {
			stats[category] = len(keys)
		}
	}

	return stats, nil
}

// 私有方法 - 生成快取鍵值

func (m *MatchingCacheService) potentialMatchesKey(userID uint) string {
	return fmt.Sprintf("matching:potential:%d", userID)
}

func (m *MatchingCacheService) userMatchesKey(userID uint, status entity.MatchStatus) string {
	return fmt.Sprintf("matching:matches:%d:%s", userID, status)
}

func (m *MatchingCacheService) matchingStatsKey(userID uint) string {
	return fmt.Sprintf("matching:stats:%d", userID)
}

func (m *MatchingCacheService) compatibilityScoreKey(user1ID, user2ID uint) string {
	// 確保鍵值的一致性，較小的ID在前
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}
	return fmt.Sprintf("matching:compatibility:%d:%d", user1ID, user2ID)
}

func (m *MatchingCacheService) nearbyUsersKey(userID uint, maxDistance int) string {
	return fmt.Sprintf("matching:nearby:%d:%d", userID, maxDistance)
}

func (m *MatchingCacheService) commonInterestUsersKey(userID uint) string {
	return fmt.Sprintf("matching:interests:%d", userID)
}
