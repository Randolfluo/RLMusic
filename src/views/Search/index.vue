<template>
  <div class="search-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
    </div>

    <div class="search-header-container">
      <div class="search-header glass-card">
        <n-input
          ref="inputRef"
          v-model:value="inputValue"
          class="main-input"
          size="large"
          round
          placeholder="搜索音乐、歌手、歌单、专辑..."
          clearable
          @input="handleInput"
          @keydown.enter="handleEnter"
        >
          <template #prefix>
            <n-icon :component="Search" class="search-icon" />
          </template>
        </n-input>
        <n-button size="large" type="primary" round class="search-btn" @click="handleEnter">
          搜索
        </n-button>
      </div>
    </div>

    <div class="search-content">
      <!-- 默认显示：历史记录 -->
      <transition name="fade" mode="out-in">
        <div v-if="!inputValue" class="default-content">
          <div class="section history" v-if="music.getSearchHistory[0] && setting.searchHistory">
            <div class="section-header">
              <h3>
                <n-icon :component="History" />
                搜索历史
              </h3>
              <n-button text size="small" class="clear-btn" @click="delHistory">
                <template #icon>
                  <n-icon :component="DeleteFour" />
                </template>
                清空
              </n-button>
            </div>
            <div class="history-tags">
              <n-tag
                v-for="item in music.getSearchHistory"
                :key="item"
                round
                checkable
                class="history-tag"
                @click="clickHistory(item)"
              >
                {{ item }}
              </n-tag>
            </div>
          </div>
          
          <!-- 热搜榜 (Optional, if you want to display it) -->
          <div class="section hot-search" v-if="searchData.hot && searchData.hot.length">
             <div class="section-header">
              <h3>
                <n-icon :component="Fire" color="#f55e55" />
                热门搜索
              </h3>
            </div>
            <div class="hot-tags">
               <div 
                v-for="(item, index) in searchData.hot" 
                :key="index" 
                class="hot-item glass-card-sm"
                @click="clickHistory(item.first)"
               >
                 <div class="rank-num" :class="`rank-${index + 1}`">{{ index + 1 }}</div>
                 <div class="hot-info">
                   <span class="keyword">{{ item.first }}</span>
                   <!-- 模拟热度值 -->
                   <span class="hot-score" v-if="index < 3">
                     <n-icon :component="Fire" color="#f55e55" size="12" />
                     {{ 100 - index * 10 }}w
                   </span>
                 </div>
               </div>
            </div>
          </div>
        </div>

        <!-- 搜索结果区域 -->
        <div v-else class="suggest-content">
          <div class="search-title-bar">
            <h2>"{{ inputValue }}" 的搜索结果</h2>
          </div>
          
          <n-tabs type="segment" animated v-model:value="searchType" class="custom-tabs">
            <!-- 单曲 Tab -->
            <n-tab-pane name="songs" tab="单曲">
                 <div class="result-container">
                    <n-skeleton v-if="loading" text :repeat="6" :sharp="false" size="medium" style="margin-bottom: 12px;" />
                    <SongList v-else :songs="searchData.suggest.songs" :loading="loading" />
                 </div>
            </n-tab-pane>

            <!-- 歌单 Tab -->
            <n-tab-pane name="playlists" tab="歌单">
                <div class="grid-container" v-if="loading">
                   <div v-for="i in 8" :key="i" class="grid-card">
                      <n-skeleton height="100%" width="100%" style="aspect-ratio: 1; border-radius: 16px;" />
                      <n-skeleton text style="margin-top: 8px; width: 60%;" />
                   </div>
                </div>
                <n-empty v-else-if="!searchData.suggest.playlists?.length" description="暂无歌单" class="empty-state" />
                <div v-else class="grid-container">
                     <div 
                        v-for="(playlist, index) in searchData.suggest.playlists" 
                        :key="playlist.id"
                        class="grid-card playlist-card"
                        :style="{ animationDelay: `${index * 0.05}s` }"
                        @click="router.push(`/playlist/${playlist.id}`)"
                     >
                         <div class="card-cover">
                             <n-image 
                                :src="playlist.coverImgUrl || playlist.picUrl" 
                                preview-disabled 
                                object-fit="cover"
                                fallback-src="/images/logo/logo.png"
                             >
                                <template #placeholder>
                                    <div class="placeholder-icon">
                                        <n-icon :component="Record" />
                                    </div>
                                </template>
                             </n-image>
                             <div class="play-overlay">
                                 <n-icon :component="PlayOne" />
                             </div>
                         </div>
                         <div class="card-info">
                             <div class="card-title">{{ playlist.name || playlist.title }}</div>
                             <div class="card-desc">{{ playlist.description || playlist.creator?.nickname || '暂无描述' }}</div>
                         </div>
                     </div>
                </div>
            </n-tab-pane>

            <!-- 歌手 Tab -->
             <n-tab-pane name="artists" tab="歌手">
                <div class="grid-container artist-grid" v-if="loading">
                   <div v-for="i in 8" :key="i" class="grid-card artist-card">
                      <n-skeleton circle height="120px" width="120px" style="margin: 0 auto 12px;" />
                      <n-skeleton text style="width: 50%; margin: 0 auto;" />
                   </div>
                </div>
                <n-empty v-else-if="!searchData.suggest.artists?.length" description="暂无歌手" class="empty-state" />
                 <div v-else class="grid-container artist-grid">
                     <div 
                        v-for="(artist, index) in searchData.suggest.artists" 
                        :key="artist.id"
                        class="grid-card artist-card"
                        :style="{ animationDelay: `${index * 0.05}s` }"
                        @click="router.push(`/artist?id=${artist.id}`)"
                     >
                         <div class="artist-avatar">
                             <n-avatar round :size="120" :src="artist.picUrl || artist.img1v1Url" fallback-src="/images/logo/logo.png" object-fit="cover" />
                         </div>
                         <div class="card-title">{{ artist.name }}</div>
                     </div>
                </div>
            </n-tab-pane>

             <!-- 专辑 Tab -->
             <n-tab-pane name="albums" tab="专辑">
                 <div class="grid-container" v-if="loading">
                   <div v-for="i in 8" :key="i" class="grid-card">
                      <n-skeleton height="100%" width="100%" style="aspect-ratio: 1; border-radius: 16px;" />
                      <n-skeleton text style="margin-top: 8px; width: 60%;" />
                   </div>
                </div>
                 <n-empty v-else-if="!searchData.suggest.albums?.length" description="暂无专辑" class="empty-state" />
                 <div v-else class="grid-container">
                      <div 
                        v-for="(album, index) in searchData.suggest.albums" 
                        :key="album.id"
                        class="grid-card album-card"
                        :style="{ animationDelay: `${index * 0.05}s` }"
                        @click="router.push(`/album?id=${album.id}`)"
                      >
                         <div class="card-cover">
                             <n-image 
                                :src="album.picUrl" 
                                preview-disabled 
                                object-fit="cover"
                                fallback-src="/images/logo/logo.png"
                             />
                             <div class="album-bg"></div>
                         </div>
                         <div class="card-info">
                             <div class="card-title">{{ album.name }}</div>
                             <div class="card-desc">{{ album.artist?.name }}</div>
                         </div>
                      </div>
                 </div>
            </n-tab-pane>

          </n-tabs>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useMessage, useDialog, NSkeleton } from "naive-ui";
