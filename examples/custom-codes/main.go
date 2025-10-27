package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/chenguowei/go-i18n"
	"github.com/chenguowei/go-i18n/response"
)

func main() {
	// 配置：不加载内置错误码，完全自定义
	config := i18n.Config{
		DefaultLanguage:  "zh-CN",
		FallbackLanguage: "en",
		LocalesPath:      "locales",
		LocaleConfig: i18n.LocaleConfig{
			Mode:      "flat",
			Languages: []string{"en", "zh-CN"},
		},
		ResponseConfig: i18n.ResponseConfig{
			LoadBuiltin: false, // 不加载内置错误码
			AutoInit:     true,  // 自动初始化
		},
		Debug: true,
	}

	// 初始化
	if err := i18n.InitWithConfig(config); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}

	// 注册自定义错误码
	registerCustomCodes()

	// 创建 Gin 路由
	r := gin.Default()

	// 添加 i18n 中间件
	r.Use(i18n.Middleware())

	// 定义路由
	r.GET("/success", successHandler)
	r.GET("/user/:id", getUserHandler)
	r.GET("/order/:id", getOrderHandler)
	r.GET("/payment/:id", getPaymentHandler)
	r.GET("/stats", statsHandler)

	// 启动服务
	fmt.Println("🚀 Custom codes server starting on :8080")
	r.Run(":8080")
}

// 注册自定义错误码
func registerCustomCodes() {
	// 方式1：单个注册
	response.RegisterCustomCode(1000, "USER_NOT_FOUND", 404)
	response.RegisterCustomCode(1001, "INVALID_USER_ID", 400)
	response.RegisterCustomCode(1002, "USER_DISABLED", 403)

	// 方式2：批量注册
	customCodes := []response.CodeDefinition{
		{Code: 2000, Message: "ORDER_NOT_FOUND", HTTPStatus: 404},
		{Code: 2001, Message: "ORDER_EXPIRED", HTTPStatus: 410},
		{Code: 2002, Message: "ORDER_CANCELLED", HTTPStatus: 422},
		{Code: 2003, Message: "INSUFFICIENT_STOCK", HTTPStatus: 409},
	}

	response.BatchRegisterCodes(customCodes)

	// 方式3：从映射表加载
	paymentCodes := map[response.Code]string{
		3000: "PAYMENT_FAILED",
		3001: "PAYMENT_TIMEOUT",
		3002: "INSUFFICIENT_BALANCE",
		3003: "PAYMENT_METHOD_INVALID",
	}

	paymentStatus := map[response.Code]int{
		3000: 402,
		3001: 408,
		3002: 402,
		3003: 400,
	}

	response.LoadCodesFromMap(paymentCodes, paymentStatus)

	fmt.Println("✅ Custom error codes registered successfully")
}

func successHandler(c *gin.Context) {
	response.JSON(c, 0, map[string]interface{}{
		"message": "Operation completed successfully",
	})
}

func getUserHandler(c *gin.Context) {
	userID := c.Param("id")

	if userID == "" {
		response.JSON(c, 1001, map[string]interface{}{
			"error": "User ID is required",
		})
		return
	}

	// 模拟用户查找
	if userID == "999" {
		response.JSON(c, 1000, map[string]interface{}{
			"error": "User not found",
			"user_id": userID,
		})
		return
	}

	if userID == "disabled" {
		response.JSON(c, 1002, map[string]interface{}{
			"error": "User account is disabled",
			"user_id": userID,
		})
		return
	}

	response.JSON(c, 0, map[string]interface{}{
		"user_id": userID,
		"name":    "John Doe",
		"email":   "john@example.com",
	})
}

func getOrderHandler(c *gin.Context) {
	orderID := c.Param("id")

	// 模拟订单查找
	if orderID == "expired" {
		response.JSON(c, 2001, map[string]interface{}{
			"error": "Order has expired",
			"order_id": orderID,
		})
		return
	}

	if orderID == "cancelled" {
		response.JSON(c, 2002, map[string]interface{}{
			"error": "Order was cancelled",
			"order_id": orderID,
		})
		return
	}

	if orderID == "nostock" {
		response.JSON(c, 2003, map[string]interface{}{
			"error": "Insufficient stock for this order",
			"order_id": orderID,
		})
		return
	}

	response.JSON(c, 0, map[string]interface{}{
		"order_id": orderID,
		"status":   "active",
		"total":   99.99,
	})
}

func getPaymentHandler(c *gin.Context) {
	paymentID := c.Param("id")

	// 模拟支付状态
	if paymentID == "failed" {
		response.JSON(c, 3000, map[string]interface{}{
			"error": "Payment processing failed",
			"payment_id": paymentID,
		})
		return
	}

	if paymentID == "timeout" {
		response.JSON(c, 3001, map[string]interface{}{
			"error": "Payment processing timeout",
			"payment_id": paymentID,
		})
		return
	}

	if paymentID == "nobalance" {
		response.JSON(c, 3002, map[string]interface{}{
			"error": "Insufficient account balance",
			"payment_id": paymentID,
		})
		return
	}

	if paymentID == "invalidmethod" {
		response.JSON(c, 3003, map[string]interface{}{
			"error": "Invalid payment method",
			"payment_id": paymentID,
		})
		return
	}

	response.JSON(c, 0, map[string]interface{}{
		"payment_id": paymentID,
		"status":     "completed",
		"amount":     99.99,
	})
}

func statsHandler(c *gin.Context) {
	stats := response.GetCodeStats()
	registeredCodes := response.GetRegisteredCodes()

	response.JSON(c, 0, map[string]interface{}{
		"stats": stats,
		"registered_codes_count": len(registeredCodes),
		"is_initialized": response.IsInitialized(),
	})
}