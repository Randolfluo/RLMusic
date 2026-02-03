<template>
  <div class="home">
    <div class="section-title">
      <h2>公共歌单</h2>
    </div>

    <n-spin :show="loading">
      <div v-if="!loading && playlists.length === 0" class="empty">
        <n-empty description="暂无公共歌单" />
      </div>
      <n-grid
        v-else
        x-gap="20"
        y-gap="20"
        cols="2 s:3 m:4 l:5 xl:6"
        responsive="screen"
      >
        <n-grid-item v-for="item in playlists" :key="item.id">
          <n-card
            hoverable
            class="playlist-card"
            content-style="padding: 0;"
            @click="router.push(`/playlist/${item.id}`)"
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
import { ref, onMounted } from "vue";
import { getPublicPlaylists } from "@/api/playlist";
import { ResultCode } from "@/utils/request";
import { Play } from "@icon-park/vue-next";
import { useMessage } from "naive-ui";
import { useRouter } from "vue-router";

const router = useRouter();
const message = useMessage();
const loading = ref(false);
const playlists = ref<any[]>([]);

onMounted(() => {
  getPlaylists();
});

const getPlaylists = async () => {
  loading.value = true;
  try {
    const res = await getPublicPlaylists();
    if (res.code === ResultCode.SUCCESS) {
      playlists.value = res.data;
    }
  } catch (error) {
    message.error("获取歌单失败");
  } finally {
    loading.value = false;
  }
};

const formatCount = (count: number) => {
  if (!count) return 0;
  if (count > 10000) return (count / 10000).toFixed(1) + "万";
  return count;
};
</script>

<style scoped lang="scss">
.home {
  padding: 24px;
  .section-title {
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    h2 {
      font-size: 24px;
      font-weight: bold;
      margin: 0;
      position: relative;
      padding-left: 10px;
      &::before {
        content: "";
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
        width: 4px;
        height: 18px;
        background-color: var(--n-color-primary);
        border-radius: 2px;
      }
    }
  }
  
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