import SongList from "@/components/DataList/SongList.vue";
import { 
    Search, DeleteFour, History, Fire, PlayOne,
    Voice, RecordDisc, Record,
} from "@icon-park/vue-next";
import { getSearchHot, getSearchSuggest } from "@/api/search";
import { useMusicDataStore } from "@/store/musicData";
import { useSettingDataStore } from "@/store/settingData";
import debounce from "@/utils/debounce";

const router = useRouter();
const route = useRoute();
const music = useMusicDataStore();
const setting = useSettingDataStore();

const dialog = useDialog();

const inputRef = ref<any>(null); // Use any to support component ref methods
const inputValue = ref("");
const loading = ref(false);
const searchType = ref("songs");

const searchData = reactive<{
    hot: any[];
    suggest: {
        songs: any[];
        artists: any[];
        albums: any[];
        playlists: any[];
    }
}>({
    hot: [],
    suggest: {
        songs: [],
        artists: [],
        albums: [],
        playlists: []
    }
});

// 初始化
onMounted(() => {
    // 聚焦输入框
    nextTick(() => {
        setTimeout(() => {
             inputRef.value?.focus();
        }, 300);
    });
    
    // 获取热搜
    getSearchHot().then((res: any) => {
        if(res.data) searchData.hot = res.data;
        // console.log("Hot search:", res.data); // Debug
    });

    // 如果路由带有 query
    if (route.query.q) {
        inputValue.value = route.query.q as string;
        performSearch(inputValue.value);
    }
});

