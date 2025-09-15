package dto

import "time"

// APIResponse 统一API响应格式
// @Description API统一响应格式
type APIResponse struct {
	Code    int         `json:"code" example:"200"`                    // 响应码
	Message string      `json:"message" example:"success"`             // 响应消息
	Data    interface{} `json:"data,omitempty"`                       // 响应数据
	Error   string      `json:"error,omitempty" example:"error info"` // 错误信息
} // @name APIResponse

// PaginationRequest 分页请求参数
// @Description 分页请求参数
type PaginationRequest struct {
	Page     int `json:"page" form:"page" example:"1"`         // 页码，从1开始
	PageSize int `json:"page_size" form:"page_size" example:"10"` // 每页数量
} // @name PaginationRequest

// PaginationResponse 分页响应数据
// @Description 分页响应数据
type PaginationResponse struct {
	Total       int64       `json:"total" example:"100"`       // 总记录数
	Page        int         `json:"page" example:"1"`          // 当前页码
	PageSize    int         `json:"page_size" example:"10"`    // 每页数量
	TotalPages  int         `json:"total_pages" example:"10"`  // 总页数
	HasNext     bool        `json:"has_next" example:"true"`   // 是否有下一页
	HasPrev     bool        `json:"has_prev" example:"false"`  // 是否有上一页
	Data        interface{} `json:"data"`                      // 数据列表
} // @name PaginationResponse

// IDRequest ID请求参数
// @Description ID请求参数
type IDRequest struct {
	ID string `json:"id" uri:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"` // 资源ID
} // @name IDRequest

// BaseEntity 基础实体信息
// @Description 基础实体信息
type BaseEntity struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"` // 唯一标识
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`         // 创建时间
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`         // 更新时间
} // @name BaseEntity

// ErrorResponse 错误响应
// @Description 错误响应格式
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`           // 错误码
	Message string `json:"message" example:"Bad Request"` // 错误消息
	Details string `json:"details,omitempty"`            // 错误详情
} // @name ErrorResponse
