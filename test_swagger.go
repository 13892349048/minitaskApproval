package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/taskflow/docs" // 导入Swagger文档
)

func mai1n() {
	// 设置为开发模式
	gin.SetMode(gin.DebugMode)

	// 创建路由
	r := gin.Default()

	// 添加基本的健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 设置Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 重定向根路径到 Swagger UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// 重定向 /docs 到 Swagger UI
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	fmt.Println("服务器启动在 http://localhost:8080")
	fmt.Println("Swagger UI 地址: http://localhost:8080/swagger/index.html")
	fmt.Println("健康检查: http://localhost:8080/health")

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
