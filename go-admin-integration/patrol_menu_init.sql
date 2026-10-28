-- 巡检模块菜单与权限初始化
-- 挂载结构:App应用(58) → 巡检管理(135) → 巡检问题(136) + 5个按钮权限(137-141)
-- 角色授权:admin(1) 与 test(2) 均可见可用

INSERT INTO admin_sys_menu (id, title, icon, path, element, redirect, permission, sort, parent_id, parent_ids, menu_type, is_keep_alive, is_affix, is_hidden, is_frame, create_by, update_by, created_at, updated_at) VALUES
(135, '巡检管理', 'RadarChartOutlined', '/app/patrol', '', '', '', 1, 58, '0,58,', 1, 1, 2, 2, 1, 1, 1, now(), now()),
(136, '巡检问题', 'FileDoneOutlined', '/app/patrol/patrol-issue', '/app/patrol/patrol-issue/index', '', '', 1, 135, '0,58,135,', 2, 1, 2, 2, 1, 1, 1, now(), now()),
(137, '巡检问题查询', '', '', '', '', 'app:patrol-issue:query', 1, 136, '0,58,135,136,', 3, 1, 2, 2, 1, 1, 1, now(), now()),
(138, '巡检问题新增', '', '', '', '', 'app:patrol-issue:add', 2, 136, '0,58,135,136,', 3, 1, 2, 2, 1, 1, 1, now(), now()),
(139, '巡检问题编辑', '', '', '', '', 'app:patrol-issue:edit', 3, 136, '0,58,135,136,', 3, 1, 2, 2, 1, 1, 1, now(), now()),
(140, '巡检问题删除', '', '', '', '', 'app:patrol-issue:del', 4, 136, '0,58,135,136,', 3, 1, 2, 2, 1, 1, 1, now(), now()),
(141, '巡检问题导出', '', '', '', '', 'app:patrol-issue:export', 5, 136, '0,58,135,136,', 3, 1, 2, 2, 1, 1, 1, now(), now());

INSERT INTO admin_sys_role_menu (role_id, menu_id) VALUES
(1, 135), (1, 136), (1, 137), (1, 138), (1, 139), (1, 140), (1, 141),
(2, 135), (2, 136), (2, 137), (2, 138), (2, 139), (2, 140), (2, 141);

-- casbin 策略:test 角色访问巡检接口(admin 角色在代码中豁免 casbin,无需策略)
INSERT INTO admin_sys_casbin_rule (p_type, v0, v1, v2) VALUES
('p', 'test', '/admin-api/v1/app/patrol/issue', 'GET'),
('p', 'test', '/admin-api/v1/app/patrol/issue', 'POST'),
('p', 'test', '/admin-api/v1/app/patrol/issue', 'DELETE'),
('p', 'test', '/admin-api/v1/app/patrol/issue/:id', 'GET'),
('p', 'test', '/admin-api/v1/app/patrol/issue/:id/status', 'PUT'),
('p', 'test', '/admin-api/v1/app/patrol/issue/report', 'POST');
