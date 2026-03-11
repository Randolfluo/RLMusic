<template>
  <div class="desktop-lyric" :class="{ locked: isLocked }" :style="containerStyle">
    <!-- Toolbar (Matching the provided image) -->
    <transition name="fade">
      <div class="toolbar-wrapper" v-show="!isLocked || showUnlock">
        <div class="toolbar no-drag">
          <div class="toolbar-group">
            <div class="tool-btn drag-region" title="拖动">
              <n-icon size="18" :component="MusicNoteFilled" />
            </div>
            <div class="tool-btn" @click="control('prev')" title="上一首">
              <n-icon size="18" :component="SkipPreviousFilled" />
            </div>
            <div class="tool-btn" @click="control('toggle')" :title="isPlaying ? '暂停' : '播放'">
              <n-icon size="20" :component="isPlaying ? PauseFilled : PlayArrowFilled" />
            </div>
            <div class="tool-btn" @click="control('next')" title="下一首">
              <n-icon size="18" :component="SkipNextFilled" />
            </div>
          </div>

          <div class="separator"></div>

          <div class="toolbar-group">
            <div class="tool-btn" @click="moveWindow('left')" title="左移">
              <n-icon size="18" :component="ArrowBackIosNewFilled" />
            </div>
            <div class="tool-btn" @click="moveWindow('right')" title="右移">
              <n-icon size="18" :component="ArrowForwardIosFilled" />
            </div>
            <div class="tool-btn" @click="toggleLock" :title="isLocked ? '解锁' : '锁定'">
              <n-icon size="18" :component="isLocked ? LockFilled : LockOpenFilled" />
            </div>
            
            <n-popover trigger="click" placement="bottom" :show-arrow="false" style="padding: 0; background: transparent;">
              <template #trigger>
                <div class="tool-btn" title="设置">
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

            <div class="tool-btn close-btn" @click="close" title="关闭">
              <n-icon size="18" :component="CloseFilled" />
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Lyric Content (Glossy/Stylized) -->
    <div class="lyric-container">
      <div class="lrc-line current" :style="currentLyricStyle">
        <span class="lrc-text">{{ currentLyric || 'RLMusic 听我想听' }}</span>
        <div class="tlyric" v-if="showTranslation && currentTlyric">{{ currentTlyric }}</div>
      </div>
      <div class="lrc-line next" :style="nextLyricStyle" v-if="nextLyric">
        <span class="lrc-text">{{ nextLyric }}</span>
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
  CloseFilled,
  ArrowBackIosNewFilled,
  ArrowForwardIosFilled
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
    '--accent-color': followTheme.value ? themeColor.value : '#0088ff'
  }
});

const currentLyricStyle = computed(() => {
  return {
    fontSize: `${fontSize.value}rem`
  }
});

const nextLyricStyle = computed(() => {
  return {
    fontSize: `${fontSize.value * 0.6}rem`
  }
});

// Actions
const control = (action: string) => {
  window.ipcRenderer.send('desktop-lyric-control', action);
};

const moveWindow = (direction: string) => {
  window.ipcRenderer.send('desktop-lyric-move', direction);
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
    if (data.title !== undefined) songTitle.value = data.title;
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

  document.addEventListener('mouseenter', () => {
    if (isLocked.value) showUnlock.value = true;
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
  align-items: center;
  justify-content: center;
  background-color: transparent;
  overflow: hidden;
  position: relative;
  font-family: "HarmonyOS Sans", "PingFang SC", "Microsoft YaHei", sans-serif;
  user-select: none;
  
  // Background highlight when not locked
  &:not(.locked):hover {
    background-color: rgba(0, 0, 0, 0.1);
  }

  &.locked {
    pointer-events: none;
    .toolbar-wrapper {
      opacity: 0;
    }
  }

  .toolbar-wrapper {
    position: absolute;
    top: 10px;
    left: 0;
    width: 100%;
    display: flex;
    justify-content: center;
    z-index: 100;
    transition: all 0.3s ease;
    pointer-events: auto;
  }

  .toolbar {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 12px;
    background: rgba(30, 30, 30, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    
    .toolbar-group {
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .separator {
      width: 1px;
      height: 16px;
      background: rgba(255, 255, 255, 0.2);
      margin: 0 8px;
    }

    .tool-btn {
      width: 32px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: rgba(255, 255, 255, 0.8);
      cursor: pointer;
      border-radius: 50%;
      transition: all 0.2s;

      &:hover {
        background: rgba(255, 255, 255, 0.15);
        color: #fff;
        transform: scale(1.1);
      }

      &:active {
        transform: scale(0.95);
      }

      &.drag-region {
        -webkit-app-region: drag;
        cursor: move;
      }

      &.close-btn:hover {
        background: #f44336;
      }
    }
  }

  .lyric-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    width: 100%;
    padding: 40px 20px 20px;
    
    .lrc-line {
      width: 100%;
      margin: 8px 0;
      transition: all 0.3s ease;
      
      .lrc-text {
        display: inline-block;
        font-weight: 900;
        letter-spacing: 1px;
        // The Glossy Effect
        background: linear-gradient(to bottom, 
          #ffffff 0%, 
          #ffffff 30%,
          var(--accent-color) 50%, 
          #0044aa 100%
        );
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.8));
        paint-order: stroke fill;
        -webkit-text-stroke: 1px rgba(0, 0, 0, 0.3);
      }

      &.current {
        .lrc-text {
          filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.9));
        }
      }

      &.next {
        opacity: 0.6;
        .lrc-text {
          background: #ccc;
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
        }
      }

      .tlyric {
        font-size: 0.5em;
        font-weight: 600;
        color: #fff;
        opacity: 0.9;
        margin-top: 4px;
        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.8);
      }
    }
  }
}

.no-drag {
  -webkit-app-region: no-drag;
}

// Settings Panel
.settings-panel {
  padding: 16px;
  width: 240px;
  color: #fff;
  border-radius: 12px;
  
  &.glass-effect {
    background: rgba(30, 30, 30, 0.95);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  }
  
  .panel-header {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 16px;
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
