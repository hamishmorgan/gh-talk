# Action Plan - Start Implementation

**Immediate actions to address critical review findings and start building**

## Critical Problems Identified

From [CRITICAL-REVIEW.md](CRITICAL-REVIEW.md):

1. 🔴 **Analysis Paralysis** - 11,882 lines of docs, 0 features
2. 🟡 **Thread ID UX** - 30-char IDs may be unusable
3. 🟢 **Missing Cobra** - Not added to go.mod
4. 🟡 **Unrealistic Coverage** - 80% too high for MVP
5. 🟢 **CI Will Fail** - Expected, need to implement code

## Immediate Fixes (Next 30 Minutes)

### Fix 1: Add Missing Dependencies

```bash
cd /Users/hamish/src/gh-talk

# Add Cobra
go get github.com/spf13/cobra@latest

# Verify
go mod tidy
go build
```

### Fix 2: Fix Go Version

```bash
# Edit go.mod
go 1.24.6 → go 1.21

# Reason: Go 1.24 doesn't exist yet
```

### Fix 3: Lower Coverage Threshold

```bash
# Edit .github/workflows/test.yml
total < 80 → total < 60

# For MVP, 60% is realistic
# Increase to 80% after Phase 1
```

### Fix 4: Add LICENSE File

```bash
# Create LICENSE file with MIT license
# SPEC says MIT but file doesn't exist
```

## Revised MVP Scope

### TRUE MVP (2 Weeks)

**Goal:** Usable for daily work, not feature-complete

**3 Commands Only:**

#### 1. `gh talk list threads`

**Functionality:**

- List review threads for current PR
- Filter: --unresolved (default), --resolved, --all
- Output: Table (TTY) or TSV (non-TTY)
- Context: Auto-detect or --pr flag

**Why First:**

- Foundation for all other commands
- Immediately useful (see what needs attention)
- Tests our architecture end-to-end

**Estimated:** 2-3 days

#### 2. `gh talk reply`

**Functionality:**

- Reply to thread by full ID
- Interactive selection if no ID
- Message as argument or --editor
- Optional --resolve flag

**Why Second:**

- Most common action (address feedback)
- Tests mutations and error handling
- Validates thread ID system

**Estimated:** 2-3 days

#### 3. `gh talk resolve`

**Functionality:**

- Resolve thread by full ID
- Interactive selection if no ID
- Bulk support (multiple IDs)
- Confirmation for bulk

**Why Third:**

- Common workflow completion action
- Tests bulk operations
- Simple mutation (good learning)

**Estimated:** 1-2 days

**Testing:** 2-3 days  
**Documentation:** 1 day  
**Polish:** 1-2 days  

**Total:** 10-15 days (2-3 weeks with buffer)

### Explicitly Deferred to Phase 2

- react command
- hide command
- show command (just use `gh pr view` for now)
- Issue support (PRs first)
- URL support
- Configuration file
- Short numeric IDs (reconsider after feedback)
- Advanced filtering (--author, --file, --since)
- Bulk confirmations

### Explicitly Deferred to Phase 3

- TUI interactive mode
- fzf integration
- Shell completion
- Advanced features

## Implementation Order (Day by Day)

### Day 1-2: Foundation

**Tasks:**

1. Add Cobra dependency
2. Fix go.mod (Go version)
3. Create `internal/commands/root.go`
4. Create `internal/api/client.go` (wrapper around go-gh)
5. Update main.go to call commands.Execute()
6. Test: `gh talk --help` works

**Deliverable:** Cobra setup, help text works

### Day 3-4: First Query

**Tasks:**

1. Create `internal/api/types.go` (Thread, Comment structs)
2. Create `internal/api/threads.go` (ListThreads method)
3. Write actual GraphQL query
4. Test with testdata/pr_full_response.json
5. Test with real PR #1

**Deliverable:** Can fetch threads from API

### Day 5-6: List Command

**Tasks:**

1. Create `internal/commands/list.go`
2. Create `internal/format/table.go`
3. Implement runListThreads
4. Add flags (--unresolved, --resolved, --all, --pr)
5. Context detection (repository.Current())
6. Test on real PR

