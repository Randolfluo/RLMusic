import { app, BrowserWindow, shell, ipcMain, screen } from 'electron'
import { createRequire } from 'node:module'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import os from 'node:os'
import { spawn } from 'child_process'

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

export const VITE_DEV_SERVER_URL = process.env['VITE_DEV_SERVER_URL']
export const MAIN_DIST = path.join(process.env.APP_ROOT, 'dist-electron')
export const RENDERER_DIST = path.join(process.env.APP_ROOT, 'dist')

process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL ? path.join(process.env.APP_ROOT, 'public') : RENDERER_DIST

// Win7 禁用 GPU 加速
if (os.release().startsWith('6.1')) app.disableHardwareAcceleration()

// Win10+ 通知设置应用ID
if (process.platform === 'win32') app.setAppUserModelId('LocalMusicPlayer')

// 屏蔽 Electron 常见的 Autofill 相关的终端报错
// "Request Autofill.enable failed", "Request Autofill.setAddresses failed"
app.commandLine.appendSwitch('disable-features', 'AutofillServerCommunication,Autofill,PasswordManager')

// 单例模式锁，防止启动多个实例
if (!app.requestSingleInstanceLock()) {
  app.quit()
  process.exit(0)
}

let win: BrowserWindow | null
let splash: BrowserWindow | null
let desktopLyricWindow: BrowserWindow | null = null
let serverProcess: any = null // 保存后端进程实例

// 启动后端服务
function startServer() {
  // 仅在生产环境尝试启动
  if (VITE_DEV_SERVER_URL) return;

  const serverName = process.platform === 'win32' ? 'server.exe' : 'server';
  // resources 目录路径 (Electron 打包后资源目录)
  const resourcesPath = process.resourcesPath;
  const serverPath = path.join(resourcesPath, serverName);
  
  const fs = require('fs');
  if (fs.existsSync(serverPath)) {
    console.log(`Starting server from: ${serverPath}`);
    // 启动服务，不显示窗口
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
  startServer(); // 尝试启动后端

  // 优先查找 .ico 文件 (Windows 最佳实践)，如果不存在则使用 .png
  let iconPath = path.join(process.env.VITE_PUBLIC as string, 'images/logo/favicon.ico')
  const fs = require('fs') // 引入 fs 模块用于检查文件是否存在
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

// Electron 初始化完成并准备创建浏览器窗口时调用
app.whenReady().then(createWindow)

app.on('window-all-closed', () => {
  win = null
  
  // 杀死后端进程
  if (serverProcess) {
    serverProcess.kill();
    serverProcess = null;
  }

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
