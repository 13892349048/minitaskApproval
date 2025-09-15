package event

import "fmt"

// DomainErrorType 领域错误类型
type DomainErrorType string

const (
	// 通用错误
	ErrInvalidInput     DomainErrorType = "INVALID_INPUT"
	ErrNotFound         DomainErrorType = "NOT_FOUND"
	ErrAlreadyExists    DomainErrorType = "ALREADY_EXISTS"
	ErrPermissionDenied DomainErrorType = "PERMISSION_DENIED"
	ErrInvalidState     DomainErrorType = "INVALID_STATE"
	ErrBusinessRule     DomainErrorType = "BUSINESS_RULE_VIOLATION"

	// 认证相关
	ErrInvalidCredentials DomainErrorType = "INVALID_CREDENTIALS"
	ErrTokenExpired       DomainErrorType = "TOKEN_EXPIRED"
	ErrTokenInvalid       DomainErrorType = "TOKEN_INVALID"

	// 用户相关
	ErrUserNotFound DomainErrorType = "USER_NOT_FOUND"
	ErrUserExists   DomainErrorType = "USER_ALREADY_EXISTS"
	ErrUserInactive DomainErrorType = "USER_INACTIVE"

	// 项目相关
	ErrProjectNotFound    DomainErrorType = "PROJECT_NOT_FOUND"
	ErrProjectExists      DomainErrorType = "PROJECT_ALREADY_EXISTS"
	ErrProjectInvalidType DomainErrorType = "PROJECT_INVALID_TYPE"

	// 任务相关
	ErrTaskNotFound      DomainErrorType = "TASK_NOT_FOUND"
	ErrTaskInvalidStatus DomainErrorType = "TASK_INVALID_STATUS"
)

// DomainError 领域错误
type DomainError struct {
	Type    DomainErrorType        `json:"type"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Cause   error                  `json:"-"`
}

// Error 实现error接口
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap 实现错误链
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewDomainError 创建领域错误
func NewDomainError(errorType DomainErrorType, message string) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// NewDomainErrorWithCause 创建带原因的领域错误
func NewDomainErrorWithCause(errorType DomainErrorType, message string, cause error) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		Details: make(map[string]interface{}),
		Cause:   cause,
	}
}

// WithDetails 添加错误详情
func (e *DomainError) WithDetails(key string, value interface{}) *DomainError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithDetail 添加单个详情（链式调用）
func (e *DomainError) WithDetail(key string, value interface{}) *DomainError {
	return e.WithDetails(key, value)
}

// IsDomainError 检查是否为领域错误
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// GetDomainError 获取领域错误
func GetDomainError(err error) *DomainError {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr
	}
	return nil
}

// IsErrorType 检查错误类型
func IsErrorType(err error, errorType DomainErrorType) bool {
	if domainErr := GetDomainError(err); domainErr != nil {
		return domainErr.Type == errorType
	}
	return false
}
