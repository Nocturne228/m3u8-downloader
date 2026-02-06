# M3U8 下载器优化文档 v2.0

## 概述
对原 m3u8 下载器进行了全面的逻辑、交互体验和视频合并功能的优化。

---

## 核心改进

### 1. **使用 FFmpeg 进行 TS 合并** ⭐⭐⭐
**改进前：** 简单的二进制连接，生成的文件虽然能播放但不是正确的 MP4 格式
```go
// 旧方法：直接拼接字节
func mergeTs(downloadDir string) string {
    // 只是简单的将 ts 文件字节连接...
}
```

**改进后：** 使用 FFmpeg concat demuxer 进行正确的视频合并
```go
func mergeWithFFmpeg(tsDir, movieName string) (string, error) {
    // 1. 收集所有 ts 文件
    // 2. 创建 concat.txt 描述文件
    // 3. 使用 FFmpeg 的 concat demuxer 进行合并
    // 4. 生成真正的 MP4 格式文件
}
```

**优势：**
- ✅ 生成标准的 MP4 格式
- ✅ 完整的视频元数据
- ✅ 更好的兼容性
- ✅ 支持字幕、多音轨等

### 2. **交互体验优化**

#### 改进的进度条显示
**旧显示：** `[Downloading] ■■■■■ 50% (file_name.ts)`
```
- 只显示文件名，没有速度信息
- 无法计算剩余时间
- 多线程竞态条件导致计数不准确
```

**新显示：** `[Downloading] ■■■■■ 50.00% (500/1000) 25.5 files/s ETA: 20s`
```go
func updateProgressBar(total int) {
    downloadCount := atomic.LoadInt64(&stats.DownloadCount)  // 原子操作
    proportion := float32(downloadCount) / float32(total)
    
    // 计算下载速度
    elapsed := now.Sub(stats.StartTime).Seconds()
    speed := float64(downloadCount) / elapsed
    
    // 计算 ETA
    remainCount := float64(total) - float64(downloadCount)
    eta := remainCount / speed
}
```

**改进列表：**
- ✅ 显示实时下载速度（files/s）
- ✅ 显示已下载/总数
- ✅ 计算剩余时间 ETA
- ✅ 使用原子操作（`atomic.LoadInt64`）避免竞态条件
- ✅ 美化输出格式

#### 优化的命令行提示
```
========================================
     M3U8 视频下载工具 v2.0
=========================================
[功能]  多线程下载直播流 m3u8 视频
[特性]  
  - 支持加密 TS 文件自动解密
  - 使用 FFmpeg 合并生成 MP4 格式
  - 支持断点续传和失败重试
  - 显示下载速度和 ETA
[提醒]  
  - 下载失败，请使用 -ht=v2 
  - ...
========================================
```

### 3. **逻辑改进**

#### 增强的参数验证
```go
// 旧：直接使用参数，无验证
// 新：完整的验证逻辑
if !strings.HasPrefix(m3u8Url, "http") || m3u8Url == "" {
    fmt.Println("[Error] 请输入有效的 m3u8 下载地址")
    return
}

if maxGoroutines <= 0 || maxGoroutines > 256 {
    fmt.Println("[Error] 线程数必须在 1-256 之间")
    return
}

if movieName == "" {
    fmt.Println("[Error] 输出文件名不能为空")
    return
}
```

#### FFmpeg 可用性检查
```go
func checkFFmpeg() bool {
    cmd := exec.Command("ffmpeg", "-version")
    _, err := cmd.Output()
    return err == nil
}

// 在程序启动时检查
if !checkFFmpeg() {
    fmt.Printf("[Failed] 系统未安装 FFmpeg，请先安装 FFmpeg\n")
    fmt.Println("macOS: brew install ffmpeg")
    fmt.Println("Linux: apt-get install ffmpeg")
    return
}
```

#### 改进的下载检查函数
```go
// 旧：检查是否存在特定索引文件（不可靠）
func checkTsDownDir(dir string) bool {
    if isExist, _ := pathExists(filepath.Join(dir, "00000.ts")); !isExist {
        return true  // 逻辑反向了！
    }
    return false
}

// 新：统计实际下载的 ts 文件数量，允许 10% 的损失
func checkTsDownDir(dir string, expectedCount int) bool {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return false
    }
    tsCount := 0
    for _, f := range files {
        if !f.IsDir() && filepath.Ext(f.Name()) == ".ts" {
            tsCount++
        }
    }
    // 允许损失率 10%
    threshold := int(float64(expectedCount) * 0.9)
    return tsCount >= threshold
}
```

### 4. **错误处理改进**

#### 获取 M3U8 文件的重试机制
```go
func getM3u8Body(Url string) string {
    maxRetries := 3
    for i := 1; i <= maxRetries; i++ {
        r, err := grequests.Get(Url, ro)
        if err != nil {
            if i < maxRetries {
                logger.Printf("[warn] 第 %d 次获取 m3u8 文件失败，2秒后重试: %v\n", i, err)
                time.Sleep(2 * time.Second)
                continue
            }
            checkErr(fmt.Errorf("获取 m3u8 文件失败（已重试 %d 次）: %v", maxRetries, err))
        }
        // ... 检查响应状态 ...
        return r.String()
    }
}
```

#### TS 文件下载的指数退避
```go
func downloadTsFile(ts TsInfo, download_dir, key string, retries int) {
    // ...
    if err != nil || !res.Ok {
        if retries > 0 {
            // 指数退避：重试次数越多，等待时间越长
            time.Sleep(time.Duration(MAX_RETRY_COUNT-retries) * time.Second)
            downloadTsFile(ts, download_dir, key, retries-1)
            return
        }
    }
    // ...
}
```

