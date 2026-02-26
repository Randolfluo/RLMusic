<!-- 顶部导航栏，常驻显示 -->
<template>
  <nav class="glass-nav">
    <div class="nav-content">
      <div class="left">
        <div class="logo-wrapper" @click="router.push('/')">
          <img src="/images/logo/favicon.png" alt="logo" class="logo-img" />
          <span class="app-name">Music</span>
        </div>
        <div class="controls">
          <n-button circle quaternary size="small" @click="router.go(-1)" class="nav-btn">
            <template #icon>
              <n-icon :component="Left" />
            </template>
          </n-button>
          <n-button circle quaternary size="small" @click="router.go(1)" class="nav-btn">
            <template #icon>
              <n-icon :component="Right" />
            </template>
          </n-button>
        </div>
      </div>
      
      <div class="center">
        <div class="nav-links">
          <router-link class="link-item" to="/" active-class="active">
            <span class="link-text">发现</span>
            <div class="active-indicator"></div>
          </router-link>
          <router-link class="link-item" to="/search" active-class="active">
            <span class="link-text">搜索</span>
            <div class="active-indicator"></div>
          </router-link>
          <router-link class="link-item" to="/system-stats" active-class="active">
            <span class="link-text">系统信息</span>
            <div class="active-indicator"></div>
          </router-link>
          <n-dropdown
            class="link-user-dropdown"
            trigger="hover"
            size="large"
            :options="userOptions"
            @select="menuSelect"
          >
            <router-link class="link-item" to="/user" active-class="active">
              <span class="link-text">我的</span>
              <div class="active-indicator"></div>
            </router-link>
          </n-dropdown>
        </div>
      </div>

      <div class="right">
        <!-- 管理员入口 -->
        <n-tooltip trigger="hover">
          <template #trigger>
            <div class="admin-trigger" @click="handleAdminClick">
              <n-icon :component="Permissions" size="20" />
            </div>
          </template>
          管理员入口
        </n-tooltip>
        <n-tooltip trigger="hover">
          <template #trigger>
            <div class="connect-trigger" @click="handleServerConfigClick">
              <n-icon :component="Connection" size="20" />
            </div>
          </template>
          连接服务
        </n-tooltip>

        <!-- 移动端汉堡菜单按钮 -->
        <div class="mobile-menu-trigger" @click="showMobileMenu = true">
          <n-icon :component="HamburgerButton" size="24" />
        </div>

        <!-- 下拉菜单 -->
        <n-dropdown
          class="user-dropdown"
          placement="bottom-end"
          :show="showDropdown"
          :show-arrow="true"
          :options="dropdownOptions"
          :on-clickoutside="closeDropdown"
          @select="dropdownSelect"
        >
          <div class="user-trigger" @click="handleAvatarClick">
            <n-avatar
              class="avatar"
              round
              :size="36"
              :src="resolveAvatarUrl(user.getUserData.avatarUrl) || '/images/ico/user-filling.svg'"
              fallback-src="/images/ico/user-filling.svg"
            />
          </div>
        </n-dropdown>
        <!-- 关于本站 -->
        <!-- <div class="about-trigger" @click="aboutSiteRef?.openAboutSite()">
            <n-icon :component="Info" size="20" />
        </div> -->
        <AboutSite ref="aboutSiteRef" />
      </div>
    </div>

    <!-- 移动端侧边菜单 -->
    <n-drawer v-model:show="showMobileMenu" :width="280" placement="right">
      <n-drawer-content title="菜单">
        <n-list clickable>
          <n-list-item @click="router.push('/'); showMobileMenu = false">
            <template #prefix><n-icon :component="Left" /></template>
            发现
          </n-list-item>
          <n-list-item @click="router.push('/search'); showMobileMenu = false">
            <template #prefix><n-icon :component="Left" /></template>
            搜索
          </n-list-item>
          <n-list-item @click="router.push('/system-stats'); showMobileMenu = false">
            <template #prefix><n-icon :component="Left" /></template>
            系统信息
          </n-list-item>
          <n-list-item @click="router.push('/user'); showMobileMenu = false">
            <template #prefix><n-icon :component="Left" /></template>
            我的
          </n-list-item>
          <n-list-item @click="handleServerConfigClick(); showMobileMenu = false">
            <template #prefix><n-icon :component="Connection" /></template>
            连接服务
          </n-list-item>
        </n-list>
      </n-drawer-content>
    </n-drawer>
  </nav>
</template>

<script setup>
import { NIcon, NAvatar, NText, NProgress, NButton, NDropdown, NTooltip } from "naive-ui";
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
  Permissions,
  Connection,
  HamburgerButton,
} from "@icon-park/vue-next";
import { userStore, settingStore } from "@/store";
import { useRouter } from "vue-router";
import AboutSite from "@/components/DataModel/AboutSite.vue";
import { useMessage, NDrawer, NDrawerContent, NList, NListItem } from "naive-ui";
import { resolveAvatarUrl } from "@/api/user";

const router = useRouter();
const user = userStore();
const setting = settingStore();
const aboutSiteRef = ref(null);
const message = useMessage();
const showMobileMenu = ref(false);

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

