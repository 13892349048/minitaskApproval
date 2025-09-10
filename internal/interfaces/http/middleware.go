package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskflow/pkg/errors"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// corsMiddleware CORS中间件
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 允许的域名列表（实际项目中应该从配置读取）
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://taskflow.example.com",
		}

		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// requestIDMiddleware 请求ID中间件
func (s *Server) requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// loggingMiddleware 日志中间件
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 使用结构化日志
		logger.Info("HTTP Request",
			zap.String("request_id", param.Keys["request_id"].(string)),
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}

// securityHeadersMiddleware 安全头中间件
func (s *Server) securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止XSS攻击
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// HSTS（仅在HTTPS时）
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// CSP（内容安全策略）
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// authMiddleware JWT认证中间件
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.RespondWithError(c, http.StatusUnauthorized, "MISSING_AUTH_HEADER", "Authorization header is required")
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errors.RespondWithError(c, http.StatusUnauthorized, "INVALID_AUTH_FORMAT", "Invalid authorization header format")
			return
		}

		token := tokenParts[1]
		if token == "" {
			errors.RespondWithError(c, http.StatusUnauthorized, "EMPTY_TOKEN", "Token cannot be empty")
			return
		}

		// 验证JWT token
		claims, err := s.jwtService.ValidateToken(token)
		if err != nil {
			switch err {
			case errors.ErrExpiredToken:
				errors.RespondWithError(c, http.StatusUnauthorized, "TOKEN_EXPIRED", "Token has expired")
			case errors.ErrInvalidTokenType:
				errors.RespondWithError(c, http.StatusUnauthorized, "INVALID_TOKEN_TYPE", "Invalid token type")
			default:
				errors.RespondWithError(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token")
			}
			return
		}

		// 设置用户上下文信息
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("user_claims", claims)

		// 记录认证成功日志
		logger.Debug("User authenticated successfully",
			zap.String("user_id", claims.UserID),
			zap.String("email", claims.Email),
			zap.Strings("roles", claims.Roles),
		)

		c.Next()
	}
}

// rateLimitMiddleware 限流中间件
func (s *Server) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现限流逻辑
		// 可以使用Redis实现滑动窗口限流
		c.Next()
	}
}

// metricsMiddleware 监控指标中间件
func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 记录请求指标
		duration := time.Since(start)

		// TODO: 发送指标到Prometheus
		logger.Debug("Request metrics",
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	}
}

// errorHandlingMiddleware 错误处理中间件
func (s *Server) errorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			logger.Error("Request error",
				zap.String("request_id", c.GetString("request_id")),
				zap.Error(err.Err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			// 根据错误类型返回适当的HTTP状态码
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid request format",
					"code":  "INVALID_REQUEST",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
					"code":  "INTERNAL_ERROR",
				})
			}
		}
	}
}
