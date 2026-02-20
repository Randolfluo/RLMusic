<template>
  <div class="artist-view">
    <n-spin :show="loading">
      <!-- 歌手信息头部 -->
      <div class="artist-header" v-if="artistInfo">
        <div class="cover-wrapper">
          <n-image
            class="cover-img"
            :src="artistInfo.cover || '/images/logo/logo.png'"
            fallback-src="/images/logo/logo.png"
            object-fit="cover"
          />
        </div>
        <div class="info-content">
          <h1 class="name">{{ artistInfo.name }}</h1>
          <div class="desc" v-if="artistInfo.description && artistInfo.description !== ''">
            <n-ellipsis :line-clamp="3" :tooltip="false">
              {{ artistInfo.description }}
              <template #tooltip>
                  <div class="desc-tooltip">{{ artistInfo.description }}</div>
              </template>
            </n-ellipsis>
          </div>
           <div class="desc" v-else>
            暂无简介
          </div>
        </div>
      </div>

      <!-- 歌手歌曲列表 -->
      <div class="artist-songs" v-if="artistInfo">
        <n-divider title-placement="left">歌曲列表</n-divider>
        <SongList :songs="artistInfo.songs || []" />
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import { useMessage, NImage, NEllipsis, NSpin, NDivider } from "naive-ui";
import { getArtistDetail } from "@/api/song";
import { ResultCode } from "@/utils/request";
import SongList from "@/components/DataList/SongList.vue";

const route = useRoute();
const message = useMessage();

const loading = ref(false);
const artistInfo = ref<any>(null);

const initData = async () => {
  const id = route.query.id as string;
  if (!id) return;

  loading.value = true;
  try {
    // 1. 获取歌手详情
    const res = await getArtistDetail(id);
    if (res.code === ResultCode.SUCCESS) {
      artistInfo.value = res.data;
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

watch(() => route.query.id, () => {
    if (route.name === 'artist') {
        initData();
    }
});

onMounted(() => {
  initData();
});
</script>

<style scoped lang="scss">
.artist-view {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;

  .artist-header {
    display: flex;
    gap: 32px;
    align-items: flex-start;

    .cover-wrapper {
        flex-shrink: 0;
        width: 240px;
        height: 240px;
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 8px 24px rgba(0,0,0,0.15);

        .cover-img {
            width: 100%;
            height: 100%;
        }
    }

    .info-content {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        min-height: 240px;

        .name {
            font-size: 32px;
            font-weight: 800;
            margin-bottom: 16px;
            color: var(--n-text-color);
        }

        .desc {
            color: var(--n-text-color-3);
            margin-bottom: 24px;
            line-height: 1.6;
            max-width: 600px;
        }
    }
  }

  .artist-songs {
    margin-top: 32px;
  }

  @media (max-width: 768px) {
    padding: 16px;
    
    .artist-header {
      flex-direction: column;
      align-items: center;
      text-align: center;
      gap: 24px;

      .cover-wrapper {
        width: 180px;
        height: 180px;
      }

      .info-content {
        min-height: auto;
        align-items: center;

        .name {
          font-size: 24px;
        }

        .desc {
          margin-bottom: 0;
        }
      }
    }
  }
}
</style>