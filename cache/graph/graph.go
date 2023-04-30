package graph

import (
	"context"
	"github.com/anaregdesign/papaya/cache"
	"time"
)

type Graph[S comparable, T any] struct {
	defaultTTL time.Duration
	vertices   *cache.Cache[S, T]
	edges      *edgeCache[S]
}

func NewGraph[S comparable, T any](ctx context.Context, defaultTTL time.Duration) *Graph[S, T] {
	return &Graph[S, T]{
		defaultTTL: defaultTTL,
		vertices:   cache.NewCache[S, T](ctx, defaultTTL),
		edges:      newEdgeCache[S](ctx, defaultTTL),
	}
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
