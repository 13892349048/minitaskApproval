-- 004_core_business_tables.sql
-- 核心业务表结构

-- ABAC权限策略表
CREATE TABLE permission_policies (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    resource_type VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    effect ENUM('allow', 'deny') NOT NULL,
    conditions JSON NOT NULL,
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_resource_action (resource_type, action),
    INDEX idx_priority (priority),
    INDEX idx_active (is_active)
);

-- 权限属性定义表
CREATE TABLE permission_attributes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    attribute_type ENUM('user', 'resource', 'environment', 'action') NOT NULL,
    data_type ENUM('string', 'number', 'boolean', 'array', 'object') NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_name (name),
    INDEX idx_type (attribute_type)
);



-- 插入基础权限属性定义
INSERT INTO permission_attributes (id, name, attribute_type, data_type, description) VALUES
-- 用户属性
('attr_user_id', 'user.id', 'user', 'string', '用户ID'),
('attr_user_role', 'user.role', 'user', 'string', '用户角色'),
('attr_user_department', 'user.department_id', 'user', 'string', '用户部门ID'),
('attr_user_manager', 'user.manager_id', 'user', 'string', '用户直属领导ID'),

-- 资源属性
('attr_project_owner', 'project.owner_id', 'resource', 'string', '项目所有者ID'),
('attr_project_manager', 'project.manager_id', 'resource', 'string', '项目管理者ID'),
('attr_project_members', 'project.member_ids', 'resource', 'array', '项目成员ID列表'),
('attr_task_creator', 'task.creator_id', 'resource', 'string', '任务创建者ID'),
('attr_task_responsible', 'task.responsible_id', 'resource', 'string', '任务负责人ID'),
('attr_task_participants', 'task.participant_ids', 'resource', 'array', '任务参与人员ID列表'),
('attr_task_project', 'task.project_id', 'resource', 'string', '任务所属项目ID'),

-- 环境属性
('attr_time_now', 'env.current_time', 'environment', 'string', '当前时间'),
('attr_ip_address', 'env.ip_address', 'environment', 'string', '请求IP地址'),

-- 动作属性
('attr_action_type', 'action.type', 'action', 'string', '操作类型'),
('attr_action_target', 'action.target_id', 'action', 'string', '操作目标ID');


