# AI 使用日志

| 字段 | 内容 |
|---|---|
| 题目 | 现场巡检问题标注(Unity AR 上报 → Go 后端保存 → React 查看 → React 改状态) |
| AI 工具 | DeepSeek 编码代理(glm-5.3-flash) |
| 使用方式 | 逐步对话:人工发指令,AI 出全部代码/文档/排查,人工在 Unity 编辑器与真机上操作验证 |
| 说明 | 带 ※ 的 Prompt 原文为按会话记录还原(早期对话轮次被上下文压缩),其余为逐字原文 |

---
## P-000.1

- 目标:本地创建数据库，工供项目使用

### Prompt 原文

> ※ 给我在我电脑的pgsql里面建一个数据库叫go-admin，有什么问题直接问我，我的数据库密码是123456，用户名是postgres

### 结果
- 数据库 go-admin 已创建，编码 UTF8
- PostgreSQL 版本：18.6,x64 Windows
- 已验证可正常连接

---
## P-000.2

- 目标:拉取二开框架项目源码

### Prompt 原文

> ※ https://github.com/bitxx/go-admin这个项目用git拉取下来到bitxx-go-admin里面，遇到问题直接问我

### 结果
- 分支：master,最新提交:ce17ea9 Merge branch 'dev'
- 位置:D:\bitxx-go-admin

## P-000.3

- 目标:执行项目里的sql脚本

### Prompt 原文

> ※ app_pgsql.sql脚本在建的数据库里面执行一下

### 结果
- 建表数据添加成功


## P-000.4

- 目标:根据项目的文档，启动后端

### Prompt 原文

> ※ <项目启动提示> 以上是该项目的启动方式，请启动后端
### 过程中解决的三个坑(记下来以后有用)：

- 生成的代码在 Windows 下非法:go generate 把路径写成 go-admin\app\...(反斜杠进了 Go 字符串)，已把 572 处 \ 修正为 /
- ⚠️ 文档里 -c=config/settings.dev.yml 的写法有 bug:这版 cobra 解析 -c=x 时路径会被截断成 config/settings,必须用空格分隔:-c config/settings.dev.yml
- 沙箱限制 Go 写 C:\Users\zaq\go 和 go env -w,所以 go generate 里的两条 go env -w 跳过了(已用环境变量 GOPROXY/GOPATH/GOCACHE 等效替代)
### 结果
- 后端启动成功

## P-000.5

- 目标:启动前端

### Prompt 原文

> ※ 准备启动前端
### 过程中解决的坑(记下来以后有用)：

- pnpm 12 崩溃(thread has overflowed its stack)—— Windows 上的已知问题，改用 npm
- npm 严格 peer 依赖冲突 —— 老项目标配，加 --legacy-peer-deps
- postinstall 脚本 EPERM —— 沙箱禁止子进程管道，用 --ignore-scripts 跳过(Vite 6 的 esbuild 二进制通过平台包分发，不依赖脚本)
- esbuild 报 "debugger" 非法标识符 —— vite.config.ts 把 "debugger" 关键字放进了 pure 选项，新版 esbuild 不允许。已通过在 .env.development
- 里加 VITE_DROP_CONSOLE = false 修复(开发环境保留 console/debugger 反而更利于调试)
.env 热重启不生效 —— 该项目的 getEnv.ts 会把环境变量写进 process.env,Vite 热重启时进程内旧值优先于文件，必须整个进程重启
### 结果
- 前端启动成功

## P-001

- 目标:发布题目,约定协作方式

### Prompt 原文

> ※ 这是评审题目....:Unity AR 放置标记并上报 → Go 后端处理业务逻辑 → React 页面查看问题 → React 修改处理状态(open/in_progress/resolved)。数据结构固定 {id, title, description, priority(low/medium/high), status, position{x,y,z}, createdAt, updatedAt},camelCase JSON。我们一步一步来先从写后端业务接口开始,然后搭建前端页面，最后做unity客户端，这是大体的，一个个分步解决，每次做完停下来等我确认。

### 结果

- 采用
- 后续每步均先确认再推进;题目要求的数据结构成为三端统一契约,最终真机验收通过

## P-001.1

- 目标:先把本地的unity项目的代码库建立起来

### Prompt 原文

> ※初始化建立unity的本地项目文件夹，git init.，关联我的项目仓库....编写必要的README.md,等等

### 结果

