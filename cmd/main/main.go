package main

import (
	"log"

	"github.com/camwhite18/ArachneGraph/pkg/graph"
)

func main() {
	persister := &graph.BoltPersist{Path: "data.db"}

	g := graph.NewGraph()
	g.AddNode(&graph.Node{ID: "A"})
	g.AddNode(&graph.Node{ID: "B"})
	g.AddEdge(&graph.Edge{ID: "AB", From: "A", To: "B", Weight: 1})

	if err := persister.Save(g); err != nil {
		log.Fatalf("could not persist graph: %v", err)
	}
	log.Print("persisted graph")

	g2, err := persister.Load()
	if err != nil {
		log.Fatalf("could not load persisted graph: %v", err)
	}
	log.Printf("loaded graph with nodes: %v", g2.Nodes())
}
