<template>
  <Provider>    
    <!-- provider组件，全局样式管理 -->
    
    <!-- 全屏页面（如登录页） -->
    <div v-if="route.meta.hideLayout" style="width: 100vw; height: 100vh;">
      <router-view />
    </div>

    <!-- 普通页面布局 -->
    <n-layout class="main-layout" v-else> 
       <!-- Naive UI的Layout组件，创建基础的页面布局框架 -->
      <n-layout-header class="nav-header">    
        <Nav />
      </n-layout-header>
      <!-- 智能布局容器 -->
      <n-layout-content
        class="layout-content"
        position="absolute"
        :class="music.getPlaylists[0] && music.showPlayBar ? 'show-player' : ''"
        :native-scrollbar="false"
        content-style="min-height: 100%; display: flex; flex-direction: column;"
      >
       <!-- 主内容包装器 -->
        <main ref="mainContent" class="main-container">
          <!-- 智能返回顶部按钮 -->
          <n-back-top
            :bottom="music.getPlaylists[0] && music.showPlayBar ? 120 : 40"
            :right="40"
            class="custom-back-top"
          >
            <n-icon :component="ToTop" size="24" />
          </n-back-top>  
          <!-- Vue Router的占位符，根据当前URL显示对应的页面 -->
          <router-view v-slot="{ Component }">  
            <keep-alive>
              <Transition name="fade-slide" mode="out-in">
                <component :is="Component" />
              </Transition>
            </keep-alive>
          </router-view>
        </main>
      </n-layout-content>
    </n-layout>
    
    <!-- 仅在非全屏模式下显示播放器 -->
    <template v-if="!route.meta.hideLayout">
      <Player />
      <BigPlayer v-if="music.getPlaylists[0]" />
    </template>
  </Provider>
</template>

<script setup lang="ts">
import { musicStore, userStore } from "@/store";
import { login } from "@/api/login"; // 引入登录api
import { aesDecrypt } from "@/utils/encrypt"; // 引入解密
import Provider from "@/components/Provider/index.vue";
import Nav from "@/components/Nav/index.vue";
import Player from "@/components/Player/index.vue";
import BigPlayer from "@/components/Player/BigPlayer.vue";
import packageJson from "@/../package.json";
import { ref, onMounted } from 'vue';
import { ResultCode } from "@/utils/request";
import { useRoute } from "vue-router"; // 引入 useRoute
import { NIcon, NBackTop } from "naive-ui";
import { ToTop } from "@icon-park/vue-next";

const route = useRoute(); // 获取当前路由信息
const music = musicStore();
const user = userStore();
// const router = useRouter();
// const setting = settingStore();

// 自动登录逻辑
onMounted(() => {
    const savedUser = localStorage.getItem("remember_user");
    const savedPass = localStorage.getItem("remember_pass");
    if (savedUser && savedPass && !user.userLogin) {
        try {
            const password = aesDecrypt(savedPass);
            login({ username: savedUser, password: password }).then(res => {
                if (res.code === ResultCode.SUCCESS) {
                    const userData = {
                        userId: res.data.id,
                        nickname: res.data.username,
                        email: res.data.email,
                        userGroup: res.data.user_group,
                        avatarUrl: res.data.avatar,
                    };
                    user.setUserData(userData);
                    localStorage.setItem("token", res.data.token); 
                    console.log("Auto login success");
                }
            }).catch(e => {
                console.error("Auto login failed", e);
            });
        } catch (e) {
            console.error("Auto login decrypt failed", e);
        }
    }
});

// 监听空格键控制播放暂停
onMounted(() => {
  window.addEventListener("keydown", (e) => {
    if (e.code === "Space") {
      // 如果当前焦点在输入框或文本域中，不触发
      if (
        document.activeElement?.tagName === "INPUT" ||
        document.activeElement?.tagName === "TEXTAREA"
      ) {
        return;
      }
      e.preventDefault();
      music.setPlayState(!music.getPlayState);
    }
  });
});

const mainContent = ref<HTMLElement | null>(null);

// 公告数据
const annShow =
  import.meta.env.VITE_ANN_TITLE && import.meta.env.VITE_ANN_CONTENT
    ? true
    : false;
const annTitle = import.meta.env.VITE_ANN_TITLE;
const annContene = import.meta.env.VITE_ANN_CONTENT;
const annDuration = Number(import.meta.env.VITE_ANN_DURATION);

