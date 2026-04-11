# GraphFS

Git-friendly filesystem graph database with one-file-per-entity storage for minimal diffs.

## Features

- 📄 **One file per entity** - Nodes stored as `nodes/{id}.json`, edges as `edges/{from}__{type}__{to}.json`
- 🔒 **Deterministic JSON** - Sorted keys and consistent formatting for clean git diffs
- 🎯 **Confidence levels** - Support for EXTRACTED (AST), INFERRED (LLM), and AMBIGUOUS relationships
- 🔌 **Pluggable storage** - `Store` interface for custom backends
- ✅ **Schema validation** - Validate nodes, edges, and referential integrity

## Installation

```bash
go get github.com/plexusone/graphfs
```

## Usage

### Creating a Graph

```go
import (
    "github.com/plexusone/graphfs/pkg/graph"
    "github.com/plexusone/graphfs/pkg/store"
)

// Create a filesystem store
fs, err := store.NewFSStore(".graphfs")
if err != nil {
    panic(err)
}

// Create nodes
node := &graph.Node{
    ID:    "func_main",
    Type:  graph.NodeTypeFunction,
    Label: "main",
    Attrs: map[string]string{"package": "main"},
}
fs.WriteNode(node)

// Create edges
edge := &graph.Edge{
    From:       "func_main",
    To:         "func_helper",
    Type:       graph.EdgeTypeCalls,
    Confidence: graph.ConfidenceExtracted,
}
fs.WriteEdge(edge)
```

### Loading a Graph

```go
g, err := fs.LoadGraph()
if err != nil {
    panic(err)
}

fmt.Printf("Nodes: %d, Edges: %d\n", g.NodeCount(), g.EdgeCount())
```

### Validation

```go
import "github.com/plexusone/graphfs/pkg/schema"

validator := schema.NewValidator()
validator.AllowedNodeTypes = []string{"function", "file", "package"}

if err := validator.ValidateNode(node); err != nil {
    fmt.Printf("Invalid node: %v\n", err)
}

// Validate entire graph
errs := validator.ValidateGraph(g)
for _, err := range errs {
    fmt.Printf("Error: %v\n", err)
}
```

## Storage Format

### Nodes

```
.graphfs/
  nodes/
    func_main.json
    func_helper.json
    pkg_mypackage.json
```

Each node file contains:

```json
{
  "attrs": {
    "package": "main"
  },
  "id": "func_main",
  "label": "main",
  "type": "function"
}
```

### Edges

```
.graphfs/
  edges/
    func_main__calls__func_helper.json
```

Each edge file contains:

```json
{
  "confidence": "EXTRACTED",
  "from": "func_main",
  "to": "func_helper",
  "type": "calls"
}
```

## Node Types

| Constant | Value |
|----------|-------|
| `NodeTypeFunction` | `function` |
| `NodeTypeMethod` | `method` |
| `NodeTypeClass` | `class` |
| `NodeTypeStruct` | `struct` |
| `NodeTypeFile` | `file` |
| `NodeTypePackage` | `package` |
| `NodeTypeModule` | `module` |
| `NodeTypeVariable` | `variable` |
| `NodeTypeConstant` | `constant` |
| `NodeTypeInterface` | `interface` |

## Edge Types

| Constant | Value |
|----------|-------|
| `EdgeTypeCalls` | `calls` |
| `EdgeTypeImports` | `imports` |
| `EdgeTypeImplements` | `implements` |
| `EdgeTypeExtends` | `extends` |
| `EdgeTypeUses` | `uses` |
| `EdgeTypeContains` | `contains` |
| `EdgeTypeDependsOn` | `depends_on` |
| `EdgeTypeReferences` | `references` |

## Confidence Levels

| Level | Description |
|-------|-------------|
| `EXTRACTED` | Directly extracted from source (AST, imports) |
| `INFERRED` | Inferred by LLM or heuristic with confidence score |
| `AMBIGUOUS` | Uncertain relationship requiring human review |

## License

MIT
