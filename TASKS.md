# GraphFS Tasks

Git-friendly filesystem graph database.

## Phase 1 - MVP (Complete)

- [x] Core types: Node, Edge, Graph
- [x] Filesystem store: one file per entity (full CRUD)
- [x] Deterministic JSON serialization (sorted keys)
- [x] Schema validation (ValidateNode, ValidateEdge, ValidateGraph)
- [x] Referential integrity checks (ValidateEdgeRefs)
- [x] Unit tests for all packages (48 tests)
- [x] CLI: `graphfs validate <path>`
- [x] CLI: `graphfs format <path>`

## Phase 2 - Documentation (Complete)

- [x] README with usage examples
- [x] MkDocs documentation site (Material theme)
- [x] Getting started guide
- [x] Storage format specification
- [x] API reference
- [x] CHANGELOG (structured JSON + Markdown)
- [x] GitHub Actions CI workflows

## Phase 3 - Query & Analysis (Current)

- [ ] Basic graph traversal (BFS, DFS)
- [ ] Path finding between nodes
- [ ] Cycle detection (find cycles in graph)
- [ ] Semantic diff: show added/removed nodes/edges
- [ ] CLI: `graphfs query <node-id>`
- [ ] CLI: `graphfs diff <path1> <path2>`

## Phase 4 - Schema Extensions

- [ ] Custom schema definitions (YAML/JSON config)
- [ ] Custom node/edge type validation
- [ ] Required attributes per type

## Phase 5 - Git Integration (Optional)

- [ ] Pre-commit hook for validation
- [ ] PR diff annotations
- [ ] CI/CD integration helpers

## Phase 6 - Advanced Features

- [ ] Graph merging (for multi-source inputs)
- [ ] Incremental updates (patch files)
- [ ] Export formats: GraphML, Cypher (Neo4j), DOT
