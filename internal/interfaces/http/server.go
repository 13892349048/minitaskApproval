package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	userAppService "github.com/taskflow/internal/application/service"
	"github.com/taskflow/internal/domain/auth/service"
	"github.com/taskflow/internal/infrastructure/config"
	"github.com/taskflow/internal/infrastructure/http/controllers"
	"github.com/taskflow/internal/interfaces/http/handler"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"

	_ "github.com/taskflow/docs" // Import generated docs
)

// Server HTTP服务器
type Server struct {
	config      *config.Config
	router      *gin.Engine
	server      *http.Server
	jwtService  service.JWTService
	userService *userAppService.UserAppService
	authHandler *handler.AuthHandler
}

// NewServer 创建新的HTTP服务器
func NewServer(cfg *config.Config, jwtService service.JWTService, userService *userAppService.UserAppService) *Server {
	// 设置Gin模式
	if cfg.App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建认证处理器
	authHandler := handler.NewAuthHandler(jwtService, userService)

	server := &Server{
		config:      cfg,
		router:      gin.New(),
		jwtService:  jwtService,
		userService: userService,
		authHandler: authHandler,
	}

	// 设置中间件
	server.setupMiddleware()

	// 设置路由
	server.setupRoutes()

	// 设置Swagger文档
	server.setupSwagger()

	return server
}

// Start 启动服务器
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.config.App.Port),
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	logger.Info("Starting HTTP server",
		zap.Int("port", s.config.App.Port),
		zap.String("mode", s.config.App.Mode))

	return s.server.ListenAndServe()
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

// setupMiddleware 设置中间件
func (s *Server) setupMiddleware() {
	// 基础中间件
	s.router.Use(gin.Recovery())
	s.router.Use(s.corsMiddleware())
	s.router.Use(s.requestIDMiddleware())
	s.router.Use(s.loggingMiddleware())

	// 安全中间件
	s.router.Use(s.securityHeadersMiddleware())
}

