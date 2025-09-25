package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis 客戶端封裝
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	MaxRetries   int           `yaml:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// DefaultRedisConfig 預設 Redis 配置
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}
}

// NewRedisClient 建立新的 Redis 客戶端
func NewRedisClient(config *RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	})

	ctx := context.Background()
	client := &RedisClient{
		client: rdb,
		ctx:    ctx,
	}

	// 測試連接
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("Redis 連接失敗: %w", err)
	}

	log.Printf("Redis 連接成功: %s:%d", config.Host, config.Port)
	return client, nil
}

// Ping 測試 Redis 連接
func (r *RedisClient) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Close 關閉 Redis 連接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Set 設置鍵值對
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get 獲取值
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// GetJSON 獲取 JSON 值
func (r *RedisClient) GetJSON(key string, v interface{}) error {
	val, err := r.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), v)
}

// SetJSON 設置 JSON 值
func (r *RedisClient) SetJSON(key string, value interface{}, expiration time.Duration) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(key, jsonBytes, expiration)
}

// Exists 檢查鍵是否存在
func (r *RedisClient) Exists(keys ...string) (int64, error) {
	return r.client.Exists(r.ctx, keys...).Result()
}

// Delete 刪除鍵
func (r *RedisClient) Delete(keys ...string) (int64, error) {
	return r.client.Del(r.ctx, keys...).Result()
}

// Expire 設置鍵的過期時間
func (r *RedisClient) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// TTL 獲取鍵的剩餘生存時間
func (r *RedisClient) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// Increment 遞增值
func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// IncrementBy 按指定值遞增
func (r *RedisClient) IncrementBy(key string, value int64) (int64, error) {
	return r.client.IncrBy(r.ctx, key, value).Result()
}

// Decrement 遞減值
func (r *RedisClient) Decrement(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// Hash 操作

// HSet 設置哈希字段
func (r *RedisClient) HSet(key string, values ...interface{}) error {
	return r.client.HSet(r.ctx, key, values...).Err()
}

// HGet 獲取哈希字段值
func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// HGetAll 獲取哈希所有字段
func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

// HExists 檢查哈希字段是否存在
func (r *RedisClient) HExists(key, field string) (bool, error) {
	return r.client.HExists(r.ctx, key, field).Result()
}

// HDel 刪除哈希字段
func (r *RedisClient) HDel(key string, fields ...string) (int64, error) {
	return r.client.HDel(r.ctx, key, fields...).Result()
}

// List 操作

// LPush 從列表左側推入元素
func (r *RedisClient) LPush(key string, values ...interface{}) (int64, error) {
	return r.client.LPush(r.ctx, key, values...).Result()
}

// RPush 從列表右側推入元素
func (r *RedisClient) RPush(key string, values ...interface{}) (int64, error) {
	return r.client.RPush(r.ctx, key, values...).Result()
}

// LPop 從列表左側彈出元素
func (r *RedisClient) LPop(key string) (string, error) {
	return r.client.LPop(r.ctx, key).Result()
}

// RPop 從列表右側彈出元素
func (r *RedisClient) RPop(key string) (string, error) {
	return r.client.RPop(r.ctx, key).Result()
}

// LLen 獲取列表長度
func (r *RedisClient) LLen(key string) (int64, error) {
	return r.client.LLen(r.ctx, key).Result()
}

// LRange 獲取列表範圍內的元素
func (r *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

// Set 操作

// SAdd 向集合添加成員
func (r *RedisClient) SAdd(key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(r.ctx, key, members...).Result()
}

// SMembers 獲取集合所有成員
func (r *RedisClient) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// SIsMember 檢查成員是否在集合中
func (r *RedisClient) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// SRem 從集合移除成員
func (r *RedisClient) SRem(key string, members ...interface{}) (int64, error) {
	return r.client.SRem(r.ctx, key, members...).Result()
}

// SCard 獲取集合成員數量
func (r *RedisClient) SCard(key string) (int64, error) {
	return r.client.SCard(r.ctx, key).Result()
}

// Sorted Set 操作

// ZAdd 向有序集合添加成員
func (r *RedisClient) ZAdd(key string, members ...redis.Z) (int64, error) {
	return r.client.ZAdd(r.ctx, key, members...).Result()
}

// ZRange 按索引範圍獲取有序集合成員
func (r *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(r.ctx, key, start, stop).Result()
}

// ZRangeByScore 按分數範圍獲取有序集合成員
func (r *RedisClient) ZRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(r.ctx, key, opt).Result()
}

// ZRem 從有序集合移除成員
func (r *RedisClient) ZRem(key string, members ...interface{}) (int64, error) {
	return r.client.ZRem(r.ctx, key, members...).Result()
}

// ZCard 獲取有序集合成員數量
func (r *RedisClient) ZCard(key string) (int64, error) {
	return r.client.ZCard(r.ctx, key).Result()
}

// Pub/Sub 操作

// Publish 發布訊息到頻道
func (r *RedisClient) Publish(channel string, message interface{}) error {
	return r.client.Publish(r.ctx, channel, message).Err()
}

// Subscribe 訂閱頻道
func (r *RedisClient) Subscribe(channels ...string) *redis.PubSub {
	return r.client.Subscribe(r.ctx, channels...)
}

// PSubscribe 模式訂閱頻道
func (r *RedisClient) PSubscribe(patterns ...string) *redis.PubSub {
	return r.client.PSubscribe(r.ctx, patterns...)
}

// Transaction 操作

// TxPipeline 創建事務管道
func (r *RedisClient) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

// Watch 監視鍵的變化
func (r *RedisClient) Watch(fn func(*redis.Tx) error, keys ...string) error {
	return r.client.Watch(r.ctx, fn, keys...)
}

// Pipeline 操作

// Pipeline 創建管道
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// Lua 腳本操作

// Eval 執行 Lua 腳本
func (r *RedisClient) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(r.ctx, script, keys, args...).Result()
}

// EvalSha 通過 SHA1 執行 Lua 腳本
func (r *RedisClient) EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.EvalSha(r.ctx, sha1, keys, args...).Result()
}

// ScriptLoad 載入 Lua 腳本
func (r *RedisClient) ScriptLoad(script string) (string, error) {
	return r.client.ScriptLoad(r.ctx, script).Result()
}

// 連接池統計

// PoolStats 獲取連接池統計
func (r *RedisClient) PoolStats() *redis.PoolStats {
	return r.client.PoolStats()
}

// 工具方法

// Keys 根據模式獲取鍵列表（生產環境慎用）
func (r *RedisClient) Keys(pattern string) ([]string, error) {
	return r.client.Keys(r.ctx, pattern).Result()
}

// Scan 掃描鍵（推薦用於生產環境）
func (r *RedisClient) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.Scan(r.ctx, cursor, match, count).Result()
}

// FlushDB 清空當前資料庫（危險操作）
func (r *RedisClient) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}

// FlushAll 清空所有資料庫（危險操作）
func (r *RedisClient) FlushAll() error {
	return r.client.FlushAll(r.ctx).Err()
}

// DBSize 獲取資料庫鍵數量
func (r *RedisClient) DBSize() (int64, error) {
	return r.client.DBSize(r.ctx).Result()
}

// Info 獲取 Redis 資訊
func (r *RedisClient) Info(section ...string) (string, error) {
	return r.client.Info(r.ctx, section...).Result()
}

// IsConnected 檢查連接狀態
func (r *RedisClient) IsConnected() bool {
	return r.Ping() == nil
}

// GetClient 獲取原始 Redis 客戶端（用於複雜操作）
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// WithContext 使用指定上下文創建新的客戶端實例
func (r *RedisClient) WithContext(ctx context.Context) *RedisClient {
	return &RedisClient{
		client: r.client,
		ctx:    ctx,
	}
}