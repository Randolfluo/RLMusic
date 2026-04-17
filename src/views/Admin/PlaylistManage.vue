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
      <header class="editorial-header" data-animate="fade-up">
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
      <section class="stats-row" data-animate="fade-up" data-delay="1">
        <div class="stat-card">
          <div class="stat-icon purple"><n-icon :component="MusicList" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ pagination.itemCount }}</span>
            <span class="stat-name">总歌单数</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon teal"><n-icon :component="FileText" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ describedCount }}</span>
            <span class="stat-name">有 AI 描述</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon coral"><n-icon :component="Music" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ totalSongs }}</span>
            <span class="stat-name">总歌曲数</span>
          </div>
        </div>
      </section>

      <!-- 搜索控制面板 -->
      <section class="control-panel" data-animate="fade-up" data-delay="2">
        <div class="search-wrapper">
          <div class="search-icon">
            <n-icon :component="Search" />
          </div>
          <input
            v-model="searchText"
            type="text"
            placeholder="搜索歌单名称..."
            class="search-input"
            @keyup.enter="handleSearch"
          />
          <button class="search-action" @click="handleSearch">
            搜索
          </button>
        </div>

        <div class="header-actions">
          <button class="ai-batch-btn" @click="handleBatchAIDesc" :disabled="loading.aiBatch">
            <n-icon :component="Lightning" />
            <span>{{ loading.aiBatch ? '生成中...' : '批量 AI 描述' }}</span>
          </button>
          <button class="create-btn" @click="openCreateModal">
            <n-icon :component="Plus" />
            <span>新建歌单</span>
          </button>
        </div>
      </section>

      <!-- 桌面端自定义表格 -->
      <Transition name="fade-scale">
        <div v-if="!isMobile" class="custom-table-view" data-animate="fade-up" data-delay="3">
          <div class="table-header">
            <div class="th cover-col">封面</div>
            <div class="th info-col">歌单信息</div>
            <div class="th desc-col">描述</div>
            <div class="th stats-col">统计</div>
            <div class="th owner-col">创建者</div>
            <div class="th action-col">操作</div>
          </div>
          <div class="table-body">
            <div
              v-for="(playlist, index) in filteredPlaylists"
              :key="playlist.id"
              class="table-row"
              :style="{ animationDelay: `${index * 40}ms` }"
            >
              <div class="td cover-col">
                <div class="table-cover-wrapper">
                  <img
                    :src="resolveCoverUrl(playlist.cover_url) || '/images/pic/default.png'"
                    alt="cover"
                    class="table-cover"
                  />
                </div>
              </div>
              <div class="td info-col">
                <div class="table-info-cell">
                  <span class="table-title">{{ playlist.title }}</span>
                  <span class="table-id">#{{ playlist.id }}</span>
                </div>
              </div>
              <div class="td desc-col">
                <div class="table-desc">
                  <span v-if="playlist.description">{{ playlist.description }}</span>
                  <span v-else class="empty-desc">暂无描述</span>
                </div>
              </div>
              <div class="td stats-col">
                <div class="table-stats">
                  <div class="table-stat-item">
                    <n-icon :component="Music" size="14" />
                    <span>{{ playlist.total_songs || 0 }} 首</span>
                  </div>
                  <div class="table-stat-item plays">
                    <n-icon :component="PlayOne" size="14" />
                    <span>{{ (playlist.play_count || 0).toLocaleString() }}</span>
                  </div>
                </div>
              </div>
              <div class="td owner-col">
                <div class="table-owner">
                  <n-icon :component="User" size="14" />
                  <span>UID {{ playlist.owner_id }}</span>
                </div>
              </div>
              <div class="td action-col">
                <div class="table-actions">
                  <button class="table-btn edit" @click="openEditModal(playlist)" title="编辑">
                    <n-icon :component="Edit" size="16" />
                  </button>
                  <button class="table-btn ai" @click="handleGenerateDesc(playlist)" :disabled="loading.ai[playlist.id]" title="生成 AI 描述">
                    <n-icon :component="Lightning" size="16" />
                  </button>
                  <n-popconfirm @positive-click="handleDelete(playlist)" positive-text="确认删除" negative-text="取消">
                    <template #trigger>
                      <button class="table-btn delete" title="删除">
                        <n-icon :component="Delete" size="16" />
                      </button>
                    </template>
                    确定要删除该公共歌单吗？此操作不可逆！
                  </n-popconfirm>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 移动端卡片视图 -->
      <Transition name="fade-scale">
        <div v-if="isMobile" class="grid-view">
          <div class="playlist-cards">
            <div
              v-for="(playlist, index) in filteredPlaylists"
              :key="playlist.id"
              class="playlist-card"
              :style="{ animationDelay: `${index * 60}ms` }"
            >
              <div class="card-shimmer"></div>
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

                <div class="playlist-desc" v-if="playlist.description">
                  {{ playlist.description }}
                </div>
                <div class="playlist-desc empty" v-else>暂无描述</div>

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

                <div class="card-actions">
                  <button class="action-btn secondary" @click="openEditModal(playlist)">
                    <n-icon :component="Edit" />
                    编辑
                  </button>
                  <button class="action-btn ai" @click="handleGenerateDesc(playlist)" :disabled="loading.ai[playlist.id]">
                    <n-icon :component="Lightning" />
                    AI
                  </button>
                  <button class="action-btn danger" @click="confirmDelete(playlist)">
                    <n-icon :component="Delete" />
                    删除
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 空状态 -->
      <div v-if="filteredPlaylists.length === 0 && !loading.list" class="empty-state">
        <div class="empty-icon"><n-icon :component="SearchEmpty" size="48" /></div>
        <p>未找到匹配的歌单</p>
      </div>

      <!-- 分页器 -->
      <div class="pagination-wrapper" v-if="filteredPlaylists.length > 0">
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

    <!-- 创建/编辑模态框 -->
    <Transition name="modal">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-panel">
          <div class="modal-header">
            <h3>{{ isEditing ? '编辑歌单' : '新建公共歌单' }}</h3>
            <button class="modal-close" @click="closeModal"><n-icon :component="Close" /></button>
          </div>
          <div class="modal-body">
            <div class="form-group" v-if="isEditing">
              <label>歌单封面</label>
              <div class="cover-upload-area" @click="triggerCoverUpload">
                <img v-if="coverPreviewUrl" :src="coverPreviewUrl" class="cover-preview" />
                <div v-else class="cover-placeholder">
                  <n-icon :component="Camera" size="24" />
                  <span>点击上传封面</span>
                </div>
                <div v-if="coverUploadLoading" class="cover-loading">
                  <span>上传中...</span>
                </div>
              </div>
              <input ref="coverInputRef" type="file" accept="image/*" style="display: none" @change="handleCoverChange" />
            </div>
            <div class="form-group">
              <label>歌单名称</label>
              <input v-model="form.title" type="text" placeholder="输入歌单名称" />
            </div>
            <div class="form-group">
              <label>描述</label>
              <textarea v-model="form.description" rows="4" placeholder="输入歌单描述"></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn-secondary" @click="closeModal">取消</button>
            <button class="btn-primary" :class="{ 'btn-dark': isEditing }" @click="submitForm" :disabled="loading.submit || !form.title.trim()">
              {{ loading.submit ? '保存中...' : (isEditing ? '保存修改' : '创建歌单') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, NPagination, NPopconfirm, NIcon } from 'naive-ui';
import { Left, Search, Delete, Music, PlayOne, User, Lightning, Plus, Edit, Close, Camera } from '@icon-park/vue-next';
import { getPublicPlaylists, deletePublicPlaylist, createPrivatePlaylist, updatePlaylist, uploadPlaylistCover } from '@/api/playlist';
import { getPlaylistAIDescription, generatePublicPlaylistDescriptions } from '@/api/ai';
import { resolveCoverUrl } from "@/api/song";
import { ResultCode } from "@/utils/request";

const SearchEmpty = Search;
const MusicList = Music;
const FileText = Edit;

const router = useRouter();
const message = useMessage();

const windowWidth = ref(window.innerWidth);
const isMobile = computed(() => windowWidth.value < 768);

const handleResize = () => {
  windowWidth.value = window.innerWidth;
};

const loading = reactive({
  list: false,
  submit: false,
  aiBatch: false,
  ai: {} as Record<number, boolean>,
});

const searchText = ref('');
const playlistList = ref<any[]>([]);

const pagination = reactive({
  page: 1,
  pageSize: 12,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [12, 24, 48],
});

const describedCount = computed(() => playlistList.value.filter(p => p.description && p.description.trim()).length);
const totalSongs = computed(() => playlistList.value.reduce((sum, p) => sum + (p.total_songs || 0), 0));
const filteredPlaylists = computed(() => {
  if (!searchText.value.trim()) return playlistList.value;
  const q = searchText.value.toLowerCase();
  return playlistList.value.filter(p =>
    (p.title && p.title.toLowerCase().includes(q)) ||
    (p.description && p.description.toLowerCase().includes(q))
  );
});

const showModal = ref(false);
const isEditing = ref(false);
const editingId = ref<number | null>(null);
const form = reactive({ title: '', description: '' });
const coverInputRef = ref<HTMLInputElement | null>(null);
const coverPreviewUrl = ref('');
const coverUploadLoading = ref(false);

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  pagination.page = 1;
  fetchPlaylists();
};

