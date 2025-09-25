package redis

import (
	"fmt"
	"time"
)

// CacheService Redis 快取服務
type CacheService struct {
	client *RedisClient
	prefix string
}

// NewCacheService 建立新的快取服務
func NewCacheService(client *RedisClient, prefix string) *CacheService {
	return &CacheService{
		client: client,
		prefix: prefix,
	}
}

// 鍵名生成器
func (c *CacheService) key(parts ...string) string {
	key := c.prefix
	for _, part := range parts {
		key += ":" + part
	}
	return key
}

// 用戶相關快取

// SetUserSession 設置用戶會話
func (c *CacheService) SetUserSession(userID uint, sessionData interface{}, expiration time.Duration) error {
	key := c.key("user", "session", fmt.Sprintf("%d", userID))
	return c.client.SetJSON(key, sessionData, expiration)
}

// GetUserSession 獲取用戶會話
func (c *CacheService) GetUserSession(userID uint, sessionData interface{}) error {
	key := c.key("user", "session", fmt.Sprintf("%d", userID))
	return c.client.GetJSON(key, sessionData)
}

// DeleteUserSession 刪除用戶會話
func (c *CacheService) DeleteUserSession(userID uint) error {
	key := c.key("user", "session", fmt.Sprintf("%d", userID))
	_, err := c.client.Delete(key)
	return err
}

// SetUserProfile 快取用戶資料
func (c *CacheService) SetUserProfile(userID uint, profile interface{}, expiration time.Duration) error {
	key := c.key("user", "profile", fmt.Sprintf("%d", userID))
	return c.client.SetJSON(key, profile, expiration)
}

// GetUserProfile 獲取用戶資料快取
func (c *CacheService) GetUserProfile(userID uint, profile interface{}) error {
	key := c.key("user", "profile", fmt.Sprintf("%d", userID))
	return c.client.GetJSON(key, profile)
}

// InvalidateUserProfile 使用戶資料快取失效
func (c *CacheService) InvalidateUserProfile(userID uint) error {
	key := c.key("user", "profile", fmt.Sprintf("%d", userID))
	_, err := c.client.Delete(key)
	return err
}

// SetUserOnlineStatus 設置用戶在線狀態
func (c *CacheService) SetUserOnlineStatus(userID uint, isOnline bool, expiration time.Duration) error {
	key := c.key("user", "online", fmt.Sprintf("%d", userID))
	return c.client.Set(key, isOnline, expiration)
}

// GetUserOnlineStatus 獲取用戶在線狀態
func (c *CacheService) GetUserOnlineStatus(userID uint) (bool, error) {
	key := c.key("user", "online", fmt.Sprintf("%d", userID))
	result, err := c.client.Get(key)
	if err != nil {
		return false, err
	}
	return result == "true", nil
}

// 配對相關快取

// SetUserSwipeHistory 設置用戶滑動歷史
func (c *CacheService) SetUserSwipeHistory(userID uint, targetUserIDs []uint, expiration time.Duration) error {
	key := c.key("match", "swipe_history", fmt.Sprintf("%d", userID))
	return c.client.SetJSON(key, targetUserIDs, expiration)
}

// GetUserSwipeHistory 獲取用戶滑動歷史
func (c *CacheService) GetUserSwipeHistory(userID uint) ([]uint, error) {
	key := c.key("match", "swipe_history", fmt.Sprintf("%d", userID))
	var history []uint
	err := c.client.GetJSON(key, &history)
	return history, err
}

// AddToSwipeHistory 添加到滑動歷史
func (c *CacheService) AddToSwipeHistory(userID, targetUserID uint) error {
	key := c.key("match", "swipe_history", fmt.Sprintf("%d", userID))

	// 獲取現有歷史
	var history []uint
	c.client.GetJSON(key, &history)

	// 檢查是否已存在
	for _, id := range history {
		if id == targetUserID {
			return nil // 已存在
		}
	}

	// 添加新的用戶ID
	history = append(history, targetUserID)

	// 限制歷史記錄數量（保留最近 1000 個）
	if len(history) > 1000 {
		history = history[len(history)-1000:]
	}

	return c.client.SetJSON(key, history, 24*time.Hour)
}

// SetMatchCandidates 設置配對候選人
func (c *CacheService) SetMatchCandidates(userID uint, candidates interface{}, expiration time.Duration) error {
	key := c.key("match", "candidates", fmt.Sprintf("%d", userID))
	return c.client.SetJSON(key, candidates, expiration)
}

