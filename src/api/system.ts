import axios from "@/utils/request";

export interface SystemStats {
    song_count: number;
    album_count: number;
    artist_count: number;
    music_duration: number;
    playlist_count: number;
    user_count: number;
    system_uptime: number;
    user_listening_duration: number;
    user_scanned_duration: number;
    cpu_usage?: number;
    mem_usage?: number;
    api_call_count?: number;
}

/**
 * 获取系统统计信息
 */
export const getSystemStats = () => {
    return axios({
        method: "GET",
        url: "/system/stats"
    }) as Promise<{ code: number, msg: string, data: SystemStats }>;
};

/**
 * 初始化基础文件夹
 */
export const initFolder = () => {
    return axios({
        method: "POST",
        url: "/file/initFolder"
    });
};

/**
 * 更新系统配置
 * @param {object} data
 */
export const updateConfig = (data: { filepath: string }) => {
    return axios({
        method: "POST",
        url: "/system/config",
        data
    });
};
