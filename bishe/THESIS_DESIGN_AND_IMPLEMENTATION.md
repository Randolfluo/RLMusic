# 系统设计与实现总结（对应论文第4、5章）

> 说明：本内容参考了你提供的目录结构（4.1~4.5, 5.1~5.5），并结合 RLMusic 项目的实际业务进行了适配。请注意，原图中的“线上考试”、“课程学习”等模块已替换为本项目核心功能（如音乐播放、曲库扫描）。

---

## 第4章 系统设计

### 4.1 总体架构设计

本系统采用基于 B/S（Browser/Server）与 C/S（Client/Server）混合模式的前后端分离架构，整体分为四层：

1.  **表现层（Presentation Layer）**：
    *   基于 **Vue 3 + TypeScript** 构建响应式 Web 界面。
    *   通过 **Electron** 封装为桌面客户端，提供原生系统托盘与本地文件访问能力。
    *   通过 **Capacitor** 封装为 Android 应用，适配移动端触控交互。
2.  **网关与接口层（Interface Layer）**：
    *   后端基于 **Gin** 框架提供 RESTful API。
    *   使用 **JWT** 进行无状态身份认证。
    *   通过 **WebSocket** 实现实时双向通信（用于“一起听”功能）。
3.  **业务逻辑层（Business Layer）**：
    *   包含用户管理、音乐扫描、元数据解析、播放控制、AI 内容生成等核心模块。
    *   利用 Go 语言的 `goroutine` 实现高性能并发文件扫描与转码。
4.  **数据持久层（Data Layer）**：
    *   使用 **Gorm** 作为 ORM 框架，支持 **SQLite**（轻量级/默认）和 **MySQL**（生产环境）。
    *   非结构化数据（音频文件、封面图、日志）直接存储于本地文件系统，并通过 Docker Bind Mount 进行持久化。

### 4.2 模块设计

系统自顶向下划分为以下核心子系统：

1.  **前台客户端子系统**：
    *   **播放模块**：音频解码、进度控制、歌词同步、封面渲染。
    *   **交互模块**：歌单管理、收藏点赞、播放历史、实时房间。
2.  **后台管理子系统**：
    *   **资源管理**：本地路径配置、曲库扫描、元数据修正。
    *   **用户管理**：用户列表、权限分配、封禁/解封。
    *   **运维管理**：系统状态监控、Excel 数据导出、日志查看。
3.  **AI 增强子系统**：
    *   **文案生成**：基于 LLM 生成歌手/专辑介绍。
    *   **语音合成**：基于 TTS 生成播客风格开场白。

### 4.3 功能模块

| 模块名称 | 功能点描述 |
| :--- | :--- |
| **用户模块** | 注册、登录、个人信息修改、头像上传、密码重置。 |
| **音乐资源** | 本地目录扫描、ID3 信息提取、封面哈希去重、音频流分块传输。 |
| **播放控制** | 播放/暂停/切歌、播放模式（循环/随机）、音量控制、倍速播放、桌面歌词。 |
| **歌单管理** | 创建/修改/删除歌单、添加歌曲、公共歌单广场、歌单订阅。 |
| **AI 服务** | 歌手/专辑智能介绍生成、歌曲情感分析、播客开场白语音合成。 |
| **系统运维** | CPU/内存监控、接口调用统计、全库 Excel 导出、数据重置。 |

### 4.4 E-R图与数据库设计

系统核心实体关系如下：

*   **User (用户)** 1 : N **Playlist (私有歌单)**
*   **User (用户)** N : M **Playlist (收藏歌单)**
*   **Playlist (歌单)** N : M **Song (歌曲)**
*   **Song (歌曲)** N : M **Artist (歌手)**
*   **Song (歌曲)** N : 1 **Album (专辑)**
*   **Song (歌曲)** N : 1 **Cover (封面)**

### 4.5 系统数据库表设计

| 表名 | 中文名 | 核心字段 | 说明 |
| :--- | :--- | :--- | :--- |
| `user` | 用户表 | id, username, password, role, avatar | 存储账号与权限 |
| `song` | 歌曲表 | id, title, file_path, duration, artist_id | 核心资源表，存储元数据 |
| `playlist` | 歌单表 | id, title, owner_id, is_public | 存储歌单信息 |
| `artist` | 歌手表 | id, name, description | 歌手基础信息 |
| `cover` | 封面表 | id, hash, path | 图片去重存储 |
| `history` | 历史表 | id, user_id, song_id, created_at | 播放记录 |

