import { app as i, ipcMain as l, BrowserWindow as d, screen as w, shell as m } from "electron";
import { createRequire as h } from "node:module";
import { fileURLToPath as y } from "node:url";
import n from "node:path";
import g from "node:os";
const I = h(import.meta.url), p = n.dirname(y(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
i.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill");
process.env.APP_ROOT = n.join(p, "..");
const r = process.env.VITE_DEV_SERVER_URL, P = n.join(process.env.APP_ROOT, "dist-electron"), u = n.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = r ? n.join(process.env.APP_ROOT, "public") : u;
g.release().startsWith("6.1") && i.disableHardwareAcceleration();
process.platform === "win32" && i.setAppUserModelId("LocalMusicPlayer");
i.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill,PasswordManager");
i.requestSingleInstanceLock() || (i.quit(), process.exit(0));
let o, c, e = null;
function f() {
  let t = n.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  I("fs").existsSync(t) || (t = n.join(process.env.VITE_PUBLIC, "images/logo/favicon.png")), c = new d({
    width: 500,
    height: 300,
    transparent: !0,
    frame: !1,
    alwaysOnTop: !0,
    icon: t
  }), c.loadFile(n.join(process.env.VITE_PUBLIC, "loading.html")), o = new d({
    title: "Local Music Player",
    show: !1,
    // 先隐藏主窗口
    icon: t,
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: n.join(p, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    }
  }), o.setMenu(null), o.webContents.on("did-finish-load", () => {
    o?.webContents.send("main-process-message", (/* @__PURE__ */ new Date()).toLocaleString());
  }), o.once("ready-to-show", () => {
    setTimeout(() => {
      c?.destroy(), c = null, o?.show(), o?.focus();
    }, 2e3);
  }), r ? (o.loadURL(r), o.webContents.openDevTools()) : o.loadFile(n.join(u, "index.html")), o.webContents.setWindowOpenHandler(({ url: a }) => (a.startsWith("https:") && m.openExternal(a), { action: "deny" }));
}
function k() {
  if (e) return;
  const { width: t, height: s } = w.getPrimaryDisplay().workAreaSize;
  if (e = new d({
    width: 800,
    height: 120,
    x: (t - 800) / 2,
    y: s - 150,
    frame: !1,
    transparent: !0,
    alwaysOnTop: !0,
    skipTaskbar: !0,
    resizable: !0,
    // 允许调整大小
    webPreferences: {
      preload: n.join(p, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    },
    backgroundColor: "#00000000"
    // Ensure transparency
  }), r) {
    const a = r.endsWith("/") ? r : `${r}/`;
    e.loadURL(`${a}#/desktop-lyric`);
  } else
    e.loadFile(n.join(u, "index.html"), { hash: "desktop-lyric" });
  e.on("closed", () => {
    e = null;
  });
}
l.on("open-desktop-lyric", () => {
  k();
});
l.on("close-desktop-lyric", () => {
  e && e.close();
});
l.on("update-desktop-lyric", (t, s) => {
  e && e.webContents.send("update-lyric", s);
});
l.on("desktop-lyric-control", (t, s) => {
  o && o.webContents.send("player-control", s);
});
l.on("lock-desktop-lyric", (t, s) => {
  e && (e.setIgnoreMouseEvents(s, { forward: !0 }), s ? (e.setFocusable(!1), e.webContents.send("desktop-lyric-locked", !0)) : (e.setFocusable(!0), e.webContents.send("desktop-lyric-locked", !1)));
});
l.on("unlock-desktop-lyric", () => {
  e && (e.setIgnoreMouseEvents(!1, { forward: !0 }), e.setFocusable(!0), e.webContents.send("desktop-lyric-locked", !1));
});
l.on("update-desktop-lyric-settings", (t, s) => {
  e && e.webContents.send("update-settings", s);
});
i.whenReady().then(f);
i.on("window-all-closed", () => {
  o = null, process.platform !== "darwin" && i.quit();
});
i.on("second-instance", () => {
  o && (o.isMinimized() && o.restore(), o.focus());
});
i.on("activate", () => {
  const t = d.getAllWindows();
  t.length ? t[0].focus() : f();
});
export {
  P as MAIN_DIST,
  u as RENDERER_DIST,
  r as VITE_DEV_SERVER_URL
};
