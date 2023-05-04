package graph

import view "github.com/anaregdesign/papaya/view/graph"

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

func (g *Graph[S, T]) Render(key2int func(k S) int, value2string func(v T) string) view.GraphView {
	var vertices []view.VertexView
	var edges []view.EdgeView

	for i, v := range g.Vertices {
		vertices = append(vertices, view.VertexView{
			ID:    key2int(i),
			Label: value2string(v),
		})
	}

	for from, tos := range g.Edges {
		for to, value := range tos {
			edges = append(edges, view.EdgeView{
				From:  key2int(from),
				To:    key2int(to),
				Value: value,
			})
		}
	}

	return view.GraphView{
		Vertices: vertices,
		Edges:    edges,
	}
}
