# GraphFS

Git-friendly filesystem graph database with one-file-per-entity storage for minimal diffs.

## Overview

GraphFS is a Go library for storing and querying graph data using the filesystem. Each node and edge is stored as a separate JSON file, making it ideal for:

- **Version control** - Clean git diffs when graph structure changes
- **Code analysis tools** - Store AST relationships, call graphs, dependency trees
- **LLM-augmented graphs** - Mix deterministic (AST) and inferred (LLM) relationships
- **Debugging** - Human-readable JSON files you can inspect directly

## Key Features

- **One file per entity** - Nodes stored as `nodes/{id}.json`, edges as `edges/{from}__{type}__{to}.json`
- **Deterministic JSON** - Sorted keys and consistent formatting for clean git diffs
- **Confidence levels** - Support for EXTRACTED (AST), INFERRED (LLM), and AMBIGUOUS relationships
- **Pluggable storage** - `Store` interface for custom backends
- **Schema validation** - Validate nodes, edges, and referential integrity
- **Graph traversal** - BFS, DFS, and path finding algorithms
- **Graph analysis** - Hub detection, community detection (Louvain), graph diff

## Quick Example

```go
import (
    "github.com/plexusone/graphfs/pkg/graph"
    "github.com/plexusone/graphfs/pkg/store"
)

// Create a filesystem store
fs, _ := store.NewFSStore(".graphfs")

// Add a node
fs.WriteNode(&graph.Node{
    ID:    "func_main",
    Type:  graph.NodeTypeFunction,
    Label: "main",
})

// Add an edge
fs.WriteEdge(&graph.Edge{
    From:       "func_main",
    To:         "func_helper",
    Type:       graph.EdgeTypeCalls,
    Confidence: graph.ConfidenceExtracted,
})

// Load the entire graph
g, _ := fs.LoadGraph()
fmt.Printf("Nodes: %d, Edges: %d\n", g.NodeCount(), g.EdgeCount())
```

## Installation

```bash
go get github.com/plexusone/graphfs
```

## Documentation

- [Getting Started](getting-started.md) - Installation and basic usage
- [Storage Format](storage-format.md) - File layout and JSON schema
- [API Reference](api-reference.md) - Types, interfaces, and constants
- [Changelog](releases/changelog.md) - Release history

## License

MIT License - see [LICENSE](https://github.com/plexusone/graphfs/blob/main/LICENSE) for details.
