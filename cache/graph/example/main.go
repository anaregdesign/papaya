package main

import (
	"context"
	"encoding/json"
	"github.com/anaregdesign/papaya/cache/graph"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()
	c := graph.NewGraphCache[int, int](1 * time.Minute)
	go c.Watch(ctx, 1*time.Second)

	//c.AddEdge("a", "b", 1)
	//c.AddEdge("b", "c", 1)
	//c.AddEdge("c", "d", 1)
	//c.AddEdge("a", "b", 1)
	//c.AddEdge("a", "c", 1)
	//c.AddEdge("a", "d", 1)
	//c.AddEdge("a", "e", 1)

	for i := 0; i < 100000; i++ {
		v1 := rand.Intn(10000)
		v2 := rand.Intn(10000)
		c.AddEdge(v1, v2, 1)
	}
	start := time.Now()
	g := c.Neighbor(1, 3, 3, true)
	end := time.Now()

	if jsonText, err := json.MarshalIndent(g, "", "\t"); err == nil {
		println(string(jsonText))
	}

	println("Time: ", end.Sub(start).String())

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
