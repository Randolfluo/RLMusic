<template>
  <div class="playlist-detail-page">
    <!-- Background Decoration -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
    </div>

    <div class="playlist-content">
      <!-- Header Section -->
      <div class="header-section">
        <div class="cover-wrapper">
          <n-image
            class="cover-img"
            :src="resolveCoverUrl(playlist.cover_url) || '/images/logo/favicon.png'"
            fallback-src="/images/logo/favicon.png"
            object-fit="cover"
            preview-disabled
            @click="music.setBigPlayerState(true)"
            style="cursor: pointer;"
          />
          <div class="cover-overlay" @click="music.setBigPlayerState(true)">
            <n-icon :component="PlayOne" size="48" />
          </div>
          <div
            v-if="canUploadCover"
            class="cover-upload-btn"
            :class="{ loading: coverUploadLoading }"
            @click.stop="triggerCoverUpload"
          >
            <n-icon :component="Camera" size="18" />
            <span v-if="coverUploadLoading">上传中...</span>
          </div>
          <input
            ref="coverInputRef"
            type="file"
            accept="image/*"
            style="display: none"
            @change="handleCoverChange"
          />
        </div>

        <div class="info-wrapper">
          <div class="tag-badge">Playlist</div>
          <h1 class="playlist-title">{{ playlist.title }}</h1>
          <p class="playlist-desc" v-if="playlist.description">
            {{ playlist.description }}
          </p>
          <div class="meta-info">
            <div class="creator" v-if="playlist.owner_id">
              <n-avatar
                round
                size="small"
                :src="playlist.owner?.avatarUrl || playlist.owner?.avatar_url || '/images/logo/favicon.png'"
                fallback-src="/images/logo/favicon.png"
              />
              <span>{{ playlist.owner?.nickname || playlist.owner?.username || `User ${playlist.owner_id}` }}</span>
            </div>
            <span class="divider" v-if="playlist.owner_id">•</span>
            <span class="song-count">{{ playlistSongCount }} 首歌曲</span>
          </div>
          <div class="actions">
            <n-button
              type="primary"
              round
              size="large"
              class="play-btn"
              :style="{
                background: `linear-gradient(135deg, ${themeColor} 0%, ${adjustColor(themeColor, -20)} 100%)`,
                boxShadow: `0 8px 20px ${themeColor}40`
              }"
              @click="playAll"
            >
              <template #icon>
                <n-icon :component="Play" />
              </template>
              播放全部
            </n-button>

            <n-button
              v-if="user.userLogin && playlist.owner_id !== user.userData.userId"
              round size="large"
              class="action-btn"
              @click="handleSubscribe"
              :loading="subLoading"
            >
              <template #icon>
                <n-icon :component="Like" :color="isSubscribed ? '#d03050' : undefined" v-if="isSubscribed" />
                <n-icon :component="Like" v-else />
              </template>
              {{ isSubscribed ? '取消收藏' : '收藏' }}
            </n-button>

            <n-button
              v-if="canGenerateIntro"
              round size="large"
              class="action-btn"
              @click="handleGenerateIntro"
            >
              <template #icon>
                <n-icon :component="Voice" />
              </template>
              生成开场白
            </n-button>

            <n-button
              v-if="canEditPlaylist"
              round size="large"
              class="action-btn"
              @click="openEditModal"
            >
              <template #icon>
                <n-icon :component="Edit" />
              </template>
              编辑歌单
            </n-button>
          </div>
        </div>
      </div>

      <!-- Songs List Section -->
      <div class="songs-section glass-panel">
        <n-spin :show="loading">
          <SongList
            :songs="playlist.songs || []"
            :loading="loading"
            :page="page"
            :page-size="limit"
            :playlist-id="playlist.id"
            :is-owner="isOwner"
            :is-public="!!playlist.is_public"
            @refresh="refreshPlaylist"
          />
          <div class="pagination-container" v-if="playlist.songs && playlist.songs.length > 0">
            <Pagination
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

    <!-- 编辑歌单弹窗 -->
    <n-modal v-model:show="showEditModal" title="编辑歌单" preset="card" style="width: 90%; max-width: 480px;">
      <n-form :model="editForm" label-placement="left" :label-width="80">
        <n-form-item label="歌单名称">
          <n-input v-model:value="editForm.title" placeholder="请输入歌单名称" maxlength="50" />
        </n-form-item>
        <n-form-item label="歌单描述">
          <n-input v-model:value="editForm.description" placeholder="请输入歌单描述" type="textarea" :autosize="{ minRows: 3, maxRows: 5 }" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showEditModal = false">取消</n-button>
          <n-button type="primary" @click="submitEdit" :loading="editLoading">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute } from "vue-router";
