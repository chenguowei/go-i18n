package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/chenguowei/go-i18n"
)

func main() {
	// 使用默认配置初始化
	if err := i18n.Init(); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(i18n.Middleware())

	// 示例1: 使用默认状态码 (200)
	r.GET("/default", func(c *gin.Context) {
		i18n.JSON(c, i18n.Success, map[string]interface{}{
			"message": "默认 HTTP 200 状态码",
			"usage":   "i18n.JSON()",
		})
	})

	// 示例2: 使用自定义成功状态码 (201)
	r.POST("/created", func(c *gin.Context) {
		i18n.JSONWithStatus(c, i18n.Success, map[string]interface{}{
			"message":  "资源创建成功",
			"usage":    "i18n.JSONWithStatus()",
			"status":   "HTTP 201 Created",
		}, http.StatusCreated)
	})

	// 示例3: 使用自定义错误状态码 (400)
	r.GET("/bad-request", func(c *gin.Context) {
		i18n.ErrorWithStatus(c, i18n.InvalidParam, http.StatusBadRequest)
	})

	// 示例4: 使用自定义错误消息和状态码 (422)
	r.GET("/unprocessable", func(c *gin.Context) {
		i18n.ErrorWithMessageAndStatus(c, i18n.InvalidParam,
			"请求参数无法处理", http.StatusUnprocessableEntity)
	})

	// 示例5: 返回带元数据的自定义状态码响应 (202)
	r.PUT("/accepted", func(c *gin.Context) {
		meta := &i18n.Meta{
			RequestID: c.GetHeader("X-Request-ID"),
			Language:  "zh-CN",
			Version:   "v1.0",
		}
		i18n.JSONWithStatusAndMeta(c, i18n.Success,
			map[string]interface{}{
				"message": "请求已接受，正在处理中",
				"usage":   "i18n.JSONWithStatusAndMeta()",
			}, http.StatusAccepted, meta)
	})

	// 示例6: RESTful API 完整示例
	api := r.Group("/api/v1")
	{
		// 获取用户列表 - 200 OK
		api.GET("/users", listUsers)

		// 创建用户 - 201 Created
		api.POST("/users", createUser)

		// 获取特定用户 - 200 OK 或 404 Not Found
		api.GET("/users/:id", getUser)

		// 更新用户 - 200 OK 或 404 Not Found
		api.PUT("/users/:id", updateUser)

		// 删除用户 - 204 No Content 或 404 Not Found
		api.DELETE("/users/:id", deleteUser)
	}

	// 示例7: 使用模板参数的响应
	r.GET("/template", func(c *gin.Context) {
		templateData := map[string]interface{}{
			"ResourceName": "用户",
			"ResourceID":   "12345",
			"Action":       "创建",
			"Timestamp":    time.Now().Format("2006-01-02 15:04:05"),
		}

		i18n.JSONWithTemplateAndStatus(c, i18n.Success,
			map[string]interface{}{
				"template_used": true,
				"message_type":  "template_based",
			},
			templateData,
			http.StatusCreated)
	})

	// 示例7-2: 多语言模板响应演示
	r.GET("/template/i18n", func(c *gin.Context) {
		templateData := map[string]interface{}{
			"ResourceType": "User",
			"ResourceID":   "12345",
			"Action":       "created",
		}

		i18n.JSONWithTemplateAndStatus(c, i18n.Success,
			map[string]interface{}{
				"i18n_enabled": true,
				"endpoint":     "/template/i18n",
			},
			templateData,
			http.StatusOK)
	})

	// 示例7-3: 错误消息的多语言模板
	r.GET("/template/error", func(c *gin.Context) {
		templateData := map[string]interface{}{
			"FieldName": "email",
			"Reason":    "format is invalid",
		}

		i18n.JSONWithTemplateAndStatus(c, i18n.InvalidParam,
			map[string]interface{}{
				"error_type": "validation_error",
				"endpoint":   "/template/error",
			},
			templateData,
			http.StatusBadRequest)
	})

	// 示例8: 不同业务场景的状态码
	r.GET("/scenarios", func(c *gin.Context) {
		scenarios := []map[string]interface{}{
			{
				"scenario":   "重定向",
				"status":     "HTTP 301 Moved Permanently",
				"usage":      "i18n.JSONWithStatus()",
				"example":    "资源永久迁移",
			},
			{
				"scenario":   "认证失败",
				"status":     "HTTP 401 Unauthorized",
				"usage":      "i18n.ErrorWithStatus()",
				"example":    "Token 无效或过期",
			},
			{
				"scenario":   "权限不足",
				"status":     "HTTP 403 Forbidden",
				"usage":      "i18n.ErrorWithStatus()",
				"example":    "用户无权限访问资源",
			},
			{
				"scenario":   "资源不存在",
				"status":     "HTTP 404 Not Found",
				"usage":      "i18n.ErrorWithStatus()",
				"example":    "请求的资源不存在",
			},
			{
				"scenario":   "请求超时",
				"status":     "HTTP 408 Request Timeout",
				"usage":      "i18n.ErrorWithStatus()",
				"example":    "服务器等待请求时超时",
			},
			{
				"scenario":   "服务器内部错误",
				"status":     "HTTP 500 Internal Server Error",
				"usage":      "i18n.ErrorWithStatus()",
				"example":    "服务器处理请求时发生意外错误",
			},
		}
		i18n.JSONWithStatus(c, i18n.Success, scenarios, http.StatusOK)
	})

	// 启动服务器
	fmt.Println("🚀 HTTP Status 码示例服务器启动在 :8080")
	fmt.Println("📋 可用的端点:")
	fmt.Println("  GET  /default             - 默认状态码 (200)")
	fmt.Println("  POST /created             - 创建成功 (201)")
	fmt.Println("  GET  /bad-request         - 错误请求 (400)")
	fmt.Println("  GET  /unprocessable       - 无法处理的实体 (422)")
	fmt.Println("  PUT  /accepted            - 请求已接受 (202)")
	fmt.Println("  GET  /template            - 模板参数响应 (201)")
	fmt.Println("  GET  /template/i18n       - 多语言模板响应 (200)")
	fmt.Println("  GET  /template/error      - 多语言错误模板 (400)")
	fmt.Println("  GET  /api/v1/users        - 用户列表 (200)")
	fmt.Println("  POST /api/v1/users        - 创建用户 (201)")
	fmt.Println("  GET  /api/v1/users/:id    - 获取用户 (200/404)")
	fmt.Println("  PUT  /api/v1/users/:id    - 更新用户 (200/404)")
	fmt.Println("  DEL  /api/v1/users/:id    - 删除用户 (204/404)")
	fmt.Println("  GET  /scenarios           - 不同场景状态码说明")
	r.Run(":8080")
}

