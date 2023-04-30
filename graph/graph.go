package graph

import (
	"context"
	"github.com/anaregdesign/papaya/cache"
	"time"
)

type Graph[S comparable, T any] struct {
	defaultTTL time.Duration
	vertices   *cache.Cache[S, T]
	edges      map[S]map[S]*weight
}

func NewGraph[S comparable, T any](ctx context.Context, defaultTTL time.Duration) *Graph[S, T] {
	return &Graph[S, T]{
		defaultTTL: defaultTTL,
		vertices:   cache.NewCache[S, T](ctx, defaultTTL),
		edges:      make(map[S]map[S]*weight),
	}
}

func (g *Graph[S, T]) GetVertex(key S) (T, bool) {
	return g.vertices.Get(key)
}

func (g *Graph[S, T]) getWeight(tail, head S) float64 {
	if _, ok := g.edges[tail]; !ok {
		return 0
	}

	if _, ok := g.edges[tail][head]; !ok {
		return 0
	}

	if g.edges[tail][head].isZero() {
		delete(g.edges[tail], head)
		return 0
	}

	return g.edges[tail][head].value()
}

func (g *Graph[S, T]) AddVertexWithTTL(key S, value T, ttl time.Duration) {
	g.vertices.SetWithTTL(key, value, ttl)
}

func (g *Graph[S, T]) AddVertex(key S, value T) {
	g.AddVertexWithTTL(key, value, g.defaultTTL)
}

func (g *Graph[S, T]) AddEdgeWithTTL(tail, head S, w float64, ttl time.Duration) {
	if _, ok := g.edges[tail]; !ok {
		g.edges[tail] = make(map[S]*weight)
	}

	if _, ok := g.edges[tail][head]; !ok {
		g.edges[tail][head] = newWeight()
	}

	g.edges[tail][head].add(w, ttl)
}

func (g *Graph[S, T]) AddEdge(tail, head S, w float64) {
	g.AddEdgeWithTTL(tail, head, w, g.defaultTTL)
}
