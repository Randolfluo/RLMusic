import request, { apiBaseURL, resolveServerUrl } from "@/utils/request";
import lyricFormat from "@/utils/lyricFormat.js";

/**
 * 检查音乐是否可用
 * @param {Number} id - 歌曲 id
 */
export function checkMusicCanUse(_id: number | string) {
  return new Promise((resolve) => {
    resolve({ success: true, message: "ok" });
  });
}

/**
 * 获取音乐 url
 * @param {Number} id - 歌曲 id
 * @param {String} level - 音质
 */
export function getMusicUrl(id: number | string, _level?: string) {
  return new Promise((resolve) => {
    resolve({
      data: [
        {
          url: `${apiBaseURL}/song/stream/${id}`,
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
  return request({
    url: `/song/lyric/${id}`,
    method: "get",
  }).then((res) => {
    // 后端返回结构为 { code, data: { lrc: { lyric: "..." }, ... } }
    // 这里的 res 即为 response.data
    const data = res.data; 
    if (data && data.lrc && data.lrc.lyric) {
      return lyricFormat(data.lrc.lyric, data.tlyric?.lyric);
    }
    return [];
  });
}

/**
 * 记录播放历史
 * @param {Number} id - 歌曲 id
 */
export const recordHistory = (id: number) => {
    return request({
        method: "POST",
        url: "/song/history",
        data: {
            song_id: id
        }
    });
};

/**
 * 获取播放历史
 * @param {Number} page
 * @param {Number} limit
 */
export const getHistoryList = (page = 1, limit = 20) => {
    return request({
        method: "GET",
        url: "/song/history",
        params: {
            page,
            limit
        }
    });
};

/**
 * 清空播放历史
 */
export const clearHistory = () => {
    return request({
        method: "DELETE",
        url: "/song/history"
    });
};

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

/**
 * 获取我喜欢的歌曲列表
 * @param {Number} page
 * @param {Number} limit
 */
export const getLikedSongs = (page = 1, limit = 20) => {
    return request({
        method: "GET",
        url: "/song/like",
        params: {
            page,
            limit
        }
    });
};

/**
 * 获取歌曲详情
 * @param {Number|String} id - 歌曲 id
 */
export function getSongDetail(id: number | string) {
  return request({
    method: "GET",
    url: `/song/detail/${id}`,
  });
}

/**
 * 获取歌曲封面
 * @param {Number|String} id - 歌曲 id
 */
export function getSongCover(id: number | string) {
  return `${apiBaseURL}/song/cover/${id}`;
}

export function resolveCoverUrl(url?: string) {
  if (!url) return url;
  if (url.startsWith("/covers/") || url.startsWith("covers/")) {
    return resolveServerUrl(url.startsWith("/") ? url : `/${url}`);
  }
  if (url.startsWith("/api/") || url.startsWith("api/")) {
    return resolveServerUrl(url.startsWith("/") ? url : `/${url}`);
  }
  if (/^https?:\/\//.test(url)) {
    return url;
  }
  if (url.startsWith("/")) {
    return resolveServerUrl(url);
  }
  // 处理数据库中存储的裸文件名，补全为 /covers/ 路径
  return resolveServerUrl(`/covers/${url}`);
}

/**
 * 扫描用户音乐
 */
export const scanMusic = () => {
    return request({
        method: "POST",
        url: "/song/scan"
    });
};

/**
 * 获取歌手详情
 * @param {Number|String} id
 */
export const getArtistDetail = (id: number | string, page = 1, limit = 30) => {
  return request({
    method: "GET",
    url: `/song/artist/${id}`,
    params: {
      page,
      limit,
    },
  });
};

/**
 * 获取专辑详情
 * @param {Number|String} id
 */
export const getAlbumDetail = (id: number | string) => {
    return request({
        method: "GET",
        url: `/song/album/${id}`
    });
};
