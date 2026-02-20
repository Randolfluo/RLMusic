<script setup>
import { storeToRefs } from "pinia";
import { settingStore, userStore } from "@/store";
import { Config, MusicOne, Effects } from "@icon-park/vue-next";
import { changePassword } from "@/api/user";

const message = useMessage();
const showPasswordModal = ref(false);
const passwordLoading = ref(false);
const passwordFormRef = ref(null);
const passwordForm = reactive({
  oldPassword: "",
  newPassword: "",
  confirmPassword: "",
});

const validateConfirmPassword = (rule, value) => {
  if (value !== passwordForm.newPassword) {
    return new Error("两次输入的密码不一致");
  }
  return true;
};

const passwordRules = {
  oldPassword: [{ required: true, message: "请输入旧密码", trigger: "blur" }],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 4, max: 20, message: "密码长度在 4 到 20 个字符", trigger: "blur" },
  ],
  confirmPassword: [
    { required: true, message: "请再次输入新密码", trigger: "blur" },
    { validator: validateConfirmPassword, trigger: "blur" },
  ],
};

const handlePasswordSubmit = (e) => {
  e.preventDefault();
  passwordFormRef.value?.validate((errors) => {
    if (!errors) {
      passwordLoading.value = true;
      changePassword({
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword,
      })
        .then(() => {
          message.success("密码修改成功");
          showPasswordModal.value = false;
          // 重置表单
          passwordForm.oldPassword = "";
          passwordForm.newPassword = "";
          passwordForm.confirmPassword = "";
        })
        .catch((err) => {
          console.error(err);
        })
        .finally(() => {
          passwordLoading.value = false;
        });
    }
  });
};

const setting = settingStore();
const user = userStore();
const {
  themeColor,
  searchHistory,
  bottomLyricShow,
  podcastMode,
  lrcMousePause,
  playerStyle,
  playerBgBlur,
  lyricsFontSize,
  lyricsBlock,
  lyricsBlur,
  musicFrequency,
  musicFrequencyScale,
  particleEffect,
  particleLimit,
  desktopLyricShow,
  desktopLyricFontSize,
  desktopLyricLock,
  desktopLyricFollowTheme,
  desktopLyricShowTranslation,
} = storeToRefs(setting);

// 桌面歌词设置
const handleDesktopLyricToggle = (val) => {
  if (window.ipcRenderer) {
    if (val) {
      window.ipcRenderer.send("open-desktop-lyric");
    } else {
      window.ipcRenderer.send("close-desktop-lyric");
    }
  }
};

const handleDesktopLyricLock = (val) => {
  if (window.ipcRenderer) {
    window.ipcRenderer.send("lock-desktop-lyric", val);
  }
};

const handleDesktopLyricSettings = () => {
  if (window.ipcRenderer) {
    window.ipcRenderer.send("update-desktop-lyric-settings", {
      fontSize: desktopLyricFontSize.value,
      followTheme: desktopLyricFollowTheme.value,
      themeColor: themeColor.value,
      showTranslation: desktopLyricShowTranslation.value,
    });
  }
};

watch(themeColor, () => {
  if (desktopLyricFollowTheme.value) {
    handleDesktopLyricSettings();
  }
});

// 歌词滚动位置
const lyricsBlockOptions = [
  {
    label: "靠近顶部",
    value: "start",
  },
  {
    label: "水平居中",
    value: "center",
  },
];

// 播放器样式
const playerStyleOptions = [
  {
    label: "封面模式",
    value: "cover",
  },
  {
    label: "唱片模式",
    value: "record",
  },
];
</script>

