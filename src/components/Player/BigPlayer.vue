<!-- 全屏播放器 -->
<template>
  <Transition name="up">
    <div
      v-show="music.showBigPlayer"
      class="bplayer"
      :style="
        bgCover
          ? 'background-image: url(' +
            bgCover +
            ')'
          : ''
      "
    >
      <div
        class="gray"
        :style="{ backdropFilter: setting.playerBgBlur ? 'blur(80px)' : 'none' }"
      />
      <canvas class="particle-canvas" ref="particleCanvas" />
      <n-icon
        class="close"
        size="40"
        :component="KeyboardArrowDownFilled"
        @click="music.setBigPlayerState(false)"
      />
      <n-icon
        class="screenfull"
        size="36"
        :component="screenfullIcon"
        @click="screenfullChange"
      />
      <div
        :class="
          music.getPlaySongLyric[0]
            ? 'all'
            : 'all noLrc'
        "
      >
        <!-- 提示文本 -->
        <Transition name="lrc">
          <div class="tip" v-show="lrcMouseStatus">
            <n-text>点击选中的歌词以调整播放进度</n-text>
          </div>
        </Transition>
        <div class="left">
          <PlayerCover v-if="setting.playerStyle === 'cover'" />
          <PlayerRecord v-else />
        </div>
        <div
          class="right"
          @mouseenter="menuShow = true"
          @mouseleave="menuShow = false"
        >
          <Transition name="lrc">
            <div
              class="lrcShow"
              v-if="
                music.getPlaySongLyric[0]
              "
            >
              <div class="data" v-show="setting.playerStyle === 'record'">
                <div class="name text-hidden">
                  <span>{{
                    music.getPlaySongData
                      ? music.getPlaySongData.name
                      : "暂无歌曲"
                  }}</span>
                  <span
                    v-if="music.getPlaySongData && music.getPlaySongData.alia"
                    >{{ music.getPlaySongData.alia[0] }}</span
                  >
                </div>
                <div
                  class="artists text-hidden"
                  v-if="music.getPlaySongData && music.getPlaySongData.artist"
                >
                  <span
                    class="artist"
                    v-for="(item, index) in music.getPlaySongData.artist"
                    :key="item"
                  >
                    <span>{{ item.name }}</span>
                    <span
                      v-if="index != music.getPlaySongData.artist.length - 1"
                      >/</span
                    >
                  </span>
                </div>
              </div>
              <div
                :class="
                  setting.playerStyle === 'cover'
                    ? 'lrc-all cover'
                    : 'lrc-all record'
                "
                v-if="music.getPlaySongLyric[0]"
                :style="
                  setting.lyricsPosition === 'center'
                    ? { textAlign: 'center', paddingRight: '0' }
                    : null
                "
                @mouseenter="
                  lrcMouseStatus = setting.lrcMousePause ? true : false
                "
                @mouseleave="lrcAllLeave"
              >
                <!-- 提示文本 -->
                <div class="tip">
                  <n-text>点击选中的歌词以调整播放进度</n-text>
                </div>
                <div class="placeholder"></div>
                <div
                  :class="
                    music.getPlaySongLyricIndex == index ? 'lrc on' : 'lrc'
                  "
                  :style="{ marginBottom: setting.lyricsFontSize - 1.6 + 'vh' }"
                  v-for="(item, index) in music.getPlaySongLyric"
                  :key="item"
                  :id="'lrc' + index"
                  @click="jumpTime(item.time)"
                >
                  <div
                    class="curr-time"
                    :style="{ fontSize: setting.lyricsFontSize * 0.5 + 'vh' }"
                  >
                    {{ formatTime(item.time) }}
                  </div>
                  <div
                    :class="setting.lyricsBlur ? 'lrc-text blur' : 'lrc-text'"
                    :style="{
                      transformOrigin:
                        setting.lyricsPosition === 'center' ? 'center' : null,
                      filter: setting.lyricsBlur
                        ? `blur(${getFilter(
                            music.getPlaySongLyricIndex,
                            index
                          )}px)`
                        : null,
                    }"
                  >
                    <span
                      class="lyric"
                      :style="{ fontSize: setting.lyricsFontSize + 'vh' }"
                    >
                      {{ item.lyric }}
                    </span>
                    <span
                      v-show="
                        music.getPlaySongTransl &&
                        setting.getShowTransl &&
                        item.lyricFy
                      "
                      :style="{ fontSize: setting.lyricsFontSize - 0.4 + 'vh' }"
                      class="lyric-fy"
                    >
                      {{ item.lyricFy }}</span
                    >
                  </div>
                </div>
                <div class="placeholder"></div>
              </div>
              <div
                :class="menuShow ? 'menu show' : 'menu'"
              >
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-icon
                      class="style-switch"
                      :component="setting.playerStyle === 'cover' ? DiscFullOutlined : ImageOutlined"
                      @click="changePlayerStyle"
                    />
                  </template>
                  {{ setting.playerStyle === 'cover' ? '切换为唱片模式' : '切换为封面模式' }}
                </n-tooltip>

                <n-tooltip trigger="hover">
                  <template #trigger>
                    <div class="speed-control">
                       <n-icon :component="SlowMotionVideoRound" />
                       <div class="speed-popup">
                          <div class="val">{{ music.getPlayRate }}x</div>
                          <n-slider
                             v-model:value="music.persistData.playRate"
                             :tooltip="false"
                             :min="0.5"
                             :max="2.0"
                             :step="0.1"
                             vertical
                             class="slider"
                             @update:value="(v) => music.setPlayRate(v)"
                             @click.stop
                           />
                       </div>
                    </div>
                  </template>
                  播放速度
                </n-tooltip>

                <n-icon
                  v-if="music.getPlaySongTransl"
                  :class="setting.getShowTransl ? 'open' : ''"
                  :component="GTranslateFilled"
                  @click="setting.setShowTransl(!setting.getShowTransl)"
                />
                <div class="lyric-offset-control">
                  <n-icon class="btn" :component="RemoveOutlined" @click="changeOffset(-0.1)" />
                  <span class="text">{{ offsetText }}</span>
                  <n-icon class="btn" :component="AddOutlined" @click="changeOffset(0.1)" />
                </div>
              </div>
            </div>
          </Transition>
        </div>
      </div>
      <div class="canvas">
        <canvas v-if="setting.musicFrequency" class="avBars" ref="avBars" />
      </div>
    </div>
  </Transition>
