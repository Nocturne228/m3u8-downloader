# M3U8 下载器 v2.1 - 交互优化说明

> 🎯 **优化聚焦**：精简输出、支持位置参数、提升用户体验

---

## 📝 主要改进

### 1. 简化命令行使用

#### 旧方式（v2.0）
```bash
./m3u8-downloader -u "http://example.com/video.m3u8" -n 32
```

#### 新方式（v2.1）✨
```bash
# 最简形式 - 直接传入 URL 作为第一个参数
./m3u8-downloader "http://example.com/video.m3u8"

# 仍然支持 -u 参数（可选）
./m3u8-downloader -u "http://example.com/video.m3u8" -n 32

# 混合使用
./m3u8-downloader "http://example.com/video.m3u8" -n 32 -o "my_video"
```

### 2. 精简终端输出

#### 改进前
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
  ...（多行冗余信息）

[Info] 创建临时目录: ./movie
待下载 ts 文件数量: 1000 个
...
[Info] 开始使用 FFmpeg 合并视频...
[Info] 已清除临时 ts 文件
```

#### 改进后（精简、干净）
```
[Preparing] 准备下载 1000 个 TS 文件...
[Downloading] ■■■■■ 50.00% (500/1000) 25.5 files/s ETA: 20s
[Merging] 合并视频中...

[Success] 视频已保存：/path/to/movie.mp4
下载耗时：45.1s
```

### 3. 智能参数处理

程序现在支持三种参数方式：

```bash
# 方式1：位置参数（最简单） ⭐ 推荐
./m3u8-downloader https://example.com/video.m3u8

# 方式2：使用 -u 参数
./m3u8-downloader -u https://example.com/video.m3u8

# 方式3：混合使用
./m3u8-downloader https://example.com/video.m3u8 -n 32 -o "movie"
```

---

## 📚 使用示例

### 基础使用

```bash
# 最简单的用法
./m3u8-downloader "https://example.com/video.m3u8"

# 查看帮助
./m3u8-downloader
```

### 高级用法

```bash
# 指定下载线程数和输出文件名
./m3u8-downloader "https://example.com/video.m3u8" -n 32 -o "my_video"

# 使用 v2 Host 分析模式
./m3u8-downloader "https://example.com/video.m3u8" -ht v2

# 添加自定义 Cookie
./m3u8-downloader "https://example.com/video.m3u8" -c "session=abc123"

# 保存到指定目录
./m3u8-downloader "https://example.com/video.m3u8" -sp "/home/user/videos"

# 不自动清除 TS 文件
./m3u8-downloader "https://example.com/video.m3u8" -r=false

# 完整示例
./m3u8-downloader "https://example.com/video.m3u8" \
  -n 32 \
  -o "downloaded_video" \
  -ht v1 \
  -c "key1=value1; key2=value2" \
  -sp "/home/user/downloads"
