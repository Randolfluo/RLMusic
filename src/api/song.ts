import request from "@/utils/request";

/**
 * 检查音乐是否可用
 * @param {Number} id - 歌曲 id
 */
export function checkMusicCanUse(id) {
  return new Promise((resolve) => {
    resolve({ success: true, message: "ok" });
  });
}

/**
 * 获取音乐 url
 * @param {Number} id - 歌曲 id
 * @param {String} level - 音质
 */
export function getMusicUrl(id, level) {
  return new Promise((resolve) => {
    resolve({
      data: [
        {
          url: `/api/song/stream/${id}`,
          fee: 0,
          freeTrialInfo: null,
          type: "flac",
        },
      ],
    });
  });
}

/**
 * 获取歌词
 * @param {Number} id - 歌曲 id
 */
export function getMusicLyric(id) {
  return new Promise((resolve) => {
    resolve({
      lrc: {
        lyric: "[00:00.000] 暂无歌词",
      },
      tlyric: {
        lyric: "",
      },
    });
  });
}