import {
  getPublicPlaylistDetail,
  getPrivatePlaylistDetail,
  subscribePlaylist,
  unsubscribePlaylist,
  checkIsSubscribed,
  updatePlaylist
} from "@/api/playlist";
import { generatePlaylistIntros } from "@/api/ai";
import { ResultCode } from "@/utils/request";
import {
  useMessage, NButton, NIcon, NImage, NSpin, NAvatar, useDialog,
  NModal, NForm, NFormItem, NInput, NSpace
} from "naive-ui";
import { Play, Like, Voice, PlayOne, Edit, Camera } from "@icon-park/vue-next";
import { musicStore, userStore, settingStore } from "@/store";
import Pagination from "@/components/Pagination/index.vue";
import SongList from "@/components/DataList/SongList.vue";
import { resolveCoverUrl } from "@/api/song";
import { uploadPlaylistCover } from "@/api/playlist";

const route = useRoute();
const message = useMessage();
const dialog = useDialog();
const music = musicStore();
const user = userStore();
const setting = settingStore();

// 主题色
const themeColor = computed(() => setting.themeColor);

const loading = ref(false);
const subLoading = ref(false);
const isSubscribed = ref(false);
const playlist = ref<any>({});
const page = ref(1);
const limit = ref(30);

const playlistSongCount = computed(() => {
  const totalSongs = Number(playlist.value.total_songs);
  if (Number.isFinite(totalSongs) && totalSongs > 0) return totalSongs;
  const total = Number(playlist.value.total);
  if (Number.isFinite(total) && total > 0) return total;
  if (Array.isArray(playlist.value.songs)) return playlist.value.songs.length;
  return 0;
});

// 编辑歌单弹窗
const showEditModal = ref(false);
const editForm = ref({ title: '', description: '' });
const editLoading = ref(false);

const openEditModal = () => {
  editForm.value = {
    title: playlist.value.title || '',
    description: playlist.value.description || ''
  };
  showEditModal.value = true;
};

const submitEdit = async () => {
  if (!editForm.value.title.trim()) {
    message.warning("歌单名称不能为空");
    return;
  }
  editLoading.value = true;
  try {
    const res = await updatePlaylist(playlist.value.id, {
      title: editForm.value.title.trim(),
      description: editForm.value.description.trim(),
      is_public: !!playlist.value.is_public
    });
    if (res.code === ResultCode.SUCCESS) {
      message.success("编辑成功");
      playlist.value.title = editForm.value.title.trim();
      playlist.value.description = editForm.value.description.trim();
      showEditModal.value = false;
    } else {
      message.error(res.message || "编辑失败");
    }
  } catch (error) {
    message.error("编辑失败");
  } finally {
    editLoading.value = false;
  }
};

const isOwner = computed(() => {
  if (!user.userLogin || !playlist.value.owner_id) return false;
  return Number(user.userData.userId) === Number(playlist.value.owner_id);
});

const isAdmin = computed(() => {
  return user.userLogin && user.userData.userGroup === 'admin';
});

const canGenerateIntro = computed(() => {
  return isOwner.value || isAdmin.value;
});

const canEditPlaylist = computed(() => {
  return isOwner.value || isAdmin.value;
});

const canUploadCover = computed(() => {
  // 公共歌单仅管理员可上传，私有歌单仅所有者可上传
  if (!playlist.value.id) return false;
  if (!!playlist.value.is_public) {
    return isAdmin.value;
  }
  return isOwner.value;
});

const coverInputRef = ref<HTMLInputElement | null>(null);
const coverUploadLoading = ref(false);

const triggerCoverUpload = () => {
  coverInputRef.value?.click();
};

