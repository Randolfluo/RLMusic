<template>
  <div class="song-view-container">
    <!-- 动态背景 -->
    <div class="dynamic-background" :style="{ backgroundImage: `url(${coverUrl || '/images/logo/logo.png'})` }">
      <div class="backdrop-blur"></div>
    </div>

    <div class="song-view-content">
      <div v-if="loading" class="loading-container">
        <n-skeleton height="260px" width="260px" style="border-radius: 20px" />
        <div class="info-placeholders">
          <n-skeleton text style="width: 50%; height: 40px; margin-bottom: 20px" />
          <n-skeleton text style="width: 30%; height: 24px; margin-bottom: 12px" />
          <n-skeleton text style="width: 30%; height: 24px; margin-bottom: 30px" />
          <div style="display: flex; gap: 20px">
             <n-skeleton width="120px" height="40px" round />
             <n-skeleton width="120px" height="40px" round />
          </div>
        </div>
      </div>
      
      <div v-else-if="song" class="content-wrapper">
        <!-- 头部信息 -->
        <div class="header-section glass-panel">
          <div class="cover-wrapper">
            <n-image
              :src="coverUrl"
              class="cover-img"
              object-fit="cover"
              fallback-src="/images/logo/logo.png"
              :preview-disabled="false"
            />
          </div>
          
          <div class="info-wrapper">
            <h1 class="song-title">{{ song.title }}</h1>
            
            <div class="meta-row">
              <div class="meta-item">
                <n-icon :component="User" class="icon" />
                <span class="label">歌手：</span>
                <span v-if="song.artists && song.artists.length > 0">
                  <template v-for="(artist, index) in song.artists" :key="artist.id">
                    <span class="value link" @click="router.push(`/artist?id=${artist.id}`)">{{ artist.name }}</span>
                    <span v-if="index < song.artists.length - 1" class="separator"> / </span>
                  </template>
                </span>
                <span class="value link" @click="router.push(`/artist?id=${song.artist_id}`)" v-else-if="song.artist_id">{{ song.artist_name }}</span>
                <span class="value" v-else>{{ song.artist_name || '未知歌手' }}</span>
              </div>
              
              <div class="meta-item">
                <n-icon :component="RecordDisc" class="icon" />
                <span class="label">专辑：</span>
                <span class="value link" @click="router.push(`/album?id=${song.album_id}`)" v-if="song.album_id">{{ song.album_name }}</span>
                <span class="value" v-else>{{ song.album_name || '未知专辑' }}</span>
              </div>

              <div class="meta-item desc-item" v-if="song.description">
                <n-icon :component="BookOne" class="icon" />
                <span class="label">简介：</span>
                <span class="value desc-text" :class="{ expanded: descExpanded }">{{ song.description }}</span>
                <span class="desc-toggle" @click="descExpanded = !descExpanded">
                  {{ descExpanded ? '收起' : '展开' }}
                </span>
              </div>
            </div>

            <div class="actions">
              <n-button type="primary" size="large" round class="action-btn play-btn" @click="playMusic">
                <template #icon>
                  <n-icon :component="PlayOne" />
                </template>
                <span class="btn-text">立即播放</span>
              </n-button>
              <n-button size="large" round class="action-btn like-btn" @click="handleLike">
                <template #icon>
                  <n-icon :component="Like" />
                </template>
                <span class="btn-text">收藏</span>
              </n-button>
            </div>
          </div>
        </div>

        <!-- 详情卡片区域 -->
        <div class="details-grid">
           <!-- 格式 -->
           <div class="detail-card glass-card">
              <div class="card-icon">
                <n-icon :component="FileCodeOne" />
              </div>
              <div class="card-content">
                <div class="card-label">音频格式</div>
                <div class="card-value">{{ song.format }}</div>
                <div class="card-tag">{{ ((song.bit_rate || 0) / 1000).toFixed(0) }} kbps</div>
              </div>
           </div>

            <!-- 时长 -->
           <div class="detail-card glass-card">
              <div class="card-icon">
                <n-icon :component="Time" />
              </div>
              <div class="card-content">
                <div class="card-label">时长</div>
                <div class="card-value">{{ formatTime(song.duration) }}</div>
                <div class="card-tag">Time</div>
              </div>
           </div>

           <!-- 采样率 -->
           <div class="detail-card glass-card">
              <div class="card-icon">
                <n-icon :component="Voice" />
              </div>
              <div class="card-content">
                <div class="card-label">采样率</div>
                <div class="card-value">{{ song.sample_rate }} Hz</div>
                <div class="card-tag">{{ song.bit_depth }} bit</div>
              </div>
           </div>
           
           <!-- 声道/大小 -->
           <div class="detail-card glass-card">
              <div class="card-icon">
                <n-icon :component="DatabaseNetwork" />
              </div>
              <div class="card-content">
                <div class="card-label">声道/大小</div>
                <div class="card-value">{{ song.channels === 2 ? 'Stereo' : (song.channels === 1 ? 'Mono' : song.channels + ' Ch') }}</div>
                <div class="card-tag">{{ (song.file_size / 1024 / 1024).toFixed(2) }} MB</div>
              </div>
           </div>

            <!-- 轨道/年份 -->
           <div class="detail-card glass-card">
              <div class="card-icon">
                <n-icon :component="Calendar" />
              </div>
              <div class="card-content">
                <div class="card-label">发行信息</div>
                <div class="card-value">{{ song.year || '未知年份' }}</div>
                <div class="card-tag">
                  <span v-if="song.disc_num">Disk {{ song.disc_num }} / </span>
                  Track {{ song.track_num }}
                </div>
              </div>
           </div>

           <!-- 播放次数 -->
           <div class="detail-card glass-card">
              <div class="card-icon">
                <n-icon :component="Play" />
              </div>
              <div class="card-content">
                <div class="card-label">播放统计</div>
                <div class="card-value">{{ song.play_count }}</div>
                <div class="card-tag">Plays</div>
              </div>
           </div>

           <!-- 文件信息 (Clickable) -->
           <div class="detail-card glass-card clickable" @click="showFileInfo = true" :title="song.file_path">
              <div class="card-icon">
                <n-icon :component="FolderCode" />
              </div>
              <div class="card-content">
                <div class="card-label">源文件</div>
                <div class="card-value file-name">{{ song.file_name }}</div>
                <div class="card-tag">点击查看详细路径</div>
              </div>
           </div>
        </div>

        <!-- AI 开场白 -->
        <div v-if="song.opening_file" class="ai-section glass-panel">
          <div class="section-header">
            <div class="section-icon">
              <n-icon :component="VolumeNotice" />
            </div>
            <div class="section-title">AI 智能开场白</div>
          </div>
          <audio :src="openingAudioUrl" controls class="audio-player"></audio>
        </div>

        <!-- 歌词展示 -->
        <div v-if="lyrics.length > 0" class="lyrics-section glass-panel">
          <div class="section-header">
            <div class="section-icon">
              <n-icon :component="BookOne" />
            </div>
            <div class="section-title">歌词</div>
            <n-button quaternary size="small" class="toggle-btn" @click="showLyrics = !showLyrics">
              {{ showLyrics ? '收起' : '展开' }}
            </n-button>
          </div>
          <n-collapse-transition :show="showLyrics">
            <div class="lyrics-content">
              <p v-for="(line, index) in lyrics" :key="index" class="lyric-line">
                <span v-if="line.lyric" class="lyric-text">{{ line.lyric }}</span>
                <span v-if="line.lyricFy" class="lyric-trans">{{ line.lyricFy }}</span>
              </p>
            </div>
          </n-collapse-transition>
        </div>
        <div v-else-if="!lyricsLoading" class="lyrics-section glass-panel empty">
          <div class="section-header">
            <div class="section-icon">
              <n-icon :component="BookOne" />
            </div>
            <div class="section-title">歌词</div>
          </div>
          <div class="empty-tip">暂无歌词</div>
        </div>

        <!-- File Detail Modal -->
        <n-modal v-model:show="showFileInfo">
            <n-card
                style="width: 600px; max-width: 90vw;"
                title="文件详细信息"
                :bordered="false"
                size="huge"
                role="dialog"
                aria-modal="true"
                class="modal-card"
            >
                <div class="file-details">
                    <div class="detail-row">
                        <span class="label">文件ID:</span>
                        <span class="value">{{ song.id }}</span>
                    </div>
                    <div class="detail-row">
                        <span class="label">文件名:</span>
                        <span class="value">{{ song.file_name }}</span>
                    </div>
                    <div class="detail-row">
                        <span class="label">完整路径:</span>
                        <span class="value path-value">{{ song.file_path }}</span>
                    </div>
                     <div class="detail-row">
                        <span class="label">文件大小:</span>
                        <span class="value">{{ song.file_size }} bytes ({{ (song.file_size / 1024 / 1024).toFixed(2) }} MB)</span>
                    </div>
                     <div class="detail-row">
                        <span class="label">哈希值:</span>
                        <span class="value">{{ song.hash || 'N/A' }}</span>
                    </div>
                </div>
            </n-card>
        </n-modal>
      </div>
      
      <div v-else class="error-container">
          <n-empty description="未找到歌曲信息">
            <template #extra>
              <n-button size="small" @click="router.back()">
                返回上一页
              </n-button>
            </template>
          </n-empty>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useMessage } from "naive-ui";
