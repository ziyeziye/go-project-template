package rdb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var _rdb *redis.Client

func Engine() *redis.Client {
	if _rdb != nil {
		return _rdb
	}

	_rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("RDB_HOST"), viper.GetString("RDB_PORT")),
		Password: viper.GetString("RDB_PASSWORD"),
		DB:       0,
	})

	return _rdb
}

// RunWithLock 使用Redis锁执行函数，避免重复执行
// key: 锁的键名
// ttl: 锁的过期时间
// fn: 需要执行的函数
func RunWithLock(key string, ttl time.Duration, fn func()) bool {
	// 尝试获取锁
	if res, _ := Engine().Get(context.Background(), key).Result(); res == "" {
		// 设置锁
		Engine().Set(context.Background(), key, "1", ttl)

		// 异步执行函数
		go func() {
			// 使用 defer 处理清理和错误捕获
			defer func() {
				Engine().Del(context.Background(), key)
				recover() // 只捕获 panic，不做任何处理
			}()

			fn()
		}()
		return true
	}
	return false
}

// AcquireLock 尝试获取 Redis 锁并返回一个释放锁的函数
func AcquireLock(ctx context.Context, key string, duration time.Duration) (func(), error) {
	locked, err := Engine().SetNX(ctx, key, "locked", duration).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock for key %s: %v", key, err)
	}
	if !locked {
		return nil, fmt.Errorf("key %s is already locked", key)
	}

	// 返回一个函数用于释放锁
	return func() {
		Engine().Del(ctx, key)
	}, nil
}
