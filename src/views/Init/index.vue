<template>
  <div class="init-page">
    <div class="init-card">
      <div class="header">
        <div class="title">首次启动向导</div>
        <div class="subtitle">完成配置后即可正常使用与局域网访问</div>
      </div>

      <n-alert v-if="modeLabel" type="info" :bordered="false" class="mode-alert">
        当前模式：{{ modeLabel }}
      </n-alert>

      <n-card v-if="isServerMode" :bordered="false" class="section-card">
        <div class="section-title">服务器端配置</div>
        <n-form label-placement="left" label-width="120">
          <n-form-item label="后端端口">
            <n-input-number v-model:value="serverForm.backendPort" :min="1" :max="65535" />
          </n-form-item>
          <n-form-item label="前端端口">
            <n-input-number v-model:value="serverForm.frontendPort" :min="1" :max="65535" />
          </n-form-item>
          <n-form-item label="基础文件夹">
            <div class="path-row">
              <n-input v-model:value="serverForm.baseFolderPath" placeholder="请选择或输入基础文件夹路径" />
              <n-button secondary type="primary" class="pick-btn" @click="pickFolder">选择</n-button>
            </div>
          </n-form-item>
          <n-form-item label="对外访问 IP">
            <n-select
              v-model:value="serverForm.accessIp"
              :options="ipOptions"
              placeholder="自动选择局域网 IP"
              :disabled="ipOptions.length === 0"
            />
          </n-form-item>
        </n-form>

        <div class="address-panel" v-if="ips.length > 0">
          <div class="address-title">对外访问地址</div>
          <div class="address-grid">
            <div class="address-item">
              <div class="address-label">前端</div>
              <div class="address-value">{{ frontendUrlPreview }}</div>
            </div>
            <div class="address-item">
              <div class="address-label">后端</div>
              <div class="address-value">{{ backendUrlPreview }}</div>
            </div>
          </div>
        </div>

        <div class="actions">
          <n-button secondary @click="checkPorts">检测端口</n-button>
          <n-button type="primary" :loading="saving" @click="applyServerConfig">保存并启动</n-button>
        </div>
      </n-card>

      <n-card v-else :bordered="false" class="section-card">
        <div class="section-title">客户端配置</div>
        <n-form label-placement="left" label-width="120">
          <n-form-item label="后端地址">
            <n-input v-model:value="clientForm.backendUrl" placeholder="http://192.168.1.10:12345" />
          </n-form-item>
        </n-form>

        <n-alert v-if="qrApiParam" type="success" :bordered="false" class="hint-alert">
          已从二维码参数识别到后端地址：{{ qrApiParam }}
        </n-alert>

        <div class="actions">
          <n-button type="primary" :loading="saving" @click="applyClientConfig">保存并进入</n-button>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useMessage, NAlert, NButton, NCard, NForm, NFormItem, NInput, NInputNumber, NSelect } from "naive-ui";
import { useRouter } from "vue-router";

const router = useRouter();
const message = useMessage();

const appMode = import.meta.env.VITE_APP_MODE as string | undefined;
const isElectron = typeof navigator !== "undefined" && navigator.userAgent.includes("Electron");
const isServerMode = computed(() => appMode === "server" && isElectron);
const modeLabel = computed(() => {
  if (!isElectron) return "";
  if (appMode === "server") return "Electron 服务器端";
  if (appMode === "client") return "Electron 客户端";
  return "Electron";
});

const saving = ref(false);

const serverForm = reactive({
  backendPort: 12345,
  frontendPort: 23456,
  baseFolderPath: "",
  accessIp: "",
});

const clientForm = reactive({
  backendUrl: "",
});

const qrApiParam = ref("");
const ips = ref<string[]>([]);
const ipOptions = computed(() => ips.value.map((ip) => ({ label: ip, value: ip })));
const currentIp = computed(() => {
  if (serverForm.accessIp) return serverForm.accessIp;
  if (ips.value.length === 0) return "";
  return ips.value[0];
});

const normalizeServerUrl = (value: string) => {
  const trimmed = value.trim();
  if (!trimmed) return "";
  return trimmed.endsWith("/") ? trimmed.slice(0, -1) : trimmed;
};

const frontendUrlPreview = computed(() => {
  if (!currentIp.value) return "";
  return `http://${currentIp.value}:${serverForm.frontendPort}/`;
});

const backendUrlPreview = computed(() => {
  if (!currentIp.value) return "";
  return `http://${currentIp.value}:${serverForm.backendPort}`;
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
      if (!serverForm.accessIp && ips.value.length > 0) {
        serverForm.accessIp = ips.value[0] || "";
      }
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
    if (cfg?.access_ip) serverForm.accessIp = String(cfg.access_ip || "");
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
      accessIp: serverForm.accessIp,
    });

    localStorage.setItem("init_done", "true");
    localStorage.setItem("backend_port", String(serverForm.backendPort));
    localStorage.setItem("frontend_port", String(serverForm.frontendPort));
    localStorage.setItem("server_url", `http://localhost:${serverForm.backendPort}`);
    if (serverForm.accessIp) {
      localStorage.setItem("access_ip", serverForm.accessIp);
    }

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
    const value = normalizeServerUrl(clientForm.backendUrl);
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
    clientForm.backendUrl = qrApiParam.value || stored;
  } else {
    await loadConfig();
    await loadIps();
  }
});
</script>

<style scoped lang="scss">
.init-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.init-card {
  width: 760px;
  max-width: 100%;
  background: var(--n-color-modal);
  border: 1px solid var(--n-divider-color);
  border-radius: 22px;
  padding: 28px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.12);
}

.header {
  margin-bottom: 14px;
}

.title {
  font-size: 22px;
  font-weight: 800;
  color: var(--n-text-color);
}

.subtitle {
  margin-top: 6px;
  font-size: 13px;
  color: var(--n-text-color-3);
}

.mode-alert {
  margin-bottom: 16px;
}

.hint-alert {
  margin-top: 10px;
}

.section-card {
  border-radius: 18px;
}

.section-title {
  font-weight: 700;
  margin-bottom: 12px;
  color: var(--n-text-color);
}

.path-row {
  display: flex;
  gap: 10px;
  width: 100%;
}

.pick-btn {
  flex: 0 0 auto;
}

.address-panel {
  margin-top: 14px;
  padding: 14px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.03);
  border: 1px solid var(--n-divider-color);
}

.address-title {
  font-weight: 700;
  margin-bottom: 10px;
}

.address-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.address-item {
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.address-label {
  font-size: 12px;
  color: var(--n-text-color-3);
  margin-bottom: 6px;
}

.address-value {
  font-family: "DM Mono", monospace;
  font-size: 12px;
  color: var(--n-text-color);
  word-break: break-all;
}

.ip-switch {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.ip-tip {
  font-size: 12px;
  color: var(--n-text-color-3);
}

.actions {
  margin-top: 18px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
