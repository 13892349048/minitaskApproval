package dto

import "time"

// TaskCreateRequest 创建任务请求
// @Description 创建任务请求参数
type TaskCreateRequest struct {
	Title         string    `json:"title" binding:"required" example:"完成项目文档"`           // 任务标题
	Description   string    `json:"description" example:"编写项目的技术文档和用户手册"`             // 任务描述
	ProjectID     string    `json:"project_id" binding:"required" example:"proj-123"`     // 项目ID
	ResponsibleID string    `json:"responsible_id" binding:"required" example:"user-456"` // 负责人ID
	TaskType      string    `json:"task_type" example:"single"`                          // 任务类型
	Priority      string    `json:"priority" example:"normal"`                           // 优先级
	DueDate       time.Time `json:"due_date" example:"2023-12-31T23:59:59Z"`             // 截止日期
} // @name TaskCreateRequest

// TaskUpdateRequest 更新任务请求
// @Description 更新任务请求参数
type TaskUpdateRequest struct {
	Title         string     `json:"title" example:"更新的任务标题"`                           // 任务标题
	Description   string     `json:"description" example:"更新的任务描述"`                     // 任务描述
	ResponsibleID string     `json:"responsible_id" example:"user-789"`                  // 负责人ID
	Priority      string     `json:"priority" example:"high"`                            // 优先级
	DueDate       *time.Time `json:"due_date" example:"2023-12-31T23:59:59Z"`            // 截止日期
	Status        string     `json:"status" example:"in_progress"`                       // 状态
} // @name TaskUpdateRequest

// TaskInfo 任务信息
// @Description 任务详细信息
type TaskInfo struct {
	BaseEntity
	Title         string    `json:"title" example:"完成项目文档"`                    // 任务标题
	Description   string    `json:"description" example:"编写项目的技术文档和用户手册"`    // 任务描述
	ProjectID     string    `json:"project_id" example:"proj-123"`            // 项目ID
	ProjectName   string    `json:"project_name" example:"TaskFlow项目"`        // 项目名称
	CreatorID     string    `json:"creator_id" example:"user-123"`            // 创建者ID
	CreatorName   string    `json:"creator_name" example:"张三"`                // 创建者姓名
	ResponsibleID string    `json:"responsible_id" example:"user-456"`        // 负责人ID
	ResponsibleName string  `json:"responsible_name" example:"李四"`            // 负责人姓名
	TaskType      string    `json:"task_type" example:"single"`               // 任务类型
	Priority      string    `json:"priority" example:"normal"`                // 优先级
	Status        string    `json:"status" example:"pending"`                 // 状态
	DueDate       time.Time `json:"due_date" example:"2023-12-31T23:59:59Z"`  // 截止日期
	CompletedAt   *time.Time `json:"completed_at,omitempty"`                  // 完成时间
} // @name TaskInfo

// TaskListRequest 任务列表请求
// @Description 任务列表查询参数
type TaskListRequest struct {
	PaginationRequest
	ProjectID     string `json:"project_id" form:"project_id" example:"proj-123"`     // 项目ID
	ResponsibleID string `json:"responsible_id" form:"responsible_id" example:"user-456"` // 负责人ID
	Status        string `json:"status" form:"status" example:"pending"`              // 状态
	Priority      string `json:"priority" form:"priority" example:"high"`             // 优先级
	TaskType      string `json:"task_type" form:"task_type" example:"single"`         // 任务类型
	Keyword       string `json:"keyword" form:"keyword" example:"文档"`                // 关键词搜索
} // @name TaskListRequest

// TaskAssignRequest 任务分配请求
// @Description 任务分配请求参数
type TaskAssignRequest struct {
	ExecutorID string `json:"executor_id" binding:"required" example:"user-789"` // 执行者ID
	Comment    string `json:"comment" example:"请按时完成"`                          // 备注
} // @name TaskAssignRequest

// TaskStatusChangeRequest 任务状态变更请求
// @Description 任务状态变更请求参数
type TaskStatusChangeRequest struct {
	Status       string `json:"status" binding:"required" example:"in_progress"` // 新状态
	ChangeReason string `json:"change_reason" example:"开始执行任务"`                 // 变更原因
} // @name TaskStatusChangeRequest

// TaskCompletionRequest 任务完成提交请求
// @Description 任务完成提交请求参数
type TaskCompletionRequest struct {
	Summary     string   `json:"summary" binding:"required" example:"任务已按要求完成"` // 完成总结
	Attachments []string `json:"attachments" example:"file1.pdf,file2.docx"`    // 附件列表
} // @name TaskCompletionRequest