// 模拟用户数据
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: "1", Name: "张三", Email: "zhangsan@example.com"},
	{ID: "2", Name: "李四", Email: "lisi@example.com"},
}

func listUsers(c *gin.Context) {
	i18n.JSONWithStatus(c, i18n.Success, map[string]interface{}{
		"users": users,
		"total": len(users),
	}, http.StatusOK)
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		i18n.ErrorWithMessageAndStatus(c, i18n.InvalidParam,
			"请求参数格式错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 模拟创建用户
	newUser.ID = fmt.Sprintf("%d", len(users)+1)
	users = append(users, newUser)

	i18n.JSONWithStatus(c, i18n.Success, newUser, http.StatusCreated)
}

func getUser(c *gin.Context) {
	userID := c.Param("id")

	for _, user := range users {
		if user.ID == userID {
			i18n.JSONWithStatus(c, i18n.Success, user, http.StatusOK)
			return
		}
	}

	i18n.ErrorWithStatus(c, i18n.NotFound, http.StatusNotFound)
}

func updateUser(c *gin.Context) {
	userID := c.Param("id")

	// 查找用户
	for i, user := range users {
		if user.ID == userID {
			var updateData User
			if err := c.ShouldBindJSON(&updateData); err != nil {
				i18n.ErrorWithMessageAndStatus(c, i18n.InvalidParam,
					"请求参数格式错误: "+err.Error(), http.StatusBadRequest)
				return
			}

			// 更新用户信息
			updateData.ID = userID
			users[i] = updateData

			i18n.JSONWithStatus(c, i18n.Success, updateData, http.StatusOK)
			return
		}
	}

	i18n.ErrorWithStatus(c, i18n.NotFound, http.StatusNotFound)
}

func deleteUser(c *gin.Context) {
	userID := c.Param("id")

	// 查找并删除用户
	for i, user := range users {
		if user.ID == userID {
			users = append(users[:i], users[i+1:]...)
			c.Status(http.StatusNoContent)
			c.Writer.WriteHeaderNow()
			return
		}
	}

	i18n.ErrorWithStatus(c, i18n.NotFound, http.StatusNotFound)
}