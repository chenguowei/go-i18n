# 🌐 GoI18n-Gin 库架构设计文档

**仓库**: https://github.com/chenguowei/go-i18n
**模块**: `github.com/chenguowei/go-i18n`

## 📋 项目概述

**GoI18n-Gin** 是一个专为 Gin 框架设计的开箱即用 i18n 库，提供高性能、易集成、生产就绪的多语言解决方案。

### 🎯 设计目标

- **开箱即用**：最小化配置，3行代码集成
- **高性能**：多层缓存 + 对象池，响应时间 < 0.1ms
- **生产就绪**：完善的错误处理和降级机制
- **易于集成**：标准 Gin 中间件，零侵入
- **监控友好**：内置指标和调试端点

---

## 🏗️ 整体架构

### 架构分层

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (Application Layer)                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Gin App   │  │  User Code  │  │   Response System   │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                  中间件层 (Middleware Layer)                  │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              I18n Gin Middleware                     │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐  │    │
│  │  │ Language│  │ Context │  │   Response Helper   │  │    │
│  │  │ Parser  │  │ Manager │  │                     │  │    │
│  │  └─────────┘  └─────────┘  └─────────────────────┘  │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   核心引擎层 (Core Engine Layer)              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Translation Engine                       │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐  │    │
│  │  │ Bundle  │  │Localizer│  │   Template Engine   │  │    │
│  │  │ Manager │  │  Pool   │  │                     │  │    │
│  │  └─────────┘  └─────────┘  └─────────────────────┘  │    │
│  └─────────────────────────────────────────────────────┐    │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    存储层 (Storage Layer)                    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Cache & Storage System                   │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐  │    │
│  │  │ L1 Cache│  │L2 Cache │  │   File Storage      │  │    │
│  │  │ (Memory)│  │ (LRU)   │  │ & Hot Reload        │  │    │
│  │  └─────────┘  └─────────┘  └─────────────────────┘  │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 核心组件

| 组件 | 职责 | 特性 |
|------|------|------|
| **I18n Middleware** | Gin 中间件，语言解析和上下文注入 | 自动语言检测、优先级管理 |
| **Translation Engine** | 核心翻译引擎 | 多层缓存、对象池、模板支持 |
| **Cache Manager** | 缓存管理 | L1(内存) + L2(LRU) + L3(文件) |
| **Pool Manager** | 对象池管理 | Localizer 复用、内存优化 |
| **Response System** | 统一响应系统 | 错误码映射、模板参数 |
| **Config Manager** | 配置管理 | 多格式支持、环境变量 |

---

## 📦 包结构设计

```
go-i18n/
├── i18n.go                      # 核心 API 和初始化（主入口）
├── middleware.go                # Gin 中间件
├── translator.go                # 翻译引擎
├── config.go                    # 配置管理
├── version.go                   # 版本信息
├── options.go                   # 选项模式配置
│
├── response/                    # 响应系统
│   ├── response.go           # 统一响应结构
│   ├── codes.go             # 错误码定义
│   └── helper.go            # 响应辅助函数
│
├── errors/                      # 错误定义
│   ├── errors.go             # 库专用错误
│   └── codes.go              # 错误码常量
│
├── internal/                    # 内部包 (不对外暴露)
│   ├── cache/                  # 缓存系统
│   │   ├── memory.go          # L1 内存缓存
│   │   ├── lru.go             # L2 LRU 缓存
│   │   ├── file.go            # L3 文件缓存
│   │   └── manager.go         # 缓存管理器
│   │
│   ├── pool/                   # 对象池
│   │   ├── localizer.go       # Localizer 对象池
│   │   └── manager.go         # 池管理器
│   │
│   ├── storage/                # 存储层
│   │   ├── bundle.go          # 语言包管理
│   │   ├── loader.go          # 文件加载器
│   │   └── watcher.go         # 热更新监听
│   │
│   ├── parser/                 # 语言解析
│   │   ├── accept_lang.go     # Accept-Language 解析
│   │   └── header.go          # Header 解析
│   │
│   └── monitor/                # 监控系统
│       ├── metrics.go         # 性能指标
│       ├── stats.go           # 统计信息
│       └── debug.go           # 调试端点
│
├── examples/                    # 使用示例
│   ├── quickstart/            # 快速开始
│   │   ├── main.go
│   │   └── locales/
│   ├── advanced/              # 高级用法
│   │   ├── main.go
│   │   ├── config.yaml
│   │   └── locales/
│   └── monitoring/            # 监控集成
│       ├── main.go
│       └── prometheus.go
│
├── configs/                     # 配置示例
│   ├── default.yaml           # 默认配置
│   ├── development.yaml       # 开发环境配置
│   └── production.yaml        # 生产环境配置
│
├── locales/                     # 默认语言文件
│   ├── en.json
│   ├── zh-CN.json
│   └── zh-TW.json
│
├── test/                        # 测试文件
│   ├── unit/                  # 单元测试
│   ├── integration/           # 集成测试
│   └── benchmark/             # 性能测试
│
├── docs/                        # 文档
│   ├── api.md                 # API 文档
│   ├── configuration.md       # 配置文档
│   └── migration.md           # 迁移指南
│
├── go.mod
├── go.sum
├── README.md
├── CHANGELOG.md
└── LICENSE
```

