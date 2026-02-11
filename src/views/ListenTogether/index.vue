<!--
  ListenTogether/index.vue
  功能：一起听歌页面（多人实时聊天室 + 播放器控制）
  说明：
    - 左侧：手机端风格播放器 UI
    - 右侧：聊天室/成员列表/房间大厅
-->
<template>
  <div class="listen-together">
    <!-- 左侧：播放器区域 -->
    <div class="lt-player">
      <div class="player-content">
        <!-- 歌曲信息 -->
        <div class="song-info">
          <div class="song-name text-hidden" :title="music.getPlaySongData?.name">
            {{ music.getPlaySongData?.name || '暂无歌曲' }}
          </div>
          <div class="artist-name text-hidden" :title="getArtistNames(music.getPlaySongData?.artist)">
            {{ getArtistNames(music.getPlaySongData?.artist) }}
          </div>
        </div>

        <!-- 封面 / 歌词 -->
        <div class="cover-wrapper" @click="showLyric = !showLyric" :class="{ 'show-lyric': showLyric }">
          <!-- 封面 -->
          <n-image
            v-show="!showLyric"
            :src="coverUrl"
            class="cover-img"
            object-fit="cover"
            fallback-src="/images/logo/logo.png"
            preview-disabled
          />
          <!-- 歌词 -->
          <div v-show="showLyric" class="lyric-full-view" ref="lyricViewRef">
             <div v-if="!music.getPlaySongLyric || music.getPlaySongLyric.length === 0" class="no-lyric">
               暂无歌词
             </div>
             <div v-else class="lyric-scroll">
               <div 
                 v-for="(line, index) in music.getPlaySongLyric" 
                 :key="index"
                 class="lyric-line"
                 :class="{ active: index === music.getPlaySongLyricIndex }"
                 :id="'lyric-' + index"
                 @click.stop="handleLyricClick(line.time)"
               >
                 <div class="text">{{ line.lyric }}</div>
                 <div class="trans" v-if="line.lyricFy && setting.getShowTransl">{{ line.lyricFy }}</div>
               </div>
             </div>
          </div>
        </div>
        
        <!-- 歌词/提示 (简化版) -->
        <div class="lyric-preview text-hidden" :style="{ opacity: showLyric ? 0 : 0.8 }">
           {{ currentLyricLine }}
        </div>

        <!-- 进度条 -->
        <div class="progress-area">
          <span class="time">{{ music.getPlaySongTime.songTimePlayed }}</span>
          <n-slider
            v-model:value="music.getPlaySongTime.barMoveDistance"
            class="progress-slider"
            :step="0.01"
            :tooltip="false"
            @update:value="handleSliderChange"
          />
          <span class="time">{{ music.getPlaySongTime.songTimeDuration }}</span>
        </div>

        <!-- 控制按钮 -->
        <div class="controls">
          <n-icon size="32" :component="SkipPreviousRound" class="control-btn" @click="music.setPlaySongIndex('prev')" />
          <n-icon 
            size="56" 
            :component="music.getPlayState ? PauseCircleFilled : PlayCircleFilled" 
            class="control-btn play-btn" 
            @click="music.setPlayState(!music.getPlayState)" 
          />
          <n-icon size="32" :component="SkipNextRound" class="control-btn" @click="music.setPlaySongIndex('next')" />
        </div>
        
        <!-- 额外功能区 -->
        <div class="extra-actions">
           <n-tooltip trigger="hover">
             <template #trigger>
               <n-icon size="24" :component="FavoriteBorderRound" class="action-icon" />
             </template>
             喜欢
           </n-tooltip>
           <n-tooltip trigger="hover">
             <template #trigger>
               <n-icon size="24" :component="PlaylistPlayRound" class="action-icon" />
             </template>
             播放列表
           </n-tooltip>
           <n-tooltip trigger="hover">
             <template #trigger>
               <n-icon 
                 size="24" 
                 :component="GTranslateFilled" 
                 class="action-icon" 
                 :class="{ active: setting.getShowTransl }" 
                 @click="setting.setShowTransl(!setting.getShowTransl)" 
               />
             </template>
             {{ setting.getShowTransl ? '关闭翻译' : '开启翻译' }}
           </n-tooltip>
           
           <n-popover trigger="click" placement="top" style="padding: 0; background: transparent;">
              <template #trigger>
                 <n-icon size="24" :component="SlowMotionVideoRound" class="action-icon" />
              </template>
              <div class="speed-popup">
                  <div class="val">{{ music.getPlayRate }}x</div>
                  <n-slider
                     v-model:value="music.persistData.playRate"
                     :tooltip="false"
                     :min="0.5"
                     :max="2.0"
                     :step="0.1"
                     vertical
                     class="speed-slider"
                     @update:value="(v) => music.setPlayRate(v)"
                     @click.stop
                   />
              </div>
           </n-popover>

           <n-tooltip trigger="hover">
             <template #trigger>
               <n-icon size="24" :component="FullscreenRound" class="action-icon" @click="music.setBigPlayerState(true)" />
             </template>
             全屏模式
           </n-tooltip>
        </div>
      </div>
    </div>

    <!-- 右侧：房间/聊天区域 -->
    <div class="lt-room">
      <!-- 顶部栏 -->
      <div class="room-header">
        <div class="left">
          <h2 v-if="currentRoomId">
            {{ currentRoomId }} 
            <n-tag size="small" type="success" round v-if="currentRoomId">在线: {{ currentRoomUsers.length }}</n-tag>
          </h2>
          <h2 v-else>大厅</h2>
        </div>
        <div class="right">
          <n-space>
             <n-button size="small" secondary @click="showDebug = !showDebug">Debug</n-button>
             <n-button v-if="!currentRoomId" type="primary" size="small" @click="showCreateRoomModal = true">创建房间</n-button>
             <n-button v-else type="error" size="small" secondary @click="leaveRoom(currentRoomId)">退出房间</n-button>
          </n-space>
        </div>
      </div>

      <!-- 内容区 -->
      <div class="room-content">
        <!-- 未加入房间：显示大厅列表 -->
        <div v-if="!currentRoomId" class="lobby-container">
          <div class="section-title">
             <n-icon :component="ExploreRound" />
             <span>活跃房间</span>
          </div>
          <div v-if="!availableRooms || availableRooms.length === 0" class="empty-rooms">
             <n-empty description="暂无活跃房间，快去创建一个吧" />
          </div>
          <div class="room-grid" v-else>
            <div 
              v-for="room in availableRooms" 
              :key="room.id" 
              class="room-card"
              @click="joinRoom(room.id)"
            >
              <div class="room-card-icon">
                <n-icon size="32" :component="MusicNoteFilled" />
              </div>
              <div class="room-card-info">
                <div class="name">{{ room.id }}</div>
                <div class="count">{{ room.count }} 人在线</div>
              </div>
              <n-button size="small" type="primary" ghost class="join-btn">加入</n-button>
            </div>
          </div>
        </div>

        <!-- 已加入房间：显示聊天和成员 -->
        <div v-else class="active-room-container">
          <div class="room-split-layout">
            <!-- Chat Area -->
            <div class="chat-section">
              <div class="messages-list" ref="messagesRef">
                 <div v-for="(msg, index) in currentRoomMessages" :key="index">
                    <!-- 系统消息 -->
                    <div v-if="msg.type === 'system'" class="system-message">
                      <span class="content">{{ msg.content }}</span>
                    </div>
                    <!-- 聊天消息 -->
                    <div v-else class="message-item" :class="{ 'my-message': msg.isMine }">
                      <n-avatar round size="small" :src="msg.avatarUrl || '/images/ico/user-filling.svg'" class="avatar" />
                      <div class="message-content-wrapper">
                        <span class="sender">{{ msg.sender }}</span>
                        <div class="message-bubble">
                          {{ msg.content }}
                        </div>
                      </div>
                    </div>
                 </div>
              </div>
              <div class="chat-input" :style="{ height: inputHeight + 'px' }">
                <div class="resize-handle-horizontal" @mousedown="startResizeInput"></div>
                <n-input
                  v-model:value="inputValue"
                  type="textarea"
                  placeholder=""
                  :autosize="false"
                  :bordered="false"
                  class="chat-textarea"
                  @keydown.enter.prevent="handleSendMessage"
                />
                <div class="action-bar">
                   <n-button type="primary" size="small" @click="handleSendMessage" class="send-btn">发送(S)</n-button>
                </div>
              </div>
            </div>

            <!-- Resize Handle Vertical -->
            <div class="resize-handle-vertical" @mousedown="startResizeMembers"></div>

            <!-- Members Sidebar -->
            <div class="members-section" :style="{ width: membersWidth + 'px' }">
              <div class="members-header">
                群聊成员 {{ currentRoomUsers.length }}
              </div>
              <div class="members-list">
                 <div v-for="u in currentRoomUsers" :key="u.id" class="member-item">
                    <n-avatar round size="medium" :src="u.avatarUrl || '/images/ico/user-filling.svg'" />
                    <div class="member-info">
                       <div class="nickname">
                         {{ u.nickname }}
                         <n-tag v-if="String(u.id) === String(currentRoomOwnerId)" type="warning" size="tiny" round>房主</n-tag>
                       </div>
                       <!-- <div class="status">在线</div> -->
                    </div>
                 </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建房间弹窗 -->
    <n-modal v-model:show="showCreateRoomModal" preset="dialog" title="创建房间">
      <n-input v-model:value="newRoomName" placeholder="请输入房间名称" @keyup.enter="handleCreateRoom" />
      <template #action>
        <n-button @click="showCreateRoomModal = false">取消</n-button>
        <n-button type="primary" @click="handleCreateRoom">确定</n-button>
      </template>
    </n-modal>

    <!-- Debug Panel -->
    <WebSocketDebug v-if="showDebug" @close="showDebug = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch, defineAsyncComponent, computed } from 'vue';
