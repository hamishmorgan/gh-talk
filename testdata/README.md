# Test Data

Real GitHub API responses captured from live testing.

## Files

### `pr_full_response.json`

**Source:** PR #1 (<https://github.com/hamishmorgan/gh-talk/pull/1>)  
**Captured:** Initial state with 2 review threads  
**Contains:**

- Complete PullRequest object
- 2 ReviewThreads with comments
- Reaction groups
- Complete metadata

**Use For:**

- Understanding complete PR structure
- Testing PR query parsing
- Example of review threads with replies
- Real ID formats

### `issue_full_response.json`

**Source:** Issue #2 (<https://github.com/hamishmorgan/gh-talk/issues/2>)  
**Captured:** Issue with comments and reactions  
**Contains:**

- Complete Issue object
- 2 issue comments (1 minimized)
- Reactions on issue body and comments
- Labels and participants

**Use For:**

- Understanding issue structure
- Testing issue comment parsing
- Example of minimized comments
- Issue-specific fields (stateReason, participants)

### `pr_with_resolved_threads.json`

**Source:** PR #1 (after additional testing)  
**Captured:** State with mixed resolved/unresolved threads  
**Contains:**

- 4 ReviewThreads total
- 2 resolved threads (with resolvedBy)
- 2 unresolved threads
- Various comment counts (1-2 per thread)

**Use For:**

- Testing thread resolution filtering
- Understanding resolved vs unresolved states
- Permission field behavior (viewerCanResolve/Unresolve)
- Display logic for mixed states

## Usage in Tests

```go
package api

import (
    "encoding/json"
    "os"
    "testing"
)

func TestParseThreads(t *testing.T) {
    // Load test fixture
    data, _ := os.ReadFile("../testdata/pr_with_resolved_threads.json")
    
    var response struct {
        Data struct {
            Repository struct {
                PullRequest struct {
                    ReviewThreads struct {
                        Nodes []Thread
                    }
                }
            }
        }
    }
    
    json.Unmarshal(data, &response)
    
    threads := response.Data.Repository.PullRequest.ReviewThreads.Nodes
    
    // Test filtering
    resolved := filterResolved(threads)
    if len(resolved) != 2 {
        t.Errorf("Expected 2 resolved threads, got %d", len(resolved))
    }
}
```

## Data Integrity

**These files contain:**

- ✅ Real production data from GitHub API
- ✅ Actual ID formats (not mocked)
- ✅ Complete field structures
- ✅ Real timestamps
- ✅ Authentic relationships

**Do NOT:**

- ❌ Manually edit (keep as real API responses)
- ❌ Commit updated versions without testing
- ❌ Use for production (test data only)

## Regenerating Test Data

To capture fresh test data:

```bash
# Query PR
gh api graphql -f query='...' > testdata/pr_full_response.json

# Query Issue
gh api graphql -f query='...' > testdata/issue_full_response.json
```

See `docs/REAL-DATA.md` for complete query examples.

---

**Last Updated**: 2025-11-02  
**Source**: Real GitHub API responses from hamishmorgan/gh-talk
