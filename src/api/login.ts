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
 * 退出登录
 * POST /auth/logout
 */
export const logout = () => {
  return axios({
    method: "POST",
    url: "/auth/logout",
  });
};

