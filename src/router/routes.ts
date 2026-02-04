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
  // 歌手详情
  {
    path: "/artist",
    name: "artist",
    meta: {
      title: "歌手详情",
    },
    component: () => import("@/views/Artist/index.vue"),
  },
  // 专辑详情
  {
    path: "/album",
    name: "album",
    meta: {
      title: "专辑详情",
    },
    component: () => import("@/views/Album/index.vue"),
  },
  // 歌单详情
  {
    path: "/playlist/:id",
    name: "playlist",
    meta: {
        title: "歌单详情",
    },
    component: () => import("@/views/Playlist/index.vue"),
  },
  // 公共歌单
  {
    path: "/playlists",
    name: "public-playlists",
    meta: {
      title: "公共歌单",
    },
    component: () => import("@/views/Playlist/PublicPlaylists.vue"),
  },
  // 私有歌单
  {
    path: "/private-playlists",
    name: "private-playlists",
    meta: {
      title: "私有歌单",
      needLogin: true,
    },
    component: () => import("@/views/Playlist/PrivatePlaylists.vue"),
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