package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"m3u8-downloader/internal/config"
	"m3u8-downloader/internal/core"
	"m3u8-downloader/internal/logger"
)

var (
	// 命令行参数
	urlFlag     = flag.String("u", "", "M3U8 下载地址")
	nFlag       = flag.Int("n", 24, "下载线程数 (1-256, 默认 24)")
	htFlag      = flag.String("ht", "v1", "主机类型 (v1: 带路径, v2: 仅主机)")
	oFlag       = flag.String("o", "movie", "输出流名 (不带后缀)")
	cFlag       = flag.String("c", "", "自定义请求 Cookie")
	sFlag       = flag.Bool("s", false, "允许不安全的 HTTPS 请求")
	spFlag      = flag.String("sp", "", "文件保存的绝对路径")
	rFlag       = flag.Bool("r", true, "下载完成后自动清除 TS 文件")
	helpFlag    = flag.Bool("help", false, "显示帮助信息")
	versionFlag = flag.Bool("v", false, "显示版本信息")
)

const Version = "2.0.0"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 解析参数
	flag.Parse()

	// 显示帮助或版本
	if *helpFlag {
		showHelp()
		return
	}

	if *versionFlag {
		fmt.Printf("m3u8-downloader v%s\n", Version)
		return
	}

	// 获取 M3U8 URL
	m3u8URL := *urlFlag
	if m3u8URL == "" && len(flag.Args()) > 0 {
		m3u8URL = flag.Args()[0]
	}

	// 验证 URL
	if m3u8URL == "" {
		fmt.Fprintf(os.Stderr, "错误: 请提供 M3U8 URL\n")
		fmt.Fprintf(os.Stderr, "用法: m3u8-downloader <url> [选项]\n")
		fmt.Fprintf(os.Stderr, "    或: m3u8-downloader -u <url> [选项]\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(m3u8URL, "http") {
		fmt.Fprintf(os.Stderr, "错误: URL 必须以 http 或 https 开头\n")
		os.Exit(1)
	}

	// 验证线程数
	if *nFlag <= 0 || *nFlag > 256 {
		fmt.Fprintf(os.Stderr, "错误: 线程数必须在 1-256 之间\n")
		os.Exit(1)
	}

	// 验证输出文件名
	if *oFlag == "" {
		fmt.Fprintf(os.Stderr, "错误: 输出文件名不能为空\n")
		os.Exit(1)
	}

	// 创建配置
	cfg := config.DefaultConfig()
	cfg.Download.MaxGoroutines = *nFlag
	cfg.Download.HostType = *htFlag
	cfg.Download.SavePath = *spFlag
	cfg.Download.AutoClear = *rFlag
	cfg.Download.InsecureSkipVerify = *sFlag
	cfg.Download.Cookie = *cFlag

	// 创建日志记录器
	log := logger.New(cfg.Log.Level)

	// 创建应用程序
	app, err := core.NewApplication(cfg, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 初始化应用程序失败: %v\n", err)
		os.Exit(1)
	}

	// 运行应用程序
	err = app.Run(m3u8URL, *oFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	help := `M3U8 下载器 v%s - 高性能视频流下载工具

用法:
  m3u8-downloader <url> [选项]
  m3u8-downloader -u <url> [选项]

参数:
  <url>                    M3U8 下载地址 (http(s)://...)

选项:
  -u string               M3U8 下载地址 (可选，推荐使用位置参数)
  -n int                  下载线程数 (默认 24，范围 1-256)
  -ht string              主机类型 (默认 v1)
                          v1: http(s)://host + 目录路径
                          v2: http(s)://host
  -o string               输出文件名，不包括后缀 (默认 movie)
  -c string               自定义 HTTP Cookie
  -s                      允许不安全的 HTTPS 请求 (默认 false)
  -sp string              文件保存路径，绝对路径 (默认当前目录)
  -r                      下载完成后自动清除 TS 文件 (默认 true)
  -help                   显示帮助信息
  -v                      显示版本信息

示例:
  # 基本用法 - 推荐
  m3u8-downloader "https://example.com/video.m3u8"

  # 指定线程数和输出文件名
  m3u8-downloader "https://example.com/video.m3u8" -n 32 -o "my_video"

  # 指定保存路径
  m3u8-downloader "https://example.com/video.m3u8" -sp "/Users/user/Videos"

  # 使用自定义 Cookie
  m3u8-downloader "https://example.com/video.m3u8" -c "session=abc123"

  # 允许不安全的 HTTPS 连接
  m3u8-downloader "https://example.com/video.m3u8" -s

`
	fmt.Printf(help, Version)
}
