package schema

import (
	"testing"

	"github.com/plexusone/graphfs/pkg/graph"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("NewValidator returned nil")
	}
}

func TestValidateNode_Valid(t *testing.T) {
	v := NewValidator()

	node := &graph.Node{
		ID:   "test_node",
		Type: graph.NodeTypeFunction,
	}

	if err := v.ValidateNode(node); err != nil {
		t.Errorf("expected valid node, got error: %v", err)
	}
}

func TestValidateNode_MissingID(t *testing.T) {
	v := NewValidator()

	node := &graph.Node{
		Type: graph.NodeTypeFunction,
	}

	err := v.ValidateNode(node)
	if err == nil {
		t.Error("expected error for missing ID")
	}

	verr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if verr.Field != "id" {
		t.Errorf("expected field 'id', got '%s'", verr.Field)
	}
}

func TestValidateNode_MissingType(t *testing.T) {
	v := NewValidator()

	node := &graph.Node{
		ID: "test_node",
	}

	err := v.ValidateNode(node)
	if err == nil {
		t.Error("expected error for missing type")
	}

	verr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if verr.Field != "type" {
		t.Errorf("expected field 'type', got '%s'", verr.Field)
	}
}

func TestValidateNode_InvalidIDCharacters(t *testing.T) {
	v := NewValidator()

	invalidIDs := []string{
		"path/to/node",
		"node\\with\\backslash",
		"node:colon",
		"node*star",
		"node?question",
		"node\"quote",
		"node<less",
		"node>greater",
		"node|pipe",
	}

	for _, id := range invalidIDs {
		node := &graph.Node{ID: id, Type: "function"}
		err := v.ValidateNode(node)
		if err == nil {
			t.Errorf("expected error for invalid ID '%s'", id)
		}
	}
}

func TestValidateNode_RequireLabel(t *testing.T) {
	v := NewValidator()
	v.RequireNodeLabel = true

	node := &graph.Node{
		ID:   "test_node",
		Type: graph.NodeTypeFunction,
	}

	err := v.ValidateNode(node)
	if err == nil {
		t.Error("expected error for missing label")
	}

	node.Label = "TestNode"
	if err := v.ValidateNode(node); err != nil {
		t.Errorf("expected valid node with label, got error: %v", err)
	}
}

func TestValidateNode_AllowedTypes(t *testing.T) {
	v := NewValidator()
	v.AllowedNodeTypes = []string{"function", "file"}

	validNode := &graph.Node{ID: "test", Type: "function"}
	if err := v.ValidateNode(validNode); err != nil {
		t.Errorf("expected valid node, got error: %v", err)
	}

	invalidNode := &graph.Node{ID: "test", Type: "class"}
	if err := v.ValidateNode(invalidNode); err == nil {
		t.Error("expected error for disallowed type")
	}
}

func TestValidateEdge_Valid(t *testing.T) {
	v := NewValidator()

	edge := &graph.Edge{
		From:       "node_a",
		To:         "node_b",
		Type:       graph.EdgeTypeCalls,
		Confidence: graph.ConfidenceExtracted,
	}

	if err := v.ValidateEdge(edge); err != nil {
		t.Errorf("expected valid edge, got error: %v", err)
	}
}

func TestValidateEdge_MissingFrom(t *testing.T) {
	v := NewValidator()

	edge := &graph.Edge{
		To:         "node_b",
		Type:       graph.EdgeTypeCalls,
		Confidence: graph.ConfidenceExtracted,
	}

	err := v.ValidateEdge(edge)
	if err == nil {
		t.Error("expected error for missing from")
	}
}

func TestValidateEdge_MissingTo(t *testing.T) {
	v := NewValidator()

	edge := &graph.Edge{
		From:       "node_a",
		Type:       graph.EdgeTypeCalls,
		Confidence: graph.ConfidenceExtracted,
	}

	err := v.ValidateEdge(edge)
	if err == nil {
		t.Error("expected error for missing to")
	}
}

func TestValidateEdge_MissingType(t *testing.T) {
	v := NewValidator()

	edge := &graph.Edge{
		From:       "node_a",
		To:         "node_b",
		Confidence: graph.ConfidenceExtracted,
	}

	err := v.ValidateEdge(edge)
	if err == nil {
		t.Error("expected error for missing type")
	}
}

func TestValidateEdge_MissingConfidence(t *testing.T) {
	v := NewValidator()

	edge := &graph.Edge{
		From: "node_a",
		To:   "node_b",
		Type: graph.EdgeTypeCalls,
	}

	err := v.ValidateEdge(edge)
	if err == nil {
		t.Error("expected error for missing confidence")
	}
}

