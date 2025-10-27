package i18n

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

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
	// 基础配置
	if val := os.Getenv("I18N_DEFAULT_LANGUAGE"); val != "" {
		config.DefaultLanguage = val
	}

	if val := os.Getenv("I18N_FALLBACK_LANGUAGE"); val != "" {
		config.FallbackLanguage = val
	}

	if val := os.Getenv("I18N_LOCALES_PATH"); val != "" {
		config.LocalesPath = val
	}

	// 调试和监控
	if val := os.Getenv("I18N_DEBUG"); val != "" {
		config.Debug = parseBool(val, false)
	}

	if val := os.Getenv("I18N_ENABLE_METRICS"); val != "" {
		config.EnableMetrics = parseBool(val, false)
	}

	if val := os.Getenv("I18N_ENABLE_WATCHER"); val != "" {
		config.EnableWatcher = parseBool(val, false)
	}

	// 缓存配置
	if val := os.Getenv("I18N_CACHE_ENABLE"); val != "" {
		config.Cache.Enable = parseBool(val, config.Cache.Enable)
	}

	if val := os.Getenv("I18N_CACHE_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil && size > 0 {
			config.Cache.Size = size
		}
	}

	if val := os.Getenv("I18N_CACHE_TTL"); val != "" {
		if ttl, err := time.ParseDuration(val); err == nil {
			// 将 time.Duration 转换为 int64 秒数
			config.Cache.TTL = int64(ttl.Seconds())
		}
	}

	if val := os.Getenv("I18N_CACHE_L2_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil && size > 0 {
			config.Cache.L2Size = size
		}
	}

	if val := os.Getenv("I18N_CACHE_ENABLE_FILE"); val != "" {
		config.Cache.EnableFile = parseBool(val, config.Cache.EnableFile)
	}

	// 池配置
	if val := os.Getenv("I18N_POOL_ENABLE"); val != "" {
		config.Pool.Enable = parseBool(val, config.Pool.Enable)
	}

	if val := os.Getenv("I18N_POOL_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil && size > 0 {
			config.Pool.Size = size
		}
	}

	if val := os.Getenv("I18N_POOL_WARMUP"); val != "" {
		config.Pool.WarmUp = parseBool(val, config.Pool.WarmUp)
	}

	if val := os.Getenv("I18N_POOL_LANGUAGES"); val != "" {
		config.Pool.Languages = strings.Split(val, ",")
		for i, lang := range config.Pool.Languages {
			config.Pool.Languages[i] = strings.TrimSpace(lang)
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

	if config.LocalesPath == "" {
		return fmt.Errorf("locales_path cannot be empty")
	}

	// 验证缓存配置
	if config.Cache.Enable {
		if config.Cache.Size <= 0 {
			return fmt.Errorf("cache size must be positive")
		}

		if config.Cache.TTL <= 0 {
			return fmt.Errorf("cache TTL must be positive")
		}

		if config.Cache.L2Size <= 0 {
			return fmt.Errorf("cache L2 size must be positive")
		}
	}

	// 验证池配置
	if config.Pool.Enable {
		if config.Pool.Size <= 0 {
			return fmt.Errorf("pool size must be positive")
		}

		if len(config.Pool.Languages) == 0 {
			config.Pool.Languages = []string{config.DefaultLanguage}
		}
	}

	return nil
}

// SaveConfigToFile 保存配置到文件
func SaveConfigToFile(config Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// MergeConfigs 合并配置（后面的配置会覆盖前面的）
func MergeConfigs(configs ...Config) Config {
	result := DefaultConfig

	for _, config := range configs {
		if config.DefaultLanguage != "" {
			result.DefaultLanguage = config.DefaultLanguage
		}
		if config.FallbackLanguage != "" {
			result.FallbackLanguage = config.FallbackLanguage
		}
		if config.LocalesPath != "" {
			result.LocalesPath = config.LocalesPath
		}

		// 合并缓存配置
		if config.Cache.Enable {
			result.Cache.Enable = config.Cache.Enable
		}
		if config.Cache.Size > 0 {
			result.Cache.Size = config.Cache.Size
		}
		if config.Cache.TTL > 0 {
			result.Cache.TTL = config.Cache.TTL
		}
		if config.Cache.L2Size > 0 {
			result.Cache.L2Size = config.Cache.L2Size
		}
		result.Cache.EnableFile = config.Cache.EnableFile

		// 合并池配置
		if config.Pool.Enable {
			result.Pool.Enable = config.Pool.Enable
		}
		if config.Pool.Size > 0 {
			result.Pool.Size = config.Pool.Size
		}
		result.Pool.WarmUp = config.Pool.WarmUp
		if len(config.Pool.Languages) > 0 {
			result.Pool.Languages = config.Pool.Languages
		}

		// 其他配置
		result.Debug = config.Debug
		result.EnableMetrics = config.EnableMetrics
		result.EnableWatcher = config.EnableWatcher
	}

	return result
}

// 辅助函数

func parseBool(s string, defaultValue bool) bool {
	switch strings.ToLower(s) {
	case "true", "1", "yes", "on", "enabled":
		return true
	case "false", "0", "no", "off", "disabled":
		return false
	default:
		return defaultValue
	}
}

// ConfigForEnvironment 根据环境返回配置
func ConfigForEnvironment(env string) Config {
	config := DefaultConfig

	switch strings.ToLower(env) {
	case "development", "dev":
		config.Debug = true
		config.EnableMetrics = false
		config.EnableWatcher = true
		config.Cache.TTL = int64(30 * time.Minute.Seconds())
		config.Pool.WarmUp = false

	case "testing", "test":
		config.Debug = false
		config.EnableMetrics = false
		config.EnableWatcher = false
		config.Cache.Size = 100
		config.Pool.Size = 10
		config.Pool.WarmUp = false

	case "production", "prod":
		config.Debug = false
		config.EnableMetrics = true
		config.EnableWatcher = false
		config.Cache.Size = 10000
		config.Cache.L2Size = 50000
		config.Pool.Size = 500
		config.Pool.WarmUp = true

	case "staging", "stage":
		config.Debug = true
		config.EnableMetrics = true
		config.EnableWatcher = true
		config.Cache.Size = 5000
		config.Pool.Size = 200
		config.Pool.WarmUp = true
	}

	return config
}

// 配置验证规则

type ValidationRule func(Config) error

var validationRules = []ValidationRule{
	validateLanguageCodes,
	validatePaths,
	validatePerformance,
}

func validateLanguageCodes(config Config) error {
	// 验证默认语言代码
	if !IsValidLanguageCode(config.DefaultLanguage) {
		return fmt.Errorf("invalid default language code: %s", config.DefaultLanguage)
	}

	// 验证降级语言代码
	if !IsValidLanguageCode(config.FallbackLanguage) {
		return fmt.Errorf("invalid fallback language code: %s", config.FallbackLanguage)
	}

	return nil
}

func validatePaths(config Config) error {
	// 检查路径是否存在
	if _, err := os.Stat(config.LocalesPath); os.IsNotExist(err) {
		return fmt.Errorf("locales path does not exist: %s", config.LocalesPath)
	}

	return nil
}

func validatePerformance(config Config) error {
	// 检查性能配置是否合理
	if config.Cache.Enable && config.Cache.Size < 100 {
		return fmt.Errorf("cache size too small for production use (minimum 100)")
	}

	if config.Pool.Enable && config.Pool.Size < 10 {
		return fmt.Errorf("pool size too small for production use (minimum 10)")
	}

	return nil
}

// ValidateConfigWithRules 使用自定义规则验证配置
func ValidateConfigWithRules(config Config, rules ...ValidationRule) error {
	// 基础验证
	if err := ValidateConfig(config); err != nil {
		return err
	}

	// 自定义规则验证
	allRules := append(validationRules, rules...)
	for _, rule := range allRules {
		if err := rule(config); err != nil {
			return err
		}
	}

	return nil
}