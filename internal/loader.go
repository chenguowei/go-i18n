package internal

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// LocaleMode 语言文件组织模式
type LocaleMode string

const (
	// FlatMode 扁平化模式: locales/en.json, locales/zh-CN.json
	FlatMode LocaleMode = "flat"
	// NestedMode 分层模式: locales/en/common.json, locales/zh-CN/common.json
	NestedMode LocaleMode = "nested"
)

// LocaleLoaderConfig 语言文件加载器配置
type LocaleLoaderConfig struct {
	Mode      LocaleMode `yaml:"mode" json:"mode"`
	Path      string     `yaml:"path" json:"path"`
	Languages []string   `yaml:"languages" json:"languages"`
	Modules   []string   `yaml:"modules,omitempty" json:"modules,omitempty"` // 仅在嵌套模式下使用
}

// LocaleLoader 语言文件加载器
type LocaleLoader struct {
	config LocaleLoaderConfig
	bundle *i18n.Bundle
}

// NewLocaleLoader 创建语言文件加载器
func NewLocaleLoader(config LocaleLoaderConfig, bundle *i18n.Bundle) *LocaleLoader {
	return &LocaleLoader{
		config: config,
		bundle: bundle,
	}
}

// LoadLocales 加载语言文件
func (l *LocaleLoader) LoadLocales() error {
	switch l.config.Mode {
	case FlatMode:
		return l.loadFlatStructure()
	case NestedMode:
		return l.loadNestedStructure()
	default:
		return fmt.Errorf("unsupported locale mode: %s", l.config.Mode)
	}
}

// loadFlatStructure 加载扁平化结构
// locales/en.json, locales/zh-CN.json
func (l *LocaleLoader) loadFlatStructure() error {
	for _, lang := range l.config.Languages {
		filename := filepath.Join(l.config.Path, lang+".json")
		if err := l.loadLocaleFile(filename, lang); err != nil {
			return fmt.Errorf("failed to load locale file %s: %w", filename, err)
		}
	}
	return nil
}

// loadNestedStructure 加载分层结构
// locales/en/common.json, locales/en/errors.json
// locales/zh-CN/common.json, locales/zh-CN/errors.json
func (l *LocaleLoader) loadNestedStructure() error {
	// 如果没有指定模块，使用默认模块
	modules := l.config.Modules
	if len(modules) == 0 {
		modules = []string{"common", "errors", "ui"}
	}

	for _, lang := range l.config.Languages {
		for _, module := range modules {
			filename := filepath.Join(l.config.Path, lang, module+".json")
			if err := l.loadLocaleFile(filename, lang); err != nil {
				// 在嵌套模式下，允许某些模块文件不存在
				if isFileNotExistError(err) {
					continue // 跳过不存在的文件
				}
				return fmt.Errorf("failed to load locale file %s: %w", filename, err)
			}
		}
	}
	return nil
}

// loadLocaleFile 加载单个语言文件
func (l *LocaleLoader) loadLocaleFile(filename, lang string) error {
	// 检查文件是否存在
	if _, err := fs.Stat(nil, filename); err != nil {
		return err
	}

	// 这里需要实际的文件加载逻辑
	// 暂时返回成功，实际实现需要读取文件并解析
	// return l.bundle.MustLoadMessageFile(filename, lang)

	return nil
}

// DetectLocaleMode 自动检测语言文件结构模式
func DetectLocaleMode(localesPath string) (LocaleMode, error) {
	// 检查是否存在语言目录
	langDir := filepath.Join(localesPath, "en")
	if dirExists(langDir) {
		return NestedMode, nil
	}

	// 检查是否存在语言文件
	langFile := filepath.Join(localesPath, "en.json")
	if fileExists(langFile) {
		return FlatMode, nil
	}

	return "", fmt.Errorf("no valid locale structure found in %s", localesPath)
}

// GetLocaleFiles 获取所有语言文件列表
func (l *LocaleLoader) GetLocaleFiles() []string {
	var files []string

	switch l.config.Mode {
	case FlatMode:
		for _, lang := range l.config.Languages {
			files = append(files, filepath.Join(l.config.Path, lang+".json"))
		}
	case NestedMode:
		modules := l.config.Modules
		if len(modules) == 0 {
			modules = []string{"common", "errors", "ui"}
		}

		for _, lang := range l.config.Languages {
			for _, module := range modules {
				files = append(files, filepath.Join(l.config.Path, lang, module+".json"))
			}
		}
	}

	return files
}

