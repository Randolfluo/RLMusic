<template>
  <div class="playlist-detail">
    <div class="header">
      <div class="cover">
        <n-image
          class="cover-img"
          :src="album.picUrl || album.cover_url || '/images/logo/favicon.png'"
          fallback-src="/images/logo/favicon.png"
          object-fit="cover"
          preview-disabled
          @click="music.setBigPlayerState(true)"
          style="cursor: pointer;"
        />
      </div>
      <div class="info">
        <div class="tag">专辑</div>
        <div class="title">{{ album.title || album.name }}</div>
        <div class="creator" v-if="album.artist_name || (album.artist && album.artist.name)">
          歌手: {{ album.artist_name || (album.artist && album.artist.name) }}
        </div>
        <div class="desc" v-if="album.description">
          {{ album.description }}
        </div>
        <div class="actions">
          <n-button type="primary" round size="large" @click="playAll">
            <template #icon>
              <n-icon :component="Play" />
            </template>
            播放全部
          </n-button>
        </div>
      </div>
    </div>

    <n-divider />

    <div class="songs-list">
      <n-spin :show="loading">
        <SongList 
          :songs="album.songs || []" 
          :loading="loading"
        />
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import { getAlbumDetail } from "@/api/song";
import { ResultCode } from "@/utils/request";
import { useMessage, NButton, NIcon, NImage, NDivider, NSpin } from "naive-ui";
import { Play } from "@icon-park/vue-next";
import { musicStore } from "@/store";
import SongList from "@/components/DataList/SongList.vue";

const route = useRoute();
const message = useMessage();
const music = musicStore();

const loading = ref(false);
const album = ref<any>({});

const fetchAlbumDetail = async (id: string) => {
  if (!id) return;
  
  loading.value = true;
  try {
    const res = await getAlbumDetail(id);
    if (res.code === ResultCode.SUCCESS) {
      album.value = res.data;
    } else {
      message.error(res.message || "获取专辑详情失败");
    }
  } catch (error) {
    message.error("获取专辑详情失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  const id = route.query.id as string;
  if (id) {
    fetchAlbumDetail(id);
  } else {
    message.error("未找到专辑ID");
  }
});

// Watch for route changes (e.g. clicking another album from recommendations if we add them later)
watch(
  () => route.query.id,
  (newId) => {
    if (newId) {
      fetchAlbumDetail(newId as string);
    }
  }
);

const playAll = () => {
    if (album.value.songs && album.value.songs.length > 0) {
        // Map songs to player format
        const tracks = album.value.songs.map((song: any) => ({
            ...song,
            name: song.title,
            artist: song.artists || [{ name: song.artist_name, id: song.artist_id }],
            album: { 
              name: album.value.title || album.value.name, 
              id: album.value.id, 
              picUrl: album.value.picUrl || album.value.cover_url 
            },
            picUrl: song.cover_url || album.value.picUrl || album.value.cover_url
        }));

        music.setPlaylists(tracks);
        music.setPlaySongIndex(0);
        music.setPlayState(true);
    }
}
</script>

<style scoped lang="scss">
.playlist-detail {
  padding: 24px;
  
  .header {
    display: flex;
    margin-bottom: 24px;
    
    .cover {
      width: 200px;
      height: 200px;
      border-radius: 8px;
      overflow: hidden;
      margin-right: 24px;
      flex-shrink: 0;
      box-shadow: 0 4px 12px rgba(0,0,0,0.1);
      
      .cover-img {
        width: 100%;
        height: 100%;
      }
    }
    
    .info {
      flex: 1;
      display: flex;
      flex-direction: column;
      
      .tag {
        display: inline-block;
        border: 1px solid var(--n-color-primary);
        color: var(--n-color-primary);
        padding: 2px 8px;
        border-radius: 4px;
        font-size: 13px;
        align-self: flex-start;
        margin-bottom: 12px;
      }
      
      .title {
        font-size: 24px;
        font-weight: bold;
        margin-bottom: 12px;
        color: var(--n-text-color);
      }
      
      .creator, .desc {
        color: var(--n-text-color-3);
        margin-bottom: 8px;
        font-size: 14px;
      }
      
      .desc {
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        margin-top: 8px;
      }
      
      .actions {
        margin-top: auto;
      }
    }
  }

  .songs-list {
    margin-top: 20px;
  }
}
</style>
