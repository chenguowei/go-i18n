# 🌐 GoI18n-Gin 棕地架构文档

## 概述

本文档捕获 **GoI18n-Gin** 代码库的当前状态，包括技术债务、变通方法和实际使用模式。它作为 AI 代理进行增强工作的参考。

### 文档范围

全面文档化整个系统，基于实际代码分析更新现有架构文档。

### 变更日志

| 日期       | 版本 | 描述                       | 作者         |
| ---------- | ---- | -------------------------- | ------------ |
| 2025-10-30 | 1.0  | 初始棕地分析和代码对齐更新 | Winston (AI) |

## 快速参考 - 关键文件和入口点

### 理解系统的关键文件

- **主入口**: `i18n.go` (核心服务和全局API)
- **Gin 中间件**: `middleware.go` (语言检测和上下文注入)
- **翻译引擎**: `translator.go` (核心翻译逻辑)
- **响应系统**: `response_functions.go`, `codes.go` (统一响应处理)
- **配置管理**: `config.go` (配置加载和验证)
- **版本信息**: `version.go` (版本和构建信息)

### 内部组件

- **缓存系统**: `internal/cache.go` (多层缓存实现)
- **对象池**: `internal/pool.go` (Localizer 对象池)
- **文件监听**: `internal/watcher.go` (热更新支持)
- **接口定义**: `internal/interfaces.go` (内部接口)

## 高层架构

### 技术总结

GoI18n-Gin 是一个为 Gin 框架设计的国际化库，采用多层缓存、对象池化和热更新等高性能设计模式。

### 实际技术栈 (基于 go.mod)

| 类别        | 技术                          | 版本    | 说明                      |
| ----------- | ----------------------------- | ------- | ------------------------- |
| 运行时      | Go                            | 1.21    | 最低要求 Go 1.21          |
| Web 框架    | Gin                           | v1.9.1  | HTTP 中间件和上下文处理   |
| i18n 核心   | nicksnyder/go-i18n/v2         | v2.4.0  | 底层翻译引擎              |
| 文本处理    | golang.org/x/text              | v0.14.0 | 语言解析和匹配            |
| 文件监听    | fsnotify                       | v1.7.0  | 热更新文件监听            |
| 配置解析    | gopkg.in/yaml.v3               | v3.0.1  | YAML 配置文件支持         |
| 测试框架    | testify                       | v1.8.3  | 单元测试和集成测试        |

### 仓库结构实际检查

- **类型**: 单一 Go 模块 (monorepo 不适用)
- **包管理**: Go modules
- **依赖管理**: go.mod/go.sum

## 源码树和模块组织

### 项目结构 (实际)

```text
go-i18n/
├── i18n.go                      # 核心 API 和服务主入口
├── middleware.go                # Gin 中间件实现
├── translator.go                # 翻译引擎核心逻辑
├── config.go                    # 配置加载和验证
├── version.go                   # 版本信息
├── response_functions.go        # 统一响应函数
├── codes.go                     # 错误码定义和管理
│
├── internal/                    # 内部包 (不对外暴露)
│   ├── interfaces.go           # 内部接口定义
│   ├── cache.go                # 缓存管理器实现
│   ├── pool.go                 # 对象池管理器
│   ├── loader.go               # 语言文件加载器
│   └── watcher.go              # 文件监听器
│
├── errors/                      # 错误定义包
│   └── errors.go              # 库专用错误类型
│
├── cmd/                         # 命令行工具
│   └── migrate/               # 迁移工具
│       └── main.go
│
├── examples/                    # 使用示例
│   ├── quickstart/            # 快速开始示例
│   ├── nested/                # 嵌套模式示例
│   ├── custom-codes/          # 自定义错误码示例
│   ├── http-status-codes/     # HTTP 状态码示例
│   └── hybrid-codes/          # 混合错误码示例
│
├── locales/                     # 默认语言文件
├── docs/                        # 文档目录
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

### 关键模块及其用途

- **核心服务**: `i18n.go` - 全局服务管理、单例模式、API 入口点
- **语言检测**: `middleware.go` - 多源语言检测、优先级管理、上下文注入
- **翻译执行**: `translator.go` - 本地化器管理、缓存集成、降级处理
- **响应统一**: `response_functions.go` - 统一 JSON 响应格式、错误码映射
- **配置系统**: `config.go` - 多源配置加载、环境变量支持、验证

## 数据模型和 API

### 核心接口定义

基于 `internal/interfaces.go` 和实际实现：

```go
// Translator 翻译器接口
type Translator interface {
    Translate(ctx context.Context, messageID string, templateData ...map[string]interface{}) string
    TranslateWithLanguage(ctx context.Context, lang, messageID string, templateData ...map[string]interface{}) string
    Localizer(ctx context.Context) *i18n.Localizer
    LocalizerWithLanguage(ctx context.Context, lang string) *i18n.Localizer
    LoadLocales(localesPath string) error
}

