import axios from "@/utils/request";

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  avatar: string;
  last_login: string;
  total_songs: number;
  total_albums: number;
  total_artists: number;
  total_duration: number;
  favorite_song: string;
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
