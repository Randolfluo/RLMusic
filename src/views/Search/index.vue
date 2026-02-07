<template>
  <div class="search-page">
    <div class="search-header">
      <n-card :bordered="false" class="header-card">
        <n-input-group>
            <!-- 移除选择框，直接展示输入框 -->
            <n-input
            ref="inputRef"
            v-model:value="inputValue"
            class="main-input"
            size="large"
            round
            placeholder="搜索音乐、歌手、歌单、专辑"
            clearable
            @input="handleInput"
            @keydown.enter="handleEnter"
            >
            <template #prefix>
                <n-icon :component="Search" />
            </template>
            </n-input>
            <n-button size="large" type="primary" ghost @click="handleEnter">
                搜索
            </n-button>
        </n-input-group>
      </n-card>
    </div>

    <div class="search-content">
      <!-- 默认显示：历史记录 -->
      <div v-if="!inputValue" class="default-content">
        <n-card class="section history" title="搜索历史" :bordered="false" size="small" v-if="music.getSearchHistory[0] && setting.searchHistory">
           <template #header-extra>
               <n-button text size="small" @click="delHistory">
                   <template #icon>
                       <n-icon :component="DeleteFour" />
                   </template>
                   清空
               </n-button>
           </template>
           <n-space>
              <n-tag
                v-for="item in music.getSearchHistory"
                :key="item"
                round
                checkable
                type="primary"
                @click="clickHistory(item)"
              >
                {{ item }}
              </n-tag>
           </n-space>
        </n-card>
      </div>

      <!-- 搜索结果区域 -->
      <div v-else class="suggest-content">
          <div class="search-title-bar">
              <h2>{{ inputValue }} 的相关搜索</h2>
          </div>
          
          <n-tabs type="line" animated v-model:value="searchType">
            <!-- 单曲 Tab -->
            <n-tab-pane name="songs" tab="单曲">
                 <SongList :songs="searchData.suggest.songs" :loading="loading" />
            </n-tab-pane>

            <!-- 歌单 Tab -->
            <n-tab-pane name="playlists" tab="歌单">
                <div class="loading" v-if="loading"><n-spin /></div>
                <n-empty v-else-if="!searchData.suggest.playlists?.length" description="暂无歌单" />
                <n-grid v-else cols="2 600:3 900:4" x-gap="20" y-gap="20">
                     <n-gi v-for="playlist in searchData.suggest.playlists" :key="playlist.id">
                         <n-card hoverable class="playlist-card" @click="router.push(`/playlist/${playlist.id}`)">
                             <template #cover>
                                 <!-- placeholder cover -->
                                 <div style="height: 120px; background: #f0f0f0; display: flex; align-items: center; justify-content: center;">
                                     <n-icon :component="Record" size="40" color="#ccc" />
                                 </div>
                             </template>
                             <n-thing :title="playlist.title">
                                 <template #description>{{ playlist.description || '暂无描述' }}</template>
                             </n-thing>
                         </n-card>
                     </n-gi>
                </n-grid>
            </n-tab-pane>

            <!-- 歌手 Tab -->
             <n-tab-pane name="artists" tab="歌手">
                <div class="loading" v-if="loading"><n-spin /></div>
                <n-empty v-else-if="!searchData.suggest.artists?.length" description="暂无歌手" />
                 <n-grid v-else cols="2 600:4 900:6" x-gap="20" y-gap="20">
                     <n-gi v-for="artist in searchData.suggest.artists" :key="artist.id">
                         <n-card hoverable @click="router.push(`/artist?id=${artist.id}`)" content-style="text-align: center;">
                             <n-avatar round :size="80" :src="artist.pic_url || undefined" fallback-src="/images/logo/logo.png">
                                 <n-icon :component="Voice" />
                             </n-avatar>
                             <div style="margin-top: 10px; font-weight: bold;">{{ artist.name }}</div>
                         </n-card>
                     </n-gi>
                </n-grid>
            </n-tab-pane>

             <!-- 专辑 Tab -->
             <n-tab-pane name="albums" tab="专辑">
                 <div class="loading" v-if="loading"><n-spin /></div>
                 <n-empty v-else-if="!searchData.suggest.albums?.length" description="暂无专辑" />
                 <n-list v-else hoverable clickable bordered>
                      <n-list-item 
                        v-for="album in searchData.suggest.albums" 
                        :key="album.id"
                        @click="router.push(`/album?id=${album.id}`)"
                      >
                         <template #prefix>
                             <n-icon :component="RecordDisc" size="40" color="#888" />
                         </template>
                         <n-thing :title="album.title">
                           <template #description>
                               {{ album.Artist?.name }}
                           </template>
                         </n-thing>
                      </n-list-item>
                 </n-list>
            </n-tab-pane>

          </n-tabs>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useMessage, useDialog } from "naive-ui";
import SongList from "@/components/DataList/SongList.vue";
import { 
    Search, DeleteFour,
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

</script>

<style lang="scss" scoped>
.search-page {
    padding: 30px;
    max-width: 1200px;
    margin: 0 auto;
    
    .search-header {
        margin-bottom: 30px;
        display: flex;
        justify-content: center;
        
        .header-card {
            width: 100%;
            max-width: 800px;
            background-color: transparent !important; // Blend in
        }

        .main-input {
            flex: 1;
        }
    }

    .search-content {
        .default-content, .suggest-content {
            animation: fadeIn 0.3s ease;
        }
        
        .search-title-bar {
            margin-bottom: 20px;
            h2 {
                font-size: 22px;
                font-weight: 600;
                color: var(--n-text-color-1);
            }
        }

        .result-section {
            margin-bottom: 20px;
        }

        .section.history {
           margin-bottom: 20px;
        }
    }
    
    .loading, .empty {
        text-align: center;
        padding: 60px 0;
        opacity: 0.8;
    }
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>
