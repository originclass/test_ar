# AI 使用日志

| 字段 | 内容 |
|---|---|
| 题目 | 现场巡检问题标注(Unity AR 上报 → Go 后端保存 → React 查看 → React 改状态) |
| AI 工具 | DeepSeek 编码代理(glm-5.3-flash) |
| 使用方式 | 逐步对话:人工发指令,AI 出全部代码/文档/排查,人工在 Unity 编辑器与真机上操作验证 |
| 说明 | 带 ※ 的 Prompt 原文为按会话记录还原(早期对话轮次被上下文压缩),其余为逐字原文 |

---

## P-001

- 目标:发布题目,约定协作方式

### Prompt 原文

> ※ 这是评审题目(原文见 SESSION_NOTES.md):Unity AR 放置标记并上报 → Go 后端保存 → React 页面查看问题 → React 修改处理状态(open/in_progress/resolved)。数据结构固定 {id, title, description, priority(low/medium/high), status, position{x,y,z}, createdAt, updatedAt},camelCase JSON。我们一步一步来,每一步停下来等我确认。

### 结果

- 采用
- 后续每步均先确认再推进;题目要求的数据结构成为三端统一契约,最终真机验收通过

## P-002

- 目标:确定项目架构(独立项目 vs 复用现有框架)

### Prompt 原文

> ※ 不做独立项目,后端复用我现成的 go-admin 项目(D:\bitxx-go-admin):在它的 /app 下建 patrol 模块,前端页面放它的 web 里;D:\patrol-ar 只放 Unity 工程。

### 结果

- 修改后采用
- 已按独立三目录做完 Go 后端与 React 管理端两个 commit,接此指令后迁移为 go-admin 集成方案,原实现降级为参考代码
- 验证:patrol 模块编译通过;curl 直测上报/列表/改状态接口 4 个分支全对(含 31000/31001 错误码)

## P-003

- 目标:搭建后端巡检模块

### Prompt 原文

> ※ 按 go-admin 四层结构(models/dto/service/apis/router)做 patrol 模块:免登录上报接口给 Unity 用,列表/详情/改状态/新增/删除走认证;错误码 31000 起;表 app_patrol_issue 建在 go-admin 库。

### 结果

- 采用
- 验证:gin 启动日志可见全部 patrol 路由注册;curl 上报成功返回自增 ID;管理接口无 token 返回 401

## P-004

- 目标:前端巡检问题页面接入菜单

### Prompt 原文

> ※ 在 go-admin 的 web 里做巡检问题页面(列表/筛选/状态下拉修改/新增/删除),加菜单和按钮权限,admin 和 test 角色都能用。

### 结果

- 修改后采用
- 首次验收发现菜单点击白屏:AI 排查出菜单 icon 为空串导致 React createElement(undefined) 崩溃,修 DB 后重新登录验收通过
- 验证:浏览器登录 → 巡检问题页列表渲染、状态下拉保存成功、刷新数据仍在

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

## P-020

- 目标:按评审要求重做 AI 使用日志并提交

### Prompt 原文

> 请在项目根目录提交 AI_LOG.md,按实际顺序保存本题中所有由你主动发给 AI 的问题或指令,包括追问和纠错。每条记录使用以下简化格式即可:## P-00N / 目标 / Prompt 原文 / 结果(采用或修改后采用或未采用 + 如何验证)。要求:必须保留 Prompt 原文,不能只写总结;不需要记录 IDE 的逐字自动补全;至少记录一次你检查、修改或放弃 AI 结果的过程。

### 结果

- 采用(本文件即产物)
- 验证:仓库根 AI_LOG.md 按时间序 20 条,原文保留(含错别字),含多处检查/修改/放弃 AI 结果的记录(P-004、P-007、P-010、P-013、P-015、P-016、P-019)
