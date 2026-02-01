<template>
  <div class="user-center">
    <n-card title="个人中心">
      <div class="user-info">
        <n-avatar
          round
          :size="100"
          :src="user.getUserData.avatarUrl || '/images/ico/user-filling.svg'"
          fallback-src="/images/ico/user-filling.svg"
        />
        <h2 class="username">{{ user.getUserData.nickname || '用户' }}</h2>
        <p class="user-id">ID: {{ user.getUserData.userId }}</p>
        <p v-if="user.getUserData.email" class="user-email">{{ user.getUserData.email }}</p>
      </div>

      <n-divider />

      <n-grid x-gap="12" :cols="4" style="margin-bottom: 24px; text-align: center;">
        <n-gi>
          <n-statistic label="累计听歌">
            {{ userInfoDetails.totalSongs || 0 }} 首
          </n-statistic>
        </n-gi>
        <n-gi>
          <n-statistic label="累计专辑">
             {{ userInfoDetails.totalAlbums || 0 }}
          </n-statistic>
        </n-gi>
        <n-gi>
          <n-statistic label="累计歌手">
             {{ userInfoDetails.totalArtists || 0 }}
          </n-statistic>
        </n-gi>
         <n-gi>
          <n-statistic label="累计时长">
             {{ formatDuration(userInfoDetails.totalDuration) }}
          </n-statistic>
        </n-gi>
      </n-grid>
      
      <div style="text-align: center; margin-bottom: 24px;">
         <n-statistic label="最喜爱的歌曲">
             {{ userInfoDetails.favoriteSong || '暂无' }}
         </n-statistic>
      </div>
      
      <n-divider />
      
      <div class="actions">
        <n-space justify="center">
          <n-button @click="$router.push('/user/like')">我的喜欢/上传</n-button>
          <n-button type="error" @click="handleLogout">退出登录</n-button>
        </n-space>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { userStore } from "@/store";
import { useRouter } from "vue-router";
import { useMessage } from "naive-ui";
import axios from "@/utils/request"; 
import { onMounted, ref } from "vue";
import { getUserInfo } from "@/api/user";
import { ResultCode } from "@/utils/request";

const user = userStore();
const router = useRouter();
const message = useMessage();
const userInfoDetails = ref<any>({});

const formatDuration = (seconds: number) => {
    if (!seconds) return "0 分钟";
    if (seconds < 60) return `${seconds} 秒`;
    const min = Math.floor(seconds / 60);
    if (min < 60) return `${min} 分钟`;
    return `${(min/60).toFixed(1)} 小时`;
}

onMounted(() => {
    getUserInfo().then(res => {
        if(res.code === ResultCode.SUCCESS) {
            userInfoDetails.value = res.data;
            // update avatar if backend has new one
            if (res.data.avatar) {
                 user.userData.avatarUrl = res.data.avatar;
            }
        }
    })
})

const handleLogout = () => {
  axios({
      method: "POST",
      url: "/auth/logout"
  }).then(() => {
      message.success("已退出登录");
  }).catch(() => {
      // Even if server fails, we clean up locally
      message.success("已退出登录");
  }).finally(() => {
      user.userLogOut();
      localStorage.removeItem("token");
      router.push("/login");
  });
};
</script>

<style scoped>
.user-center {
  padding: 40px 20px;
  max-width: 800px;
  margin: 0 auto;
}
.user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}
.username {
  margin-top: 16px;
  font-size: 24px;
}
.user-id {
  color: #666;
  margin-top: 4px;
}
.user-email {
  color: #888;
  font-size: 14px;
}
</style>
