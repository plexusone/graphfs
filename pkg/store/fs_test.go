package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/plexusone/graphfs/pkg/graph"
)

func TestNewFSStore(t *testing.T) {
	dir := t.TempDir()

	fs, err := NewFSStore(dir)
	if err != nil {
		t.Fatalf("NewFSStore failed: %v", err)
	}

	if fs.Root != dir {
		t.Errorf("expected root '%s', got '%s'", dir, fs.Root)
	}

	// Check directories were created
	nodesDir := filepath.Join(dir, "nodes")
	edgesDir := filepath.Join(dir, "edges")

	if _, err := os.Stat(nodesDir); os.IsNotExist(err) {
		t.Error("nodes directory not created")
	}
	if _, err := os.Stat(edgesDir); os.IsNotExist(err) {
		t.Error("edges directory not created")
	}
}

func TestFSStore_WriteAndGetNode(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	node := &graph.Node{
		ID:    "test_node",
		Type:  graph.NodeTypeFunction,
		Label: "TestNode",
		Attrs: map[string]string{"key": "value"},
	}

	if err := fs.WriteNode(node); err != nil {
		t.Fatalf("WriteNode failed: %v", err)
	}

	// Verify file exists
	path := filepath.Join(dir, "nodes", "test_node.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("node file not created")
	}

	// Read it back
	got, err := fs.GetNode("test_node")
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}

	if got.ID != node.ID {
		t.Errorf("expected ID '%s', got '%s'", node.ID, got.ID)
	}
	if got.Type != node.Type {
		t.Errorf("expected Type '%s', got '%s'", node.Type, got.Type)
	}
	if got.Label != node.Label {
		t.Errorf("expected Label '%s', got '%s'", node.Label, got.Label)
	}
	if got.Attrs["key"] != "value" {
		t.Errorf("expected Attrs[key]='value', got '%s'", got.Attrs["key"])
	}
}

func TestFSStore_GetNode_NotFound(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	_, err := fs.GetNode("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent node")
	}
}

func TestFSStore_WriteAndGetEdge(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	edge := &graph.Edge{
		From:            "node_a",
		To:              "node_b",
		Type:            graph.EdgeTypeCalls,
		Confidence:      graph.ConfidenceInferred,
		ConfidenceScore: 0.85,
		Attrs:           map[string]string{"reason": "test"},
	}

	if err := fs.WriteEdge(edge); err != nil {
		t.Fatalf("WriteEdge failed: %v", err)
	}

	// Verify file exists
	path := filepath.Join(dir, "edges", "node_a__calls__node_b.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("edge file not created")
	}

	// Read it back
	got, err := fs.GetEdge("node_a", "calls", "node_b")
	if err != nil {
		t.Fatalf("GetEdge failed: %v", err)
	}

	if got.From != edge.From {
		t.Errorf("expected From '%s', got '%s'", edge.From, got.From)
	}
	if got.To != edge.To {
		t.Errorf("expected To '%s', got '%s'", edge.To, got.To)
	}
	if got.Type != edge.Type {
		t.Errorf("expected Type '%s', got '%s'", edge.Type, got.Type)
	}
	if got.Confidence != edge.Confidence {
		t.Errorf("expected Confidence '%s', got '%s'", edge.Confidence, got.Confidence)
	}
	if got.ConfidenceScore != edge.ConfidenceScore {
		t.Errorf("expected ConfidenceScore %f, got %f", edge.ConfidenceScore, got.ConfidenceScore)
	}
}

func TestFSStore_GetEdge_NotFound(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	_, err := fs.GetEdge("a", "calls", "b")
	if err == nil {
		t.Error("expected error for nonexistent edge")
	}
}

func TestFSStore_DeleteNode(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	node := &graph.Node{ID: "to_delete", Type: "function"}
	_ = fs.WriteNode(node)

	if err := fs.DeleteNode("to_delete"); err != nil {
		t.Fatalf("DeleteNode failed: %v", err)
	}

	_, err := fs.GetNode("to_delete")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestFSStore_DeleteNode_NonExistent(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	// Should not error for non-existent node
	if err := fs.DeleteNode("nonexistent"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFSStore_DeleteEdge(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	edge := &graph.Edge{
		From:       "a",
		To:         "b",
		Type:       "calls",
		Confidence: graph.ConfidenceExtracted,
	}
	_ = fs.WriteEdge(edge)

	if err := fs.DeleteEdge("a", "calls", "b"); err != nil {
		t.Fatalf("DeleteEdge failed: %v", err)
	}

	_, err := fs.GetEdge("a", "calls", "b")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestFSStore_ListNodes(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	nodes := []*graph.Node{
		{ID: "node_1", Type: "function"},
		{ID: "node_2", Type: "file"},
		{ID: "node_3", Type: "package"},
	}

	for _, n := range nodes {
		_ = fs.WriteNode(n)
	}

	got, err := fs.ListNodes()
	if err != nil {
		t.Fatalf("ListNodes failed: %v", err)
	}

	if len(got) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(got))
	}
}

func TestFSStore_ListEdges(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	edges := []*graph.Edge{
		{From: "a", To: "b", Type: "calls", Confidence: graph.ConfidenceExtracted},
		{From: "b", To: "c", Type: "imports", Confidence: graph.ConfidenceExtracted},
	}

	for _, e := range edges {
		_ = fs.WriteEdge(e)
	}

	got, err := fs.ListEdges()
	if err != nil {
		t.Fatalf("ListEdges failed: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("expected 2 edges, got %d", len(got))
	}
}

func TestFSStore_SaveAndLoadGraph(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	g := graph.NewGraph()
	g.AddNode(&graph.Node{ID: "a", Type: "function", Label: "FuncA"})
	g.AddNode(&graph.Node{ID: "b", Type: "function", Label: "FuncB"})
	g.AddEdge(&graph.Edge{
		From:       "a",
		To:         "b",
		Type:       "calls",
		Confidence: graph.ConfidenceExtracted,
	})

	if err := fs.SaveGraph(g); err != nil {
		t.Fatalf("SaveGraph failed: %v", err)
	}

	loaded, err := fs.LoadGraph()
	if err != nil {
		t.Fatalf("LoadGraph failed: %v", err)
	}

	if loaded.NodeCount() != 2 {
		t.Errorf("expected 2 nodes, got %d", loaded.NodeCount())
	}
	if loaded.EdgeCount() != 1 {
		t.Errorf("expected 1 edge, got %d", loaded.EdgeCount())
	}

	nodeA := loaded.GetNode("a")
	if nodeA == nil {
		t.Error("node 'a' not found")
	} else if nodeA.Label != "FuncA" {
		t.Errorf("expected Label 'FuncA', got '%s'", nodeA.Label)
	}
}

func TestFSStore_LoadGraph_Empty(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFSStore(dir)

	g, err := fs.LoadGraph()
	if err != nil {
		t.Fatalf("LoadGraph failed: %v", err)
	}

	if g.NodeCount() != 0 {
		t.Errorf("expected 0 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Errorf("expected 0 edges, got %d", g.EdgeCount())
	}
}
