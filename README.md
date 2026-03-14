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

## 目录结构

项目采用 Monorepo 风格的目录组织，将前端 Web、Electron 主进程、Go 后端以及 Android 原生工程集中在一个仓库中进行管理。

```text
.
├─ src/                  # 前端 Vue 3 源码目录
│  ├─ api/               # API 接口统一封装，使用 Axios 进行网络请求
│  ├─ components/        # 可复用的 Vue 组件库 (如 Player、Nav、Admin 等)
│  ├─ core/              # 核心业务逻辑 (包含 Websocket 通信、Timeline 调度等)
│  ├─ router/            # Vue Router 路由配置及导航守卫
│  ├─ store/             # Pinia 状态管理 (音乐播放状态、用户数据、设置等)
│  ├─ style/             # 全局样式文件 (SCSS)
│  ├─ utils/             # 通用工具函数 (时间格式化、加密解密、防抖节流等)
│  └─ views/             # 页面级视图组件 (首页、歌单、搜索、管理后台等)
│
├─ server/               # Go (Gin) 后端源码目录
│  ├─ cmd/               # 服务端入口文件 (main.go)
│  ├─ internal/          # 内部核心逻辑，按职责划分
│  │  ├─ global/         # 全局变量、配置结构体、统一返回格式
│  │  ├─ handle/         # 控制器层，处理 HTTP 请求逻辑
│  │  ├─ middleware/     # Gin 中间件 (Auth鉴权、统计拦截等)
│  │  ├─ model/          # GORM 数据库模型定义及基础 DB 操作
│  │  └─ ws/             # WebSocket 服务端实现，处理长连接通信
│  ├─ utils/             # 后端工具包 (AI调用、音频解析、JWT、加密等)
│  ├─ prompts/           # LLM (大语言模型) 提示词 Markdown 模板
│  ├─ config.yml         # 后端服务主配置文件
│  └─ Dockerfile         # 后端 Docker 构建文件
│
├─ electron/             # Electron 桌面端主进程目录
│  ├─ main.ts            # 主进程入口，管理窗口、系统托盘、IPC 通信
│  └─ preload.ts         # 预加载脚本，向渲染进程安全暴露 Node.js API
│
├─ android/              # Capacitor 生成的 Android 原生工程目录
│
├─ public/               # 静态资源目录 (不经过 Vite 编译直接复制，如 favicon、默认头像)
│
├─ docker-compose.yml    # Docker 容器编排配置文件 (前后端一键部署)
├─ Dockerfile.web        # 前端 Nginx 部署构建文件
├─ vite.config.ts        # Vite 构建配置文件
├─ capacitor.config.ts   # Capacitor 跨平台配置文件
└─ package.json          # 项目依赖及 npm scripts 脚本定义
```

## 快速开始

### 1. 开发环境搭建 (Development)

适用于需要修改代码或参与贡献的开发者。