</template>


<script setup>
import {
  KeyboardArrowDownFilled,
  GTranslateFilled,
  FullscreenRound,
  FullscreenExitRound,
  AddOutlined,
  RemoveOutlined,
  DiscFullOutlined,
  ImageOutlined,
  SlowMotionVideoRound,
} from "@vicons/material";
import { musicStore, settingStore } from "@/store";
import { useRouter } from "vue-router";
import MusicFrequency from "@/utils/MusicFrequency";
import ParticleEffect from "@/utils/ParticleEffect";
import PlayerRecord from "./PlayerRecord.vue";
import PlayerCover from "./PlayerCover.vue";
import screenfull from "screenfull";
import { getSongCover } from "@/api/song"; // 1. 导入 getSongCover

const router = useRouter();
const music = musicStore();
const setting = settingStore();

// 计算背景图
const bgCover = computed(() => {
  if (music.getPlaySongData?.id) {
    return getSongCover(music.getPlaySongData.id);
  }
  return "";
});

// 工具栏显隐
const menuShow = ref(false);

// ... rest of script ...


// 切换播放器样式
const changePlayerStyle = () => {
  setting.playerStyle = setting.playerStyle === "cover" ? "record" : "cover";
};

// 音乐频谱
const avBars = ref(null);
const musicFrequency = ref(null);

