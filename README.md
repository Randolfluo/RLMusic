# 基于VUE的音乐播放器




### 🛠️ 环境与部署 (Environment & Deployment)
- [ ] Docker 部署：通过挂载目录实现目录共享


- [ ] 桌面歌词：创建背景透明 (`transparent: true`) 且忽略鼠标事件 (`ignoreMouseEvents`) 的置顶窗口

















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
pnpm exec cap open android
方式 B：命令行构建 (适合 CI/CD)
cd android
./gradlew assembleDebug
```

electron部署
```
pnpm build:client
pnpm build:server
```


添加初始化配置页（首次启动向导）
对于 electron 服务器端：配置后端端口、对外访问地址（前端/后端）、基础文件夹路径；仅首次启动弹出并做连通性校验
对于 electron 客户端：配置后端地址/端口；仅首次启动弹出


对于electron服务端，管理员界面支持修改后端端口、对外访问地址（前端/后端）、基础文件夹路径

当无法访问后端服务时，跳转进入配置后端服务页，提示用户配置后端服务器地址




