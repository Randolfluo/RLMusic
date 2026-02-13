<template>
  <div class="private-playlists">
    <div class="header">
      <div class="cover">
        <n-image
          class="cover-img"
          :src="userStore.userData.avatarUrl || '/images/logo/favicon.png'"
          fallback-src="/images/logo/favicon.png"
          object-fit="cover"
          preview-disabled
        />
      </div>
      <div class="info">
        <div class="tag">私有歌单</div>
        <div class="title">{{ userStore.userData.nickname }} 的歌单</div>
        <div class="creator">
          <n-avatar 
            round 
            size="small" 
            :src="userStore.userData.avatarUrl || '/images/logo/favicon.png'" 
            style="margin-right: 8px; vertical-align: middle;"
          />
          <span style="vertical-align: middle;">{{ userStore.userData.nickname }}</span>
        </div>
        <div class="desc">
            共 {{ total }} 个歌单
        </div>
      </div>
    </div>

    <n-divider />

    <PlaylistGrid :loading="loading" :playlists="playlists" empty-text="暂无私有歌单" @refresh="getPlaylists" />

    <div class="pagination-container" style="display: flex; justify-content: center; margin-top: 20px;">
      <Pagination
        v-if="playlists.length > 0"
        :totalCount="total"
        :pageNumber="page"
        :showSizePicker="true"
        @pageNumberChange="onPageChange"
        @pageSizeChange="onPageSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getUserPrivatePlaylists } from "@/api/playlist";
import { ResultCode } from "@/utils/request";
import { useMessage, NImage, NDivider, NAvatar } from "naive-ui";
import PlaylistGrid from "@/components/DataList/PlaylistGrid.vue";
import Pagination from "@/components/Pagination/index.vue";
import { useUserDataStore } from "@/store/userData";

const message = useMessage();
const userStore = useUserDataStore();
const loading = ref(false);
const playlists = ref<any[]>([]);
const page = ref(1);
const limit = ref(20);
const total = ref(0);

onMounted(() => {
  getPlaylists();
});

const getPlaylists = async () => {
  loading.value = true;
  try {
    const res = await getUserPrivatePlaylists(page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      if (Array.isArray(res.data)) {
          playlists.value = res.data;
          total.value = res.data.length;
      } else {
          playlists.value = res.data.list;
          total.value = res.data.total;
      }
    }
  } catch (error) {
    message.error("获取歌单失败");
  } finally {
    loading.value = false;
  }
};

const onPageChange = (val: number) => {
  page.value = val;
  getPlaylists();
};

const onPageSizeChange = (val: number) => {
  limit.value = val;
  page.value = 1;
  getPlaylists();
};
</script>

<style scoped lang="scss">
.private-playlists {
  padding: 24px;
  
  .header {
    display: flex;
    margin-bottom: 24px;
    
    .cover {
      width: 200px;
      height: 200px;
      border-radius: 8px;
      overflow: hidden;
      margin-right: 24px;
      flex-shrink: 0;
      box-shadow: 0 4px 12px rgba(0,0,0,0.1);
      
      .cover-img {
        width: 100%;
        height: 100%;
      }
    }
    
    .info {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      
      .tag {
        display: inline-block;
        width: fit-content;
        padding: 4px 12px;
        border: 1px solid var(--n-color-primary);
        color: var(--n-color-primary);
        border-radius: 4px;
        font-size: 14px;
        margin-bottom: 12px;
      }
      
      .title {
        font-size: 32px;
        font-weight: bold;
        margin-bottom: 12px;
        line-height: 1.2;
      }
      
      .creator {
        display: flex;
        align-items: center;
        margin-bottom: 16px;
        font-size: 14px;
        opacity: 0.8;
      }
      
      .desc {
        font-size: 14px;
        opacity: 0.6;
        line-height: 1.6;
        margin-bottom: 16px;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }
  }
}
</style>
