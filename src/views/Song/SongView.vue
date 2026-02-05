<template>
  <div class="song-view">
    <div v-if="loading" class="loading-container">
      <n-skeleton height="200px" width="200px" />
      <div class="info-placeholders">
        <n-skeleton text style="width: 50%" />
        <n-skeleton text style="width: 30%" />
        <n-skeleton text style="width: 80%" :repeat="3" />
      </div>
    </div>
    
    <div v-else-if="song" class="content">
      <!-- 头部信息 -->
      <div class="header-section">
        <div class="cover-wrapper">
          <n-image
            :src="coverUrl"
            class="cover-img"
            object-fit="contain"
            fallback-src="/images/logo/logo.png"
            :preview-disabled="false"
          />
        </div>
        
        <div class="info-wrapper">
          <h1 class="song-title">{{ song.title }}</h1>
          
          <div class="meta-row">
            <div class="meta-item">
              <n-icon :component="User" />
              <span class="label">歌手：</span>
              <!-- Logic for displaying Artists list or fallback -->
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
              <n-icon :component="RecordDisc" />
              <span class="label">专辑：</span>
              <span class="value link" @click="router.push(`/album?id=${song.album_id}`)" v-if="song.album_id">{{ song.album_name }}</span>
              <span class="value" v-else>{{ song.album_name || '未知专辑' }}</span>
            </div>
          </div>

          <div class="actions">
            <n-button type="primary" size="large" round @click="playMusic">
              <template #icon>
                <n-icon :component="PlayOne" />
              </template>
              播放
            </n-button>
             <n-button size="large" round @click="handleLike">
               <template #icon>
                <n-icon :component="Like" />
              </template>
              收藏
            </n-button>
          </div>
        </div>
      </div>

      <!-- 简介与详情卡片 -->
      <div class="details-section">
        
        <div class="info-cards">
           <!-- 格式 -->
           <div class="info-card">
              <div class="card-label">音频格式</div>
              <div class="card-value">{{ song.format }}</div>
              <div class="card-tag">{{ ((song.bit_rate || 0) / 1000).toFixed(0) }} kbps</div>
           </div>

            <!-- 时长 -->
           <div class="info-card">
              <div class="card-label">时长</div>
              <div class="card-value">{{ formatTime(song.duration) }}</div>
              <div class="card-tag">Time</div>
           </div>

           <!-- 采样率 -->
           <div class="info-card">
              <div class="card-label">采样率</div>
              <div class="card-value">{{ song.sample_rate }} Hz</div>
              <div class="card-tag">{{ song.bit_depth }} bit</div>
           </div>
           
           <!-- 声道/大小 -->
           <div class="info-card">
              <div class="card-label">声道/大小</div>
              <div class="card-value">{{ song.channels === 2 ? 'Stereo' : (song.channels === 1 ? 'Mono' : song.channels + ' Ch') }}</div>
              <div class="card-tag">{{ (song.file_size / 1024 / 1024).toFixed(2) }} MB</div>
           </div>

            <!-- 轨道/年份 -->
           <div class="info-card">
              <div class="card-label">发行信息</div>
              <div class="card-value">{{ song.year || '未知年份' }}</div>
              <div class="card-tag">
                  <span v-if="song.disc_num">Disk {{ song.disc_num }} / </span>
                  Track {{ song.track_num }}
              </div>
           </div>

           <!-- 播放次数 -->
           <div class="info-card">
              <div class="card-label">播放统计</div>
              <div class="card-value">{{ song.play_count }}</div>
              <div class="card-tag">Plays</div>
           </div>

           <!-- 文件信息 -->
           <div class="info-card" :title="song.file_path">
              <div class="card-label">源文件</div>
              <div class="card-value" style="font-size: 14px; line-height: 1.5; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ song.file_name }}</div>
              <div class="card-tag">ID: {{ song.id }} | Path: {{ song.file_path }}</div>
           </div>
        </div>
        
        <div class="intro-block" v-if="false">
          <!-- 预留简介区域，目前数据库没有 -->
          <h3>音乐简介</h3>
          <p class="intro-text">暂无简介...</p>
        </div>

      </div>

    </div>
    
    <div v-else class="error-container">
        <n-empty description="未找到歌曲信息" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useMessage } from "naive-ui";
