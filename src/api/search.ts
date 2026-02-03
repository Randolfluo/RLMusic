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
export const getSearchSuggest = (keywords) => {
  return Promise.resolve({ result: {} });
};

/**
 * 搜索结果
 * @param {string} keywords - 搜索关键词
 */
export const getSearchData = (keywords, limit = 30, offset = 0, type = 1) => {
  return Promise.resolve({ data: { result: { songs: [] } } });
};
