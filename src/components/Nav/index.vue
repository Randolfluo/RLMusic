<!-- 顶部导航栏，常驻显示 -->
<template>
  <nav>
    <div class="left">
      <div class="logo" @click="router.push('/')">
        <img src="/images/logo/favicon.png" alt="logo" />
      </div>
      <div class="controls">
        <n-icon size="22" :component="Left" @click="router.go(-1)" />
        <n-icon size="22" :component="Right" @click="router.go(1)" />
      </div>
    </div>
    <div class="center">
      <router-link class="link link-home" to="/">发现</router-link>
      <router-link class="link link-search" to="/search">搜索</router-link>
      <n-dropdown
        class="link-discover-wrap"
        trigger="hover"
        size="large"
        :options="discoverOptions"
        @select="menuSelect"
      >
        <router-link class="link" to="/discover">一起听歌</router-link>
      </n-dropdown>
      <router-link class="link link-system" to="/system-stats">系统信息</router-link>
      <n-dropdown
        class="link-user-wrap"
        trigger="hover"
        size="large"
        :options="userOptions"
        @select="menuSelect"
      >
        <router-link class="link" to="/user">我的</router-link>
      </n-dropdown>
    </div>
    <div class="right">
      <!-- 下拉菜单 -->
      <n-dropdown
        class="dropdown"
        placement="bottom-end"
        :show="showDropdown"
        :show-arrow="true"
        :options="dropdownOptions"
        :on-clickoutside="closeDropdown"
        @select="dropdownSelect"
      >
        <n-avatar
          class="avatar"
          round
          size="small"
          :src="
            user.getUserData.avatarUrl
              ? user.getUserData.avatarUrl
              : '/images/ico/user-filling.svg'
          "
          :img-props="{ class: 'avatarImg' }"
          fallback-src="/images/ico/user-filling.svg"
          @click="handleAvatarClick"
        />
      </n-dropdown>
      <!-- 关于本站 -->
      <AboutSite ref="aboutSiteRef" />
    </div>
  </nav>
</template>

<script setup>
import { NIcon, NAvatar, NText, NProgress } from "naive-ui";
import {
  Left,
  Right,
  Login,
  Logout,
  Info,
  SettingTwo,
  History,
  SunOne,
  Moon,
} from "@icon-park/vue-next";
import { userStore, settingStore } from "@/store";
import { useRouter } from "vue-router";
import AboutSite from "@/components/DataModel/AboutSite.vue";

const router = useRouter();
const user = userStore();
const setting = settingStore();
const aboutSiteRef = ref(null);

// 下拉菜单显隐
const showDropdown = ref(false);
const closeDropdown = (event) => {
  // 解决点击头像无法关闭
  if (event.target.className == "avatarImg") {
    // handled in avatar click
  } else {
    showDropdown.value = false;
  }
};

const handleAvatarClick = () => {
  if (!user.userLogin) {
    router.push("/login");
  } else {
    showDropdown.value = !showDropdown.value;
  }
};

// 用户数据模块
const userDataRender = () => {
  return h(
    "div",
    {
      style:
        "display: flex; align-items: center; padding: 8px 12px; cursor: pointer",
      onclick: () => {
        user.userLogin ? router.push("/user") : router.push("/login");
        showDropdown.value = false;
      },
    },
    [
      h(NAvatar, {
        round: true,
        style: "margin-right: 12px",
        src: user.userLogin
          ? user.getUserData.avatarUrl
          : "/images/ico/user-filling.svg",
        fallbackSrc: "/images/ico/user-filling.svg",
      }),
      h("div", null, [
        h("div", null, [
          h(
            NText,
            { depth: 2 },
            {
              default: () =>
                user.userLogin ? user.getUserData.nickname : "未登录",
            }
          ),
        ]),
        h("div", { style: "font-size: 12px;" }, [
          h(
            NText,
            { depth: 3 },
            {
              default: () =>
                user.userLogin
                  ? user.getUserData.email  // 显示邮箱
                  : "登录后享受完整功能",
            }
          ),
        ]),
      ]),
    ]
  );
};

// 下拉框数据
const discoverOptions = ref([
  
]);
const userOptions = computed(() => [
  {
    label: "我喜欢的音乐",
    key: "/like",
  },  
{
    label: "收藏的歌单",
    key: "/likeplaylist",
  },  
  {
    label: "播放历史",
    key: "/history",
  },

]);
const dropdownOptions = computed(() => [
  {
    key: "header",
    type: "render",
    render: userDataRender,
  },
  {
    key: "header-divider",
    type: "divider",
  },
  {
    label: () => {
      return h(NText, null, {
        default: () =>
          setting.getSiteTheme == "light" ? "深色模式" : "浅色模式",
      });
    },
    key: "changeTheme",
    icon: () => {
      return h(NIcon, null, {
        default: () => (setting.getSiteTheme == "light" ? h(Moon) : h(SunOne)),
      });
    },
  },
  {
    label: "全局设置",
    key: "setting",
    icon: () => {
      return h(NIcon, null, {
        default: () => h(SettingTwo),
      });
    },
  },
  {
    label: () => {
      return h(NText, null, {
        default: () => (user.userLogin ? "退出登录" : "登录账号"),
      });
    },
    key: "user",
    icon: () => {
      return h(NIcon, null, {
        default: () => (user.userLogin ? h(Logout) : h(Login)),
      });
    },
  },
  {
    label: "关于本站",
    key: "about",
    icon: () => {
      return h(NIcon, null, {
        default: () => h(Info),
      });
    },
  },
]);

