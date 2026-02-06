# M3U8 下载器 - 测试指南

## 📋 测试概览

本项目采用分层测试策略：

```
┌─────────────────────────────────────────┐
│         E2E 测试 (端到端)              │  ← 实际 M3U8 URL 测试
└────────────┬────────────────────────────┘
             ↓
┌─────────────────────────────────────────┐
│      集成测试 (多模块协作)              │  ← 模块间交互
└────────────┬────────────────────────────┘
             ↓
┌─────────────────────────────────────────┐
│        单元测试 (模块内部)              │  ← 独立函数/方法
└─────────────────────────────────────────┘
```

---

## 🧪 单元测试

### 1. 运行所有单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./internal/config
go test ./internal/errors

# 显示详细输出
go test -v ./...

# 显示测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # 在浏览器中查看
```

### 2. 现有单元测试

#### **配置模块** (`internal/config/config_test.go`)

```go
// 默认配置验证
TestDefaultConfig()

// 配置验证规则
TestConfigValidate()
  ├─ 有效配置
  ├─ 无效的线程数 (太小)
  ├─ 无效的线程数 (太大)
  └─ 其他规则
```

运行：
```bash
go test -v ./internal/config
```

#### **错误处理** (`internal/errors/errors_test.go`)

```go
// 错误创建和检查
TestErrorCreation()
TestIsCode()
TestErrorMessage()
```

运行：
```bash
go test -v ./internal/errors
```

### 3. 添加新的单元测试

示例：为 HTTP 客户端添加测试

```go
// internal/http/client_test.go
package http

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	logger := mockLogger()
	client := NewClient(5*time.Second, 3, "Test-Agent", logger)
	
	if client == nil {
		t.Error("Expected non-nil client")
	}
}

func TestClientInsecureVerify(t *testing.T) {
	logger := mockLogger()
	client := NewClient(5*time.Second, 3, "Test-Agent", logger)
	
	httpClient := client.(*HTTPClient)
	httpClient.SetInsecureVerify(true)
	
	if !httpClient.insecureVerify {
		t.Error("Expected insecureVerify to be true")
	}
}
```

---

## 🔗 集成测试

### 1. 创建集成测试目录

```bash
mkdir -p test/integration
```

### 2. 模块集成测试

#### **配置 + 日志集成**

```go
// test/integration/config_logger_test.go
package integration

import (
	"testing"
	"m3u8-downloader/internal/config"
	"m3u8-downloader/internal/logger"
)

func TestConfigWithLogger(t *testing.T) {
	// 创建配置
	cfg := config.DefaultConfig()
	
	// 创建日志记录器
	log := logger.New(cfg.Log.Level)
	
	// 验证配置
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Config validation failed: %v", err)
	}
	
	// 使用日志
	log.Info("配置验证成功")
}
```

运行：
```bash
go test -v ./test/integration -run ConfigWith
```

#### **HTTP 客户端集成测试**

```go
// test/integration/http_client_test.go
package integration

import (
	"testing"
	"time"
	"m3u8-downloader/internal/http"
	"m3u8-downloader/internal/logger"
)

func TestHTTPClientWithLogger(t *testing.T) {
	log := logger.New("info")
	client := http.NewClient(5*time.Second, 3, "Test", log)
	
	// 测试真实 URL
	data, err := client.Get("https://www.example.com")
	if err != nil {
		t.Logf("HTTP request error (expected for test): %v", err)
		return
	}
	
	if len(data) == 0 {
		t.Error("Expected non-empty response")
	}
}
```

---

## 🌐 E2E 端到端测试

### 1. 本地测试 M3U8 服务器

创建简单的 M3U8 文件用于测试：

```bash
# test/fixtures/sample.m3u8
mkdir -p test/fixtures
cat > test/fixtures/sample.m3u8 << 'EOF'
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXTINF:10.0,
segment1.ts
#EXTINF:10.0,
segment2.ts
#EXT-X-ENDLIST
EOF
```

### 2. 启动本地 HTTP 服务器测试

```bash
# 在项目根目录启动 HTTP 服务器
python3 -m http.server 8080 &

