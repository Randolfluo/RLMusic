# 基于 Vue3 与 Go 的多端本地音乐系统设计与实现（毕业设计论文初稿）

> 说明：本稿为“可直接二次修改”的毕业设计初稿，已结合当前项目真实实现撰写。学校格式要求（封面、页眉、参考文献格式、致谢字数）请按学院模板调整。

---

## 摘要

随着数字音乐消费场景从单一终端向多终端协同演进，传统本地播放器在跨平台一致性、数据管理能力和功能扩展性方面逐渐暴露不足。针对上述问题，本文设计并实现了一个面向 Web、桌面端（Electron）与移动端（Android）的多端本地音乐系统 RLMusic。系统采用前后端分离架构，前端基于 Vue 3 + TypeScript + Vite 构建交互界面，后端基于 Go + Gin + Gorm 提供接口与业务服务，数据库支持 SQLite / MySQL。  

系统围绕“本地曲库管理、跨端播放控制、后台运营管理、AI 内容增强”四个核心能力展开：一方面，支持对本地音频目录进行扫描、元数据解析、封面提取与去重入库；另一方面，提供播放控制、歌词展示、历史记录、歌单管理、用户权限管理等完整业务闭环。针对管理侧的数据可运维需求，系统实现了数据库全表 Excel 导出，并增加导出自检机制（响应类型与文件头签名校验），解决导出文件损坏风险。部署层面，系统提供 Docker 一键部署方案，支持 Bind Mount 将宿主机目录映射至容器，实现音乐资源和业务数据的长期持久化。  

实验与功能验证结果表明，该系统具备良好的跨端一致性、可维护性与可扩展性，能够满足中小规模本地音乐管理和播放场景需求，并为后续推荐算法、观测性建设和多媒体能力扩展提供了可行基础。

**关键词**：多端音乐系统；Vue 3；Go；Gin；Gorm；Docker；JWT；Excel 导出；AI 增强

---

## Abstract

With the evolution of digital music consumption from single-device usage to multi-device collaboration, traditional local music players show limitations in cross-platform consistency, data management, and feature extensibility. To address these issues, this thesis designs and implements RLMusic, a multi-platform local music system targeting Web, desktop (Electron), and Android. The system adopts a frontend-backend separated architecture. The frontend is built with Vue 3, TypeScript, and Vite, while the backend is implemented with Go, Gin, and Gorm, with support for SQLite/MySQL.  

The system focuses on four core capabilities: local music library management, cross-platform playback control, administrative operations, and AI-enhanced content generation. It supports local directory scanning, metadata extraction, cover parsing and deduplication, as well as playback, lyrics, history, playlist management, and role-based user administration. For operational requirements, a full-database Excel export mechanism is implemented with a self-check process (content-type and file signature validation), reducing the risk of corrupted export files. The deployment layer provides Docker-based one-click deployment and Bind Mount support for persistent media/data storage.  

Experimental validation indicates that the system achieves good cross-platform consistency, maintainability, and extensibility, and is suitable for small-to-medium local music management scenarios, while also providing a solid base for future recommendation algorithms, observability, and multimedia extensions.

**Key Words**: Multi-platform music system; Vue 3; Go; Gin; Gorm; Docker; JWT; Excel export; AI enhancement

---

## 第1章 绪论

### 1.1 研究背景

在流媒体音乐服务广泛普及的背景下，本地音乐管理并未消失，反而在高音质收藏、离线播放、私有化数据控制等场景中持续存在。尤其是个人 NAS、本地硬盘曲库与家庭局域网播放器等使用模式，对系统提出了三类核心需求：  
1) 多终端一致访问（电脑、移动端、Web）；  
2) 本地数据可控（音乐文件、数据库、日志可持续保存）；  
3) 功能可扩展（从单纯播放向管理、分析、AI 辅助过渡）。  

传统单机播放器在交互一致性、管理能力和平台覆盖方面难以满足当前需求，具备跨端架构能力、可部署能力和可维护能力的新型本地音乐系统具有较高研究价值与工程价值。

