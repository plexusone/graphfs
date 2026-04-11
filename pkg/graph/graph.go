// Package graph provides core types for the GraphFS graph database.
package graph

// Graph represents a complete graph with nodes and edges.
type Graph struct {
	// Nodes maps node ID to Node.
	Nodes map[string]*Node `json:"nodes"`

	// Edges holds all edges in the graph.
	Edges []*Edge `json:"edges"`
}

// NewGraph creates an empty graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make([]*Edge, 0),
	}
}

// AddNode adds a node to the graph. If a node with the same ID exists,
// it will be overwritten.
func (g *Graph) AddNode(n *Node) {
	g.Nodes[n.ID] = n
}

// AddEdge adds an edge to the graph.
func (g *Graph) AddEdge(e *Edge) {
	g.Edges = append(g.Edges, e)
}

// GetNode returns a node by ID, or nil if not found.
func (g *Graph) GetNode(id string) *Node {
	return g.Nodes[id]
}

// NodeCount returns the number of nodes in the graph.
func (g *Graph) NodeCount() int {
	return len(g.Nodes)
}

// EdgeCount returns the number of edges in the graph.
func (g *Graph) EdgeCount() int {
	return len(g.Edges)
}
