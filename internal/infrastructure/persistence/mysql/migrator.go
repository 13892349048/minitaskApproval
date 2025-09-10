package mysql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Migrator 数据库迁移管理器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移管理器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// ValidateModels 验证GORM模型与数据库结构是否一致
func (m *Migrator) ValidateModels() error {
	logger.Info("开始验证GORM模型与数据库结构...")

	models := []interface{}{
		&User{}, &Role{}, &Permission{}, &UserRole{}, &PermissionPolicy{},
		&Project{}, &ProjectMember{},
		&Task{}, &TaskParticipant{}, &RecurrenceRule{}, &TaskExecution{}, &ParticipantCompletion{},
		&ApprovalRecord{}, &ExtensionRequest{},
		&DomainEvent{}, &OperationLog{},
		&File{}, &FileAssociation{},
	}

	var errors []string

	for _, model := range models {
		if err := m.validateModel(model); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		logger.Error("GORM模型验证失败", zap.Strings("errors", errors))
		return fmt.Errorf("模型验证失败: %s", strings.Join(errors, "; "))
	}

	logger.Info("GORM模型验证通过")
	return nil
}

// validateModel 验证单个模型
func (m *Migrator) validateModel(model interface{}) error {
	modelType := reflect.TypeOf(model).Elem()
	tableName := m.getTableName(model)

	// 检查表是否存在
	if !m.db.Migrator().HasTable(model) {
		return fmt.Errorf("表 %s 不存在", tableName)
	}

	// 检查字段
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// 跳过关联字段和特殊字段
		if m.shouldSkipField(field) {
			continue
		}

		columnName := m.getColumnName(field)
		if !m.db.Migrator().HasColumn(model, columnName) {
			return fmt.Errorf("表 %s 缺少字段 %s", tableName, columnName)
		}
	}

	logger.Debug("模型验证通过", zap.String("table", tableName))
	return nil
}

// getTableName 获取表名
func (m *Migrator) getTableName(model interface{}) string {
	if tabler, ok := model.(interface{ TableName() string }); ok {
		return tabler.TableName()
	}
	return m.db.NamingStrategy.TableName(reflect.TypeOf(model).Elem().Name())
}

// getColumnName 获取字段名
func (m *Migrator) getColumnName(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")
	if gormTag != "" {
		// 解析gorm标签中的column名
		parts := strings.Split(gormTag, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}

	// 使用GORM的命名策略
	return m.db.NamingStrategy.ColumnName("", field.Name)
}

// shouldSkipField 判断是否跳过字段验证
func (m *Migrator) shouldSkipField(field reflect.StructField) bool {
	gormTag := field.Tag.Get("gorm")

	// 跳过关联字段
	if strings.Contains(gormTag, "foreignKey") ||
		strings.Contains(gormTag, "many2many") ||
		strings.Contains(gormTag, "hasMany") ||
		strings.Contains(gormTag, "hasOne") ||
		strings.Contains(gormTag, "belongsTo") {
		return true
	}

	// 跳过忽略的字段
	if strings.Contains(gormTag, "-") {
		return true
	}

	// 跳过切片类型（通常是关联）
	if field.Type.Kind() == reflect.Slice {
		return true
	}

	// 跳过指针类型的结构体（通常是关联）
	if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
		// 除非是时间类型
		if field.Type.Elem().Name() != "Time" {
			return true
		}
	}

	return false
}

// SyncModels 同步模型到数据库（仅开发环境）
func (m *Migrator) SyncModels(isDevelopment bool) error {
	if !isDevelopment {
		return fmt.Errorf("生产环境不允许自动同步模型")
	}

	logger.Warn("开发环境：正在同步GORM模型到数据库...")

	models := []interface{}{
		&User{}, &Role{}, &Permission{}, &UserRole{}, &PermissionPolicy{},
		&Project{}, &ProjectMember{},
		&Task{}, &TaskParticipant{}, &RecurrenceRule{}, &TaskExecution{}, &ParticipantCompletion{},
		&ApprovalRecord{}, &ExtensionRequest{},
		&DomainEvent{}, &OperationLog{},
		&File{}, &FileAssociation{},
	}

	for _, model := range models {
		if err := m.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("同步模型失败 %T: %w", model, err)
		}
	}

	logger.Info("GORM模型同步完成")
	return nil
}

// CreateMigrationFromModel 从GORM模型生成SQL迁移脚本（开发工具）
func (m *Migrator) CreateMigrationFromModel(model interface{}, migrationName string) (string, error) {
	tableName := m.getTableName(model)
	modelType := reflect.TypeOf(model).Elem()

	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("-- 迁移: %s\n", migrationName))
	sql.WriteString(fmt.Sprintf("-- 生成时间: %s\n\n", "{{.Timestamp}}"))
	sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", tableName))

	var fields []string
	var indexes []string

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		if m.shouldSkipField(field) {
			continue
		}

		columnName := m.getColumnName(field)
		columnType := m.getColumnType(field)
		columnDef := fmt.Sprintf("    `%s` %s", columnName, columnType)

		// 添加注释
		if comment := field.Tag.Get("comment"); comment != "" {
			columnDef += fmt.Sprintf(" COMMENT '%s'", comment)
		}

		fields = append(fields, columnDef)

		// 检查索引
		if idx := field.Tag.Get("index"); idx != "" {
			indexes = append(indexes, fmt.Sprintf("    INDEX `idx_%s` (`%s`)", columnName, columnName))
		}
	}

	sql.WriteString(strings.Join(fields, ",\n"))

	if len(indexes) > 0 {
		sql.WriteString(",\n\n")
		sql.WriteString(strings.Join(indexes, ",\n"))
	}

	sql.WriteString("\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	sql.WriteString(fmt.Sprintf(" COMMENT='%s表';\n", tableName))

	return sql.String(), nil
}

// getColumnType 获取字段的SQL类型
func (m *Migrator) getColumnType(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")

	// 如果gorm标签中指定了类型，直接使用
	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "type:") {
			return strings.TrimPrefix(part, "type:")
		}
	}

	// 根据Go类型推断SQL类型
	switch field.Type.Kind() {
	case reflect.String:
		if strings.Contains(gormTag, "primaryKey") {
			return "VARCHAR(36) NOT NULL PRIMARY KEY"
		}
		return "VARCHAR(255)"
	case reflect.Int, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Bool:
		return "BOOLEAN DEFAULT FALSE"
	case reflect.Ptr:
		if field.Type.Elem().Name() == "Time" {
			return "TIMESTAMP NULL"
		}
		return "VARCHAR(255)"
	default:
		if field.Type.Name() == "Time" {
			if strings.Contains(gormTag, "autoCreateTime") || strings.Contains(gormTag, "autoUpdateTime") {
				return "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
			}
			return "TIMESTAMP"
		}
		return "TEXT"
	}
}

// CheckMigrationStatus 检查迁移状态
func (m *Migrator) CheckMigrationStatus() error {
	logger.Info("检查数据库迁移状态...")

	// 创建迁移状态表
	if err := m.createMigrationTable(); err != nil {
		return err
	}

	// 这里可以扩展检查已执行的迁移
	logger.Info("迁移状态检查完成")
	return nil
}

// createMigrationTable 创建迁移状态表
func (m *Migrator) createMigrationTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库迁移记录表';
	`

	return m.db.Exec(sql).Error
}
