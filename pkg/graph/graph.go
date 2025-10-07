package graph

import (
	"container/heap"
	"fmt"
	"math"
	"sync"

	pq "github.com/camwhite18/ArachneGraph/pkg/priorityQueue"
)

type Graph interface {
	AddNode(node *Node) error
	GetNode(id string) (*Node, error)
	DeleteNode(id string) error
	AddEdge(edge *Edge) error
	GetEdge(from, to string) (*Edge, error)
	DeleteEdge(from, to string) error
	Neighbors(id string) ([]*Edge, error)
	Incoming(id string) ([]*Edge, error)
	DFS(start string) ([]string, error)
	BFS(start string) ([]string, error)
	ShortestPath(source, target string) ([]string, float64, error)
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

func (g *graphImpl) AddNode(node *Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.nodes[node.ID]; exists {
		return fmt.Errorf("node %q already exists", node.ID)
	}
	g.nodes[node.ID] = node
	return nil
}

func (g *graphImpl) GetNode(id string) (*Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	node, exists := g.nodes[id]
	if !exists {
		return nil, fmt.Errorf("node %q not found", id)
	}
	return node, nil
}

func (g *graphImpl) DeleteNode(id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.nodes[id]; !exists {
		return fmt.Errorf("node %q not found", id)
	}
	delete(g.nodes, id)
	delete(g.out, id)
	delete(g.in, id)
	for from := range g.out {
		delete(g.out[from], id)
	}
	for to := range g.in {
		delete(g.in[to], id)
	}
	return nil
}

func (g *graphImpl) AddEdge(edge *Edge) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.nodes[edge.From]; !exists {
		return fmt.Errorf("node %q not found", edge.From)
	}
	if _, exists := g.nodes[edge.To]; !exists {
		return fmt.Errorf("node %q not found", edge.To)
	}
	if g.out[edge.From] == nil {
		g.out[edge.From] = make(map[string]*Edge)
	}
	if g.in[edge.To] == nil {
		g.in[edge.To] = make(map[string]*Edge)
	}
	if _, exists := g.out[edge.From][edge.To]; exists {
		return fmt.Errorf("edge from %q to %q already exists", edge.From, edge.To)
	}
	g.out[edge.From][edge.To] = edge
	g.in[edge.To][edge.From] = edge
	return nil
}

func (g *graphImpl) GetEdge(from, to string) (*Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	edge, exists := g.out[from][to]
	if !exists {
		return nil, fmt.Errorf("edge from %q to %q not found", from, to)
	}
	return edge, nil
}

func (g *graphImpl) DeleteEdge(from, to string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.out[from][to]; !exists {
		return fmt.Errorf("edge from %q to %q not found", from, to)
	}
	delete(g.out[from], to)
	delete(g.in[to], from)
	return nil
}

func (g *graphImpl) Neighbors(id string) ([]*Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.nodes[id]; !exists {
		return nil, fmt.Errorf("node %q not found", id)
	}
	edges := make([]*Edge, 0, len(g.out[id]))
	for _, edge := range g.out[id] {
		edges = append(edges, edge)
	}
	return edges, nil
}

func (g *graphImpl) Incoming(id string) ([]*Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.nodes[id]; !exists {
		return nil, fmt.Errorf("node %q not found", id)
	}
	edges := make([]*Edge, 0, len(g.in[id]))
	for _, edge := range g.in[id] {
		edges = append(edges, edge)
	}
	return edges, nil
}

// DFS performs depth-first search starting from the given node
func (g *graphImpl) DFS(start string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.nodes[start]; !exists {
		return nil, fmt.Errorf("node %q not found", start)
	}

	visited := make(map[string]bool)
	var result []string

	var dfs func(id string)
	dfs = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true
		result = append(result, id)

		// Visit neighbors in consistent order (sorted by target ID)
		for _, edge := range g.out[id] {
			dfs(edge.To)
		}
	}

	dfs(start)
	return result, nil
}

// BFS performs breadth-first search starting from the given node
func (g *graphImpl) BFS(start string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.nodes[start]; !exists {
		return nil, fmt.Errorf("node %q not found", start)
	}

	visited := make(map[string]bool)
	var result []string
	queue := []string{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		for _, edge := range g.out[current] {
			if !visited[edge.To] {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
	}

	return result, nil
}

// ShortestPath computes the shortest path using Dijkstra's algorithm
func (g *graphImpl) ShortestPath(source, target string) ([]string, float64, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.nodes[source]; !exists {
		return nil, 0, fmt.Errorf("source node %q not found", source)
	}
	if _, exists := g.nodes[target]; !exists {
		return nil, 0, fmt.Errorf("target node %q not found", target)
	}

	// Initialise distances and previous nodes
	dist := make(map[string]float64)
	prev := make(map[string]string)
	for id := range g.nodes {
		dist[id] = math.Inf(1)
	}
	dist[source] = 0

	// Priority queue for Dijkstra's
	priorityQueue := make(pq.PriorityQueue, 0)
	heap.Init(&priorityQueue)
	heap.Push(&priorityQueue, &pq.PriorityQueueItem{NodeID: source, Distance: 0})

	visited := make(map[string]bool)

	for priorityQueue.Len() > 0 {
		item := heap.Pop(&priorityQueue).(*pq.PriorityQueueItem)
		current := item.NodeID

		if visited[current] {
			continue
		}
		visited[current] = true

		// If we reached the target, reconstruct path
		if current == target {
			var path []string
			for node := target; node != ""; node = prev[node] {
				path = append([]string{node}, path...)
				if node == source {
					break
				}
			}
			return path, dist[target], nil
		}

		// Check neighbours
		for _, edge := range g.out[current] {
			neighbor := edge.To
			if visited[neighbor] {
				continue
			}

			alt := dist[current] + edge.Weight
			if alt < dist[neighbor] {
				dist[neighbor] = alt
				prev[neighbor] = current
				heap.Push(&priorityQueue, &pq.PriorityQueueItem{NodeID: neighbor, Distance: alt})
			}
		}
	}

	return nil, 0, fmt.Errorf("no path found from %q to %q", source, target)
}
