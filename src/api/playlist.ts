import axios from "@/utils/request";

/**
 * 获取公共歌单列表
 */
export const getPublicPlaylists = () => {
    return axios({
        method: "GET",
        url: "/song/playlists/public",
    });
};

/**
 * 获取歌单详情
 * @param {string|number} id - 歌单ID
 */
export const getPlaylistDetail = (id: string | number) => {
    return axios({
        method: "GET",
        url: `/song/playlist/${id}`,
    });
};


