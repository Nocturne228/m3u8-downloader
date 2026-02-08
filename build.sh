#!/bin/bash

# M3U8 下载器编译和使用脚本
# 依赖: Go 1.16+ 和 FFmpeg

echo "=========================================="
echo "  M3U8 下载器 - 编译和使用说明"
echo "=========================================="
echo ""

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.16+"
    exit 1
fi

# 检查FFmpeg是否安装
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ FFmpeg 未安装"
    echo ""
    echo "请根据你的操作系统安装 FFmpeg："
    echo "macOS:        brew install ffmpeg"
    echo "Ubuntu:       sudo apt-get install ffmpeg"
    echo "CentOS:       sudo yum install ffmpeg"
    echo "Windows:      https://ffmpeg.org/download.html"
    echo ""
    exit 1
fi

echo "✅ Go 已安装: $(go version)"
echo "✅ FFmpeg 已安装: $(ffmpeg -version 2>&1 | head -n 1)"
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 下载依赖
echo "📦 下载依赖..."
go mod download

# 编译
echo "🔨 编译中..."
if [ "$(uname)" == "Darwin" ]; then
    # macOS
    go build -o m3u8-downloader cmd/m3u8-downloader/main.go
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    # Linux
    go build -o m3u8-downloader cmd/m3u8-downloader/main.go
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ] || [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
    # Windows
    go build -o m3u8-downloader.exe cmd/m3u8-downloader/main.go
else
    go build -o m3u8-downloader cmd/m3u8-downloader/main.go
fi

if [ $? -eq 0 ]; then
    echo "✅ 编译成功！"
    echo ""
    echo "📖 使用示例："
    echo ""
echo "基本使用:"
echo "  ./m3u8-downloader 'https://example.com/index.m3u8'"
echo ""
echo "指定输出文件名:"
echo "  ./m3u8-downloader 'https://example.com/index.m3u8' -o 'my_video.mp4'"
echo ""
echo "自定义参数:"
echo "  ./m3u8-downloader 'https://example.com/index.m3u8' -c 20 -t 15 -r 5"
echo "  （-c: 并发数，-t: 超时秒数，-r: 重试次数）"
    echo ""
    echo "完整参数说明，请执行:"
    echo "  ./m3u8-downloader -h"
    echo ""
    echo "🚀 现在你可以开始下载视频了！"
else
    echo "❌ 编译失败！"
    exit 1
fi
