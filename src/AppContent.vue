<template>
  <!-- 全屏页面（如登录页） -->
  <div v-if="route.meta.hideLayout" style="width: 100vw; height: 100vh;">
    <router-view />
  </div>

  <!-- 普通页面布局 -->
  <n-layout class="main-layout" v-else> 
      <!-- Naive UI的Layout组件，创建基础的页面布局框架 -->
    <n-layout-header class="nav-header">    
      <Nav />
    </n-layout-header>
    <!-- 智能布局容器 -->
    <n-layout-content
      class="layout-content"
      position="absolute"
      :class="music.getPlaylists[0] && music.showPlayBar ? 'show-player' : ''"
      :native-scrollbar="false"
      content-style="min-height: 100%; display: flex; flex-direction: column;"
    >
      <!-- 主内容包装器 -->
      <main ref="mainContent" class="main-container">
        <!-- 智能返回顶部按钮 -->
        <n-back-top
          :bottom="music.getPlaylists[0] && music.showPlayBar ? 120 : 40"
          :right="40"
          class="custom-back-top"
        >
          <n-icon :component="ToTop" size="24" />
        </n-back-top>  
        <!-- Vue Router的占位符，根据当前URL显示对应的页面 -->
        <router-view v-slot="{ Component }">  
          <keep-alive>
            <Transition name="fade-slide" mode="out-in">
              <component :is="Component" v-if="Component"/>
            </Transition>
          </keep-alive>
        </router-view>
      </main>
    </n-layout-content>
  </n-layout>
  
  <!-- 仅在非全屏模式下显示播放器 -->
  <template v-if="!route.meta.hideLayout">
    <Player />
    <BigPlayer v-if="music.getPlaylists[0]" />
  </template>
  <n-modal v-model:show="showServerConfig" class="server-modal" :mask-closable="false">
    <div class="server-card" :class="{ 'is-connected': serverConnectionState === 'connected', 'is-checking': serverConnectionState === 'checking' }">
      <div class="server-ambient"></div>
      <div class="server-grid"></div>
      <div class="server-glow"></div>
      <div class="server-shell">
        <div class="server-aside">
          <div class="server-badge">{{ serverBadgeText }}</div>
          <div class="server-icon">
            <n-icon :component="Connection" size="30" />
          </div>
          <div class="server-aside-title">{{ serverAsideTitle }}</div>
          <div class="server-aside-subtitle">{{ serverAsideSubtitle }}</div>
          <div class="server-pulses">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
        <div class="server-main">
          <div class="server-header">
            <div class="server-title">
              <span>服务器状态</span>
            </div>
            <n-tag :type="serverConnectionState === 'connected' ? 'success' : serverConnectionState === 'checking' ? 'info' : 'warning'" size="small" round>
              {{ serverConnectionState === 'connected' ? '已连接' : serverConnectionState === 'checking' ? '检测中' : '需要配置' }}
            </n-tag>
          </div>
          <div class="server-status">
            <span class="server-dot"></span>
            <span>{{ serverStatusText }}</span>
          </div>
          <div class="server-desc">
            {{ serverStatusDesc }}
          </div>
          <div class="server-meta" v-if="resolvedServerOrigin">
            <div class="server-meta-label">当前地址</div>
            <div class="server-meta-value">{{ resolvedServerOrigin }}</div>
          </div>
          <div class="server-input">
            <div class="server-label">服务器地址</div>
            <n-input-group>
              <n-input
                v-model:value="serverUrlInput"
                placeholder="http://127.0.0.1:12345"
                size="large"
                clearable
                autofocus
                class="server-input-inner"
                @keyup.enter="applyServerUrl"
              />
              <n-button v-if="isCapacitor" type="primary" size="large" secondary @click="scanQrCode" :loading="scanning">
                <template #icon>
                  <n-icon :component="QrCodeScannerOutlined" />
                </template>
              </n-button>
            </n-input-group>
          </div>
          <div class="server-actions">
            <n-button secondary size="large" @click="applyServerUrl">保存并重连</n-button>
            <n-button quaternary size="large" @click="dismissServerConfig">稍后设置</n-button>
          </div>
          <div class="server-hint">
            需保证后端服务已启动且端口可达
          </div>
        </div>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { musicStore, userStore } from "@/store";
