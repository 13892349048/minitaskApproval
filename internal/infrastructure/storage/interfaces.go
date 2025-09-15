package storage

import (
	"context"
	"database/sql"
	"time"
)

// DatabaseInterface 数据库接口抽象
type DatabaseInterface interface {
	// 事务管理
	BeginTx(ctx context.Context, opts *sql.TxOptions) (TransactionInterface, error)
	
	// 查询操作
	QueryContext(ctx context.Context, query string, args ...interface{}) (RowsInterface, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) RowInterface
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	
	// 连接管理
	Ping() error
	Close() error
	Stats() sql.DBStats
}

// TransactionInterface 事务接口
type TransactionInterface interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (RowsInterface, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) RowInterface
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

// RowsInterface 查询结果集接口
type RowsInterface interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// RowInterface 单行查询结果接口
type RowInterface interface {
	Scan(dest ...interface{}) error
}

// CacheInterface 缓存接口抽象
type CacheInterface interface {
	// 基本操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	
	// 批量操作
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	MSet(ctx context.Context, pairs ...interface{}) error
	
	// 过期管理
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// 连接管理
	Ping(ctx context.Context) error
	Close() error
}

// StorageConfig 存储配置
type StorageConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Cache    CacheConfig    `yaml:"cache"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `yaml:"driver"`
	DSN             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type     string        `yaml:"type"` // redis, memory, etc.
	Address  string        `yaml:"address"`
	Password string        `yaml:"password"`
	DB       int           `yaml:"db"`
	PoolSize int           `yaml:"pool_size"`
	TTL      time.Duration `yaml:"default_ttl"`
}
