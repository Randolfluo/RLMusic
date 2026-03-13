<template>
  <div class="init-page">
    <div class="glass-bg"></div>

    <n-button 
      quaternary 
      circle 
      class="back-btn" 
      @click="router.push('/')"
    >
      <template #icon>
        <n-icon :component="ArrowBackOutlined" size="24" />
      </template>
    </n-button>

    <div class="init-container">
      <div class="header">
        <div class="logo-wrapper">
          <img src="/images/logo/favicon.png" alt="Logo" class="logo" />
        </div>
        <div class="title-group">
          <h1 class="app-title">RLMusic</h1>
          <p class="subtitle">您的私有云音乐库</p>
        </div>
      </div>

      <n-alert v-if="modeLabel" type="info" :bordered="false" class="mode-alert glass-panel">
        <template #icon>
          <n-icon size="18" :component="SettingsOutlined" />
        </template>
        当前模式：{{ modeLabel }}
        <n-button
          v-if="isElectron"
          text
          type="error"
          size="tiny"
          style="margin-left: 12px; vertical-align: middle;"
          @click="handleClearData"
        >
          <template #icon>
            <n-icon :component="DeleteOutlined" />
          </template>
          清除数据
        </n-button>
      </n-alert>

      <div class="config-card glass-panel" v-if="isServerMode">
        <div class="section-header">
          <n-icon size="20" :component="DnsOutlined" color="#009688" />
          <span class="section-title">服务端配置</span>
        </div>
        
        <n-form class="custom-form" label-placement="top" :show-feedback="false">
          <div class="form-grid">
            <n-form-item label="后端端口">
              <n-input-number v-model:value="serverForm.backendPort" :min="1" :max="65535" placeholder="12345" class="custom-input" />
            </n-form-item>
            <n-form-item label="前端端口">
              <n-input-number v-model:value="serverForm.frontendPort" :min="1" :max="65535" placeholder="23456" class="custom-input" />
            </n-form-item>
          </div>
          
          <n-form-item label="基础文件夹">
            <n-input-group>
              <n-input v-model:value="serverForm.baseFolderPath" placeholder="音乐库根目录" readonly class="custom-input" />
              <n-button type="primary" ghost @click="pickFolder" class="folder-btn">
                <template #icon><n-icon :component="FolderOpenOutlined" /></template>
              </n-button>
            </n-input-group>
          </n-form-item>
        </n-form>

        <div class="address-preview" v-if="ips.length > 0">
          <div class="preview-header">
            <n-icon :component="CheckCircleOutlined" color="#009688" />
            <span>局域网访问地址已就绪</span>
          </div>
          <div class="preview-scroll">
            <div v-for="ip in ips" :key="ip" class="preview-item">
              <span class="label">IP: {{ ip }}</span>
              <span class="value">{{ `http://${ip}:${serverForm.backendPort}` }}</span>
            </div>
          </div>
        </div>

        <div class="actions">
          <n-button secondary type="primary" @click="checkPorts" class="action-btn secondary">检测端口</n-button>
          <n-button type="primary" :loading="saving" @click="applyServerConfig" class="action-btn primary">启动服务</n-button>
        </div>
      </div>

      <div class="config-card glass-panel" v-else>
        <div class="section-header">
          <n-icon size="20" :component="ConnectWithoutContactOutlined" color="#009688" />
          <span class="section-title">连接到服务器</span>
        </div>
        
        <n-form class="custom-form" label-placement="top" :show-feedback="false">
          <n-form-item label="服务器 IP">
            <n-input v-model:value="clientForm.backendIp" placeholder="例如: 192.168.1.10" class="custom-input">
              <template #prefix>
                <n-icon :component="DnsOutlined" />
              </template>
            </n-input>
          </n-form-item>
          
          <n-form-item label="端口号">
            <n-input-number v-model:value="clientForm.backendPort" :min="1" :max="65535" placeholder="12345" class="custom-input" :show-button="false">
              <template #prefix>
                <n-icon :component="NumbersOutlined" />
              </template>
            </n-input-number>
          </n-form-item>
        </n-form>

        <div class="scan-section" v-if="isCapacitor">
          <div class="divider">
            <span>或者</span>
          </div>
          <n-button class="scan-btn" :loading="scanning" @click="scanQrCode">
            <template #icon>
              <n-icon :component="QrCodeScannerOutlined" />
            </template>
            扫描二维码连接
          </n-button>
        </div>

        <transition name="fade-slide">
          <div v-if="qrApiParam || scanHint" class="status-card success">
            <n-icon :component="CheckCircleOutlined" color="#009688" />
            <div class="status-text">
              <div class="status-title">已识别服务器</div>
              <div class="status-desc">{{ clientPreviewUrl }}</div>
            </div>
          </div>
        </transition>

        <div class="actions">
          <n-button type="primary" :loading="saving" @click="applyClientConfig" class="action-btn primary">
            立即连接
            <template #icon>
              <n-icon :component="ArrowForwardOutlined" />
            </template>
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useMessage, NAlert, NButton, NForm, NFormItem, NInput, NInputNumber, NIcon, NInputGroup, useDialog } from "naive-ui";
import { useRouter } from "vue-router";
import { Capacitor } from "@capacitor/core";
import { BarcodeScanner, SupportedFormat } from "@capacitor-community/barcode-scanner";
import { 
  SettingsOutlined, 
  FolderOpenOutlined, 
  DnsOutlined, 
  NumbersOutlined, 
  QrCodeScannerOutlined, 
  ArrowForwardOutlined,
  CheckCircleOutlined,
  ConnectWithoutContactOutlined,
  DeleteOutlined,
  ArrowBackOutlined
} from "@vicons/material";

