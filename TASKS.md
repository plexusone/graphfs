# GraphFS Tasks

Git-friendly filesystem graph database.

## Phase 1 - MVP (Complete)

- [x] Core types: Node, Edge, Graph
- [x] Filesystem store: one file per entity (full CRUD)
- [x] Deterministic JSON serialization (sorted keys)
- [x] Schema validation (ValidateNode, ValidateEdge, ValidateGraph)
- [x] Unit tests for all packages (40 tests)
- [x] CLI: `graphfs validate <path>`
- [x] CLI: `graphfs format <path>`

## Phase 2 - Enhanced Validation

- [ ] Custom schema definitions (YAML/JSON config)
- [x] Referential integrity checks (ValidateEdgeRefs)
- [ ] Duplicate detection
- [ ] Cycle detection (optional)

## Phase 3 - Diff & Query

- [ ] Semantic diff: show added/removed nodes/edges
- [ ] Basic graph traversal queries
- [ ] Path finding between nodes

## Phase 4 - Git Integration (Optional)

- [ ] Pre-commit hook for validation
- [ ] PR diff annotations
- [ ] CI/CD integration helpers

## Phase 5 - Advanced Features

- [ ] Graph merging (for multi-source inputs)
- [ ] Incremental updates (patch files)
- [ ] Export formats: GraphML, Cypher (Neo4j), DOT
