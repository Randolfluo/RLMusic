<template>
  <div class="playlist-detail">
    <div class="header">
      <div class="cover">
        <n-image
          class="cover-img"
          :src="playlist.cover_url || '/images/logo/favicon.png'"
          fallback-src="/images/logo/favicon.png"
          object-fit="cover"
          preview-disabled
          @click="music.setBigPlayerState(true)"
          style="cursor: pointer;"
        />
      </div>
      <div class="info">
        <div class="tag">歌单</div>
        <div class="title">{{ playlist.title }}</div>
        <div class="creator" v-if="playlist.owner_id">
          Created by User {{ playlist.owner_id }}
        </div>
        <div class="desc" v-if="playlist.description">
          {{ playlist.description }}
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
          :songs="playlist.songs || []" 
          :loading="loading"
          :page="page"
          :page-size="limit"
        />
        <div class="pagination-container" style="display: flex; justify-content: center; margin-top: 20px;">
          <Pagination
            v-if="playlist.songs && playlist.songs.length > 0"
            :totalCount="playlist.total || 0"
            :pageNumber="page"
            :showSizePicker="true"
            @pageNumberChange="onPageChange"
            @pageSizeChange="onPageSizeChange"
          />
        </div>
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { getPublicPlaylistDetail, getPrivatePlaylistDetail } from "@/api/playlist";
import { ResultCode } from "@/utils/request";
import { useMessage, NButton, NIcon, NImage, NDivider, NSpin } from "naive-ui";
import { Play } from "@icon-park/vue-next";
import { musicStore } from "@/store";
import Pagination from "@/components/Pagination/index.vue";
import SongList from "@/components/DataList/SongList.vue";

const route = useRoute();
const message = useMessage();
const music = musicStore();

const loading = ref(false);
const playlist = ref<any>({});
const page = ref(1);
const limit = ref(30);

onMounted(() => {
  const id = route.params.id as string;
  if (id) {
    fetchPlaylistDetail(id);
  } else {
    message.error("未找到歌单ID");
  }
});

const fetchPlaylistDetail = async (id: string) => {
  loading.value = true;
  try {
    // 尝试获取公开歌单详情
    try {
      const res = await getPublicPlaylistDetail(id, page.value, limit.value);
      if (res.code === ResultCode.SUCCESS) {
        playlist.value = res.data;
        return;
      }
    } catch (e) {
      // 获取公开详情失败（可能是私有歌单），继续尝试获取私有详情
    }

    // 尝试获取私有歌单详情
    const res = await getPrivatePlaylistDetail(id, page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      playlist.value = res.data;
    } else {
      // 如果都失败了
      message.error(res.message || "获取歌单详情失败");
    }
  } catch (error) {
    message.error("获取歌单详情失败");
  } finally {
    loading.value = false;
  }
};

const onPageChange = (val: number) => {
  page.value = val;
  fetchPlaylistDetail(route.params.id as string);
};

const onPageSizeChange = (val: number) => {
  limit.value = val;
  page.value = 1;
  fetchPlaylistDetail(route.params.id as string);
};

const playAll = () => {
    if (playlist.value.songs && playlist.value.songs.length > 0) {
        // 重构歌曲数据结构以适配 store??
        // 假设 store 需要标准结构，目前后端返回的字段已经在 Metadata 中包含
        // 我们可能需要适配一下字段名，例如 song.name -> song.title
        // 暂时直接传，视 backend 返回的 json 结构而定
        
        // 适配 backend song -> frontend song
        // Backend: title, artist_name, album_title, id, cover_url
        // Frontend Player usually expects: name, artist (array), album object...
        // 让我们先按后端其实返回了完整的 Song 对象来处理
       
        // 这里做一个简单的映射，防止前端播放器报错，具体视 Player 组件实现而定
        const tracks = playlist.value.songs.map((song: any) => ({
            ...song,
            name: song.title, // 适配 name
            artist: [{ name: song.artist_name, id: song.artist_id }], // 适配 artist array
            album: { name: song.album_title, id: song.album_id, picUrl: song.cover_url || playlist.value.cover_url } // 适配 album
        }));

        music.setPlaylists(tracks);
        music.setPlaySongIndex(0);
        music.setPlayState(true);
    }
}

const rowProps = (_row: any, index: number) => {
  return {
    style: "cursor: pointer;",
    onClick: () => {
        // 单击播放（原双击逻辑改为单击）
        const tracks = playlist.value.songs.map((song: any) => ({
            ...song,
            name: song.title, 
            artist: [{ name: song.artist_name, id: song.artist_id }],
            album: { name: song.album_title, id: song.album_id, picUrl: song.cover_url || playlist.value.cover_url }
        }));
        music.setPlaylists(tracks);
        music.setPlaySongIndex(index);
        music.setPlayState(true);
    }
  };
};

</script>

<style scoped lang="scss">
:deep(.n-data-table .n-data-table-td) {
  background-color: #fff !important;
}
:deep(.n-data-table .n-data-table-tr:hover .n-data-table-td) {
  background-color: #f0f0f0 !important;
}
:deep(.n-data-table .n-data-table-th) {
  background-color: #fff !important;
}
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
      }
      
      .actions {
        margin-top: auto;
      }
    }
  }

  .songs-list {
    margin-top: 20px;
    
    :deep(.song-title-link) {
      color: inherit;
      cursor: pointer;
      text-decoration: none;
      transition: color 0.3s var(--n-bezier);
      
      &:hover {
        color: var(--n-color-primary);
      }
    }
  }
}
</style>