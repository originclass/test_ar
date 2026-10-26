# 现场巡检问题标注(Patrol AR)

园区巡检人员用手机在现场放置 AR 标记并上报问题(墙面破损 / 设备异常 / 安全隐患),管理人员在网页端查看问题并更新处理状态。

**流程:** Unity AR 放置标记并上报 → Go 后端保存(PostgreSQL) → React 管理端查看 → React 修改处理状态

## 目录结构

```
patrol-ar/
├── backend/   # Go 后端(API + 持久化)
├── web/       # React + TypeScript 管理端
├── unity/     # Unity AR 移动端(Android + ARCore)
└── docs/      # 搭建文档 / 演示脚本 / AI 使用记录
```

## 快速启动(后端)

```powershell
cd backend
go build -o patrol-api.exe .      # 首次编译
.\patrol-api.exe                  # 默认监听 :8080,数据库连接见 store.go / PATROL_DB_DSN 环境变量
```

- 状态检查:`GET http://localhost:8080/health`
- 上报问题(Unity):`POST /api/issues`
- 问题列表:`GET /api/issues`
- 修改状态(React):`PATCH /api/issues/{id}/status`,body:`{"status":"open|in_progress|resolved"}`

统一响应格式:`{"code":0,"msg":"data"}`,业务失败 code 非 0 并配合 HTTP 400/404/500。

### React 管理端

```powershell
cd web
npm install        # 首次安装依赖
npm run dev        # 开发模式,http://localhost:5173
npm run build      # 类型检查 + 产物构建(dist/)
```

- 后端地址默认 `http://localhost:8080`,可用环境变量 `VITE_API_BASE` 覆盖
- 功能:问题列表(标题/优先级/状态/坐标/时间)、状态下拉即时修改、后端健康指示(10 秒探测)、状态/优先级筛选、后端不可用时显示错误横幅(不白屏,可点重试)

### Unity AR 端(搭建中)

见 `docs/unity-setup.md`。

## 技术栈

| 部分 | 技术 |
|---|---|
| 移动端 | Unity LTS + AR Foundation + ARCore(Android) |
| 管理端 | React 18 + TypeScript + Vite |
| 后端 | Go + PostgreSQL(本机 PostgreSQL 18,端口 3306) |

## 当前进度

- [x] 第 0 步:项目初始化
- [x] 第 1 步:Go 后端
- [x] 第 2 步:React 管理端
- [ ] 第 3 步:Unity AR 端
- [ ] 第 4 步:三端联调
- [ ] 第 5 步:文档与交付
