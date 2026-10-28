# AI 协作日志(P-001)

| 字段 | 内容 |
|---|---|
| 任务编号 | P-001 |
| 任务描述 | Unity AR 放置标记并上报 → Go 后端保存 → React 页面查看问题 → React 修改处理状态 |
| 使用的 AI 工具 | DeepSeek(编码代理,glm-5.3-flash) |
| 使用时段 | 2026-10-26 ~ 2026-10-28 |
| 生成内容占比 | 约 85%(全部业务代码、SQL、Unity 脚本、文档均由 AI 生成;人工完成:架构方案确认、Unity 编辑器操作、真机操作与验收) |

## 一、任务拆解

AI 将题目拆解为 6 个阶段,逐步推进、每阶段经人工确认后进入下一步:

1. **架构决策**:题目允许自由选型。初期按"独立项目"推进(backend/web/unity 三目录,已完成 Go 后端与 React 管理端两个 commit);人工提出"为什么不复用现有的 go-admin 框架项目"后,AI 对比两方案,确认**迁移到 go-admin 集成**——业务代码作为框架内独立模块 `/app/app/patrol`,前端作为框架 React 项目的 `/api/app/patrol/patrol-issue` + `/views/app/patrol/patrol-issue` 页面。原独立实现保留为参考。
2. **数据库**:`app_patrol_issue` 表(position 拍平为 pos_x/y/z 三列;API 层保持嵌套对象,模型层转换)。
3. **后端模块**:models/dto/service/apis/router 四层,错误码 31000-31004,免认证上报 + 认证管理双路由组。
4. **管理端**:ProTable 列表 + 状态下拉(PUT 修改)+ 新增弹窗 + 删除;菜单/权限/casbin SQL(目录"巡检管理" + 页面"巡检问题" + 5 个按钮权限,admin/test 双角色授权)。
5. **Unity 客户端**:AI 产出全部 C# 脚本与保姆级搭建文档,人工在 Unity 编辑器操作(建工程、装包、搭场景、接 UI、打 APK),AI 通过无线 adb 远程参与真机调试。
6. **真机联调**:AI 主导排查 4 个连环问题(见踩坑),最终三端链路真机验收通过。

## 二、关键技术点(AI 实现)

- **数据结构契约**:三端统一 camelCase JSON;`position{x,y,z}` 嵌套对象在 gorm 层拆为三列存储。
- **参数校验**:标题为空 → 31000;优先级枚举外 → 31001;状态枚举外 → 31002;唯一 ID/默认 open/时间戳均后端生成。
- **鉴权分层**:Unity 上报走免登录路由(题目硬性验收点);管理端全接口 JWT + casbin(admin 角色代码级豁免,test 角色由 SQL 策略授权)。
- **Unity 交互细节**:`EventSystem.IsPointerOverGameObject` 拦截 UI 上的点击(题目要求 6);双输入系统兼容(条件编译适配 Input System 与 legacy);Toast + 常驻引导语双提示机制;标记颜色语义(橙=待提交,绿=已提交)。
- **AI 远程调试**:手机无线 adb 配对后,AI 直接完成 APK 安装、启动、logcat 抓取、数据库验证,人工只负责"看屏幕现象 + 点手机"。

## 三、踩坑记录(全部由 AI 通过日志/清单解析定位)

| # | 现象 | 根因 | 修复 |
|---|---|---|---|
| 1 | APP 打开全黑 | APK 清单缺 CAMERA 权限与 ARCore meta-data(XR Plug-in Management 未在 Android 标签勾选 ARCore) | 勾选后重新打包 |
| 2 | 打包报错 minSdk 29 | ARCore+Vulkan 要求 API≥29 | Minimum API Level → 29 |
| 3 | 相机正常但 AR 会话异常,logcat 刷 `InvalidOperationException: ARCommandBufferSupportRendererFeature must be added` | Unity 6 + URP + Vulkan 下 ARCore 新增要求:渲染器必须挂 AR Command Buffer Support 特性 | Mobile/PC 两个 URP Renderer 资产均添加 |
| 4 | 相机画面正常但点平面无反应 | `activeInputHandler=1`(仅新 Input System),脚本用的旧 `Input.GetTouch` 永远读不到触摸 | 脚本改双系统兼容(条件编译),Player 设置改 Both |
| 5 | 提交无反应,后端零请求日志 | Windows 防火墙默认拦 8888 入站(手机 ping 通但 TCP 被吞) | 防火墙放行 TCP 8888 |
| 6 | 网络通但 APP 请求仍失败,后端仍零日志 | Android 9+ 默认禁止明文 HTTP;Unity 打包默认 `Allow downloads over HTTP=Auto` | 改 `Always allowed`(清单加 usesCleartextTraffic) |
| 7 | 前端菜单点了没页面/整页白屏崩溃 | ① 框架对菜单项无条件 `addIcon(icon)`,菜单图标留空 → React createElement(undefined) 崩溃;② Vite dev server 启动早于新页面文件创建,`import.meta.glob` 路由表不含新页面 | 菜单图标补 FileDoneOutlined;重启 dev server |
| 8 | gin 启动 panic 重复注册 | 同一 method+path 在免认证/认证两组重复注册 | 免认证上报改独立路径 `/issue/report` |

## 四、验收记录

- **功能**:真机扫描出水平平面(淡蓝可视化)→ 点击放置橙色标记 → 表单填写(标题必填校验)→ 提交成功 Toast(含后端生成的 ID)→ 管理端列表实时可见 → 状态修改 → 后端重启数据仍在。
- **数据**:数据库 5 条记录,4 条来自真机提交(坐标为 AR 真实追踪值),1 条网页改状态为 resolved。
- **错误路径**:标题为空/优先级非法(Unity 前端拦截 + 后端 31000/31001 双保险)、无 token 访问管理接口 401、服务不可达时 Unity 显示具体失败原因。
