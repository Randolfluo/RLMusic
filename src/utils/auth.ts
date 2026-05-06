import { useUserDataStore } from "@/store/userData";

export function checkLogin(showDialog = true): boolean {
  const user = useUserDataStore();
  if (user.userLogin) return true;

  if (showDialog && window.$dialog) {
    window.$dialog.warning({
      title: "需要登录",
      content: "登录后即可收藏喜欢的歌曲",
      positiveText: "去登录",
      negativeText: "取消",
      onPositiveClick: () => {
        window.location.href = "/login";
      },
    });
  }
  return false;
}
