package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/chenguowei/go-i18n"
)

func main() {
	// 使用分层结构配置
	config := i18n.Config{
		DefaultLanguage:  "zh-CN",
		FallbackLanguage: "en",
		LocalesPath:      "locales",
		LocaleConfig: i18n.LocaleConfig{
			Mode:      "nested",
			Languages: []string{"en", "zh-CN", "zh-TW"},
			Modules:   []string{"common", "errors", "ui", "emails"},
		},
		Cache: i18n.CacheConfig{
			Enable: true,
			Size:   1000,
			TTL:    2 * time.Hour,
		},
		Pool: i18n.PoolConfig{
			Enable:    true,
			Size:      100,
			WarmUp:    true,
			Languages: []string{"en", "zh-CN", "zh-TW"},
		},
		Debug: true,
	}

	// 初始化
	if err := i18n.InitWithConfig(config); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 添加 i18n 中间件
	r.Use(i18n.Middleware())

	// 定义路由
	r.GET("/welcome", welcomeHandler)
	r.GET("/hello", helloHandler)
	r.GET("/error", errorHandler)
	r.GET("/user-profile", userProfileHandler)

	// 启动服务
	fmt.Println("🚀 Nested structure server starting on :8080")
	r.Run(":8080")
}

func welcomeHandler(c *gin.Context) {
	// 使用通用模块的翻译
	message := i18n.TFromGin(c, "WELCOME")

	i18n.JSON(c, i18n.Success, map[string]interface{}{
		"message": message,
		"lang":    i18n.GetLanguageFromGin(c),
	})
}

func helloHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "World")

	// 使用通用模块的翻译
	message := i18n.TFromGin(c, "HELLO_USER", map[string]interface{}{
		"name": name,
	})

	i18n.JSON(c, i18n.Success, map[string]interface{}{
		"message": message,
		"lang":    i18n.GetLanguageFromGin(c),
	})
}

func errorHandler(c *gin.Context) {
	// 使用错误模块的翻译
	i18n.JSON(c, i18n.ErrUserNotFound, nil)
}

func userProfileHandler(c *gin.Context) {
	// 使用UI模块的翻译
	data := map[string]interface{}{
		"title":       i18n.TFromGin(c, "PROFILE_TITLE"),
		"description": i18n.TFromGin(c, "PROFILE_DESCRIPTION"),
		"edit":        i18n.TFromGin(c, "EDIT_PROFILE"),
		"save":        i18n.TFromGin(c, "SAVE_CHANGES"),
	}

	i18n.JSON(c, i18n.Success, data)
}