**Deliverable:** `gh talk list threads` works!

### Day 7-8: Reply Command

**Tasks:**

1. Create `internal/api/mutations.go` (ReplyToThread)
2. Create `internal/commands/reply.go`
3. Implement runReply
4. Interactive selection (go-gh prompter)
5. Editor integration
6. Test on real PR

**Deliverable:** `gh talk reply` works!

### Day 9-10: Resolve Command

**Tasks:**

1. Add ResolveThread to mutations.go
2. Create `internal/commands/resolve.go`
3. Implement runResolve
4. Bulk support (multiple args)
5. Interactive multi-select
6. Test on real PR

**Deliverable:** `gh talk resolve` works!

### Day 11-12: Testing

**Tasks:**

1. Unit tests for api package
2. Integration tests for commands
3. Mock client implementation
4. Aim for 60%+ coverage
5. Fix any bugs discovered

**Deliverable:** Solid test coverage

### Day 13-14: Polish

**Tasks:**

1. Error message improvements
2. Documentation updates
3. README examples
4. Bug fixes
5. Performance testing

**Deliverable:** Usable MVP

### Day 15: Release

**Tasks:**

1. Tag v0.1.0
2. Trigger release workflow
3. Test installation
4. Announce
5. Gather feedback

**Deliverable:** v0.1.0 released!

## Thread ID Decision Revisited

### Reconsideration: Add Short IDs

**Why Reconsider:**

- UX is critical for adoption
- 30-char IDs are unusable
- Interactive mode friction
- Every list should be copyable

**Simple Solution:**

```go
// Session-scoped cache (in-memory only)
type ThreadCache struct {
    PR       int
    Threads  []Thread
    Mapping  map[int]string  // 1 → PRRT_kwDO...
    FetchedAt time.Time
}

var sessionCache *ThreadCache  // Global, OK for CLI

func getThreadID(arg string) (string, error) {
    // Try short ID first
    if num, err := strconv.Atoi(arg); err == nil {
        if sessionCache != nil && sessionCache.Mapping != nil {
            if id, ok := sessionCache.Mapping[num]; ok {
                return id, nil
            }
        }
        return "", fmt.Errorf("no thread #%d in current session\nRun 'gh talk list threads' first", num)
    }
    
    // Fall back to full ID
    if strings.HasPrefix(arg, "PRRT_") {
        return arg, nil
    }
    
    // Interactive
    if arg == "" {
        return selectInteractively()
    }
    
    return "", fmt.Errorf("invalid thread ID")
}
```

**Rules:**

- Cache populated by `list threads`
- Valid for current session only
- Clear on PR context change
- Fallback to full ID always works

**Implementation:**

- 30-60 minutes of work
- Huge UX improvement
- Simple (no persistent storage)

**Recommendation:** **ADD IT** to MVP

## Dependency Injection Pattern

### Solution for Testing Commands

```go
// internal/api/client.go
type Client interface {
    ListThreads(owner, name string, pr int) ([]Thread, error)
    ResolveThread(id string) error
    // ... other methods
}

type graphQLClient struct {
    client *api.GraphQLClient
}

func NewClient() (Client, error) {
    gql, err := api.DefaultGraphQLClient()
    if err != nil {
        return nil, err
    }
    return &graphQLClient{client: gql}, nil
}

// internal/api/mock.go (for tests)
type MockClient struct {
    ListThreadsFunc func(owner, name string, pr int) ([]Thread, error)
    // ... other funcs
}

func (m *MockClient) ListThreads(owner, name string, pr int) ([]Thread, error) {
    if m.ListThreadsFunc != nil {
        return m.ListThreadsFunc(owner, name, pr)
    }
    return nil, nil
}
```

