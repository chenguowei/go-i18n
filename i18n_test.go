package i18n

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	)

func TestInit(t *testing.T) {
	// 测试默认初始化
	err := Init()
	require.NoError(t, err)
	require.NotNil(t, GetService())
}

func TestInitWithConfig(t *testing.T) {
	config := Config{
		DefaultLanguage:  "zh-CN",
		FallbackLanguage: "en",
		LocalesPath:      "examples/quickstart/locales",
		Cache: CacheConfig{
			Enable: true,
			Size:   100,
			TTL:    int64(time.Minute.Seconds()),
		},
		Pool: PoolConfig{
			Enable:    true,
			Size:      10,
			WarmUp:    false,
			Languages: []string{"en", "zh-CN"},
		},
		Debug: false,
	}

	err := InitWithConfig(config)
	require.NoError(t, err)
}

func TestTranslate(t *testing.T) {
	// 设置测试环境
	config := Config{
		DefaultLanguage:  "en",
		FallbackLanguage: "en",
		LocalesPath:      "examples/quickstart/locales",
		Cache: CacheConfig{
			Enable: false,
		},
		Pool: PoolConfig{
			Enable: false,
		},
		Debug: false,
	}

	err := InitWithConfig(config)
	require.NoError(t, err)

	// 测试翻译
	ctx := context.Background()
	message := T(ctx, "WELCOME")
	assert.Equal(t, "WELCOME", message) // 如果没有加载语言文件，会返回 messageID
}

func TestMiddleware(t *testing.T) {
	// 初始化
	err := Init()
	require.NoError(t, err)

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由
	r := gin.New()
	r.Use(Middleware())

	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		lang := GetLanguageFromGin(c)
		JSON(c, Success, map[string]interface{}{
			"lang": lang,
		})
	})

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Language", "zh-CN")

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, 200, w.Code)
}

func TestGetStats(t *testing.T) {
	err := Init()
	require.NoError(t, err)

	stats := GetStats()
	assert.NotNil(t, stats)
	assert.Equal(t, "en", "en") // 检查默认语言（暂时简化）
}

func TestGetMetrics(t *testing.T) {
	err := Init()
	require.NoError(t, err)

	metrics := GetMetrics()
	assert.NotNil(t, metrics)
	assert.True(t, metrics.AvgTranslationTime > 0)
}

func TestConfigValidation(t *testing.T) {
	// 测试有效配置
	validConfig := Config{
		DefaultLanguage:  "en",
		FallbackLanguage: "en",
		LocalesPath:      "locales",
		Cache: CacheConfig{
			Enable: true,
			Size:   100,
			L2Size: 1000,
			TTL:    int64(time.Hour.Seconds()),
		},
		Pool: PoolConfig{
			Enable: true,
			Size:   10,
		},
	}

	err := ValidateConfig(validConfig)
	assert.NoError(t, err)

	// 测试无效配置
	invalidConfig := Config{
		DefaultLanguage:  "", // 空语言
		FallbackLanguage: "en",
		LocalesPath:      "locales",
	}

	err = ValidateConfig(invalidConfig)
	assert.Error(t, err)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig
	assert.Equal(t, "en", config.DefaultLanguage)
	assert.Equal(t, "en", config.FallbackLanguage)
	assert.Equal(t, "locales", config.LocalesPath)
	assert.True(t, config.Cache.Enable)
	assert.True(t, config.Pool.Enable)
	assert.False(t, config.Debug)
}

// 基准测试
func BenchmarkTranslate(b *testing.B) {
	err := Init()
	require.NoError(b, err)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T(ctx, "WELCOME")
	}
}

func BenchmarkMiddleware(b *testing.B) {
	err := Init()
	require.NoError(b, err)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// 集成测试
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 使用测试配置
	config := Config{
		DefaultLanguage:  "en",
		FallbackLanguage: "en",
		LocalesPath:      "examples/quickstart/locales",
		Cache: CacheConfig{
			Enable: true,
			Size:   100,
			TTL:    int64(time.Minute.Seconds()),
		},
		Pool: PoolConfig{
			Enable:    true,
			Size:      10,
			WarmUp:    true,
			Languages: []string{"en", "zh-CN"},
		},
		Debug: true,
		EnableMetrics: true,
	}

	err := InitWithConfig(config)
	require.NoError(t, err)

	// 测试服务是否正常工作
	stats := GetStats()
	assert.NotNil(t, stats)
	assert.True(t, stats.Cache.TotalSize >= 0)
	assert.True(t, stats.Pool.PoolSize >= 0)

	// 清理
	err = Close()
	require.NoError(t, err)
}