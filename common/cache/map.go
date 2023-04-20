package cache

import (
	"context"
	"sync"
	"time"
)

type volatile[T any] struct {
	value T
	ttl   int64
}

func (v *volatile[T]) IsExpired() bool {
	return v.ttl <= time.Now().Unix()
}

type Cache[S comparable, T any] struct {
	defaultTTL time.Duration
	cache      map[S]volatile[T]
	mu         sync.RWMutex
}

func NewCache[S comparable, T any](ctx context.Context, defaultTTL time.Duration) *Cache[S, T] {
	cache := &Cache[S, T]{
		defaultTTL: defaultTTL,
		cache:      make(map[S]volatile[T]),
	}
	go cache.FlushInterval(ctx, time.Minute)
	return cache
}

func (c *Cache[S, T]) Get(key S) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if v, ok := c.cache[key]; ok {
		if v.IsExpired() {
			var noop T
			return noop, false
		}
		return v.value, true
	}
	var noop T
	return noop, false
}

func (c *Cache[S, T]) Set(key S, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = volatile[T]{
		value: value,
		ttl:   time.Now().Add(c.defaultTTL).Unix(),
	}
}

func (c *Cache[S, T]) SetWithTTL(key S, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = volatile[T]{
		value: value,
		ttl:   time.Now().Add(ttl).Unix(),
	}
}

func (c *Cache[S, T]) Delete(key S) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cache[key]; ok {
		delete(c.cache, key)
	}
}
func (c *Cache[S, T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[S]volatile[T])
}

func (c *Cache[S, T]) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}

func (c *Cache[S, T]) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.cache {
		if v.IsExpired() {
			delete(c.cache, k)
		}
	}
}

func (c *Cache[S, T]) FlushInterval(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.Flush()
		case <-ctx.Done():
			return
		}
	}
}
