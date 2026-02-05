<template>
  <div class="history-view">
    <div class="header">
      <h2>播放历史</h2>
      <div class="actions">
        <n-button size="small" type="error" ghost @click="handleClearHistory" :disabled="loading || historyList.length === 0">
          <template #icon>
            <n-icon :component="Delete" />
          </template>
          清空历史
        </n-button>
        <n-button size="small" type="primary" @click="handlePlayAll" :disabled="loading || historyList.length === 0" style="margin-left: 12px">
           <template #icon>
            <n-icon :component="Play" />
          </template>
          播放全部
        </n-button>
      </div>
    </div>

    <n-spin :show="loading">
      <div v-if="historyList.length === 0 && !loading" class="empty">
        <n-empty description="暂无播放记录" />
      </div>
      <n-data-table
        v-else
        :columns="columns"
        :data="historyList"
        :bordered="false"
        :row-props="rowProps"
        :row-class-name="() => 'song-row'"
      />
      <div class="pagination" v-if="total > 0">
        <n-pagination
          v-model:page="page"
          v-model:page-size="limit"
          :item-count="total"
          show-size-picker
          :page-sizes="[20, 50, 100]"
          @update:page="fetchHistory"
          @update:page-size="onPageSizeChange"
        />
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from "vue";
import { getHistoryList, clearHistory } from "@/api/song";
import { ResultCode } from "@/utils/request";
import { useMessage, NImage, NTime } from "naive-ui";
import { useMusicDataStore } from "@/store/musicData";
import { Delete, Play } from "@icon-park/vue-next";
import { useUserDataStore } from "@/store/userData";

const message = useMessage();
const music = useMusicDataStore();
const userStore = useUserDataStore();

const loading = ref(false);
const historyList = ref<any[]>([]);
const page = ref(1);
const limit = ref(20);
const total = ref(0);

const columns = [
  {
    title: "#",
    key: "index",
    width: 60,
    render: (_: any, index: number) => index + 1 + (page.value - 1) * limit.value,
  },
  {
    title: "封面",
    key: "cover",
    width: 60,
    render: (row: any) => {
      return h(NImage, {
        src: row.cover_url || '/images/logo/favicon.png',
        width: 40,
        height: 40,
        objectFit: 'cover',
        style: { borderRadius: '4px', verticalAlign: 'middle' },
        previewDisabled: true
      });
    }
  },
  {
    title: "标题",
    key: "title",
    ellipsis: { tooltip: true }
  },
  {
    title: "歌手",
    key: "artist_name",
    render: (row: any) => row.artist || "Unknown",
    ellipsis: { tooltip: true }
  },
  {
    title: "专辑",
    key: "album_title",
    render: (row: any) => row.album || "Unknown",
    ellipsis: { tooltip: true }
  },
  {
    title: "播放时间",
    key: "played_at",
    width: 180,
    render: (row: any) => h(NTime, { time: new Date(row.played_at), format: 'yyyy-MM-dd HH:mm:ss' })
  },
];

onMounted(() => {
  if (userStore.userLogin) {
    fetchHistory();
  } else {
    // If not logged in, maybe show local history?
    // For now, let's just stick to backend history as requested.
    // Or we could mix them? The request is "frontend and backend add history",
    // strongly suggesting the backend one.
    // I will show a message or empty state if not logged in, but better to fetch backend.
    // Actually, local logic was user.persistData.playHistory.
    // The requirement implies syncing.
    // Let's assume user is logged in for this view or redirect/show empty.
  }
});

const fetchHistory = async () => {
  loading.value = true;
  try {
    const res: any = await getHistoryList(page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      historyList.value = res.data.list;
      total.value = res.data.total;
    }
  } catch (error) {
    message.error("获取播放历史失败");
  } finally {
    loading.value = false;
  }
};

const onPageSizeChange = () => {
  page.value = 1;
  fetchHistory();
};

const handleClearHistory = async () => {
  try {
    const res: any = await clearHistory();
    if (res.code === ResultCode.SUCCESS) {
      message.success("清空成功");
      historyList.value = [];
      total.value = 0;
    }
  } catch (error) {
    message.error("清空失败");
  }
};

const handlePlayAll = () => {
  if (historyList.value.length > 0) {
      // Need to map structure if Player expects different fields
      // History items from backend: { id, title, artist, album, cover_url, ... }
      // Player typically expects: { name, artist: [{name}], album: {name, picUrl} }
      // The backend returns convenient strings for artist/album names in handle_history.go
      // We need to reconstruct objects for the Player to work best

      const tracks = historyList.value.map(item => ({
          ...item,
          name: item.title,
          artist: [{ name: item.artist, id: 0 }],
          album: { name: item.album, id: 0, picUrl: item.cover_url }
      }));

      music.setPlaylists(tracks);
      music.setPlaySongIndex(0);
      music.setPlayState(true);
  }
};

const rowProps = (row: any, index: number) => {
  return {
    style: "cursor: pointer;",
    onClick: () => {
       const tracks = historyList.value.map(item => ({
          ...item,
          name: item.title,
          artist: [{ name: item.artist, id: 0 }],
          album: { name: item.album, id: 0, picUrl: item.cover_url }
      }));
      music.setPlaylists(tracks);
      music.setPlaySongIndex(index);
      music.setPlayState(true);
    }
  };
};
</script>

<style scoped lang="scss">
.history-view {
  padding: 24px;
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    h2 {
      font-size: 24px;
      font-weight: bold;
      margin: 0;
    }
  }
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }
}

:deep(.n-data-table .n-data-table-td) {
  background-color: transparent;
}
:deep(.song-row:hover .n-data-table-td) {
  background-color: var(--n-color-hover) !important;
}
</style>
