# go-admin(完整工程,含巡检模块)

这是 [go-admin](https://github.com/go-admin-team/go-admin) 框架的完整源码 + 本项目新增的「巡检问题标注」业务模块。**本目录已包含前后端全部代码,拉取后无需再找框架仓库**。

## 本项目新增了什么

框架本体之外,新增的全部业务代码:

```
app/app/patrol/                     # 后端巡检模块(models/dto/service/apis/router 四层)
config/base/lang/patrol.go          # 业务错误码 31000-31004
web/src/api/app/patrol/             # 前端 API 客户端
web/src/views/app/patrol/           # 前端巡检问题页面(列表/状态修改/新增/删除)
```

根目录的 `patrol_menu_init.sql` 是菜单/角色/casbin 授权初始化(执行一次即可)。

## 快速运行

### 0. 前置

- Go 1.19+,Node.js 16+,PostgreSQL
- 数据库:创建 `go-admin` 库,先导入框架自带 `app_pgsql.sql`(框架全部基础表),再执行下面的初始化 SQL

### 1. 数据库初始化

```powershell
psql -U postgres -d go-admin -f db_schema.sql          # 巡检问题表(见下)
psql -U postgres -d go-admin -f patrol_menu_init.sql   # 菜单 + 角色授权 + casbin
```

`db_schema.sql`(巡检问题表):

```sql
CREATE TABLE IF NOT EXISTS app_patrol_issue (
	id          BIGSERIAL PRIMARY KEY,
	title       VARCHAR(255) NOT NULL,
	description VARCHAR(1000) DEFAULT '',
	priority    VARCHAR(10) NOT NULL,
	status      VARCHAR(20) NOT NULL,
	pos_x       DOUBLE PRECISION DEFAULT 0,
	pos_y       DOUBLE PRECISION DEFAULT 0,
	pos_z       DOUBLE PRECISION DEFAULT 0,
	create_by   INT DEFAULT 0,
	update_by   INT DEFAULT 0,
	created_at  TIMESTAMPTZ DEFAULT now(),
	updated_at  TIMESTAMPTZ DEFAULT now()
);
```

> 本地 PostgreSQL 端口用了 3306(配置在 `config/db.sh` 或 settings.yml 里,按你的环境改)。

### 2. 启动后端(:8888)

```powershell
go build -o go-admin-api.exe main.go
.\go-admin-api.exe server --config config/settings.dev.yml
```

### 3. 启动前端(:1688)

```powershell
cd web
npm install
npm run dev
```

登录 http://localhost:1688(admin/123456)→ 菜单「App应用 → 巡检管理 → 巡检问题」。

### 4. Unity AR 客户端

用 Unity 6 打开仓库根的 `ar_poj/AR_prj`,菜单 `Patrol → Setup Guide` 查看完整搭建/打包说明;客户端上报地址在 `Assets/Scripts/PatrolApi.cs` 的 `BaseUrl`(改成你电脑的局域网 IP:8888)。
