<template>
  <div class="history-page">
    <!-- Background Decoration -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
    </div>

    <div class="history-content">
      <!-- Header Section -->
      <div class="header-section glass-panel">
        <div class="cover-wrapper">
          <div class="cover-placeholder">
            <n-icon :component="History" size="48" />
          </div>
        </div>
        <div class="info-wrapper">
          <div class="tag-badge">History</div>
          <h1 class="page-title">播放历史</h1>
          <div class="meta-info">
            <n-avatar
              round
              size="small"
              :src="resolveAvatarUrl(userStore.userData.avatarUrl) || '/images/logo/favicon.png'"
              fallback-src="/images/logo/favicon.png"
            />
            <span class="user-name">{{ userStore.userData.nickname || '用户' }}</span>
            <span class="divider">•</span>
            <span class="song-count">{{ total }} 首歌曲</span>
          </div>
          <div class="actions">
            <n-button type="primary" round size="large" class="play-btn" @click="handlePlayAll" :disabled="loading || historyList.length === 0">
              <template #icon>
                <n-icon :component="Play" />
              </template>
              播放全部
            </n-button>
            <n-button round size="large" class="action-btn" @click="handleClearHistory" :disabled="loading || historyList.length === 0">
              <template #icon>
                <n-icon :component="Delete" />
              </template>
              清空历史
            </n-button>
          </div>
        </div>
      </div>

      <!-- Content Section -->
      <div class="content-section glass-panel">
        <n-spin :show="loading">
          <div v-if="historyList.length === 0 && !loading" class="empty-state">
            <n-empty description="暂无播放记录">
              <template #icon>
                <n-icon :component="Music" size="64" />
              </template>
            </n-empty>
          </div>
          <n-data-table
            v-else
            :columns="columns"
            :data="historyList"
            :bordered="false"
            :row-props="rowProps"
            :row-class-name="() => 'song-row'"
          />
          <div class="pagination-container" v-if="total > 0">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from "vue";
import { getHistoryList, clearHistory, getSongCover, resolveCoverUrl } from "@/api/song";
import { resolveAvatarUrl } from "@/api/user";
import { ResultCode } from "@/utils/request";
import { useMessage, NImage, NTime, NTag, NAvatar, NIcon } from "naive-ui";
import { useMusicDataStore } from "@/store/musicData";
import { Delete, Play, History, Music } from "@icon-park/vue-next";
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
    width: 70,
    render: (row: any) => {
      return h(NImage, {
        src: row.cover_url ? resolveCoverUrl(row.cover_url) : row.song_id || row.id ? getSongCover(row.song_id || row.id) : '/images/logo/favicon.png',
        width: 44,
        height: 44,
        objectFit: 'cover',
        style: { borderRadius: '8px', verticalAlign: 'middle', boxShadow: '0 4px 12px rgba(0,0,0,0.1)' },
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
        style: { fontWeight: '600', fontSize: '15px' }
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
          album: { name: item.album, id: 0, picUrl: item.cover_url ? resolveCoverUrl(item.cover_url) : item.song_id || item.id ? getSongCover(item.song_id || item.id) : '' }
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
          album: { name: item.album, id: 0, picUrl: item.cover_url ? resolveCoverUrl(item.cover_url) : item.song_id || item.id ? getSongCover(item.song_id || item.id) : '' }
      }));
      music.setPlaylists(tracks);
      music.setPlaySongIndex(index);
      music.setPlayState(true);
    }
  };
};
</script>

