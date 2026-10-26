import { useCallback, useEffect, useState } from "react";
import { checkHealth, fetchIssues, updateIssueStatus } from "./api";
import type { Issue, IssueStatus } from "./types";

type Health = "checking" | "up" | "down";

const STATUS_OPTIONS: { value: IssueStatus; label: string }[] = [
	{ value: "open", label: "待处理" },
	{ value: "in_progress", label: "处理中" },
	{ value: "resolved", label: "已解决" }
];

const PRIORITY_LABEL: Record<string, string> = { low: "低", medium: "中", high: "高" };

const formatTime = (iso: string) => {
	const d = new Date(iso);
	return isNaN(d.getTime()) ? "-" : d.toLocaleString("zh-CN", { hour12: false });
};

export default function App() {
	const [issues, setIssues] = useState<Issue[]>([]);
	const [loading, setLoading] = useState(true);
	const [errorMsg, setErrorMsg] = useState("");
	const [health, setHealth] = useState<Health>("checking");
	const [updatingId, setUpdatingId] = useState("");
	// 客户端筛选(加分项):按状态/优先级过滤
	const [statusFilter, setStatusFilter] = useState<"all" | IssueStatus>("all");
	const [priorityFilter, setPriorityFilter] = useState<"all" | string>("all");

	// 拉取问题列表;失败时设置错误信息(页面保留骨架,不白屏)
	const load = useCallback(async () => {
		setErrorMsg("");
		try {
			const data = await fetchIssues();
			setIssues(data);
		} catch (err) {
			setErrorMsg(err instanceof Error ? err.message : String(err));
		} finally {
			setLoading(false);
		}
	}, []);

	// 首次加载 + 每 10 秒探测后端健康状态
	useEffect(() => {
		void load();
		const timer = setInterval(async () => {
			try {
				await checkHealth();
				setHealth("up");
			} catch {
				setHealth("down");
			}
		}, 10_000);
		return () => clearInterval(timer);
	}, [load]);

	// 修改状态:成功后重新拉取列表(顺带验证持久化),失败则提示
	const handleStatusChange = async (issue: Issue, status: IssueStatus) => {
		if (status === issue.status) return;
		setUpdatingId(issue.id);
		try {
			await updateIssueStatus(issue.id, status);
			await load();
		} catch (err) {
			setErrorMsg(err instanceof Error ? `状态修改失败:${err.message}` : String(err));
		} finally {
			setUpdatingId("");
		}
	};

	const filtered = issues.filter(
		it =>
			(statusFilter === "all" || it.status === statusFilter) &&
			(priorityFilter === "all" || it.priority === priorityFilter)
	);

	const healthNode = {
		checking: <span className="health checking">● 检测中…</span>,
		up: <span className="health up">● 后端在线</span>,
		down: <span className="health down">● 后端离线</span>
	}[health];

	return (
		<div className="page">
			<header className="header">
				<h1>巡检问题管理</h1>
				{healthNode}
			</header>

			{errorMsg && (
				<div className="error-banner">
					<span>⚠ {errorMsg}</span>
					<button onClick={() => void load()}>重试</button>
				</div>
			)}

			<div className="toolbar">
				<label>
					状态筛选{" "}
					<select value={statusFilter} onChange={e => setStatusFilter(e.target.value as typeof statusFilter)}>
						<option value="all">全部</option>
						{STATUS_OPTIONS.map(o => (
							<option key={o.value} value={o.value}>
								{o.label}
							</option>
						))}
					</select>
				</label>
				<label>
					优先级{" "}
					<select value={priorityFilter} onChange={e => setPriorityFilter(e.target.value)}>
						<option value="all">全部</option>
						<option value="low">低</option>
						<option value="medium">中</option>
						<option value="high">高</option>
					</select>
				</label>
				<button className="ghost" onClick={() => void load()}>
					刷新
				</button>
			</div>

			<main className="table-wrap">
				{loading ? (
					<div className="empty">加载中…</div>
				) : filtered.length === 0 ? (
					<div className="empty">
						{errorMsg ? "暂无数据(后端异常)" : "还没有问题记录,用 Unity AR 端上报一条试试"}
					</div>
				) : (
					<table>
						<thead>
							<tr>
								<th>标题</th>
								<th>优先级</th>
								<th>状态</th>
								<th>位置 (x, y, z)</th>
								<th>上报时间</th>
								<th>更新时间</th>
							</tr>
						</thead>
						<tbody>
							{filtered.map(it => (
								<tr key={it.id}>
									<td>
										<div className="title">{it.title}</div>
										{it.description && <div className="desc">{it.description}</div>}
									</td>
									<td>
										<span className={`tag priority-${it.priority}`}>{PRIORITY_LABEL[it.priority] ?? it.priority}</span>
									</td>
									<td>
										<select
											className={`status-select status-${it.status}`}
											value={it.status}
											disabled={updatingId === it.id}
											onChange={e => void handleStatusChange(it, e.target.value as IssueStatus)}
										>
											{STATUS_OPTIONS.map(o => (
												<option key={o.value} value={o.value}>
													{o.label}
												</option>
											))}
										</select>
									</td>
									<td className="mono">
										({it.position.x.toFixed(2)}, {it.position.y.toFixed(2)}, {it.position.z.toFixed(2)})
									</td>
									<td>{formatTime(it.createdAt)}</td>
									<td>{formatTime(it.updatedAt)}</td>
								</tr>
							))}
						</tbody>
					</table>
				)}
			</main>

			<footer className="footer">共 {filtered.length} 条记录 · 数据来自 Go 后端 + PostgreSQL</footer>
		</div>
	);
}