<template>
  <div class="setting-view">
    <div class="setting-header">
      <h1 class="page-title">全局设置</h1>
      <p class="page-subtitle">自定义您的音乐体验</p>
    </div>

    <div class="setting-content">
      <div class="setting-section">
        <h2 class="section-title">
          <n-icon :component="Config" /> 基础设置
        </h2>
        <div class="setting-grid">
          <n-card class="setting-card full-width" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">主题颜色</div>
                <div class="desc">个性化应用的主题色调</div>
              </div>
              <n-color-picker
                class="control color-picker"
                v-model:value="themeColor"
                :show-alpha="false"
                :swatches="[
                  '#009688',
                  '#18a058',
                  '#2080f0',
                  '#f0a020',
                  '#d03050',
                  '#ffc0cb',
                ]"
                size="small"
              />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">搜索历史</div>
                <div class="desc">保存并显示最近的搜索记录</div>
              </div>
              <n-switch v-model:value="searchHistory" :round="false" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">底栏歌词</div>
                <div class="desc">播放时在底部显示歌词</div>
              </div>
              <n-switch v-model:value="bottomLyricShow" :round="false" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">账号安全</div>
                <div class="desc">修改您的登录密码</div>
              </div>
              <n-button secondary type="primary" size="small" @click="showPasswordModal = true">
                修改密码
              </n-button>
            </div>
          </n-card>
        </div>
      </div>

      <div class="setting-section">
        <h2 class="section-title">
          <n-icon :component="MusicOne" /> 歌词与播放器
        </h2>
        <div class="setting-grid">
          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">播客模式</div>
                <div class="desc">启用播客优化的播放体验，在歌曲播放前会播放开场白</div>
              </div>
              <n-switch v-model:value="podcastMode" :round="false" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">智能暂停滚动</div>
                <div class="desc">鼠标悬停时暂停歌词滚动</div>
              </div>
              <n-switch v-model:value="lrcMousePause" :round="false" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">播放器样式</div>
                <div class="desc">选择全屏播放器的视觉风格</div>
              </div>
              <n-select
                class="control"
                v-model:value="playerStyle"
                :options="playerStyleOptions"
                size="small"
              />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">背景模糊</div>
                <div class="desc">全屏播放器背景高斯模糊</div>
              </div>
              <n-switch v-model:value="playerBgBlur" :round="false" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">歌词滚动位置</div>
                <div class="desc">高亮歌词在屏幕中的位置</div>
              </div>
              <n-select
                class="control"
                v-model:value="lyricsBlock"
                :options="lyricsBlockOptions"
                size="small"
              />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">歌词模糊</div>
                <div class="desc">非高亮歌词应用模糊效果</div>
              </div>
              <n-switch v-model:value="lyricsBlur" :round="false" />
            </div>
          </n-card>

          <n-card class="setting-card full-width lyric-preview-card" :bordered="false">
             <div class="card-inner vertical">
               <div class="header">
                 <div class="info">
                   <div class="name">歌词文本大小</div>
                   <div class="desc">调整歌词显示的字体大小</div>
                 </div>
                 <n-slider
                   v-model:value="lyricsFontSize"
                   :tooltip="false"
                   :max="3.4"
                   :min="2.2"
                   :step="0.01"
                   class="slider-control"
                 />
               </div>
               
               <div class="preview-box" :class="{ blur: lyricsBlur }">
                 <div class="lrc-line prev">这是一句歌词</div>
                 <div class="lrc-line current" :style="{ fontSize: lyricsFontSize + 'vh' }">
                   This is a lyric
                 </div>
                 <div class="lrc-line next" :style="{ fontSize: (lyricsFontSize - 0.4) + 'vh' }">
                   下一句歌词预览
                 </div>
               </div>
             </div>
          </n-card>
        </div>
      </div>

      <div class="setting-section">
        <h2 class="section-title">
          <n-icon :component="Effects" /> 视觉效果
        </h2>
        <div class="setting-grid">
          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">音乐频谱</div>
                <div class="desc">播放时显示可视化频谱</div>
              </div>
              <n-switch v-model:value="musicFrequency" :round="false" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">粒子效果</div>
                <div class="desc">背景动态粒子动画</div>
              </div>
              <n-switch v-model:value="particleEffect" :round="false" />
            </div>
          </n-card>

          <transition name="fade">
            <n-card class="setting-card full-width" :bordered="false" v-if="particleEffect">
              <div class="card-inner vertical">
                <div class="info">
                  <div class="name">粒子密度</div>
                  <div class="desc">调整背景粒子的数量 ({{ particleLimit }})</div>
                </div>
                <n-slider
                  v-model:value="particleLimit"
                  :tooltip="false"
                  :max="200"
                  :min="10"
                  :step="10"
                  class="slider-control mt-4"
                />
              </div>
            </n-card>
          </transition>

          <transition name="fade">
            <n-card class="setting-card full-width" :bordered="false" v-if="musicFrequency">
              <div class="card-inner vertical">
                <div class="info">
                  <div class="name">频谱跳动幅度</div>
                  <div class="desc">调整频谱的灵敏度 ({{ musicFrequencyScale }})</div>
                </div>
                <n-slider
                  v-model:value="musicFrequencyScale"
                  :tooltip="false"
                  :max="200"
                  :min="10"
                  :step="10"
                  class="slider-control mt-4"
                />
              </div>
            </n-card>
          </transition>
        </div>
      </div>

      <div class="setting-section">
        <h2 class="section-title">
          <n-icon :component="Config" /> 桌面歌词
        </h2>
        <div class="setting-grid">
          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">开启桌面歌词</div>
                <div class="desc">在屏幕上显示悬浮歌词</div>
              </div>
              <n-switch v-model:value="desktopLyricShow" :round="false" @update:value="handleDesktopLyricToggle" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">锁定歌词窗口</div>
                <div class="desc">锁定后无法拖动和调整大小</div>
              </div>
              <n-switch v-model:value="desktopLyricLock" :round="false" @update:value="handleDesktopLyricLock" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">歌词跟随主题</div>
                <div class="desc">歌词颜色跟随应用主题色</div>
              </div>
              <n-switch v-model:value="desktopLyricFollowTheme" :round="false" @update:value="handleDesktopLyricSettings" />
            </div>
          </n-card>

          <n-card class="setting-card" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">显示双语歌词</div>
                <div class="desc">在桌面歌词中显示翻译</div>
              </div>
              <n-switch v-model:value="desktopLyricShowTranslation" :round="false" @update:value="handleDesktopLyricSettings" />
            </div>
          </n-card>

          <n-card class="setting-card full-width" :bordered="false">
            <div class="card-inner">
              <div class="info">
                <div class="name">歌词字体大小</div>
                <div class="desc">调整桌面歌词的显示大小: {{ desktopLyricFontSize }}rem</div>
              </div>
              <n-slider
                class="control slider"
                v-model:value="desktopLyricFontSize"
                :step="0.1"
                :min="1.0"
                :max="5.0"
                :tooltip="false"
                @update:value="handleDesktopLyricSettings"
              />
            </div>
          </n-card>
        </div>
      </div>
    </div>

    <n-modal v-model:show="showPasswordModal" preset="card" title="修改密码" style="width: 400px">
      <n-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-placement="left"
        label-width="80"
      >
        <n-form-item label="旧密码" path="oldPassword">
          <n-input
            v-model:value="passwordForm.oldPassword"
            type="password"
            placeholder="请输入旧密码"
            show-password-on="click"
          />
        </n-form-item>
        <n-form-item label="新密码" path="newPassword">
          <n-input
            v-model:value="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password-on="click"
          />
        </n-form-item>
        <n-form-item label="确认密码" path="confirmPassword">
          <n-input
            v-model:value="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password-on="click"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px">
          <n-button @click="showPasswordModal = false">取消</n-button>
          <n-button type="primary" :loading="passwordLoading" @click="handlePasswordSubmit">
            确认修改
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<style lang="scss" scoped>
.setting-view {
  padding: 40px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  
  @media (max-width: 768px) {
    padding: 24px;
  }
}

