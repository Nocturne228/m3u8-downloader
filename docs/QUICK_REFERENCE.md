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
# 快速参考 (Quick Reference)

本文件是 `README.md` 的快速摘要，详细信息请参阅主说明文件。

常用命令

```bash
# 简单下载（支持 -o 在 URL 前或后）
./m3u8-downloader "https://example.com/video.m3u8"

# 指定输出和并发数
./m3u8-downloader "https://example.com/video.m3u8" -o my_video -n 32

# 使用 fragment 自动取名
./m3u8-downloader "https://sut.com/.../video.m3u8#FC2-PPV-4048938"
```

开发与测试

```bash
# 构建
go build -o m3u8-downloader cmd/m3u8-downloader/main.go

# 运行测试
go test ./...

# 格式化与静态检查
go fmt ./...
go vet ./...
```

文档导航

- 主文档：`README.md`
- 开发指南：`docs/DEVELOPMENT.md`
- 测试指南：`docs/TEST.md`
- 发布说明：`docs/RELEASE_NOTES.md`

（本文件仅作速查）
│   ├── http/                   # HTTP客户端
