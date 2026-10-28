package models

import (
	"time"
)

// Position 巡检问题在 AR 场景中的三维坐标。
// 注意:数据库中拍平为 pos_x/pos_y/pos_z 三列,API JSON 中保持嵌套对象(与题目约定的数据结构一致),
// 通过 FillPosition / TakePosition 在两种形态之间转换。
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// PatrolIssue 巡检问题。三端(Unity / React / Go)的 JSON 字段名保持一致。
type PatrolIssue struct {
	Id          int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Title       string     `json:"title" gorm:"column:title;type:varchar(255);comment:问题标题"`
	Description string     `json:"description" gorm:"column:description;type:varchar(1000);comment:问题描述"`
	Priority    string     `json:"priority" gorm:"column:priority;type:varchar(10);comment:优先级(low/medium/high)"`
	Status      string     `json:"status" gorm:"column:status;type:varchar(20);comment:处理状态(open/in_progress/resolved)"`
	PosX        float64    `json:"-" gorm:"column:pos_x;comment:标记位置X"`
	PosY        float64    `json:"-" gorm:"column:pos_y;comment:标记位置Y"`
	PosZ        float64    `json:"-" gorm:"column:pos_z;comment:标记位置Z"`
	Position    Position   `json:"position" gorm:"-"`
	CreateBy    int64      `json:"-" gorm:"column:create_by;type:int;comment:创建者(0表示AR客户端上报)"`
	UpdateBy    int64      `json:"-" gorm:"column:update_by;type:int;comment:更新者"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"column:created_at;type:datetime;comment:上报时间"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"column:updated_at;type:datetime;comment:更新时间"`
}

func (PatrolIssue) TableName() string {
	return "app_patrol_issue"
}

// FillPosition 把拍平的坐标列组装为嵌套的 Position 对象(查询后调用)
func (e *PatrolIssue) FillPosition() {
	e.Position = Position{X: e.PosX, Y: e.PosY, Z: e.PosZ}
}

// TakePosition 把嵌套的 Position 对象拆到拍平的坐标列(入库前调用)
func (e *PatrolIssue) TakePosition() {
	e.PosX = e.Position.X
	e.PosY = e.Position.Y
	e.PosZ = e.Position.Z
}

// 优先级枚举
const (
	PatrolPriorityLow    = "low"
	PatrolPriorityMedium = "medium"
	PatrolPriorityHigh   = "high"
)

// 处理状态枚举
const (
	PatrolStatusOpen       = "open"
	PatrolStatusInProgress = "in_progress"
	PatrolStatusResolved   = "resolved"
)

// IsValidPriority 优先级枚举校验
func IsValidPriority(p string) bool {
	return p == PatrolPriorityLow || p == PatrolPriorityMedium || p == PatrolPriorityHigh
}

// IsValidStatus 状态枚举校验
func IsValidStatus(s string) bool {
	return s == PatrolStatusOpen || s == PatrolStatusInProgress || s == PatrolStatusResolved
}