import { NInput, NButton, NAvatar, NModal, NEmpty, NSpace, NSlider, NIcon, NTabs, NTabPane, NTag, NTooltip, NImage, useThemeVars } from 'naive-ui';
import { 
  PlayCircleFilled, 
  PauseCircleFilled, 
  SkipPreviousRound, 
  SkipNextRound,
  MusicNoteFilled,
  FavoriteBorderRound,
  PlaylistPlayRound,
  GTranslateFilled,
  FullscreenRound,
  ExploreRound,
  SlowMotionVideoRound,
} from "@vicons/material";
import { chatStore, musicStore, settingStore } from '@/store';
import { storeToRefs } from 'pinia';
import { getSongCover } from "@/api/song";

const WebSocketDebug = defineAsyncComponent(() => import('@/components/WebSocketDebug.vue'));

const chat = chatStore();
const music = musicStore();
const setting = settingStore();
const themeVars = useThemeVars();

const { 
  availableRooms, 
  currentRoomId, 
  currentRoomMessages, 
  currentRoomUsers, 
  currentRoomOwnerId 
} = storeToRefs(chat);

// 状态管理
const showCreateRoomModal = ref(false);
const showDebug = ref(false);
const newRoomName = ref('');
const inputValue = ref('');
const messagesRef = ref<HTMLElement | null>(null);
const showLyric = ref(false);
const lyricViewRef = ref<HTMLElement | null>(null);

