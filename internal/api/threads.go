package api

import (
	"context"
	"fmt"
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

// ListThreads fetches all review threads for a pull request
func (c *Client) ListThreads(ctx context.Context, owner, name string, pr int) ([]Thread, error) {
	var query struct {
		Repository struct {
			PullRequest struct {
				ReviewThreads struct {
					Nodes []struct {
						ID          graphql.String
						IsResolved  graphql.Boolean
						IsCollapsed graphql.Boolean
						IsOutdated  graphql.Boolean
						Path        graphql.String
						Line        *graphql.Int
						StartLine   *graphql.Int
						DiffSide    graphql.String
						SubjectType graphql.String
						ResolvedBy  *struct {
							Login graphql.String
						}
						ViewerCanResolve   graphql.Boolean
						ViewerCanUnresolve graphql.Boolean
						ViewerCanReply     graphql.Boolean
						Comments           struct {
							TotalCount graphql.Int
							Nodes      []struct {
								ID        graphql.String
								Body      graphql.String
								CreatedAt string
								Author    struct {
									Login graphql.String
								}
								ReactionGroups []struct {
									Content graphql.String
									Users   struct {
										TotalCount graphql.Int
									} `graphql:"users(first: 1)"`
									ViewerHasReacted graphql.Boolean
								}
							}
						} `graphql:"comments(first: 50)"`
					}
				} `graphql:"reviewThreads(first: 100)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  graphQLString(owner),
		"name":   graphQLString(name),
		"number": graphQLInt(pr),
	}

	err := c.queryWithContext(ctx, "ListThreads", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("query threads: %w", err)
	}

	// Convert to our types
	threads := make([]Thread, 0, len(query.Repository.PullRequest.ReviewThreads.Nodes))
	for _, node := range query.Repository.PullRequest.ReviewThreads.Nodes {
		thread := Thread{
			ID:                 string(node.ID),
			IsResolved:         bool(node.IsResolved),
			IsCollapsed:        bool(node.IsCollapsed),
			IsOutdated:         bool(node.IsOutdated),
			Path:               string(node.Path),
			DiffSide:           string(node.DiffSide),
			SubjectType:        string(node.SubjectType),
			ViewerCanResolve:   bool(node.ViewerCanResolve),
			ViewerCanUnresolve: bool(node.ViewerCanUnresolve),
			ViewerCanReply:     bool(node.ViewerCanReply),
		}

		if node.Line != nil {
			thread.Line = int(*node.Line)
		}
		if node.StartLine != nil {
			thread.StartLine = int(*node.StartLine)
		}

		if node.ResolvedBy != nil {
			thread.ResolvedBy = &User{
				Login: string(node.ResolvedBy.Login),
			}
		}

		// Convert comments
		thread.Comments = make([]Comment, 0, len(node.Comments.Nodes))
		for _, c := range node.Comments.Nodes {
			comment := Comment{
				ID:   string(c.ID),
				Body: string(c.Body),
				Author: User{
					Login: string(c.Author.Login),
				},
			}

			// Parse CreatedAt
			if c.CreatedAt != "" {
				if t, err := time.Parse(time.RFC3339, c.CreatedAt); err == nil {
					comment.CreatedAt = t
				}
			}

			// Convert reaction groups
			comment.ReactionGroups = make([]ReactionGroup, 0, len(c.ReactionGroups))
			for _, rg := range c.ReactionGroups {
				if rg.Users.TotalCount > 0 {
					comment.ReactionGroups = append(comment.ReactionGroups, ReactionGroup{
						Content:          string(rg.Content),
						ViewerHasReacted: bool(rg.ViewerHasReacted),
						Users: ReactionUsers{
							TotalCount: int(rg.Users.TotalCount),
						},
					})
				}
			}

			thread.Comments = append(thread.Comments, comment)
		}

		threads = append(threads, thread)
	}

	return threads, nil
}