### 1.2 研究目的与意义

本文目标是实现一个工程化、可部署、可扩展的多端本地音乐系统，形成完整的软件工程交付闭环，包括需求分析、架构设计、数据库设计、接口设计、功能实现与测试验证。  

研究意义体现在：  
- **工程实践价值**：覆盖前后端分离、跨端容器、鉴权、运维导出等完整链路；  
- **教学实践价值**：适合软件工程毕业设计“从设计到实现”的全流程训练；  
- **扩展研究价值**：为推荐系统、智能标签、观测性与媒体中台能力提供基础。

### 1.3 国内外研究现状（简述）

国外流媒体平台重在推荐与版权分发，国内播放器产品强调生态整合与社交化体验。相较之下，本地音乐系统研究常聚焦于文件播放能力，缺少“管理后台 + 跨端一致 + AI 扩展 + 可部署性”的综合实现。本文工作重点在于将这些能力整合为统一系统，并落地为可运行工程。

### 1.4 本文主要工作

1. 设计并实现前后端分离、多端统一的本地音乐系统架构；  
2. 实现本地曲库扫描、元数据解析、封面去重、关系数据落库；  
3. 实现用户权限、歌单管理、播放历史、后台管理能力；  
4. 实现 AI 描述和播客开场白生成链路；  
5. 实现全库 Excel 导出与导出自检，提升数据可运维性；  
6. 提供 Docker 部署方案并完成功能与构建验证。

### 1.5 论文结构安排

- 第1章：绪论  
- 第2章：需求分析  
- 第3章：系统总体设计  
- 第4章：详细设计  
- 第5章：系统实现  
- 第6章：系统测试  
- 第7章：总结与展望

**本章小结**：本章明确了研究背景、目标与论文组织结构，为后续需求与设计章节提供问题导向。

---

## 第2章 需求分析

### 2.1 用户角色分析

系统角色划分为：  
- **游客用户**：可浏览公开资源，不能进行私有操作；  
- **普通用户**：可登录后使用私有歌单、收藏、历史等能力；  
- **管理员**：可访问后台，执行用户管理、公共歌单管理、曲库扫描、数据导出等操作。

### 2.2 功能需求分析

1. **曲库扫描需求**：支持配置本地目录，批量扫描并自动更新曲库；  
2. **播放需求**：支持播放控制、歌词同步、封面展示、播放模式切换；  
3. **账户需求**：支持注册登录、权限校验、头像与密码管理；  
4. **歌单需求**：支持创建/修改/删除私有歌单，浏览与管理公共歌单；  
5. **历史需求**：自动记录播放历史，支持查询和清空；  
6. **管理需求**：支持系统统计监控、用户管理、数据导出；  
7. **AI 需求**：支持歌手/专辑/歌单文案与播客开场白生成。

### 2.3 非功能需求分析

- **性能**：千级曲目数据下保持可接受响应；  
- **可扩展性**：模块化设计，便于新增业务与平台；  
- **可靠性**：异常兜底与容错机制（默认封面、导出自检）；  
- **可部署性**：Docker 一键启动，路径挂载可配置；  
- **安全性**：JWT 鉴权、管理员二次权限校验、敏感操作隔离。

### 2.4 用例概述

- 用例 A：管理员执行曲库扫描  
- 用例 B：用户创建并维护私有歌单  
- 用例 C：用户播放歌曲并查看歌词  
- 用例 D：管理员导出数据库并保存本地  

**本章小结**：本章从角色与功能两方面定义系统边界，为总体架构和数据库设计提供约束依据。

---

## 第3章 系统总体设计

### 3.1 架构设计

系统采用四层架构：  
1) 表现层（Vue 3 页面与组件）；  
2) 容器层（浏览器 / Electron / Capacitor）；  
3) 服务层（Gin REST API + 中间件）；  
4) 数据层（Gorm + SQLite/MySQL + 本地文件系统）。  

