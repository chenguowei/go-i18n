package main

import (
	"log"

	"github.com/chenguowei/go-i18n"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化 i18n（使用默认配置）
	// if err := i18n.Init(); err != nil {
	// 	log.Fatal("Failed to initialize i18n:", err)
	// }

	if err := i18n.InitFromConfigFile("config.yaml"); err != nil {
		log.Fatal("Failed to initialize i18n from config file:", err)
	}

	// 2. 创建 Gin 路由
	r := gin.Default()

	// 3. 添加 i18n 中间件
	r.Use(i18n.Middleware())

	// 4. 定义路由
	r.GET("/welcome", welcomeHandler)
	r.GET("/hello", helloHandler)
	r.GET("/user/:id", getUserHandler)
	r.GET("/error", errorHandler)

	// 5. 添加调试端点
	r.GET("/debug/stats", debugHandler)

	// 6. 启动服务
	log.Println("🚀 Server starting on :8080")
	log.Println("📊 Debug stats: http://localhost:8080/debug/stats")
	r.Run(":8080")
}

func welcomeHandler(c *gin.Context) {
	// 简单翻译
	// message := i18n.TFromGin(c, "WELCOME")
	// lang := i18n.GetLanguageFromGin(c)

	i18n.JSON(c, i18n.Success, nil)
}

func helloHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "World")

	// 带参数的翻译
	message := i18n.TFromGin(c, "HELLO_USER", map[string]interface{}{
		"name": name,
	})

	i18n.JSON(c, i18n.Success, map[string]interface{}{
		"message": message,
		"lang":    i18n.GetLanguageFromGin(c),
	})
}

func getUserHandler(c *gin.Context) {
	userID := c.Param("id")

	// 模拟用户查找
	if userID == "404" {
		i18n.JSON(c, i18n.InvalidParam, nil)
		return
	}

	// 成功响应
	i18n.JSON(c, i18n.Success, map[string]interface{}{
		"id":    userID,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}

func errorHandler(c *gin.Context) {
	// 使用预定义的错误码
	i18n.JSON(c, i18n.InternalError, nil)
}

func debugHandler(c *gin.Context) {
	stats := i18n.GetStats()
	metrics := i18n.GetMetrics()

	i18n.JSON(c, i18n.Success, map[string]interface{}{
		"stats":   stats,
		"metrics": metrics,
		"version": i18n.GetVersion(),
	})
}