// GetMatchCandidates 獲取配對候選人
func (c *CacheService) GetMatchCandidates(userID uint, candidates interface{}) error {
	key := c.key("match", "candidates", fmt.Sprintf("%d", userID))
	return c.client.GetJSON(key, candidates)
}

// InvalidateMatchCandidates 使配對候選人快取失效
func (c *CacheService) InvalidateMatchCandidates(userID uint) error {
	key := c.key("match", "candidates", fmt.Sprintf("%d", userID))
	_, err := c.client.Delete(key)
	return err
}

// 聊天相關快取

// SetChatMessages 快取聊天訊息
func (c *CacheService) SetChatMessages(chatID uint, messages interface{}, expiration time.Duration) error {
	key := c.key("chat", "messages", fmt.Sprintf("%d", chatID))
	return c.client.SetJSON(key, messages, expiration)
}

// GetChatMessages 獲取聊天訊息快取
func (c *CacheService) GetChatMessages(chatID uint, messages interface{}) error {
	key := c.key("chat", "messages", fmt.Sprintf("%d", chatID))
	return c.client.GetJSON(key, messages)
}

// AddChatMessage 添加聊天訊息到快取
func (c *CacheService) AddChatMessage(chatID uint, message interface{}) error {
	key := c.key("chat", "messages", fmt.Sprintf("%d", chatID))

	// 使用列表推入新訊息
	listKey := c.key("chat", "list", fmt.Sprintf("%d", chatID))
	messageJSON, err := c.client.client.Get(c.client.ctx, key).Result()
	if err == nil {
		// 如果快取存在，更新快取
		var messages []interface{}
		if err := c.client.GetJSON(key, &messages); err == nil {
			messages = append(messages, message)
			// 限制訊息數量（保留最近 100 條）
			if len(messages) > 100 {
				messages = messages[len(messages)-100:]
			}
			return c.client.SetJSON(key, messages, time.Hour)
		}
	}

	// 如果快取不存在，直接推入列表
	_, err = c.client.LPush(listKey, messageJSON)
	if err != nil {
		return err
	}

	// 設置列表過期時間
	return c.client.Expire(listKey, time.Hour)
}

// SetUnreadMessageCount 設置未讀訊息數量
func (c *CacheService) SetUnreadMessageCount(userID, chatID uint, count int64) error {
	key := c.key("chat", "unread", fmt.Sprintf("%d", userID), fmt.Sprintf("%d", chatID))
	return c.client.Set(key, count, 24*time.Hour)
}

// GetUnreadMessageCount 獲取未讀訊息數量
func (c *CacheService) GetUnreadMessageCount(userID, chatID uint) (int64, error) {
	key := c.key("chat", "unread", fmt.Sprintf("%d", userID), fmt.Sprintf("%d", chatID))
	result, err := c.client.Get(key)
	if err != nil {
		return 0, err
	}

	// 轉換字符串為數字
	var count int64
	fmt.Sscanf(result, "%d", &count)
	return count, nil
}

// IncrementUnreadMessageCount 增加未讀訊息數量
func (c *CacheService) IncrementUnreadMessageCount(userID, chatID uint) (int64, error) {
	key := c.key("chat", "unread", fmt.Sprintf("%d", userID), fmt.Sprintf("%d", chatID))
	count, err := c.client.Increment(key)
	if err != nil {
		return 0, err
	}

	// 設置過期時間
	c.client.Expire(key, 24*time.Hour)
	return count, nil
}

// ClearUnreadMessageCount 清空未讀訊息數量
func (c *CacheService) ClearUnreadMessageCount(userID, chatID uint) error {
	key := c.key("chat", "unread", fmt.Sprintf("%d", userID), fmt.Sprintf("%d", chatID))
	_, err := c.client.Delete(key)
	return err
}

// 速率限制相關快取

// SetRateLimit 設置速率限制
func (c *CacheService) SetRateLimit(identifier string, count int64, expiration time.Duration) error {
	key := c.key("rate_limit", identifier)
	return c.client.Set(key, count, expiration)
}

// IncrementRateLimit 增加速率限制計數
func (c *CacheService) IncrementRateLimit(identifier string, expiration time.Duration) (int64, error) {
	key := c.key("rate_limit", identifier)
	count, err := c.client.Increment(key)
	if err != nil {
		return 0, err
	}

	// 如果是新鍵，設置過期時間
	if count == 1 {
		c.client.Expire(key, expiration)
	}

	return count, nil
}

