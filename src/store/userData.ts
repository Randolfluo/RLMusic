import { defineStore } from "pinia";
import { useChatDataStore } from "./chatData";

export const useUserDataStore = defineStore("userData", {
  state: () => {
    return {
      userLogin: false,
      userData: {
        avatarUrl: "",
        nickname: "",
        userId: 0,
        email: "",
        userGroup: "",
      },
      userOtherData: {},
    };
  },
  getters: {
    getUserData(state) {
      return state.userData;
    },
    getUserOtherData(state) {
      return state.userOtherData;
    },
  },
  actions: {
    userLogOut() {
      const chatStore = useChatDataStore();
      chatStore.stopListen();
      
      this.userLogin = false;
      this.userData = {
        avatarUrl: "",
        nickname: "",
        userId: 0,
        email: "",
        userGroup: "",
      };
      this.userOtherData = {};
    },
    setUserData(data: any) {
        this.userData = data;
        this.userLogin = true;
    },
    setUserOtherData(data: any = {}) {
        this.userOtherData = data;
    }
  },
  persist: {
    key: 'userData',
    storage: sessionStorage,
  }
});
