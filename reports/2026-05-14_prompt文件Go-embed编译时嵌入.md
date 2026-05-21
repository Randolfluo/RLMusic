# 执行报告：prompt 文件改用 Go embed 编译时嵌入

## 任务摘要

修复 `go run` 开发模式下 prompt 文件读取失败的问题。原 `prompt.Read()` 依赖运行时文件系统路径查找，当通过 `go run` 启动时，可执行文件位于临时目录，CWD 也未必是 `server/`，导致所有候选路径均不存在，9 种 prompt 模板全部读取失败。

**修复方案：** 使用 Go 1.16+ 的 `embed` 包将 prompt 文件嵌入二进制，运行时无需依赖任何文件系统路径。同时保留文件系统回退以支持开发时快速迭代 prompt 措辞。

## 修改文件清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `server/prompts/embed.go` | 新增 | 使用 `//go:embed prompt_*.md` 将 9 个 prompt 文件嵌入二进制 |
| `server/internal/utils/prompt/Prompt.go` | 修改 | 优先从嵌入 FS 读取，文件系统作为回退 |
| `server/Dockerfile` | 修改 | 移除不再需要的 `COPY --from=builder /app/prompts` |
| `package.json` | 修改 | `build:go` 移除 `resources/prompts` 复制步骤 |

## 关键变更说明

### 1. embed.go — 编译时嵌入

```go
//go:embed prompt_*.md
var FS embed.FS
```

- 位置：`server/prompts/embed.go`，package `prompts`
- 模式 `prompt_*.md` 匹配 9 个文件：`prompt_{album,artist,playlist,podcast}{1,2,3}.md`
- 排除 `README.md`（不以 `prompt_` 开头）

### 2. Prompt.go — 读取顺序调整

读取顺序从「文件系统查找」改为「嵌入 FS 优先 + 文件系统回退」：
1. `fs.ReadFile(prompts.FS, name)` — 嵌入 FS，编译时确定，100% 可靠
2. 文件系统候选路径（`<exeDir>/prompts/` → `<exeDir>/../prompts/` → `CWD/prompts/`）— 开发回退

### 3. Dockerfile / package.json — 清理冗余步骤

- Dockerfile：嵌入后无需单独 `COPY prompts/`
- package.json `build:go`：嵌入后无需 `cpSync prompts → resources/prompts`

## 验证结果

- `go build ./...` — 通过
- `go vet ./...` — 通过
- `go test ./...` — 全部通过
- 嵌入的 9 个 prompt 文件均匹配正确，README.md 被排除

## 遗留问题

无。原文件系统回退路径保留，不影响现有部署方式（Docker、Electron、`go build` 均可用）。
