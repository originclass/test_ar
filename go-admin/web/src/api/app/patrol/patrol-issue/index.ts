import request from "@/utils/request";
import { ReqPage, ResPage } from "@/utils/request/interface";

/** 巡检问题(与后端 models.PatrolIssue 的 JSON 字段一致,三端共用同一结构) */
export interface PatrolPosition {
	x?: number;
	y?: number;
	z?: number;
}

export interface PatrolIssueModel {
	id?: number;
	title?: string;
	description?: string;
	/** low | medium | high */
	priority?: string;
	/** open | in_progress | resolved */
	status?: string;
	position?: PatrolPosition;
	createdAt?: string;
	updatedAt?: string;
}

/** 分页查询问题列表(React 管理端,需登录) */
export const getPatrolIssuePageApi = (params: ReqPage & Record<string, any>) => {
	return request.get<ResPage<PatrolIssueModel>>(`/admin-api/v1/app/patrol/issue`, {
		...params,
		pageIndex: params?.current
	});
};

/** 获取问题详情 */
export const getPatrolIssueApi = (id: number) => {
	return request.get<PatrolIssueModel>(`/admin-api/v1/app/patrol/issue/` + id);
};

/** 新增问题(管理端) */
export const addPatrolIssueApi = (data: object) => {
	return request.post<object>(`/admin-api/v1/app/patrol/issue`, data);
};

/** 按问题 ID 修改处理状态(open/in_progress/resolved) */
export const updatePatrolIssueStatusApi = (id: number, status: string) => {
	return request.put<object>(`/admin-api/v1/app/patrol/issue/${id}/status`, { status });
};

/** 删除问题(支持批量) */
export const delPatrolIssueApi = (ids: number[]) => {
	return request.delete<object>(`/admin-api/v1/app/patrol/issue`, { ids });
};
