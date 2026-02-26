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

electron部署



添加初始化配置页（首次启动向导）
对于 electron 服务器端：配置后端端口、对外访问地址（前端/后端）、基础文件夹路径；仅首次启动弹出并做连通性校验
对于 electron 客户端：配置后端地址/端口；仅首次启动弹出


对于electron服务端，管理员界面支持修改后端端口、对外访问地址（前端/后端）、基础文件夹路径

当无法访问后端服务时，跳转进入配置后端服务页，提示用户配置后端服务器地址




config.yml会被导入到electron包里吗，portable软件是如何存储数据的


优化portable client的连接服务器页




config.yml 是否会被导入到 Electron 包里

服务器端包会导入：在 server 构建配置里明确把 server/config.yml 作为 extraResources 拷贝到安装包资源目录：electron-builder.server.json
客户端包不导入：client 构建配置里没有 config.yml 的 extraResources 条目：electron-builder.client.json
Portable 版本如何存储数据

Electron 默认数据仍写入系统用户目录 app.getPath("userData")，与是否 portable 无关；除非你显式改写 userData 路径
你当前已引入的运行时配置写入就是 app.getPath("userData")/app-config.json：main.ts
后端数据库与封面数据实际由 config.yml 里的 BasicPath 决定，SQLite DSN 会落到 BasicPath/FileName/data 下：config.go
结论：portable 并不会自动“随包存储数据”，数据仍在用户目录或你配置的 BasicPath；如果希望“随包存储”，需要把 BasicPath 指到应用目录（或改写 userData）
评估：管理员界面支持修改后端端口 / 对外访问地址 / 基础文件夹路径的合理性

合理，但需要明确约束与流程
后端端口：服务端读取端口在启动时生效，修改后需重启服务进程才会生效（当前启动逻辑是启动时读取配置）
对外访问地址（前端/后端）：用于 LAN 暴露与二维码生成，属于“展示与接入控制”，放在管理员界面是合理的
基础文件夹路径：直接影响数据库与封面目录，修改后应提示迁移/重建或重新初始化，放在管理员界面合理但必须加“风险提示与校验”
建议的最小约束（保持合理性）

端口修改后提示“需要重启服务”并提供一键重启
基础路径修改需校验可写并提示对已有数据的影响
对外访问地址修改后同步更新二维码/局域网入口的显示
如果你希望，我可以把这三项直接做进管理员界面并加上重启与校验流程。