---

## 🔧 核心 API 设计

### 主 API (i18n.go)

```go
package i18n

import (
    "context"
    "time"
    "github.com/gin-gonic/gin"
)

// Config 配置结构
type Config struct {
    // 基础配置
    DefaultLanguage  string        `yaml:"default_language" json:"default_language"`
    FallbackLanguage string        `yaml:"fallback_language" json:"fallback_language"`
    LocalesPath      string        `yaml:"locales_path" json:"locales_path"`

    // 缓存配置
    Cache            CacheConfig   `yaml:"cache" json:"cache"`
    Pool             PoolConfig    `yaml:"pool" json:"pool"`

    // 调试和监控
    Debug            bool          `yaml:"debug" json:"debug"`
    EnableMetrics    bool          `yaml:"enable_metrics" json:"enable_metrics"`
    EnableWatcher    bool          `yaml:"enable_watcher" json:"enable_watcher"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
    Enable     bool          `yaml:"enable" json:"enable"`
    Size       int           `yaml:"size" json:"size"`
    TTL        time.Duration `yaml:"ttl" json:"ttl"`
    L2Size     int           `yaml:"l2_size" json:"l2_size"`
    EnableFile bool          `yaml:"enable_file" json:"enable_file"`
}

// PoolConfig 对象池配置
type PoolConfig struct {
    Enable    bool   `yaml:"enable" json:"enable"`
    Size      int    `yaml:"size" json:"size"`
    WarmUp    bool   `yaml:"warm_up" json:"warm_up"`
    Languages []string `yaml:"languages" json:"languages"`
}

// Init 初始化 i18n 系统（使用默认配置）
func Init() error

// InitWithConfig 使用自定义配置初始化
func InitWithConfig(config Config) error

// InitFromConfigFile 从配置文件初始化
func InitFromConfigFile(configPath string) error

// Middleware 返回 Gin 中间件
func Middleware() gin.HandlerFunc

// MiddlewareWithConfig 返回带配置的 Gin 中间件
func MiddlewareWithConfig(config Config) gin.HandlerFunc

// T 翻译函数（便捷方法）
func T(ctx context.Context, messageID string, templateData ...map[string]interface{}) string

// TFromGin 从 Gin Context 翻译
func TFromGin(c *gin.Context, messageID string, templateData ...map[string]interface{}) string

// GetLanguage 获取当前语言
func GetLanguage(ctx context.Context) string

// GetLanguageFromGin 从 Gin Context 获取语言
func GetLanguageFromGin(c *gin.Context) string

// SetLanguage 设置语言（用于测试或特殊场景）
func SetLanguage(ctx context.Context, language string) context.Context

