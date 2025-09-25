package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 速率限制器接口
type RateLimiter interface {
	Allow(key string) bool
	Reset(key string)
	GetStats(key string) RateLimitStats
}

// RateLimitStats 速率限制統計
type RateLimitStats struct {
	Requests   int           `json:"requests"`
	Remaining  int           `json:"remaining"`
	ResetTime  time.Time     `json:"reset_time"`
	RetryAfter time.Duration `json:"retry_after"`
}

// TokenBucketLimiter 令牌桶限制器
type TokenBucketLimiter struct {
	buckets  map[string]*TokenBucket
	mu       sync.RWMutex
	rate     int           // 每秒補充的令牌數
	capacity int           // 桶的容量
	window   time.Duration // 時間窗口
}

// TokenBucket 令牌桶
type TokenBucket struct {
	tokens     int       // 當前令牌數
	capacity   int       // 桶容量
	refillRate int       // 每秒補充的令牌數
	lastRefill time.Time // 最後補充時間
}

// NewTokenBucketLimiter 建立新的令牌桶限制器
func NewTokenBucketLimiter(rate, capacity int, window time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		buckets:  make(map[string]*TokenBucket),
		rate:     rate,
		capacity: capacity,
		window:   window,
	}
}

// Allow 檢查是否允許請求
func (t *TokenBucketLimiter) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	bucket, exists := t.buckets[key]
	if !exists {
		bucket = &TokenBucket{
			tokens:     t.capacity,
			capacity:   t.capacity,
			refillRate: t.rate,
			lastRefill: time.Now(),
		}
		t.buckets[key] = bucket
	}

	// 補充令牌
	now := time.Now()
	duration := now.Sub(bucket.lastRefill)
	tokensToAdd := int(duration.Seconds()) * bucket.refillRate

	if tokensToAdd > 0 {
		bucket.tokens += tokensToAdd
		if bucket.tokens > bucket.capacity {
			bucket.tokens = bucket.capacity
		}
		bucket.lastRefill = now
	}

	// 檢查是否有令牌可用
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// Reset 重置指定鍵的限制
func (t *TokenBucketLimiter) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.buckets, key)
}

// GetStats 獲取速率限制統計
func (t *TokenBucketLimiter) GetStats(key string) RateLimitStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	bucket, exists := t.buckets[key]
	if !exists {
		return RateLimitStats{
			Requests:   0,
			Remaining:  t.capacity,
			ResetTime:  time.Now().Add(t.window),
			RetryAfter: 0,
		}
	}

	// 計算重試時間
	retryAfter := time.Duration(0)
	if bucket.tokens == 0 {
		retryAfter = time.Second / time.Duration(bucket.refillRate)
	}

	return RateLimitStats{
		Requests:   t.capacity - bucket.tokens,
		Remaining:  bucket.tokens,
		ResetTime:  bucket.lastRefill.Add(t.window),
		RetryAfter: retryAfter,
	}
}

// SlidingWindowLimiter 滑動時間窗口限制器
type SlidingWindowLimiter struct {
	windows map[string]*SlidingWindow
	mu      sync.RWMutex
	limit   int           // 時間窗口內的請求限制
	window  time.Duration // 時間窗口大小
}

// SlidingWindow 滑動時間窗口
type SlidingWindow struct {
	requests []time.Time
	limit    int
	window   time.Duration
}

// NewSlidingWindowLimiter 建立新的滑動時間窗口限制器
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windows: make(map[string]*SlidingWindow),
		limit:   limit,
		window:  window,
	}
}

// Allow 檢查是否允許請求
func (s *SlidingWindowLimiter) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	window, exists := s.windows[key]
	if !exists {
		window = &SlidingWindow{
			requests: make([]time.Time, 0),
			limit:    s.limit,
			window:   s.window,
		}
		s.windows[key] = window
	}

	now := time.Now()
	cutoff := now.Add(-window.window)

	// 移除過期的請求記錄
	validRequests := make([]time.Time, 0)
	for _, req := range window.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	window.requests = validRequests

	// 檢查是否超過限制
	if len(window.requests) >= window.limit {
		return false
	}

	// 記錄新請求
	window.requests = append(window.requests, now)
	return true
}

// Reset 重置指定鍵的限制
func (s *SlidingWindowLimiter) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.windows, key)
}

// GetStats 獲取速率限制統計
func (s *SlidingWindowLimiter) GetStats(key string) RateLimitStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	window, exists := s.windows[key]
	if !exists {
		return RateLimitStats{
			Requests:   0,
			Remaining:  s.limit,
			ResetTime:  time.Now().Add(s.window),
			RetryAfter: 0,
		}
	}

	now := time.Now()
	cutoff := now.Add(-window.window)

	// 計算有效請求數
	validCount := 0
	oldestRequest := now
	for _, req := range window.requests {
		if req.After(cutoff) {
			validCount++
			if req.Before(oldestRequest) {
				oldestRequest = req
			}
		}
	}

	remaining := s.limit - validCount
	if remaining < 0 {
		remaining = 0
	}

	// 計算重置時間和重試時間
	resetTime := oldestRequest.Add(s.window)
	retryAfter := time.Duration(0)
	if remaining == 0 {
		retryAfter = time.Until(resetTime)
		if retryAfter < 0 {
			retryAfter = 0
		}
	}

	return RateLimitStats{
		Requests:   validCount,
		Remaining:  remaining,
		ResetTime:  resetTime,
		RetryAfter: retryAfter,
	}
}