// 下拉框事件
const menuSelect = (key) => {
  router.push(key);
};
const dropdownSelect = (key) => {
  showDropdown.value = false;
  switch (key) {
    // 明暗切换
    case "changeTheme":
      setting.getSiteTheme == "light"
        ? setting.setSiteTheme("dark")
        : setting.setSiteTheme("light");
      break;
    // 播放历史
    case "history":
      router.push("/history");
      break;
    // 设置
    case "setting":
      router.push("/setting");
      break;
    // 用户
    case "user":
      if (user.userLogin) {
        // 退出登录
        $dialog.warning({
          class: "s-dialog",
          title: "退出登录",
          content: "确认退出当前用户登录？",
          positiveText: "退出登录",
          negativeText: "取消",
          onPositiveClick: () => {
            user.userLogOut();
            $message.success("已退出登录");
          },
        });
      } else {
        // 登录
        router.push("/login");
      }
      break;
    // 关于
    case "about":
      aboutSiteRef.value.openAboutSite();
      break;
    default:
      break;
  }
};
</script>

<style lang="scss" scoped>
nav {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  max-width: 1400px;
  margin: 0 auto;
  .left {
    flex: 1;
    max-width: 300px;
    display: flex;
    flex-direction: row;
    align-items: center;
    @media (max-width: 990px) {
      flex: initial;
    }
    // 移动端固定宽度，确保中间内容绝对居中
    @media (max-width: 768px) {
      flex: 0 0 auto; // 不再锁定宽度，允许根据内容收缩，或者设为 auto
      min-width: 0;   // 允许缩小
      margin-right: 8px; // 给一点间距
      justify-content: flex-start;
    }
    .logo {
      width: 30px;
      height: 30px;
      margin-right: 12px;
      cursor: pointer;
      img {
        width: 100%;
        height: 100%;
      }
    }
    .controls {
      display: flex;
      flex-direction: row;
      align-items: center;
      @media (max-width: 520px) {
        display: none;
      }
      .n-icon {
        margin: 0 4px;
        border-radius: 8px;
        padding: 4px;
        cursor: pointer;
        transition: all 0.3s;
        &:hover {
          background-color: var(--n-border-color);
        }
        &:active {
          transform: scale(0.95);
        }
      }
    }
  }
  .center {
    flex: 1;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    // 允许内容收缩，防止撑开
    min-width: 0; 

    @media (max-width: 768px) {
      // 改为靠左对齐，这是解决左侧项目被遮挡的关键
      // 配合 overflow-x: auto，用户可以滑动查看所有内容
      justify-content: flex-start; 
      flex: 1;                 
      overflow-x: auto;        
      scrollbar-width: none;   
      &::-webkit-scrollbar {
        display: none;         
      }
    }
    
    .link {
      display: block;
      text-decoration: none;
      color: var(--n-text-color);
      padding: 6px 16px;
      margin: 0 4px; // 增加一点间距
      border-radius: 8px;
      transition: all 0.3s;
      cursor: pointer;
      white-space: nowrap; // 强制不换行，防止文字竖排
      flex-shrink: 0; // 防止在小屏幕被挤压消失

      &:hover {
        background-color: $mainColor;
        color: var(--n-color);
      }
      &:active {
        transform: scale(0.95);
      }
      /* 移动端缩小间距和字号 */
      @media (max-width: 768px) {
        padding: 6px 8px;  // 稍微减小内边距，以便放下更多内容
        font-size: 15px;
        font-weight: bold;
        margin: 0 1px;     // 减小外边距
        // 选中状态添加底部边框效果，看起来更像移动端 App
        &.router-link-active {
            background-color: transparent !important;
            color: $mainColor !important;
            position: relative;
            &::after {
                content: "";
                position: absolute;
                bottom: 2px;
                left: 50%;
                transform: translateX(-50%);
                width: 16px;
                height: 3px;
                background-color: $mainColor;
                border-radius: 2px;
            }
        }
      }
    }

    .router-link-active {
      background-color: $mainColor;
      color: var(--n-color);
    }
  }
  .right {
    flex: 1;
    max-width: 300px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-end;
    // 移动端固定宽度，与左侧保持一致，实现平衡
    @media (max-width: 768px) {
      flex: 0 0 auto; // 同上，取消固定宽度
      min-width: 0;
      margin-left: 8px; // 给一点间距
    }
    .avatar {
      width: 30px;
      min-width: 30px;
      height: 30px;
      margin-left: 12px;
      box-shadow: 0 4px 12px -2px rgb(0 0 0 / 10%);
      cursor: pointer;
    }
  }
}
</style>