- 首次提交在本地，还没有 push 到 GitHub —— 等项目有实质内容后(比如第 1 步后端完成)再一起推，避免空仓库推来推去；推送时如果弹出 GitHub 登录授权，按提示操作即可。
- 你机器上现在有两个库:go-admin(之前的项目)和 patrol_ar(本次巡检项目)，互不干扰。


## P-001.2

- 目标:编写后端代码

### Prompt 原文

> ※ 开始写后端代码

### 结果
- 后端代码写完


## P-001.3

- 目标:编写前端代码

### Prompt 原文

> ※ 开始写前端代码

### 结果
- 钱端代码写完，前后端接口打通联调成功

## P-001.3

- 目标:重改项目结构

### Prompt 原文

> ※ 你为什么不在D:\bitxx-go-admin\web这个前端项目里面完成这个第2步前端需求？而是新启动一个项目，还有你后端是使用的是拉取的后端项目做基于她做二次开发的吗

### 结果
- agent向我确认信息


## P-001.3

- 目标:提供信息

### Prompt 原文

> ※ 我觉得你应该在后端单独建一个模块，用于实现本次的所有业务功能，也就是需要的各种后端接口后端就在/app/下面建一个模块就是，然后就是各层，就是各种增删改查接口，然后前端在之前启动的teact项目里，有页面，api对接后端返回的问题列表信息，根据问题 ID 修改处理状态（就是每个问题的状态可被更新，调用后端接口更新数据库数据），还有几个接口是用于客户端的，比如（2. 接收并保存客户端 上报的问题；也就是往失数据库插入数据，为新问题生成唯一 ID，默认状态为 open，客户端提交问题时候6. 标题为空或枚举值错误时，返回明确的参数错误；）对吧，这些后端接口就是增删改查，前端有一个页面，有接口请求响应后端，D:\bitxx-go-admin\web\src\api\app放新模块的前端接口，D:\bitxx-go-admin\web\src\views\app接受后端数据展示并且可以修改问题的状态就是
### 结果
- 修改项目结构，成功迁移业务代码，后端接口跳通


## P-002

- 目标:修改前端出现的页面没有显示问题

### Prompt 原文

> ※ 我没有在前端页面看到问题这个功能的页面啊，你的前端页面写了吗，前后端接口联调了嘛

### 结果
页面写了、接口也联调验证过了 —— 问题很可能出在菜单是登录时一次性拉取并缓存在本地的：你当前这个登录会话是在我插入菜单数据之前建立的，所以侧边栏里没有新菜单。我先验证这个推断，查一下前端菜单的加载与缓存机制：


找到了。LoginForm.tsx 第 53 行：菜单树是登录那一刻拉取并存进 redux 的(getMenuRoleApi → setRouteList),你当前会话建立于我插菜单数据之前，所以侧边栏没有。确认一下 GetMenuRole 的角色过滤逻辑，排除其他可能：


原因确认了，页面写了、接口也联调通过，你没看到菜单的原因是：菜单数据是登录那一刻从后端拉取并缓存在前端的(LoginForm.tsx 登录成功后调 getMenuRoleApi 存入 redux),而你的这个登录会话建立在我插入菜单数据之前 —— 后端当时的菜单树里自然没有“巡检问题”。

后端代码我核实过：admin 角色的菜单是每次登录实时查库返回的(不走服务端缓存)，所以解决方式非常简单：

## P-003

- 目标:修改前端

### Prompt 原文

> ※ 我点击巡检管理直接没有任何页面，为什么呢

### 结果
我逐层核实了一遍，页面代码、菜单数据、接口全链路都是好的：

