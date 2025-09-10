package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse 成功响应结构
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// respondWithError 统一错误响应
func RespondWithError(c *gin.Context, statusCode int, code, message string) {
	// 记录错误日志
	logger.Error("HTTP Error Response",
		zap.Int("status_code", statusCode),
		zap.String("error_code", code),
		zap.String("message", message),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("remote_addr", c.ClientIP()),
	)

	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    code,
		Message: message,
	}

	c.JSON(statusCode, response)
	c.Abort()
}

// respondWithSuccess 统一成功响应
func RespondWithSuccess(c *gin.Context, data interface{}, message string) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}

	c.JSON(http.StatusOK, response)
}

// respondWithCreated 创建成功响应
func RespondWithCreated(c *gin.Context, data interface{}, message string) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}

	c.JSON(http.StatusCreated, response)
}

// 为什么这样设计？
//
// 1. 统一响应格式：
//    - 前端可以依赖固定的响应结构
//    - 错误处理标准化
//    - 便于API文档生成
//
// 2. 错误日志记录：
//    - 自动记录所有错误响应
//    - 包含请求上下文信息
//    - 便于问题排查和监控
//
// 3. HTTP状态码标准化：
//    - 使用标准HTTP状态码
//    - 业务错误码与HTTP状态码分离
//    - 便于HTTP客户端处理