// CacheManager 缓存管理器接口 (internal/interfaces.go)
type CacheManager interface {
    Get(key string) (string, bool)
    Set(key string, value string)
    Clear()
    GetStats() CacheStats
    Close() error
}

// PoolManager 对象池管理器接口 (internal/interfaces.go)
type PoolManager interface {
    Get(language string) *i18n.Localizer
    WarmUp(languages []string)
    GetStats() PoolStats
    Close() error
}
```

### API 规范

**核心 API** (从 `i18n.go` 提取):

- `Init() error` - 使用默认配置初始化
- `InitWithConfig(config Config) error` - 使用自定义配置初始化
- `InitFromConfigFile(configPath string) error` - 从配置文件初始化
- `Middleware() gin.HandlerFunc` - 返回 Gin 中间件
- `T(ctx context.Context, messageID string, templateData ...map[string]interface{}) string` - 翻译函数
- `TFromGin(c *gin.Context, messageID string, templateData ...map[string]interface{}) string` - 从 Gin Context 翻译
- `GetLanguage(ctx context.Context) string` - 获取当前语言
- `Reload() error` - 重新加载语言文件
- `Close() error` - 关闭系统

**响应 API** (从 `response_functions.go` 提取):

- `JSON(c *gin.Context, code Code, data interface{})` - 标准 JSON 响应
- `JSONWithTemplate(c *gin.Context, code Code, data interface{}, templateData map[string]interface{})` - 支持模板参数的响应
- `Error(c *gin.Context, code Code, message string)` - 错误响应
- `SuccessResponse(c *gin.Context, data interface{})` - 成功响应便捷方法

## 技术债务和已知问题

### 关键技术债务

1. **硬编码语言列表**: `translator.go` 中硬编码了支持的语言文件列表，而不是动态扫描目录
2. **缓存实现简化**: `internal/cache.go` 中的缓存实现相对基础，缺少高级特性如分布式缓存
3. **全局状态依赖**: 严重依赖全局单例，可能在测试和多应用场景中造成问题
4. **错误处理不统一**: 不同模块的错误处理模式不一致

### 变通方法和注意事项

- **语言检测标准化**: `middleware.go` 中实现了语言代码标准化，处理常见别名如 `zh_cn` -> `zh-cn`
- **降级翻译机制**: `translator.go` 实现了多层降级：指定语言 -> 降级语言 -> 格式化 messageID
- **调试日志分散**: 调试信息分布在各个文件中，使用不同的日志格式

### 集成约束

- **Gin 框架紧耦合**: 中间件直接依赖 Gin Context，难以适配其他框架
- **文件系统依赖**: 语言文件加载依赖传统文件系统，不支持嵌入式资源
- **JSON 格式锁定**: 主要支持 JSON 格式的语言文件，其他格式支持有限

## 集成点和外部依赖

### 外部服务

| 服务  | 用途           | 集成类型     | 关键文件              |
| ----- | -------------- | ------------ | --------------------- |
| Gin   | Web 框架       | 直接依赖     | `middleware.go`       |
| fsnotify | 文件监听     | SDK 集成     | `internal/watcher.go` |

### 内部集成点

- **应用层集成**: 通过 Gin 中间件 `middleware.go` 集成
- **响应系统集成**: 通过 `response_functions.go` 与业务逻辑集成
- **配置系统集成**: 支持环境变量、YAML 文件配置

## 开发和部署

### 本地开发设置

实际可工作的步骤：

1. **克隆仓库**: `git clone https://github.com/chenguowei/go-i18n`
2. **安装依赖**: `go mod download`
3. **准备语言文件**: 确保 `locales/` 目录包含所需语言文件
4. **运行示例**: `cd examples/quickstart && go run main.go`

