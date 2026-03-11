# RLMusic

基于 Vue 3 的多端音乐播放器，支持 Web、Electron 与 Android（Capacitor）构建。内置播放列表、歌词、播放控制、歌手/专辑/歌曲页面等功能，强调体验与界面一致性。

## 功能亮点
- **多端一致体验**：支持 Web、Windows/macOS/Linux (Electron)、Android，保持统一的设计语言。
- **AI 智能特性**：
    - **播客生成**：基于歌曲信息，利用 LLM (Qwen) 生成深度解析文案，并结合 TTS (CosyVoice) 自动合成电台级播客开场白。
    - **智能推荐**：根据听歌习惯与偏好，生成个性化推荐与“猜你喜欢”列表。
    - **歌词/文案优化**：AI 辅助生成或优化歌词翻译、专辑介绍与艺术家背景故事。
- **沉浸式播放**：
    - 多视图歌曲列表（缩略图/简洁模式）。
    - 桌面级歌词显示、动态背景与粒子特效。
    - 完美支持无缝播放与高解析度音频。
- **社交与互动**：
    - “一起听”房间：支持多端实时同步播放进度，内置 WebSocket 聊天室。
    - 评论与弹幕：歌曲评论区支持 Markdown 与表情互动。
- **管理与扩展**：
    - 强大的后台管理系统，支持用户、歌单、评论管理。
    - 插件化架构，易于扩展新的音频源与数据服务。

## 技术栈
- **前端**：Vue 3 + Vite + Pinia + Naive UI + TypeScript
- **桌面端**：Electron + Electron Builder
- **移动端**：Capacitor (Android/iOS)
- **后端**：Go (Gin) + Gorm + SQLite/MySQL
- **AI 服务**：
    - LLM: Qwen (通义千问) / SiliconFlow
    - TTS: CosyVoice / EdgeTTS
- **实时通信**：WebSocket (Gorilla WebSocket)
- **数据存储**：Redis (缓存/消息队列) + MongoDB (日志/非结构化数据)

## 目录结构
```
.
├─ src/                # 前端源码
├─ server/             # Go 后端
├─ electron/           # Electron 主进程
├─ public/             # 静态资源
├─ docker-compose.yml  # Redis/Mongo 组件
└─ package.json
```

## 快速开始

### 前端开发
```bash
pnpm install
pnpm dev
```

### Web 构建
```bash
pnpm build:web
```

构建产物在 dist/，可使用 Nginx / Caddy / 任意静态服务器部署。

### Electron 构建
```bash
pnpm build:client
pnpm build:server
```

### Android 构建
```bash
pnpm build:web
pnpm exec cap sync android
pnpm exec cap open android
```

命令行方式（适合 CI）：
```bash
cd android
./gradlew assembleDebug
```

## Docker（基础设施）
当前 docker-compose 提供 Redis 与 MongoDB：
```bash
docker compose up -d
```

## 计划
- Web 前端 Nginx Docker 镜像与部署模式
- 前后端拆分容器部署方案
- 修复桌面歌词的显示问题，置顶透明窗口
- 更新操作文档
- 向量数据库（如 Milvus）可以将歌曲与用户兴趣表示为高维向量，并通过高效的相似度检索（如 HNSW / IVF）快速返回 Top-K 相似歌曲，从而实现可扩展、实时的个性化音乐推荐系统。
- 移动端添加扫描二维码功能。
## License
MIT
