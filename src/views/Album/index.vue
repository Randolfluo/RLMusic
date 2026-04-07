<template>
  <div class="album-page">
    <!-- Background Decoration -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
    </div>

    <div class="album-content">
      <n-spin :show="loading">
        <!-- Album Header -->
        <div class="header-section glass-panel">
          <div class="cover-wrapper">
            <n-image
              class="cover-img"
              :src="
                resolveCoverUrl(album.cover) ||
                (album.cover_song_id ? getSongCover(album.cover_song_id) : null) ||
                resolveCoverUrl(album.picUrl) ||
                resolveCoverUrl(album.cover_url) ||
                '/images/logo/favicon.png'
              "
              fallback-src="/images/logo/favicon.png"
              object-fit="cover"
              preview-disabled
              @click="music.setBigPlayerState(true)"
            />
            <div class="cover-overlay">
              <n-icon :component="PlayOne" size="48" />
            </div>
          </div>
          <div class="info-wrapper">
            <div class="tag-badge">Album</div>
            <h1 class="album-title">{{ album.title || album.name }}</h1>
            <div class="artist-info" v-if="album.artist_name || (album.artist && album.artist.name)">
              <n-icon :component="User" size="16" />
              <span>{{ album.artist_name || (album.artist && album.artist.name) }}</span>
            </div>
            <div class="album-desc" v-if="album.description">
              {{ album.description }}
            </div>
            <div class="meta-info">
              <span class="song-count">{{ (album.songs || []).length }} 首歌曲</span>
              <span class="divider" v-if="album.publish_time">•</span>
              <span class="publish-time" v-if="album.publish_time">{{ album.publish_time }}</span>
            </div>
            <div class="actions">
              <n-button type="primary" round size="large" class="play-btn" @click="playAll">
                <template #icon>
                  <n-icon :component="Play" />
                </template>
                播放全部
              </n-button>
            </div>
          </div>
        </div>

        <!-- Songs Section -->
        <div class="songs-section glass-panel">
          <div class="section-header">
            <div class="section-icon">
              <n-icon :component="Music" size="20" />
            </div>
            <h2 class="section-title">专辑歌曲</h2>
          </div>
          <SongList
            :songs="album.songs || []"
            :loading="loading"
          />
        </div>
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import { getAlbumDetail, getSongCover, resolveCoverUrl } from "@/api/song";
import { ResultCode } from "@/utils/request";
import { useMessage, NButton, NIcon, NImage, NSpin } from "naive-ui";
import { Play, PlayOne, User, Music } from "@icon-park/vue-next";
import { musicStore } from "@/store";
import SongList from "@/components/DataList/SongList.vue";

const route = useRoute();
const message = useMessage();
const music = musicStore();

const loading = ref(false);
const album = ref<any>({});

const fetchAlbumDetail = async (id: string) => {
  if (!id) return;

  loading.value = true;
  try {
    const res = await getAlbumDetail(id);
    if (res.code === ResultCode.SUCCESS) {
      album.value = res.data;
    } else {
      message.error(res.message || "获取专辑详情失败");
    }
  } catch (error) {
    message.error("获取专辑详情失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  const id = route.query.id as string;
  if (id) {
    fetchAlbumDetail(id);
  } else {
    message.error("未找到专辑ID");
  }
});

watch(
  () => route.query.id,
  (newId) => {
    if (newId) {
      fetchAlbumDetail(newId as string);
    }
  }
);

const playAll = () => {
  if (album.value.songs && album.value.songs.length > 0) {
    const tracks = album.value.songs.map((song: any) => ({
      ...song,
      name: song.title,
      artist: song.artists || [{ name: song.artist_name, id: song.artist_id }],
      album: {
        name: album.value.title || album.value.name,
        id: album.value.id,
        picUrl: resolveCoverUrl(album.value.picUrl) || resolveCoverUrl(album.value.cover_url)
      },
      picUrl: song.cover_url ? resolveCoverUrl(song.cover_url) : resolveCoverUrl(album.value.picUrl) || resolveCoverUrl(album.value.cover_url)
    }));

    music.setPlaylists(tracks);
    music.setPlaySongIndex(0);
    music.setPlayState(true);
  }
}
</script>

<style scoped lang="scss">
.album-page {
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
      background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
      top: -100px;
      right: -100px;
      animation-delay: 0s;
    }

    .blob-2 {
      width: 400px;
      height: 400px;
      background: linear-gradient(135deg, #06b6d4 0%, #22d3ee 100%);
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

.album-content {
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
    width: 240px;
    height: 240px;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
    position: relative;
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    cursor: pointer;

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
      background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
      color: white;
      width: fit-content;
      margin-bottom: 16px;
      box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
    }

    .album-title {
      font-family: 'Plus Jakarta Sans', sans-serif;
      font-size: 36px;
      font-weight: 800;
      margin: 0 0 12px 0;
      color: var(--n-text-color);
      line-height: 1.2;
      letter-spacing: -0.02em;
    }

    .artist-info {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 15px;
      color: var(--n-text-color-2);
      margin-bottom: 12px;
      font-weight: 500;
    }

    .album-desc {
      font-size: 14px;
      color: var(--n-text-color-3);
      line-height: 1.6;
      margin: 0 0 16px 0;
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

      .song-count {
        font-weight: 600;
        padding: 6px 14px;
        background: rgba(6, 182, 212, 0.1);
        border-radius: 100px;
        color: #06b6d4;
      }

      .divider {
        opacity: 0.5;
      }
    }

    .actions {
      display: flex;
      gap: 12px;

      .play-btn {
        background: linear-gradient(135deg, #10b981 0%, #059669 100%);
        border: none;
        padding: 0 28px;
        height: 48px;
        font-size: 16px;
        font-weight: 600;
        box-shadow: 0 8px 20px rgba(16, 185, 129, 0.35);
        transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 12px 28px rgba(16, 185, 129, 0.45);
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

  .section-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;

    .section-icon {
      width: 40px;
      height: 40px;
      border-radius: 12px;
      background: linear-gradient(135deg, #06b6d4 0%, #22d3ee 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      box-shadow: 0 4px 12px rgba(6, 182, 212, 0.3);
    }

    .section-title {
      font-family: 'Plus Jakarta Sans', sans-serif;
      font-size: 20px;
      font-weight: 700;
      margin: 0;
      color: var(--n-text-color);
    }
  }
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .album-page {
    padding: 20px;
  }

  .header-section {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 24px;
    gap: 24px;

    .cover-wrapper {
      width: 180px;
      height: 180px;
    }

    .info-wrapper {
      min-height: auto;
      align-items: center;

      .tag-badge {
        align-self: center;
      }

      .album-title {
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
  .album-page {
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
