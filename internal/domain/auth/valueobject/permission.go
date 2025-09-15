package valueobject

// PermissionID 权限ID值对象
type PermissionID string

func (id PermissionID) String() string {
	return string(id)
}

// ResourceType 资源类型值对象
type ResourceType string

const (
	ResourceTypeProject ResourceType = "project"
	ResourceTypeTask    ResourceType = "task"
	ResourceTypeUser    ResourceType = "user"
	ResourceTypeFile    ResourceType = "file"
)

// ActionType 操作类型值对象
type ActionType string

const (
	ActionTypeCreate  ActionType = "create"
	ActionTypeRead    ActionType = "read"
	ActionTypeUpdate  ActionType = "update"
	ActionTypeDelete  ActionType = "delete"
	ActionTypeAssign  ActionType = "assign"
	ActionTypeApprove ActionType = "approve"
	ActionTypeExecute ActionType = "execute"
)