// 粒子效果
const particleCanvas = ref(null);
const particleEffect = ref(null);

// 格式化时间
const formatTime = (time) => {
  const min = Math.floor(time / 60);
  const sec = Math.floor(time % 60).toString().padStart(2, '0');
  return `${min}:${sec}`;
};

// 歌词模糊数值
const getFilter = (lrcIndex, index) => {
  if (lrcIndex >= index) {
    return lrcIndex - index;
  } else {
    return index - lrcIndex;
  }
};

// 点击歌词跳转
const jumpTime = (time) => {
  lrcMouseStatus.value = false;
  if (typeof $player !== 'undefined' && $player) $player.currentTime = time;
};

// 鼠标移出歌词区域
const lrcMouseStatus = ref(false);
const lrcAllLeave = () => {
  lrcMouseStatus.value = false;
  lyricsScroll(music.getPlaySongLyricIndex);
};

// 全屏切换
const screenfullIcon = shallowRef(FullscreenRound);
const screenfullChange = () => {
  if (screenfull.isEnabled) {
    screenfull.toggle();
    screenfullIcon.value = screenfull.isFullscreen
      ? FullscreenRound
      : FullscreenExitRound;
    // 延迟一段时间执行列表滚动
    setTimeout(() => {
      lrcMouseStatus.value = false;
      lyricsScroll(music.getPlaySongLyricIndex);
    }, 500);
  }
};

// 前往评论
const toComment = () => {
  music.setBigPlayerState(false);
  router.push({
    path: "/comment",
    query: {
      id: music.getPlaySongData ? music.getPlaySongData.id : null,
    },
  });
};

// 歌词滚动
const lyricsScroll = (index) => {
  const type = setting.lyricsBlock;
  const el = document.getElementById(
    `lrc${type === "center" ? index : index - 2}`
  );
  if (el && !lrcMouseStatus.value) {
    el.scrollIntoView({
      behavior: "smooth",
      block: type,
    });
  }
};

// 歌词偏移控制
const offsetText = computed(() => {
  const offset = music.getLyricOffset;
  if (offset === 0) return "0.0s";
  return offset > 0 ? `延迟 ${offset.toFixed(1)}s` : `提前 ${Math.abs(offset).toFixed(1)}s`;
});

const changeOffset = (val) => {
  const current = music.getLyricOffset;
  const next = parseFloat((current + val).toFixed(1));
  music.setLyricOffset(next);
};


// 初始化频谱
const initMusicFrequency = () => {
  if (typeof $player !== 'undefined' && $player) {
    if (musicFrequency.value) musicFrequency.value.dispose();
    $player.crossOrigin = "anonymous";
    musicFrequency.value = new MusicFrequency(
      avBars.value,
      $player,
      setting.getThemeColor,
      null,
      50,
      null,
      5,
      setting.musicFrequencyScale
    );
    musicFrequency.value.drawSpectrum();
  }
};

// 初始化粒子效果
const initParticleEffect = () => {
  if (typeof window !== 'undefined') {
     if (particleEffect.value) {
       particleEffect.value.dispose();
       particleEffect.value = null;
     }
     if (setting.particleEffect && particleCanvas.value) {
        particleEffect.value = new ParticleEffect(particleCanvas.value, setting.particleLimit);
        particleEffect.value.start();
     }
  }
};

onMounted(() => {
  nextTick(() => {
    if (setting.musicFrequency) {
      initMusicFrequency();
    }
    if (music.showBigPlayer) {
      initParticleEffect();
    }
    // 滚动歌词
    lyricsScroll(music.getPlaySongLyricIndex);
  });
});

onUnmounted(() => {
  if (particleEffect.value) {
    particleEffect.value.dispose();
  }
});

// 监听主题颜色变化
watch(
  () => setting.getThemeColor,
  (val) => {
    if (musicFrequency.value) {
      musicFrequency.value.setColor(val);
    }
  }
);

