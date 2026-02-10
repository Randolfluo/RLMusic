<template>
  <div class="home">
    <div class="banner-section" :style="bannerStyle" @click="router.push('/listen-together')">
      <div class="banner-content">
        <h3>一起听歌</h3>
        <p>与好友实时同步听歌，分享音乐的快乐</p>
      </div>
      <div class="banner-icon">
        <n-icon :component="MusicOne" size="48" />
      </div>
    </div>

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
import { Right, MusicOne } from "@icon-park/vue-next";
import { useMessage, NIcon } from "naive-ui";
import { useRouter } from "vue-router";
import { useUserDataStore } from "@/store/userData";
import { useThemeVars } from "naive-ui";
import { computed } from "vue";

const router = useRouter();
const message = useMessage();
const userStore = useUserDataStore();
const themeVars = useThemeVars();

const bannerStyle = computed(() => {
  return {
    background: `linear-gradient(135deg, ${themeVars.value.primaryColor} 0%, ${themeVars.value.infoColor} 100%)`
  };
});

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
}

.banner-section {
  /* background: linear-gradient(135deg, var(--n-color-primary) 0%, color-mix(in srgb, var(--n-color-primary), black 10%) 100%); */
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  color: white;
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
  }

  .banner-content {
    h3 {
      font-size: 24px;
      margin: 0 0 8px 0;
      font-weight: bold;
    }
    p {
      margin: 0;
      opacity: 0.9;
      font-size: 14px;
    }
  }

  .banner-icon {
    opacity: 0.8;
  }
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  h2 {
    font-size: 20px;
    font-weight: bold;
    margin: 0;
    position: relative;
    padding-left: 12px;
    
    &::before {
      content: "";
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 4px;
      height: 16px;
      background-color: var(--n-color-primary);
      border-radius: 2px;
    }
  }
}
</style>
