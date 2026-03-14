import { app, BrowserWindow, shell, ipcMain, screen, dialog } from 'electron'
import { createRequire } from 'node:module'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import os from 'node:os'
import { spawn } from 'child_process'
import * as http from 'node:http'
import * as net from 'node:net'
import fs from 'node:fs'

const require = createRequire(import.meta.url)
const __dirname = path.dirname(fileURLToPath(import.meta.url))

// 屏蔽 Electron 的安全警告
// 这个警告通常出现在开发环境中，因为 Vite 等构建工具需要使用 'unsafe-eval'
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = 'true'

// 屏蔽 Autofill 相关的终端报错
app.commandLine.appendSwitch('disable-features', 'AutofillServerCommunication,Autofill')

// 目录结构说明
//
// ├─┬ dist-electron
// │ ├─┬ main
// │ │ └── index.js    > 主进程入口
// │ └─┬ preload
// │   └── index.mjs   > 预加载脚本
// ├─┬ dist
// │ └── index.html    > 渲染进程入口
//
process.env.APP_ROOT = path.join(__dirname, '..')

// 定义应用模式和名称常量
const APP_NAME_BASE = 'RLMusic'
const getAppId = (mode: string) => `${APP_NAME_BASE}-${mode}`

export const VITE_DEV_SERVER_URL = process.env['VITE_DEV_SERVER_URL']
const appMode = process.env.VITE_APP_MODE || (fs.existsSync(path.join(process.resourcesPath, process.platform === 'win32' ? 'server.exe' : 'server')) ? 'server' : 'client')
export const MAIN_DIST = path.join(process.env.APP_ROOT, 'dist-electron')
export const RENDERER_DIST = path.join(process.env.APP_ROOT, 'dist')

process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL ? path.join(process.env.APP_ROOT, 'public') : RENDERER_DIST

// Win7 禁用 GPU 加速
if (os.release().startsWith('6.1')) app.disableHardwareAcceleration()

if (process.platform === 'win32') app.setAppUserModelId(getAppId(appMode))

// 屏蔽 Electron 常见的 Autofill 相关的终端报错
// "Request Autofill.enable failed", "Request Autofill.setAddresses failed"
app.commandLine.appendSwitch('disable-features', 'AutofillServerCommunication,Autofill,PasswordManager')

// 设置用户数据目录，使用统一的应用名称
app.setPath('userData', path.join(app.getPath('appData'), getAppId(appMode)))
if (!app.requestSingleInstanceLock(appMode)) {
  app.quit()
  process.exit(0)
}

let win: BrowserWindow | null
let splash: BrowserWindow | null
let desktopLyricWindow: BrowserWindow | null = null
let serverProcess: any = null // 保存后端进程实例
let frontendServer: http.Server | null = null

type AppConfig = {
  init_done: boolean
  backend_port: number
  frontend_port: number
  base_folder: string
  access_ip: string
}

const getAppConfigPath = () => path.join(app.getPath('userData'), 'app-config.json')

const readAppConfig = (): AppConfig => {
  try {
    const raw = fs.readFileSync(getAppConfigPath(), 'utf-8')
    const parsed = JSON.parse(raw || '{}')
    return {
      init_done: !!parsed.init_done,
      backend_port: Number(parsed.backend_port) || 12345,
      frontend_port: Number(parsed.frontend_port) || 23456,
      base_folder: String(parsed.base_folder || ''),
      access_ip: String(parsed.access_ip || ''),
    }
  } catch {
    return { init_done: false, backend_port: 12345, frontend_port: 23456, base_folder: '', access_ip: '' }
  }
}

const writeAppConfig = (cfg: Partial<AppConfig>) => {
  const next = { ...readAppConfig(), ...cfg }
  fs.mkdirSync(path.dirname(getAppConfigPath()), { recursive: true })
  fs.writeFileSync(getAppConfigPath(), JSON.stringify(next, null, 2), 'utf-8')
  return next
}

