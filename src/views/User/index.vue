<template>
  <div class="user-center">
    <n-card title="个人中心" hoverable class="profile-card">
      <div class="user-info">
        <n-upload
          :show-file-list="false"
          :custom-request="handleAvatarUpload"
          accept="image/png,image/jpeg,image/gif,image/webp"
          class="avatar-uploader"
        >
            <div class="avatar-wrapper">
                <n-avatar
                round
                :size="120"
                :src="user.getUserData.avatarUrl || '/images/ico/user-filling.svg'"
                fallback-src="/images/ico/user-filling.svg"
                class="user-avatar"
                />
                <div class="avatar-overlay">
                    <n-icon size="32" :component="Camera" color="#ffffff" />
                </div>
            </div>
        </n-upload>

        <h2 class="username">{{ user.getUserData.nickname || '用户' }}</h2>
        <div class="user-meta">
            <n-tag size="small" :bordered="false" type="primary" round>ID: {{ user.getUserData.userId }}</n-tag>
        </div>
        <p v-if="user.getUserData.email" class="user-email">{{ user.getUserData.email }}</p>
      </div>

      <n-divider />

      <n-grid x-gap="24" y-gap="24" cols="2 600:4" style="margin-bottom: 24px;">
        <n-gi>
          <n-card size="small" :bordered="false" class="stat-card">
              <n-statistic label="累计歌曲">
                <template #prefix>
                    <n-icon :component="Music" color="#18a058" />
                </template>
                {{ userInfoDetails.total_songs || 0 }}
                <template #suffix>首</template>
              </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small" :bordered="false" class="stat-card">
            <n-statistic label="累计专辑">
                <template #prefix>
                    <n-icon :component="RecordDisc" color="#2080f0" />
                </template>
                {{ userInfoDetails.total_albums || 0 }}
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small" :bordered="false" class="stat-card">
             <n-statistic label="累计歌手">
                <template #prefix>
                    <n-icon :component="Voice" color="#f0a020" />
                </template>
                {{ userInfoDetails.total_artists || 0 }}
            </n-statistic>
          </n-card>
        </n-gi>
         <n-gi>
           <n-card size="small" :bordered="false" class="stat-card">
              <n-statistic label="累计时长">
                <template #prefix>
                    <n-icon :component="Time" color="#d03050" />
                </template>
                {{ formatDuration(userInfoDetails.total_duration) }}
              </n-statistic>
           </n-card>
        </n-gi>
      </n-grid>
      
      <div class="favorite-section">
         <n-card size="small" embedded :bordered="false" class="fav-card">
             <n-space align="center" justify="center" vertical>
                 <n-text depth="3">最喜爱的歌曲</n-text>
                 <n-space align="center">
                    <n-icon :component="Like" color="#d03050" size="24" />
                    <span class="fav-song-name">{{ userInfoDetails.favorite_song || '暂无数据' }}</span>
                 </n-space>
             </n-space>
         </n-card>
      </div>
      
      <n-divider />
      
      <div class="actions">
        <n-space justify="center" size="large">
          <n-button strong secondary type="primary" size="large" @click="$router.push('/user/like')">
              我的喜欢/上传
          </n-button>
          <n-button strong secondary type="error" size="large" @click="handleLogout">
              退出登录
          </n-button>
        </n-space>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { userStore } from "@/store";
import { useRouter } from "vue-router";
import { useMessage, type UploadCustomRequestOptions } from "naive-ui";
import { 
    Camera, 
    Music, 
    RecordDisc, 
    Voice, 
    Time,
    Like
} from "@icon-park/vue-next";
import axios from "@/utils/request"; 
import { onMounted, ref } from "vue";
import { getUserInfo, uploadAvatar } from "@/api/user";
import { ResultCode } from "@/utils/request";

const user = userStore();
const router = useRouter();
const message = useMessage();
const userInfoDetails = ref<any>({});

const formatDuration = (seconds: number) => {
    if (!seconds) return "0分钟";
    if (seconds < 60) return `${seconds}秒`;
    const min = Math.floor(seconds / 60);
    if (min < 60) return `${min}分钟`;
    return `${(min/60).toFixed(1)}小时`;
}

const handleAvatarUpload = ({ file }: UploadCustomRequestOptions) => {
    const formData = new FormData();
    formData.append("file", file.file as File);
    uploadAvatar(formData).then(res => {
         if(res.code === ResultCode.SUCCESS) {
             message.success("头像上传成功");
             user.userData.avatarUrl = res.data;
         } else {
             message.error(res.message || "上传失败");
         }
    }).catch(() => {
        message.error("上传出错");
    });
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
  max-width: 900px;
  margin: 0 auto;
}

.profile-card {
    border-radius: 16px;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.avatar-uploader {
    display: flex;
    justify-content: center;
}

.avatar-wrapper {
    position: relative;
    cursor: pointer;
    border-radius: 50%;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    transition: transform 0.3s;
}

.avatar-wrapper:hover {
    transform: scale(1.05);
}

.user-avatar {
    display: block;
}

.avatar-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.4);
    display: flex;
    justify-content: center;
    align-items: center;
    opacity: 0;
    transition: opacity 0.3s;
    backdrop-filter: blur(2px);
}

.avatar-wrapper:hover .avatar-overlay {
    opacity: 1;
}

.username {
  margin-top: 20px;
  margin-bottom: 8px;
  font-size: 28px;
  font-weight: 700;
  color: var(--n-text-color);
}

.user-meta {
    margin-bottom: 8px;
}

.user-email {
  color: var(--n-text-color-3);
  font-size: 14px;
}

.stat-card {
    background: var(--n-action-color);
    text-align: center;
    border-radius: 12px;
    transition: all 0.3s;
}

.stat-card:hover {
    background: var(--n-color-target);
    transform: translateY(-2px);
}

.favorite-section {
    text-align: center; 
    margin-bottom: 24px;
    padding: 0 20px;
}

.fav-card {
    background-color: rgba(240, 160, 32, 0.05); 
    border-radius: 12px;
}

.fav-song-name {
    font-size: 18px;
    font-weight: 600;
    color: #d03050; 
}
</style>
