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

	// Framework-specific edge types

	// EdgeTypeInjects represents dependency injection (Spring @Autowired, etc.)
	EdgeTypeInjects = "injects"

	// EdgeTypeHandlesRoute represents a controller method handling an HTTP route.
	EdgeTypeHandlesRoute = "handles_route"

	// EdgeTypeHasMany represents a one-to-many relationship (JPA @OneToMany).
	EdgeTypeHasMany = "has_many"

	// EdgeTypeBelongsTo represents a many-to-one relationship (JPA @ManyToOne).
	EdgeTypeBelongsTo = "belongs_to"

	// EdgeTypeAnnotatedWith represents a class/method annotated with a specific annotation.
	EdgeTypeAnnotatedWith = "annotated_with"

	// EdgeTypeMethodOf represents a method belonging to a class/struct (non-containment).
	EdgeTypeMethodOf = "method_of"
)
