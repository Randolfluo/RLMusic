<template>
  <div class="player-record" :class="{ playing: music.getPlayState }">
    <div class="record-wrapper">
      <div class="disc">
        <div class="cover">
            <img :src="coverUrl" alt="cover" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { musicStore } from "@/store";
import { getSongCover } from "@/api/song";

const music = musicStore();

const coverUrl = computed(() => {
  if (music.getPlaySongData?.id) {
    return getSongCover(music.getPlaySongData.id);
  }
  return "/images/logo/logo.png";
});
</script>

<style lang="scss" scoped>
.player-record {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;

  .record-wrapper {
    position: relative;
    width: 45vh;
    height: 45vh;
    max-width: 75vw;
    max-height: 75vw;
    border-radius: 50%;
    background: #000;
    background-image: radial-gradient(circle, #444 0%, #111 100%);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    animation: rotate 20s linear infinite;
    animation-play-state: paused;

    .disc {
      width: 100%;
      height: 100%;
      border-radius: 50%;
      background: linear-gradient(30deg, transparent 40%, rgba(255, 255, 255, 0.1) 45%, rgba(255, 255, 255, 0.1) 55%, transparent 60%);
      display: flex;
      align-items: center;
      justify-content: center;
      
      .cover {
          width: 70%;
          height: 70%;
          border-radius: 50%;
          overflow: hidden;
          border: 4px solid #111;
          
          img {
              width: 100%;
              height: 100%;
              object-fit: cover;
          }
      }
    }
  }

  &.playing {
    .record-wrapper {
      animation-play-state: running;
    }
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