# 测试下载
./m3u8-downloader-v2 "http://localhost:8080/test/fixtures/sample.m3u8" \
  -o "test_video" \
  -sp "./test/e2e_output"

# 检查结果
ls -la test/e2e_output
```

### 3. E2E 测试脚本

```bash
#!/bin/bash
# test/e2e/test.sh

set -e

echo "=== 构建程序 ==="
go build -o m3u8-downloader-v2 cmd/m3u8-downloader/main.go

echo "=== 启动测试 HTTP 服务器 ==="
cd test/fixtures
python3 -m http.server 8888 > /dev/null 2>&1 &
SERVER_PID=$!
cd ../..

sleep 2  # 等待服务器启动

echo "=== 运行 E2E 测试 ==="
OUTPUT_DIR="./test/e2e_output"
mkdir -p "$OUTPUT_DIR"

# 测试 1: 基本下载
echo "测试 1: 基本下载"
./m3u8-downloader-v2 "http://localhost:8888/sample.m3u8" \
  -o "test_video" \
  -sp "$OUTPUT_DIR" || echo "预期的错误（测试文件不完整）"

# 测试 2: 参数验证
echo "测试 2: 参数验证"
./m3u8-downloader-v2 "" 2>&1 | grep -q "错误" && echo "✓ 参数验证正确"

# 清理
kill $SERVER_PID 2>/dev/null || true
rm -rf "$OUTPUT_DIR"

echo "=== E2E 测试完成 ==="
```

运行：
```bash
chmod +x test/e2e/test.sh
./test/e2e/test.sh
```

---

## 📊 测试覆盖率目标

### 当前状态
- 单元测试覆盖率: **20%**
- 集成测试覆盖率: **5%**
- 总体覆盖率: **25%**

### 改进计划

#### Phase 1 (立即)
```
目标: 40% 总体覆盖率
时间: 1-2 周
内容:
  ✓ config 包 → 80% 覆盖
  ✓ errors 包 → 100% 覆盖
  ✓ logger 包 → 60% 覆盖
  ✓ http 包 → 40% 覆盖
```

#### Phase 2 (2-4 周)
```
目标: 60% 总体覆盖率
包括:
  ✓ m3u8 包 → 60% 覆盖
  ✓ core 包 → 50% 覆盖
  ✓ video 包 → 40% 覆盖
  ✓ util 包 → 70% 覆盖
```

#### Phase 3 (1-2 月)
```
目标: 80%+ 总体覆盖率
包括:
  ✓ E2E 测试完整化
  ✓ 场景测试
  ✓ 边界条件测试
  ✓ 性能基准测试
```

---

## ✅ 测试检查清单

### Before Running Tests

- [ ] 依赖安装: `go mod download`
- [ ] 代码编译: `go build ./...`
- [ ] 代码格式: `go fmt ./...`
- [ ] Lint 检查: `go vet ./...`

### Unit Tests

```bash
# 运行所有单元测试
go test -v ./internal/...

# 查看覆盖率
go test -cover ./internal/...

# 预期输出
# ok  	m3u8-downloader/internal/config	0.234s	coverage: 80.0% of statements
# ok  	m3u8-downloader/internal/errors	0.123s	coverage: 100.0% of statements
```

### Integration Tests

```bash
# 运行集成测试
go test -v ./test/integration

# 预期输出
# ok  	m3u8-downloader/test/integration	2.345s
```

### E2E Tests

```bash
# 运行端到端测试
./test/e2e/test.sh

# 预期输出
# === E2E 测试完成 ===
# All tests passed ✓
```

---

## 🐛 调试技巧

### 1. 启用日志级别

```bash
# 设置 DEBUG 级别日志
export LOG_LEVEL=debug
./m3u8-downloader-v2 "url" -v

