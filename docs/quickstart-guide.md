# 🚀 GoI18n-Gin 快速开始指南

**仓库**: https://github.com/chenguowei/go-i18n
**模块**: `github.com/chenguowei/go-i18n`

## 📋 概述

GoI18n-Gin 是一个专为 Gin 框架设计的开箱即用多语言库。本指南将帮助您在 5 分钟内完成集成和部署。

## 🎯 适用场景

- ✅ Web 应用多语言支持
- ✅ API 服务国际化
- ✅ 微服务多语言响应
- ✅ 用户界面本地化

---

## 📦 安装

### 方式一：Go Modules（推荐）

```bash
go get github.com/chenguowei/go-i18n@latest
```

### 方式二：源码安装

```bash
git clone https://github.com/chenguowei/go-i18n.git
cd go-i18n
go install ./...
```

---

## 🏃‍♂️ 5分钟快速集成

### 第一步：创建语言文件

在您的项目中创建 `locales` 目录：

```
your-project/
├── main.go
├── locales/
│   ├── en.json
│   ├── zh-CN.json
│   └── zh-TW.json
└── go.mod
```

**locales/en.json**
```json
[
  {
    "id": "WELCOME",
    "translation": "Welcome"
  },
  {
    "id": "USER_NOT_FOUND",
    "translation": "User not found"
  },
  {
    "id": "INVALID_PARAMS",
    "translation": "Invalid parameters"
  },
  {
    "id": "HELLO_USER",
    "translation": "Hello, {{.name}}!"
  }
]
```

**locales/zh-CN.json**
```json
[
  {
    "id": "WELCOME",
    "translation": "欢迎"
  },
  {
    "id": "USER_NOT_FOUND",
    "translation": "用户不存在"
  },
  {
    "id": "INVALID_PARAMS",
    "translation": "参数错误"
  },
  {
    "id": "HELLO_USER",
    "translation": "你好，{{.name}}！"
  }
]
```

### 第二步：编写主程序

**main.go**
```go
package main

import (
    "github.com/gin-gonic/gin"
    i18n "github.com/chenguowei/go-i18n"
)

func main() {
    // 1️⃣ 初始化 i18n（默认配置）
    if err := i18n.Init(); err != nil {
        panic("Failed to initialize i18n: " + err.Error())
    }

    // 2️⃣ 创建 Gin 路由
    r := gin.Default()

    // 3️⃣ 添加 i18n 中间件
    r.Use(i18n.Middleware())

    // 4️⃣ 定义路由
    r.GET("/welcome", welcomeHandler)
    r.GET("/hello", helloHandler)
    r.GET("/error", errorHandler)

    // 5️⃣ 启动服务
    r.Run(":8080")
}

func welcomeHandler(c *gin.Context) {
    // 简单翻译
    message := i18n.TFromGin(c, "WELCOME")

    i18n.SuccessResponse(c, map[string]interface{}{
        "message": message,
        "lang":    i18n.GetLanguageFromGin(c),
    })
}

func helloHandler(c *gin.Context) {
    name := c.DefaultQuery("name", "World")

    // 带参数的翻译
    message := i18n.TFromGin(c, "HELLO_USER", map[string]interface{}{
        "name": name,
    })

    i18n.SuccessResponse(c, map[string]interface{}{
        "message": message,
        "lang":    i18n.GetLanguageFromGin(c),
    })
}

func errorHandler(c *gin.Context) {
    // 使用预定义的错误码
    i18n.Error(c, i18n.UserNotFound)
}
```

### 第三步：运行和测试

```bash
# 运行服务
go run main.go

# 测试不同语言
curl http://localhost:8080/welcome
curl -H "Accept-Language: zh-CN" http://localhost:8080/welcome
curl -H "Accept-Language: zh-TW" http://localhost:8080/welcome

# 测试参数化翻译
curl "http://localhost:8080/hello?name=Alice"
curl -H "Accept-Language: zh-CN" "http://localhost:8080/hello?name=Alice"

# 测试错误响应
curl -H "Accept-Language: zh-CN" http://localhost:8080/error
```

**预期响应：**

