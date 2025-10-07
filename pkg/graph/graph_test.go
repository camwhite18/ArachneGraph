package graph

import (
	"testing"
)

// TestAddNode tests adding nodes to the graph
func TestAddNode(t *testing.T) {
	g := NewGraph()

	node := &Node{ID: "A"}
	err := g.AddNode(node)
	if err != nil {
		t.Fatalf("Failed to add node: %v", err)
	}

	// Verify node was added
	retrieved, err := g.GetNode("A")
	if err != nil {
		t.Fatalf("Failed to get node: %v", err)
	}
	if retrieved.ID != "A" {
		t.Errorf("Expected node ID 'A', got '%s'", retrieved.ID)
	}
}

// TestAddDuplicateNode tests adding duplicate nodes
func TestAddDuplicateNode(t *testing.T) {
	g := NewGraph()

	node1 := &Node{ID: "A"}
	err := g.AddNode(node1)
	if err != nil {
		t.Fatalf("Failed to add first node: %v", err)
	}

	node2 := &Node{ID: "A"}
	err = g.AddNode(node2)
	if err == nil {
		t.Error("Expected error when adding duplicate node, got nil")
	}
}

// TestGetNonExistentNode tests retrieving a non-existent node
func TestGetNonExistentNode(t *testing.T) {
	g := NewGraph()

	_, err := g.GetNode("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent node, got nil")
	}
}

// TestDeleteNode tests deleting nodes
func TestDeleteNode(t *testing.T) {
	g := NewGraph()

	node := &Node{ID: "A"}
	g.AddNode(node)

	err := g.DeleteNode("A")
	if err != nil {
		t.Fatalf("Failed to delete node: %v", err)
	}

	// Verify node was deleted
	_, err = g.GetNode("A")
	if err == nil {
		t.Error("Expected error when getting deleted node, got nil")
	}
}

// TestDeleteNodeWithEdges tests that deleting a node removes its edges
func TestDeleteNodeWithEdges(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	nodeC := &Node{ID: "C"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)

	edge1 := &Edge{ID: "e1", From: "A", To: "B", Weight: 1.0}
	edge2 := &Edge{ID: "e2", From: "B", To: "C", Weight: 2.0}
	g.AddEdge(edge1)
	g.AddEdge(edge2)

	// Delete node B
	err := g.DeleteNode("B")
	if err != nil {
		t.Fatalf("Failed to delete node: %v", err)
	}

	// Verify edges involving B are gone
	_, err = g.GetEdge("A", "B")
	if err == nil {
		t.Error("Expected error when getting edge to deleted node, got nil")
	}

	_, err = g.GetEdge("B", "C")
	if err == nil {
		t.Error("Expected error when getting edge from deleted node, got nil")
	}
}

// TestAddEdge tests adding edges to the graph
func TestAddEdge(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)

	edge := &Edge{ID: "e1", From: "A", To: "B", Weight: 1.0}
	err := g.AddEdge(edge)
	if err != nil {
		t.Fatalf("Failed to add edge: %v", err)
	}

	// Verify edge was added
	retrieved, err := g.GetEdge("A", "B")
	if err != nil {
		t.Fatalf("Failed to get edge: %v", err)
	}
	if retrieved.From != "A" || retrieved.To != "B" {
		t.Errorf("Expected edge A->B, got %s->%s", retrieved.From, retrieved.To)
	}
}

// TestAddEdgeToNonExistentNode tests adding edges with non-existent nodes
func TestAddEdgeToNonExistentNode(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	g.AddNode(nodeA)

	edge := &Edge{ID: "e1", From: "A", To: "B", Weight: 1.0}
	err := g.AddEdge(edge)
	if err == nil {
		t.Error("Expected error when adding edge to non-existent node, got nil")
	}
}

// TestAddDuplicateEdge tests adding duplicate edges
func TestAddDuplicateEdge(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)

	edge1 := &Edge{ID: "e1", From: "A", To: "B", Weight: 1.0}
	err := g.AddEdge(edge1)
	if err != nil {
		t.Fatalf("Failed to add first edge: %v", err)
	}

	edge2 := &Edge{ID: "e2", From: "A", To: "B", Weight: 2.0}
	err = g.AddEdge(edge2)
	if err == nil {
		t.Error("Expected error when adding duplicate edge, got nil")
	}
}

// TestSelfLoop tests self-loop edges
func TestSelfLoop(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	g.AddNode(nodeA)

	edge := &Edge{ID: "e1", From: "A", To: "A", Weight: 1.0}
	err := g.AddEdge(edge)
	if err != nil {
		t.Fatalf("Failed to add self-loop: %v", err)
	}

	// Verify self-loop was added
	retrieved, err := g.GetEdge("A", "A")
	if err != nil {
		t.Fatalf("Failed to get self-loop: %v", err)
	}
	if retrieved.From != "A" || retrieved.To != "A" {
		t.Errorf("Expected self-loop A->A, got %s->%s", retrieved.From, retrieved.To)
	}
}

// TestDeleteEdge tests deleting edges
func TestDeleteEdge(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)

	edge := &Edge{ID: "e1", From: "A", To: "B", Weight: 1.0}
	g.AddEdge(edge)

	err := g.DeleteEdge("A", "B")
	if err != nil {
		t.Fatalf("Failed to delete edge: %v", err)
	}

	// Verify edge was deleted
	_, err = g.GetEdge("A", "B")
	if err == nil {
		t.Error("Expected error when getting deleted edge, got nil")
	}
}

