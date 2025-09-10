-- ================================================
-- TaskFlow 多层级项目管理系统 - 初始数据
-- 版本: 1.0.0
-- 创建时间: 2024-01-15
-- ================================================

SET NAMES utf8mb4;

-- ================================================
-- 插入系统角色
-- ================================================

INSERT INTO `roles` (`id`, `name`, `display_name`, `description`, `is_system`) VALUES
('role-super-admin', 'super_admin', '超级管理员', '系统超级管理员，拥有所有权限', TRUE),
('role-admin', 'admin', '系统管理员', '系统管理员，负责用户和权限管理', TRUE),
('role-project-owner', 'project_owner', '项目所有者', '项目所有者，可以创建和管理项目', TRUE),
('role-project-manager', 'project_manager', '项目经理', '项目经理，负责项目的日常管理', TRUE),
('role-team-leader', 'team_leader', '团队负责人', '团队负责人，可以管理团队成员和任务', TRUE),
('role-employee', 'employee', '普通员工', '普通员工，可以参与项目和任务', TRUE);

-- ================================================
-- 插入系统权限
-- ================================================

-- 用户管理权限
INSERT INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
('perm-user-create', 'user:create', 'user', 'create', '创建用户'),
('perm-user-read', 'user:read', 'user', 'read', '查看用户信息'),
('perm-user-update', 'user:update', 'user', 'update', '更新用户信息'),
('perm-user-delete', 'user:delete', 'user', 'delete', '删除用户'),
('perm-user-list', 'user:list', 'user', 'list', '查看用户列表'),

-- 角色和权限管理权限
('perm-role-create', 'role:create', 'role', 'create', '创建角色'),
('perm-role-read', 'role:read', 'role', 'read', '查看角色信息'),
('perm-role-update', 'role:update', 'role', 'update', '更新角色信息'),
('perm-role-delete', 'role:delete', 'role', 'delete', '删除角色'),
('perm-role-assign', 'role:assign', 'role', 'assign', '分配角色'),

-- 项目管理权限
('perm-project-create', 'project:create', 'project', 'create', '创建项目'),
('perm-project-read', 'project:read', 'project', 'read', '查看项目信息'),
('perm-project-update', 'project:update', 'project', 'update', '更新项目信息'),
('perm-project-delete', 'project:delete', 'project', 'delete', '删除项目'),
('perm-project-list', 'project:list', 'project', 'list', '查看项目列表'),
('perm-project-manage-members', 'project:manage_members', 'project', 'manage_members', '管理项目成员'),
('perm-project-assign-manager', 'project:assign_manager', 'project', 'assign_manager', '分配项目经理'),

-- 任务管理权限
('perm-task-create', 'task:create', 'task', 'create', '创建任务'),
('perm-task-read', 'task:read', 'task', 'read', '查看任务信息'),
('perm-task-update', 'task:update', 'task', 'update', '更新任务信息'),
('perm-task-delete', 'task:delete', 'task', 'delete', '删除任务'),
('perm-task-list', 'task:list', 'task', 'list', '查看任务列表'),
('perm-task-assign', 'task:assign', 'task', 'assign', '分配任务'),
('perm-task-approve', 'task:approve', 'task', 'approve', '审批任务'),
('perm-task-submit', 'task:submit', 'task', 'submit', '提交任务'),
('perm-task-execute', 'task:execute', 'task', 'execute', '执行任务'),
('perm-task-review', 'task:review', 'task', 'review', '审查任务工作成果'),
('perm-task-manage-participants', 'task:manage_participants', 'task', 'manage_participants', '管理任务参与者'),

-- 延期申请权限
('perm-extension-request', 'extension:request', 'extension', 'request', '申请延期'),
('perm-extension-approve', 'extension:approve', 'extension', 'approve', '审批延期申请'),

-- 文件管理权限
('perm-file-upload', 'file:upload', 'file', 'upload', '上传文件'),
('perm-file-download', 'file:download', 'file', 'download', '下载文件'),
('perm-file-delete', 'file:delete', 'file', 'delete', '删除文件'),

-- 统计分析权限
('perm-stats-view', 'stats:view', 'stats', 'view', '查看统计数据'),
('perm-stats-export', 'stats:export', 'stats', 'export', '导出统计数据');

