package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
)

// WebSocket 消息类型定义
const (
	MsgTypeHello        = "HELLO"         // 客户端连接初始化
	MsgTypeTimeSyncReq  = "TIME_SYNC_REQ" // 客户端请求时间同步
	MsgTypeTimeSyncRes  = "TIME_SYNC_RES" // 服务器返回时间同步响应
	MsgTypeHeartbeat    = "HEARTBEAT"     // 心跳包，用于保活和RTT测量
	MsgTypeChat         = "CHAT"          // 聊天消息
	MsgTypeJoinRoom     = "JOIN_ROOM"     // 加入房间
	MsgTypeLeaveRoom    = "LEAVE_ROOM"    // 离开房间
	MsgTypeRoomMembers  = "ROOM_MEMBERS"  // 房间成员列表
	MsgTypeRoomInfo     = "ROOM_INFO"     // 房间信息（包含房主、时间线等）
	MsgTypeGetRoomList  = "GET_ROOM_LIST" // 获取房间列表请求
	MsgTypeRoomListRes  = "ROOM_LIST_RES" // 房间列表响应
	MsgTypeTimelineInit = "TIMELINE_INIT"
	MsgTypePlay         = "PLAY"
	MsgTypePause        = "PAUSE"
	MsgTypeSeek         = "SEEK"
	MsgTypeChangeSong   = "CHANGE_SONG"
	MsgTypeSetSpeed     = "SET_SPEED"
)

// RoomTimeline 房间权威时间线
type RoomTimeline struct {
	SongID          string  `json:"song_id"`           // 当前播放的歌曲唯一ID
	StartTimestamp  int64   `json:"start_timestamp"`   // 服务器时间：这首歌“从0ms开始播放”的时间点
	Paused          bool    `json:"paused"`            // 当前是否处于暂停状态
	PausePositionMs float64 `json:"pause_position_ms"` // 暂停时，歌曲精确停在的位置
	Speed           float64 `json:"speed"`             // 播放速度，默认 1.0
}

// ApplyPlay 应用播放操作
func (t *RoomTimeline) ApplyPlay(serverTime int64) {
	t.StartTimestamp = serverTime - int64(t.PausePositionMs)
	t.Paused = false
}

// ApplyPause 应用暂停操作
func (t *RoomTimeline) ApplyPause(serverTime int64, positionMs float64) {
	t.PausePositionMs = positionMs
	t.Paused = true
}

// ApplySeek 应用跳转操作
func (t *RoomTimeline) ApplySeek(serverTime int64, positionMs float64) {
	t.StartTimestamp = serverTime - int64(positionMs)
	t.Paused = false
}

// ApplyChangeSong 应用切歌操作
func (t *RoomTimeline) ApplyChangeSong(serverTime int64, songId string) {
	t.SongID = songId
	t.StartTimestamp = serverTime
	t.PausePositionMs = 0
	t.Paused = false
}

// ApplySetSpeed 应用倍速操作
func (t *RoomTimeline) ApplySetSpeed(serverTime int64, newSpeed float64) {
	if !t.Paused {
		oldSpeed := t.Speed
		if oldSpeed == 0 {
			oldSpeed = 1.0
		}
		// 计算当前播放了多少时间（基于旧速度）
		// pos = (server_time - start_timestamp) * old_speed
		pos := float64(serverTime-t.StartTimestamp) * oldSpeed

		// 反推新的 start_timestamp，使得在新速度下当前位置不变
		// new_start = server_time - pos / new_speed
		t.StartTimestamp = serverTime - int64(pos/newSpeed)
	}
	t.Speed = newSpeed
}

// SaveTimeline 保存时间线到 Redis
func (s *WSServer) SaveTimeline(roomId string, timeline *RoomTimeline) error {
	jsonBytes, err := json.Marshal(timeline)
	if err != nil {
		return err
	}
	return s.RDB.Set(context.Background(), fmt.Sprintf("room:%s:timeline", roomId), jsonBytes, 0).Err()
}

