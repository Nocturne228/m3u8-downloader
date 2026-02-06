# M3U8 Downloader v2.0 - Makefile
# 便捷的构建和开发工具

.PHONY: help build test clean release run cover lint fmt

# 变量定义
BINARY_NAME=m3u8-downloader
VERSION=2.0.0
MAIN_PATH=cmd/m3u8-downloader/main.go
BUILD_DIR=build
RELEASE_DIR=releases

# 默认目标
help:
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║           m3u8-downloader v2.0 - Build Tasks               ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "常见命令:"
	@echo "  make build          - 编译本地平台执行文件"
	@echo "  make release        - 交叉编译多平台发布包"
	@echo "  make test           - 运行单元测试"
	@echo "  make cover          - 生成测试覆盖率报告"
	@echo "  make clean          - 清理构建文件"
	@echo "  make run            - 编译并运行程序"
	@echo "  make lint           - 代码静态检查"
	@echo "  make fmt            - 代码格式化"
	@echo "  make verify         - 验证编译和测试"
	@echo ""
	@echo "用途示例:"
	@echo "  make build          # 快速编译开发版本"
	@echo "  make release        # 构建多平台发布包"
	@echo "  make test           # 运行所有单元测试"
	@echo ""

# 编译本地版本
build: deps
	@echo "🔨 编译中 ($(BINARY_NAME) v$(VERSION))..."
	@mkdir -p $(BUILD_DIR)
	GO111MODULE=on go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ 编译完成: ./$(BUILD_DIR)/$(BINARY_NAME)"

# 编译开发版本（带调试信息）
build-dev: deps
	@echo "🔨 编译开发版本（带调试信息）..."
	@mkdir -p $(BUILD_DIR)
	GO111MODULE=on go build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME)-dev $(MAIN_PATH)
	@echo "✅ 开发版本: ./$(BUILD_DIR)/$(BINARY_NAME)-dev"

# 交叉编译发布包
release: clean deps
	@echo "🚀 编译多平台发布包..."
	@mkdir -p $(RELEASE_DIR)
	
	@echo "  • darwin/amd64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -ldflags="-s -w" -o $(RELEASE_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	
	@echo "  • darwin/arm64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 GO111MODULE=on go build -ldflags="-s -w" -o $(RELEASE_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	@echo "  • linux/amd64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-s -w" -o $(RELEASE_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	
	@echo "  • linux/arm64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -ldflags="-s -w" -o $(RELEASE_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	@echo "  • windows/amd64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build -ldflags="-s -w" -o $(RELEASE_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "  • windows/arm64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 GO111MODULE=on go build -ldflags="-s -w" -o $(RELEASE_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)
	
	@echo "✅ 发布包构建完成:"
	@ls -lh $(RELEASE_DIR)/ | tail -n +2

# 下载依赖
deps:
	@echo "📦 下载依赖..."
	@go mod download
	@go mod verify

# 运行程序（编译+运行）
run: build
	@echo ""
	@./$(BUILD_DIR)/$(BINARY_NAME) -h

# 运行所有测试
test: deps
	@echo "🧪 运行单元测试..."
	@go test -v -timeout=5m ./...

# 测试特定包
test-pkg:
	@echo "🧪 运行特定包的测试..."
	@read -p "输入包名 (例如: internal/logger): " PKG; \
	go test -v ./$$PKG

# 生成测试覆盖率
cover: deps
	@echo "📊 生成测试覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告: coverage.html"

# 代码格式化
fmt:
	@echo "🎨 代码格式化..."
	@go fmt ./...
	@echo "✅ 格式化完成"

# 代码静态检查
lint:
	@echo "🔍 静态代码检查..."
	@go vet ./...
	@echo "✅ 检查完成"

# 竞态条件检查（仅限并发代码）
race: deps
	@echo "🏃 竞态条件检查..."
	@go test -race -timeout=5m ./internal/core ./internal/logger

# 验证编译和测试
verify: fmt lint test
	@echo ""
	@echo "✅ 所有验证通过！"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf $(BUILD_DIR) $(RELEASE_DIR) coverage.* *.prof
	@go clean ./...
	@echo "✅ 清理完成"

# 检查依赖更新
check-updates:
	@echo "📦 检查依赖更新..."
	@go list -u -m all

# 更新依赖
update-deps:
	@echo "📦 更新所有依赖..."
	@go get -u ./...
	@go mod tidy

# 显示项目信息
info:
	@echo "📋 项目信息:"
	@echo "  名称: $(BINARY_NAME)"
	@echo "  版本: $(VERSION)"
	@echo "  入口: $(MAIN_PATH)"
	@echo "  Go版本: $$(go version | cut -d' ' -f3)"
	@echo ""
	@echo "📊 代码统计:"
	@echo "  Go文件:"
	@wc -l $$(find . -name "*.go" -type f) | tail -1
	@echo ""
	@echo "🎯 所有目标: make help"

.DEFAULT_GOAL := help