-- ================================================
-- 角色权限关联
-- ================================================

-- 超级管理员：拥有所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 'role-super-admin', `id` FROM `permissions`;

-- 系统管理员：用户和角色管理权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
('role-admin', 'perm-user-create'),
('role-admin', 'perm-user-read'),
('role-admin', 'perm-user-update'),
('role-admin', 'perm-user-delete'),
('role-admin', 'perm-user-list'),
('role-admin', 'perm-role-create'),
('role-admin', 'perm-role-read'),
('role-admin', 'perm-role-update'),
('role-admin', 'perm-role-delete'),
('role-admin', 'perm-role-assign'),
('role-admin', 'perm-project-read'),
('role-admin', 'perm-project-list'),
('role-admin', 'perm-task-read'),
('role-admin', 'perm-task-list'),
('role-admin', 'perm-stats-view'),
('role-admin', 'perm-stats-export');

-- 项目所有者：项目和任务管理权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
('role-project-owner', 'perm-user-read'),
('role-project-owner', 'perm-user-list'),
('role-project-owner', 'perm-project-create'),
('role-project-owner', 'perm-project-read'),
('role-project-owner', 'perm-project-update'),
('role-project-owner', 'perm-project-delete'),
('role-project-owner', 'perm-project-list'),
('role-project-owner', 'perm-project-manage-members'),
('role-project-owner', 'perm-project-assign-manager'),
('role-project-owner', 'perm-task-create'),
('role-project-owner', 'perm-task-read'),
('role-project-owner', 'perm-task-update'),
('role-project-owner', 'perm-task-delete'),
('role-project-owner', 'perm-task-list'),
('role-project-owner', 'perm-task-assign'),
('role-project-owner', 'perm-task-approve'),
('role-project-owner', 'perm-task-review'),
('role-project-owner', 'perm-task-manage-participants'),
('role-project-owner', 'perm-extension-approve'),
('role-project-owner', 'perm-file-upload'),
('role-project-owner', 'perm-file-download'),
('role-project-owner', 'perm-file-delete'),
('role-project-owner', 'perm-stats-view'),
('role-project-owner', 'perm-stats-export');

-- 项目经理：项目管理和任务管理权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
('role-project-manager', 'perm-user-read'),
('role-project-manager', 'perm-user-list'),
('role-project-manager', 'perm-project-read'),
('role-project-manager', 'perm-project-update'),
('role-project-manager', 'perm-project-list'),
('role-project-manager', 'perm-project-manage-members'),
('role-project-manager', 'perm-task-create'),
('role-project-manager', 'perm-task-read'),
('role-project-manager', 'perm-task-update'),
('role-project-manager', 'perm-task-list'),
('role-project-manager', 'perm-task-assign'),
('role-project-manager', 'perm-task-approve'),
('role-project-manager', 'perm-task-review'),
('role-project-manager', 'perm-task-manage-participants'),
('role-project-manager', 'perm-extension-approve'),
('role-project-manager', 'perm-file-upload'),
('role-project-manager', 'perm-file-download'),
('role-project-manager', 'perm-stats-view');

-- 团队负责人：任务管理权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
('role-team-leader', 'perm-user-read'),
('role-team-leader', 'perm-user-list'),
('role-team-leader', 'perm-project-read'),
('role-team-leader', 'perm-project-list'),
('role-team-leader', 'perm-task-create'),
('role-team-leader', 'perm-task-read'),
('role-team-leader', 'perm-task-update'),
('role-team-leader', 'perm-task-list'),
('role-team-leader', 'perm-task-assign'),
('role-team-leader', 'perm-task-approve'),
('role-team-leader', 'perm-task-review'),
('role-team-leader', 'perm-task-manage-participants'),
('role-team-leader', 'perm-extension-approve'),
('role-team-leader', 'perm-file-upload'),
('role-team-leader', 'perm-file-download'),
('role-team-leader', 'perm-stats-view');

-- 普通员工：基础权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
('role-employee', 'perm-user-read'),
('role-employee', 'perm-project-read'),
('role-employee', 'perm-project-list'),
('role-employee', 'perm-task-read'),
('role-employee', 'perm-task-list'),
('role-employee', 'perm-task-submit'),
('role-employee', 'perm-task-execute'),
('role-employee', 'perm-extension-request'),
('role-employee', 'perm-file-upload'),
('role-employee', 'perm-file-download');