const router = useRouter();
const message = useMessage();
const dialog = useDialog();

const appMode = import.meta.env.VITE_APP_MODE as string | undefined;
const isElectron = typeof navigator !== "undefined" && navigator.userAgent.includes("Electron");
const isCapacitor = typeof window !== "undefined" && Capacitor.isNativePlatform();
const isServerMode = computed(() => appMode === "server" && isElectron);
const modeLabel = computed(() => {
  if (!isElectron) return "";
  if (appMode === "server") return "Electron 服务器端";
  if (appMode === "client") return "Electron 客户端";
  return "Electron";
});

const saving = ref(false);
const scanning = ref(false);
const scanHint = ref("");

const serverForm = reactive({
  backendPort: 12345,
  frontendPort: 23456,
  baseFolderPath: "",
});

const clientForm = reactive({
  backendIp: "",
  backendPort: 12345,
});

const qrApiParam = ref("");
const ips = ref<string[]>([]);

const normalizeServerUrl = (value: string) => {
  const trimmed = value.trim();
  if (!trimmed) return "";
  return trimmed.endsWith("/") ? trimmed.slice(0, -1) : trimmed;
};

const parseServerFromUrl = (value: string) => {
  const normalized = normalizeServerUrl(value);
  if (!normalized) return { ip: "", port: 12345 };
  try {
    const url = new URL(normalized);
    return {
      ip: url.hostname || "",
      port: Number(url.port) || 12345,
    };
  } catch {
    return { ip: normalized, port: 12345 };
  }
};

const fillClientServer = (value: string) => {
  const parsed = parseServerFromUrl(value);
  clientForm.backendIp = parsed.ip;
  clientForm.backendPort = parsed.port;
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

const scanQrCode = async () => {
  if (!isCapacitor) {
    message.warning("仅移动端 App 支持扫码");
    return;
  }
  scanning.value = true;
  scanHint.value = "";
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
    await startScannerOverlay();
    const result = await BarcodeScanner.startScan({
      targetedFormats: [SupportedFormat.QR_CODE],
    });
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
    fillClientServer(apiValue);
    qrApiParam.value = apiValue;
    scanHint.value = `已识别并填充：${apiValue}`;
    message.success("二维码识别成功，请确认后保存");
  } catch (err: any) {
    const prettyMessage = buildScanErrorMessage(err);
    scanHint.value = prettyMessage;
    message.error(prettyMessage);
  } finally {
    await stopScannerOverlay();
    scanning.value = false;
  }
};

const handleClearData = () => {
  dialog.warning({
    title: "清除所有数据",
    content: "此操作将清除应用在本地的所有配置、缓存和数据库，并重置应用。确定要继续吗？",
    positiveText: "确定清除",
    negativeText: "取消",
    onPositiveClick: async () => {
      // 1. 清除 LocalStorage
      localStorage.clear();
      
      // 2. 如果是 Electron，调用 IPC 清除 UserData
      if (isElectron && window.ipcRenderer) {
        try {
          const res = await window.ipcRenderer.invoke('app-clear-data');
          if (!res.success) {
            message.error("清除数据失败: " + res.error);
            return;
          }
          // Electron 主进程会处理重启，这里不需要做额外操作
        } catch (e: any) {
          message.error("清除数据请求失败: " + (e.message || e));
        }
      } else {
        // Web 或 Capacitor 环境，直接刷新重置
        message.success("数据已清除，正在重启...");
        setTimeout(() => {
          window.location.reload();
        }, 1000);
      }
    }
  });
};

const clientPreviewUrl = computed(() => {
  const ip = clientForm.backendIp.trim();
  const port = Number(clientForm.backendPort);
  if (!ip || !Number.isFinite(port) || port <= 0) return "";
  return `http://${ip}:${port}`;
});

const invoke = async (channel: string, payload?: any) => {
  const ipc = (window as any)?.ipcRenderer;
  if (!ipc?.invoke) throw new Error("IPC unavailable");
  return ipc.invoke(channel, payload);
};

