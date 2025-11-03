# GitHub API Capabilities for gh-talk

**Comprehensive guide to GitHub API features relevant to PR and Issue conversation management**

## Overview

`gh-talk` leverages the GitHub GraphQL API (v4) for all interactions with GitHub. This document details the available API capabilities, data structures, and limitations.

## API Choice: GraphQL vs REST

### Why GraphQL?

- **Efficient Queries**: Fetch exactly the data needed in a single request
- **Nested Data**: Retrieve threads, comments, reactions, and metadata in one query
- **Type Safety**: Strong schema with introspection support
- **Flexible Filtering**: Built-in support for complex filters
- **Better Rate Limits**: 5,000 points/hour vs REST's 5,000 requests/hour

### REST API Limitations

- Multiple requests needed for nested data (threads → comments → reactions)
- Over-fetching of data
- Less flexible filtering
- More complex pagination

## Authentication

### Methods

1. **Personal Access Tokens (PATs)** - Recommended for gh-talk
   - Scopes required: `repo`, `read:org`
   - Used via `gh` CLI's existing authentication

2. **GitHub Apps**
   - Fine-grained permissions
   - Independent of user account
   - Better for organization-wide tools

### Rate Limits

- **Authenticated Requests**: 5,000 points/hour
- **Unauthenticated**: 60 requests/hour (not usable for gh-talk)
- **GraphQL Cost**: Each query costs points based on complexity
  - Simple query: ~1 point
  - Complex nested query: ~50-100 points
  - Mutations: ~1 point each

### Rate Limit Headers

```
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4999
X-RateLimit-Reset: 1372700873
X-RateLimit-Used: 1
X-RateLimit-Resource: graphql
```

## Core Data Types

### PullRequestReviewThread

**Description**: A threaded list of comments for a given pull request.

**Key Fields:**

```graphql
type PullRequestReviewThread {
  id: ID!
  comments(first: Int): PullRequestReviewCommentConnection!
  isResolved: Boolean!
  isCollapsed: Boolean!
  isOutdated: Boolean!
  path: String!
  line: Int
  startLine: Int
  diffSide: DiffSide!
  resolvedBy: User
  viewerCanResolve: Boolean!
  viewerCanUnresolve: Boolean!
  viewerCanReply: Boolean!
  pullRequest: PullRequest!
  repository: Repository!
}
```

**Important Notes:**

- `id` is a Node ID (base64-encoded), not a database ID
- `isResolved` indicates if the thread is marked as resolved
- `isCollapsed` shows if the thread is visually collapsed
- `isOutdated` means the diff has changed since the comment
- Permissions are per-viewer (`viewerCan*` fields)

### PullRequestReviewComment

**Description**: A comment on a pull request review.

**Key Fields:**

```graphql
type PullRequestReviewComment implements Comment {
  id: ID!
  author: Actor
  body: String!
  createdAt: DateTime!
  updatedAt: DateTime!
  reactions(first: Int): ReactionConnection
  reactionGroups: [ReactionGroup!]
  isMinimized: Boolean!
  minimizedReason: String
  viewerCanReact: Boolean!
  viewerCanUpdate: Boolean!
  viewerCanDelete: Boolean!
  viewerCanMinimize: Boolean!
  pullRequestReview: PullRequestReview
}
```

### Reaction

**Description**: An emoji reaction to a comment.

**Available Reactions (ReactionContent enum):**

- `THUMBS_UP` - 👍 (`:+1:`)
- `THUMBS_DOWN` - 👎 (`:-1:`)
- `LAUGH` - 😄 (`:laugh:`)
- `HOORAY` - 🎉 (`:hooray:`)
- `CONFUSED` - 😕 (`:confused:`)
- `HEART` - ❤️ (`:heart:`)
- `ROCKET` - 🚀 (`:rocket:`)
- `EYES` - 👀 (`:eyes:`)

**ReactionGroup:**

```graphql
type ReactionGroup {
  content: ReactionContent!
  users(first: Int): ReactingUserConnection!
  viewerHasReacted: Boolean!
  createdAt: DateTime
}
```

### PullRequestReview

**Description**: A review on a pull request.

**Review States:**

