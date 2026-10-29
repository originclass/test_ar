package lang

import "go-admin/core/lang"

// 巡检模块业务错误码(31000 段,避免与现有 30000 段冲突)
const (
	PatrolTitleEmptyCode       = 31000 // 标题为空
	PatrolPriorityInvalidCode  = 31001 // 优先级枚举值错误
	PatrolStatusInvalidCode    = 31002 // 状态枚举值错误
	PatrolStatusUpdateErrCode  = 31003 // 状态更新失败
	PatrolPositionParseErrCode = 31004 // 位置坐标解析失败
)

func init() {
	if lang.MsgInfo == nil {
		return
	}
	lang.MsgInfo[PatrolTitleEmptyCode] = "问题标题不能为空"
	lang.MsgInfo[PatrolPriorityInvalidCode] = "优先级必须是 low/medium/high 之一"
	lang.MsgInfo[PatrolStatusInvalidCode] = "状态必须是 open/in_progress/resolved 之一"
	lang.MsgInfo[PatrolStatusUpdateErrCode] = "问题状态更新失败：%s"
	lang.MsgInfo[PatrolPositionParseErrCode] = "位置坐标格式错误"
}