已知设置问题：
- 语言文件必须使用 JSON 格式
- 降级语言必须明确配置
- 调试模式需要在配置中显式启用

### 构建和部署过程

- **构建命令**: `go build ./...` (标准 Go 构建)
- **测试命令**: `go test ./...` (单元测试和集成测试)
- **跨平台构建**: 支持标准的 Go 交叉编译

## 测试实际情况

### 当前测试覆盖

- **单元测试**: 基础覆盖，主要在 `*_test.go` 文件中
- **集成测试**: 有限，主要在 `examples/` 目录中体现
- **基准测试**: 未发现专门的基准测试文件
- **手动测试**: 主要通过运行示例进行验证

### 运行测试

```bash
go test ./...                   # 运行所有测试
go run examples/quickstart/main.go  # 运行快速开始示例
go run examples/custom-codes/main.go # 运行自定义错误码示例
```

## 架构不一致和优化机会

### 文档与代码不一致

1. **缓存架构**: 现有架构文档描述的三层缓存 (L1/L2/L3) 与实际 `internal/cache.go` 实现不符
2. **对象池架构**: 文档描述的复杂池化策略与 `internal/pool.go` 的简化实现有差距
3. **监控系统**: 文档描述的监控端点在实际代码中未完全实现

### 性能优化机会

1. **内存分配**: 减少翻译过程中的字符串拼接和哈希计算
2. **并发安全**: 优化全局状态的并发访问控制
3. **缓存策略**: 实现更智能的缓存失效和预热机制

## 代码质量和设计模式

### 设计模式使用

- **单例模式**: 全局服务管理 (`i18n.go`)
- **中间件模式**: Gin 集成 (`middleware.go`)
- **策略模式**: 语言检测策略
- **工厂模式**: 配置和服务创建

### 代码组织原则

- **关注点分离**: 核心、缓存、池化等职责分离良好
- **接口抽象**: 内部组件通过接口解耦
- **配置驱动**: 通过配置控制行为

## 附录 - 实用命令和脚本

### 常用命令

```bash
# 开发
go mod tidy                      # 整理依赖
go run examples/quickstart/main.go  # 快速测试
go test ./...                    # 运行测试

# 构建
go build -v ./...               # 详细构建
go install .                    # 本地安装

# 调试
GI18N_DEBUG=true go run examples/quickstart/main.go  # 启用调试模式
```

### 调试和故障排除

- **调试模式**: 设置配置中 `Debug: true` 或环境变量 `GI18N_DEBUG=true`
- **日志输出**: 调试信息会输出到标准控制台，包含语言检测、文件加载、翻译时间等信息
- **常见问题**:
  - 语言文件未找到：检查 `LocalesPath` 配置
  - 翻译失败：检查 messageID 是否在语言文件中定义
  - 缓存问题：禁用缓存进行调试

## 配置参考

### 默认配置 (从 `config.go` 提取)

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
  ttl: 3600  # 1小时 (秒)
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

### 环境变量支持

支持的环境变量 (从代码中提取):
- `GI18N_DEFAULT_LANGUAGE`
- `GI18N_FALLBACK_LANGUAGE`
- `GI18N_LOCALES_PATH`
- `GI18N_DEBUG`
- `GI18N_ENABLE_METRICS`
- `GI18N_ENABLE_WATCHER`
- `GI18N_CACHE_ENABLE`
- `GI18N_CACHE_SIZE`
- `GI18N_CACHE_TTL`
- `GI18N_POOL_ENABLE`
- `GI18N_POOL_SIZE`

---

**注意**: 此文档反映代码库的实际状态，包括所有技术债务、变通方法和约束。用于指导 AI 代理进行有效的工作而不破坏现有功能。