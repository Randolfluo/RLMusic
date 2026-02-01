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
            alignItems: lyricsPosition == 'center' ? 'center' : null,
            transformOrigin:
              lyricsPosition == 'center' ? 'center' : 'center left',
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
      <div class="name">默认歌词位置</div>
      <n-select
        class="set"
        v-model:value="lyricsPosition"
        :options="lyricsPositionOptions"
      />
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
        @click="changeMusicFrequency"
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
  searchHistory,
  bottomLyricShow,
  lrcMousePause,
  playerStyle,
  lyricsFontSize,
  lyricsPosition,
  lyricsBlock,
  lyricsBlur,
  musicFrequency,
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



// 歌词位置
const lyricsPositionOptions = [
  {
    label: "居左",
    value: "left",
  },
  {
    label: "居中",
    value: "center",
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

// 音乐频谱提醒
const changeMusicFrequency = () => {
  if (musicFrequency.value) {
    $dialog.warning({
      class: "s-dialog",
      title: "实验性功能",
      content: "确认开启音乐频谱？将于刷新后生效",
      positiveText: "开启",
      negativeText: "取消",
      onMaskClick: () => {
        musicFrequency.value = false;
      },
      onPositiveClick: () => {
        musicFrequency.value = true;
        location.reload();
      },
      onNegativeClick: () => {
        musicFrequency.value = false;
      },
    });
  }
};
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
