package internal

import (
	"fmt"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// LocalizerPool Localizer 对象池实现
type LocalizerPool struct {
	mu         sync.RWMutex
	pools      map[string]*sync.Pool
	poolMap    map[string]int
	maxPoolSize int
	stats      PoolStats
	bundle     *i18n.Bundle
	fallbackLang string
}

// NewLocalizerPool 创建 Localizer 池
func NewLocalizerPool(maxPoolSize int, bundle *i18n.Bundle, fallbackLang string) PoolManager {
	return &LocalizerPool{
		pools:        make(map[string]*sync.Pool),
		poolMap:      make(map[string]int),
		maxPoolSize:  maxPoolSize,
		stats:        PoolStats{},
		bundle:       bundle,
		fallbackLang: fallbackLang,
	}
}

// Get 获取 Localizer
func (p *LocalizerPool) Get(lang string) *i18n.Localizer {
	p.mu.RLock()
	pool, exists := p.pools[lang]
	p.mu.RUnlock()

	if exists {
		if localizer := pool.Get(); localizer != nil {
			p.stats.recordGet()
			return localizer.(*i18n.Localizer)
		}
	}

	// 池中没有或为空，创建新的
	newLocalizer := i18n.NewLocalizer(p.bundle, lang, p.fallbackLang)
	p.stats.recordCreate()
	return newLocalizer
}

// Put 归还 Localizer
func (p *LocalizerPool) Put(lang string, localizer *i18n.Localizer) {
	if localizer == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果池还没创建，先创建
	if _, exists := p.pools[lang]; !exists {
		p.pools[lang] = &sync.Pool{
			New: func() interface{} {
				return i18n.NewLocalizer(p.bundle, lang, p.fallbackLang)
			},
		}
		p.poolMap[lang] = 0
	}

	// 检查池大小限制
	currentSize := p.poolMap[lang]
	if currentSize >= p.maxPoolSize {
		// 池已满，不归还（让 GC 处理）
		return
	}

	// 归还到池中
	p.pools[lang].Put(localizer)
	p.poolMap[lang]++
	p.stats.recordPut()
}

// WarmUp 预热常用语言的池
func (p *LocalizerPool) WarmUp(languages []string) {
	for _, lang := range languages {
		p.mu.Lock()
		if _, exists := p.pools[lang]; !exists {
			p.pools[lang] = &sync.Pool{
				New: func() interface{} {
					return i18n.NewLocalizer(p.bundle, lang, p.fallbackLang)
				},
			}
			p.poolMap[lang] = 0
		}
		p.mu.Unlock()

		// 预创建一些 Localizer
		for i := 0; i < 5; i++ {
			localizer := i18n.NewLocalizer(p.bundle, lang, p.fallbackLang)
			p.Put(lang, localizer)
		}
	}

	if len(languages) > 0 {
		fmt.Printf("[i18n] Pool warmed up for languages: %v\n", languages)
	}
}

// GetStats 获取池统计
func (p *LocalizerPool) GetStats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// 计算总池大小
	totalPoolSize := int64(0)
	for _, size := range p.poolMap {
		totalPoolSize += int64(size)
	}

	return PoolStats{
		Gets:     p.stats.Gets,
		Puts:     p.stats.Puts,
		Creates:  p.stats.Creates,
		PoolSize: totalPoolSize,
	}
}

// Close 关闭池
func (p *LocalizerPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 清空所有池
	for lang, pool := range p.pools {
		// 清空池中的对象
		for i := 0; i < p.poolMap[lang]; i++ {
			if obj := pool.Get(); obj != nil {
				// 让 GC 处理
			}
		}
		p.poolMap[lang] = 0
	}

	p.pools = make(map[string]*sync.Pool)
	p.poolMap = make(map[string]int)

	return nil
}

// 记录统计方法
func (s *PoolStats) recordGet() {
	s.Gets++
}

func (s *PoolStats) recordPut() {
	s.Puts++
}

func (s *PoolStats) recordCreate() {
	s.Creates++
}

// NoOpPool 空操作池
type NoOpPool struct{}

func (p *NoOpPool) Get(lang string) *i18n.Localizer {
	return nil
}

func (p *NoOpPool) Put(lang string, localizer *i18n.Localizer) {
}

func (p *NoOpPool) WarmUp(languages []string) {
}

func (p *NoOpPool) GetStats() PoolStats {
	return PoolStats{}
}

func (p *NoOpPool) Close() error {
	return nil
}

// PoolManagerFactory 池管理器工厂
type PoolManagerFactory struct{}

// NewPoolManager 创建池管理器
func NewPoolManager(config PoolConfig, bundle *i18n.Bundle) PoolManager {
	if !config.Enable {
		return &NoOpPool{}
	}

	return NewLocalizerPool(config.Size, bundle, config.FallbackLanguage)
}

// PoolConfig 池配置（重新定义以避免循环依赖）
type PoolConfig struct {
	Enable    bool     `yaml:"enable" json:"enable"`
	Size      int      `yaml:"size" json:"size"`
	WarmUp    bool     `yaml:"warm_up" json:"warm_up"`
	Languages []string `yaml:"languages" json:"languages"`
	FallbackLanguage string `yaml:"fallback_language" json:"fallback_language"`
}