const isPortAvailable = (port: number) =>
  new Promise<boolean>((resolve) => {
    const server = net.createServer()
    server.once('error', () => resolve(false))
    server.once('listening', () => server.close(() => resolve(true)))
    server.listen(port, '0.0.0.0')
  })

const getLocalIPs = () => {
  const nets = os.networkInterfaces()
  const ips: string[] = []
  Object.values(nets).forEach((items) => {
    (items || []).forEach((n) => {
      if (n.family === 'IPv4' && !n.internal) ips.push(n.address)
    })
  })
  return ips
}

const updateServerConfigYml = (backendPort: number, baseFolderPath: string) => {
  const configPath = path.join(process.resourcesPath, 'config.yml')
  if (!fs.existsSync(configPath)) return
  const escapedPath = baseFolderPath.replace(/'/g, "''")
  const raw = fs.readFileSync(configPath, 'utf-8')
  const withPort = raw.replace(/(^\s*Port:\s*).*/m, `$1${backendPort}`)
  const withPath = withPort.replace(/(^\s*FilePath:\s*).*/m, `$1'${escapedPath}'`)
  const withName = withPath.replace(/(^\s*FileName:\s*).*/m, `$1''`)
  fs.writeFileSync(configPath, withName, 'utf-8')
}

const guessMime = (p: string) => {
  const ext = path.extname(p).toLowerCase()
  if (ext === '.html') return 'text/html; charset=utf-8'
  if (ext === '.js' || ext === '.mjs') return 'text/javascript; charset=utf-8'
  if (ext === '.css') return 'text/css; charset=utf-8'
  if (ext === '.json') return 'application/json; charset=utf-8'
  if (ext === '.svg') return 'image/svg+xml'
  if (ext === '.png') return 'image/png'
  if (ext === '.jpg' || ext === '.jpeg') return 'image/jpeg'
  if (ext === '.ico') return 'image/x-icon'
  if (ext === '.woff') return 'font/woff'
  if (ext === '.woff2') return 'font/woff2'
  if (ext === '.ttf') return 'font/ttf'
  return 'application/octet-stream'
}

const stopFrontendServer = async () => {
  if (!frontendServer) return
  await new Promise<void>((resolve) => frontendServer?.close(() => resolve()))
  frontendServer = null
}

const startFrontendServer = async (port: number) => {
  await stopFrontendServer()
  const root = RENDERER_DIST
  const indexPath = path.join(root, 'index.html')
  frontendServer = http.createServer((req, res) => {
    const method = req.method || 'GET'
    if (method !== 'GET' && method !== 'HEAD') {
      res.statusCode = 405
      res.end()
      return
    }
    try {
      const url = new URL(req.url || '/', `http://${req.headers.host || 'localhost'}`)
      const pathname = decodeURIComponent(url.pathname || '/')
      const safePath = pathname.replace(/\\/g, '/')
      const requested = safePath === '/' ? '/index.html' : safePath
      const resolved = path.resolve(path.join(root, requested))
      const rootResolved = path.resolve(root)
      const targetPath = resolved.startsWith(rootResolved) ? resolved : indexPath
      const exists = fs.existsSync(targetPath) && fs.statSync(targetPath).isFile()
      const finalPath = exists ? targetPath : indexPath
      const data = fs.readFileSync(finalPath)
      res.setHeader('Content-Type', guessMime(finalPath))
      res.statusCode = 200
      if (method === 'HEAD') {
        res.end()
      } else {
        res.end(data)
      }
    } catch {
      res.statusCode = 500
      res.end()
    }
  })
  await new Promise<void>((resolve, reject) => {
    frontendServer?.once('error', reject)
    frontendServer?.listen(port, '0.0.0.0', () => resolve())
  })
}

// 启动后端服务
function startServer() {
  // 仅在生产环境尝试启动
  if (VITE_DEV_SERVER_URL) return;

  const serverName = process.platform === 'win32' ? 'server.exe' : 'server';
  // resources 目录路径 (Electron 打包后资源目录)
  const resourcesPath = process.resourcesPath;
  const serverPath = path.join(resourcesPath, serverName);
  
  if (fs.existsSync(serverPath)) {
    console.log(`Starting server from: ${serverPath}`);
    // 启动服务，不显示窗口
    if (serverProcess) {
      serverProcess.kill()
      serverProcess = null
    }
    serverProcess = spawn(serverPath, [], {
      cwd: resourcesPath, // 设置工作目录为 resources
      windowsHide: true,
      stdio: 'ignore' // 忽略输出，避免缓冲区填满挂起
    });

    serverProcess.on('error', (err: any) => {
      console.error('Failed to start server:', err);
    });

    serverProcess.on('close', (code: any) => {
      console.log(`Server process exited with code ${code}`);
      serverProcess = null;
    });
  } else {
    console.log('Server binary not found, running in client-only mode.');
  }
}

/**
 * 创建主窗口
 */
function createWindow() {
  const cfg = readAppConfig()
  if (cfg.init_done) {
    startServer()
    startFrontendServer(cfg.frontend_port).catch(() => {})
  }

  // 优先查找 .ico 文件 (Windows 最佳实践)，如果不存在则使用 .png
  let iconPath = path.join(process.env.VITE_PUBLIC as string, 'images/logo/favicon.ico')
  if (!fs.existsSync(iconPath)) {
    iconPath = path.join(process.env.VITE_PUBLIC as string, 'images/logo/favicon.png')
  }

  // 创建启动页
  splash = new BrowserWindow({
    width: 500,
    height: 300,
    transparent: true,
    frame: false,
    alwaysOnTop: true,
    icon: iconPath,
  })
  splash.loadFile(path.join(process.env.VITE_PUBLIC as string, 'loading.html'))

  win = new BrowserWindow({
    title: 'Local Music Player',
    show: false, // 先隐藏主窗口
    icon: iconPath,
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: path.join(__dirname, '../dist-electron/preload.mjs'), 
      nodeIntegration: true,
      contextIsolation: true,
    },
  })
  
  // 隐藏菜单栏 (File, Edit, etc.)
  win.setMenu(null)

  // 页面加载完成后，发送当前时间给渲染进程（测试用）
  win.webContents.on('did-finish-load', () => {
    win?.webContents.send('main-process-message', (new Date).toLocaleString())
  })

  // 等待页面加载完成 (ready-to-show) 后再显示主窗口并关闭启动页
  win.once('ready-to-show', () => {
    setTimeout(() => {
      splash?.destroy()
      splash = null
      win?.show()
      win?.focus()
    }, 2000)
  })

  if (VITE_DEV_SERVER_URL) {
    win.loadURL(VITE_DEV_SERVER_URL)
    // 开发模式下打开开发者工具
    win.webContents.openDevTools()
  } else {
    win.loadFile(path.join(RENDERER_DIST, 'index.html'))
  }

  // 让所有 https 链接在默认浏览器打开，而不是在应用内
  win.webContents.setWindowOpenHandler(({ url }: { url: string }) => {
    if (url.startsWith('https:')) shell.openExternal(url)
    return { action: 'deny' }
  })
}

