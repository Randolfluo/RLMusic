import axios from "@/utils/request";

export interface SystemStats {
    song_count: number;
    song_volume?: number;
    album_count: number;
    artist_count: number;
    music_duration: number;
    playlist_count: number;
    user_count: number;
    system_uptime: number;
    user_listening_duration?: number; // Optional now as it's removed from main stats
    user_scanned_duration?: number;   // Optional now as it's removed from main stats
    cpu_usage?: number;
    mem_usage?: number;
    api_call_count?: number;
}

export interface SystemStatus {
    cpu_usage: number;
    mem_usage: number;
    api_call_count: number;
    system_uptime: number;
    go_routines: number;
    total_volume: number;
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
 * 获取系统实时状态
 */
export const getSystemStatus = () => {
    return axios({
        method: "GET",
        url: "/system/status",
        hiddenBar: true // 不显示顶部进度条
    }) as Promise<{ code: number, msg: string, data: SystemStatus }>;
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

/**
 * 获取局域网 IP 列表
 * @param port 前端端口，用于生成完整 URL
 */
export const getLocalIPs = (port?: string) => {
    return axios({
        method: "GET",
        url: "/system/local-ips",
        params: { port }
    }) as Promise<{ code: number, msg: string, data: { ips: string[], port: string, urls: string[] } }>;
};

/**
 * 重置系统数据 (Admin)
 */
export const resetSystem = () => {
    return axios({
        method: "DELETE",
        url: "/system/reset"
    });
};
