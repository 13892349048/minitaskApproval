// TaskFlow API Server
// @title TaskFlow API
// @version 1.0
// @description 多层级项目管理系统 API 文档
// @description
// @description 这是一个基于 DDD 和事件驱动架构的企业级项目管理系统。
// @description
// @description 主要功能包括：
// @description - 用户认证和授权
// @description - 项目管理
// @description - 任务管理
// @description - 权限控制
// @description - 事件驱动的业务流程
// @termsOfService http://swagger.io/terms/
// @contact.name TaskFlow API Support
// @contact.url http://www.taskflow.com/support
// @contact.email support@taskflow.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"log"

	_ "github.com/taskflow/docs" // 导入Swagger文档
	"github.com/taskflow/internal/app"
)

/*
优先级（高到低）：
挂载的my.cnf文件 (最高优先级)
command中的命令行参数
environment环境变量
MySQL默认配置（最低优先级）

环境变量可以覆盖 config.yaml
但不会影响MySQL/Redis容器配置
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
*/

func main() {
	// 创建应用程序实例
	application, err := app.NewApp("./configs")
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// 运行应用程序
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
