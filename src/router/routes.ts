// 路由配置文件
import Home from "@/views/Home/index.vue";

const routes: any = [
  {
    path: "/",
    name: "home",
    component: Home,
  },
  {
    path: "/login",
    name: "login",
    meta: {
      title: "登录",
    },
    component: () => import("@/views/Login/LoginView.vue"),
  },
  {
    path: "/user",
    name: "user",
    meta: {
      title: "用户中心",
      needLogin: true,
    },
    component: () => import("@/views/User/index.vue"),
  },
  // 全局设置设置
  {
    path: "/setting",
    name: "setting",
    meta: {
      title: "全局设置",
    },
    component: () => import("@/views/Setting/SettingView.vue"),
  },
 // 403
  {
    path: "/403",
    name: "403",
    meta: {
      title: "403",
    },
    component: () => import("@/views/State/403.vue"),
  },
  // 500
  {
    path: "/500",
    name: "500",
    meta: {
      title: "500",
    },
    component: () => import("@/views/State/500.vue"),
  },
  // 404
  {
    path: "/404",
    name: "404",
    meta: {
      title: "404",
    },
    component: () => import("@/views/State/404.vue"),
  },

  {
    path: "/:pathMatch(.*)",
    redirect: "/404",
  },
];











export default routes;