func TestValidateEdge_InvalidConfidence(t *testing.T) {
	v := NewValidator()

	edge := &graph.Edge{
		From:       "node_a",
		To:         "node_b",
		Type:       graph.EdgeTypeCalls,
		Confidence: graph.Confidence("INVALID"),
	}

	err := v.ValidateEdge(edge)
	if err == nil {
		t.Error("expected error for invalid confidence")
	}
}

func TestValidateEdge_InferredScoreRange(t *testing.T) {
	v := NewValidator()

	// Valid score
	edge := &graph.Edge{
		From:            "node_a",
		To:              "node_b",
		Type:            graph.EdgeTypeCalls,
		Confidence:      graph.ConfidenceInferred,
		ConfidenceScore: 0.75,
	}

	if err := v.ValidateEdge(edge); err != nil {
		t.Errorf("expected valid edge, got error: %v", err)
	}

	// Invalid score (negative)
	edge.ConfidenceScore = -0.1
	if err := v.ValidateEdge(edge); err == nil {
		t.Error("expected error for negative score")
	}

	// Invalid score (> 1)
	edge.ConfidenceScore = 1.5
	if err := v.ValidateEdge(edge); err == nil {
		t.Error("expected error for score > 1")
	}
}

func TestValidateEdge_AllowedTypes(t *testing.T) {
	v := NewValidator()
	v.AllowedEdgeTypes = []string{"calls", "imports"}

	validEdge := &graph.Edge{
		From:       "a",
		To:         "b",
		Type:       "calls",
		Confidence: graph.ConfidenceExtracted,
	}
	if err := v.ValidateEdge(validEdge); err != nil {
		t.Errorf("expected valid edge, got error: %v", err)
	}

	invalidEdge := &graph.Edge{
		From:       "a",
		To:         "b",
		Type:       "extends",
		Confidence: graph.ConfidenceExtracted,
	}
	if err := v.ValidateEdge(invalidEdge); err == nil {
		t.Error("expected error for disallowed type")
	}
}

func TestValidateEdgeRefs_Valid(t *testing.T) {
	v := NewValidator()

	nodes := map[string]*graph.Node{
		"node_a": {ID: "node_a", Type: "function"},
		"node_b": {ID: "node_b", Type: "function"},
	}

	edge := &graph.Edge{From: "node_a", To: "node_b"}

	if err := v.ValidateEdgeRefs(edge, nodes); err != nil {
		t.Errorf("expected valid refs, got error: %v", err)
	}
}

func TestValidateEdgeRefs_MissingFrom(t *testing.T) {
	v := NewValidator()

	nodes := map[string]*graph.Node{
		"node_b": {ID: "node_b", Type: "function"},
	}

	edge := &graph.Edge{From: "node_a", To: "node_b"}

	err := v.ValidateEdgeRefs(edge, nodes)
	if err == nil {
		t.Error("expected error for missing from node")
	}
}

func TestValidateEdgeRefs_MissingTo(t *testing.T) {
	v := NewValidator()

	nodes := map[string]*graph.Node{
		"node_a": {ID: "node_a", Type: "function"},
	}

	edge := &graph.Edge{From: "node_a", To: "node_b"}

	err := v.ValidateEdgeRefs(edge, nodes)
	if err == nil {
		t.Error("expected error for missing to node")
	}
}

func TestValidateGraph(t *testing.T) {
	v := NewValidator()

	g := graph.NewGraph()
	g.AddNode(&graph.Node{ID: "a", Type: "function"})
	g.AddNode(&graph.Node{ID: "b", Type: "function"})
	g.AddEdge(&graph.Edge{
		From:       "a",
		To:         "b",
		Type:       "calls",
		Confidence: graph.ConfidenceExtracted,
	})

	errs := v.ValidateGraph(g)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateGraph_WithErrors(t *testing.T) {
	v := NewValidator()

	g := graph.NewGraph()
	g.AddNode(&graph.Node{ID: "a", Type: "function"})
	// Edge references non-existent node
	g.AddEdge(&graph.Edge{
		From:       "a",
		To:         "nonexistent",
		Type:       "calls",
		Confidence: graph.ConfidenceExtracted,
	})

	errs := v.ValidateGraph(g)
	if len(errs) == 0 {
		t.Error("expected errors for dangling reference")
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{Field: "id", Message: "is required"}
	expected := "id: is required"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}
