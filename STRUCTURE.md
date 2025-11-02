# Project Structure

## Overview

`gh-talk` follows a **minimal packages** approach that balances simplicity with maintainability. This structure is appropriate for a GitHub CLI extension of moderate complexity.

## Directory Layout

```
gh-talk/
├── main.go              # Entry point only
├── internal/            # All implementation (private)
│   ├── api/            # GitHub GraphQL API client
│   ├── commands/       # Cobra command implementations
│   ├── filter/         # Thread/comment filtering logic
│   ├── format/         # Output formatting (table, JSON, markdown)
│   ├── config/         # Configuration management
│   ├── cache/          # Caching layer for API responses
│   └── tui/            # Terminal UI for interactive mode
├── go.mod
├── go.sum
├── SPEC.md             # Complete specification
├── README.md           # User documentation
├── AGENTS.md           # AI agent guidelines
└── structure_test.go   # Integration test
```

## Design Decisions

### Why `main.go` in Root?

- **GitHub CLI Extension Convention**: Extensions are invoked as `gh-EXTENSION-NAME`, and having `main.go` in the root follows the standard pattern for gh extensions.
- **Simplicity**: No extra navigation through `cmd/` directories for a single binary.

### Why `internal/` for Everything?

- **Correct Semantics**: Nothing in this codebase is meant to be imported by external projects. We're building a CLI tool, not a library.
- **Compiler Enforcement**: The Go compiler prevents external imports from `internal/`, ensuring encapsulation.
- **Clear Intent**: Signals to developers that this is application code, not reusable library code.

### Why NOT `pkg/`?

- **Misleading**: `pkg/` signals "safe for external import," which doesn't apply to a CLI extension.
- **Not a Library**: If we wanted to extract reusable components later, we'd create a separate module.

### Why NOT `cmd/`?

- **Single Binary**: The `cmd/` directory is for projects with MULTIPLE executables (like `kubectl` which has `kubectl`, `kubelet`, `kube-proxy`).
- **Unnecessary Layer**: For a single-binary project, it adds navigation overhead without benefit.

## Package Responsibilities

### `internal/api`
- GraphQL client for GitHub API
- Query construction and execution
- Response parsing
- Error handling and retries

### `internal/commands`
- Cobra command definitions
- Command-line flag handling
- Command execution logic
- Help text and examples

### `internal/filter`
- Thread/comment filtering by status, author, date, file, etc.
- Filter composition and logic
- Predicate functions

### `internal/format`
- Table output formatting
- JSON serialization
- Markdown generation
- Terminal styling and colors

### `internal/config`
- Configuration file management
- Environment variable handling
- Default values
- Aliases

### `internal/cache`
- API response caching
- Cache invalidation
- TTL management

### `internal/tui`
- Interactive terminal UI (Bubble Tea)
- Keyboard navigation
- Real-time updates
- UI components

## Comparison to Alternatives

### vs. Flat Structure
**Flat** (all files in root):
- ✅ Simpler for very small projects
- ❌ Becomes cluttered at scale
- ❌ Harder to enforce boundaries

**Current** (minimal packages):
- ✅ Clear separation of concerns
- ✅ Scalable to 10,000+ LOC
- ✅ Easy to test individual packages

### vs. Standard Go Layout
**Standard** (cmd/, pkg/, internal/):
- ✅ Familiar to Go developers
- ✅ Scales to very large projects
- ❌ Overkill for CLI extensions
- ❌ `pkg/` is misleading for non-library code

**Current** (minimal packages):
- ✅ Right-sized for the project
- ✅ Follows gh extension conventions
- ✅ Room to grow without reorganizing

## Evolution Path

As the project grows, the structure can evolve:

1. **Small** (< 3,000 LOC): Current structure is perfect
2. **Medium** (3,000-10,000 LOC): Add more packages under `internal/` as needed
3. **Large** (10,000+ LOC): Consider extracting reusable components to separate modules

The beauty of this structure is that it doesn't need major refactoring as the project scales—just add more packages under `internal/` as domains emerge.

## Testing Strategy

- Each package has its own test file (`*_test.go`)
- `structure_test.go` validates the overall package organization
- All packages are importable and compile successfully

## References

- [GitHub CLI Extensions Documentation](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions)
- [Go's Internal Packages](https://go.dev/doc/go1.4#internalpackages)
- [Structuring Go Applications](https://www.gobeyond.dev/packages-as-layers/)

---

**Last Updated**: 2025-11-02

