package api

import (
	"context"
	"fmt"
	"strings"

	graphql "github.com/cli/shurcooL-graphql"
)

// MinimizeComment hides/minimizes a comment
func (c *Client) MinimizeComment(ctx context.Context, commentID, classifier string) error {
	var mutation struct {
		MinimizeComment struct {
			MinimizedComment struct {
				IsMinimized     graphql.Boolean
				MinimizedReason graphql.String
			}
		} `graphql:"minimizeComment(input: $input)"`
	}

	type MinimizeCommentInput struct {
		SubjectID  graphql.ID     `json:"subjectId"`
		Classifier graphql.String `json:"classifier"`
	}

	variables := map[string]interface{}{
		"input": MinimizeCommentInput{
			SubjectID:  graphQLID(commentID),
			Classifier: graphQLString(classifier),
		},
	}

	err := c.mutateWithContext(ctx, "MinimizeComment", &mutation, variables)
	if err != nil {
		return fmt.Errorf("minimize comment: %w", err)
	}

	return nil
}

// UnminimizeComment unhides a comment
func (c *Client) UnminimizeComment(ctx context.Context, commentID string) error {
	var mutation struct {
		UnminimizeComment struct {
			UnminimizedComment struct {
				IsMinimized graphql.Boolean
			}
		} `graphql:"unminimizeComment(input: $input)"`
	}

	type UnminimizeCommentInput struct {
		SubjectID graphql.ID `json:"subjectId"`
	}

	variables := map[string]interface{}{
		"input": UnminimizeCommentInput{
			SubjectID: graphQLID(commentID),
		},
	}

	err := c.mutateWithContext(ctx, "UnminimizeComment", &mutation, variables)
	if err != nil {
		return fmt.Errorf("unminimize comment: %w", err)
	}

	return nil
}

// ParseClassifier converts user-friendly reason to GraphQL classifier
func ParseClassifier(reason string) (string, error) {
	reason = strings.ToUpper(strings.ReplaceAll(reason, "-", "_"))

	valid := []string{"SPAM", "ABUSE", "OFF_TOPIC", "OUTDATED", "DUPLICATE", "RESOLVED"}
	for _, v := range valid {
		if reason == v {
			return v, nil
		}
	}

	return "", fmt.Errorf("invalid reason: %s\n\nValid reasons: spam, abuse, off-topic, outdated, duplicate, resolved", reason)
}


