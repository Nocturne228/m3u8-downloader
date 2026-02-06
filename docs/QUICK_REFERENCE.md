# 快速参考 (Quick Reference)

本文件提供m3u8-downloader项目的快速参考指南。

## 🚀 快速命令

### 基本使用

```bash
# 简单下载
./m3u8-downloader "https://example.com/video.m3u8"

# 指定输出文件
./m3u8-downloader "https://example.com/video.m3u8" -o my_video.mp4

# 设置并发数
./m3u8-downloader "https://example.com/video.m3u8" -c 20

# 完整参数
./m3u8-downloader "URL" -o output.mp4 -c 10 -t 15 -r 5
```

---

## 🔨 开发命令

### 构建

```bash
# 编译可执行文件
go build -o m3u8-downloader cmd/m3u8-downloader/main.go

# 带调试信息的构建
go build -gcflags="all=-N -l" -o m3u8-downloader cmd/m3u8-downloader/main.go

# 发布构建（优化大小）
go build -ldflags="-s -w" -o m3u8-downloader cmd/m3u8-downloader/main.go
```

### 测试

```bash
# 运行所有测试
go test ./...

# 特定包的测试
go test ./internal/config -v
go test ./internal/logger -v
go test ./internal/errors -v

# 显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 竞态条件检查
go test -race ./...
```

### 代码检查

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 代码风格检查（需安装golangci-lint）
golangci-lint run ./...
```

---

## 📂 项目结构

```
m3u8-downloader/
├── cmd/m3u8-downloader/        # CLI入口
├── internal/                    # 内部包
│   ├── config/                 # 配置管理
│   ├── logger/                 # 日志系统
│   ├── errors/                 # 错误处理
│   ├── http/                   # HTTP客户端
│   ├── m3u8/                   # M3U8解析
│   ├── core/                   # 核心引擎
│   ├── video/                  # 视频处理
│   └── util/                   # 工具函数
├── test/                       # 测试
├── docs/                       # 文档
└── README.md, go.mod, etc.    # 配置文件
```

---

## 📚 文档导航

| 文档 | 用途 | 行数 |
|------|------|------|
| [README.md](/README.md) | 项目概览和基本使用 | 200+ |
| [docs/TEST.md](/docs/TEST.md) | 详细的测试指南 | 400+ |
| [docs/DEVELOPMENT.md](/docs/DEVELOPMENT.md) | 开发和扩展指南 | 600+ |
| [docs/CONTRIBUTING.md](/docs/CONTRIBUTING.md) | 贡献规范和流程 | 400+ |
| [docs/CHANGELOG.md](/docs/CHANGELOG.md) | 版本历史和更新 | 300+ |
| [docs/REFACTORING_REPORT_V2.0.md](/docs/REFACTORING_REPORT_V2.0.md) | 架构改进详解 | 600+ |
| [docs/PROJECT_FILES.md](/docs/PROJECT_FILES.md) | 文件清单和说明 | 400+ |

---

## 🎯 常见任务

### 添加新功能

1. **修改配置**（如果需要）
   ```go
   // internal/config/config.go
   type Config struct {
       NewOption string
   }
   ```

2. **编写功能代码**
   ```go
   // internal/[package]/file.go
   // 实现相关功能
   ```

3. **添加测试**
   ```go
   // internal/[package]/*_test.go
   func TestNewFeature(t *testing.T) {
       // 测试代码
   }
   ```

4. **更新文档**
   - 更新 README.md 或相关文档

5. **提交代码**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   git push
   ```

### 修复Bug

1. 创建bugfix分支
   ```bash
   git checkout -b fix/issue-description
   ```

2. 创建测试用例再现bug
   ```go
   func TestBugFix(t *testing.T) {
       // 验证bug和修复
   }
   ```

3. 修复代码

4. 验证所有测试通过
   ```bash
   go test ./...
   ```

5. 提交PR

### 性能优化

参考 [DEVELOPMENT.md](/docs/DEVELOPMENT.md) 中的"性能优化指南"部分

```bash
# 基准测试
go test -bench=. -benchmem ./...

# CPU分析
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# 内存分析
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

---

## 🐛 常见问题解决

### 编译错误

```bash
# 清理缓存
go clean -cache
go clean -modcache

# 重新获取依赖
go mod tidy
go get -u ./...

# 重新编译
go build -o m3u8-downloader cmd/m3u8-downloader/main.go
```

### 测试失败

```bash
# 运行并显示详细信息
go test -v ./...

# 运行特定测试
go test -run TestName -v ./package

# 显示输出
go test -v -count=1 ./...
```

### 导入问题

```bash
# 格式化导入
goimports -w .

# 或手动：
go fmt ./...
```

---

## ✅ 验证清单

在提交PR前，确保：

- [ ] `go fmt ./...` 通过
- [ ] `go vet ./...` 无错误
- [ ] `go test ./...` 全部通过
- [ ] 代码添加了注释说明
- [ ] 相关文档已更新
- [ ] 提交信息符合约定
- [ ] 无未提交的临时文件

---

## 🔗 关键链接

- **GitHub**: [m3u8-downloader](https://github.com/your-username/m3u8-downloader)
- **Issues**: [问题跟踪](https://github.com/your-username/m3u8-downloader/issues)
- **Discussions**: [讨论区](https://github.com/your-username/m3u8-downloader/discussions)

---

## 📞 获取帮助

### 问题解答顺序

1. 查看 [README.md](/README.md) 的FAQ部分
2. 查看 [docs/DEVELOPMENT.md](/docs/DEVELOPMENT.md) 中的常见问题
3. 搜索 [现有Issue](https://github.com/your-username/m3u8-downloader/issues)
4. 在 [讨论区](https://github.com/your-username/m3u8-downloader/discussions) 提问
5. 提交新 [Issue](https://github.com/your-username/m3u8-downloader/issues/new)

### 报告问题时包含

- Go版本: `go version`
- FFmpeg版本: `ffmpeg -version`
- 操作系统
- 完整的错误信息
- 复现步骤

---

## 🚀 发布流程

```bash
# 1. 更新版本
# internal/config/config.go 中更新版本号（如果需要）

# 2. 更新CHANGELOG
# docs/CHANGELOG.md 中添加[X.X.X]部分

# 3. 编译发布版本
./build-release.sh

# 4. 标记Git标签
git tag -a vX.X.X -m "Release version X.X.X"
git push origin vX.X.X

# 5. 创建GitHub Release
# 在GitHub页面上传编译的二进制文件
```

---

## 💡 有用的Go命令

| 命令 | 说明 |
|------|------|
| `go run ./cmd/m3u8-downloader` | 直接运行（不编译） |
| `go get -u ./...` | 更新所有依赖 |
| `go mod tidy` | 清理未使用的依赖 |
| `go env` | 显示Go环境变量 |
| `go list ./...` | 列出所有包 |
| `go doc [package]` | 查看包文档 |
| `go test -cover ./...` | 覆盖率统计 |
| `go fmt -l .` | 检查需要格式化的文件 |

---

## 🎓 学习资源

- [Go官方文档](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [FFmpeg文档](https://ffmpeg.org/documentation.html)

---

最后更新: 2024年  
版本: v2.0  
快速参考版 1.0
