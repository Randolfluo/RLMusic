import { useSettingDataStore } from "./settingData";
import { useMusicDataStore } from "./musicData";
import { useUserDataStore } from "./userData";
import { useChatDataStore } from "./chatData";

export const settingStore = () => useSettingDataStore();
export const musicStore = () => useMusicDataStore();
export const userStore = () => useUserDataStore();
export const chatStore = () => useChatDataStore();
