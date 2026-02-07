<template>
  <div class="home">

    <div class="section-title">
      <h2>公共歌单</h2>
      <n-button text @click="router.push('/playlists')" style="font-size: 14px">
        更多
        <template #icon>
          <n-icon :component="Right" />
        </template>
      </n-button>
    </div>

    <PlaylistGrid 
      :loading="publicLoading" 
      :playlists="publicPlaylists" 
      empty-text="暂无公共歌单" 
      collapsed 
      :collapsed-rows="2" 
    />

    <div class="section-title" v-if="userStore.userLogin">
      <h2>私有歌单</h2>
      <n-button text @click="router.push('/private-playlists')" style="font-size: 14px">
        更多
        <template #icon>
          <n-icon :component="Right" />
        </template>
      </n-button>
    </div>

    <PlaylistGrid 
      v-if="userStore.userLogin"
      :loading="privateLoading" 
      :playlists="privatePlaylists" 
      empty-text="暂无私有歌单" 
      collapsed 
      :collapsed-rows="2" 
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getPublicPlaylists, getUserPrivatePlaylists } from "@/api/playlist";
import SystemStats from "@/components/Home/SystemStats.vue";
import PlaylistGrid from "@/components/DataList/PlaylistGrid.vue";
import { ResultCode } from "@/utils/request";
import { Right } from "@icon-park/vue-next";
import { useMessage } from "naive-ui";
import { useRouter } from "vue-router";
import { useUserDataStore } from "@/store/userData";

const router = useRouter();
const message = useMessage();
const userStore = useUserDataStore();

const publicLoading = ref(false);
const publicPlaylists = ref<any[]>([]);

const privateLoading = ref(false);
const privatePlaylists = ref<any[]>([]);

onMounted(() => {
  getPublicList();
  if (userStore.userLogin) {
    getPrivateList();
  }
});

const getPublicList = async () => {
  publicLoading.value = true;
  try {
    const res = await getPublicPlaylists();
    if (res.code === ResultCode.SUCCESS) {
      publicPlaylists.value = res.data || [];
    }
  } catch (error) {
    message.error("获取公共歌单失败");
  } finally {
    publicLoading.value = false;
  }
};

const getPrivateList = async () => {
  privateLoading.value = true;
  try {
    const res = await getUserPrivatePlaylists();
    if (res.code === ResultCode.SUCCESS) {
      privatePlaylists.value = res.data || [];
    }
  } catch (error) {
    // 可能是未登录或权限问题，这里简单处理
    console.error(error);
  } finally {
    privateLoading.value = false;
  }
};
</script>

<style scoped lang="scss">
.home {
  padding: 24px;
  .section-title {
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
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
