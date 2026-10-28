# Unity AR 巡检客户端 — 保姆级搭建文档

> 你不需要写任何代码。脚本已全部写好,放在 `Assets/Scripts/`(3 个文件):
> - `PatrolApi.cs` —— 调后端接口上报问题(已对接 go-admin)
> - `TapToPlaceReporter.cs` —— 点平面放标记 + 弹表单 + 提交 + 提示(核心逻辑)
> - `ArPlaneAnchor.cs` —— 让检测到的平面显示为淡蓝色、可点击
>
> 目标设备:支持 ARCore 的安卓手机(设备列表见 support.google.com/ar)
> 开发机要求:Unity 2022 LTS(含 Android Build Support)

---

## 第 1 步:新建工程(10 分钟)

1. 打开 Unity Hub → New project
2. 模板选 **3D (URP)**(没有就选 3D)
3. 工程名:`PatrolAR`,位置随意 → Create
4. 菜单 **File → Build Settings** → 平台列表选 **Android** → 右下 **Switch Platform**(等右下角进度条转完)
5. 把本目录 `Assets/Scripts/` 下的 3 个 `.cs` 文件复制进工程的 `Assets/Scripts/`(文件夹不存在就新建)。Unity 会自动编译,底部 Console 不应有红色报错

## 第 2 步:安装 AR 包 + 开启 ARCore(10 分钟)

1. 菜单 **Window → Package Manager** → 左上角下拉选 **Unity Registry**
2. 搜索 `AR Foundation` → Install(5.x 版本,2022 LTS 一般自带)
3. 再搜索 `ARCore XR Plugin` → Install(版本和 AR Foundation 一致)
4. 菜单 **Edit → Project Settings** → 左侧 **XR Plug-in Management**:
   - 第一次进入会提示安装,点 Yes
   - 切到顶部的 **Android 机器人图标** 标签页 → 勾选 **ARCore**
5. 左侧 **Player**(就是其他设置的入口):
   - **Other Settings → Identification → Minimum API Level** = **Android 7.0 (API 24)**
   - **Other Settings → Configuration → Internet Access** = **Require**
   - **Other Settings → Configuration → Active Input Handling** = **Both**(若选项存在;改完 Unity 会重启)

## 第 3 步:搭场景(20 分钟)

Hierarchy 里默认有 `Main Camera` 和 `Directional Light`。

1. **删掉 Main Camera**(AR 会用自己的相机)
2. 菜单 **GameObject → XR**:
   - 点 **AR Session**(AR 系统:追踪、相机权限)
   - 点 **AR Session Origin**(AR 相机 + 管理器容器)
   > 若菜单里没有 XR 子项:Hierarchy 左上 `+` 号 → XR 里找;再没有说明 AR Foundation 没装好,回第 2 步
3. **AR Session Origin** 的 Inspector 确认有 3 个组件:`AR Session Origin`、`AR Plane Manager`、`AR Raycast Manager`(5.x 若 Raycast Manager 不自动存在:点 Add Component 搜 `AR Raycast Manager` 添加)
4. **AR Plane Manager** 组件上:
   - **Detection Mode** = **Horizontal**(只检测水平面,对应题目要求 2)
5. 制作平面预制体(让平面可见可点):
   1. 菜单 **GameObject → XR → AR Default Plane**,在 Hierarchy 里生成一个 "AR Default Plane"
   2. 把它从 Hierarchy **拖进 Project 的 Assets 文件夹** → 变成预制体(蓝色图标)
   3. 双击进入这个预制体,点 **Add Component** → 搜 `ArPlaneAnchor` 添加(就是我们的脚本)
   4. 回到场景,把 Hierarchy 里那个 "AR Default Plane" **删掉**(预制体已在 Assets 里)
   5. 选中 Hierarchy 的 **AR Session Origin** → AR Plane Manager 的 **Plane Prefab** 槽 → 把 Assets 里的 AR Default Plane 预制体拖进去
6. 放置交互物体:
   1. Hierarchy 右键 → **Create Empty**,命名 `ArInteraction`
   2. 选中它 → Add Component → 搜 `TapToPlaceReporter` 添加

## 第 4 步:搭 UI(30 分钟,数值照抄)

在 `ArInteraction` 上,我们需要:1 个提示文本 + 1 个填写面板(标题框/描述框/优先级下拉/提交/取消)。

1. Hierarchy 右键 → **UI → Canvas**,命名 `ReportCanvas`
   - Canvas Scaler 组件:UI Scale Mode = **Scale With Screen Size**,Reference Resolution = **1080 x 1920**
2. 在 ReportCanvas 上右键 → **UI → Text - TextMeshPro**?——**不要**,选 **UI → Legacy → Text**,命名 `StatusText`,RectTransform 设置:
   - 锚点顶部居中,Pivot (0.5, 1),Pos Y = -60;Width 900,Height 80
   - Font Size 36,Alignment 居中,Color 白色
3. 在 ReportCanvas 上右键 → **UI → Image**,命名 `ReportPanel`:
   - 锚点底部居中,Pivot (0.5, 0),Pos Y = 120;Width 950,Height 700
   - Image 的 Color 改成黑色、Alpha 200(半透明黑底)
4. 右键 ReportPanel → **UI → Legacy → Input Field**,命名 `TitleInput`(占位符 Placeholder 文本改成 `问题标题(必填)`,Font Size 30)
5. 右键 ReportPanel → **UI → Legacy → Input Field**,命名 `DescInput`:
   - Content Type = Standard;把它的 Text 子物体勾上 **Multi Line**(或者直接在 Input Field 组件勾 Multi Line)
   - Placeholder 文本改成 `问题描述(可选)`