// 创建桌面歌词窗口
function createDesktopLyricWindow() {
  if (desktopLyricWindow) return

  const { width, height } = screen.getPrimaryDisplay().workAreaSize

  desktopLyricWindow = new BrowserWindow({
    width: 800,
    height: 120,
    x: (width - 800) / 2,
    y: height - 150,
    frame: false,
    transparent: true,
    alwaysOnTop: true,
    skipTaskbar: true,
    resizable: true, // 允许调整大小
    webPreferences: {
      preload: path.join(__dirname, '../dist-electron/preload.mjs'),
      nodeIntegration: true,
      contextIsolation: true,
    },
    backgroundColor: '#00000000', // Ensure transparency
  })

  // 加载路由
  if (VITE_DEV_SERVER_URL) {
    const url = VITE_DEV_SERVER_URL.endsWith('/') ? VITE_DEV_SERVER_URL : `${VITE_DEV_SERVER_URL}/`
    desktopLyricWindow.loadURL(`${url}#/desktop-lyric`)
  } else {
    // 生产环境下加载 hash 路由
    desktopLyricWindow.loadFile(path.join(RENDERER_DIST, 'index.html'), { hash: 'desktop-lyric' })
  }

  desktopLyricWindow.on('closed', () => {
    desktopLyricWindow = null
  })
}

