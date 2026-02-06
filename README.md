# m3u8-downloader

> 🚀 **v2.0 版本重大升级：使用 FFmpeg 进行专业级视频合并，提升用户体验**

golang 多线程下载直播流m3u8格式的视屏，跨平台。 你只需指定必要的 flag (`u`、`o`、`n`、`ht`) 来运行, 工具就会自动帮你解析 M3U8 文件，并将 TS 片段下载下来合并成一个文件。

## ✨ v2.0 新特性

- 🎬 **使用 FFmpeg 进行视频合并** - 生成标准 MP4 格式，而不是原始 TS 流
- 📊 **改进的进度显示** - 显示实时下载速度和 ETA（剩余时间）
- 🔄 **智能重试机制** - 指数退避算法，更好的容错能力
- ⚡ **并发安全改进** - 使用原子操作，确保数据准确性
- 🛡️ **增强的错误处理** - 完整的参数验证和日志记录

详见 [优化文档](./IMPROVEMENTS.md) 和 [代码对比](./CODE_COMPARISON.md)

## 功能介绍

1. ✅ 下载和解析 M3U8
2. ✅ 下载 TS 失败重试（加密的同步解密）
3. ✅ **使用 FFmpeg 合并 TS 片段为 MP4**
4. ✅ 显示下载速度和 ETA
5. ✅ 完整的错误恢复机制

> 可以下载岛国小电影  
> 可以下载岛国小电影  
> 可以下载岛国小电影    
> 重要的事情说三遍......

## 效果展示
![demo](./demo.gif)

---

## 📋 参数说明：

```
-u   m3u8下载地址(http(s)://url/xx/xx/index.m3u8) [必需]
-o   输出文件名（默认：movie）
-n   下载线程数（默认：24，建议：16-32）
-ht  Host类型（v1 或 v2，默认：v1）
     v1: http(s):// + url.Host + filepath.Dir(url.Path)
     v2: http(s):// + url.Host
-c   自定义请求Cookie（例如：key1=v1; key2=v2）
-r   是否自动清除ts文件（默认：true）
-s   是否允许不安全的请求（默认：0）
-sp  文件保存的绝对路径（默认：当前路径）
```

> **⚠️ 重要前置条件：需要系统安装 FFmpeg**
> 
> - **macOS**: `brew install ffmpeg`
> - **Ubuntu/Debian**: `sudo apt-get install ffmpeg`
> - **CentOS/RHEL**: `sudo yum install ffmpeg`
> - **Windows**: https://ffmpeg.org/download.html

默认情况只需要传 `-u` 参数，其他参数保持默认即可。 部分链接可能限制请求频率，可根据实际情况调整 `-n` 参数的值。

## 下载

已经编译好的平台有： [点击下载](https://github.com/llychao/m3u8-downloader/releases)

- m3u8-darwin-amd64
- m3u8-darwin-arm64
- m3u8-linux-386
- m3u8-linux-amd64
- m3u8-linux-arm64
- m3u8-windows-386.exe
- m3u8-windows-amd64.exe
- m3u8-windows-arm64.exe

## 用法

### 编译方式（推荐）

```bash
自己编译：go build -o m3u8-downloader
简洁使用：./m3u8-downloader  -u=http://example.com/index.m3u8
完整使用：./m3u8-downloader  -u=http://example.com/index.m3u8 -o=example -n=16 -ht=v1 -c="key1=v1; key2=v2"
```

### 二进制方式:

Linux 和 MacOS 和 Windows PowerShell

```
简洁使用：
./m3u8-linux-amd64 -u=http://example.com/index.m3u8
./m3u8-darwin-amd64 -u=http://example.com/index.m3u8 
.\m3u8-windows-amd64.exe -u=http://example.com/index.m3u8

完整使用：
./m3u8-linux-amd64 -u=http://example.com/index.m3u8 -o=example -n=16 -ht=v1 -c="key1=v1; key2=v2"
./m3u8-darwin-amd64 -u=http://example.com/index.m3u8 -o=example -n=16 -ht=v1 -c="key1=v1; key2=v2"
.\m3u8-windows-amd64.exe -u=http://example.com/index.m3u8 -o=example -n=16 -ht=v1 -c="key1=v1; key2=v2"
```

## 问题说明

1.在Linux或者mac平台，如果显示无运行权限，请用chmod 命令进行添加权限
```bash
 # Linux amd64平台
 chmod 0755 m3u8-linux-amd64
 # Mac darwin amd64平台
 chmod 0755 m3u8-darwin-amd64
 ```
2.下载失败的情况,请设置 -ht="v1" 或者 -ht="v2" （默认为 v1）
```golang
func get_host(Url string, ht string) string {
    u, err := url.Parse(Url)
    var host string
    checkErr(err)
    switch ht {
    case "v1":
        host = u.Scheme + "://" + u.Host + path.Dir(u.Path)
    case "v2":
        host = u.Scheme + "://" + u.Host
    }
    return host
}
```
