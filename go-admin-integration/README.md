# go-admin 集成说明

本目录是"巡检问题标注"模块在 go-admin 框架中的**全部新增文件**,按目标路径镜像存放(相对 go-admin 仓库根目录)。集成为纯增量:不修改框架任何既有文件,只动两个挂载点。

## 文件 → 目标路径

```
go-admin-integration/                          go-admin 仓库内路径
├── app/app/patrol/models/patrol_issue.go   →  app/app/patrol/models/patrol_issue.go
├── app/app/patrol/service/dto/...          →  app/app/patrol/service/dto/patrol_issue.go
├── app/app/patrol/service/patrol_issue.go  →  app/app/patrol/service/patrol_issue.go
├── app/app/patrol/apis/patrol_issue.go     →  app/app/patrol/apis/patrol_issue.go
├── app/app/patrol/router/patrol_issue.go   →  app/app/patrol/router/patrol_issue.go
├── app/app/patrol/router/router.go         →  app/app/patrol/router/router.go
├── config/base/lang/patrol.go              →  config/base/lang/patrol.go
├── db_schema.sql                           →  (psql 导入 go-admin 库)
├── patrol_menu_init.sql                    →  (psql 导入 go-admin 库)
└── init.go.patch                           →  app/app/init.go 增量修改(见下)
```

## 挂载点 1:路由注册(app/app/init.go)

`AllRouter()` 中追加两行(import + append),完整 diff 见 `init.go.patch`:

```go
import (
	patrolRouter "go-admin/app/app/patrol/router"
	userRouter "go-admin/app/app/user/router"
)

func AllRouter() []func() {
	var routers []func()
	routers = append(routers, userRouter.InitRouter)
	routers = append(routers, patrolRouter.InitRouter)   // ← 新增
	return routers
}
```

## 挂载点 2:菜单与权限(patrol_menu_init.sql)

- `admin_sys_menu`:135 巡检管理(目录)→ 136 巡检问题(页面,element=`/app/patrol/patrol-issue/index`)→ 137-141 按钮权限(query/add/edit/del/export)
- `admin_sys_role_menu`:角色 1(admin)与 2(test)全量授权
- `admin_sys_casbin_rule`:test 角色 6 条接口策略(admin 角色代码级豁免 casbin,无需策略)

## 编译与运行

```powershell
# go generate(重新生成 API 描述;Windows 下若报 unknown escape,把
# app/admin/sys/models/sys_api_gen_desc.go 里的反斜杠 '\' 批量替换为 '/')
go generate ./...

go build -o go-admin-api.exe main.go
.\go-admin-api.exe server --config config/settings.dev.yml   # 监听 :8888
```

## 设计要点

- **免认证上报**:`POST /admin-api/v1/app/patrol/issue/report` 注册在 `routerNoCheckRole`,供 Unity AR 客户端免登录上报;管理端同业务走认证组的 `POST /app/patrol/issue`。两条路径是刻意区分的 —— gin 不允许同 method+path 注册两次。
- **错误码**:31000-31004 定义于 `config/base/lang/patrol.go`(独立号段,不与框架冲突)。
- **position 存储**:数据库拍平三列 `pos_x/pos_y/pos_z`,API JSON 保持 `{x,y,z}` 嵌套(题目数据结构),model 层 `FillPosition/TakePosition` 转换。
- **服务返回约定**:与框架一致 `(data, respCode, err)`,respCode 用 `baseLang.SuccessCode` 等。
