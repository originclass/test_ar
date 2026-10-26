// 与 Go 后端约定的问题数据结构(三端字段名一致)
export type Priority = "low" | "medium" | "high";
export type IssueStatus = "open" | "in_progress" | "resolved";

export interface Position {
	x: number;
	y: number;
	z: number;
}

export interface Issue {
	id: string;
	title: string;
	description: string;
	priority: Priority;
	status: IssueStatus;
	position: Position;
	createdAt: string;
	updatedAt: string;
}

/** 后端统一响应:{code, msg, data},code=0 表示成功 */
export interface ApiResp<T> {
	code: number;
	msg: string;
	data: T;
}
