# 贡献指南 (Contributing)

感谢您对 m3u8-downloader 项目的兴趣！本文档提供了贡献项目的指导。

## 📋 目录

1. [行为准则](#行为准则)
2. [我想贡献](#我想贡献)
3. [报告错误](#报告错误)
4. [建议功能](#建议功能)
5. [拉取请求流程](#拉取请求流程)
6. [代码风格](#代码风格)
7. [提交信息规范](#提交信息规范)
8. [文档贡献](#文档贡献)

---

## 行为准则

### 我们的承诺

为了促进一个开放且热烈的社区环境，我们作为贡献者和维护者承诺：

- 使用包容和尊重的语言
- 欢迎不同的观点和经验
- 优雅接受建设性批评
- 专注于对社区最好的事情
- 向其他社区成员示范

### 令人无法接受的行为

令人无法接受的行为包括：

- 使用有性别含义的语言或意象
- 个人攻击或責低性评论
- 公开或私人骚扰
- 发布他人的隐私信息

维护者有权删除、编辑或拒绝任何不符合本行为准则的评论、提交、代码、wiki编辑或问题。

---

## 我想贡献

### 不确定从哪里开始？

从以下开始：

1. **查看开放的 Issue** - 寻找标记为 `good first issue` 或 `help wanted` 的问题
2. **阅读文档** - 确保充分理解项目架构和功能
3. **运行测试** - 验证开发环境设置正确
4. **选择任务** - 从小的、明确的贡献开始

### 贡献类型

我们欢迎以下类型的贡献：

- 🐛 **错误修复** - 修复已知问题
- ✨ **功能实现** - 实现计划中的功能
- 📚 **文档改进** - 改进或扩展文档
- 🧪 **测试增强** - 编写或改进单元测试
- ⚙️ **性能优化** - 优化代码性能
- 🔍 **代码审查** - 审查他人的拉取请求

---

## 报告错误

### 提交错误报告前

检查问题跟踪器以确保尚未报告此错误。

### 如何提交良好的错误报告

**标题**: 简洁明确的描述

**描述**: 尽可能详细地描述问题

**再现步骤**:
```
1. 运行命令 ...
2. 使用参数 ...
3. 观察到的行为 ...
```

**预期行为**: 应该发生什么

**实际行为**: 实际发生了什么

**环境信息**:
- Go版本: `go version`
- FFmpeg版本: `ffmpeg -version`
- 操作系统和版本
- 项目提交ID: `git log --oneline -1`

**附加上下文**: 任何其他相关信息

### 错误报告模板

```markdown
# [Bug标题]

## 环境
- Go版本: 
- FFmpeg版本: 
- OS: 

## 问题描述
[清晰描述问题]

## 再现步骤
1. 
2. 
3. 

## 预期行为
[描述预期行为]

## 实际行为
[描述实际行为]

## 日志输出
[粘贴相关日志]

## 附加信息
[任何其他信息]
```

---

## 建议功能

### 建议新功能前

- 检查是否已有相关Issue或Pull Request
- 检查 [CHANGELOG.md](./CHANGELOG.md) 中的计划路线图
- 确保功能与项目目标一致

### 如何建议功能

**标题**: 清晰的功能标题

**描述**: 功能的详细描述

**用例**: 为什么需要此功能？

**伪代码**: 预期的API或用法示例

**备选方案**: 考虑过的替代方案

### 功能建议模板

```markdown
# [功能标题]

## 功能描述
[清晰描述该功能]

## 问题所解决的
[描述当前问题]

## 预期行为
[功能应该如何工作的详细说明]

## 用例示例
```bash
./m3u8-downloader "url" [新参数]
```

## 附加上下文
[其他相关上下文]
```

---

## 拉取请求流程

### 1. Fork项目

```bash
# Fork后克隆您的fork
git clone https://github.com/your-username/m3u8-downloader.git
cd m3u8-downloader

# 添加上游远程
git remote add upstream https://github.com/original-owner/m3u8-downloader.git
```

### 2. 创建功能分支

```bash
# 从latest development分支创建
git fetch upstream
git checkout -b feature/your-feature upstream/develop

# 或修复bug:
git checkout -b fix/issue-description upstream/develop
```

**分支命名规范**:
- `feature/description` - 新功能
- `fix/description` - 错误修复
- `docs/description` - 文档更新
- `test/description` - 测试增强
- `perf/description` - 性能优化

### 3. 进行更改

- 遵循 [代码风格](#代码风格) 指南
- 添加或修改测试代码
- 更新相关文档
- 确保代码通过检查

```bash
# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 运行测试
go test ./...

# 检查竞态条件（并发代码）
go test -race ./...
```

### 4. 提交更改

遵循 [提交信息规范](#提交信息规范)

```bash
# 查看变化
git diff

# 分阶段提交
git add [files]
git commit -m "Your message"

# 或提交所有更改
git add .
git commit -m "Your message"
```

### 5. 保持与上游同步

```bash
# 获取最新更改
git fetch upstream

# 重新基于最新develop分支
git rebase upstream/develop

# 如有冲突，解决后继续
git rebase --continue
```

### 6. 推送到您的Fork

```bash
git push origin feature/your-feature
```

### 7. 创建Pull Request

**PR标题**: 清晰、符合约定的标题

**PR描述**:
```markdown
## 描述
[清晰描述您的更改]

## 相关Issue
Fixes #issue-number (如果适用)

## 更改类型
- [ ] 错误修复
- [ ] 新功能
- [ ] 破坏性更改
- [ ] 文档更新

## 检查清单
- [ ] 我已遵循项目的代码风格
- [ ] 我已自行审查了代码
- [ ] 我已对更改添加注释
- [ ] 我已添加必要的测试
- [ ] 所有新增和现有的单元测试都通过
- [ ] 我已更新相关文档

## 截图或日志（如适用）
[粘贴截图或日志输出]

## 附加信息
[任何其他相关信息]
```

### 8. Code Review和合并

- 维护者将审查您的PR
- 可能要求更改
- 所有检查通过后，您的PR将被合并

---

## 代码风格

### Go代码风格

遵循 [Effective Go](https://golang.org/doc/effective_go) 和 [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### 命名规范

```go
// 包名：小写，与目录同名
package config

// 常量：PascalCase
const (
    DefaultTimeout = 10 * time.Second
    MaxRetries     = 3
)

// 变量：camelCase（私有）或PascalCase（公开）
var (
    defaultConfig *Config  // 私有
    GlobalConfig  *Config  // 公开 - 应避免
)

// 函数：PascalCase（公开）或camelCase（私有）
func NewConfig() *Config { }      // 公开
func newDefaultConfig() *Config { } // 私有

// 结构体：PascalCase
type DownloadManager struct { }

// 接口：以"er"结尾（可选）
type Logger interface { }
type Reader interface { }
```

### 代码格式

```bash
# 自动格式化
go fmt ./...

# 使用gofmt或goimports
gofmt -s -w .
goimports -w .
```

### 错误处理

```go
// ✅ 好的做法
if err != nil {
    return errors.New(HTTPRequest, "request failed", err)
}

// ❌ 避免
if err != nil {
    panic(err)  // 不要 panic
}

// ❌ 避免
_ = err  // 忽略错误
```

### 注释

```go
// Package config 提供配置管理功能
package config

// DefaultConfig 返回默认的配置对象
func DefaultConfig() *Config {
    // ...
}
```

---

## 提交信息规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)

### 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 类型

- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档更新
- `style`: 代码风格更改（不影响功能）
- `refactor`: 代码重构（不改变功能）
- `perf`: 性能改进
- `test`: 测试相关更改
- `chore`: 依赖更新、工具配置等

### 示例

```bash
# 新功能
git commit -m "feat(downloader): add pause/resume support"

# 错误修复
git commit -m "fix(logger): prevent concurrent write race condition

Previously, multiple goroutines could write to the log simultaneously,
causing garbled output. Use mutex to synchronize access.

Fixes #123"

# 文档更新
git commit -m "docs: update installation instructions for Windows"

# 性能优化
git commit -m "perf(merger): optimize FFmpeg concat performance

Reduce memory usage by 30% when merging large videos."
```

### 提交信息指南

1. 使用英文（国际项目）或中文（保持项目语言一致）
2. 第一行简洁，不超过50个字符
3. 在第一行后留一空行
4. 换行的内容不超过72个字符
5. 在正文中详细说明**为什么**做出这些改变
6. 使用 "Fix #issue-number" 自动关闭相关Issue

---

## 文档贡献

### 文档位置

- `README.md` - 项目概览
- `docs/TEST.md` - 测试指南
- `docs/DEVELOPMENT.md` - 开发指南
- `docs/REFACTORING_REPORT_V2.0.md` - 架构报告
- `docs/CHANGELOG.md` - 更新日志

### 文档风格

- 使用清晰的、活泼的语言
- 包含代码示例
- 对复杂概念进行解释
- 包含目录和链接
- 使用恰当的Markdown格式

### Markdown风格

```markdown
# 一级标题

## 二级标题

### 三级标题

**加粗**，*斜体*，`代码`

- 列表项
- 列表项

1. 编号列表
2. 编号列表

> 引用

\`\`\`go
// 代码块
fmt.Println("Hello")
\`\`\`

[链接](https://example.com)
```

---

## 获取帮助

如有疑问，您可以：

1. 在GitHub Issues中提问
2. 检查 [FAQ](docs/DEVELOPMENT.md#常见问题)
3. 查看现有PR以了解贡献流程
4. 在讨论区提出问题

---

## 许可

通过贡献，您同意在MIT许可证下许可您的贡献。

---

感谢您的贡献！🎉

最后更新: 2024年
版本: 1.0
