<template>
  <div class="home">
    <div class="page-header">
      <h1 class="welcome-text">
        <span class="greeting">Hello,</span>
        <span class="highlight">{{ userStore.userLogin ? userStore.userData?.nickname : 'Music Lover' }}</span>
      </h1>
      <p class="subtitle">Enjoy your local music journey</p>
    </div>

    <div class="banner-grid">
      <!-- 局域网访问二维码 -->
      <div class="banner-card access-card" :style="bannerStyle" @click="nextIp">
        <div class="card-bg-deco"></div>
        <div class="card-content">
          <div class="text-info">
            <h3>局域网访问</h3>
            <p v-if="urls.length > 0" class="ip-text">
              <span class="url-value">{{ currentUrl }}</span>
              <n-icon v-if="urls.length > 1" :component="Refresh" class="refresh-icon" />
            </p>
            <p v-else class="ip-text">正在获取 IP...</p>
          </div>
          <div class="qr-code-wrapper" v-if="urls.length > 0" @click.stop="openQrModal">
            <qrcode-vue :value="currentShareUrl" :size="80" level="H" class="qrcode" />
          </div>
        </div>
      </div>

      <!-- 二维码放大弹窗 -->
      <n-modal v-model:show="showQrModal">
        <div class="qr-modal-container">
          <div class="modal-header">
            <h3>局域网访问</h3>
            <p v-if="urls.length > 1" class="subtitle">
              点击二维码切换网络 ({{ currentIpIndex + 1 }}/{{ urls.length }})
            </p>
          </div>

          <div
            class="modal-qr-wrapper"
            :class="{ clickable: urls.length > 1 }"
            @click="nextIp"
          >
            <qrcode-vue :value="currentShareUrl" :size="260" level="H" class="modal-qrcode" />
            <div class="switch-hint" v-if="urls.length > 1">
              <n-icon :component="Refresh" size="32" color="white" />
            </div>
          </div>

          <div class="modal-footer">
            <div class="url-pill">
              <span>{{ currentShareUrl }}</span>
            </div>
            <n-button circle secondary type="primary" @click="copyUrl" class="copy-btn">
              <template #icon>
                <n-icon :component="Copy" />
              </template>
            </n-button>
          </div>
        </div>
      </n-modal>

      <!-- 一起听歌入口 -->
      <!-- <div class="banner-card listen-card" @click="router.push('/listen-together')">
        <div class="card-bg-deco-listen"></div>
        <div class="card-content">
          <div class="text-info">
            <h3>一起听歌</h3>
            <p>与好友同步听歌</p>
          </div>
          <div class="icon-wrapper">
            <n-icon :component="MusicOne" size="48" />
          </div>
        </div>
      </div> -->
    </div>

    <div class="content-section">
      <div class="section-header">
        <div class="title-group">
          <h2>公共歌单</h2>
          <span class="badge">Public</span>
        </div>
        <n-button text class="more-btn" @click="router.push('/playlists')">
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
        cols="2 s:3 m:4 l:6 xl:6"
        @refresh="getPublicList"
        class="fade-in-section"
      />
    </div>

    <div class="content-section" v-if="userStore.userLogin">
      <div class="section-header">
        <div class="title-group">
          <h2>私有歌单</h2>
          <span class="badge private">Private</span>
        </div>
        <n-button text class="more-btn" @click="router.push('/private-playlists')">
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
        cols="2 s:3 m:4 l:6 xl:6"
        @refresh="getPrivateList"
        class="fade-in-section delay-1"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getPublicPlaylists, getUserPrivatePlaylists } from "@/api/playlist";
import { getLocalIPs } from "@/api/system";
import QrcodeVue from 'qrcode.vue';
import PlaylistGrid from "@/components/DataList/PlaylistGrid.vue";
import { ResultCode } from "@/utils/request";
import { Right, Refresh, Copy } from "@icon-park/vue-next";
import { useMessage, NIcon, NModal, NButton } from "naive-ui";
import { useRouter } from "vue-router";
import { useUserDataStore } from "@/store/userData";
import { useThemeVars } from "naive-ui";
import { computed } from "vue";

const router = useRouter();
const message = useMessage();
const userStore = useUserDataStore();
const themeVars = useThemeVars();

const urls = ref<string[]>([]);
const currentIpIndex = ref(0);
const showQrModal = ref(false);

const currentUrl = computed(() => {
  if (urls.value.length === 0) return '';
  return urls.value[currentIpIndex.value];
});

