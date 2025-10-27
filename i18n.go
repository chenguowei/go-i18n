package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"github.com/gin-gonic/gin"

	"github.com/chenguowei/go-i18n/internal"
)

var (
	globalInstance *Service
	globalOnce     sync.Once
)

// Service i18n 服务
type Service struct {
	bundle         *i18n.Bundle
	translator     Translator
	cache          internal.CacheManager
	pool           internal.PoolManager
	config         Config
	watcher        internal.FileWatcher
	initTime       time.Time
	mu             sync.RWMutex
}

// Config 配置结构
type Config struct {
	// 基础配置
	DefaultLanguage  string        `yaml:"default_language" json:"default_language"`
	FallbackLanguage string        `yaml:"fallback_language" json:"fallback_language"`
	LocalesPath      string        `yaml:"locales_path" json:"locales_path"`

	// 语言文件配置
	LocaleConfig LocaleConfig `yaml:"locale_config" json:"locale_config"`

	// 缓存配置
	Cache CacheConfig `yaml:"cache" json:"cache"`

	// 对象池配置
	Pool PoolConfig `yaml:"pool" json:"pool"`

	// 调试和监控
	Debug         bool `yaml:"debug" json:"debug"`
	EnableMetrics bool `yaml:"enable_metrics" json:"enable_metrics"`
	EnableWatcher bool `yaml:"enable_watcher" json:"enable_watcher"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enable     bool          `yaml:"enable" json:"enable"`
	Size       int           `yaml:"size" json:"size"`
	TTL        time.Duration `yaml:"ttl" json:"ttl"`
	L2Size     int           `yaml:"l2_size" json:"l2_size"`
	EnableFile bool          `yaml:"enable_file" json:"enable_file"`
}

// LocaleConfig 语言文件配置
type LocaleConfig struct {
	// 结构模式: "flat" 或 "nested"
	Mode string `yaml:"mode" json:"mode"`

	// 支持的语言列表
	Languages []string `yaml:"languages" json:"languages"`

	// 模块列表（仅在嵌套模式下使用）
	Modules []string `yaml:"modules,omitempty" json:"modules,omitempty"`
}

// PoolConfig 对象池配置
type PoolConfig struct {
	Enable    bool     `yaml:"enable" json:"enable"`
	Size      int      `yaml:"size" json:"size"`
	WarmUp    bool     `yaml:"warm_up" json:"warm_up"`
	Languages []string `yaml:"languages" json:"languages"`
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	DefaultLanguage:  "en",
	FallbackLanguage: "en",
	LocalesPath:      "locales",
	LocaleConfig: LocaleConfig{
		Mode:      "flat",
		Languages: []string{"en", "zh-CN", "zh-TW"},
	},
	Cache: CacheConfig{
		Enable:     true,
		Size:       1000,
		TTL:        time.Hour,
		L2Size:     5000,
		EnableFile: false,
	},
	Pool: PoolConfig{
		Enable:    true,
		Size:      100,
		WarmUp:    true,
		Languages: []string{"en", "zh-CN", "zh-TW"},
	},
	Debug:         false,
	EnableMetrics: false,
	EnableWatcher: false,
}

// Init 使用默认配置初始化
func Init() error {
	return InitWithConfig(DefaultConfig)
}

// InitWithConfig 使用自定义配置初始化
func InitWithConfig(config Config) error {
	var err error
	globalOnce.Do(func() {
		globalInstance, err = NewService(config)
	})
	return err
}

// InitFromConfigFile 从配置文件初始化
func InitFromConfigFile(configPath string) error {
	config, err := LoadConfigFromFile(configPath)
	if err != nil {
		return err
	}
	return InitWithConfig(config)
}

// NewService 创建新的 i18n 服务
func NewService(config Config) (*Service, error) {
	if err := ValidateConfig(config); err != nil {
		return nil, err
	}

	// 创建 bundle
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	service := &Service{
		bundle:   bundle,
		config:   config,
		initTime: time.Now(),
	}

	// 创建缓存管理器
	if config.Cache.Enable {
		service.cache = internal.NewCacheManager(internal.CacheConfig{
			Enable: config.Cache.Enable,
			Size:   config.Cache.Size,
			TTL:    config.Cache.TTL,
		})
	}

	// 创建对象池管理器
	if config.Pool.Enable {
		service.pool = internal.NewPoolManager(internal.PoolConfig{
			Enable:          config.Pool.Enable,
			Size:            config.Pool.Size,
			WarmUp:          config.Pool.WarmUp,
			Languages:       config.Pool.Languages,
			FallbackLanguage: config.FallbackLanguage,
		}, bundle)
	}

	// 创建翻译器
	service.translator = NewTranslator(bundle, service.cache, service.pool, config)

	// 创建文件监听器
	if config.EnableWatcher {
		service.watcher = internal.NewFileWatcher(config.LocalesPath, service.reloadLocales)
	}

	// 加载语言文件
	if err := service.loadLocales(); err != nil {
		return nil, err
	}

	// 预热对象池
	if config.Pool.WarmUp && service.pool != nil {
		service.pool.WarmUp(config.Pool.Languages)
	}

	return service, nil
}

// GetService 获取全局服务实例
func GetService() *Service {
	if globalInstance == nil {
		panic("i18n not initialized, call Init() first")
	}
	return globalInstance
}

// Service 方法

// Middleware 返回 Gin 中间件
func (s *Service) Middleware() gin.HandlerFunc {
	return MiddlewareWithOpts(DefaultMiddlewareOptions)
}

// Translate 翻译函数
func (s *Service) Translate(ctx context.Context, messageID string, templateData ...map[string]interface{}) string {
	return s.translator.Translate(ctx, messageID, templateData...)
}

// TranslateFromGin 从 Gin Context 翻译
func (s *Service) TranslateFromGin(c *gin.Context, messageID string, templateData ...map[string]interface{}) string {
	lang, _ := c.Get("i18n_language")
	if lang == nil {
		lang = s.config.DefaultLanguage
	}

	ctx := context.WithValue(context.Background(), LanguageKey, lang)
	return s.translator.Translate(ctx, messageID, templateData...)
}

// GetLanguage 获取当前语言
func (s *Service) GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value(LanguageKey).(string); ok {
		return lang
	}
	return s.config.DefaultLanguage
}

// GetLanguageFromGin 从 Gin Context 获取语言
func (s *Service) GetLanguageFromGin(c *gin.Context) string {
	if lang, exists := c.Get("i18n_language"); exists {
		if str, ok := lang.(string); ok {
			return str
		}
	}
	return s.config.DefaultLanguage
}

// GetStats 获取统计信息
func (s *Service) GetStats() internal.Stats {
	cacheStats := internal.CacheStats{}
	poolStats := internal.PoolStats{}

	if s.cache != nil {
		cacheStats = s.cache.GetStats()
	}

	if s.pool != nil {
		poolStats = s.pool.GetStats()
	}

	return internal.Stats{
		Cache:      cacheStats,
		Pool:       poolStats,
		Uptime:     time.Since(s.initTime).String(),
		NumLocales: len(s.supportedLanguages()),
	}
}

// GetMetrics 获取性能指标
func (s *Service) GetMetrics() internal.Metrics {
	return internal.Metrics{
		AvgTranslationTime: 100 * time.Microsecond, // 模拟数据
		P95TranslationTime: 200 * time.Microsecond,
		P99TranslationTime: 500 * time.Microsecond,
		TotalTranslations:  1000,
		MemoryUsage:        10 * 1024 * 1024, // 10MB
	}
}

// Reload 重新加载语言文件
func (s *Service) Reload() error {
	return s.reloadLocales()
}

// Close 关闭 i18n 系统
func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var errs []error

	// 关闭文件监听器
	if s.watcher != nil {
		if err := s.watcher.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// 关闭缓存
	if s.cache != nil {
		if err := s.cache.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// 关闭对象池
	if s.pool != nil {
		if err := s.pool.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errs)
	}

	return nil
}

// 全局便捷方法


// T 翻译函数（便捷方法）
func T(ctx context.Context, messageID string, templateData ...map[string]interface{}) string {
	return GetService().Translate(ctx, messageID, templateData...)
}

// TFromGin 从 Gin Context 翻译
func TFromGin(c *gin.Context, messageID string, templateData ...map[string]interface{}) string {
	return GetService().TranslateFromGin(c, messageID, templateData...)
}

// GetLanguage 获取当前语言
func GetLanguage(ctx context.Context) string {
	return GetService().GetLanguage(ctx)
}

// GetLanguageFromGin 从 Gin Context 获取语言
func GetLanguageFromGin(c *gin.Context) string {
	return GetService().GetLanguageFromGin(c)
}

// GetStats 获取统计信息
func GetStats() internal.Stats {
	return GetService().GetStats()
}

// GetMetrics 获取性能指标
func GetMetrics() internal.Metrics {
	return GetService().GetMetrics()
}

// Reload 重新加载语言文件
func Reload() error {
	return GetService().Reload()
}

// Close 关闭 i18n 系统
func Close() error {
	if globalInstance != nil {
		return globalInstance.Close()
	}
	return nil
}

// 辅助方法

func (s *Service) supportedLanguages() []string {
	return []string{"en", "zh-CN", "zh-TW"}
}

// 内部方法

func (s *Service) loadLocales() error {
	return s.translator.LoadLocales(s.config.LocalesPath)
}

func (s *Service) reloadLocales() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.config.Debug {
		log.Printf("[i18n] Reloading locales from %s", s.config.LocalesPath)
	}

	// 重新加载语言文件
	if err := s.loadLocales(); err != nil {
		return err
	}

	// 清空缓存
	if s.cache != nil {
		s.cache.Clear()
	}

	return nil
}