const fetchPlaylists = async () => {
  loading.list = true;
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
    loading.list = false;
  }
};

const handleSearch = () => {
  pagination.page = 1;
  fetchPlaylists();
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

const openCreateModal = () => {
  isEditing.value = false;
  editingId.value = null;
  form.title = '';
  form.description = '';
  coverPreviewUrl.value = '';
  showModal.value = true;
};

const openEditModal = (row: any) => {
  isEditing.value = true;
  editingId.value = row.id;
  form.title = row.title || '';
  form.description = row.description || '';
  coverPreviewUrl.value = resolveCoverUrl(row.cover_url) || '/images/pic/default.png';
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
  coverPreviewUrl.value = '';
  if (coverInputRef.value) coverInputRef.value.value = '';
};

const triggerCoverUpload = () => {
  coverInputRef.value?.click();
};

const handleCoverChange = async (e: Event) => {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file || !editingId.value) return;

  // 本地预览
  coverPreviewUrl.value = URL.createObjectURL(file);
  coverUploadLoading.value = true;

  try {
    const formData = new FormData();
    formData.append('file', file);
    const res: any = await uploadPlaylistCover(editingId.value, formData);
    if (res.code === ResultCode.SUCCESS) {
      message.success('封面上传成功');
      fetchPlaylists();
    } else {
      message.error(res.message || '上传失败');
    }
  } catch (error) {
    message.error('上传失败');
  } finally {
    coverUploadLoading.value = false;
    target.value = '';
  }
};

