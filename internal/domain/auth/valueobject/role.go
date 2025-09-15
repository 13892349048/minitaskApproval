package valueobject

// RoleID 角色ID值对象
type RoleID string

func (id RoleID) String() string {
	return string(id)
}

// 预定义系统角色
const (
	RoleSuperAdmin     RoleID = "super_admin"
	RoleAdmin          RoleID = "admin"
	RoleProjectOwner   RoleID = "project_owner"
	RoleProjectManager RoleID = "project_manager"
	RoleTeamLeader     RoleID = "team_leader"
	RoleEmployee       RoleID = "employee"
)
