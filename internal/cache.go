package internal

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	mu       sync.RWMutex
	items    map[string]*CacheEntry
	maxSize  int
	ttl      time.Duration
	stats    CacheStats
}

// CacheEntry 缓存条目
type CacheEntry struct {
	Value     string
	ExpiresAt time.Time
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(maxSize int, ttl time.Duration) CacheManager {
	return &MemoryCache{
		items:   make(map[string]*CacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
		stats: CacheStats{
			Hits:      0,
			Misses:    0,
			Evictions: 0,
		},
	}
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, exists := c.items[key]; exists {
		if !entry.IsExpired() {
			c.stats.recordHit()
			return entry.Value, true
		}
		// 过期了，删除
		delete(c.items, key)
	}

	c.stats.recordMiss()
	return "", false
}

// Set 设置缓存
func (c *MemoryCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果缓存已满，清理过期条目
	if len(c.items) >= c.maxSize {
		c.evictExpired()
	}

	// 如果还是满的，随机删除一些条目
	if len(c.items) >= c.maxSize {
		c.evictRandom(int(float64(c.maxSize) * 0.2)) // 删除 20%
	}

	c.items[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	c.stats.recordSet(len(c.items))
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear 清空缓存
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheEntry)
	c.stats.recordSet(0)
}

// Close 关闭缓存
func (c *MemoryCache) Close() error {
	c.Clear()
	return nil
}

// GetStats 获取统计信息
func (c *MemoryCache) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		Hits:      c.stats.Hits,
		Misses:    c.stats.Misses,
		HitRate:   c.stats.CalculateHitRate(),
		TotalSize: int64(len(c.items)),
		Evictions: c.stats.Evictions,
	}
}

// IsExpired 检查是否过期
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// evictExpired 清理过期条目
func (c *MemoryCache) evictExpired() {
	now := time.Now()
	for key, entry := range c.items {
		if now.After(entry.ExpiresAt) {
			delete(c.items, key)
			c.stats.recordEviction()
		}
	}
}

// evictRandom 随机删除条目
func (c *MemoryCache) evictRandom(count int) {
	for key := range c.items {
		if count <= 0 {
			break
		}
		delete(c.items, key)
		c.stats.recordEviction()
		count--
	}
}

// 记录统计方法
func (s *CacheStats) recordHit() {
	s.Hits++
}

func (s *CacheStats) recordMiss() {
	s.Misses++
}

func (s *CacheStats) recordSet(size int) {
	s.TotalSize = int64(size)
}

func (s *CacheStats) recordEviction() {
	s.Evictions++
}

// CalculateHitRate 计算命中率
func (s *CacheStats) CalculateHitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total)
}

// CacheManagerFactory 缓存管理器工厂
type CacheManagerFactory struct{}

// NewCacheManager 创建缓存管理器
func NewCacheManager(config CacheConfig) CacheManager {
	if !config.Enable {
		return &NoOpCache{}
	}

	return NewMemoryCache(config.Size, config.TTL)
}

// NoOpCache 空操作缓存
type NoOpCache struct{}

func (c *NoOpCache) Get(key string) (string, bool) {
	return "", false
}

func (c *NoOpCache) Set(key, value string) {
}

func (c *NoOpCache) Delete(key string) {
}

func (c *NoOpCache) Clear() {
}

func (c *NoOpCache) Close() error {
	return nil
}

func (c *NoOpCache) GetStats() CacheStats {
	return CacheStats{}
}

// BuildCacheKey 构建缓存键的辅助函数
func BuildCacheKey(lang, messageID string, templateData []map[string]interface{}) string {
	if len(templateData) == 0 {
		return fmt.Sprintf("%s:%s", lang, messageID)
	}

	// 对模板数据进行哈希处理
	templateHash := md5.Sum([]byte(fmt.Sprintf("%v", templateData)))
	return fmt.Sprintf("%s:%s:%x", lang, messageID, templateHash)
}

// CacheConfig 缓存配置（重新定义以避免循环依赖）
type CacheConfig struct {
	Enable     bool          `yaml:"enable" json:"enable"`
	Size       int           `yaml:"size" json:"size"`
	TTL        time.Duration `yaml:"ttl" json:"ttl"`
	L2Size     int           `yaml:"l2_size" json:"l2_size"`
	EnableFile bool          `yaml:"enable_file" json:"enable_file"`
}