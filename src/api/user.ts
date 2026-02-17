import axios from "@/utils/request";

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  avatar: string;
  last_login: string;
  ip_src: string;
  total_duration: number;
  user_group: string;
}

/**
 * 获取用户信息
 */
export const getUserInfo = () => {
  return axios({
    method: "GET",
    url: "/user/info",
  });
};

/**
 * 管理员获取所有用户列表
 */
export const adminGetAllUsers = (params: { page: number; limit: number; query?: string }) => {
  return axios({
    method: "GET",
    url: "/admin/user/list",
    params,
  });
};

/**
 * 管理员删除用户
 */
export const adminDeleteUser = (id: number) => {
  return axios({
    method: "DELETE",
    url: `/admin/user/delete/${id}`,
  });
};

/**
 * 管理员修改用户权限
 */
export const adminUpdateUserRole = (id: number, user_group: string) => {
  return axios({
    method: "PUT",
    url: `/admin/user/role/${id}`,
    data: { user_group },
  });
};

/**
 * 检查是否为管理员
 */
export const checkIsAdmin = () => {
  return axios({
    method: "GET",
    url: "/user/is-admin",
  });
};

/**
 * 上传头像
 */
export const uploadAvatar = (formData: FormData) => {
    return axios({
        method: "POST",
        url: "/user/avatar",
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

/**
 * 注销用户
 */
export const deleteUser = () => {
    return axios({
        method: "DELETE",
        url: "/user",
    });
};

/**
 * 更新用户听歌时长
 */
export const updateListeningDuration = (duration: number) => {
    return axios({
        method: "POST",
        url: "/user/duration",
        data: { duration }
    });
};
