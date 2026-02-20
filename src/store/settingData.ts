import { defineStore } from "pinia";

export const useSettingDataStore = defineStore("settingData", {
  state: () => {
    return {
      // 主题色
      themeColor: "#009688",
      // 搜索历史
      searchHistory: true,
      // 轮播图显示
      bannerShow: true,
      // 自动签到
      autoSignIn: true,
      // 列表点击方式
      listClickMode: "dblclick",
      // 播放器样式
      playerStyle: "cover",
      // 播放器背景模糊
      playerBgBlur: true,
      // 底栏歌词显示
      bottomLyricShow: true,
      // 是否显示歌词翻译
      showTransl: true,
      // 歌曲音质
      songLevel: "exhigh",
      // 歌词滚动位置
      lyricsBlock: "center",
      // 歌词大小
      lyricsFontSize: 2.8,
      // 歌词模糊
      lyricsBlur: false,
      // 音乐频谱
      musicFrequency: false,
      // 粒子效果
      particleEffect: true,
      // 粒子数量
      particleLimit: 50,
      // 频谱跳动幅度
      musicFrequencyScale: 90,
      // 鼠标移入歌词区域暂停滚动
      lrcMousePause: true,
      // 播客模式
      podcastMode: false,
      // 桌面歌词
      desktopLyricShow: false,
      desktopLyricFontSize: 3.0,
      desktopLyricLock: false,
      desktopLyricFollowTheme: true,
      desktopLyricShowTranslation: true,
    };
  },
  getters: {
    // 获取是否开启翻译
    getShowTransl(state) {
      return state.showTransl;
    },
    // 获取主题色
    getThemeColor(state) {
      return state.themeColor;
    },
  },
  actions: {
    // 更改翻译开启选项
    setShowTransl(value: boolean) {
      this.showTransl = value;
    },
    // 设置主题色
    setThemeColor(value: string) {
      this.themeColor = value;
    },
    // 设置播客模式
    setPodcastMode(value: boolean) {
      this.podcastMode = value;
    },
    // 设置桌面歌词显示
    setDesktopLyricShow(value: boolean) {
      this.desktopLyricShow = value;
    },
    // 设置桌面歌词大小
    setDesktopLyricFontSize(value: number) {
      this.desktopLyricFontSize = value;
    },
    // 设置桌面歌词锁定
    setDesktopLyricLock(value: boolean) {
      this.desktopLyricLock = value;
    },
    // 设置桌面歌词跟随主题
    setDesktopLyricFollowTheme(value: boolean) {
      this.desktopLyricFollowTheme = value;
    },
    // 设置桌面歌词显示翻译
    setDesktopLyricShowTranslation(value: boolean) {
      this.desktopLyricShowTranslation = value;
    },
  },
  // 开启数据持久化
  persist: [
    {
      storage: localStorage,
    },
  ],
});


