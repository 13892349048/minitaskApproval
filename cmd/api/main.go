package main

import (
	"log"

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
