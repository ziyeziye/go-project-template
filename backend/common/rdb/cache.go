package rdb

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/go-redis/cache/v9"
)

type RdsCache struct {
	Engine *cache.Cache
	Ctx    context.Context
}

var _cache *RdsCache

func Cache() *RdsCache {
	if _cache != nil {
		return _cache
	}

	_cache = &RdsCache{
		Engine: cache.New(&cache.Options{
			Redis: Engine(),
			// LocalCache: cache.NewTinyLFU(1000, time.Minute), // 优先使用本地缓存
		}),
		Ctx: context.Background(),
	}

	return _cache
}

func (c *RdsCache) Set(key string, value interface{}, ttl ...time.Duration) error {
	var t time.Duration
	t = 0
	if len(ttl) > 0 {
		t = ttl[0]
	}

	return c.Engine.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   t,
	})
}

func (c *RdsCache) Get(key string, value interface{}) error {
	return c.Engine.Get(c.Ctx, key, value)
}

func (c *RdsCache) Del(key string) error {
	return c.Engine.Delete(c.Ctx, key)
}

func (c *RdsCache) Exists(key string) bool {
	return c.Engine.Exists(c.Ctx, key)
}

func Keys(keys ...string) string {
	return strings.Join(keys, "_")
}

// SGet函数用于从缓存中获取数据，如果缓存中没有数据，则调用callback函数获取数据，并将数据存入缓存中
func SGet[T any](callback func() (T, error), lifeSpan time.Duration, keys ...string) (item T, err error) {
	// 生成缓存键
	key := Keys(keys...)
	// 从缓存中获取数据
	err = Cache().Get(key, &item)
	// 如果缓存中有数据，则直接返回
	if err == nil {
		return item, nil
	}
	// 调用callback函数获取数据
	item, err = callback()

	// 如果获取的数据为空，则直接返回
	if isZeroRef(item) {
		return item, nil
	}

	// 如果callback函数返回错误，则直接返回
	if err != nil {
		return item, err
	}

	// 将数据存入缓存中
	err = Cache().Set(key, item, lifeSpan)
	return item, err
}

func MustGet[T any](key string) T {
	var item T
	err := Cache().Get(key, &item)
	if err != nil {
		return *new(T)
	}
	return item
}

func isZeroRef[T any](v T) bool {
	return reflect.ValueOf(&v).Elem().IsZero()
}