```go
// internal/commands/list.go
var clientFactory = api.NewClient  // Can override in tests

func runListThreads(cmd *cobra.Command, args []string) error {
    client, err := clientFactory()  // Use factory
    // ...
}

// internal/commands/list_test.go
func TestListThreads(t *testing.T) {
    // Override factory
    originalFactory := clientFactory
    defer func() { clientFactory = originalFactory }()
    
    clientFactory = func() (api.Client, error) {
        return &api.MockClient{
            ListThreadsFunc: func(owner, name string, pr int) ([]Thread, error) {
                // Return test data
                return loadThreadsFromFixture(t, "testdata/pr_full_response.json")
            },
        }, nil
    }
    
    // Test command
    cmd := NewListCommand()
    err := cmd.Execute()
    // Assertions...
}
```

**Estimated:** 1 hour to set up pattern

## Success Criteria

### MVP Success = Can Use Daily

**Must Be Able To:**

1. See unresolved threads in current PR
2. Reply to threads
3. Resolve threads after addressing

**If This Works:**

- Saves time vs browser
- Actually used in workflow
- Validated the concept

**Then Add:**

- More commands based on what's missing
- Features users actually request
- Improvements from real usage

## Risk Mitigation

### For Each Risk in Critical Review

**Analysis Paralysis:**

- ✅ STOP researching after this
- ✅ Commit to implementation timeline
- ✅ 2-week MVP deadline

**Thread ID UX:**

- 🔄 ADD short IDs (session-scoped)
- ✅ 30 min investment, huge UX win
- ✅ Best of both worlds

**GraphQL Complexity:**

- ✅ Start with simplest query
- ✅ Test against fixtures first
- ✅ Iterate on structure

**Coverage Goals:**

- ✅ Lower to 60% for MVP
- ✅ Focus on API package
- ✅ Increase coverage in Phase 2

**Scope Creep:**

- ✅ TRUE MVP: 3 commands only
- ✅ Defer everything else
- ✅ User feedback drives roadmap

## Commit Strategy

### This Commit

**Commit Message:**

```
Add engineering infrastructure and critical review

Infrastructure:
- CI/CD workflows (test, build)
- Linting configuration
- Makefile and tooling
- Documentation (ENGINEERING, ENVIRONMENT, etc.)

Critical Review:
- Identified 5 critical problems
- Analysis paralysis (too much planning)
- Thread ID UX concerns (reconsidering short IDs)
- Unrealistic coverage goals (lowering to 60%)
- Clear action plan to start implementation

Action Plan:
- TRUE MVP: 3 commands in 2 weeks
- Add short IDs (session-scoped)
- Focus on usability over perfection
- Iterate based on real usage

Status: Research complete, problems identified, ready to build
```

### Next Commit

**Commit Message:**

```
Add Cobra and implement foundation

- Add Cobra dependency
- Fix Go version (1.21)
- Create root command
- Set up command structure
- Update main.go
```

### Third Commit

**Commit Message:**

```
Implement list threads command (MVP #1)

- API client wrapper
- GraphQL query for threads
- Table formatting
- Basic filtering
- Tests with fixtures

First working feature!
```

## Summary

### Problems Found

**Critical:**

1. Too much planning, not enough building
2. Thread ID system may have poor UX
3. Scope too large for MVP

**Medium:**
4. Coverage goals unrealistic
5. Some design decisions untested
6. Missing dependency injection pattern

**Minor:**
7. Go version wrong (1.24.6 doesn't exist)
8. Small documentation inconsistencies
9. Missing LICENSE file

### Actions Required

**Before Next Coding Session:**

1. ✅ Fix go.mod (Go version, add Cobra)
2. ✅ Lower coverage target (60%)
3. ✅ Add LICENSE file
4. ✅ Reconsider short IDs (add them!)

**First Implementation Sprint:**

1. Root command setup
2. API client wrapper
3. List threads command
4. Test on real PR
5. Iterate

**Success Metric:**

- Working `gh talk list threads` in 3-4 days
- Usable for real code review
- Validated our architecture

## The Hard Truth

**We've done:**

- ✅ Exceptional research
- ✅ Thorough planning
- ✅ Quality infrastructure

**We haven't done:**

- ❌ Built anything
- ❌ Validated with real code
- ❌ Proven it works

**Next step is clear:** **Build the first command** 🚀

---

**Last Updated**: 2025-11-02  
**Status**: Problems identified, action plan created  
**Next Action**: Fix go.mod and start implementing