const normalizeServerUrl = (value: string) => {
  const trimmed = value.trim();
  if (!trimmed) return "";
  return trimmed.endsWith("/") ? trimmed.slice(0, -1) : trimmed;
};

const getBackendUrl = () => {
  const stored = localStorage.getItem("server_url");
  const normalizedStored = stored ? normalizeServerUrl(stored) : "";
  const isLocalHost = normalizedStored.includes("localhost") || normalizedStored.includes("127.0.0.1");
  if (normalizedStored && !isLocalHost) return normalizedStored;
  const backendPortRaw = localStorage.getItem("backend_port") || "";
  const backendPort = Number(backendPortRaw) > 0 ? Number(backendPortRaw) : 12345;
  if (currentUrl.value) {
    try {
      const url = new URL(currentUrl.value);
      return `${url.protocol}//${url.hostname}:${backendPort}`;
    } catch {
      return `${window.location.protocol}//${window.location.hostname}:${backendPort}`;
    }
  }
  return `${window.location.protocol}//${window.location.hostname}:${backendPort}`;
};

const currentShareUrl = computed(() => {
  if (!currentUrl.value) return "";
  const backendUrl = getBackendUrl();
  if (!backendUrl) return currentUrl.value;
  try {
    const url = new URL(currentUrl.value);
    url.searchParams.set("api", backendUrl);
    return url.toString();
  } catch {
    const encoded = encodeURIComponent(backendUrl);
    return currentUrl.value.includes("?")
      ? `${currentUrl.value}&api=${encoded}`
      : `${currentUrl.value}?api=${encoded}`;
  }
});

const openQrModal = () => {
  showQrModal.value = true;
};

const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(currentShareUrl.value || '');
    message.success("链接已复制");
  } catch (err) {
    message.error("复制失败");
  }
};

const nextIp = () => {
  if (urls.value.length > 1) {
    currentIpIndex.value = (currentIpIndex.value + 1) % urls.value.length;
    message.success("已切换 IP 地址");
  }
};

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
  fetchLocalIPs();
});

const fetchLocalIPs = async () => {
  try {
    const isElectron = typeof navigator !== "undefined" && navigator.userAgent.includes("Electron");
    const accessIp = localStorage.getItem("access_ip") || "";
    const storedPort = localStorage.getItem("frontend_port") || "";
    const port = window.location.port || (isElectron ? (storedPort || "23456") : "80");
    if (isElectron && accessIp) {
      urls.value = [`http://${accessIp}:${port}/`];
      return;
    }
    const res = await getLocalIPs(port);
    if (res.code === ResultCode.SUCCESS) {
      urls.value = res.data.urls || [];
    }
  } catch (e) {
    console.error(e);
  }
};

const getPublicList = async () => {
  publicLoading.value = true;
  try {
    const res = await getPublicPlaylists(1, 12);
    if (res.code === ResultCode.SUCCESS) {
        if (Array.isArray(res.data)) {
            publicPlaylists.value = res.data;
        } else {
            publicPlaylists.value = res.data.list;
        }
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
    const res = await getUserPrivatePlaylists(1, 12);
    if (res.code === ResultCode.SUCCESS) {
        if (Array.isArray(res.data)) {
            privatePlaylists.value = res.data;
        } else {
            privatePlaylists.value = res.data.list;
        }
    }
  } catch (error) {
    console.error(error);
  } finally {
    privateLoading.value = false;
  }
};
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@600;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

.home {
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;
  min-height: 100vh;
  position: relative;
  overflow: hidden;
  font-family: 'Plus Jakarta Sans', sans-serif;

  /* 温暖米色调背景 */
  background: #faf8f5;

  /* 动态渐变装饰 */
  &::before {
    content: '';
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: -1;
    pointer-events: none;
    background:
      radial-gradient(ellipse 80% 50% at 20% -10%, rgba(224, 122, 95, 0.15), transparent),
      radial-gradient(ellipse 80% 50% at 80% 120%, rgba(61, 139, 139, 0.12), transparent),
      radial-gradient(ellipse 60% 40% at 90% 20%, rgba(212, 165, 116, 0.1), transparent);
    animation: bg-shift 20s ease-in-out infinite alternate;
  }

  @keyframes bg-shift {
    0% { transform: translateX(0) scale(1); }
    100% { transform: translateX(-20px) scale(1.05); }
  }

  @media (max-width: 768px) {
    padding: 16px;
  }
}

/* Page Header */
.page-header {
  margin-bottom: 40px;
  animation: fade-in-down 0.8s cubic-bezier(0.2, 0.8, 0.2, 1);

  @media (max-width: 768px) {
    margin-bottom: 24px;

    .welcome-text {
      font-size: 32px;
    }
  }

  .welcome-text {
    font-size: 42px;
    font-weight: 700;
    line-height: 1.2;
    margin: 0 0 8px 0;
    letter-spacing: -1px;
    color: #1a1a1a;
    font-family: 'Plus Jakarta Sans', sans-serif;

    .greeting {
      opacity: 0.5;
      font-weight: 400;
      margin-right: 12px;
    }

    .highlight {
      background: linear-gradient(120deg, #e07a5f 0%, #d4a574 100%);
      -webkit-background-clip: text;
      background-clip: text;
      -webkit-text-fill-color: transparent;
      display: inline-block;
    }
  }

  .subtitle {
    font-size: 16px;
    color: #666666;
    margin: 0;
    font-weight: 500;
    opacity: 0.8;
  }
}

/* Banner Grid */
.banner-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  margin-bottom: 48px;
  animation: fade-in-up 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.1s backwards;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 16px;
    margin-bottom: 32px;
  }
}

