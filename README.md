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

本仓库包含**完整的 Unity AR 客户端工程**与**参考实现**(early prototype),业务后端与管理端集成在 [go-admin](https://github.com/go-admin-team/go-admin) 框架(集成代码副本见 `go-admin-integration/`):

```
patrol-ar/
├── ar_poj/AR_prj/            # ★ Unity 6 AR 客户端工程(真机验证通过)
│   ├── Assets/Scripts/       #    全部 C# 源码(3 个脚本,见下文)
│   ├── Assets/Scenes/        #    SampleScene(AR Session + 平面检测 + UI)
│   └── Assets/Settings/      #    URP 渲染配置(含 AR 渲染特性,见 AI_LOG 踩坑记录)
├── go-admin-integration/     # ★ 集成到 go-admin 的全部新增代码(含目标路径说明)
├── unity/                    # Unity 保姆级搭建文档(SETUP_GUIDE.md)与脚本副本
├── backend/                  # 参考实现:独立 Go 后端(prototype)
├── web/                      # 参考实现:独立 React 管理端(prototype)
├── docs/                     # 文档(搭建指南 / AI 协作日志)
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

## go-admin 集成(`go-admin-integration/`)

业务代码按 go-admin 四层结构(模块:`patrol`)新增,不改动框架既有文件(仅注册路由与菜单 SQL 两处挂载点)。文件清单与集成目标路径:

| 本仓库文件 | 集成目标路径 | 说明 |
|---|---|---|
| `app/app/patrol/models/patrol_issue.go` | `app/app/patrol/models/` | gorm model:嵌套 `position` 对象 ↔ `pos_x/y/z` 三列转换;枚举常量与校验 |
| `app/app/patrol/service/dto/patrol_issue.go` | `app/app/patrol/service/dto/` | 请求 DTO:分页查询(search 标签)/ 上报 / 状态更新 |
| `app/app/patrol/service/patrol_issue.go` | `app/app/patrol/service/` | 业务:Insert(标题/优先级校验、状态强制 open)、UpdateStatus(枚举校验)、GetPage/Get/Delete |
| `app/app/patrol/apis/patrol_issue.go` | `app/app/patrol/apis/` | gin 处理器(框架标准链式调用) |
| `app/app/patrol/router/*.go` | `app/app/patrol/router/` | 免认证组(`POST /app/patrol/issue/report`,供 Unity 上报)+ JWT+casbin 认证组(管理端) |
| `config/base/lang/patrol.go` | `config/base/lang/` | 业务错误码 31000-31004 注册 |
| `init.go.patch` | `app/app/init.go` | 模块路由注册(一处 append) |
| `patrol_menu_init.sql` | (psql 执行) | 菜单/按钮权限/admin+test 角色授权/casbin 规则 |
| `db_schema.sql` | (psql 执行) | `app_patrol_issue` 表结构 |

### 快速启动(go-admin 后端 + React 管理端)

```powershell
# 1. 数据库(PostgreSQL,库名 go-admin)
psql -U postgres -d go-admin -f go-admin-integration/db_schema.sql
psql -U postgres -d go-admin -f go-admin-integration/patrol_menu_init.sql

# 2. 后端(go-admin 框架,把 go-admin-integration/ 下文件按上表拷入框架源码)
go build -o go-admin-api.exe main.go
.\go-admin-api.exe server --config config/settings.dev.yml     # :8888

# 3. 管理端
cd web && npm install && npm run dev
# 登录 admin → 菜单「App应用 → 巡检管理 → 巡检问题」
```

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

见 [docs/AI_LOG.md](docs/AI_LOG.md)(任务拆解、架构决策、8 项踩坑记录、验收记录)。
