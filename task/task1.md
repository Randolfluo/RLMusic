请确定以下内容已经存在：

时间线数学模型已确定

WebSocket 事件协议已定义

前后端模块划分已完成

Redis 结构已确定

Timeline 字段定义已确定

🎯 接下来的目标

从“架构设计”进入“可运行的一起听系统”

🪜 实现顺序（必须按顺序，不可打乱）
Step 1️⃣：后端实现 RoomTimeline 数据结构（最先做）

实现一个结构体：

RoomTimeline {
    song_id
    start_timestamp
    pause_position
    paused
    speed
}


并实现 5 个函数：

ApplyPlay(serverTime)
ApplyPause(serverTime)
ApplySeek(serverTime, position)
ApplyChangeSong(serverTime, songId)
ApplySetSpeed(serverTime, newSpeed)


这些函数只修改 Timeline，不做任何网络操作。

Step 2️⃣：Redis 持久化 Timeline 和 Members

实现：

SaveTimeline(roomId, timeline)
LoadTimeline(roomId) -> timeline
AddMember(roomId, userId)
RemoveMember(roomId, userId)
GetMembers(roomId)


此时仍然没有 WebSocket。

Step 3️⃣：实现 WebSocket 连接和房间管理

实现能力：

用户连接 WS

加入房间

离开房间

广播消息到房间

此时 不要接入 Timeline，只验证广播正常。

Step 4️⃣：实现 TIMELINE_INIT

当用户加入房间时：

timeline = LoadTimeline(roomId)
发送 TIMELINE_INIT 给该用户


此时后端完成。

Step 5️⃣：前端实现 Player Core（独立于 WS）

必须先能本地做到：

load(song)
play()
pause()
seek(ms)
setSpeed(rate)
currentTime()


且完全不依赖网络。

Step 6️⃣：前端实现 Timeline Engine（最关键）

实现一个模块，输入：

timeline 状态
server_time


输出：

应该 seek 到的位置
应该 play 还是 pause


先写死数据测试，不要接 WS。

Step 7️⃣：前端实现时间校准 TIME_SYNC

实现：

offset = server_time - Date.now()
server_now = Date.now() + offset


验证时间计算正确。

Step 8️⃣：WebSocket 接入 Timeline Engine

现在才把 WS 接进来：

收到：

TIMELINE_INIT

PLAY

PAUSE

SEEK

CHANGE_SONG

SET_SPEED

全部交给 Timeline Engine。

WS 不允许直接操作 Player。

Step 9️⃣：把 Host 的操作接入事件链路

Host 点击：

播放

暂停

拖动进度条

切歌

改倍速

行为必须是：

发送事件 → 等待服务器广播 → Timeline Engine 驱动 Player


Host 不允许本地直接播放。

Step 🔟：测试“加入即同步”（最终验证）

测试流程：

A 进入房间播放 2 分钟

B 加入房间

B 必须自动播放到 A 的进度，无任何额外同步

如果成立，“一起听”功能完成。

🧪 最终判定标准（必须满足）

系统中：

没有任何“进度广播”

没有任何定时同步

所有同步只依赖 Timeline 事件

新用户加入自动对齐进度