# GoI18n-Gin 🌐

[![GoDoc](https://pkg.go.dev/badge/github.com/chenguowei/go-i18n.svg)](https://pkg.go.dev/github.com/chenguowei/go-i18n)
[![Build Status](https://github.com/chenguowei/go-i18n/workflows/CI/badge.svg)](https://github.com/chenguowei/go-i18n/actions)
[![Coverage](https://codecov.io/gh/chenguowei/go-i18n/branch/main/graph/badge.svg)](https://codecov.io/gh/chenguowei/go-i18n)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenguowei/go-i18n)](https://goreportcard.com/report/github.com/chenguowei/go-i18n)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

🌍 一个专为 Gin 框架设计的开箱即用多语言库，提供高性能、易集成、生产就绪的国际化解决方案。

## ✨ 特性

- 🚀 **开箱即用** - 3行代码完成集成，零配置启动
- ⚡ **高性能** - 多层缓存 + 对象池，响应时间 < 0.1ms
- 🛡️ **生产就绪** - 完善的错误处理和降级机制
- 🔥 **热更新** - 无需重启，动态加载语言文件 (可选)
- 📊 **监控友好** - 内置性能指标和调试支持
- 🌐 **多语言检测** - Header、Cookie、Query、Accept-Language 多种方式
- 🎯 **零侵入** - 标准 Gin 中间件，完全兼容现有代码
- 📝 **统一响应** - 内置多语言响应码系统

## 📦 安装

```bash
go get github.com/chenguowei/go-i18n@latest
```

## 🚀 快速开始

### 基础使用

```go
package main

import (
    "github.com/gin-gonic/gin"
    i18n "github.com/chenguowei/go-i18n"
)

func main() {
    // 1. 初始化 i18n (使用默认配置)
    if err := i18n.Init(); err != nil {
        panic(err)
    }

    // 2. 创建 Gin 路由并添加中间件
    r := gin.Default()
    r.Use(i18n.Middleware())

    // 3. 使用翻译和统一响应
    r.GET("/hello", func(c *gin.Context) {
        name := c.DefaultQuery("name", "World")

        // 翻译消息
        message := i18n.TFromGin(c, "WELCOME_MESSAGE", map[string]interface{}{
            "name": name,
        })

        // 统一响应格式
        i18n.SuccessResponse(c, map[string]interface{}{
            "message": message,
            "lang":    i18n.GetLanguageFromGin(c),
        })
    })

    r.Run(":8080")
}
```

### 高级配置

```go
package main

import (
    "time"
    "github.com/gin-gonic/gin"
    i18n "github.com/chenguowei/go-i18n"
)

func main() {
    // 自定义配置
    config := i18n.Config{
        DefaultLanguage:  "zh-CN",
        FallbackLanguage: "en",
        LocalesPath:      "./locales",

        // 语言文件配置
        LocaleConfig: i18n.LocaleConfig{
            Mode:      "flat", // 支持 "flat" 或 "nested"
            Languages: []string{"en", "zh-CN", "zh-TW", "ja"},
        },

        // 缓存配置
        Cache: i18n.CacheConfig{
            Enable:     true,
            Size:       5000,
            TTL:        int64((2 * time.Hour).Seconds()),
            L2Size:     10000,
            EnableFile: false,
        },

        // 对象池配置
        Pool: i18n.PoolConfig{
            Enable:    true,
            Size:      200,
            WarmUp:    true,
            Languages: []string{"en", "zh-CN", "zh-TW", "ja"},
        },

        // 调试和监控
        Debug:         true,
        EnableMetrics: true,
        EnableWatcher: true, // 热更新
    }

    // 使用自定义配置初始化
    if err := i18n.InitWithConfig(config); err != nil {
        panic(err)
    }

    r := gin.Default()

    // 自定义中间件选项
    middlewareOpts := i18n.DefaultMiddlewareOptions
    middlewareOpts.SupportedLangs = []string{"en", "zh-CN", "zh-TW", "ja"}
    middlewareOpts.EnableCookie = true
    middlewareOpts.EnableQuery = true

    r.Use(i18n.MiddlewareWithOpts(middlewareOpts))

    // 业务路由
    r.GET("/api/user/:id", getUserHandler)
    r.POST("/api/user", createUserHandler)

    r.Run(":8080)
}

func getUserHandler(c *gin.Context) {
    userID := c.Param("id")

    if userID == "404" {
        // 使用模板参数的错误响应
        i18n.JSONWithTemplate(c, i18n.ErrUserNotFound, nil, map[string]interface{}{
            "userID": userID,
        })
        return
    }

    i18n.SuccessResponse(c, map[string]interface{}{
        "id":   userID,
        "name": "John Doe",
        "lang": i18n.GetLanguageFromGin(c),
    })
}

func createUserHandler(c *gin.Context) {
    // 创建用户逻辑...

    i18n.SuccessResponse(c, map[string]interface{}{
        "message": i18n.TFromGin(c, "USER_CREATED"),
        "id":      "12345",
    })
}
```

## 📝 语言文件

### 创建语言文件

在 `locales/` 目录下创建 JSON 格式的语言文件：

**locales/en.json**
```json
{
  "WELCOME_MESSAGE": "Hello, {{.name}}!",
  "USER_CREATED": "User created successfully",
  "USER_NOT_FOUND": "User with ID {{.userID}} not found",
  "INVALID_PARAM": "Invalid parameters provided"
}
```

**locales/zh-CN.json**
```json
{
  "WELCOME_MESSAGE": "你好，{{.name}}！",
  "USER_CREATED": "用户创建成功",
  "USER_NOT_FOUND": "ID为{{.userID}}的用户未找到",
  "INVALID_PARAM": "提供的参数无效"
}
```

**locales/zh-TW.json**
```json
{
  "WELCOME_MESSAGE": "你好，{{.name}}！",
  "USER_CREATED": "用戶創建成功",
  "USER_NOT_FOUND": "ID為{{.userID}}的用戶未找到",
  "INVALID_PARAM": "提供的參數無效"
}
```

## 🌍 语言检测

库支持多种语言检测方式，按优先级顺序：

1. **Header**: `X-Language: zh-CN`
2. **Cookie**: `lang=zh-CN`
3. **Query Parameter**: `?lang=zh-CN`
4. **Accept-Language**: `Accept-Language: zh-CN,zh;q=0.9,en;q=0.8`
5. **Default**: 配置的默认语言

### 使用示例

```bash
# 通过 Header 指定语言
curl -H "X-Language: zh-CN" http://localhost:8080/api/hello

# 通过 Cookie 指定语言
curl -b "lang=zh-CN" http://localhost:8080/api/hello

# 通过 Query 参数指定语言
curl "http://localhost:8080/api/hello?lang=zh-CN"

# 使用 Accept-Language
curl -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8" http://localhost:8080/api/hello
```

## 📊 响应系统

库提供了统一的响应系统，支持多语言错误消息：

### 标准响应格式

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {...},
  "meta": {
    "timestamp": "2025-10-30T10:30:00Z",
    "language": "zh-CN",
    "request_id": "req-123",
    "trace_id": "trace-456"
  }
}
```

### 响应函数

```go
// 成功响应
i18n.SuccessResponse(c, data)

// 错误响应
i18n.Error(c, i18n.ErrInvalidParam)
i18n.ErrorWithTemplate(c, i18n.ErrUserNotFound, map[string]interface{}{
    "userID": userID,
})

// 自定义响应
i18n.JSON(c, customCode, data)
i18n.JSONWithStatus(c, customCode, data, http.StatusBadRequest)

// 分页响应
i18n.ListResponse(c, i18n.Success, items, total, page, perPage)
```

## 🔧 配置参考

### 默认配置

```yaml
default_language: "en"
fallback_language: "en"
locales_path: "locales"

locale_config:
  mode: "flat"
  languages: ["en", "zh-CN", "zh-TW"]

response_config:
  load_builtin: true
  auto_init: true

cache:
  enable: true
  size: 1000
  ttl: 3600
  l2_size: 5000
  enable_file: false

pool:
  enable: true
  size: 100
  warm_up: true
  languages: ["en", "zh-CN", "zh-TW"]

debug: false
enable_metrics: false
enable_watcher: false
```

### 环境变量配置

支持通过环境变量覆盖配置：

```bash
export GI18N_DEFAULT_LANGUAGE="zh-CN"
export GI18N_DEBUG="true"
export GI18N_ENABLE_METRICS="true"
export GI18N_CACHE_SIZE="5000"
export GI18N_POOL_SIZE="200"
```

### 从配置文件加载

```go
// 从 YAML 文件加载配置
err := i18n.InitFromConfigFile("config/i18n.yaml")

// 或者手动加载配置
config, err := i18n.LoadConfigFromFile("config/i18n.yaml")
if err != nil {
    panic(err)
}
err = i18n.InitWithConfig(config)
```

## 🚀 中间件选项

```go
// 自定义中间件选项
opts := i18n.MiddlewareOptions{
    HeaderKey:      "X-Language",           // 自定义 Header 键
    CookieName:     "lang",                  // 自定义 Cookie 名称
    QueryKey:       "lang",                  // 自定义查询参数键
    SupportedLangs: []string{"en", "zh-CN"}, // 支持的语言列表
    EnableCookie:   true,                    // 启用 Cookie 检测
    EnableQuery:    true,                    // 启用查询参数检测
}

r.Use(i18n.MiddlewareWithOpts(opts))
```

## 📈 性能优化

### 缓存策略

- **L1 缓存**: 内存缓存，最快访问
- **L2 缓存**: LRU 缓存，自动淘汰
- **对象池**: Localizer 对象复用，减少内存分配

### 调试和监控

```go
// 启用调试模式
config := i18n.Config{
    Debug:         true,  // 启用详细日志
    EnableMetrics: true,  // 启用性能指标
}

// 获取统计信息
stats := i18n.GetStats()
metrics := i18n.GetMetrics()

// 热更新
config.EnableWatcher = true
```

## 🔧 高级用法

### 热更新

```go
config := i18n.Config{
    EnableWatcher: true,  // 启用文件监听
}

// 手动重新加载
i18n.Reload()
```

### 多种翻译方式

```go
// 从 Gin Context 翻译
message := i18n.TFromGin(c, "WELCOME", data)

// 从 context.Context 翻译
message := i18n.T(ctx, "WELCOME", data)

// 使用指定语言翻译
message := i18n.GetService().TranslateWithLanguage(ctx, "zh-CN", "WELCOME", data)

// 获取当前语言
lang := i18n.GetLanguageFromGin(c)
lang = i18n.GetLanguage(ctx)
```

### 复数翻译

```go
// 语言文件中定义复数形式
{
  "ITEMS_COUNT": {
    "one": "{{.count}} item",
    "other": "{{.count}} items"
  }
}

// 使用复数翻译
message := i18n.GetService().Pluralize(ctx, "ITEMS_COUNT", count, data)
```

## 📖 文档

- [🚀 快速开始指南](docs/quickstart-guide.md)
- [🏗️ 架构设计](docs/library-architecture.md)
- [🔍 棕地架构分析](docs/brownfield-architecture.md)
- [📁 项目结构](docs/project-structure.md)
- [📊 API 文档](docs/api.md)
- [⚙️ 配置参考](docs/configuration.md)

## 🧪 示例项目

查看 `examples/` 目录中的完整示例：

- [quickstart](examples/quickstart/) - 基础使用示例
- [nested](examples/nested/) - 嵌套模式示例
- [custom-codes](examples/custom-codes/) - 自定义错误码示例
- [http-status-codes](examples/http-status-codes/) - HTTP 状态码示例
- [hybrid-codes](examples/hybrid-codes/) - 混合错误码示例

运行示例：

```bash
cd examples/quickstart
go run main.go
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/chenguowei/go-i18n.git
cd go-i18n

# 安装依赖
go mod download

# 运行测试
go test ./...

# 运行示例
go run examples/quickstart/main.go
```

### 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n) - 核心 i18n 引擎
- [gin-gonic/gin](https://github.com/gin-gonic/gin) - 优秀的 Web 框架
- [fsnotify](https://github.com/fsnotify/fsnotify) - 文件监听库