- `PENDING` - Review is in draft
- `COMMENTED` - General feedback without approval/rejection
- `APPROVED` - Approved the changes
- `CHANGES_REQUESTED` - Requested changes
- `DISMISSED` - Review was dismissed

## GraphQL Queries

### List Review Threads

**Query:**

```graphql
query ListReviewThreads($owner: String!, $repo: String!, $pr: Int!) {
  repository(owner: $owner, name: $repo) {
    pullRequest(number: $pr) {
      reviewThreads(first: 100) {
        nodes {
          id
          isResolved
          isOutdated
          path
          line
          comments(first: 50) {
            nodes {
              id
              author {
                login
              }
              body
              createdAt
              reactions(first: 10) {
                nodes {
                  content
                  user {
                    login
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

**Filtering:**

- By resolution status (requires client-side filtering)
- By file path (client-side)
- By author (client-side)
- By date range (client-side)

**Cost**: ~50-100 points (depends on number of threads)

### Get Thread Details

**Query:**

```graphql
query GetThread($threadId: ID!) {
  node(id: $threadId) {
    ... on PullRequestReviewThread {
      id
      isResolved
      path
      line
      startLine
      diffSide
      resolvedBy {
        login
      }
      comments(first: 100) {
        nodes {
          id
          author {
            login
          }
          body
          createdAt
          updatedAt
          reactionGroups {
            content
            users(first: 10) {
              totalCount
              nodes {
                login
              }
            }
          }
        }
      }
      pullRequest {
        number
        title
        repository {
          nameWithOwner
        }
      }
    }
  }
}
```

**Cost**: ~20-40 points

## GraphQL Mutations

### 1. Resolve Review Thread

**Mutation:**

```graphql
mutation ResolveThread($threadId: ID!) {
  resolveReviewThread(input: {threadId: $threadId}) {
    thread {
      id
      isResolved
      resolvedBy {
        login
      }
    }
  }
}
```

**Requirements:**

- `threadId`: Node ID of the thread
- Viewer must have `viewerCanResolve: true`

**Cost**: 1 point

### 2. Unresolve Review Thread

**Mutation:**

```graphql
mutation UnresolveThread($threadId: ID!) {
  unresolveReviewThread(input: {threadId: $threadId}) {
    thread {
      id
      isResolved
    }
  }
}
```

**Requirements:**

- `threadId`: Node ID of the thread
- Viewer must have `viewerCanUnresolve: true`

**Cost**: 1 point

### 3. Reply to Review Thread

**Mutation:**

```graphql
mutation ReplyToThread($threadId: ID!, $body: String!) {
  addPullRequestReviewThreadReply(input: {
    pullRequestReviewThreadId: $threadId
    body: $body
  }) {
    comment {
      id
      body
      createdAt
      author {
        login
      }
    }
  }
}
```

**Requirements:**

- `pullRequestReviewThreadId`: Node ID of the thread
- `body`: Comment text (markdown supported)
- Viewer must have `viewerCanReply: true`
- Optional: `pullRequestReviewId` for pending reviews

**Cost**: 1 point

**Note**: If you want to reply AND resolve, you need two separate mutations.

### 4. Add Reaction

**Mutation:**

```graphql
mutation AddReaction($subjectId: ID!, $content: ReactionContent!) {
  addReaction(input: {
    subjectId: $subjectId
    content: $content
  }) {
    reaction {
      id
      content
      user {
        login
      }
    }
    subject {
      id
    }
  }
}
```

**Requirements:**

- `subjectId`: Node ID of comment or issue
- `content`: One of the 8 ReactionContent enum values
- Viewer must have `viewerCanReact: true`

**Cost**: 1 point

### 5. Remove Reaction

**Mutation:**

```graphql
mutation RemoveReaction($subjectId: ID!, $content: ReactionContent!) {
  removeReaction(input: {
    subjectId: $subjectId
    content: $content
  }) {
    reaction {
      id
      content
    }
    subject {
      id
    }
  }
}
```

**Cost**: 1 point

### 6. Minimize Comment

**Mutation:**

```graphql
mutation MinimizeComment($subjectId: ID!, $classifier: ReportedContentClassifiers!) {
  minimizeComment(input: {
    subjectId: $subjectId
    classifier: $classifier
  }) {
    minimizedComment {
      isMinimized
      minimizedReason
    }
  }
}
```

**Classifiers:**

- `SPAM`
- `ABUSE`
- `OFF_TOPIC`
- `OUTDATED`
- `DUPLICATE`
- `RESOLVED`

**Requirements:**

- `subjectId`: Node ID of the comment
- Viewer must have `viewerCanMinimize: true` (usually repo write access)

**Cost**: 1 point

### 7. Unminimize Comment

**Mutation:**

```graphql
mutation UnminimizeComment($subjectId: ID!) {
  unminimizeComment(input: {
    subjectId: $subjectId
  }) {
    unminimizedComment {
      isMinimized
    }
  }
}
```

**Cost**: 1 point

### 8. Dismiss Pull Request Review

**Mutation:**

```graphql
mutation DismissReview($reviewId: ID!, $message: String!) {
  dismissPullRequestReview(input: {
    pullRequestReviewId: $reviewId
    message: $message
  }) {
    review {
      id
      state
    }
  }
}
```

**Requirements:**

- `pullRequestReviewId`: Node ID of the review
- `message`: Reason for dismissal (required)
- Viewer must have write access to repository

**Cost**: 1 point

**Note**: This dismisses a review, not individual comments. Useful when all threads are resolved but review still shows as "Changes Requested".

## Limitations & Constraints

### GraphQL Limitations

1. **No Server-Side Filtering for Threads**
   - Cannot filter `reviewThreads` by `isResolved` in the query
   - Must fetch all threads and filter client-side
   - Workaround: Use pagination and cache results

2. **Introspection Limits**
   - Maximum 2 uses of introspection fields per query
   - Prevents dynamic schema exploration in single query

3. **Pagination Complexity**
   - Must use cursor-based pagination (`after`, `before`)
   - Cannot use offset-based pagination
   - Requires managing cursors for each nested connection

4. **No Bulk Mutations**
   - Cannot resolve multiple threads in one mutation
   - Must make separate mutations for each operation
   - Rate limits apply per mutation

### Permission Model

**Thread Resolution:**

- PR author can resolve any thread
- Comment author can resolve their own threads
- Repo admins can resolve any thread
- Others cannot resolve (even with write access)

**Comment Hiding:**

- Requires repository write access
- Cannot hide own comments
- Minimized comments still appear (just collapsed)

**Review Dismissal:**

- Requires repository write access
- Cannot dismiss own review
- Dismissed reviews remain visible in timeline

### Rate Limit Considerations

**Query Costs:**

- Simple thread list: ~50 points
- Detailed thread with full conversation: ~100 points
- 100 threads with details: ~1000 points (manageable)

**Mutation Costs:**

- Each operation: 1 point
- Resolving 50 threads: 50 points
- Adding 50 reactions: 50 points

**Strategy:**

- Cache thread data for 5 minutes
- Batch read operations where possible
- Avoid unnecessary refetches after mutations
- Monitor rate limit headers

### Data Consistency

**Eventually Consistent:**

- Reactions may not appear immediately
- Thread resolution status may lag slightly
- Comments appear quickly but reactions/updates lag

**No Transactional Mutations:**

- Cannot reply + resolve atomically
- Must handle partial failures
- No rollback mechanism

## Best Practices

### Query Optimization

1. **Request Only Needed Fields**

   ```graphql
   # Bad - over-fetching
   comments { nodes { ...AllFields } }
   
   # Good - selective
   comments { nodes { id author { login } body } }
   ```

2. **Use Fragments for Reusability**

   ```graphql
   fragment ThreadSummary on PullRequestReviewThread {
     id
     isResolved
     path
     line
   }
   ```

3. **Implement Pagination**

   ```graphql
   reviewThreads(first: 50, after: $cursor) {
     pageInfo {
       hasNextPage
       endCursor
     }
     nodes { ... }
   }
   ```

### Error Handling

**Common Errors:**

```json
{
  "errors": [{
    "type": "NOT_FOUND",
    "message": "Could not resolve to a node with the global id"
  }]
}
```

**Error Types:**

- `NOT_FOUND` - Invalid ID or deleted resource
- `FORBIDDEN` - Permission denied
- `UNPROCESSABLE` - Invalid input (e.g., empty body)
- `SERVICE_UNAVAILABLE` - GitHub API is down
- `RATE_LIMIT` - Exceeded rate limit

**Handling Strategy:**

- Check `errors` array in response
- Validate IDs before mutations
- Retry with exponential backoff for rate limits
- Fall back gracefully for permission errors

### Caching Strategy

**What to Cache:**

- Thread lists (5 minute TTL)
- Comment contents (10 minute TTL)
- Reaction counts (1 minute TTL)
- User information (1 hour TTL)

**What NOT to Cache:**

- `isResolved` status (changes frequently)
- `viewerCan*` permissions (context-dependent)
- Reaction lists (users may change frequently)

**Invalidation:**

- Clear cache after mutations
- Use ETags when available
- Implement lazy refresh for background updates

## Future API Capabilities

### Upcoming Features (GitHub Roadmap)

- **Suggestion Comments**: Apply code suggestions from CLI
- **Code Review Assignments**: Automated assignment via API
- **Thread Filters**: Server-side filtering by status
- **Bulk Operations**: Resolve multiple threads in one mutation

### Not Available (Limitations)

- **Thread Creation**: Cannot create new review threads via API (only via code review)
- **Review Submission**: Cannot submit reviews (approve/request changes) via API without web UI
- **File Annotations**: Cannot add file-level comments (only line-specific)
- **PR Creation from CLI**: Limited compared to web UI workflow

## Example Workflows

### Complete Thread Management Flow

```
1. List unresolved threads
   → Query: ListReviewThreads(isResolved: false)
   → Cost: ~50 points
   
