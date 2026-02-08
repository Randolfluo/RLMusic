# 基于VUE的音乐播放器



## TODO List

### 🛠️ 环境与部署 (Environment & Deployment)
- [ ] 检查 `sqlite` 依赖使用情况（是否冗余）
- [ ] Docker 部署：使用 MongoDB，通过挂载目录实现目录共享
- [ ] 后端容器化：支持 Electron 带有服务器和不带服务器（纯本地）的版本选择

### 💾 后端与数据 (Backend & Data)
- [ ] Redis 集成：维护发现页公共数据、实现 Token 黑名单机制
- [ ] 数据库优化：仅歌单存储封面，其他封面信息直接读取文件，减少数据库体积
- [ ] 文件管理：实现文件夹视图，默认属性为 private

### ✨ 功能特性 (Features)
- [ ] 歌曲详情：显示采样率、专辑组、同类歌曲推荐
- [ ] 局域网协同：支持多人听歌、聊天室功能
- [ ] 移动端联动：二维码扫描快速访问网页
- [ ] 效率工具：番茄钟功能（配备动画）
- [ ] 国际化：添加简繁体切换
- [ ] QWen3TTS 集成

### 🎨 界面与交互 (UI/UX)
- [ ] 视觉效果：添加启动动画
- [ ] 动态主题：主题色自动跟随歌曲封面
- [ ] 交互优化：添加长按呼出旋转轮盘
- [ ] 桌面交互：双击放大等基础 Electron 交互体验

### 🖥️ 桌面端增强 (Desktop Specific)
- [ ] 桌面歌词：创建背景透明 (`transparent: true`) 且忽略鼠标事件 (`ignoreMouseEvents`) 的置顶窗口
- [ ] 系统集成：全局快捷键管理（后台切歌）、Windows 任务栏缩略图控制（播放/暂停/切歌）
- [ ] 窗口管理：最小化到托盘、关闭时提示是否最小化到系统 Tray

### ⚙️ 设置与自定义 (Settings)
- [ ] 外观设置：增加主题色设置选项、自定义字体选择
- [ ] 音频设置：手动指定音频输出设备
- [ ] 行为设置：完善最小化到托盘和关闭应用的相关选项

//支持更多音乐拓展格式
// 专辑播放模式



服务端运行
cd server
go run .\cmd\main.go


客户端开发
pnpm dev
客户端构建
pnpm 

用户添加的我喜欢的歌单字段，扫描时创建的我喜欢的歌单自动加入收藏歌单




//添加歌曲播放组件跳转，组件在comoenent/player，接口位于request.http。bigplayer.vue是全屏播放组件


//添加歌单管理模式，支持物理方式删除歌单

//total_duartion和listening_duration字段值修复

//添加获取详细封面图接口,歌单封面,大屏播放时调用


//添加总播放次数，添加api请求次数
//添加cpu信息，内存信息等，将系统概述分为一个单独的页面


采用naiveui设计，尽量与项目风格统一，   写成更合适的prompt


//支持多关键词匹配

//修复自动登录功能




聊天室模型：
房主点击播放
    ↓
发送 NewTimeline 给服务器
    ↓
服务器写 Redis (version++)
    ↓
服务器 PUBLISH timeline
    ↓
所有客户端（含房主）收到
    ↓
Sync Scheduler 计算播放时间
    ↓
在未来精确时间点启动

新用户加入房间：

1. 从 Redis GET room timeline
2. 立即开始对齐播放
3. 再加入 websocket 广播

同步数据结构
RoomTimeline {
  song_id: "xxx",
  start_timestamp: 1700000000.235,   // 服务器时间
  offset: 12.5,
  paused: false,
  owner_id: "u1001",
  version: 42
}
不是请求-响应模型（Request/Response Thinking），而是状态广播模型（State Broadcast Model）
房主 提交 新时间线
服务器 覆盖 时间线
服务器 广播 时间线

HTTP与WebSocket
HTTP 是“请求-响应”
WebSocket 是“长连接、双向实时通信”

Redis 时间线读写
PubSub 广播:允许消息的发送者（发布者）将消息发送给多个接收者（订阅者）。用于服务器间同步
房主操作
   ↓
服务器写 Timeline 到 Redis
   ↓
PUBLISH 通知
   ↓
所有 WebSocket 服务器收到
   ↓
读取 Timeline
   ↓
推给客户端


| 领域           | 技术                       | 用途     | 是否关键 |
| ------------ | ------------------------ | ------ | ---- |
| WebSocket    | Gorilla / ws / socket.io | 实时连接   | ✅ 核心 |
| Redis        | Timeline + Pub/Sub       | 时间线权威源 | ✅ 核心 |
| MongoDB      | 歌曲 / 聊天持久化               | 数据存储   | ✅    |
| NTP 时间校准     | 自实现                      | 播放精确同步 | ✅ 核心 |
| 本地 Scheduler | 播放调度器                    | 精确启动播放 | ✅ 核心 |
| JWT          | 鉴权                       | 房间权限   | ✅    |
| Docker       | 部署                       | 环境一致   | ✅    |
| k8s（未来）      | 扩展                       | 横向扩容   | ⭐    |

1. 房间系统（Room）
2. 时间线系统（Timeline）⭐最核心
3. 聊天系统（Chat）
4. 播放调度系统（Client Scheduler）⭐最难

https://chatgpt.com/s/t_69875e32fe7c8191a6c1594adc92ccb5


异步实现