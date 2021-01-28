package infra

import (
	"context"
	"time"
)

// KV interface represents the key value storage
type KV interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, tab, key string, val interface{}) error
	List(ctx context.Context, tab string, keys []string, val interface{}) error
	Put(ctx context.Context, tab, key string, val interface{}) error
	PutEx(ctx context.Context, tab, key string, val interface{}, d time.Duration) error
	PutNxEx(ctx context.Context, tab, key string, val interface{}, d time.Duration) error
	PutNx(ctx context.Context, tab, key string, val interface{}) error
	Inc(ctx context.Context, tab, key string, val int) error
	SetEx(ctx context.Context, tab, key string, exp time.Duration) error
	GetEx(ctx context.Context, tab, key string) (time.Duration, error)
	Del(ctx context.Context, tab, key string) error
	DeleteByPattern(ctx context.Context, tab, pattern string) error
	AddToSet(ctx context.Context, tab, key string, members ...Member) error
	RemoveFromSet(ctx context.Context, tab, key string, members ...string) error
	RemoveFromAllSet(ctx context.Context, tab string, members ...string) error
	ListDataFromSet(ctx context.Context, tab, key string, skip, limit int64, data interface{}, asc bool) error
	ListDataFromSetWithRetrival(ctx context.Context, listTab, retrivalTab, key string, skip, limit int64, data interface{}, asc bool) error
}

type Member struct {
	Score int64
	Val   interface{}
}
