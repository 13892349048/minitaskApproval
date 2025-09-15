package domainerror

import "fmt"

// DomainErrorType 领域错误类型
type DomainErrorType string

const (
	ErrPermissionNotFound   DomainErrorType = "PERMISSION_NOT_FOUND"
	ErrRoleNotFound         DomainErrorType = "ROLE_NOT_FOUND"
	ErrPolicyNotFound       DomainErrorType = "POLICY_NOT_FOUND"
	ErrRoleAlreadyAssigned  DomainErrorType = "ROLE_ALREADY_ASSIGNED"
	ErrRoleNotAssigned      DomainErrorType = "ROLE_NOT_ASSIGNED"
	ErrSystemRoleImmutable  DomainErrorType = "SYSTEM_ROLE_IMMUTABLE"
	ErrInvalidPermission    DomainErrorType = "INVALID_PERMISSION"
	ErrInvalidRole          DomainErrorType = "INVALID_ROLE"
	ErrInvalidPolicy        DomainErrorType = "INVALID_POLICY"
	ErrPermissionDenied     DomainErrorType = "PERMISSION_DENIED"
	ErrInvalidEvaluationCtx DomainErrorType = "INVALID_EVALUATION_CONTEXT"
)

// DomainError 领域错误
type DomainError struct {
	Type    DomainErrorType        `json:"type"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error 实现error接口
func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewDomainError 创建领域错误
func NewDomainError(errorType DomainErrorType, message string) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		Details: make(map[string]interface{}),
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