const handleCoverChange = async (e: Event) => {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file || !playlist.value.id) return;

  coverUploadLoading.value = true;
  try {
    const formData = new FormData();
    formData.append("file", file);
    const res: any = await uploadPlaylistCover(playlist.value.id, formData);
    if (res.code === ResultCode.SUCCESS) {
      message.success("封面上传成功");
      refreshPlaylist();
    } else {
      message.error(res.message || "上传失败");
    }
  } catch (error) {
    message.error("上传失败");
  } finally {
    coverUploadLoading.value = false;
    target.value = "";
  }
};

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
    try {
      const res = await getPublicPlaylistDetail(id, page.value, limit.value);
      if (res.code === ResultCode.SUCCESS) {
        playlist.value = res.data;
        checkSubscribedStatus(id);
        return;
      }
    } catch (e) {
      // Continue to try private
    }

    const res = await getPrivatePlaylistDetail(id, page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      playlist.value = res.data;
      checkSubscribedStatus(id);
    } else {
      message.error(res.message || "获取歌单详情失败");
    }
  } catch (error) {
    message.error("获取歌单详情失败");
  } finally {
    loading.value = false;
  }
};

const checkSubscribedStatus = async (id: string) => {
  if (!user.userLogin) return;
  try {
    const res = await checkIsSubscribed(id);
    if (res.code === ResultCode.SUCCESS) {
      isSubscribed.value = res.data.is_subscribed;
    }
  } catch (e) {
    console.error(e);
  }
};

const handleSubscribe = async () => {
  if (!user.userLogin) {
    message.warning("请先登录");
    return;
  }
  const id = playlist.value.id;
  subLoading.value = true;
  try {
    let res;
    if (isSubscribed.value) {
      res = await unsubscribePlaylist(id);
    } else {
      res = await subscribePlaylist(id);
    }

    if (res.code === ResultCode.SUCCESS) {
      isSubscribed.value = !isSubscribed.value;
      message.success(isSubscribed.value ? "收藏成功" : "已取消收藏");
    } else {
      message.error(res.message || "操作失败");
    }
  } catch (e) {
    message.error("操作失败");
  } finally {
    subLoading.value = false;
  }
};

