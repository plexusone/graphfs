package graph

import "testing"

func TestNewGraph(t *testing.T) {
	g := NewGraph()

	if g == nil {
		t.Fatal("NewGraph returned nil")
	}
	if g.Nodes == nil {
		t.Error("Nodes map is nil")
	}
	if g.Edges == nil {
		t.Error("Edges slice is nil")
	}
	if g.NodeCount() != 0 {
		t.Errorf("expected 0 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Errorf("expected 0 edges, got %d", g.EdgeCount())
	}
}

func TestGraph_AddNode(t *testing.T) {
	g := NewGraph()

	node := &Node{
		ID:    "test_node",
		Type:  NodeTypeFunction,
		Label: "TestNode",
	}

	g.AddNode(node)

	if g.NodeCount() != 1 {
		t.Errorf("expected 1 node, got %d", g.NodeCount())
	}

	got := g.GetNode("test_node")
	if got == nil {
		t.Fatal("GetNode returned nil")
	}
	if got.ID != "test_node" {
		t.Errorf("expected ID 'test_node', got '%s'", got.ID)
	}
	if got.Type != NodeTypeFunction {
		t.Errorf("expected Type 'function', got '%s'", got.Type)
	}
}

func TestGraph_AddNode_Overwrite(t *testing.T) {
	g := NewGraph()

	node1 := &Node{ID: "same_id", Type: NodeTypeFunction, Label: "First"}
	node2 := &Node{ID: "same_id", Type: NodeTypeMethod, Label: "Second"}

	g.AddNode(node1)
	g.AddNode(node2)

	if g.NodeCount() != 1 {
		t.Errorf("expected 1 node after overwrite, got %d", g.NodeCount())
	}

	got := g.GetNode("same_id")
	if got.Label != "Second" {
		t.Errorf("expected Label 'Second', got '%s'", got.Label)
	}
}

func TestGraph_AddEdge(t *testing.T) {
	g := NewGraph()

	edge := &Edge{
		From:       "node_a",
		To:         "node_b",
		Type:       EdgeTypeCalls,
		Confidence: ConfidenceExtracted,
	}

	g.AddEdge(edge)

	if g.EdgeCount() != 1 {
		t.Errorf("expected 1 edge, got %d", g.EdgeCount())
	}
}

func TestGraph_GetNode_NotFound(t *testing.T) {
	g := NewGraph()

	got := g.GetNode("nonexistent")
	if got != nil {
		t.Error("expected nil for nonexistent node")
	}
}

func TestNodeTypeConstants(t *testing.T) {
	tests := []struct {
		constant string
		value    string
	}{
		{NodeTypeFunction, "function"},
		{NodeTypeMethod, "method"},
		{NodeTypeClass, "class"},
		{NodeTypeStruct, "struct"},
		{NodeTypeFile, "file"},
		{NodeTypePackage, "package"},
		{NodeTypeModule, "module"},
		{NodeTypeVariable, "variable"},
		{NodeTypeConstant, "constant"},
		{NodeTypeInterface, "interface"},
	}

	for _, tt := range tests {
		if tt.constant != tt.value {
			t.Errorf("expected %s, got %s", tt.value, tt.constant)
		}
	}
}

func TestEdgeTypeConstants(t *testing.T) {
	tests := []struct {
		constant string
		value    string
	}{
		{EdgeTypeCalls, "calls"},
		{EdgeTypeImports, "imports"},
		{EdgeTypeImplements, "implements"},
		{EdgeTypeExtends, "extends"},
		{EdgeTypeUses, "uses"},
		{EdgeTypeContains, "contains"},
		{EdgeTypeDependsOn, "depends_on"},
		{EdgeTypeReferences, "references"},
	}

	for _, tt := range tests {
		if tt.constant != tt.value {
			t.Errorf("expected %s, got %s", tt.value, tt.constant)
		}
	}
}

func TestConfidenceConstants(t *testing.T) {
	tests := []struct {
		constant Confidence
		value    string
	}{
		{ConfidenceExtracted, "EXTRACTED"},
		{ConfidenceInferred, "INFERRED"},
		{ConfidenceAmbiguous, "AMBIGUOUS"},
	}

	for _, tt := range tests {
		if string(tt.constant) != tt.value {
			t.Errorf("expected %s, got %s", tt.value, tt.constant)
		}
	}
}