const submitForm = async () => {
  if (!form.title.trim()) return;
  loading.submit = true;
  try {
    if (isEditing.value && editingId.value) {
      const res = await updatePlaylist(editingId.value, {
        title: form.title.trim(),
        description: form.description.trim(),
        is_public: true,
      });
      if (res.code === ResultCode.SUCCESS) {
        message.success('保存成功');
        closeModal();
        fetchPlaylists();
      } else {
        message.error(res.message || '保存失败');
      }
    } else {
      // 先创建私有歌单，再改为公共
      const createRes = await createPrivatePlaylist({
        title: form.title.trim(),
        description: form.description.trim(),
      });
      if (createRes.code === ResultCode.SUCCESS && createRes.data?.id) {
        await updatePlaylist(createRes.data.id, { is_public: true });
        message.success('创建成功');
        closeModal();
        fetchPlaylists();
      } else {
        message.error(createRes.message || '创建失败');
      }
    }
  } catch (error) {
    message.error(isEditing.value ? '保存失败' : '创建失败');
  } finally {
    loading.submit = false;
  }
};

const handleGenerateDesc = async (row: any) => {
  loading.ai[row.id] = true;
  try {
    const res = await getPlaylistAIDescription(row.id);
    if (res.code === ResultCode.SUCCESS) {
      message.success('AI 描述生成成功');
      fetchPlaylists();
    } else {
      message.error(res.message || '生成失败');
    }
  } catch (error) {
    message.error('生成失败');
  } finally {
    loading.ai[row.id] = false;
  }
};

