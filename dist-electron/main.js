import { app as i, ipcMain as f, dialog as R, BrowserWindow as S, shell as U } from "electron";
import { createRequire as V } from "node:module";
import { fileURLToPath as $ } from "node:url";
import o from "node:path";
import k from "node:os";
import { spawn as M } from "child_process";
import * as W from "node:http";
import * as H from "node:net";
import a from "node:fs";
V(import.meta.url);
const C = o.dirname($(import.meta.url));
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = "true";
i.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill");
process.env.APP_ROOT = o.join(C, "..");
const q = "RLMusic", T = (n) => `${q}-${n}`, g = process.env.VITE_DEV_SERVER_URL, y = process.env.VITE_APP_MODE || (a.existsSync(o.join(process.resourcesPath, process.platform === "win32" ? "server.exe" : "server")) ? "server" : "client"), re = o.join(process.env.APP_ROOT, "dist-electron"), b = o.join(process.env.APP_ROOT, "dist");
process.env.VITE_PUBLIC = g ? o.join(process.env.APP_ROOT, "public") : b;
k.release().startsWith("6.1") && i.disableHardwareAcceleration();
process.platform === "win32" && i.setAppUserModelId(T(y));
i.commandLine.appendSwitch("disable-features", "AutofillServerCommunication,Autofill,PasswordManager");
i.setPath("userData", o.join(i.getPath("appData"), T(y)));
i.requestSingleInstanceLock({ mode: y }) || (i.quit(), process.exit(0));
let s, w, l = null, d = null;
const _ = () => o.join(i.getPath("userData"), "app-config.json"), x = () => {
  try {
    const n = a.readFileSync(_(), "utf-8"), e = JSON.parse(n || "{}");
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
}, I = (n) => {
  const e = { ...x(), ...n };
  return a.mkdirSync(o.dirname(_()), { recursive: !0 }), a.writeFileSync(_(), JSON.stringify(e, null, 2), "utf-8"), e;
}, P = (n) => new Promise((e) => {
  const t = H.createServer();
  t.once("error", () => e(!1)), t.once("listening", () => t.close(() => e(!0))), t.listen(n, "0.0.0.0");
}), G = () => {
  const n = k.networkInterfaces(), e = [];
  return Object.values(n).forEach((t) => {
    (t || []).forEach((r) => {
      r.family === "IPv4" && !r.internal && e.push(r.address);
    });
  }), e;
}, J = (n, e) => {
  const t = o.join(process.resourcesPath, "config.yml");
  if (!a.existsSync(t)) return;
  const r = e.replace(/'/g, "''"), h = a.readFileSync(t, "utf-8").replace(/(^\s*Port:\s*).*/m, `$1${n}`).replace(/(^\s*FilePath:\s*).*/m, `$1'${r}'`).replace(/(^\s*FileName:\s*).*/m, "$1''");
  a.writeFileSync(t, h, "utf-8");
}, Y = (n) => {
  const e = o.extname(n).toLowerCase();
  return e === ".html" ? "text/html; charset=utf-8" : e === ".js" || e === ".mjs" ? "text/javascript; charset=utf-8" : e === ".css" ? "text/css; charset=utf-8" : e === ".json" ? "application/json; charset=utf-8" : e === ".svg" ? "image/svg+xml" : e === ".png" ? "image/png" : e === ".jpg" || e === ".jpeg" ? "image/jpeg" : e === ".ico" ? "image/x-icon" : e === ".woff" ? "font/woff" : e === ".woff2" ? "font/woff2" : e === ".ttf" ? "font/ttf" : "application/octet-stream";
}, E = async () => {
  d && (await new Promise((n) => d?.close(() => n())), d = null);
}, D = async (n) => {
  await E();
  const e = b, t = o.join(e, "index.html");
  d = W.createServer((r, c) => {
    const u = r.method || "GET";
    if (u !== "GET" && u !== "HEAD") {
      c.statusCode = 405, c.end();
      return;
    }
    try {
      const p = new URL(r.url || "/", `http://${r.headers.host || "localhost"}`), m = decodeURIComponent(p.pathname || "/").replace(/\\/g, "/"), O = m === "/" ? "/index.html" : m, A = o.resolve(o.join(e, O)), N = o.resolve(e), v = A.startsWith(N) ? A : t, j = a.existsSync(v) && a.statSync(v).isFile() ? v : t, B = a.readFileSync(j);
      c.setHeader("Content-Type", Y(j)), c.statusCode = 200, u === "HEAD" ? c.end() : c.end(B);
    } catch {
      c.statusCode = 500, c.end();
    }
  }), await new Promise((r, c) => {
    d?.once("error", c), d?.listen(n, "0.0.0.0", () => r());
  });
};
function L() {
  if (g) return;
  const n = process.platform === "win32" ? "server.exe" : "server", e = process.resourcesPath, t = o.join(e, n);
  a.existsSync(t) ? (console.log(`Starting server from: ${t}`), l && (l.kill(), l = null), l = M(t, [], {
    cwd: e,
    // 设置工作目录为 resources
    windowsHide: !0,
    stdio: "ignore"
    // 忽略输出，避免缓冲区填满挂起
  }), l.on("error", (r) => {
    console.error("Failed to start server:", r);
  }), l.on("close", (r) => {
    console.log(`Server process exited with code ${r}`), l = null;
  })) : console.log("Server binary not found, running in client-only mode.");
}
function F() {
  const n = x();
  n.init_done && (L(), D(n.frontend_port).catch(() => {
  }));
  let e = o.join(process.env.VITE_PUBLIC, "images/logo/favicon.ico");
  a.existsSync(e) || (e = o.join(process.env.VITE_PUBLIC, "images/logo/favicon.png")), w = new S({
    width: 500,
    height: 300,
    transparent: !0,
    frame: !1,
    alwaysOnTop: !0,
    icon: e
  }), w.loadFile(o.join(process.env.VITE_PUBLIC, "loading.html")), s = new S({
    title: "RLmusic",
    show: !1,
    // 先隐藏主窗口
    icon: e,
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: o.join(C, "../dist-electron/preload.mjs"),
      nodeIntegration: !0,
      contextIsolation: !0
    }
  }), s.setMenu(null), s.webContents.on("did-finish-load", () => {
    s?.webContents.send("main-process-message", (/* @__PURE__ */ new Date()).toLocaleString());
  }), s.once("ready-to-show", () => {
    setTimeout(() => {
      w?.destroy(), w = null, s?.show(), s?.focus();
    }, 2e3);
  }), g ? (s.loadURL(g), s.webContents.openDevTools()) : s.loadFile(o.join(b, "index.html")), s.webContents.setWindowOpenHandler(({ url: t }) => (t.startsWith("https:") && U.openExternal(t), { action: "deny" }));
}
f.handle("app-clear-data", async () => {
  try {
    s && await s.webContents.session.clearStorageData();
    const n = i.getPath("userData");
    if (a.existsSync(n)) {
      const e = a.readdirSync(n);
      for (const t of e)
        if (!(t === "Lockfile" || t.startsWith("Singleton") || t === "TransportSecurity"))
          try {
            const r = o.join(n, t);
            a.rmSync(r, { recursive: !0, force: !0 });
          } catch (r) {
            console.warn(`Failed to delete ${t}:`, r.message);
          }
    }
    return l && (l.kill(), l = null), await E(), i.relaunch(), i.exit(0), { success: !0 };
  } catch (n) {
    return console.error("Failed to clear app data:", n), { success: !1, error: n.message };
  }
});
f.handle("app-config-get", async () => x());
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
    return Buffer.isBuffer(t) ? r = t : t instanceof ArrayBuffer ? r = Buffer.from(new Uint8Array(t)) : ArrayBuffer.isView(t) ? r = Buffer.from(t.buffer, t.byteOffset, t.byteLength) : r = Buffer.from(t), a.writeFileSync(e, r), { success: !0 };
  } catch (r) {
    return { success: !1, error: r.message };
  }
});
f.handle("select-directory", async () => await R.showOpenDialog({
  properties: ["openDirectory", "createDirectory"]
}));
f.handle("get-local-ips", async () => ({ ips: G() }));
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
    const h = await P(r), m = await P(c);
    if (!h) throw new Error("backend port unavailable");
    if (!m) throw new Error("frontend port unavailable");
    return I({
      init_done: !0,
      backend_port: r,
      frontend_port: c,
      base_folder: u,
      access_ip: p
    }), J(r, u), L(), await D(c), { ok: !0 };
  }
  return I({ init_done: !0 }), { ok: !0 };
});
i.whenReady().then(F);
i.on("window-all-closed", () => {
  s = null, l && (l.kill(), l = null), E().catch(() => {
  }), process.platform !== "darwin" && i.quit();
});
i.on("second-instance", () => {
  s && (s.isMinimized() && s.restore(), s.focus());
});
i.on("activate", () => {
  const n = S.getAllWindows();
  n.length ? n[0].focus() : F();
});
export {
  re as MAIN_DIST,
  b as RENDERER_DIST,
  g as VITE_DEV_SERVER_URL
};
