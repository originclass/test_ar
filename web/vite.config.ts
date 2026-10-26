import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// 后端地址可通过环境变量 VITE_API_BASE 覆盖,默认本机 8080
export default defineConfig({
	plugins: [react()],
	server: {
		host: "0.0.0.0",
		port: 5173
	}
});
