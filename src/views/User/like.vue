<template>
  <div class="user-like">
    <div class="playlist-header">
      <div class="cover">
        <div class="cover-placeholder">
          <n-icon size="60" color="#ffffff">
            <Like theme="filled" />
          </n-icon>
        </div>
      </div>
      <div class="info">
        <div class="tag">
          <n-tag type="error" size="small" round>
            <template #icon>
              <n-icon :component="Like" />
            </template>
            我喜欢的
          </n-tag>
        </div>
        <h1 class="title">我喜欢的歌曲</h1>
        <div class="meta">
          <div class="creator">
            <n-avatar
              round
              size="small"
              :src="resolveAvatarUrl(user.userData.avatarUrl) || '/images/logo/favicon.png'"
              fallback-src="/images/logo/favicon.png"
            />
            <span class="name">{{ user.userData.nickname || '用户' }}</span>
          </div>
        </div>
        <div class="stats">
          共 {{ total }} 首歌曲
        </div>
        <div class="actions">
          <n-button type="primary" round size="large" @click="handlePlayAll" :disabled="loading || songs.length === 0">
            <template #icon>
              <n-icon :component="Play" />
            </template>
            播放全部
          </n-button>
        </div>
      </div>
    </div>

    <div class="content">
      <SongList 
        :songs="songs" 
        :loading="loading" 
        :page="page" 
        :page-size="limit"
        empty-description="暂无喜欢的歌曲"
      />
      
      <div class="pagination" v-if="total > 0">
        <n-pagination
          v-model:page="page"
          v-model:page-size="limit"
          :item-count="total"
          show-size-picker
          :page-sizes="[20, 50, 100]"
          @update:page="fetchSongs"
          @update:page-size="onPageSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getLikedSongs, getSongCover, resolveCoverUrl } from "@/api/song";
import { resolveAvatarUrl } from "@/api/user";
import { ResultCode } from "@/utils/request";
import { useMessage, NPagination, NButton, NIcon, NTag, NAvatar } from "naive-ui";
import { Play, Like } from "@icon-park/vue-next";
import SongList from "@/components/DataList/SongList.vue";
import { musicStore, userStore } from "@/store";

const message = useMessage();
const music = musicStore();
const user = userStore();

const loading = ref(false);
const songs = ref<any[]>([]);
const page = ref(1);
const limit = ref(20);
const total = ref(0);

const fetchSongs = async () => {
  loading.value = true;
  try {
    const res: any = await getLikedSongs(page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      songs.value = res.data.list;
      total.value = res.data.total;
    }
  } catch (error) {
    message.error("获取喜欢的歌曲失败");
  } finally {
    loading.value = false;
  }
};

const onPageSizeChange = () => {
  page.value = 1;
  fetchSongs();
};

const handlePlayAll = () => {
  if (songs.value.length > 0) {
      const tracks = songs.value.map(song => ({
        ...song,
        name: song.title,
        artist: song.artists || [{ name: song.artist_name, id: song.artist_id }],
        album: song.album || { 
            name: song.album_name || song.album_title, 
            id: song.album_id, 
            picUrl: song.cover_url ? resolveCoverUrl(song.cover_url) : getSongCover(song.id) 
        },
        picUrl: song.cover_url ? resolveCoverUrl(song.cover_url) : getSongCover(song.id)
      }));

      music.setPlaylists(tracks);
      music.setPlaySongIndex(0);
      music.setPlayState(true);
  }
};

onMounted(() => {
  fetchSongs();
});
</script>

<style scoped lang="scss">
.user-like {
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
        background: linear-gradient(135deg, #ef4444 0%, #ec4899 100%);
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
        margin-bottom: 24px;
      }

      .actions {
        display: flex;
        gap: 16px;
      }
    }
  }

  .content {
    background-color: var(--n-card-color);
    border-radius: 12px;
    
    .pagination {
      display: flex;
      justify-content: center;
      padding: 24px 0;
    }
  }
}

@media (max-width: 768px) {
  .user-like {
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
