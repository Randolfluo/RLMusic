import axios from "@/utils/request";

/**
 * 登录/注册参数定义
 */
export interface LoginParams {
  username: string;
  password: string;
}

export interface RegisterParams {
  username: string;
  password: string;
  email: string;
}

/**
 * 用户注册
 * POST /auth/register
 */
export const register = (data: RegisterParams) => {
  return axios({
    method: "POST",
    url: "/auth/register",
    data,
  });
};

/**
 * 用户登录
 * POST /auth/login
 */
export const login = (data: LoginParams) => {
  return axios({
    method: "POST",
    url: "/auth/login",
    data,
  });
};

/**
 * 获取登录状态
 * Note: 此接口可能需要后端适配，此处保留以兼容前端调用
 */
export const getLoginState = () => {
  return axios({
    method: "GET",
    hiddenBar: true,
    url: "/login/status",
    params: {
      timestamp: new Date().getTime(),
    },
  });
};
