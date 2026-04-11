package analyze

import (
	"testing"

	"github.com/plexusone/graphfs/pkg/graph"
)

func TestFindHubs(t *testing.T) {
	nodes := []*graph.Node{
		{ID: "A", Type: "func", Label: "FuncA"},
		{ID: "B", Type: "func", Label: "FuncB"},
		{ID: "C", Type: "func", Label: "FuncC"},
		{ID: "D", Type: "file", Label: "File"},                                // Excluded type
		{ID: "E", Type: "func", Attrs: map[string]string{"external": "true"}}, // External
	}

	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
		{From: "A", To: "C", Type: "calls"},
		{From: "B", To: "C", Type: "calls"},
		{From: "D", To: "A", Type: "contains"}, // From excluded type
	}

	t.Run("find top 2 hubs", func(t *testing.T) {
		hubs := FindHubs(nodes, edges, 2, []string{"file"})
		if len(hubs) != 2 {
			t.Errorf("Expected 2 hubs, got %d", len(hubs))
		}
		if hubs[0].ID != "A" {
			t.Errorf("Expected top hub to be A (3 connections), got %s", hubs[0].ID)
		}
	})

	t.Run("excludes external nodes", func(t *testing.T) {
		hubs := FindHubs(nodes, edges, 10, nil)
		for _, h := range hubs {
			if h.ID == "E" {
				t.Error("External node E should be excluded")
			}
		}
	})
}

func TestIsolatedNodes(t *testing.T) {
	nodes := []*graph.Node{
		{ID: "A", Type: "func"},
		{ID: "B", Type: "func"},
		{ID: "C", Type: "func"},
		{ID: "D", Type: "func"},
	}

	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
	}

	isolated := IsolatedNodes(nodes, edges, 0, nil)
	if len(isolated) != 2 { // C and D have 0 degree
		t.Errorf("Expected 2 isolated nodes, got %d", len(isolated))
	}
}

func TestNodesByType(t *testing.T) {
	nodes := []*graph.Node{
		{ID: "A", Type: "func"},
		{ID: "B", Type: "func"},
		{ID: "C", Type: "type"},
	}

	byType := NodesByType(nodes)
	if len(byType["func"]) != 2 {
		t.Errorf("Expected 2 func nodes, got %d", len(byType["func"]))
	}
	if len(byType["type"]) != 1 {
		t.Errorf("Expected 1 type node, got %d", len(byType["type"]))
	}
}

func TestEdgesByType(t *testing.T) {
	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
		{From: "B", To: "C", Type: "calls"},
		{From: "A", To: "C", Type: "imports"},
	}

	byType := EdgesByType(edges)
	if len(byType["calls"]) != 2 {
		t.Errorf("Expected 2 calls edges, got %d", len(byType["calls"]))
	}
	if len(byType["imports"]) != 1 {
		t.Errorf("Expected 1 imports edge, got %d", len(byType["imports"]))
	}
}

func TestHubScore(t *testing.T) {
	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
		{From: "A", To: "C", Type: "calls"},
		{From: "B", To: "C", Type: "calls"},
	}

	score := HubScore("A", edges)
	if score != 2 {
		t.Errorf("Expected hub score 2 for A, got %d", score)
	}
}

func TestAuthorityScore(t *testing.T) {
	edges := []*graph.Edge{
		{From: "A", To: "C", Type: "calls"},
		{From: "B", To: "C", Type: "calls"},
	}

	score := AuthorityScore("C", edges)
	if score != 2 {
		t.Errorf("Expected authority score 2 for C, got %d", score)
	}
}

func TestDiffGraphs(t *testing.T) {
	oldNodes := []*graph.Node{
		{ID: "A", Type: "func", Label: "FuncA"},
		{ID: "B", Type: "func", Label: "FuncB"},
	}
	oldEdges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
	}

	newNodes := []*graph.Node{
		{ID: "A", Type: "func", Label: "FuncA"},
		{ID: "C", Type: "func", Label: "FuncC"}, // Added
	}
	newEdges := []*graph.Edge{
		{From: "A", To: "C", Type: "calls"}, // Added
	}

	diff := DiffGraphs(oldNodes, newNodes, oldEdges, newEdges)

	if len(diff.NewNodes) != 1 || diff.NewNodes[0].ID != "C" {
		t.Errorf("Expected 1 new node C, got %v", diff.NewNodes)
	}
	if len(diff.RemovedNodes) != 1 || diff.RemovedNodes[0].ID != "B" {
		t.Errorf("Expected 1 removed node B, got %v", diff.RemovedNodes)
	}
	if len(diff.NewEdges) != 1 {
		t.Errorf("Expected 1 new edge, got %d", len(diff.NewEdges))
	}
	if len(diff.RemovedEdges) != 1 {
		t.Errorf("Expected 1 removed edge, got %d", len(diff.RemovedEdges))
	}
}

