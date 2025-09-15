package mysql

import (
	"time"

	"gorm.io/gorm"
)

// UserModel 用户持久化模型 - 重命名避免与Domain User冲突
type UserModel struct {
	ID           string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	FullName     string         `gorm:"type:varchar(100);not null" json:"full_name"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Avatar       *string        `gorm:"type:varchar(500)" json:"avatar"`
	Role         string         `gorm:"type:enum('employee','manager','director','admin');default:'employee'" json:"role"`
	Status       string         `gorm:"type:enum('active','inactive','suspended');default:'active'" json:"status"`
	Phone        *string        `gorm:"type:varchar(20)" json:"phone"`
	Department   *string        `gorm:"type:varchar(100)" json:"department"`
	Position     *string        `gorm:"type:varchar(100)" json:"position"`
	JoinDate     *time.Time     `gorm:"type:date" json:"join_date"`
	DepartmentID *string        `gorm:"type:varchar(36)" json:"department_id"`
	ManagerID    *string        `gorm:"type:varchar(36)" json:"manager_id"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Roles              []Role            `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	OwnedProjects      []Project         `gorm:"foreignKey:OwnerID" json:"owned_projects,omitempty"`
	ManagedProjects    []Project         `gorm:"foreignKey:ManagerID" json:"managed_projects,omitempty"`
	ProjectMemberships []ProjectMember   `gorm:"foreignKey:UserID" json:"project_memberships,omitempty"`
	CreatedTasks       []Task            `gorm:"foreignKey:CreatorID" json:"created_tasks,omitempty"`
	ResponsibleTasks   []Task            `gorm:"foreignKey:ResponsibleID" json:"responsible_tasks,omitempty"`
	TaskParticipations []TaskParticipant `gorm:"foreignKey:UserID" json:"task_participations,omitempty"`
	UploadedFiles      []File            `gorm:"foreignKey:UploaderID" json:"uploaded_files,omitempty"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}
