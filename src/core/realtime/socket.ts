import { ref } from 'vue';

// Message Types
export enum MsgType {
  HELLO = 'HELLO',
  TIME_SYNC_REQ = 'TIME_SYNC_REQ',
  TIME_SYNC_RES = 'TIME_SYNC_RES',
  HEARTBEAT = 'HEARTBEAT',
  CHAT = 'CHAT',
  JOIN_ROOM = 'JOIN_ROOM',
  LEAVE_ROOM = 'LEAVE_ROOM',
  ROOM_MEMBERS = 'ROOM_MEMBERS',
  ROOM_INFO = 'ROOM_INFO',
  GET_ROOM_LIST = 'GET_ROOM_LIST',
  ROOM_LIST_RES = 'ROOM_LIST_RES',
  TIMELINE_INIT = 'TIMELINE_INIT',
  PLAY = 'PLAY',
  PAUSE = 'PAUSE',
  SEEK = 'SEEK',
  CHANGE_SONG = 'CHANGE_SONG',
  SET_SPEED = 'SET_SPEED'
}

export interface WSMessage {
  type: string;
  payload: any;
}

class SocketClient {
  private ws: WebSocket | null = null;
  private url: string;
  public isConnected = ref(false);
  public lastLatency = ref(0);
  public timeOffset = ref(0);
  
  // Callbacks
  private onMessageCallbacks: ((msg: WSMessage) => void)[] = [];

  constructor() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    this.url = `${protocol}//${window.location.host}/api/ws/chat`;
  }

  connect() {
    if (this.ws) return;

    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log('WS Connected');
      this.isConnected.value = true;
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as WSMessage;
        this.handleSystemMessage(data);
        this.onMessageCallbacks.forEach(cb => cb(data));
      } catch (e) {
        console.error('WS Parse Error', e);
      }
    };

    this.ws.onclose = () => {
      console.log('WS Disconnected');
      this.isConnected.value = false;
      this.ws = null;
    };

    this.ws.onerror = (err) => {
      console.error('WS Error', err);
    };
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  send(type: string, payload: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }));
    } else {
      console.warn('WS not connected, cannot send', type);
    }
  }

  // --- Specific Events ---

  sendHello(user: { id: number; nickname: string }) {
    this.send(MsgType.HELLO, user);
  }

  sendTimeSyncReq() {
    this.send(MsgType.TIME_SYNC_REQ, { client_now: Date.now() });
  }

  sendHeartbeat() {
    this.send(MsgType.HEARTBEAT, { client_now: Date.now() });
  }

  sendChat(content: string, sender: string, avatarUrl: string, targetRoomId: string) {
    this.send(MsgType.CHAT, { content, sender, avatarUrl, targetRoomId });
  }

  sendGetRoomList() {
    this.send(MsgType.GET_ROOM_LIST, {});
  }

  sendJoinRoom(roomId: string, user: any) {
    this.send(MsgType.JOIN_ROOM, { roomId, user });
  }

  sendLeaveRoom(roomId: string, user: any) {
    this.send(MsgType.LEAVE_ROOM, { roomId, user });
  }

  sendPlay(roomId: string) {
    this.send(MsgType.PLAY, { roomId });
  }

  sendPause(roomId: string, positionMs: number) {
    this.send(MsgType.PAUSE, { roomId, position_ms: positionMs });
  }

  sendSeek(roomId: string, positionMs: number) {
    this.send(MsgType.SEEK, { roomId, position_ms: positionMs });
  }

  sendChangeSong(roomId: string, songId: string) {
    this.send(MsgType.CHANGE_SONG, { roomId, song_id: songId });
  }

  sendSetSpeed(roomId: string, speed: number) {
    this.send(MsgType.SET_SPEED, { roomId, speed });
  }

  // --- Internal Handling ---

  private handleSystemMessage(msg: WSMessage) {
    switch (msg.type) {
      case MsgType.TIME_SYNC_RES:
        // Handle Time Sync
        // msg.payload.server_now
        if (msg.payload && msg.payload.server_now) {
          const clientNow = Date.now();
          const rtt = this.lastLatency.value * 2;
          const serverNow = msg.payload.server_now;
          // Offset = server_time - local_time
          // Correct formula: server_now approx (client_now + rtt/2) + offset
          // offset = server_now - (client_now - rtt/2)? 
          // task.md says: offset = server_time - Date.now(). 
          // Assuming server_time is "time when server received".
          // More precise: offset = server_now - (client_req + rtt/2)
          // But simplified: offset = server_now - client_now
          this.timeOffset.value = serverNow - clientNow;
          console.log('[TimeSync] Offset:', this.timeOffset.value, 'ms');
        }
        break;
      case MsgType.HEARTBEAT:
        // Handle Heartbeat RTT
        if (msg.payload && msg.payload.client_now) {
          const rtt = Date.now() - msg.payload.client_now;
          this.lastLatency.value = Math.round(rtt / 2);
        }
        break;
    }
  }

  // --- API ---

  onMessage(cb: (msg: WSMessage) => void) {
    this.onMessageCallbacks.push(cb);
  }

  offMessage(cb: (msg: WSMessage) => void) {
    this.onMessageCallbacks = this.onMessageCallbacks.filter(c => c !== cb);
  }
}

export const socket = new SocketClient();
