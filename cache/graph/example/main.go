package main

import (
	"context"
	"encoding/json"
	"github.com/anaregdesign/papaya/cache/graph"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	c := graph.NewGraphCache[string, string](ctx, 1*time.Minute)

	c.AddEdge("a", "b", 1)
	c.AddEdge("b", "c", 1)
	c.AddEdge("c", "d", 1)
	c.AddEdge("a", "b", 1)
	c.AddEdge("a", "c", 1)
	c.AddEdge("a", "d", 1)
	c.AddEdge("a", "e", 1)

	g := c.Neighbor("a", 5, 3, true)

	if jsonText, err := json.MarshalIndent(g, "", "\t"); err == nil {
		println(string(jsonText))
	}

	/*
		Output:
		{
		        "vertices": {
		                "a": "",
		                "b": "",
		                "c": "",
		                "d": "",
		                "e": ""
		        },
		        "edges": {
		                "a": {
		                        "b": 0.6309297535714574
		                },
		                "b": {
		                        "c": 1
		                },
		                "c": {
		                        "d": 1,
		                        "e": 1
		                }
		        }
		}

	*/
}
