// Package store provides filesystem-backed persistence for graphs.
package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/plexusone/graphfs/pkg/format"
	"github.com/plexusone/graphfs/pkg/graph"
)

// FSStore implements Store using the filesystem.
// It stores one file per node and one file per edge for git-friendly diffs.
type FSStore struct {
	// Root is the base directory for the graph store.
	Root string
}

// NewFSStore creates a new filesystem-backed store.
func NewFSStore(root string) (*FSStore, error) {
	// Create directory structure
	dirs := []string{
		filepath.Join(root, "nodes"),
		filepath.Join(root, "edges"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}
	return &FSStore{Root: root}, nil
}

// WriteNode writes a node to the filesystem.
func (s *FSStore) WriteNode(n *graph.Node) error {
	path := s.nodePath(n.ID)
	data, err := format.MarshalCanonical(n)
	if err != nil {
		return fmt.Errorf("marshaling node %s: %w", n.ID, err)
	}
	return os.WriteFile(path, data, 0600)
}

// WriteEdge writes an edge to the filesystem.
func (s *FSStore) WriteEdge(e *graph.Edge) error {
	path := s.edgePath(e.From, e.Type, e.To)
	data, err := format.MarshalCanonical(e)
	if err != nil {
		return fmt.Errorf("marshaling edge %s->%s: %w", e.From, e.To, err)
	}
	return os.WriteFile(path, data, 0600)
}

// GetNode retrieves a node by ID.
func (s *FSStore) GetNode(id string) (*graph.Node, error) {
	path := s.nodePath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("node %s not found", id)
		}
		return nil, fmt.Errorf("reading node %s: %w", id, err)
	}
	var n graph.Node
	if err := format.UnmarshalCanonical(data, &n); err != nil {
		return nil, fmt.Errorf("unmarshaling node %s: %w", id, err)
	}
	return &n, nil
}

// GetEdge retrieves an edge by its composite key.
func (s *FSStore) GetEdge(from, edgeType, to string) (*graph.Edge, error) {
	path := s.edgePath(from, edgeType, to)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("edge %s--%s-->%s not found", from, edgeType, to)
		}
		return nil, fmt.Errorf("reading edge: %w", err)
	}
	var e graph.Edge
	if err := format.UnmarshalCanonical(data, &e); err != nil {
		return nil, fmt.Errorf("unmarshaling edge: %w", err)
	}
	return &e, nil
}

// ListNodes returns all nodes in the store.
func (s *FSStore) ListNodes() ([]*graph.Node, error) {
	dir := filepath.Join(s.Root, "nodes")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading nodes directory: %w", err)
	}

	var nodes []*graph.Node
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".json")
		n, err := s.GetNode(id)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

// ListEdges returns all edges in the store.
func (s *FSStore) ListEdges() ([]*graph.Edge, error) {
	dir := filepath.Join(s.Root, "edges")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading edges directory: %w", err)
	}

	var edges []*graph.Edge
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading edge file %s: %w", entry.Name(), err)
		}
		var e graph.Edge
		if err := format.UnmarshalCanonical(data, &e); err != nil {
			return nil, fmt.Errorf("unmarshaling edge file %s: %w", entry.Name(), err)
		}
		edges = append(edges, &e)
	}
	return edges, nil
}

// DeleteNode removes a node by ID.
func (s *FSStore) DeleteNode(id string) error {
	path := s.nodePath(id)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("deleting node %s: %w", id, err)
	}
	return nil
}

// DeleteEdge removes an edge.
func (s *FSStore) DeleteEdge(from, edgeType, to string) error {
	path := s.edgePath(from, edgeType, to)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("deleting edge: %w", err)
	}
	return nil
}

// LoadGraph loads the entire graph from storage.
func (s *FSStore) LoadGraph() (*graph.Graph, error) {
	g := graph.NewGraph()

	nodes, err := s.ListNodes()
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		g.AddNode(n)
	}

	edges, err := s.ListEdges()
	if err != nil {
		return nil, err
	}
	for _, e := range edges {
		g.AddEdge(e)
	}

	return g, nil
}

// SaveGraph saves an entire graph to storage.
func (s *FSStore) SaveGraph(g *graph.Graph) error {
	for _, n := range g.Nodes {
		if err := s.WriteNode(n); err != nil {
			return err
		}
	}
	for _, e := range g.Edges {
		if err := s.WriteEdge(e); err != nil {
			return err
		}
	}
	return nil
}

// nodePath returns the file path for a node.
func (s *FSStore) nodePath(id string) string {
	return filepath.Join(s.Root, "nodes", id+".json")
}

// edgePath returns the file path for an edge.
// Format: from__type__to.json
func (s *FSStore) edgePath(from, edgeType, to string) string {
	filename := fmt.Sprintf("%s__%s__%s.json", from, edgeType, to)
	return filepath.Join(s.Root, "edges", filename)
}
