package graph

import (
	"context"
	"github.com/anaregdesign/papaya/cache"
	"github.com/anaregdesign/papaya/collection/pq"
	"github.com/anaregdesign/papaya/collection/set"
	"math"
	"sync"
	"time"
)

type GraphCache[S comparable, T any] struct {
	mu         sync.RWMutex
	defaultTTL time.Duration
	vertices   *cache.Cache[S, T]
	edges      *edgeCache[S]
}

func NewGraphCache[S comparable, T any](ctx context.Context, defaultTTL time.Duration) *GraphCache[S, T] {
	g := &GraphCache[S, T]{
		defaultTTL: defaultTTL,
		vertices:   cache.NewCache[S, T](ctx, defaultTTL),
		edges:      newEdgeCache[S](ctx, defaultTTL),
	}
	go g.watch(ctx, time.Minute)
	return g
}

func (c *GraphCache[S, T]) GetVertex(key S) (T, bool) {
	return c.vertices.Get(key)
}

func (c *GraphCache[S, T]) getWeight(tail, head S) float64 {
	return c.edges.get(tail, head)
}

func (c *GraphCache[S, T]) AddVertexWithTTL(key S, value T, ttl time.Duration) {
	c.vertices.SetWithTTL(key, value, ttl)
}

func (c *GraphCache[S, T]) AddVertex(key S, value T) {
	c.AddVertexWithTTL(key, value, c.defaultTTL)
}

func (c *GraphCache[S, T]) AddEdgeWithTTL(tail, head S, w float64, ttl time.Duration) {
	if !c.vertices.Has(tail) {
		var noop T
		c.AddVertexWithTTL(tail, noop, ttl)
	}
	if !c.vertices.Has(head) {
		var noop T
		c.AddVertexWithTTL(head, noop, ttl)
	}
	c.edges.setWithTTL(tail, head, w, ttl)
}

func (c *GraphCache[S, T]) AddEdge(tail, head S, w float64) {
	c.AddEdgeWithTTL(tail, head, w, c.defaultTTL)
}

func (c *GraphCache[S, T]) flush() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for tail, heads := range c.edges.tf {
		if !c.vertices.Has(tail) {
			c.mu.Lock()
			delete(c.edges.tf, tail)
			c.mu.Unlock()
			continue
		}
		for head := range heads {
			if !c.vertices.Has(head) {
				c.edges.delete(tail, head)
			}
		}
	}
}

func (c *GraphCache[S, T]) Neighbor(seed S, step int, top int, tfidf bool) *Graph[S, T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	g := NewGraph[S, T]()

	if v, ok := c.vertices.Get(seed); !ok {
		return g
	} else {
		g.Vertices[seed] = v
	}

	targets := set.NewSet[S]()
	targets.Add(seed)
	seen := set.NewSet[S]()
	for i := 0; i < step; i++ {
		for _, tail := range targets.Values() {
			// Skip if already seen
			if seen.Has(tail) {
				continue
			}

			// Add edges to the graph
			edges := make(map[S]float64)
			for head, w := range c.edges.tf[tail] {
				if tfidf {
					edges[head] = w.value() / math.Log2(float64(1+c.edges.df[head]))
				} else {
					edges[head] = w.value()
				}
			}

			// Filter light edges
			if len(c.edges.tf[tail]) > top {
				g.Edges[tail] = pq.FilterMap(edges, top)
			}

			// Add new vertices to targets
			for head := range g.Edges[tail] {
				g.Vertices[head], _ = c.vertices.Get(head)
				targets.Add(head)
			}

			// Mark as seen
			seen.Add(tail)
		}
	}

	return g
}

func (c *GraphCache[S, T]) watch(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.flush()

		case <-ctx.Done():
			return
		}
	}
}
