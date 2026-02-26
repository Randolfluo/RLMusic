import request, { resolveServerUrl } from "@/utils/request";

/**
 * 语音合成 (TTS)
 * @param {Object} data
 * @param {string} data.text - 文本内容
 * @param {string} [data.voice] - 音色
 * @param {string} [data.format] - 格式 (flac/mp3/wav)
 * @param {number} [data.sample_rate] - 采样率
 * @param {number} [data.volume] - 音量 (0-100)
 * @param {number} [data.rate] - 语速 (0.5-2.0)
 * @param {number} [data.pitch] - 音调 (0.5-2.0)
 */
export const generateTTS = (data: {
    text: string;
    voice?: string;
    format?: string;
    sample_rate?: number;
    volume?: number;
    rate?: number;
    pitch?: number;
}) => {
    return request({
        method: "POST",
        url: "/ai/tts",
        data
    });
};

/**
 * 硅基流动 Qwen 对话
 * @param {string} prompt - 提示词
 */
export const generateChat = (prompt: string) => {
    return request({
        method: "POST",
        url: "/ai/chat",
        data: {
            prompt
        }
    });
};

/**
 * 获取歌单AI描述
 * @param {number|string} id - 歌单ID
 */
export const getPlaylistAIDescription = (id: number | string) => {
    return request({
        method: "GET",
        url: `/song/playlist/description/${id}`
    });
};

/**
 * 获取艺术家AI描述
 * @param {number|string} id - 艺术家ID
 */
export const getArtistAIDescription = (id: number | string) => {
    return request({
        method: "GET",
        url: `/song/artist/description/${id}`
    });
};

/**
 * 获取专辑AI描述
 * @param {number|string} id - 专辑ID
 */
export const getAlbumAIDescription = (id: number | string) => {
    return request({
        method: "GET",
        url: `/song/album/description/${id}`
    });
};

/**
 * 生成歌曲AI最终开场白并转语音 (Step 4: TTS)
 * @param {number|string} id - 歌曲ID
 */
export const getSongOpeningTTS = (id: number | string) => {
    return request({
        method: "GET",
        url: `/song/podcast/opening-tts/${id}`
    });
};

/**
 * 批量生成歌单内歌曲的开场白
 * @param {number|string} id - 歌单ID
 * @param {string} [voice] - 音色 (默认 Neil)
 */
export const generatePlaylistIntros = (id: number | string, voice: string = "Neil") => {
    return request({
        method: "POST",
        url: `/song/podcast/generate-playlist-intros/${id}`,
        data: {
            voice
        }
    });
};

/**
 * 批量生成所有公共歌单的开场白
 * @param {string} [voice] - 音色 (默认 Neil)
 */
export const generateAllPublicPlaylistIntros = (voice: string = "Neil") => {
    return request({
        method: "POST",
        url: "/song/podcast/generate-all-public-playlist-intros",
        data: {
            voice
        }
    });
};

/**
 * 批量生成公共歌单描述 (Admin)
 */
export const generatePublicPlaylistDescriptions = () => {
    return request({
        method: "POST",
        url: "/song/playlist/generate-public-descriptions"
    });
};

/**
 * 批量生成艺术家描述 (Admin)
 */
export const generateArtistDescriptions = () => {
    return request({
        method: "POST",
        url: "/song/artist/generate-descriptions"
    });
};

/**
 * 批量生成专辑描述 (Admin)
 */
export const generateAlbumDescriptions = () => {
    return request({
        method: "POST",
        url: "/song/album/generate-descriptions"
    });
};

/**
 * 获取歌曲开场白文本
 * @param {number|string} id - 歌曲ID
 */
export const getSongOpeningText = (id: number | string) => {
    return request({
        method: "GET",
        url: `/song/opening/text/${id}`
    });
};

/**
 * 获取歌曲开场白语音地址
 * @param {number|string} id - 歌曲ID
 */
export const getSongOpeningAudioUrl = (id: number | string) => {
    return resolveServerUrl(`/api/song/opening/${id}`);
};
