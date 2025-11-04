package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up resolved thread comments",
	Long: `Hide comments in all resolved threads.

This command finds all resolved threads in a PR and hides their original
comments to declutter the conversation view.

Flags:
  --pr NUMBER      PR number (required if not in PR branch)
  --yes           Skip confirmation prompt

Examples:
  # Clean up resolved threads in PR 137
  gh talk cleanup --pr 137

  # Skip confirmation prompt
  gh talk cleanup --pr 137 --yes

  # Clean up current PR (if on PR branch)
  gh talk cleanup`,
	RunE: runCleanup,
}

func init() {
	cleanupCmd.Flags().Int("pr", 0, "PR number")
	cleanupCmd.Flags().Bool("yes", false, "Skip confirmation prompt")
}

func runCleanup(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get repository
	owner, name, err := getRepository(cmd)
	if err != nil {
		return err
	}

	// Get PR number
	prNum, err := getCurrentPR(cmd)
	if err != nil {
		return err
	}

	// Get yes flag
	skipConfirm, _ := cmd.Flags().GetBool("yes")

	// Create client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// List all threads
	threads, err := client.ListThreads(ctx, owner, name, prNum)
	if err != nil {
		return fmt.Errorf("failed to list threads: %w", err)
	}

	// Find resolved threads and collect comment IDs to hide
	var commentIDs []string
	resolvedCount := 0

	for _, thread := range threads {
		if thread.IsResolved {
			resolvedCount++
			// Add the first comment (original comment) from each resolved thread
			if len(thread.Comments) > 0 {
				commentIDs = append(commentIDs, thread.Comments[0].ID)
			}
		}
	}

	// Check if there's anything to clean up
	if len(commentIDs) == 0 {
		if resolvedCount == 0 {
			fmt.Println("No resolved threads found")
		} else {
			fmt.Println("No comments to hide (resolved threads may already be cleaned up)")
		}
		return nil
	}

	// Show summary
	fmt.Printf("Found %d resolved thread(s) with %d comment(s) to hide\n", resolvedCount, len(commentIDs))

	// Confirm unless --yes flag
	if !skipConfirm {
		fmt.Printf("\nAbout to hide %d comment(s) in resolved threads\n", len(commentIDs))
		fmt.Print("Continue? (Y/n): ")

		var response string
		// Ignore scan errors (empty input is acceptable)
		fmt.Scanln(&response) //nolint:errcheck

		if response != "" && response != "Y" && response != "y" && response != "yes" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	// Hide all comments
	hiddenCount := 0
	for _, commentID := range commentIDs {
		err := client.MinimizeComment(ctx, commentID, "RESOLVED")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to hide %s: %v\n", commentID, err)
			continue
		}
		hiddenCount++
	}

	// Show results
	fmt.Printf("\n✓ Hidden %d/%d comment(s) in resolved threads\n", hiddenCount, len(commentIDs))

	if hiddenCount < len(commentIDs) {
		fmt.Fprintf(os.Stderr, "\nWarning: %d comment(s) could not be hidden\n", len(commentIDs)-hiddenCount)
	}

	return nil
}
