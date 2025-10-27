# 📁 GoI18n-Gin 项目目录结构

## 完整项目结构

**模块**: `github.com/chenguowei/go-i18n`

```
go-i18n/
├── 📄 go.mod                          # Go 模块定义
├── 📄 go.sum                          # 依赖锁定文件
├── 📄 README.md                       # 项目说明文档
├── 📄 CHANGELOG.md                    # 版本变更日志
├── 📄 LICENSE                         # 开源协议
├── 📄 Makefile                        # 构建脚本
│
├── 📄 i18n.go                         # 核心 API 和初始化（主入口）
├── 📄 middleware.go                   # Gin 中间件
├── 📄 translator.go                   # 翻译引擎
├── 📄 config.go                       # 配置管理
├── 📄 version.go                      # 版本信息
├── 📄 options.go                      # 选项模式配置
│
├── 📁 response/                       # 统一响应系统
│   ├── 📄 response.go                # 响应结构体定义
│   ├── 📄 codes.go                  # 错误码常量
│   ├── 📄 helper.go                 # 响应辅助函数
│   └── 📄 doc.go                    # 包文档
│
├── 📁 errors/                         # 错误定义包
│   ├── 📄 errors.go                 # 库专用错误类型
│   ├── 📄 codes.go                  # 错误码定义
│   └── 📄 doc.go                    # 包文档
│
├── 📁 internal/                       # 内部包（不对外暴露）
│   ├── 📁 cache/                      # 缓存系统
│   │   ├── 📄 memory.go              # L1 内存缓存
│   │   ├── 📄 lru.go                 # L2 LRU 缓存
│   │   ├── 📄 file.go                # L3 文件缓存
│   │   ├── 📄 manager.go             # 缓存管理器
│   │   └── 📄 interface.go           # 缓存接口定义
│   │
│   ├── 📁 pool/                       # 对象池系统
│   │   ├── 📄 localizer.go           # Localizer 对象池
│   │   ├── 📄 manager.go             # 池管理器
│   │   └── 📄 interface.go           # 池接口定义
│   │
│   ├── 📁 storage/                    # 存储层
│   │   ├── 📄 bundle.go              # 语言包管理
│   │   ├── 📄 loader.go              # 文件加载器
│   │   ├── 📄 watcher.go             # 热更新监听器
│   │   └── 📄 interface.go           # 存储接口定义
│   │
│   ├── 📁 parser/                     # 语言解析器
│   │   ├── 📄 accept_lang.go         # Accept-Language 解析
│   │   ├── 📄 header.go              # HTTP Header 解析
│   │   └── 📄 parser.go              # 解析器接口
│   │
│   ├── 📁 monitor/                    # 监控系统
│   │   ├── 📄 metrics.go             # 性能指标收集
│   │   ├── 📄 stats.go               # 统计信息
│   │   ├── 📄 debug.go               # 调试端点
│   │   └── 📄 prometheus.go          # Prometheus 集成
│   │
│   └── 📁 utils/                      # 内部工具
│       ├── 📄 logger.go              # 日志工具
│       ├── � validator.go           # 验证工具
│       └── 📄 file.go                # 文件工具
│
├── 📁 examples/                       # 使用示例
│   ├── 📁 quickstart/                 # 快速开始示例
│   │   ├── 📄 main.go                # 示例主程序
│   │   ├── 📄 go.mod                 # 示例模块
│   │   └── 📁 locales/               # 示例语言文件
│   │       ├── 📄 en.json
│   │       ├── 📄 zh-CN.json
│   │       └── 📄 zh-TW.json
│   │
│   ├── 📁 advanced/                   # 高级用法示例
│   │   ├── 📄 main.go                # 高级示例主程序
│   │   ├── 📄 config.yaml            # 配置文件示例
│   │   ├── 📄 go.mod                 # 示例模块
│   │   └── 📁 locales/               # 高级语言文件
│   │       ├── 📄 en.json
│   │       ├── 📄 zh-CN.json
│   │       ├── 📄 zh-TW.json
│   │       └── 📄 ja.json
│   │
│   ├── 📁 monitoring/                 # 监控集成示例
│   │   ├── 📄 main.go                # 监控示例主程序
│   │   ├── 📄 prometheus.go          # Prometheus 集成
│   │   ├── 📄 go.mod                 # 示例模块
│   │   └── 📁 locales/               # 监控语言文件
│   │       └── 📄 en.json
│   │
│   └── 📁 migration/                  # 迁移示例
│       ├── 📄 v1_to_v2.go           # 版本迁移示例
│       ├── 📄 legacy/main.go         # 旧版本示例
│       └── 📄 modern/main.go         # 新版本示例
│
├── 📁 configs/                        # 配置文件模板
│   ├── 📄 default.yaml               # 默认配置模板
│   ├── 📄 development.yaml           # 开发环境配置
│   ├── 📄 production.yaml            # 生产环境配置
│   ├── 📄 testing.yaml               # 测试环境配置
│   └── 📄 docker.yaml                # Docker 环境配置
│
├── 📁 locales/                        # 默认语言文件
│   ├── 📄 en.json                    # 英文语言包
│   ├── 📄 zh-CN.json                 # 简体中文语言包
│   ├── 📄 zh-TW.json                 # 繁体中文语言包
│   └── 📄 README.md                  # 语言文件说明
│
├── 📁 test/                           # 测试文件
│   ├── 📁 unit/                       # 单元测试
│   │   ├── 📄 gi18n_test.go          # 核心功能测试
│   │   ├── 📄 middleware_test.go     # 中间件测试
│   │   ├── 📄 translator_test.go     # 翻译器测试
│   │   ├── 📄 cache_test.go          # 缓存测试
│   │   ├── 📄 pool_test.go           # 对象池测试
│   │   └── 📄 config_test.go         # 配置测试
│   │
│   ├── 📁 integration/                # 集成测试
│   │   ├── 📄 gin_integration_test.go # Gin 集成测试
│   │   ├── 📄 hot_reload_test.go     # 热更新测试
│   │   └── 📁 fixtures/              # 测试数据
│   │       ├── 📁 locales/
│   │       └── 📄 configs/
│   │
│   ├── 📁 benchmark/                  # 性能测试
│   │   ├── 📄 translation_test.go    # 翻译性能测试
│   │   ├── 📄 cache_test.go          # 缓存性能测试
│   │   ├── 📄 pool_test.go           # 池性能测试
│   │   └── 📄 concurrency_test.go    # 并发测试
│   │
│   └── 📁 testdata/                   # 测试数据
│       ├── 📁 valid_locales/         # 有效语言文件
│       ├── 📁 invalid_locales/       # 无效语言文件
│       └── 📁 configs/               # 测试配置
│
├── 📁 docs/                           # 文档目录
│   ├── 📄 README.md                  # 文档首页
│   ├── 📄 api.md                     # API 文档
│   ├── 📄 configuration.md           # 配置文档
│   ├── 📄 migration.md               # 迁移指南
│   ├── 📄 performance.md             # 性能指南
│   ├── 📄 troubleshooting.md         # 故障排除
│   ├── 📄 contributing.md            # 贡献指南
│   ├── 📄 architecture.md            # 架构文档
│   └── 📁 images/                    # 文档图片
│       ├── 📄 architecture-diagram.svg
│       ├── 📄 performance-chart.png
│       └── 📄 workflow-diagram.svg
│
├── 📁 scripts/                        # 脚本目录
│   ├── 📄 build.sh                   # 构建脚本
│   ├── 📄 test.sh                    # 测试脚本
│   ├── 📄 release.sh                 # 发布脚本
│   ├── 📄 benchmark.sh               # 性能测试脚本
│   └── 📄 generate_locales.py        # 语言文件生成脚本
│
├── 📁 .github/                        # GitHub 配置
│   ├── 📁 workflows/                 # GitHub Actions
│   │   ├── 📄 ci.yml                 # 持续集成
│   │   ├── 📄 release.yml            # 自动发布
│   │   └── 📄 benchmark.yml          # 性能测试
│   │
│   ├── 📄 ISSUE_TEMPLATE/            # Issue 模板
│   │   ├── 📄 bug_report.md          # Bug 报告模板
│   │   └── 📄 feature_request.md     # 功能请求模板
│   │
│   └── 📄 PULL_REQUEST_TEMPLATE.md   # PR 模板
│
├── 📁 .gitignore                      # Git 忽略文件
├── 📁 .golangci.yml                   # GoLint 配置
├── 📁 Dockerfile                      # Docker 构建文件
├── 📁 docker-compose.yml              # Docker Compose 配置
└── 📁 .gitattributes                  # Git 属性配置
```

