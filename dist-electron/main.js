import { app as a, ipcMain as d, dialog as R, BrowserWindow as P, screen as T, shell as U } from "electron";
import { createRequire as V } from "node:module";
import { fileURLToPath as B } from "node:url";
import s from "node:path";
import D from "node:os";
import { spawn as H } from "child_process";
import * as z from "node:http";
import * as q from "node:net";
import l from "node:fs";
V(import.meta.url);
const k = s.dirname(B(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
a.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill");
process.env.APP_ROOT = s.join(k, "..");
const G = "RLMusic", F = (t) => `${G}-${t}`, h = process.env.VITE_DEV_SERVER_URL, x = process.env.VITE_APP_MODE || (l.existsSync(s.join(process.resourcesPath, process.platform === "win32" ? "server.exe" : "server")) ? "server" : "client"), ie = s.join(process.env.APP_ROOT, "dist-electron"), b = s.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = h ? s.join(process.env.APP_ROOT, "public") : b;
D.release().startsWith("6.1") && a.disableHardwareAcceleration();
process.platform === "win32" && a.setAppUserModelId(F(x));
a.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill,PasswordManager");
a.setPath("userData", s.join(a.getPath("appData"), F(x)));
a.requestSingleInstanceLock(x) || (a.quit(), process.exit(0));
let i, v, r = null, f = null, m = null;
const _ = () => s.join(a.getPath("userData"), "app-config.json"), E = () => {
  try {
    const t = l.readFileSync(_(), "utf-8"), e = JSON.parse(t || "{}");
    return {
      init_done: !!e.init_done,
      backend_port: Number(e.backend_port) || 12345,
      frontend_port: Number(e.frontend_port) || 23456,
      base_folder: String(e.base_folder || ""),
      access_ip: String(e.access_ip || "")
    };
  } catch {
    return { init_done: !1, backend_port: 12345, frontend_port: 23456, base_folder: "", access_ip: "" };
  }
}, C = (t) => {
  const e = { ...E(), ...t };
  return l.mkdirSync(s.dirname(_()), { recursive: !0 }), l.writeFileSync(_(), JSON.stringify(e, null, 2), "utf-8"), e;
}, y = (t) => new Promise((e) => {
  const n = q.createServer();
  n.once("error", () => e(!1)), n.once("listening", () => n.close(() => e(!0))), n.listen(t, "0.0.0.0");
}), J = () => {
  const t = D.networkInterfaces(), e = [];
  return Object.values(t).forEach((n) => {
    (n || []).forEach((o) => {
      o.family === "IPv4" && !o.internal && e.push(o.address);
    });
  }), e;
}, Y = (t, e) => {
  const n = s.join(process.resourcesPath, "config.yml");
  if (!l.existsSync(n)) return;
  const o = e.replace(/'/g, "''"), w = l.readFileSync(n, "utf-8").replace(/(^\s*Port:\s*).*/m, `$1${t}`).replace(/(^\s*FilePath:\s*).*/m, `$1'${o}'`).replace(/(^\s*FileName:\s*).*/m, "$1''");
  l.writeFileSync(n, w, "utf-8");
}, X = (t) => {
  const e = s.extname(t).toLowerCase();
  return e === ".html" ? "text/html; charset=utf-8" : e === ".js" || e === ".mjs" ? "text/javascript; charset=utf-8" : e === ".css" ? "text/css; charset=utf-8" : e === ".json" ? "application/json; charset=utf-8" : e === ".svg" ? "image/svg+xml" : e === ".png" ? "image/png" : e === ".jpg" || e === ".jpeg" ? "image/jpeg" : e === ".ico" ? "image/x-icon" : e === ".woff" ? "font/woff" : e === ".woff2" ? "font/woff2" : e === ".ttf" ? "font/ttf" : "application/octet-stream";
}, I = async () => {
  m && (await new Promise((t) => m?.close(() => t())), m = null);
}, L = async (t) => {
  await I();
  const e = b, n = s.join(e, "index.html");
  m = z.createServer((o, c) => {
    const p = o.method || "GET";
    if (p !== "GET" && p !== "HEAD") {
      c.statusCode = 405, c.end();
      return;
    }
    try {
      const u = new URL(o.url || "/", `http://${o.headers.host || "localhost"}`), g = decodeURIComponent(u.pathname || "/").replace(/\\/g, "/"), M = g === "/" ? "/index.html" : g, j = s.resolve(s.join(e, M)), W = s.resolve(e), S = j.startsWith(W) ? j : n, A = l.existsSync(S) && l.statSync(S).isFile() ? S : n, $ = l.readFileSync(A);
      c.setHeader("Content-Type", X(A)), c.statusCode = 200, p === "HEAD" ? c.end() : c.end($);
    } catch {
      c.statusCode = 500, c.end();
    }
  }), await new Promise((o, c) => {
    m?.once("error", c), m?.listen(t, "0.0.0.0", () => o());
  });
};
function O() {
  if (h) return;
  const t = process.platform === "win32" ? "server.exe" : "server", e = process.resourcesPath, n = s.join(e, t);
  l.existsSync(n) ? (console.log(`Starting server from: ${n}`), f && (f.kill(), f = null), f = H(n, [], {
    cwd: e,
    // 设置工作目录为 resources
    windowsHide: !0,
    stdio: "ignore"
    // 忽略输出，避免缓冲区填满挂起
  }), f.on("error", (o) => {
    console.error("Failed to start server:", o);
  }), f.on("close", (o) => {
    console.log(`Server process exited with code ${o}`), f = null;
  })) : console.log("Server binary not found, running in client-only mode.");
}
function N() {
  const t = E();
  t.init_done && (O(), L(t.frontend_port).catch(() => {
  }));
  let e = s.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  l.existsSync(e) || (e = s.join(process.env.VITE_PUBLIC, "images/logo/favicon.png")), v = new P({
    width: 500,
    height: 300,
    transparent: !0,
    frame: !1,
    alwaysOnTop: !0,
    icon: e
  }), v.loadFile(s.join(process.env.VITE_PUBLIC, "loading.html")), i = new P({
    title: "Local Music Player",
    show: !1,
    // 先隐藏主窗口
    icon: e,
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: s.join(k, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    }
  }), i.setMenu(null), i.webContents.on("did-finish-load", () => {
    i?.webContents.send("main-process-message", (/* @__PURE__ */ new Date()).toLocaleString());
  }), i.once("ready-to-show", () => {
    setTimeout(() => {
      v?.destroy(), v = null, i?.show(), i?.focus();
    }, 2e3);
  }), h ? (i.loadURL(h), i.webContents.openDevTools()) : i.loadFile(s.join(b, "index.html")), i.webContents.setWindowOpenHandler(({ url: n }) => (n.startsWith("https:") && U.openExternal(n), { action: "deny" }));
}
function K() {
  if (r) return;
  const { width: t, height: e } = T.getPrimaryDisplay().workAreaSize;
  if (r = new P({
    width: 800,
    height: 120,
    x: (t - 800) / 2,
    y: e - 150,
    frame: !1,
    transparent: !0,
    alwaysOnTop: !0,
    skipTaskbar: !0,
    resizable: !0,
    // 允许调整大小
    webPreferences: {
      preload: s.join(k, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    },
    backgroundColor: "#00000000"
    // Ensure transparency
  }), h) {
    const n = h.endsWith("/") ? h : `${h}/`;
    r.loadURL(`${n}#/desktop-lyric`);
  } else
    r.loadFile(s.join(b, "index.html"), { hash: "desktop-lyric" });
  r.on("closed", () => {
    r = null;
  });
}
d.on("open-desktop-lyric", () => {
  K();
});
d.on("close-desktop-lyric", () => {
  r && r.close();
});
d.on("update-desktop-lyric", (t, e) => {
  r && r.webContents.send("update-lyric", e);
});
d.on("desktop-lyric-control", (t, e) => {
  i && i.webContents.send("player-control", e);
});
d.on("desktop-lyric-move", (t, e) => {
  if (r) {
    const [n, o] = r.getPosition(), { width: c } = T.getPrimaryDisplay().workAreaSize;
    r.getBounds().width;
    const p = 50;
    let u = n;
    e === "left" ? u = n - p : e === "right" && (u = n + p), r.setPosition(u, o);
  }
});
d.on("lock-desktop-lyric", (t, e) => {
  r && (r.setIgnoreMouseEvents(e, { forward: !0 }), e ? (r.setFocusable(!1), r.webContents.send("desktop-lyric-locked", !0)) : (r.setFocusable(!0), r.webContents.send("desktop-lyric-locked", !1)));
});
d.on("unlock-desktop-lyric", () => {
  r && (r.setIgnoreMouseEvents(!1, { forward: !0 }), r.setFocusable(!0), r.webContents.send("desktop-lyric-locked", !1));
});
d.on("update-desktop-lyric-settings", (t, e) => {
  r && r.webContents.send("update-settings", e);
});
d.handle("app-clear-data", async () => {
  try {
    i && await i.webContents.session.clearStorageData();
    const t = a.getPath("userData");
    if (l.existsSync(t)) {
      const e = l.readdirSync(t);
      for (const n of e)
        if (!(n === "Lockfile" || n.startsWith("Singleton") || n === "TransportSecurity"))
          try {
            const o = s.join(t, n);
            l.rmSync(o, { recursive: !0, force: !0 });
          } catch (o) {
            console.warn(`Failed to delete ${n}:`, o.message);
          }
    }
    return f && (f.kill(), f = null), await I(), a.relaunch(), a.exit(0), { success: !0 };
  } catch (t) {
    return console.error("Failed to clear app data:", t), { success: !1, error: t.message };
  }
});
d.handle("app-config-get", async () => E());
d.handle("show-save-dialog", async (t, e) => {
  const { filePath: n, canceled: o } = await R.showSaveDialog({
    ...e,
    filters: [
      { name: "Excel Files", extensions: ["xlsx"] },
      { name: "All Files", extensions: ["*"] }
    ]
  });
  return { filePath: n, canceled: o };
});
d.handle("save-file", async (t, { path: e, data: n }) => {
  try {
    return l.writeFileSync(e, Buffer.from(n)), { success: !0 };
  } catch (o) {
    return { success: !1, error: o.message };
  }
});
d.handle("select-directory", async () => await R.showOpenDialog({
  properties: ["openDirectory", "createDirectory"]
}));
d.handle("get-local-ips", async () => ({ ips: J() }));
d.handle("check-ports", async (t, e) => {
  const n = Number(e?.backendPort) || 0, o = Number(e?.frontendPort) || 0;
  return {
    backendAvailable: n > 0 ? await y(n) : !1,
    frontendAvailable: o > 0 ? await y(o) : !1
  };
});
d.handle("apply-initial-config", async (t, e) => {
  const n = String(e?.mode || ""), o = Number(e?.backendPort) || 12345, c = Number(e?.frontendPort) || 23456, p = String(e?.baseFolderPath || ""), u = String(e?.accessIp || "");
  if (n === "server") {
    const w = await y(o), g = await y(c);
    if (!w) throw new Error("backend port unavailable");
    if (!g) throw new Error("frontend port unavailable");
    return C({
      init_done: !0,
      backend_port: o,
      frontend_port: c,
      base_folder: p,
      access_ip: u
    }), Y(o, p), O(), await L(c), { ok: !0 };
  }
  return C({ init_done: !0 }), { ok: !0 };
});
a.whenReady().then(N);
a.on("window-all-closed", () => {
  i = null, f && (f.kill(), f = null), I().catch(() => {
  }), process.platform !== "darwin" && a.quit();
});
a.on("second-instance", () => {
  i && (i.isMinimized() && i.restore(), i.focus());
});
a.on("activate", () => {
  const t = P.getAllWindows();
  t.length ? t[0].focus() : N();
});
export {
  ie as MAIN_DIST,
  b as RENDERER_DIST,
  h as VITE_DEV_SERVER_URL
};
