<template>
  <div class="history-view">
    <div class="playlist-header">
      <div class="cover">
        <div class="cover-placeholder">
          <n-icon size="60" color="#ffffff">
            <History theme="filled" />
          </n-icon>
        </div>
      </div>
      <div class="info">
        <div class="tag">
          <n-tag type="info" size="small" round>
            <template #icon>
              <n-icon :component="History" />
            </template>
            最近播放
          </n-tag>
        </div>
        <h1 class="title">播放历史</h1>
        <div class="meta">
          <div class="creator">
            <n-avatar
              round
              size="small"
              :src="userStore.userData.avatarUrl || '/images/logo/favicon.png'"
              fallback-src="/images/logo/favicon.png"
            />
            <span class="name">{{ userStore.userData.nickname || '用户' }}</span>
          </div>
        </div>
        <div class="stats">
          共 {{ total }} 首歌曲
        </div>
        <div class="actions">
          <n-button type="primary" round size="large" @click="handlePlayAll" :disabled="loading || historyList.length === 0">
            <template #icon>
              <n-icon :component="Play" />
            </template>
            播放全部
          </n-button>
          <n-button round size="large" @click="handleClearHistory" :disabled="loading || historyList.length === 0" style="margin-left: 12px">
            <template #icon>
              <n-icon :component="Delete" />
            </template>
            清空历史
          </n-button>
        </div>
      </div>
    </div>

    <div class="content">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from "vue";
import { getHistoryList, clearHistory } from "@/api/song";
import { ResultCode } from "@/utils/request";
import { useMessage, NImage, NTime, NTag, NAvatar, NIcon } from "naive-ui";
import { useMusicDataStore } from "@/store/musicData";
import { Delete, Play, History } from "@icon-park/vue-next";
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
    align: 'center',
    render: (_: any, index: number) => index + 1 + (page.value - 1) * limit.value,
  },
  {
    title: "",
    key: "cover",
    width: 60,
    render: (row: any) => {
      return h(NImage, {
        src: row.cover_url || '/images/logo/favicon.png',
        width: 40,
        height: 40,
        objectFit: 'cover',
        style: { borderRadius: '6px', verticalAlign: 'middle', boxShadow: '0 2px 6px rgba(0,0,0,0.1)' },
        previewDisabled: true
      });
    }
  },
  {
    title: "标题",
    key: "title",
    ellipsis: { tooltip: true },
    render: (row: any) => {
      return h('span', { 
        style: { fontWeight: '500', fontSize: '15px' } 
      }, row.title)
    }
  },
  {
    title: "歌手",
    key: "artist_name",
    render: (row: any) => h('span', { style: { opacity: 0.8 } }, row.artist || "Unknown"),
    ellipsis: { tooltip: true }
  },
  {
    title: "专辑",
    key: "album_title",
    render: (row: any) => h('span', { style: { opacity: 0.6 } }, row.album || "Unknown"),
    ellipsis: { tooltip: true }
  },
  {
    title: "播放时间",
    key: "played_at",
    width: 180,
    render: (row: any) => h(NTime, { 
      time: new Date(row.played_at), 
      format: 'yyyy-MM-dd HH:mm:ss',
      style: { opacity: 0.6, fontSize: '13px' }
    })
  },
];

onMounted(() => {
  if (userStore.userLogin) {
    fetchHistory();
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

const rowProps = (_row: any, index: number) => {
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
        background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
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
    padding: 24px;
    min-height: 400px;

    .pagination {
      display: flex;
      justify-content: center;
      margin-top: 24px;
    }
  }
}

@media (max-width: 768px) {
  .history-view {
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

:deep(.n-data-table .n-data-table-td) {
  background-color: transparent;
  padding: 12px 8px;
  vertical-align: middle;
}

:deep(.n-data-table .n-data-table-th) {
  background-color: transparent;
  font-weight: 600;
  opacity: 0.6;
}

:deep(.song-row) {
  transition: all 0.2s;
  border-radius: 8px;
  
  &:hover {
    background-color: var(--n-color-hover) !important;
    
    .n-data-table-td {
      background-color: transparent !important;
    }
  }
}
</style>