// 拖拽调整大小逻辑
const membersWidth = ref(240);
const inputHeight = ref(160);
const isResizingMembers = ref(false);
const isResizingInput = ref(false);

const startResizeMembers = () => {
  isResizingMembers.value = true;
  document.addEventListener('mousemove', handleResizeMembers);
  document.addEventListener('mouseup', stopResizeMembers);
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';
};

const handleResizeMembers = (e: MouseEvent) => {
  if (!isResizingMembers.value) return;
  const container = document.querySelector('.room-split-layout');
  if (container) {
    const rect = container.getBoundingClientRect();
    const newWidth = rect.right - e.clientX;
    // 限制最小和最大宽度
    if (newWidth >= 180 && newWidth <= 400) {
      membersWidth.value = newWidth;
    }
  }
};

const stopResizeMembers = () => {
  isResizingMembers.value = false;
  document.removeEventListener('mousemove', handleResizeMembers);
  document.removeEventListener('mouseup', stopResizeMembers);
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
};

const startResizeInput = () => {
  isResizingInput.value = true;
  document.addEventListener('mousemove', handleResizeInput);
  document.addEventListener('mouseup', stopResizeInput);
  document.body.style.cursor = 'row-resize';
  document.body.style.userSelect = 'none';
};

const handleResizeInput = (e: MouseEvent) => {
  if (!isResizingInput.value) return;
  const container = document.querySelector('.chat-section');
  if (container) {
    const rect = container.getBoundingClientRect();
    const newHeight = rect.bottom - e.clientY;
    // 限制最小和最大高度
    if (newHeight >= 100 && newHeight <= 500) {
      inputHeight.value = newHeight;
    }
  }
};

