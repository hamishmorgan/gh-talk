# Real GitHub API Data Structures

**Documentation of actual API responses from live testing**

## Overview

This document captures real API responses from GitHub's GraphQL API, tested on PR #1 in the gh-talk repository. All IDs, structures, and behaviors are from production data.

**Test PR:** https://github.com/hamishmorgan/gh-talk/pull/1

## Real ID Formats

### ID Format Discovery

All GitHub IDs follow a pattern: `PREFIX_base64EncodedData`

**Actual ID Examples from Testing:**

| Type | Prefix | Example | Notes |
|------|--------|---------|-------|
| Pull Request | `PR_` | `PR_kwDOQN97u86xFPR4` | Node ID for PR |
| Review Thread | `PRRT_` | `PRRT_kwDOQN97u85gQeTN` | Thread container |
| Review Comment | `PRRC_` | `PRRC_kwDOQN97u86UHqK7` | Individual comment in thread |
| Review | `PRR_` | `PRR_kwDOQN97u87LMeCy` | Review submission |
| Reaction | `REA_` | `REA_lATOQN97u86UHqK7zhQo4hY` | Emoji reaction |
| Issue Comment | `IC_` | `IC_kwDOQN97u87PVA8l` | Top-level PR/Issue comment |
| User | `MDQ6` | `MDQ6VXNlcjU1OTUwOA==` | User ID (legacy format) |

**Key Findings:**
- ✅ Our spec's assumed format `PRRT_abc123` was close to reality
- ✅ Prefixes match our expectations
- ⚠️ Real IDs are 20-30 characters (much longer than assumed)
- ⚠️ IDs are opaque - no way to derive them from URLs
- ⚠️ Must store or query to get IDs

### Database IDs vs Node IDs

**Both Exist:**
```json
{
  "id": "PRRC_kwDOQN97u86UHqK7",        // Node ID (GraphQL)
  "databaseId": 2485035707                // Database ID (legacy, integer)
}
```

**Use Node IDs:**
- GraphQL API requires Node IDs
- Database IDs are legacy (REST API)
- Node IDs are stable and portable

## PullRequestReviewThread Structure

### Complete Thread Object

**Real Response from Testing:**
```json
{
  "id": "PRRT_kwDOQN97u85gQeTN",
  "isResolved": false,
  "isCollapsed": false,
  "isOutdated": false,
  "path": "test_file.go",
  "line": 7,
  "startLine": 7,
  "diffSide": "RIGHT",
  "startDiffSide": null,
  "subjectType": "LINE",
  "resolvedBy": null,
  "viewerCanResolve": true,
  "viewerCanUnresolve": false,
  "viewerCanReply": true,
  "comments": {
    "totalCount": 2,
    "nodes": [...]
  }
}
```

### Field Analysis

**Resolution Fields:**
- `isResolved: false` - Thread is not resolved
- `isCollapsed: false` - Thread is visible (not collapsed in UI)
- `resolvedBy: null` - No user has resolved it yet
- When resolved: `resolvedBy` contains user object

**Line Information:**
- `line: 7` - Current line number in file
- `startLine: 7` - Start of multi-line comment (same for single-line)
- `diffSide: "RIGHT"` - Comment on new code (LEFT = old, RIGHT = new)
- `startDiffSide: null` - Only set for multi-line comments
- `subjectType: "LINE"` - Type is LINE (vs FILE for file-level comments)

**Position Info:**
- `path: "test_file.go"` - File path in repository
- `position` in comments - Position in diff (not in file)

