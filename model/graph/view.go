package graph

type vertexView struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Value string `json:"value,omitempty"`
}

type EdgeView struct {
	From  int     `json:"from"`
	To    int     `json:"to"`
	Value float32 `json:"value,omitempty"`
}

type GraphView struct {
	Vertices []vertexView `json:"nodes"`
	Edges    []EdgeView   `json:"edges"`
}

func View(g Graph[int, string]) GraphView {
	var vertices []vertexView
	var edges []EdgeView

	for i, v := range g.Vertices {
		vertices = append(vertices, vertexView{
			ID:    i,
			Label: v,
		})
	}

	for from, tos := range g.Edges {
		for to, value := range tos {
			edges = append(edges, EdgeView{
				From:  from,
				To:    to,
				Value: value,
			})
		}
	}

	return GraphView{
		Vertices: vertices,
		Edges:    edges,
	}
}
