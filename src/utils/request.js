import axios from "axios";

// 创建 axios 实例
const service = axios.create({
  baseURL: "/api", // 基础路径，通过 vite 代理转发
  timeout: 30000, // 请求超时时间
  withCredentials: true, // 跨域请求时发送 cookies
});

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 如果没有配置 hiddenBar，则显示加载条
    if (!config.hiddenBar && window.$loadingBar) {
      window.$loadingBar.start();
    }
    return config;
  },
  (error) => {
    if (window.$loadingBar) {
      window.$loadingBar.error();
    }
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response) => {
    // 关闭加载条
    if (!response.config.hiddenBar && window.$loadingBar) {
      window.$loadingBar.finish();
    }
    return response.data;
  },
  (error) => {
    if (window.$loadingBar) {
      window.$loadingBar.error();
    }
    
    // 处理错误信息
    let message = "请求失败";
    if (error.response) {
      switch (error.response.status) {
        case 301:
          message = "需登录";
          break;
        case 400:
          message = "请求错误";
          break;
        case 401:
          message = "未授权，请登录";
          break;
        case 403:
          message = "拒绝访问";
          break;
        case 404:
          message = "请求地址出错";
          break;
        case 500:
          message = "服务器内部错误";
          break;
        default:
          message = `连接错误 ${error.response.status}`;
      }
    } else {
        message = error.message;
    }

    if (window.$message) {
        window.$message.error(message);
    }
    
    return Promise.reject(error);
  }
);

export default service;
