# Storage Format

GraphFS stores each node and edge as a separate JSON file, designed for git-friendly diffs and human readability.

## Directory Structure

```
.graphfs/
├── nodes/
│   ├── func_main.json
│   ├── func_helper.json
│   ├── pkg_mypackage.json
│   └── ...
└── edges/
    ├── func_main__calls__func_helper.json
    ├── pkg_mypackage__contains__func_main.json
    └── ...
```

## Node Files

Nodes are stored in `nodes/{id}.json`:

```json
{
  "attrs": {
    "line": "10",
    "package": "main",
    "source_file": "main.go"
  },
  "id": "func_main",
  "label": "main",
  "type": "function"
}
```

### Node Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier (filesystem-safe) |
| `type` | string | Yes | Node category (function, file, etc.) |
| `label` | string | No | Human-readable display name |
| `attrs` | object | No | Additional key-value metadata |

### Node ID Conventions

Node IDs must be filesystem-safe (no `/\:*?"<>|` characters). Recommended conventions:

| Node Type | ID Pattern | Example |
|-----------|------------|---------|
| Function | `func_{file}.{name}` | `func_main.go.HandleRequest` |
| Method | `method_{receiver}.{name}` | `method_Server.Start` |
| Type/Struct | `type_{name}` | `type_User` |
| Package | `pkg_{name}` | `pkg_auth` |
| File | `file_{path}` | `file_cmd_main.go` |

## Edge Files

Edges are stored in `edges/{from}__{type}__{to}.json`:

```json
{
  "confidence": "EXTRACTED",
  "from": "func_main",
  "to": "func_helper",
  "type": "calls"
}
```

With confidence score (for inferred edges):

```json
{
  "attrs": {
    "reason": "Both handle authentication flow"
  },
  "confidence": "INFERRED",
  "confidence_score": 0.85,
  "from": "pkg_auth",
  "to": "pkg_session",
  "type": "depends_on"
}
```

### Edge Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `from` | string | Yes | Source node ID |
| `to` | string | Yes | Target node ID |
| `type` | string | Yes | Relationship type (calls, imports, etc.) |
| `confidence` | string | Yes | EXTRACTED, INFERRED, or AMBIGUOUS |
| `confidence_score` | float | No | 0.0-1.0 score for INFERRED edges |
| `attrs` | object | No | Additional key-value metadata |

### Edge Filename Convention

Edge filenames use double underscores as separators:

```
{from}__{type}__{to}.json
```

Examples:

- `func_main__calls__func_helper.json`
- `pkg_auth__imports__pkg_crypto.json`
- `type_User__extends__type_BaseModel.json`

## Canonical JSON Format

GraphFS uses deterministic JSON serialization for clean git diffs:

1. **Sorted keys** - Object keys are alphabetically ordered
2. **Consistent indentation** - 2-space indentation
3. **No trailing newline** - Files end without a trailing newline
4. **No HTML escaping** - Characters like `<` and `>` are not escaped

This ensures that:

- Same data always produces identical output
- Git diffs show only actual changes
- Files are human-readable and editable

### Example Diff

When adding an attribute to a node:

```diff
 {
   "attrs": {
+    "doc": "Entry point for the application",
     "package": "main"
   },
   "id": "func_main",
```

## Confidence Levels

### EXTRACTED

Deterministic relationships extracted directly from source code:

- Import statements
- Function calls (AST analysis)
- Type definitions
- Method receivers

```json
{
  "confidence": "EXTRACTED",
  "from": "file_main.go",
  "to": "pkg_fmt",
  "type": "imports"
}
```

### INFERRED

Relationships discovered by LLM analysis or heuristics:

- Implicit dependencies
- Semantic similarity
- Design pattern detection
- Cross-cutting concerns

```json
{
  "confidence": "INFERRED",
  "confidence_score": 0.75,
  "from": "pkg_handlers",
  "to": "pkg_middleware",
  "type": "depends_on",
  "attrs": {
    "reason": "Handler functions use middleware for auth"
  }
}
```

### AMBIGUOUS

Uncertain relationships requiring human review:

```json
{
  "confidence": "AMBIGUOUS",
  "confidence_score": 0.25,
  "from": "func_processData",
  "to": "func_validateInput",
  "type": "calls",
  "attrs": {
    "note": "Indirect call through interface, needs verification"
  }
}
```

## Best Practices

### 1. Use Meaningful IDs

Bad:

```json
{"id": "n1", "type": "function", "label": "main"}
```

Good:

```json
{"id": "func_main.go.main", "type": "function", "label": "main"}
```

### 2. Include Source Location

```json
{
  "attrs": {
    "source_file": "pkg/auth/handler.go",
    "line": "42",
    "end_line": "58"
  }
}
```

### 3. Document Inferred Relationships

```json
{
  "confidence": "INFERRED",
  "confidence_score": 0.8,
  "attrs": {
    "reason": "Both packages handle user authentication",
    "extracted_by": "claude-3-opus"
  }
}
```

### 4. Use Appropriate Edge Types

- `calls` - Direct function/method invocation
- `imports` - Package import
- `contains` - Hierarchical containment (package contains file)
- `references` - Type reference (field type, return type)
- `implements` - Interface implementation
- `extends` - Struct embedding
- `depends_on` - Inferred dependency
- `uses` - General usage relationship