/* Header */
.setting-header {
  margin-bottom: 48px;
  animation: fade-in-down 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
  
  .page-title {
    font-size: 42px;
    font-weight: 800;
    margin: 0 0 12px 0;
    background: linear-gradient(120deg, var(--n-color-primary) 0%, #a78bfa 100%);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: -1px;
  }
  
  .page-subtitle {
    font-size: 16px;
    color: var(--n-text-color-3);
    margin: 0;
    opacity: 0.8;
  }
}

/* Content Layout */
.setting-content {
  display: flex;
  flex-direction: column;
  gap: 56px;
  animation: fade-in-up 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.1s backwards;
}

.setting-section {
  .section-title {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 20px;
    font-weight: 700;
    margin: 0 0 24px 0;
    color: var(--n-text-color);
    opacity: 0.9;
    
    .n-icon {
      color: var(--n-color-primary);
      background: rgba(var(--n-color-primary-rgb), 0.1);
      padding: 8px;
      border-radius: 12px;
    }
  }
}

.setting-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 20px;
}

/* Setting Card */
.setting-card {
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.6) !important;
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.3) !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  overflow: hidden;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
    background: rgba(255, 255, 255, 0.8) !important;
    border-color: var(--n-color-primary) !important;
  }
  
  &.full-width {
    grid-column: 1 / -1;
  }
  
  .card-inner {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px;
    height: 100%;
    
    &.vertical {
      flex-direction: column;
      align-items: stretch;
      gap: 20px;
      
      .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }
    }
  }
  
  .info {
    .name {
      font-size: 16px;
      font-weight: 600;
      color: var(--n-text-color);
      margin-bottom: 4px;
    }
    
    .desc {
      font-size: 13px;
      color: var(--n-text-color-3);
      line-height: 1.4;
    }
  }
  
  .control {
    width: 140px;
  }
  
  .color-picker {
    width: auto;
  }
}

/* Lyric Preview */
.lyric-preview-card {
  .preview-box {
    background: rgba(0, 0, 0, 0.03);
    border-radius: 16px;
    padding: 32px;
    text-align: center;
    transition: all 0.3s;
    border: 1px solid rgba(0, 0, 0, 0.05);
    
    &.blur {
      .lrc-line:not(.current) {
        filter: blur(3px);
        opacity: 0.4;
      }
    }
    
    .lrc-line {
      transition: all 0.4s ease;
      margin: 12px 0;
      
      &.prev, &.next {
        font-size: 16px;
        color: var(--n-text-color-3);
        opacity: 0.6;
      }
      
      &.current {
        font-weight: 800;
        color: var(--n-color-primary);
        background: linear-gradient(120deg, var(--n-color-primary) 0%, #a78bfa 100%);
        -webkit-background-clip: text;
        background-clip: text;
        -webkit-text-fill-color: transparent;
        transform: scale(1.05);
      }
    }
  }
}

.slider-control {
  margin-top: 8px;
  width: 100%;
}

.mt-4 {
  margin-top: 16px;
}

/* Animations */
@keyframes fade-in-down {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
