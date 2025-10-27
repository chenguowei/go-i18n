.PHONY: build test clean benchmark lint docker example

# 变量定义
BINARY_NAME=go-i18n
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
LDFLAGS=-ldflags "-X github.com/chenguowei/go-i18n/version.Version=$(VERSION) -X github.com/chenguowei/go-i18n/version.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

# 默认目标
all: build

# 构建
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/server

# 测试
test:
	@echo "Running tests..."
	go test -v ./...

# 测试覆盖率
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 性能测试
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./test/benchmark/

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# 代码格式化
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.log

# 运行示例
example:
	@echo "Running quickstart example..."
	cd examples/quickstart && go run -o ../../bin/example .

# 运行示例并测试不同语言
example-test:
	@echo "Testing example with different languages..."
	@echo "Testing English:"
	curl -s http://localhost:8080/welcome | jq .
	@echo "Testing Chinese:"
	curl -s -H "Accept-Language: zh-CN" http://localhost:8080/welcome | jq .
	@echo "Testing with parameters:"
	curl -s "http://localhost:8080/hello?name=Alice" | jq .

# Docker 构建
docker-build:
	@echo "Building Docker image..."
	docker build -t go-i18n:$(VERSION) .

# Docker 运行
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 go-i18n:$(VERSION)

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
	go install golang.org/x/tools/cmd/goimports@latest

# 热重开发
dev:
	@echo "Starting development server with hot reload..."
	air -c .air.toml

# 发布
release:
	@echo "Creating release..."
	./scripts/release.sh $(VERSION)

# 创建标签
tag:
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# 文档生成
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation available at http://localhost:6060"

# 安装
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/

# 卸载
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

# 静态分析
static:
	@echo "Running static analysis..."
	gosec ./...
	govulncheck ./...

# 安全扫描
security:
	@echo "Running security scan..."
	nancy sleuth

# 集成测试
integration-test:
	@echo "Running integration tests..."
	go test -v -tags=integration ./test/integration/...

# 端到端测试
e2e-test:
	@echo "Running end-to-end tests..."
	go test -v -tags=e2e ./test/e2e/

# 性能分析
profile:
	@echo "Running CPU profiling..."
	go test -cpuprofile=cpu.prof -bench=. ./test/benchmark/
	go tool pprof cpu.prof

# 内存分析
memory-profile:
	@echo "Running memory profiling..."
	go test -memprofile=mem.prof -bench=. ./test/benchmark/
	go tool pprof mem.prof

# 竞争检测
race:
	@echo "Running race detection tests..."
	go test -race -v ./...

# 代码质量检查
quality: lint static security

# 完整的 CI 流程
ci: deps-check fmt lint test security

# 预提交检查
pre-commit: fmt lint test

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  benchmark      - Run benchmarks"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  clean          - Clean build artifacts"
	@echo "  example        - Run quickstart example"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  deps-update    - Update dependencies"
	@echo "  dev-tools      - Install development tools"
	@echo "  dev            - Start development server"
	@echo "  release        - Create release"
	@echo "  tag            - Create git tag"
	@echo "  docs           - Generate documentation"
	@echo "  install        - Install binary"
	@echo "  uninstall      - Uninstall binary"
	@echo "  static         - Run static analysis"
	@echo "  security       - Run security scan"
	@echo "  integration-test - Run integration tests"
	@echo "  e2e-test       - Run end-to-end tests"
	@echo "  profile        - Run CPU profiling"
	@echo "  memory-profile - Run memory profiling"
	@echo "  race           - Run race detection"
	@echo "  quality        - Run all quality checks"
	@echo "  ci             - Run CI pipeline"
	@echo "  pre-commit     - Run pre-commit checks"
	@echo "  help           - Show this help message"