// GetStats 获取统计信息
func GetStats() Stats

// GetMetrics 获取性能指标
func GetMetrics() Metrics

// Reload 重新加载语言文件
func Reload() error

// Close 关闭 i18n 系统（清理资源）
func Close() error
```

### Gin 中间件 (middleware.go)

```go
package i18n

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/text/language"
)

// LanguageSource 语言来源
type LanguageSource int

const (
    SourceHeader LanguageSource = iota
    SourceAcceptLanguage
    SourceCookie
    SourceQuery
    SourceDefault
)

// LanguageInfo 语言信息
type LanguageInfo struct {
    Language string        `json:"language"`
    Source   LanguageSource `json:"source"`
    Quality  float64       `json:"quality"`
}

// MiddlewareOptions 中间件选项
type MiddlewareOptions struct {
    HeaderKey        string   `yaml:"header_key" json:"header_key"`
    CookieName       string   `yaml:"cookie_name" json:"cookie_name"`
    QueryKey         string   `yaml:"query_key" json:"query_key"`
    SupportedLangs   []string `yaml:"supported_langs" json:"supported_langs"`
    EnableCookie     bool     `yaml:"enable_cookie" json:"enable_cookie"`
    EnableQuery      bool     `yaml:"enable_query" json:"enable_query"`
}

// DefaultMiddlewareOptions 默认中间件选项
var DefaultMiddlewareOptions = MiddlewareOptions{
    HeaderKey:      "X-Language",
    CookieName:     "lang",
    QueryKey:       "lang",
    SupportedLangs: []string{"en", "zh-CN", "zh-TW"},
    EnableCookie:   true,
    EnableQuery:    true,
}

// Middleware 创建 i18n 中间件
func Middleware() gin.HandlerFunc {
    return MiddlewareWithOpts(DefaultMiddlewareOptions)
}

// MiddlewareWithOpts 创建带选项的 i18n 中间件
func MiddlewareWithOpts(opts MiddlewareOptions) gin.HandlerFunc {
    return func(c *gin.Context) {
        lang := detectLanguage(c, opts)

        // 设置语言到上下文
        c.Set("gi18n_language", lang)
        c.Set("gi18n_language_info", LanguageInfo{
            Language: lang,
            Source:   detectLanguageSource(c, opts),
            Quality:  getLanguageQuality(c, lang),
        })

        // 设置响应头
        c.Header("Content-Language", lang)

        c.Next()
    }
}

// detectLanguage 检测语言
func detectLanguage(c *gin.Context, opts MiddlewareOptions) string {
    // 1. Header 优先级最高
    if header := c.GetHeader(opts.HeaderKey); header != "" {
        if isValidLanguage(header, opts.SupportedLangs) {
            return header
        }
    }

    // 2. Cookie
    if opts.EnableCookie {
        if cookie, err := c.Cookie(opts.CookieName); err == nil {
            if isValidLanguage(cookie, opts.SupportedLangs) {
                return cookie
            }
        }
    }

    // 3. Query Parameter
    if opts.EnableQuery {
        if query := c.Query(opts.QueryKey); query != "" {
            if isValidLanguage(query, opts.SupportedLangs) {
                return query
            }
        }
    }

    // 4. Accept-Language Header
    if accept := c.GetHeader("Accept-Language"); accept != "" {
        if lang := parseAcceptLanguage(accept, opts.SupportedLangs); lang != "" {
            return lang
        }
    }

    // 5. 默认语言
    return getDefaultLanguage()
}
```

### 翻译引擎 (translator.go)

```go
package i18n

import (
    "context"
    "crypto/md5"
    "fmt"
    "time"
    "github.com/nicksnyder/go-i18n/v2/i18n"
)

