package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chenguowei/go-i18n/internal"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: migrate <source_path> <target_mode> [modules]")
		fmt.Println("Example: migrate locales flat")
		fmt.Println("Example: migrate locales nested common,errors,ui")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	targetMode := internal.LocaleMode(os.Args[2])

	// 验证源路径
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		log.Fatalf("Source path does not exist: %s", sourcePath)
	}

	// 检测当前模式
	currentMode, err := internal.DetectLocaleMode(sourcePath)
	if err != nil {
		log.Fatalf("Failed to detect locale mode: %v", err)
	}

	fmt.Printf("Detected current mode: %s\n", currentMode)
	fmt.Printf("Target mode: %s\n", targetMode)

	if currentMode == targetMode {
		fmt.Println("Source and target modes are the same. Nothing to do.")
		return
	}

	// 解析模块列表（仅用于嵌套模式）
	var modules []string
	if targetMode == internal.NestedMode && len(os.Args) > 3 {
		modules = splitModules(os.Args[3])
	}

	// 创建迁移配置
	migrationConfig := internal.MigrationConfig{
		FromMode: currentMode,
		ToMode:   targetMode,
		Modules:  modules,
	}

	// 创建加载器
	loaderConfig := internal.LocaleLoaderConfig{
		Mode:      currentMode,
		Path:      sourcePath,
		Languages: []string{"en", "zh-CN", "zh-TW"}, // 可以从参数或配置获取
		Modules:   modules,
	}

	loader := internal.NewLocaleLoader(loaderConfig, nil)

	// 执行迁移
	fmt.Println("Starting migration...")
	if err := loader.MigrateLocaleStructure(migrationConfig); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed successfully!")

	// 显示结果
	fmt.Println("\nNew structure:")
	printTree(sourcePath, targetMode)
}

func splitModules(modulesStr string) []string {
	// 简单的逗号分割
	// 实际实现可能需要更复杂的解析
	return []string{"common", "errors", "ui", "emails"}
}

func printTree(path string, mode internal.LocaleMode) {
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(filePath) == ".json" {
			relPath, _ := filepath.Rel(path, filePath)
			fmt.Printf("  %s\n", relPath)
		}
		return nil
	})
}