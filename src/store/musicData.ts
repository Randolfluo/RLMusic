import { defineStore } from "pinia";
import { getSongTime } from "@/utils/timeTools";
import { useUserDataStore } from "./userData";
import { recordHistory, toggleLike, getLikedSongs } from "@/api/song";
import { ResultCode } from "@/utils/request";

export const useMusicDataStore = defineStore("musicData", {
  state: () => {
    return {
      // 是否展示播放界面
      showBigPlayer: false,
      // 是否展示播放控制条
      showPlayBar: true,
      // 是否展示播放列表
      showPlayList: false,
      // 播放状态
      playState: false,
      // 当前歌曲播放链接
      playSongLink: null as string | null,
      // 歌词偏移时间 (秒)
      lyricOffset: 0,
      // 当前歌曲歌词
      playSongLyric: [] as any[],
      // 当前歌曲歌词播放索引
      playSongLyricIndex: 0,
      // 当前歌曲是否拥有翻译
      playSongTransl: false,
      // 每日推荐
      dailySongsData: [],
      // 歌单分类
      catList: {},
      // 精品歌单分类
      highqualityCatList: [],
      // 持久化数据
      persistData: {
        // 搜索历史
        searchHistory: [] as string[],
        // 是否处于私人 FM 模式
        personalFmMode: false,
        // 私人 FM 数据
        personalFmData: {},
        // 播放列表类型
        playListMode: "list",
        // 喜欢音乐列表
        likeList: [] as number[],
        // 播放列表
        playlists: [] as any[],
        // 当前歌曲索引
        playSongIndex: 0,
        // 当前播放模式
        // normal-顺序播放 random-随机播放 single-单曲循环
        playSongMode: "normal",
        // 当前播放时间
        playSongTime: {
          currentTime: 0,
          duration: 0,
          barMoveDistance: 0,
          songTimePlayed: "00:00",
          songTimeDuration: "00:00",
        },
        // 播放音量
        playVolume: 0.7,
        // 播放速率
        playRate: 1.0,
        // 静音前音量
        playVolumeMute: 0,
        // 列表状态
        playlistState: 0, // 0 顺序 1 单曲循环 2 随机
        // 播放历史
        playHistory: [] as any[],
      },
    };
  },
  getters: {
    // 获取是否处于私人FM模式
    getPersonalFmMode(state) {
      return state.persistData.personalFmMode;
    },
    // 获取私人FM模式数据
    getPersonalFmData(state) {
      return state.persistData.personalFmData;
    },
    // 获取是否拥有翻译
    getPlaySongTransl(state) {
      return state.playSongTransl;
    },
    // 获取歌词偏移
    getLyricOffset(state) {
      return state.lyricOffset;
    },
    // 获取每日推荐
    getDailySongs(state) {
      return state.dailySongsData;
    },
    // 获取播放列表
    getPlaylists(state) {
      return state.persistData.playlists;
    },
    // 获取播放模式
    getPlaySongMode(state) {
      return state.persistData.playSongMode;
    },
    // 获取播放速率
    getPlayRate(state) {
      return state.persistData.playRate || 1.0;
    },
    // 获取播放时间获取当前歌曲索引
    getPlaySongIndex(state) {
        return state.persistData.playSongIndex;
    },
    // 获取当前歌曲
    getPlaySongData(state): any {
      return state.persistData.playlists[state.persistData.playSongIndex];
    },
    // 获取当前歌词
    getPlaySongLyric(state) {
      return state.playSongLyric;
    },
    // 获取当前歌词索引
    getPlaySongLyricIndex(state) {
      return state.playSongLyricIndex;
    },
    // 获取当前播放时间
    getPlaySongTime(state) {
      return state.persistData.playSongTime;
    },
    // 获取播放状态
    getPlayState(state) {
      return state.playState;
    },
    // 获取播放链接
    getPlaySongLink(state) {
      return state.playSongLink;
    },
    // 获取喜欢音乐列表
    getLikeList(state) {
      return state.persistData.likeList;
    },
    // 获取播放历史
    getPlayHistory(state) {
      return state.persistData.playHistory;
    },
    // 获取播放列表模式
    getPlayListMode(state) {
      return state.persistData.playListMode;
    },
    // 获取搜索历史
    getSearchHistory(state) {
      return state.persistData.searchHistory;
    },
    // 判断歌曲是否在红心列表
    getSongIsLike: (state) => {
      return (id: number) => {
        return state.persistData.likeList.includes(id);
      };
    },
  },
  actions: {
    // 更改播放界面显隐
    setBigPlayerState(value: boolean) {
      this.showBigPlayer = value;
    },
    // 更改播放状态
    setPlayState(value: boolean) {
      this.playState = value;
    },
    // 更改播放列表
    setPlaylists(value: any[]) {
      this.persistData.playlists = value;
      this.persistData.playSongIndex = 0;
    },
    // 更改当前歌曲索引
    setPlaySongIndex(value: "next" | "prev" | number) {
      if (typeof value === "number") {
        this.persistData.playSongIndex = value;
      } else {
        const len = this.persistData.playlists.length;
        if (len === 0) return;

        // 随机播放模式
        if (this.persistData.playSongMode === "random") {
            if (len === 1) return;
            // 随机选取一个新的索引
            let newIndex = this.persistData.playSongIndex;
            while (newIndex === this.persistData.playSongIndex) {
                newIndex = Math.floor(Math.random() * len);
            }
            this.persistData.playSongIndex = newIndex;
            return;
        }

        if (value === "next") {
          this.persistData.playSongIndex =
            (this.persistData.playSongIndex + 1) % len;
        } else {
          this.persistData.playSongIndex =
            (this.persistData.playSongIndex - 1 + len) % len;
        }
      }
    },
    // 添加到下一首播放
    addSongToNext(song: any) {
        const list = this.persistData.playlists;
        // 如果列表为空，直接播放
        if (list.length === 0) {
            this.setPlaylists([song]);
            this.setPlayState(true);
            return;
        }
        
        // 检查歌曲是否已在列表中
        const index = list.findIndex(item => item.id === song.id);
        const currentIndex = this.persistData.playSongIndex;
        
        // 插入位置：当前播放歌曲的下一位
        const nextIndex = currentIndex + 1;
        
        if (index !== -1) {
            // 如果已在列表中
            if (index === currentIndex) return; // 就是当前正在播放的，不做处理
            
            // 移动到下一位
            // 先删除
            list.splice(index, 1);
            // 如果删除的元素在当前播放之前，currentIndex 需要减 1，插入位置相应调整
            // 但因为我们是用 splice 修改原数组，直接重新计算位置比较麻烦
            // 简单处理：重新获取修改后的列表
            
            // 修正：如果删除的元素在当前播放之前，当前索引会发生变化，需要修正
            if (index < currentIndex) {
                this.persistData.playSongIndex--;
            }
            // 重新计算插入位置（因为删除了一个元素）
            const newNextIndex = this.persistData.playSongIndex + 1;
            list.splice(newNextIndex, 0, song);
        } else {
            // 不在列表中，直接插入
            list.splice(nextIndex, 0, song);
        }
        this.persistData.playlists = list;
    },
    // 更改音量
    setPlayVolume(value: number) {
      this.persistData.playVolume = value;
    },
    // 更改播放时间
    setPlaySongTime(value: any) {
      this.persistData.playSongTime.currentTime = value.currentTime;
      this.persistData.playSongTime.duration = value.duration;
      // 计算进度条位置
      if (value.duration > 0) {
        this.persistData.playSongTime.barMoveDistance =
          (value.currentTime / value.duration) * 100;
        this.persistData.playSongTime.songTimePlayed = getSongTime(
          value.currentTime * 1000
        );
        this.persistData.playSongTime.songTimeDuration = getSongTime(
           value.duration * 1000
        );
      }
      // 计算歌词位置
      // lyricOffset > 0 代表歌词延迟显示（即需要在播放进度更靠后时才显示当前行），意味着“有效时间”要减去 offset
      const effectiveTime = value.currentTime - this.lyricOffset;
      const index = this.playSongLyric.findIndex(
        (item: any) => item.time > effectiveTime
      );
      if (index === -1) {
        this.playSongLyricIndex = this.playSongLyric.length - 1;
      } else {
        this.playSongLyricIndex = index - 1;
      }
    },
    // 设置歌词偏移
    setLyricOffset(value: number) {
      this.lyricOffset = value;
    },
    // 添加播放历史
    setPlayHistory(data: any) {
      const list = this.persistData.playHistory;
      const index = list.findIndex((item: any) => item.id === data.id);
      if (index !== -1) {
        list.splice(index, 1);
      }
      list.unshift(data);
      if (list.length > 100) {
        list.pop();
      }
      this.persistData.playHistory = list;

      // 同步到后端
      const userStore = useUserDataStore();
      if (userStore.userLogin && data && data.id) {
        recordHistory(data.id);
      }
    },
    // 更改当前歌曲播放链接
    setPlaySongLink(value: string) {
      this.playSongLink = value as string | null;
    },
    // 更改当前歌曲歌词
    setPlaySongLyric(value: any[]) {
      this.playSongLyric = value;
      this.playSongTransl = value.some((item) => item.lyricFy);
    },
    // 私人FM不感冒
    setFmDislike(id: number) {
        // TODO: Implement API logic
    },
    // 更改喜欢列表
    async changeLikeList(id: number, like: boolean) {
      const list = this.persistData.likeList;
      // 乐观更新
      if (like) {
        if (!list.includes(id)) list.push(id);
      } else {
        const index = list.indexOf(id);
        if (index !== -1) list.splice(index, 1);
      }
      this.persistData.likeList = list;
      
      // 调用 API
      try {
        await toggleLike(id);
      } catch (error) {
        console.error("Toggle like failed", error);
        // 回滚 (可选，这里简化处理暂不回滚，或者提示用户)
        $message.error("操作失败");
        // 实际应用中应该回滚状态
        if (like) {
            const index = list.indexOf(id);
            if (index !== -1) list.splice(index, 1);
        } else {
            list.push(id);
        }
        this.persistData.likeList = list;
      }
    },
    // 获取喜欢列表
    async setLikeList() {
        try {
            // 获取所有喜欢的歌曲 (假设 limit 足够大)
            const res: any = await getLikedSongs(1, 10000);
            if (res.code === ResultCode.SUCCESS) {
                const ids = res.data.list.map((song: any) => song.id);
                this.persistData.likeList = ids;
            }
        } catch (error) {
            console.error("Fetch liked songs failed", error);
        }
    },
    setSearchHistory(text: string | null, clean: boolean = false) {
      if (clean) {
        this.persistData.searchHistory = [];
        return;
      }
      if (!text) return;
      const history = this.persistData.searchHistory;
      if (history.includes(text)) {
        history.splice(history.indexOf(text), 1);
      }
      history.unshift(text);
      if (history.length > 20) {
        history.pop();
      }
      this.persistData.searchHistory = history;
    },
    // 设置播放模式
    setPlaySongMode() {
      if (this.persistData.playSongMode === "normal") {
        this.persistData.playSongMode = "random";
        window.$message.info("随机播放");
      } else if (this.persistData.playSongMode === "random") {
        this.persistData.playSongMode = "single";
        window.$message.info("单曲循环");
      } else {
        this.persistData.playSongMode = "normal";
        window.$message.info("列表循环");
      }
    },
    // 设置播放速率
    setPlayRate(rate: number) {
      this.persistData.playRate = rate;
    },
    // 设置播放列表状态
  }
});