**前置依赖**：
*   [Node.js](https://nodejs.org/) (推荐 v18+)
*   [Go](https://go.dev/) (推荐 v1.20+)
*   [pnpm](https://pnpm.io/) (`npm install -g pnpm`)

**步骤**：

1.  **安装依赖**
    ```bash
    pnpm install
    ```

2.  **启动后端服务 (Server)**
    ```bash
    cd server
    go mod tidy
    air  # 使用 air 进行热重载开发，或者使用 `go run cmd/main.go`
    ```
    后端服务默认运行在 `http://localhost:12345`。

3.  **启动前端 (Web)**
    ```bash
    # 新开一个终端窗口
    pnpm dev:web
    ```
    前端页面默认运行在 `http://localhost:23456`。

4.  **启动 Electron 客户端**
    ```bash
    # 如果想调试 Electron 环境
    pnpm dev
    ```

### 2. 构建与部署 (Production)

适用于生产环境部署或生成安装包。

#### Web 部署 (Nginx)
构建纯静态资源，可部署在任何 Web 服务器上。
```bash
pnpm build:web
# 构建产物位于 dist/ 目录
```

#### Electron 安装包
构建 Windows/macOS/Linux 安装包。
```bash
# 构建客户端 (连接远程 Server)
pnpm build:client

# 构建服务端 (自带本地 Server，单机版)
pnpm build:server
```
构建产物位于 `release/` 目录。

#### Android APK
构建 Android 安装包 (需配置 Android Studio 环境)。
```bash
pnpm build:android
```

### 3. Docker 一键部署 (推荐)

项目支持前后端分离的 Docker 容器化部署。

1.  **准备环境**
    确保已安装 Docker 和 Docker Compose。

2.  **构建并启动**
    在项目根目录下执行：
    ```bash
    docker-compose up -d --build
    ```

3.  **访问应用**
    -   前端访问地址: `http://localhost`
    -   后端 API 地址: `http://localhost:12345`

4.  **数据持久化与音乐文件**
    -   数据库和日志文件会自动映射到当前目录下的 `data` 和 `log` 文件夹。
    -   **音乐文件**: 默认使用名为 `music_data` 的 Docker 卷。你需要将音乐文件复制到该卷中，或者修改 `docker-compose.yml` 将本地音乐目录映射到容器内的 `/music` 目录。
        例如，将本地 `C:\MyMusic` 映射进去：
        ```yaml
        volumes:
          - ./data:/app/data
          - ./log:/app/log
          - C:\MyMusic:/music  # 修改此处
        ```

## 核心功能详解

### 🎵 沉浸式播放体验
- **动态背景**：播放器背景根据专辑封面主色调实时生成高斯模糊与动态流光效果。
- **桌面歌词**：Electron 端支持独立的桌面歌词窗口，支持锁定、调整大小和双行显示。
- **无缝切换**：Web 端与移动端布局自适应，支持 PWA 安装，提供原生 App 般的体验。

### 🤝 一起听 (Listen Together)
- **实时同步**：基于 WebSocket 实现多端播放进度毫秒级同步。
- **互动聊天**：内置聊天室，支持发送表情和文本，与好友边听边聊。
- **房间管理**：支持创建私密房间，邀请制加入。

### 🤖 AI 实验室
- **智能描述**：利用 LLM 分析歌单内的歌曲风格、情感和流派，自动生成富有文采的歌单介绍。
- **播客合成**：将生成的文案通过 TTS 引擎转为语音，并在播放歌单时作为“电台开场白”自动播放，营造电台氛围。
- **本地扫描**：智能识别本地音乐文件的元数据（ID3 Tags），自动归类专辑和艺术家。

### 📱 移动端特性
- **扫码连接**：在 Web/Electron 端生成二维码，移动端 App 扫码即可直连局域网内的服务器，无需手动输入 IP。
- **手势操作**：支持侧滑返回、点击穿透等原生手势体验。



## 配置说明

### 环境变量 (.env)
在项目根目录创建 `.env` 文件以配置前端环境变量：

| 变量名 | 说明 | 默认值 |
| :--- | :--- | :--- |
| `VITE_MUSIC_API` | 后端 API 地址 | `http://localhost:12345` |
| `VITE_APP_MODE` | 应用模式 (web/client/server) | `web` |
| `VITE_ANN_TITLE` | 首页公告标题 | - |
| `VITE_ANN_CONTENT` | 首页公告内容 | - |

### 后端配置 (server/config.yml)
后端服务启动时会读取 `config.yml`，主要配置项包括：

```yaml
Server:
  Port: :12345        # 服务端口
  DbType: "sqlite"    # 数据库类型
Sqlite:
  Dsn: "data.db"      # 数据库文件路径
BasicPath:
  FilePath: "./music" # 音乐文件存储根目录
QwenTTS:              # 阿里云通义千问 TTS 配置
  ApiKey: "sk-..."
SiliconFlow:          # 硅基流动 LLM 配置
  ApiKey: "sk-..."
```

## 开发计划 (Roadmap)

- [x] **基础功能**
    - [x] 音乐播放与控制
    - [x] 歌单/专辑/艺术家管理
    - [x] 歌词显示与滚动
- [x] **多端支持**
    - [x] Web 端响应式布局
    - [x] Electron 桌面端 (Client/Server 模式)
    - [x] Docker 容器化部署
    - [x] Android 客户端 (Capacitor)
        - [x] 二维码扫描连接
        - [x] 原生返回手势
- [ ] **AI 增强**
    - [x] 歌单智能描述生成
    - [x] 播客开场白合成
    - [ ] 歌词智能翻译与纠错
    - [ ] 基于内容的智能推荐
- [ ] **高级特性**
    - [ ] 支持更多音频格式 (FLAC/WAV/APE)
    - [ ] 音乐标签 (Tag) 编辑器
    - [ ] DLNA/Chromecast 投屏支持

## 致谢

本项目使用了以下开源库：

- [Vue 3](https://vuejs.org/) & [Vite](https://vitejs.dev/)
- [Naive UI](https://www.naiveui.com/)
- [Pinia](https://pinia.vuejs.org/)
- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [Capacitor](https://capacitorjs.com/)
- [Electron](https://www.electronjs.org/)

## License
MIT
