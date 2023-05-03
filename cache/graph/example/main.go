package main

import (
	"context"
	"encoding/json"
	"github.com/anaregdesign/papaya/cache/graph"
	"time"
)

func main() {
	ctx := context.Background()
	c := graph.NewGraphCache[string, string](1 * time.Minute)
	go c.Watch(ctx, 1*time.Second)

	c.AddEdge("a", "b", 1)
	c.AddEdge("b", "c", 1)
	c.AddEdge("c", "d", 1)
	c.AddEdge("a", "b", 1)
	c.AddEdge("a", "c", 1)
	c.AddEdge("a", "d", 1)
	c.AddEdge("a", "e", 1)

	g := c.Neighbor("a", 2, 3, true)

	if jsonText, err := json.MarshalIndent(g, "", "\t"); err == nil {
		println(string(jsonText))
	}

	/*
		## Output:
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
		                        "b": 2,
		                        "d": 0.6309297535714574,
		                        "e": 1
		                },
		                "b": {
		                        "c": 0.6309297535714574
		                }
		        }
		}
	*/
}