// LoadTimeline 从 Redis 加载时间线
func (s *WSServer) LoadTimeline(roomId string) (*RoomTimeline, error) {
	jsonStr, err := s.RDB.Get(context.Background(), fmt.Sprintf("room:%s:timeline", roomId)).Result()
	if err != nil {
		return nil, err
	}
	var timeline RoomTimeline
	err = json.Unmarshal([]byte(jsonStr), &timeline)
	if err != nil {
		return nil, err
	}
	return &timeline, nil
}

// AddMember 添加成员到房间
func (s *WSServer) AddMember(roomId string, userId string, userJson string) error {
	return s.RDB.HSet(context.Background(), fmt.Sprintf("room:%s:members", roomId), userId, userJson).Err()
}

// RemoveMember 从房间移除成员
func (s *WSServer) RemoveMember(roomId string, userId string) error {
	return s.RDB.HDel(context.Background(), fmt.Sprintf("room:%s:members", roomId), userId).Err()
}

// GetMembers 获取房间成员列表
func (s *WSServer) GetMembers(roomId string) ([]interface{}, error) {
	membersMap, err := s.RDB.HGetAll(context.Background(), fmt.Sprintf("room:%s:members", roomId)).Result()
	if err != nil {
		return nil, err
	}
	var members []interface{}
	for _, jsonStr := range membersMap {
		var member interface{}
		if err := json.Unmarshal([]byte(jsonStr), &member); err == nil {
			members = append(members, member)
		}
	}
	return members, nil
}

// WSMessage 定义 WebSocket 消息结构
type WSMessage struct {
	Type    string      `json:"type"`    // 消息类型
	Payload interface{} `json:"payload"` // 消息负载
}

// WSServer WebSocket 服务端结构体
type WSServer struct {
	M   *melody.Melody // Melody 实例，用于处理 WebSocket 连接和消息
	RDB *redis.Client  // Redis 客户端
}

// leaveRoom 处理用户离开房间的逻辑
func (wsServer *WSServer) leaveRoom(s *melody.Session, roomIdStr string, userObj interface{}) {
	userId, _ := s.Get("userId")
	userIdStr := fmt.Sprintf("%v", userId)

	// 从 Redis 中移除用户
	wsServer.RDB.HDel(context.Background(), fmt.Sprintf("room:%s:members", roomIdStr), userIdStr)

	// 检查是否是房主离开
	ownerKey := fmt.Sprintf("room:%s:owner", roomIdStr)
	currentOwner, _ := wsServer.RDB.Get(context.Background(), ownerKey).Result()
	if currentOwner == userIdStr {
		// 房主离开，删除房主键
		wsServer.RDB.Del(context.Background(), ownerKey)

		// 尝试移交房主给房间内其他人
		// 获取所有成员
		membersMap, _ := wsServer.RDB.HGetAll(context.Background(), fmt.Sprintf("room:%s:members", roomIdStr)).Result()
		if len(membersMap) > 0 {
			// 随机选取一个新房主
			for newOwnerId := range membersMap {
				wsServer.RDB.Set(context.Background(), ownerKey, newOwnerId, 0)
				break
			}
		} else {
			// 房间没人了，清理时间线和活跃房间列表
			wsServer.RDB.Del(context.Background(), fmt.Sprintf("room:%s:timeline", roomIdStr))
			wsServer.RDB.SRem(context.Background(), "active_rooms", roomIdStr)
		}
	}

	// 构造离开房间消息
	leaveMsg := WSMessage{
		Type: MsgTypeLeaveRoom,
		Payload: map[string]interface{}{
			"roomId": roomIdStr,
			"user":   userObj,
		},
	}

	// 广播离开消息
	if msgBytes, err := json.Marshal(leaveMsg); err == nil {
		wsServer.M.Broadcast(msgBytes)
	}

	// 广播更新后的成员列表和房间信息
	wsServer.BroadcastRoomMembers(roomIdStr)
	wsServer.BroadcastRoomInfo(roomIdStr)
}

