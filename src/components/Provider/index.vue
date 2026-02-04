<!-- 提供全局的UI主题、样式和配置，最外层，影响所有子组件 -->
<template>
  <!-- 全局配置组件 -->
   <!-- 全局配置组件：这是最外层的容器，控制语言、主题、断点等核心配置 -->
  <n-config-provider
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme="theme"
    :theme-overrides="themeOverrides"
    :breakpoints="{   
      //  定义响应式布局的断点，决定 grid 布局在不同屏幕下的表现
      xs: 0,
      mb: 480,
      s: 640,
      m: 1024,
      l: 1280,
      xl: 1536,
      xxl: 1920,
    }"
    abstract
    inline-theme-disabled
  >  <!-- 以下是功能性 Provider，必须包裹在 n-config-provider 内部 -->
    <n-global-style />
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <slot></slot>
            <NaiveProviderContent />
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup>
import { defineComponent, h, ref, watch, onMounted } from "vue";
import {
  zhCN,         // 中文语言包
  dateZhCN,     // 中文日期包
  darkTheme,    // Naive UI 的深色主题对象
  useOsTheme,   // 钩子：获取操作系统的深色/浅色偏好
  useLoadingBar, // 钩子：获取加载条 API
  useDialog,     // 钩子：获取对话框 API
  useMessage,    // 钩子：获取消息 API
  useNotification,// 钩子：获取通知 API
} from "naive-ui";
import { settingStore } from "@/store";  // 引入全局设置

const setting = settingStore(); // 获取 Pinia 中的设置数据
const osThemeRef = useOsTheme();  // 获取当前系统的主题状态 (返回 'dark' 或 'light')

// 明暗切换
const theme = ref(null);
const changeTheme = () => {
   // 从 Store 中读取用户的设置
  if (setting.getSiteTheme == "light") {
    theme.value = null;
  } else if (setting.getSiteTheme == "dark") {
    theme.value = darkTheme;
  }
};

//自定义主题色配置
const themeOverrides = {
  common: {
    primaryColor: "#009688",
    primaryColorHover: "#26A69A",
    primaryColorSuppl: "#26A69A",
    primaryColorPressed: "#F64B41",
  },
};

// 挂载 naive 组件的方法，解决作用域问题
const setupNaiveTools = () => {
  window.$loadingBar = useLoadingBar(); // 进度条
  window.$notification = useNotification(); // 通知
  window.$message = useMessage(); // 信息
  window.$dialog = useDialog(); // 对话框
};

//使得在非组件上下文中，也能使用 Naive UI 的全局反馈组件
const NaiveProviderContent = defineComponent({
  setup() {
    setupNaiveTools();
  },
  render() {
    return h("div", {
      class: {
        tools: true,
      },
    });
  },
});


// 页面挂载时执行一次，初始化主题
onMounted(() => {
  changeTheme();
});

// 监听 Store 中的主题明暗变化
watch(
  () => setting.getSiteTheme,
  () => {
    changeTheme();
  }
);

// 监听操作系统的颜色方案
watch(
  () => osThemeRef.value,
  (value) => {
    if (setting.themeAuto) {
      value == "dark"
        ? setting.setSiteTheme("dark")
        : setting.setSiteTheme("light");
    }
  }
);


</script>