import { login } from "@/api/login"; // 引入登录api
import { aesDecrypt } from "@/utils/encrypt"; // 引入解密
import Nav from "@/components/Nav/index.vue";
import Player from "@/components/Player/index.vue";
import BigPlayer from "@/components/Player/BigPlayer.vue";
import packageJson from "@/../package.json";
import { ref, onMounted, onUnmounted, computed, watch } from 'vue';
import { ResultCode, apiBaseURL } from "@/utils/request";
import { useRoute, useRouter } from "vue-router"; // 引入 useRoute
import { NIcon, NBackTop, NModal, NInput, NButton, NTag, NInputGroup, NLayout, NLayoutHeader, NLayoutContent } from "naive-ui";
import { BarcodeScanner, SupportedFormat } from "@capacitor-community/barcode-scanner";
import { Capacitor } from "@capacitor/core";
import { App as CapacitorApp } from "@capacitor/app";
import { QrCodeScannerOutlined } from "@vicons/material";
import { ToTop, Connection } from "@icon-park/vue-next";
import { useMessage } from "naive-ui";

const message = useMessage();
const route = useRoute(); // 获取当前路由信息
const router = useRouter();
const music = musicStore();
const user = userStore();
// const router = useRouter();
// const setting = settingStore();
const showServerConfig = ref(false);
const isCapacitor = typeof window !== "undefined" && Capacitor.isNativePlatform();
const scanning = ref(false);
const serverUrlInput = ref(localStorage.getItem("server_url") || "");
const serverConnectionState = ref<"checking" | "connected" | "disconnected">("checking");
const lastCheckedOrigin = ref("");
const normalizeServerUrl = (value: string) => {
  const trimmed = value.trim();
  if (!trimmed) return "";
  return trimmed.endsWith("/") ? trimmed.slice(0, -1) : trimmed;
};