const handleBatchAIDesc = async () => {
  loading.aiBatch = true;
  try {
    const res = await generatePublicPlaylistDescriptions();
    if (res.code === ResultCode.SUCCESS) {
      message.success(res.message || '批量生成任务已启动');
    } else {
      message.error(res.message || '启动失败');
    }
  } catch (error) {
    message.error('启动失败');
  } finally {
    loading.aiBatch = false;
  }
};

onMounted(() => {
  fetchPlaylists();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=DM+Serif+Display:ital@0;1&family=Outfit:wght@300;400;500;600;700&display=swap');

:root {
  --bg-cream: #fdfbf7;
  --bg-warm: #f7f3ec;
  --bg-card: #ffffff;
  --text-ink: #1f140e;
  --text-body: #5c4d41;
  --text-muted: #998e84;
  --accent-coral: #c75b39;
  --accent-teal: #2d6a4f;
  --accent-purple: #7c6fae;
  --accent-gold: #c9a227;
  --border-light: rgba(31, 20, 14, 0.08);
  --shadow-sm: 0 2px 8px rgba(31, 20, 14, 0.04);
  --shadow-md: 0 8px 24px rgba(31, 20, 14, 0.08);
  --shadow-lg: 0 16px 40px rgba(31, 20, 14, 0.12);
}

.playlist-manage-container {
  min-height: 100vh;
  background: var(--bg-cream);
  font-family: 'Outfit', sans-serif;
  position: relative;
  overflow-x: hidden;
}

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
  filter: blur(90px);
  opacity: 0.45;
  animation: blob-float 22s ease-in-out infinite;

  &.blob-1 {
    width: 650px;
    height: 650px;
    background: linear-gradient(135deg, rgba(124, 111, 174, 0.24), rgba(201, 162, 39, 0.14));
    top: -220px;
    right: -120px;
    animation-delay: 0s;
  }

  &.blob-2 {
    width: 550px;
    height: 550px;
    background: linear-gradient(135deg, rgba(199, 91, 57, 0.18), rgba(45, 106, 79, 0.12));
    bottom: -180px;
    left: -180px;
    animation-delay: -8s;
  }

  &.blob-3 {
    width: 450px;
    height: 450px;
    background: linear-gradient(135deg, rgba(45, 106, 79, 0.14), rgba(124, 111, 174, 0.1));
    top: 45%;
    left: 55%;
    transform: translate(-50%, -50%);
    animation-delay: -16s;
  }
}

@keyframes blob-float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(25px, -35px) scale(1.06); }
  66% { transform: translate(-15px, 25px) scale(0.94); }
}

.grain-overlay {
  position: absolute;
  inset: 0;
  opacity: 0.025;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}

.content-wrapper {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: 48px 32px 120px;

  @media (max-width: 768px) {
    padding: 24px 16px 100px;
  }
}

[data-animate="fade-up"] {
  opacity: 0;
  transform: translateY(24px);
  animation: fadeUp 0.7s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
[data-animate="fade-up"][data-delay="1"] { animation-delay: 80ms; }
[data-animate="fade-up"][data-delay="2"] { animation-delay: 160ms; }
[data-animate="fade-up"][data-delay="3"] { animation-delay: 240ms; }

@keyframes fadeUp {
  to { opacity: 1; transform: translateY(0); }
}

.editorial-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--border-light);

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    margin-bottom: 24px;
  }
}

.header-left {
  display: flex;
  align-items: flex-start;
  gap: 18px;
}

