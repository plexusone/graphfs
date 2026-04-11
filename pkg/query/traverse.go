// Package query provides graph traversal and search capabilities.
package query

import (
	"github.com/plexusone/graphfs/pkg/graph"
)

// Direction specifies the traversal direction.
type Direction int

const (
	// Outgoing follows edges from source to target (what does X call?)
	Outgoing Direction = iota
	// Incoming follows edges from target to source (what calls X?)
	Incoming
	// Both follows edges in both directions
	Both
)

// TraversalResult holds the result of a graph traversal.
type TraversalResult struct {
	// StartNode is the node where traversal began.
	StartNode string

	// Visited contains all visited node IDs in order.
	Visited []string

	// Edges contains all traversed edges.
	Edges []*graph.Edge

	// Depth maps node ID to its depth from start (0 = start node).
	Depth map[string]int

	// Parents maps node ID to the edge that led to it.
	Parents map[string]*graph.Edge
}

// Traverser performs graph traversals.
type Traverser struct {
	// Edges indexed by source node.
	outgoing map[string][]*graph.Edge
	// Edges indexed by target node.
	incoming map[string][]*graph.Edge
	// All nodes.
	nodes map[string]*graph.Node
}

// NewTraverser creates a traverser from a graph.
func NewTraverser(g *graph.Graph) *Traverser {
	t := &Traverser{
		outgoing: make(map[string][]*graph.Edge),
		incoming: make(map[string][]*graph.Edge),
		nodes:    g.Nodes,
	}

	for _, e := range g.Edges {
		t.outgoing[e.From] = append(t.outgoing[e.From], e)
		t.incoming[e.To] = append(t.incoming[e.To], e)
	}

	return t
}

// NewTraverserFromEdges creates a traverser from a slice of edges.
func NewTraverserFromEdges(edges []*graph.Edge, nodes map[string]*graph.Node) *Traverser {
	t := &Traverser{
		outgoing: make(map[string][]*graph.Edge),
		incoming: make(map[string][]*graph.Edge),
		nodes:    nodes,
	}

	for _, e := range edges {
		t.outgoing[e.From] = append(t.outgoing[e.From], e)
		t.incoming[e.To] = append(t.incoming[e.To], e)
	}

	return t
}

// BFS performs breadth-first search from a starting node.
func (t *Traverser) BFS(start string, dir Direction, maxDepth int, edgeTypes []string) *TraversalResult {
	result := &TraversalResult{
		StartNode: start,
		Visited:   []string{start},
		Edges:     make([]*graph.Edge, 0),
		Depth:     map[string]int{start: 0},
		Parents:   make(map[string]*graph.Edge),
	}

	if maxDepth == 0 {
		maxDepth = 100 // Default max depth
	}

	// Edge type filter
	typeFilter := make(map[string]bool)
	for _, et := range edgeTypes {
		typeFilter[et] = true
	}

	// BFS queue: (nodeID, depth)
	type queueItem struct {
		node  string
		depth int
	}
	queue := []queueItem{{start, 0}}
	visited := map[string]bool{start: true}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.depth >= maxDepth {
			continue
		}

		// Get edges based on direction
		var edges []*graph.Edge
		switch dir {
		case Outgoing:
			edges = t.outgoing[current.node]
		case Incoming:
			edges = t.incoming[current.node]
		case Both:
			edges = append(t.outgoing[current.node], t.incoming[current.node]...)
		}

		for _, e := range edges {
			// Apply edge type filter
			if len(typeFilter) > 0 && !typeFilter[e.Type] {
				continue
			}

			// Determine the neighbor node
			neighbor := e.To
			if dir == Incoming || (dir == Both && e.To == current.node) {
				neighbor = e.From
			}

			if visited[neighbor] {
				continue
			}

			visited[neighbor] = true
			result.Visited = append(result.Visited, neighbor)
			result.Edges = append(result.Edges, e)
			result.Depth[neighbor] = current.depth + 1
			result.Parents[neighbor] = e

			queue = append(queue, queueItem{neighbor, current.depth + 1})
		}
	}

	return result
}

// DFS performs depth-first search from a starting node.
func (t *Traverser) DFS(start string, dir Direction, maxDepth int, edgeTypes []string) *TraversalResult {
	result := &TraversalResult{
		StartNode: start,
		Visited:   []string{},
		Edges:     make([]*graph.Edge, 0),
		Depth:     make(map[string]int),
		Parents:   make(map[string]*graph.Edge),
	}

	if maxDepth == 0 {
		maxDepth = 100
	}

	// Edge type filter
	typeFilter := make(map[string]bool)
	for _, et := range edgeTypes {
		typeFilter[et] = true
	}

	visited := make(map[string]bool)

	var dfs func(node string, depth int)
	dfs = func(node string, depth int) {
		if visited[node] || depth > maxDepth {
			return
		}

		visited[node] = true
		result.Visited = append(result.Visited, node)
		result.Depth[node] = depth

		// Get edges based on direction
		var edges []*graph.Edge
		switch dir {
		case Outgoing:
			edges = t.outgoing[node]
		case Incoming:
			edges = t.incoming[node]
		case Both:
			edges = append(t.outgoing[node], t.incoming[node]...)
		}

		for _, e := range edges {
			// Apply edge type filter
			if len(typeFilter) > 0 && !typeFilter[e.Type] {
				continue
			}

			// Determine the neighbor node
			neighbor := e.To
			if dir == Incoming || (dir == Both && e.To == node) {
				neighbor = e.From
			}

			if !visited[neighbor] {
				result.Edges = append(result.Edges, e)
				result.Parents[neighbor] = e
				dfs(neighbor, depth+1)
			}
		}
	}

	dfs(start, 0)
	return result
}

// FindPath finds a path between two nodes using BFS.
func (t *Traverser) FindPath(from, to string, edgeTypes []string) *TraversalResult {
	result := &TraversalResult{
		StartNode: from,
		Visited:   []string{},
		Edges:     make([]*graph.Edge, 0),
		Depth:     make(map[string]int),
		Parents:   make(map[string]*graph.Edge),
	}

	// Edge type filter
	typeFilter := make(map[string]bool)
	for _, et := range edgeTypes {
		typeFilter[et] = true
	}

	// BFS to find shortest path
	type queueItem struct {
		node string
	}
	queue := []queueItem{{from}}
	visited := map[string]bool{from: true}
	parent := make(map[string]*graph.Edge)

	found := false
	for len(queue) > 0 && !found {
		current := queue[0]
		queue = queue[1:]

		// Check both directions for path finding
		edges := append(t.outgoing[current.node], t.incoming[current.node]...)

		for _, e := range edges {
			if len(typeFilter) > 0 && !typeFilter[e.Type] {
				continue
			}

			neighbor := e.To
			if e.To == current.node {
				neighbor = e.From
			}

			if visited[neighbor] {
				continue
			}

			visited[neighbor] = true
			parent[neighbor] = e

			if neighbor == to {
				found = true
				break
			}

			queue = append(queue, queueItem{neighbor})
		}
	}

	if !found {
		return result // Empty result, no path found
	}

	// Reconstruct path
	path := []string{to}
	pathEdges := []*graph.Edge{}
	current := to

	for current != from {
		e := parent[current]
		pathEdges = append([]*graph.Edge{e}, pathEdges...)
		if e.To == current {
			current = e.From
		} else {
			current = e.To
		}
		path = append([]string{current}, path...)
	}

	result.Visited = path
	result.Edges = pathEdges
	for i, n := range path {
		result.Depth[n] = i
	}

	return result
}
