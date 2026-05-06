# RLMusic 项目指令

## 计划模式（Plan Mode）

对于涉及以下任一条件的任务，**必须先进入计划模式**（使用 `/plan` 命令），在获得用户明确批准后方可实施：

- **修改任何文件**（包括创建、删除）
- 涉及数据库 Schema 变更
- 涉及 API 接口变更（前后端契约）
- 引入新的依赖包
- 重构核心模块（store、router、api 层）
- 影响用户数据或登录状态的变更

例外（无需计划模式可直接执行）：
- 纯信息查询（代码解释、文件查找、日志查看）
- 运行已有的测试/构建命令
- 查看或读取文件内容而不修改

### 计划内容要求

计划必须包含：
1. **修改的文件清单** — 每个文件的修改/创建/删除操作
2. **关键代码变更** — 涉及的具体函数、组件、接口
3. **回滚方案** — 如果变更失败如何恢复
4. **测试验证项** — 需要验证的功能点

### 计划批准后

- 按步骤逐一实施
- 每完成一步提交一次 commit
- 遇到未预见到的问题时暂停并重新请求指导

## 技术栈

- 前端：Vue 3 + TypeScript + Pinia + Naive UI + Vite
- 后端：Go + Gin + GORM 
- 构建：Electron（桌面端），capacitor(移动端)

## 代码规范

- 所有新代码使用 TypeScript，避免 `any` 类型
- 后端 API 返回统一格式：`{ code, message, data }`
- 错误码定义在 `src/utils/request.ts`（前端）和对应常量文件（后端）
- 图片/封面 URL 使用 `resolveCoverUrl()` 统一处理



## 提交规范

遵循 Conventional Commits：
```
feat(scope): 描述
fix(scope): 描述
docs(scope): 描述
refactor(scope): 描述
```
