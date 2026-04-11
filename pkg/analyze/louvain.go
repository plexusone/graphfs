package analyze

import (
	"math/rand/v2"

	"github.com/plexusone/graphfs/pkg/graph"
	gonumgraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/community"
	"gonum.org/v1/gonum/graph/simple"
)

// LouvainOptions configures the Louvain algorithm.
type LouvainOptions struct {
	// Resolution controls community granularity. Higher values = smaller communities.
	// Default: 1.0 (standard modularity)
	Resolution float64

	// Seed for random number generator. Use 0 for non-deterministic.
	Seed uint64

	// ExcludeEdgeTypes lists edge types to exclude from community detection.
	// Default: ["contains", "imports"] (structural edges)
	ExcludeEdgeTypes []string

	// ExcludeNodeTypes lists node types to exclude from community detection.
	// Default: ["package", "file"] (hub nodes)
	ExcludeNodeTypes []string
}

// DefaultLouvainOptions returns sensible defaults for Louvain.
func DefaultLouvainOptions() LouvainOptions {
	return LouvainOptions{
		Resolution:       1.0,
		Seed:             42, // Deterministic by default for reproducibility
		ExcludeEdgeTypes: []string{"contains", "imports"},
		ExcludeNodeTypes: []string{"package", "file"},
	}
}

// LouvainResult contains the results of Louvain community detection.
type LouvainResult struct {
	// Communities maps community ID to member node IDs.
	Communities map[int][]string

	// NodeCommunity maps node ID to community ID.
	NodeCommunity map[string]int

	// Modularity is the Q score of the detected communities.
	Modularity float64

	// NumLevels is the number of hierarchical levels in the dendrogram.
	NumLevels int
}

// DetectCommunitiesLouvain performs community detection using the Louvain algorithm.
func DetectCommunitiesLouvain(nodes []*graph.Node, edges []*graph.Edge, opts LouvainOptions) *LouvainResult {
	// Build exclusion sets
	excludeEdgeTypes := make(map[string]bool)
	for _, t := range opts.ExcludeEdgeTypes {
		excludeEdgeTypes[t] = true
	}

	excludeNodeTypes := make(map[string]bool)
	for _, t := range opts.ExcludeNodeTypes {
		excludeNodeTypes[t] = true
	}

	// Filter nodes
	var filteredNodes []*graph.Node
	nodeIDSet := make(map[string]bool)
	for _, n := range nodes {
		if excludeNodeTypes[n.Type] {
			continue
		}
		if n.Attrs != nil && n.Attrs["external"] == "true" {
			continue
		}
		filteredNodes = append(filteredNodes, n)
		nodeIDSet[n.ID] = true
	}

	// Filter edges (only keep edges between included nodes)
	var filteredEdges []*graph.Edge
	for _, e := range edges {
		if excludeEdgeTypes[e.Type] {
			continue
		}
		if !nodeIDSet[e.From] || !nodeIDSet[e.To] {
			continue
		}
		filteredEdges = append(filteredEdges, e)
	}

	// Handle empty or trivial graphs
	if len(filteredNodes) == 0 {
		return &LouvainResult{
			Communities:   make(map[int][]string),
			NodeCommunity: make(map[string]int),
			Modularity:    0,
		}
	}

	if len(filteredEdges) == 0 {
		// No edges = each node is its own community
		communities := make(map[int][]string)
		nodeCommunity := make(map[string]int)
		for i, n := range filteredNodes {
			communities[i] = []string{n.ID}
			nodeCommunity[n.ID] = i
		}
		return &LouvainResult{
			Communities:   communities,
			NodeCommunity: nodeCommunity,
			Modularity:    0,
		}
	}

	// Build gonum graph
	g, _, intToNodeID := buildGonumGraph(filteredNodes, filteredEdges)

	// Create random source
	var src rand.Source
	if opts.Seed != 0 {
		src = rand.NewPCG(opts.Seed, opts.Seed)
	}

	// Run Louvain algorithm
	reduced := community.Modularize(g, opts.Resolution, src)

	// Extract communities
	gonumCommunities := reduced.Communities()

	// Note: We skip counting hierarchical levels as it requires type assertion
	// and the number of levels is not critical for our use case.
	numLevels := 1

	// Convert gonum communities to our format
	communities := make(map[int][]string)
	nodeCommunity := make(map[string]int)

	for cid, members := range gonumCommunities {
		for _, gonumNode := range members {
			nodeID := intToNodeID[gonumNode.ID()]
			communities[cid] = append(communities[cid], nodeID)
			nodeCommunity[nodeID] = cid
		}
	}

	// Calculate modularity
	modularity := community.Q(g, gonumCommunities, opts.Resolution)

	return &LouvainResult{
		Communities:   communities,
		NodeCommunity: nodeCommunity,
		Modularity:    modularity,
		NumLevels:     numLevels,
	}
}

