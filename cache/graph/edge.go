package graph

import (
	"context"
	"sync"
	"time"
)

type weightValue struct {
	value      float64
	expiration time.Time
}

func (w weightValue) expired() bool {
	return time.Now().After(w.expiration)
}

type weight struct {
	values []weightValue
}

func newWeight() *weight {
	return &weight{
		values: make([]weightValue, 0),
	}
}

func (w *weight) value() float64 {
	w.flush()
	var sum float64
	for _, v := range w.values {
		sum += v.value
	}
	return sum
}

func (w *weight) add(value float64, ttl time.Duration) {
	w.values = append(w.values, weightValue{
		value:      value,
		expiration: time.Now().Add(ttl),
	})
}

func (w *weight) isZero() bool {
	return w.value() == 0
}

func (w *weight) flush() {
	v := make([]weightValue, 0)
	for _, value := range w.values {
		if !value.expired() {
			v = append(v, value)
		}
	}
	w.values = v
}

type edgeCache[S comparable] struct {
	mu         sync.RWMutex
	defaultTTL time.Duration
	cache      map[S]map[S]*weight
}

func newEdgeCache[S comparable](ctx context.Context, defaultTTL time.Duration) *edgeCache[S] {
	c := &edgeCache[S]{
		defaultTTL: defaultTTL,
		cache:      make(map[S]map[S]*weight),
	}
	go c.watch(ctx, time.Minute)
	return c
}

func (c *edgeCache[S]) get(tail, head S) float64 {
	if _, ok := c.cache[tail]; !ok {
		return 0
	}

	if _, ok := c.cache[tail][head]; !ok {
		return 0
	}

	if w := c.cache[tail][head]; w.isZero() {
		go c.delete(tail, head)
		return 0
	} else {
		return w.value()
	}
}

func (c *edgeCache[S]) setWithTTL(tail, head S, w float64, ttl time.Duration) {
	if _, ok := c.cache[tail]; !ok {
		c.cache[tail] = make(map[S]*weight)
	}

	if _, ok := c.cache[tail][head]; !ok {
		c.cache[tail][head] = newWeight()
	}

	c.cache[tail][head].add(w, ttl)
}

func (c *edgeCache[S]) set(tail, head S, w float64) {
	c.setWithTTL(tail, head, w, c.defaultTTL)
}

func (c *edgeCache[S]) delete(tail, head S) {
	if _, ok := c.cache[tail]; !ok {
		return
	}

	if _, ok := c.cache[tail][head]; !ok {
		return
	}

	delete(c.cache[tail], head)
	if len(c.cache[tail]) == 0 {
		delete(c.cache, tail)
	}
}

func (c *edgeCache[S]) flush() {
	for tail, heads := range c.cache {
		for head, w := range heads {
			if w.isZero() {
				delete(c.cache[tail], head)
				if len(c.cache[tail]) == 0 {
					delete(c.cache, tail)
				}
			}
		}
	}
}

func (c *edgeCache[S]) watch(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			c.flush()

		case <-ctx.Done():
			return

		}
	}
}
