# 执行报告：Electron 端 AI 智能套件 prompt 文件打包修复

## 任务摘要

用户反馈 Electron 桌面端管理员 "AI 智能套件"（歌单洞察、艺术家画像、专辑纪事、全站开场白）无法生成内容。根因为 `electron-builder.server.json` 的 `extraResources` 中遗漏了 `server/prompts/` 目录，导致 9 个 prompt markdown 模板文件未被纳入 Electron 安装包。后端 `Prompt.go` 在运行时查找 `prompts/` 目录失败 → `os.ReadFile` 返回 ENOENT → 所有依赖 prompt 模板的 AI 生成功能全部失败。

## 修改文件清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `electron-builder.server.json` | 修改 | `extraResources` 新增 `{"from": "server/prompts", "to": "prompts"}` |
| `package.json` | 修改 | `build:go` 脚本增加同步拷贝 `server/prompts/` → `resources/prompts/` |

## 关键变更说明

### 1. `electron-builder.server.json` L20-23 — 将 prompts 打入安装包

```diff
  "extraResources": [
    { "from": "resources/server.exe", "to": "server.exe" },
    { "from": "server/config.yml", "to": "config.yml" },
+   { "from": "server/prompts", "to": "prompts" }
  ],
```

### 2. `package.json` L14 — build:go 同步拷贝 prompts 到 resources

```diff
- "build:go": "cd server && go mod tidy && go build -o ../resources/server.exe cmd/main.go",
+ "build:go": "cd server && go mod tidy && go build -o ../resources/server.exe cmd/main.go && node -e \"const{cpSync}=require('fs');cpSync('prompts','../resources/prompts',{recursive:true})\"",
```

`build:go` 执行时 CWD 已在 `server/`，因此源路径为 `prompts`，目标为 `../resources/prompts`。

### 3. 影响范围

后端 `Prompt.go` `Read()` 函数按三级路径查找 prompt 文件：
1. `<exeDir>/prompts/<name>` — 部署形态：`server.exe` 同级目录
2. `<exeDir>/../prompts/<name>` — 后备：`bin/` 子目录形态
3. `./prompts/<name>` — 开发形态：CWD

修复前 Electron 打包产物中 `resources/` 仅有 `server.exe` 和 `config.yml`，无 `prompts/`。三种查找均失败。

修复后第 1 条路径（exeDir = `resources/`）可直接命中 `resources/prompts/`。

涉及的 9 个 prompt 文件：

| 文件 | 用途 |
|------|------|
| `prompt_playlist1.md` | 歌单 AI 分析（Step 1） |
| `prompt_playlist2.md` | 歌单 AI 描述（Step 2） |
| `prompt_artist1.md` | 艺术家 AI 分析（Step 1） |
| `prompt_artist2.md` | 艺术家 AI 描述（Step 2） |
| `prompt_album1.md` | 专辑 AI 分析（Step 1） |
| `prompt_album2.md` | 专辑 AI 描述（Step 2） |
| `prompt_podcast1.md` | 歌曲开场白分析 |
| `prompt_podcast2.md` | 歌曲开场白草稿 |
| `prompt_podcast3.md` | 歌曲开场白终稿 |

## 验证结果

- **本地 `resources/prompts/`**：9 个文件完整拷贝，路径正确
- **生产构建 `pnpm build:server`**：exit code 0，NSIS 安装包 + 便携版均成功打包
- **安装包内 `resources/prompts/`**：`release/server/1.0.0/win-unpacked/resources/prompts/` 包含全部 9 个 prompt 文件

## 回滚方案

```bash
git checkout electron-builder.server.json package.json
rm -rf resources/prompts
```

## 遗留问题

无。本修复仅涉及 Electron 打包清单，不改变任何运行时逻辑、API 接口或数据库 Schema。