const handleAdminClick = () => {
  if (user.userLogin) {
    // 假设后端返回了 userGroup，并在登录时保存到了 userData
    if (user.getUserData.userGroup === 'admin') {
      message.success("欢迎管理员");
      router.push("/admin"); 
    } else {
      message.warning("当前账号非管理员");
    }
  } else {
    router.push("/login");
  }
};

const handleServerConfigClick = () => {
  window.dispatchEvent(new CustomEvent("open-server-config"));
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
          ? resolveAvatarUrl(user.getUserData.avatarUrl)
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
  // {
  //   label: () => {
  //     return h(NText, null, {
  //       default: () =>
  //         setting.getSiteTheme == "light" ? "深色模式" : "浅色模式",
  //     });
  //   },
  //   key: "changeTheme",
  //   icon: () => {
  //     return h(NIcon, null, {
  //       default: () => (setting.getSiteTheme == "light" ? h(Moon) : h(SunOne)),
  //     });
  //   },
  // },
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
    // // 明暗切换
    // case "changeTheme":
    //   setting.getSiteTheme == "light"
    //     ? setting.setSiteTheme("dark")
    //     : setting.setSiteTheme("light");
    //   break;
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
.glass-nav {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  // Glassmorphism background is usually handled by the parent or global header style, 
  // but we can add a subtle background here if needed. 
  // For now, we assume the parent container might have some background, 
  // or we make this nav transparent to blend with the page background.
  background: rgba(255, 255, 255, 0.4);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;

  :global(.dark) & {
    background: rgba(30, 30, 30, 0.4);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  }

  .nav-content {
    width: 100%;
    max-width: 1400px;
    height: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 24px;

    @media (max-width: 768px) {
      padding: 0 16px;
    }
  }

  .left {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 24px;

    .logo-wrapper {
      display: flex;
      align-items: center;
      gap: 12px;
      cursor: pointer;
      transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);

      &:hover {
        transform: scale(1.05);
      }

      .logo-img {
        width: 32px;
        height: 32px;
        filter: drop-shadow(0 4px 6px rgba(0,0,0,0.1));
      }

      .app-name {
        font-size: 20px;
        font-weight: 800;
        background: linear-gradient(120deg, var(--n-primary-color), #4facfe);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: -0.5px;
        
        @media (max-width: 640px) {
          display: none;
        }
      }
    }

    .controls {
      display: flex;
      gap: 8px;
      
      @media (max-width: 768px) {
        display: none;
      }
      
      .nav-btn {
        transition: all 0.2s;
        &:hover {
          background-color: rgba(0,0,0,0.05);
          transform: translateY(-1px);
        }
      }
    }
  }

  .center {
    flex: 2;
    display: flex;
    justify-content: center;
    
    .nav-links {
      display: flex;
      align-items: center;
      background: rgba(0, 0, 0, 0.03);
      padding: 4px;
      border-radius: 99px;
      border: 1px solid rgba(255, 255, 255, 0.2);
      
      :global(.dark) & {
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.05);
      }

      @media (max-width: 768px) {
      display: none !important;
    }
    }

    .link-item {
      position: relative;
      text-decoration: none;
      color: var(--n-text-color-3);
      padding: 8px 24px;
      border-radius: 99px;
      font-size: 15px;
      font-weight: 600;
      transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;

      @media (max-width: 768px) {
        padding: 6px 12px;
        font-size: 14px;
      }

      .active-indicator {
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: var(--n-primary-color);
        opacity: 0;
        z-index: -1;
        border-radius: 99px;
        transition: all 0.3s ease;
        transform: scale(0.8);
      }

      &:hover {
        color: var(--n-text-color);
        background: rgba(0,0,0,0.03);
        
        :global(.dark) & {
            background: rgba(255,255,255,0.05);
        }
      }

      &.active {
        color: var(--n-primary-color);
        text-shadow: none;

        .active-indicator {
          opacity: 0.15;
          transform: scale(1);
          box-shadow: none;
        }
        
        &:hover {
            background: transparent;
        }
      }
    }
  }

  .right {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 16px;

    .about-trigger,
    .admin-trigger,
    .connect-trigger,
    .mobile-menu-trigger {
        width: 36px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        cursor: pointer;
        color: var(--n-text-color-3);
        transition: all 0.2s;
        
        &:hover {
            background: rgba(0,0,0,0.05);
            color: var(--n-text-color);
            transform: scale(1.1);
        }
    }
    
    .mobile-menu-trigger {
      display: none;
      @media (max-width: 768px) {
        display: flex;
      }
    }

    .user-trigger {
      cursor: pointer;
      border-radius: 50%;
      padding: 2px;
      border: 2px solid transparent;
      transition: all 0.3s;
      display: flex;
      align-items: center;

      &:hover {
        border-color: var(--n-primary-color);
        transform: scale(1.05);
      }

      .avatar {
        box-shadow: 0 4px 12px rgba(0,0,0,0.1);
      }
    }
  }
}
</style>