import {
  User, RecordDisc, PlayOne, Like, FileCodeOne, Time,
  Voice, DatabaseNetwork, Calendar, Play, FolderCode, BookOne,
  VolumeNotice
} from "@icon-park/vue-next";
import { getSongDetail, getSongCover, toggleLike, getMusicLyric } from "@/api/song";
import { useMusicDataStore } from "@/store/musicData";

const route = useRoute();
const router = useRouter();
const message = useMessage();
const musicStore = useMusicDataStore();

const songId = route.params.id as string;
const loading = ref(true);
const showFileInfo = ref(false); // Controls file info modal
// Use loose type with specific array definition for artists to ensure v-for index is number
const song = ref<{ artists?: any[]; [key: string]: any } | null>(null);

// Lyrics
const lyrics = ref<any[]>([]);
const lyricsLoading = ref(false);
const showLyrics = ref(typeof window !== 'undefined' ? window.innerWidth > 768 : true);

// Mobile description expand
const descExpanded = ref(false);

const coverUrl = computed(() => {
    if (!songId) return '';
    return getSongCover(songId);
});

const openingAudioUrl = computed(() => {
    if (!song.value?.opening_file) return '';
    return `/api/song/opening/${song.value.id}`;
});

const loadData = async () => {
    loading.value = true;
    try {
        const res: any = await getSongDetail(songId);
        // Server returns code 1000 for success
        if (res.code === 1000 && res.data) {
             song.value = res.data;
        } else if (res.id) {
            // Fallback if structure is different
            song.value = res;
        } else {
            message.warning(res.message || "获取信息失败");
        }
    } catch (e) {
        console.error(e);
        message.error("加载歌曲信息失败");
    } finally {
        loading.value = false;
    }
}