const handleInput = (val: string) => {
    if (!val.trim()) {
        loading.value = false;
        return;
    }
    loading.value = true;
    debounce(() => {
        performSearch(val);
    }, 500);
};

const performSearch = async (val: string) => {
    if(!val.trim()) return;
    try {
        const res = await getSearchSuggest(val.trim());
        loading.value = false;
        if(res) {
            // Check if result structure matches expected, adjust if needed based on API response
            // Assuming res returns objects with properties like 'songs', 'playlists' etc.
            // If the API returns differently, mapping logic should be here.
            searchData.suggest = {
                songs: res.songs || [],
                artists: res.artists || [],
                albums: res.albums || [],
                playlists: res.playlists || []
            };
        }
    } catch (e) {
        loading.value = false;
        console.error(e);
    }
}

const handleEnter = () => {
    if(inputValue.value.trim()) {
        music.setSearchHistory(inputValue.value.trim());
        performSearch(inputValue.value.trim());
    }
}

const clickHistory = (val: string) => {
    inputValue.value = val;
    performSearch(val);
}

const delHistory = () => {
    dialog.warning({
        title: "删除历史",
        content: "确认删除全部历史记录?",
        positiveText: "删除",
        negativeText: "取消",
        onPositiveClick: () => {
            music.setSearchHistory(null, true);
        }
    })
}

</script>

