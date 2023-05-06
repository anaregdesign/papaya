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

func (g *Graph[S, T]) PutVertex(key S, value T) {
	g.Vertices[key] = value
}

func (g *Graph[S, T]) PutEdge(tail, head S, weight float32) {
	if _, ok := g.Vertices[tail]; !ok {
		var noop T
		g.Vertices[tail] = noop
	}

	if _, ok := g.Vertices[head]; !ok {
		var noop T
		g.Vertices[head] = noop
	}

	if _, ok := g.Edges[tail]; !ok {
		g.Edges[tail] = make(map[S]float32)
	}
	g.Edges[tail][head] = weight
}

func (g *Graph[S, T]) ConnectedGraph(seed S) *Graph[S, T] {
	targets := set.NewSet[S]()
	seen := set.NewSet[S]()
	connected := NewGraph[S, T]()
	connected.PutVertex(seed, g.Vertices[seed])

	targets.Add(seed)
	for {
		for _, tail := range targets.Values() {
			if seen.Has(tail) {
				continue
			}

			for head, weight := range g.Edges[tail] {
				connected.PutVertex(head, g.Vertices[head])
				connected.PutEdge(tail, head, weight)
			}
			seen.Add(tail)
		}
		for _, heads := range connected.Edges {
			for head := range heads {
				targets.Add(head)
			}
		}

		if targets.Size() == seen.Size() {
			break
		}
	}

	return connected
}

// MinimumSpanningTree
/*
 * MinimumSpanningTree returns a minimum spanning tree of the graph.
 * The seed is the starting point of the tree.
 * If negate is true, the maximum spanning tree is returned.
 */
func (g *Graph[S, T]) MinimumSpanningTree(seed S, negate bool) *Graph[S, T] {
	connected := g.ConnectedGraph(seed)

	type edge struct {
		tail   S
		head   S
		weight float32
	}

	mst := NewGraph[S, T]()
	q := make(pq.PriorityQueue[*edge, float32], 0)
	heap.Init(&q)
	seen := set.NewSet[S]()

	mst.PutVertex(seed, connected.Vertices[seed])
	for {
		if len(mst.Vertices) == len(connected.Vertices) {
			break
		}

		for tail := range mst.Vertices {
			if seen.Has(tail) {
				continue
			}

			for head, weight := range connected.Edges[tail] {
				var w float32
				if negate {
					w = weight
				} else {
					w = -weight
				}

				item := &pq.Item[*edge, float32]{
					Value:    &edge{tail, head, w},
					Priority: w,
				}
				heap.Push(&q, item)
			}
			seen.Add(tail)
		}

		var pickedUp *edge
		for {
			if q.Len() == 0 {
				break
			}
			pickedUp = heap.Pop(&q).(*pq.Item[*edge, float32]).Value
			if _, ok := mst.Vertices[pickedUp.head]; !ok {
				break
			}
		}
		if pickedUp == nil {
			continue
		}
		mst.PutVertex(pickedUp.head, connected.Vertices[pickedUp.head])
		mst.PutEdge(pickedUp.tail, pickedUp.head, connected.Edges[pickedUp.tail][pickedUp.head])

	}
	return mst
}

func (g *Graph[S, T]) ShortestPathTree(seed S, negate bool) *Graph[S, T] {
	connected := g.ConnectedGraph(seed)
	spt := NewGraph[S, T]()
	spt.PutVertex(seed, connected.Vertices[seed])

	type edge struct {
		tail   S
		head   S
		weight float32
	}

	seen := set.NewSet[S]()
	q := make(pq.PriorityQueue[*edge, float32], 0)
	heap.Init(&q)
	pivot := seed
	position := float32(0.0)
	for {
		if len(spt.Vertices) == len(connected.Vertices) {
			break
		}

		for head, weight := range connected.Edges[pivot] {
			if seen.Has(head) {
				continue
			}
			var w float32
			if negate {
				w = weight
			} else {
				w = -weight
			}

			heap.Push(&q, &pq.Item[*edge, float32]{
				Value: &edge{
					tail:   pivot,
					head:   head,
					weight: weight,
				},
				Priority: position + w,
			})
		}

		pickedUp := heap.Pop(&q).(*pq.Item[*edge, float32])
		spt.PutVertex(pickedUp.Value.head, connected.Vertices[pickedUp.Value.head])
		spt.PutEdge(pickedUp.Value.tail, pickedUp.Value.head, pickedUp.Value.weight)

		seen.Add(pivot)
		pivot = pickedUp.Value.head
		position = pickedUp.Priority
	}

	return spt
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