// IPC 监听
ipcMain.on('open-desktop-lyric', () => {
  createDesktopLyricWindow()
})

ipcMain.on('close-desktop-lyric', () => {
  if (desktopLyricWindow) {
    desktopLyricWindow.close()
  }
})

ipcMain.on('update-desktop-lyric', (event, data) => {
  if (desktopLyricWindow) {
    desktopLyricWindow.webContents.send('update-lyric', data)
  }
})

ipcMain.on('desktop-lyric-control', (event, action) => {
  if (win) {
    win.webContents.send('player-control', action)
  }
})

ipcMain.on('desktop-lyric-move', (event, direction) => {
  if (desktopLyricWindow) {
    const [x, y] = desktopLyricWindow.getPosition()
    const { width } = screen.getPrimaryDisplay().workAreaSize
    const winWidth = desktopLyricWindow.getBounds().width
    
    // 每次移动 50px
    const step = 50
    let newX = x

    if (direction === 'left') {
      newX = x - step
    } else if (direction === 'right') {
      newX = x + step
    }

    // 边界检查（可选）
    // if (newX < 0) newX = 0
    // if (newX + winWidth > width) newX = width - winWidth

    desktopLyricWindow.setPosition(newX, y)
  }
})

ipcMain.on('lock-desktop-lyric', (event, locked) => {
  if (desktopLyricWindow) {
    desktopLyricWindow.setIgnoreMouseEvents(locked, { forward: true })
    if (locked) {
      desktopLyricWindow.setFocusable(false)
      // 告诉渲染进程已锁定
      desktopLyricWindow.webContents.send('desktop-lyric-locked', true)
    } else {
      desktopLyricWindow.setFocusable(true)
      // 告诉渲染进程已解锁
      desktopLyricWindow.webContents.send('desktop-lyric-locked', false)
    }
  }
})

// 添加解锁监听
ipcMain.on('unlock-desktop-lyric', () => {
  if (desktopLyricWindow) {
    desktopLyricWindow.setIgnoreMouseEvents(false, { forward: true })
    desktopLyricWindow.setFocusable(true)
    desktopLyricWindow.webContents.send('desktop-lyric-locked', false)
  }
})

ipcMain.on('update-desktop-lyric-settings', (event, settings) => {
  if (desktopLyricWindow) {
    desktopLyricWindow.webContents.send('update-settings', settings)
  }
})

// 清除应用所有数据
ipcMain.handle('app-clear-data', async () => {
  try {
    // 1. 清除 Session 数据（缓存、Cookie、LocalStorage 等）
    if (win) {
      await win.webContents.session.clearStorageData()
    }

    // 2. 清除 Electron UserData 目录下的文件（跳过锁文件）
    const userDataPath = app.getPath('userData')
    if (fs.existsSync(userDataPath)) {
      const files = fs.readdirSync(userDataPath)
      for (const file of files) {
        // 跳过锁文件和 Singleton 文件，防止 EPERM 错误
        if (file === 'Lockfile' || file.startsWith('Singleton') || file === 'TransportSecurity') continue
        
        try {
          const curPath = path.join(userDataPath, file)
          fs.rmSync(curPath, { recursive: true, force: true })
        } catch (e: any) {
          console.warn(`Failed to delete ${file}:`, e.message)
          // 忽略无法删除的文件（通常是运行时锁定的，不影响重置效果）
        }
      }
    }

    // 3. 停止所有服务
    if (serverProcess) {
      serverProcess.kill()
      serverProcess = null
    }
    await stopFrontendServer()

    // 4. 重启应用
    app.relaunch()
    app.exit(0)
    return { success: true }
  } catch (error: any) {
    console.error('Failed to clear app data:', error)
    return { success: false, error: error.message }
  }
})