// Translator 翻译器接口
type Translator interface {
    Translate(ctx context.Context, messageID string, templateData ...map[string]interface{}) string
    TranslateWithLanguage(ctx context.Context, lang, messageID string, templateData ...map[string]interface{}) string
    Localizer(ctx context.Context) *i18n.Localizer
    LocalizerWithLanguage(ctx context.Context, lang string) *i18n.Localizer
}

// translator 翻译器实现
type translator struct {
    bundle *i18n.Bundle
    cache  internal.CacheManager
    pool   internal.PoolManager
    config Config
}

// NewTranslator 创建翻译器
func NewTranslator(config Config) (Translator, error) {
    bundle := i18n.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("json", i18n.UnmarshalJSON)

    // 创建缓存和池
    cache := internal.NewCacheManager(config.Cache)
    pool := internal.NewPoolManager(config.Pool, bundle)

    t := &translator{
        bundle: bundle,
        cache:  cache,
        pool:   pool,
        config: config,
    }

    // 加载语言文件
    if err := t.loadLocales(); err != nil {
        return nil, fmt.Errorf("failed to load locales: %w", err)
    }

    return t, nil
}

// Translate 翻译文本
func (t *translator) Translate(ctx context.Context, messageID string, templateData ...map[string]interface{}) string {
    lang := GetLanguage(ctx)
    return t.TranslateWithLanguage(ctx, lang, messageID, templateData...)
}

// TranslateWithLanguage 使用指定语言翻译
func (t *translator) TranslateWithLanguage(ctx context.Context, lang, messageID string, templateData ...map[string]interface{}) string {
    start := time.Now()
    defer func() {
        if t.config.EnableMetrics {
            recordTranslationTime(time.Since(start))
        }
    }()

    // 构建缓存键
    cacheKey := t.buildCacheKey(lang, messageID, templateData)

    // 尝试从缓存获取
    if cached, found := t.cache.Get(cacheKey); found {
        recordCacheHit()
        return cached
    }

    recordCacheMiss()

    // 获取 Localizer
    loc := t.getLocalizer(lang)

    // 执行翻译
    result := t.doTranslate(loc, messageID, templateData...)

    // 存入缓存
    t.cache.Set(cacheKey, result)

    return result
}

// Localizer 获取 Localizer
func (t *translator) Localizer(ctx context.Context) *i18n.Localizer {
    lang := GetLanguage(ctx)
    return t.LocalizerWithLanguage(ctx, lang)
}

// LocalizerWithLanguage 获取指定语言的 Localizer
func (t *translator) LocalizerWithLanguage(ctx context.Context, lang string) *i18n.Localizer {
    return t.getLocalizer(lang)
}

// buildCacheKey 构建缓存键
func (t *translator) buildCacheKey(lang, messageID string, templateData []map[string]interface{}) string {
    if len(templateData) == 0 {
        return fmt.Sprintf("%s:%s", lang, messageID)
    }

    // 对模板数据进行哈希
    templateHash := md5.Sum([]byte(fmt.Sprintf("%v", templateData)))
    return fmt.Sprintf("%s:%s:%x", lang, messageID, templateHash)
}

// getLocalizer 获取 Localizer（带池化）
func (t *translator) getLocalizer(lang string) *i18n.Localizer {
    if t.config.Pool.Enable {
        return t.pool.Get(lang)
    }
    return i18n.NewLocalizer(t.bundle, lang, t.config.FallbackLanguage)
}

// doTranslate 执行实际翻译
func (t *translator) doTranslate(loc *i18n.Localizer, messageID string, templateData ...map[string]interface{}) string {
    config := &i18n.LocalizeConfig{
        MessageID: messageID,
    }

    if len(templateData) > 0 {
        config.TemplateData = templateData[0]
    }

    if translated, err := loc.Localize(config); err == nil {
        return translated
    }

    // 翻译失败，返回 messageID 或使用降级语言
    if t.config.Debug {
        log.Printf("[gi18n] Translation failed for %s: %v", messageID, err)
    }

    // 尝试使用降级语言
    if t.config.FallbackLanguage != "" && t.config.FallbackLanguage != getLanguageFromContext(ctx) {
        fallbackLoc := i18n.NewLocalizer(t.bundle, t.config.FallbackLanguage)
        if translated, err := fallbackLoc.Localize(config); err == nil {
            return translated
        }
    }

    // 最后返回 messageID
    return messageID
}
```

### 配置管理 (config.go)

```go
package i18n