## 核心文件说明

### 📄 go.mod
```go
module github.com/your-org/go-i18n

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/nicksnyder/go-i18n/v2 v2.4.0
    golang.org/x/text v0.14.0
    gopkg.in/yaml.v3 v3.0.1
    github.com/fsnotify/fsnotify v1.7.0
    github.com/prometheus/client_golang v1.17.0
)
```

### 📄 Makefile
```makefile
.PHONY: build test clean benchmark lint docker

# 变量定义
BINARY_NAME=go-i18n
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# 构建
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/server

# 测试
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 性能测试
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./test/benchmark/

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t go-i18n:$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 go-i18n:$(VERSION)

# 示例
examples:
	@echo "Running examples..."
	cd examples/quickstart && go run .
	cd examples/advanced && go run .

# 发布
release:
	@echo "Creating release..."
	./scripts/release.sh $(VERSION)

# 依赖管理
deps-update:
	@echo "Updating dependencies..."
	go mod tidy
	go mod vendor

deps-check:
	@echo "Checking dependencies..."
	go mod verify

# 代码生成
generate:
	@echo "Generating code..."
	go generate ./...

# 安装开发工具
dev-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
```

### 📄 .gitignore
```
# Binaries
*.exe
*.dll
*.so
*.dylib
bin/
dist/

# Go specific
*.test
*.prof
coverage.out
coverage.html
vendor/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Temporary files
tmp/
temp/
*.tmp

# Config files with secrets
config/local.yaml
config/production.yaml
.env
.env.local

# Build artifacts
*.tar.gz
*.zip

# Documentation build
docs/build/
docs/_build/

# Examples build
examples/*/bin/
examples/*/dist/
```

### 📄 .golangci.yml
```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gocyclo
    - gofmt
    - goimports
    - misspell
    - unconvert
    - unparam
    - nakedret
    - prealloc
    - scopelint
    - gocritic

linters-settings:
  gocyclo:
    min-complexity: 15
  goimports:
    local-prefixes: github.com/your-org/go-i18n
  misspell:
    locale: US

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

## 使用方式

### 1. 作为库使用
```bash
go get github.com/your-org/go-i18n@v1.0.0
```

### 2. 开发环境设置
```bash
git clone https://github.com/your-org/go-i18n.git
cd go-i18n
make dev-tools
make deps-update
make test
```

### 3. 运行示例
```bash
make examples
# 或单独运行
cd examples/quickstart
go run .
```

### 4. 构建
```bash
make build
make docker-build
```

### 5. 运行测试
```bash
make test
make benchmark
make test-coverage
```

这个目录结构遵循了 Go 项目的最佳实践，清晰地分离了公共 API 和内部实现，便于维护和扩展。