func (s *Server) setupRoutes() {
	// 创建健康检查控制器
	healthController := controllers.NewHealthController()

	// 健康检查（无需认证，根路径）
	s.router.GET("/health", healthController.HealthCheck)
	s.router.GET("/version", s.versionInfo)

	// API版本分组
	v1 := s.router.Group("/api/v1")
	{
		// 健康检查（API版本路径）
		v1.GET("/health", healthController.HealthCheck)
		v1.GET("/version", s.versionInfo)

		// 认证相关（无需token）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/refresh", s.authHandler.RefreshToken)
		}

		// 需要认证的认证接口
		authProtected := v1.Group("/auth")
		authProtected.Use(s.authMiddleware())
		{
			authProtected.POST("/logout", s.authHandler.Logout)
			authProtected.GET("/profile", s.authHandler.GetProfile)
		}

		// 需要认证的接口
		protected := v1.Group("")
		protected.Use(s.authMiddleware()) // JWT认证中间件
		{
			// 用户管理
			users := protected.Group("/users")
			{
				users.GET("", handler.ListUsers)
				users.GET("/:id", handler.GetUser)
				users.PUT("/:id", handler.UpdateUser)
				users.DELETE("/:id", handler.DeleteUser)
			}
			// 项目管理
			projects := protected.Group("/projects")
			{
				projects.GET("", handler.ListProjects)
				projects.POST("", handler.CreateProject)
				projects.GET("/:id", handler.GetProject)
				projects.PUT("/:id", handler.UpdateProject)
				projects.DELETE("/:id", handler.DeleteProject)

				// 项目成员管理
				projects.GET("/:id/members", handler.GetProjectMembers)
				projects.POST("/:id/members", handler.AddProjectMember)
				projects.DELETE("/:id/members/:user_id", handler.RemoveProjectMember)

				// 项目层级管理
				projects.GET("/:id/children", handler.GetSubProjects)
				projects.POST("/:id/children", handler.CreateSubProject)
			}

			// 任务管理
			tasks := protected.Group("/tasks")
			{
				tasks.GET("", handler.ListTasks)
				tasks.POST("", handler.CreateTask)
				tasks.GET("/:id", handler.GetTask)
				tasks.PUT("/:id", handler.UpdateTask)
				tasks.DELETE("/:id", handler.DeleteTask)

				// 任务状态管理
				tasks.POST("/:id/submit", handler.SubmitTask)
				tasks.POST("/:id/approve", handler.ApproveTask)
				tasks.POST("/:id/reject", handler.RejectTask)
				tasks.POST("/:id/assign", handler.AssignTask)

				// 任务参与者管理
				tasks.GET("/:id/participants", handler.GetTaskParticipants)
				tasks.POST("/:id/participants", handler.AddTaskParticipant)
				tasks.DELETE("/:id/participants/:user_id", handler.RemoveTaskParticipant)

				// 任务执行管理
				tasks.POST("/:id/executions", handler.CreateTaskExecution)
				tasks.GET("/:id/executions", handler.GetTaskExecutions)
				tasks.POST("/:id/executions/:exec_id/work", handler.SubmitWork)
				tasks.POST("/:id/executions/:exec_id/review", handler.ReviewWork)

				// 延期申请
				tasks.POST("/:id/extensions", handler.RequestExtension)
				tasks.GET("/:id/extensions", handler.GetTaskExtensions)
				tasks.PUT("/extensions/:ext_id/approve", handler.ApproveExtension)
				tasks.PUT("/extensions/:ext_id/reject", handler.RejectExtension)
			}
			// 文件管理
			files := protected.Group("/files")
			{
				files.POST("/upload/init", handler.InitFileUpload)
				files.PUT("/upload/:upload_id/chunks/:chunk", handler.UploadChunk)
				files.POST("/upload/:upload_id/complete", handler.CompleteUpload)
				files.GET("/upload/:upload_id/status", handler.GetUploadStatus)
				files.GET("/:id", handler.GetFile)
				files.DELETE("/:id", handler.DeleteFile)
			}

			// 统计分析
			stats := protected.Group("/stats")
			{
				stats.GET("/dashboard", handler.GetDashboard)
				stats.GET("/projects/:id/stats", handler.GetProjectStats)
				stats.GET("/users/:id/workload", handler.GetUserWorkload)
				stats.GET("/tasks/completion-rate", handler.GetTaskCompletionRate)
			}

			// 搜索
			search := protected.Group("/search")
			{
				search.GET("/tasks", handler.SearchTasks)
				search.GET("/projects", handler.SearchProjects)
				search.GET("/users", handler.SearchUsers)
			}
		}
	}
}

// setupSwagger 设置Swagger文档路由
func (s *Server) setupSwagger() {
	// 只在开发和测试环境启用Swagger
	if s.config.App.Mode != "production" {
		// 添加CORS中间件以防止跨域问题
		s.router.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		// 添加调试端点检查Swagger文档注册状态
		s.router.GET("/swagger-debug", func(c *gin.Context) {
			spec := swag.GetSwagger("swagger")
			c.JSON(http.StatusOK, gin.H{
				"swagger_enabled": true,
				"mode":            s.config.App.Mode,
				"spec_found":      spec != nil,
				"message":         "Swagger is enabled in non-production mode",
			})
		})

		// Swagger UI 路由 - 配置CSP以允许内联脚本
		s.router.GET("/swagger/*any", func(c *gin.Context) {
			// 设置Content Security Policy以允许Swagger UI的内联脚本
			c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		})

		// 重定向根路径到 Swagger UI
		s.router.GET("/docs", func(c *gin.Context) {
			c.Redirect(302, "/swagger/index.html")
		})

		// API文档根路径重定向
		s.router.GET("/", func(c *gin.Context) {
			c.Redirect(302, "/swagger/index.html")
		})

		logger.Info("Swagger documentation enabled",
			zap.String("url", fmt.Sprintf("http://localhost:%d/swagger/index.html", s.config.App.Port)),
			zap.String("mode", s.config.App.Mode))
	} else {
		logger.Info("Swagger documentation disabled in production mode")
	}
}

// 版本信息处理器
func (s *Server) versionInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app":        s.config.App.Name,
		"version":    s.config.App.Version,
		"build_time": "wait for build insert", // 实际项目中可以通过构建时注入
		"go_version": "1.21",
	})
}