import (
    "fmt"
    "os"
    "time"
    "gopkg.in/yaml.v3"
)

// DefaultConfig 默认配置
var DefaultConfig = Config{
    DefaultLanguage:  "en",
    FallbackLanguage: "en",
    LocalesPath:      "locales",
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

// LoadConfig 加载配置
func LoadConfig() (Config, error) {
    config := DefaultConfig

    // 从环境变量加载
    if err := loadFromEnv(&config); err != nil {
        return config, fmt.Errorf("failed to load from env: %w", err)
    }

    return config, nil
}

// LoadConfigFromFile 从文件加载配置
func LoadConfigFromFile(filename string) (Config, error) {
    config := DefaultConfig

    data, err := os.ReadFile(filename)
    if err != nil {
        if os.IsNotExist(err) {
            // 文件不存在，使用默认配置
            return config, nil
        }
        return config, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return config, fmt.Errorf("failed to parse config file: %w", err)
    }

    // 从环境变量覆盖（优先级更高）
    if err := loadFromEnv(&config); err != nil {
        return config, fmt.Errorf("failed to load from env: %w", err)
    }

    return config, nil
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config *Config) error {
    if val := os.Getenv("GI18N_DEFAULT_LANGUAGE"); val != "" {
        config.DefaultLanguage = val
    }

    if val := os.Getenv("GI18N_FALLBACK_LANGUAGE"); val != "" {
        config.FallbackLanguage = val
    }

    if val := os.Getenv("GI18N_LOCALES_PATH"); val != "" {
        config.LocalesPath = val
    }

    if val := os.Getenv("GI18N_DEBUG"); val != "" {
        config.Debug = val == "true" || val == "1"
    }

    if val := os.Getenv("GI18N_ENABLE_METRICS"); val != "" {
        config.EnableMetrics = val == "true" || val == "1"
    }

    if val := os.Getenv("GI18N_ENABLE_WATCHER"); val != "" {
        config.EnableWatcher = val == "true" || val == "1"
    }

    // 缓存配置
    if val := os.Getenv("GI18N_CACHE_ENABLE"); val != "" {
        config.Cache.Enable = val == "true" || val == "1"
    }

    if val := os.Getenv("GI18N_CACHE_SIZE"); val != "" {
        if size, err := parseInt(val); err == nil {
            config.Cache.Size = size
        }
    }

    if val := os.Getenv("GI18N_CACHE_TTL"); val != "" {
        if ttl, err := time.ParseDuration(val); err == nil {
            config.Cache.TTL = ttl
        }
    }

    // 池配置
    if val := os.Getenv("GI18N_POOL_ENABLE"); val != "" {
        config.Pool.Enable = val == "true" || val == "1"
    }

    if val := os.Getenv("GI18N_POOL_SIZE"); val != "" {
        if size, err := parseInt(val); err == nil {
            config.Pool.Size = size
        }
    }

    return nil
}

// ValidateConfig 验证配置
func ValidateConfig(config Config) error {
    if config.DefaultLanguage == "" {
        return fmt.Errorf("default_language cannot be empty")
    }

    if config.FallbackLanguage == "" {
        return fmt.Errorf("fallback_language cannot be empty")
    }

    if config.Cache.Size <= 0 {
        return fmt.Errorf("cache size must be positive")
    }

    if config.Cache.TTL <= 0 {
        return fmt.Errorf("cache TTL must be positive")
    }

    if config.Pool.Size <= 0 {
        return fmt.Errorf("pool size must be positive")
    }

    return nil
}

// parseInt 解析整数
func parseInt(s string) (int, error) {
    var result int
    _, err := fmt.Sscanf(s, "%d", &result)
    return result, err
}
```

---

## 🚀 使用示例

### 快速开始 (examples/quickstart/main.go)

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
    "github.com/chenguowei/go-i18n/response"
)

func main() {
    // 1. 初始化 i18n（使用默认配置）
    if err := i18n.Init(); err != nil {
        panic(err)
    }

    // 2. 创建 Gin 路由
    r := gin.Default()

    // 3. 添加 i18n 中间件
    r.Use(i18n.Middleware())

    // 4. 定义路由
    r.GET("/hello", func(c *gin.Context) {
        name := c.Query("name")
        if name == "" {
            name = "World"
        }

        // 使用翻译函数
        message := i18n.TFromGin(c, "HELLO_MESSAGE", map[string]interface{}{
            "name": name,
        })

        response.JSON(c, response.Success, map[string]interface{}{
            "message": message,
        })
    })

    r.GET("/error", func(c *gin.Context) {
        response.JSON(c, response.ErrNotFound, nil)
    })

    // 5. 启动服务
    r.Run(":8080")
}
```

### 高级配置 (examples/advanced/main.go)

```go
package main

import (
    "log"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
    "github.com/chenguowei/go-i18n/response"
)

func main() {
    // 自定义配置
    config := i18n.Config{
        DefaultLanguage:  "zh-CN",
        FallbackLanguage: "en",
        LocalesPath:      "./locales",
        Cache: i18n.CacheConfig{
            Enable:     true,
            Size:       5000,
            TTL:        2 * time.Hour,
            L2Size:     10000,
            EnableFile: true,
        },
        Pool: i18n.PoolConfig{
            Enable:    true,
            Size:      200,
            WarmUp:    true,
            Languages: []string{"en", "zh-CN", "zh-TW", "ja"},
        },
        Debug:         true,
        EnableMetrics: true,
        EnableWatcher: true,
    }

    // 使用自定义配置初始化
    if err := i18n.InitWithConfig(config); err != nil {
        log.Fatal("Failed to initialize i18n:", err)
    }

    r := gin.Default()

    // 使用自定义中间件选项
    middlewareOpts := i18n.DefaultMiddlewareOptions
    middlewareOpts.SupportedLangs = []string{"en", "zh-CN", "zh-TW", "ja"}
    middlewareOpts.EnableCookie = true
    middlewareOpts.EnableQuery = true

    r.Use(i18n.MiddlewareWithOpts(middlewareOpts))

    // 添加调试端点
    r.GET("/debug/i18n/stats", func(c *gin.Context) {
        stats := i18n.GetStats()
        metrics := i18n.GetMetrics()

        c.JSON(200, gin.H{
            "stats":   stats,
            "metrics": metrics,
        })
    })

    // 业务路由
    r.GET("/user/:id", getUserHandler)
    r.POST("/user", createUserHandler)

    r.Run(":8080")
}

func getUserHandler(c *gin.Context) {
    userID := c.Param("id")

    // 模拟用户查找
    if userID == "404" {
        response.JSONWithTemplate(c, response.ErrUserNotFound, nil, map[string]interface{}{
            "userID": userID,
        })
        return
    }

    response.JSON(c, response.Success, map[string]interface{}{
        "id":   userID,
        "name": "John Doe",
    })
}

func createUserHandler(c *gin.Context) {
    // 创建用户逻辑...

    response.JSON(c, response.Success, map[string]interface{}{
        "message": i18n.TFromGin(c, "USER_CREATED"),
        "id":      "12345",
    })
}
```

### 监控集成 (examples/monitoring/main.go)

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/chenguowei/go-i18n"
    "github.com/chenguowei/go-i18n/internal/monitor"
)

