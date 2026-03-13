<template>
  <div class="playlist-manage-container">
    <!-- Ambient Background Effects -->
    <div class="ambient-glows">
      <div class="glow-orb orb-1"></div>
      <div class="glow-orb orb-2"></div>
    </div>
    <div class="noise-overlay"></div>

    <div class="content-wrapper">
      <div class="header-section">
        <div class="title-group">
          <button class="nav-btn" @click="router.back()">
            <n-icon :component="Left" />
          </button>
          <div class="text-content">
            <h1 class="page-title">公共歌单库</h1>
            <p class="page-subtitle">管理全站公开的音乐收藏集</p>
          </div>
        </div>
        <div class="header-actions">
          <div class="stat-badge">
            <span class="label">Total</span>
            <span class="value">{{ pagination.itemCount }}</span>
          </div>
        </div>
      </div>

      <div class="table-container glass-panel">
        <n-data-table
          :columns="columns"
          :data="playlistList"
          :loading="loading"
          :pagination="pagination"
          :row-key="row => row.id"
          :row-class-name="'playlist-row'"
          :scroll-x="800"
          remote
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, NPopconfirm, NImage, NIcon, NDataTable, type DataTableColumns } from 'naive-ui';
import { Left, Delete, Music, PlayOne, User } from '@icon-park/vue-next';
import { getPublicPlaylists, deletePublicPlaylist } from '@/api/playlist';
import { resolveCoverUrl } from "@/api/song";
// Force reload dependency
import { ResultCode } from "@/utils/request";

const router = useRouter();
const message = useMessage();

const loading = ref(false);
const playlistList = ref<any[]>([]);

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: ({ itemCount }: { itemCount?: number }) => `共 ${itemCount || 0} 个歌单`,
  onChange: (page: number) => {
    pagination.page = page;
    fetchPlaylists();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    fetchPlaylists();
  }
});

const columns: DataTableColumns = [
  {
    title: '封面',
    key: 'cover',
    width: 80,
    render(row: any) {
      const src = resolveCoverUrl(row.cover_url);
      return h('div', { class: 'cover-wrapper' }, [
        h(NImage, {
            width: 48,
            height: 48,
            src: src,
            fallbackSrc: "/images/pic/default.png",
            objectFit: 'cover',
          showToolbar: false,
          intersectionObserverOptions: { rootMargin: '100px' }
        })
      ]);
    }
  },
  {
    title: '歌单信息',
    key: 'title',
    render(row: any) {
      return h('div', { class: 'info-cell' }, [
        h('span', { class: 'playlist-title' }, row.title),
        h('span', { class: 'playlist-id' }, `ID: ${row.id}`)
      ]);
    }
  },
  {
    title: '统计',
    key: 'stats',
    width: 200,
    render(row: any) {
      return h('div', { class: 'stats-cell' }, [
        h('div', { class: 'stat-item' }, [
          h(NIcon, { component: Music, size: 14 }),
          h('span', `${row.total_songs || 0}`)
        ]),
        h('div', { class: 'stat-item' }, [
          h(NIcon, { component: PlayOne, size: 14 }),
          h('span', (row.play_count || 0).toLocaleString())
        ])
      ]);
    }
  },
  {
    title: '创建者',
    key: 'owner_id',
    width: 150,
    render(row: any) {
      return h('div', { class: 'user-cell' }, [
        h(NIcon, { component: User, size: 14 }),
        h('span', `UID ${row.owner_id}`)
      ]);
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    align: 'center',
    fixed: 'right',
    render(row: any) {
      return h(
        NPopconfirm,
        {
          onPositiveClick: () => handleDelete(row),
          negativeText: '取消',
          positiveText: '确认删除'
        },
        {
          trigger: () => h(
            'button',
            { class: 'action-btn delete', title: '删除歌单' },
            h(NIcon, { component: Delete, size: 18 })
          ),
          default: () => '确定要删除该公共歌单吗？此操作不可逆！'
        }
      );
    }
  }
];

onMounted(() => {
  fetchPlaylists();
});

const fetchPlaylists = async () => {
  loading.value = true;
  try {
    const res = await getPublicPlaylists(pagination.page, pagination.pageSize);
    if (res.code === ResultCode.SUCCESS) {
      playlistList.value = res.data.list;
      pagination.itemCount = res.data.total;
    } else {
      message.error(res.message || '获取歌单列表失败');
    }
  } catch (error) {
    message.error('获取歌单列表失败');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

const handleDelete = async (row: any) => {
  try {
    const res = await deletePublicPlaylist(row.id);
    if (res.code === ResultCode.SUCCESS) {
      message.success('删除歌单成功');
      if (playlistList.value.length === 1 && pagination.page > 1) {
        pagination.page--;
      }
      fetchPlaylists();
    } else {
      message.error(res.message || '删除歌单失败');
    }
  } catch (error) {
    message.error('删除歌单失败');
  }
};
</script>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

.playlist-manage-container {
  position: relative;
  min-height: 100vh;
  width: 100%;
  padding: 40px;
  background: #f8f9fc;
  font-family: 'Plus Jakarta Sans', sans-serif;
  overflow: hidden;

  :global(.dark) & {
    background: #0f1115;
    color: #fff;
  }
}

/* Ambient Background */
.ambient-glows {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;

  .glow-orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(100px);
    opacity: 0.4;
  }

  .orb-1 {
    top: -10%;
    right: -5%;
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(139, 92, 246, 0.3) 0%, rgba(139, 92, 246, 0) 70%);
  }

  .orb-2 {
    bottom: -10%;
    left: -10%;
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, rgba(59, 130, 246, 0.25) 0%, rgba(59, 130, 246, 0) 70%);
  }
}

