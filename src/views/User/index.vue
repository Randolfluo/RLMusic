<template>
  <div class="user-center">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
    </div>

    <div class="profile-container">
      <div class="profile-header glass-card">
        <div class="avatar-section">
          <n-upload
            :show-file-list="false"
            :custom-request="handleAvatarUpload"
            accept="image/png,image/jpeg,image/gif,image/webp"
            class="avatar-uploader"
          >
              <div class="avatar-wrapper">
                  <n-avatar
                    round
                    :size="140"
                    :src="resolveAvatarUrl(user.getUserData.avatarUrl) || 'images/ico/user-filling.svg'"
                    fallback-src="images/ico/user-filling.svg"
                    class="user-avatar"
                    object-fit="cover"
                  />
                  <div class="avatar-overlay">
                      <n-icon size="40" :component="Camera" color="#ffffff" />
                      <span class="upload-text">更换头像</span>
                  </div>
              </div>
          </n-upload>
        </div>

        <div class="info-section">
          <div class="user-main">
            <h2 class="username">{{ user.getUserData.nickname || '用户' }}</h2>
            <n-tag size="small" :bordered="false" type="primary" round class="id-tag">
              ID: {{ user.getUserData.userId }}
            </n-tag>
          </div>
          <p v-if="user.getUserData.email" class="user-email">{{ user.getUserData.email }}</p>
          <div class="user-bio">
            生活不止眼前的苟且，还有诗和远方的田野。
          </div>
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-card glass-card-sm">
          <div class="stat-icon-wrapper time-icon">
            <n-icon :component="Time" size="24" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ formatDuration(userInfoDetails.total_duration) }}</div>
            <div class="stat-label">累计听歌时长</div>
          </div>
        </div>
        
        <div class="stat-card glass-card-sm">
          <div class="stat-icon-wrapper ip-icon">
            <n-icon :component="Connection" size="24" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ userInfoDetails.ip_src || '未知' }}</div>
            <div class="stat-label">当前 IP 地址</div>
          </div>
        </div>
        
        <div class="stat-card glass-card-sm">
           <div class="stat-icon-wrapper like-icon">
            <n-icon :component="Like" size="24" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ userInfoDetails.like_count || 0 }}</div>
            <div class="stat-label">喜欢的歌曲</div>
          </div>
        </div>
      </div>

      <div class="actions-section">
        <n-button 
          class="logout-btn" 
          type="error" 
          secondary 
          round 
          size="large" 
          @click="handleLogout"
        >
          <template #icon>
            <n-icon :component="Power" />
          </template>
          退出登录
        </n-button>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { userStore } from "@/store";
import { useRouter } from "vue-router";
import { useMessage, type UploadCustomRequestOptions } from "naive-ui";
import { 
    Camera, 
    Time,
    Connection,
    Like,
    Power
} from "@icon-park/vue-next";
import axios from "@/utils/request"; 
import { onMounted, ref } from "vue";
import { getUserInfo, uploadAvatar, resolveAvatarUrl } from "@/api/user";
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

<style scoped lang="scss">
.user-center {
  padding: 40px 20px;
  min-height: 90vh;
  position: relative;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 80px;

  .bg-decoration {
      position: fixed;
      top: 0;
      left: 0;
      width: 100vw;
      height: 100vh;
      z-index: -1;
      pointer-events: none;
      overflow: hidden;

      .blob {
          position: absolute;
          border-radius: 50%;
          filter: blur(80px);
          opacity: 0.4;
          animation: float 20s infinite ease-in-out;
      }

      .blob-1 {
          width: 600px;
          height: 600px;
          background: var(--n-primary-color);
          top: -200px;
          right: -100px;
          animation-delay: 0s;
      }

      .blob-2 {
          width: 500px;
          height: 500px;
          background: #4facfe;
          bottom: -100px;
          left: -150px;
          animation-delay: -5s;
      }
  }
}

.profile-container {
    width: 100%;
    max-width: 800px;
    display: flex;
    flex-direction: column;
    gap: 32px;
    position: relative;
    z-index: 10;
}

.glass-card {
    background: rgba(255, 255, 255, 0.65);
    backdrop-filter: blur(20px);
    border-radius: 32px;
    box-shadow: 0 10px 40px rgba(0,0,0,0.08);
    border: 1px solid rgba(255, 255, 255, 0.5);
    
    @media (prefers-color-scheme: dark) {
        background: rgba(30, 30, 30, 0.65);
        border: 1px solid rgba(255, 255, 255, 0.08);
        box-shadow: 0 10px 40px rgba(0,0,0,0.3);
    }
}

