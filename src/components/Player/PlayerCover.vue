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
    
    <PlayerControl v-if="showControls" />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { musicStore } from "@/store";
import { getSongCover } from "@/api/song";
import PlayerControl from "./PlayerControl.vue";

const props = defineProps({
  showControls: {
    type: Boolean,
    default: true
  }
});

const music = musicStore();

const songData = computed(() => {
    return music.getPlaySongData as any;
});

const coverUrl = computed(() => {
  if (songData.value?.id) {
    return getSongCover(songData.value.id);
  }
  return "/images/logo/logo.png";
});
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
}
</style>