const parseApiFromQr = (raw: string) => {
  const value = raw.trim();
  if (!value) return "";
  try {
    const directUrl = new URL(value);
    const api = directUrl.searchParams.get("api");
    if (api) return normalizeServerUrl(decodeURIComponent(api));
    return normalizeServerUrl(value);
  } catch {}

  const apiMatch = value.match(/(?:\?|&|#)api=([^&#]+)/i);
  if (apiMatch?.[1]) {
    return normalizeServerUrl(decodeURIComponent(apiMatch[1]));
  }

  try {
    const obj = JSON.parse(value);
    if (obj?.api) return normalizeServerUrl(String(obj.api));
  } catch {}

  return normalizeServerUrl(value);
};

const buildScanErrorMessage = (error: any) => {
  const raw = String(error?.message || error || "未知错误");
  if (/User cancelled|canceled|cancelled/i.test(raw)) {
    return "你已取消扫码。";
  }
  if (/permission|denied/i.test(raw)) {
    return "相机权限被拒绝，请在系统设置中开启相机权限后重试。";
  }
  if (/timeout|超时/i.test(raw)) {
    return "扫码超时，请将二维码放入取景框中央后重试。";
  }
  return `扫码失败：${raw}`;
};

const startScannerOverlay = async () => {
  await BarcodeScanner.hideBackground();
  document.body.classList.add("scanner-active");
  document.documentElement.classList.add("scanner-active");
};

const stopScannerOverlay = async () => {
  try {
    await BarcodeScanner.showBackground();
  } catch {}
  try {
    await BarcodeScanner.stopScan();
  } catch {}
  document.body.classList.remove("scanner-active");
  document.documentElement.classList.remove("scanner-active");
};

onMounted(() => {
  if (isCapacitor) {
    CapacitorApp.addListener('backButton', async () => {
      // 1. QR Scanner
      if (scanning.value) {
        await stopScannerOverlay();
        scanning.value = false;
        return;
      }
      
      // 2. Server Config Modal
      if (showServerConfig.value) {
        showServerConfig.value = false;
        return;
      }
      
      // 3. Big Player
      if (music.showBigPlayer) {
        music.setBigPlayerState(false);
        return;
      }
      
      // 4. Router Back or Exit
      if (router.currentRoute.value.path !== '/' && router.currentRoute.value.path !== '/login') {
        router.back();
      } else {
        CapacitorApp.exitApp();
      }
    });
  }
});

onUnmounted(() => {
  if (isCapacitor) {
    CapacitorApp.removeAllListeners();
  }
});

const scanQrCode = async () => {
  if (!isCapacitor) {
    message.warning("仅移动端 App 支持扫码");
    return;
  }
  scanning.value = true;
  try {
    const permission = await BarcodeScanner.checkPermission({ force: true });
    if (permission.denied) {
      await BarcodeScanner.openAppSettings();
      message.error("相机权限已被永久拒绝，请在系统设置中手动开启");
      return;
    }
    if (!permission.granted) {
      message.error("未获得相机权限");
      return;
    }

    await BarcodeScanner.prepare({
      targetedFormats: [SupportedFormat.QR_CODE],
    });
    
    // 暂时关闭弹窗，避免遮挡
    const originalShowConfig = showServerConfig.value;
    showServerConfig.value = false;
    
    // 确保之前任何正在进行的扫描都停止
    await stopScannerOverlay();

    await startScannerOverlay();
    const result = await BarcodeScanner.startScan({
      targetedFormats: [SupportedFormat.QR_CODE],
    });
    
    // 恢复弹窗
    showServerConfig.value = originalShowConfig;
    
    const content = result.hasContent ? String(result.content || "") : "";
    if (!content) {
      message.warning("未识别到二维码内容");
      return;
    }
    const apiValue = parseApiFromQr(content);
    if (!apiValue) {
      message.error("二维码中未识别到 api 地址");
      return;
    }
    
    serverUrlInput.value = apiValue;
    message.success("已识别并填充服务器地址");
    
    // 自动尝试连接
    applyServerUrl();
  } catch (err: any) {
    const prettyMessage = buildScanErrorMessage(err);
    message.error(prettyMessage);
    showServerConfig.value = true; // 确保弹窗恢复
  } finally {
    await stopScannerOverlay();
    scanning.value = false;
  }
};

const applyServerUrl = () => {
  const normalized = normalizeServerUrl(serverUrlInput.value);
  if (!normalized) {
    localStorage.removeItem("server_url");
    location.reload();
    return;
  }
  localStorage.setItem("server_url", normalized);
  location.reload();
};

const dismissServerConfig = () => {
  showServerConfig.value = false;
};

const normalizeOrigin = (value: string) => {
  if (!value) return "";
  const trimmed = value.trim();
  if (!trimmed) return "";
  return trimmed.endsWith("/api") ? trimmed.replace(/\/api$/, "") : trimmed;
};

const resolvedServerOrigin = computed(() => {
  const fromInput = normalizeOrigin(serverUrlInput.value);
  if (fromInput) return fromInput;
  const fromApi = normalizeOrigin(apiBaseURL);
  if (fromApi) return fromApi;
  if (typeof window !== "undefined") return window.location.origin;
  return "";
});

const serverStatusText = computed(() => {
  if (serverConnectionState.value === "checking") return "检测中";
  if (serverConnectionState.value === "connected") return "服务器已连接";
  return "无法访问服务器";
});

const serverStatusDesc = computed(() => {
  if (serverConnectionState.value === "connected") return "连接正常，可继续使用";
  if (serverConnectionState.value === "checking") return "正在校验后端连通性";
  return "请更新后端地址或检查服务状态";
});

const serverAsideTitle = computed(() => {
  if (serverConnectionState.value === "connected") return "服务器已连接";
  if (serverConnectionState.value === "checking") return "正在检测";
  return "服务器未连接";
});

const serverAsideSubtitle = computed(() => {
  if (serverConnectionState.value === "connected") return "当前连接稳定，可继续使用";
  if (serverConnectionState.value === "checking") return "正在尝试与后端建立连接";
  return "请输入可访问的后端地址";
});

const serverBadgeText = computed(() => {
  if (serverConnectionState.value === "connected") return "连接正常";
  if (serverConnectionState.value === "checking") return "连接检测";
  return "连接诊断";
});

const checkServerConnection = async () => {
  if (typeof window === "undefined") return;
  const origin = resolvedServerOrigin.value;
  lastCheckedOrigin.value = origin;
  if (!origin || !origin.startsWith("http")) {
    serverConnectionState.value = "disconnected";
    return;
  }
  serverConnectionState.value = "checking";
  const controller = new AbortController();
  const timer = window.setTimeout(() => controller.abort(), 3000);
  try {
    const res = await fetch(`${origin}/api/system/stats`, {
      credentials: "include",
      signal: controller.signal,
    });
    serverConnectionState.value = res.ok ? "connected" : "disconnected";
  } catch {
    serverConnectionState.value = "disconnected";
  } finally {
    window.clearTimeout(timer);
  }
};

let serverCheckTimer: number | undefined;
const scheduleServerCheck = () => {
  if (!showServerConfig.value) return;
  if (serverCheckTimer) window.clearTimeout(serverCheckTimer);
  serverCheckTimer = window.setTimeout(() => {
    checkServerConnection();
  }, 400);
};

// 自动登录逻辑
onMounted(() => {
    const savedUser = localStorage.getItem("remember_user");
    const savedPass = localStorage.getItem("remember_pass");
    if (savedUser && savedPass && !user.userLogin) {
        try {
            const password = aesDecrypt(savedPass);
            login({ username: savedUser, password: password }).then(res => {
                if (res.code === ResultCode.SUCCESS) {
                    const userData = {
                        userId: res.data.id,
                        nickname: res.data.username,
                        email: res.data.email,
                        userGroup: res.data.user_group,
                        avatarUrl: res.data.avatar,
                    };
                    user.setUserData(userData);
                    localStorage.setItem("token", res.data.token); 
                    console.log("Auto login success");
                }
            }).catch(e => {
                console.error("Auto login failed", e);
            });
        } catch (e) {
            console.error("Auto login decrypt failed", e);
        }
    }
});

// 监听空格键控制播放暂停
onMounted(() => {
  window.addEventListener("keydown", (e) => {
    if (e.code === "Space") {
      // 如果当前焦点在输入框或文本域中，不触发
      if (
        document.activeElement?.tagName === "INPUT" ||
        document.activeElement?.tagName === "TEXTAREA"
      ) {
        return;
      }
      e.preventDefault();
      music.setPlayState(!music.getPlayState);
    }
  });
});

const handleServerConnectFail = (event: Event) => {
  const detail = (event as CustomEvent).detail;
  if (detail?.baseURL && !serverUrlInput.value) {
    const base = String(detail.baseURL).replace(/\/api$/, "");
    if (base && base.startsWith("http")) {
      serverUrlInput.value = base;
    }
  }
  serverConnectionState.value = "disconnected";
  showServerConfig.value = true;
};

onMounted(() => {
  window.addEventListener("server-connection-failed", handleServerConnectFail);
});

onUnmounted(() => {
  window.removeEventListener("server-connection-failed", handleServerConnectFail);
});

const openServerConfig = () => {
  showServerConfig.value = true;
};

onMounted(() => {
  window.addEventListener("open-server-config", openServerConfig);
});

onUnmounted(() => {
  window.removeEventListener("open-server-config", openServerConfig);
});

watch(showServerConfig, (value) => {
  if (value) {
    checkServerConnection();
  }
});

watch(serverUrlInput, () => {
  scheduleServerCheck();
});

const mainContent = ref<HTMLElement | null>(null);

// 公告数据
const annShow =
  import.meta.env.VITE_ANN_TITLE && import.meta.env.VITE_ANN_CONTENT
    ? true
    : false;
const annTitle = import.meta.env.VITE_ANN_TITLE;
const annContene = import.meta.env.VITE_ANN_CONTENT;
const annDuration = Number(import.meta.env.VITE_ANN_DURATION);

// 系统重置
const cleanAll = () => {
  window.$message ? window.$message.success("重置成功") : alert("重置成功");
  localStorage.clear();
  window.location.href = "/";
};


onMounted(() => {
  // 将应用内部功能暴露到全局window对象的技术，允许在浏览器控制台直接调用这些功能
  (window as any).$mainContent = mainContent.value;  // 将Vue组件的DOM引用暴露到全局window对象
  (window as any).$cleanAll = cleanAll;   // 暴露重置函数

  // 公告，如果有则显示
  if (annShow) {
    if(typeof $notification !== 'undefined') {
       $notification["info"]({
        content: annTitle,
        meta: annContene,
        duration: annDuration,
      });
    }
  }
  
  // 版权声明
  const logoText = `${packageJson.name}`;
  const copyrightNotice = `\n\n版本: ${packageJson.version}\n作者: ${packageJson.author}\n作者主页: ${packageJson.home}\n`;
  console.info(
    `%c${logoText} %c ${copyrightNotice}`,
    "color:#f55e55;font-size:26px;font-weight:bold;",
    "font-size:16px"
  );
  console.info(
    "若站点出现异常，可尝试在下方输入 %c$cleanAll()%c 然后按回车来重置",
    "background: #eaeffd;color:#f55e55;padding: 4px 6px;border-radius:8px;",
    "background:unset;color:unset;"
  );
});

</script>

<style lang="scss" scoped>
.main-layout {
  height: 100vh;
  /* 柔和的浅色渐变背景 */
  background: linear-gradient(135deg, #fdfbfb 0%, #ebedee 100%);
  transition: background 0.3s ease;

  :global(.dark) & {
    /* 深色模式下的深邃背景 */
    background: linear-gradient(135deg, #1a1a1a 0%, #0d0d0d 100%);
  }
}

.nav-header {
  height: 60px;
  background: transparent !important;
  z-index: 100;
  /* 移除边框，由 Nav 组件处理 */
  border: none !important;
}

.layout-content {
  top: 60px;
  bottom: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: transparent !important;

  &.show-player {
    bottom: 90px; /* 留出播放器空间 */
  }

  :deep(.n-scrollbar-rail--vertical) {
    right: 4px;
    width: 6px;
  }
}

.main-container {
  width: 100%;
  max-width: 1440px; /* 宽屏适配 */
  margin: 0 auto;
  padding: 0; /* Nav已包含padding，这里内容由各页面自行控制padding，或者统一加 */
  /* 
     由于各页面（如 Home, Playlist）已有 padding，这里设为 0 以免双重 padding。
     但考虑到一致性，如果大部分页面都需要 padding，可以在这里加。
     查看之前的页面代码，大多都有 padding: 24px。
     所以这里设为 0，让页面自己控制。
  */
  min-height: 100%;
  box-sizing: border-box;
}

/* 优化返回顶部按钮 */
:deep(.n-back-top) {
  background-color: var(--n-color-primary) !important;
  color: #fff !important;
  box-shadow: 0 4px 16px rgba(var(--n-color-primary-rgb), 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 50%;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:hover {
    transform: translateY(-4px) scale(1.05);
    box-shadow: 0 8px 24px rgba(var(--n-color-primary-rgb), 0.5);
  }
}

/* 路由跳转动画 - 渐变滑动 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.server-modal {
  display: flex;
  align-items: center;
  justify-content: center;
}

.server-modal :deep(.n-modal-body) {
  padding: 0;
}

.server-card {
  width: min(720px, 94vw);
  padding: 22px;
  border-radius: 26px;
  background: radial-gradient(120% 120% at 0% 0%, rgba(255, 244, 229, 0.95) 0%, rgba(255, 255, 255, 0.98) 55%, rgba(255, 255, 255, 0.94) 100%);
  border: 1px solid rgba(255, 178, 102, 0.35);
  box-shadow: 0 28px 90px rgba(15, 23, 42, 0.22);
  backdrop-filter: blur(16px);
  position: relative;
  overflow: hidden;
}

:global(.dark) .server-card {
  background: radial-gradient(120% 120% at 0% 0%, rgba(42, 32, 20, 0.92) 0%, rgba(22, 22, 24, 0.96) 55%, rgba(18, 18, 20, 0.92) 100%);
  border: 1px solid rgba(255, 178, 102, 0.18);
  box-shadow: 0 28px 90px rgba(0, 0, 0, 0.5);
}

.server-card.is-connected {
  border: 1px solid rgba(80, 200, 140, 0.4);
  background: radial-gradient(120% 120% at 0% 0%, rgba(226, 255, 238, 0.95) 0%, rgba(255, 255, 255, 0.98) 55%, rgba(248, 255, 252, 0.92) 100%);
}

:global(.dark) .server-card.is-connected {
  border: 1px solid rgba(80, 200, 140, 0.25);
  background: radial-gradient(120% 120% at 0% 0%, rgba(28, 40, 32, 0.92) 0%, rgba(18, 20, 20, 0.96) 55%, rgba(18, 22, 20, 0.92) 100%);
}

.server-shell {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 18px;
  position: relative;
  z-index: 1;
}

.server-aside {
  padding: 18px 16px;
  border-radius: 18px;
  background: linear-gradient(170deg, rgba(255, 255, 255, 0.75), rgba(255, 255, 255, 0.4));
  border: 1px solid rgba(255, 179, 102, 0.25);
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: flex-start;
}

:global(.dark) .server-aside {
  background: linear-gradient(170deg, rgba(35, 30, 26, 0.9), rgba(20, 20, 22, 0.6));
  border: 1px solid rgba(255, 179, 102, 0.18);
}

.server-badge {
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: rgba(255, 140, 60, 0.95);
  font-weight: 700;
}

.server-card.is-connected .server-badge {
  color: rgba(24, 168, 96, 0.95);
}

.server-icon {
  width: 50px;
  height: 50px;
  border-radius: 16px;
  background: rgba(255, 165, 90, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ff914d;
  box-shadow: inset 0 0 0 1px rgba(255, 165, 90, 0.2);
}

.server-card.is-connected .server-icon {
  background: rgba(72, 201, 132, 0.18);
  color: #1f9d63;
  box-shadow: inset 0 0 0 1px rgba(72, 201, 132, 0.3);
}

.server-aside-title {
  font-size: 18px;
  font-weight: 800;
}

.server-aside-subtitle {
  font-size: 12px;
  line-height: 1.6;
  opacity: 0.7;
}

.server-pulses {
  display: flex;
  gap: 6px;
  margin-top: auto;
}

.server-pulses span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 145, 77, 0.45);
  animation: server-pulse 1.8s infinite;
}

.server-card.is-connected .server-pulses span {
  background: rgba(72, 201, 132, 0.55);
}

.server-pulses span:nth-child(2) {
  animation-delay: 0.3s;
}

.server-pulses span:nth-child(3) {
  animation-delay: 0.6s;
}

.server-main {
  padding: 8px 6px 6px 2px;
}

.server-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 12px;
  position: relative;
  z-index: 1;
}

.server-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0.3px;
}

.server-status {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 1.4px;
  font-weight: 700;
  margin-bottom: 10px;
  position: relative;
  z-index: 1;
}

.server-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #ff914d;
  box-shadow: 0 0 12px rgba(255, 145, 77, 0.7);
}

.server-card.is-connected .server-dot {
  background: #22c55e;
  box-shadow: 0 0 14px rgba(34, 197, 94, 0.6);
}

.server-card.is-checking .server-dot {
  background: #38bdf8;
  box-shadow: 0 0 14px rgba(56, 189, 248, 0.6);
}

.server-desc {
  font-size: 14px;
  opacity: 0.7;
  line-height: 1.6;
  position: relative;
  z-index: 1;
}

.server-meta {
  margin-top: 14px;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(255, 168, 76, 0.2);
  display: grid;
  gap: 6px;
  position: relative;
  z-index: 1;
}

.server-card.is-connected .server-meta {
  border: 1px solid rgba(72, 201, 132, 0.28);
}

:global(.dark) .server-meta {
  background: rgba(18, 18, 20, 0.75);
  border: 1px solid rgba(255, 168, 76, 0.16);
}

.server-meta-label {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  opacity: 0.6;
}

.server-meta-value {
  font-size: 13px;
  font-weight: 600;
  word-break: break-all;
}

.server-input {
  margin-top: 18px;
  position: relative;
  z-index: 1;
}

.server-label {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  opacity: 0.6;
  margin-bottom: 8px;
}

.server-input-inner {
  --n-border-radius: 16px;
  --n-height: 48px;
}

.server-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
  position: relative;
  z-index: 1;
  flex-wrap: wrap;
}

.server-hint {
  margin-top: 16px;
  font-size: 12px;
  opacity: 0.6;
  position: relative;
  z-index: 1;
}

.server-ambient {
  position: absolute;
  inset: -40% -10% auto auto;
  width: 280px;
  height: 280px;
  background: radial-gradient(circle, rgba(255, 168, 76, 0.5), transparent 70%);
  filter: blur(0px);
  opacity: 0.7;
  pointer-events: none;
}

.server-grid {
  position: absolute;
  inset: 0;
  background-image:
    repeating-linear-gradient(90deg, rgba(255, 160, 90, 0.08), rgba(255, 160, 90, 0.08) 1px, transparent 1px, transparent 18px),
    repeating-linear-gradient(0deg, rgba(20, 20, 20, 0.06), rgba(20, 20, 20, 0.06) 1px, transparent 1px, transparent 18px);
  opacity: 0.55;
  mix-blend-mode: multiply;
  pointer-events: none;
}

.server-glow {
  position: absolute;
  inset: auto -20% -25% -20%;
  height: 200px;
  background: radial-gradient(circle, rgba(255, 150, 72, 0.25), transparent 70%);
  filter: blur(40px);
  opacity: 0.9;
  pointer-events: none;
}

.server-card :deep(.n-input) {
  background: rgba(255, 255, 255, 0.85);
  border-radius: 16px;
  transition: all 0.2s ease;
}

:global(.dark) .server-card :deep(.n-input) {
  background: rgba(20, 20, 22, 0.85);
}

.server-card :deep(.n-input--focus) {
  box-shadow: 0 0 0 2px rgba(255, 145, 77, 0.25);
}

.server-card :deep(.n-button--secondary) {
  background: rgba(255, 164, 88, 0.18);
  border-color: rgba(255, 164, 88, 0.35);
  color: #c25115;
}

:global(.dark) .server-card :deep(.n-button--secondary) {
  color: #ffb27a;
}

.server-card :deep(.n-button--secondary:hover) {
  transform: translateY(-1px);
}

.server-card :deep(.n-button--quaternary:hover) {
  background: rgba(255, 164, 88, 0.12);
}

@keyframes server-pulse {
  0% {
    transform: scale(0.7);
    opacity: 0.4;
  }
  60% {
    transform: scale(1);
    opacity: 0.95;
  }
  100% {
    transform: scale(0.7);
    opacity: 0.4;
  }
}

.server-card :deep(.n-input__input-el) {
  -webkit-user-select: text;
  user-select: text;
}

.server-card :deep(.n-input__textarea-el) {
  -webkit-user-select: text;
  user-select: text;
}

@media (max-width: 720px) {
  .server-card {
    width: min(520px, 94vw);
    padding: 18px;
  }

  .server-shell {
    grid-template-columns: 1fr;
  }

  .server-aside {
    align-items: center;
    text-align: center;
  }

  .server-main {
    padding: 0;
  }
}
</style>
