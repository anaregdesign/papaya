package cache

import (
	"context"
	"github.com/anaregdesign/papaya/model"
	"time"
)

type LoadingCache[S comparable, T any] struct {
	cache  *Cache[S, T]
	loader model.Loader[S, T]
}

func NewLoadingCache[S comparable, T any](ctx context.Context, loader model.Loader[S, T], defaultTTL time.Duration) *LoadingCache[S, T] {
	return &LoadingCache[S, T]{
		cache:  NewCache[S, T](ctx, defaultTTL),
		loader: loader,
	}
}

func (c *LoadingCache[S, T]) Get(key S) (T, bool) {
	if value, ok := c.cache.Get(key); ok {
		return value, true
	}
	if value, ok := c.loader(key); ok {
		c.cache.Set(key, value)
		return value, true
	}
	var noop T
	return noop, false
}

func (c *LoadingCache[S, T]) Set(key S, value T) {
	c.cache.Set(key, value)
}

func (c *LoadingCache[S, T]) SetWithTTL(key S, value T, ttl time.Duration) {
	c.cache.SetWithTTL(key, value, ttl)
}
