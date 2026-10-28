package dto

import (
	"time"

	"go-admin/app/app/patrol/models"
	"go-admin/core/dto"
)

// PatrolIssueQueryReq 分页查询请求(标题模糊搜索、优先级/状态精确、时间范围)
type PatrolIssueQueryReq struct {
	dto.Pagination `search:"-"`
	Title          string `form:"title" search:"type:contains;column:title;table:app_patrol_issue" comment:"问题标题"`
	Priority       string `form:"priority" search:"type:exact;column:priority;table:app_patrol_issue" comment:"优先级"`
	Status         string `form:"status" search:"type:exact;column:status;table:app_patrol_issue" comment:"处理状态"`
	BeginCreatedAt string `form:"beginCreatedAt" search:"type:gte;column:created_at;table:app_patrol_issue" comment:"上报时间起"`
	EndCreatedAt   string `form:"endCreatedAt" search:"type:lte;column:created_at;table:app_patrol_issue" comment:"上报时间止"`
	PatrolIssueOrder
}

type PatrolIssueOrder struct {
	IdOrder        int64      `form:"idOrder" search:"type:order;column:id;table:app_patrol_issue"`
	PriorityOrder  string     `form:"priorityOrder" search:"type:order;column:priority;table:app_patrol_issue"`
	StatusOrder    string     `form:"statusOrder" search:"type:order;column:status;table:app_patrol_issue"`
	CreatedAtOrder *time.Time `form:"createdAtOrder" search:"type:order;column:created_at;table:app_patrol_issue"`
	UpdatedAtOrder *time.Time `form:"updatedAtOrder" search:"type:order;column:updated_at;table:app_patrol_issue"`
}

func (m *PatrolIssueQueryReq) GetNeedSearch() interface{} {
	return *m
}

// PatrolIssueInsertReq 上报/新增问题请求(Unity 客户端与管理端共用)。
// id、status、createdAt、updatedAt 均由后端生成,不接受客户端传入。
type PatrolIssueInsertReq struct {
	Title       string          `json:"title" comment:"问题标题"`
	Description string          `json:"description" comment:"问题描述"`
	Priority    string          `json:"priority" comment:"优先级"`
	Position    models.Position `json:"position" comment:"标记位置"`
	CurrUserId  int64           `json:"-" comment:"当前登录用户(客户端上报时为0)"`
}

// PatrolIssueUpdateStatusReq 按 ID 修改处理状态请求
type PatrolIssueUpdateStatusReq struct {
	Id         int64  `json:"-" uri:"id" comment:"问题编号"`
	Status     string `json:"status" comment:"处理状态"`
	CurrUserId int64  `json:"-" comment:"当前登录用户"`
}

// PatrolIssueGetReq 获取问题详情请求
type PatrolIssueGetReq struct {
	Id int64 `uri:"id"`
}

// PatrolIssueDeleteReq 删除请求参数
type PatrolIssueDeleteReq struct {
	Ids []int64 `json:"ids"`
}
