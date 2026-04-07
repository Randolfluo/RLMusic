<template>
  <div class="playlist-manage-container">
    <!-- 动态背景层 -->
    <div class="ambient-background">
      <div class="gradient-blob blob-1"></div>
      <div class="gradient-blob blob-2"></div>
      <div class="gradient-blob blob-3"></div>
      <div class="grain-overlay"></div>
    </div>

    <div class="content-wrapper">
      <!-- 杂志风格标题区 -->
      <header class="editorial-header">
        <div class="header-left">
          <button class="back-btn" @click="router.back()">
            <n-icon :component="Left" />
          </button>
          <div class="title-stack">
            <span class="kicker">Music Collection</span>
            <h1 class="headline">公共歌单库</h1>
            <p class="deck">管理全站公开的音乐收藏集</p>
          </div>
        </div>
        <div class="header-right">
          <div class="stat-pill">
            <span class="stat-number">{{ pagination.itemCount }}</span>
            <span class="stat-label">个歌单</span>
          </div>
        </div>
      </header>

      <!-- 数据概览卡片 -->
      <section class="metrics-bar">
        <div class="metric-item">
          <span class="metric-value">{{ pagination.page }}</span>
          <span class="metric-label">当前页</span>
        </div>
        <div class="metric-divider"></div>
        <div class="metric-item">
          <span class="metric-value">{{ pagination.pageSize }}</span>
          <span class="metric-label">每页显示</span>
        </div>
        <div class="metric-divider"></div>
        <div class="metric-item">
          <span class="metric-value">{{ playlistList.length }}</span>
          <span class="metric-label">本页歌单</span>
        </div>
      </section>

      <!-- 桌面端表格视图 -->
      <div v-if="!isMobile" class="table-view">
        <div class="table-container">
          <n-data-table
            :columns="columns"
            :data="playlistList"
            :loading="loading"
            :row-key="row => row.id"
            :row-class-name="'playlist-row'"
            :scroll-x="900"
          />
        </div>
      </div>

      <!-- 移动端卡片视图 -->
      <div v-else class="grid-view">
        <div class="playlist-cards">
          <div
            v-for="playlist in playlistList"
            :key="playlist.id"
            class="playlist-card"
          >
            <div class="card-accent"></div>
            <div class="card-content">
              <div class="playlist-header">
                <div class="cover-wrapper">
                  <img
                    :src="resolveCoverUrl(playlist.cover_url) || '/images/pic/default.png'"
                    alt="cover"
                    class="cover-img"
                  />
                </div>
                <div class="playlist-meta">
                  <h3 class="playlist-title">{{ playlist.title }}</h3>
                  <span class="playlist-id">#{{ playlist.id }}</span>
                </div>
              </div>

              <div class="playlist-stats">
                <div class="stat-badge">
                  <n-icon :component="Music" />
                  <span>{{ playlist.total_songs || 0 }} 首歌曲</span>
                </div>
                <div class="stat-badge plays">
                  <n-icon :component="PlayOne" />
                  <span>{{ (playlist.play_count || 0).toLocaleString() }} 播放</span>
                </div>
              </div>

              <div class="playlist-owner">
                <n-icon :component="User" />
                <span>创建者 UID {{ playlist.owner_id }}</span>
              </div>

              <button class="delete-btn" @click="confirmDelete(playlist)">
                <n-icon :component="Delete" />
                删除歌单
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页器 -->
      <div class="pagination-wrapper">
        <n-pagination
          v-model:page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :item-count="pagination.itemCount"
          :show-size-picker="true"
          :page-sizes="pagination.pageSizes"
          @update:page="fetchPlaylists"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, NPopconfirm, NImage, NIcon, NDataTable, NPagination, type DataTableColumns } from 'naive-ui';
import { Left, Delete, Music, PlayOne, User } from '@icon-park/vue-next';
import { getPublicPlaylists, deletePublicPlaylist } from '@/api/playlist';
import { resolveCoverUrl } from "@/api/song";
import { ResultCode } from "@/utils/request";

const router = useRouter();
const message = useMessage();

// 响应式窗口大小检测
const windowWidth = ref(window.innerWidth);
const isMobile = computed(() => windowWidth.value < 768);

const handleResize = () => {
  windowWidth.value = window.innerWidth;
};

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

const loading = ref(false);
const playlistList = ref<any[]>([]);

const pagination = reactive({
  page: 1,
  pageSize: 12,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [12, 24, 48],
});

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  pagination.page = 1;
  fetchPlaylists();
};