.back-btn {
  width: 46px;
  height: 46px;
  border-radius: 14px;
  border: 1px solid var(--border-light);
  background: var(--bg-warm);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.25s ease;
  color: var(--text-body);

  &:hover {
    background: #fff;
    transform: translateX(-3px);
    color: var(--text-ink);
    box-shadow: var(--shadow-sm);
  }
}

.title-stack {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.kicker {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.18em;
  color: var(--accent-purple);
}

.headline {
  font-family: 'DM Serif Display', serif;
  font-size: 44px;
  font-weight: 400;
  color: var(--text-ink);
  line-height: 1.05;
  margin: 0;

  @media (max-width: 768px) {
    font-size: 34px;
  }
}

.deck {
  font-size: 14px;
  color: var(--text-body);
  margin: 6px 0 0;
  max-width: 280px;
  line-height: 1.5;
}

.stat-pill {
  display: flex;
  align-items: baseline;
  gap: 6px;
  padding: 12px 22px;
  background: linear-gradient(135deg, #2c2420, #4a3f38);
  border-radius: 100px;
  color: white;
  box-shadow: var(--shadow-md);

  .stat-number {
    font-size: 26px;
    font-weight: 600;
  }

  .stat-label {
    font-size: 13px;
    opacity: 0.85;
  }
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 28px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  border-radius: 20px;
  padding: 18px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  border: 1px solid var(--border-light);
  box-shadow: var(--shadow-sm);
  transition: transform 0.25s ease, box-shadow 0.25s ease;

  &:hover {
    transform: translateY(-3px);
    box-shadow: var(--shadow-md);
  }
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;

  &.coral { background: rgba(199, 91, 57, 0.12); color: var(--accent-coral); }
  &.purple { background: rgba(124, 111, 174, 0.14); color: var(--accent-purple); }
  &.teal { background: rgba(45, 106, 79, 0.12); color: var(--accent-teal); }
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-value {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-ink);
}

.stat-name {
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.control-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 22px;
  gap: 16px;

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
  }
}

.search-wrapper {
  flex: 1;
  max-width: 420px;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid var(--border-light);
  border-radius: 16px;
  padding: 4px;
  transition: all 0.25s ease;

  &:focus-within {
    border-color: rgba(124, 111, 174, 0.35);
    box-shadow: 0 0 0 4px rgba(124, 111, 174, 0.08);
    background: #fff;
  }

  @media (max-width: 768px) {
    max-width: none;
  }
}

.search-icon {
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  padding: 10px 8px;
  font-size: 14px;
  color: var(--text-ink);
  outline: none;
  font-family: 'Outfit', sans-serif;

  &::placeholder {
    color: var(--text-muted);
  }
}

.search-action {
  padding: 10px 20px;
  background: var(--text-ink);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #3d2e25;
    transform: translateY(-1px);
  }
}

.header-actions {
  display: flex;
  gap: 10px;

  @media (max-width: 768px) {
    width: 100%;
  }
}

.ai-batch-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 12px;
  border: 1px solid rgba(199, 91, 57, 0.25);
  background: rgba(199, 91, 57, 0.06);
  color: var(--accent-coral);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover:not(:disabled) {
    background: rgba(199, 91, 57, 0.12);
    transform: translateY(-1px);
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  @media (max-width: 768px) {
    flex: 1;
    justify-content: center;
  }
}

.create-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  border-radius: 12px;
  border: none;
  background: var(--text-ink);
  color: white;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #3d2e25;
    transform: translateY(-1px);
  }

  @media (max-width: 768px) {
    flex: 1;
    justify-content: center;
  }
}

.custom-table-view {
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(16px);
  border-radius: 24px;
  border: 1px solid var(--border-light);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  margin-bottom: 24px;
}

.table-header {
  display: grid;
  grid-template-columns: 80px 200px 1fr 160px 130px 130px;
  align-items: center;
  padding: 14px 20px;
  background: rgba(247, 243, 236, 0.6);
  border-bottom: 1px solid var(--border-light);

  @media (max-width: 1100px) {
    grid-template-columns: 72px 160px 1fr 140px 110px 110px;
  }
}

