import axios from "@/utils/request";

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  avatar: string;
  last_login: string;
  ip_src: string;
  total_duration: number;
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
