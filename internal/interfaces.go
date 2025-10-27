package internal

import (
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// CacheManager 缓存管理器接口
type CacheManager interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
	Clear()
	Close() error
	GetStats() CacheStats
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits       int64   `json:"hits"`
	Misses     int64   `json:"misses"`
	HitRate    float64 `json:"hit_rate"`
	TotalSize  int64   `json:"total_size"`
	Evictions  int64   `json:"evictions"`
}

// PoolManager 对象池管理器接口
type PoolManager interface {
	Get(lang string) *i18n.Localizer
	Put(lang string, localizer *i18n.Localizer)
	WarmUp(languages []string)
	GetStats() PoolStats
	Close() error
}

// PoolStats 池统计信息
type PoolStats struct {
	Gets     int64 `json:"gets"`
	Puts     int64 `json:"puts"`
	Creates  int64 `json:"creates"`
	PoolSize int64 `json:"pool_size"`
}

// FileWatcher 文件监听器接口
type FileWatcher interface {
	Close() error
}

// Stats 统计信息
type Stats struct {
	Cache      CacheStats `json:"cache"`
	Pool       PoolStats  `json:"pool"`
	Uptime     string     `json:"uptime"`
	NumLocales int        `json:"num_locales"`
}

// Metrics 性能指标
type Metrics struct {
	AvgTranslationTime time.Duration `json:"avg_translation_time"`
	P95TranslationTime time.Duration `json:"p95_translation_time"`
	P99TranslationTime time.Duration `json:"p99_translation_time"`
	TotalTranslations  int64         `json:"total_translations"`
	MemoryUsage        int64         `json:"memory_usage"`
}

// 性能记录函数
var (
	RecordTranslationTime func(time.Duration)
	RecordCacheHit        func()
	RecordCacheMiss       func()
)

// 设置性能记录函数
func SetMetricsFuncs(
	recordTranslationTime func(time.Duration),
	recordCacheHit func(),
	recordCacheMiss func(),
) {
	RecordTranslationTime = recordTranslationTime
	RecordCacheHit = recordCacheHit
	RecordCacheMiss = recordCacheMiss
}

// 默认的空实现
func init() {
	RecordTranslationTime = func(time.Duration) {}
	RecordCacheHit = func() {}
	RecordCacheMiss = func() {}
}