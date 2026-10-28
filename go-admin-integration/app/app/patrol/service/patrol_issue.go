package service

import (
	"errors"
	"time"

	"go-admin/app/app/patrol/models"
	"go-admin/app/app/patrol/service/dto"
	baseLang "go-admin/config/base/lang"
	cDto "go-admin/core/dto"
	"go-admin/core/dto/service"
	"go-admin/core/lang"
	"go-admin/core/middleware"
	"gorm.io/gorm"
)

// PatrolIssue 巡检问题业务逻辑
type PatrolIssue struct {
	service.Service
}

// NewPatrolIssueService 实例化巡检问题服务
func NewPatrolIssueService(s *service.Service) *PatrolIssue {
	var srv = new(PatrolIssue)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 分页查询问题列表(按上报时间倒序)
func (e *PatrolIssue) GetPage(c *dto.PatrolIssueQueryReq, p *middleware.DataPermission) ([]models.PatrolIssue, int64, int, error) {
	var data models.PatrolIssue
	var list []models.PatrolIssue
	var count int64

	err := e.Orm.Order("created_at desc").Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			middleware.Permission(data.TableName(), p),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		return nil, 0, baseLang.DataQueryLogCode, lang.MsgLogErrf(e.Log, e.Lang, baseLang.DataQueryCode, baseLang.DataQueryLogCode, err)
	}
	for i := range list {
		list[i].FillPosition()
	}
	return list, count, baseLang.SuccessCode, nil
}

// Get 按 ID 获取问题详情
func (e *PatrolIssue) Get(id int64, p *middleware.DataPermission) (*models.PatrolIssue, int, error) {
	if id <= 0 {
		return nil, baseLang.ParamErrCode, lang.MsgErr(baseLang.ParamErrCode, e.Lang)
	}
	data := &models.PatrolIssue{}
	err := e.Orm.Scopes(
		middleware.Permission(data.TableName(), p),
	).First(data, id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, baseLang.DataQueryLogCode, lang.MsgLogErrf(e.Log, e.Lang, baseLang.DataQueryCode, baseLang.DataQueryLogCode, err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, baseLang.DataNotFoundCode, lang.MsgErr(baseLang.DataNotFoundCode, e.Lang)
	}
	data.FillPosition()
	return data, baseLang.SuccessCode, nil
}

// Insert 保存问题(Unity 客户端上报与管理端新增共用)。
// 参数错误时返回明确的错误码:标题为空 / 优先级枚举非法。
// 新问题唯一 ID 由数据库自增生成,默认状态为 open。
func (e *PatrolIssue) Insert(c *dto.PatrolIssueInsertReq) (int64, int, error) {
	if c.Title == "" {
		return 0, baseLang.PatrolTitleEmptyCode, lang.MsgErr(baseLang.PatrolTitleEmptyCode, e.Lang)
	}
	if !models.IsValidPriority(c.Priority) {
		return 0, baseLang.PatrolPriorityInvalidCode, lang.MsgErr(baseLang.PatrolPriorityInvalidCode, e.Lang)
	}

	now := time.Now()
	var data models.PatrolIssue
	data.Title = c.Title
	data.Description = c.Description
	data.Priority = c.Priority
	data.Status = models.PatrolStatusOpen // 新问题默认 open
	data.Position = c.Position
	data.TakePosition()
	data.CreateBy = c.CurrUserId
	data.UpdateBy = c.CurrUserId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err := e.Orm.Create(&data).Error
	if err != nil {
		return 0, baseLang.DataInsertLogCode, lang.MsgLogErrf(e.Log, e.Lang, baseLang.DataInsertCode, baseLang.DataInsertLogCode, err)
	}
	return data.Id, baseLang.SuccessCode, nil
}

// UpdateStatus 按问题 ID 修改处理状态;返回是否发生了变更
func (e *PatrolIssue) UpdateStatus(c *dto.PatrolIssueUpdateStatusReq, p *middleware.DataPermission) (bool, int, error) {
	if c.Id <= 0 {
		return false, baseLang.ParamErrCode, lang.MsgErr(baseLang.ParamErrCode, e.Lang)
	}
	if !models.IsValidStatus(c.Status) {
		return false, baseLang.PatrolStatusInvalidCode, lang.MsgErr(baseLang.PatrolStatusInvalidCode, e.Lang)
	}
	data, respCode, err := e.Get(c.Id, p)
	if err != nil {
		return false, respCode, err
	}
	if data.Status == c.Status { // 状态未变化,无需更新
		return false, baseLang.SuccessCode, nil
	}
	updates := map[string]interface{}{
		"status":     c.Status,
		"updated_at": time.Now(),
		"update_by":  c.CurrUserId,
	}
	err = e.Orm.Model(data).Where("id=?", data.Id).Updates(&updates).Error
	if err != nil {
		return false, baseLang.DataUpdateLogCode, lang.MsgLogErrf(e.Log, e.Lang, baseLang.DataUpdateCode, baseLang.DataUpdateLogCode, err)
	}
	return true, baseLang.SuccessCode, nil
}

// Delete 删除问题(支持批量)
func (e *PatrolIssue) Delete(ids []int64, p *middleware.DataPermission) (int, error) {
	if len(ids) <= 0 {
		return baseLang.ParamErrCode, lang.MsgErr(baseLang.ParamErrCode, e.Lang)
	}
	var data models.PatrolIssue
	err := e.Orm.Scopes(
		middleware.Permission(data.TableName(), p),
	).Delete(&data, ids).Error
	if err != nil {
		return baseLang.DataDeleteLogCode, lang.MsgLogErrf(e.Log, e.Lang, baseLang.DataDeleteCode, baseLang.DataDeleteLogCode, err)
	}
	return baseLang.SuccessCode, nil
}

// Count 数据总数(按查询条件)
func (e *PatrolIssue) Count(queryCondition *dto.PatrolIssueQueryReq) (int64, int, error) {
	var err error
	var count int64
	err = e.Orm.Model(&models.PatrolIssue{}).
		Scopes(
			cDto.MakeCondition(queryCondition.GetNeedSearch()),
		).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, baseLang.DataQueryLogCode, lang.MsgLogErrf(e.Log, e.Lang, baseLang.DataQueryCode, baseLang.DataQueryLogCode, err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, baseLang.DataNotFoundCode, lang.MsgErr(baseLang.DataNotFoundCode, e.Lang)
	}
	return count, baseLang.SuccessCode, nil
}
