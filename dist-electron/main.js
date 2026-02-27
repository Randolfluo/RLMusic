import { app as c, ipcMain as d, dialog as N, BrowserWindow as v, screen as U, shell as W } from "electron";
import { createRequire as $ } from "node:module";
import { fileURLToPath as V } from "node:url";
import o from "node:path";
import A from "node:os";
import { spawn as H } from "child_process";
import * as B from "node:http";
import * as q from "node:net";
import l from "node:fs";
$(import.meta.url);
const k = o.dirname(V(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
c.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill");
process.env.APP_ROOT = o.join(k, "..");
const f = process.env.VITE_DEV_SERVER_URL, S = process.env.VITE_APP_MODE || (l.existsSync(o.join(process.resourcesPath, process.platform === "win32" ? "server.exe" : "server")) ? "server" : "client"), re = o.join(process.env.APP_ROOT, "dist-electron"), y = o.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = f ? o.join(process.env.APP_ROOT, "public") : y;
A.release().startsWith("6.1") && c.disableHardwareAcceleration();
process.platform === "win32" && c.setAppUserModelId(`LocalMusicPlayer-${S}`);
c.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill,PasswordManager");
c.setPath("userData", o.join(c.getPath("appData"), `LocalMusicPlayer-${S}`));
c.requestSingleInstanceLock(S) || (c.quit(), process.exit(0));
let i, g, s = null, p = null, h = null;
const _ = () => o.join(c.getPath("userData"), "app-config.json"), E = () => {
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
  return l.mkdirSync(o.dirname(_()), { recursive: !0 }), l.writeFileSync(_(), JSON.stringify(e, null, 2), "utf-8"), e;
}, P = (t) => new Promise((e) => {
  const n = q.createServer();
  n.once("error", () => e(!1)), n.once("listening", () => n.close(() => e(!0))), n.listen(t, "0.0.0.0");
}), z = () => {
  const t = A.networkInterfaces(), e = [];
  return Object.values(t).forEach((n) => {
    (n || []).forEach((r) => {
      r.family === "IPv4" && !r.internal && e.push(r.address);
    });
  }), e;
}, G = (t, e) => {
  const n = o.join(process.resourcesPath, "config.yml");
  if (!l.existsSync(n)) return;
  const r = e.replace(/'/g, "''"), m = l.readFileSync(n, "utf-8").replace(/(^\s*Port:\s*).*/m, `$1:${t}`).replace(/(^\s*FilePath:\s*).*/m, `$1'${r}'`);
  l.writeFileSync(n, m, "utf-8");
}, J = (t) => {
  const e = o.extname(t).toLowerCase();
  return e === ".html" ? "text/html; charset=utf-8" : e === ".js" || e === ".mjs" ? "text/javascript; charset=utf-8" : e === ".css" ? "text/css; charset=utf-8" : e === ".json" ? "application/json; charset=utf-8" : e === ".svg" ? "image/svg+xml" : e === ".png" ? "image/png" : e === ".jpg" || e === ".jpeg" ? "image/jpeg" : e === ".ico" ? "image/x-icon" : e === ".woff" ? "font/woff" : e === ".woff2" ? "font/woff2" : e === ".ttf" ? "font/ttf" : "application/octet-stream";
}, R = async () => {
  h && (await new Promise((t) => h?.close(() => t())), h = null);
}, T = async (t) => {
  await R();
  const e = y, n = o.join(e, "index.html");
  h = B.createServer((r, a) => {
    const u = r.method || "GET";
    if (u !== "GET" && u !== "HEAD") {
      a.statusCode = 405, a.end();
      return;
    }
    try {
      const m = new URL(r.url || "/", `http://${r.headers.host || "localhost"}`), w = decodeURIComponent(m.pathname || "/").replace(/\\/g, "/"), O = w === "/" ? "/index.html" : w, j = o.resolve(o.join(e, O)), F = o.resolve(e), b = j.startsWith(F) ? j : n, x = l.existsSync(b) && l.statSync(b).isFile() ? b : n, M = l.readFileSync(x);
      a.setHeader("Content-Type", J(x)), a.statusCode = 200, u === "HEAD" ? a.end() : a.end(M);
    } catch {
      a.statusCode = 500, a.end();
    }
  }), await new Promise((r, a) => {
    h?.once("error", a), h?.listen(t, "0.0.0.0", () => r());
  });
};
function L() {
  if (f) return;
  const t = process.platform === "win32" ? "server.exe" : "server", e = process.resourcesPath, n = o.join(e, t);
  l.existsSync(n) ? (console.log(`Starting server from: ${n}`), p && (p.kill(), p = null), p = H(n, [], {
    cwd: e,
    // 设置工作目录为 resources
    windowsHide: !0,
    stdio: "ignore"
    // 忽略输出，避免缓冲区填满挂起
  }), p.on("error", (r) => {
    console.error("Failed to start server:", r);
  }), p.on("close", (r) => {
    console.log(`Server process exited with code ${r}`), p = null;
  })) : console.log("Server binary not found, running in client-only mode.");
}
function D() {
  const t = E();
  t.init_done && (L(), T(t.frontend_port).catch(() => {
  }));
  let e = o.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  l.existsSync(e) || (e = o.join(process.env.VITE_PUBLIC, "images/logo/favicon.png")), g = new v({
    width: 500,
    height: 300,
    transparent: !0,
    frame: !1,
    alwaysOnTop: !0,
    icon: e
  }), g.loadFile(o.join(process.env.VITE_PUBLIC, "loading.html")), i = new v({
    title: "Local Music Player",
    show: !1,
    // 先隐藏主窗口
    icon: e,
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: o.join(k, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    }
  }), i.setMenu(null), i.webContents.on("did-finish-load", () => {
    i?.webContents.send("main-process-message", (/* @__PURE__ */ new Date()).toLocaleString());
  }), i.once("ready-to-show", () => {
    setTimeout(() => {
      g?.destroy(), g = null, i?.show(), i?.focus();
    }, 2e3);
  }), f ? (i.loadURL(f), i.webContents.openDevTools()) : i.loadFile(o.join(y, "index.html")), i.webContents.setWindowOpenHandler(({ url: n }) => (n.startsWith("https:") && W.openExternal(n), { action: "deny" }));
}
function Y() {
  if (s) return;
  const { width: t, height: e } = U.getPrimaryDisplay().workAreaSize;
  if (s = new v({
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
      preload: o.join(k, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    },
    backgroundColor: "#00000000"
    // Ensure transparency
  }), f) {
    const n = f.endsWith("/") ? f : `${f}/`;
    s.loadURL(`${n}#/desktop-lyric`);
  } else
    s.loadFile(o.join(y, "index.html"), { hash: "desktop-lyric" });
  s.on("closed", () => {
    s = null;
  });
}
d.on("open-desktop-lyric", () => {
  Y();
});
d.on("close-desktop-lyric", () => {
  s && s.close();
});
d.on("update-desktop-lyric", (t, e) => {
  s && s.webContents.send("update-lyric", e);
});
d.on("desktop-lyric-control", (t, e) => {
  i && i.webContents.send("player-control", e);
});
d.on("lock-desktop-lyric", (t, e) => {
  s && (s.setIgnoreMouseEvents(e, { forward: !0 }), e ? (s.setFocusable(!1), s.webContents.send("desktop-lyric-locked", !0)) : (s.setFocusable(!0), s.webContents.send("desktop-lyric-locked", !1)));
});
d.on("unlock-desktop-lyric", () => {
  s && (s.setIgnoreMouseEvents(!1, { forward: !0 }), s.setFocusable(!0), s.webContents.send("desktop-lyric-locked", !1));
});
d.on("update-desktop-lyric-settings", (t, e) => {
  s && s.webContents.send("update-settings", e);
});
d.handle("app-config-get", async () => E());
d.handle("select-directory", async () => await N.showOpenDialog({
  properties: ["openDirectory", "createDirectory"]
}));
d.handle("get-local-ips", async () => ({ ips: z() }));
d.handle("check-ports", async (t, e) => {
  const n = Number(e?.backendPort) || 0, r = Number(e?.frontendPort) || 0;
  return {
    backendAvailable: n > 0 ? await P(n) : !1,
    frontendAvailable: r > 0 ? await P(r) : !1
  };
});
d.handle("apply-initial-config", async (t, e) => {
  const n = String(e?.mode || ""), r = Number(e?.backendPort) || 12345, a = Number(e?.frontendPort) || 23456, u = String(e?.baseFolderPath || ""), m = String(e?.accessIp || "");
  if (n === "server") {
    const I = await P(r), w = await P(a);
    if (!I) throw new Error("backend port unavailable");
    if (!w) throw new Error("frontend port unavailable");
    return C({
      init_done: !0,
      backend_port: r,
      frontend_port: a,
      base_folder: u,
      access_ip: m
    }), G(r, u), L(), await T(a), { ok: !0 };
  }
  return C({ init_done: !0 }), { ok: !0 };
});
c.whenReady().then(D);
c.on("window-all-closed", () => {
  i = null, p && (p.kill(), p = null), R().catch(() => {
  }), process.platform !== "darwin" && c.quit();
});
c.on("second-instance", () => {
  i && (i.isMinimized() && i.restore(), i.focus());
});
c.on("activate", () => {
  const t = v.getAllWindows();
  t.length ? t[0].focus() : D();
});
export {
  re as MAIN_DIST,
  y as RENDERER_DIST,
  f as VITE_DEV_SERVER_URL
};
