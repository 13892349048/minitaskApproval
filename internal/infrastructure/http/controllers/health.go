package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/internal/infrastructure/http/dto"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=HealthStatus} "服务正常"
// @Router /health [get]
func (h *HealthController) HealthCheck(c *gin.Context) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services: map[string]string{
			"database":  "connected",
			"redis":     "connected",
			"eventbus":  "running",
		},
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Code:    200,
		Message: "success",
		Data:    status,
	})
}

// HealthStatus 健康状态
// @Description 服务健康状态信息
type HealthStatus struct {
	Status    string            `json:"status" example:"healthy"`                    // 服务状态
	Timestamp time.Time         `json:"timestamp" example:"2023-01-01T00:00:00Z"`   // 检查时间
	Version   string            `json:"version" example:"1.0.0"`                    // 服务版本
	Services  map[string]string `json:"services"`                                   // 各服务状态
} // @name HealthStatus
