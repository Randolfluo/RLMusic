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
          <Transition name="fade" mode="out-in">
             <div :key="music.getPlaySongData?.id" class="song-info-inner">
                <div class="song-name text-hidden" :title="music.getPlaySongData?.name">
                  {{ music.getPlaySongData?.name || '暂无歌曲' }}
                </div>
                <div class="artist-name text-hidden" :title="getArtistNames(music.getPlaySongData?.artist)">
                  {{ getArtistNames(music.getPlaySongData?.artist) }}
                </div>
             </div>
          </Transition>
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
          <div class="slider-wrapper">
            <n-slider
              v-model:value="music.getPlaySongTime.barMoveDistance"
              class="progress-slider"
              :step="0.01"
              :tooltip="false"
              @update:value="handleSliderChange"
            >
              <template #thumb>
                <div class="custom-thumb">
                  <div class="time-tooltip">{{ tooltipTime }}</div>
                </div>
              </template>
            </n-slider>
          </div>
        </div>

        <!-- 控制按钮 -->
        <div class="controls">
          <n-icon size="32" :component="SkipPreviousRound" class="control-btn" @click="changeSong('prev')" />
          <n-icon 
            size="56" 
            :component="music.getPlayState ? PauseCircleFilled : PlayCircleFilled" 
            class="control-btn play-btn" 
            @click.stop="togglePlay" 
          />
          <n-icon size="32" :component="SkipNextRound" class="control-btn" @click="changeSong('next')" />
        </div>
        
        <!-- 额外功能区 -->
        <div class="extra-actions">
           <!-- 音量调节 -->
           <div class="volume-control">
             <n-icon 
                size="24" 
                :component="
                  music.persistData.playVolume == 0
                    ? VolumeOffRound
                    : music.persistData.playVolume < 0.4
                    ? VolumeMuteRound
                    : music.persistData.playVolume < 0.7
                    ? VolumeDownRound
                    : VolumeUpRound
                " 
                class="action-icon"
                @click="volumeMute" 
              />
              <div class="volume-popup">
                  <div class="val">{{ Math.round(music.persistData.playVolume * 100) }}%</div>
                  <n-slider
                     v-model:value="music.persistData.playVolume"
                     :tooltip="false"
                     :min="0"
                     :max="1"
                     :step="0.01"
                     vertical
                     class="volume-slider"
                   />
              </div>
           </div>

           <n-tooltip trigger="hover">
             <template #trigger>
               <n-icon size="24" :component="PlaylistPlayRound" class="action-icon" @click.stop="music.showPlayList = !music.showPlayList" />
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
           
           <div class="speed-control">
              <n-icon size="24" :component="SlowMotionVideoRound" class="action-icon" />
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
                     @update:value="handleSpeedChange"
                     @click.stop
                   />
              </div>
           </div>

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
                      <n-avatar round size="small" :src="resolveAvatarUrl(msg.avatarUrl) || '/images/ico/user-filling.svg'" class="avatar" />
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
                    <n-avatar round size="medium" :src="resolveAvatarUrl(u.avatarUrl) || defaultAvatar" />
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

    <!-- 播放列表组件 -->
    <div class="playlist-container" v-if="music.showPlayList">
       <PlayList class="fixed-mode" />
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
import { NInput, NButton, NAvatar, NModal, NEmpty, NSpace, NSlider, NIcon, NTag, NTooltip, NImage, useThemeVars } from 'naive-ui';
import { 
  PlayCircleFilled, 
  PauseCircleFilled, 
  SkipPreviousRound, 
  SkipNextRound,
  MusicNoteFilled,
  PlaylistPlayRound,
  GTranslateFilled,
  FullscreenRound,
  VolumeOffRound,
  VolumeMuteRound,
  VolumeDownRound,
  VolumeUpRound,
  SlowMotionVideoRound,
  ExploreRound,
} from "@vicons/material";
import { chatStore, musicStore, settingStore } from '@/store';
import { storeToRefs } from 'pinia';
import { getSongCover } from "@/api/song";
import { resolveAvatarUrl } from "@/api/user";
import { timelineEngine } from "@/core/realtime/timeline";
import PlayList from "@/components/DataList/PlayList.vue";