**Permissions:**
- `viewerCanResolve: true` - Current user can resolve
- `viewerCanUnresolve: false` - Cannot unresolve (until it's resolved)
- `viewerCanReply: true` - Can add replies

**Outdated Status:**
- `isOutdated: false` - Diff hasn't changed
- When true: diff changed since comment was made
- Outdated threads still appear but with indicator

## PullRequestReviewComment Structure

### Complete Comment Object

**Original Comment in Thread:**
```json
{
  "id": "PRRC_kwDOQN97u86UHqK7",
  "databaseId": 2485035707,
  "body": "Consider using a constant for the TODO comment",
  "path": "test_file.go",
  "position": 7,
  "originalPosition": 7,
  "diffHunk": "@@ -0,0 +1,25 @@\n+package main\n+\n+import \"fmt\"\n+\n+// TestFunction demonstrates some code for review testing\n+func TestFunction() {\n+\t// TODO: This could be improved",
  "createdAt": "2025-11-02T21:43:10Z",
  "updatedAt": "2025-11-02T21:43:10Z",
  "author": {
    "login": "hamishmorgan"
  },
  "authorAssociation": "OWNER",
  "replyTo": null,
  "isMinimized": false,
  "minimizedReason": null,
  "viewerCanReact": true,
  "viewerCanUpdate": true,
  "viewerCanDelete": true,
  "viewerCanMinimize": true,
  "reactionGroups": [...]
}
```

**Reply Comment in Thread:**
```json
{
  "id": "PRRC_kwDOQN97u86UHqOo",
  "databaseId": 2485035944,
  "body": "Good point! I will refactor this to use a constant.",
  "path": "test_file.go",
  "position": 7,
  "originalPosition": 7,
  "diffHunk": "...",
  "createdAt": "2025-11-02T21:43:42Z",
  "updatedAt": "2025-11-02T21:43:42Z",
  "author": {
    "login": "hamishmorgan"
  },
  "authorAssociation": "OWNER",
  "replyTo": {
    "id": "PRRC_kwDOQN97u86UHqK7"
  },
  "isMinimized": false,
  "minimizedReason": null,
  "viewerCanReact": true,
  "viewerCanUpdate": true,
  "viewerCanDelete": true,
  "viewerCanMinimize": true,
  "reactionGroups": [...]
}
```

### Key Differences: Original vs Reply

**Original Comment:**
- `replyTo: null` - First in thread
- Creates the thread

**Reply Comment:**
- `replyTo: { "id": "PRRC_..." }` - Points to parent comment
- Same position/path as original
- Same diffHunk

**Threading:**
- Thread contains multiple comments
- Comments linked via `replyTo`
- All comments share same thread ID
- Position determines thread grouping

### Diff Hunk Format

**Example:**
```
@@ -0,0 +1,25 @@
+package main
+
+import \"fmt\"
+
+// TestFunction demonstrates some code for review testing
+func TestFunction() {
+\t// TODO: This could be improved
```

**Format:**
- Standard unified diff format
- Shows context around the commented line
- Includes line numbers
- Escaped special characters

### Author Association Values

**From Testing:**
- `OWNER` - Repository owner
- Other values (not tested):
  - `COLLABORATOR` - Has write access
  - `CONTRIBUTOR` - Has contributed before
  - `FIRST_TIME_CONTRIBUTOR` - First contribution
  - `FIRST_TIMER` - First GitHub contribution
  - `MEMBER` - Organization member
  - `NONE` - No association

## Reaction Groups Structure

### Complete reactionGroups Array

**Always Returns All 8 Reaction Types:**
```json
"reactionGroups": [
  {
    "content": "THUMBS_UP",
    "createdAt": null,
    "users": {
      "totalCount": 0,
      "nodes": []
    },
    "viewerHasReacted": false
  },
  {
    "content": "ROCKET",
    "createdAt": "2025-11-02T21:43:33Z",
    "users": {
      "totalCount": 1,
      "nodes": [
        {
          "login": "hamishmorgan"
        }
      ]
    },
    "viewerHasReacted": true
  },
  // ... 6 more reaction types
]
```

### Key Findings

**Always Present:**
- ✅ ALL 8 reaction types in every response
- ✅ Even reactions with 0 count are included
- ✅ `createdAt` is null when totalCount is 0
- ✅ `createdAt` is timestamp of first reaction when count > 0

**Viewer Context:**
- `viewerHasReacted: true` - Current user has reacted with this emoji
- `viewerHasReacted: false` - Current user hasn't reacted

**User Lists:**
- `users.totalCount` - Total number of reactions
- `users.nodes` - Array of users who reacted
- Need `first` parameter for pagination

**Implications for gh-talk:**
- Can show reaction counts without extra queries
- Easy to detect if user has reacted
- Can display all reactions or filter to non-zero
- createdAt useful for "first to react" scenarios

## Mutation Responses

### addReaction Response

**Request:**
```graphql
mutation {
  addReaction(input: {
    subjectId: "IC_kwDOQN97u87PVA8l"
    content: THUMBS_UP
  }) { ... }
}
```

**Response:**
```json
{
  "data": {
    "addReaction": {
      "reaction": {
        "id": "REA_lALOQN97u87PVA8lzhKY8PA",
        "content": "THUMBS_UP",
        "user": {
          "login": "hamishmorgan"
        },
        "createdAt": "2025-11-02T21:42:52Z"
      },
      "subject": {
        "id": "IC_kwDOQN97u87PVA8l"
      }
    }
  }
}
```

**Key Points:**
- ✅ Returns new reaction ID
- ✅ Confirms content type
- ✅ Includes timestamp
- ✅ Echoes subject ID
- ⚠️ Does NOT return updated reactionGroups (must query separately)

### addPullRequestReviewThreadReply Response

**Request:**
```graphql
mutation {
  addPullRequestReviewThreadReply(input: {
    pullRequestReviewThreadId: "PRRT_kwDOQN97u85gQeTN"
    body: "Good point! I will refactor this to use a constant."
  }) { ... }
}
```

**Response:**
```json
{
  "data": {
    "addPullRequestReviewThreadReply": {
      "comment": {
        "id": "PRRC_kwDOQN97u86UHqOo",
        "body": "Good point! I will refactor this to use a constant.",
        "author": {
          "login": "hamishmorgan"
        },
        "createdAt": "2025-11-02T21:43:42Z",
        "pullRequestReview": {
          "id": "PRR_kwDOQN97u87LMeGr"
        },
        "replyTo": {
          "id": "PRRC_kwDOQN97u86UHqK7"
        }
      }
    }
  }
}
```

**Key Points:**
- ✅ Returns new comment with ID
- ✅ Creates a new review (PRR_) automatically
- ✅ Sets replyTo to link comments
- ⚠️ Creates review even for simple replies
- 💡 This is why there can be many reviews per PR

### resolveReviewThread Response

**Request:**
```graphql
mutation {
  resolveReviewThread(input: {
    threadId: "PRRT_kwDOQN97u85gQeTN"
  }) { ... }
}
```

**Response:**
```json
{
  "data": {
    "resolveReviewThread": {
      "thread": {
        "id": "PRRT_kwDOQN97u85gQeTN",
        "isResolved": true,
        "resolvedBy": {
          "login": "hamishmorgan"
        },
        "comments": {
          "nodes": [
            {
              "id": "PRRC_kwDOQN97u86UHqK7",
              "body": "Consider using a constant for the TODO comment"
            },
            {
              "id": "PRRC_kwDOQN97u86UHqOo",
              "body": "Good point! I will refactor this to use a constant."
            }
          ]
        }
      }
    }
  }
}
```

**Key Points:**
- ✅ Returns updated thread
- ✅ Sets `isResolved: true`
- ✅ Records `resolvedBy` user
- ✅ Includes all comments in thread
- 💡 Instant effect (no delay)

### unresolveReviewThread Response

**Request:**
```graphql
mutation {
  unresolveReviewThread(input: {
    threadId: "PRRT_kwDOQN97u85gQeTN"
  }) { ... }
}
```

**Response:**
```json
{
  "data": {
    "unresolveReviewThread": {
      "thread": {
        "id": "PRRT_kwDOQN97u85gQeTN",
        "isResolved": false,
        "resolvedBy": null
      }
    }
  }
}
```

**Key Points:**
- ✅ Sets `isResolved: false`
- ✅ Clears `resolvedBy` to null
- ⚠️ viewerCanUnresolve becomes true only after resolution

### minimizeComment Response

**Request:**
```graphql
mutation {
  minimizeComment(input: {
    subjectId: "IC_kwDOQN97u87PVA8l"
    classifier: OUTDATED
  }) { ... }
}
```

**Response:**
```json
{
  "data": {
    "minimizeComment": {
      "minimizedComment": {
        "isMinimized": true,
        "minimizedReason": "outdated",
        "viewerCanMinimize": true
      }
    }
  }
}
```

**Key Points:**
- ✅ Sets `isMinimized: true`
- ✅ Records reason (lowercase: "outdated", "spam", etc.)
- ⚠️ Comment still exists (just collapsed in UI)
- ⚠️ Still counts in totals

### addPullRequestReview Response

**Request:**
```graphql
mutation {
  addPullRequestReview(input: {
    pullRequestId: "PR_kwDOQN97u86xFPR4"
    event: COMMENT
    body: "Testing review thread creation"
    comments: [{
      path: "test_file.go"
      position: 7
      body: "Consider using a constant"
    }]
  }) { ... }
}
```

**Response:**
```json
{
  "data": {
    "addPullRequestReview": {
      "pullRequestReview": {
        "id": "PRR_kwDOQN97u87LMeCy",
        "state": "COMMENTED",
        "body": "Testing review thread creation",
        "author": {
          "login": "hamishmorgan"
        },
        "comments": {
          "nodes": [
            {
              "id": "PRRC_kwDOQN97u86UHqK7",
              "body": "Consider using a constant",
              "path": "test_file.go",
              "position": 7,
              "diffHunk": "...",
              "pullRequestReview": {
                "id": "PRR_kwDOQN97u87LMeCy"
              }
            }
          ]
        }
      }
    }
  }
}
```

**Key Points:**
- ✅ Creates review AND comments atomically
- ✅ Review gets an ID
- ✅ Each comment gets an ID
- ✅ Comments automatically linked to review
- ⚠️ Position is diff position, not line number
- 💡 Creating thread requires creating a review first

## Error Responses

### Invalid ID Error

**Request:**
```graphql
mutation {
  resolveReviewThread(input: {
    threadId: "PRRT_invalid123"
  }) { ... }
}
```

**Response:**
```json
{
  "data": {
    "resolveReviewThread": null
  },
  "errors": [
    {
      "type": "NOT_FOUND",
      "path": ["resolveReviewThread"],
      "locations": [{"line": 3, "column": 3}],
      "message": "Could not resolve to a node with the global id of 'PRRT_invalid123'"
    }
  ]
}
```

**Error Structure:**
- `data.<mutation>: null` - Operation failed
- `errors` array with error objects
- `type` field: `NOT_FOUND`, `FORBIDDEN`, `UNPROCESSABLE`, etc.
- `path` shows which field failed
- `message` is human-readable error

### Common Error Types (Expected)

**NOT_FOUND:**
```json
{
  "type": "NOT_FOUND",
  "message": "Could not resolve to a node with the global id of 'X'"
}
```

**FORBIDDEN:**
```json
{
  "type": "FORBIDDEN",
  "message": "Resource not accessible by personal access token"
}
```

**UNPROCESSABLE:**
```json
{
  "type": "UNPROCESSABLE",
  "message": "Body can't be blank"
}
```

## Reviews vs Threads vs Comments

### Relationship Hierarchy

**Discovered Structure:**
```
PullRequest
├── reviews[] (PRR_)
│   ├── id
│   ├── state (COMMENTED, APPROVED, CHANGES_REQUESTED)
│   ├── body (review-level comment)
│   └── comments[] (PRRC_) - line comments in this review
│
└── reviewThreads[] (PRRT_)
    ├── id
    ├── isResolved
    ├── path + line
    └── comments[] (PRRC_) - all comments in thread
```

**Key Insight:**
- **Reviews** (PRR_) - Container for submission (approve, request changes, comment)
- **Threads** (PRRT_) - Grouping of comments by file location
- **Comments** (PRRC_) - Individual messages

**Multiple Reviews Can Contribute to One Thread:**
```
Thread PRRT_123 on line 7:
  ├── Comment PRRC_456 (from Review PRR_789)
  ├── Comment PRRC_101 (from Review PRR_111) - Reply
  └── Comment PRRC_202 (from Review PRR_333) - Another reply
```

**Why This Matters:**
- Each reply creates a new review
- Many reviews per PR is normal
- Threads group comments regardless of review
- gh-talk should focus on threads, not reviews

## Top-Level Comments vs Review Comments

### Top-Level PR Comment (IC_)

**Structure:**
```json
{
  "id": "IC_kwDOQN97u87PVA8l",
  "body": "This is a test comment...",
  "author": {
    "login": "hamishmorgan"
  },
  "isMinimized": true,
  "minimizedReason": "outdated",
  "reactionGroups": [...]
}
```

**Characteristics:**
- ID prefix: `IC_` (IssueComment)
- No `path` or `position` (not tied to code)
- No `replyTo` (top-level only)
- No `diffHunk`
- Can be minimized
- Can have reactions
- Same reactionGroups structure

**Differences from Review Comments:**
- Not part of a review thread
- Not in `reviewThreads` collection
- In `comments` collection instead
- Cannot be resolved (no thread)
- Created via different API call

## Position vs Line Number

### Critical Discovery

**Position ≠ Line Number**

**From Testing:**
```json
{
  "position": 7,           // Position in diff (1-based)
  "originalPosition": 7,   // Original position when created
  "line": 7,              // Current line in file
  "diffHunk": "..."       // Context showing the code
}
```

**What's the Difference?**
- `position` - Index in the diff (changes as diff changes)
- `line` - Actual line number in the file
- `originalPosition` - Position when comment was created

**When They Differ:**
- File is modified after comment
- Lines added/removed above comment
- `position` updates, `line` updates
- `originalPosition` stays same

**Implications:**
- Use `line` for display (user-facing)
- Use `position` for API calls (creating comments)
- Track `originalPosition` for outdated detection

## Permission Fields

### Viewer Permissions on Threads

**From Testing:**
```json
{
  "viewerCanResolve": true,      // Can mark as resolved
  "viewerCanUnresolve": false,   // Can unresolve (changes after resolving)
  "viewerCanReply": true         // Can add comments
}
```

**When Testing Resolved Thread:**
```json
{
  "viewerCanResolve": false,     // Already resolved
  "viewerCanUnresolve": true,    // Now can unresolve
  "viewerCanReply": true         // Still can reply
}
```

**Permission Rules:**
- Can resolve: PR author, comment author, or repo admin
- Can unresolve: Same as resolve (but only if resolved)
- Can reply: Anyone with read access
- Permissions are dynamic based on state

### Viewer Permissions on Comments

**From Testing:**
```json
{
  "viewerCanReact": true,
  "viewerCanUpdate": true,
  "viewerCanDelete": true,
  "viewerCanMinimize": true
}
```

**Permission Logic:**
- Can react: Always true for authenticated users
- Can update: Comment author only
- Can delete: Comment author or repo admin
- Can minimize: Repo admin or maintainer

## Query Performance & Costs

### Query Cost Analysis

**Simple Thread List (tested):**
```graphql
query {
  repository(...) {
    pullRequest(...) {
      reviewThreads(first: 10) {
        nodes { id isResolved path line }
      }
    }
  }
}
```
**Estimated Cost:** ~10-20 points

**Detailed Thread Query (tested):**
```graphql
query {
  repository(...) {
    pullRequest(...) {
      reviewThreads(first: 20) {
        nodes {
          id
          # ... all fields
          comments(first: 20) {
            nodes {
              # ... all fields
              reactionGroups {
                users(first: 10) { ... }
              }
            }
          }
        }
      }
    }
  }
}
```
**Estimated Cost:** ~50-100 points (depends on data)

### Pagination Behavior

**Tested:**
```json
"reviewThreads": {
  "totalCount": 2,
  "nodes": [...]
}
```

**Key Fields:**
- `totalCount` - Total number available
- `nodes` - Array of current page
- `pageInfo` - Pagination cursors (not tested, but required for next page)

**Required Pattern:**
```graphql
reviewThreads(first: 20, after: $cursor) {
  pageInfo {
    hasNextPage
    endCursor
  }
  nodes { ... }
}
```

## Important Discoveries

### 1. **Thread ID is NOT Derivable**

❌ Cannot generate thread ID from:
- PR number
- File path
- Line number
- Comment URL

✅ Must query to get thread IDs

**Implication for gh-talk:**
- Must cache thread data
- Need short ID mapping (thread #1, #2, etc.)
- Or accept long IDs in commands

### 2. **Every Reply Creates a Review**

**Unexpected Behavior:**
- Replying to a thread creates a new review
- One PR can have 50+ reviews from replies
- Reviews are just containers

**Implication:**
- Reviews are less important than threads
- gh-talk should focus on threads
- Don't confuse users with review counts

### 3. **ReactionGroups Always Complete**

**Behavior:**
- ALL 8 reaction types always present
- Zero-count reactions included
- Simplifies logic (no null checks)

**Implication:**
- Easy to display "👍 0  🎉 0  ❤️ 3"
- Can show all or filter non-zero
- viewerHasReacted is convenient

### 4. **Minimize is Not Hide**

**Behavior:**
- Minimized comments still exist
- Still in query results
- Just `isMinimized: true`
- Collapsed in UI only

**Implication:**
- Need to filter client-side if hiding
- Or include minimized status in display
- Can't actually "delete" comments

### 5. **Position System is Complex**

**Fields:**
- `position` - Current position in diff
- `originalPosition` - Position when created
- `line` - Current line in file
- Can all be different!

**Implication:**
- Display `line` to users
- Track when they differ (outdated indicator)
- Understand diff position for API calls

## Test Data Summary

**Created in PR #1:**
- ✅ 1 Pull Request
- ✅ 3 Reviews
- ✅ 2 Review Threads
- ✅ 3 Review Comments (1 original + 1 reply + 1 new thread)
- ✅ 1 Top-Level Comment (minimized)
- ✅ 2 Reactions (THUMBS_UP on IC_, ROCKET on PRRC_)
- ✅ 1 Resolved thread (then unresolved for testing)

**IDs Collected:**
```
PR:      PR_kwDOQN97u86xFPR4
PRRT:    PRRT_kwDOQN97u85gQeTN (thread 1)
PRRT:    PRRT_kwDOQN97u85gQecu (thread 2)
PRRC:    PRRC_kwDOQN97u86UHqK7 (comment 1)
PRRC:    PRRC_kwDOQN97u86UHqOo (reply to comment 1)
PRRC:    PRRC_kwDOQN97u86UHqWJ (comment 2)
PRR:     PRR_kwDOQN97u87LMeCy (review 1)
PRR:     PRR_kwDOQN97u87LMeGr (review 2, auto-created)
PRR:     PRR_kwDOQN97u87LMePc (review 3)
IC:      IC_kwDOQN97u87PVA8l (top-level comment)
REA:     REA_lALOQN97u87PVA8lzhKY8PA (thumbs up)
REA:     REA_lATOQN97u86UHqK7zhQo4hY (rocket)
```

## Design Implications

### For ID Handling

**Options:**
1. **Use Full Node IDs** (cumbersome but accurate)
   ```bash
   gh talk reply PRRT_kwDOQN97u85gQeTN "message"
   ```

2. **Short ID Mapping** (user-friendly)
   ```bash
   gh talk list threads
   # 1. test_file.go:7  - Consider using constant
   # 2. test_file.go:14 - Loop optimization
   
   gh talk reply 1 "message"  # Maps to PRRT_...
   ```

3. **Interactive Selection** (best UX)
   ```bash
   gh talk reply
   # ? Select thread:
   #   > 1. test_file.go:7  - Consider using constant
   #     2. test_file.go:14 - Loop optimization
   
   # Type message: ...
   ```

**Recommendation:** Support all three
- Short IDs for interactive use
- Full IDs for scripting
- Interactive for exploratory use

### For Data Storage

**Must Cache:**
- Thread ID to short ID mapping
- Thread metadata (path, line, preview)
- Cache per PR (invalidate on changes)

**Cache Structure:**
```json
{
  "pr": 1,
  "threads": {
    "1": "PRRT_kwDOQN97u85gQeTN",
    "2": "PRRT_kwDOQN97u85gQecu"
  },
  "metadata": {
    "PRRT_kwDOQN97u85gQeTN": {
      "path": "test_file.go",
      "line": 7,
      "preview": "Consider using constant"
    }
  }
}
```

### For Display

**Reaction Display:**
```
# Option A: Show all (noisy)
👍 0  👎 0  😄 0  🎉 0  😕 0  ❤️ 0  🚀 1  👀 0

# Option B: Show non-zero (clean)
🚀 1

# Option C: Show with names (verbose)
🚀 hamishmorgan
```

**Recommendation:** Option B by default, Option C with `--verbose`

### For Threading

**Display Thread:**
```
Thread PRRT_...gQeTN (test_file.go:7) [resolved]
│
├─ hamishmorgan: Consider using a constant for the TODO comment
│  🚀 1
│
└─ hamishmorgan: Good point! I will refactor this to use a constant.
   (no reactions)
```

## Limitations Confirmed

### 1. **No Server-Side Filtering**

**Tested:**
- Cannot filter `reviewThreads(isResolved: false)` in query
- Must fetch all and filter client-side
- Confirmed limitation from API docs

### 2. **Pagination Required**

**Tested:**
- Must include `first` or `last` for all connections
- Error if missing: `MISSING_PAGINATION_BOUNDARIES`
- Cannot fetch "all" in one query

### 3. **Nested Pagination Complexity**

**Structure:**
```graphql
reviewThreads(first: 20) {
  nodes {
    comments(first: 50) {
      nodes {
        reactionGroups {
          users(first: 10) {  // Must paginate here too
            nodes { ... }
          }
        }
      }
    }
  }
}
```

**Implication:**
- Three levels of pagination
- Complex cursor management
- May need multiple queries for large threads

## Best Practices from Testing

### 1. **Always Request permissions**

```graphql
{
  viewerCanResolve
  viewerCanUnresolve
  viewerCanReply
}
```

**Why:** Prevents errors, enables graceful degradation

### 2. **Request Both IDs**

```graphql
{
  id          // Node ID for GraphQL
  databaseId  // For reference/debugging
}
```

**Why:** Helps with debugging, migration, and API compatibility

### 3. **Include Pagination Info**

```graphql
{
  totalCount
  pageInfo {
    hasNextPage
    endCursor
  }
}
```

**Why:** Know when more data exists, how to fetch it

### 4. **Get Diff Context**

```graphql
{
  diffHunk  // Code context
  path
  line
}
```

**Why:** Users need context to understand comment location

### 5. **Check Minimized Status**

```graphql
{
  isMinimized
  minimizedReason
}
```

**Why:** Hidden comments affect UX, need to handle appropriately

## Issue Data Structures

### Issue vs Pull Request Differences

**Key Insight:** Issues are simpler than PRs - no review threads, no diff context, but otherwise very similar comment/reaction model.

### Real Issue ID Format

**From Testing:**

| Type | Prefix | Example | Notes |
|------|--------|---------|-------|
| Issue | `I_` | `I_kwDOQN97u87VYpUq` | Issue Node ID |
| Issue Comment | `IC_` | `IC_kwDOQN97u87PVCb0` | Same as top-level PR comment |
| Label | `LA_` | `LA_kwDOQN97u88AAAACOo1ePQ` | Label ID |

**Observation:**
- Issue IDs use `I_` prefix (vs `PR_` for pull requests)
- Issue comments use `IC_` prefix (same as top-level PR comments)
- Same base64-encoded format (20-30 characters)

### Complete Issue Object

**Real Response from Testing:**
```json
{
  "id": "I_kwDOQN97u87VYpUq",
  "number": 2,
  "title": "Test Issue for API Exploration",
  "url": "https://github.com/hamishmorgan/gh-talk/issues/2",
  "state": "OPEN",
  "body": "This issue is for testing...",
  "createdAt": "2025-11-02T21:51:08Z",
  "updatedAt": "2025-11-02T21:52:08Z",
  "closed": false,
  "closedAt": null,
  "stateReason": null,
  "author": {
    "login": "hamishmorgan"
  },
  "authorAssociation": "OWNER",
  "viewerCanReact": true,
  "viewerCanUpdate": true,
  "viewerSubscription": "SUBSCRIBED",
  "labels": {
    "totalCount": 1,
    "nodes": [
      {
        "id": "LA_kwDOQN97u88AAAACOo1ePQ",
        "name": "documentation",
        "color": "0075ca",
        "description": "Improvements or additions to documentation"
      }
    ]
  },
  "assignees": {
    "totalCount": 0,
    "nodes": []
  },
  "participants": {
    "totalCount": 1,
    "nodes": [
      {
        "login": "hamishmorgan"
      }
    ]
  },
  "reactionGroups": [...],
  "comments": {
    "totalCount": 2,
    "nodes": [...]
  }
}
```

### Issue-Specific Fields

**Not in PRs:**
- `stateReason` - Why issue was closed
  - `COMPLETED` - Work is done
  - `NOT_PLANNED` - Won't be worked on
  - `REOPENED` - Issue was reopened
  - `null` - Issue is open
- `participants` - Users who commented or were mentioned
- No `reviewThreads` (issues don't have code review)
- No `files` or `diff` fields

**State Transitions:**
```
OPEN → CLOSED (stateReason: COMPLETED or NOT_PLANNED)
CLOSED → OPEN (stateReason: REOPENED)
```

### Issue Comment Object

**Structure (Identical to Top-Level PR Comment):**
```json
{
  "id": "IC_kwDOQN97u87PVCb0",
  "databaseId": 3478398708,
  "body": "This is the first test comment...",
  "createdAt": "2025-11-02T21:51:27Z",
  "updatedAt": "2025-11-02T21:51:27Z",
  "author": {
    "login": "hamishmorgan"
  },
  "authorAssociation": "OWNER",
  "isMinimized": false,
  "minimizedReason": null,
  "viewerCanReact": true,
  "viewerCanUpdate": true,
  "viewerCanDelete": true,
  "viewerCanMinimize": true,
  "reactionGroups": [
    {
      "content": "HEART",
      "createdAt": "2025-11-02T21:51:57Z",
      "users": {
        "totalCount": 1,
        "nodes": [
          {
            "login": "hamishmorgan"
          }
        ]
      },
      "viewerHasReacted": true
    },
    // ... all 8 reaction types
  ]
}
```

**Key Observations:**
- ✅ Same structure as top-level PR comments
- ✅ Same `IC_` prefix for comments
- ✅ Same reaction system (8 types, always present)
- ✅ Same minimization capability
- ✅ Same permission fields
- ❌ No `path` or `position` (no code context)
- ❌ No `diffHunk` (issues aren't diffs)
- ❌ No `replyTo` field (no threading in issues)

### Issue Reactions

**Issue Body Can Have Reactions:**
```json
{
  "id": "I_kwDOQN97u87VYpUq",
  "reactionGroups": [
    {
      "content": "EYES",
      "createdAt": "2025-11-02T21:52:06Z",
      "users": {
        "totalCount": 1,
        "nodes": [
          {
            "login": "hamishmorgan"
          }
        ]
      },
      "viewerHasReacted": true
    },
    // ... all 8 types
  ]
}
```

**What Can Be Reacted To:**
- ✅ Issue body (the main issue description)
- ✅ Issue comments
- ✅ Same 8 reaction types as PRs
- ✅ Same structure and behavior

### Issue Comments vs PR Review Comments

**Similarities:**
- ✅ Same ID format (`IC_` prefix)
- ✅ Same reaction system
- ✅ Same minimization
- ✅ Same permission model
- ✅ Same creation/update timestamps
- ✅ Same author association

**Differences:**

| Feature | Issue Comments | PR Review Comments |
|---------|---------------|-------------------|
| ID Prefix | `IC_` | `IC_` (top-level) or `PRRC_` (review) |
| Code Context | ❌ No path/position | ✅ path, position, diffHunk |
| Threading | ❌ Flat list only | ✅ Review threads (PRRT_) |
| Resolution | ❌ Cannot resolve | ✅ Can resolve threads |
| Location | ❌ No file reference | ✅ File + line number |
| Review | ❌ No review parent | ✅ Part of review (PRR_) |

### Issue Mutations

#### Add Issue Comment

**Mutation:**
```graphql
mutation {
  addComment(input: {
    subjectId: "I_kwDOQN97u87VYpUq"
    body: "Comment text"
  }) {
    commentEdge {
      node {
        id
        body
        createdAt
      }
    }
  }
}
```

**Note:** Uses generic `addComment` mutation (works for issues and PRs)

#### Update Issue Comment

**Mutation:**
```graphql
mutation {
  updateIssueComment(input: {
    id: "IC_kwDOQN97u87PVCb0"
    body: "Updated comment text"
  }) {
    issueComment {
      id
      body
      updatedAt
    }
  }
}
```

#### Delete Issue Comment

**Mutation:**
```graphql
mutation {
  deleteIssueComment(input: {
    id: "IC_kwDOQN97u87PVCb0"
  }) {
    clientMutationId
  }
}
```

**Note:** Permanent deletion (unlike minimize which just hides)

#### Close/Reopen Issue

**Close:**
```graphql
mutation {
  closeIssue(input: {
    issueId: "I_kwDOQN97u87VYpUq"
    stateReason: COMPLETED  # or NOT_PLANNED
  }) {
    issue {
      id
      state
      stateReason
      closedAt
    }
  }
}
```

**Reopen:**
```graphql
mutation {
  reopenIssue(input: {
    issueId: "I_kwDOQN97u87VYpUq"
  }) {
    issue {
      id
      state
      stateReason
    }
  }
}
```

**State After Close:**
```json
{
  "state": "CLOSED",
  "closed": true,
  "closedAt": "2025-11-02T21:52:34Z",
  "stateReason": "COMPLETED"
}
```

**State After Reopen:**
```json
{
  "state": "OPEN",
  "closed": false,
  "closedAt": null,
  "stateReason": "REOPENED"
}
```

### Issue Labels

**Label Structure:**
```json
{
  "id": "LA_kwDOQN97u88AAAACOo1ePQ",
  "name": "documentation",
  "color": "0075ca",
  "description": "Improvements or additions to documentation"
}
```

**Key Fields:**
- `id` - Label Node ID (LA_ prefix)
- `name` - Label text
- `color` - Hex color (6 digits, no #)
- `description` - Optional description

### Participants Field

**Unique to Issues:**
```json
{
  "participants": {
    "totalCount": 1,
    "nodes": [
      {
        "login": "hamishmorgan"
      }
    ]
  }
}
```

**Definition:** Users who have:
- Created the issue
- Commented on the issue
- Been mentioned in the issue
- Been assigned to the issue

**Use Case:** Tracking who's involved in the conversation

### Subscription Status

**viewerSubscription Field:**
- `SUBSCRIBED` - Receiving notifications
- `UNSUBSCRIBED` - Not receiving notifications
- `IGNORED` - Explicitly ignoring

**From Testing:**
```json
{
  "viewerSubscription": "SUBSCRIBED"
}
```

**Note:** Auto-subscribed when you create/comment on issue

## Issue vs PR: What gh-talk Supports

### Issue Commands

**Supported:**
```bash
gh talk list comments --issue 2           # List issue comments
gh talk react IC_... 👍                   # React to issue comment
gh talk react I_... 🎉                    # React to issue body
gh talk hide IC_... --reason spam         # Minimize comment
gh talk show 2 --type issue               # Show issue details
```

**NOT Supported (Issue Limitations):**
```bash
# ❌ No review threads on issues
gh talk list threads --issue 2            # N/A - issues don't have threads

# ❌ No thread resolution
gh talk resolve ...                       # N/A - only for PR threads

# ❌ No code-specific comments
# Issues are discussions, not code reviews
```

### Unified Comment Model

**Critical Discovery:**
Both Issues and PRs use `IC_` for top-level comments:

```
Issue Comment (IC_):
- Created via: gh issue comment
- Structure: Same as PR top-level comment
- Capabilities: React, minimize, update, delete

PR Top-Level Comment (IC_):
- Created via: gh pr comment
- Structure: Identical to issue comment
- Capabilities: Same

PR Review Comment (PRRC_):
- Created via: addPullRequestReview mutation
- Structure: Extended with path/position/diffHunk
- Capabilities: React, minimize, AND part of threads
```

**Implication for gh-talk:**
- Single code path for issue/PR comments
- Check parent type (Issue vs PullRequest)
- Treat issue comments like PR top-level comments
- Review comments (PRRC_) are special case

## Resolved vs Unresolved Thread Examples

### Current Test State

**PR #1 has 4 threads with mixed resolution status:**

**Resolved Threads (2):**
```json
{
  "id": "PRRT_kwDOQN97u85gQfgh",
  "path": "test_file.go",
  "line": 10,
  "isResolved": true,
  "resolvedBy": {
    "login": "hamishmorgan"
  },
  "viewerCanResolve": false,
  "viewerCanUnresolve": true,
  "comments": {
    "totalCount": 2,
    "nodes": [
      {
        "id": "PRRC_kwDOQN97u86UHrl7",
        "body": "Good naming for variables x, y, z",
        "replyTo": null
      },
      {
        "id": "PRRC_kwDOQN97u86UHroG",
        "body": "Thanks for the positive feedback! 👍",
        "replyTo": {
          "id": "PRRC_kwDOQN97u86UHrl7"
        }
      }
    ]
  }
}
```

```json
{
  "id": "PRRT_kwDOQN97u85gQfgi",
  "path": "test_file.go",
  "line": 18,
  "isResolved": true,
  "resolvedBy": {
    "login": "hamishmorgan"
  },
  "viewerCanResolve": false,
  "viewerCanUnresolve": true,
  "comments": {
    "totalCount": 1,
    "nodes": [
      {
        "id": "PRRC_kwDOQN97u86UHrl9",
        "body": "Consider extracting this condition to a named variable for clarity",
        "replyTo": null
      }
    ]
  }
}
```

**Unresolved Threads (2):**
```json
{
  "id": "PRRT_kwDOQN97u85gQeTN",
  "path": "test_file.go",
  "line": 7,
  "isResolved": false,
  "resolvedBy": null,
  "viewerCanResolve": true,
  "viewerCanUnresolve": false,
  "comments": {
    "totalCount": 2,
    "nodes": [
      {
        "id": "PRRC_kwDOQN97u86UHqK7",
        "body": "Consider using a constant for the TODO comment",
        "replyTo": null
      },
      {
        "id": "PRRC_kwDOQN97u86UHqOo",
        "body": "Good point! I will refactor this to use a constant.",
        "replyTo": {
          "id": "PRRC_kwDOQN97u86UHqK7"
        }
      }
    ]
  }
}
```

```json
{
  "id": "PRRT_kwDOQN97u85gQecu",
  "path": "test_file.go",
  "line": 14,
  "isResolved": false,
  "resolvedBy": null,
  "viewerCanResolve": true,
  "viewerCanUnresolve": false,
  "comments": {
    "totalCount": 1,
    "nodes": [
      {
        "id": "PRRC_kwDOQN97u86UHqWJ",
        "body": "This loop could be optimized using a range",
        "replyTo": null
      }
    ]
  }
}
```

### Key Observations: Resolved vs Unresolved

**Resolved Thread Characteristics:**
- ✅ `isResolved: true`
- ✅ `resolvedBy` contains user who resolved it
- ✅ `viewerCanResolve: false` (already resolved)
- ✅ `viewerCanUnresolve: true` (can reopen)
- ✅ Comments remain accessible
- ✅ Thread data unchanged except resolution status

**Unresolved Thread Characteristics:**
- ✅ `isResolved: false`
- ✅ `resolvedBy: null` (no resolver)
- ✅ `viewerCanResolve: true` (can resolve)
- ✅ `viewerCanUnresolve: false` (nothing to unresolve)
- ✅ Active conversation state

**Important:** Resolution is just a status flag - comments and data remain intact

### Display Implications

**List View Should Show:**
```
UNRESOLVED THREADS (2):
  1. test_file.go:7   Consider using a constant... (2 comments)
  2. test_file.go:14  This loop could be optimized... (1 comment)

RESOLVED THREADS (2):
  3. test_file.go:10  Good naming for variables... (2 comments) ✓ hamishmorgan
  4. test_file.go:18  Consider extracting this... (1 comment) ✓ hamishmorgan
```

**Filtering:**
```bash
# Show only unresolved (default for active work)
gh talk list threads --unresolved

# Show only resolved (to verify fixes)
gh talk list threads --resolved

# Show all (see everything)
gh talk list threads --all
```

## Test Data Summary (Complete)

**Created in Testing:**

**PR #1:**
- 1 Pull Request (`PR_kwDOQN97u86xFPR4`)
- 5 Reviews (`PRR_...`) - including auto-created from replies
- 4 Review Threads (`PRRT_...`) - 2 resolved, 2 unresolved
- 6 Review Comments (`PRRC_...`)
- 1 Top-Level Comment (`IC_...`, minimized)
- 2 Reactions on PR comments

**Issue #2:**
- 1 Issue (`I_kwDOQN97u87VYpUq`)
- 2 Issue Comments (`IC_...`)
- 3 Reactions (1 on issue body, 2 on comments)
- 1 Label (`LA_...`)
- 1 Minimized comment
- State changes (OPEN → CLOSED → OPEN)

**All IDs Collected:**
```
# Pull Requests
PR:      PR_kwDOQN97u86xFPR4

# Review Threads (4 total: 2 unresolved, 2 resolved)
PRRT:    PRRT_kwDOQN97u85gQeTN (thread 1, line 7, UNRESOLVED)
PRRT:    PRRT_kwDOQN97u85gQecu (thread 2, line 14, UNRESOLVED)
PRRT:    PRRT_kwDOQN97u85gQfgh (thread 3, line 10, RESOLVED)
PRRT:    PRRT_kwDOQN97u85gQfgi (thread 4, line 18, RESOLVED)

# Review Comments (6 total)
PRRC:    PRRC_kwDOQN97u86UHqK7 (thread 1, original comment)
PRRC:    PRRC_kwDOQN97u86UHqOo (thread 1, reply)
PRRC:    PRRC_kwDOQN97u86UHqWJ (thread 2, original comment)
PRRC:    PRRC_kwDOQN97u86UHrl7 (thread 3, original comment)
PRRC:    PRRC_kwDOQN97u86UHroG (thread 3, reply)
PRRC:    PRRC_kwDOQN97u86UHrl9 (thread 4, original comment)

# Reviews (5 total - some auto-created from replies)
PRR:     PRR_kwDOQN97u87LMeCy (review 1, with thread 1 original)
PRR:     PRR_kwDOQN97u87LMeGr (review 2, auto-created from thread 1 reply)
PRR:     PRR_kwDOQN97u87LMePc (review 3, with thread 2)
PRR:     PRR_kwDOQN97u87LMfVg (review 4, with threads 3 and 4)
PRR:     [additional auto-created review from thread 3 reply]

# Top-Level PR Comment
IC:      IC_kwDOQN97u87PVA8l (top-level PR comment, minimized)

# Issues
I:       I_kwDOQN97u87VYpUq (issue #2)
IC:      IC_kwDOQN97u87PVCb0 (issue comment 1)
IC:      IC_kwDOQN97u87PVCcO (issue comment 2, minimized)
LA:      LA_kwDOQN97u88AAAACOo1ePQ (label: documentation)

# Reactions
REA:     REA_lALOQN97u87PVA8lzhKY8PA (THUMBS_UP on PR comment)
REA:     REA_lATOQN97u86UHqK7zhQo4hY (ROCKET on review comment)
REA:     REA_lALOQN97u87PVCb0zhKY8tQ (HEART on issue comment)
REA:     REA_lALOQN97u87PVCcOzhKY8tU (HOORAY on issue comment)
REA:     REA_lAHOQN97u87VYpUqzg3x5G8 (EYES on issue body)
```

## Key Differences: Issues vs PRs

### What's the Same

**Comment System:**
- ✅ Same `IC_` prefix for comments
- ✅ Same reaction system (8 types)
- ✅ Same minimization capability
- ✅ Same permission model
- ✅ Same creation/update timestamps
- ✅ Same author association

**Reactions:**
- ✅ Same for issue body and comments
- ✅ All 8 types always in reactionGroups
- ✅ Same mutation (addReaction/removeReaction)
- ✅ Same structure and behavior

### What's Different

**Issues DO NOT Have:**
- ❌ Review threads (`reviewThreads` field doesn't exist)
- ❌ Review comments (`PRRC_` type)
- ❌ Reviews (`PRR_` type)
- ❌ Code context (path, position, diffHunk)
- ❌ Thread resolution
- ❌ Diff-related fields

**Issues Have Unique:**
- ✅ `stateReason` field (COMPLETED, NOT_PLANNED, REOPENED)
- ✅ `participants` field (conversation participants)
- ✅ Simpler close/reopen workflow
- ✅ Can be transferred between repos

**PRs Have Unique:**
- ✅ Review threads and thread resolution
- ✅ Review comments with code context
- ✅ Reviews (approve, request changes)
- ✅ Diff information
- ✅ File and line numbers
- ✅ CI/CD status checks

## Implications for gh-talk

### Unified vs Separate Commands

**Option A: Unified Commands**
```bash
gh talk list comments              # Auto-detect PR or Issue
gh talk react IC_... 👍            # Works for both
gh talk hide IC_... --reason spam  # Works for both
```

**Option B: Separate Commands**
```bash
gh talk list threads               # PR only
gh talk list comments --pr 1       # PR comments
gh talk list comments --issue 2    # Issue comments
```

**Recommendation:** Hybrid approach
- Commands that work for both: unified (react, hide, show)
- PR-specific commands: explicit (list threads, resolve)
- Auto-detect when possible, allow explicit override

### Comment Type Detection

**Strategy:**
```go
// Detect comment type from ID prefix
func GetCommentType(id string) CommentType {
    switch {
    case strings.HasPrefix(id, "IC_"):
        return IssueComment  // Could be issue OR PR top-level
    case strings.HasPrefix(id, "PRRC_"):
        return ReviewComment  // PR review comment only
    case strings.HasPrefix(id, "I_"):
        return Issue
    case strings.HasPrefix(id, "PR_"):
        return PullRequest
    case strings.HasPrefix(id, "PRRT_"):
        return ReviewThread
    default:
        return Unknown
    }
}
```

**Challenge:** `IC_` is ambiguous (issue or PR)
- Must query to determine parent
- Or require context flag (--issue vs --pr)

### Issue-Specific Features

**Close with Reason:**
```bash
gh talk close issue 2 --reason completed
gh talk close issue 2 --reason "not planned"
```

**List Participants:**
```bash
gh talk show issue 2 --participants
# Shows: Users involved in conversation
```

**Filter by Label:**
```bash
gh talk list comments --issue 2 --label documentation
```

### What gh-talk Should Support for Issues

**Phase 1 (MVP):**
- ✅ List issue comments
- ✅ Add reactions to issue comments
- ✅ Add reactions to issue body
- ✅ Minimize/hide issue comments
- ✅ View issue details

**Phase 2:**
- ✅ Add comments to issues
- ✅ Edit comments
- ✅ Filter comments by author, date
- ✅ Show participants
- ✅ Close/reopen issues

**Phase 3:**
- ✅ Bulk comment operations
- ✅ Issue templates
- ✅ Label management
- ✅ Assignee management

**NOT Supported (Issue Limitations):**
- ❌ Review threads (issues don't have them)
- ❌ Thread resolution (no threads)
- ❌ Code-specific comments (no diff)
- ❌ File/line filtering (no code context)

## References

- **Test PR:** https://github.com/hamishmorgan/gh-talk/pull/1
- **Test Issue:** https://github.com/hamishmorgan/gh-talk/issues/2
- **Full PR Response:** `testdata/pr_full_response.json`
- **Full Issue Response:** `testdata/issue_full_response.json`
- **PR with Mixed States:** `testdata/pr_with_resolved_threads.json`
- **Test Date:** 2025-11-02
- **gh Version:** 2.x
- **GraphQL API Version:** v4

---

**Last Updated**: 2025-11-02  
**Test Environment**: Real GitHub repository with live API  
**Data Sources**: PR #1 and Issue #2 in hamishmorgan/gh-talk  
**Test Fixtures**: `testdata/` directory contains real API responses

