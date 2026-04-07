import { app, ipcMain, dialog, BrowserWindow, shell } from "electron";
import { createRequire } from "node:module";
import { fileURLToPath } from "node:url";
import path from "node:path";
import os from "node:os";
import { spawn } from "child_process";
import * as http from "node:http";
import * as net from "node:net";
import fs from "node:fs";
createRequire(import.meta.url);
const __dirname$1 = path.dirname(fileURLToPath(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
app.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill");
process.env.APP_ROOT = path.join(__dirname$1, "..");
const APP_NAME_BASE = "RLMusic";
const getAppId = (mode) => `${APP_NAME_BASE}-${mode}`;
const VITE_DEV_SERVER_URL = process.env["VITE_DEV_SERVER_URL"];
const appMode = process.env.VITE_APP_MODE || (fs.existsSync(path.join(process.resourcesPath, process.platform === "win32" ? "server.exe" : "server")) ? "server" : "client");
const MAIN_DIST = path.join(process.env.APP_ROOT, "dist-electron");
const RENDERER_DIST = path.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL ? path.join(process.env.APP_ROOT, "public") : RENDERER_DIST;
if (os.release().startsWith("6.1")) app.disableHardwareAcceleration();
if (process.platform === "win32") app.setAppUserModelId(getAppId(appMode));
app.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill,PasswordManager");
app.setPath("userData", path.join(app.getPath("appData"), getAppId(appMode)));
if (!app.requestSingleInstanceLock({ mode: appMode })) {
  app.quit();
  process.exit(0);
}
let win;
let splash;
let serverProcess = null;
let frontendServer = null;
const getAppConfigPath = () => path.join(app.getPath("userData"), "app-config.json");
const readAppConfig = () => {
  try {
    const raw = fs.readFileSync(getAppConfigPath(), "utf-8");
    const parsed = JSON.parse(raw || "{}");
    return {
      init_done: !!parsed.init_done,
      backend_port: Number(parsed.backend_port) || 12345,
      frontend_port: Number(parsed.frontend_port) || 23456,
      base_folder: String(parsed.base_folder || ""),
      access_ip: String(parsed.access_ip || "")
    };
  } catch {
    return { init_done: false, backend_port: 12345, frontend_port: 23456, base_folder: "", access_ip: "" };
  }
};
const writeAppConfig = (cfg) => {
  const next = { ...readAppConfig(), ...cfg };
  fs.mkdirSync(path.dirname(getAppConfigPath()), { recursive: true });
  fs.writeFileSync(getAppConfigPath(), JSON.stringify(next, null, 2), "utf-8");
  return next;
};
const isPortAvailable = (port) => new Promise((resolve) => {
  const server = net.createServer();
  server.once("error", () => resolve(false));
  server.once("listening", () => server.close(() => resolve(true)));
  server.listen(port, "0.0.0.0");
});
const getLocalIPs = () => {
  const nets = os.networkInterfaces();
  const ips = [];
  Object.values(nets).forEach((items) => {
    (items || []).forEach((n) => {
      if (n.family === "IPv4" && !n.internal) ips.push(n.address);
    });
  });
  return ips;
};
const updateServerConfigYml = (backendPort, baseFolderPath) => {
  const configPath = path.join(process.resourcesPath, "config.yml");
  if (!fs.existsSync(configPath)) return;
  const escapedPath = baseFolderPath.replace(/'/g, "''");
  const raw = fs.readFileSync(configPath, "utf-8");
  const withPort = raw.replace(/(^\s*Port:\s*).*/m, `$1${backendPort}`);
  const withPath = withPort.replace(/(^\s*FilePath:\s*).*/m, `$1'${escapedPath}'`);
  const withName = withPath.replace(/(^\s*FileName:\s*).*/m, `$1''`);
  fs.writeFileSync(configPath, withName, "utf-8");
};
const guessMime = (p) => {
  const ext = path.extname(p).toLowerCase();
  if (ext === ".html") return "text/html; charset=utf-8";
  if (ext === ".js" || ext === ".mjs") return "text/javascript; charset=utf-8";
  if (ext === ".css") return "text/css; charset=utf-8";
  if (ext === ".json") return "application/json; charset=utf-8";
  if (ext === ".svg") return "image/svg+xml";
  if (ext === ".png") return "image/png";
  if (ext === ".jpg" || ext === ".jpeg") return "image/jpeg";
  if (ext === ".ico") return "image/x-icon";
  if (ext === ".woff") return "font/woff";
  if (ext === ".woff2") return "font/woff2";
  if (ext === ".ttf") return "font/ttf";
  return "application/octet-stream";
};
const stopFrontendServer = async () => {
  if (!frontendServer) return;
  await new Promise((resolve) => frontendServer?.close(() => resolve()));
  frontendServer = null;
};
const startFrontendServer = async (port) => {
  await stopFrontendServer();
  const root = RENDERER_DIST;
  const indexPath = path.join(root, "index.html");
  frontendServer = http.createServer((req, res) => {
    const method = req.method || "GET";
    if (method !== "GET" && method !== "HEAD") {
      res.statusCode = 405;
      res.end();
      return;
    }
    try {
      const url = new URL(req.url || "/", `http://${req.headers.host || "localhost"}`);
      const pathname = decodeURIComponent(url.pathname || "/");
      const safePath = pathname.replace(/\\/g, "/");
      const requested = safePath === "/" ? "/index.html" : safePath;
      const resolved = path.resolve(path.join(root, requested));
      const rootResolved = path.resolve(root);
      const targetPath = resolved.startsWith(rootResolved) ? resolved : indexPath;
      const exists = fs.existsSync(targetPath) && fs.statSync(targetPath).isFile();
      const finalPath = exists ? targetPath : indexPath;
      const data = fs.readFileSync(finalPath);
      res.setHeader("Content-Type", guessMime(finalPath));
      res.statusCode = 200;
      if (method === "HEAD") {
        res.end();
      } else {
        res.end(data);
      }
    } catch {
      res.statusCode = 500;
      res.end();
    }
  });
  await new Promise((resolve, reject) => {
    frontendServer?.once("error", reject);
    frontendServer?.listen(port, "0.0.0.0", () => resolve());
  });
};
function startServer() {
  if (VITE_DEV_SERVER_URL) return;
  const serverName = process.platform === "win32" ? "server.exe" : "server";
  const resourcesPath = process.resourcesPath;
  const serverPath = path.join(resourcesPath, serverName);
  if (fs.existsSync(serverPath)) {
    console.log(`Starting server from: ${serverPath}`);
    if (serverProcess) {
      serverProcess.kill();
      serverProcess = null;
    }
    serverProcess = spawn(serverPath, [], {
      cwd: resourcesPath,
      // 设置工作目录为 resources
      windowsHide: true,
      stdio: "ignore"
      // 忽略输出，避免缓冲区填满挂起
    });
    serverProcess.on("error", (err) => {
      console.error("Failed to start server:", err);
    });
    serverProcess.on("close", (code) => {
      console.log(`Server process exited with code ${code}`);
      serverProcess = null;
    });
  } else {
    console.log("Server binary not found, running in client-only mode.");
  }
}
function createWindow() {
  const cfg = readAppConfig();
  if (cfg.init_done) {
    startServer();
    startFrontendServer(cfg.frontend_port).catch(() => {
    });
  }
  let iconPath = path.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  if (!fs.existsSync(iconPath)) {
    iconPath = path.join(process.env.VITE_PUBLIC, "images/logo/favicon.png");
  }
  splash = new BrowserWindow({
    width: 500,
    height: 300,
    transparent: true,
    frame: false,
    alwaysOnTop: true,
    icon: iconPath
  });
  splash.loadFile(path.join(process.env.VITE_PUBLIC, "loading.html"));
  win = new BrowserWindow({
    title: "RLmusic",
    show: false,
    // 先隐藏主窗口
    icon: iconPath,
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: path.join(__dirname$1, "../dist-electron/preload.mjs"),
      nodeIntegration: true,
      contextIsolation: true
    }
  });
  win.setMenu(null);
  win.webContents.on("did-finish-load", () => {
    win?.webContents.send("main-process-message", (/* @__PURE__ */ new Date()).toLocaleString());
  });
  win.once("ready-to-show", () => {
    setTimeout(() => {
      splash?.destroy();
      splash = null;
      win?.show();
      win?.focus();
    }, 2e3);
  });
  if (VITE_DEV_SERVER_URL) {
    win.loadURL(VITE_DEV_SERVER_URL);
    win.webContents.openDevTools();
  } else {
    win.loadFile(path.join(RENDERER_DIST, "index.html"));
  }
  win.webContents.setWindowOpenHandler(({ url }) => {
    if (url.startsWith("https:")) shell.openExternal(url);
    return { action: "deny" };
  });
}
ipcMain.handle("app-clear-data", async () => {
  try {
    if (win) {
      await win.webContents.session.clearStorageData();
    }
    const userDataPath = app.getPath("userData");
    if (fs.existsSync(userDataPath)) {
      const files = fs.readdirSync(userDataPath);
      for (const file of files) {
        if (file === "Lockfile" || file.startsWith("Singleton") || file === "TransportSecurity") continue;
        try {
          const curPath = path.join(userDataPath, file);
          fs.rmSync(curPath, { recursive: true, force: true });
        } catch (e) {
          console.warn(`Failed to delete ${file}:`, e.message);
        }
      }
    }
    if (serverProcess) {
      serverProcess.kill();
      serverProcess = null;
    }
    await stopFrontendServer();
    app.relaunch();
    app.exit(0);
    return { success: true };
  } catch (error) {
    console.error("Failed to clear app data:", error);
    return { success: false, error: error.message };
  }
});
ipcMain.handle("app-config-get", async () => {
  return readAppConfig();
});
ipcMain.handle("show-save-dialog", async (event, options) => {
  const { filePath, canceled } = await dialog.showSaveDialog({
    ...options,
    filters: [
      { name: "Excel Files", extensions: ["xlsx"] },
      { name: "All Files", extensions: ["*"] }
    ]
  });
  return { filePath, canceled };
});
ipcMain.handle("save-file", async (event, { path: filePath, data }) => {
  try {
    let buffer;
    if (Buffer.isBuffer(data)) {
      buffer = data;
    } else if (data instanceof ArrayBuffer) {
      buffer = Buffer.from(new Uint8Array(data));
    } else if (ArrayBuffer.isView(data)) {
      buffer = Buffer.from(data.buffer, data.byteOffset, data.byteLength);
    } else {
      buffer = Buffer.from(data);
    }
    fs.writeFileSync(filePath, buffer);
    return { success: true };
  } catch (err) {
    return { success: false, error: err.message };
  }
});
ipcMain.handle("select-directory", async () => {
  const res = await dialog.showOpenDialog({
    properties: ["openDirectory", "createDirectory"]
  });
  return res;
});
ipcMain.handle("get-local-ips", async () => {
  return { ips: getLocalIPs() };
});
ipcMain.handle("check-ports", async (_event, payload) => {
  const backendPort = Number(payload?.backendPort) || 0;
  const frontendPort = Number(payload?.frontendPort) || 0;
  return {
    backendAvailable: backendPort > 0 ? await isPortAvailable(backendPort) : false,
    frontendAvailable: frontendPort > 0 ? await isPortAvailable(frontendPort) : false
  };
});
ipcMain.handle("apply-initial-config", async (_event, payload) => {
  const mode = String(payload?.mode || "");
  const backendPort = Number(payload?.backendPort) || 12345;
  const frontendPort = Number(payload?.frontendPort) || 23456;
  const baseFolderPath = String(payload?.baseFolderPath || "");
  const accessIp = String(payload?.accessIp || "");
  if (mode === "server") {
    const backendOk = await isPortAvailable(backendPort);
    const frontendOk = await isPortAvailable(frontendPort);
    if (!backendOk) throw new Error("backend port unavailable");
    if (!frontendOk) throw new Error("frontend port unavailable");
    writeAppConfig({
      init_done: true,
      backend_port: backendPort,
      frontend_port: frontendPort,
      base_folder: baseFolderPath,
      access_ip: accessIp
    });
    updateServerConfigYml(backendPort, baseFolderPath);
    startServer();
    await startFrontendServer(frontendPort);
    return { ok: true };
  }
  writeAppConfig({ init_done: true });
  return { ok: true };
});
app.whenReady().then(createWindow);
app.on("window-all-closed", () => {
  win = null;
  if (serverProcess) {
    serverProcess.kill();
    serverProcess = null;
  }
  stopFrontendServer().catch(() => {
  });
  if (process.platform !== "darwin") app.quit();
});
app.on("second-instance", () => {
  if (win) {
    if (win.isMinimized()) win.restore();
    win.focus();
  }
});
app.on("activate", () => {
  const allWindows = BrowserWindow.getAllWindows();
  if (allWindows.length) {
    allWindows[0].focus();
  } else {
    createWindow();
  }
});
export {
  MAIN_DIST,
  RENDERER_DIST,
  VITE_DEV_SERVER_URL
};
