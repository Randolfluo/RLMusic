import request from "@/utils/request";

/**
 * 检查音乐是否可用
 * @param {Number} id - 歌曲 id
 */
export function checkMusicCanUse(id: number | string) {
  return new Promise((resolve) => {
    resolve({ success: true, message: "ok" });
  });
}

/**
 * 获取音乐 url
 * @param {Number} id - 歌曲 id
 * @param {String} level - 音质
 */
export function getMusicUrl(id: number | string, level?: string) {
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
export function getMusicLyric(id: number | string) {
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

/**
 * 点赞/取消点赞歌曲
 * @param {Number|String} id - 歌曲 id
 */
export function toggleLike(id: number | string) {
  return request({
    method: "POST",
    url: `/song/like/${id}`,
  });
}
