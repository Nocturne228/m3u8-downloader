# Development Guide

## 项目开发指南

本文档提供了对 m3u8-downloader 项目进行开发和贡献的指导。

### 环境设置

#### 前置需求
- Go 1.16 或更高版本
- FFmpeg 2.8 或更高版本（用于视频合并）
- Git（版本控制）

#### 项目结构概览

```
m3u8-downloader/
├── cmd/                          # 应用程序入口
│   └── m3u8-downloader/
│       └── main.go              # CLI主入口
├── internal/                     # 内部包（不对外暴露）
│   ├── config/                  # 配置管理
│   │   ├── config.go
│   │   └── config_test.go
│   ├── logger/                  # 日志系统
│   │   ├── logger.go
│   │   └── logger_test.go
│   ├── errors/                  # 错误处理
│   │   ├── errors.go
│   │   └── errors_test.go
│   ├── http/                    # HTTP 客户端（重试机制）
│   │   └── client.go
│   ├── m3u8/                    # M3U8 播放列表解析
│   │   ├── parser.go
│   │   └── fetcher.go
│   ├── core/                    # 核心下载逻辑
│   │   ├── manager.go           # 下载管理器
│   │   └── application.go       # 应用协调
│   ├── video/                   # 视频处理
│   │   └── merger.go            # FFmpeg 合并
│   └── util/                    # 工具函数
│       └── util.go
├── test/                        # 测试文件
│   └── unit/
│       └── config_test.go
├── go.mod                       # 模块定义
├── go.sum                       # 依赖哈希
├── README.md                    # 项目说明
├── docs/                        # 文档目录
│   ├── TEST.md                  # 测试指南
│   ├── REFACTORING_REPORT_V2.0.md  # 重构报告
│   └── DEVELOPMENT.md          # 本文件
└── build-release.sh            # 发布脚本
```

### 核心模块说明

#### config 配置管理
- **文件**: `internal/config/config.go`
- **职责**: 应用程序配置读取和校验
- **关键类型**:
  - `Config`: 主配置结构体 (HTTP、下载、FFmpeg、日志配置)
  - `ConfigError`: 配置错误类型
- **主要方法**:
  - `DefaultConfig()`: 获取默认配置
  - `Validate()`: 配置校验

**扩展建议**:
- 支持从YAML/JSON文件读取配置
- 环境变量覆盖配置值

#### logger 日志系统
- **文件**: `internal/logger/logger.go`
- **职责**: 结构化日志记录
- **关键接口**:
  - `Logger`: 日志接口 (Debug, Info, Warn, Error, Fatal)
  - `ConsoleLogger`: 控制台实现
- **日志级别**: Debug < Info < Warn < Error < Fatal

**扩展建议**:
- 添加文件日志后端
- 支持JSON日志格式
- 添加日志轮转(rotation)

#### errors 错误处理
- **文件**: `internal/errors/errors.go`
- **职责**: 统一的错误定义和处理
- **错误代码** (13种):
  - HTTPRequest, M3U8Parse, FileIO
  - VideoMerge, InvalidConfig 等
- **错误支持**: 错误链(wrapping)、错误代码匹配

**扩展建议**:
- 错误国际化(i18n)
- 错误恢复建议提示

#### http HTTP客户端
- **文件**: `internal/http/client.go`
- **职责**: HTTP请求with重试机制
- **特性**:
  - 指数退避(exponential backoff)重试
  - 网络错误检测
  - 自定义超时和头部
  - Cookie支持

**扩展建议**:
- SOCKS5代理支持
- HTTP/2推送
- 请求签名支持

#### m3u8 播放列表处理
- **文件**: `internal/m3u8/parser.go`, `internal/m3u8/fetcher.go`
- **职责**: M3U8 清单解析和获取
- **关键类型**:
  - `TsSegment`: 媒体片段信息
  - `EncryptionKey`: 加密密钥
  - `Manifest`: 完整播放列表
- **特殊支持**: 加密分片、HLS变体流

**扩展建议**:
- 支持DASH (MPEG-DASH)
- 支持Master播放列表(多码率自适应)

#### core 核心下载引擎
- **文件**: `internal/core/manager.go`, `internal/core/application.go`
- **职责**: 整体下载流程协调
- **关键类**:
  - `DownloadManager`: 单个下载任务管理
  - `Application`: 应用级协调 (依赖注入)
- **特性**:
  - Goroutine并发控制
  - 实时进度显示
  - 自动重试失败片段

**扩展建议**:
- 断点续传(resume)
- 频段加速下载
- 下载历史记录

#### video 视频处理
- **文件**: `internal/video/merger.go`
- **职责**: FFmpeg视频合并
- **特性**:
  - concat demuxer合并
  - 输出验证
- **支持格式**: TS → MP4

**扩展建议**:
- 多种输出格式(MKV, WebM等)
- 视频转码选项
- 字幕处理

#### util 工具函数
- **文件**: `internal/util/util.go`
- **功能**:
  - 文件操作: PathExists, EnsureDir, ReadFile, WriteFile
  - AES加密: PKCS7Padding, AesDecrypt, AesEncrypt
  - TS处理: RemoveTSPadding, ListTSFiles

### 开发工作流

