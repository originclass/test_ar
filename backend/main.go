// 现场巡检问题标注系统 - Go 后端
// 为 Unity AR 移动端和 React 管理端提供统一的 REST API,数据持久化到 PostgreSQL。
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

// App 组合数据层,HTTP 方法挂在上面
type App struct {
	store *Store
}

func main() {
	// 初始化数据库(15 秒超时,连不上直接失败退出,便于发现问题)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	store, err := NewStore(ctx)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	app := &App{store: store}

	// 路由:Go 1.22+ 的标准库支持 "METHOD /path/{param}" 模式,无需第三方路由
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.healthHandler)
	mux.HandleFunc("POST /api/issues", app.createIssueHandler)
	mux.HandleFunc("GET /api/issues", app.listIssuesHandler)
	mux.HandleFunc("PATCH /api/issues/{id}/status", app.updateStatusHandler)

	addr := os.Getenv("PATROL_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           corsMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("巡检后端已启动: http://localhost%s (GET /health 检查状态)", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("服务退出: %v", err)
	}
}