ipcMain.handle('app-config-get', async () => {
  return readAppConfig()
})

ipcMain.handle('show-save-dialog', async (event, options) => {
  const { filePath, canceled } = await dialog.showSaveDialog({
    ...options,
    filters: [
      { name: 'Excel Files', extensions: ['xlsx'] },
      { name: 'All Files', extensions: ['*'] }
    ]
  })
  return { filePath, canceled }
})

ipcMain.handle('save-file', async (event, { path: filePath, data }) => {
  try {
    // data comes as Uint8Array from renderer
    fs.writeFileSync(filePath, Buffer.from(data))
    return { success: true }
  } catch (err: any) {
    return { success: false, error: err.message }
  }
})

ipcMain.handle('select-directory', async () => {
  const res = await dialog.showOpenDialog({
    properties: ['openDirectory', 'createDirectory'],
  })
  return res
})

ipcMain.handle('get-local-ips', async () => {
  return { ips: getLocalIPs() }
})

ipcMain.handle('check-ports', async (_event, payload) => {
  const backendPort = Number(payload?.backendPort) || 0
  const frontendPort = Number(payload?.frontendPort) || 0
  return {
    backendAvailable: backendPort > 0 ? await isPortAvailable(backendPort) : false,
    frontendAvailable: frontendPort > 0 ? await isPortAvailable(frontendPort) : false,
  }
})

ipcMain.handle('apply-initial-config', async (_event, payload) => {
  const mode = String(payload?.mode || '')
  const backendPort = Number(payload?.backendPort) || 12345
  const frontendPort = Number(payload?.frontendPort) || 23456
  const baseFolderPath = String(payload?.baseFolderPath || '')
  const accessIp = String(payload?.accessIp || '')

  if (mode === 'server') {
    const backendOk = await isPortAvailable(backendPort)
    const frontendOk = await isPortAvailable(frontendPort)
    if (!backendOk) throw new Error('backend port unavailable')
    if (!frontendOk) throw new Error('frontend port unavailable')

    writeAppConfig({
      init_done: true,
      backend_port: backendPort,
      frontend_port: frontendPort,
      base_folder: baseFolderPath,
      access_ip: accessIp,
    })
    updateServerConfigYml(backendPort, baseFolderPath)
    startServer()
    await startFrontendServer(frontendPort)
    return { ok: true }
  }

  writeAppConfig({ init_done: true })
  return { ok: true }
})

// Electron 初始化完成并准备创建浏览器窗口时调用
app.whenReady().then(createWindow)

app.on('window-all-closed', () => {
  win = null
  
  // 杀死后端进程
  if (serverProcess) {
    serverProcess.kill();
    serverProcess = null;
  }

  stopFrontendServer().catch(() => {})

  // 除了 macOS 外，所有窗口关闭时退出应用
  if (process.platform !== 'darwin') app.quit()
})

app.on('second-instance', () => {
  if (win) {
    // 试图启动第二个实例时，聚焦到主窗口
    if (win.isMinimized()) win.restore()
    win.focus()
  }
})

app.on('activate', () => {
  const allWindows = BrowserWindow.getAllWindows()
  if (allWindows.length) {
    allWindows[0].focus()
  } else {
    // macOS 上点击 dock 图标如果没有窗口则重新创建
    createWindow()
  }
})