// 监听频谱幅度变化
watch(
  () => setting.musicFrequencyScale,
  (val) => {
    if (musicFrequency.value) {
      musicFrequency.value.setScale(val);
    }
  }
);

// 监听是否显示频谱
watch(
  () => setting.musicFrequency,
  (val) => {
    if (val) {
      nextTick(() => {
        initMusicFrequency();
      });
    } else {
      if (musicFrequency.value) {
        musicFrequency.value.dispose();
        musicFrequency.value = null;
      }
    }
  }
);

// 监听是否显示粒子效果
watch(
  () => setting.particleEffect,
  (val) => {
    if (val) {
      nextTick(() => {
        if (music.showBigPlayer) {
          initParticleEffect();
        }
      });
    } else {
      if (particleEffect.value) {
        particleEffect.value.dispose();
        particleEffect.value = null;
      }
    }
  }
);

// 监听粒子数量变化
watch(
  () => setting.particleLimit,
  (val) => {
    if (particleEffect.value) {
      particleEffect.value.setLimit(val);
    }
  }
);

// 监听页面是否打开
watch(
  () => music.showBigPlayer,
  (val) => {
    if (val) {
      console.log("开启播放器", music.getPlaySongLyricIndex);
      nextTick(() => {
        lyricsScroll(music.getPlaySongLyricIndex);
        initParticleEffect();
      });
    } else {
      if (particleEffect.value) {
        particleEffect.value.stop();
      }
    }
  }
);

// 监听歌词滚动
watch(
  () => music.getPlaySongLyricIndex,
  (val) => lyricsScroll(val)
);
</script>