const pickFolder = async () => {
  try {
    const res = await invoke("select-directory");
    if (res?.canceled) return;
    if (Array.isArray(res?.filePaths) && res.filePaths[0]) {
      serverForm.baseFolderPath = res.filePaths[0];
    }
  } catch {
    message.error("选择文件夹失败");
  }
};

const loadIps = async () => {
  if (!isElectron) return;
  try {
    const res = await invoke("get-local-ips");
    if (Array.isArray(res?.ips)) {
      ips.value = res.ips;
    }
  } catch {
    ips.value = [];
  }
};

const loadConfig = async () => {
  if (!isElectron) return;
  try {
    const cfg = await invoke("app-config-get");
    if (cfg?.backend_port) serverForm.backendPort = Number(cfg.backend_port) || serverForm.backendPort;
    if (cfg?.frontend_port) serverForm.frontendPort = Number(cfg.frontend_port) || serverForm.frontendPort;
    if (cfg?.base_folder) serverForm.baseFolderPath = String(cfg.base_folder || "");
  } catch {}
};

const checkPorts = async () => {
  saving.value = true;
  try {
    const res = await invoke("check-ports", {
      backendPort: serverForm.backendPort,
      frontendPort: serverForm.frontendPort,
    });
    if (res?.backendAvailable === false) {
      message.error(`后端端口 ${serverForm.backendPort} 不可用`);
      return;
    }
    if (res?.frontendAvailable === false) {
      message.error(`前端端口 ${serverForm.frontendPort} 不可用`);
      return;
    }
    message.success("端口可用");
  } catch {
    message.error("端口检测失败");
  } finally {
    saving.value = false;
  }
};

const waitForBackend = async (origin: string) => {
  const url = `${origin}/api/system/stats`;
  for (let i = 0; i < 20; i++) {
    try {
      const res = await fetch(url, { credentials: "include" });
      if (res.ok) return true;
    } catch {}
    await new Promise((r) => setTimeout(r, 500));
  }
  return false;
};

const applyServerConfig = async () => {
  saving.value = true;
  try {
    const baseFolder = serverForm.baseFolderPath.trim();
    if (!baseFolder) {
      message.error("请选择基础文件夹");
      return;
    }
    await invoke("apply-initial-config", {
      mode: "server",
      backendPort: serverForm.backendPort,
      frontendPort: serverForm.frontendPort,
      baseFolderPath: baseFolder,
      accessIp: "",
    });

    localStorage.setItem("init_done", "true");
    localStorage.setItem("backend_port", String(serverForm.backendPort));
    localStorage.setItem("frontend_port", String(serverForm.frontendPort));
    localStorage.setItem("server_url", `http://localhost:${serverForm.backendPort}`);

    const backendOrigin = `http://localhost:${serverForm.backendPort}`;
    const ok = await waitForBackend(backendOrigin);
    if (!ok) {
      message.warning("后端启动中，请稍后刷新重试");
      router.push("/");
      return;
    }
    try {
      await fetch(`${backendOrigin}/api/file/initFolder`, { method: "POST" });
    } catch {}
    message.success("配置完成");
    router.push("/");
    location.reload();
  } catch {
    message.error("保存失败");
  } finally {
    saving.value = false;
  }
};

const applyClientConfig = async () => {
  saving.value = true;
  try {
    const ip = clientForm.backendIp.trim();
    const port = Number(clientForm.backendPort);
    if (!ip) {
      message.error("请输入后端 IP");
      return;
    }
    if (!Number.isFinite(port) || port <= 0 || port > 65535) {
      message.error("请输入正确的端口");
      return;
    }
    const value = normalizeServerUrl(`http://${ip}:${port}`);
    if (!value) {
      message.error("请输入后端地址");
      return;
    }
    localStorage.setItem("server_url", value);
    localStorage.setItem("init_done", "true");

    try {
      const url = new URL(value);
      if (url.port) localStorage.setItem("backend_port", url.port);
    } catch {}

    message.success("配置完成");
    router.push("/");
    location.reload();
  } finally {
    saving.value = false;
  }
};

onMounted(async () => {
  const apiParam = new URLSearchParams(window.location.search).get("api") || "";
  qrApiParam.value = apiParam ? normalizeServerUrl(apiParam) : "";

  if (!isServerMode.value) {
    const stored = localStorage.getItem("server_url") || "";
    const value = qrApiParam.value || stored;
    if (!value) return;
    fillClientServer(value);
  } else {
    await loadConfig();
    await loadIps();
  }
});
</script>

<style scoped lang="scss">
.init-page {
  height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  position: relative;
  background-color: #f0f2f5;
  overflow: hidden;
}