```json
// 中文响应
{
  "code": 0,
  "message": "欢迎",
  "data": {
    "message": "欢迎",
    "lang": "zh-CN"
  },
  "meta": {
    "language": "zh-CN",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

---

## 🔧 高级配置

### 自定义配置

```go
package main

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
)

func main() {
    // 自定义配置
    config := i18n.Config{
        DefaultLanguage:  "zh-CN",
        FallbackLanguage: "en",
        LocalesPath:      "./locales",
        Cache: i18n.CacheConfig{
            Enable: true,
            Size:   5000,
            TTL:    int64((2 * time.Hour).Seconds()),
        },
        Pool: i18n.PoolConfig{
            Enable: true,
            Size:   200,
            WarmUp: true,
        },
        Debug: true,
    }

    // 使用自定义配置初始化
    if err := i18n.InitWithConfig(config); err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(i18n.Middleware())

    // ... 路由定义
    r.Run(":8080")
}
```

### 配置文件方式

**config/i18n.yaml**
```yaml
default_language: "zh-CN"
fallback_language: "en"
locales_path: "./locales"

cache:
  enable: true
  size: 5000
  ttl: "2h"
  l2_size: 10000

pool:
  enable: true
  size: 200
  warm_up: true
  languages: ["en", "zh-CN", "zh-TW"]

debug: true
enable_metrics: true
enable_watcher: true
```

**main.go**
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
)
```

func main() {
    // 从配置文件初始化
    if err := i18n.InitFromConfigFile("config/i18n.yaml"); err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(i18n.Middleware())

    // ... 路由定义
    r.Run(":8080")
}
```

---

## 🌐 语言指定方式

### 优先级顺序

1. **`X-Language` Header** - 最高优先级
2. **Cookie** - 用户偏好存储
3. **Query Parameter** - URL 参数
4. **`Accept-Language` Header** - 浏览器标准
5. **默认语言** - 兜底方案

### 使用示例

```bash
# 方式1：X-Language Header（最高优先级）
curl -H "X-Language: zh-CN" http://localhost:8080/welcome

# 方式2：Accept-Language Header（浏览器标准）
curl -H "Accept-Language: zh-CN,en-US;q=0.9" http://localhost:8080/welcome

# 方式3：Cookie
curl -b "lang=zh-CN" http://localhost:8080/welcome

# 方式4：Query Parameter
curl "http://localhost:8080/welcome?lang=zh-CN"

# 方式5：不指定，使用默认语言
curl http://localhost:8080/welcome
```

---

## 📊 监控和调试

### 添加监控端点

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
)
```

func main() {
    config := i18n.Config{
        EnableMetrics: true,
        Debug: true,
    }

    if err := i18n.InitWithConfig(config); err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(i18n.Middleware())

    // 业务路由
    r.GET("/api/hello", helloHandler)

    // 监控端点
    r.GET("/debug/i18n/stats", func(c *gin.Context) {
        stats := i18n.GetStats()
        metrics := i18n.GetMetrics()

        c.JSON(200, gin.H{
            "stats":   stats,
            "metrics": metrics,
        })
    })

    r.Run(":8080")
}
```

### 查看统计信息

```bash
curl http://localhost:8080/debug/i18n/stats
```

**响应示例：**
```json
{
  "stats": {
    "cache_hits": 1250,
    "cache_misses": 45,
    "cache_hit_rate": 0.965,
    "total_translations": 1295,
    "pool_hits": 800,
    "pool_misses": 20
  },
  "metrics": {
    "avg_translation_time": "85μs",
    "p95_translation_time": "120μs",
    "p99_translation_time": "200μs",
    "memory_usage": "12.5MB"
  }
}
```

---

## 🔥 热更新

### 启用热更新

```go
config := i18n.Config{
    EnableWatcher: true,
    LocalesPath:   "./locales",
}

i18n.InitWithConfig(config)
```

### 测试热更新

1. 启动服务
2. 修改 `locales/zh-CN.json` 文件
3. 保存文件，服务自动重载
4. 再次请求，看到新的翻译内容

---

## 🚨 错误处理

### 自定义错误码

**response/codes.go**
```go
package response

type Code int

const (
    Success         Code = 0
    ErrUserNotFound Code = 1001
    ErrInvalidParam Code = 1002
    ErrUnauthorized Code = 1003
)

