import { app, BrowserWindow, shell } from "electron";
import { createRequire } from "node:module";
import { fileURLToPath } from "node:url";
import path from "node:path";
import os from "node:os";
const require$1 = createRequire(import.meta.url);
const __dirname$1 = path.dirname(fileURLToPath(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
process.env.APP_ROOT = path.join(__dirname$1, "..");
const VITE_DEV_SERVER_URL = process.env["VITE_DEV_SERVER_URL"];
const MAIN_DIST = path.join(process.env.APP_ROOT, "dist-electron");
const RENDERER_DIST = path.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL ? path.join(process.env.APP_ROOT, "public") : RENDERER_DIST;
if (os.release().startsWith("6.1")) app.disableHardwareAcceleration();
if (process.platform === "win32") app.setAppUserModelId("LocalMusicPlayer");
if (!app.requestSingleInstanceLock()) {
  app.quit();
  process.exit(0);
}
let win;
let splash;
function createWindow() {
  let iconPath = path.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  const fs = require$1("fs");
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
    title: "Local Music Player",
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
app.whenReady().then(createWindow);
app.on("window-all-closed", () => {
  win = null;
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