.noise-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/200' opacity='0.03'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.8' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
  z-index: 1;
  pointer-events: none;
}

.content-wrapper {
  position: relative;
  z-index: 2;
  max-width: 1200px;
  margin: 0 auto;
}

/* Header */
.header-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;

  .title-group {
    display: flex;
    align-items: center;
    gap: 20px;

    .nav-btn {
      width: 48px;
      height: 48px;
      border-radius: 14px;
      border: 1px solid rgba(0, 0, 0, 0.05);
      background: rgba(255, 255, 255, 0.6);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24px;
      color: #64748b;
      cursor: pointer;
      transition: all 0.2s ease;
      backdrop-filter: blur(8px);

      &:hover {
        background: #fff;
        transform: translateX(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
        color: #1e293b;
      }
    }

    .text-content {
      .page-title {
        font-size: 32px;
        font-weight: 800;
        color: #1e293b;
        margin: 0;
        line-height: 1.1;
        letter-spacing: -0.02em;
      }

      .page-subtitle {
        font-size: 14px;
        color: #64748b;
        margin: 6px 0 0;
        font-weight: 500;
      }
    }
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 16px;

    .stat-badge {
      padding: 8px 16px;
      background: rgba(255, 255, 255, 0.5);
      border: 1px solid rgba(255, 255, 255, 0.6);
      border-radius: 100px;
      display: flex;
      align-items: center;
      gap: 8px;
      backdrop-filter: blur(4px);

      .label {
        font-size: 12px;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: #94a3b8;
        font-weight: 700;
      }

      .value {
        font-size: 16px;
        font-weight: 700;
        color: #3b82f6;
      }
    }
  }
}

/* Glass Panel Table */
.glass-panel {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.6);
  border-radius: 24px;
  padding: 8px; /* Padding for the table container */
  box-shadow: 0 20px 40px -10px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

/* Custom Table Styles */
:deep(.n-data-table) {
  --n-th-font-weight: 700 !important;
  --n-th-text-color: #64748b !important;
  
  .n-data-table-th {
    background: transparent !important;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
    padding: 16px 24px !important;
    font-size: 13px;
    text-transform: uppercase;
    letter-spacing: 0.03em;
  }

  .n-data-table-td {
    background: transparent !important;
    border-bottom: 1px solid rgba(0, 0, 0, 0.03) !important;
    padding: 16px 24px !important;
    transition: background 0.2s;
  }

  .n-data-table-tr:last-child .n-data-table-td {
    border-bottom: none !important;
  }

  .n-data-table-tr:hover .n-data-table-td {
    background: rgba(255, 255, 255, 0.5) !important;
  }
}

/* Cell Renderers */
.cover-wrapper {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s;
  
  &:hover {
    transform: scale(1.05);
  }
}

.info-cell {
  display: flex;
  flex-direction: column;
  
  .playlist-title {
    font-size: 15px;
    font-weight: 700;
    color: #1e293b;
    margin-bottom: 4px;
  }
  
  .playlist-id {
    font-size: 12px;
    font-family: monospace;
    color: #94a3b8;
    background: rgba(0,0,0,0.03);
    padding: 2px 6px;
    border-radius: 4px;
    width: fit-content;
  }
}

.stats-cell {
  display: flex;
  gap: 16px;
  
  .stat-item {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #64748b;
    font-size: 13px;
    font-weight: 500;
    
    .n-icon {
      color: #94a3b8;
    }
  }
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #475569;
  font-weight: 600;
  font-size: 13px;
  background: rgba(255, 255, 255, 0.5);
  padding: 6px 12px;
  border-radius: 8px;
  border: 1px solid rgba(0,0,0,0.03);
  width: fit-content;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: #94a3b8;

  &.delete:hover {
    background: #fee2e2;
    color: #ef4444;
    transform: scale(1.1);
  }
}

@media (max-width: 768px) {
  .playlist-manage-container {
    padding: 20px 16px;
  }

  .header-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    
    .header-actions {
      width: 100%;
      justify-content: space-between;
    }
  }
}
</style>