var CodeMessage = map[Code]string{
    Success:         "SUCCESS",
    ErrUserNotFound: "USER_NOT_FOUND",
    ErrInvalidParam: "INVALID_PARAM",
    ErrUnauthorized: "UNAUTHORIZED",
}
```

### 使用错误响应

```go
func getUserHandler(c *gin.Context) {
    userID := c.Param("id")

    if userID == "" {
        response.JSON(c, response.ErrInvalidParam, nil)
        return
    }

    // 模拟用户查找
    user, err := findUser(userID)
    if err != nil {
        response.JSON(c, response.ErrUserNotFound, nil)
        return
    }

    response.JSON(c, response.Success, user)
}
```

---

## 🎨 最佳实践

### 1. 语言文件组织

```
locales/
├── en/
│   ├── common.json    # 通用翻译
│   ├── errors.json    # 错误信息
│   └── ui.json        # 界面文本
├── zh-CN/
│   ├── common.json
│   ├── errors.json
│   └── ui.json
└── zh-TW/
    ├── common.json
    ├── errors.json
    └── ui.json
```

### 2. 翻译 ID 命名规范

```json
{
  "id": "MODULE_ACTION_ENTITY",
  "translation": "Translation text"
}
```

示例：
- `USER_CREATE_SUCCESS` - 用户创建成功
- `ORDER_NOT_FOUND` - 订单不存在
- `VALIDATION_EMAIL_REQUIRED` - 邮箱必填验证

### 3. 模板参数使用

```json
{
  "id": "USER_WELCOME",
  "translation": "Welcome, {{.name}}! Your account {{.status}}."
}
```

```go
i18n.TFromGin(c, "USER_WELCOME", map[string]interface{}{
    "name": "Alice",
    "status": "is active",
})
```

### 4. 配置文件分层

```
config/
├── i18n.yaml          # 基础配置
├── i18n.dev.yaml      # 开发环境覆盖
├── i18n.prod.yaml     # 生产环境覆盖
└── i18n.test.yaml     # 测试环境覆盖
```

### 5. 错误处理策略

```go
// 翻译失败时的处理策略
func safeTranslate(c *gin.Context, messageID string, fallback string) string {
    translated := i18n.TFromGin(c, messageID)
    if translated == messageID {
        // 翻译失败，使用降级文本
        return fallback
    }
    return translated
}
```

---

## 🔧 常见问题

### Q: 如何添加新语言？

A: 在 `locales` 目录添加新的语言文件，如 `ja.json`，重启服务即可。

### Q: 翻译不生效怎么办？

A: 检查以下几点：
1. 语言文件格式是否正确
2. messageID 是否匹配
3. 语言代码是否标准（如 zh-CN 而不是 zh_cn）

### Q: 如何处理复数形式？

A: 使用 go-i18n 的复数语法：

```json
{
  "id": "ITEM_COUNT",
  "translation": {
    "one": "{{$count}} item",
    "other": "{{$count}} items"
  }
}
```

### Q: 如何提高性能？

A: 启用缓存和对象池：

```go
config := i18n.Config{
    Cache: i18n.CacheConfig{
        Enable: true,
        Size:   10000,
        TTL:    time.Hour,
    },
    Pool: i18n.PoolConfig{
        Enable: true,
        Size:   500,
        WarmUp: true,
    },
}
```

---

## 📚 下一步

- 📖 查看 [完整 API 文档](api.md)
- 🔧 了解 [高级配置](configuration.md)
- 📊 学习 [性能优化](performance.md)
- 🐛 查看 [故障排除](troubleshooting.md)
- 🤝 参与 [项目贡献](contributing.md)

---

## 💬 获取帮助

- 📋 [GitHub Issues](https://github.com/your-org/go-i18n/issues)
- 💬 [Discord 社区](https://discord.gg/go-i18n)
- 📖 [在线文档](https://go-i18n.dev)
- 🐦 [Twitter](https://twitter.com/goi18n)

---

🎉 **恭喜！您已经成功集成了 GoI18n-Gin 库！**

现在您的 Gin 应用已经具备了完整的多语言支持能力。享受开发吧！ 🚀