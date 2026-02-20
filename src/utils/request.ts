import axios from "axios";
import { aesDecrypt } from "@/utils/encrypt";

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
    _retry?: boolean;
  }
  // 扩展 AxiosResponse 以使返回值通过 TS 检查 (因为我们在 interceptor 中返回了 response.data)
  // 如果您想覆盖默认的 AxiosResponse 行为，这可能有点 tricky，
  // 通常更推荐下面的这种方式：不要在这里改 AxiosResponse，而是直接定义接口。
}

/**
 * 统一的 API 响应结构
 */
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 覆盖 axios 默认导出，强制让它通过泛型返回 ApiResponse
// 实际上我们这里直接返回 axios 实例，但在业务调用时需要注意类型转化。
// 为了让调用者方便，我们可以通过包装一层或者让 TS 认为 service 返回的是 data。

// 修正：Response interceptor 已经返回了 response.data。
// 所以 axios.get<T>() 返回的 Promise<AxiosResponse<T>> 实际上变成了 Promise<T>。
// 我们需要自定义一个类型或扩展 axios 的类型定义。



// 重新声明 AxiosInstance 的类型以匹配我们的拦截器行为
// 拦截器现在返回的是 ApiResponse 结构，而不是 AxiosResponse
declare module 'axios' {
  interface AxiosInstance {
    (config: AxiosRequestConfig): Promise<any>;
    request<T = any>(config: AxiosRequestConfig): Promise<T>;
    get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>;
    delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>;
    head<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>;
    post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>;
    put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>;
    patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>;
  }
}

// 创建 axios 实例
const appMode = import.meta.env.VITE_APP_MODE;
let baseURL = "/api";

if (appMode === 'server') {
    baseURL = "http://localhost:12345/api";
} else if (appMode === 'client') {
     const storedUrl = localStorage.getItem('server_url');
     if (storedUrl) {
         baseURL = storedUrl.endsWith('/') ? `${storedUrl}api` : `${storedUrl}/api`;
     } else {
         baseURL = "http://localhost:12345/api"; // Default fallback
     }
}

const service = axios.create({
  baseURL: baseURL, // 基础路径，通过 vite 代理转发
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
    // 改为 sessionStorage
    const token = sessionStorage.getItem("token");
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
        ResultCode.ERR_USER_NOT_EXIST,
      ];
      if (tokenErrors.includes(res.code)) {
        const config = response.config;
        // 尝试自动登录
        if (!config._retry && localStorage.getItem("auto_login") === 'true') {
           const savedUser = localStorage.getItem("remember_user");
           const savedPass = localStorage.getItem("remember_pass");
           
           if (savedUser && savedPass) {
               config._retry = true;
               try {
                   const password = aesDecrypt(savedPass);
                   // 使用默认 axios 实例发送请求，避免死循环
                   // 假设 api 代理为 /api
                   return axios.post(`${baseURL}/auth/login`, { 
                       username: savedUser, 
                       password: password 
                   }).then(loginRes => {
                       const loginData = loginRes.data;
                       if (loginData.code === ResultCode.SUCCESS) {
                           const newToken = loginData.data.token;
                           sessionStorage.setItem("token", newToken);
                           if (window.$message) window.$message.success('已自动重新登录');
                           
                           // 更新 Authorization 头并重试原请求
                           config.headers.Authorization = `Bearer ${newToken}`;
                           return service(config);
                       }
                       throw new Error("Auto login failed");
                   }).catch(() => {
                       sessionStorage.removeItem('token');
                       window.location.href = '/login';
                       return Promise.reject(response.data); // 返回原错误
                   });
               } catch (e) {
                   console.error("Auto login error", e);
               }
           }
        }

        if (window.$message) window.$message.error('登录失效，请重新登录');
        sessionStorage.removeItem('token');
        window.location.href = '/login'; 
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