.th {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
}

.table-body {
  .table-row {
    display: grid;
    grid-template-columns: 80px 200px 1fr 160px 130px 130px;
    align-items: center;
    padding: 14px 20px;
    border-bottom: 1px solid var(--border-light);
    transition: background 0.2s ease;
    opacity: 0;
    transform: translateY(12px);
    animation: rowIn 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;

    @media (max-width: 1100px) {
      grid-template-columns: 72px 160px 1fr 140px 110px 110px;
    }

    &:last-child {
      border-bottom: none;
    }

    &:hover {
      background: rgba(255, 255, 255, 0.9);
    }
  }
}

@keyframes rowIn {
  to { opacity: 1; transform: translateY(0); }
}

.table-cover-wrapper {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: transform 0.25s ease;

  &:hover {
    transform: scale(1.06) rotate(-2deg);
  }

  .table-cover {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.table-info-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.table-title {
  font-weight: 500;
  color: var(--text-ink);
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.table-id {
  font-size: 11px;
  color: var(--text-muted);
  font-family: 'DM Mono', monospace;
}

.table-desc {
  font-size: 13px;
  color: var(--text-body);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;

  .empty-desc {
    color: var(--text-muted);
    font-style: italic;
  }
}

.table-stats {
  display: flex;
  gap: 8px;
}

.table-stat-item {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  background: rgba(124, 111, 174, 0.1);
  border-radius: 20px;
  font-size: 11px;
  color: var(--accent-purple);
  font-weight: 600;

  &.plays {
    background: rgba(199, 91, 57, 0.1);
    color: var(--accent-coral);
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
  color: var(--text-body);

  .n-icon {
    color: var(--text-muted);
  }
}

.table-actions {
  display: flex;
  gap: 6px;
  justify-content: flex-end;
}

.table-btn {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-muted);

  &.edit:hover {
    background: rgba(45, 106, 79, 0.1);
    color: var(--accent-teal);
  }

  &.ai:hover {
    background: rgba(201, 162, 39, 0.12);
    color: #a67c00;
  }

  &.ai:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &.delete:hover {
    background: rgba(199, 91, 57, 0.1);
    color: var(--accent-coral);
  }
}

.grid-view {
  margin-bottom: 24px;
}

.playlist-cards {
  display: grid;
  grid-template-columns: 1fr;
  gap: 14px;
}

.playlist-card {
  position: relative;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 22px;
  overflow: hidden;
  transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  opacity: 0;
  transform: translateY(16px);
  animation: cardIn 0.55s cubic-bezier(0.16, 1, 0.3, 1) forwards;

  &:active {
    transform: scale(0.985);
  }
}

@keyframes cardIn {
  to { opacity: 1; transform: translateY(0); }
}

.card-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(105deg, transparent 40%, rgba(255,255,255,0.6) 50%, transparent 60%);
  transform: translateX(-100%);
  transition: transform 0s;
  pointer-events: none;
}

.playlist-card:hover .card-shimmer {
  transform: translateX(100%);
  transition: transform 0.8s ease;
}

.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent-purple), var(--accent-coral));
}

.card-content {
  padding: 18px;
}

.playlist-header {
  display: flex;
  gap: 14px;
  margin-bottom: 14px;
}

.cover-wrapper {
  width: 86px;
  height: 86px;
  border-radius: 16px;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.12);
  transition: transform 0.3s ease;
}

.playlist-card:hover .cover-wrapper {
  transform: scale(1.04) rotate(-2deg);
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
  gap: 4px;
}

