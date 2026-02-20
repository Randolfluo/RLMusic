# 基于VUE的音乐播放器




### 🛠️ 环境与部署 (Environment & Deployment)
- [ ] Docker 部署：通过挂载目录实现目录共享
- [ ] 后端容器化：支持 Electron 带有服务器和不带服务器（纯本地）的版本选择


- [ ] 桌面歌词：创建背景透明 (`transparent: true`) 且忽略鼠标事件 (`ignoreMouseEvents`) 的置顶窗口












web构建分为两种






检查setting.view的每个配置是否都能正常工作


添加修改密码api，添加移动端适配，美化ui。


web端部署（需要先）
``` shell
pnpm build:web
```
然后将问价


安卓端部署
``` shell
pnpm build:web
方式 A：使用 Android Studio (推荐，可视化的)
pnpm exec cap sync android
方式 B：命令行构建 (适合 CI/CD)
cd android
./gradlew assembleDebug
```

electron