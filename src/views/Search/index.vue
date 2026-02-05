<template>
  <div class="search-page">
    <div class="search-header">
      <n-input-group style="max-width: 600px; width: 100%">
        <n-select
          :style="{ width: '25%' }"
          v-model:value="searchType"
          :options="typeOptions"
          size="large"
          class="type-select"
        />
        <n-input
          ref="inputRef"
          v-model:value="inputValue"
          class="main-input"
          size="large"
          placeholder="搜索音乐、歌手、歌单、专辑"
          clearable
          @input="handleInput"
          @keydown.enter="handleEnter"
        >
          <template #prefix>
            <n-icon :component="Search" />
          </template>
        </n-input>
      </n-input-group>
    </div>

    <div class="search-content">
      <!-- 默认显示：历史记录和热搜 -->
      <div v-if="!inputValue" class="default-content">
        <!-- 历史记录 -->
        <div class="section history" v-if="music.getSearchHistory[0] && setting.searchHistory">
           <div class="section-title">
              <div class="left">
                  <n-icon :component="History" />
                  <span>搜索历史</span>
              </div>
               <n-button text size="small" @click="delHistory">
                   <template #icon>
                       <n-icon :component="DeleteFour" />
                   </template>
                   清空
               </n-button>
           </div>
           <n-space>
              <n-tag
                v-for="item in music.getSearchHistory"
                :key="item"
                round
                checkable
                @click="clickHistory(item)"
              >
                {{ item }}
              </n-tag>
           </n-space>
        </div>

        <!-- 热搜 (Placeholder for now as API returns empty) -->
        <!-- <div class="section hot" v-if="searchData.hot.length > 0">
             <div class="section-title">
                  <div class="left">
                      <n-icon :component="Fire" color="#d03050" />
                      <span>热搜榜</span>
                  </div>
             </div>
             <div class="hot-grid">
                 <div
                  v-for="(item, index) in searchData.hot"
                  :key="index"
                  class="hot-item"
                  @click="clickHistory(item.searchWord)"
                 >
                    <span :class="['index', index < 3 ? 'top' : '']">{{ index + 1 }}</span>
                    <div class="info">
                        <span class="word">{{ item.searchWord }}</span>
                        <span class="desc" v-if="item.content">{{ item.content }}</span>
                    </div>
                 </div>
             </div>
        </div> -->
      </div>

      <!-- 搜索建议结果 -->
      <div v-else class="suggest-content">
          <div class="loading" v-if="loading">
              <n-spin size="small" />
              <span>搜索中...</span>
          </div>
          
          <div v-else-if="isEmptyResult" class="empty">
              <n-empty description="未找到相关结果" />
          </div>

          <div v-else class="result-list">
              <!-- 单曲 -->
              <div class="result-section" v-if="shouldShow('songs') && searchData.suggest.songs?.length">
                  <div class="section-header">
                      <n-icon :component="MusicOne" />
                      <span>单曲</span>
                  </div>
                  <div class="list-items">
                      <div 
                        class="item song-item" 
                        v-for="song in searchData.suggest.songs" 
                        :key="song.id"
                        @click="playSong(song)"
                      >
                          <div class="name">{{ song.title }}</div>
                          <div class="artist">{{ song.artist_name || (song.Artist?.name) || '未知歌手' }}</div>
                          <div class="album" v-if="song.album_name || song.Album?.title">
                              - {{ song.album_name || song.Album?.title }}
                          </div>
                      </div>
                  </div>
              </div>

               <!-- 歌手 -->
              <div class="result-section" v-if="shouldShow('artists') && searchData.suggest.artists?.length">
                  <div class="section-header">
                      <n-icon :component="Voice" />
                      <span>歌手</span>
                  </div>
                  <div class="list-items">
                      <div 
                        class="item" 
                        v-for="artist in searchData.suggest.artists" 
                        :key="artist.id"
                        @click="router.push(`/artist?id=${artist.id}`)"
                      >
                          <div class="name" v-html="artist.name"></div>
                      </div>
                  </div>
              </div>

               <!-- 专辑 -->
              <div class="result-section" v-if="shouldShow('albums') && searchData.suggest.albums?.length">
                  <div class="section-header">
                      <n-icon :component="RecordDisc" />
                      <span>专辑</span>
                  </div>
                  <div class="list-items">
                      <div 
                        class="item" 
                        v-for="album in searchData.suggest.albums" 
                        :key="album.id"
                        @click="router.push(`/album?id=${album.id}`)"
                      >
                          <div class="name">{{ album.title }}</div>
                          <div class="sub">- {{ album.Artist?.name }}</div>
                      </div>
                  </div>
              </div>
              
                <!-- 歌单 -->
               <div class="result-section" v-if="shouldShow('playlists') && searchData.suggest.playlists?.length">
                  <div class="section-header">
                      <n-icon :component="Record" />
                      <span>歌单</span>
                  </div>
                  <div class="list-items">
                      <div 
                        class="item" 
                        v-for="playlist in searchData.suggest.playlists" 
                        :key="playlist.id"
                        @click="router.push(`/playlist/${playlist.id}`)"
                      >
                          <div class="name">{{ playlist.title }}</div>
                      </div>
                  </div>
              </div>
          </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useMessage, useDialog } from "naive-ui";
