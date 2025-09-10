-- ================================================
-- TaskFlow 多层级项目管理系统 - 初始数据库结构
-- 版本: 1.0.0
-- 创建时间: 2024-01-15
-- ================================================

-- 设置字符集和排序规则
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ================================================
-- 用户和权限相关表
-- ================================================

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '用户ID',
    `email` VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱地址',
    `name` VARCHAR(100) NOT NULL COMMENT '用户姓名',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
    `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
    `status` ENUM('active', 'inactive', 'suspended') DEFAULT 'active' COMMENT '用户状态',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号码',
    `department` VARCHAR(100) DEFAULT NULL COMMENT '部门',
    `position` VARCHAR(100) DEFAULT NULL COMMENT '职位',
    `join_date` DATE DEFAULT NULL COMMENT '入职日期',
    `department_id` VARCHAR(36) DEFAULT NULL COMMENT '部门ID',
    `manager_id` VARCHAR(36) DEFAULT NULL COMMENT '上级管理者ID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX `idx_email` (`email`),
    INDEX `idx_status` (`status`),
    INDEX `idx_department` (`department_id`),
    INDEX `idx_manager` (`manager_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS `roles` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '角色ID',
    `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
    `display_name` VARCHAR(100) NOT NULL COMMENT '显示名称',
    `description` TEXT COMMENT '角色描述',
    `is_system` BOOLEAN DEFAULT FALSE COMMENT '是否系统角色',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX `idx_name` (`name`),
    INDEX `idx_system` (`is_system`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `assigned_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
    `assigned_by` VARCHAR(36) DEFAULT NULL COMMENT '分配人ID',
    
    PRIMARY KEY (`user_id`, `role_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`assigned_by`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    
    INDEX `idx_user` (`user_id`),
    INDEX `idx_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 权限表
CREATE TABLE IF NOT EXISTS `permissions` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '权限ID',
    `name` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限名称',
    `resource` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `action` VARCHAR(50) NOT NULL COMMENT '操作类型',
    `description` TEXT COMMENT '权限描述',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX `idx_resource_action` (`resource`, `action`),
    INDEX `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS `role_permissions` (
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `permission_id` VARCHAR(36) NOT NULL COMMENT '权限ID',
    
    PRIMARY KEY (`role_id`, `permission_id`),
    FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ABAC权限策略表
CREATE TABLE IF NOT EXISTS `permission_policies` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '策略ID',
    `name` VARCHAR(200) NOT NULL COMMENT '策略名称',
    `description` TEXT COMMENT '策略描述',
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `action` VARCHAR(50) NOT NULL COMMENT '操作类型',
    `effect` ENUM('allow', 'deny') NOT NULL COMMENT '效果',
    `conditions` JSON NOT NULL COMMENT '条件规则',
    `priority` INT DEFAULT 0 COMMENT '优先级',
    `is_active` BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX `idx_resource_action` (`resource_type`, `action`),
    INDEX `idx_priority` (`priority`),
    INDEX `idx_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ABAC权限策略表';

-- ================================================
-- 项目管理相关表
-- ================================================

-- 项目表
CREATE TABLE IF NOT EXISTS `projects` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '项目ID',
    `name` VARCHAR(200) NOT NULL COMMENT '项目名称',
    `description` TEXT COMMENT '项目描述',
    `project_type` ENUM('master', 'sub', 'temporary') NOT NULL COMMENT '项目类型',
    `parent_project_id` VARCHAR(36) DEFAULT NULL COMMENT '父项目ID',
    `owner_id` VARCHAR(36) NOT NULL COMMENT '项目所有者ID',
    `manager_id` VARCHAR(36) DEFAULT NULL COMMENT '项目管理者ID',
    `status` ENUM('draft', 'active', 'paused', 'completed', 'cancelled') DEFAULT 'draft' COMMENT '项目状态',
    `start_date` DATE DEFAULT NULL COMMENT '开始日期',
    `end_date` DATE DEFAULT NULL COMMENT '结束日期',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (`parent_project_id`) REFERENCES `projects`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`owner_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    FOREIGN KEY (`manager_id`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    
    INDEX `idx_parent_project` (`parent_project_id`),
    INDEX `idx_owner` (`owner_id`),
    INDEX `idx_manager` (`manager_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_type_status` (`project_type`, `status`),
    INDEX `idx_dates` (`start_date`, `end_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

-- 项目成员表
CREATE TABLE IF NOT EXISTS `project_members` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '记录ID',
    `project_id` VARCHAR(36) NOT NULL COMMENT '项目ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `role` ENUM('manager', 'member') NOT NULL COMMENT '角色',
    `joined_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    `added_by` VARCHAR(36) DEFAULT NULL COMMENT '添加人ID',
    
    UNIQUE KEY `uk_project_user` (`project_id`, `user_id`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`added_by`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    
    INDEX `idx_project` (`project_id`),
    INDEX `idx_user` (`user_id`),
    INDEX `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目成员表';

-- ================================================
-- 任务管理相关表
-- ================================================

-- 任务表
CREATE TABLE IF NOT EXISTS `tasks` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '任务ID',
    `title` VARCHAR(300) NOT NULL COMMENT '任务标题',
    `description` TEXT COMMENT '任务描述',
    `task_type` ENUM('single_execution', 'recurring') NOT NULL COMMENT '任务类型',
    `priority` ENUM('low', 'normal', 'high', 'urgent') DEFAULT 'normal' COMMENT '优先级',
    `project_id` VARCHAR(36) NOT NULL COMMENT '项目ID',
    `creator_id` VARCHAR(36) NOT NULL COMMENT '创建者ID',
    `responsible_id` VARCHAR(36) NOT NULL COMMENT '负责人ID',
    `status` ENUM('draft', 'pending_approval', 'approved', 'in_progress', 'pending_final_review', 'completed', 'rejected', 'cancelled', 'paused') DEFAULT 'draft' COMMENT '任务状态',
    
    -- 时间管理
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `start_date` TIMESTAMP NULL COMMENT '开始时间',
    `due_date` TIMESTAMP NULL COMMENT '截止时间',
    `completed_at` TIMESTAMP NULL COMMENT '完成时间',
    `estimated_hours` INT DEFAULT 0 COMMENT '预估工时',
    
    -- 重复任务相关
    `workflow_id` VARCHAR(36) DEFAULT NULL COMMENT '工作流ID',
    
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE RESTRICT,
    FOREIGN KEY (`creator_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    FOREIGN KEY (`responsible_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    
    INDEX `idx_project_status` (`project_id`, `status`),
    INDEX `idx_responsible_status` (`responsible_id`, `status`),
    INDEX `idx_creator` (`creator_id`),
    INDEX `idx_due_date` (`due_date`),
    INDEX `idx_status` (`status`),
    INDEX `idx_task_type` (`task_type`),
    INDEX `idx_priority` (`priority`),
    INDEX `idx_created_at` (`created_at`),
    
    -- 全文搜索索引
    FULLTEXT INDEX `ft_title_description` (`title`, `description`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务表';

-- 任务参与人员表
CREATE TABLE IF NOT EXISTS `task_participants` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '记录ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `added_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `added_by` VARCHAR(36) NOT NULL COMMENT '添加人ID',
    
    UNIQUE KEY `uk_task_user` (`task_id`, `user_id`),
    FOREIGN KEY (`task_id`) REFERENCES `tasks`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`added_by`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    
    INDEX `idx_task` (`task_id`),
    INDEX `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务参与人员表';

-- 重复任务规则表
CREATE TABLE IF NOT EXISTS `recurrence_rules` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '规则ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID',
    `frequency` ENUM('daily', 'weekly', 'monthly') NOT NULL COMMENT '重复频率',
    `interval_value` INT DEFAULT 1 COMMENT '间隔值',
    `end_date` TIMESTAMP NULL COMMENT '结束日期',
    `max_executions` INT DEFAULT NULL COMMENT '最大执行次数',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    UNIQUE KEY `uk_task` (`task_id`),
    FOREIGN KEY (`task_id`) REFERENCES `tasks`(`id`) ON DELETE CASCADE,
    
    INDEX `idx_frequency` (`frequency`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='重复任务规则表';

-- ================================================
-- 任务执行相关表
-- ================================================

-- 任务执行记录表
CREATE TABLE IF NOT EXISTS `task_executions` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '执行ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID',
    `execution_date` TIMESTAMP NOT NULL COMMENT '执行日期',
    `status` ENUM('pending', 'in_progress', 'pending_review', 'pending_final_review', 'completed', 'rejected', 'cancelled') DEFAULT 'pending' COMMENT '执行状态',
    `started_at` TIMESTAMP NULL COMMENT '开始时间',
    `submitted_at` TIMESTAMP NULL COMMENT '提交时间',
    `completed_at` TIMESTAMP NULL COMMENT '完成时间',
    `result` TEXT COMMENT '执行结果',
    
    INDEX `idx_task` (`task_id`),
    INDEX `idx_execution_date` (`execution_date`),
    INDEX `idx_status` (`status`),
    INDEX `idx_task_status` (`task_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务执行记录表';

-- 参与人员完成记录表
CREATE TABLE IF NOT EXISTS `participant_completions` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '完成记录ID',
    `execution_id` VARCHAR(36) NOT NULL COMMENT '执行ID',
    `participant_id` VARCHAR(36) NOT NULL COMMENT '参与者ID',
    `work_result` TEXT COMMENT '工作成果',
    `status` ENUM('pending', 'submitted', 'approved', 'rejected') DEFAULT 'pending' COMMENT '状态',
    `submitted_at` TIMESTAMP NULL COMMENT '提交时间',
    `reviewed_at` TIMESTAMP NULL COMMENT '审批时间',
    `reviewer_id` VARCHAR(36) DEFAULT NULL COMMENT '审批人ID',
    `review_comment` TEXT COMMENT '审批意见',
    
    UNIQUE KEY `uk_execution_participant` (`execution_id`, `participant_id`),
    FOREIGN KEY (`execution_id`) REFERENCES `task_executions`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`participant_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`reviewer_id`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    
    INDEX `idx_execution` (`execution_id`),
    INDEX `idx_participant` (`participant_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='参与人员完成记录表';

-- ================================================
-- 审批和延期相关表
-- ================================================

-- 审批记录表
CREATE TABLE IF NOT EXISTS `approval_records` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '审批记录ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID',
    `execution_id` VARCHAR(36) DEFAULT NULL COMMENT '执行ID',
    `approver_id` VARCHAR(36) NOT NULL COMMENT '审批人ID',
    `approval_type` ENUM('task_creation', 'task_completion', 'extension_request') NOT NULL COMMENT '审批类型',
    `action` ENUM('approve', 'reject') NOT NULL COMMENT '审批动作',
    `comment` TEXT COMMENT '审批意见',
    `approved_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '审批时间',
    
    FOREIGN KEY (`task_id`) REFERENCES `tasks`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`execution_id`) REFERENCES `task_executions`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`approver_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    
    INDEX `idx_task` (`task_id`),
    INDEX `idx_execution` (`execution_id`),
    INDEX `idx_approver` (`approver_id`),
    INDEX `idx_type` (`approval_type`),
    INDEX `idx_approved_at` (`approved_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审批记录表';

-- 延期申请表
CREATE TABLE IF NOT EXISTS `extension_requests` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '延期申请ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID',
    `requester_id` VARCHAR(36) NOT NULL COMMENT '申请人ID',
    `original_due_date` TIMESTAMP NOT NULL COMMENT '原截止时间',
    `requested_due_date` TIMESTAMP NOT NULL COMMENT '申请截止时间',
    `reason` TEXT NOT NULL COMMENT '申请理由',
    `status` ENUM('pending', 'approved', 'rejected') DEFAULT 'pending' COMMENT '申请状态',
    `requested_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    `reviewed_at` TIMESTAMP NULL COMMENT '审批时间',
    `reviewer_id` VARCHAR(36) DEFAULT NULL COMMENT '审批人ID',
    `review_comment` TEXT COMMENT '审批意见',
    
    FOREIGN KEY (`task_id`) REFERENCES `tasks`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`requester_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    FOREIGN KEY (`reviewer_id`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    
    INDEX `idx_task` (`task_id`),
    INDEX `idx_requester` (`requester_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_requested_at` (`requested_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='延期申请表';

-- ================================================
-- 事件和日志相关表
-- ================================================

-- 领域事件表（事件溯源）
CREATE TABLE IF NOT EXISTS `domain_events` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '事件ID',
    `event_type` VARCHAR(100) NOT NULL COMMENT '事件类型',
    `aggregate_id` VARCHAR(36) NOT NULL COMMENT '聚合根ID',
    `aggregate_type` VARCHAR(50) NOT NULL COMMENT '聚合根类型',
    `event_data` JSON NOT NULL COMMENT '事件数据',
    `event_version` INT DEFAULT 1 COMMENT '事件版本',
    `occurred_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发生时间',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '触发用户ID',
    
    INDEX `idx_aggregate` (`aggregate_type`, `aggregate_id`),
    INDEX `idx_event_type` (`event_type`),
    INDEX `idx_occurred_at` (`occurred_at`),
    INDEX `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='领域事件表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '日志ID',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '用户ID',
    `operation` VARCHAR(100) NOT NULL COMMENT '操作类型',
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(36) NOT NULL COMMENT '资源ID',
    `ip_address` VARCHAR(45) COMMENT 'IP地址',
    `user_agent` TEXT COMMENT '用户代理',
    `request_data` JSON COMMENT '请求数据',
    `response_status` INT COMMENT '响应状态码',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX `idx_user` (`user_id`),
    INDEX `idx_operation` (`operation`),
    INDEX `idx_resource` (`resource_type`, `resource_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- ================================================
-- 文件管理相关表
-- ================================================

-- 文件表
CREATE TABLE IF NOT EXISTS `files` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '文件ID',
    `filename` VARCHAR(255) NOT NULL COMMENT '文件名',
    `original_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
    `file_type` VARCHAR(50) NOT NULL COMMENT '文件类型',
    `file_size` BIGINT NOT NULL COMMENT '文件大小',
    `file_path` VARCHAR(500) NOT NULL COMMENT '文件路径',
    `mime_type` VARCHAR(100) NOT NULL COMMENT 'MIME类型',
    `md5_hash` VARCHAR(32) NOT NULL COMMENT 'MD5哈希',
    `uploader_id` VARCHAR(36) NOT NULL COMMENT '上传者ID',
    `upload_status` ENUM('uploading', 'completed', 'failed') DEFAULT 'uploading' COMMENT '上传状态',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (`uploader_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    
    INDEX `idx_uploader` (`uploader_id`),
    INDEX `idx_type` (`file_type`),
    INDEX `idx_status` (`upload_status`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_md5` (`md5_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- 文件关联表
CREATE TABLE IF NOT EXISTS `file_associations` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '关联ID',
    `file_id` VARCHAR(36) NOT NULL COMMENT '文件ID',
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(36) NOT NULL COMMENT '资源ID',
    `association_type` ENUM('attachment', 'avatar', 'document') NOT NULL COMMENT '关联类型',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (`file_id`) REFERENCES `files`(`id`) ON DELETE CASCADE,
    
    INDEX `idx_file` (`file_id`),
    INDEX `idx_resource` (`resource_type`, `resource_id`),
    INDEX `idx_type` (`association_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件关联表';

SET FOREIGN_KEY_CHECKS = 1;

-- ================================================
-- 创建数据库完成
-- ================================================