// // 空格暂停与播放
// const spacePlayOrPause = (e: KeyboardEvent) => {
//   if (e.code === "Space") {
//     // console.log(e.target.tagName);
//     if (router.currentRoute.value.name == "video") return false;
//     const target = e.target as HTMLElement;
//     if (target.tagName === "BODY") {    //在输入框内时，不执行空格暂停与播放
//       e.preventDefault();
//       music.setPlayState(!music.getPlayState);
//     } else {
//       return false;
//     }
//   }
// };

// 系统重置
const cleanAll = () => {
  $message ? $message.success("重置成功") : alert("重置成功");
  localStorage.clear();
  window.location.href = "/";
};


onMounted(() => {
  // 将应用内部功能暴露到全局window对象的技术，允许在浏览器控制台直接调用这些功能
  (window as any).$mainContent = mainContent.value;  // 将Vue组件的DOM引用暴露到全局window对象
  (window as any).$cleanAll = cleanAll;   // 暴露重置函数

  // 公告，如果有则显示
  if (annShow) {
    if(typeof $notification !== 'undefined') {
       $notification["info"]({
        content: annTitle,
        meta: annContene,
        duration: annDuration,
      });
    }
  }
  
  // 版权声明
  const logoText = `${packageJson.name}`;
  const copyrightNotice = `\n\n版本: ${packageJson.version}\n作者: ${packageJson.author}\n作者主页: ${packageJson.home}\n`;
  console.info(
    `%c${logoText} %c ${copyrightNotice}`,
    "color:#f55e55;font-size:26px;font-weight:bold;",
    "font-size:16px"
  );
  console.info(
    "若站点出现异常，可尝试在下方输入 %c$cleanAll()%c 然后按回车来重置",
    "background: #eaeffd;color:#f55e55;padding: 4px 6px;border-radius:8px;",
    "background:unset;color:unset;"
  );

  // // 检查账号登录状态
  // getLoginState()
  //   .then((res) => {
  //     if (res.data.profile && user.userLogin) {
  //       // 保存登录信息
  //       user.userLogin = true;
  //       user.setUserData(res.data.profile);
  //     } else {
  //       user.userLogOut();
  //       if (music.getPlayListMode === "cloud") {
  //         $message.info("登录已失效，请重新登录");
  //         music.setPlaylists([]);
  //       }
  //     }
  //   })
  //   .catch((err) => {
  //     $message.error("请求发生错误");
  //     console.error("请求发生错误" + err);
  //     router.push("/500");
  //     return false;
  //   });

  // // 获取喜欢音乐列表
  // music.setLikeList();

  // // 键盘监听
  // window.addEventListener("keydown", spacePlayOrPause);
});

</script>



<style lang="scss" scoped>
.main-layout {
  height: 100vh;
  /* 柔和的浅色渐变背景 */
  background: linear-gradient(135deg, #fdfbfb 0%, #ebedee 100%);
  transition: background 0.3s ease;

  :global(.dark) & {
    /* 深色模式下的深邃背景 */
    background: linear-gradient(135deg, #1a1a1a 0%, #0d0d0d 100%);
  }
}

.nav-header {
  height: 60px;
  background: transparent !important;
  z-index: 100;
  /* 移除边框，由 Nav 组件处理 */
  border: none !important;
}

.layout-content {
  top: 60px;
  bottom: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: transparent !important;

  &.show-player {
    bottom: 90px; /* 留出播放器空间 */
  }

  :deep(.n-scrollbar-rail--vertical) {
    right: 4px;
    width: 6px;
  }
}

.main-container {
  width: 100%;
  max-width: 1440px; /* 宽屏适配 */
  margin: 0 auto;
  padding: 0; /* Nav已包含padding，这里内容由各页面自行控制padding，或者统一加 */
  /* 
     由于各页面（如 Home, Playlist）已有 padding，这里设为 0 以免双重 padding。
     但考虑到一致性，如果大部分页面都需要 padding，可以在这里加。
     查看之前的页面代码，大多都有 padding: 24px。
     所以这里设为 0，让页面自己控制。
  */
  min-height: 100%;
  box-sizing: border-box;
}

/* 优化返回顶部按钮 */
:deep(.n-back-top) {
  background-color: var(--n-color-primary) !important;
  color: #fff !important;
  box-shadow: 0 4px 16px rgba(var(--n-color-primary-rgb), 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 50%;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:hover {
    transform: translateY(-4px) scale(1.05);
    box-shadow: 0 8px 24px rgba(var(--n-color-primary-rgb), 0.5);
  }
}

/* 路由跳转动画 - 渐变滑动 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
