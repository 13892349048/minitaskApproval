// 简单的Swagger测试程序
// 用于验证Swagger文档是否正确加载
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	_ "github.com/taskflow/docs" // 导入生成的Swagger文档
)

func main1() {
	// 检查Swagger文档是否已注册
	spec := swag.GetSwagger("swagger")
	if spec == nil {
		log.Fatal("Swagger文档未注册！请检查docs包导入")
	}

	fmt.Printf("✅ Swagger文档已成功注册\n")

	// 创建Gin路由
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 添加CORS中间件以防止跨域问题
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"swagger":   "enabled",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// Swagger文档调试端点
	r.GET("/swagger-debug", func(c *gin.Context) {
		spec := swag.GetSwagger("swagger")
		if spec == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Swagger spec not found"})
			return
		}
	})

	// 原始Swagger JSON端点
	r.GET("/swagger-json", func(c *gin.Context) {
		spec := swag.GetSwagger("swagger")
		if spec == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Swagger spec not found"})
			return
		}

		// 生成完整的Swagger JSON
		swaggerJSON, err := json.Marshal(spec)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, string(swaggerJSON))
	})

	// Swagger UI路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 根路径重定向到Swagger UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	// 文档路径重定向
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	fmt.Println("\n🚀 服务器启动成功！")
	fmt.Println("📍 服务地址: http://localhost:8080")
	fmt.Println("📖 Swagger UI: http://localhost:8080/swagger/index.html")
	fmt.Println("🔍 调试信息: http://localhost:8080/swagger-debug")
	fmt.Println("📄 Swagger JSON: http://localhost:8080/swagger-json")
	fmt.Println("❤️  健康检查: http://localhost:8080/health")
	fmt.Println("\n按 Ctrl+C 停止服务器")

	// 启动HTTP服务器
	log.Fatal(r.Run(":8080"))
}
