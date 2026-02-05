import axios from "@/utils/request";

/**
 * 搜索部分
 */

/**
 * 获取热门搜索列表
 */
export const getSearchHot = () => {
  return Promise.resolve({ data: [] });
};

export const getSearchSongs = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/song", { params: { keywords, limit, offset } });
};

export const getSearchArtists = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/artist", { params: { keywords, limit, offset } });
};

export const getSearchAlbums = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/album", { params: { keywords, limit, offset } });
};

export const getSearchPlaylists = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/playlist", { params: { keywords, limit, offset } });
};

