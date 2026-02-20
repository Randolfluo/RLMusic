<template>
  <div class="desktop-lyric" :class="{ locked: isLocked }" :style="containerStyle">
    <!-- Header / Controls (Hidden when locked, shown on hover) -->
    <transition name="fade">
      <div class="header" v-show="!isLocked || showUnlock">
        <div class="left drag-region" v-show="!isLocked">
          <div class="app-icon">
            <n-icon size="16" :component="MusicNoteFilled" />
          </div>
          <span class="title">{{ songTitle || 'Local Music Player' }}</span>
        </div>
        
        <!-- 锁定状态下的解锁按钮 -->
        <div class="unlock-btn no-drag" v-if="isLocked" @click="toggleLock" title="解锁">
          <n-icon size="18" :component="LockFilled" />
        </div>

      <div class="center no-drag" v-if="!isLocked">
        <div class="control-btn" @click="control('prev')" title="上一首">
          <n-icon size="22" :component="SkipPreviousFilled" />
        </div>
        <div class="control-btn play-btn" @click="control('toggle')" :title="isPlaying ? '暂停' : '播放'">
          <n-icon size="28" :component="isPlaying ? PauseFilled : PlayArrowFilled" />
        </div>
        <div class="control-btn" @click="control('next')" title="下一首">
          <n-icon size="22" :component="SkipNextFilled" />
        </div>
      </div>
      
      <div class="right no-drag" v-if="!isLocked">
        <!-- Settings (FontSize) -->
        <n-popover trigger="click" placement="bottom-end" :show-arrow="false" style="padding: 0; background: transparent;">
          <template #trigger>
            <div class="action-btn" title="设置">
              <n-icon size="18" :component="SettingsFilled" />
            </div>
          </template>
          <div class="settings-panel glass-effect">
            <div class="panel-header">设置</div>
            <div class="setting-item">
              <span class="label">字体大小</span>
              <n-slider v-model:value="fontSize" :min="1.5" :max="5" :step="0.1" class="custom-slider" @update:value="updateSettings" />
            </div>
            <div class="setting-item">
              <span class="label">跟随主题</span>
              <n-switch v-model:value="followTheme" size="small" @update:value="updateSettings" />
            </div>
            <div class="setting-item">
              <span class="label">双语歌词</span>
              <n-switch v-model:value="showTranslation" size="small" @update:value="updateSettings" />
            </div>
          </div>
        </n-popover>

        <div class="action-btn" @click="toggleLock" :title="isLocked ? '解锁' : '锁定'">
          <n-icon size="18" :component="isLocked ? LockFilled : LockOpenFilled" />
        </div>
        <div class="action-btn close-btn" @click="close" title="关闭">
          <n-icon size="18" :component="CloseFilled" />
        </div>
      </div>
      </div>
    </transition>

    <!-- Lyric Content -->
    <div class="lyric-content">
      <div class="lrc-line current" :style="currentLyricStyle">
        {{ currentLyric || '暂无歌词' }}
        <div class="tlyric" v-if="showTranslation && currentTlyric">{{ currentTlyric }}</div>
      </div>
      <div class="lrc-line next" :style="nextLyricStyle">
        {{ nextLyric }}
        <div class="tlyric" v-if="showTranslation && nextTlyric">{{ nextTlyric }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { NIcon, NPopover, NSlider, NSwitch } from 'naive-ui';
import { 
  MusicNoteFilled, 
  SkipPreviousFilled, 
  PlayArrowFilled, 
  PauseFilled, 
  SkipNextFilled,
  SettingsFilled,
  LockFilled,
  LockOpenFilled,
  CloseFilled 
} from '@vicons/material';

const songTitle = ref('');
const currentLyric = ref('');
const currentTlyric = ref('');
const nextLyric = ref('');
const nextTlyric = ref('');
const isPlaying = ref(false);
const isLocked = ref(false);
const fontSize = ref(3.0);
const followTheme = ref(true);
const themeColor = ref('#009688'); // 默认主题色
const showTranslation = ref(true);
const showUnlock = ref(false);

// Styles
const containerStyle = computed(() => {
  return {
    // '--theme-color': followTheme.value ? themeColor.value : '#ffffff'
  }
});

const currentLyricStyle = computed(() => {
  return {
    fontSize: `${fontSize.value}rem`,
    color: followTheme.value ? themeColor.value : '#ffffff',
    textShadow: '0 2px 4px rgba(0,0,0,0.5), 0 0 2px rgba(0,0,0,0.8)'
  }
});

const nextLyricStyle = computed(() => {
  return {
    fontSize: `${fontSize.value * 0.6}rem`,
    color: followTheme.value ? themeColor.value : '#ffffff',
    opacity: 0.6,
    textShadow: '0 2px 4px rgba(0,0,0,0.5)'
  }
});

// Actions
const control = (action: string) => {
  window.ipcRenderer.send('desktop-lyric-control', action);
};

const toggleLock = () => {
  isLocked.value = !isLocked.value;
  window.ipcRenderer.send('lock-desktop-lyric', isLocked.value);
};

const close = () => {
  window.ipcRenderer.send('close-desktop-lyric');
};

const updateSettings = () => {
  window.ipcRenderer.send('update-desktop-lyric-settings', {
    fontSize: fontSize.value,
    followTheme: followTheme.value,
    showTranslation: showTranslation.value
  });
};

// IPC Listeners
onMounted(() => {
  // Ensure transparent background
  document.documentElement.style.backgroundColor = 'transparent';
  document.body.style.backgroundColor = 'transparent';
  const app = document.getElementById('app');
  if (app) app.style.backgroundColor = 'transparent';

  window.ipcRenderer.on('update-lyric', (_event: any, data: any) => {
    if (data.current !== undefined) currentLyric.value = data.current;
    if (data.next !== undefined) nextLyric.value = data.next;
    if (data.currentTlyric !== undefined) currentTlyric.value = data.currentTlyric;
    if (data.nextTlyric !== undefined) nextTlyric.value = data.nextTlyric;
    if (data.isPlaying !== undefined) isPlaying.value = data.isPlaying;
  });
  
  window.ipcRenderer.on('update-settings', (_event: any, data: any) => {
    if (data.fontSize) fontSize.value = data.fontSize;
    if (data.followTheme !== undefined) followTheme.value = data.followTheme;
    if (data.themeColor) themeColor.value = data.themeColor;
    if (data.showTranslation !== undefined) showTranslation.value = data.showTranslation;
  });

  window.ipcRenderer.on('desktop-lyric-locked', (_event: any, locked: boolean) => {
    isLocked.value = locked;
  });

  // 监听鼠标移入移出，用于在锁定状态下显示解锁按钮
  document.addEventListener('mouseenter', () => {
    if (isLocked.value) {
      showUnlock.value = true;
    }
  });

  document.addEventListener('mouseleave', () => {
    showUnlock.value = false;
  });
});

onUnmounted(() => {
  window.ipcRenderer.removeAllListeners('update-lyric');
  window.ipcRenderer.removeAllListeners('update-settings');
});
</script>

<style lang="scss" scoped>
.desktop-lyric {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background-color: transparent;
  overflow: hidden;
  position: relative;
  font-family: "HarmonyOS Sans", "PingFang SC", "Microsoft YaHei", sans-serif;
  user-select: none;
  transition: background-color 0.3s ease;

  &:hover {
    .header {
      opacity: 1;
      transform: translateY(0);
    }
  }

  &.locked {
    pointer-events: none; // Allow click-through when locked
    .header {
      display: none;
    }
  }

  .header {
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    color: #fff;
    opacity: 0;
    transform: translateY(-10px);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    background: rgba(40, 40, 40, 0.65);
    backdrop-filter: blur(16px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 100;
    margin: 8px 8px 0;
    border-radius: 12px;
    /* 防止header被拉伸 */
    flex-shrink: 0; 

    .drag-region {
      -webkit-app-region: drag;
      cursor: move;
    }
    
    .no-drag {
      -webkit-app-region: no-drag;
    }

    // 锁定按钮样式
    .unlock-btn {
      width: 28px;
      height: 28px;
      background: rgba(255, 255, 255, 0.1);
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      position: absolute;
      left: 50%;
      transform: translateX(-50%);
      pointer-events: auto; // 确保在锁定状态下可点击
      transition: all 0.2s;
      
      &:hover {
        background: rgba(255, 255, 255, 0.2);
        transform: translateX(-50%) scale(1.1);
      }
    }

    .left {
      display: flex;
      align-items: center;
      gap: 10px;
      flex: 1;
      
      .app-icon {
        width: 28px;
        height: 28px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .title {
        font-size: 13px;
        font-weight: 500;
        opacity: 0.9;
        max-width: 180px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    .center {
      display: flex;
      align-items: center;
      gap: 16px;
      position: absolute;
      left: 50%;
      transform: translateX(-50%);
      
      .control-btn {
        width: 36px;
        height: 36px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.2s;
        color: rgba(255, 255, 255, 0.85);

        &:hover {
          background: rgba(255, 255, 255, 0.1);
          color: #fff;
          transform: scale(1.05);
        }
        
        &:active {
          transform: scale(0.95);
        }

        &.play-btn {
          width: 44px;
          height: 44px;
          background: rgba(255, 255, 255, 0.15);
          color: #fff;
          
          &:hover {
            background: rgba(255, 255, 255, 0.25);
            transform: scale(1.1);
          }
          
          &:active {
            transform: scale(0.95);
          }
        }
      }
    }

    .right {
      display: flex;
      align-items: center;
      gap: 8px;
      flex: 1;
      justify-content: flex-end;

      .action-btn {
        width: 32px;
        height: 32px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.2s;
        color: rgba(255, 255, 255, 0.7);

        &:hover {
          background: rgba(255, 255, 255, 0.1);
          color: #fff;
        }

        &.close-btn:hover {
          background: rgba(239, 68, 68, 0.8);
          color: #fff;
        }
      }
    }
  }

  .lyric-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding-bottom: 20px;
    cursor: default;
    
    // Add a subtle gradient mask at the top/bottom if needed, 
    // but for 2 lines it might be overkill.
    
    .lrc-line {
      line-height: 1.4;
      transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
      width: 100%;
      padding: 0 32px;
      letter-spacing: 0.5px;
      
      &.current {
        margin-bottom: 12px;
        font-weight: 800;
        // Text shadow handles visibility against any background
        text-shadow: 0 2px 10px rgba(0,0,0,0.5), 0 0 2px rgba(0,0,0,0.3);
      }

      &.next {
        font-weight: 500;
        text-shadow: 0 2px 4px rgba(0,0,0,0.5);
      }

      .tlyric {
        font-size: 0.6em;
        opacity: 0.8;
        font-weight: 400;
        margin-top: 4px;
      }
    }
  }
}

// Global styles for popover content
.settings-panel {
  padding: 16px;
  width: 240px;
  color: #fff;
  border-radius: 12px;
  
  &.glass-effect {
    background: rgba(30, 30, 30, 0.85);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }
  
  .panel-header {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 16px;
    opacity: 0.9;
    padding-bottom: 8px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .setting-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    font-size: 13px;
    
    &:last-child { margin-bottom: 0; }
    
    .label {
      opacity: 0.8;
    }
    
    .custom-slider {
      width: 120px;
    }
  }
}

// Transitions
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
