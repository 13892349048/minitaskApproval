package errors

import (
	"fmt"
	"net/http"
)

// 领域错误类型
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	return e.Message
}

// 应用错误类型
type AppError struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// 错误构造函数
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Type:       "validation_error",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewPermissionDeniedError(message string) *AppError {
	return &AppError{
		Type:       "permission_denied",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:       "not_found",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:       "internal_error",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewInvalidTokenError(message string) *AppError {
	return &AppError{
		Type:       "invalid_token",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewExpiredTokenError(message string) *AppError {
	return &AppError{
		Type:       "expired_token",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewInvalidTokenTypeError(message string) *AppError {
	return &AppError{
		Type:       "invalid_token_type",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewTokenRevokedError(message string) *AppError {
	return &AppError{
		Type:       "token_revoked",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// 业务错误常量
var (
	ErrUserNotFound     = NewDomainError("USER_NOT_FOUND", "用户不存在")
	ErrUserExists       = NewDomainError("USER_EXISTS", "用户已存在")
	ErrInvalidPassword  = NewDomainError("INVALID_PASSWORD", "密码错误")
	ErrProjectNotFound  = NewDomainError("PROJECT_NOT_FOUND", "项目不存在")
	ErrTaskNotFound     = NewDomainError("TASK_NOT_FOUND", "任务不存在")
	ErrInvalidStatus    = NewDomainError("INVALID_STATUS", "无效的状态转换")
	ErrPermissionDenied = NewDomainError("PERMISSION_DENIED", "权限不足")
	// 常用错误
	ErrInvalidToken     = NewInvalidTokenError("无效的token")
	ErrExpiredToken     = NewExpiredTokenError("过期的token")
	ErrInvalidTokenType = NewInvalidTokenTypeError("无效的token类型")
	ErrTokenRevoked     = NewTokenRevokedError("撤销的token")
)
