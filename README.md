# 现场巡检问题标注(Patrol AR)

园区巡检人员用手机在现场放置 AR 标记并上报问题(墙面破损 / 设备异常 / 安全隐患),管理人员在网页管理后台查看问题列表并更新处理状态。

**业务闭环(已全部打通并真机验收):**

```
Unity AR(Android 真机)                Go 后端(go-admin 框架)            React 管理端
┌──────────────────────┐  POST /app/patrol/issue/report   ┌──────────────┐
│ ARCore 检测水平平面    │ ──────────────────────────────▶ │ 参数校验      │
│ 点击平面放置 AR 标记   │  (免登录,标题/优先级后端强校验)   │ 生成唯一 ID   │
│ 填写标题/描述/优先级    │                                 │ 默认 open    │
│ 提交 → 成功/失败提示   │ ◀─────── {code:200, data:问题ID} │ PostgreSQL   │
└──────────────────────┘                                 └──────┬───────┘
                                                                 │ JWT + casbin 鉴权
                                                          ┌──────▼───────┐
                                                          │ 巡检问题列表   │
                                                          │ 分页/筛选     │
                                                          │ 按 ID 改状态  │
                                                          └──────────────┘
```

## 交付内容

本仓库包含**可直接运行的完整代码**:Unity AR 客户端工程 + go-admin 框架完整源码(前后端,含巡检模块)+ 早期参考实现:

```
patrol-ar/
├── ar_poj/AR_prj/            # ★ Unity 6 AR 客户端工程(完整工程,真机验证通过)
│   ├── Assets/Scripts/       #    全部 C# 源码(3 个脚本,见下文)
│   ├── Assets/Scenes/        #    SampleScene(AR Session + 平面检测 + UI)
│   └── Assets/Settings/      #    URP 渲染配置(含 AR 渲染特性,见 AI_LOG 踩坑记录)
├── go-admin/                 # ★ go-admin 框架完整源码(前后端)+ 巡检模块(开箱即跑)
│   ├── app/app/patrol/       #    后端巡检模块(models/dto/service/apis/router 四层)
│   ├── web/                  #    React 管理端完整工程(含巡检问题页面)
│   └── PROJECT_NOTES.md      #    本工程运行说明(数据库初始化/启动/Unity 对接)
├── unity/                    # Unity 保姆级搭建文档(SETUP_GUIDE.md)与脚本副本
├── backend/                  # 参考实现:独立 Go 后端(prototype)
├── web/                      # 参考实现:独立 React 管理端(prototype)
└── docs/                     # 文档(搭建指南 / AI 协作日志)
```

## Unity AR 客户端(`ar_poj/AR_prj`)

| 文件 | 职责 |
|---|---|
| `Assets/Scripts/TapToPlaceReporter.cs` | 核心交互:点击平面放标记 → 弹表单 → 提交 → 成功/失败提示;双输入系统兼容(旧 Input / 新 Input System 均可);`IsPointerOverGameObject` 保证点 UI 不误放标记 |
| `Assets/Scripts/PatrolApi.cs` | 上报接口客户端:组装 JSON、`UnityWebRequest` POST、解析 `{code,msg,data}`、网络层/业务层错误分类转成用户可读提示 |
| `Assets/Scripts/ArPlaneAnchor.cs` | 平面可视化:检测到的水平平面渲染为淡蓝色半透明 |

**真机环境(OnePlus Ace 3 Pro / ARCore 1.57):已验收平面检测、放标记、表单提交、失败提示、UI 不误触发全部用例。**

### 上报接口约定(题目规定的数据结构)

```json
POST {host}:8888/admin-api/v1/app/patrol/issue/report
{
  "title":       "入口墙面破损",
  "description": "左侧墙体存在裂缝",
  "priority":    "high",
  "position":    { "x": 0.42, "y": 0.03, "z": 1.26 }
}
成功响应:{ "code": 200, "msg": "操作成功", "data": 2 }
```

- `title` 必填(空 → code 31000);`priority` 枚举 low/medium/high(非法 → 31001)
- 唯一 ID、默认状态 `open`、时间戳全部由**后端**生成
- 管理端同结构(camelCase),列表接口返回 `position` 嵌套对象