const columns: DataTableColumns = [
  {
    title: '封面',
    key: 'cover',
    width: 90,
    render(row: any) {
      const src = resolveCoverUrl(row.cover_url);
      return h('div', { class: 'table-cover-wrapper' }, [
        h(NImage, {
          width: 56,
          height: 56,
          src: src,
          fallbackSrc: "/images/pic/default.png",
          objectFit: 'cover',
          showToolbar: false,
          style: { borderRadius: '12px' }
        })
      ]);
    }
  },
  {
    title: '歌单信息',
    key: 'title',
    render(row: any) {
      return h('div', { class: 'table-info-cell' }, [
        h('span', { class: 'table-playlist-title' }, row.title),
        h('span', { class: 'table-playlist-id' }, `#${row.id}`)
      ]);
    }
  },
  {
    title: '统计',
    key: 'stats',
    width: 220,
    render(row: any) {
      return h('div', { class: 'table-stats' }, [
        h('div', { class: 'table-stat-item' }, [
          h(NIcon, { component: Music, size: 14 }),
          h('span', `${row.total_songs || 0} 首`)
        ]),
        h('div', { class: 'table-stat-item plays' }, [
          h(NIcon, { component: PlayOne, size: 14 }),
          h('span', (row.play_count || 0).toLocaleString())
        ])
      ]);
    }
  },
  {
    title: '创建者',
    key: 'owner_id',
    width: 160,
    render(row: any) {
      return h('div', { class: 'table-owner' }, [
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
            { class: 'table-delete-btn', title: '删除歌单' },
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

const confirmDelete = (row: any) => {
  if (window.confirm('确定要删除该公共歌单吗？此操作不可逆！')) {
    handleDelete(row);
  }
};
</script>

<style lang="scss">
// 导入特色字体
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@600;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

:root {
  // 编辑风格配色 - 温暖米色调
  --bg-primary: #faf8f5;
  --bg-secondary: #f5f2ed;
  --bg-tertiary: #ebe7e0;
  --text-primary: #1a1a1a;
  --text-secondary: #666666;
  --text-muted: #999999;
  --accent-coral: #e07a5f;
  --accent-teal: #3d8b8b;
  --accent-gold: #d4a574;
  --accent-ink: #2c3e50;
  --music-purple: #7c6fae;
  --play-orange: #e07a5f;
  --danger-red: #c75b5b;
  --shadow-soft: 0 4px 20px rgba(0, 0, 0, 0.06);
  --shadow-medium: 0 8px 30px rgba(0, 0, 0, 0.1);
  --shadow-deep: 0 12px 40px rgba(0, 0, 0, 0.14);
}

.playlist-manage-container {
  min-height: 100vh;
  background: var(--bg-primary);
  font-family: 'Plus Jakarta Sans', sans-serif;
  position: relative;
  overflow-x: hidden;
}

// 动态背景
.ambient-background {
  position: fixed;
  inset: 0;
  z-index: 0;
  overflow: hidden;
  pointer-events: none;
}

.gradient-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: blob-float 20s ease-in-out infinite;

  &.blob-1 {
    width: 600px;
    height: 600px;
    background: linear-gradient(135deg, rgba(124, 111, 174, 0.25), rgba(212, 165, 116, 0.2));
    top: -200px;
    right: -100px;
    animation-delay: 0s;
  }

  &.blob-2 {
    width: 500px;
    height: 500px;
    background: linear-gradient(135deg, rgba(224, 122, 95, 0.2), rgba(61, 139, 139, 0.15));
    bottom: -150px;
    left: -150px;
    animation-delay: -7s;
  }

  &.blob-3 {
    width: 400px;
    height: 400px;
    background: linear-gradient(135deg, rgba(61, 139, 139, 0.15), rgba(124, 111, 174, 0.1));
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    animation-delay: -14s;
  }
}

@keyframes blob-float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.05); }
  66% { transform: translate(-20px, 20px) scale(0.95); }
}

.grain-overlay {
  position: absolute;
  inset: 0;
  opacity: 0.03;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}

// 内容包装器
.content-wrapper {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: 48px 32px;

  @media (max-width: 768px) {
    padding: 24px 16px;
  }
}

// 编辑风格标题区
.editorial-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 40px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--bg-tertiary);

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
    margin-bottom: 28px;
  }
}

.header-left {
  display: flex;
  align-items: flex-start;
  gap: 20px;
}

