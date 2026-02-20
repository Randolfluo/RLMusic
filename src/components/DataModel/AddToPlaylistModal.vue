<template>
  <n-modal
    v-model:show="showModal"
    preset="card"
    title="添加到歌单"
    class="add-to-playlist-modal"
    :style="{ width: '400px' }"
    :bordered="false"
  >
    <div class="playlist-list">
      <!-- Create New Playlist Item -->
      <div class="playlist-item create-new" @click="openCreateModal">
        <div class="cover-box">
          <n-icon size="24" :component="Plus" />
        </div>
        <div class="info">
          <span class="title">创建新歌单</span>
        </div>
      </div>

      <n-divider style="margin: 8px 0" />

      <!-- Existing Playlists -->
      <n-spin :show="loading">
        <div v-if="playlists.length > 0">
            <div
            v-for="playlist in playlists"
            :key="playlist.id"
            class="playlist-item"
            @click="handleAddToPlaylist(playlist)"
            >
            <div class="cover-box">
                <n-image
                    v-if="playlist.cover_url"
                    :src="playlist.cover_url"
                    fallback-src="/images/default_cover.png"
                    preview-disabled
                    object-fit="cover"
                    width="48"
                    height="48"
                />
                <n-icon v-else size="24" :component="Music" color="#999" />
            </div>
            <div class="info">
                <span class="title">{{ playlist.title }}</span>
                <span class="count">{{ playlist.total_songs || 0 }} 首音乐</span>
            </div>
            </div>
        </div>
        <div v-else class="empty-tip">
            暂无私有歌单
        </div>
      </n-spin>
    </div>

    <!-- Create Playlist Modal -->
    <CreatePlaylistModal
      v-model:show="showCreateModal"
      @created="handlePlaylistCreated"
    />
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useMessage, NModal, NIcon, NDivider, NSpin, NImage } from 'naive-ui';
import { Plus, Music } from '@icon-park/vue-next';
import { getUserPrivatePlaylists, addSongsToPlaylist } from '@/api/playlist';
import CreatePlaylistModal from './CreatePlaylistModal.vue';

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  songIds: {
    type: Array as () => number[],
    default: () => []
  }
});

const emit = defineEmits(['update:show', 'success']);

const message = useMessage();
const showModal = ref(false);
const showCreateModal = ref(false);
const loading = ref(false);
const playlists = ref<any[]>([]);

watch(() => props.show, (val) => {
  showModal.value = val;
  if (val) {
    fetchPlaylists();
  }
});

watch(showModal, (val) => {
  emit('update:show', val);
});

const fetchPlaylists = () => {
  loading.value = true;
  // 获取全部歌单（或足够多），这里使用较大 limit
  getUserPrivatePlaylists(1, 1000)
    .then((res: any) => {
      if (res.code === 200 || res.code === 1000) {
        if (Array.isArray(res.data)) {
            playlists.value = res.data;
        } else if (res.data && Array.isArray(res.data.list)) {
            playlists.value = res.data.list;
        } else {
            playlists.value = [];
        }
      }
    })
    .catch((err) => {
      console.error(err);
      message.error('获取歌单列表失败');
    })
    .finally(() => {
      loading.value = false;
    });
};

const openCreateModal = () => {
  showCreateModal.value = true;
};

const handlePlaylistCreated = (_newPlaylist: any) => {
  // Refresh list and optionally add to it immediately?
  // For now, just refresh list
  fetchPlaylists();
  // Optional: if user wants to immediately add to the new playlist, we could do that here.
  // But usually users expect to see it in the list first.
};

const handleAddToPlaylist = (playlist: any) => {
  if (props.songIds.length === 0) {
    message.warning('未选择歌曲');
    return;
  }

  const loadingMsg = message.loading('添加中...', { duration: 0 });
  
  addSongsToPlaylist({
    playlist_id: playlist.id,
    song_ids: props.songIds
  })
    .then((res: any) => {
      loadingMsg.destroy();
      if (res.code === 200 || res.code === 1000) {
        if (res.data && res.data.count === 0) {
            message.info(res.data.message || '所选歌曲已在歌单中');
        } else {
            message.success(res.data?.message || '添加成功');
            emit('success');
        }
        showModal.value = false;
      } else {
        message.error(res.message || '添加失败');
      }
    })
    .catch((err) => {
      loadingMsg.destroy();
      message.error(err.message || '添加失败');
    });
};
</script>

<style scoped lang="scss">
.playlist-list {
  max-height: 400px;
  overflow-y: auto;
}

.playlist-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-bottom: 4px;

  &:hover {
    background-color: var(--n-action-color); // Use Naive UI var or fallback
    background-color: rgba(0, 0, 0, 0.05);
  }

  .cover-box {
    width: 48px;
    height: 48px;
    border-radius: 6px;
    background-color: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .info {
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .title {
      font-size: 14px;
      font-weight: 500;
      margin-bottom: 4px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .count {
      font-size: 12px;
      color: #999;
    }
  }
}

.create-new {
  .cover-box {
    background-color: #f0f0f0;
    color: #666;
  }
}

.empty-tip {
    text-align: center;
    padding: 20px;
    color: #999;
}
</style>