## 后端 + 管理端(`go-admin/`,开箱即跑)

完整 go-admin 工程(框架源码 + 巡检模块已就位),运行步骤详见 [`go-admin/PROJECT_NOTES.md`](go-admin/PROJECT_NOTES.md),概要:

```powershell
# 1. 数据库(PostgreSQL,库名 go-admin;先导入框架的 app_pgsql.sql 基础表)
psql -U postgres -d go-admin -f go-admin/db_schema.sql          # 巡检问题表
psql -U postgres -d go-admin -f go-admin/patrol_menu_init.sql   # 菜单 + 角色授权 + casbin

# 2. 后端 :8888
cd go-admin
go build -o go-admin-api.exe main.go
.\go-admin-api.exe server --config config/settings.dev.yml

# 3. 管理端 :1688
cd web && npm install && npm run dev
# 登录 admin/123456 → 菜单「App应用 → 巡检管理 → 巡检问题」
```

本项目在框架中新增的业务代码(其余均为框架本体):

| 路径 | 说明 |
|---|---|
| `app/app/patrol/models/patrol_issue.go` | gorm model:嵌套 `position` 对象 ↔ `pos_x/y/z` 三列转换;枚举常量与校验 |
| `app/app/patrol/service/dto/patrol_issue.go` | 请求 DTO:分页查询(search 标签)/ 上报 / 状态更新 |
| `app/app/patrol/service/patrol_issue.go` | 业务:Insert(标题/优先级校验、状态强制 open)、UpdateStatus(枚举校验)、GetPage/Get/Delete |
| `app/app/patrol/apis/patrol_issue.go` | gin 处理器(框架标准链式调用) |
| `app/app/patrol/router/*.go` | 免认证组(`POST /app/patrol/issue/report`,供 Unity 上报)+ JWT+casbin 认证组(管理端) |
| `config/base/lang/patrol.go` | 业务错误码 31000-31004 注册 |
| `web/src/api/app/patrol/patrol-issue/index.ts` | 前端 API 客户端 |
| `web/src/views/app/patrol/patrol-issue/` | 前端页面(列表/状态修改/新增/删除) |
| `patrol_menu_init.sql` | 菜单/按钮权限/admin+test 角色授权/casbin 规则 |
| `db_schema.sql` | `app_patrol_issue` 表结构 |

接口全景(题目 3.3 全覆盖):

| 方法 | 路径 | 鉴权 | 用途 |
|---|---|---|---|
| POST | `/admin-api/v1/app/patrol/issue/report` | 免登录(AR 客户端) | 上报问题,生成唯一 ID、默认 open |
| GET | `/admin-api/v1/app/patrol/issue` | JWT+casbin | 分页列表(标题/优先级/状态/时间筛选) |
| GET | `/admin-api/v1/app/patrol/issue/{id}` | JWT+casbin | 详情 |
| PUT | `/admin-api/v1/app/patrol/issue/{id}/status` | JWT+casbin | 按 ID 修改处理状态 |
| POST | `/admin-api/v1/app/patrol/issue` | JWT+casbin | 管理端新增 |
| DELETE | `/admin-api/v1/app/patrol/issue` | JWT+casbin | 批量删除 |

参数错误返回明确错误码:标题空 `31000`、优先级非法 `31001`、状态非法 `31002`。

## 参考实现(prototype)

开发前期为快速验证接口契约,先做了独立轻量版(`backend/` Go+pgx,`web/` React+Vite):`POST/GET /api/issues`、`PATCH /api/issues/{id}/status`、健康指示与错误横幅。**已停止演进**,保留作为设计参照;正式实现以上表 go-admin 集成为准。

## 技术栈

| 部分 | 技术 |
|---|---|
| 移动端 | Unity 6000.3 LTS + AR Foundation 6.x + ARCore(Vulkan) |
| 后端 | go-admin 框架(Gin + GORM)+ PostgreSQL 18 |
| 管理端 | go-admin web(React 18 + TypeScript + Vite + Ant Design Pro) |
| 联调 | adb 无线调试 + 局域网 HTTP + 防火墙放行 |

## 实现过程

见 [AI_LOG.md](AI_LOG.md)(20 条人工指令原文记录、任务拆解、8 项踩坑、验收记录)。
