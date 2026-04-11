package analyze

import (
	"sort"

	"github.com/plexusone/graphfs/pkg/graph"
)

// Community represents a detected community in the graph.
type Community struct {
	ID       int      `json:"id"`
	Size     int      `json:"size"`
	Cohesion float64  `json:"cohesion"`
	Members  []string `json:"members"`
	Label    string   `json:"label,omitempty"`
}

// ClusterResult contains the results of community detection.
type ClusterResult struct {
	Communities   []Community    `json:"communities"`
	NodeCommunity map[string]int `json:"node_community"`
	Modularity    float64        `json:"modularity,omitempty"`
}

// ClusterAlgorithm specifies the community detection algorithm to use.
type ClusterAlgorithm string

const (
	// AlgorithmLouvain uses the Louvain modularity optimization algorithm.
	AlgorithmLouvain ClusterAlgorithm = "louvain"

	// AlgorithmConnectedComponents uses connected components.
	AlgorithmConnectedComponents ClusterAlgorithm = "components"
)

// ClusterOptions configures community detection.
type ClusterOptions struct {
	// Algorithm specifies which algorithm to use.
	// Default: AlgorithmLouvain
	Algorithm ClusterAlgorithm

	// Resolution for Louvain algorithm (higher = smaller communities).
	// Default: 1.0
	Resolution float64

	// ExcludeEdgeTypes lists edge types to exclude from community detection.
	ExcludeEdgeTypes []string

	// ExcludeNodeTypes lists node types to exclude from community detection.
	ExcludeNodeTypes []string
}

// DefaultClusterOptions returns sensible defaults.
func DefaultClusterOptions() ClusterOptions {
	return ClusterOptions{
		Algorithm:        AlgorithmLouvain,
		Resolution:       1.0,
		ExcludeEdgeTypes: []string{"contains", "imports"},
		ExcludeNodeTypes: []string{"package", "file"},
	}
}

// DetectCommunities performs community detection using the Louvain algorithm by default.
func DetectCommunities(nodes []*graph.Node, edges []*graph.Edge) *ClusterResult {
	return DetectCommunitiesWithOptions(nodes, edges, DefaultClusterOptions())
}

// DetectCommunitiesWithOptions performs community detection with configurable algorithm.
func DetectCommunitiesWithOptions(nodes []*graph.Node, edges []*graph.Edge, opts ClusterOptions) *ClusterResult {
	switch opts.Algorithm {
	case AlgorithmLouvain:
		louvainOpts := DefaultLouvainOptions()
		louvainOpts.Resolution = opts.Resolution
		louvainOpts.ExcludeEdgeTypes = opts.ExcludeEdgeTypes
		louvainOpts.ExcludeNodeTypes = opts.ExcludeNodeTypes
		result := DetectCommunitiesLouvain(nodes, edges, louvainOpts)
		return LouvainToClusters(result, nodes, edges)

	case AlgorithmConnectedComponents:
		adj := buildAdjacencyList(nodes, edges, opts.ExcludeEdgeTypes, opts.ExcludeNodeTypes)
		components := connectedComponents(nodes, adj, opts.ExcludeNodeTypes)
		return buildClusterResult(nodes, edges, components)

	default:
		// Fall back to Louvain
		louvainOpts := DefaultLouvainOptions()
		result := DetectCommunitiesLouvain(nodes, edges, louvainOpts)
		return LouvainToClusters(result, nodes, edges)
	}
}

// buildAdjacencyList creates an undirected adjacency list from edges.
func buildAdjacencyList(nodes []*graph.Node, edges []*graph.Edge, excludeEdgeTypes, excludeNodeTypes []string) map[string][]string {
	excludeEdges := make(map[string]bool)
	for _, t := range excludeEdgeTypes {
		excludeEdges[t] = true
	}

	excludeNodes := make(map[string]bool)
	for _, t := range excludeNodeTypes {
		excludeNodes[t] = true
	}

	// Build node type map
	nodeType := make(map[string]string)
	for _, n := range nodes {
		nodeType[n.ID] = n.Type
	}

	adj := make(map[string][]string)
	for _, e := range edges {
		if excludeEdges[e.Type] {
			continue
		}
		if excludeNodes[nodeType[e.From]] || excludeNodes[nodeType[e.To]] {
			continue
		}
		adj[e.From] = append(adj[e.From], e.To)
		adj[e.To] = append(adj[e.To], e.From)
	}
	return adj
}

