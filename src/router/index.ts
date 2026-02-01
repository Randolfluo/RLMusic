import { createRouter, createWebHistory } from "vue-router";
import routes from "./routes";
import { getLoginState } from "@/api/login";
import { userStore } from "@/store";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
});

// 路由守卫
router.beforeEach((to, _from, next) => {
  const user = userStore();
  $loadingBar.start();    // 开始加载进度条
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
  $loadingBar.finish();   // 结束加载进度条
});

export default router;