const WebSocketDebug = defineAsyncComponent(() => import('@/components/WebSocketDebug.vue'));

const chat = chatStore();
const music = musicStore();
const setting = settingStore();
// const user = userStore();
const themeVars = useThemeVars();
// const { persistData } = storeToRefs(music);

const defaultAvatar = computed(() => `${import.meta.env.BASE_URL}images/ico/user-filling.svg`);

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
// const lyricViewRef = ref<HTMLElement | null>(null);

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

// 辅助函数：格式化秒数
const formatSeconds = (seconds: number) => {
  if (!seconds && seconds !== 0) return '00:00';
  const m = Math.floor(seconds / 60);
  const s = Math.floor(seconds % 60);
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
};

// 实时计算拖拽时的显示时间
const tooltipTime = computed(() => {
   const duration = music.getPlaySongTime.duration || 0;
   const val = music.getPlaySongTime.barMoveDistance || 0;
   return formatSeconds((duration / 100) * val);
});

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

// 音量逻辑
const volumeMute = () => {
  if (music.persistData.playVolume > 0) {
    music.persistData.playVolumeMute = music.persistData.playVolume;
    music.persistData.playVolume = 0;
  } else {
    music.persistData.playVolume = music.persistData.playVolumeMute || 0.5;
  }
};

// 播放器逻辑
const handleSliderChange = (val: number) => {
   if (music.getPlaySongTime.duration) {
      const time = (music.getPlaySongTime.duration / 100) * val;
      // 立即更新本地播放器进度
      const player = (window as any).$player;
      if (player) {
        player.currentTime = time;
      }
      timelineEngine.requestSeek(time);
   }
};

const handleLyricClick = (time: number) => {
   timelineEngine.requestSeek(time);
};

const togglePlay = () => {
  // 乐观更新：立即切换图标状态
  music.setPlayState(!music.getPlayState);
  
  // 实际执行逻辑
  timelineEngine.togglePlay();
};

const changeSong = (dir: 'prev' | 'next') => {
    const list = music.getPlaylists;
    const len = list.length;
    if (len === 0) return;
    let index = music.getPlaySongIndex;
    if (dir === 'next') {
        index = (index + 1) % len;
    } else {
        index = (index - 1 + len) % len;
    }
    const song = list[index];
    if (song && song.id) {
        // 先乐观更新 UI，防止点击无反应
        music.setPlaySongIndex(index);
        timelineEngine.requestChangeSong(String(song.id));
    }
};

const handleSpeedChange = (v: number) => {
    // 立即应用到本地，提升响应速度
    music.setPlayRate(v);
    const player = (window as any).$player;
    if (player) player.playbackRate = v;
    
    // 发送请求
    timelineEngine.requestSetSpeed(v);
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
  
  // 清理可能存在的全局事件监听
  document.removeEventListener('mousemove', handleResizeMembers);
  document.removeEventListener('mouseup', stopResizeMembers);
  document.removeEventListener('mousemove', handleResizeInput);
  document.removeEventListener('mouseup', stopResizeInput);
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
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
  position: relative; // 确保播放列表可以绝对定位在容器内
}

  /* 播放列表容器 */
  .playlist-container {
     position: absolute;
     right: 24px;
     bottom: 24px;
     z-index: 1000;
     
     :deep(.play-list) {
        position: relative;
        bottom: auto;
        right: auto;
        height: 500px;
        background: var(--n-color-modal);
        border: 1px solid var(--n-border-color);
        border-radius: 12px;
        box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
     }
  }