.back-btn {
  position: absolute;
  top: 24px;
  left: 24px;
  z-index: 10;
  background: rgba(255, 255, 255, 0.5);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  
  &:hover {
    background: white;
    transform: scale(1.1);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
  }
}

.glass-bg {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle at 50% 50%, rgba(0, 150, 136, 0.15), rgba(255, 255, 255, 0) 60%),
              linear-gradient(135deg, #e0f7fa 0%, #f1f8e9 100%);
  z-index: 0;
  animation: bg-drift 20s ease-in-out infinite alternate;
}

.init-container {
  width: 480px;
  max-width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 4px;

  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.header {
  text-align: center;
  margin-bottom: 8px;
  
  .logo-wrapper {
    width: 80px;
    height: 80px;
    margin: 0 auto 16px;
    background: white;
    border-radius: 20px;
    padding: 12px;
    box-shadow: 0 12px 24px rgba(0, 150, 136, 0.15);
    display: flex;
    align-items: center;
    justify-content: center;
    
    .logo {
      width: 100%;
      height: 100%;
      object-fit: contain;
    }
  }
  
  .title-group {
    .app-title {
      font-size: 28px;
      font-weight: 800;
      color: #1a1a1a;
      margin: 0;
      letter-spacing: -0.5px;
      background: linear-gradient(135deg, #009688, #00796b);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
    
    .subtitle {
      font-size: 14px;
      color: #666;
      margin: 4px 0 0;
      font-weight: 500;
    }
  }
}

.glass-panel {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.6);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.04);
  padding: 24px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.06);
  }
}

.mode-alert {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: #00796b;
  font-weight: 600;
  font-size: 13px;
}

.config-card {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  
  .section-title {
    font-size: 16px;
    font-weight: 700;
    color: #333;
  }
}

.custom-form {
  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }
  
  :deep(.n-form-item-label) {
    font-size: 13px;
    color: #666;
    font-weight: 500;
  }
}

.custom-input, .custom-select {
  :deep(.n-input), :deep(.n-base-selection) {
    border-radius: 12px;
    background-color: rgba(255, 255, 255, 0.5);
    transition: all 0.2s;
    
    &:hover, &.n-input--focus {
      background-color: white;
      box-shadow: 0 0 0 1px #00968820;
    }
  }
}

.folder-btn {
  border-radius: 0 12px 12px 0;
}

.address-preview {
  background: rgba(0, 150, 136, 0.04);
  border-radius: 16px;
  padding: 16px;
  display: grid;
  gap: 12px;
  border: 1px dashed rgba(0, 150, 136, 0.2);
  
  .preview-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 13px;
    
    .label {
      color: #666;
    }
    
    .value {
      font-family: 'JetBrains Mono', monospace;
      font-weight: 600;
      color: #00796b;
      background: rgba(255, 255, 255, 0.5);
      padding: 2px 8px;
      border-radius: 6px;
    }
  }
}

.scan-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  
  .divider {
    display: flex;
    align-items: center;
    color: #999;
    font-size: 12px;
    
    &::before, &::after {
      content: '';
      flex: 1;
      height: 1px;
      background: rgba(0, 0, 0, 0.06);
    }
    
    span {
      padding: 0 12px;
    }
  }
  
  .scan-btn {
    width: 100%;
    height: 48px;
    border-radius: 14px;
    font-weight: 600;
    border: 1px solid rgba(0, 0, 0, 0.08);
    background: white;
    transition: all 0.2s;
    
    &:active {
      background: #f5f5f5;
      transform: scale(0.98);
    }
  }
}

.status-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: #e0f2f1;
  border-radius: 14px;
  border: 1px solid rgba(0, 150, 136, 0.1);
  
  .status-text {
    flex: 1;
    overflow: hidden;
    
    .status-title {
      font-size: 13px;
      font-weight: 700;
      color: #00695c;
    }
    
    .status-desc {
      font-size: 12px;
      color: #00897b;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      margin-top: 2px;
      font-family: 'JetBrains Mono', monospace;
    }
  }
}

.actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
  
  .action-btn {
    height: 48px;
    border-radius: 14px;
    font-weight: 600;
    font-size: 15px;
    
    &.primary {
      flex: 2;
      box-shadow: 0 8px 20px rgba(0, 150, 136, 0.25);
    }
    
    &.secondary {
      flex: 1;
    }
  }
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@keyframes bg-drift {
  0% { transform: translate(0, 0); }
  100% { transform: translate(20px, 20px); }
}

:global(html.scanner-active),
:global(body.scanner-active) {
  background: transparent !important;
}

:global(body.scanner-active .init-page) {
  opacity: 0;
}

@media (max-width: 768px) {
  .init-page {
    padding: 16px;
    align-items: flex-start;
    padding-top: 40px;
  }

  .init-container {
    width: 100%;
  }
  
  .glass-panel {
    padding: 20px;
  }
}
</style>
