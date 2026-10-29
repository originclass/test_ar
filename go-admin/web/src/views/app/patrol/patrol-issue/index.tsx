import {
	delPatrolIssueApi,
	getPatrolIssuePageApi,
	PatrolIssueModel,
	updatePatrolIssueStatusApi
} from "@/api/app/patrol/patrol-issue";
import HocAuth from "@/components/HocAuth";
import LoadingButton from "@/components/LoadingButton";
import { pagination } from "@/config/proTable";
import { ResultEnum } from "@/enums/httpEnum";
import { message, modal } from "@/hooks/useMessage";
import { formatDataForProTable } from "@/utils";
import { DeleteOutlined, PlusCircleOutlined } from "@ant-design/icons";
import type { ActionType, ProColumns } from "@ant-design/pro-components";
import { ProTable } from "@ant-design/pro-components";
import { Badge, Select, Space, Tag } from "antd";
import React, { useRef } from "react";
import FormModal, { FormModalRef } from "./components/FormModal";

const STATUS_ENUM: Record<string, { text: string; color: string }> = {
	open: { text: "待处理", color: "orange" },
	in_progress: { text: "处理中", color: "blue" },
	resolved: { text: "已解决", color: "green" }
};

const PRIORITY_ENUM: Record<string, { text: string; color: string }> = {
	low: { text: "低", color: "green" },
	medium: { text: "中", color: "gold" },
	high: { text: "高", color: "red" }
};

const PatrolIssue: React.FC = () => {
	const actionRef = React.useRef<ActionType>();
	const formModalRef = useRef<FormModalRef>(null);

	const handleStatusChange = async (id: number, status: string) => {
		try {
			const { code, msg } = await updatePatrolIssueStatusApi(id, status);
			if (code !== ResultEnum.SUCCESS) {
				message.error(msg);
				return;
			}
			message.success("状态已更新");
			actionRef.current?.reload(false);
		} catch {
			message.error("状态更新失败");
		}
	};

	const handleDelete = (id: number) => {
		modal.confirm({
			title: "提示",
			content: "是否确认删除编号为 " + id + " 的问题记录?",
			okText: "确认",
			cancelText: "取消",
			onOk: async () => {
				const { code, msg } = await delPatrolIssueApi([id]);
				if (code !== ResultEnum.SUCCESS) {
					message.error(msg);
					return;
				}
				message.success(msg);
				actionRef.current?.reload(false);
			}
		});
	};

	// 列定义:标题/优先级/状态为题目必展示字段,另附坐标与时间
	const columns: ProColumns<PatrolIssueModel>[] = [
		{
			title: "编号",
			dataIndex: "id",
			hideInSearch: true,
			width: 70,
			align: "center"
		},
		{
			title: "问题标题",
			dataIndex: "title",
			width: 220,
			render: (_, data) => (
				<div>
					<div style={{ fontWeight: 600 }}>{data.title}</div>
					{data.description && (
						<div style={{ color: "#999", fontSize: 12, whiteSpace: "pre-wrap" }}>{data.description}</div>
					)}
				</div>
			)
		},
		{
			title: "优先级",
			dataIndex: "priority",
			valueType: "select",
			valueEnum: PRIORITY_ENUM,
			width: 90,
			align: "center",
			render: (_, data) => {
				const it = PRIORITY_ENUM[data.priority ?? ""];
				return it ? <Tag color={it.color}>{it.text}</Tag> : data.priority;
			}
		},
		{
			title: "处理状态",
			dataIndex: "status",
			valueType: "select",
			valueEnum: STATUS_ENUM,
			width: 130,
			align: "center",
			render: (_, data) => (
				<Select
					size="small"
					style={{ width: 110 }}
					value={data.status}
					onChange={val => void handleStatusChange(data.id!, val)}
					options={Object.entries(STATUS_ENUM).map(([value, it]) => ({ value, label: it.text }))}
				/>
			)
		},
		{
			title: "标记位置 (x, y, z)",
			dataIndex: "position",
			hideInSearch: true,
			width: 160,
			align: "center",
			render: (_, data) => {
				const p = data.position;
				return p ? `(${(p.x ?? 0).toFixed(2)}, ${(p.y ?? 0).toFixed(2)}, ${(p.z ?? 0).toFixed(2)})` : "-";
			}
		},
		{
			title: "上报时间",
			dataIndex: "createdAt",
			hideInSearch: true,
			valueType: "dateTime",
			width: 170
		},
		{
			title: "上报时间",
			dataIndex: "createdAt",
			valueType: "dateTimeRange",
			hideInTable: true,
			search: { transform: value => ({ beginCreatedAt: value[0], endCreatedAt: value[1] }) }
		},
		{
			title: "操作",
			valueType: "option",
			align: "center",
			fixed: "right",
			width: 100,
			render: (_, data) => (
				<Space>
					<HocAuth permission={["app:patrol-issue:del"]}>
						<LoadingButton key="delete" type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => handleDelete(data.id!)}>
							删除
						</LoadingButton>
					</HocAuth>
				</Space>
			)
		}
	];

	return (
		<>
			<ProTable<PatrolIssueModel>
				className="ant-pro-table-scroll"
				columns={columns}
				actionRef={actionRef}
				bordered
				cardBordered
				defaultSize="small"
				scroll={{ x: 1100, y: "100%" }}
				request={async params => {
					const { data } = await getPatrolIssuePageApi(params);
					return formatDataForProTable<PatrolIssueModel>(data);
				}}
				options={{ reload: true, density: true, fullScreen: true }}
				rowKey="id"
				search={{ labelWidth: "auto" }}
				pagination={pagination}
				dateFormatter="string"
				headerTitle={
					<Space>
						巡检问题
						<Badge count="AR" style={{ backgroundColor: "#1677ff" }} title="该列表数据来自 Unity AR 客户端上报" />
					</Space>
				}
				toolBarRender={() => [
					<HocAuth permission={["app:patrol-issue:add"]}>
						<LoadingButton type="primary" key="add" icon={<PlusCircleOutlined />} onClick={() => formModalRef.current?.showAddFormModal()}>
							新增
						</LoadingButton>
					</HocAuth>
				]}
			/>
			<FormModal ref={formModalRef} onConfirm={() => actionRef.current?.reload(false)} />
		</>
	);
};

export default PatrolIssue;
