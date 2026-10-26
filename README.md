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

## 快速启动(待补充)

> 各部分详细启动方式将在实现完成后补充于此。

## 技术栈

| 部分 | 技术 |
|---|---|
| 移动端 | Unity LTS + AR Foundation + ARCore(Android) |
| 管理端 | React 18 + TypeScript + Vite |
| 后端 | Go + PostgreSQL(本机 PostgreSQL 18,端口 3306) |

## 当前进度

- [ ] 第 0 步:项目初始化
- [ ] 第 1 步:Go 后端
- [ ] 第 2 步:React 管理端
- [ ] 第 3 步:Unity AR 端
- [ ] 第 4 步:三端联调
- [ ] 第 5 步:文档与交付
