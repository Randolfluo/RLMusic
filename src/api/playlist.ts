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
 * 获取用户公开歌单列表
 */
export const getUserPublicPlaylists = () => {
    return axios({
        method: "GET",
        url: "/song/playlists/user/public",
    });
};

/**
 * 获取用户私有歌单列表
 */
export const getUserPrivatePlaylists = () => {
    return axios({
        method: "GET",
        url: "/song/playlists/user/private",
    });
};

/**
 * 获取公共歌单详情
 * @param {string|number} id - 歌单ID
 * @param {number} page - 页码
 * @param {number} limit - 每页条数
 */
export const getPublicPlaylistDetail = (id: string | number, page: number = 1, limit: number = 20) => {
    return axios({
        method: "GET",
        url: `/song/playlist/public/${id}`,
        params: {
            page,
            limit
        }
    });
};

/**
 * 获取私有歌单详情
 * @param {string|number} id - 歌单ID
 * @param {number} page - 页码
 * @param {number} limit - 每页条数
 */
export const getPrivatePlaylistDetail = (id: string | number, page: number = 1, limit: number = 20) => {
    return axios({
        method: "GET",
        url: `/song/playlist/private/${id}`,
        params: {
            page,
            limit
        }
    });
};

/**
 * 更新歌单信息
 * @param {string|number} id
 * @param {object} data
 */
export const updatePlaylist = (id: string | number, data: any) => {
    return axios({
        method: "PUT",
        url: `/song/playlist/${id}`,
        data
    });
};


