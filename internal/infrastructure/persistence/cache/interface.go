package cache

import (
	"context"
	"time"
)

// Interface 缓存接口 - 在现有架构基础上添加缓存抽象
type Interface interface {
	// 基本操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
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
