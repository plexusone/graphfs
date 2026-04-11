# API Reference

## Package `graph`

Core types for representing graph data.

### Node

```go
type Node struct {
    ID    string            `json:"id"`
    Type  string            `json:"type"`
    Label string            `json:"label,omitempty"`
    Attrs map[string]string `json:"attrs,omitempty"`
}
```

| Field | Description |
|-------|-------------|
| `ID` | Unique, stable identifier (filesystem-safe) |
| `Type` | Node category from `NodeType*` constants |
| `Label` | Human-readable display name |
| `Attrs` | Extensible key-value metadata |

### Edge

```go
type Edge struct {
    From            string            `json:"from"`
    To              string            `json:"to"`
    Type            string            `json:"type"`
    Confidence      Confidence        `json:"confidence"`
    ConfidenceScore float64           `json:"confidence_score,omitempty"`
    Attrs           map[string]string `json:"attrs,omitempty"`
}
```

| Field | Description |
|-------|-------------|
| `From` | Source node ID |
| `To` | Target node ID |
| `Type` | Relationship type from `EdgeType*` constants |
| `Confidence` | How the relationship was determined |
| `ConfidenceScore` | 0.0-1.0 score for INFERRED edges |
| `Attrs` | Extensible key-value metadata |

### Graph

```go
type Graph struct {
    Nodes map[string]*Node `json:"nodes"`
    Edges []*Edge          `json:"edges"`
}
```

#### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `NewGraph` | `func NewGraph() *Graph` | Create empty graph |
| `AddNode` | `func (g *Graph) AddNode(n *Node)` | Add or replace node |
| `AddEdge` | `func (g *Graph) AddEdge(e *Edge)` | Append edge |
| `GetNode` | `func (g *Graph) GetNode(id string) *Node` | Get node by ID |
| `NodeCount` | `func (g *Graph) NodeCount() int` | Count nodes |
| `EdgeCount` | `func (g *Graph) EdgeCount() int` | Count edges |

### Confidence

```go
type Confidence string

const (
    ConfidenceExtracted Confidence = "EXTRACTED"
    ConfidenceInferred  Confidence = "INFERRED"
    ConfidenceAmbiguous Confidence = "AMBIGUOUS"
)
```

| Value | Description |
|-------|-------------|
| `EXTRACTED` | Directly extracted from source (AST, imports) |
| `INFERRED` | Inferred by LLM or heuristic with confidence score |
| `AMBIGUOUS` | Uncertain relationship requiring human review |

### Node Type Constants

```go
const (
    NodeTypeFunction  = "function"
    NodeTypeMethod    = "method"
    NodeTypeClass     = "class"
    NodeTypeStruct    = "struct"
    NodeTypeFile      = "file"
    NodeTypePackage   = "package"
    NodeTypeModule    = "module"
    NodeTypeVariable  = "variable"
    NodeTypeConstant  = "constant"
    NodeTypeInterface = "interface"
)
```

### Edge Type Constants

```go
const (
    EdgeTypeCalls      = "calls"
    EdgeTypeImports    = "imports"
    EdgeTypeImplements = "implements"
    EdgeTypeExtends    = "extends"
    EdgeTypeUses       = "uses"
    EdgeTypeContains   = "contains"
    EdgeTypeDependsOn  = "depends_on"
    EdgeTypeReferences = "references"
)
```

---

## Package `store`

Filesystem-backed persistence for graphs.

### Store Interface

```go
type Store interface {
    WriteNode(n *graph.Node) error
    WriteEdge(e *graph.Edge) error
    GetNode(id string) (*graph.Node, error)
    GetEdge(from, edgeType, to string) (*graph.Edge, error)
    ListNodes() ([]*graph.Node, error)
    ListEdges() ([]*graph.Edge, error)
    DeleteNode(id string) error
    DeleteEdge(from, edgeType, to string) error
    LoadGraph() (*graph.Graph, error)
    SaveGraph(g *graph.Graph) error
}
```

### FSStore

```go
type FSStore struct {
    Root string
}

func NewFSStore(root string) (*FSStore, error)
```

Creates a filesystem-backed store at the given path. Automatically creates `nodes/` and `edges/` subdirectories.

#### File Layout

| Entity | Path |
|--------|------|
| Node | `{root}/nodes/{id}.json` |
| Edge | `{root}/edges/{from}__{type}__{to}.json` |

---

## Package `format`

Deterministic JSON serialization.

### MarshalCanonical

```go
func MarshalCanonical(v any) ([]byte, error)
```

Produces deterministic JSON output with:

- Alphabetically sorted keys
- 2-space indentation
- No trailing newline
- No HTML escaping

### UnmarshalCanonical

```go
func UnmarshalCanonical(data []byte, v any) error
```

Standard JSON unmarshaling (wrapper around `json.Unmarshal`).

---

## Package `schema`

Graph validation.

### Validator

```go
type Validator struct {
    AllowedNodeTypes []string
    AllowedEdgeTypes []string
    RequireNodeLabel bool
}

func NewValidator() *Validator
```

| Field | Description |
|-------|-------------|
| `AllowedNodeTypes` | Restrict node types (empty = allow all) |
| `AllowedEdgeTypes` | Restrict edge types (empty = allow all) |
| `RequireNodeLabel` | Require non-empty label on nodes |

### Validation Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `ValidateNode` | `func (v *Validator) ValidateNode(n *graph.Node) error` | Validate single node |
| `ValidateEdge` | `func (v *Validator) ValidateEdge(e *graph.Edge) error` | Validate single edge |
| `ValidateEdgeRefs` | `func (v *Validator) ValidateEdgeRefs(e *graph.Edge, nodes map[string]*graph.Node) error` | Check referential integrity |
| `ValidateGraph` | `func (v *Validator) ValidateGraph(g *graph.Graph) []error` | Validate entire graph |