2. User selects thread to address
   → Query: GetThread(threadId)
   → Cost: ~30 points
   
3. User replies to thread
   → Mutation: ReplyToThread(threadId, body)
   → Cost: 1 point
   
4. User resolves thread
   → Mutation: ResolveThread(threadId)
   → Cost: 1 point
   
Total: ~82 points per complete workflow
```

### Bulk Reaction Workflow

```
1. List comments needing reactions
   → Query: ListReviewThreads
   → Cost: ~50 points
   
2. Add reactions to multiple comments (n comments)
   → Mutation: AddReaction × n
   → Cost: n points
   
Total: ~50 + n points
```

## API Client Implementation

### Using go-gh Library

```go
package api

import (
    "github.com/cli/go-gh/v2/pkg/api"
)

// Get authenticated GraphQL client
func NewClient() (*api.GraphQLClient, error) {
    client, err := api.DefaultGraphQLClient()
    if err != nil {
        return nil, fmt.Errorf("create client: %w", err)
    }
    return client, nil
}

// Execute query
func ExecuteQuery(client *api.GraphQLClient, query string, vars map[string]interface{}) error {
    var response struct {
        // Define response structure
    }
    
    err := client.Query("repo", query, &response, vars)
    if err != nil {
        return fmt.Errorf("execute query: %w", err)
    }
    
    return nil
}
```

### Error Handling Pattern

```go
type GraphQLError struct {
    Type    string `json:"type"`
    Message string `json:"message"`
}

func handleGraphQLErrors(errs []GraphQLError) error {
    if len(errs) == 0 {
        return nil
    }
    
    for _, err := range errs {
        switch err.Type {
        case "NOT_FOUND":
            return &NotFoundError{Message: err.Message}
        case "FORBIDDEN":
            return &PermissionError{Message: err.Message}
        case "RATE_LIMIT":
            return &RateLimitError{Message: err.Message}
        }
    }
    
    return fmt.Errorf("graphql error: %s", errs[0].Message)
}
```

## References

- [GitHub GraphQL API Documentation](https://docs.github.com/en/graphql)
- [GitHub GraphQL Explorer](https://docs.github.com/en/graphql/overview/explorer)
- [GraphQL Rate Limits](https://docs.github.com/en/graphql/overview/resource-limitations)
- [go-gh Library Documentation](https://github.com/cli/go-gh)
- [GitHub API Best Practices](https://docs.github.com/en/rest/guides/best-practices-for-using-the-rest-api)

---

**Last Updated**: 2025-11-02  
**API Version**: GitHub GraphQL API v4  
**Schema Introspection Date**: 2025-11-02