// buildGonumGraph converts our graph to gonum's undirected weighted graph.
func buildGonumGraph(nodes []*graph.Node, edges []*graph.Edge) (*simple.WeightedUndirectedGraph, map[string]int64, map[int64]string) {
	g := simple.NewWeightedUndirectedGraph(0, 0)

	// Create node ID mappings
	nodeIDToInt := make(map[string]int64)
	intToNodeID := make(map[int64]string)

	for i, n := range nodes {
		id := int64(i)
		nodeIDToInt[n.ID] = id
		intToNodeID[id] = n.ID
		g.AddNode(simple.Node(id))
	}

	// Add edges (combine weights for parallel edges)
	edgeWeights := make(map[[2]int64]float64)
	for _, e := range edges {
		fromID := nodeIDToInt[e.From]
		toID := nodeIDToInt[e.To]

		// Skip self-loops (gonum doesn't allow them)
		if fromID == toID {
			continue
		}

		// Normalize edge direction for undirected graph
		key := [2]int64{fromID, toID}
		if fromID > toID {
			key = [2]int64{toID, fromID}
		}

		// Accumulate weight (each edge adds 1.0)
		edgeWeights[key] += 1.0
	}

	// Add weighted edges to graph
	for key, weight := range edgeWeights {
		g.SetWeightedEdge(simple.WeightedEdge{
			F: simple.Node(key[0]),
			T: simple.Node(key[1]),
			W: weight,
		})
	}

	return g, nodeIDToInt, intToNodeID
}

// LouvainToClusters converts LouvainResult to ClusterResult for compatibility.
func LouvainToClusters(result *LouvainResult, nodes []*graph.Node, edges []*graph.Edge) *ClusterResult {
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

	// Build community list with cohesion scores
	var communityList []Community
	for cid, members := range result.Communities {
		cohesion := CohesionScore(members, adj)
		communityList = append(communityList, Community{
			ID:       cid,
			Size:     len(members),
			Cohesion: cohesion,
			Members:  members,
		})
	}

	// Sort by size descending and re-number
	sortAndRenumberCommunities(communityList, result.NodeCommunity)

	return &ClusterResult{
		Communities:   communityList,
		NodeCommunity: result.NodeCommunity,
		Modularity:    result.Modularity,
	}
}

// sortAndRenumberCommunities sorts communities by size and updates IDs.
func sortAndRenumberCommunities(communities []Community, nodeCommunity map[string]int) {
	// Create old ID to new ID mapping
	type idPair struct {
		oldID int
		size  int
	}
	pairs := make([]idPair, len(communities))
	for i, c := range communities {
		pairs[i] = idPair{oldID: c.ID, size: c.Size}
	}

	// Sort by size descending
	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].size > pairs[i].size {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	// Build old to new ID mapping
	oldToNew := make(map[int]int)
	for newID, p := range pairs {
		oldToNew[p.oldID] = newID
	}

	// Update community IDs
	for i := range communities {
		communities[i].ID = oldToNew[communities[i].ID]
	}

	// Sort communities slice by new ID
	for i := 0; i < len(communities)-1; i++ {
		for j := i + 1; j < len(communities); j++ {
			if communities[j].ID < communities[i].ID {
				communities[i], communities[j] = communities[j], communities[i]
			}
		}
	}

	// Update node community mapping
	for nodeID, oldCID := range nodeCommunity {
		nodeCommunity[nodeID] = oldToNew[oldCID]
	}
}

// Ensure simple.Node implements gonumgraph.Node
var _ gonumgraph.Node = simple.Node(0)
