// HTTP 处理层:参数校验、统一响应、CORS。
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func writeOK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, Response{Code: 0, Msg: "ok", Data: data})
}

func writeError(w http.ResponseWriter, httpStatus, code int, msg string) {
	writeJSON(w, httpStatus, Response{Code: code, Msg: msg})
}

// healthHandler 服务运行状态检查
func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]string{"status": "up", "time": time.Now().Format(time.RFC3339)})
}

// createIssueHandler Unity 上报问题:校验 → 生成 ID/默认状态 → 落库
func (app *App) createIssueHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateIssueReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, 40001, "请求体不是合法的 JSON: "+err.Error())
		return
	}
	if req.Title == "" { // 标题必填
		writeError(w, http.StatusBadRequest, 40002, "标题不能为空")
		return
	}
	if !isValidPriority(req.Priority) { // 优先级枚举校验
		writeError(w, http.StatusBadRequest, 40003, "优先级必须是 low/medium/high 之一")
		return
	}

	now := time.Now()
	issue := Issue{
		ID:          uuid.NewString(), // 后端生成唯一 ID
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      StatusOpen, // 新问题默认 open
		Position:    req.Position,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := app.store.Create(r.Context(), &issue); err != nil {
		writeError(w, http.StatusInternalServerError, 50001, "保存失败: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, Response{Code: 0, Msg: "ok", Data: issue})
}

// listIssuesHandler 返回问题列表(React 与 Unity 共用同一份数据)
func (app *App) listIssuesHandler(w http.ResponseWriter, r *http.Request) {
	issues, err := app.store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, 50002, "查询失败: "+err.Error())
		return
	}
	writeOK(w, issues)
}

// updateStatusHandler 按 ID 修改处理状态
func (app *App) updateStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, 40004, "请求体不是合法的 JSON: "+err.Error())
		return
	}
	if !isValidStatus(req.Status) { // 状态枚举校验
		writeError(w, http.StatusBadRequest, 40005, "状态必须是 open/in_progress/resolved 之一")
		return
	}
	updated, err := app.store.UpdateStatus(r.Context(), id, req.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, 50003, "更新失败: "+err.Error())
		return
	}
	if !updated {
		writeError(w, http.StatusNotFound, 40401, "问题不存在: "+id)
		return
	}
	writeOK(w, map[string]string{"id": id, "status": req.Status})
}

// corsMiddleware 浏览器管理端跨域访问后端,必须放行 CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions { // 预检请求直接放行
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