### ValidationError

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string
```

Structured validation error with field name and message.

### Validation Rules

**Node validation:**

- `id` required, filesystem-safe (no `/\:*?"<>|`)
- `type` required
- `label` required if `RequireNodeLabel` is true
- `type` must be in `AllowedNodeTypes` if set

**Edge validation:**

- `from` required
- `to` required
- `type` required, must be in `AllowedEdgeTypes` if set
- `confidence` required, must be EXTRACTED, INFERRED, or AMBIGUOUS
- `confidence_score` must be 0.0-1.0 for INFERRED edges

**Referential integrity:**

- `from` node must exist in graph
- `to` node must exist in graph

---

## Package `query`

Graph traversal algorithms.

### Traverser

```go
type Traverser struct {
    // private fields
}

func NewTraverser(g *graph.Graph) *Traverser
func NewTraverserFromEdges(edges []*graph.Edge, nodes map[string]*graph.Node) *Traverser
```

Creates a traverser for graph exploration.

### Direction

```go
type Direction int

const (
    Outgoing Direction = iota  // Follow edges from source to target
    Incoming                   // Follow edges from target to source
    Both                       // Follow edges in both directions
)
```

### TraversalResult

```go
type TraversalResult struct {
    StartNode string              // Node where traversal began
    Visited   []string            // Visited node IDs in order
    Edges     []*graph.Edge       // Traversed edges
    Depth     map[string]int      // Node ID to depth from start
    Parents   map[string]*graph.Edge // Node ID to edge that led to it
}
```

### Traversal Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `BFS` | `func (t *Traverser) BFS(start string, dir Direction, maxDepth int, edgeTypes []string) *TraversalResult` | Breadth-first search |
| `DFS` | `func (t *Traverser) DFS(start string, dir Direction, maxDepth int, edgeTypes []string) *TraversalResult` | Depth-first search |
| `FindPath` | `func (t *Traverser) FindPath(from, to string, edgeTypes []string) *TraversalResult` | Find shortest path |

**Parameters:**

| Parameter | Description |
|-----------|-------------|
| `start` | Starting node ID |
| `dir` | Traversal direction (Outgoing, Incoming, Both) |
| `maxDepth` | Maximum traversal depth (0 = default 100) |
| `edgeTypes` | Filter by edge types (nil = all types) |

---

## Package `analyze`

Graph analysis algorithms.

### Hub Detection

```go
type HubNode struct {
    ID        string
    Label     string
    Type      string
    InDegree  int
    OutDegree int
    Total     int
}

func FindHubs(nodes []*graph.Node, edges []*graph.Edge, topN int, excludeTypes []string) []HubNode
func IsolatedNodes(nodes []*graph.Node, edges []*graph.Edge, threshold int, excludeTypes []string) []*graph.Node
```

| Function | Description |
|----------|-------------|
| `FindHubs` | Returns top N most connected nodes |
| `IsolatedNodes` | Returns nodes with degree <= threshold |

### Community Detection

```go
type Community struct {
    ID       int
    Size     int
    Cohesion float64
    Members  []string
    Label    string
}

type ClusterResult struct {
    Communities   []Community
    NodeCommunity map[string]int
    Modularity    float64
}

func DetectCommunities(nodes []*graph.Node, edges []*graph.Edge) *ClusterResult
func DetectCommunitiesWithOptions(nodes []*graph.Node, edges []*graph.Edge, opts ClusterOptions) *ClusterResult
```

Uses the Louvain algorithm for modularity optimization.

### ClusterOptions

```go
type ClusterOptions struct {
    Algorithm        ClusterAlgorithm  // "louvain" or "components"
    Resolution       float64           // Higher = smaller communities (default 1.0)
    ExcludeEdgeTypes []string          // Edge types to ignore
    ExcludeNodeTypes []string          // Node types to ignore
}

func DefaultClusterOptions() ClusterOptions
```

### Graph Diff

```go
type GraphDiff struct {
    NewNodes     []NodeChange
    RemovedNodes []NodeChange
    NewEdges     []EdgeChange
    RemovedEdges []EdgeChange
    Summary      string
}

type NodeChange struct {
    ID    string
    Label string
    Type  string
}

type EdgeChange struct {
    From       string
    To         string
    Type       string
    Confidence string
}

func DiffGraphs(oldNodes, newNodes []*graph.Node, oldEdges, newEdges []*graph.Edge) *GraphDiff
```

Compares two graph snapshots and returns changes.

### Utility Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NodesByType` | `func NodesByType(nodes []*graph.Node) map[string][]*graph.Node` | Group nodes by type |
| `EdgesByType` | `func EdgesByType(edges []*graph.Edge) map[string][]*graph.Edge` | Group edges by type |
| `EdgesByConfidence` | `func EdgesByConfidence(edges []*graph.Edge) map[graph.Confidence][]*graph.Edge` | Group edges by confidence |
| `CohesionScore` | `func CohesionScore(members []string, adj map[string]map[string]bool) float64` | Calculate community density |
| `HubScore` | `func HubScore(nodeID string, edges []*graph.Edge) int` | Calculate out-degree |
| `AuthorityScore` | `func AuthorityScore(nodeID string, edges []*graph.Edge) int` | Calculate in-degree |
| `InferredEdges` | `func InferredEdges(edges []*graph.Edge) []*graph.Edge` | Get INFERRED/AMBIGUOUS edges |
