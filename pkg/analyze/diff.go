// Package analyze provides graph analysis functions.
package analyze

import (
	"fmt"

	"github.com/plexusone/graphfs/pkg/graph"
)

// GraphDiff represents the difference between two graph snapshots.
type GraphDiff struct {
	NewNodes     []NodeChange `json:"new_nodes"`
	RemovedNodes []NodeChange `json:"removed_nodes"`
	NewEdges     []EdgeChange `json:"new_edges"`
	RemovedEdges []EdgeChange `json:"removed_edges"`
	Summary      string       `json:"summary"`
}

// NodeChange represents a node that was added or removed.
type NodeChange struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

// EdgeChange represents an edge that was added or removed.
type EdgeChange struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Type       string `json:"type"`
	Confidence string `json:"confidence"`
}

// DiffGraphs compares two graph snapshots and returns what changed.
func DiffGraphs(oldNodes, newNodes []*graph.Node, oldEdges, newEdges []*graph.Edge) *GraphDiff {
	// Build node sets
	oldNodeSet := make(map[string]*graph.Node)
	newNodeSet := make(map[string]*graph.Node)
	for _, n := range oldNodes {
		oldNodeSet[n.ID] = n
	}
	for _, n := range newNodes {
		newNodeSet[n.ID] = n
	}

	// Find added and removed nodes
	var addedNodes, removedNodes []NodeChange
	for id, n := range newNodeSet {
		if _, exists := oldNodeSet[id]; !exists {
			addedNodes = append(addedNodes, NodeChange{
				ID:    n.ID,
				Label: n.Label,
				Type:  n.Type,
			})
		}
	}
	for id, n := range oldNodeSet {
		if _, exists := newNodeSet[id]; !exists {
			removedNodes = append(removedNodes, NodeChange{
				ID:    n.ID,
				Label: n.Label,
				Type:  n.Type,
			})
		}
	}

	// Build edge sets using composite keys
	oldEdgeSet := make(map[string]*graph.Edge)
	newEdgeSet := make(map[string]*graph.Edge)
	for _, e := range oldEdges {
		key := edgeKey(e)
		oldEdgeSet[key] = e
	}
	for _, e := range newEdges {
		key := edgeKey(e)
		newEdgeSet[key] = e
	}

	// Find added and removed edges
	var addedEdges, removedEdges []EdgeChange
	for key, e := range newEdgeSet {
		if _, exists := oldEdgeSet[key]; !exists {
			addedEdges = append(addedEdges, EdgeChange{
				From:       e.From,
				To:         e.To,
				Type:       e.Type,
				Confidence: string(e.Confidence),
			})
		}
	}
	for key, e := range oldEdgeSet {
		if _, exists := newEdgeSet[key]; !exists {
			removedEdges = append(removedEdges, EdgeChange{
				From:       e.From,
				To:         e.To,
				Type:       e.Type,
				Confidence: string(e.Confidence),
			})
		}
	}

	return &GraphDiff{
		NewNodes:     addedNodes,
		RemovedNodes: removedNodes,
		NewEdges:     addedEdges,
		RemovedEdges: removedEdges,
		Summary:      formatDiffSummary(len(addedNodes), len(removedNodes), len(addedEdges), len(removedEdges)),
	}
}

func edgeKey(e *graph.Edge) string {
	// Use sorted endpoints for undirected comparison
	if e.From < e.To {
		return e.From + "|" + e.To + "|" + e.Type
	}
	return e.To + "|" + e.From + "|" + e.Type
}

func formatDiffSummary(newNodes, removedNodes, newEdges, removedEdges int) string {
	if newNodes == 0 && removedNodes == 0 && newEdges == 0 && removedEdges == 0 {
		return "No changes detected"
	}

	result := ""
	if newNodes > 0 {
		result += pluralize(newNodes, "new node", "new nodes")
	}
	if removedNodes > 0 {
		if result != "" {
			result += ", "
		}
		result += pluralize(removedNodes, "node removed", "nodes removed")
	}
	if newEdges > 0 {
		if result != "" {
			result += ", "
		}
		result += pluralize(newEdges, "new edge", "new edges")
	}
	if removedEdges > 0 {
		if result != "" {
			result += ", "
		}
		result += pluralize(removedEdges, "edge removed", "edges removed")
	}
	return result
}

func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return "1 " + singular
	}
	return fmt.Sprintf("%d %s", n, plural)
}
