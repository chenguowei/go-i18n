package i18n

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// MiddlewareOptions 中间件选项
type MiddlewareOptions struct {
	HeaderKey      string   `yaml:"header_key" json:"header_key"`
	CookieName     string   `yaml:"cookie_name" json:"cookie_name"`
	QueryKey       string   `yaml:"query_key" json:"query_key"`
	SupportedLangs []string `yaml:"supported_langs" json:"supported_langs"`
	EnableCookie   bool     `yaml:"enable_cookie" json:"enable_cookie"`
	EnableQuery    bool     `yaml:"enable_query" json:"enable_query"`
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

// Middleware 返回 Gin 中间件
func Middleware() gin.HandlerFunc {
	return MiddlewareWithOpts(DefaultMiddlewareOptions)
}

// MiddlewareWithOpts 返回带选项的 Gin 中间件
func MiddlewareWithOpts(opts MiddlewareOptions) gin.HandlerFunc {
	// 预编译支持的语言列表
	supportedTags := make([]language.Tag, len(opts.SupportedLangs))
	for i, lang := range opts.SupportedLangs {
		if tag, err := language.Parse(lang); err == nil {
			supportedTags[i] = tag
		}
	}

	// 创建语言匹配器
	matcher := language.NewMatcher(supportedTags)

	return func(c *gin.Context) {
		lang := detectLanguage(c, opts, matcher)

		// 设置语言到上下文
		c.Set("i18n_language", lang)
		c.Set("i18n_language_source", getLanguageSource(c, opts))
		c.Set("i18n_language_quality", getLanguageQuality(c, lang))

		// 设置响应头
		c.Header("Content-Language", lang)

		// 记录调试信息
		if service := GetService(); service.config.Debug {
			start := time.Now()
			c.Set("i18n_start_time", start)
			fmt.Printf("[i18n] Request language: %s (source: %v)\n",
				lang, getLanguageSource(c, opts))
		}

		c.Next()

		// 记录处理时间
		if service := GetService(); service.config.Debug {
			if startTime, exists := c.Get("i18n_start_time"); exists {
				duration := time.Since(startTime.(time.Time))
				fmt.Printf("[i18n] Processing time: %v\n", duration)
			}
		}
	}
}

// detectLanguage 检测语言
func detectLanguage(c *gin.Context, opts MiddlewareOptions, matcher language.Matcher) string {
	// 1. Header 优先级最高
	if header := c.GetHeader(opts.HeaderKey); header != "" {
		if isValidLanguage(header, opts.SupportedLangs) {
			return normalizeLanguage(header)
		}
	}

	// 2. Cookie
	if opts.EnableCookie {
		if cookie, err := c.Cookie(opts.CookieName); err == nil {
			if isValidLanguage(cookie, opts.SupportedLangs) {
				return normalizeLanguage(cookie)
			}
		}
	}

	// 3. Query Parameter
	if opts.EnableQuery {
		if query := c.Query(opts.QueryKey); query != "" {
			if isValidLanguage(query, opts.SupportedLangs) {
				return normalizeLanguage(query)
			}
		}
	}

	// 4. Accept-Language Header
	if accept := c.GetHeader("Accept-Language"); accept != "" {
		if lang := parseAcceptLanguage(accept, matcher); lang != "" {
			return lang
		}
	}

	// 5. 默认语言
	return GetService().config.DefaultLanguage
}

// parseAcceptLanguage 解析 Accept-Language Header
func parseAcceptLanguage(accept string, matcher language.Matcher) string {
	tags, _, err := language.ParseAcceptLanguage(accept)
	if err != nil || len(tags) == 0 {
		return ""
	}

	if tag, _, conf := matcher.Match(tags...); conf > language.No {
		return tag.String()
	}

	return ""
}

// getLanguageSource 获取语言来源
func getLanguageSource(c *gin.Context, opts MiddlewareOptions) string {
	// Header
	if header := c.GetHeader(opts.HeaderKey); header != "" {
		return "header"
	}

	// Cookie
	if opts.EnableCookie {
		if _, err := c.Cookie(opts.CookieName); err == nil {
			return "cookie"
		}
	}

	// Query
	if opts.EnableQuery {
		if query := c.Query(opts.QueryKey); query != "" {
			return "query"
		}
	}

	// Accept-Language
	if accept := c.GetHeader("Accept-Language"); accept != "" {
		return "accept-language"
	}

	return "default"
}

// getLanguageQuality 获取语言质量值
func getLanguageQuality(c *gin.Context, lang string) float64 {
	accept := c.GetHeader("Accept-Language")
	if accept == "" {
		return 1.0
	}

	tags, q, err := language.ParseAcceptLanguage(accept)
	if err != nil {
		return 1.0
	}

	for i, tag := range tags {
		if tag.String() == lang {
			if i < len(q) {
				return float64(q[i])
			}
			return 1.0
		}
	}

	return 0.0
}

// isValidLanguage 验证语言代码是否有效
func isValidLanguage(lang string, supportedLangs []string) bool {
	if lang == "" || len(lang) > 10 {
		return false
	}

	// 基本格式验证
	for _, r := range lang {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '-') {
			return false
		}
	}

	// 检查是否在支持的语言列表中
	normalized := normalizeLanguage(lang)
	for _, supported := range supportedLangs {
		if normalizeLanguage(supported) == normalized {
			return true
		}
	}

	return false
}

// normalizeLanguage 标准化语言代码
func normalizeLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))

	// 处理常见的别名
	switch lang {
	case "zh", "zh_cn", "zh-chs":
		return "zh-cn"
	case "zh_tw", "zh-cht":
		return "zh-tw"
	case "en":
		return "en"
	default:
		return lang
	}
}

// LanguageContextKey 上下文键
type LanguageContextKey string

const (
	LanguageKey     LanguageContextKey = "i18n_language"
	LanguageSource  LanguageContextKey = "i18n_language_source"
	LanguageQuality LanguageContextKey = "i18n_language_quality"
)

// GetLanguageFromContext 从上下文获取语言
func GetLanguageFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(LanguageKey).(string); ok {
		return lang
	}
	return GetService().config.DefaultLanguage
}

// SetLanguageToContext 设置语言到上下文
func SetLanguageToContext(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, LanguageKey, language)
}