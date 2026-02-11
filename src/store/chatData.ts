import { defineStore } from "pinia";
import { socket, MsgType, type WSMessage } from "@/core/realtime/socket";
import { useUserDataStore } from "./userData";

interface RoomUser {
  id: number | string;
  nickname: string;
  avatarUrl: string;
}

interface ChatMessage {
  sender: string;
  content: string;
  avatarUrl?: string;
  isMine: boolean;
  type?: 'system' | 'chat';
}

interface RoomData {
  messages: ChatMessage[];
  users: RoomUser[];
  ownerId: string | number;
}

export const useChatDataStore = defineStore("chatData", {
  state: () => ({
    availableRooms: [] as any[],
    joinedRooms: [] as string[],
    currentRoomId: '',
    roomData: {} as Record<string, RoomData>,
    heartbeatInterval: null as any,
    isListening: false, // 是否已经注册了 socket 监听
  }),
  getters: {
    currentRoomMessages: (state) => {
      if (!state.currentRoomId) return [];
      return state.roomData[state.currentRoomId]?.messages || [];
    },
    currentRoomUsers: (state) => {
      if (!state.currentRoomId) return [];
      return state.roomData[state.currentRoomId]?.users || [];
    },
    currentRoomOwnerId: (state) => {
      if (!state.currentRoomId) return '';
      return state.roomData[state.currentRoomId]?.ownerId || '';
    },
  },
  actions: {
    initRoomData(roomId: string) {
      if (!this.roomData[roomId]) {
        this.roomData[roomId] = {
          messages: [],
          users: [],
          ownerId: ''
        };
      }
    },
    
    // 启动监听
    startListen() {
      if (this.isListening) return;
      socket.onMessage(this.handleMessage.bind(this));
      this.isListening = true;
      
      if (!socket.isConnected.value) {
        socket.connect();
      }
      this.startHeartbeatLoop();
    },

    // 停止监听 (完全断开)
    stopListen() {
      socket.offMessage(this.handleMessage.bind(this)); 
      
      socket.disconnect();
      if (this.heartbeatInterval) {
        clearInterval(this.heartbeatInterval);
        this.heartbeatInterval = null;
      }
      this.isListening = false;
      this.joinedRooms = [];
      this.roomData = {};
      this.currentRoomId = '';
    },

    handleMessage(msg: WSMessage) {
      const user = useUserDataStore();
      
      switch (msg.type) {
        case MsgType.ROOM_LIST_RES:
          this.availableRooms = Array.isArray(msg.payload) ? msg.payload : [];
          break;

        case MsgType.CHAT:
          if (msg.payload && msg.payload.targetRoomId) {
            const rid = msg.payload.targetRoomId;
            this.initRoomData(rid);
            this.roomData[rid].messages.push({
              sender: msg.payload.sender || 'Unknown',
              content: msg.payload.content,
              avatarUrl: msg.payload.avatarUrl,
              isMine: msg.payload.sender === user.getUserData.nickname,
              type: 'chat'
            });
          }
          break;

        case MsgType.ROOM_MEMBERS:
          if (msg.payload && msg.payload.roomId && Array.isArray(msg.payload.members)) {
            const rid = msg.payload.roomId;
            this.initRoomData(rid);
            this.roomData[rid].users = msg.payload.members;
          }
          break;

        case MsgType.ROOM_INFO:
          if (msg.payload && msg.payload.roomId) {
            const rid = msg.payload.roomId;
            this.initRoomData(rid);
            this.roomData[rid].ownerId = msg.payload.ownerId;
          }
          break;

        case MsgType.JOIN_ROOM:
          if (msg.payload && msg.payload.roomId) {
            const rid = msg.payload.roomId;
            this.initRoomData(rid);
            this.roomData[rid].messages.push({
              type: 'system',
              content: `${msg.payload.user?.nickname || 'Someone'} 加入了群聊`,
              sender: 'System',
              isMine: false
            });
          }
          break;

        case MsgType.LEAVE_ROOM:
          if (msg.payload && msg.payload.roomId) {
            const rid = msg.payload.roomId;
            const leavingUser = msg.payload.user;
            const myId = user.getUserData.userId;

            // 检查是否是自己离开
            if (leavingUser && String(leavingUser.id) === String(myId)) {
               this.joinedRooms = this.joinedRooms.filter(id => id !== rid);
               delete this.roomData[rid];
               if (this.currentRoomId === rid) {
                   this.currentRoomId = '';
               }
               // if (window.$message) window.$message.info(`已退出房间: ${rid}`);
            } else {
                // 别人离开
                this.initRoomData(rid);
                this.roomData[rid].messages.push({
                  type: 'system',
                  content: `${leavingUser?.nickname || 'Someone'} 离开了群聊`,
                  sender: 'System',
                  isMine: false
                });
            }
          }
          break;
          
        case MsgType.HELLO:
           if (window.$message) window.$message.success(`Server: ${msg.payload.msg || 'Welcome'}`);
           socket.sendGetRoomList();
           break;
      }
    },

    createRoom(name: string) {
      if (!name.trim()) return;
      this.joinRoom(name);
    },

    getRoomList() {
      if (socket.isConnected.value) {
        socket.sendGetRoomList();
      }
    },

    joinRoom(roomId: string) {
      const user = useUserDataStore();
      if (!socket.isConnected.value) {
        socket.connect();
        this.startHeartbeatLoop();
        setTimeout(() => {
            this.doJoin(roomId, user);
        }, 500);
      } else {
        this.doJoin(roomId, user);
      }
    },

    doJoin(roomId: string, userStore: any) {
      if (this.joinedRooms.includes(roomId)) {
        this.currentRoomId = roomId;
        return;
      }
      
      socket.sendJoinRoom(roomId, {
        id: userStore.getUserData.userId,
        nickname: userStore.getUserData.nickname,
        avatarUrl: userStore.getUserData.avatarUrl
      });
      
      // 单房间模式：清理其他房间数据
      this.joinedRooms.forEach(r => {
        if (r !== roomId) {
          delete this.roomData[r];
        }
      });

      this.joinedRooms = [roomId];
      this.currentRoomId = roomId;
      this.initRoomData(roomId);
      if (window.$message) window.$message.success(`加入房间: ${roomId}`);
      
      setTimeout(() => socket.sendGetRoomList(), 500);
    },

    leaveRoom(roomId: string) {
      const user = useUserDataStore();
      socket.sendLeaveRoom(roomId, {
        id: user.getUserData.userId,
        nickname: user.getUserData.nickname
      });
      
      this.joinedRooms = this.joinedRooms.filter(id => id !== roomId);
      delete this.roomData[roomId];
      
      if (this.currentRoomId === roomId) {
        this.currentRoomId = this.joinedRooms.length > 0 ? this.joinedRooms[0] : '';
      }
      
      if (window.$message) window.$message.info(`已退出房间: ${roomId}`);
      
      if (this.joinedRooms.length === 0) {
          // 不再强制断开，或者根据需求决定
          // 用户需求：只有关闭页面或注销才断开。
          // 所以这里我们不断开连接，只停止心跳? 不，心跳也应该保持以维持连接。
          // 只有 logout 时才调用 stopListen
      } else {
          setTimeout(() => socket.sendGetRoomList(), 500);
      }
    },

    sendMessage(content: string) {
      const user = useUserDataStore();
      if (!content.trim() || !this.currentRoomId) return;
      if (!socket.isConnected.value) {
        if (window.$message) window.$message.error('未连接到聊天服务器');
        return;
      }

      socket.sendChat(
        content, 
        user.getUserData.nickname || 'Guest', 
        user.getUserData.avatarUrl,
        this.currentRoomId
      );
    },

    startHeartbeatLoop() {
      if (this.heartbeatInterval) clearInterval(this.heartbeatInterval);
      // Initial Time Sync
      if (socket.isConnected.value) socket.sendTimeSyncReq();
      
      this.heartbeatInterval = setInterval(() => {
        if (socket.isConnected.value) {
          socket.sendHeartbeat();
          socket.sendGetRoomList();
          // Periodically sync time? Maybe every minute?
          // For now, let's just do it.
          socket.sendTimeSyncReq();
        }
      }, 10000); // 10s
    }
  },
  // 持久化配置，可选，如果需要在刷新页面后恢复状态
  // persist: {
  //   key: 'chatData',
  //   storage: sessionStorage,
  // }
});