const stopResizeInput = () => {
  isResizingInput.value = false;
  document.removeEventListener('mousemove', handleResizeInput);
  document.removeEventListener('mouseup', stopResizeInput);
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
};

// 辅助函数：获取歌手名
const getArtistNames = (artists: any[]) => {
  if (!artists || artists.length === 0) return 'Unknown Artist';
  return artists.map(a => a.name).join(' / ');
};

// 辅助函数：获取当前歌词
const currentLyricLine = computed(() => {
  const lyric = music.getPlaySongLyric;
  const index = music.getPlaySongLyricIndex;
  if (lyric && lyric.length > 0 && index >= 0 && index < lyric.length) {
    return lyric[index].lyric || '...';
  }
  return '...';
});

const coverUrl = computed(() => {
  if (music.getPlaySongData?.id) {
    return getSongCover(music.getPlaySongData.id);
  }
  return "/images/logo/logo.png";
});

// 播放器逻辑
const handleSliderChange = (val: number) => {
   if ((window as any).$player && music.getPlaySongTime.duration) {
      (window as any).$player.currentTime = (music.getPlaySongTime.duration / 100) * val;
   }
};

const handleLyricClick = (time: number) => {
   if ((window as any).$player) {
      (window as any).$player.currentTime = time;
   }
};

// 监听歌词滚动
watch(() => music.getPlaySongLyricIndex, (index) => {
  if (showLyric.value && index >= 0) {
    nextTick(() => {
      const el = document.getElementById(`lyric-${index}`);
      if (el) {
        el.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
    });
  }
});

// 房间逻辑
const handleCreateRoom = () => {
  if (!newRoomName.value.trim()) return;
  chat.createRoom(newRoomName.value);
  showCreateRoomModal.value = false;
  newRoomName.value = '';
};

const joinRoom = (roomId: string) => {
  chat.joinRoom(roomId);
};

const leaveRoom = (roomId: string) => {
  chat.leaveRoom(roomId);
};

const handleSendMessage = () => {
  if (!inputValue.value.trim()) return;
  chat.sendMessage(inputValue.value);
  inputValue.value = '';
};

// 滚动逻辑
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight;
    }
  });
};

watch(currentRoomId, scrollToBottom);
watch(() => currentRoomMessages.value.length, () => {
  if (currentRoomId.value) scrollToBottom();
});

onMounted(() => {
  chat.startListen();
  if (currentRoomId.value) {
    scrollToBottom();
  } else {
    // 未加入房间时，主动获取一次房间列表
    chat.getRoomList();
  }
  // 隐藏底部播放栏
  music.showPlayBar = false;
});

onUnmounted(() => {
  // chat.stopListen(); // 根据需求决定是否断开
  // 恢复底部播放栏
  music.showPlayBar = true;
});
</script>

