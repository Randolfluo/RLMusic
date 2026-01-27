import { createApp } from 'vue'
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
// 引入pinia数据持久化的官方插件，在store开启持久化选项
import App from './App.vue'
import router from './router'

import "@/style/global.scss"; // 全局样式
const app = createApp(App);

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate); // 使用数据持久化插件

app.use(pinia);     // 使用pinia

app.use(router);    // 注册路由

app.mount("#app");  //将应用挂载到index.html中的<div id="app"></div> 


if (navigator.serviceWorker) {
  navigator.serviceWorker.addEventListener('controllerchange', () => {
    // 弹出更新提醒
    console.log("站点已更新，刷新后生效");
    // $message 是 naive-ui 的全局对象，需要确保此时已经挂载
    if (window.$message) {
      window.$message.info("站点已更新，刷新后生效", {
        closable: true,
        duration: 0,
      });
    }
  });
}