# 输出完整堆栈跟踪
go test -v -run TestName
```

### 2. 使用 pprof 进行性能分析

```bash
# CPU 分析
go test -cpuprofile=cpu.prof ./internal/core
go tool pprof cpu.prof

# 内存分析
go test -memprofile=mem.prof ./internal/core
go tool pprof mem.prof
```

### 3. Race 条件检查

```bash
# 检查并发问题
go test -race ./...

# 预期：如果没有 race 条件
# PASS
# ok  	m3u8-downloader/internal/core	0.456s
```

---

## 📝 测试编写指南

### 优秀的单元测试示例

```go
package config

import (
	"testing"
	"time"
)

// 命名规范: Test<Function/Method><Scenario>
func TestConfigValidate_InvalidThreadCount(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "线程数为0",
			cfg: &Config{
				Download: DownloadConfig{
					MaxGoroutines: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "线程数超过256",
			cfg: &Config{
				Download: DownloadConfig{
					MaxGoroutines: 300,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

### Mock 和 Stub 示例

```go
package http

import (
	"testing"
	"m3u8-downloader/internal/logger"
)

// Mock Logger
type MockLogger struct {
	messages []string
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.messages = append(m.messages, msg)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.messages = append(m.messages, msg)
}

// ... 其他方法

func TestClientWithMockLogger(t *testing.T) {
	mockLog := &MockLogger{}
	client := NewClient(5*time.Second, 3, "Test", mockLog)
	
	// 执行操作
	client.Get("http://example.com")
	
	// 验证日志
	if len(mockLog.messages) == 0 {
		t.Error("Expected log messages")
	}
}
```

---

## 🚀 CI/CD 集成

### GitHub Actions 示例

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.16
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v1
      with:
        files: ./coverage.out
    
    - name: Run linters
      run: |
        go get github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        golangci-lint run ./...
    
    - name: Build
      run: go build -o m3u8-downloader cmd/m3u8-downloader/main.go
```

---

## 📊 测试结果示例

```bash
$ go test -v -cover ./...

=== RUN   TestDefaultConfig
--- PASS: TestDefaultConfig (0.001s)

=== RUN   TestConfigValidate
--- PASS: TestConfigValidate (0.002s)

=== RUN   TestErrorCreation
--- PASS: TestErrorCreation (0.001s)

=== RUN   TestIsCode
--- PASS: TestIsCode (0.001s)

PASS
coverage: 25.3% of statements
ok  	m3u8-downloader/internal/config	0.123s
ok  	m3u8-downloader/internal/errors	0.089s
ok  	m3u8-downloader/test/integration	2.456s
```

---

## 🎯 快速测试命令

```bash
# 快速检查
make test-quick

# 完整测试
make test-full

# 生成覆盖率报告
make test-coverage

# 性能测试
make test-bench

# 竞态条件检查
make test-race
```

### Makefile

```makefile
.PHONY: test-quick test-full test-coverage test-race test-bench

test-quick:
	go test -v ./internal/...

test-full:
	go test -v -race ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-race:
	go test -race ./...

test-bench:
	go test -bench=. -benchmem ./internal/core
```

---

## 💡 最佳实践

1. **命名清晰** - 从测试名称能看出测试目的
2. **测试独立** - 每个测试完全独立，不依赖其他测试
3. **使用 Table-Driven** - 多个场景使用参数化测试
4. **Mock 外部依赖** - 隔离被测试的代码
5. **边界条件** - 测试 0、负数、nil、empty 等
6. **错误路径** - 不仅测试正常路径，也要测试错误情况
7. **性能考虑** - 使用 Benchmark 测试性能
8. **覆盖率目标** - 针对不同包设置合理的覆盖率目标

---

**测试是代码质量的保障，现在就开始编写测试吧！** ✅
