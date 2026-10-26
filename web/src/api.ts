import type { ApiResp, Issue, IssueStatus } from "./types";

// 后端地址:开发/部署时可用 VITE_API_BASE 环境变量覆盖(Unity 真机演示时 PC 局域网 IP 也可用这里改)
const API_BASE: string = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";

/**
 * 统一请求封装:超时控制 + 错误归一化。
 * 任何失败(网络不通/超时/业务错误)都抛出带可读信息的 Error,由页面统一展示,
 * 保证后端不可用时页面给出错误提示而不是白屏。
 */
async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const ctrl = new AbortController();
	const timer = setTimeout(() => ctrl.abort(), 8000); // 8 秒超时
	try {
		const res = await fetch(`${API_BASE}${path}`, {
			...init,
			signal: ctrl.signal,
			headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) }
		});
		let body: ApiResp<T>;
		try {
			body = (await res.json()) as ApiResp<T>;
		} catch {
			throw new Error(`后端响应异常(HTTP ${res.status})`);
		}
		if (!res.ok || body.code !== 0) {
			throw new Error(body?.msg || `请求失败(HTTP ${res.status})`);
		}
		return body.data;
	} catch (err) {
		if (err instanceof DOMException && err.name === "AbortError") {
			throw new Error("请求超时,请检查后端是否在运行");
		}
		if (err instanceof TypeError) {
			// fetch 网络层失败(连不上/断网/CORS)
			throw new Error(`无法连接后端服务(${API_BASE}),请确认后端已启动`);
		}
		throw err;
	} finally {
		clearTimeout(timer);
	}
}

export const fetchIssues = () => request<Issue[]>("/api/issues");

export const updateIssueStatus = (id: string, status: IssueStatus) =>
	request<{ id: string; status: IssueStatus }>(`/api/issues/${id}/status`, {
		method: "PATCH",
		body: JSON.stringify({ status })
	});

export const checkHealth = () => request<{ status: string }>("/health");
