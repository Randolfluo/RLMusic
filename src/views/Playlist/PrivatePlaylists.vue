<template>
  <div class="private-playlists">
    <div class="section-title">
      <h2>私有歌单</h2>
    </div>

    <PlaylistGrid :loading="loading" :playlists="playlists" empty-text="暂无私有歌单" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getUserPrivatePlaylists } from "@/api/playlist";
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
    const res = await getUserPrivatePlaylists();
    if (res.code === ResultCode.SUCCESS) {
      playlists.value = res.data;
    }
  } catch (error) {
    message.error("获取歌单失败");
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped lang="scss">
.private-playlists {
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