const loadLyrics = async () => {
    lyricsLoading.value = true;
    try {
        const data = await getMusicLyric(songId);
        lyrics.value = data || [];
    } catch (e) {
        console.error(e);
    } finally {
        lyricsLoading.value = false;
    }
}


const playMusic = () => {
    if (song.value) {
        const musicItem = {
            id: song.value.id,
            name: song.value.title,
            artist: song.value.artist_name,
            album: song.value.album_name,
            cover: getSongCover(song.value.id),
            src: `/api/song/stream/${song.value.id}`,
            ...song.value
        };
        // Add to list and play
        musicStore.setPlaylists([musicItem]);
        musicStore.setPlaySongIndex(0);
        musicStore.setPlayState(true);
    }
}

const handleLike = async () => {
    try {
        await toggleLike(songId);
        message.success("操作成功");
    } catch (error) {
         message.error("操作失败");
    }
}

const formatTime = (seconds: number) => {
    if (!seconds) return "00:00";
    const min = Math.floor(seconds / 60);
    const sec = Math.floor(seconds % 60);
    return `${min.toString().padStart(2, '0')}:${sec.toString().padStart(2, '0')}`;
}

onMounted(() => {
    loadData();
    loadLyrics();
})
</script>

<style lang="scss" scoped>
.song-view-container {
  background: #faf8f5;
  position: relative;
  min-height: 100vh;
  width: 100%;
  overflow: hidden;

  .dynamic-background {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-size: cover;
    background-position: center;
    z-index: 0;
    
    .backdrop-blur {
      width: 100%;
      height: 100%;
      backdrop-filter: blur(60px) saturate(180%);
      background: rgba(255, 255, 255, 0.5);
      mask-image: linear-gradient(to bottom, rgba(0,0,0,1) 0%, rgba(0,0,0,0.6) 100%);
      
      @media (prefers-color-scheme: dark) {
        background: rgba(0, 0, 0, 0.6);
      }
    }
  }

  .song-view-content {
    position: relative;
    z-index: 1;
    padding: 40px;
    max-width: 1200px;
    margin: 0 auto;

    @media (max-width: 768px) {
      padding: 16px;
    }
  }
}

