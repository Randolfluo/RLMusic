核心原则

不广播播放进度

不做定时同步

只通过 Timeline 事件 同步

客户端统一公式：

if paused:
  pos = pause_position
else:
  pos = (server_now - start_timestamp) * speed

后端权威字段：

song_id
start_timestamp
pause_position
paused
speed
必须支持的事件
TIMELINE_INIT
PLAY
PAUSE
SEEK
CHANGE_SONG
SET_SPEED
服务器对事件的处理（只改 Timeline）
PLAY
start_timestamp = T - pause_position
paused = false
PAUSE
pause_position = (T - start_timestamp) * speed
paused = true
SEEK
start_timestamp = T - position_ms
paused = false
CHANGE_SONG
song_id = x
start_timestamp = T
pause_position = 0
paused = false
SET_SPEED
pos = (T - start_timestamp) * old_speed
start_timestamp = T - pos / new_speed
speed = new_speed
Redis
room:{id}:timeline  (Hash)
room:{id}:members   (Set)
前端必须实现
Player Core
load / play / pause / seek / setSpeed
Timeline Engine

收到任意事件：

更新 timeline
计算 pos
player.seek(pos)
按 paused 决定 play/pause
时间校准
offset = server_time - Date.now()
server_now = Date.now() + offset
正确链路
Host 操作
 → 发事件
 → 服务器重建 Timeline
 → 广播
 → 所有客户端自行 seek

Host 也必须等广播。

完成判定

新用户加入自动对齐进度

无任何进度同步代码

同步仅依赖 Timeline 事件