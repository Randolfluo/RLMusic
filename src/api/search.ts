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

/**
 * 搜索建议
 * @param {string} keywords - 搜索关键词
 */
export const getSearchSuggest = (keywords: string) => {
  return axios.get("/search/suggest", { params: { keywords } });
};

/**
 * 搜索结果
 * @param {string} keywords - 搜索关键词
 * @param {number} limit - 每页限制
 * @param {number} offset - 分页偏移量 (页码)
 * @param {number} type - 搜索类型 (1: 歌曲)
 */
export const getSearchData = (keywords: string, limit: number = 30, offset: number = 1, type: number = 1) => {
    return axios.get("/search/detail", {
        params: {
            keywords,
            type,
            offset,
            limit
        }
    });
};
