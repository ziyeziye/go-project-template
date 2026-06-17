package cachex

import (
	"reflect"
	"strings"
	"time"

	"github.com/muesli/cache2go"
)

var engine *cache2go.CacheTable

func Engine() *cache2go.CacheTable {
	if engine == nil {
		engine = cache2go.Cache("CacheGo")
	}
	return engine
}

func New(table string) *cache2go.CacheTable {
	return cache2go.Cache(table)
}

func GetOrAdd[T any](callback func() (T, error), lifeSpan time.Duration, keys ...string) (item T, err error) {
	key := strings.Join(keys, "_")
	res, err := Engine().Value(key)
	if err == nil {
		return res.Data().(T), nil
	}
	item, err = callback()

	if isZeroRef(item) {
		return item, nil
	}

	if err != nil {
		return item, err
	}

	add := Engine().Add(key, lifeSpan, item)
	return add.Data().(T), nil
}

func MustGet[T any](key string) T {
	res, err := Engine().Value(key)
	if err != nil {
		return *new(T)
	}
	return res.Data().(T)
}

func isZeroRef[T any](v T) bool {
	return reflect.ValueOf(&v).Elem().IsZero()
}