func main() {
    // 启用监控的配置
    config := i18n.Config{
        DefaultLanguage:  "en",
        FallbackLanguage: "en",
        EnableMetrics:    true,
        Debug:           false,
    }

    if err := i18n.InitWithConfig(config); err != nil {
        panic(err)
    }

    // 注册 Prometheus 指标
    monitor.RegisterPrometheusMetrics()

    r := gin.Default()
    r.Use(i18n.Middleware())

    // 业务路由
    r.GET("/api/hello", helloHandler)

    // 监控端点
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))
    r.GET("/debug/i18n", monitor.DebugHandler())

    r.Run(":8080")
}

func helloHandler(c *gin.Context) {
    name := c.DefaultQuery("name", "World")

    c.JSON(200, gin.H{
        "message": i18n.TFromGin(c, "HELLO", map[string]interface{}{
            "name": name,
        }),
        "lang": i18n.GetLanguageFromGin(c),
    })
}
```

---

## 📊 性能指标

### 基准测试结果

```bash
go test -bench=. -benchmem ./test/benchmark/
```

预期结果：

```
BenchmarkTranslation-8         5000000    120 ns/op    16 B/op    1 allocs/op
BenchmarkTranslationWithCache-8  10000000    85 ns/op     0 B/op    0 allocs/op
BenchmarkConcurrency-8          2000000    200 ns/op    32 B/op    2 allocs/op
BenchmarkLocalizerPool-8       10000000    60 ns/op      0 B/op    0 allocs/op
```

### 性能目标

| 指标 | 目标值 | 测试条件 |
|------|--------|----------|
| 翻译响应时间 | < 0.1ms | 缓存命中 |
| 缓存命中率 | > 85% | 热点数据 |
| 并发处理能力 | 10K+ QPS | 8 核 CPU |
| 内存分配 | 减少 80% | 对象池 vs 无池 |
| 启动时间 | < 100ms | 冷启动 |

---

## 🔄 版本发布策略

### 版本规范 (SemVer 2.0.0)

- **MAJOR**: 不兼容的 API 变更
- **MINOR**: 向后兼容的功能新增
- **PATCH**: 向后兼容的问题修正

### 发布周期

- **主版本**: 每年 1-2 次
- **次版本**: 每季度 1-2 次
- **补丁版本**: 根据需要随时发布

### 兼容性保证

1. **API 稳定性**: 公共 API 向后兼容
2. **配置兼容**: 配置格式向后兼容
3. **数据格式**: 语言文件格式向后兼容

### 升级路径

```bash
# v1.x -> v2.x 升级指南
# 详见 docs/migration.md
```

---

## 📝 集成指南

### 1. 快速集成

```bash
go get github.com/your-org/go-i18n@latest
```

```go
import "github.com/your-org/go-i18n/pkg/gi18n"