.profile-header {
    display: flex;
    align-items: center;
    padding: 40px;
    gap: 40px;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    
    &:hover {
        transform: translateY(-5px);
        box-shadow: 0 15px 50px rgba(0,0,0,0.12);
    }

    @media (max-width: 600px) {
        flex-direction: column;
        text-align: center;
        gap: 20px;
        padding: 30px 20px;
    }

    .avatar-section {
        flex-shrink: 0;
    }

    .info-section {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        
        .user-main {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 8px;
            
            @media (max-width: 600px) {
                justify-content: center;
                flex-wrap: wrap;
            }
            
            .username {
                margin: 0;
                font-size: 32px;
                font-weight: 800;
                color: var(--n-text-color);
                line-height: 1.2;
            }
            
            .id-tag {
                font-weight: 600;
            }
        }
        
        .user-email {
            color: var(--n-text-color-3);
            font-size: 15px;
            margin: 0 0 16px 0;
        }
        
        .user-bio {
            font-size: 16px;
            color: var(--n-text-color-2);
            line-height: 1.6;
            opacity: 0.8;
            font-style: italic;
        }
    }
}

.avatar-wrapper {
    position: relative;
    cursor: pointer;
    border-radius: 50%;
    overflow: hidden;
    box-shadow: 0 8px 24px rgba(0,0,0,0.15);
    transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    border: 4px solid rgba(255, 255, 255, 0.8);
    
    @media (prefers-color-scheme: dark) {
        border: 4px solid rgba(255, 255, 255, 0.1);
    }
    
    &:hover {
        transform: scale(1.05) rotate(2deg);
        box-shadow: 0 12px 32px rgba(var(--n-primary-color-rgb), 0.3);
    }

    .user-avatar {
        display: block;
        background-color: var(--n-card-color);
    }

    .avatar-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        opacity: 0;
        transition: all 0.3s ease;
        backdrop-filter: blur(4px);
        gap: 8px;
        
        .upload-text {
            color: white;
            font-size: 12px;
            font-weight: 600;
        }
    }

    &:hover .avatar-overlay {
        opacity: 1;
    }
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
    
    @media (max-width: 600px) {
        grid-template-columns: 1fr;
    }
    
    .stat-card {
        padding: 24px;
        border-radius: 24px;
        display: flex;
        align-items: center;
        gap: 20px;
        transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
        cursor: default;
        
        &.glass-card-sm {
            background: rgba(255, 255, 255, 0.5);
            backdrop-filter: blur(12px);
            border: 1px solid rgba(255, 255, 255, 0.3);
            box-shadow: 0 4px 15px rgba(0,0,0,0.03);
            
            @media (prefers-color-scheme: dark) {
                background: rgba(255, 255, 255, 0.05);
                border: 1px solid rgba(255, 255, 255, 0.05);
            }
        }
        
        &:hover {
            transform: translateY(-5px);
            background: rgba(255, 255, 255, 0.8);
            box-shadow: 0 10px 25px rgba(0,0,0,0.08);
            
            @media (prefers-color-scheme: dark) {
                background: rgba(255, 255, 255, 0.1);
            }
        }
        
        .stat-icon-wrapper {
            width: 56px;
            height: 56px;
            border-radius: 16px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 24px;
            box-shadow: 0 8px 16px rgba(0,0,0,0.15);
            
            &.time-icon {
                background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 99%, #fecfef 100%);
                color: #d03050;
            }
            
            &.ip-icon {
                background: linear-gradient(120deg, #89f7fe 0%, #66a6ff 100%);
                color: #0052d4;
            }
            
            &.like-icon {
                background: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
                color: #555;
            }
        }
        
        .stat-content {
            .stat-value {
                font-size: 20px;
                font-weight: 800;
                color: var(--n-text-color);
                line-height: 1.2;
                margin-bottom: 4px;
                font-family: 'DM Mono', monospace;
            }
            
            .stat-label {
                font-size: 13px;
                color: var(--n-text-color-3);
                font-weight: 500;
            }
        }
    }
}

.actions-section {
    display: flex;
    justify-content: center;
    margin-top: 20px;
    
    .logout-btn {
        padding: 0 40px;
        height: 50px;
        font-weight: 600;
        font-size: 16px;
        box-shadow: 0 4px 15px rgba(208, 48, 80, 0.3);
        transition: all 0.3s;
        
        &:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 25px rgba(208, 48, 80, 0.4);
        }
    }
}

@keyframes float {
    0% { transform: translate(0, 0) rotate(0deg); }
    33% { transform: translate(30px, -50px) rotate(10deg); }
    66% { transform: translate(-20px, 20px) rotate(-5deg); }
    100% { transform: translate(0, 0) rotate(0deg); }
}
</style>
