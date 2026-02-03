import { fileURLToPath, URL } from "node:url";
import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import electron from "vite-plugin-electron/simple";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { NaiveUiResolver } from "unplugin-vue-components/resolvers";
import { VitePWA } from "vite-plugin-pwa";
import type { ConfigEnv } from "vite";


// https://vite.dev/config/
export default ({ mode }: ConfigEnv) => {
  const VITE_MUSIC_API = loadEnv(mode, process.cwd()).VITE_MUSIC_API;
  console.log("VITE_MUSIC_API:", VITE_MUSIC_API); // 在终端输出环境变量，方便调试
  return defineConfig({
    plugins: [
      vue(),
      electron({
        main: {
          entry: "electron/main.ts",
        },
        preload: {
          input: "electron/preload.ts",
        },
        // Optional: Use Node.js API in the Renderer-process
        renderer: {},
      }),
      AutoImport({
        imports: [
          "vue",  // 自动导入 Vue 的 API（ref, reactive, computed 等）
          {
            "naive-ui": [  // 自动导入 Naive UI 的常用组件
              "useDialog",    // 对话框
              "useMessage",   // 消息提示
              "useNotification", // 通知
              "useLoadingBar",   // 加载条
            ],
          },
        ],
      }),
      Components({
        resolvers: [NaiveUiResolver()], // 自动解析 Naive UI 组件
      }),
      VitePWA({
        registerType: "autoUpdate",  // 自动更新 Service Worker

        workbox: {
          clientsClaim: true,     // 立即控制页面
          skipWaiting: true,      // 跳过等待，立即激活新 Service Worker
          cleanupOutdatedCaches: true,  // 清理过期缓存

          runtimeCaching: [  // 运行时缓存策略
            {
              urlPattern: /(.*?)\.(woff2|woff|ttf)/,  // 字体文件
              handler: "CacheFirst",  // 缓存优先
              options: { cacheName: "file-cache" }
            },
            {
              urlPattern: /(.*?)\.(webp|png|jpe?g|svg|gif|bmp|psd|tiff|tga|eps)/,
              handler: "CacheFirst",  // 图片文件缓存优先
              options: { cacheName: "image-cache" }
            }
          ]
        },

        manifest: {  // Web App Manifest
          name: "localmusicplayer",           // 应用名称
          short_name: "localmusicplayer",     // 短名称
          description: "基于VUE的音乐播放器设计与实现",  // 描述
          display: "standalone",     // 显示模式：独立应用
          start_url: "/",           // 启动URL
          theme_color: "#fff",      // 主题色
          background_color: "#efefef",  // 背景色
          icons: [  // 应用图标
            {
              src: "/images/logo/favicon.png",
              sizes: "200x200",
              type: "image/png"
            }
          ]
        }
      })
    ],
    server: {
      port: 8888,      // 开发服务器端口
      open: true,      // 自动打开浏览器
      host: true,     // 如果需要外部访问，可以添加 host
      proxy: {  // API 代理配置
        "/api": {
          target: VITE_MUSIC_API,  // 从环境变量读取API地址
          changeOrigin: true,  // 改变请求源
          // rewrite: (path) => path.replace(/^\/api/, "")  // 重写路径
        },
        "/covers": {
           target: VITE_MUSIC_API,
           changeOrigin: true,
        }
      }
    },
    // 添加 ssr 配置
    ssr: {
      noExternal: true,  // 或者根据需求配置
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: '@use "@/style/index.scss" as *;',  // 全局 SCSS 导入
        }
      }
    },
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),  // @ 指向 src 目录
      }
    },
    build: {
      minify: "terser",  // 使用 terser 压缩
      terserOptions: {
        compress: {
          pure_funcs: ["console.log"],  // 删除 console.log
        }
      },
      sourcemap: false,  // 不生成 sourcemap
    }
  })
}