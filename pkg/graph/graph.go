package graph

import "sync"

type Graph interface {
	AddNode(node *Node) error
	GetNode(id string) (*Node, error)
	DeleteNode(id string) error
	AddEdge(edge *Edge) error
	GetEdge(from, to string) (*Edge, error)
	DeleteEdge(from, to string) error
	Neighbors(id string) ([]*Edge, error)
	Incoming(id string) ([]*Edge, error)
}

type graphImpl struct {
	nodes map[string]*Node
	out   map[string]map[string]*Edge
	in    map[string]map[string]*Edge
	mu    sync.RWMutex
}

func NewGraph() Graph {
	return &graphImpl{
		nodes: make(map[string]*Node),
		out:   make(map[string]map[string]*Edge),
		in:    make(map[string]map[string]*Edge),
		mu:    sync.RWMutex{},
	}
}
