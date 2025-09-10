package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/taskflow/internal/infrastructure/config"
	"github.com/taskflow/internal/infrastructure/persistence/mysql"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	var (
		configPath = flag.String("config", "./configs", "配置文件路径")
		command    = flag.String("cmd", "validate", "命令: validate, sync, generate")
		modelName  = flag.String("model", "", "模型名称（用于generate命令）")
		force      = flag.Bool("force", false, "强制执行（用于sync命令）")
	)
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.InitLogger(&logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		FilePath:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
	}); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 连接数据库
	db, err := mysql.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 创建迁移管理器
	migrator := mysql.NewMigrator(db)

	// 执行命令
	switch *command {
	case "validate":
		if err := validateModels(migrator); err != nil {
			os.Exit(1)
		}
	case "sync":
		if err := syncModels(migrator, cfg.App.Mode == "development", *force); err != nil {
			os.Exit(1)
		}
	case "generate":
		if *modelName == "" {
			logger.Error("generate命令需要指定model参数")
			os.Exit(1)
		}
		if err := generateMigration(migrator, *modelName); err != nil {
			os.Exit(1)
		}
	case "status":
		if err := checkStatus(migrator); err != nil {
			os.Exit(1)
		}
	default:
		logger.Error("未知命令", zap.String("command", *command))
		fmt.Println("可用命令: validate, sync, generate, status")
		os.Exit(1)
	}
}

func validateModels(migrator *mysql.Migrator) error {
	logger.Info("开始验证GORM模型...")

	if err := migrator.ValidateModels(); err != nil {
		logger.Error("模型验证失败", zap.Error(err))
		return err
	}

	logger.Info("✅ 所有模型验证通过")
	return nil
}

func syncModels(migrator *mysql.Migrator, isDevelopment, force bool) error {
	if !isDevelopment && !force {
		logger.Error("非开发环境不允许自动同步模型，使用 -force 参数强制执行")
		return fmt.Errorf("production environment sync not allowed")
	}

	if !isDevelopment && force {
		logger.Warn("⚠️  强制在非开发环境同步模型，请确保你知道自己在做什么！")
		time.Sleep(3 * time.Second)
	}

	logger.Info("开始同步GORM模型到数据库...")

	if err := migrator.SyncModels(isDevelopment || force); err != nil {
		logger.Error("模型同步失败", zap.Error(err))
		return err
	}

	logger.Info("✅ 模型同步完成")
	return nil
}

func generateMigration(migrator *mysql.Migrator, modelName string) error {
	logger.Info("生成迁移脚本功能待实现", zap.String("model", modelName))
	// TODO: 实现从GORM模型生成SQL迁移脚本的功能
	return nil
}

func checkStatus(migrator *mysql.Migrator) error {
	logger.Info("检查迁移状态...")

	if err := migrator.CheckMigrationStatus(); err != nil {
		logger.Error("状态检查失败", zap.Error(err))
		return err
	}

	logger.Info("✅ 迁移状态检查完成")
	return nil
}
