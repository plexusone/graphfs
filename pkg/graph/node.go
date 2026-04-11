// Package graph provides core types for the GraphFS graph database.
package graph

// Node represents an entity in the graph.
type Node struct {
	// ID is the unique, stable identifier for this node.
	// Should be deterministic (e.g., based on content hash or path+symbol).
	ID string `json:"id"`

	// Type categorizes the node (e.g., "function", "file", "class", "module").
	Type string `json:"type"`

	// Label is a human-readable name for display.
	Label string `json:"label,omitempty"`

	// Attrs holds additional attributes as key-value pairs.
	Attrs map[string]string `json:"attrs,omitempty"`
}

// NodeType constants for common node types.
const (
	NodeTypeFunction = "function"
	NodeTypeMethod   = "method"
	NodeTypeClass    = "class"
	NodeTypeStruct   = "struct"
	NodeTypeFile     = "file"
	NodeTypePackage  = "package"
	NodeTypeModule   = "module"
	NodeTypeVariable = "variable"
	NodeTypeConstant = "constant"
	NodeTypeInterface = "interface"
)
