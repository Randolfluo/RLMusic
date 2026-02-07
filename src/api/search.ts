import axios from "@/utils/request";

/**
 * 搜索部分 API (Updated)
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

/**
 * 聚合搜索建议 (For Search/index.vue)
 */
export const getSearchSuggest = async (keywords: string) => {
    // Parallel requests
    const [songs, artists, albums, playlists] = await Promise.all([
        getSearchSongs(keywords, 10),
        getSearchArtists(keywords, 5),
        getSearchAlbums(keywords, 5),
        getSearchPlaylists(keywords, 5)
    ]);
    
    // Extract data from backend response structure: { code: 1000, data: { result: { songs: [] } } }
    // Note: Request interceptor returns response.data, so code/msg/data are available directly.
    return {
        songs: (songs as any).data?.result?.songs || [],
        artists: (artists as any).data?.result?.artists || [],
        albums: (albums as any).data?.result?.albums || [],
        playlists: (playlists as any).data?.result?.playlists || []
    };
}

