import { socket, MsgType, type WSMessage } from './socket';
import { useMusicDataStore } from '@/store/musicData';
import { useChatDataStore } from '@/store/chatData';

class TimelineEngine {
  private _music: any = null;
  private _chat: any = null;

  get music() {
    if (!this._music) this._music = useMusicDataStore();
    return this._music;
  }

  get chat() {
    if (!this._chat) this._chat = useChatDataStore();
    return this._chat;
  }
  
  constructor() {
    socket.onMessage(this.handleMessage.bind(this));
  }

  // Handle messages from Server
  handleMessage(msg: WSMessage) {
    // Only process if we are in a room? 
    // Yes, but MsgType.ROOM_INFO implies we are involved.
    
    switch (msg.type) {
      case MsgType.ROOM_INFO:
        if (msg.payload && msg.payload.timeline) {
            this.applyTimeline(msg.payload.timeline);
        }
        break;
      // We rely on ROOM_INFO for state updates as it is broadcasted after every change.
      // Specific events can be used for UI notifications if needed.
    }
  }

  applyTimeline(timeline: any) {
    if (!timeline) return;
    
    const player = (window as any).$player;
    if (!player) return;

    // 1. Check Song
    // If song ID differs, we might need to change it. 
    // For now, assuming the user has the same playlist context.
    // If timeline.song_id is implemented...
    
    // 2. Calculate Target Position
    const serverNow = Date.now() + socket.timeOffset.value;
    let targetPosition = 0;
    
    if (timeline.paused) {
        targetPosition = timeline.pause_position_ms / 1000; // ms to s
    } else {
        const elapsed = serverNow - timeline.start_timestamp;
        targetPosition = (elapsed * timeline.speed) / 1000;
    }
    
    // Prevent negative time
    if (targetPosition < 0) targetPosition = 0;

    // 3. Apply to Player (Sync)
    const currentPosition = player.currentTime;
    const drift = Math.abs(targetPosition - currentPosition);
    
    // Threshold: 0.2s
    if (drift > 0.2) {
        console.log(`[Timeline] Syncing: Drift ${drift.toFixed(3)}s. Seek to ${targetPosition.toFixed(3)}`);
        player.currentTime = targetPosition;
    }

    // 4. Play/Pause State
    if (timeline.paused) {
        if (!player.paused) {
            player.pause();
            this.music.setPlayState(false);
        }
    } else {
        if (player.paused) {
            player.play(); // This might fail if no user interaction, but usually fine in app
            this.music.setPlayState(true);
        }
    }
    
    // 5. Speed
    if (timeline.speed && Math.abs(player.playbackRate - timeline.speed) > 0.01) {
        player.playbackRate = timeline.speed;
        this.music.setPlayRate(timeline.speed);
    }
  }
  
  // --- Public methods for UI (Client -> Server) ---
  
  requestPlay() {
     if (!this.chat.currentRoomId) return;
     socket.sendPlay(this.chat.currentRoomId);
  }
  
  requestPause() {
     if (!this.chat.currentRoomId) return;
     const player = (window as any).$player;
     const pos = player ? player.currentTime * 1000 : 0;
     socket.sendPause(this.chat.currentRoomId, pos);
  }
  
  requestSeek(time: number) {
     if (!this.chat.currentRoomId) return;
     socket.sendSeek(this.chat.currentRoomId, time * 1000);
  }
  
  requestChangeSong(songId: string) {
     if (!this.chat.currentRoomId) return;
     socket.sendChangeSong(this.chat.currentRoomId, songId);
  }
  
  requestSetSpeed(speed: number) {
     if (!this.chat.currentRoomId) return;
     socket.sendSetSpeed(this.chat.currentRoomId, speed);
  }
  
  // Toggle Play/Pause convenience method
  togglePlay() {
      const player = (window as any).$player;
      if (player && !player.paused) {
          this.requestPause();
      } else {
          this.requestPlay();
      }
  }
}

export const timelineEngine = new TimelineEngine();