.back-btn {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  border: 1px solid var(--bg-tertiary);
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  color: var(--text-secondary);

  &:hover {
    background: var(--bg-tertiary);
    transform: translateX(-2px);
    color: var(--text-primary);
  }
}

.title-stack {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.kicker {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.15em;
  color: var(--music-purple);
}

.headline {
  font-family: 'Playfair Display', serif;
  font-size: 42px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.1;
  margin: 0;

  @media (max-width: 768px) {
    font-size: 32px;
  }
}

.deck {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 8px 0 0;
  max-width: 300px;
  line-height: 1.5;
}

.stat-pill {
  display: flex;
  align-items: baseline;
  gap: 6px;
  padding: 12px 20px;
  background: linear-gradient(135deg, var(--accent-ink), #3d4f5f);
  border-radius: 100px;
  color: white;

  .stat-number {
    font-size: 24px;
    font-weight: 700;
  }

  .stat-label {
    font-size: 13px;
    opacity: 0.8;
  }
}

// 数据指标栏
.metrics-bar {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 16px 24px;
  background: var(--bg-secondary);
  border-radius: 16px;
  margin-bottom: 24px;

  @media (max-width: 768px) {
    justify-content: space-around;
    gap: 0;
  }
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  font-family: 'Playfair Display', serif;
}

.metric-label {
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.metric-divider {
  width: 1px;
  height: 32px;
  background: var(--bg-tertiary);

  @media (max-width: 768px) {
    display: none;
  }
}

// 表格视图样式
.table-view {
  margin-bottom: 24px;
}

.table-container {
  background: var(--bg-secondary);
  border-radius: 20px;
  padding: 8px;
  box-shadow: var(--shadow-soft);
}

:deep(.n-data-table) {
  --n-th-font-weight: 600 !important;
  --n-th-text-color: var(--text-secondary) !important;
  --n-border-color: transparent !important;

  .n-data-table-th {
    background: transparent !important;
    padding: 16px 20px !important;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  .n-data-table-td {
    background: transparent !important;
    padding: 16px 20px !important;
    border-bottom: 1px solid var(--bg-tertiary) !important;
  }

  .n-data-table-tr:last-child .n-data-table-td {
    border-bottom: none !important;
  }

  .n-data-table-tr:hover .n-data-table-td {
    background: rgba(255, 255, 255, 0.5) !important;
  }
}

.table-cover-wrapper {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s ease;

  &:hover {
    transform: scale(1.05);
  }
}

.table-info-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.table-playlist-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.table-playlist-id {
  font-size: 12px;
  color: var(--text-muted);
  font-family: monospace;
}

.table-stats {
  display: flex;
  gap: 10px;
}

.table-stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(124, 111, 174, 0.1);
  border-radius: 20px;
  font-size: 12px;
  color: var(--music-purple);
  font-weight: 600;

  &.plays {
    background: rgba(224, 122, 95, 0.1);
    color: var(--play-orange);
  }

  .n-icon {
    opacity: 0.7;
  }
}

.table-owner {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-secondary);

  .n-icon {
    color: var(--text-muted);
  }
}

.table-delete-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-muted);

  &:hover {
    background: rgba(199, 91, 91, 0.1);
    color: var(--danger-red);
  }
}

// 网格/卡片视图（移动端）
.grid-view {
  margin-bottom: 24px;
}

.playlist-cards {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.playlist-card {
  position: relative;
  background: var(--bg-secondary);
  border-radius: 20px;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: var(--shadow-soft);

  &:active {
    transform: scale(0.98);
  }
}

.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--music-purple), var(--accent-coral));
}

.card-content {
  padding: 20px;
}

.playlist-header {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.cover-wrapper {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.playlist-meta {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 6px;
}

.playlist-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.playlist-id {
  font-size: 12px;
  color: var(--text-muted);
  font-family: monospace;
}

.playlist-stats {
  display: flex;
  gap: 10px;
  margin-bottom: 12px;
}

.stat-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: rgba(124, 111, 174, 0.1);
  border-radius: 20px;
  font-size: 13px;
  color: var(--music-purple);
  font-weight: 600;

  &.plays {
    background: rgba(224, 122, 95, 0.1);
    color: var(--play-orange);
  }

  .n-icon {
    opacity: 0.7;
  }
}

.playlist-owner {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 16px;

  .n-icon {
    color: var(--text-muted);
  }
}

.delete-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: rgba(199, 91, 91, 0.1);
  color: var(--danger-red);
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: var(--danger-red);
    color: white;
  }
}

// 分页器
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}
</style>