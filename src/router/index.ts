import { createRouter, createWebHistory } from "vue-router";
import routes from "./routes";


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
});

// 全局路由守卫



export default router;