/* 左侧播放器 */
  .lt-player {
    z-index: 2;
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
    height: 50px; /* 固定高度防止跳动 */
    
    .song-info-inner {
       width: 100%;
       display: flex;
       flex-direction: column;
       align-items: center;
    }

    .song-name {
      font-size: 20px;
      font-weight: bold;
      margin-bottom: 4px;
      color: var(--n-text-color);
      width: 100%;
    }

    .artist-name {
      font-size: 14px;
      color: var(--n-text-color-3);
      width: 100%;
    }
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.3s ease;
  }

  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
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

    .slider-wrapper {
      flex: 1;
      padding: 8px 0; /* 增加点击热区 */
      cursor: pointer;
      position: relative;
      display: flex;
      align-items: center;
    }

    .progress-slider {
      width: 100%;
      --n-rail-height: 4px;
      --n-handle-size: 16px;
      
      :deep(.n-slider-rail) {
        background-color: #E0E0E0;
        border-radius: 2px;
        
        .n-slider-rail__fill {
          background-color: var(--n-color-primary);
          border-radius: 2px;
          transition: width 150ms cubic-bezier(0.165, 0.84, 0.44, 1);
        }
      }
      
      /* 自定义 Thumb 容器 */
      :deep(.n-slider-handle) {
        background-color: transparent;
        box-shadow: none;
        width: 0;
        height: 0;
        
        /* 实际显示的圆点在 slot 中 */
        .custom-thumb {
           width: 16px;
           height: 16px;
           background-color: #ffffff;
           border-radius: 50%;
           box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
           position: relative;
           display: flex;
           justify-content: center;
           align-items: center;
           transition: transform 0.2s;
           
           /* 悬停放大 */
           &:hover {
              transform: scale(1.2);
              
              .time-tooltip {
                 opacity: 1;
                 transform: translateY(-8px);
              }
           }
           
           /* 拖动时显示时间提示 */
           .time-tooltip {
              position: absolute;
              bottom: 100%;
              background-color: rgba(0, 0, 0, 0.8);
              color: #fff;
              padding: 4px 8px;
              border-radius: 4px;
              font-size: 12px;
              white-space: nowrap;
              opacity: 0;
              transform: translateY(0);
              transition: all 0.2s;
              pointer-events: none;
              margin-bottom: 8px;
           }
        }
      }
    }

    /* 高对比度模式适配 */
    @media (prefers-contrast: high) {
      .progress-slider {
        :deep(.n-slider-rail) {
          background-color: #FFF;
          border: 1px solid #000;
          .n-slider-rail__fill {
            background-color: #000;
          }
        }
      }
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
    align-items: center; // 确保子元素垂直居中
    
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

    .volume-control,
    .speed-control {
       position: relative;
       display: flex;
       align-items: center;
       justify-content: center;
       height: 100%; /* 继承父高度 */

       /* 弹窗样式调整为绝对定位 */
       .volume-popup,
       .speed-popup {
          position: absolute;
          bottom: 100%; /* 在上方显示 */
          left: 50%;
          transform: translateX(-50%);
          margin-bottom: 4px; /* 间距 */
          
          /* 初始隐藏 */
          opacity: 0;
          visibility: hidden;
          transition: all 0.2s ease-in-out;
          pointer-events: none; /* 防止未显示时阻挡点击 */
       }

       /* Hover 显示逻辑 */
       &:hover {
          .volume-popup,
          .speed-popup {
             opacity: 1;
             visibility: visible;
             pointer-events: auto;
             bottom: calc(100% + 4px); /* 稍微上浮动画 */
          }
       }
       
       /* 连接桥，防止鼠标移出间隙导致弹窗消失 */
       &::after {
          content: '';
          position: absolute;
          top: -20px;
          left: 0;
          width: 100%;
          height: 30px;
          background: transparent;
       }
    }
  }
}

.volume-popup,
.speed-popup {
     /* 保持原有样式，仅移除位置相关属性（如果之前写死的话） */
     width: 44px;
     height: 140px;
     background: var(--n-color-modal);
     backdrop-filter: blur(10px);
     border-radius: 18px;
     padding: 12px 0;
     display: flex;
     flex-direction: column;
     align-items: center;
     box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
     border: 1px solid var(--n-border-color);
     z-index: 10000;
     
     .val {
        font-size: 10px;
        margin-bottom: 8px;
        font-weight: bold;
     }
     
     .volume-slider,
     .speed-slider {
        height: 100%;
        --n-handle-size: 12px;
        --n-rail-width: 4px;
        
        :deep(.n-slider-rail) {
          background-color: #E0E0E0;
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

/* 确保内容不会被遮挡 */
.lt-room {
  flex: 1;
  background: var(--n-color-modal);
  border-radius: 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--n-border-color);
  z-index: 1; /* 确保低于播放列表 */

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