<style lang="scss" scoped>
.up-enter-active,
.up-leave-active {
  transform: translateY(0);
  transition: all 0.5s cubic-bezier(0.65, 0.05, 0.36, 1);
}
.up-enter-from,
.up-leave-to {
  transform: translateY(100%);
}
.lrc-enter-active,
.lrc-leave-active {
  transition: opacity 0.3s ease;
}
.lrc-enter-active {
  transition-delay: 0.3s;
}
.lrc-enter-from,
.lrc-leave-to {
  opacity: 0;
}
.bplayer {
  :deep(.n-slider-rail) {
    background-color: rgba(255, 255, 255, 0.2);
    .n-slider-rail__fill {
      background-color: #fff;
    }
  }
  :deep(.n-slider-handle) {
    background-color: #fff;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  }

  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999;
  overflow: hidden;
  color: #ffffff;
  background-repeat: no-repeat;
  background-size: 150% 150%;
  background-position: center;
  display: flex;
  justify-content: center;
  &::after {
    // content: "";
    position: absolute;
    top: 0;
    left: calc(50% - 2px);
    height: 100%;
    width: 4px;
    background-color: red;
  }
  .gray {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: #00000060;
    // backdrop-filter: blur(80px);
    z-index: -1;
    transition: backdrop-filter 0.3s;
  }
  .particle-canvas {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 0;
    pointer-events: none;
  }
  .close,
  .screenfull {
    position: absolute;
    top: 24px;
    right: 24px;
    opacity: 0.3;
    color: #fff;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.3s;
    z-index: 2;
    &:hover {
      background-color: #ffffff20;
      transform: scale(1.05);
      opacity: 1;
    }
    &:active {
      transform: scale(1);
    }
  }
  .screenfull {
    right: 80px;
    padding: 2px;
    @media (max-width: 768px) {
      display: none;
    }
  }
  .all {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    transition: all 0.3s ease-in-out;
    position: relative;
    &.noLrc {
      .left {
        width: 100%;
        padding-right: 0;
        transform: none;
        align-items: center;
      }
      @media (max-width: 768px) {
        .left {
          width: 100%;
          display: flex !important;
          transform: none;
          align-items: center;
        }
        .right {
          display: none !important;
        }
      }
    }
    @media (max-width: 768px) {
      .left {
        display: none !important;
      }
      .right {
        padding: 0 2vw;
        .lrcShow {
          .lrc-all {
            height: 70vh !important;
            // padding-right: 16% !important;
            margin-right: 0 !important;
          }
          .data,
          .menu {
            display: block !important;
            opacity: 1 !important;
          }
        }
      }
    }
    .tip {
      position: absolute;
      top: 24px;
      left: calc(50% - 150px);
      width: 300px;
      height: 40px;
      border-radius: 25px;
      background-color: #ffffff20;
      backdrop-filter: blur(20px);
      display: flex;
      align-items: center;
      justify-content: center;
      span {
        color: #ffffffc7;
      }
    }
    .left {
      // flex: 1;
      // padding: 0 4vw;
      width: 50%;
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      justify-content: center;
      transition: all 0.3s ease-in-out;
      padding-right: 3.8vw;
      box-sizing: border-box;
      cursor: pointer; // 添加鼠标指针样式
    }
    .right {
      flex: 1;
      height: 100%;
      padding-left: 2vw;
      .lrcShow {
        height: 100%;
        display: flex;
        justify-content: center;
        flex-direction: column;
        .data {
          padding: 0 20px;
          margin-bottom: 8px;
          .name {
            font-size: 3vh;
            line-clamp: 2;
            -webkit-line-clamp: 2;
            padding-right: 26px;
            span {
              &:nth-of-type(2) {
                margin-left: 12px;
                font-size: 2.3vh;
                opacity: 0.6;
              }
            }
          }
          .artists {
            margin-top: 4px;
            opacity: 0.6;
            font-size: 1.8vh;
            .artist {
              span {
                &:nth-of-type(2) {
                  margin: 0 2px;
                }
              }
            }
          }
        }
        .lrc-all {
          // margin-right: 20%;
          margin-right: 4vw;
          scrollbar-width: none;
          // max-width: 460px;
          // max-width: 52vh;
          width: 90%;
          overflow: auto;
          mask: linear-gradient(
            180deg,
            hsla(0, 0%, 100%, 0) 0,
            hsla(0, 0%, 100%, 0.6) 15%,
            #fff 25%,
            #fff 75%,
            hsla(0, 0%, 100%, 0.6) 85%,
            hsla(0, 0%, 100%, 0)
          );
          -webkit-mask: linear-gradient(
            180deg,
            hsla(0, 0%, 100%, 0) 0,
            hsla(0, 0%, 100%, 0.6) 15%,
            #fff 25%,
            #fff 75%,
            hsla(0, 0%, 100%, 0.6) 85%,
            hsla(0, 0%, 100%, 0)
          );
          &::-webkit-scrollbar {
            display: none;
          }
          &.cover {
            height: 80vh;
          }
          &.record {
            height: 60vh;
          }
          &:hover {
            .lrc-text {
              &.blur {
                filter: blur(0) !important;
              }
            }
          }
          .placeholder {
            width: 100%;
            height: 50%;
          }
          .lrc {
            opacity: 0.4;
            transition: all 0.3s;
            // display: flex;
            // flex-direction: column;
            // margin-bottom: 4px;
            // padding: 12px 20px;
            margin-bottom: 0.8vh;
            padding: 1.8vh 3vh;
            padding-left: 5vh;
            border-radius: 8px;
            transition: all 0.3s;
            transform-origin: left center;
            cursor: pointer;
            position: relative;
            .curr-time {
              position: absolute;
              left: 1vh;
              top: 50%;
              transform: translateY(-50%);
              // font-size: 1.2vh;
              opacity: 0;
              transition: all 0.3s;
              font-variant-numeric: tabular-nums;
            }
            .lrc-text {
              display: flex;
              flex-direction: column;
              transition: all 0.35s ease-in-out;
              transform: scale(0.95);
              transform-origin: left center;
              .lyric {
                transition: all 0.3s;
                // font-size: 2.4vh;
              }
              .lyric-fy {
                margin-top: 2px;
                transition: all 0.3s;
                opacity: 0.8;
                // font-size: 2vh;
              }
            }
            &.on {
              opacity: 1;
              .lrc-text {
                transform: scale(1.3);
                .lyric, .lyric-fy {
                  color: $mainColor !important;
                  font-weight: 900 !important;
                  text-shadow: 0 0 20px var(--main-primary-color-dim, rgba(0, 150, 136, 0.4)) !important;
                  opacity: 1 !important;
                  filter: none !important;
                }
              }
            }
            &:hover {
              @media (min-width: 768px) {
                background-color: #ffffff20;
                .curr-time {
                  opacity: 0.6;
                }
              }
            }
            &:active {
              transform: scale(0.95);
            }
          }
        }
        .menu {
          opacity: 0;
          padding: 0 20px;
          margin-top: 20px;
          display: flex;
          flex-direction: row;
          align-items: center;
          transition: all 0.3s;
          &.show {
            opacity: 1;
          }
          .lyric-offset-control {
            display: flex;
            align-items: center;
            background-color: #ffffff20;
            backdrop-filter: blur(10px);
            padding: 4px 10px;
            border-radius: 20px;
            margin-left: 12px;
            .btn {
              padding: 0 !important;
              margin: 0 !important;
              font-size: 20px !important;
              opacity: 0.8;
              &:hover {
                background-color: transparent !important;
                opacity: 1;
                transform: scale(1.1);
              }
            }
            .text {
              margin: 0 10px;
              font-size: 1.6vh;
              min-width: 40px;
              text-align: center;
            }
          }
          .speed-control {
            position: relative;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 8px;
            
            .n-icon {
               margin-right: 0 !important;
            }

            .speed-popup {
               position: absolute;
               bottom: 40px;
               left: 50%;
               transform: translateX(-50%);
               width: 36px;
               height: 120px;
               background: rgba(0,0,0,0.6);
               backdrop-filter: blur(10px);
               border-radius: 18px;
               padding: 12px 0;
               display: flex;
               flex-direction: column;
               align-items: center;
               opacity: 0;
               visibility: hidden;
               transition: all 0.3s;
               
               .val {
                  font-size: 10px;
                  margin-bottom: 8px;
                  font-weight: bold;
               }
               
               .slider {
                  height: 100%;
                  --n-handle-size: 12px;
                  --n-rail-width: 4px;
               }
            }
            
            &:hover .speed-popup {
               opacity: 1;
               visibility: visible;
            }
          }
          .n-icon {
            margin-right: 8px;
            font-size: 24px;
            cursor: pointer;
            padding: 8px;
            border-radius: 8px;
            opacity: 0.4;
            transition: all 0.3s;
            &:hover {
              background-color: #ffffff30;
            }
            &:active {
              transform: scale(0.95);
            }
            &.open {
              opacity: 1;
            }
          }
        }
      }
    }
  }
  .canvas {
    display: flex;
    justify-content: center;
    align-items: flex-end;
    max-width: 1600px;
    z-index: -1;
    position: absolute;
    bottom: 0;
    -webkit-mask: linear-gradient(
      to right,
      hsla(0deg, 0%, 100%, 0) 0,
      hsla(0deg, 0%, 100%, 0.6) 15%,
      #fff 30%,
      #fff 70%,
      hsla(0deg, 0%, 100%, 0.6) 85%,
      hsla(0deg, 0%, 100%, 0)
    );
    mask: linear-gradient(
      to right,
      hsla(0deg, 0%, 100%, 0) 0,
      hsla(0deg, 0%, 100%, 0.6) 15%,
      #fff 30%,
      #fff 70%,
      hsla(0deg, 0%, 100%, 0.6) 85%,
      hsla(0deg, 0%, 100%, 0)
    );
    .avBars {
      max-width: 1600px;
      opacity: 0.6;
    }
  }
}
</style>
