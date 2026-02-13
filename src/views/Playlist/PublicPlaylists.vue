<template>
  <div class="public-playlists">
    <div class="section-title">
      <h2>公共歌单</h2>
    </div>

    <PlaylistGrid :loading="loading" :playlists="playlists" empty-text="暂无公共歌单" />
    
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
import { getPublicPlaylists } from "@/api/playlist";
import { ResultCode } from "@/utils/request";
import { useMessage } from "naive-ui";
import PlaylistGrid from "@/components/DataList/PlaylistGrid.vue";
import Pagination from "@/components/Pagination/index.vue";

const message = useMessage();
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
    const res = await getPublicPlaylists(page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      // 兼容旧接口直接返回数组的情况，如果后端返回 { list, total } 则取 list
      if (Array.isArray(res.data)) {
         playlists.value = res.data;
         total.value = res.data.length; // 这种情况通常没有总数，或者是全部数据
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
.public-playlists {
  padding: 24px;
  .section-title {
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    h2 {
      font-size: 24px;
      font-weight: bold;
      margin: 0;
      position: relative;
      padding-left: 10px;
      &::before {
        content: "";
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
        width: 4px;
        height: 18px;
        background-color: var(--n-color-primary);
        border-radius: 2px;
      }
    }
  }
}
</style>
