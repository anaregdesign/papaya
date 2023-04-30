package graph

import (
	"context"
	"github.com/anaregdesign/papaya/cache"
	"sync"
	"time"
)

type Graph[S comparable, T any] struct {
	mu         sync.RWMutex
	defaultTTL time.Duration
	vertices   *cache.Cache[S, T]
	edges      *edgeCache[S]
}

func NewGraph[S comparable, T any](ctx context.Context, defaultTTL time.Duration) *Graph[S, T] {
	g := &Graph[S, T]{
		defaultTTL: defaultTTL,
		vertices:   cache.NewCache[S, T](ctx, defaultTTL),
		edges:      newEdgeCache[S](ctx, defaultTTL),
	}
	go g.watch(ctx, time.Minute)
	return g
}

func (g *Graph[S, T]) GetVertex(key S) (T, bool) {
	return g.vertices.Get(key)
}

func (g *Graph[S, T]) getWeight(tail, head S) float64 {
	return g.edges.get(tail, head)
}

func (g *Graph[S, T]) AddVertexWithTTL(key S, value T, ttl time.Duration) {
	g.vertices.SetWithTTL(key, value, ttl)
}

func (g *Graph[S, T]) AddVertex(key S, value T) {
	g.AddVertexWithTTL(key, value, g.defaultTTL)
}

func (g *Graph[S, T]) AddEdgeWithTTL(tail, head S, w float64, ttl time.Duration) {
	g.edges.setWithTTL(tail, head, w, ttl)
}

func (g *Graph[S, T]) AddEdge(tail, head S, w float64) {
	g.AddEdgeWithTTL(tail, head, w, g.defaultTTL)
}

func (g *Graph[S, T]) flush() {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for tail, heads := range g.edges.cache {
		if !g.vertices.Has(tail) {
			g.mu.Lock()
			delete(g.edges.cache, tail)
			g.mu.Unlock()
			continue
		}
		for head, _ := range heads {
			if !g.vertices.Has(head) {
				g.edges.delete(tail, head)
			}
		}
	}
}

func (g *Graph[S, T]) watch(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.flush()

		case <-ctx.Done():
			return
		}
	}
}
