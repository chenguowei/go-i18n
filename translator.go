package i18n

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/chenguowei/go-i18n/internal"
)

// Translator 翻译器接口
type Translator interface {
	Translate(ctx context.Context, messageID string, templateData ...map[string]interface{}) string
	TranslateWithLanguage(ctx context.Context, lang, messageID string, templateData ...map[string]interface{}) string
	Localizer(ctx context.Context) *i18n.Localizer
	LocalizerWithLanguage(ctx context.Context, lang string) *i18n.Localizer
	LoadLocales(localesPath string) error
}

// translator 翻译器实现
type translator struct {
	bundle *i18n.Bundle
	cache  internal.CacheManager
	pool   internal.PoolManager
	config Config
}

// NewTranslator 创建翻译器
func NewTranslator(bundle *i18n.Bundle, cache internal.CacheManager, pool internal.PoolManager, config Config) Translator {
	return &translator{
		bundle: bundle,
		cache:  cache,
		pool:   pool,
		config: config,
	}
}

// Translate 翻译文本
func (t *translator) Translate(ctx context.Context, messageID string, templateData ...map[string]interface{}) string {
	lang := GetLanguageFromContext(ctx)
	return t.TranslateWithLanguage(ctx, lang, messageID, templateData...)
}

// TranslateWithLanguage 使用指定语言翻译
func (t *translator) TranslateWithLanguage(ctx context.Context, lang, messageID string, templateData ...map[string]interface{}) string {
	start := time.Now()
	defer func() {
		if t.config.EnableMetrics {
			internal.RecordTranslationTime(time.Since(start))
		}
	}()

	// 构建缓存键
	cacheKey := t.buildCacheKey(lang, messageID, templateData)

	// 尝试从缓存获取
	if t.cache != nil {
		if cached, found := t.cache.Get(cacheKey); found {
			internal.RecordCacheHit()
			return cached
		}
	}

	internal.RecordCacheMiss()

	// 获取 Localizer
	loc := t.getLocalizer(lang)

	// 执行翻译
	result := t.doTranslate(loc, messageID, templateData...)

	// 存入缓存
	if t.cache != nil {
		t.cache.Set(cacheKey, result)
	}

	return result
}

// Localizer 获取 Localizer
func (t *translator) Localizer(ctx context.Context) *i18n.Localizer {
	lang := GetLanguageFromContext(ctx)
	return t.LocalizerWithLanguage(ctx, lang)
}

// LocalizerWithLanguage 获取指定语言的 Localizer
func (t *translator) LocalizerWithLanguage(ctx context.Context, lang string) *i18n.Localizer {
	return t.getLocalizer(lang)
}

// LoadLocales 加载语言文件
func (t *translator) LoadLocales(localesPath string) error {
	// 这里简化实现，实际项目中应该遍历目录
	// 支持的语言文件
	supportedFiles := []string{
		"en.json",
		"zh-CN.json",
		"zh-TW.json",
		"ja.json",
		"ko.json",
		"fr.json",
		"de.json",
		"es.json",
		"ru.json",
	}

	loadedCount := 0
	for _, filename := range supportedFiles {
		filePath := filepath.Join(localesPath, filename)
		if _, err := t.bundle.LoadMessageFile(filePath); err == nil {
			loadedCount++
			if t.config.Debug {
				log.Printf("[i18n] Loaded locale file: %s", filename)
			}
		} else if t.config.Debug {
			log.Printf("[i18n] Failed to load %s: %v", filename, err)
		}
	}

	if t.config.Debug {
		log.Printf("[i18n] Loaded %d locale files from %s", loadedCount, localesPath)
	}

	return nil
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
	if t.pool != nil {
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

	translated, err := loc.Localize(config)
	if err == nil {
		return translated
	}

	// 翻译失败处理
	if t.config.Debug {
		log.Printf("[i18n] Translation failed for %s: %v", messageID, err)
	}

	// 尝试使用降级语言
	if t.config.FallbackLanguage != "" {
		fallbackLoc := i18n.NewLocalizer(t.bundle, t.config.FallbackLanguage)
		if translated, err := fallbackLoc.Localize(config); err == nil {
			if t.config.Debug {
				log.Printf("[i18n] Used fallback translation for %s", messageID)
			}
			return translated
		}
	}

	// 最后返回 messageID
	return t.fallbackMessage(messageID)
}

// fallbackMessage 降级消息处理
func (t *translator) fallbackMessage(messageID string) string {
	// 移除下划线，转换为更友好的格式
	result := strings.ReplaceAll(messageID, "_", " ")

	// 首字母大写
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}

	return result
}

// 翻译辅助函数

// Pluralize 复数翻译
func (t *translator) Pluralize(ctx context.Context, messageID string, count interface{}, templateData ...map[string]interface{}) string {
	lang := GetLanguageFromContext(ctx)
	loc := t.getLocalizer(lang)

	config := &i18n.LocalizeConfig{
		MessageID:    messageID,
		PluralCount:  count,
	}

	if len(templateData) > 0 {
		config.TemplateData = templateData[0]
	}

	if translated, err := loc.Localize(config); err == nil {
		return translated
	}

	return t.fallbackMessage(messageID)
}

// TranslateTemplate 翻译模板字符串
func (t *translator) TranslateTemplate(ctx context.Context, template string, templateData ...map[string]interface{}) string {
	lang := GetLanguageFromContext(ctx)
	loc := t.getLocalizer(lang)

	config := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TEMPLATE",
			Other: template,
		},
	}

	if len(templateData) > 0 {
		config.TemplateData = templateData[0]
	}

	if translated, err := loc.Localize(config); err == nil {
		return translated
	}

	// 尝试简单的模板替换
	result := template
	if len(templateData) > 0 {
		for key, value := range templateData[0] {
			placeholder := fmt.Sprintf("{{.%s}}", key)
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
		}
	}

	return result
}

// 支持语言列表
var SupportedLanguages = []struct {
	Code string
	Name string
}{
	{"en", "English"},
	{"zh-CN", "简体中文"},
	{"zh-TW", "繁體中文"},
	{"ja", "日本語"},
	{"ko", "한국어"},
	{"fr", "Français"},
	{"de", "Deutsch"},
	{"es", "Español"},
	{"ru", "Русский"},
	{"pt", "Português"},
	{"it", "Italiano"},
	{"ar", "العربية"},
	{"hi", "हिन्दी"},
	{"th", "ไทย"},
	{"vi", "Tiếng Việt"},
}

// GetSupportedLanguages 获取支持的语言列表
func GetSupportedLanguages() []map[string]interface{} {
	result := make([]map[string]interface{}, len(SupportedLanguages))
	for i, lang := range SupportedLanguages {
		result[i] = map[string]interface{}{
			"code": lang.Code,
			"name": lang.Name,
		}
	}
	return result
}

// IsValidLanguageCode 检查是否为有效的语言代码
func IsValidLanguageCode(code string) bool {
	_, err := language.Parse(code)
	return err == nil
}

// NormalizeLanguageCode 标准化语言代码
func NormalizeLanguageCode(code string) string {
	tag, err := language.Parse(code)
	if err != nil {
		return code
	}

	base, _ := tag.Base()
	script, _ := tag.Script()
	region, _ := tag.Region()

	var result strings.Builder
	result.WriteString(base.String())

	if script.String() != "" {
		result.WriteString("-")
		result.WriteString(script.String())
	}

	if region.String() != "" {
		result.WriteString("-")
		result.WriteString(region.String())
	}

	return result.String()
}