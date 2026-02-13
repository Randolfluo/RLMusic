<template>
  <div class="user-like">
    <div class="header">
      <div class="title-section">
        <h2>我喜欢的歌曲</h2>
        <div class="actions">
          <n-button size="small" type="primary" @click="handlePlayAll" :disabled="loading || songs.length === 0">
             <template #icon>
              <n-icon :component="Play" />
            </template>
            播放全部
          </n-button>
        </div>
      </div>
    </div>

    <SongList 
      :songs="songs" 
      :loading="loading" 
      :page="page" 
      :page-size="limit"
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
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getLikedSongs } from "@/api/song";
import { ResultCode } from "@/utils/request";
import { useMessage, NPagination, NButton, NIcon } from "naive-ui";
import { Play } from "@icon-park/vue-next";
import SongList from "@/components/DataList/SongList.vue";
import { musicStore } from "@/store";

const message = useMessage();
const music = musicStore();

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
            picUrl: song.cover_url || `/api/song/cover/${song.id}` 
        },
        picUrl: song.cover_url || `/api/song/cover/${song.id}`
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
  .header {
    margin-bottom: 24px;
    
    .title-section {
      display: flex;
      align-items: center;
      gap: 16px;
      
      h2 {
        font-size: 24px;
        font-weight: bold;
        margin: 0;
      }
    }
  }
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }
}
</style>