// GetRateLimit 獲取速率限制計數
func (c *CacheService) GetRateLimit(identifier string) (int64, error) {
	key := c.key("rate_limit", identifier)
	result, err := c.client.Get(key)
	if err != nil {
		return 0, err
	}

	var count int64
	fmt.Sscanf(result, "%d", &count)
	return count, nil
}

// JWT Token 黑名單

// AddTokenToBlacklist 添加 token 到黑名單
func (c *CacheService) AddTokenToBlacklist(tokenID string, expiration time.Duration) error {
	key := c.key("jwt", "blacklist", tokenID)
	return c.client.Set(key, "1", expiration)
}

// IsTokenBlacklisted 檢查 token 是否在黑名單中
func (c *CacheService) IsTokenBlacklisted(tokenID string) (bool, error) {
	key := c.key("jwt", "blacklist", tokenID)
	exists, err := c.client.Exists(key)
	return exists > 0, err
}

// 驗證碼快取

// SetVerificationCode 設置驗證碼
func (c *CacheService) SetVerificationCode(identifier, code string, expiration time.Duration) error {
	key := c.key("verification", identifier)
	return c.client.Set(key, code, expiration)
}

// GetVerificationCode 獲取驗證碼
func (c *CacheService) GetVerificationCode(identifier string) (string, error) {
	key := c.key("verification", identifier)
	return c.client.Get(key)
}

// DeleteVerificationCode 刪除驗證碼
func (c *CacheService) DeleteVerificationCode(identifier string) error {
	key := c.key("verification", identifier)
	_, err := c.client.Delete(key)
	return err
}

// 地理位置快取

// SetUserLocation 設置用戶位置
func (c *CacheService) SetUserLocation(userID uint, latitude, longitude float64) error {
	key := c.key("location", "users")
	// 暫時用 hash 存儲位置資訊
	locationData := fmt.Sprintf("%.6f,%.6f", latitude, longitude)
	return c.client.HSet(key, fmt.Sprintf("%d", userID), locationData)
}

// GetNearbyUsers 獲取附近用戶
func (c *CacheService) GetNearbyUsers(userID uint, radiusKM float64) ([]uint, error) {
	// 這裡需要實作 Redis GEO 查詢
	// 暫時返回空列表
	return []uint{}, nil
}

// 統計數據快取

// IncrementDailyStats 增加每日統計
func (c *CacheService) IncrementDailyStats(statType string, date string) (int64, error) {
	key := c.key("stats", "daily", date, statType)
	count, err := c.client.Increment(key)
	if err != nil {
		return 0, err
	}

	// 設置過期時間為 30 天
	c.client.Expire(key, 30*24*time.Hour)
	return count, nil
}

// GetDailyStats 獲取每日統計
func (c *CacheService) GetDailyStats(statType string, date string) (int64, error) {
	key := c.key("stats", "daily", date, statType)
	result, err := c.client.Get(key)
	if err != nil {
		return 0, err
	}

	var count int64
	fmt.Sscanf(result, "%d", &count)
	return count, nil
}

// 清理方法

// ClearUserCache 清空用戶相關快取
func (c *CacheService) ClearUserCache(userID uint) error {
	patterns := []string{
		c.key("user", "session", fmt.Sprintf("%d", userID)),
		c.key("user", "profile", fmt.Sprintf("%d", userID)),
		c.key("user", "online", fmt.Sprintf("%d", userID)),
		c.key("match", "swipe_history", fmt.Sprintf("%d", userID)),
		c.key("match", "candidates", fmt.Sprintf("%d", userID)),
	}

	for _, pattern := range patterns {
		c.client.Delete(pattern)
	}

	return nil
}

// ClearChatCache 清空聊天相關快取
func (c *CacheService) ClearChatCache(chatID uint) error {
	patterns := []string{
		c.key("chat", "messages", fmt.Sprintf("%d", chatID)),
		c.key("chat", "list", fmt.Sprintf("%d", chatID)),
	}

	for _, pattern := range patterns {
		c.client.Delete(pattern)
	}

	return nil
}

// Health 檢查快取服務健康狀態
func (c *CacheService) Health() error {
	return c.client.Ping()
}
