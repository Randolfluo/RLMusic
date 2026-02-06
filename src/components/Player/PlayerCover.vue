<template>
  <div class="player-cover">
    <div class="cover-wrapper">
      <n-image
        :src="coverUrl"
        class="cover-img"
        object-fit="contain"
        fallback-src="/images/logo/logo.png"
        preview-disabled
      />
    </div>
    
    <div class="info-section">
      <div class="song-info">
        <div class="title" :title="songData?.name">
          {{ songData?.name || "暂无歌曲" }}
        </div>
        <div class="artist" :title="songData?.artist_name">
          <span
            v-for="(artist, index) in songData?.artist || []"
            :key="index"
          >
            {{ artist.name }}
            <span v-if="index < (Number(songData?.artist?.length) || 0) - 1"> / </span>
          </span>
          <span v-if="!songData?.artist && songData?.artist_name">
            {{ songData.artist_name }}
          </span>
        </div>
      </div>

      <div class="progress-section">
        <div class="time-track">
          <span class="time">{{ music.getPlaySongTime.songTimePlayed }}</span>
          <n-slider
            v-model:value="music.getPlaySongTime.barMoveDistance"
            :step="0.01"
            :tooltip="false"
            class="progress-slider"
            @update:value="handleProgressUpdate"
          />
          <span class="time">{{ music.getPlaySongTime.songTimeDuration }}</span>
        </div>
      </div>

      <div class="controls-section">
        <n-icon 
            size="24" 
            class="mode-icon btn" 
            :component="modeIcon" 
            @click="music.setPlaySongMode()"
        />
        
        <div class="main-controls">
          <n-icon 
            size="32" 
            class="btn" 
            :component="SkipPreviousRound" 
            @click="music.setPlaySongIndex('prev')"
          />
          
          <div class="play-pause-btn" @click="music.setPlayState(!music.getPlayState)">
            <n-icon 
                size="48" 
                color="white"
                :component="music.getPlayState ? PauseCircleFilled : PlayCircleFilled" 
            />
          </div>
          
          <n-icon 
            size="32" 
            class="btn" 
            :component="SkipNextRound" 
            @click="music.setPlaySongIndex('next')"
          />
        </div>

        <n-icon 
            size="24" 
            class="menu-icon btn" 
            :component="PlaylistPlayRound"
            @click.stop="music.showPlayList = !music.showPlayList"
        />
      </div>
      
      <div class="volume-section">
         <n-icon
          size="20"
          class="volume-icon"
          :component="volumeIcon"
          @click="toggleMute"
        />
        <n-slider
          v-model:value="music.persistData.playVolume"
          :step="0.01"
          :min="0"
          :max="1"
          :tooltip="false"
          class="volume-slider"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { musicStore } from "@/store";
import { getSongCover } from "@/api/song";
import {
  PlayCircleFilled,
  PauseCircleFilled,
  SkipPreviousRound,
  SkipNextRound,
  PlaylistPlayRound,
  VolumeUpRound,
  VolumeDownRound,
  VolumeMuteRound,
  VolumeOffRound,
} from "@vicons/material";
import { PlayCycle, PlayOnce, ShuffleOne } from "@icon-park/vue-next";

const music = musicStore();

const songData = computed(() => {
    return music.getPlaySongData as any;
});

const handleProgressUpdate = (val: number) => {
    if ((window as any).$player && music.getPlaySongTime.duration) {
        (window as any).$player.currentTime = (music.getPlaySongTime.duration / 100) * val;
    }
};

const coverUrl = computed(() => {
  if (songData.value?.id) {
    return getSongCover(songData.value.id);
  }
  return "/images/logo/logo.png";
});

const modeIcon = computed(() => {
    const mode = music.persistData.playSongMode;
    if (mode === 'random') return ShuffleOne;
    if (mode === 'single') return PlayOnce;
    return PlayCycle;
});

const volumeIcon = computed(() => {
    const vol = Number(music.persistData.playVolume);
    if (vol === 0) return VolumeOffRound;
    if (vol < 0.4) return VolumeMuteRound;
    if (vol < 0.7) return VolumeDownRound;
    return VolumeUpRound;
});

const toggleMute = () => {
    if (Number(music.persistData.playVolume) > 0) {
        music.persistData.playVolumeMute = music.persistData.playVolume;
        music.persistData.playVolume = 0;
    } else {
        music.persistData.playVolume = music.persistData.playVolumeMute || 0.5;
    }
};

</script>

<style lang="scss" scoped>
.player-cover {
  display: flex;
  flex-direction: column;
  align-items: center; // Center horizontally
  justify-content: center;
  width: 100%;
  height: 100%;
  padding-bottom: 2rem;
  
  .cover-wrapper {
    width: 45vh;
    height: 45vh;
    max-width: 80vw;
    max-height: 80vw;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    margin-bottom: 24px;

    .cover-img {
      width: 100%;
      height: 100%;
      display: flex; /* Naive UI n-image might render wrapper */
      justify-content: center;
      align-items: center;
      :deep(img) {
        width: 100%;
        height: 100%;
        object-fit: cover; 
      }
    }
  }

  .info-section {
    width: 45vh;
    max-width: 80vw;
    display: flex;
    flex-direction: column;
    gap: 12px;

    .song-info {
        text-align: left;
        margin-bottom: 8px;

        .title {
            font-size: 24px;
            font-weight: bold;
            color: #fff;
            margin-bottom: 8px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }

        .artist {
            font-size: 16px;
            color: rgba(255, 255, 255, 0.7);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }
    }

    .progress-section {
        .time-track {
            display: flex;
            align-items: center;
            gap: 12px;
            
            .time {
                font-size: 12px;
                color: rgba(255, 255, 255, 0.6);
                width: 40px;
                text-align: center;
                font-variant-numeric: tabular-nums;
            }

            .progress-slider {
                flex: 1;
                :deep(.n-slider-rail) {
                    background-color: rgba(255, 255, 255, 0.2);
                    height: 4px;
                }
                :deep(.n-slider-rail__fill) {
                    background-color: #fff;
                }
                :deep(.n-slider-handle) {
                    box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.5);
                }
            }
        }
    }

    .controls-section {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 8px;

        .btn {
            cursor: pointer;
            color: rgba(255, 255, 255, 0.8);
            transition: all 0.2s;
            &:hover {
                color: #fff;
                transform: scale(1.1);
            }
            &:active {
                transform: scale(0.95);
            }
        }

        .main-controls {
            display: flex;
            align-items: center;
            gap: 24px;

            .play-pause-btn {
                cursor: pointer;
                transition: transform 0.2s;
                opacity: 0.9;
                &:hover {
                    transform: scale(1.1);
                    opacity: 1;
                }
                 &:active {
                    transform: scale(0.95);
                }
            }
        }
    }
    
    .volume-section {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-top: 12px;
        
        .volume-icon {
            opacity: 0.7;
            cursor: pointer;
            &:hover { opacity: 1; }
        }
        
        .volume-slider {
            width: 100px;
             :deep(.n-slider-rail) {
                    background-color: rgba(255, 255, 255, 0.2);
                    height: 4px;
            }
            :deep(.n-slider-rail__fill) {
                background-color: rgba(255, 255, 255, 0.8);
            }
        }
    }
  }
}
</style>