// 最简使用
i18n.Init()
r.Use(i18n.Middleware())
```

### 2. 自定义集成

```go
// 使用配置文件
i18n.InitFromConfigFile("config/i18n.yaml")

// 使用自定义配置
config := i18n.Config{...}
i18n.InitWithConfig(config)
```

### 3. 现有项目集成

```go
// 在现有的 Gin 项目中添加
import "github.com/your-org/go-i18n/pkg/response"

// 替换原有的 JSON 响应
// c.JSON(200, gin.H{"code": 0, "message": "success"})
response.JSON(c, response.Success, data)
```

---

## 🎯 总结

GoI18n-Gin 库提供了完整的多语言解决方案，具备以下核心优势：

### 🚀 **开箱即用**
- 3 行代码完成集成
- 默认配置适合大多数场景
- 完整的使用示例和文档

### ⚡ **高性能**
- 多层缓存架构
- Localizer 对象池
- 响应时间 < 0.1ms

### 🛡️ **生产就绪**
- 完善的错误处理
- 降级机制保证可用性
- 详细的监控和调试支持

### 🔧 **易于扩展**
- 模块化设计
- 插件式架构
- 丰富的配置选项

### 📈 **监控友好**
- 内置性能指标
- Prometheus 集成
- 实时统计信息

这个库设计充分考虑了易用性、性能和可维护性，可以作为 Go 生态系统的标准 i18n 解决方案。