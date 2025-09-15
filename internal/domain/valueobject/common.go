package valueobject

import (
	"time"
)

// 通用分页和排序值对象

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// DefaultPagination 默认分页参数
func DefaultPagination() PaginationRequest {
	return PaginationRequest{
		Page:     1,
		PageSize: 20,
	}
}

// Offset 计算偏移量
func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// SortRequest 排序请求
type SortRequest struct {
	Field     string    `json:"field"`
	Direction SortOrder `json:"direction"`
}

// SortOrder 排序方向
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(page, pageSize int, total int64) PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// 通用时间范围值对象

// DateRange 日期范围
type DateRange struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

// IsValid 验证日期范围是否有效
func (dr DateRange) IsValid() bool {
	if dr.StartDate != nil && dr.EndDate != nil {
		return dr.StartDate.Before(*dr.EndDate) || dr.StartDate.Equal(*dr.EndDate)
	}
	return true
}

// 通用状态值对象

// EntityStatus 实体状态
type EntityStatus string

const (
	EntityStatusActive   EntityStatus = "active"
	EntityStatusInactive EntityStatus = "inactive"
	EntityStatusDeleted  EntityStatus = "deleted"
)

// 通用审计信息

// AuditInfo 审计信息
type AuditInfo struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedBy UserID     `json:"created_by"`
	UpdatedBy UserID     `json:"updated_by"`
}

// 通用搜索条件

// SearchCriteria 通用搜索条件接口
type SearchCriteria interface {
	GetPagination() PaginationRequest
	GetSort() []SortRequest
}

// BaseSearchCriteria 基础搜索条件
type BaseSearchCriteria struct {
	Pagination PaginationRequest `json:"pagination"`
	Sort       []SortRequest     `json:"sort"`
	DateRange  *DateRange        `json:"date_range,omitempty"`
}

// GetPagination 获取分页参数
func (b BaseSearchCriteria) GetPagination() PaginationRequest {
	if b.Pagination.Page == 0 {
		return DefaultPagination()
	}
	return b.Pagination
}

// GetSort 获取排序参数
func (b BaseSearchCriteria) GetSort() []SortRequest {
	return b.Sort
}

// 通用响应包装

// Response 通用响应结构
type Response[T any] struct {
	Success    bool                `json:"success"`
	Data       T                   `json:"data,omitempty"`
	Error      *ErrorInfo          `json:"error,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	Timestamp  time.Time           `json:"timestamp"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse[T any](data T) Response[T] {
	return Response[T]{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewSuccessResponseWithPagination 创建带分页的成功响应
func NewSuccessResponseWithPagination[T any](data T, pagination PaginationResponse) Response[T] {
	return Response[T]{
		Success:    true,
		Data:       data,
		Pagination: &pagination,
		Timestamp:  time.Now(),
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse[T any](code, message, details string) Response[T] {
	return Response[T]{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

// 通用ID生成器接口

// IDGenerator ID生成器接口
type IDGenerator interface {
	GenerateUserID() UserID
	GenerateProjectID() ProjectID
	GenerateTaskID() TaskID
}

// 通用验证器接口

// Validator 通用验证器接口
type Validator interface {
	Validate(interface{}) error
}

// 通用缓存键

// CacheKey 缓存键类型
type CacheKey string

// UserCacheKeys 用户相关缓存键
const (
	UserCacheKeyPrefix     CacheKey = "user:"
	UserRolesCacheKey      CacheKey = "user:roles:"
	UserPermissionsCacheKey CacheKey = "user:permissions:"
)

// BuildUserCacheKey 构建用户缓存键
func BuildUserCacheKey(userID UserID) CacheKey {
	return CacheKey(string(UserCacheKeyPrefix) + string(userID))
}

// BuildUserRolesCacheKey 构建用户角色缓存键
func BuildUserRolesCacheKey(userID UserID) CacheKey {
	return CacheKey(string(UserRolesCacheKey) + string(userID))
}
