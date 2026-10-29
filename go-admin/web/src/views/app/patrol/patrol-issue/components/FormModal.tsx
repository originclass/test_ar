import { addPatrolIssueApi, PatrolIssueModel } from "@/api/app/patrol/patrol-issue";
import LoadingButton from "@/components/LoadingButton";
import { ResultEnum } from "@/enums/httpEnum";
import { message } from "@/hooks/useMessage";
import { Form, Input, Modal, Select } from "antd";
import { forwardRef, useImperativeHandle, useState } from "react";

export interface FormModalRef {
	showAddFormModal: () => void;
}

interface ModalProps {
	onConfirm: () => void;
}

const PRIORITY_OPTIONS = [
	{ value: "low", label: "低" },
	{ value: "medium", label: "中" },
	{ value: "high", label: "高" }
];

/** 新增巡检问题弹窗(管理端模拟客户端上报用) */
const FormModal = forwardRef<FormModalRef, ModalProps>(({ onConfirm }, ref) => {
	const [form] = Form.useForm();
	const [isModalOpen, setIsModalOpen] = useState(false);
	const [model, setModel] = useState<PatrolIssueModel>({});

	useImperativeHandle(ref, () => ({
		showAddFormModal() {
			reset();
			setIsModalOpen(true);
		}
	}));

	const reset = () => {
		setModel({ id: 0 });
		setTimeout(() => form.resetFields(), 100);
	};

	const handleConfirm = (done: () => void) => {
		form
			.validateFields()
			.then(async values => {
				try {
					const payload = { ...values, position: values.position ?? { x: 0, y: 0, z: 0 } };
					const { msg, code } = await addPatrolIssueApi(payload);
					if (code !== ResultEnum.SUCCESS) {
						message.error(msg);
						return;
					}
					message.success(msg);
					reset();
					setIsModalOpen(false);
					onConfirm();
				} finally {
					done();
				}
			})
			.catch(() => {
				message.error("表单校验失败");
				done();
			});
	};

	return (
		<Modal
			title="新增巡检问题"
			getContainer={false}
			width={520}
			open={isModalOpen}
			maskClosable={false}
			keyboard={false}
			onCancel={() => {
				reset();
				setIsModalOpen(false);
			}}
			destroyOnClose
			footer={[
				<LoadingButton
					key="cancel"
					onClick={done => {
						reset();
						setIsModalOpen(false);
						done();
					}}
				>
					取消
				</LoadingButton>,
				<LoadingButton key="confirm" type="primary" onClick={done => handleConfirm(done)}>
					确定
				</LoadingButton>
			]}
		>
			<Form form={form} layout="vertical" initialValues={model}>
				<Form.Item name="title" label="问题标题" rules={[{ required: true, message: "问题标题不能为空" }]}>
					<Input placeholder="例如:入口墙面破损" maxLength={255} />
				</Form.Item>
				<Form.Item name="description" label="问题描述(可选)">
					<Input.TextArea placeholder="补充说明(可选)" rows={3} maxLength={1000} />
				</Form.Item>
				<Form.Item name="priority" label="优先级" rules={[{ required: true, message: "请选择优先级" }]}>
					<Select placeholder="请选择" options={PRIORITY_OPTIONS} />
				</Form.Item>
				<Form.Item
					name="position"
					label="标记位置 (x, y, z)"
					rules={[
						{
							validator: (_, value) => {
								if (!value) return Promise.resolve();
								const ok = ["x", "y", "z"].every(k => typeof value[k] === "number" && !isNaN(value[k]));
								return ok ? Promise.resolve() : Promise.reject(new Error("坐标必须是数字"));
							}
						}
					]}
				>
					<Input.Group compact>
						<Form.Item name={["position", "x"]} noStyle>
							<Input style={{ width: "32%" }} placeholder="x" type="number" step="0.01" />
						</Form.Item>
						<Form.Item name={["position", "y"]} noStyle>
							<Input style={{ width: "32%", marginLeft: "2%" }} placeholder="y" type="number" step="0.01" />
						</Form.Item>
						<Form.Item name={["position", "z"]} noStyle>
							<Input style={{ width: "32%", marginLeft: "2%" }} placeholder="z" type="number" step="0.01" />
						</Form.Item>
					</Input.Group>
				</Form.Item>
			</Form>
		</Modal>
	);
});

export default FormModal;