-- ================================================
-- 插入ABAC权限策略示例
-- ================================================

-- 用户只能查看自己的信息
INSERT INTO `permission_policies` (`id`, `name`, `description`, `resource_type`, `action`, `effect`, `conditions`, `priority`, `is_active`) VALUES
('policy-user-self-read', '用户查看自己信息', '用户只能查看自己的个人信息', 'user', 'read', 'allow', 
'{"subject.id": {"$eq": "resource.id"}}', 100, TRUE);

-- 项目成员可以查看项目信息
INSERT INTO `permission_policies` (`id`, `name`, `description`, `resource_type`, `action`, `effect`, `conditions`, `priority`, `is_active`) VALUES
('policy-project-member-read', '项目成员查看项目', '项目成员可以查看所属项目信息', 'project', 'read', 'allow',
'{"subject.id": {"$in": "resource.member_ids"}}', 90, TRUE);

-- 任务负责人可以管理任务
INSERT INTO `permission_policies` (`id`, `name`, `description`, `resource_type`, `action`, `effect`, `conditions`, `priority`, `is_active`) VALUES
('policy-task-responsible-manage', '任务负责人管理任务', '任务负责人可以管理自己负责的任务', 'task', 'update', 'allow',
'{"subject.id": {"$eq": "resource.responsible_id"}}', 80, TRUE);

-- 项目所有者可以管理项目下的所有任务
INSERT INTO `permission_policies` (`id`, `name`, `description`, `resource_type`, `action`, `effect`, `conditions`, `priority`, `is_active`) VALUES
('policy-project-owner-manage-tasks', '项目所有者管理任务', '项目所有者可以管理项目下的所有任务', 'task', '*', 'allow',
'{"subject.id": {"$eq": "resource.project.owner_id"}}', 70, TRUE);

-- ================================================
-- 创建默认超级管理员用户
-- ================================================

INSERT INTO `users` (`id`, `email`, `name`, `password_hash`, `status`, `department`, `position`) VALUES
('user-super-admin', 'admin@taskflow.com', '系统管理员', 
'$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- 密码: password
'active', 'IT部', '系统管理员');

-- 分配超级管理员角色
INSERT INTO `user_roles` (`user_id`, `role_id`, `assigned_by`) VALUES
('user-super-admin', 'role-super-admin', 'user-super-admin');

-- ================================================
-- 创建示例项目和任务数据（用于测试）
-- ================================================

-- 示例项目
INSERT INTO `projects` (`id`, `name`, `description`, `project_type`, `owner_id`, `status`, `start_date`) VALUES
('proj-demo-master', '企业数字化转型项目', '公司整体数字化转型的主项目', 'master', 'user-super-admin', 'active', CURDATE());

INSERT INTO `projects` (`id`, `name`, `description`, `project_type`, `parent_project_id`, `owner_id`, `status`, `start_date`) VALUES
('proj-demo-sub1', '系统架构设计', '设计新系统的整体架构', 'sub', 'proj-demo-master', 'user-super-admin', 'active', CURDATE()),
('proj-demo-sub2', '用户界面开发', '开发用户界面和交互功能', 'sub', 'proj-demo-master', 'user-super-admin', 'active', CURDATE());

-- 示例任务
INSERT INTO `tasks` (`id`, `title`, `description`, `task_type`, `priority`, `project_id`, `creator_id`, `responsible_id`, `status`, `due_date`) VALUES
('task-demo-1', '需求调研', '调研用户需求和业务流程', 'single_execution', 'high', 'proj-demo-sub1', 'user-super-admin', 'user-super-admin', 'in_progress', DATE_ADD(CURDATE(), INTERVAL 7 DAY)),
('task-demo-2', '技术选型', '选择合适的技术栈和框架', 'single_execution', 'normal', 'proj-demo-sub1', 'user-super-admin', 'user-super-admin', 'draft', DATE_ADD(CURDATE(), INTERVAL 14 DAY));

-- ================================================
-- 数据初始化完成
-- ================================================
