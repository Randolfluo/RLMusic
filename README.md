# 基于VUE的音乐播放器



## TODO List

### 🛠️ 环境与部署 (Environment & Deployment)
- [ ] Docker 部署：通过挂载目录实现目录共享
- [ ] 后端容器化：支持 Electron 带有服务器和不带服务器（纯本地）的版本选择

### 💾 后端与数据 (Backend & Data)
- [ ] Redis 集成：实现 Token 黑名单机制

### ✨ 功能特性 (Features)
- [ ] 局域网协同：支持多人听歌、聊天室功能
- [ ] 移动端联动：二维码扫描快速访问网页
- [ ] 效率工具：番茄钟功能（配备动画）
- [ ] 国际化：添加简繁体切换
- [ ] QWen3TTS 集成：实现语音合成功能

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
// 添加ai推荐和聊天功能


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


https://chatgpt.com/s/t_69875e32fe7c8191a6c1594adc92ccb5


异步实现



8. 优化歌曲页ui设计。
10. 修改私有歌单的歌单为为软链接，存储歌曲id和开场白文件名，播放时根据i文件路径。添加api实现开场白功能，根据歌手名，专辑名，歌曲名生成开场白，并写入到静态文件夹data/Podcast/。
10. 支持歌单描述、专辑描述、歌手描述的获取和修改。scan后开启一个线程修改数据库中的描述信息（描述信息通过POST {{baseUrl}}/ai/chat接口生成）
11. 进入歌单页显示详细的封面描述图

修改扫描用户音乐为管理员组内才能执行，并移动api到system.http。

添加搜索推荐栏

生活不止，修改为随机谚语
修复个人信息错误的问题


ai pipeline 实现

两阶段 AI Pipeline（Analyze → Generate）

修复重置系统数据库不完全

检查setting.view的每个配置是否都能正常工作