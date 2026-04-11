// Package schema provides validation for graph data.
package schema

import (
	"fmt"
	"strings"

	"github.com/plexusone/graphfs/pkg/graph"
)

// Validator validates graph data against schema rules.
type Validator struct {
	// AllowedNodeTypes restricts node types. If empty, all types are allowed.
	AllowedNodeTypes []string

	// AllowedEdgeTypes restricts edge types. If empty, all types are allowed.
	AllowedEdgeTypes []string

	// RequireNodeLabel requires all nodes to have a non-empty label.
	RequireNodeLabel bool
}

// NewValidator creates a validator with default settings.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidationError represents a validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateNode validates a single node.
func (v *Validator) ValidateNode(n *graph.Node) error {
	if n.ID == "" {
		return &ValidationError{Field: "id", Message: "node ID is required"}
	}

	// Check for invalid characters in ID (must be filesystem-safe)
	if strings.ContainsAny(n.ID, "/\\:*?\"<>|") {
		return &ValidationError{Field: "id", Message: "node ID contains invalid characters"}
	}

	if n.Type == "" {
		return &ValidationError{Field: "type", Message: "node type is required"}
	}

	if v.RequireNodeLabel && n.Label == "" {
		return &ValidationError{Field: "label", Message: "node label is required"}
	}

	if len(v.AllowedNodeTypes) > 0 && !contains(v.AllowedNodeTypes, n.Type) {
		return &ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("node type %q is not allowed", n.Type),
		}
	}

	return nil
}

// ValidateEdge validates a single edge.
func (v *Validator) ValidateEdge(e *graph.Edge) error {
	if e.From == "" {
		return &ValidationError{Field: "from", Message: "edge source is required"}
	}

	if e.To == "" {
		return &ValidationError{Field: "to", Message: "edge target is required"}
	}

	if e.Type == "" {
		return &ValidationError{Field: "type", Message: "edge type is required"}
	}

	if e.Confidence == "" {
		return &ValidationError{Field: "confidence", Message: "edge confidence is required"}
	}

	// Validate confidence value
	switch e.Confidence {
	case graph.ConfidenceExtracted, graph.ConfidenceInferred, graph.ConfidenceAmbiguous:
		// Valid
	default:
		return &ValidationError{
			Field:   "confidence",
			Message: fmt.Sprintf("invalid confidence value %q", e.Confidence),
		}
	}

	// Confidence score only meaningful for inferred edges
	if e.Confidence == graph.ConfidenceInferred {
		if e.ConfidenceScore < 0 || e.ConfidenceScore > 1 {
			return &ValidationError{
				Field:   "confidence_score",
				Message: "confidence score must be between 0 and 1",
			}
		}
	}

	if len(v.AllowedEdgeTypes) > 0 && !contains(v.AllowedEdgeTypes, e.Type) {
		return &ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("edge type %q is not allowed", e.Type),
		}
	}

	return nil
}

// ValidateEdgeRefs validates that edge references exist in the node set.
func (v *Validator) ValidateEdgeRefs(e *graph.Edge, nodes map[string]*graph.Node) error {
	if _, ok := nodes[e.From]; !ok {
		return &ValidationError{
			Field:   "from",
			Message: fmt.Sprintf("source node %q does not exist", e.From),
		}
	}

	if _, ok := nodes[e.To]; !ok {
		return &ValidationError{
			Field:   "to",
			Message: fmt.Sprintf("target node %q does not exist", e.To),
		}
	}

	return nil
}

// ValidateGraph validates an entire graph.
func (v *Validator) ValidateGraph(g *graph.Graph) []error {
	var errs []error

	// Validate all nodes
	for id, n := range g.Nodes {
		if err := v.ValidateNode(n); err != nil {
			errs = append(errs, fmt.Errorf("node %s: %w", id, err))
		}
	}

	// Validate all edges
	for i, e := range g.Edges {
		if err := v.ValidateEdge(e); err != nil {
			errs = append(errs, fmt.Errorf("edge[%d]: %w", i, err))
		}
		if err := v.ValidateEdgeRefs(e, g.Nodes); err != nil {
			errs = append(errs, fmt.Errorf("edge[%d]: %w", i, err))
		}
	}

	return errs
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