检查项	结果
菜单树(后端实时查询)	✅ App应用 → 巡检管理(135) → 巡检问题(136) 结构正确，排序正确，admin 全量可见
页面组件文件与路径映射	✅ 菜单 element=/app/patrol/patrol-issue/index → /src/views/app/patrol/patrol-issue/index.tsx 精确对应，与能正常打开的“用户等级”页面机制完全一致
页面所有 import(动态图标/分页配置/message/modal/工具函数)	✅ 逐个核对了真实导出，无缺失
Vite 编译	✅ 页面、弹窗、API 三个模块编译均 200 无错误
后端接口	✅ 上报/校验/鉴权四分支此前已实测通过
找到的问题： Vite 的路由表(import.meta.glob)是在 dev server 启动那一刻扫描的 —— 我的服务器是在写完页面文件之前启动的，运行中扫不到“后创建的文件”，所以菜单点了路由却找不到组件 → 空白。前端已经用新路由表重启完毕(http://localhost:1688)。
## P-004

- 目标:继续排查前端问题

### Prompt 原文

> ※ 给出浏览器的报错信息

### 结果
根因已修复。 你控制台里那行致命报错

代码块
React.createElement: type is invalid ... got: undefined
  at LayoutMenu (Menu/index.tsx)
正是侧边栏在渲染菜单图标时炸的：框架对每个菜单项无条件执行 addIcon(item.icon),而我插入“巡检问题”(菜单 136)时图标留了空字符串 → createElement(undefined) → React 整树崩溃，表现就是点击后白屏/无页面。(其他几条只是框架原有的废弃 API 警告，与本次无关。)
## P-005

- 目标:明确 Unity AR 客户端需求

### Prompt 原文

> 真机运行、检测水平平面、点击平面放标记、填写并提交问题(标题必填/描述可选/优先级三档/位置坐标)、提交成功或失败明确提示、点按钮或输入框不误放标记;标记用基础几何体即可

### 结果

- 采用
- AI 产出 3 个 C# 脚本 + 8 步保姆级搭建文档,人工照文档完成工程搭建
- 验证:真机全流程验收(见 P-016)

## P-006

- 目标:解决 Unity 脚本报错 CS0234/CS0246(找不到 ARFoundation 命名空间)

### Prompt 原文

> ※ 脚本复制进工程后一堆红色报错 CS0234/CS0246,Unity 里 ARFoundation 命名空间不存在。

### 结果

- 采用
- 原因:先复制了脚本、还没装 AR Foundation 包
- 验证:Package Manager 装 AR Foundation + ARCore XR Plugin 6.3.5 后编译通过

## P-007

- 目标:解决编译错误 CS0103:'TrackableTypes' does not exist

### Prompt 原文

> ※ 还是有报错 CS0103:当前上下文中不存在"TrackableTypes"。

### 结果

- 修改后采用
- AI 首版按 AR Foundation 5.x 写的复数枚举,6.x 改为单数 TrackableType.PlaneWithinBounds
- 验证:改后 Unity Console 无编译错误

## P-008

- 目标:解决打包报错(minSdk 不满足 ARCore+Vulkan 要求)

### Prompt 原文

> ※ 打包失败:ARCore Required apps using Vulkan require a minimum SDK version of AndroidApiLevel29. Currently set to AndroidApiLevel25。

### 结果

- 采用
- Player → Minimum API Level 改 Android 9.0(API 29)
- 验证:APK 构建成功
### 中间跟着agent给出的教程在unity搭建场景，配置AR相关模块，打包构建apk的步骤省略，这些均为手动操作
### 以下都是跟AI询问，apk打开黑屏，没有相机，没有标记区域，点击标记区域不出现自己写的UI的，表单数据无法上传，接口不通，采用手机无线调试根据日志排查问题，都在agent指导下给出解决方案，最终完成功能模块，对话轮数较多，不一一展示了。
## P-009

- 目标:排查真机 APP 打开黑屏(第 1 版 APK)

### Prompt 原文

> ※ 装到手机上了,打开是全黑的,什么都没有。

### 结果

- 修改后采用
- AI 让人导出 APK 用 aapt 检查清单:缺 CAMERA 权限与 ar.core meta-data → XR Plug-in Management 的 Android 标签勾选 ARCore 后重打
- 验证:新清单含权限与 meta-data;aapt dump 确认

## P-010

- 目标:排查真机黑屏(第 2 版 APK,配置已全对)

### Prompt 原文

> ※ 这回权限那些都配了,打包也过了,手机上还是黑屏。

### 结果

- 修改后采用
- AI 用无线 adb 抓 logcat,发现 InvalidOperationException:Vulkan 下必须给 URP Renderer 加 AR Command Buffer Support 特性 → 两个渲染器资产都加上
- 验证:重打后 logcat 该异常计数为 0,相机画面出现

## P-011

- 目标:解决电脑连不上手机(USB 调试授权反复失败)

### Prompt 原文

> ※ USB 插上手机一直弹授权框,点了允许 adb devices 还是空的,系统里能看到一加 Ace 3 Pro。

### 结果

- 采用
- 改用无线调试:手机开发者选项生成配对码,adb pair <ip>:<port> <code> 配对成功
- 验证:adb devices 显示 device;后续安装/启动/日志全部远程完成

## P-012

- 目标:排查"相机画面出来了但看不到 UI"

### Prompt 原文

> ※ 加了 AR Command Buffer Support 重新打包安装后,相机画面已经出来了(黑屏解决了),但是没有我做的 UI 界面。

### 结果

- 修改后采用
- AI 先解释了设计行为(提示文字 3 秒消失、表单面板默认隐藏),并把提示改为常驻引导;事后证明另有输入系统问题(见 P-013)
- 验证:logcat 确认 AR 会话无异常;常驻提示改动编译通过

## P-013

- 目标:排查"点了淡蓝平面但不弹表单"

### Prompt 原文

> 淡蓝色出现后我点他,还是没有出现提交相关UI

### 结果

- 修改后采用
- AI 查 ProjectSettings 发现 activeInputHandler=1(仅新 Input System),脚本用的旧 Input.GetTouch 永远读不到触摸 → 脚本改双系统兼容(条件编译 EnhancedTouch/GetTouch),并指示 Player 设置改 Both
- 验证:重新打包后点平面弹出橙色方块 + 表单(下一轮用户确认)

## P-014

- 目标:真机提交联调(装新 APK、起前后端)

### Prompt 原文

> 现在点击蓝色出现UI了, @C:\Users\zaq\Downloads\we3.apk 这是现在的apk,来测试下提交问题,启动前后端

### 结果

- 采用
- AI 无线 adb 安装 APK、启动 go-admin 后端(8888)与前端(1688)、记录数据库基准(当时 1 条)
- 验证:adb install Success;后端 /captcha 200;APP 启动 pid 正常

## P-015

- 目标:排查"手机提交了但后台没有数据"(第 1 次)

### Prompt 原文

> 我点击提交,但是没看到后台添加了数据啊

### 结果

- 修改后采用
- AI 三路排查(数据库/手机日志/后端日志):后端只有浏览器请求,手机请求全部未达 → Windows 防火墙拦 8888 入站;UAC 提权加放行规则后,手机 curl 后端成功返回 JSON
- 验证:netsh 显示规则 go-admin-8888 Enabled;手机侧 curl 拿到 {"code":200}

## P-016

- 目标:排查"提交无提示、无数据"(第 2 次,防火墙已通)

### Prompt 原文

> 我点击提交l,没有提示成功,但是也没看到数据库新增数据

### 结果

- 修改后采用
- 决定性实验:手机 curl 直调上报接口能穿透(网络已通)→ 锁定 APP 内部;aapt 查 APK 清单发现缺 usesCleartextTraffic —— Android 9+ 默认禁明文 http → Player 设置 Allow downloads over HTTP 改 Always allowed,重打
- 验证:重打后真机提交成功,Toast 显示"提交成功!问题编号 #N";数据库新增 4 条真机记录(带 AR 真实坐标);管理端列表同步可见

## P-017

- 目标:提交 GitHub

### Prompt 原文

> 这下都成了,你确实很厉害,这些技术我从没接触过就做了一个app,来提交github吧

### 结果

- 修改后采用
- AI 产出:README 重写(架构/接口/启动)、AI_LOG、go-admin 集成包、Unity 工程源码,分 4 个主题 commit;直连 GitHub 被重置,改走本机代理(7890)推送成功
- 验证:git rev-parse 本地 = origin/main;GitHub 上 7 个提交可见

## P-018

- 目标:获取前后端启停命令

### Prompt 原文

> 这个后端和前端的关停和启动命令是什么

### 结果

- 采用
- AI 给出启动(后端 exe + npm run dev)/关停(按端口杀进程)/健康检查三组命令
- 验证:同会话内多次用这些命令起停服务,均有效

## P-019

- 目标:纠正交付不完整(必须有完整可运行代码)

### Prompt 原文

> 不对啊,我看你提交的代码,为什么只有这次做的这个功能的这些代码,我需要完整的前后端代码都提交,包括完整的unity客户端的代码啊,你只是提交部分,别人拉取下来怎么运行呢

### 结果

- 未采用(首次交付范围)→ 修改后重新交付
- AI 首次只提交了"新增模块代码",被否定;随后把 go-admin 框架完整源码(前后端,排除 node_modules/exe/缓存)并入仓库,删除冗余集成目录,补 db_schema.sql 与菜单 SQL 到工程根,加 PROJECT_NOTES.md 运行说明,README 同步
- 验证:471 文件入库并推送;核对仓库含 app/app/patrol 全部四层、web 巡检页面、两份 SQL;Unity 工程确认此前已完整(Assets/Packages/ProjectSettings,Library 为可重建缓存属标准排除)