// connectedComponents finds connected components using DFS.
func connectedComponents(nodes []*graph.Node, adj map[string][]string, excludeNodeTypes []string) map[int][]string {
	excludeNodes := make(map[string]bool)
	for _, t := range excludeNodeTypes {
		excludeNodes[t] = true
	}

	visited := make(map[string]bool)
	components := make(map[int][]string)
	componentID := 0

	var dfs func(nodeID string, members *[]string)
	dfs = func(nodeID string, members *[]string) {
		if visited[nodeID] {
			return
		}
		visited[nodeID] = true
		*members = append(*members, nodeID)
		for _, neighbor := range adj[nodeID] {
			dfs(neighbor, members)
		}
	}

	for _, n := range nodes {
		if excludeNodes[n.Type] {
			continue
		}
		if n.Attrs != nil && n.Attrs["external"] == "true" {
			continue
		}
		if !visited[n.ID] {
			var members []string
			dfs(n.ID, &members)
			if len(members) > 0 {
				components[componentID] = members
				componentID++
			}
		}
	}

	return components
}

// buildClusterResult creates a ClusterResult from community membership.
func buildClusterResult(_ []*graph.Node, edges []*graph.Edge, communities map[int][]string) *ClusterResult {
	// Build edge adjacency for cohesion calculation
	adj := make(map[string]map[string]bool)
	for _, e := range edges {
		if adj[e.From] == nil {
			adj[e.From] = make(map[string]bool)
		}
		if adj[e.To] == nil {
			adj[e.To] = make(map[string]bool)
		}
		adj[e.From][e.To] = true
		adj[e.To][e.From] = true
	}

	// Build node community map
	nodeCommunity := make(map[string]int)
	for cid, members := range communities {
		for _, nodeID := range members {
			nodeCommunity[nodeID] = cid
		}
	}

	// Build community list with cohesion scores
	var communityList []Community
	for cid, members := range communities {
		cohesion := CohesionScore(members, adj)
		communityList = append(communityList, Community{
			ID:       cid,
			Size:     len(members),
			Cohesion: cohesion,
			Members:  members,
		})
	}

	// Sort by size descending
	sort.Slice(communityList, func(i, j int) bool {
		return communityList[i].Size > communityList[j].Size
	})

	// Re-number communities by size
	for i := range communityList {
		oldID := communityList[i].ID
		communityList[i].ID = i
		for _, member := range communityList[i].Members {
			if nodeCommunity[member] == oldID {
				nodeCommunity[member] = i
			}
		}
	}

	return &ClusterResult{
		Communities:   communityList,
		NodeCommunity: nodeCommunity,
	}
}

// CohesionScore calculates the ratio of actual intra-community edges to maximum possible.
func CohesionScore(members []string, adj map[string]map[string]bool) float64 {
	n := len(members)
	if n <= 1 {
		return 1.0
	}

	// Count actual edges within community
	memberSet := make(map[string]bool)
	for _, m := range members {
		memberSet[m] = true
	}

	actual := 0
	for _, m := range members {
		for neighbor := range adj[m] {
			if memberSet[neighbor] {
				actual++
			}
		}
	}
	// Divide by 2 since we count each edge twice
	actual = actual / 2

	// Maximum possible edges in a clique
	possible := n * (n - 1) / 2

	if possible == 0 {
		return 0.0
	}

	cohesion := float64(actual) / float64(possible)
	// Round to 2 decimal places
	return float64(int(cohesion*100)) / 100
}
