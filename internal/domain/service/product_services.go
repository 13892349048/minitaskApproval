package service

import (
	"fmt"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/valueobject"
)

/*
明确的职责分工：

领域服务负责：

跨聚合根的复杂业务规则（项目合并、所有权转移）
需要访问多个仓储的操作（层级关系验证）
复杂的业务策略（权限继承算法）
聚合根内部保留：

单聚合根的业务规则（成员管理、状态变更）
数据完整性验证（字段验证、状态一致性）
简单的权限检查（用户访问权限）
*/
// ProjectDomainService 项目领域服务接口
// type ProjectDomainService interface {
// 	// 权限检查
// 	CanUserManageProject(project aggregate.ProjectAggregate, userID valueobject.UserID) bool
// 	CanUserCreateSubProject(parent aggregate.ProjectAggregate, userID valueobject.UserID) bool
// 	CanUserAddMember(project aggregate.ProjectAggregate, userID valueobject.UserID) bool

// 	// 业务规则验证
// 	ValidateProjectHierarchy(parent aggregate.ProjectAggregate, child aggregate.ProjectAggregate) error
// 	ValidateProjectDeletion(project aggregate.ProjectAggregate) error
// 	ValidateProjectCompletion(project aggregate.ProjectAggregate) error

// 	// 复杂业务逻辑
// 	TransferProjectOwnership(project aggregate.ProjectAggregate, newOwnerID valueobject.UserID, transferredBy valueobject.UserID) error
// 	MergeProjects(sourceProject, targetProject aggregate.ProjectAggregate, mergedBy valueobject.UserID) error
// 	ArchiveProject(project aggregate.ProjectAggregate, archivedBy valueobject.UserID) error
// }

// ProjectDomainService 项目领域服务实现
// type ProjectDomainServiceImpl struct {
// 	// 可以注入其他领域服务或仓储
// }

// // NewProjectDomainService 创建项目领域服务
// func NewProjectDomainService() *ProjectDomainServiceImpl {
// 	return &ProjectDomainServiceImpl{}
// }

// ===========================================
// 跨聚合根的复杂业务规则（属于领域服务）
// ===========================================

// TransferProjectOwnership 转移项目所有权 - 涉及复杂业务规则
func (s *ProjectDomainServiceImpl) TransferProjectOwnership(
	project *aggregate.Project,
	newOwnerID valueobject.UserID,
	transferredBy valueobject.UserID,
) error {
	// 复杂业务规则：
	// 1. 只有当前所有者可以转移
	// 2. 新所有者不能是当前管理者
	// 3. 需要检查新所有者的权限级别
	// 4. 转移后需要重新分配管理者角色

	if transferredBy != project.OwnerID {
		return fmt.Errorf("only current owner can transfer ownership")
	}

	if project.ManagerID != nil && newOwnerID == *project.ManagerID {
		return fmt.Errorf("new owner cannot be current manager")
	}

	// 这里可能需要调用用户领域服务检查权限级别
	// if !s.userService.HasSufficientLevel(newOwnerID) { ... }

	return nil
}


// MergeProjects 合并项目 - 复杂的跨聚合根操作
func (s *ProjectDomainServiceImpl) MergeProjects(
	sourceProject, targetProject *aggregate.Project,
	mergedBy valueobject.UserID,
) error {
	// 复杂合并规则：
	// 1. 权限检查
	// 2. 状态兼容性
	// 3. 成员合并策略
	// 4. 任务迁移规则

	if !s.canMergeProjects(sourceProject, targetProject, mergedBy) {
		return fmt.Errorf("insufficient permission to merge projects")
	}

	if !s.areProjectsCompatible(sourceProject, targetProject) {
		return fmt.Errorf("projects are not compatible for merging")
	}

	return nil
}

// ===========================================
// 单聚合根内的业务规则（应该在聚合根内部）
// ===========================================

// 以下方法应该移到Project聚合根内部：
// - CanUserAccess (已在聚合根中)
// - CanUserManageProject (已在聚合根中)
// - ValidateProjectCompletion (应该在聚合根中)

// ===========================================
// 私有辅助方法
// ===========================================

func (s *ProjectDomainServiceImpl) wouldCreateCycle(parent, child *aggregate.Project) bool {
	// 实现循环检测算法
	// 这里需要访问项目仓储来检查完整的层级关系
	return false
}

func (s *ProjectDomainServiceImpl) canMergeProjects(source, target *aggregate.Project, userID valueobject.UserID) bool {
	// 检查用户是否有权限合并这两个项目
	return source.CanUserAccess(userID) && target.CanUserAccess(userID)
}

func (s *ProjectDomainServiceImpl) areProjectsCompatible(source, target *aggregate.Project) bool {
	// 检查项目是否可以合并（状态、类型等）
	return source.ProjectType == target.ProjectType
}