该结构兼顾跨端复用与职责隔离，前后端通过 HTTP API 解耦，部署通过 Docker 标准化。

### 3.2 技术选型说明

- 前端：Vue 3 生态成熟，组合式 API 便于复杂状态管理；  
- 后端：Go 在并发和部署体积方面优势明显；  
- ORM：Gorm 提升实体与查询维护效率；  
- 桌面端：Electron 便于接入文件系统和本地进程控制；  
- 移动端：Capacitor 降低多端开发复杂度；  
- 部署：Docker 降低环境差异风险。

### 3.3 目录结构设计

系统源码按“前端 / 后端 / 桌面端 / 部署”划分：  
- 前端：`src/`  
- 后端：`server/`  
- 桌面端：`electron/`  
- 部署：`docker-compose.yml`、`Dockerfile.web`、`server/Dockerfile`

### 3.4 运行拓扑

标准部署下：  
- `frontend` 容器对外暴露 `23456`，负责静态资源与反向代理；  
- `backend` 容器对外暴露 `12345`，负责业务 API；  
- 宿主机通过 Bind Mount 提供 `/music` 数据源与 `data/log` 持久化目录。

**本章小结**：本章完成技术路线与宏观架构落地，为后续详细设计提供实施框架。

---

## 第4章 详细设计

### 4.1 前端详细设计

#### 4.1.1 路由设计
- 由 `src/router/routes.ts` 统一维护页面路由；  
- `meta.needLogin` 控制登录态访问；  
- `src/router/index.ts` 执行初始化流程与权限守卫。

#### 4.1.2 状态管理设计

Pinia 划分为四类核心状态：
- `musicData.ts`：播放列表、当前歌曲、进度、模式、音量；
- `userData.ts`：用户信息、登录态、权限组；
- `settingData.ts`：主题、偏好与客户端参数；
- `chatData.ts`：实时互动消息状态。  

使用 `pinia-plugin-persistedstate` 持久化关键状态，减少刷新丢失。

#### 4.1.3 请求链路设计
- `src/utils/request.ts` 统一 Axios 实例；
- 请求拦截器注入 Token；
- 响应拦截器处理业务码、自动重登、失效跳转。

### 4.2 后端详细设计

#### 4.2.1 路由与中间件
- 路由注册中心：`server/internal/Manager.go`；
- 鉴权中间件：`server/internal/middleware/auth.go`；
- 统计中间件：`server/internal/middleware/stats.go`；
- CORS 与 DB 注入：`server/internal/middleware/base.go`。

#### 4.2.2 模块划分
- 认证模块：`handle_auth.go`  
- 歌曲模块：`handle_song.go`  
- 系统模块：`handle_system.go`  
- 用户模块：`handle_user.go`  
- 搜索模块：`handle_search.go`  
- AI 模块：`handle_ai.go`

### 4.3 数据库详细设计

核心实体：`user`、`song`、`playlist`、`artist`、`album`、`cover`、`history`、`system_info`。  
核心中间表：`playlist_songs`、`song_artists`、`user_subscribed_playlists`。  

通过 `server/internal/model/ZBase.go` 完成自动迁移，降低版本演进成本。

### 4.4 关键时序设计（简述）

- 登录时序：登录请求 -> 密码校验 -> JWT 下发 -> 前端注入 Token；  
- 播放时序：请求歌曲流 -> 播放器加载 -> 进度事件更新 -> 历史上报；  
- 导出时序：请求导出 -> 客户端自检 -> 服务端生成 -> Web 下载/IPC 落盘。

**本章小结**：本章从路由、状态、模块、数据模型与时序层面完成了系统级详细设计。

---

## 第5章 系统实现

### 5.1 本地曲库扫描实现

系统通过扫描用户指定目录，识别音频元数据并入库：  
- 解析字段：歌名、时长、艺术家、专辑、封面；  
- 封面处理：哈希去重并落地文件；  
- 数据同步：更新歌曲和统计信息。  

对应代码：  
- `server/internal/handle/handle_song.go`  
- `server/internal/model/Song.go`  
- `server/internal/model/Metadata.go`

