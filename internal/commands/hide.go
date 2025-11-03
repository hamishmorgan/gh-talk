package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var hideCmd = &cobra.Command{
	Use:   "hide <comment-id...>",
	Short: "Minimize/hide comments",
	Long: `Minimize (hide) one or more comments with a reason.

Arguments:
  comment-id...  One or more comment IDs (PRRC_... or IC_...)

Examples:
  # Hide single comment as spam
  gh talk hide IC_kwDOQN97u87PVA8l --reason spam

  # Hide multiple comments (bulk operation)
  gh talk hide PRRC_aaa PRRC_bbb PRRC_ccc --reason resolved

  # Hide as off-topic
  gh talk hide PRRC_kwDOQN97u86UHqK7 --reason off-topic`,
	Args: cobra.MinimumNArgs(1),
	RunE: runHide,
}

var unhideCmd = &cobra.Command{
	Use:   "unhide <comment-id>",
	Short: "Unhide a comment",
	Long: `Unhide (unminimize) a previously hidden comment.

Arguments:
  comment-id  Comment ID (PRRC_... or IC_...)

Examples:
  # Unhide a comment
  gh talk unhide IC_kwDOQN97u87PVA8l`,
	Args: cobra.ExactArgs(1),
	RunE: runUnhide,
}

func init() {
	hideCmd.Flags().String("reason", "off-topic", "Reason (spam, abuse, off-topic, outdated, duplicate, resolved)")
}

func runHide(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	commentIDs := args

	// Validate all comment IDs
	for _, commentID := range commentIDs {
		if !strings.HasPrefix(commentID, "PRRC_") && !strings.HasPrefix(commentID, "IC_") {
			return fmt.Errorf("invalid comment ID: %s\n\nExpected format: PRRC_... or IC_...", commentID)
		}
	}

	// Parse reason
	reason, _ := cmd.Flags().GetString("reason")
	classifier, err := api.ParseClassifier(reason)
	if err != nil {
		return err
	}

	// Create client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Hide each comment
	for _, commentID := range commentIDs {
		err = client.MinimizeComment(ctx, commentID, classifier)
		if err != nil {
			return fmt.Errorf("failed to hide %s: %w", commentID, err)
		}
		fmt.Printf("✓ Hidden comment %s (reason: %s)\n", commentID, strings.ToLower(classifier))
	}

	if len(commentIDs) > 1 {
		fmt.Printf("\n✓ Hidden %d comments\n", len(commentIDs))
	}

	return nil
}

func runUnhide(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	commentID := args[0]

	// Validate comment ID
	if !strings.HasPrefix(commentID, "PRRC_") && !strings.HasPrefix(commentID, "IC_") {
		return fmt.Errorf("invalid comment ID: %s\n\nExpected format: PRRC_... or IC_...", commentID)
	}

	// Create client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Unhide comment
	err = client.UnminimizeComment(ctx, commentID)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Unhidden comment %s\n", commentID)
	return nil
}
