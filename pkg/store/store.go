// Package store provides filesystem-backed persistence for graphs.
package store

import (
	"github.com/plexusone/graphfs/pkg/graph"
)

// Store defines the interface for graph persistence.
type Store interface {
	// WriteNode writes a node to storage.
	WriteNode(n *graph.Node) error

	// WriteEdge writes an edge to storage.
	WriteEdge(e *graph.Edge) error

	// GetNode retrieves a node by ID.
	GetNode(id string) (*graph.Node, error)

	// GetEdge retrieves an edge by its composite key (from__type__to).
	GetEdge(from, edgeType, to string) (*graph.Edge, error)

	// ListNodes returns all nodes in the store.
	ListNodes() ([]*graph.Node, error)

	// ListEdges returns all edges in the store.
	ListEdges() ([]*graph.Edge, error)

	// DeleteNode removes a node by ID.
	DeleteNode(id string) error

	// DeleteEdge removes an edge.
	DeleteEdge(from, edgeType, to string) error

	// LoadGraph loads the entire graph from storage.
	LoadGraph() (*graph.Graph, error)

	// SaveGraph saves an entire graph to storage.
	SaveGraph(g *graph.Graph) error
}
