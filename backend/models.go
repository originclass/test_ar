// 数据结构与常量定义。三端(Unity / React / Go)的 JSON 字段名保持一致。
package main

import "time"

// 问题优先级枚举
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

// 问题处理状态枚举
const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusResolved   = "resolved"
)

// Position AR 标记放置在真实世界中的位置坐标
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Issue 一条巡检问题。time.Time 在 JSON 序列化时默认输出 RFC3339 格式字符串。
type Issue struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	Position    Position  `json:"position"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateIssueReq Unity 上报问题的请求体(不含 id / status,由后端生成)
type CreateIssueReq struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Position    Position `json:"position"`
}

// UpdateStatusReq React 修改处理状态的请求体
type UpdateStatusReq struct {
	Status string `json:"status"`
}

// Response 统一响应格式:{code, msg, data},code 为 0 表示成功
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func isValidPriority(p string) bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}

func isValidStatus(s string) bool {
	return s == StatusOpen || s == StatusInProgress || s == StatusResolved
}
