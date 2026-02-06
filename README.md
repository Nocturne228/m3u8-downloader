# m3u8-downloader

[![Go Version](https://img.shields.io/badge/go-1.16+-blue)]()
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

M3U8播放列表下载器，支持HLS视频下载、自动合并和加密处理。

## 📋 前置需求

- **Go**: 1.16 或更高版本
- **FFmpeg**: 2.8 或更高版本（用于视频合并）
  - macOS: `brew install ffmpeg`
  - Ubuntu: `sudo apt-get install ffmpeg`
  - Windows: [FFmpeg官方下载](https://ffmpeg.org/download.html)

## 🚀 快速开始

### 安装

```bash
# 克隆项目
git clone https://github.com/your-username/m3u8-downloader.git
cd m3u8-downloader

# 构建可执行文件
go build -o ./build/m3u8-downloader cmd/m3u8-downloader/main.go
# make build    # Makefile 方式
# ./build.sh    # 编写好的构建脚本

# 验证安装
./build/m3u8-downloader -h
```



### 基本使用

最简单的用法 - 直接提供M3U8 URL：

```bash
./m3u8-downloader "https://example.com/playlist.m3u8"
```

完整的参数用法：

```bash
./m3u8-downloader [选项] <M3U8_URL>

选项:
  -u string           M3U8播放列表URL（可选，支持位置参数）
  -o string           输出文件路径（默认: 同URL名称）
  -t int              HTTP请求超时秒数（默认: 10）
  -c int              并发下载数（默认: 10）
  -r int              失败重试次数（默认: 3）
  -h, -help           显示帮助信息
  -v, -version        显示版本信息
```

### 使用示例

```bash
# 最简单的方式
./m3u8-downloader "https://example.com/video.m3u8"

# 指定输出文件名
./m3u8-downloader "https://example.com/video.m3u8" -o my_video.mp4

# 设置自定义参数
./m3u8-downloader "https://example.com/video.m3u8" \
  -o output.mp4 \
  -c 20 \              # 20个并发
  -t 15 \              # 超时15秒
  -r 5                 # 重试5次
```

## 📁 项目结构

```
m3u8-downloader/
├── cmd/                      # 应用程序入口点
│   └── m3u8-downloader/
│       └── main.go          # CLI主程序
├── internal/                # 内部包（不对外暴露）
│   ├── config/             # ⚙️  配置管理
│   ├── logger/             # 📝 日志系统
│   ├── errors/             # ⚠️  错误处理
│   ├── http/               # 🌐 HTTP客户端
│   ├── m3u8/               # 📋 M3U8解析
│   ├── core/               # 🔧 核心下载引擎
│   ├── video/              # 🎬 视频处理
│   └── util/               # 🛠️  工具函数
├── test/                    # 📊 测试文件
├── docs/                    # 📚 文档
│   ├── TEST.md             # 测试指南
│   ├── DEVELOPMENT.md      # 开发指南
│   └── REFACTORING_REPORT_V2.0.md
├── go.mod                  # 模块定义
├── go.sum                  # 依赖校验和
├── README.md              # 本文件
└── build-release.sh       # 发布构建脚本
```

## 🏗️ 架构设计

### 模块职责

| 模块 | 职责 | 关键功能 |
|------|------|--------|
| **config** | 配置管理 | 默认配置、参数验证 |
| **logger** | 日志系统 | 结构化日志、多级别输出 |
| **errors** | 错误处理 | 统一错误定义、错误链 |
| **http** | HTTP客户端 | 请求控制、指数退避重试 |
| **m3u8** | 播放列表处理 | M3U8解析、密钥提取 |
| **core** | 核心下载 | 任务管理、并发控制 |
| **video** | 视频处理 | FFmpeg合并、格式转换 |
| **util** | 工具函数 | 文件I/O、加密解密 |

### 工作流程

```
┌─────────────────────────────────────────────┐
│ 1. 解析命令行参数                            │
│    ↓                                         │
│ 2. 读取M3U8播放列表                         │
│    ↓                                         │
│ 3. 解析媒体清单（提取片段、密钥等）        │
│    ↓                                         │
│ 4. 并发下载所有TS片段                       │
│    ├─ 自动检测加密                          │
│    ├─ 失败自动重试                          │
│    └─ 实时显示进度                          │
│    ↓                                         │
│ 5. 使用FFmpeg合并TS片段为MP4                │
│    ↓                                         │
│ 6. 验证输出文件和清理临时文件               │
└─────────────────────────────────────────────┘
```

## 📝 日志输出

程序会输出关键操作信息：

```
✓ 成功解析M3U8: 1080个片段
⟳ 正在下载... [████████░░] 80% | 速度: 5.2 MB/s | 剩余: 45s
✓ 下载完毕: 850 MB
⟳ 正在合并...
✓ 视频合并完成: output.mp4 (850 MB)
```

## 🧪 测试

### 运行单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/logger -v

# 显示覆盖率
go test -cover ./...

# 生成HTML覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试包含内容

- ✅ 配置验证测试
- ✅ 错误处理测试
- ✅ 日志输出测试
- ✅ HTTP客户端测试

详细信息见 [TEST.md](docs/TEST.md)

## 🔨 开发

### 构建

```bash
# 标准构建
go build -o m3u8-downloader cmd/m3u8-downloader/main.go

# 带调试符号的构建
go build -gcflags="all=-N -l" -o m3u8-downloader cmd/m3u8-downloader/main.go

# 交叉编译（示例：为Linux构建）
GOOS=linux GOARCH=amd64 go build -o m3u8-downloader cmd/m3u8-downloader/main.go
```

### 代码检查

```bash
# 格式化代码
go fmt ./...

# 代码静态检查
go vet ./...

# 竞态条件检查（仅限有并发的代码）
go test -race ./...
```

### 发布构建

```bash
# 使用发布脚本
./build-release.sh

# 或手动：
go build -ldflags="-s -w" -o releases/m3u8-downloader cmd/m3u8-downloader/main.go
```

## 📚 详细文档

- **[测试指南](docs/TEST.md)** - 单元测试、集成测试、E2E测试详解
- **[开发指南](docs/DEVELOPMENT.md)** - 开发流程、贡献指南、扩展方向
- **[重构报告](docs/REFACTORING_REPORT_V2.0.md)** - v2.0架构设计和优化细节

## ⚙️ 配置说明

配置通过命令行参数传递（目前支持）。计划扩展支持YAML/JSON配置文件。

### 关键配置项

```go
HTTP:
  Timeout: 10s          # HTTP请求超时
  MaxRetries: 3         # 最大重试次数
  RetryDelay: 1s        # 重试延迟

Download:
  MaxConcurrent: 10     # 并发下载数
  ChunkSize: 1MB        # 每个片段大小
  OutDir: .             # 输出目录

FFmpeg:
  Path: ffmpeg          # FFmpeg可执行文件路径
  Timeout: 300s         # FFmpeg合并超时

Log:
  Level: INFO           # info/debug/warn/error
  Format: console       # 输出格式
```

## 🐛 常见问题

### Q: 下载速度慢
**A:** 尝试增加并发数 `-c 20` 或更大，同时检查网络连接。

### Q: 找不到ffmpeg
**A:** 确保FFmpeg已安装并在PATH中。验证：`ffmpeg -version`

### Q: 下载中断
**A:** 程序有自动重试机制，可用 `-r 5` 增加重试次数。

### Q: 内存占用过高
**A:** 较大文件的下载可能占用内存。可通过减少 `-c` 参数降低并发数。

### Q: 输出文件播放异常
**A:** 某些特殊的M3U8格式可能需要特殊处理。检查日志输出是否有错误提示。

## 📄 许可证

本项目采用MIT许可证 - 详见 [LICENSE](LICENSE) 文件

---

**项目版本**: v2.0  
**最后更新**: 2025年  
**Go版本**: 1.16+  