.playlist-title {
  font-size: 17px;
  font-weight: 500;
  color: var(--text-ink);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.playlist-id {
  font-size: 11px;
  color: var(--text-muted);
  font-family: 'DM Mono', monospace;
}

.playlist-desc {
  font-size: 13px;
  color: var(--text-body);
  line-height: 1.55;
  margin-bottom: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;

  &.empty {
    color: var(--text-muted);
    font-style: italic;
  }
}

.playlist-stats {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.stat-badge {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 7px 12px;
  background: rgba(124, 111, 174, 0.1);
  border-radius: 20px;
  font-size: 12px;
  color: var(--accent-purple);
  font-weight: 600;

  &.plays {
    background: rgba(199, 91, 57, 0.1);
    color: var(--accent-coral);
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
  color: var(--text-body);
  margin-bottom: 14px;

  .n-icon {
    color: var(--text-muted);
  }
}

.card-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  padding: 10px 12px;
  border-radius: 12px;
  border: none;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &.secondary {
    background: var(--bg-warm);
    color: var(--text-body);

    &:hover {
      background: var(--text-ink);
      color: white;
    }
  }

  &.ai {
    background: rgba(201, 162, 39, 0.1);
    color: #8f6d1a;

    &:hover:not(:disabled) {
      background: #c9a227;
      color: white;
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  &.danger {
    background: rgba(199, 91, 57, 0.08);
    color: var(--accent-coral);

    &:hover {
      background: var(--accent-coral);
      color: white;
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-muted);
  gap: 12px;

  .empty-icon {
    opacity: 0.5;
  }

  p {
    font-size: 14px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

// 模态框
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 200;
  background: rgba(31, 20, 14, 0.35);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.modal-panel {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid var(--border-light);
  box-shadow: var(--shadow-lg);
  width: 100%;
  max-width: 480px;
  overflow: hidden;
  animation: modalIn 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalIn {
  from { opacity: 0; transform: translateY(20px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 22px;
  border-bottom: 1px solid var(--border-light);

  h3 {
    font-family: 'DM Serif Display', serif;
    font-size: 22px;
    font-weight: 400;
    color: var(--text-ink);
    margin: 0;
  }
}

.modal-close {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  border: none;
  background: var(--bg-warm);
  color: var(--text-body);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: var(--border-light);
    color: var(--text-ink);
  }
}

.modal-body {
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;

  label {
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--text-muted);
  }

  input, textarea {
    padding: 12px 14px;
    border: 1px solid var(--border-light);
    border-radius: 14px;
    background: rgba(255,255,255,0.8);
    font-size: 14px;
    color: var(--text-ink);
    font-family: 'Outfit', sans-serif;
    outline: none;
    transition: all 0.2s ease;

    &:focus {
      border-color: rgba(124, 111, 174, 0.4);
      box-shadow: 0 0 0 3px rgba(124, 111, 174, 0.08);
    }

    &::placeholder {
      color: var(--text-muted);
    }
  }

  textarea {
    resize: vertical;
    min-height: 90px;
  }
}

.cover-upload-area {
  width: 120px;
  height: 120px;
  border-radius: 16px;
  border: 2px dashed var(--border-light);
  background: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  position: relative;
  transition: all 0.2s ease;

  &:hover {
    border-color: rgba(124, 111, 174, 0.5);
    background: rgba(255, 255, 255, 0.9);
  }

  .cover-preview {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 14px;
  }

  .cover-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    color: var(--text-muted);
    font-size: 12px;
  }

  .cover-loading {
    position: absolute;
    inset: 0;
    background: rgba(255, 255, 255, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    color: var(--text-body);
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 22px;
  border-top: 1px solid var(--border-light);
}

.btn-secondary {
  padding: 10px 18px;
  border-radius: 12px;
  border: 1px solid var(--border-light);
  background: transparent;
  color: var(--text-body);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: var(--bg-warm);
  }
}

.btn-primary {
  padding: 10px 20px;
  border-radius: 12px;
  border: none;
  background: var(--text-ink);
  color: white;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover:not(:disabled) {
    background: #3d2e25;
    transform: translateY(-1px);
  }

  &.btn-dark {
    background: #000000;

    &:hover:not(:disabled) {
      background: #1a1a1a;
    }
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

// 过渡动画
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.3s ease;
}
.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.98);
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.25s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
