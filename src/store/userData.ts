import { defineStore } from "pinia";

export const useUserDataStore = defineStore("userData", {
  state: () => {
    return {
      userLogin: false,
      userData: {
        avatarUrl: "",
        nickname: "",
        userId: 0,
      },
    };
  },
  getters: {
    getUserData(state) {
      return state.userData;
    },
  },
  actions: {
    userLogOut() {
      this.userLogin = false;
      this.userData = {
        avatarUrl: "",
        nickname: "",
        userId: 0,
      };
    },
    setUserData(data: any) {
        this.userData = data;
        this.userLogin = true;
    }
  },
  persist: {
    key: 'userData',
    storage: localStorage,
  }
});
