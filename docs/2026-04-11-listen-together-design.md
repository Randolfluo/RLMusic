# 一起听歌直播间功能设计文档

## 1. 概述

### 1.1 功能目标
实现公开直播间模式的一起听歌功能，支持多直播间、用户限加入一个、房主控制播放、智能时间线同步。

### 1.2 核心特性
- 公开直播间：所有登录用户可加入
- 单房间限制：用户同一时间只能在一个直播间
- 房主控制：仅房主可切歌、暂停、跳转
- 智能同步：服务器维护时间线，新用户瞬间同步当前进度
- 实时聊天：文字聊天 + 系统消息
- 移动端适配：响应式设计，适配手机和桌面

## 2. 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        前端 (Vue 3)                         │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ 直播间列表页 │  │   直播间内页  │  │    聊天组件     │   │
│  │  - 简单列表  │  │  - 播放器控制 │  │  - 文字消息     │   │
│  │  - 显示人数  │  │  - 成员列表   │  │  - 系统消息     │   │
│  │  - 当前歌曲  │  │  - 房主权限   │  │  - 输入框       │   │
│  └──────────────┘  └──────────────┘  └─────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ WebSocket (ws://)
┌─────────────────────────────────────────────────────────────┐
│                      后端 (Go + Gin)                        │
├─────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────┐  │
│  │            WebSocket Handler (保留核心)                │  │
│  │  - JOIN_ROOM / LEAVE_ROOM                             │  │
│  │  - CHAT / SYSTEM_MSG                                  │  │
│  │  - PLAY / PAUSE / CHANGE_SONG / SEEK                  │  │
│  │  - TIMELINE_SYNC                                      │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ 房间HTTP API │  │   时间线引擎  │  │    Redis存储    │   │
│  │  - 列表查询  │  │  - 权威时间   │  │  - 房间状态     │   │
│  │  - 创建房间  │  │  - 偏移计算   │  │  - 用户房间映射 │   │
│  └──────────────┘  └──────────────┘  └─────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2.1 核心决策
1. 保留 `server/ws/server.go` 的时间线同步引擎（高精度同步）
2. 简化房间管理，只保留活跃房间（Redis），不持久化到数据库
3. 前端重新设计，简化现有1361行的复杂页面

## 3. 数据模型

### 3.1 房间实体（Redis存储）

```typescript
interface Room {
  id: string;                    // 房间ID
  name: string;                  // 房间名称
  ownerId: string;               // 房主ID
  ownerName: string;             // 房主昵称
  ownerAvatar: string;           // 房主头像
  currentSong: SongInfo | null;  // 当前播放歌曲
  memberCount: number;           // 在线人数
  createdAt: number;             // 创建时间戳
  isActive: boolean;             // 是否活跃
}

interface RoomTimeline {
  songId: string;                // 当前歌曲ID
  startTimestamp: number;        // 服务器起始时间（毫秒）
  paused: boolean;               // 是否暂停
  pausePositionMs: number;       // 暂停位置
  speed: number;                 // 播放速度（默认1.0）
}
```

### 3.2 WebSocket统一包络

所有WebSocket消息使用统一包络结构，便于错误处理、日志追踪和请求-响应关联：

```typescript
interface WsEnvelope<T> {
  type: string;                    // 消息类型
  requestId?: string;              // 用于请求-响应关联（可选）
  roomId?: string;                 // 房间ID（可选）
  clientSendTs?: number;           // 客户端发送时间戳（客户端发时填写）
  serverTs: number;                // 服务器时间戳（毫秒，服务端填写）
  payload: T;                      // 具体消息内容
  error?: {                        // 错误信息（可选）
    code: string;                  // 错误码
    message: string;               // 错误描述
  };
}

// TIME_SYNC 专用消息（用于计算RTT和时钟偏移）
interface TimeSyncRequest {
  type: 'TIME_SYNC_REQ';
  clientSendTs: number;            // T1: 客户端发送时间
}

interface TimeSyncResponse {
  type: 'TIME_SYNC_RES';
  clientSendTs: number;            // T1: 原样返回
  serverTs: number;                // T2: 服务器接收时间
  serverSendTs: number;            // T3: 服务器发送时间
}

// RTT 计算: (T4 - T1) - (T3 - T2)
// 时钟偏移: ((T2 - T1) + (T3 - T4)) / 2

// 错误码定义
enum ErrorCode {
  NOT_ROOM_OWNER = 'NOT_ROOM_OWNER',       // 非房主执行房主操作
  ALREADY_IN_ROOM = 'ALREADY_IN_ROOM',     // 用户已在其他房间
  ROOM_NOT_FOUND = 'ROOM_NOT_FOUND',       // 房间不存在
  ROOM_CLOSED = 'ROOM_CLOSED',             // 房间已关闭
  ROOM_FULL = 'ROOM_FULL',                 // 房间人数已满(500人)
  RATE_LIMITED = 'RATE_LIMITED',           // 操作过于频繁(聊天频控)
  MESSAGE_TOO_LARGE = 'MESSAGE_TOO_LARGE', // 消息大小超限(>8KB)
  INVALID_OPERATION = 'INVALID_OPERATION', // 无效操作
  UNAUTHORIZED = 'UNAUTHORIZED',           // 未登录
}

// 客户端消息示例
interface JoinRoomPayload {
  roomId: string;
}

interface ChatPayload {
  content: string;
}

// 使用示例
const joinRoomMsg: WsEnvelope<JoinRoomPayload> = {
  type: 'JOIN_ROOM',
  requestId: 'uuid-123',
  roomId: 'room-456',
  clientSendTs: Date.now(),  // 客户端填写发送时间
  serverTs: 0,               // 服务端填写，客户端发时填0
  payload: { roomId: 'room-456' }
};
```

### 3.3 Redis Key设计规范

| Key | 类型 | 说明 | TTL |
|-----|------|------|-----|
| `room:{roomId}:meta` | Hash | 房间基本信息（name, ownerId, ownerName, createdAt, isActive） | 房间关闭后30s |
| `room:{roomId}:timeline` | Hash | 房间时间线状态（songId, startTimestamp, paused, pausePositionMs, speed） | 房间关闭后30s |
| `room:{roomId}:members` | Set | 房间成员ID集合 | 房间关闭后30s |
| `user:{userId}:currentRoom` | String | 用户当前所在房间ID（单房间限制关键） | 见下方TTL策略 |
| `rooms:active` | ZSet | 活跃房间列表，score为最后活跃时间戳 | - |

**TTL策略（user:{userId}:currentRoom）：**
| 场景 | TTL | 说明 |
|------|-----|------|
| 正常在线 | 24h（持续心跳续租） | 每次心跳重置TTL |
| 异常断开 | 5min（grace period） | WebSocket异常断开，保留5分钟便于重连恢复 |
| 主动离开 | 立即删除 | 用户点击离开，立即清理 |
| 房间关闭 | 立即删除 | 房主离开导致房间关闭，清理所有成员 |

**清理策略：**
- 房主离开时触发房间关闭流程
- 使用 Lua 脚本原子清理：`room:{id}:*` + `rooms:active` + 所有成员的 `user:{id}:currentRoom`
- 成员心跳超时（60s无心跳）自动踢出

**单房间限制原子实现（Lua）：**
```lua
-- join_room.lua
-- KEYS[1]: user:{userId}:currentRoom    (用户当前房间)
-- KEYS[2]: room:{roomId}:members         (成员ID集合)
-- KEYS[3]: room:{roomId}:member:{userId} (用户信息hash)
-- ARGV[1]: roomId
-- ARGV[2]: userId
-- ARGV[3]: ttlSeconds (过期秒数)
-- ARGV[4]: nowMs      (当前时间戳，用于zset排序)

-- 1. 检查用户是否已在其他房间
local currentRoom = redis.call('GET', KEYS[1])
if currentRoom and currentRoom ~= ARGV[1] then
  return {err = 'ALREADY_IN_ROOM', currentRoom = currentRoom}
end

-- 2. 原子设置用户当前房间（带过期时间）
redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[3])

-- 3. 添加成员ID到集合（只存userId，便于去重）
redis.call('SADD', KEYS[2], ARGV[2])

-- 4. 存储用户信息到独立hash
redis.call('HMSET', KEYS[3], 'joinTime', ARGV[4])
redis.call('EXPIRE', KEYS[3], ARGV[3])

-- 5. 更新房间活跃时间（score用当前时间戳，用于排序）
redis.call('ZADD', 'rooms:active', ARGV[4], ARGV[1])

return {ok = 'SUCCESS'}
```

## 4. 状态机定义

### 4.1 房间状态机

```
                    创建房间
        ┌─────────────────────────────┐
        │                             ▼
   ┌────┴────┐                  ┌──────────┐
   │  IDLE   │◄─────────────────┤  ACTIVE  │
   │ (初始)  │   房主返回/重建  │ (运行中) │
   └────┬────┘                  └────┬─────┘
        │                            │
        │ 房间过期                    │ 房主离开/房间超时
        ▼                            ▼
   ┌─────────┐                  ┌──────────┐
   │ CLOSED  │◄─────────────────┤ CLOSING  │
   │ (关闭)  │   清理完成         │ (关闭中) │
   └─────────┘                  └──────────┘
                                        │
                                        │ 广播关闭通知
                                        ▼
                                   踢出所有成员
```

**状态说明：**
| 状态 | 说明 | 可接受操作 |
|------|------|-----------|
| `IDLE` | 房间未创建或已完全关闭 | 创建房间 |
| `ACTIVE` | 房间正常运行 | 加入、离开、聊天、播放控制 |
| `CLOSING` | 房主离开，正在清理 | 只读查询（短暂状态） |
| `CLOSED` | 房间已关闭，资源已释放 | 无 |

**状态迁移事件：**
| 事件 | 触发条件 | 状态变更 | 回执处理 |
|------|---------|---------|---------|
| `ROOM_CREATE` | 房主创建房间 | IDLE → ACTIVE | 成功：返回房间信息；失败：返回错误码 |
| `OWNER_LEAVE` | 房主断开连接 | ACTIVE → CLOSING | 广播系统消息，开始清理 |
| `ROOM_CLEANUP` | 清理完成 | CLOSING → CLOSED | 无 |
| `ROOM_EXPIRE` | 房间超时（24h无活动） | ACTIVE → CLOSING | 同 OWNER_LEAVE |

### 4.2 播放器状态机

```
              切歌/房主播放
        ┌─────────────────────────────┐
        │                             ▼
   ┌────┴────┐   房主暂停      ┌──────────┐
   │ STOPPED │◄────────────────┤ PLAYING  │
   │ (停止)  │────────────────►│ (播放中) │
   └────┬────┘   播放完成/切歌  └────┬─────┘
        │                            │
        │ 网络缓冲                    │ 房主暂停
        ▼                            ▼
   ┌──────────┐                ┌──────────┐
   │BUFFERING │                │ PAUSED   │
   │ (缓冲中) │                │ (暂停)   │
   └────┬─────┘                └──────────┘
        │                            │
        └──────────缓冲完成──────────►│
                                     │
                                     │ 房主继续
                                     ▼
                              返回 PLAYING
```

**状态说明：**
| 状态 | 说明 | 触发条件 |
|------|------|---------|
| `STOPPED` | 无歌曲或已停止 | 初始状态、播放完成 |
| `PLAYING` | 正在播放 | 房主点击播放/继续 |
| `PAUSED` | 已暂停 | 房主点击暂停 |
| `BUFFERING` | 客户端缓冲中 | 网络延迟、刚切歌 |

**状态同步策略：**
- 服务器维护权威状态，客户端以服务器状态为准
- 房主操作后，服务器先更新状态，再广播给所有成员
- 客户端检测到状态不一致时，立即同步到服务器状态

## 6. 页面设计

### 5.1 直播间列表页

**桌面端：**
- 表格形式展示
- 列：直播间名称、当前播放、人数
- 点击行进入直播间

**移动端适配：**
- 列表项垂直排列，间距加大（手指可点击）
- 隐藏部分信息，只显示：房间名、当前歌曲（单行截断）、人数
- 全宽卡片设计，增加触摸区域

### 5.2 直播间内页

**桌面端（左右分栏）：**
- 左侧：播放器（封面、歌名、进度条、房主控制按钮）
- 右侧：聊天区域（消息列表、输入框）+ 在线成员

**移动端适配（单栏 + 底部播放器）：**
- 使用底部固定播放器，上方显示聊天
- 可滑动切换"聊天"和"成员"标签页
- 房主控制按钮在播放器下方
- 聊天输入框固定在键盘上方

## 7. 关键交互流程

### 5.1 用户进入直播间流程

```
用户点击房间
    │
    ▼
检查登录状态 ──未登录──→ 跳转登录页
    │
    已登录
    ▼
WebSocket连接 ──失败──→ 显示错误提示
    │
    成功
    ▼
发送 JOIN_ROOM 消息
    │
    ▼
接收 ROOM_INFO 消息
    │
    ▼
计算播放偏移量 ──→ 从正确位置开始播放
    │
    ▼
显示聊天和成员列表
```

### 5.2 时间线同步机制

```
服务器维护权威时间线
    │
    ├── 新用户加入 ──→ 发送当前时间线状态
    │                    客户端计算偏移开始播放
    │
    ├── 房主切歌 ────→ 广播新时间线
    │                    所有客户端同步切换
    │
    ├── 房主暂停 ────→ 记录暂停位置
    │                    广播暂停状态
    │
    └── 房主继续 ────→ 重新计算startTimestamp
                         广播继续状态
```

### 5.2.1 时间同步校正阈值策略

客户端定期与服务器同步时间，根据偏差大小采取不同校正策略：

| 偏差范围 | 处理策略 | 用户体验 |
|---------|---------|---------|
| `<= 300ms` | 不处理，容忍范围内 | 用户无感知 |
| `300ms ~ 1500ms` | 平滑追帧：微调 `playbackRate`（0.95x ~ 1.05x） | 轻微变速，2-3秒内追上 |
| `> 1500ms` | 直接 seek 校正 | 可能有轻微跳跃感 |

**实现细节：**
```typescript
class TimelineSync {
  private readonly TOLERABLE_DRIFT = 300;      // 可容忍偏差（ms）
  private readonly SMOOTH_THRESHOLD = 1500;    // 平滑追帧上限（ms）
  private readonly CORRECTION_RATE = 0.05;     // 追帧速率调整量

  correctDrift(currentTime: number, serverTime: number): void {
    const drift = serverTime - currentTime;
    const absDrift = Math.abs(drift);

    if (absDrift <= this.TOLERABLE_DRIFT) {
      // 不处理
      return;
    }

    if (absDrift <= this.SMOOTH_THRESHOLD) {
      // 平滑追帧
      const rate = drift > 0 
        ? 1 + this.CORRECTION_RATE   // 落后，加速
        : 1 - this.CORRECTION_RATE;  // 超前，减速
      this.player.playbackRate = rate;
      
      // 2秒后恢复正常速度
      setTimeout(() => this.player.playbackRate = 1, 2000);
    } else {
      // 直接校正
      this.player.currentTime = serverTime / 1000;
    }
  }
}
```

**同步频率：**
- 初始连接：立即同步
- 正常播放：每30秒同步一次
- 切歌/暂停/继续：立即同步
- 检测到偏差>1秒时：立即同步

### 5.3 权限控制

| 操作 | 房主 | 普通听众 |
|------|------|----------|
| 播放/暂停 | ✅ | ❌ |
| 切歌 | ✅ | ❌ |
| 进度跳转 | ✅ | ❌ |
| 发送聊天 | ✅ | ✅ |
| 查看成员 | ✅ | ✅ |
| 离开房间 | ✅ | ✅ |

### 5.4 单房间限制的并发处理

**问题场景：** 用户在两台设备上同时点击加入不同房间。

**解决方案 - 原子占位：**

```lua
-- join_room.lua
-- KEYS[1]: user:{userId}:currentRoom    (用户当前房间)
-- KEYS[2]: room:{roomId}:members         (成员ID集合)
-- KEYS[3]: room:{roomId}:member:{userId} (用户信息hash)
-- ARGV[1]: roomId
-- ARGV[2]: userId
-- ARGV[3]: ttlSeconds (过期秒数)
-- ARGV[4]: nowMs      (当前时间戳，用于zset排序)

-- 返回值数组：[status, data/error]
-- status: 0=成功, 1=业务错误, 2=系统错误

-- 1. 检查用户是否已在其他房间
local currentRoom = redis.call('GET', KEYS[1])
if currentRoom and currentRoom ~= ARGV[1] then
  -- 返回错误: 已在其他房间
  return {1, cjson.encode({code='ALREADY_IN_ROOM', currentRoom=currentRoom})}
end

-- 2. 原子设置用户当前房间（带过期时间）
redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[3])

-- 3. 添加成员ID到集合（只存userId，便于去重）
redis.call('SADD', KEYS[2], ARGV[2])

-- 4. 存储用户加入时间到独立hash (使用 HSET，Redis新版本推荐)
redis.call('HSET', KEYS[3], 'joinTime', ARGV[4])
redis.call('EXPIRE', KEYS[3], ARGV[3])

-- 5. 更新房间活跃时间（score用当前时间戳，用于排序）
redis.call('ZADD', 'rooms:active', ARGV[4], ARGV[1])

-- 返回成功
return {0, cjson.encode({code='SUCCESS', roomId=ARGV[1]})}
```

**客户端处理：**
| 场景 | 服务端行为 | 客户端提示 |
|------|-----------|-----------|
| 用户已在房间A，尝试加入房间A | 幂等处理，返回成功 | 正常进入 |
| 用户已在房间A，尝试加入房间B | 返回 `ALREADY_IN_ROOM` | "您已在其他直播间，请先离开" |
| WebSocket异常断开 | 保留 `user:{id}:currentRoom` 5分钟 | 重连后自动恢复 |
| 主动离开房间 | 删除 `user:{id}:currentRoom` | 返回列表页 |
| 心跳超时（60s） | 踢出成员，清理房间关系 | 提示"连接超时" |
| 房主踢出用户 | 删除关系，强制断开WebSocket | 提示"您已被移出房间" |

**断线恢复策略：**
1. 用户重新连接 WebSocket
2. 检查 `user:{id}:currentRoom` 是否存在
3. 存在则自动发送 `REJOIN_ROOM` 消息
4. 服务端验证后返回当前房间状态

**补充：room:{roomId}:member:{userId} 生命周期管理**
| 触发条件 | 操作 |
|---------|------|
| 用户主动离开 | 立即 `DEL room:{id}:member:{uid}` |
| 用户被踢出 | 立即删除，并广播KICK消息 |
| 心跳超时（60s） | 踢出流程中删除 |
| 房间关闭 | Lua脚本批量删除 `room:{id}:member:*` |
| 定期清理 | 扫描孤儿key（member存在但user:{id}:currentRoom不存在） |

### 5.5 幂等规则表

| 消息类型 | 是否幂等 | 重复请求处理策略 | 说明 |
|---------|---------|-----------------|------|
| `JOIN_ROOM` | ✅ 是 | 已在该房间：返回当前房间状态 | 用户重复点击或重连时触发 |
| `LEAVE_ROOM` | ✅ 是 | 已不在房间：返回成功 | 多次发送或断线后触发 |
| `CHAT` | ❌ 否 | 正常处理（消息ID去重） | 需业务层去重（如用clientSendTs） |
| `HEARTBEAT` | ✅ 是 | 更新心跳时间 | 多次发送无影响 |
| `PLAY` | ⚠️ 条件 | 已是PLAYING：返回成功 | 状态不变时幂等 |
| `PAUSE` | ⚠️ 条件 | 已是PAUSED：返回成功 | 状态不变时幂等 |
| `SEEK` | ⚠️ 条件 | 同一位置：返回成功 | 位置相同则幂等 |
| `CHANGE_SONG` | ❌ 否 | 正常处理（切到同一首也处理） | 房主可能想从头播放 |
| `TIME_SYNC_REQ` | ✅ 是 | 返回时间同步响应 | 多次请求无影响 |

**幂等实现规范：**

1. **requestId 去重缓存**
   - Key: `req:{userId}:{requestId}`
   - Value: 处理结果（JSON）
   - TTL: 300s
   - 示例: `req:10086:uuid-abc-123`

2. **状态类操作（PLAY/PAUSE/SEEK）**
   - 先检查当前状态，状态相同直接返回成功
   - 不依赖requestId缓存

3. **CHAT消息去重**
   - 客户端生成 `messageId`（建议: `userId + clientSendTs + seq`）
   - 服务端缓存 `msg:{messageId}`，TTL=60s
   - 避免同毫秒冲突（seq为客户端自增序号）

### 5.6 安全底线

#### 5.6.1 聊天消息限制

| 限制项 | 值 | 说明 |
|-------|---|------|
| 消息长度 | 1 ~ 300字符 | 空消息拒绝，超长截断 |
| 发送频率 | 2条/秒 | 超过则提示"发送太频繁" |
| 相同内容 | 5分钟内禁止重复 | 防刷屏 |
| 敏感词 | 需接入敏感词过滤 | 可先用本地简单规则 |

#### 5.6.2 WebSocket消息大小限制

| 限制项 | 值 | 超限处理 |
|-------|---|---------|
| 单消息大小 | 8KB | 关闭连接并记录日志 |
| 消息队列 | 100条/连接 | 超限丢弃旧消息 |
| 连接数/用户 | 3个 | 踢出最早连接 |
| 连接数/房间 | 500人 | 加入时返回 ROOM_FULL |

### 5.7 错误处理

| 场景 | 处理方式 |
|------|----------|
| WebSocket断开 | 显示"连接断开，正在重连..."，3次失败后提示手动刷新 |
| 房主离开 | 房间自动关闭，所有用户收到"直播间已结束"提示，返回列表页 |
| 同步偏差过大 | 客户端自动校正，用户无感知 |
| 加入房间失败（如已在其他房间） | 提示"请先离开当前直播间" |

## 6. 技术要点

### 6.1 复用部分
- `server/ws/server.go` - 时间线同步引擎
- `src/core/realtime/timeline.ts` - 客户端时间线计算
- WebSocket连接管理逻辑

### 6.2 新建/修改部分
- `src/views/ListenTogether/index.vue` - 重新设计直播间页面
- `src/views/ListenTogether/RoomList.vue` - 新增直播间列表页
- `src/components/Chat/` - 聊天组件
- 路由配置 - 添加直播间路由

## 8. 待办事项

- [ ] 重新设计直播间列表页
- [ ] 重新设计直播间内页（支持移动端）
- [ ] 实现聊天组件（文字+系统消息）
- [ ] 调整权限控制逻辑
- [ ] 添加移动端响应式适配
- [ ] 测试时间线同步精度
- [ ] 测试WebSocket重连机制
