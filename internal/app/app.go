package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/taskflow/internal/application/shared"
	"github.com/taskflow/internal/application/user"
	"github.com/taskflow/internal/infrastructure/auth"
	"github.com/taskflow/internal/infrastructure/config"
	"github.com/taskflow/internal/infrastructure/persistence/mysql"
	httpServer "github.com/taskflow/internal/interfaces/http"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用程序结构
type App struct {
	config         *config.Config
	db             *gorm.DB
	httpServer     *httpServer.Server
	transactionMgr shared.TransactionManager
	jwtService     shared.JWTService
	userAppService *user.UserAppService
}

// NewApp 创建新的应用程序实例
func NewApp(configPath string) (*App, error) {
	// 1. 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. 初始化日志
	if err := logger.InitLogger(&logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		FilePath:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger.Info("Application initializing...",
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("mode", cfg.App.Mode))

	// 3. 连接数据库
	db, err := mysql.NewDatabase(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 4. 验证数据库模型一致性
	migrator := mysql.NewMigrator(db)
	if err := migrator.CheckMigrationStatus(); err != nil {
		logger.Warn("Migration status check failed", zap.Error(err))
	}

	if err := migrator.ValidateModels(); err != nil {
		if cfg.App.Mode == "development" {
			logger.Warn("Model validation failed in development mode", zap.Error(err))
			// 开发环境下可以选择自动同步（谨慎使用）
			// if err := migrator.SyncModels(true); err != nil {
			//     logger.Error("Failed to sync models", zap.Error(err))
			// }
		} else {
			logger.Error("Model validation failed in production", zap.Error(err))
			return nil, fmt.Errorf("database model validation failed: %w", err)
		}
	}

	// 5. 创建事务管理器
	transactionMgr := mysql.NewTransactionManager(db)

	// 6. 创建JWT服务
	jwtService := auth.NewJWTService(shared.JWTConfig{
		Secret:             cfg.JWT.Secret,
		AccessTokenExpiry:  time.Duration(cfg.JWT.ExpireHours) * time.Hour,
		RefreshTokenExpiry: time.Duration(cfg.JWT.RefreshExpireHours) * time.Hour,
		Issuer:             cfg.App.Name,
	})

	// 7. 创建仓储层
	userRepo := mysql.NewUserRepository(db)

	// 8. 创建应用服务层
	userAppService := user.NewUserAppService(userRepo, transactionMgr)

	// 9. 创建HTTP服务器
	httpSrv := httpServer.NewServer(cfg, jwtService, userAppService)

	app := &App{
		config:         cfg,
		db:             db,
		httpServer:     httpSrv,
		transactionMgr: transactionMgr,
		jwtService:     jwtService,
		userAppService: userAppService,
	}

	return app, nil
}

// Run 运行应用程序
func (a *App) Run() error {
	logger.Info("Starting TaskFlow application...")

	// 启动HTTP服务器
	go func() {
		if err := a.httpServer.Start(); err != nil {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// 等待中断信号
	return a.gracefulShutdown()
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down application...")

	// 创建关闭上下文，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := a.httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	// 关闭数据库连接
	if err := a.closeDatabase(); err != nil {
		logger.Error("Database shutdown error", zap.Error(err))
	}

	logger.Info("Application shutdown complete")
	return nil
}

// closeDatabase 关闭数据库连接
func (a *App) closeDatabase() error {
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB 获取数据库实例（用于依赖注入）
func (a *App) GetDB() *gorm.DB {
	return a.db
}

// GetConfig 获取配置实例
func (a *App) GetConfig() *config.Config {
	return a.config
}
