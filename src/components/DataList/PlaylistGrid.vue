<template>
  <div class="playlist-grid">
    <n-spin :show="loading">
      <div v-if="!loading && playlists.length === 0" class="empty">
        <n-empty :description="emptyText" />
      </div>
      <n-grid
        v-else
        x-gap="20"
        y-gap="20"
        :cols="cols"
        responsive="screen"
        :collapsed="collapsed"
        :collapsed-rows="collapsedRows"
      >
        <n-grid-item v-for="item in playlists" :key="item.id">
          <n-card
            hoverable
            class="playlist-card"
            content-style="padding: 0;"
            @click="onPlaylistClick(item.id)"
          >
            <div class="cover-container">
              <n-image
                preview-disabled
                class="cover-img"
                object-fit="cover"
                :src="item.cover_url || '/images/logo/favicon.png'"
                fallback-src="/images/logo/favicon.png"
              />
              <div class="play-count">
                <n-icon :component="Play" size="12" />
                <span>{{ formatCount(item.play_count) }}</span>
              </div>
            </div>
            <div class="info">
              <div class="title">{{ item.title }}</div>
            </div>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { Play } from "@icon-park/vue-next";
import { useRouter } from "vue-router";

const router = useRouter();

defineProps({
  playlists: {
    type: Array as () => any[],
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
  emptyText: {
    type: String,
    default: "暂无数据",
  },
  cols: {
    type: String,
    default: "2 s:3 m:4 l:5 xl:6",
  },
  collapsed: {
    type: Boolean,
    default: false,
  },
  collapsedRows: {
    type: Number,
    default: 1,
  },
});

const onPlaylistClick = (id: number) => {
  router.push(`/playlist/${id}`);
};

const formatCount = (count: number) => {
  if (!count) return 0;
  if (count > 10000) return (count / 10000).toFixed(1) + "万";
  return count;
};
</script>

<style scoped lang="scss">
.playlist-grid {
  .playlist-card {
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.3s;
    
    &:hover {
      transform: translateY(-5px);
    }
    
    .cover-container {
      position: relative;
      width: 100%;
      padding-top: 100%;
      background-color: #f5f5f5;
      
      .cover-img {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        
        :deep(img) {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
        }
      }
      
      .play-count {
        position: absolute;
        right: 8px;
        top: 8px;
        background: rgba(0, 0, 0, 0.4);
        backdrop-filter: blur(4px);
        color: #fff;
        padding: 2px 8px;
        border-radius: 12px;
        font-size: 12px;
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
    
    .info {
      padding: 12px 10px;
      .title {
        font-size: 14px;
        line-height: 1.4;
        height: 40px;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        font-weight: 500;
        color: var(--n-text-color);
      }
    }
  }
}
</style>
