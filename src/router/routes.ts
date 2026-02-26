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
      hideLayout: true,
    },
    component: () => import("@/views/Login/LoginView.vue"),
  },
  {
    path: "/init",
    name: "init",
    meta: {
      title: "初始化配置",
      hideLayout: true,
    },
    component: () => import("@/views/Init/index.vue"),
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
  // 收藏歌单
  {
    path: "/likeplaylist",
    name: "like-playlist",
    meta: {
      title: "收藏歌单",
      needLogin: true,
    },
    component: () => import("@/views/Playlist/SubscribedPlaylists.vue"),
  },
  // 我喜欢的歌曲
  {
    path: "/like",
    name: "like",
    meta: {
      title: "我喜欢的歌曲",
      needLogin: true,
    },
    component: () => import("@/views/User/like.vue"),
  },
  // 一起听歌
  // {
  //   path: "/listen-together",
  //   name: "listen-together",
  //   meta: {
  //     title: "一起听歌",
  //     needLogin: true,
  //   },
  //   component: () => import("@/views/ListenTogether/index.vue"),
  // },
  // 历史记录
  {
    path: "/history",
    name: "history",
    meta: {
      title: "播放历史",
      needLogin: true,
    },
    component: () => import("@/views/History/HistoryView.vue"),
  },
  // 搜索页
  {
    path: "/search",
    name: "search",
    meta: {
      title: "搜索",
    },
    component: () => import("@/views/Search/index.vue"),
  },
  // 系统信息
  {
    path: "/system-stats",
    name: "system-stats",
    meta: {
      title: "系统信息",
    },
    component: () => import("@/views/System/StatsView.vue"),
  },
  // 歌曲详情
  {
    path: "/song/:id",
    name: "song",
    meta: {
      title: "歌曲详情",
    },
    component: () => import("@/views/Song/SongView.vue"),
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
  // 桌面歌词
  {
    path: "/desktop-lyric",
    name: "desktop-lyric",
    meta: {
      title: "桌面歌词",
      hideLayout: true,
    },
    component: () => import("@/views/DesktopLyric/index.vue"),
  },
  // 管理员后台
  {
    path: "/admin",
    name: "admin",
    meta: {
      title: "管理员后台",
      needLogin: true,
    },
    component: () => import("@/views/Admin/index.vue"),
  },
  // 管理员用户管理
  {
    path: "/admin/users",
    name: "admin-users",
    meta: {
      title: "用户管理",
      needLogin: true,
    },
    component: () => import("@/views/Admin/UserManage.vue"),
  },
  // 管理员公共歌单管理
  {
    path: "/admin/playlists",
    name: "admin-playlists",
    meta: {
      title: "公共歌单管理",
      needLogin: true,
    },
    component: () => import("@/views/Admin/PlaylistManage.vue"),
  },

  {
    path: "/:pathMatch(.*)",
    redirect: "/404",
  },
];











export default routes;
