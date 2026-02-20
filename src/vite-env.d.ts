/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_MODE: string
  readonly VITE_MUSIC_API: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare const $loadingBar: any;
declare const $message: any;
declare const $notification: any;
declare const $dialog: any;

interface Window {
  ipcRenderer: any;
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