.banner-card {
  border-radius: 24px;
  padding: 24px 32px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  color: white;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  position: relative;
  overflow: hidden;
  height: 140px;
  border: 1px solid rgba(255, 255, 255, 0.15);

  &:hover {
    transform: translateY(-6px) scale(1.01);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);

    .card-bg-deco, .card-bg-deco-listen {
      transform: scale(1.2) rotate(10deg);
      opacity: 0.8;
    }
  }

  /* Decorative Backgrounds */
  .card-bg-deco {
    position: absolute;
    top: -50%;
    right: -20%;
    width: 300px;
    height: 300px;
    background: radial-gradient(circle, rgba(255,255,255,0.2) 0%, rgba(255,255,255,0) 70%);
    border-radius: 50%;
    opacity: 0.4;
    transition: all 0.6s ease;
    pointer-events: none;
  }

  .card-bg-deco-listen {
    position: absolute;
    bottom: -40%;
    left: -10%;
    width: 250px;
    height: 250px;
    background: radial-gradient(circle, rgba(255,255,255,0.25) 0%, rgba(255,255,255,0) 60%);
    border-radius: 50%;
    opacity: 0.3;
    transition: all 0.6s ease;
    pointer-events: none;
  }

  .card-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    z-index: 1;
    position: relative;
  }

  @media (max-width: 768px) {
    padding: 16px 20px;
    height: auto;
    min-height: 120px;

    .text-info {
      h3 {
        font-size: 20px;
      }
      p {
        font-size: 13px;
      }
    }
  }

  .text-info {
    flex: 1;
    min-width: 0;
    h3 {
      font-size: 24px;
      margin: 0 0 8px 0;
      font-weight: 700;
      letter-spacing: -0.5px;
      text-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    p {
      margin: 0;
      opacity: 0.95;
      font-size: 14px;
      line-height: 1.5;
      font-weight: 500;
    }
    .ip-text {
      font-family: 'DM Mono', monospace;
      display: inline-flex;
      align-items: center;
      gap: 8px;
      background: rgba(0, 0, 0, 0.2);
      padding: 6px 12px;
      border-radius: 8px;
      width: fit-content;
      max-width: 100%;
      font-weight: 500;
      backdrop-filter: blur(4px);
      border: 1px solid rgba(255,255,255,0.1);

      @media (max-width: 768px) {
        font-size: 12px;
        padding: 4px 8px;
      }

      .url-value {
        letter-spacing: 0.5px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
    .refresh-icon {
      font-size: 14px;
      opacity: 0.8;
      cursor: pointer;
      transition: transform 0.3s;
      &:hover {
        transform: rotate(180deg);
        opacity: 1;
      }
    }
  }

  .icon-wrapper {
    opacity: 0.9;
    filter: drop-shadow(0 4px 8px rgba(0,0,0,0.2));
    background: rgba(255, 255, 255, 0.15);
    width: 64px;
    height: 64px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(255,255,255,0.2);
    backdrop-filter: blur(4px);
  }
}

.access-card {
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(20px) saturate(180%);
  background: linear-gradient(135deg, #2c3e50 0%, #3d8b8b 100%) !important;

  .qr-code-wrapper {
    background: rgba(255, 255, 255, 0.95);
    padding: 10px;
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
    margin-left: 24px;
    cursor: zoom-in;
    flex-shrink: 0;
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    border: 1px solid rgba(0, 0, 0, 0.06);

    @media (max-width: 768px) {
      margin-left: 12px;
      padding: 8px;
      border-radius: 12px;

      .qrcode {
        width: 56px !important;
        height: 56px !important;
      }
    }

    &:hover {
      transform: scale(1.12) rotate(3deg);
      box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
    }

    .qrcode {
      display: block;
      border-radius: 8px;
    }
  }
}

.listen-card {
  background: linear-gradient(120deg, #e07a5f 0%, #d4a574 100%);
  position: relative;

  &::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url("data:image/svg+xml,%3Csvg width='20' height='20' viewBox='0 0 20 20' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='%23ffffff' fill-opacity='0.05' fill-rule='evenodd'%3E%3Ccircle cx='3' cy='3' r='3'/%3E%3Ccircle cx='13' cy='13' r='3'/%3E%3C/g%3E%3C/svg%3E");
    opacity: 0.3;
  }
}

/* Content Sections */
.content-section {
  margin-bottom: 48px;

  &.fade-in-section {
    animation: fade-in-up 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.2s backwards;

    &.delay-1 {
      animation-delay: 0.3s;
    }
  }
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0 4px;

  .title-group {
    display: flex;
    align-items: center;
    gap: 12px;

    h2 {
      font-size: 24px;
      font-weight: 700;
      margin: 0;
      letter-spacing: -0.5px;
      color: #1a1a1a;
      font-family: 'Plus Jakarta Sans', sans-serif;
    }

    .badge {
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 1px;
      padding: 4px 8px;
      border-radius: 6px;
      background: rgba(61, 139, 139, 0.1);
      color: #3d8b8b;

      &.private {
        background: rgba(224, 122, 95, 0.1);
        color: #e07a5f;
      }
    }
  }

  .more-btn {
    font-size: 14px;
    font-weight: 600;
    color: #666666;
    padding: 6px 12px;
    border-radius: 20px;
    transition: all 0.2s;

    &:hover {
      background: #f5f2ed;
      color: #e07a5f;
    }
  }
}

/* QR Modal Styles */
.qr-modal-container {
  background: #faf8f5;
  border: 1px solid #ebe7e0;
  border-radius: 24px;
  padding: 32px;
  width: 360px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.14);
  animation: modal-pop 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  position: relative;
  overflow: hidden;

  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 6px;
    background: linear-gradient(90deg, #e07a5f 0%, #d4a574 100%);
  }

  .modal-header {
    text-align: center;
    margin-bottom: 24px;

    h3 {
      font-size: 20px;
      font-weight: 700;
      margin: 0;
      color: #1a1a1a;
      font-family: 'Plus Jakarta Sans', sans-serif;
    }

    .subtitle {
      margin: 6px 0 0 0;
      font-size: 13px;
      color: #666666;
    }
  }

  .modal-qr-wrapper {
    background: white;
    padding: 12px;
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    position: relative;
    transition: all 0.3s ease;

    &.clickable {
      cursor: pointer;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 30px rgba(224, 122, 95, 0.2);

        .switch-hint {
          opacity: 1;
        }
      }

      &:active {
        transform: scale(0.98);
      }
    }

    .modal-qrcode {
      display: block;
      border-radius: 8px;
    }

    .switch-hint {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: rgba(224, 122, 95, 0.8);
      width: 64px;
      height: 64px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: opacity 0.3s ease;
      backdrop-filter: blur(4px);
      box-shadow: 0 4px 12px rgba(224, 122, 95, 0.3);
    }
  }

  .modal-footer {
    margin-top: 24px;
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;

    .url-pill {
      flex: 1;
      background: #f5f2ed;
      padding: 10px 16px;
      border-radius: 12px;
      font-family: 'DM Mono', monospace;
      font-size: 13px;
      color: #1a1a1a;
      text-align: center;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      border: 1px solid #ebe7e0;
      transition: all 0.2s;

      &:hover {
        border-color: #d4a574;
        background: white;
      }
    }

    .copy-btn {
      transition: all 0.2s;

      &:hover {
        transform: rotate(15deg) scale(1.1);
      }
    }
  }
}

/* Animations */
@keyframes fade-in-down {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes modal-pop {
  0% { opacity: 0; transform: scale(0.95); }
  100% { opacity: 1; transform: scale(1); }
}
</style>
