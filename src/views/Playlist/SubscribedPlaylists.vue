<template>
  <div class="subscribed-playlists">
    <div class="playlist-header">
      <div class="cover">
        <div class="cover-placeholder">
          <n-icon size="60" color="#ffffff">
            <Star theme="filled" />
          </n-icon>
        </div>
      </div>
      <div class="info">
        <div class="tag">
          <n-tag type="warning" size="small" round>
            <template #icon>
              <n-icon :component="Star" />
            </template>
            收藏歌单
          </n-tag>
        </div>
        <h1 class="title">收藏的歌单</h1>
        <div class="meta">
          <div class="creator">
            <n-avatar
              round
              size="small"
              :src="user.userData.avatarUrl || '/images/logo/favicon.png'"
              fallback-src="/images/logo/favicon.png"
            />
            <span class="name">{{ user.userData.nickname || '用户' }}</span>
          </div>
        </div>
        <div class="stats">
          共 {{ playlists.length }} 个歌单
        </div>
      </div>
    </div>

    <div class="content">
      <PlaylistGrid 
        :loading="loading" 
        :playlists="playlists" 
        empty-text="暂无收藏歌单" 
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getSubscribedPlaylists } from "@/api/playlist";
import { ResultCode } from "@/utils/request";
import { useMessage, NTag, NAvatar, NIcon } from "naive-ui";
import { Star } from "@icon-park/vue-next";
import PlaylistGrid from "@/components/DataList/PlaylistGrid.vue";
import { useUserDataStore } from "@/store/userData";

const message = useMessage();
const user = useUserDataStore();
const loading = ref(false);
const playlists = ref<any[]>([]);

onMounted(() => {
  getPlaylists();
});

const getPlaylists = async () => {
  loading.value = true;
  try {
    const res = await getSubscribedPlaylists();
    if (res.code === ResultCode.SUCCESS) {
      if (Array.isArray(res.data)) {
          playlists.value = res.data;
      } else if (res.data && Array.isArray(res.data.list)) {
          playlists.value = res.data.list;
      } else {
          playlists.value = [];
      }
    }
  } catch (error) {
    message.error("获取歌单失败");
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped lang="scss">
.subscribed-playlists {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;

  .playlist-header {
    display: flex;
    gap: 32px;
    margin-bottom: 32px;
    
    .cover {
      flex-shrink: 0;
      width: 200px;
      height: 200px;
      border-radius: 12px;
      overflow: hidden;
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
      
      .cover-placeholder {
        width: 100%;
        height: 100%;
        background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .info {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: center;
      
      .tag {
        margin-bottom: 12px;
      }

      .title {
        font-size: 32px;
        font-weight: 800;
        margin: 0 0 16px 0;
        line-height: 1.2;
      }

      .meta {
        display: flex;
        align-items: center;
        gap: 16px;
        margin-bottom: 12px;
        font-size: 13px;
        color: var(--n-text-color-3);

        .creator {
          display: flex;
          align-items: center;
          gap: 8px;
          
          .name {
            color: var(--n-text-color-2);
            font-weight: 500;
          }
        }
      }
      
      .stats {
        font-size: 13px;
        color: var(--n-text-color-3);
      }
    }
  }

  .content {
    background-color: var(--n-card-color);
    border-radius: 12px;
    padding: 24px;
    min-height: 400px;
  }
}

@media (max-width: 768px) {
  .subscribed-playlists {
    padding: 16px;

    .playlist-header {
      flex-direction: column;
      align-items: center;
      text-align: center;
      gap: 20px;

      .cover {
        width: 160px;
        height: 160px;
      }

      .info {
        align-items: center;
        
        .meta {
          justify-content: center;
        }
      }
    }
  }
}
</style>