func TestDiffGraphsNoChanges(t *testing.T) {
	nodes := []*graph.Node{
		{ID: "A", Type: "func"},
	}
	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
	}

	diff := DiffGraphs(nodes, nodes, edges, edges)
	if diff.Summary != "No changes detected" {
		t.Errorf("Expected 'No changes detected', got %s", diff.Summary)
	}
}

func TestCohesionScore(t *testing.T) {
	adj := map[string]map[string]bool{
		"A": {"B": true, "C": true},
		"B": {"A": true, "C": true},
		"C": {"A": true, "B": true},
	}

	// Perfect clique
	score := CohesionScore([]string{"A", "B", "C"}, adj)
	if score != 1.0 {
		t.Errorf("Expected cohesion 1.0 for complete clique, got %f", score)
	}

	// Single node
	score = CohesionScore([]string{"A"}, adj)
	if score != 1.0 {
		t.Errorf("Expected cohesion 1.0 for single node, got %f", score)
	}
}

func TestDetectCommunities(t *testing.T) {
	nodes := []*graph.Node{
		{ID: "A", Type: "func"},
		{ID: "B", Type: "func"},
		{ID: "C", Type: "func"},
	}
	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
		{From: "B", To: "C", Type: "calls"},
	}

	result := DetectCommunities(nodes, edges)
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if len(result.Communities) == 0 {
		t.Error("Expected at least one community")
	}
}

func TestDetectCommunitiesLouvain(t *testing.T) {
	nodes := []*graph.Node{
		{ID: "A", Type: "func"},
		{ID: "B", Type: "func"},
		{ID: "C", Type: "func"},
		{ID: "D", Type: "func"},
	}
	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls"},
		{From: "A", To: "C", Type: "calls"},
		{From: "B", To: "C", Type: "calls"},
		{From: "C", To: "D", Type: "calls"},
	}

	opts := DefaultLouvainOptions()
	result := DetectCommunitiesLouvain(nodes, edges, opts)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if len(result.Communities) == 0 {
		t.Error("Expected at least one community")
	}
	if len(result.NodeCommunity) != 4 {
		t.Errorf("Expected 4 node community mappings, got %d", len(result.NodeCommunity))
	}
}

func TestDetectCommunitiesLouvainEmptyGraph(t *testing.T) {
	var nodes []*graph.Node
	var edges []*graph.Edge

	opts := DefaultLouvainOptions()
	result := DetectCommunitiesLouvain(nodes, edges, opts)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if len(result.Communities) != 0 {
		t.Errorf("Expected 0 communities for empty graph, got %d", len(result.Communities))
	}
}

func TestDetectCommunitiesLouvainNoEdges(t *testing.T) {
	nodes := []*graph.Node{
		{ID: "A", Type: "func"},
		{ID: "B", Type: "func"},
	}
	var edges []*graph.Edge

	opts := DefaultLouvainOptions()
	result := DetectCommunitiesLouvain(nodes, edges, opts)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	// Each node should be its own community
	if len(result.Communities) != 2 {
		t.Errorf("Expected 2 communities (one per node), got %d", len(result.Communities))
	}
}

func TestInferredEdges(t *testing.T) {
	edges := []*graph.Edge{
		{From: "A", To: "B", Type: "calls", Confidence: graph.ConfidenceExtracted},
		{From: "A", To: "C", Type: "calls", Confidence: graph.ConfidenceInferred},
		{From: "B", To: "D", Type: "calls", Confidence: graph.ConfidenceAmbiguous},
	}

	inferred := InferredEdges(edges)
	if len(inferred) != 2 {
		t.Errorf("Expected 2 inferred/ambiguous edges, got %d", len(inferred))
	}
}