<style scoped lang="scss">
 .history-page {
  position: relative;
  min-height: 100vh;
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;
  background: #faf8f5;

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
      background: linear-gradient(135deg, rgba(224, 122, 95, 0.4) 0%, rgba(212, 165, 116, 0.3) 100%);
      top: -100px;
      right: -100px;
      animation-delay: 0s;
    }

    .blob-2 {
      width: 400px;
      height: 400px;
      background: linear-gradient(135deg, rgba(61, 139, 139, 0.35) 0%, rgba(124, 111, 174, 0.3) 100%);
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

.history-content {
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
    width: 200px;
    height: 200px;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);

    .cover-placeholder {
      width: 100%;
      height: 100%;
      background: linear-gradient(135deg, #e07a5f 0%, #d4a574 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
    }
  }

  .info-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-height: 200px;

    .tag-badge {
      display: inline-block;
      font-size: 12px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 1px;
      padding: 6px 14px;
      border-radius: 100px;
      background: linear-gradient(135deg, #e07a5f 0%, #d4a574 100%);
      color: white;
      width: fit-content;
      margin-bottom: 16px;
      box-shadow: 0 4px 12px rgba(224, 122, 95, 0.3);
    }

    .page-title {
      font-family: 'Plus Jakarta Sans', sans-serif;
      font-size: 36px;
      font-weight: 800;
      margin: 0 0 16px 0;
      color: var(--n-text-color);
      line-height: 1.2;
      letter-spacing: -0.02em;
    }

    .meta-info {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 24px;
      font-size: 14px;
      color: var(--n-text-color-3);

      .user-name {
        font-weight: 500;
        color: var(--n-text-color-2);
      }

      .divider {
        opacity: 0.5;
      }

      .song-count {
        font-weight: 600;
        padding: 6px 14px;
        background: rgba(224, 122, 95, 0.1);
        border-radius: 100px;
        color: #e07a5f;
      }
    }

    .actions {
      display: flex;
      gap: 12px;

      .play-btn {
        background: linear-gradient(135deg, #e07a5f 0%, #d4a574 100%);
        border: none;
        padding: 0 28px;
        height: 48px;
        font-size: 16px;
        font-weight: 600;
        box-shadow: 0 8px 20px rgba(224, 122, 95, 0.35);
        transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 12px 28px rgba(224, 122, 95, 0.45);
        }
      }

      .action-btn {
        padding: 0 24px;
        height: 48px;
        font-size: 15px;
        font-weight: 600;
        background: rgba(255, 255, 255, 0.6);
        border: 1px solid rgba(0, 0, 0, 0.08);
        transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);

        &:hover {
          background: rgba(255, 255, 255, 0.9);
          transform: translateY(-2px);
          box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
        }
      }
    }
  }
}

/* Content Section */
.content-section {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.04);
  min-height: 400px;

  .empty-state {
    padding: 80px 0;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }
}

/* Data Table Styles */
:deep(.n-data-table .n-data-table-td) {
  background-color: transparent;
  padding: 14px 8px;
  vertical-align: middle;
}

:deep(.n-data-table .n-data-table-th) {
  background-color: transparent;
  font-weight: 600;
  opacity: 0.6;
  padding: 14px 8px;
}

:deep(.song-row) {
  transition: all 0.2s;
  border-radius: 12px;

  &:hover {
    background-color: rgba(224, 122, 95, 0.08) !important;

    .n-data-table-td {
      background-color: transparent !important;
    }
  }
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .history-page {
    padding: 20px;
  }

  .header-section {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 24px;
    gap: 24px;

    .cover-wrapper {
      width: 160px;
      height: 160px;
    }

    .info-wrapper {
      min-height: auto;
      align-items: center;

      .tag-badge {
        align-self: center;
      }

      .page-title {
        font-size: 24px;
      }

      .meta-info {
        justify-content: center;
      }

      .actions {
        justify-content: center;
      }
    }
  }
}

/* Dark Mode Support */
:global(.dark) {
  .history-page {
    .header-section {
      background: rgba(30, 30, 30, 0.7);
      border-color: rgba(255, 255, 255, 0.1);
    }

    .content-section {
      background: rgba(30, 30, 30, 0.6);
      border-color: rgba(255, 255, 255, 0.08);
    }

    :deep(.song-row:hover) {
      background-color: rgba(224, 122, 95, 0.15) !important;
    }
  }
}
</style>