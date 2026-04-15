<template>
  <div class="artist-page">
    <!-- Background Decoration -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
    </div>

    <div class="artist-content">
      <n-spin :show="loading">
        <!-- Artist Header -->
        <div class="header-section glass-panel" v-if="artistInfo">
          <div class="cover-wrapper">
            <n-image
              class="cover-img"
              :src="resolveCoverUrl(artistInfo.cover) || '/images/logo/favicon.png'"
              fallback-src="/images/logo/favicon.png"
              object-fit="cover"
              preview-disabled
            />
            <div class="cover-overlay">
              <n-icon :component="Microphone" size="48" />
            </div>
          </div>
          <div class="info-wrapper">
            <div class="tag-badge">Artist</div>
            <h1 class="artist-name">{{ artistInfo.name }}</h1>
            <div class="desc-wrapper" v-if="artistInfo.description && artistInfo.description !== ''">
              <n-ellipsis :line-clamp="3" :tooltip="false">
                {{ artistInfo.description }}
              </n-ellipsis>
            </div>
            <div class="desc-wrapper empty" v-else>
              暂无简介
            </div>
            <div class="meta-info">
              <span class="song-count">{{ total }} 首歌曲</span>
            </div>
          </div>
        </div>

        <!-- Songs Section -->
        <div class="songs-section glass-panel" v-if="artistInfo">
          <div class="section-header">
            <div class="section-icon">
              <n-icon :component="Music" size="20" />
            </div>
            <h2 class="section-title">歌曲列表</h2>
          </div>
          <SongList :songs="songs" />
          <div class="pagination-container" v-if="total > 0">
            <Pagination
              :totalCount="total"
              :pageNumber="page"
              :showSizePicker="true"
              @pageNumberChange="onPageChange"
              @pageSizeChange="onPageSizeChange"
            />
          </div>
        </div>
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import { useMessage, NImage, NEllipsis, NSpin, NIcon } from "naive-ui";
import { getArtistDetail, resolveCoverUrl } from "@/api/song";
import { ResultCode } from "@/utils/request";
import SongList from "@/components/DataList/SongList.vue";
import Pagination from "@/components/Pagination/index.vue";
import { Microphone, Music } from "@icon-park/vue-next";

const route = useRoute();
const message = useMessage();

const loading = ref(false);
const artistInfo = ref<any>(null);
const songs = ref<any[]>([]);
const page = ref(1);
const limit = ref(30);
const total = ref(0);

