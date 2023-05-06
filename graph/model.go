package graph

import (
	"container/heap"
	"github.com/anaregdesign/papaya/collection/pq"
	"github.com/anaregdesign/papaya/collection/set"
)

type Graph[S comparable, T any] struct {
	Vertices map[S]T             `json:"vertices,omitempty"`
	Edges    map[S]map[S]float32 `json:"edges,omitempty"`
}

func NewGraph[S comparable, T any]() *Graph[S, T] {
	return &Graph[S, T]{
		Vertices: make(map[S]T),
		Edges:    make(map[S]map[S]float32),
	}
}

func (g *Graph[S, T]) AddVertex(key S, value T) {
	g.Vertices[key] = value
}

func (g *Graph[S, T]) AddEdge(tail, head S, value float32) {
	if _, ok := g.Edges[tail]; !ok {
		g.Edges[tail] = make(map[S]float32)
	}
	g.Edges[tail][head] = value
}

func (g *Graph[S, T]) MinimumSpanningTree(negate bool) *Graph[S, T] {

	type edge struct {
		tail   S
		head   S
		weight float32
	}

	mst := NewGraph[S, T]()
	q := make(pq.PriorityQueue[S, float32], len(g.Vertices))
	seen := set.NewSet[S]()

	// get a seed vertex
	var seed S
	for k := range g.Vertices {
		seed = k
		break
	}
	mst.AddVertex(seed, g.Vertices[seed])
	for {
		if len(mst.Vertices) == len(g.Vertices) {
			break
		}

		for tail := range mst.Vertices {
			if seen.Has(tail) {
				continue
			}

			for head, weight := range g.Edges[tail] {
				var w float32
				if negate {
					w = -weight
				} else {
					w = weight
				}

				heap.Push(&q, &pq.Item[edge, float32]{
					Value: edge{
						tail:   tail,
						head:   head,
						weight: w,
					},
					Priority: w,
				})
			}
			seen.Add(tail)
		}
		pickedUp := heap.Pop(&q).(*pq.Item[edge, float32]).Value
		mst.AddVertex(pickedUp.tail, g.Vertices[pickedUp.tail])
		mst.AddEdge(pickedUp.tail, pickedUp.head, pickedUp.weight)
	}
	return mst
}

func (g *Graph[S, T]) Render(key2int func(k S) int, value2string func(v T) string) GraphView {
	var vertices []VertexView
	var edges []EdgeView

	for i, v := range g.Vertices {
		vertices = append(vertices, VertexView{
			ID:    key2int(i),
			Label: value2string(v),
		})
	}

	for from, tos := range g.Edges {
		for to, value := range tos {
			edges = append(edges, EdgeView{
				From:  key2int(from),
				To:    key2int(to),
				Value: value,
			})
		}
	}

	return GraphView{
		Vertices: vertices,
		Edges:    edges,
	}
}
