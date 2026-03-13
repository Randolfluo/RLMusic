import { createRouter, createWebHashHistory } from "vue-router";
import routes from "./routes";
import { userStore } from "@/store";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routes,
});

// 路由守卫
router.beforeEach((to, _from, next) => {
  const user = userStore();
  if (typeof window.$loadingBar !== "undefined") {
    window.$loadingBar.start();    // 开始加载进度条
  }

  const isElectron = typeof navigator !== "undefined" && navigator.userAgent.includes("Electron");
  const isCapacitor = typeof (window as any).Capacitor !== "undefined";
  const initDone = typeof localStorage !== "undefined" && localStorage.getItem("init_done") === "true";
  
  if (isElectron || isCapacitor) {
    const allowList = ["/init", "/desktop-lyric"];
    if (!initDone && !allowList.includes(to.path)) {
      next("/init");
      return;
    }
    // 允许已初始化的 Electron Server 用户手动访问 /init
    const isServerMode = import.meta.env.VITE_APP_MODE === "server" && isElectron;
    if (initDone && to.path === "/init" && !isServerMode) {
      next("/");
      return;
    }
  }

  // 判断是否需要登录
  if (to.meta.needLogin) {
    if (user.userLogin) {
      next();
    } else {
      $message.error("请登录账号后使用");
      next("/login");
    }
    // Remove the old getLoginState check which caused 404 or 500 error when not logged in
    /*
    getLoginState()
      .then((res) => {
        if (res.data?.profile && user.userLogin) {
          user.setUserData(res.data.profile);
          if (!Object.keys(user.getUserOtherData).length) {
            user.setUserOtherData();
          }
          next();
        } else {
          $message.error(
            localStorage.getItem("cookie")
              ? "登录过期，请重新登录"
              : "请登录账号后使用"
          );
          user.userLogOut();
          next("/login");
        }
      })
      .catch((err) => {
        $message.error("请求发生错误");
        console.error("请求发生错误" + err);
        next("/500");
        return false;
      });
      */
  } else {
    if (!Object.keys(user.getUserOtherData).length) user.setUserOtherData();
    next();
  }
});

router.afterEach(() => {
  if (typeof window.$loadingBar !== "undefined") {
    window.$loadingBar.finish();   // 结束加载进度条
  }
});

export default router;
