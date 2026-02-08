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

# m3u8-downloader

M3U8 HLS 下载器 — 并发下载、自动解密与 FFmpeg 合并，带彩色终端输出与可复用主题。

快速亮点
- 支持并发下载与重试策略
- 自动处理 AES-128 加密的 TS 段
- 彩色终端日志（Catppuccin Mocha 主题）
- 支持 `m3u8#fragment` 格式自动提取保存名

## 快速开始

构建并运行：

```bash
go build -o m3u8-downloader cmd/m3u8-downloader/main.go
./m3u8-downloader "https://example.com/video.m3u8" -o my_video
```

支持把 `-o` 放在 URL 之前或之后；也支持在 URL 后追加 `#name` 片段自动作为输出名：

```bash
./m3u8-downloader "https://sut.com/.../video.m3u8#FC2-PPV-4048938"
```

## CLI 参数（主要）

- `-u` string : 指定 M3U8 URL（可选，通常使用位置参数）
- `-o` string : 输出文件名（不含后缀），若不指定会从 URL 或 `#fragment` 解析
- `-n` int    : 并发下載线程数（默认 24）
- `-ht` string: 主机类型 (v1/v2)
- `-sp` string: 保存目录（默认当前目录）
- `-s`        : 允许不安全 HTTPS（跳过证书验证）
- `-c` string : 自定义 Cookie
- `-r` bool   : 下载后自动清理 TS（默认 true）
- `-v`        : 显示版本

示例：

```bash
./m3u8-downloader "https://example.com/video.m3u8" -n 32 -o my_video -sp ~/Movies
```

## 项目结构（简要）

```
cmd/              # CLI入口
internal/         # 内部包（config, logger, http, m3u8, core, video, util, theme）
docs/             # 文档
build/            # 构建产物
```

## 主题与颜色

配色放在 `internal/theme`，供 `logger`、`core`、`video` 等模块复用以保持一致的终端风格。

## 文档与开发

- 详细使用与开发说明请参阅 `docs/DEVELOPMENT.md` 与 `docs/TEST.md`。
- `docs/QUICK_REFERENCE.md` 已简化为快速引用，请以 `README.md` 为主说明。

## 测试与构建

```bash
go test ./...
go build ./...
```

## 变更与发布

请参阅 `docs/RELEASE_NOTES.md`。

## 许可证

MIT — 详见 LICENSE
