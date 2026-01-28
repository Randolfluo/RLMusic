/// <reference types="vite/client" />

declare const $loadingBar: any;
declare const $message: any;
declare const $notification: any;
declare const $dialog: any;

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