import { User, RecordDisc, PlayOne, Like } from "@icon-park/vue-next";
import { getSongDetail, getSongCover, toggleLike } from "@/api/song";
import { useMusicDataStore } from "@/store/musicData";

const route = useRoute();
const router = useRouter();
const message = useMessage();
const musicStore = useMusicDataStore();

const songId = route.params.id as string;
const loading = ref(true);
// Use loose type with specific array definition for artists to ensure v-for index is number
const song = ref<{ artists?: any[]; [key: string]: any } | null>(null);

const coverUrl = computed(() => {
    if (!songId) return '';
    return getSongCover(songId);
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
})
</script>

<style lang="scss" scoped>
.song-view {
  padding: 40px;
  max-width: 1200px;
  margin: 0 auto;
  
  .loading-container {
    display: flex;
    gap: 40px;
    padding-top: 20px;
    .info-placeholders {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 20px;
    }
  }

  .header-section {
    display: flex;
    gap: 50px;
    margin-bottom: 60px;
    
    .cover-wrapper {
      flex-shrink: 0;
      width: 260px;
      height: 260px;
      border-radius: 20px;
      overflow: hidden;
      box-shadow: 0 10px 30px rgba(0,0,0,0.2);
      
      .cover-img {
        width: 100%;
        height: 100%;
        display: block; // Remove extra space
        
        :deep(img) {
             width: 100%;
             height: 100%;
             object-fit: contain !important;
        }
      }
    }

    .info-wrapper {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: center;
      
      .song-title {
        font-size: 36px;
        font-weight: 800;
        margin: 0 0 20px;
        line-height: 1.2;
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
          color: var(--n-text-color-2); // Use Naive UI var if available, or fallback
          
          .n-icon {
            margin-right: 8px;
            font-size: 18px;
          }
          
          .value {
            margin-left: 5px;
            font-weight: 500;
            &.link {
              cursor: pointer;
              color: var(--n-primary-color);
              &:hover {
                text-decoration: underline;
              }
            }
          }
        }
      }

      .actions {
        display: flex;
        gap: 20px;
      }
    }
  }

  .details-section {
    
    .info-cards {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 20px;
      margin-bottom: 40px;
      
      .info-card {
        background-color: var(--n-card-color); // Adaptation
        border-radius: 12px;
        padding: 20px;
        display: flex;
        flex-direction: column;
        // Naive UI variables usually available if using n-config-provider, else default colors
        background: #f5f5f7; 
        
        // Dark mode adaptation hint (if simple CSS)
        @media (prefers-color-scheme: dark) {
           background: #2c2c2e;
        }
        
        // If app uses specific class/attribute for dark mode, adjust accordingly. 
        // Assuming light/default for now or relying on transparent/inherit if inside n-card.
        
        .card-label {
            font-size: 13px;
            color: #86868b;
            margin-bottom: 8px;
        }
        
        .card-value {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 4px;
        }
        
        .card-tag {
            font-size: 12px;
            padding: 2px 8px;
            background: rgba(0,0,0,0.05);
            border-radius: 4px;
            width: fit-content;
        }
      }
    }

    .intro-block {
        h3 {
            font-size: 20px;
            font-weight: 700;
            margin-bottom: 15px;
            position: relative;
            padding-left: 12px;
             &::before {
                content: '';
                position: absolute;
                left: 0;
                top: 4px;
                bottom: 4px;
                width: 4px;
                background-color: var(--n-primary-color);
                border-radius: 2px;
            }
        }
        .intro-text {
            line-height: 1.8;
            color: var(--n-text-color-2);
            font-size: 15px;
        }
    }
  }
}

// Dark mode overrides if Naive UI classes available globally on body/container
:root[theme='dark'] .info-card {
    background: #333 !important;
    .card-label { color: #aaa; }
    .card-tag { background: rgba(255,255,255,0.1); }
}
</style>
