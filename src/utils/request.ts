import axios, { type AxiosRequestConfig } from "axios";

// 对应 server/internal/global/result.go
export const ResultCode = {
  SUCCESS: 1000,
  ERROR: 1001,
  ERR_REQUEST: 9001,
  ERR_DB_OP: 9004,
  ERR_REDIS_OP: 9005,
  ERR_USER_AUTH: 9006,
  
  ERR_TOKEN_NOT_EXIST: 1201,
  ERR_TOKEN_RUNTIME: 1202,
  ERR_TOKEN_WRONG: 1203,
  ERR_TOKEN_TYPE: 1204,
  ERR_TOKEN_CREATE: 1205,
  
  ERR_PASSWORD: 1002,
  ERR_USER_NOT_EXIST: 1003,
  ERR_USER_EXIST: 1004,
  ERR_PERMISSION: 1005,
  
  ERR_FILE_NOT_EXIST: 1101,
} as const;

// 扩展 AxiosRequestConfig 类型以包含自定义属性
declare module 'axios' {
  export interface AxiosRequestConfig {
    hiddenBar?: boolean;
  }
}

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
    const token = localStorage.getItem("token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`; // 添加 token
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
    
    // 业务错误拦截
    const res = response.data;
    if (res.code && res.code !== ResultCode.SUCCESS) {
      // Token 过期等处理
      const tokenErrors = [
        ResultCode.ERR_TOKEN_NOT_EXIST,
        ResultCode.ERR_TOKEN_RUNTIME,
        ResultCode.ERR_TOKEN_WRONG,
        ResultCode.ERR_TOKEN_TYPE,
      ];
      if (tokenErrors.includes(res.code)) {
        if (window.$message) window.$message.error('登录失效，请重新登录');
        localStorage.removeItem('token');
        // 可以选择跳转登录页，或者依赖路由守卫
        // window.location.href = '/login'; 
      }
      // 这里不统一 reject，交给业务层自己判断 res.code
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
          message = "请求出错(404)";
          break;
        case 408:
          message = "请求超时";
          break;
        case 500:
          message = "服务器错误";
          break;
        case 501:
          message = "服务未实现";
          break;
        case 502:
          message = "网络错误";
          break;
        case 503:
          message = "服务不可用";
          break;
        case 504:
          message = "网络超时";
          break;
        case 505:
          message = "HTTP版本不受支持";
          break;
        default:
          message = `连接出错(${error.response.status})!`;
      }
    } else {
        message = "网络连接异常,请检查网络!";
    }
    
    if (window.$message) {
        window.$message.error(message);
    } else {
        console.error(message);
    }
    
    return Promise.reject(error);
  }
);

export default service;
