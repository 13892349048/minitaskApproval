-- ================================================
-- 添加软删除字段
-- 版本: 003
-- 创建时间: 2024-01-15
-- 描述: 为支持GORM软删除的表添加deleted_at字段
-- ================================================

SET NAMES utf8mb4;

-- 为用户表添加软删除字段
ALTER TABLE `users` 
ADD COLUMN `deleted_at` TIMESTAMP NULL COMMENT '软删除时间',
ADD INDEX `idx_deleted_at` (`deleted_at`);

-- 为权限策略表添加软删除字段
ALTER TABLE `permission_policies` 
ADD COLUMN `deleted_at` TIMESTAMP NULL COMMENT '软删除时间',
ADD INDEX `idx_deleted_at` (`deleted_at`);

-- 为项目表添加软删除字段
ALTER TABLE `projects` 
ADD COLUMN `deleted_at` TIMESTAMP NULL COMMENT '软删除时间',
ADD INDEX `idx_deleted_at` (`deleted_at`);

-- 为任务表添加软删除字段
ALTER TABLE `tasks` 
ADD COLUMN `deleted_at` TIMESTAMP NULL COMMENT '软删除时间',
ADD INDEX `idx_deleted_at` (`deleted_at`);

-- 为文件表添加软删除字段
ALTER TABLE `files` 
ADD COLUMN `deleted_at` TIMESTAMP NULL COMMENT '软删除时间',
ADD INDEX `idx_deleted_at` (`deleted_at`);

-- ================================================
-- 迁移完成
-- ================================================
