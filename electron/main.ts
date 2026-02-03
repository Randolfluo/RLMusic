import { app, BrowserWindow, shell } from 'electron'
import { createRequire } from 'node:module'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import os from 'node:os'

const require = createRequire(import.meta.url)
const __dirname = path.dirname(fileURLToPath(import.meta.url))

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
if (process.platform === 'win32') app.setAppUserModelId(app.getName())

// 单例模式锁，防止启动多个实例
if (!app.requestSingleInstanceLock()) {
  app.quit()
  process.exit(0)
}

let win: BrowserWindow | null
let splash: BrowserWindow | null

/**
 * 创建主窗口
 */
function createWindow() {
  // 创建启动页
  splash = new BrowserWindow({
    width: 500,
    height: 300,
    transparent: true,
    frame: false,
    alwaysOnTop: true,
    icon: path.join(process.env.VITE_PUBLIC as string, 'images/logo/favicon.png'),
  })
  splash.loadFile(path.join(process.env.VITE_PUBLIC as string, 'loading.html'))

  win = new BrowserWindow({
    title: 'Local Music Player',
    show: false, // 先隐藏主窗口
    // 强制转换为 string 避免 TS 报错，或者使用 || '' 兼容
    icon: path.join(process.env.VITE_PUBLIC as string, 'images/logo/favicon.png'),
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

// Electron 初始化完成并准备创建浏览器窗口时调用
app.whenReady().then(createWindow)

app.on('window-all-closed', () => {
  win = null
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
