// Package analyze provides graph analysis features including hub detection,
// community detection, and graph comparison.
package analyze

import (
	"sort"

	"github.com/plexusone/graphfs/pkg/graph"
)

// HubNode represents a highly connected node in the graph.
type HubNode struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Type      string `json:"type"`
	InDegree  int    `json:"in_degree"`
	OutDegree int    `json:"out_degree"`
	Total     int    `json:"total"`
}

// FindHubs returns the top N most connected nodes in the graph.
// Use excludeTypes to filter out structural hub nodes (e.g., packages, files).
func FindHubs(nodes []*graph.Node, edges []*graph.Edge, topN int, excludeTypes []string) []HubNode {
	// Build exclusion set
	excludeSet := make(map[string]bool)
	for _, t := range excludeTypes {
		excludeSet[t] = true
	}

	// Build adjacency counts
	inDegree := make(map[string]int)
	outDegree := make(map[string]int)

	for _, e := range edges {
		outDegree[e.From]++
		inDegree[e.To]++
	}

	// Build node map for labels and types
	nodeMap := make(map[string]*graph.Node)
	for _, n := range nodes {
		nodeMap[n.ID] = n
	}

	// Calculate totals and filter
	type nodeDegree struct {
		id    string
		node  *graph.Node
		in    int
		out   int
		total int
	}

	var candidates []nodeDegree
	for id, node := range nodeMap {
		// Skip excluded types
		if excludeSet[node.Type] {
			continue
		}
		// Skip external nodes
		if node.Attrs != nil && node.Attrs["external"] == "true" {
			continue
		}

		in := inDegree[id]
		out := outDegree[id]
		total := in + out

		if total > 0 {
			candidates = append(candidates, nodeDegree{
				id:    id,
				node:  node,
				in:    in,
				out:   out,
				total: total,
			})
		}
	}

	// Sort by total degree descending
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].total > candidates[j].total
	})

	// Take top N
	if topN > len(candidates) {
		topN = len(candidates)
	}

	result := make([]HubNode, topN)
	for i := 0; i < topN; i++ {
		c := candidates[i]
		label := c.node.Label
		if label == "" {
			label = c.id
		}
		result[i] = HubNode{
			ID:        c.id,
			Label:     label,
			Type:      c.node.Type,
			InDegree:  c.in,
			OutDegree: c.out,
			Total:     c.total,
		}
	}

	return result
}

// IsolatedNodes returns nodes with degree <= threshold.
// Use excludeTypes to filter out structural nodes.
func IsolatedNodes(nodes []*graph.Node, edges []*graph.Edge, threshold int, excludeTypes []string) []*graph.Node {
	// Build exclusion set
	excludeSet := make(map[string]bool)
	for _, t := range excludeTypes {
		excludeSet[t] = true
	}

	// Build degree map
	degree := make(map[string]int)
	for _, e := range edges {
		degree[e.From]++
		degree[e.To]++
	}

	var isolated []*graph.Node
	for _, n := range nodes {
		if excludeSet[n.Type] {
			continue
		}
		if n.Attrs != nil && n.Attrs["external"] == "true" {
			continue
		}
		if degree[n.ID] <= threshold {
			isolated = append(isolated, n)
		}
	}

	return isolated
}

// NodesByType groups nodes by their type for analysis.
func NodesByType(nodes []*graph.Node) map[string][]*graph.Node {
	byType := make(map[string][]*graph.Node)
	for _, n := range nodes {
		byType[n.Type] = append(byType[n.Type], n)
	}
	return byType
}

// EdgesByType groups edges by their type for analysis.
func EdgesByType(edges []*graph.Edge) map[string][]*graph.Edge {
	byType := make(map[string][]*graph.Edge)
	for _, e := range edges {
		byType[e.Type] = append(byType[e.Type], e)
	}
	return byType
}

// EdgesByConfidence groups edges by their confidence level.
func EdgesByConfidence(edges []*graph.Edge) map[graph.Confidence][]*graph.Edge {
	byConf := make(map[graph.Confidence][]*graph.Edge)
	for _, e := range edges {
		conf := e.Confidence
		if conf == "" {
			conf = graph.ConfidenceExtracted
		}
		byConf[conf] = append(byConf[conf], e)
	}
	return byConf
}

// HubScore calculates a simple hub score based on out-degree.
// High out-degree nodes are potential architectural hubs.
func HubScore(nodeID string, edges []*graph.Edge) int {
	score := 0
	for _, e := range edges {
		if e.From == nodeID {
			score++
		}
	}
	return score
}

// AuthorityScore calculates a simple authority score based on in-degree.
// High in-degree nodes are frequently referenced entities.
func AuthorityScore(nodeID string, edges []*graph.Edge) int {
	score := 0
	for _, e := range edges {
		if e.To == nodeID {
			score++
		}
	}
	return score
}

// InferredEdges returns edges with INFERRED or AMBIGUOUS confidence.
func InferredEdges(edges []*graph.Edge) []*graph.Edge {
	var inferred []*graph.Edge
	for _, e := range edges {
		if e.Confidence == graph.ConfidenceInferred || e.Confidence == graph.ConfidenceAmbiguous {
			inferred = append(inferred, e)
		}
	}
	return inferred
}