6. 右键 ReportPanel → **UI → Legacy → Dropdown**,命名 `PriorityDropdown`:
   - 组件里的 Options 列表只保留 3 项,内容分别改成:`low`、`medium`、`high`(**顺序必须是 low 在第 0 位**)
7. 右键 ReportPanel → **UI → Button - TextMeshPro**?——**不要**,选 **UI → Legacy → Button**,命名 `SubmitButton`,子文本改成 `提交`
8. 同样再建一个 Button,命名 `CancelButton`,子文本改成 `取消`
9. 在 ReportPanel 里把 5 个控件摆整齐(建议纵向排:标题框 → 描述框 → 下拉 → 一行两个按钮)。会摆 UI 就随便摆,不会摆就右键 ReportPanel → **UI → Vertical Layout Group**,把 5 个控件直接拖进去自动排列
10. **把 ReportPanel 取消激活**:Inspector 顶上名字左边的勾选框 **去掉**(运行时由脚本弹出)

## 第 5 步:接线(5 分钟,最容易漏的一步)

选中 Hierarchy 的 `ArInteraction`,Inspector 里 `TapToPlaceReporter` 脚本有 9 个空槽,从 Hierarchy **把对应物体拖进槽里**:

| 槽位 | 拖什么 |
|---|---|
| Raycast Manager | AR Session Origin(它身上有该组件) |
| Plane Manager | AR Session Origin |
| Marker Prefab | **留空**(留空自动用橙色立方体) |
| Status Text | ReportCanvas/StatusText |
| Report Panel | ReportCanvas/ReportPanel |
| Title Input | ReportPanel/TitleInput |
| Description Input | ReportPanel/DescInput |
| Priority Dropdown | ReportPanel/PriorityDropdown |
| Submit Button | ReportPanel/SubmitButton |
| Cancel Button | ReportPanel/CancelButton |

拖完 **Ctrl+S 保存场景**(场景命名随意,如 `Main`)。

## 第 6 步:改后端地址(1 分钟,最关键的一步)

1. 打开 `Assets/Scripts/PatrolApi.cs`,第 22 行:
   ```csharp
   public static string BaseUrl = "http://192.168.10.26:8888";
   ```
2. 把 `192.168.10.26` 换成**你电脑的局域网 IP**:Windows 按 `Win+R` 输 `cmd` 回车,输入 `ipconfig`,找"无线局域网适配器 WLAN"下的 **IPv4 地址**
3. 前提:手机和电脑连**同一个 WiFi**;电脑上 go-admin 后端(`go-admin-api.exe server`)已启动

> 验证连通性:手机浏览器打开 `http://电脑IP:8888/admin-api/v1/captcha`,能看到一串 JSON 就通了。不通 → 检查 Windows 防火墙(临时关闭测试)/ 两个设备是否同网段。

## 第 7 步:打包安装(15 分钟)

1. 菜单 **File → Build Settings**:
   - Scenes In Build:Add Open Scenes(把 Main 场景加进来)
   - 点 **Build**,选个输出位置(如桌面),生成 `PatrolAR.apk`
2. 把 APK 传到手机(微信/QQ/数据线都行)→ 安装(需要允许"安装未知应用")
3. 手机**未开启开发者选项**也能装,只是安装时要授权一次
4. 首次启动会弹"相机权限" → 允许

## 第 8 步:真机验收清单(对照题目要求逐条勾)

| # | 题目要求 | 怎么验证 |
|---|---|---|
| 1 | 真机运行 | APP 打开,相机画面出现 |
| 2 | 检测水平平面 | 缓慢扫地面,出现**淡蓝色网格面** |
| 3 | 点平面放标记 | 点蓝色平面任意位置 → 出现**橙色小方块** + 弹出填写面板 |
| 4 | 填写提交 | 标题留空点提交 → 提示"标题不能为空";填好后提交 |
| 5 | 成功/失败提示 | 成功:方块变绿 + Toast"提交成功!问题编号 #N";后端没开/断网:红色错误提示(原因会写明) |
| 6 | UI 不误放标记 | 点提交/取消/输入框/下拉,**不会**在场景里生成新方块 |
| 7 | 数据互通 | 打开 PC 浏览器 http://localhost:1688 登录 → 巡检问题列表出现刚提交的记录、坐标一致 |

## 常见问题

| 现象 | 处理 |
|---|---|
| 打开 APP 提示"设备不支持 AR" | 手机不在 ARCore 支持列表,换手机 |
| 相机画面黑屏 | 没给相机权限,系统设置里允许 |
| 平面扫不出来 | 光线太暗/纯色地板;换有纹理的地面,缓慢移动 |
| 提交一直失败 | ① 后端没启动 ② IP 写错(必须局域网 IP,不能 127.0.0.1)③ 手机用了流量没连 WiFi ④ Windows 防火墙拦了 8888 |
| Unity Console 红字 `Active Input Handling` | 第 2 步第 5 条没做,改成 Both |
| 编译报 `ARPlaneAnchor` / `TapToPlaceReporter` 缺失 | 脚本没复制全,确认 Assets/Scripts 有 3 个 .cs |

---

全部勾完 → 回到电脑告诉我"Unity 端验收通过",进入第 4 步(三端联调报告 + 演示视频脚本)。
