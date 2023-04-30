package graph

type Graph[S comparable, T any] struct {
	Vertices map[S]T
	Edges    map[S]map[S]float64
}
