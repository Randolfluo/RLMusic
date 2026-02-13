<template>
  <div class="subscribed-playlists">
    <div class="section-title">
      <h2>收藏歌单</h2>
    </div>

    <PlaylistGrid :loading="loading" :playlists="playlists" empty-text="暂无收藏歌单" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getSubscribedPlaylists } from "@/api/playlist";
import { ResultCode } from "@/utils/request";
import { useMessage } from "naive-ui";
import PlaylistGrid from "@/components/DataList/PlaylistGrid.vue";

const message = useMessage();
const loading = ref(false);
const playlists = ref<any[]>([]);

onMounted(() => {
  getPlaylists();
});

const getPlaylists = async () => {
  loading.value = true;
  try {
    const res = await getSubscribedPlaylists();
    if (res.code === ResultCode.SUCCESS) {
      if (Array.isArray(res.data)) {
          playlists.value = res.data;
      } else if (res.data && Array.isArray(res.data.list)) {
          playlists.value = res.data.list;
      } else {
          playlists.value = [];
      }
    }
  } catch (error) {
    message.error("获取歌单失败");
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped lang="scss">
.subscribed-playlists {
  padding: 24px;
  .section-title {
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    h2 {
      font-size: 24px;
      font-weight: bold;
    }
  }
}
</style>