const initData = async () => {
  const id = route.query.id as string;
  if (!id) return;

  loading.value = true;
  try {
    const res = await getArtistDetail(id, page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      const data = res.data || {};
      if (data.artist) {
        artistInfo.value = data.artist;
        songs.value = Array.isArray(data.list) ? data.list : data.songs || [];
        total.value = data.total ?? 0;
      } else {
        artistInfo.value = data;
        songs.value = data.songs || [];
        total.value = Array.isArray(data.songs) ? data.songs.length : 0;
      }
    } else {
      message.error(res.message || "获取歌手信息失败");
    }
  } catch (error) {
    message.error("加载失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

const onPageChange = (val: number) => {
  page.value = val;
  initData();
};

const onPageSizeChange = (val: number) => {
  limit.value = val;
  page.value = 1;
  initData();
};

watch(() => route.query.id, () => {
  if (route.name === 'artist') {
    page.value = 1;
    initData();
  }
});

onMounted(() => {
  initData();
});
</script>

<style scoped lang="scss">
.artist-page {
  background: #faf8f5;
  position: relative;
  min-height: 100vh;
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;

  .bg-decoration {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: -1;
    pointer-events: none;
    overflow: hidden;

    .blob {
      position: absolute;
      border-radius: 50%;
      filter: blur(80px);
      opacity: 0.35;
      animation: blob-float 20s infinite ease-in-out;
    }

    .blob-1 {
      width: 500px;
      height: 500px;
      background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
      top: -100px;
      right: -100px;
      animation-delay: 0s;
    }

    .blob-2 {
      width: 400px;
      height: 400px;
      background: linear-gradient(135deg, #ec4899 0%, #f472b6 100%);
      bottom: -100px;
      left: -100px;
      animation-delay: -5s;
    }
  }

  @keyframes blob-float {
    0%, 100% { transform: translate(0, 0) scale(1); }
    25% { transform: translate(20px, -30px) scale(1.05); }
    50% { transform: translate(-10px, 20px) scale(0.95); }
    75% { transform: translate(15px, 10px) scale(1.02); }
  }
}

.artist-content {
  position: relative;
  z-index: 1;
  animation: fade-in-up 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Header Section */
.header-section {
  display: flex;
  gap: 40px;
  margin-bottom: 32px;
  padding: 32px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px) saturate(180%);
  border-radius: 28px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06);

  .cover-wrapper {
    flex-shrink: 0;
    width: 240px;
    height: 240px;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
    position: relative;
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);

    &:hover {
      transform: scale(1.03) rotate(2deg);

      .cover-overlay {
        opacity: 1;
      }
    }

    .cover-img {
      width: 100%;
      height: 100%;
      display: block;
      transition: transform 0.6s ease;

      :deep(img) {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .cover-overlay {
      position: absolute;
      inset: 0;
      background: rgba(0, 0, 0, 0.4);
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: all 0.3s ease;
      color: white;
      backdrop-filter: blur(4px);
    }
  }

  .info-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-height: 240px;

    .tag-badge {
      display: inline-block;
      font-size: 12px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 1px;
      padding: 6px 14px;
      border-radius: 100px;
      background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
      color: white;
      width: fit-content;
      margin-bottom: 16px;
      box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
    }

    .artist-name {
      font-family: 'Plus Jakarta Sans', sans-serif;
      font-size: 36px;
      font-weight: 800;
      margin: 0 0 16px 0;
      color: var(--n-text-color);
      line-height: 1.2;
      letter-spacing: -0.02em;
    }

    .desc-wrapper {
      font-size: 15px;
      color: var(--n-text-color-3);
      line-height: 1.7;
      margin: 0 0 20px 0;
      max-width: 600px;

      &.empty {
        opacity: 0.6;
        font-style: italic;
      }
    }

    .meta-info {
      display: flex;
      align-items: center;
      gap: 12px;
      font-size: 14px;
      color: var(--n-text-color-3);

      .song-count {
        font-weight: 600;
        padding: 6px 14px;
        background: rgba(236, 72, 153, 0.1);
        border-radius: 100px;
        color: #ec4899;
      }
    }
  }
}

/* Songs Section */
.songs-section {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.04);

  .section-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;

    .section-icon {
      width: 40px;
      height: 40px;
      border-radius: 12px;
      background: linear-gradient(135deg, #ec4899 0%, #f472b6 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      box-shadow: 0 4px 12px rgba(236, 72, 153, 0.3);
    }

    .section-title {
      font-family: 'Plus Jakarta Sans', sans-serif;
      font-size: 20px;
      font-weight: 700;
      margin: 0;
      color: var(--n-text-color);
    }
  }

  .pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .artist-page {
    padding: 20px;
  }

  .header-section {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 24px;
    gap: 24px;

    .cover-wrapper {
      width: 180px;
      height: 180px;
    }

    .info-wrapper {
      min-height: auto;
      align-items: center;

      .tag-badge {
        align-self: center;
      }

      .artist-name {
        font-size: 24px;
      }

      .desc-wrapper {
        max-width: 100%;
      }
    }
  }
}

/* Dark Mode Support */
:global(.dark) {
  .artist-page {
    .header-section {
      background: rgba(30, 30, 30, 0.7);
      border-color: rgba(255, 255, 255, 0.1);
    }

    .songs-section {
      background: rgba(30, 30, 30, 0.6);
      border-color: rgba(255, 255, 255, 0.08);
    }
  }
}
</style>