<style scoped lang="scss">
.listen-together {
  height: calc(100vh - 100px); // 减去顶部导航和底部播放条的高度
  display: flex;
  gap: 24px;
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

/* 左侧播放器 */
.lt-player {
  width: 380px;
  flex-shrink: 0;
  background: var(--n-color-modal);
  border-radius: 24px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--n-border-color);

  .player-content {
    flex: 1;
    padding: 32px 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
  }

  .cover-wrapper {
    width: 260px;
    height: 260px;
    margin-bottom: 24px;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    position: relative;
    background: #000;
    transition: all 0.3s;
    
    &.show-lyric {
      background: transparent;
      box-shadow: none;
      border: 1px solid var(--n-border-color);
      width: 100%;
      height: 360px;
    }
    
    .cover-img {
       width: 100%;
       height: 100%;
       display: flex;
       justify-content: center;
       align-items: center;
       
       :deep(img) {
         width: 100%;
         height: 100%;
         transition: opacity 0.3s;
       }
    }

    .lyric-full-view {
      width: 100%;
      height: 100%;
      overflow-y: auto;
      padding: 20px;
      // 隐藏滚动条
      scrollbar-width: none; /* Firefox */
      -ms-overflow-style: none; /* IE/Edge */
      &::-webkit-scrollbar {
        display: none; /* Chrome/Safari */
      }
      
      .no-lyric {
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        color: var(--n-text-color-3);
      }

      .lyric-scroll {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 16px;
        padding: 20px 0; // 留出上下空间以便滚动
        
        .lyric-line {
          text-align: center;
          cursor: pointer;
          transition: all 0.2s;
          opacity: 0.6;
          
          &:hover {
            opacity: 0.8;
          }
          
          &.active {
            opacity: 1;
            transform: scale(1.1);
            color: var(--n-color-primary);
            font-weight: bold;
          }
          
          .text {
            font-size: 14px;
            line-height: 1.5;
          }
          
          .trans {
            font-size: 12px;
            margin-top: 4px;
            opacity: 0.8;
          }
        }
      }
    }
  }

  .song-info {
    text-align: center;
    width: 100%;
    margin-bottom: 16px;

    .song-name {
      font-size: 20px;
      font-weight: bold;
      margin-bottom: 8px;
      color: var(--n-text-color);
    }

    .artist-name {
      font-size: 14px;
      color: var(--n-text-color-3);
    }
  }

  .lyric-preview {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--n-text-color-2);
    font-size: 14px;
    margin-bottom: 24px;
    text-align: center;
    opacity: 0.8;
  }

  .progress-area {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 32px;

    .time {
      font-size: 12px;
      color: var(--n-text-color-3);
      min-width: 40px;
      text-align: center;
    }

    .progress-slider {
      flex: 1;
    }
  }

  .controls {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 32px;
    margin-bottom: 32px;

    .control-btn {
      cursor: pointer;
      opacity: 0.8;
      transition: all 0.2s;
      
      &:hover {
        opacity: 1;
        transform: scale(1.1);
      }
      
      &.play-btn {
        color: var(--n-color-primary);
        opacity: 1;
      }
    }
  }

  .extra-actions {
    display: flex;
    gap: 24px;
    
    .action-icon {
      cursor: pointer;
      color: var(--n-text-color-3);
      transition: all 0.3s;
      
      &:hover {
        color: var(--n-text-color);
        transform: scale(1.1);
      }
      
      &.active {
        color: var(--n-color-primary);
      }
    }
  }
  
  .speed-popup {
     width: 44px;
     height: 140px;
     background: var(--n-color-modal);
     backdrop-filter: blur(10px);
     border-radius: 18px;
     padding: 12px 0;
     display: flex;
     flex-direction: column;
     align-items: center;
     box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
     border: 1px solid var(--n-border-color);
     z-index: 10000;
     
     .val {
        font-size: 10px;
        margin-bottom: 8px;
        font-weight: bold;
     }
     
     .speed-slider {
        height: 100%;
        --n-handle-size: 12px;
        --n-rail-width: 4px;
        
        :deep(.n-slider-rail) {
          background-color: var(--n-border-color);
          .n-slider-rail__fill {
            background-color: var(--n-color-primary);
          }
        }
        :deep(.n-slider-handle) {
          background-color: #fff;
          box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
        }
     }
  }
}

/* 右侧房间区域 */
.lt-room {
  flex: 1;
  background: var(--n-color-modal);
  border-radius: 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--n-border-color);

  .room-header {
    padding: 20px 24px;
    border-bottom: 1px solid var(--n-border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;

    h2 {
      margin: 0;
      font-size: 18px;
      display: flex;
      align-items: center;
      gap: 12px;
    }
  }

  .room-content {
    flex: 1;
    overflow: hidden;
    position: relative;
  }
}