// TestNeighborsAndIncoming tests Neighbors and Incoming methods
func TestNeighborsAndIncoming(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	nodeC := &Node{ID: "C"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)

	edge1 := &Edge{ID: "e1", From: "A", To: "B", Weight: 1.0}
	edge2 := &Edge{ID: "e2", From: "A", To: "C", Weight: 2.0}
	edge3 := &Edge{ID: "e3", From: "B", To: "C", Weight: 3.0}
	g.AddEdge(edge1)
	g.AddEdge(edge2)
	g.AddEdge(edge3)

	// Test Neighbors (outgoing edges from A)
	neighbors, err := g.Neighbors("A")
	if err != nil {
		t.Fatalf("Failed to get neighbors: %v", err)
	}
	if len(neighbors) != 2 {
		t.Errorf("Expected 2 neighbors, got %d", len(neighbors))
	}

	// Test Incoming (incoming edges to C)
	incoming, err := g.Incoming("C")
	if err != nil {
		t.Fatalf("Failed to get incoming edges: %v", err)
	}
	if len(incoming) != 2 {
		t.Errorf("Expected 2 incoming edges, got %d", len(incoming))
	}
}

// TestDFS tests depth-first search
func TestDFS(t *testing.T) {
	g := NewGraph()

	// Create a simple graph: A -> B -> C
	//                        A -> D
	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	nodeC := &Node{ID: "C"}
	nodeD := &Node{ID: "D"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)
	g.AddNode(nodeD)

	g.AddEdge(&Edge{ID: "e1", From: "A", To: "B", Weight: 1.0})
	g.AddEdge(&Edge{ID: "e2", From: "B", To: "C", Weight: 1.0})
	g.AddEdge(&Edge{ID: "e3", From: "A", To: "D", Weight: 1.0})

	result, err := g.DFS("A")
	if err != nil {
		t.Fatalf("DFS failed: %v", err)
	}

	// Verify DFS visited all reachable nodes
	if len(result) != 4 {
		t.Errorf("Expected 4 nodes visited, got %d", len(result))
	}

	// Verify A is first
	if result[0] != "A" {
		t.Errorf("Expected first node to be 'A', got '%s'", result[0])
	}
}

// TestBFS tests breadth-first search
func TestBFS(t *testing.T) {
	g := NewGraph()

	// Create a simple graph: A -> B -> C
	//                        A -> D
	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	nodeC := &Node{ID: "C"}
	nodeD := &Node{ID: "D"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)
	g.AddNode(nodeD)

	g.AddEdge(&Edge{ID: "e1", From: "A", To: "B", Weight: 1.0})
	g.AddEdge(&Edge{ID: "e2", From: "B", To: "C", Weight: 1.0})
	g.AddEdge(&Edge{ID: "e3", From: "A", To: "D", Weight: 1.0})

	result, err := g.BFS("A")
	if err != nil {
		t.Fatalf("BFS failed: %v", err)
	}

	// Verify BFS visited all reachable nodes
	if len(result) != 4 {
		t.Errorf("Expected 4 nodes visited, got %d", len(result))
	}

	// Verify A is first
	if result[0] != "A" {
		t.Errorf("Expected first node to be 'A', got '%s'", result[0])
	}

	// B and D should be at level 1 (positions 1 and 2)
	if !((result[1] == "B" || result[1] == "D") && (result[2] == "B" || result[2] == "D")) {
		t.Errorf("Expected B and D at positions 1 and 2")
	}

	// C should be at level 2 (position 3)
	if result[3] != "C" {
		t.Errorf("Expected C at position 3, got '%s'", result[3])
	}
}

// TestDijkstraShortestPath tests Dijkstra's algorithm
func TestDijkstraShortestPath(t *testing.T) {
	g := NewGraph()

	// Create test graph:
	// A --1--> B --2--> C
	// A --4--> C
	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	nodeC := &Node{ID: "C"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)

	g.AddEdge(&Edge{ID: "e1", From: "A", To: "B", Weight: 1.0})
	g.AddEdge(&Edge{ID: "e2", From: "B", To: "C", Weight: 2.0})
	g.AddEdge(&Edge{ID: "e3", From: "A", To: "C", Weight: 4.0})

	path, distance, err := g.ShortestPath("A", "C")
	if err != nil {
		t.Fatalf("ShortestPath failed: %v", err)
	}

	// Verify shortest path is A -> B -> C (cost 3)
	expectedPath := []string{"A", "B", "C"}
	if len(path) != len(expectedPath) {
		t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
	}

	for i, node := range expectedPath {
		if path[i] != node {
			t.Errorf("Expected node '%s' at position %d, got '%s'", node, i, path[i])
		}
	}

	if distance != 3.0 {
		t.Errorf("Expected distance 3.0, got %f", distance)
	}
}

// TestDijkstraNoPath tests Dijkstra's algorithm when no path exists
func TestDijkstraNoPath(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)

	// No edges, so no path from A to B
	_, _, err := g.ShortestPath("A", "B")
	if err == nil {
		t.Error("Expected error when no path exists, got nil")
	}
}

// TestDijkstraWithSelfLoop tests Dijkstra's algorithm with self-loops
func TestDijkstraWithSelfLoop(t *testing.T) {
	g := NewGraph()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	g.AddNode(nodeA)
	g.AddNode(nodeB)

	g.AddEdge(&Edge{ID: "e1", From: "A", To: "A", Weight: 1.0})
	g.AddEdge(&Edge{ID: "e2", From: "A", To: "B", Weight: 2.0})

	path, distance, err := g.ShortestPath("A", "B")
	if err != nil {
		t.Fatalf("ShortestPath failed: %v", err)
	}

	expectedPath := []string{"A", "B"}
	if len(path) != len(expectedPath) {
		t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
	}

	if distance != 2.0 {
		t.Errorf("Expected distance 2.0, got %f", distance)
	}
}