---

## 第5章 系统实现

### 5.1 用户登录注册

**功能描述**：
实现用户的身份验证与授权。用户提交账号密码，后端校验通过后签发 JWT Token，前端持久化存储并自动续期。

**核心代码映射**：
*   **前端**：`src/views/Login/LoginView.vue` (表单提交)、`src/utils/request.ts` (Token 拦截注入)。
*   **后端**：`server/internal/handle/handle_auth.go` (Login/Register 接口)、`server/internal/utils/jwt/Jwt.go` (令牌生成)。

**实现细节**：
1.  密码使用 bcrypt 算法加密存储。
2.  登录成功返回 Token，前端存入 SessionStorage。
3.  路由守卫 (`src/router/index.ts`) 拦截未登录访问 `meta.needLogin` 的页面。

### 5.2 菜单权限管理（后台管理）

**功能描述**：
根据用户角色（Admin/User）控制后台菜单的可见性与接口访问权限。

**核心代码映射**：
*   **前端**：`src/components/Nav/index.vue` (根据 `userStore.isAdmin` 显隐管理入口)。
*   **后端**：`server/internal/middleware/auth.go` (中间件校验 `UserGroup == 'admin'`)。

**实现细节**：
1.  前端：Pinia 状态管理存储 `userInfo.role`，动态渲染导航栏。
2.  后端：敏感接口（如扫描、导出）增加权限校验中间件，越权请求直接返回 403。

### 5.3 实时通讯（一起听）

**功能描述**：
利用 WebSocket 实现多用户同时收听同一首歌，同步播放进度与状态。

**核心代码映射**：
*   **前端**：`src/core/realtime/socket.ts` (Socket 连接与事件监听)。
*   **后端**：`server/ws/server.go` (连接池管理、消息广播)。

**实现细节**：
1.  建立 WS 连接：`ws://host/api/ws/chat`。
2.  定义协议：`SYNC_PLAY` (同步播放), `SYNC_PAUSE` (同步暂停), `SYNC_SEEK` (同步进度)。
3.  心跳保活：客户端定时发送 Ping，服务端回复 Pong，超时断连。

### 5.4 音乐播放与控制（对应原图“线上考试”）

**功能描述**：
实现音频文件的流式加载、解码播放及播放器 UI 交互。

**核心代码映射**：
*   **前端**：`src/components/Player/index.vue` (Audio 标签封装)、`src/store/musicData.ts` (播放状态管理)。
*   **后端**：`server/internal/handle/handle_song.go` (`StreamSong` 接口)。

**实现细节**：
1.  后端使用 `c.File()` 支持 HTTP Range 请求，实现音频拖拽缓冲。
2.  前端监听 `timeupdate` 事件驱动进度条与歌词滚动。
3.  利用 `MediaSession API` 适配系统原生媒体控制中心（通知栏切歌）。

### 5.5 本地资源扫描与管理（对应原图“课程学习”）

**功能描述**：
遍历宿主机指定目录，解析音频文件元数据并建立数据库索引。

**核心代码映射**：
*   **前端**：`src/components/Admin/AdminDashboard.vue` (触发扫描)。
*   **后端**：`server/internal/handle/handle_song.go` (`ScanUserMusic`)。

**实现细节**：
1.  使用 `github.com/dhowden/tag` 库解析 ID3/FLAC 标签。
2.  计算封面图片 MD5 哈希，实现“多首歌共用一张图”的去重存储。
3.  使用信号量（Semaphore）限制并发扫描协程数，防止磁盘 I/O 过载。

### 5.6 AI 文案与播客生成（特色功能）

**功能描述**：
集成大模型 API 生成歌手介绍，并利用 TTS 生成语音开场白。

**核心代码映射**：
*   **前端**：`src/api/ai.ts`。
*   **后端**：`server/internal/handle/handle_ai.go`、`server/internal/utils/ai/QwenTTS.go`。

**实现细节**：
1.  后端组装 Prompt（如“你是一个深夜电台主持人...”）。
2.  调用阿里 Qwen 或 SiliconFlow API 获取文本。
3.  调用 TTS 接口生成 WAV 文件并保存，播放时优先加载该音频。
