package redis

import (
	"encoding/json"
	"fmt"
	"time"
)

// SessionCacheService 會話快取服務
type SessionCacheService struct {
	client *RedisClient
	ttl    time.Duration
}

// NewSessionCacheService 建立新的會話快取服務
func NewSessionCacheService(client *RedisClient) *SessionCacheService {
	return &SessionCacheService{
		client: client,
		ttl:    24 * time.Hour, // 預設 24 小時過期
	}
}

// SessionData 會話資料結構
type SessionData struct {
	UserID      uint      `json:"user_id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	IsVerified  bool      `json:"is_verified"`
	LoginTime   time.Time `json:"login_time"`
	LastActive  time.Time `json:"last_active"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
}

// StoreSession 儲存用戶會話
func (s *SessionCacheService) StoreSession(sessionID string, data *SessionData) error {
	key := s.sessionKey(sessionID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化會話資料失敗: %w", err)
	}

	return s.client.Set(key, string(jsonData), s.ttl)
}

// GetSession 獲取用戶會話
func (s *SessionCacheService) GetSession(sessionID string) (*SessionData, error) {
	key := s.sessionKey(sessionID)

	result, err := s.client.Get(key)
	if err != nil {
		return nil, fmt.Errorf("獲取會話失敗: %w", err)
	}

	var data SessionData
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return nil, fmt.Errorf("反序列化會話資料失敗: %w", err)
	}

	return &data, nil
}

// DeleteSession 刪除用戶會話
func (s *SessionCacheService) DeleteSession(sessionID string) error {
	key := s.sessionKey(sessionID)
	_, err := s.client.Delete(key)
	return err
}

// RefreshSession 刷新會話過期時間
func (s *SessionCacheService) RefreshSession(sessionID string) error {
	key := s.sessionKey(sessionID)
	return s.client.Expire(key, s.ttl)
}

// UpdateLastActive 更新最後活躍時間
func (s *SessionCacheService) UpdateLastActive(sessionID string) error {
	// 獲取現有會話資料
	data, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	// 更新最後活躍時間
	data.LastActive = time.Now()

	// 重新儲存
	return s.StoreSession(sessionID, data)
}

// GetActiveUserSessions 獲取用戶的所有活躍會話
func (s *SessionCacheService) GetActiveUserSessions(userID uint) ([]string, error) {
	pattern := s.userSessionPattern(userID)
	return s.client.Keys(pattern)
}

// DeleteAllUserSessions 刪除用戶的所有會話
func (s *SessionCacheService) DeleteAllUserSessions(userID uint) error {
	sessions, err := s.GetActiveUserSessions(userID)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		return nil
	}

	_, err = s.client.Delete(sessions...)
	return err
}

// SetSessionTTL 設定會話過期時間
func (s *SessionCacheService) SetSessionTTL(ttl time.Duration) {
	s.ttl = ttl
}

// IsSessionActive 檢查會話是否活躍
func (s *SessionCacheService) IsSessionActive(sessionID string) bool {
	key := s.sessionKey(sessionID)
	count, err := s.client.Exists(key)
	return err == nil && count > 0
}

// GetSessionTTL 獲取會話剩餘時間
func (s *SessionCacheService) GetSessionTTL(sessionID string) (time.Duration, error) {
	key := s.sessionKey(sessionID)
	return s.client.TTL(key)
}

// CleanupExpiredSessions 清理過期會話（定期執行）
func (s *SessionCacheService) CleanupExpiredSessions() error {
	// Redis 會自動處理過期鍵，這個方法主要用於統計
	pattern := s.sessionKey("*")
	allSessions, err := s.client.Keys(pattern)
	if err != nil {
		return err
	}

	activeCount := 0
	for _, sessionKey := range allSessions {
		count, err := s.client.Exists(sessionKey)
		if err == nil && count > 0 {
			activeCount++
		}
	}

	fmt.Printf("會話清理完成，目前活躍會話數: %d\n", activeCount)
	return nil
}

// GetSessionStats 獲取會話統計資訊
func (s *SessionCacheService) GetSessionStats() (map[string]int, error) {
	pattern := s.sessionKey("*")
	allSessions, err := s.client.Keys(pattern)
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"total_sessions":  len(allSessions),
		"active_sessions": 0,
	}

	for _, sessionKey := range allSessions {
		count, err := s.client.Exists(sessionKey)
		if err == nil && count > 0 {
			stats["active_sessions"]++
		}
	}

	return stats, nil
}

// 私有方法

func (s *SessionCacheService) sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func (s *SessionCacheService) userSessionKey(userID uint, sessionID string) string {
	return fmt.Sprintf("user:%d:session:%s", userID, sessionID)
}

func (s *SessionCacheService) userSessionPattern(userID uint) string {
	return fmt.Sprintf("user:%d:session:*", userID)
}

// StoreUserSession 儲存帶用戶 ID 關聯的會話
func (s *SessionCacheService) StoreUserSession(userID uint, sessionID string, data *SessionData) error {
	// 儲存主會話
	if err := s.StoreSession(sessionID, data); err != nil {
		return err
	}

	// 建立用戶到會話的關聯
	userSessionKey := s.userSessionKey(userID, sessionID)
	return s.client.Set(userSessionKey, sessionID, s.ttl)
}

// ValidateSession 驗證會話有效性
func (s *SessionCacheService) ValidateSession(sessionID string) (*SessionData, bool, error) {
	data, err := s.GetSession(sessionID)
	if err != nil {
		return nil, false, err
	}

	// 檢查會話是否過期（額外檢查）
	if time.Since(data.LastActive) > s.ttl {
		s.DeleteSession(sessionID)
		return nil, false, nil
	}

	return data, true, nil
}