<style lang="scss" scoped>
.search-page {
    padding: 40px;
    max-width: 1200px;
    margin: 0 auto;
    min-height: 80vh;
    position: relative;
    overflow: hidden;

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
            opacity: 0.4;
            animation: float 20s infinite ease-in-out;
        }

        .blob-1 {
            width: 500px;
            height: 500px;
            background: var(--n-primary-color);
            top: -100px;
            left: -100px;
            animation-delay: 0s;
        }

        .blob-2 {
            width: 400px;
            height: 400px;
            background: #4facfe;
            bottom: -50px;
            right: -50px;
            animation-delay: -5s;
        }
    }
    
    .search-header-container {
        display: flex;
        justify-content: center;
        margin-bottom: 40px;
        position: relative;
        z-index: 10;
        
        .search-header {
            width: 100%;
            max-width: 700px;
            display: flex;
            gap: 16px;
            padding: 10px;
            transition: all 0.3s ease;
            
            &.glass-card {
                background: rgba(255, 255, 255, 0.6);
                backdrop-filter: blur(20px);
                border-radius: 50px;
                box-shadow: 0 10px 40px rgba(0,0,0,0.1);
                border: 1px solid rgba(255, 255, 255, 0.4);
                
                @media (prefers-color-scheme: dark) {
                    background: rgba(30, 30, 30, 0.6);
                    border: 1px solid rgba(255, 255, 255, 0.08);
                    box-shadow: 0 10px 40px rgba(0,0,0,0.3);
                }
            }
            
            &:focus-within {
                transform: scale(1.02);
                box-shadow: 0 15px 50px rgba(var(--n-primary-color-rgb), 0.15);
                border-color: var(--n-primary-color);
            }

            .main-input {
                flex: 1;
                background: transparent;
                
                :deep(.n-input-wrapper) {
                    padding: 0 10px;
                }
                :deep(.n-input__input-el) {
                    height: 48px;
                    font-size: 16px;
                }
                .search-icon {
                    font-size: 20px;
                    color: var(--n-text-color-3);
                }
            }
            
            .search-btn {
                height: 48px;
                padding: 0 32px;
                font-size: 16px;
                font-weight: bold;
                box-shadow: 0 4px 15px rgba(var(--n-primary-color-rgb), 0.3);
            }
        }
    }

    .search-content {
        position: relative;
        z-index: 5;

        .section {
            margin-bottom: 50px;
            
            .section-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                margin-bottom: 24px;
                
                h3 {
                    font-size: 20px;
                    font-weight: 800;
                    display: flex;
                    align-items: center;
                    gap: 10px;
                    margin: 0;
                    color: var(--n-text-color);
                    letter-spacing: 0.5px;
                }
                
                .clear-btn {
                    color: var(--n-text-color-3);
                    &:hover {
                        color: var(--n-error-color);
                    }
                }
            }
            
            .history-tags {
                display: flex;
                flex-wrap: wrap;
                gap: 12px;
                
                .history-tag {
                    padding: 6px 18px;
                    font-size: 14px;
                    cursor: pointer;
                    transition: all 0.2s cubic-bezier(0.25, 0.8, 0.25, 1);
                    border: none;
                    background: rgba(255, 255, 255, 0.5);
                    backdrop-filter: blur(10px);
                    box-shadow: 0 2px 8px rgba(0,0,0,0.05);
                    
                    @media (prefers-color-scheme: dark) {
                         background: rgba(255, 255, 255, 0.08);
                    }

                    &:hover {
                        transform: translateY(-3px);
                        color: var(--n-primary-color);
                        box-shadow: 0 5px 15px rgba(var(--n-primary-color-rgb), 0.2);
                    }
                }
            }
            
            .hot-tags {
                display: grid;
                grid-template-columns: repeat(2, 1fr);
                gap: 20px;
                
                @media (min-width: 768px) {
                    grid-template-columns: repeat(3, 1fr);
                }
                @media (min-width: 1024px) {
                    grid-template-columns: repeat(4, 1fr);
                }
                
                .hot-item {
                    display: flex;
                    align-items: center;
                    gap: 16px;
                    padding: 16px;
                    border-radius: 16px;
                    cursor: pointer;
                    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
                    position: relative;
                    overflow: hidden;
                    
                    &.glass-card-sm {
                        background: rgba(255, 255, 255, 0.4);
                        backdrop-filter: blur(10px);
                        border: 1px solid rgba(255, 255, 255, 0.3);
                        box-shadow: 0 4px 15px rgba(0,0,0,0.03);
                        
                        @media (prefers-color-scheme: dark) {
                            background: rgba(255, 255, 255, 0.05);
                            border: 1px solid rgba(255, 255, 255, 0.05);
                        }
                    }
                    
                    &:hover {
                        background: rgba(255, 255, 255, 0.8);
                        transform: translateY(-5px) scale(1.02);
                        box-shadow: 0 10px 25px rgba(0,0,0,0.08);
                        
                        @media (prefers-color-scheme: dark) {
                            background: rgba(255, 255, 255, 0.1);
                        }
                    }
                    
                    .rank-num {
                        font-weight: 900;
                        font-size: 24px;
                        color: var(--n-text-color-3);
                        min-width: 30px;
                        text-align: center;
                        font-style: italic;
                        opacity: 0.5;
                        
                        &.rank-1 { color: #f55e55; opacity: 1; font-size: 28px; }
                        &.rank-2 { color: #ff9f43; opacity: 1; font-size: 26px; }
                        &.rank-3 { color: #feca57; opacity: 1; }
                    }
                    
                    .hot-info {
                        display: flex;
                        flex-direction: column;
                        gap: 4px;
                        
                        .keyword {
                            font-weight: 600;
                            font-size: 16px;
                            color: var(--n-text-color);
                        }
                        
                        .hot-score {
                            font-size: 12px;
                            color: var(--n-text-color-3);
                            display: flex;
                            align-items: center;
                            gap: 4px;
                        }
                    }
                }
            }
        }
        
        .suggest-content {
            .search-title-bar {
                margin-bottom: 24px;
                h2 {
                    font-size: 24px;
                    font-weight: 700;
                    color: var(--n-text-color);
                }
            }
            
            .custom-tabs {
                :deep(.n-tabs-rail) {
                    border-radius: 30px;
                    padding: 4px;
                    background-color: rgba(0, 0, 0, 0.04);
                    @media (prefers-color-scheme: dark) {
                        background-color: rgba(255, 255, 255, 0.08);
                    }
                }
                :deep(.n-tabs-tab) {
                    border-radius: 26px;
                    font-weight: 600;
                    transition: all 0.3s;
                }
                :deep(.n-tabs-tab--active) {
                    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
                }
            }
            
            .result-container {
                margin-top: 24px;
            }
            
            .loading-container, .empty-state {
                margin: 60px 0;
                text-align: center;
            }
            
            .grid-container {
                display: grid;
                grid-template-columns: repeat(2, 1fr);
                gap: 24px;
                margin-top: 24px;
                
                @media (min-width: 640px) {
                    grid-template-columns: repeat(3, 1fr);
                }
                @media (min-width: 1024px) {
                    grid-template-columns: repeat(4, 1fr);
                }
                @media (min-width: 1280px) {
                    grid-template-columns: repeat(5, 1fr);
                }
                
                &.artist-grid {
                    @media (min-width: 640px) {
                        grid-template-columns: repeat(4, 1fr);
                    }
                    @media (min-width: 1024px) {
                        grid-template-columns: repeat(6, 1fr);
                    }
                }
                
                .grid-card {
                    cursor: pointer;
                    transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
                    animation: fade-up 0.6s ease backwards;
                    
                    &:hover {
                        transform: translateY(-10px);
                        .card-cover .play-overlay {
                            opacity: 1;
                            transform: scale(1);
                        }
                    }
                    
                    .card-cover {
                        position: relative;
                        aspect-ratio: 1;
                        border-radius: 16px;
                        overflow: hidden;
                        box-shadow: 0 8px 24px rgba(0,0,0,0.1);
                        margin-bottom: 12px;
                        
                        .n-image {
                            width: 100%;
                            height: 100%;
                            display: block;
                            
                            :deep(img) {
                                width: 100%;
                                height: 100%;
                                object-fit: cover;
                                transition: transform 0.6s ease;
                            }
                        }
                        
                        &:hover :deep(img) {
                            transform: scale(1.12);
                        }
                        
                        .play-overlay {
                            position: absolute;
                            inset: 0;
                            background: rgba(0,0,0,0.3);
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            opacity: 0;
                            transform: scale(0.8);
                            transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
                            color: white;
                            font-size: 48px;
                            backdrop-filter: blur(4px);
                        }
                        
                        .placeholder-icon {
                            width: 100%;
                            height: 100%;
                            background: var(--n-action-color);
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            color: var(--n-text-color-3);
                            font-size: 48px;
                        }
                        
                        .album-bg {
                            position: absolute;
                            top: 0;
                            right: -10px;
                            width: 100%;
                            height: 100%;
                            background: #1a1a1a;
                            border-radius: 50%;
                            z-index: -1;
                            transition: right 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
                            box-shadow: 2px 0 10px rgba(0,0,0,0.3);
                            
                            @media (prefers-color-scheme: dark) {
                                background: #000;
                            }
                        }
                    }
                    
                    &.album-card:hover .album-bg {
                        right: -25px;
                    }
                    
                    .card-info {
                        .card-title {
                            font-weight: 700;
                            font-size: 16px;
                            margin-bottom: 4px;
                            white-space: nowrap;
                            overflow: hidden;
                            text-overflow: ellipsis;
                            color: var(--n-text-color);
                        }
                        .card-desc {
                            font-size: 13px;
                            color: var(--n-text-color-3);
                            white-space: nowrap;
                            overflow: hidden;
                            text-overflow: ellipsis;
                        }
                    }
                    
                    &.artist-card {
                        text-align: center;
                        .artist-avatar {
                            margin-bottom: 16px;
                            transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
                            border-radius: 50%;
                            overflow: hidden;
                            box-shadow: 0 8px 20px rgba(0,0,0,0.1);
                            
                            :deep(.n-avatar) {
                                display: block;
                            }
                        }
                        
                        &:hover .artist-avatar {
                            transform: scale(1.08) rotate(3deg);
                            box-shadow: 0 15px 30px rgba(0,0,0,0.15);
                        }
                    }
                }
            }
        }
    }
}

@keyframes float {
    0% { transform: translate(0, 0) rotate(0deg); }
    33% { transform: translate(30px, -50px) rotate(10deg); }
    66% { transform: translate(-20px, 20px) rotate(-5deg); }
    100% { transform: translate(0, 0) rotate(0deg); }
}

@keyframes fade-up {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
