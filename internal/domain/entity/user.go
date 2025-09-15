package entity

import "time"

// UserData 用户数据传输对象（用于持久化和恢复）
type UserData struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	PasswordHash string     `json:"password_hash"`
	Avatar       *string    `json:"avatar"`
	Status       string     `json:"status"`
	Phone        *string    `json:"phone"`
	Department   *string    `json:"department"`
	Position     *string    `json:"position"`
	JoinDate     *time.Time `json:"join_date"`
	DepartmentID *string    `json:"department_id"`
	ManagerID    *string    `json:"manager_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	Roles        []string   `json:"roles"`
}