// RateLimitMiddleware 速率限制中間件
type RateLimitMiddleware struct {
	limiter RateLimiter
	keyFunc func(*gin.Context) string
	onLimit func(*gin.Context, RateLimitStats)
}

// NewRateLimitMiddleware 建立新的速率限制中間件
func NewRateLimitMiddleware(limiter RateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: limiter,
		keyFunc: defaultKeyFunc,
		onLimit: defaultOnLimit,
	}
}

// WithKeyFunc 設置自訂鍵生成函數
func (r *RateLimitMiddleware) WithKeyFunc(keyFunc func(*gin.Context) string) *RateLimitMiddleware {
	r.keyFunc = keyFunc
	return r
}

// WithOnLimit 設置達到限制時的處理函數
func (r *RateLimitMiddleware) WithOnLimit(onLimit func(*gin.Context, RateLimitStats)) *RateLimitMiddleware {
	r.onLimit = onLimit
	return r
}

// Handler 返回中間件處理函數
func (r *RateLimitMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := r.keyFunc(c)

		if !r.limiter.Allow(key) {
			stats := r.limiter.GetStats(key)
			r.onLimit(c, stats)
			return
		}

		// 設置速率限制標頭
		stats := r.limiter.GetStats(key)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", stats.Requests+stats.Remaining))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", stats.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", stats.ResetTime.Unix()))

		c.Next()
	}
}

// defaultKeyFunc 預設鍵生成函數（使用 IP 地址）
func defaultKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

// defaultOnLimit 預設限制處理函數
func defaultOnLimit(c *gin.Context, stats RateLimitStats) {
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", stats.Requests+stats.Remaining))
	c.Header("X-RateLimit-Remaining", "0")
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", stats.ResetTime.Unix()))

	if stats.RetryAfter > 0 {
		c.Header("Retry-After", fmt.Sprintf("%.0f", stats.RetryAfter.Seconds()))
	}

	c.JSON(http.StatusTooManyRequests, gin.H{
		"error":       "請求過於頻繁",
		"code":        "RATE_LIMIT_EXCEEDED",
		"retry_after": stats.RetryAfter.Seconds(),
	})
	c.Abort()
}

// UserKeyFunc 基於用戶 ID 的鍵生成函數
func UserKeyFunc(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%v", userID)
	}
	return fmt.Sprintf("ip:%s", c.ClientIP())
}

// RouteKeyFunc 基於路由的鍵生成函數
func RouteKeyFunc(c *gin.Context) string {
	return fmt.Sprintf("%s:%s:%s", c.ClientIP(), c.Request.Method, c.FullPath())
}

// CreateAPIRateLimiter 建立 API 速率限制器
func CreateAPIRateLimiter() *RateLimitMiddleware {
	// 每分鐘 60 個請求
	limiter := NewSlidingWindowLimiter(60, time.Minute)
	return NewRateLimitMiddleware(limiter)
}

// CreateLoginRateLimiter 建立登入速率限制器
func CreateLoginRateLimiter() *RateLimitMiddleware {
	// 每 15 分鐘 5 次登入嘗試
	limiter := NewSlidingWindowLimiter(5, 15*time.Minute)
	return NewRateLimitMiddleware(limiter).WithKeyFunc(func(c *gin.Context) string {
		return fmt.Sprintf("login:%s", c.ClientIP())
	})
}

// CreateMessageRateLimiter 建立訊息發送速率限制器
func CreateMessageRateLimiter() *RateLimitMiddleware {
	// 每分鐘 30 條訊息
	limiter := NewTokenBucketLimiter(30, 30, time.Minute)
	return NewRateLimitMiddleware(limiter).WithKeyFunc(UserKeyFunc)
}

// CreateSwipeRateLimiter 建立滑動速率限制器
func CreateSwipeRateLimiter() *RateLimitMiddleware {
	// 每小時 100 次滑動
	limiter := NewSlidingWindowLimiter(100, time.Hour)
	return NewRateLimitMiddleware(limiter).WithKeyFunc(UserKeyFunc)
}

// CreatePhotoUploadRateLimiter 建立照片上傳速率限制器
func CreatePhotoUploadRateLimiter() *RateLimitMiddleware {
	// 每天 20 次照片上傳
	limiter := NewSlidingWindowLimiter(20, 24*time.Hour)
	return NewRateLimitMiddleware(limiter).WithKeyFunc(UserKeyFunc)
}

// CreateRegistrationRateLimiter 建立註冊速率限制器
func CreateRegistrationRateLimiter() *RateLimitMiddleware {
	// 每小時每 IP 3 次註冊
	limiter := NewSlidingWindowLimiter(3, time.Hour)
	return NewRateLimitMiddleware(limiter).WithKeyFunc(func(c *gin.Context) string {
		return fmt.Sprintf("register:%s", c.ClientIP())
	})
}

// CleanupExpired 清理過期的限制記錄
func CleanupExpired(limiter RateLimiter) {
	// 這個函數需要根據具體的限制器實現來清理過期記錄
	// 可以定期調用來釋放記憶體
}

// 啟動清理協程
func StartRateLimitCleanup(limiters []RateLimiter, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			for _, limiter := range limiters {
				CleanupExpired(limiter)
			}
		}
	}()
}