// NewWSServer 创建并初始化一个新的 WebSocket 服务端
func NewWSServer(rdb *redis.Client) *WSServer {
	m := melody.New()
	wsServer := &WSServer{M: m, RDB: rdb}

	// 处理 WebSocket 连接建立事件
	m.HandleConnect(func(s *melody.Session) {
		// 这里可以记录连接日志或初始化会话状态
	})

	// 处理 WebSocket 连接断开事件
	m.HandleDisconnect(func(s *melody.Session) {
		// 获取用户信息和房间列表
		userId, _ := s.Get("userId")
		// roomId, _ := s.Get("roomId") // 废弃单房间模式
		roomsInterface, exists := s.Get("rooms")
		userJson, _ := s.Get("userJson")

		if userId != nil && exists {
			rooms := roomsInterface.([]string)

			var userObj interface{}
			if userJson != nil {
				json.Unmarshal([]byte(userJson.(string)), &userObj)
			}

			// 遍历用户所在的所有房间，执行离开逻辑
			for _, roomIdStr := range rooms {
				wsServer.leaveRoom(s, roomIdStr, userObj)
			}
		}
	})

	// 处理接收到的 WebSocket 消息
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var incoming WSMessage
		// 解析 JSON 消息
		if err := json.Unmarshal(msg, &incoming); err != nil {
			return
		}

		switch incoming.Type {
		case MsgTypeHello:
			// 客户端初始化处理
			// s.Set("info", incoming.Payload) // 可在此处存储用户信息到 session

			// 发送 HELLO 响应 (HELLO_ACK)
			response := WSMessage{
				Type: "HELLO", // 回显 HELLO 类型，前端据此确认连接成功
				Payload: map[string]string{
					"status": "connected",
					"msg":    "Welcome to Listen Together",
				},
			}
			if respBytes, err := json.Marshal(response); err == nil {
				s.Write(respBytes)
			}

		case MsgTypeTimeSyncReq:
			// 处理时间同步请求
			response := WSMessage{
				Type: MsgTypeTimeSyncRes,
				Payload: map[string]int64{
					"server_now": time.Now().UnixMilli(), // 返回服务器当前时间戳 (毫秒)
					"client_req": 0,                      // 可选：回显客户端请求时间戳
				},
			}

			// 如果负载中包含 client_now，可以用于计算 RTT (Round Trip Time)
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				if clientTime, ok := payloadMap["client_now"]; ok {
					// 这里可以选择回显 clientTime，或者直接让客户端根据发送时间自行计算
					_ = clientTime
				}
			}

			if respBytes, err := json.Marshal(response); err == nil {
				s.Write(respBytes)
			}

		case MsgTypeHeartbeat:
			// 处理心跳包：原样返回消息，用于客户端测量延迟 (RTT)
			response := WSMessage{
				Type:    MsgTypeHeartbeat,
				Payload: incoming.Payload,
			}
			if respBytes, err := json.Marshal(response); err == nil {
				s.Write(respBytes)
			}

		case MsgTypeChat:
			// 处理聊天消息：只有在房间内的用户才能发送消息
			// payload 必须包含 targetRoomId
			payloadMap, ok := incoming.Payload.(map[string]interface{})
			if !ok {
				return
			}
			targetRoomId, ok := payloadMap["targetRoomId"].(string)
			if !ok {
				// 兼容旧逻辑：尝试从 session 获取（如果只在一个房间）
				// 但现在支持多房间，最好强制要求 targetRoomId
				// 这里暂时返回，要求前端必须传
				return
			}

			userId, exists2 := s.Get("userId")
			roomsInterface, exists := s.Get("rooms")

			if exists && exists2 {
				rooms := roomsInterface.([]string)
				userIdStr := fmt.Sprintf("%v", userId)

				// 检查用户 session 是否包含该房间
				inSession := false
				for _, r := range rooms {
					if r == targetRoomId {
						inSession = true
						break
					}
				}

				if inSession {
					// 双重检查 Redis (可选，但更安全)
					inRoom, err := wsServer.RDB.HExists(context.Background(), fmt.Sprintf("room:%s:members", targetRoomId), userIdStr).Result()
					if err == nil && inRoom {
						// 广播消息
						m.Broadcast(msg)
					}
				}
			}

		case MsgTypeGetRoomList:
			// 获取房间列表
			activeRooms, err := wsServer.RDB.SMembers(context.Background(), "active_rooms").Result()
			if err != nil {
				activeRooms = []string{}
			}

			// 构造房间列表响应，包含人数
			var roomList []map[string]interface{}
			for _, roomId := range activeRooms {
				count, _ := wsServer.RDB.HLen(context.Background(), fmt.Sprintf("room:%s:members", roomId)).Result()
				roomList = append(roomList, map[string]interface{}{
					"id":    roomId,
					"count": count,
				})
			}

			response := WSMessage{
				Type:    MsgTypeRoomListRes,
				Payload: roomList,
			}
			if respBytes, err := json.Marshal(response); err == nil {
				s.Write(respBytes)
			}

		case MsgTypeJoinRoom:
			// 加入房间逻辑
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				if roomId, ok := payloadMap["roomId"].(string); ok {
					if user, ok := payloadMap["user"].(map[string]interface{}); ok {
						userId := user["id"]
						// 存储 Session 信息 (多房间支持)
						// s.Set("roomId", roomId) // 旧逻辑

						var rooms []string
						if existingRooms, exists := s.Get("rooms"); exists {
							rooms = existingRooms.([]string)
						} else {
							rooms = []string{}
						}

						// 单房间模式：如果已经在其他房间，先离开
						for _, r := range rooms {
							if r != roomId {
								wsServer.leaveRoom(s, r, user)
							}
						}

						// 重置房间列表，仅包含当前新加入的房间
						rooms = []string{roomId}
						s.Set("rooms", rooms)

						s.Set("userId", userId)

						userJson, _ := json.Marshal(user)
						s.Set("userJson", string(userJson))

						// 添加到活跃房间列表
						wsServer.RDB.SAdd(context.Background(), "active_rooms", roomId)

						// 添加到 Redis 成员列表
						userIdStr := fmt.Sprintf("%v", userId)
						wsServer.RDB.HSet(context.Background(), fmt.Sprintf("room:%s:members", roomId), userIdStr, string(userJson))

						// 检查是否需要设置房主
						ownerKey := fmt.Sprintf("room:%s:owner", roomId)
						// 使用 SetNX 尝试设置房主，如果键不存在则设置成功
						setOwner, _ := wsServer.RDB.SetNX(context.Background(), ownerKey, userIdStr, 0).Result()

						// 广播加入消息
						m.Broadcast(msg)

						// 广播更新后的成员列表
						wsServer.BroadcastRoomMembers(roomId)

						// 广播房间信息（包含房主）
						wsServer.BroadcastRoomInfo(roomId)

						// 如果是房主（新建房间或接管），且需要初始化时间线
						if setOwner {
							// 初始化空的时间线
							timeline := RoomTimeline{
								SongID:          "",
								StartTimestamp:  0,
								Paused:          true,
								PausePositionMs: 0,
								Speed:           1.0,
							}
							timelineJson, _ := json.Marshal(timeline)
							wsServer.RDB.Set(context.Background(), fmt.Sprintf("room:%s:timeline", roomId), timelineJson, 0)
							// 再次广播房间信息以包含时间线
							wsServer.BroadcastRoomInfo(roomId)
						}
					}
				}
			}

		case MsgTypeLeaveRoom:
			// 离开房间逻辑 (主动离开指定房间)
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				if roomId, ok := payloadMap["roomId"].(string); ok {
					if user, ok := payloadMap["user"].(map[string]interface{}); ok {
						userId := user["id"]
						userIdStr := fmt.Sprintf("%v", userId)

						// 从 Session 的房间列表中移除
						if existingRooms, exists := s.Get("rooms"); exists {
							rooms := existingRooms.([]string)
							newRooms := []string{}
							for _, r := range rooms {
								if r != roomId {
									newRooms = append(newRooms, r)
								}
							}
							s.Set("rooms", newRooms)
						}

						// 从 Redis 移除
						wsServer.RDB.HDel(context.Background(), fmt.Sprintf("room:%s:members", roomId), userIdStr)

						// 检查是否是房主离开
						ownerKey := fmt.Sprintf("room:%s:owner", roomId)
						currentOwner, _ := wsServer.RDB.Get(context.Background(), ownerKey).Result()
						if currentOwner == userIdStr {
							// 房主离开，删除房主键
							wsServer.RDB.Del(context.Background(), ownerKey)

							// 尝试移交房主给房间内其他人
							membersMap, _ := wsServer.RDB.HGetAll(context.Background(), fmt.Sprintf("room:%s:members", roomId)).Result()
							if len(membersMap) > 0 {
								for newOwnerId := range membersMap {
									wsServer.RDB.Set(context.Background(), ownerKey, newOwnerId, 0)
									break
								}
							} else {
								// 房间没人了
								wsServer.RDB.Del(context.Background(), fmt.Sprintf("room:%s:timeline", roomId))
								wsServer.RDB.SRem(context.Background(), "active_rooms", roomId)
							}
						}

						// 广播离开消息
						m.Broadcast(msg)

						// 广播更新后的成员列表
						wsServer.BroadcastRoomMembers(roomId)

						// 广播更新后的房间信息
						wsServer.BroadcastRoomInfo(roomId)
					}
				}
			}

		case MsgTypePlay:
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				roomId, _ := payloadMap["roomId"].(string)
				if roomId == "" {
					return
				}
				// 权限检查：仅房主可操作
				userId, _ := s.Get("userId")
				userIdStr := fmt.Sprintf("%v", userId)
				ownerId, _ := wsServer.RDB.Get(context.Background(), fmt.Sprintf("room:%s:owner", roomId)).Result()
				if userIdStr != ownerId {
					return
				}

				serverTime := time.Now().UnixMilli()
				timeline, err := wsServer.LoadTimeline(roomId)
				if err != nil || timeline == nil {
					return
				}

				timeline.ApplyPlay(serverTime)

				wsServer.SaveTimeline(roomId, timeline)

				resp := WSMessage{
					Type: MsgTypePlay,
					Payload: map[string]interface{}{
						"event":       "PLAY",
						"server_time": serverTime,
						"roomId":      roomId,
					},
				}
				if b, err := json.Marshal(resp); err == nil {
					wsServer.M.Broadcast(b)
				}
				wsServer.BroadcastRoomInfo(roomId)
			}

		case MsgTypePause:
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				roomId, _ := payloadMap["roomId"].(string)
				if roomId == "" {
					return
				}
				userId, _ := s.Get("userId")
				userIdStr := fmt.Sprintf("%v", userId)
				ownerId, _ := wsServer.RDB.Get(context.Background(), fmt.Sprintf("room:%s:owner", roomId)).Result()
				if userIdStr != ownerId {
					return
				}

				serverTime := time.Now().UnixMilli()
				positionMs, _ := payloadMap["position_ms"].(float64)

				timeline, err := wsServer.LoadTimeline(roomId)
				if err != nil || timeline == nil {
					return
				}

				timeline.ApplyPause(serverTime, positionMs)

				wsServer.SaveTimeline(roomId, timeline)

				resp := WSMessage{
					Type: MsgTypePause,
					Payload: map[string]interface{}{
						"event":       "PAUSE",
						"server_time": serverTime,
						"position_ms": positionMs,
						"roomId":      roomId,
					},
				}
				if b, err := json.Marshal(resp); err == nil {
					wsServer.M.Broadcast(b)
				}
				wsServer.BroadcastRoomInfo(roomId)
			}

		case MsgTypeSeek:
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				roomId, _ := payloadMap["roomId"].(string)
				if roomId == "" {
					return
				}
				userId, _ := s.Get("userId")
				userIdStr := fmt.Sprintf("%v", userId)
				ownerId, _ := wsServer.RDB.Get(context.Background(), fmt.Sprintf("room:%s:owner", roomId)).Result()
				if userIdStr != ownerId {
					return
				}

				serverTime := time.Now().UnixMilli()
				positionMs, _ := payloadMap["position_ms"].(float64)

				timeline, err := wsServer.LoadTimeline(roomId)
				if err != nil || timeline == nil {
					return
				}

				timeline.ApplySeek(serverTime, positionMs)

				wsServer.SaveTimeline(roomId, timeline)

				resp := WSMessage{
					Type: MsgTypeSeek,
					Payload: map[string]interface{}{
						"event":       "SEEK",
						"server_time": serverTime,
						"position_ms": positionMs,
						"roomId":      roomId,
					},
				}
				if b, err := json.Marshal(resp); err == nil {
					wsServer.M.Broadcast(b)
				}
				wsServer.BroadcastRoomInfo(roomId)
			}

		case MsgTypeChangeSong:
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				roomId, _ := payloadMap["roomId"].(string)
				songId, _ := payloadMap["song_id"].(string)
				if roomId == "" {
					return
				}
				userId, _ := s.Get("userId")
				userIdStr := fmt.Sprintf("%v", userId)
				ownerId, _ := wsServer.RDB.Get(context.Background(), fmt.Sprintf("room:%s:owner", roomId)).Result()
				if userIdStr != ownerId {
					return
				}

				serverTime := time.Now().UnixMilli()

				timeline, err := wsServer.LoadTimeline(roomId)
				if err != nil || timeline == nil {
					return
				}

				timeline.ApplyChangeSong(serverTime, songId)

				wsServer.SaveTimeline(roomId, timeline)

				resp := WSMessage{
					Type: MsgTypeChangeSong,
					Payload: map[string]interface{}{
						"event":       "CHANGE_SONG",
						"song_id":     songId,
						"server_time": serverTime,
						"roomId":      roomId,
					},
				}
				if b, err := json.Marshal(resp); err == nil {
					wsServer.M.Broadcast(b)
				}
				wsServer.BroadcastRoomInfo(roomId)
			}

		case MsgTypeSetSpeed:
			if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
				roomId, _ := payloadMap["roomId"].(string)
				newSpeed, _ := payloadMap["speed"].(float64)
				if roomId == "" {
					return
				}
				userId, _ := s.Get("userId")
				userIdStr := fmt.Sprintf("%v", userId)
				ownerId, _ := wsServer.RDB.Get(context.Background(), fmt.Sprintf("room:%s:owner", roomId)).Result()
				if userIdStr != ownerId {
					return
				}

				serverTime := time.Now().UnixMilli()

				timeline, err := wsServer.LoadTimeline(roomId)
				if err != nil || timeline == nil {
					return
				}

				timeline.ApplySetSpeed(serverTime, newSpeed)

				wsServer.SaveTimeline(roomId, timeline)

				resp := WSMessage{
					Type: MsgTypeSetSpeed,
					Payload: map[string]interface{}{
						"event":       "SET_SPEED",
						"speed":       newSpeed,
						"server_time": serverTime,
						"roomId":      roomId,
					},
				}
				if b, err := json.Marshal(resp); err == nil {
					wsServer.M.Broadcast(b)
				}
				wsServer.BroadcastRoomInfo(roomId)
			}
		}
	})

	return wsServer
}

// BroadcastRoomInfo 广播房间信息（房主、时间线等）
func (s *WSServer) BroadcastRoomInfo(roomId string) {
	// 获取房主 ID
	ownerId, _ := s.RDB.Get(context.Background(), fmt.Sprintf("room:%s:owner", roomId)).Result()

	// 获取时间线
	timeline, _ := s.LoadTimeline(roomId)

	msg := WSMessage{
		Type: MsgTypeRoomInfo,
		Payload: map[string]interface{}{
			"roomId":      roomId,
			"ownerId":     ownerId,
			"timeline":    timeline,
			"server_time": time.Now().UnixMilli(),
		},
	}

	if msgBytes, err := json.Marshal(msg); err == nil {
		s.M.Broadcast(msgBytes)
	}
}

// BroadcastRoomMembers 广播房间成员列表
func (s *WSServer) BroadcastRoomMembers(roomId string) {
	members, err := s.GetMembers(roomId)
	if err != nil {
		return
	}

	msg := WSMessage{
		Type: MsgTypeRoomMembers,
		Payload: map[string]interface{}{
			"roomId":  roomId,
			"members": members,
		},
	}

	if msgBytes, err := json.Marshal(msg); err == nil {
		s.M.Broadcast(msgBytes)
	}
}
