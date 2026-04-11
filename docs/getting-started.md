# Getting Started

## Installation

```bash
go get github.com/plexusone/graphfs
```

## Creating a Store

GraphFS uses a filesystem-backed store where each node and edge is a separate JSON file:

```go
import "github.com/plexusone/graphfs/pkg/store"

fs, err := store.NewFSStore(".graphfs")
if err != nil {
    panic(err)
}
```

This creates the following directory structure:

```
.graphfs/
  nodes/    # One JSON file per node
  edges/    # One JSON file per edge
```

## Working with Nodes

### Creating Nodes

```go
import "github.com/plexusone/graphfs/pkg/graph"

node := &graph.Node{
    ID:    "func_main",
    Type:  graph.NodeTypeFunction,
    Label: "main",
    Attrs: map[string]string{
        "package":     "main",
        "source_file": "main.go",
        "line":        "10",
    },
}

err := fs.WriteNode(node)
```

### Reading Nodes

```go
node, err := fs.GetNode("func_main")
if err != nil {
    // Node not found or read error
}
fmt.Printf("Found: %s (%s)\n", node.Label, node.Type)
```

### Listing All Nodes

```go
nodes, err := fs.ListNodes()
for _, n := range nodes {
    fmt.Printf("%s: %s\n", n.ID, n.Label)
}
```

### Deleting Nodes

```go
err := fs.DeleteNode("func_main")
```

## Working with Edges

### Creating Edges

```go
edge := &graph.Edge{
    From:       "func_main",
    To:         "func_helper",
    Type:       graph.EdgeTypeCalls,
    Confidence: graph.ConfidenceExtracted,
}

err := fs.WriteEdge(edge)
```

For LLM-inferred relationships, include a confidence score:

```go
edge := &graph.Edge{
    From:            "pkg_auth",
    To:              "pkg_crypto",
    Type:            graph.EdgeTypeDependsOn,
    Confidence:      graph.ConfidenceInferred,
    ConfidenceScore: 0.85,
    Attrs: map[string]string{
        "reason": "auth package likely uses crypto for password hashing",
    },
}
```

### Reading Edges

```go
edge, err := fs.GetEdge("func_main", "calls", "func_helper")
```

### Listing All Edges

```go
edges, err := fs.ListEdges()
for _, e := range edges {
    fmt.Printf("%s --%s--> %s\n", e.From, e.Type, e.To)
}
```

## Working with Graphs

### Loading a Complete Graph

```go
g, err := fs.LoadGraph()
if err != nil {
    panic(err)
}

fmt.Printf("Loaded %d nodes and %d edges\n", g.NodeCount(), g.EdgeCount())
```

### Saving a Complete Graph

```go
g := graph.NewGraph()

g.AddNode(&graph.Node{ID: "a", Type: "function", Label: "funcA"})
g.AddNode(&graph.Node{ID: "b", Type: "function", Label: "funcB"})
g.AddEdge(&graph.Edge{
    From: "a", To: "b",
    Type: "calls",
    Confidence: graph.ConfidenceExtracted,
})

err := fs.SaveGraph(g)
```

### In-Memory Operations

```go
g := graph.NewGraph()

// Add nodes
g.AddNode(&graph.Node{ID: "x", Type: "file", Label: "main.go"})

// Look up nodes
node := g.GetNode("x")

// Count entities
fmt.Printf("Nodes: %d, Edges: %d\n", g.NodeCount(), g.EdgeCount())
```

## Validation

### Validating Individual Entities

```go
import "github.com/plexusone/graphfs/pkg/schema"

validator := schema.NewValidator()

// Validate a node
if err := validator.ValidateNode(node); err != nil {
    fmt.Printf("Invalid node: %v\n", err)
}

// Validate an edge
if err := validator.ValidateEdge(edge); err != nil {
    fmt.Printf("Invalid edge: %v\n", err)
}
```

### Restricting Allowed Types

```go
validator := schema.NewValidator()
validator.AllowedNodeTypes = []string{"function", "file", "package"}
validator.AllowedEdgeTypes = []string{"calls", "contains", "imports"}
validator.RequireNodeLabel = true

if err := validator.ValidateNode(node); err != nil {
    // Error if node type not in allowed list
}
```

### Validating Referential Integrity

```go
// Check that edge endpoints exist
err := validator.ValidateEdgeRefs(edge, g.Nodes)
if err != nil {
    fmt.Printf("Dangling reference: %v\n", err)
}
```

### Validating an Entire Graph

```go
errs := validator.ValidateGraph(g)
if len(errs) > 0 {
    for _, err := range errs {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Next Steps

- [Storage Format](storage-format.md) - Understand the file layout
- [API Reference](api-reference.md) - Complete type and constant reference
