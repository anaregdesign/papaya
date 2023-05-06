package main

import (
	"encoding/json"
	"github.com/anaregdesign/papaya/graph"
	"log"
)

func main() {
	g := graph.NewGraph[string, string]()
	g.AddVertex("a", "A")
	g.AddVertex("b", "B")
	g.AddVertex("c", "C")
	g.AddVertex("d", "D")
	g.AddEdge("d", "d", 1)
	g.AddEdge("a", "b", 1)
	g.AddEdge("a", "c", 2)
	g.AddEdge("b", "a", 3)
	g.AddEdge("b", "c", 4)
	g.AddEdge("c", "a", 100)
	g.AddEdge("c", "b", 6)

	log.Println("Original graph")
	if jsonText, err := json.MarshalIndent(g, "", "\t"); err == nil {
		log.Println(string(jsonText))
	}
	/*
		{
		        "vertices": {
		                "a": "A",
		                "b": "B",
		                "c": "C",
		                "d": "D"
		        },
		        "edges": {
		                "a": {
		                        "b": 1,
		                        "c": 2
		                },
		                "b": {
		                        "a": 3,
		                        "c": 4
		                },
		                "c": {
		                        "a": 100,
		                        "b": 6
		                },
		                "d": {
		                        "d": 1
		                }
		        }
		}
	*/

	log.Println("Connected graph from vertex 'a'")
	connected := g.ConnectedGraph("a")
	if jsonText, err := json.MarshalIndent(connected, "", "\t"); err == nil {
		log.Println(string(jsonText))
	}

	/*
		 {
		        "vertices": {
		                "a": "A",
		                "b": "B",
		                "c": "C"
		        },
		        "edges": {
		                "a": {
		                        "b": 1,
		                        "c": 2
		                },
		                "b": {
		                        "a": 3,
		                        "c": 4
		                },
		                "c": {
		                        "a": 100,
		                        "b": 6
		                }
		        }
		}
	*/

	log.Println("Minimum Spanning Tree from vertex 'a'")
	mst := g.MinimumSpanningTree("a", false)
	if jsonText, err := json.MarshalIndent(mst, "", "\t"); err == nil {
		log.Println(string(jsonText))
	}

	/*
		{
		        "vertices": {
		                "a": "A",
		                "b": "B",
		                "c": "C"
		        },
		        "edges": {
		                "a": {
		                        "c": 2
		                },
		                "c": {
		                        "b": 6
		                }
		        }
		}

	*/

}
