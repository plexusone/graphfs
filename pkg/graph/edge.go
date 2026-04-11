// Package graph provides core types for the GraphFS graph database.
package graph

// Edge represents a relationship between two nodes.
type Edge struct {
	// From is the source node ID.
	From string `json:"from"`

	// To is the target node ID.
	To string `json:"to"`

	// Type categorizes the relationship (e.g., "calls", "imports", "implements").
	Type string `json:"type"`

	// Confidence indicates how the edge was determined.
	Confidence Confidence `json:"confidence"`

	// ConfidenceScore is a 0.0-1.0 score for INFERRED edges.
	// Only meaningful when Confidence is ConfidenceInferred.
	ConfidenceScore float64 `json:"confidence_score,omitempty"`

	// Attrs holds additional attributes as key-value pairs.
	Attrs map[string]string `json:"attrs,omitempty"`
}

// Confidence indicates how an edge relationship was determined.
type Confidence string

const (
	// ConfidenceExtracted means the edge was directly extracted from source
	// (e.g., import statement, function call in AST).
	ConfidenceExtracted Confidence = "EXTRACTED"

	// ConfidenceInferred means the edge was inferred by an LLM or heuristic,
	// with an associated confidence score.
	ConfidenceInferred Confidence = "INFERRED"

	// ConfidenceAmbiguous means the relationship is uncertain and should be
	// reviewed by a human.
	ConfidenceAmbiguous Confidence = "AMBIGUOUS"
)

// EdgeType constants for common relationship types.
const (
	EdgeTypeCalls      = "calls"
	EdgeTypeImports    = "imports"
	EdgeTypeImplements = "implements"
	EdgeTypeExtends    = "extends"
	EdgeTypeUses       = "uses"
	EdgeTypeContains   = "contains"
	EdgeTypeDependsOn  = "depends_on"
	EdgeTypeReferences = "references"
)