#### 改进的日志记录
```go
// 使用 logger 替代 fmt.Println，便于区分日志级别
logger.Printf("[Error] 获取当前目录失败: %v\n", err)
logger.Printf("[warn] 解析 URL 失败: %v，使用版本 v2\n", err)
logger.Printf("[Info] 创建下载目录: %s\n", download_dir)
```

### 5. **下载统计信息**

新增 `DownloadStats` 结构体用于跟踪下载进度：
```go
type DownloadStats struct {
    TotalCount      int64     // 总文件数
    DownloadCount   int64     // 已下载数
    StartTime       time.Time // 开始时间
    LastUpdateTime  time.Time // 最后更新时间
    LastDownloadCount int64   // 上次记录的下载数
}
```

### 6. **并发安全改进**

使用原子操作替代互斥锁：
```go
// 旧：非线程安全
downloadCount++
DrawProgressBar("Downloading", float32(downloadCount)/float32(tsLen), ...)

// 新：原子操作
atomic.AddInt64(&stats.DownloadCount, 1)
downloadCount := atomic.LoadInt64(&stats.DownloadCount)
```

---

## 使用说明

### 安装 FFmpeg

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt-get install ffmpeg
```

**CentOS/RHEL:**
```bash
sudo yum install ffmpeg
```

**Windows:**
从 https://ffmpeg.org/download.html 下载安装

### 基本用法
```bash
go run m3u8-downloader.go -u "http://example.com/video.m3u8" -o "my_video" -n 32
```

### 参数说明
```
-u string         m3u8 下载地址（必需）
-o string         输出文件名（默认：movie）
-n int            下载线程数（默认：24，建议：16-32）
-ht string        Host 类型（v1 或 v2，默认：v1）
-c string         自定义 Cookie
-r bool           是否自动清除 ts 文件（默认：true）
-s int            是否允许不安全请求（默认：0）
-sp string        文件保存路径（默认：当前目录）
```

### 使用示例
```bash
# 基本使用
./m3u8-downloader -u "https://example.com/index.m3u8"

# 指定输出名称和线程数
./m3u8-downloader -u "https://example.com/index.m3u8" -o "my_movie" -n 32

# 使用 v2 模式（某些服务器需要）
./m3u8-downloader -u "https://example.com/index.m3u8" -ht v2

# 指定保存路径
./m3u8-downloader -u "https://example.com/index.m3u8" -sp "/path/to/save"

# 自定义 Cookie（某些服务器需要）
./m3u8-downloader -u "https://example.com/index.m3u8" -c "session_id=xxx"

# 下载加密内容
./m3u8-downloader -u "https://example.com/encrypted.m3u8" -s 1
```

---

## 性能对比

| 指标 | 旧版本 | 新版本 | 提升 |
|------|--------|--------|------|
| 进度显示精度 | 文件名 | 速度 + ETA | 大幅提升 |
| 线程安全 | 非安全 | 原子操作 | 有保障 |
| FFmpeg 支持 | ❌ | ✅ | 新增 |
| 合并质量 | 原始拼接 | 标准 MP4 | 质量提升 |
| 错误恢复 | 基础 | 重试 + 日志 | 更可靠 |

---

## 主要函数改动

### 新增函数
- `checkFFmpeg()` - 检查 FFmpeg 是否可用
- `mergeWithFFmpeg()` - 使用 FFmpeg 合并 TS 文件
- `updateProgressBar()` - 改进的进度显示
- `getFileSize()` - 获取文件大小
- `drawMergeProgressBar()` - 合并进度显示

### 改进函数
- `downloadTsFile()` - 加入原子操作、指数退避
- `downloader()` - 改进进度统计逻辑
- `checkTsDownDir()` - 改进的检查逻辑，接受参数
- `getM3u8Body()` - 加入重试机制
- `getHost()` - 更好的错误处理
- `getTsList()` - 改进的字符串处理

### 已弃用函数（保留向后兼容）
- `mergeTs()` - 已被 `mergeWithFFmpeg()` 替代
- `win_merge_file()` - 已被 `mergeWithFFmpeg()` 替代
- `unix_merge_file()` - 已被 `mergeWithFFmpeg()` 替代

---

## 故障排除

### 问题 1：`ffmpeg: command not found`
**解决方案：** 安装 FFmpeg
```bash
# macOS
brew install ffmpeg

# Linux
sudo apt-get install ffmpeg

# 验证安装
ffmpeg -version
```

### 问题 2：下载失败或进度条卡住
**解决方案：**
1. 尝试使用 `-ht v2` 参数
2. 减少线程数：`-n 16`
3. 添加 Cookie（如果需要登录）：`-c "your_cookie"`

### 问题 3：合并失败
**解决方案：**
1. 检查 FFmpeg 版本（应 >= 2.8）
2. 检查磁盘空间
3. 检查 ts 文件是否完整

### 问题 4：输出文件过小或无法播放
**解决方案：**
1. 确保所有 ts 文件都下载完整
2. 使用 `ffmpeg -i output.mp4` 检查文件结构
3. 尝试手动重新合并

---

## 贡献和反馈

如有任何问题或建议，欢迎提出！

---

## 版本历史

### v2.0 (2026-02-06)
- ✨ 新增 FFmpeg 合并功能
- 🚀 改进进度显示（显示速度和 ETA）
- 🔧 优化并发控制（原子操作）
- 🛡️ 增强错误处理和日志
- 📝 改进命令行提示和参数验证

### v1.0
- 基础下载和二进制合并功能