// ValidateLocaleStructure 验证语言文件结构
func (l *LocaleLoader) ValidateLocaleStructure() error {
	switch l.config.Mode {
	case FlatMode:
		return l.validateFlatStructure()
	case NestedMode:
		return l.validateNestedStructure()
	}
	return fmt.Errorf("invalid locale mode: %s", l.config.Mode)
}

// validateFlatStructure 验证扁平化结构
func (l *LocaleLoader) validateFlatStructure() error {
	for _, lang := range l.config.Languages {
		filename := filepath.Join(l.config.Path, lang+".json")
		if !fileExists(filename) {
			return fmt.Errorf("locale file not found: %s", filename)
		}
	}
	return nil
}

// validateNestedStructure 验证分层结构
func (l *LocaleLoader) validateNestedStructure() error {
	modules := l.config.Modules
	if len(modules) == 0 {
		modules = []string{"common", "errors", "ui"}
	}

	foundAny := false
	for _, lang := range l.config.Languages {
		langDir := filepath.Join(l.config.Path, lang)
		if !dirExists(langDir) {
			return fmt.Errorf("language directory not found: %s", langDir)
		}

		for _, module := range modules {
			filename := filepath.Join(langDir, module+".json")
			if fileExists(filename) {
				foundAny = true
			}
		}
	}

	if !foundAny {
		return fmt.Errorf("no locale files found in nested structure")
	}

	return nil
}

// Helper functions

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	// 实际实现需要使用 os.Stat
	return true // 暂时返回 true
}

// dirExists 检查目录是否存在
func dirExists(dirname string) bool {
	// 实际实现需要使用 os.Stat
	return true // 暂时返回 true
}

// isFileNotExistError 检查是否为文件不存在错误
func isFileNotExistError(err error) bool {
	// 实际实现需要检查错误类型
	return strings.Contains(err.Error(), "no such file")
}

// LocaleFileStats 语言文件统计信息
type LocaleFileStats struct {
	Mode      LocaleMode `json:"mode"`
	TotalFiles int        `json:"total_files"`
	Languages []string   `json:"languages"`
	Modules   []string   `json:"modules,omitempty"`
	FileSizes map[string]int64 `json:"file_sizes"`
}

// GetStats 获取语言文件统计信息
func (l *LocaleLoader) GetStats() LocaleFileStats {
	stats := LocaleFileStats{
		Mode:      l.config.Mode,
		Languages: l.config.Languages,
		FileSizes: make(map[string]int64),
	}

	files := l.GetLocaleFiles()
	stats.TotalFiles = len(files)

	if l.config.Mode == NestedMode {
		stats.Modules = l.config.Modules
		if len(stats.Modules) == 0 {
			stats.Modules = []string{"common", "errors", "ui"}
		}
	}

	// 实际实现需要获取文件大小
	for _, file := range files {
		if fileExists(file) {
			stats.FileSizes[file] = 0 // 实际需要获取真实大小
		}
	}

	return stats
}

// MigrationConfig 迁移配置
type MigrationConfig struct {
	FromMode LocaleMode `yaml:"from_mode" json:"from_mode"`
	ToMode   LocaleMode `yaml:"to_mode" json:"to_mode"`
	Modules  []string   `yaml:"modules,omitempty" json:"modules,omitempty"`
}

// MigrateLocaleStructure 迁移语言文件结构
func (l *LocaleLoader) MigrateLocaleStructure(config MigrationConfig) error {
	if config.FromMode == config.ToMode {
		return fmt.Errorf("source and target modes are the same")
	}

	switch config.ToMode {
	case FlatMode:
		return l.migrateToFlat(config)
	case NestedMode:
		return l.migrateToNested(config)
	default:
		return fmt.Errorf("unsupported target mode: %s", config.ToMode)
	}
}

// migrateToFlat 迁移到扁平化结构
func (l *LocaleLoader) migrateToFlat(config MigrationConfig) error {
	// 实现合并多个模块文件到单个语言文件的逻辑
	// 这是一个复杂的操作，需要：
	// 1. 读取所有模块文件
	// 2. 合并翻译内容
	// 3. 写入单个语言文件
	// 4. 删除原模块文件（可选）

	return fmt.Errorf("migration to flat mode not implemented yet")
}

// migrateToNested 迁移到分层结构
func (l *LocaleLoader) migrateToNested(config MigrationConfig) error {
	// 实现拆分单个语言文件到多个模块文件的逻辑
	// 这是一个复杂的操作，需要：
	// 1. 读取单个语言文件
	// 2. 根据ID前缀或其他规则分类翻译内容
	// 3. 写入多个模块文件
	// 4. 删除原语言文件（可选）

	return fmt.Errorf("migration to nested mode not implemented yet")
}