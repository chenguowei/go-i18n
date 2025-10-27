# GoI18n-Gin 🌐

[![GoDoc](https://godoc.org/github.com/chenguowei/go-i18n?status.svg)](https://godoc.org/github.com/chenguowei/go-i18n)
[![Build Status](https://github.com/chenguowei/go-i18n/workflows/CI/badge.svg)](https://github.com/chenguowei/go-i18n/actions)
[![Coverage](https://codecov.io/gh/chenguowei/go-i18n/branch/main/graph/badge.svg)](https://codecov.io/gh/chenguowei/go-i18n)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenguowei/go-i18n)](https://goreportcard.com/report/github.com/chenguowei/go-i18n)

🌍 一个专为 Gin 框架设计的开箱即用多语言库，提供高性能、易集成、生产就绪的国际化解决方案。

## ✨ 特性

- 🚀 **开箱即用** - 3行代码完成集成
- ⚡ **高性能** - 多层缓存 + 对象池，响应时间 < 0.1ms
- 🛡️ **生产就绪** - 完善的错误处理和降级机制
- 🔥 **热更新** - 无需重启，动态加载语言文件
- 📊 **监控友好** - 内置指标和调试端点
- 🌐 **多语言支持** - 支持多种语言指定方式
- 🎯 **零侵入** - 标准Gin中间件，不影响现有代码

## 🚀 快速开始

### 安装

```bash
go get github.com/chenguowei/go-i18n@latest
```

### 基础使用

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
    "github.com/chenguowei/go-i18n/response"
)

func main() {
    // 1. 初始化 i18n
    if err := i18n.Init(); err != nil {
        panic(err)
    }

    // 2. 添加中间件
    r := gin.Default()
    r.Use(i18n.Middleware())

    // 3. 使用翻译
    r.GET("/hello", func(c *gin.Context) {
        message := i18n.TFromGin(c, "WELCOME")
        response.JSON(c, response.Success, map[string]interface{}{
            "message": message,
        })
    })

    r.Run(":8080")
}
```

### 高级配置

```go
config := i18n.Config{
    DefaultLanguage:  "zh-CN",
    FallbackLanguage: "en",
    LocalesPath:      "./locales",
    Cache: i18n.CacheConfig{
        Enable: true,
        Size:   5000,
        TTL:    2 * time.Hour,
    },
    Pool: i18n.PoolConfig{
        Enable: true,
        Size:   200,
        WarmUp: true,
    },
}

i18n.InitWithConfig(config)
```

## 📖 文档

- [快速开始指南](docs/quickstart-guide.md)
- [架构设计](docs/library-architecture.md)
- [API 文档](docs/api.md)
- [配置参考](docs/configuration.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！详见 [贡献指南](CONTRIBUTING.md)。

## 📚 更多文档

- [项目架构](docs/library-architecture.md)
- [项目总结](PROJECT_SUMMARY.md)
- [快速开始指南](docs/quickstart-guide.md)
- [项目结构](docs/project-structure.md)

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件。