/* 大厅视图 */
.lobby-container {
  padding: 24px;
  height: 100%;
  overflow-y: auto;

  .section-title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 16px;
    color: var(--n-text-color-2);
  }

  .room-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
  }

  .room-card {
    background: var(--n-card-color);
    border: 1px solid var(--n-border-color);
    border-radius: 12px;
    padding: 16px;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      border-color: var(--n-color-primary);
    }

    .room-card-icon {
      width: 48px;
      height: 48px;
      background: rgba(var(--n-primary-color-rgb), 0.1);
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 12px;
      color: var(--n-color-primary);
    }

    .room-card-info {
      margin-bottom: 16px;
      .name {
        font-weight: bold;
        font-size: 16px;
        margin-bottom: 4px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        max-width: 100%;
      }
      .count {
        font-size: 12px;
        color: var(--n-text-color-3);
      }
    }
    
    .join-btn {
      width: 100%;
    }
  }
  
  .empty-rooms {
     display: flex;
     justify-content: center;
     align-items: center;
     height: 300px;
  }
}

/* 活跃房间视图 */
.active-room-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.room-split-layout {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.chat-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  // border-right: 3px solid var(--n-border-color); // Replaced by resizer
  
  .messages-list {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    background: rgba(0, 0, 0, 0.01);
  }

  .chat-input {
    padding: 12px 16px;
    background: var(--n-card-color);
    // border-top: 3px solid var(--n-border-color); // Replaced by resizer
    display: flex;
    flex-direction: column;
    gap: 8px;
    position: relative;
    min-height: 100px;
    
    .resize-handle-horizontal {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 3px;
      background-color: v-bind('themeVars.primaryColor');
      cursor: row-resize;
      transition: height 0.2s;
      z-index: 10;
      
      &:hover {
        height: 5px;
      }
    }

    .chat-textarea {
      background: transparent;
      padding: 0;
      flex: 1;
      
      :deep(.n-input-wrapper) {
         padding: 0;
         height: 100%;
      }
      :deep(.n-input__textarea-el) {
         padding: 0;
         height: 100%;
      }
      :deep(.n-input__placeholder) {
         padding: 0;
      }
    }

    .action-bar {
       display: flex;
       justify-content: flex-end;
       flex-shrink: 0;
       
       .send-btn {
          padding: 0 20px;
       }
    }
  }
}

.resize-handle-vertical {
  width: 3px;
  background-color: v-bind('themeVars.primaryColor');
  cursor: col-resize;
  transition: width 0.2s;
  z-index: 10;
  
  &:hover {
    width: 5px;
  }
}

.members-section {
  width: 240px; // Default, overridden by inline style
  display: flex;
  flex-direction: column;
  background: var(--n-card-color);
  
  .members-header {
    padding: 16px;
    font-weight: bold;
    border-bottom: 1px solid var(--n-border-color);
    color: var(--n-text-color);
  }

  .members-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    
    .member-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 8px;
      border-radius: 6px;
      margin-bottom: 4px;
      cursor: pointer;
      
      &:hover {
        background: var(--n-color-embedded);
      }

      .member-info {
        flex: 1;
        overflow: hidden;
        
        .nickname {
          font-size: 14px;
          display: flex;
          align-items: center;
          gap: 6px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
}

/* 消息气泡 */
.system-message {
  text-align: center;
  margin: 16px 0;
  .content {
    font-size: 12px;
    color: var(--n-text-color-3);
    background: rgba(0, 0, 0, 0.05);
    padding: 2px 10px;
    border-radius: 10px;
  }
}

.message-item {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: flex-start;

  &.my-message {
    flex-direction: row-reverse;
    
    .message-content-wrapper {
      align-items: flex-end;
      
      .message-bubble {
        background: v-bind('themeVars.primaryColor');
        color: white;
        border-top-left-radius: 12px;
        border-top-right-radius: 2px;
      }
    }
  }

  .message-content-wrapper {
    display: flex;
    flex-direction: column;
    max-width: 70%;
    
    .sender {
      font-size: 12px;
      color: var(--n-text-color-3);
      margin-bottom: 4px;
      padding: 0 4px;
    }

    .message-bubble {
      background: var(--n-card-color);
      border: 1px solid var(--n-border-color);
      padding: 10px 14px;
      border-radius: 12px;
      border-top-left-radius: 2px;
      line-height: 1.5;
      word-break: break-all;
      box-shadow: 0 2px 4px rgba(0,0,0,0.05);
    }
  }
}

.text-hidden {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
