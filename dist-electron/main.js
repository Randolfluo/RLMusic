import { app as a, ipcMain as f, dialog as R, BrowserWindow as v, screen as T, shell as W } from "electron";
import { createRequire as $ } from "node:module";
import { fileURLToPath as V } from "node:url";
import s from "node:path";
import D from "node:os";
import { spawn as H } from "child_process";
import * as z from "node:http";
import * as q from "node:net";
import l from "node:fs";
$(import.meta.url);
const k = s.dirname(V(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
a.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill");
process.env.APP_ROOT = s.join(k, "..");
const G = "RLMusic", F = (n) => `${G}-${n}`, h = process.env.VITE_DEV_SERVER_URL, x = process.env.VITE_APP_MODE || (l.existsSync(s.join(process.resourcesPath, process.platform === "win32" ? "server.exe" : "server")) ? "server" : "client"), ie = s.join(process.env.APP_ROOT, "dist-electron"), b = s.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = h ? s.join(process.env.APP_ROOT, "public") : b;
D.release().startsWith("6.1") && a.disableHardwareAcceleration();
process.platform === "win32" && a.setAppUserModelId(F(x));
a.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill,PasswordManager");
a.setPath("userData", s.join(a.getPath("appData"), F(x)));
a.requestSingleInstanceLock({ mode: x }) || (a.quit(), process.exit(0));
let i, y, o = null, d = null, m = null;
const _ = () => s.join(a.getPath("userData"), "app-config.json"), E = () => {
  try {
    const n = l.readFileSync(_(), "utf-8"), e = JSON.parse(n || "{}");
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
}, C = (n) => {
  const e = { ...E(), ...n };
  return l.mkdirSync(s.dirname(_()), { recursive: !0 }), l.writeFileSync(_(), JSON.stringify(e, null, 2), "utf-8"), e;
}, P = (n) => new Promise((e) => {
  const t = q.createServer();
  t.once("error", () => e(!1)), t.once("listening", () => t.close(() => e(!0))), t.listen(n, "0.0.0.0");
}), J = () => {
  const n = D.networkInterfaces(), e = [];
  return Object.values(n).forEach((t) => {
    (t || []).forEach((r) => {
      r.family === "IPv4" && !r.internal && e.push(r.address);
    });
  }), e;
}, Y = (n, e) => {
  const t = s.join(process.resourcesPath, "config.yml");
  if (!l.existsSync(t)) return;
  const r = e.replace(/'/g, "''"), w = l.readFileSync(t, "utf-8").replace(/(^\s*Port:\s*).*/m, `$1${n}`).replace(/(^\s*FilePath:\s*).*/m, `$1'${r}'`).replace(/(^\s*FileName:\s*).*/m, "$1''");
  l.writeFileSync(t, w, "utf-8");
}, X = (n) => {
  const e = s.extname(n).toLowerCase();
  return e === ".html" ? "text/html; charset=utf-8" : e === ".js" || e === ".mjs" ? "text/javascript; charset=utf-8" : e === ".css" ? "text/css; charset=utf-8" : e === ".json" ? "application/json; charset=utf-8" : e === ".svg" ? "image/svg+xml" : e === ".png" ? "image/png" : e === ".jpg" || e === ".jpeg" ? "image/jpeg" : e === ".ico" ? "image/x-icon" : e === ".woff" ? "font/woff" : e === ".woff2" ? "font/woff2" : e === ".ttf" ? "font/ttf" : "application/octet-stream";
}, A = async () => {
  m && (await new Promise((n) => m?.close(() => n())), m = null);
}, L = async (n) => {
  await A();
  const e = b, t = s.join(e, "index.html");
  m = z.createServer((r, c) => {
    const u = r.method || "GET";
    if (u !== "GET" && u !== "HEAD") {
      c.statusCode = 405, c.end();
      return;
    }
    try {
      const p = new URL(r.url || "/", `http://${r.headers.host || "localhost"}`), g = decodeURIComponent(p.pathname || "/").replace(/\\/g, "/"), B = g === "/" ? "/index.html" : g, I = s.resolve(s.join(e, B)), M = s.resolve(e), S = I.startsWith(M) ? I : t, j = l.existsSync(S) && l.statSync(S).isFile() ? S : t, U = l.readFileSync(j);
      c.setHeader("Content-Type", X(j)), c.statusCode = 200, u === "HEAD" ? c.end() : c.end(U);
    } catch {
      c.statusCode = 500, c.end();
    }
  }), await new Promise((r, c) => {
    m?.once("error", c), m?.listen(n, "0.0.0.0", () => r());
  });
};
function O() {
  if (h) return;
  const n = process.platform === "win32" ? "server.exe" : "server", e = process.resourcesPath, t = s.join(e, n);
  l.existsSync(t) ? (console.log(`Starting server from: ${t}`), d && (d.kill(), d = null), d = H(t, [], {
    cwd: e,
    // 设置工作目录为 resources
    windowsHide: !0,
    stdio: "ignore"
    // 忽略输出，避免缓冲区填满挂起
  }), d.on("error", (r) => {
    console.error("Failed to start server:", r);
  }), d.on("close", (r) => {
    console.log(`Server process exited with code ${r}`), d = null;
  })) : console.log("Server binary not found, running in client-only mode.");
}
function N() {
  const n = E();
  n.init_done && (O(), L(n.frontend_port).catch(() => {
  }));
  let e = s.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  l.existsSync(e) || (e = s.join(process.env.VITE_PUBLIC, "images/logo/favicon.png")), y = new v({
    width: 500,
    height: 300,
    transparent: !0,
    frame: !1,
    alwaysOnTop: !0,
    icon: e
  }), y.loadFile(s.join(process.env.VITE_PUBLIC, "loading.html")), i = new v({
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
      y?.destroy(), y = null, i?.show(), i?.focus();
    }, 2e3);
  }), h ? (i.loadURL(h), i.webContents.openDevTools()) : i.loadFile(s.join(b, "index.html")), i.webContents.setWindowOpenHandler(({ url: t }) => (t.startsWith("https:") && W.openExternal(t), { action: "deny" }));
}
function K() {
  if (o) return;
  const { width: n, height: e } = T.getPrimaryDisplay().workAreaSize;
  if (o = new v({
    width: 800,
    height: 120,
    x: (n - 800) / 2,
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
    const t = h.endsWith("/") ? h : `${h}/`;
    o.loadURL(`${t}#/desktop-lyric`);
  } else
    o.loadFile(s.join(b, "index.html"), { hash: "desktop-lyric" });
  o.on("closed", () => {
    o = null;
  });
}
f.on("open-desktop-lyric", () => {
  K();
});
f.on("close-desktop-lyric", () => {
  o && o.close();
});
f.on("update-desktop-lyric", (n, e) => {
  o && o.webContents.send("update-lyric", e);
});
f.on("desktop-lyric-control", (n, e) => {
  i && i.webContents.send("player-control", e);
});
f.on("desktop-lyric-move", (n, e) => {
  if (o) {
    const [t, r] = o.getPosition(), { width: c } = T.getPrimaryDisplay().workAreaSize;
    o.getBounds().width;
    const u = 50;
    let p = t;
    e === "left" ? p = t - u : e === "right" && (p = t + u), o.setPosition(p, r);
  }
});
f.on("lock-desktop-lyric", (n, e) => {
  o && (o.setIgnoreMouseEvents(e, { forward: !0 }), e ? (o.setFocusable(!1), o.webContents.send("desktop-lyric-locked", !0)) : (o.setFocusable(!0), o.webContents.send("desktop-lyric-locked", !1)));
});
f.on("unlock-desktop-lyric", () => {
  o && (o.setIgnoreMouseEvents(!1, { forward: !0 }), o.setFocusable(!0), o.webContents.send("desktop-lyric-locked", !1));
});
f.on("update-desktop-lyric-settings", (n, e) => {
  o && o.webContents.send("update-settings", e);
});
f.handle("app-clear-data", async () => {
  try {
    i && await i.webContents.session.clearStorageData();
    const n = a.getPath("userData");
    if (l.existsSync(n)) {
      const e = l.readdirSync(n);
      for (const t of e)
        if (!(t === "Lockfile" || t.startsWith("Singleton") || t === "TransportSecurity"))
          try {
            const r = s.join(n, t);
            l.rmSync(r, { recursive: !0, force: !0 });
          } catch (r) {
            console.warn(`Failed to delete ${t}:`, r.message);
          }
    }
    return d && (d.kill(), d = null), await A(), a.relaunch(), a.exit(0), { success: !0 };
  } catch (n) {
    return console.error("Failed to clear app data:", n), { success: !1, error: n.message };
  }
});
f.handle("app-config-get", async () => E());
f.handle("show-save-dialog", async (n, e) => {
  const { filePath: t, canceled: r } = await R.showSaveDialog({
    ...e,
    filters: [
      { name: "Excel Files", extensions: ["xlsx"] },
      { name: "All Files", extensions: ["*"] }
    ]
  });
  return { filePath: t, canceled: r };
});
f.handle("save-file", async (n, { path: e, data: t }) => {
  try {
    let r;
    return Buffer.isBuffer(t) ? r = t : t instanceof ArrayBuffer ? r = Buffer.from(new Uint8Array(t)) : ArrayBuffer.isView(t) ? r = Buffer.from(t.buffer, t.byteOffset, t.byteLength) : r = Buffer.from(t), l.writeFileSync(e, r), { success: !0 };
  } catch (r) {
    return { success: !1, error: r.message };
  }
});
f.handle("select-directory", async () => await R.showOpenDialog({
  properties: ["openDirectory", "createDirectory"]
}));
f.handle("get-local-ips", async () => ({ ips: J() }));
f.handle("check-ports", async (n, e) => {
  const t = Number(e?.backendPort) || 0, r = Number(e?.frontendPort) || 0;
  return {
    backendAvailable: t > 0 ? await P(t) : !1,
    frontendAvailable: r > 0 ? await P(r) : !1
  };
});
f.handle("apply-initial-config", async (n, e) => {
  const t = String(e?.mode || ""), r = Number(e?.backendPort) || 12345, c = Number(e?.frontendPort) || 23456, u = String(e?.baseFolderPath || ""), p = String(e?.accessIp || "");
  if (t === "server") {
    const w = await P(r), g = await P(c);
    if (!w) throw new Error("backend port unavailable");
    if (!g) throw new Error("frontend port unavailable");
    return C({
      init_done: !0,
      backend_port: r,
      frontend_port: c,
      base_folder: u,
      access_ip: p
    }), Y(r, u), O(), await L(c), { ok: !0 };
  }
  return C({ init_done: !0 }), { ok: !0 };
});
a.whenReady().then(N);
a.on("window-all-closed", () => {
  i = null, d && (d.kill(), d = null), A().catch(() => {
  }), process.platform !== "darwin" && a.quit();
});
a.on("second-instance", () => {
  i && (i.isMinimized() && i.restore(), i.focus());
});
a.on("activate", () => {
  const n = v.getAllWindows();
  n.length ? n[0].focus() : N();
});
export {
  ie as MAIN_DIST,
  b as RENDERER_DIST,
  h as VITE_DEV_SERVER_URL
};