const handleGenerateIntro = () => {
  dialog.info({
    title: "生成开场白",
    content: `确定要为歌单 "${playlist.value.title}" 生成开场白吗？`,
    positiveText: "生成",
    negativeText: "取消",
    onPositiveClick: () => {
      generatePlaylistIntros(playlist.value.id)
        .then((res) => {
          if (res.code === ResultCode.SUCCESS) {
            message.success("已开始生成开场白，请稍后查看");
          } else {
            message.error(res.message || "生成失败");
          }
        })
        .catch(() => {
          message.error("请求失败");
        });
    },
  });
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

// 辅助函数：调整颜色亮度
const adjustColor = (color: string, amount: number): string => {
  const hex = color.replace('#', '');
  const r = Math.max(0, Math.min(255, parseInt(hex.slice(0, 2), 16) + amount));
  const g = Math.max(0, Math.min(255, parseInt(hex.slice(2, 4), 16) + amount));
  const b = Math.max(0, Math.min(255, parseInt(hex.slice(4, 6), 16) + amount));
  return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
};

const refreshPlaylist = () => {
  fetchPlaylistDetail(route.params.id as string);
};

const playAll = () => {
  if (playlist.value.songs && playlist.value.songs.length > 0) {
    const tracks = playlist.value.songs.map((song: any) => ({
      ...song,
      name: song.title,
      artist: song.artists && song.artists.length > 0
        ? song.artists
        : [{ name: song.artist_name, id: song.artist_id }],
      album: {
        name: song.album_title,
        id: song.album_id,
        picUrl: song.cover_url ? resolveCoverUrl(song.cover_url) : resolveCoverUrl(playlist.value.cover_url)
      }
    }));

    music.setPlaylists(tracks);
    music.setPlaySongIndex(0);
    music.setPlayState(true);
  }
};
</script>

<style scoped lang="scss">
.playlist-detail-page {
  background: #faf8f5;
  position: relative;
  min-height: 100vh;
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;

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
      background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
      top: -100px;
      right: -100px;
      animation-delay: 0s;
    }

    .blob-2 {
      width: 400px;
      height: 400px;
      background: linear-gradient(135deg, #14b8a6 0%, #5eead4 100%);
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

.playlist-content {
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
  margin-bottom: 40px;
  padding: 32px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px) saturate(180%);
  border-radius: 28px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06);

  .cover-wrapper {
    flex-shrink: 0;
    width: 240px;
    height: 240px;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
    position: relative;
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);

    &:hover {
      transform: scale(1.03) rotate(2deg);

      .cover-overlay {
        opacity: 1;
      }
    }

    .cover-img {
      width: 100%;
      height: 100%;
      transition: transform 0.6s ease;
    }

    .cover-overlay {
      position: absolute;
      inset: 0;
      background: rgba(0, 0, 0, 0.4);
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: all 0.3s ease;
      color: white;
      backdrop-filter: blur(4px);
    }

    .cover-upload-btn {
      position: absolute;
      right: 10px;
      bottom: 10px;
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 6px 12px;
      border-radius: 20px;
      background: rgba(0, 0, 0, 0.6);
      color: white;
      font-size: 12px;
      cursor: pointer;
      opacity: 0;
      transition: all 0.3s ease;
      backdrop-filter: blur(4px);

      &:hover:not(.loading) {
        background: rgba(0, 0, 0, 0.8);
        transform: scale(1.05);
      }

      &.loading {
        opacity: 1;
        cursor: default;
      }
    }

    &:hover .cover-upload-btn {
      opacity: 1;
    }
  }

  .info-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-height: 240px;

    .tag-badge {
      display: inline-block;
      font-size: 12px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 1px;
      padding: 6px 14px;
      border-radius: 100px;
      background: v-bind('`linear-gradient(135deg, ${themeColor} 0%, ${adjustColor(themeColor, 20)} 100%)`');
      color: white;
      width: fit-content;
      margin-bottom: 16px;
      box-shadow: 0 4px 12px v-bind('`${themeColor}40`');
    }

    .playlist-title {
      font-family: 'Plus Jakarta Sans', sans-serif;
      font-size: 36px;
      font-weight: 800;
      margin: 0 0 12px 0;
      color: var(--n-text-color);
      line-height: 1.2;
      letter-spacing: -0.02em;
    }

    .playlist-desc {
      font-size: 15px;
      color: var(--n-text-color-3);
      line-height: 1.6;
      margin: 0 0 20px 0;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .meta-info {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 24px;
      font-size: 14px;
      color: var(--n-text-color-3);

      .creator {
        display: flex;
        align-items: center;
        gap: 8px;
      }

      .divider {
        opacity: 0.5;
      }

      .song-count {
        font-weight: 500;
      }
    }

    .actions {
      display: flex;
      gap: 12px;

      .play-btn {
        border: none;
        padding: 0 28px;
        height: 48px;
        font-size: 16px;
        font-weight: 600;
        transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);

        &:hover {
          transform: translateY(-2px);
          filter: brightness(1.1);
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

/* Songs Section */
.songs-section {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.04);

  .pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .playlist-detail-page {
    padding: 12px;
  }

  .header-section {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 16px;
    gap: 16px;
    margin-bottom: 16px;

    .cover-wrapper {
      width: 120px;
      height: 120px;
      border-radius: 16px;

      .cover-upload-btn {
        opacity: 1;
        right: 6px;
        bottom: 6px;
        padding: 4px 8px;
        font-size: 10px;
      }
    }

    .info-wrapper {
      min-height: auto;
      align-items: center;

      .tag-badge {
        align-self: center;
        margin-bottom: 8px;
        padding: 4px 10px;
        font-size: 10px;
      }

      .playlist-title {
        font-size: 18px;
        margin-bottom: 8px;
      }

      .playlist-desc {
        font-size: 13px;
        margin-bottom: 10px;
        -webkit-line-clamp: 1;
      }

      .meta-info {
        justify-content: center;
        margin-bottom: 12px;
        font-size: 13px;
        gap: 8px;

        .n-avatar {
          width: 20px;
          height: 20px;
        }
      }

      .actions {
        flex-wrap: wrap;
        justify-content: center;
        gap: 8px;

        .play-btn {
          height: 40px;
          padding: 0 20px;
          font-size: 14px;
        }

        .action-btn {
          height: 40px;
          padding: 0 16px;
          font-size: 13px;
        }
      }
    }
  }

  .songs-section {
    padding: 12px;
    border-radius: 16px;

    .pagination-container {
      margin-top: 16px;
    }
  }
}

/* Dark Mode Support */
:global(.dark) {
  .playlist-detail-page {
    .header-section {
      background: rgba(30, 30, 30, 0.7);
      border-color: rgba(255, 255, 255, 0.1);
    }

    .songs-section {
      background: rgba(30, 30, 30, 0.6);
      border-color: rgba(255, 255, 255, 0.08);
    }
  }
}
</style>
