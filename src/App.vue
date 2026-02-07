<template>
  <Provider>    
 <!-- provider组件，全局样式管理 -->
    <n-layout style="height: 100vh"> 
       <!-- Naive UI的Layout组件，创建基础的页面布局框架 -->
      <n-layout-header bordered>    
        <Nav />
      </n-layout-header>
      <!-- 智能布局容器 -->
      <n-layout-content
        position="absolute"
        :class="music.getPlaylists[0] && music.showPlayBar ? 'show' : ''"
        :native-scrollbar="false"
        embedded
      >
       <!-- 主内容包装器 -->
        <main ref="mainContent" class="main">
          <!-- 智能返回顶部按钮 -->
          <n-back-top
            :bottom="music.getPlaylists[0] && music.showPlayBar ? 100 : 40"
            style="transition: all 0.3s"
          />  
          <!-- Vue Router的占位符，根据当前URL显示对应的页面 -->
          <router-view v-slot="{ Component }">  
            <keep-alive>
              <Transition name="scale" mode="out-in">
                <component :is="Component" />
              </Transition>
            </keep-alive>
          </router-view>
        </main>
      </n-layout-content>
    </n-layout>
    <Player />
    <BigPlayer v-if="music.getPlaylists[0]" />
  </Provider>
</template>

<script setup lang="ts">
import { musicStore, userStore, settingStore } from "@/store";
import { useRouter } from "vue-router";
import { login } from "@/api/login"; // 引入登录api
import { aesDecrypt } from "@/utils/encrypt"; // 引入解密
import Provider from "@/components/Provider/index.vue";
import Nav from "@/components/Nav/index.vue";
import Player from "@/components/Player/index.vue";
import BigPlayer from "@/components/Player/BigPlayer.vue";
import packageJson from "@/../package.json";
import { ref, onMounted } from 'vue';
import { ResultCode } from "@/utils/request";

const music = musicStore();
const user = userStore();
const router = useRouter();
const setting = settingStore();

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
                        avatarUrl: res.data.avatar || "" 
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
.n-layout-header {
  height: 60px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}
.n-layout-content {
  top: 60px;
  transition: all 0.3s;
  &.show {  //有播放器时留出空间
    bottom: 70px;
  }
  :deep(.n-scrollbar-rail--vertical) {
    right: 0;
  }
  .main {
    max-width: 1400px;
    margin: 0 auto;
  }
}

// 路由跳转动画
.scale-enter-active,
.scale-leave-active {
  transition: all 0.2s ease;
}

.scale-enter-from,
.scale-leave-to {
  opacity: 0;
  transform: scale(0.98);
}
</style>