import { 
    Search, History, DeleteFour, Fire, 
    MusicOne, Voice, RecordDisc, Record 
} from "@icon-park/vue-next";
import { getSearchHot, getSearchSuggest } from "@/api/search";
import { useMusicDataStore } from "@/store/musicData";
import { useSettingDataStore } from "@/store/settingData";
import debounce from "@/utils/debounce";

const router = useRouter();
const route = useRoute();
const music = useMusicDataStore();
const setting = useSettingDataStore();
const message = useMessage();
const dialog = useDialog();

const inputRef = ref(null);
const inputValue = ref("");
const loading = ref(false);

const searchData = reactive({
    hot: [],
    suggest: {
        songs: [],
        artists: [],
        albums: [],
        playlists: []
    }
});

const isEmptyResult = computed(() => {
    const s = searchData.suggest;
    return !s.songs?.length && !s.artists?.length && !s.albums?.length && !s.playlists?.length;
});

// 初始化
onMounted(() => {
    // 聚焦输入框
    nextTick(() => {
        inputRef.value?.focus();
    });
    
    // 获取热搜
    getSearchHot().then(res => {
        if(res.data) searchData.hot = res.data;
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
            searchData.suggest = res || {};
            // Ensure arrays
            if(!searchData.suggest.songs) searchData.suggest.songs = [];
        }
    } catch (e) {
        loading.value = false;
    }
}

const handleEnter = () => {
    if(inputValue.value.trim()) {
        music.setSearchHistory(inputValue.value.trim());
        // 如果有单曲结果，直接播放第一首? 或者不做操作仅保存历史
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

const playSong = (song: any) => {
    // 兼容 pinia store 方法名: setPlaylists vs setPlayList
    if (typeof music.setPlayList === 'function') {
        music.setPlayList([song]);
    } else if (typeof music.setPlaylists === 'function') {
        music.setPlaylists([song]);
    }
    music.setPlaySongIndex(0);
    // 可选：记录搜索历史
    music.setSearchHistory(inputValue.value.trim());
    message.success(`开始播放: ${song.title}`);
}

</script>

<style lang="scss" scoped>
.search-page {
    padding: 20px;
    max-width: 1000px;
    margin: 0 auto;
    
    .search-header {
        margin-bottom: 30px;
        display: flex;
        justify-content: center;
        .main-input {
            width: 100%;
            max-width: 600px;
            font-size: 16px;
        }
    }

    .section {
        margin-bottom: 30px;
        .section-title {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 12px;
            .left {
                display: flex;
                align-items: center;
                gap: 6px;
                font-weight: bold;
                font-size: 16px;
            }
        }
    }
    
    .hot-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        gap: 16px;
        .hot-item {
            display: flex;
            align-items: center;
            cursor: pointer;
            padding: 8px;
            border-radius: 8px;
            &:hover {
                background-color: var(--n-color-hover);
            }
            .index {
                font-weight: bold;
                margin-right: 12px;
                width: 20px;
                text-align: center;
                &.top {
                    color: #d03050;
                }
            }
            .info {
                display: flex;
                flex-direction: column;
                .word {
                    font-weight: 500;
                }
                .desc {
                    font-size: 12px;
                    opacity: 0.6;
                }
            }
        }
    }

    .result-list {
        .result-section {
            margin-bottom: 24px;
            background-color: var(--n-card-color);
            border-radius: 12px;
            padding: 16px;
            
            .section-header {
                display: flex;
                align-items: center;
                gap: 8px;
                font-size: 16px;
                font-weight: bold;
                margin-bottom: 12px;
                color: var(--n-text-color-2);
                padding-bottom: 8px;
                border-bottom: 1px solid var(--n-border-color);
            }
            
            .list-items {
                .item {
                    padding: 10px 8px;
                    cursor: pointer;
                    border-radius: 6px;
                    display: flex;
                    align-items: center;
                    gap: 8px;
                    
                    &:hover {
                        background-color: var(--n-color-hover);
                    }
                    
                    &.song-item {
                         .artist, .album {
                             opacity: 0.6;
                             font-size: 13px;
                         }
                    }
                }
            }
        }
    }
    
    .loading, .empty {
        text-align: center;
        padding: 40px 0;
        opacity: 0.7;
    }
}
</style>