### 5.2 播放与歌词实现

前端播放器组件驱动 `<audio>` 元素，并将播放状态同步到 Pinia；  
歌词通过时间戳映射实现滚动高亮。  

对应代码：  
- `src/components/Player/index.vue`  
- `src/components/Player/PlayerControl.vue`  
- `src/utils/lyricFormat.js`

### 5.3 管理后台实现

后台提供系统统计、用户管理、公共歌单管理、导出等能力。  

对应代码：  
- `src/components/Admin/AdminDashboard.vue`  
- `src/views/Admin/UserManage.vue`  
- `src/views/Admin/PlaylistManage.vue`

### 5.4 Excel 导出实现

实现要点：  
1. 后端动态遍历数据库所有表并逐表写入 xlsx；  
2. 前端导出前进行 Content-Type 与文件头签名自检；  
3. Electron 环境通过 IPC 保存文件并兼容多种二进制输入类型。  

对应代码：  
- `server/internal/handle/handle_system.go`  
- `src/components/Admin/AdminDashboard.vue`  
- `electron/main.ts`

### 5.5 Docker 部署实现

采用双容器编排（frontend/backend），并通过 Bind Mount 挂载音乐目录。  
默认音乐目录 `C:/RLMusic`，支持通过 `.env` 自定义。  

对应文件：  
- `docker-compose.yml`  
- `Dockerfile.web`  
- `server/Dockerfile`

**本章小结**：本章完成了系统关键功能的工程实现描述，证明设计方案具备可落地性。

---

## 第6章 系统测试

### 6.1 测试环境

- 操作系统：Windows  
- 前端构建工具：Vite + vue-tsc  
- 后端语言环境：Go  
- 部署工具：Docker Desktop  

### 6.2 功能测试

测试覆盖登录、播放、扫描、歌单、导出、后台管理等主流程，结果满足预期。  
建议在正式稿中附测试截图和接口返回示例。

### 6.3 安全测试

重点验证：  
- 越权访问管理员接口是否被拦截；  
- 失效 Token 是否被拒绝；  
- 非法路径请求是否被过滤。  

### 6.4 性能与稳定性测试

重点指标：  
- 千级曲库扫描耗时；  
- 首屏加载耗时；  
- 全库导出耗时与成功率。

### 6.5 测试结论

系统在主要业务路径上表现稳定，关键功能均可正常运行，满足毕业设计交付要求。

**本章小结**：通过功能、安全和性能测试，系统可用性与稳定性达到预期目标。

---

## 第7章 总结与展望

### 7.1 工作总结

本文完成了一个多端本地音乐系统的完整实现，覆盖需求分析、架构设计、数据库建模、功能实现、测试与部署，形成了具有工程可交付性的毕业设计成果。系统在跨端一致性、后台运维能力、AI 增强能力方面具有明显实践价值。

### 7.2 不足分析

1. 推荐能力仍以规则为主，个性化深度不足；  
2. 监控体系尚未形成完整观测闭环；  
3. 移动端能力仍可继续原生增强。

### 7.3 后续展望

1. 引入混合推荐模型（协同过滤 + 内容理解）；  
2. 引入 Prometheus/Grafana 构建可观测体系；  
3. 增加更多音频格式与外部设备（DLNA/Chromecast）支持；  
4. 提供更完整的数据备份与恢复策略。

**本章小结**：系统已具备实用能力和扩展基础，未来可继续向智能化与平台化演进。

---

## 参考文献（示例占位，请按学校格式替换）

[1] Vue.js Documentation.  
[2] Gin Web Framework Documentation.  
[3] Gorm Documentation.  
[4] Electron Documentation.  
[5] Docker Documentation.  
[6] 软件工程课程相关教材与论文资料。  

---

## 致谢（示例占位）

感谢指导教师在课题选题、系统设计和论文撰写过程中的耐心指导；感谢同学在联调与测试过程中的帮助；感谢开源社区提供的高质量框架与工具支持。