#### 1. 本地构建
```bash
# 构建可执行文件
go build -o m3u8-downloader cmd/m3u8-downloader/main.go

# 构建带调试信息的版本
go build -gcflags="all=-N -l" -o m3u8-downloader cmd/m3u8-downloader/main.go

# 快速测试
./m3u8-downloader -h
```

#### 2. 运行测试
```bash
# 运行所有单元测试
go test ./...

# 运行特定包的测试
go test ./internal/logger

# 显示测试覆盖率
go test -cover ./...

# 生成HTML覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### 3. 代码规范检查
```bash
# 格式化代码
go fmt ./...

# 代码静态检查 (需安装golangci-lint)
golangci-lint run ./...

# 查找错误的构造
go vet ./...

# 查找竞态条件 (运行速度慢)
go test -race ./...
```

### 添加新功能的步骤

#### 示例：添加代理支持

1. **修改配置** (`internal/config/config.go`):
   ```go
   type Config struct {
       HTTP struct {
           // ... 现有字段
           Proxy string `json:"proxy"`
       }
   }
   ```

2. **更新HTTP客户端** (`internal/http/client.go`):
   ```go
   func NewHTTPClient(cfg *config.Config) *HTTPClient {
       // 使用配置中的代理设立HTTP client
   }
   ```

3. **添加测试** (`internal/http/client_test.go`):
   ```go
   func TestProxySupport(t *testing.T) {
       // 测试代理功能
   }
   ```

4. **更新CLI参数** (`cmd/m3u8-downloader/main.go`):
   ```go
   flag.StringVar(&proxy, "proxy", "", "HTTP代理地址")
   ```

5. **集成测试验证**:
   ```bash
   go test -v ./...
   ```

### 性能优化指南

#### 1. 并发优化
- 当前: 默认10个并发下载goroutine
- 优化点: 根据网络状况动态调整并发数

#### 2. 内存优化
- 当前: 片段缓存在内存
- 优化点: 实现基于文件的缓存队列

#### 3. CPU优化
- 当前: FFmpeg外部进程合并
- 优化点: 自定义TS合并库(无transcoding)

#### 性能测试方法
```bash
# 基准测试模板
go test -bench=. -benchmem ./...

# CPU性能分析
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# 内存性能分析
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### 常见开发任务

#### 添加新的错误类型
```go
// 在 internal/errors/errors.go 中
const (
    NewErrorCode ErrorCode = "NEW_ERROR_CODE"
)

// 使用
return errors.New(NewErrorCode, "错误描述", err)
```

#### 扩展日志功能
```go
// 在 internal/logger/logger.go 中
// 添加新方法到Logger接口和ConsoleLogger实现
func (l *ConsoleLogger) Trace(msg string) {
    // 实现
}
```

#### 添加配置选项
```go
// 1. 在 Config 结构体中添加字段
type Config struct {
    NewOption string
}

// 2. 在 DefaultConfig() 中设置默认值
func DefaultConfig() *Config {
    // ...
    NewOption: "default",
}

// 3. 在 Validate() 中添加验证规则
func (c *Config) Validate() error {
    if c.NewOption == "" {
        return errors.New(InvalidConfig, "NewOption 不能为空")
    }
}

// 4. 在 main.go 中添加CLI标志
flag.StringVar(&cfg.NewOption, "new-option", "", "新选项说明")
```

### 代码审查检查清单

提交代码前，请确保：

- [ ] 代码通过 `go fmt`
- [ ] 代码通过 `go test ./...` 所有测试
- [ ] 新功能有相应的单元测试
- [ ] 没有运行时错误 (`go vet ./...`)
- [ ] 没有竞态条件 (`go test -race ./...` for concurrent code)
- [ ] 遵循项目命名规范
- [ ] 添加了必要的注释文档
- [ ] 更新了相关文档 (README.md, 本文件等)
- [ ] 提交信息清晰描述了改动内容

### 调试技巧

#### 1. 使用 Delve 调试器
```bash
# 安装 dlv
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/m3u8-downloader/main.go

# 在dlv REPL中的常用命令
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print variable
(dlv) stack
```

#### 2. 添加调试日志
```go
// 配置日志级别为Debug
logger.Debug("调试信息", map[string]interface{}{
    "variable": value,
    "status": currentStatus,
})
```

#### 3. 性能瓶颈分析
```bash
# 以pprof服务方式运行应用（如果支持）
# 或设置环境变量CPUPROFILE
env CPUPROFILE=cpu.prof ./m3u8-downloader ...
go tool pprof cpu.prof
```

### 向项目贡献

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启Pull Request

## 相关文档

- [测试指南](./TEST.md) - 详细的测试策略和方法
- [重构报告](./REFACTORING_REPORT_V2.0.md) - v2.0 架构改进详解
- 主项目README - 快速开始和使用指南

## 常见问题

**Q: 如何增加下载速度？**
A: 调整 `internal/core/manager.go` 中的 `maxConcurrent` 参数，以及 `internal/config/config.go` 中的超时设置。

**Q: 如何支持新的加密方式？**
A: 在 `internal/m3u8/parser.go` 中扩展 `parseKey()` 方法，在 `internal/util/util.go` 中添加新的解密函数。

**Q: 如何添加代理支持？**
A: 参见上面的"添加新功能的步骤"示例。

**Q: 项目使用了哪些外部依赖？**
A: 查看 `go.mod` 文件，主要依赖为 grequests (HTTP客户端库)。

---

最后更新: 2024年
版本: v2.0
