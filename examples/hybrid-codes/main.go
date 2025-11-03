package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/chenguowei/go-i18n"
)

func main() {
	// 配置：加载内置错误码 + 自定义错误码
	config := i18n.Config{
		DefaultLanguage:  "zh-CN",
		FallbackLanguage: "en",
		LocalesPath:      "locales",
		LocaleConfig: i18n.LocaleConfig{
			Mode:      "flat",
			Languages: []string{"en", "zh-CN"},
		},
		ResponseConfig: i18n.ResponseConfig{
			LoadBuiltin: true, // 加载内置错误码
			AutoInit:     true, // 自动初始化
		},
		Debug: true,
	}

	// 初始化
	if err := i18n.InitWithConfig(config); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}

	// 添加自定义错误码（在内置基础上扩展）
	addCustomCodes()

	// 创建 Gin 路由
	r := gin.Default()

	// 添加 i18n 中间件
	r.Use(i18n.Middleware())

	// 定义路由
	r.GET("/builtin", builtinHandler)
	r.GET("/custom", customHandler)
	r.GET("/overview", overviewHandler)

	// 启动服务
	fmt.Println("🚀 Hybrid codes server starting on :8080")
	r.Run(":8080")
}

// 添加自定义错误码
func addCustomCodes() {
	// 业务相关的自定义错误码
	businessCodes := []i18n.CodeDefinition{
		{Code: 5000, Message: "PRODUCT_OUT_OF_STOCK", HTTPStatus: 422},
		{Code: 5001, Message: "PROMOTION_EXPIRED", HTTPStatus: 410},
		{Code: 5002, Message: "COUPON_ALREADY_USED", HTTPStatus: 409},
		{Code: 5003, Message: "MEMBERSHIP_REQUIRED", HTTPStatus: 402},
	}

	i18n.BatchRegisterCodes(businessCodes)

	// 设置自定义消息覆盖内置错误码
	i18n.SetCustomMessage(i18n.InvalidParam, "CUSTOM_INVALID_PARAM")

	fmt.Println("✅ Hybrid error codes setup completed")
}

func builtinHandler(c *gin.Context) {
	// 使用内置错误码
	i18n.JSON(c, i18n.Success, map[string]interface{}{
		"message": "Using built-in error codes",
		"codes": []string{
			"InvalidParam",
			"UserNotFound",
			"Unauthorized",
			"InternalError",
		},
	})

	// 测试内置错误码
	if c.Query("error") == "param" {
		i18n.JSON(c, i18n.InvalidParam, map[string]interface{}{
			"error": "Invalid parameter provided",
		})
	}

	if c.Query("error") == "user" {
		i18n.JSON(c, i18n.InvalidParam, map[string]interface{}{
			"error": "User not found",
		})
	}
}

func customHandler(c *gin.Context) {
	// 使用自定义错误码
	i18n.JSON(c, 0, map[string]interface{}{
		"message": "Using custom error codes",
		"codes": []string{
			"PRODUCT_OUT_OF_STOCK",
			"PROMOTION_EXPIRED",
			"COUPON_ALREADY_USED",
			"MEMBERSHIP_REQUIRED",
		},
	})

	// 测试自定义错误码
	if c.Query("error") == "stock" {
		i18n.JSON(c, 5000, map[string]interface{}{
			"error": "Product is out of stock",
		})
	}

	if c.Query("error") == "promotion" {
		i18n.JSON(c, 5001, map[string]interface{}{
			"error": "Promotion has expired",
		})
	}
}

func overviewHandler(c *gin.Context) {
	stats := i18n.GetCodeStats()
	registeredCodes := i18n.GetRegisteredCodes()

	// 分类统计
	customCodes := make(map[i18n.Code]string)
	builtinCodes := make(map[i18n.Code]string)

	for code, message := range registeredCodes {
		// 简单判断是否为内置错误码（基于错误码范围）
		if code >= 1000 && code < 3000 {
			builtinCodes[code] = message
		} else {
			customCodes[code] = message
		}
	}

	i18n.JSON(c, 0, map[string]interface{}{
		"statistics": stats,
		"total_registered": len(registeredCodes),
		"builtin_codes": builtinCodes,
		"custom_codes": customCodes,
		"is_initialized": i18n.IsInitialized(),
	})
}