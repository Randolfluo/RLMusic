<template>
  <div class="setting">


    <div class="title">全局设置</div>
    <n-h6 prefix="bar"> 基础设置 </n-h6>
    <n-card class="set-item">
      <div class="name">明暗模式</div>
      <n-select class="set" v-model:value="theme" :options="darkOptions" />
    </n-card>
    <n-card class="set-item">
      <div class="name">明暗模式跟随系统</div>
      <n-switch v-model:value="themeAuto" :round="false" />
    </n-card>
    <n-card class="set-item">
      <div class="name">主题颜色</div>
      <n-color-picker
        class="set"
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
      />
    </n-card>
    <n-card class="set-item">
      <div class="name">显示搜索历史</div>
      <n-switch v-model:value="searchHistory" :round="false" />
    </n-card>
    <n-card class="set-item">
      <div class="name">
        显示底栏歌词
        <span class="tip">是否在播放时显示歌词</span>
      </div>
      <n-switch v-model:value="bottomLyricShow" :round="false" />
    </n-card>


    <n-h6 prefix="bar"> 歌词设置 </n-h6>
    <n-card class="set-item">
      <div class="name">
        智能暂停滚动
        <span class="tip">鼠标移入歌词区域是否暂停滚动</span>
      </div>
      <n-switch v-model:value="lrcMousePause" :round="false" />
    </n-card>
    <n-card class="set-item">
      <div class="name">播放器样式</div>
      <n-select
        class="set"
        v-model:value="playerStyle"
        :options="playerStyleOptions"
      />
    </n-card>
    <n-card class="set-item">
      <div class="name">
        播放器背景模糊
        <span class="tip">开启后全屏播放器背景将进行高斯模糊处理</span>
      </div>
      <n-switch v-model:value="playerBgBlur" :round="false" />
    </n-card>
    <n-card
      class="set-item"
      :content-style="{
        flexDirection: 'column',
        alignItems: 'flex-start',
      }"
    >
      <div class="name">歌词文本大小</div>
      <n-slider
        v-model:value="lyricsFontSize"
        :tooltip="false"
        :max="3.4"
        :min="2.2"
        :step="0.01"
        :marks="{
          2.2: '最小',
          2.8: '默认',
          3.4: '最大',
        }"
      />
      <div :class="lyricsBlur ? 'more blur' : 'more'">
        <div
          v-for="n in 3"
          :key="n"
          :class="n === 2 ? 'lrc on' : 'lrc'"
          :style="{
            margin: n === 2 ? '12px 0' : null,
          }"
        >
          <span :style="{ fontSize: lyricsFontSize + 'vh' }"
            >这是一句歌词
          </span>
          <span :style="{ fontSize: lyricsFontSize - 0.4 + 'vh' }"
            >This is a lyric
          </span>
        </div>
      </div>
    </n-card>
    <n-card class="set-item">
      <div class="name">
        歌词滚动位置
        <span class="tip">歌词高亮时所处的位置</span>
      </div>
      <n-select
        class="set"
        v-model:value="lyricsBlock"
        :options="lyricsBlockOptions"
      />
    </n-card>
    <n-card class="set-item">
      <div class="name">
        歌词模糊
      </div>
      <n-switch v-model:value="lyricsBlur" :round="false" />
    </n-card>
    <n-card class="set-item">
      <div class="name">
        显示音乐频谱
      </div>
      <n-switch
        v-model:value="musicFrequency"
        :round="false"
      />
    </n-card>
    <n-card class="set-item">
      <div class="name">
        显示粒子效果
      </div>
      <n-switch
        v-model:value="particleEffect"
        :round="false"
      />
    </n-card>
    <n-card
      class="set-item"
      :content-style="{
        'flex-direction': 'column',
        'align-items': 'flex-start',
      }"
      v-if="particleEffect"
    >
      <div class="name">
        粒子数量
        <span class="tip">调整背景粒子的密度</span>
      </div>
      <n-slider
        v-model:value="particleLimit"
        :tooltip="false"
        :max="200"
        :min="10"
        :step="10"
        :marks="{
          10: '少',
          50: '默认',
          200: '多',
        }"
      />
    </n-card>
    <n-card
      class="set-item"
      :content-style="{
        'flex-direction': 'column',
        'align-items': 'flex-start',
      }"
      v-if="musicFrequency"
    >
      <div class="name">
        频谱跳动幅度
        <span class="tip">调整频谱显示的跳动高度</span>
      </div>
      <n-slider
        v-model:value="musicFrequencyScale"
        :tooltip="false"
        :max="200"
        :min="10"
        :step="10"
        :marks="{
          10: '低',
          90: '默认',
          200: '高',
        }"
      />
    </n-card>
  </div>
</template>

<script setup>
import { storeToRefs } from "pinia";
import { settingStore, userStore } from "@/store";
const setting = settingStore();
const user = userStore();
const {
  theme,
  themeAuto,
  themeColor,
  searchHistory,
  bottomLyricShow,
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
} = storeToRefs(setting);

// 深浅模式
const darkOptions = [
  {
    label: "浅色模式",
    value: "light",
  },
  {
    label: "深色模式",
    value: "dark",
  },
];

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

<style lang="scss" scoped>
.setting {
  padding: 0 10vw;
  max-width: 1200px;
  margin: 0 auto;
  @media (max-width: 768px) {
    padding: 0;
  }
  .title {
    margin-top: 30px;
    margin-bottom: 20px;
    font-size: 40px;
    font-weight: bold;
  }
  .n-h {
    padding-left: 16px;
    font-size: 20px;
    margin-left: 4px;
  }
  .set-item {
    width: 100%;
    border-radius: 8px;
    margin-bottom: 12px;
    :deep(.n-card__content) {
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: space-between;
      .name {
        font-size: 16px;
        display: flex;
        flex-direction: column;
        padding-right: 20px;
        .tip {
          font-size: 12px;
          opacity: 0.8;
        }
      }
      .set {
        width: 200px;
        @media (max-width: 768px) {
          width: 140px;
          min-width: 140px;
        }
      }
      .more {
        padding: 12px;
        border-radius: 8px;
        background-color: var(--n-border-color);
        width: 100%;
        margin-top: 12px;
        box-sizing: border-box;
        &.blur {
          .lrc {
            filter: blur(2px);
            &.on {
              filter: blur(0);
            }
          }
        }
        .lrc {
          opacity: 0.6;
          display: flex;
          flex-direction: column;
          transform: scale(0.95);
          transform-origin: left;
          transition: all 0.3s;
          &.on {
            font-weight: bold;
            opacity: 1;
            transform: scale(1.05);
          }
        }
      }
    }
  }
}
</style>