```

---

## 🔧 参数说明

```
位置参数:
  <URL>  m3u8 下载地址 (http(s)://url/.../index.m3u8)

命名参数:
  -u string
        m3u8 下载地址（可选，优先使用位置参数）
  
  -n int
        下载线程数（默认：24），范围 1-256
  
  -ht string
        Host 分析方式（默认：v1）
        v1: 保留路径目录部分
        v2: 仅使用 Host（某些服务器需要）
  
  -o string
        输出文件名，不含后缀（默认：movie）
  
  -c string
        自定义 Cookie（例如：key1=v1; key2=v2）
  
  -r bool
        是否自动清除下载的 TS 文件（默认：true）
  
  -s int
        允许不安全的 HTTPS 请求（默认：0/否）
  
  -sp string
        文件保存的绝对路径（默认：当前目录）
```

---

## 📊 输出说明

### 下载阶段

```
[Preparing] 准备下载 1000 个 TS 文件...
```
- 表示正在解析 M3U8 文件
- 统计需要下载的 TS 片段总数

### 进度显示

```
[Downloading] ■■■■■ 50.00% (500/1000) 25.5 files/s ETA: 20s
```

| 字段 | 说明 |
|------|------|
| `■■■■■` | 进度条（实时更新） |
| `50.00%` | 完成百分比 |
| `(500/1000)` | 已下载/总数 |
| `25.5 files/s` | 下载速度 |
| `ETA: 20s` | 预计剩余时间 |

### 合并阶段

```
[Merging] 合并视频中...
```
- 使用 FFmpeg 合并 TS 文件
- 生成标准 MP4 格式

### 完成

```
[Success] 视频已保存：/path/to/movie.mp4
下载耗时：45.1s
```
- 显示最终文件路径
- 显示总耗时

---

## 🛠️ 故障排除

### 问题：找不到 m3u8-downloader 命令

**解决方案：** 确保程序在当前目录或 PATH 中
```bash
# 方式1：使用相对路径
/path/to/m3u8-downloader "url"

# 方式2：添加到 PATH
export PATH=$PATH:/path/to/m3u8-downloader-directory
m3u8-downloader "url"

# 方式3：创建软链接（Linux/macOS）
sudo ln -s /path/to/m3u8-downloader /usr/local/bin/
m3u8-downloader "url"
```

### 问题：参数不被识别

**解决方案：** 确保 URL 有双引号
```bash
# ❌ 错误
./m3u8-downloader https://example.com/video.m3u8?key=value&foo=bar

# ✅ 正确
./m3u8-downloader "https://example.com/video.m3u8?key=value&foo=bar"
```

### 问题：线程数设置不生效

**解决方案：** 确保 `-n` 参数在正确位置
```bash
# ✅ 正确用法
./m3u8-downloader "url" -n 64
./m3u8-downloader -n 64 "url"
```

---

## 💡 最佳实践

### 1. 合理选择线程数

```bash
# 少于 50 个文件
./m3u8-downloader "url" -n 8

# 50-200 个文件
./m3u8-downloader "url" -n 16

# 200-1000 个文件
./m3u8-downloader "url" -n 32

# 1000+ 个文件
./m3u8-downloader "url" -n 48

# 特别大的文件
./m3u8-downloader "url" -n 64
```

### 2. 处理特殊 URL

```bash
# URL 包含特殊字符
./m3u8-downloader "https://example.com/video?token=abc&id=123"

# URL 中有空格（使用引号）
./m3u8-downloader "https://example.com/my video.m3u8"

# 本地文件（使用完整路径）
./m3u8-downloader "http://localhost:8080/video.m3u8"
```

### 3. 下载失败处理

```bash
# 第一次使用 v1 模式（默认）
./m3u8-downloader "url"

# 失败后尝试 v2 模式
./m3u8-downloader "url" -ht v2

# 如果还是不行，减少线程数
./m3u8-downloader "url" -ht v2 -n 8
```

---

## 🎯 性能建议

| 场景 | 推荐设置 |
|------|---------|
| 快速网络 | `-n 48-64` |
| 正常网络 | `-n 24-32` |
| 慢速网络 | `-n 8-16` |
| 不稳定网络 | `-n 4-8` |
| 服务器限制请求 | `-n 2-4` |

---

## 📌 版本信息

| 版本 | 日期 | 主要改进 |
|------|------|---------|
| v2.1 | 2026-02-06 | 支持位置参数、精简输出 |
| v2.0 | 2026-02-06 | FFmpeg 合并、进度显示 |
| v1.0 | - | 基础下载功能 |

---

## ✅ 快速检查清单

下载前检查：
- [ ] FFmpeg 已安装 (`ffmpeg -version`)
- [ ] URL 有效且可访问
- [ ] 输出目录可写
- [ ] 磁盘空间充足（预留 1.5x 视频大小）
- [ ] 网络连接稳定

下载中：
- [ ] 进度条实时更新
- [ ] 显示下载速度
- [ ] ETA 在合理范围内

下载后：
- [ ] MP4 文件已生成
- [ ] 文件大小合理
- [ ] 播放器可正常播放

---

## 📞 常见问题

**Q: 可以同时下载多个视频吗？**
A: 可以，但需要在不同的终端窗口中运行，或使用不同的输出文件名。

**Q: 如何恢复中断的下载？**
A: 重新运行相同的命令，程序会跳过已下载的 TS 文件。

**Q: 如何获取更详细的日志？**
A: 当前版本不支持详细日志，但可以通过降低线程数来减少并发问题。

**Q: 支持 m3u8 嵌套吗？**
A: 不支持，需要手动转换为单层 m3u8。

---

**祝你使用愉快！** 🎉
