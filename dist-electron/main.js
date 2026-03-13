import { app as a, ipcMain as d, dialog as U, BrowserWindow as P, screen as R, shell as $ } from "electron";
import { createRequire as V } from "node:module";
import { fileURLToPath as B } from "node:url";
import s from "node:path";
import T from "node:os";
import { spawn as H } from "child_process";
import * as z from "node:http";
import * as q from "node:net";
import l from "node:fs";
V(import.meta.url);
const _ = s.dirname(B(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
a.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill");
process.env.APP_ROOT = s.join(_, "..");
const G = "RLMusic", D = (t) => `${G}-${t}`, h = process.env.VITE_DEV_SERVER_URL, k = process.env.VITE_APP_MODE || (l.existsSync(s.join(process.resourcesPath, process.platform === "win32" ? "server.exe" : "server")) ? "server" : "client"), ie = s.join(process.env.APP_ROOT, "dist-electron"), y = s.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = h ? s.join(process.env.APP_ROOT, "public") : y;
T.release().startsWith("6.1") && a.disableHardwareAcceleration();
process.platform === "win32" && a.setAppUserModelId(D(k));
a.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill,PasswordManager");
a.setPath("userData", s.join(a.getPath("appData"), D(k)));
a.requestSingleInstanceLock(k) || (a.quit(), process.exit(0));
let i, g, r = null, f = null, m = null;
const S = () => s.join(a.getPath("userData"), "app-config.json"), E = () => {
  try {
    const t = l.readFileSync(S(), "utf-8"), e = JSON.parse(t || "{}");
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
}, A = (t) => {
  const e = { ...E(), ...t };
  return l.mkdirSync(s.dirname(S()), { recursive: !0 }), l.writeFileSync(S(), JSON.stringify(e, null, 2), "utf-8"), e;
}, v = (t) => new Promise((e) => {
  const n = q.createServer();
  n.once("error", () => e(!1)), n.once("listening", () => n.close(() => e(!0))), n.listen(t, "0.0.0.0");
}), J = () => {
  const t = T.networkInterfaces(), e = [];
  return Object.values(t).forEach((n) => {
    (n || []).forEach((o) => {
      o.family === "IPv4" && !o.internal && e.push(o.address);
    });
  }), e;
}, Y = (t, e) => {
  const n = s.join(process.resourcesPath, "config.yml");
  if (!l.existsSync(n)) return;
  const o = e.replace(/'/g, "''"), u = l.readFileSync(n, "utf-8").replace(/(^\s*Port:\s*).*/m, `$1:${t}`).replace(/(^\s*FilePath:\s*).*/m, `$1'${o}'`);
  l.writeFileSync(n, u, "utf-8");
}, X = (t) => {
  const e = s.extname(t).toLowerCase();
  return e === ".html" ? "text/html; charset=utf-8" : e === ".js" || e === ".mjs" ? "text/javascript; charset=utf-8" : e === ".css" ? "text/css; charset=utf-8" : e === ".json" ? "application/json; charset=utf-8" : e === ".svg" ? "image/svg+xml" : e === ".png" ? "image/png" : e === ".jpg" || e === ".jpeg" ? "image/jpeg" : e === ".ico" ? "image/x-icon" : e === ".woff" ? "font/woff" : e === ".woff2" ? "font/woff2" : e === ".ttf" ? "font/ttf" : "application/octet-stream";
}, x = async () => {
  m && (await new Promise((t) => m?.close(() => t())), m = null);
}, L = async (t) => {
  await x();
  const e = y, n = s.join(e, "index.html");
  m = z.createServer((o, c) => {
    const p = o.method || "GET";
    if (p !== "GET" && p !== "HEAD") {
      c.statusCode = 405, c.end();
      return;
    }
    try {
      const u = new URL(o.url || "/", `http://${o.headers.host || "localhost"}`), w = decodeURIComponent(u.pathname || "/").replace(/\\/g, "/"), N = w === "/" ? "/index.html" : w, j = s.resolve(s.join(e, N)), M = s.resolve(e), b = j.startsWith(M) ? j : n, C = l.existsSync(b) && l.statSync(b).isFile() ? b : n, W = l.readFileSync(C);
      c.setHeader("Content-Type", X(C)), c.statusCode = 200, p === "HEAD" ? c.end() : c.end(W);
    } catch {
      c.statusCode = 500, c.end();
    }
  }), await new Promise((o, c) => {
    m?.once("error", c), m?.listen(t, "0.0.0.0", () => o());
  });
};
function F() {
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
function O() {
  const t = E();
  t.init_done && (F(), L(t.frontend_port).catch(() => {
  }));
  let e = s.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  l.existsSync(e) || (e = s.join(process.env.VITE_PUBLIC, "images/logo/favicon.png")), g = new P({
    width: 500,
    height: 300,
    transparent: !0,
    frame: !1,
    alwaysOnTop: !0,
    icon: e
  }), g.loadFile(s.join(process.env.VITE_PUBLIC, "loading.html")), i = new P({
    title: "Local Music Player",
    show: !1,
    // 先隐藏主窗口
    icon: e,
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: s.join(_, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    }
  }), i.setMenu(null), i.webContents.on("did-finish-load", () => {
    i?.webContents.send("main-process-message", (/* @__PURE__ */ new Date()).toLocaleString());
  }), i.once("ready-to-show", () => {
    setTimeout(() => {
      g?.destroy(), g = null, i?.show(), i?.focus();
    }, 2e3);
  }), h ? (i.loadURL(h), i.webContents.openDevTools()) : i.loadFile(s.join(y, "index.html")), i.webContents.setWindowOpenHandler(({ url: n }) => (n.startsWith("https:") && $.openExternal(n), { action: "deny" }));
}
function K() {
  if (r) return;
  const { width: t, height: e } = R.getPrimaryDisplay().workAreaSize;
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
      preload: s.join(_, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    },
    backgroundColor: "#00000000"
    // Ensure transparency
  }), h) {
    const n = h.endsWith("/") ? h : `${h}/`;
    r.loadURL(`${n}#/desktop-lyric`);
  } else
    r.loadFile(s.join(y, "index.html"), { hash: "desktop-lyric" });
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
    const [n, o] = r.getPosition(), { width: c } = R.getPrimaryDisplay().workAreaSize;
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
    return f && (f.kill(), f = null), await x(), a.relaunch(), a.exit(0), { success: !0 };
  } catch (t) {
    return console.error("Failed to clear app data:", t), { success: !1, error: t.message };
  }
});
d.handle("app-config-get", async () => E());
d.handle("select-directory", async () => await U.showOpenDialog({
  properties: ["openDirectory", "createDirectory"]
}));
d.handle("get-local-ips", async () => ({ ips: J() }));
d.handle("check-ports", async (t, e) => {
  const n = Number(e?.backendPort) || 0, o = Number(e?.frontendPort) || 0;
  return {
    backendAvailable: n > 0 ? await v(n) : !1,
    frontendAvailable: o > 0 ? await v(o) : !1
  };
});
d.handle("apply-initial-config", async (t, e) => {
  const n = String(e?.mode || ""), o = Number(e?.backendPort) || 12345, c = Number(e?.frontendPort) || 23456, p = String(e?.baseFolderPath || ""), u = String(e?.accessIp || "");
  if (n === "server") {
    const I = await v(o), w = await v(c);
    if (!I) throw new Error("backend port unavailable");
    if (!w) throw new Error("frontend port unavailable");
    return A({
      init_done: !0,
      backend_port: o,
      frontend_port: c,
      base_folder: p,
      access_ip: u
    }), Y(o, p), F(), await L(c), { ok: !0 };
  }
  return A({ init_done: !0 }), { ok: !0 };
});
a.whenReady().then(O);
a.on("window-all-closed", () => {
  i = null, f && (f.kill(), f = null), x().catch(() => {
  }), process.platform !== "darwin" && a.quit();
});
a.on("second-instance", () => {
  i && (i.isMinimized() && i.restore(), i.focus());
});
a.on("activate", () => {
  const t = P.getAllWindows();
  t.length ? t[0].focus() : O();
});
export {
  ie as MAIN_DIST,
  y as RENDERER_DIST,
  h as VITE_DEV_SERVER_URL
};
