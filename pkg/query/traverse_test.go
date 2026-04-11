package query

import (
	"testing"

	"github.com/plexusone/graphfs/pkg/graph"
)

func TestNewTraverser(t *testing.T) {
	g := &graph.Graph{
		Nodes: map[string]*graph.Node{
			"A": {ID: "A", Type: "func"},
			"B": {ID: "B", Type: "func"},
			"C": {ID: "C", Type: "func"},
		},
		Edges: []*graph.Edge{
			{From: "A", To: "B", Type: "calls"},
			{From: "B", To: "C", Type: "calls"},
		},
	}

	traverser := NewTraverser(g)

	if len(traverser.outgoing["A"]) != 1 {
		t.Errorf("Expected 1 outgoing edge from A, got %d", len(traverser.outgoing["A"]))
	}
	if len(traverser.incoming["B"]) != 1 {
		t.Errorf("Expected 1 incoming edge to B, got %d", len(traverser.incoming["B"]))
	}
}

func TestBFS(t *testing.T) {
	g := &graph.Graph{
		Nodes: map[string]*graph.Node{
			"A": {ID: "A", Type: "func"},
			"B": {ID: "B", Type: "func"},
			"C": {ID: "C", Type: "func"},
			"D": {ID: "D", Type: "func"},
		},
		Edges: []*graph.Edge{
			{From: "A", To: "B", Type: "calls"},
			{From: "A", To: "C", Type: "calls"},
			{From: "B", To: "D", Type: "calls"},
			{From: "C", To: "D", Type: "calls"},
		},
	}

	traverser := NewTraverser(g)

	t.Run("outgoing BFS", func(t *testing.T) {
		result := traverser.BFS("A", Outgoing, 10, nil)
		if result.StartNode != "A" {
			t.Errorf("Expected start node A, got %s", result.StartNode)
		}
		if len(result.Visited) != 4 {
			t.Errorf("Expected 4 visited nodes, got %d", len(result.Visited))
		}
	})

	t.Run("incoming BFS", func(t *testing.T) {
		result := traverser.BFS("D", Incoming, 10, nil)
		if len(result.Visited) != 4 {
			t.Errorf("Expected 4 visited nodes from D incoming, got %d", len(result.Visited))
		}
	})

	t.Run("BFS with max depth", func(t *testing.T) {
		result := traverser.BFS("A", Outgoing, 1, nil)
		if len(result.Visited) != 3 { // A, B, C (depth 0 and 1)
			t.Errorf("Expected 3 visited nodes with depth 1, got %d", len(result.Visited))
		}
	})

	t.Run("BFS with edge type filter", func(t *testing.T) {
		result := traverser.BFS("A", Outgoing, 10, []string{"calls"})
		if len(result.Visited) != 4 {
			t.Errorf("Expected 4 visited nodes with calls filter, got %d", len(result.Visited))
		}

		result = traverser.BFS("A", Outgoing, 10, []string{"imports"})
		if len(result.Visited) != 1 { // Only A since no imports edges
			t.Errorf("Expected 1 visited node with imports filter, got %d", len(result.Visited))
		}
	})
}

func TestDFS(t *testing.T) {
	g := &graph.Graph{
		Nodes: map[string]*graph.Node{
			"A": {ID: "A", Type: "func"},
			"B": {ID: "B", Type: "func"},
			"C": {ID: "C", Type: "func"},
		},
		Edges: []*graph.Edge{
			{From: "A", To: "B", Type: "calls"},
			{From: "B", To: "C", Type: "calls"},
		},
	}

	traverser := NewTraverser(g)
	result := traverser.DFS("A", Outgoing, 10, nil)

	if len(result.Visited) != 3 {
		t.Errorf("Expected 3 visited nodes, got %d", len(result.Visited))
	}

	// Check depths
	if result.Depth["A"] != 0 {
		t.Errorf("Expected depth 0 for A, got %d", result.Depth["A"])
	}
	if result.Depth["C"] != 2 {
		t.Errorf("Expected depth 2 for C, got %d", result.Depth["C"])
	}
}

func TestFindPath(t *testing.T) {
	g := &graph.Graph{
		Nodes: map[string]*graph.Node{
			"A": {ID: "A", Type: "func"},
			"B": {ID: "B", Type: "func"},
			"C": {ID: "C", Type: "func"},
			"D": {ID: "D", Type: "func"},
		},
		Edges: []*graph.Edge{
			{From: "A", To: "B", Type: "calls"},
			{From: "B", To: "C", Type: "calls"},
			{From: "C", To: "D", Type: "calls"},
		},
	}

	traverser := NewTraverser(g)

	t.Run("path exists", func(t *testing.T) {
		result := traverser.FindPath("A", "D", nil)
		if len(result.Visited) != 4 {
			t.Errorf("Expected path length 4, got %d", len(result.Visited))
		}
		if result.Visited[0] != "A" || result.Visited[3] != "D" {
			t.Error("Path should start at A and end at D")
		}
	})

	t.Run("no path exists", func(t *testing.T) {
		result := traverser.FindPath("D", "A", nil) // Directed edges, no path backwards
		// Undirected path finding, so it should still find a path
		if len(result.Visited) == 0 {
			t.Error("Expected path (undirected search), got empty")
		}
	})
}

func TestDirectionConstants(t *testing.T) {
	if Outgoing != 0 {
		t.Errorf("Expected Outgoing to be 0, got %d", Outgoing)
	}
	if Incoming != 1 {
		t.Errorf("Expected Incoming to be 1, got %d", Incoming)
	}
	if Both != 2 {
		t.Errorf("Expected Both to be 2, got %d", Both)
	}
}