.loading-container {
  display: flex;
  gap: 40px;
  padding-top: 40px;
  
  .info-placeholders {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
}

.content-wrapper {
  animation: slideUp 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.glass-panel {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 30px;
  
  @media (prefers-color-scheme: dark) {
    background: rgba(30, 30, 30, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
}

.header-section {
  display: flex;
  gap: 40px;
  margin-bottom: 30px;
  
  .cover-wrapper {
    flex-shrink: 0;
    width: 260px;
    height: 260px;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(0,0,0,0.2);
    transition: transform 0.3s ease;
    
    &:hover {
      transform: scale(1.02);
    }
    
    .cover-img {
      width: 100%;
      height: 100%;
      display: block;
      
      :deep(img) {
         width: 100%;
         height: 100%;
         object-fit: cover !important;
      }
    }
  }

  .info-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    
    .song-title {
      font-size: 42px;
      font-weight: 800;
      margin: 0 0 20px;
      line-height: 1.1;
      color: var(--n-text-color);
      text-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
  }

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
    text-align: left;
    gap: 16px;
    padding: 16px;

    .cover-wrapper {
      width: 140px;
      height: 140px;
      margin: 0 auto;
    }

    .info-wrapper {
      width: 100%;
      align-items: flex-start;

      .song-title {
        font-size: 22px;
        margin-bottom: 10px;
        text-align: left;
        width: 100%;
      }

      .meta-row {
        align-items: flex-start;
        width: 100%;
        margin-bottom: 16px;
        gap: 8px;

        .meta-item {
          font-size: 13px;
          flex-wrap: wrap;
          line-height: 1.4;

          &.desc-item {
            flex-direction: column;
            align-items: flex-start;
            gap: 2px;

            .icon, .label {
              margin-bottom: 0;
            }

            .desc-text {
              width: 100%;
              -webkit-line-clamp: 5;
              display: -webkit-box;
              -webkit-box-orient: vertical;
              overflow: hidden;
              font-size: 12px;
              line-height: 1.6;
              color: var(--n-text-color-2);

              &.expanded {
                -webkit-line-clamp: unset;
                display: block;
              }
            }

            .desc-toggle {
              font-size: 11px;
              color: var(--n-primary-color);
              cursor: pointer;
              font-weight: 600;
              margin-top: 2px;
            }
          }
        }
      }

      .actions {
        justify-content: flex-start;
        width: 100%;
        gap: 10px;

        .action-btn {
          height: 40px;
          font-size: 14px;
          padding: 0 16px;
          flex: 1;
          max-width: none;
        }
      }
    }
  }

    .meta-row {
      display: flex;
      flex-direction: column;
      gap: 12px;
      margin-bottom: 30px;
      
      .meta-item {
        display: flex;
        align-items: center;
        font-size: 16px;
        color: var(--n-text-color-2);
        
        .icon {
          margin-right: 8px;
          font-size: 20px;
          color: var(--n-primary-color);
          flex-shrink: 0;
        }
        
        .label {
          margin-right: 4px;
          font-weight: 500;
          flex-shrink: 0;
        }
        
        .value {
          font-weight: 600;
          color: var(--n-text-color-1);

          &.desc-text {
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
            font-size: 14px;
            font-weight: normal;
            line-height: 1.6;

            &.expanded {
              -webkit-line-clamp: unset;
              display: block;
            }
          }

          &.link {
            cursor: pointer;
            transition: color 0.2s;
            &:hover {
              color: var(--n-primary-color);
              text-decoration: underline;
            }
          }
        }

        .desc-toggle {
          font-size: 12px;
          color: var(--n-primary-color);
          cursor: pointer;
          font-weight: 600;
          margin-left: 8px;
          flex-shrink: 0;

          &:hover {
            text-decoration: underline;
          }
        }

        &.desc-item {
          align-items: flex-start;
          flex-wrap: wrap;
          row-gap: 4px;

          .desc-text {
            flex: 1 1 100%;
            margin-top: 4px;
          }

          .desc-toggle {
            margin-left: 0;
            margin-top: 4px;
          }
        }
        
        .separator {
          margin: 0 4px;
          color: var(--n-text-color-3);
        }
      }
    }

    .actions {
      display: flex;
      gap: 16px;
      
      .action-btn {
        padding: 0 24px;
        font-weight: bold;
        height: 48px;
        font-size: 16px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        display: flex;
        align-items: center;
        
        .btn-text {
          color: var(--n-text-color);
        }
        
        &.play-btn {
          background: linear-gradient(135deg, var(--n-primary-color) 0%, var(--n-primary-color-hover) 100%);
          border: none;
          
          .btn-text {
            color: #040404;
          }
        }
      }
    }
  }

.details-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;

  @media (max-width: 1024px) {
    grid-template-columns: repeat(3, 1fr);
  }
  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  @media (max-width: 480px) {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }
  
  .detail-card {
    background: rgba(255, 255, 255, 0.5);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    
    @media (prefers-color-scheme: dark) {
      background: rgba(40, 40, 40, 0.5);
      border: 1px solid rgba(255, 255, 255, 0.05);
    }
    
    &:hover {
      transform: translateY(-5px);
      background: rgba(255, 255, 255, 0.7);
      box-shadow: 0 10px 20px rgba(0,0,0,0.1);
      
      @media (prefers-color-scheme: dark) {
        background: rgba(60, 60, 60, 0.6);
      }
      
      .card-icon {
        transform: scale(1.1) rotate(5deg);
        background: var(--n-primary-color);
        color: white;
      }
    }
    
    &.clickable {
      cursor: pointer;
      &:active {
        transform: scale(0.98);
      }
    }
    
    .card-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      background: rgba(var(--n-primary-color-rgb), 0.1);
      color: var(--n-primary-color);
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 24px;
      transition: all 0.3s ease;
      flex-shrink: 0;
    }
    
    .card-content {
      flex: 1;
      min-width: 0;
      
      .card-label {
        font-size: 12px;
        color: var(--n-text-color-3);
        margin-bottom: 4px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
      }
      
      .card-value {
        font-size: 16px;
        font-weight: 700;
        color: var(--n-text-color);
        overflow: hidden;
        text-overflow: ellipsis;
        margin-bottom: 2px;

        &.file-name {
          font-size: 14px;
        }
      }
      
      .card-tag {
        font-size: 12px;
        color: var(--n-text-color-3);
      }
    }

    @media (max-width: 768px) {
      padding: 12px;
      gap: 8px;
      border-radius: 12px;

      .card-icon {
        width: 34px;
        height: 34px;
        border-radius: 8px;
        font-size: 16px;
      }

      .card-content {
        .card-label {
          font-size: 10px;
          margin-bottom: 2px;
        }
        .card-value {
          font-size: 12px;
          margin-bottom: 0;
        }
        .card-tag {
          font-size: 10px;
        }
      }
    }

    @media (max-width: 480px) {
      padding: 10px;
      gap: 6px;
      border-radius: 10px;

      .card-icon {
        width: 30px;
        height: 30px;
        border-radius: 6px;
        font-size: 14px;
      }

      .card-content {
        .card-label {
          font-size: 9px;
          margin-bottom: 1px;
        }
        .card-value {
          font-size: 11px;
        }
        .card-tag {
          font-size: 9px;
        }
      }
    }
  }
}

.ai-section {
  margin-top: 16px;
  padding: 20px;

  @media (max-width: 768px) {
    padding: 14px;
    margin-top: 10px;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 14px;

    .section-icon {
      width: 36px;
      height: 36px;
      border-radius: 10px;
      background: rgba(var(--n-primary-color-rgb), 0.1);
      color: var(--n-primary-color);
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 18px;
    }

    .section-title {
      font-size: 16px;
      font-weight: 700;
      color: var(--n-text-color);
    }
  }

  .audio-player {
    width: 100%;
    height: 40px;
    border-radius: 8px;
  }
}

.lyrics-section {
  margin-top: 16px;
  padding: 20px;

  @media (max-width: 768px) {
    padding: 14px;
    margin-top: 10px;
  }

  &.empty {
    .empty-tip {
      text-align: center;
      color: var(--n-text-color-3);
      font-size: 14px;
      padding: 20px 0;
    }
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 14px;

    .section-icon {
      width: 36px;
      height: 36px;
      border-radius: 10px;
      background: rgba(var(--n-primary-color-rgb), 0.1);
      color: var(--n-primary-color);
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 18px;
    }

    .section-title {
      flex: 1;
      font-size: 16px;
      font-weight: 700;
      color: var(--n-text-color);
    }

    .toggle-btn {
      margin-left: auto;
    }
  }

  .lyrics-content {
    max-height: 50vh;
    overflow-y: auto;
    padding-right: 8px;
    line-height: 1.8;
    font-size: 15px;
    color: var(--n-text-color-1);
    mask: linear-gradient(180deg, transparent 0, #fff 5%, #fff 95%, transparent 100%);
    -webkit-mask: linear-gradient(180deg, transparent 0, #fff 5%, #fff 95%, transparent 100%);

    @media (max-width: 768px) {
      font-size: 14px;
      max-height: 45vh;
    }

    .lyric-line {
      margin: 10px 0;
      display: flex;
      flex-direction: column;
      gap: 2px;

      .lyric-text {
        color: var(--n-text-color-1);
      }

      .lyric-trans {
        font-size: 13px;
        color: var(--n-text-color-3);

        @media (max-width: 768px) {
          font-size: 12px;
        }
      }
    }
  }
}

.file-details {
  display: flex;
  flex-direction: column;
  gap: 16px;
  
  .detail-row {
    display: flex;
    flex-direction: column;
    gap: 6px;
    
    .label {
      font-size: 13px;
      color: var(--n-text-color-3);
      font-weight: bold;
    }
    
    .value {
      font-size: 14px;
      color: var(--n-text-color-1);
      word-break: break-all;
      background: rgba(128, 128, 128, 0.1);
      padding: 12px;
      border-radius: 8px;
      font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
      border: 1px solid rgba(128, 128, 128, 0.1);
      
      &.path-value {
        max-height: 120px;
        overflow-y: auto;
      }
    }
  }
}

.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60vh;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(40px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
