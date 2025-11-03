package api

import (
	"context"
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"
)

// ReplyToThread adds a reply to a review thread
func (c *Client) ReplyToThread(ctx context.Context, threadID, body string) error {
	var mutation struct {
		AddPullRequestReviewThreadReply struct {
			Comment struct {
				ID graphql.String
			}
		} `graphql:"addPullRequestReviewThreadReply(input: $input)"`
	}

	type AddPullRequestReviewThreadReplyInput struct {
		PullRequestReviewThreadID graphql.ID     `json:"pullRequestReviewThreadId"`
		Body                      graphql.String `json:"body"`
	}

	variables := map[string]interface{}{
		"input": AddPullRequestReviewThreadReplyInput{
			PullRequestReviewThreadID: graphQLID(threadID),
			Body:                      graphQLString(body),
		},
	}

	err := c.mutateWithContext(ctx, "AddReply", &mutation, variables)
	if err != nil {
		return fmt.Errorf("add reply: %w", err)
	}

	return nil
}

// ResolveThread marks a review thread as resolved
func (c *Client) ResolveThread(ctx context.Context, threadID string) error {
	var mutation struct {
		ResolveReviewThread struct {
			Thread struct {
				ID         graphql.String
				IsResolved graphql.Boolean
			}
		} `graphql:"resolveReviewThread(input: $input)"`
	}

	type ResolveReviewThreadInput struct {
		ThreadID graphql.ID `json:"threadId"`
	}

	variables := map[string]interface{}{
		"input": ResolveReviewThreadInput{
			ThreadID: graphQLID(threadID),
		},
	}

	err := c.mutateWithContext(ctx, "ResolveThread", &mutation, variables)
	if err != nil {
		return fmt.Errorf("resolve thread: %w", err)
	}

	return nil
}

// UnresolveThread marks a review thread as unresolved
func (c *Client) UnresolveThread(ctx context.Context, threadID string) error {
	var mutation struct {
		UnresolveReviewThread struct {
			Thread struct {
				ID         graphql.String
				IsResolved graphql.Boolean
			}
		} `graphql:"unresolveReviewThread(input: $input)"`
	}

	type UnresolveReviewThreadInput struct {
		ThreadID graphql.ID `json:"threadId"`
	}

	variables := map[string]interface{}{
		"input": UnresolveReviewThreadInput{
			ThreadID: graphQLID(threadID),
		},
	}

	err := c.mutateWithContext(ctx, "UnresolveThread", &mutation, variables)
	if err != nil {
		return fmt.Errorf("unresolve thread: %w", err)
	}

	return nil
}

// AddReaction adds an emoji reaction to a comment
func (c *Client) AddReaction(ctx context.Context, subjectID, content string) error {
	var mutation struct {
		AddReaction struct {
			Reaction struct {
				ID graphql.String
			}
		} `graphql:"addReaction(input: $input)"`
	}

	type AddReactionInput struct {
		SubjectID graphql.ID     `json:"subjectId"`
		Content   graphql.String `json:"content"`
	}

	variables := map[string]interface{}{
		"input": AddReactionInput{
			SubjectID: graphQLID(subjectID),
			Content:   graphQLString(content),
		},
	}

	err := c.mutateWithContext(ctx, "AddReaction", &mutation, variables)
	if err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}

	return nil
}

// RemoveReaction removes an emoji reaction from a comment
func (c *Client) RemoveReaction(ctx context.Context, subjectID, content string) error {
	var mutation struct {
		RemoveReaction struct {
			Reaction struct {
				ID graphql.String
			}
		} `graphql:"removeReaction(input: $input)"`
	}

	type RemoveReactionInput struct {
		SubjectID graphql.ID     `json:"subjectId"`
		Content   graphql.String `json:"content"`
	}

	variables := map[string]interface{}{
		"input": RemoveReactionInput{
			SubjectID: graphQLID(subjectID),
			Content:   graphQLString(content),
		},
	}

	err := c.mutateWithContext(ctx, "RemoveReaction", &mutation, variables)
	if err != nil {
		return fmt.Errorf("remove reaction: %w", err)
	}

	return nil
}


