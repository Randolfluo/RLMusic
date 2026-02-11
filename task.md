
# 🎯 目标

实现一个：

* C/S 架构音乐播放器
* 支持聊天室
* 支持多人严格同步听歌
* 不广播播放进度
* 仅通过时间线事件实现同步

---

# 🧠 核心数学模型（最高优先级）

客户端必须始终用此公式计算播放进度：

```
if paused:
    current_position = pause_position
else:
    current_position = (server_now - start_timestamp) * speed
```

后端维护的权威状态仅包含：

```
start_timestamp
pause_position
paused
speed
song_id
```

严禁后端或客户端实现“进度广播”。

---

# 🏗️ 系统模块划分

## 前端模块

```
core/player          # 音频播放
core/realtime        # WebSocket + Timeline Engine
store                # UI 状态
```

## 后端模块（Go）

```
ws/                  # WebSocket 连接
room/                # 房间管理
timeline/            # 时间线计算
redis/               # 房间状态存储
```

---

# 🔌 WebSocket 事件（唯一允许的同步方式）

AI 只能实现以下事件：

```
TIMELINE_INIT
PLAY
PAUSE
SEEK
CHANGE_SONG
SET_SPEED
CHAT
MEMBER_JOIN
MEMBER_LEAVE
TIME_SYNC
```

禁止添加“进度同步”类事件。

---

# 🧾 事件协议定义

## TIMELINE_INIT（加入房间时）

```json
{
  "event": "TIMELINE_INIT",
  "song_id": "xxx.mp3",
  "start_timestamp": 1700000000000,
  "pause_position": 0,
  "paused": false,
  "speed": 1.0,
  "server_time": 1700000000100
}
```

客户端收到后必须：

1. load(song)
2. 计算 current_position
3. seek
4. play / pause

---

## PLAY

```json
{
  "event": "PLAY",
  "server_time": 1700000123456
}
```

服务器：

```
start_timestamp = server_time - pause_position
paused = false
```

---

## PAUSE

```json
{
  "event": "PAUSE",
  "server_time": 1700000999999,
  "position_ms": 52342
}
```

服务器：

```
pause_position = position_ms
paused = true
```

---

## SEEK

```json
{
  "event": "SEEK",
  "server_time": 1700001111111,
  "position_ms": 120000
}
```

服务器：

```
start_timestamp = server_time - position_ms
paused = false
```

---

## CHANGE_SONG

```json
{
  "event": "CHANGE_SONG",
  "song_id": "xxx.mp3",
  "server_time": 1700002222222
}
```

服务器：

```
song_id = xxx
start_timestamp = server_time
pause_position = 0
paused = false
```

---

## SET_SPEED

```json
{
  "event": "SET_SPEED",
  "speed": 1.5,
  "server_time": 1700003333333
}
```

服务器：

```
pos = (server_time - start_timestamp) * old_speed
start_timestamp = server_time - pos / new_speed
speed = new_speed
```

---

# 🗄️ Redis 数据结构（必须按此实现）

```
room:{id}:timeline   -> Hash
room:{id}:members    -> Set
```

timeline Hash 字段：

```
song_id
start_timestamp
pause_position
paused
speed
```

---

# 🧩 前端 Timeline Engine 行为规范

WebSocket 收到任意时间线事件时：

```
更新本地 timeline 状态
计算 current_position
调用 player.seek(current_position)
根据 paused 决定 play/pause
```

UI 不允许直接操作 player。

---

# ⏱️ 时间校准（必须实现）

客户端必须维护：

```
offset = server_time - Date.now()
server_now = Date.now() + offset
```

用于代入时间公式。

---

# ❌ 严禁实现的错误方式

AI 不允许生成以下逻辑：

* 广播当前播放进度
* 定时同步进度
* 客户端互相对齐进度
* Host 直接控制其他客户端播放器

---

# ✅ 正确的控制链路

```
Host 操作
   ↓
发送事件到服务器
   ↓
服务器重建 Timeline
   ↓
广播事件
   ↓
所有客户端根据 Timeline 自行 seek
```

Host 也必须等待广播，不可本地直接播放。

---

# 📌 最终原则（AI 必须遵守）

> 这是一个 Timeline 驱动的播放器系统
> 而不是“远程控制播放器”的系统
