package commands

import (
	"context"
	"fmt"

	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show PR review status summary",
	Long: `Show an overview of review thread status for a pull request.

Displays:
  - Total threads and resolution status
  - Comment counts
  - Recent activity summary
  - Overall completion status

Examples:
  # Show status for current PR
  gh talk status

  # Show status for specific PR
  gh talk status --pr 137

  # Compact output
  gh talk status --compact`,
	RunE: runStatus,
}

func init() {
	statusCmd.Flags().Bool("compact", false, "Show compact one-line summary")
}

func runStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get repository context
	owner, name, err := getRepository(cmd)
	if err != nil {
		return err
	}

	// Get PR number
	prNum, err := getCurrentPR(cmd)
	if err != nil {
		return err
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Fetch threads
	threads, err := client.ListThreads(ctx, owner, name, prNum)
	if err != nil {
		return err
	}

	// Calculate statistics
	totalThreads := len(threads)
	resolvedThreads := 0
	unresolvedThreads := 0
	totalComments := 0
	totalReactions := 0

	for _, t := range threads {
		totalComments += len(t.Comments)
		if t.IsResolved {
			resolvedThreads++
		} else {
			unresolvedThreads++
		}

		// Count reactions
		for _, c := range t.Comments {
			for _, rg := range c.ReactionGroups {
				totalReactions += rg.Users.TotalCount
			}
		}
	}

	// Display
	compact, _ := cmd.Flags().GetBool("compact")

	if compact {
		// One-line summary
		fmt.Printf("PR %s/%s#%d: %d threads (%d resolved, %d unresolved), %d comments, %d reactions\n",
			owner, name, prNum, totalThreads, resolvedThreads, unresolvedThreads, totalComments, totalReactions)
		return nil
	}

	// Detailed output
	fmt.Printf("PR: %s/%s#%d\n\n", owner, name, prNum)

	fmt.Printf("Threads:\n")
	fmt.Printf("  Total:      %d\n", totalThreads)
	fmt.Printf("  Resolved:   %d", resolvedThreads)
	if resolvedThreads == totalThreads && totalThreads > 0 {
		fmt.Printf(" ✓ All resolved!\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("  Unresolved: %d", unresolvedThreads)
	if unresolvedThreads > 0 {
		fmt.Printf(" ⚠️  Needs attention\n")
	} else if totalThreads > 0 {
		fmt.Printf(" ✓\n")
	} else {
		fmt.Printf("\n")
	}

	fmt.Printf("\nComments:\n")
	fmt.Printf("  Total:      %d\n", totalComments)

	fmt.Printf("\nReactions:\n")
	fmt.Printf("  Total:      %d\n", totalReactions)

	// Overall status
	fmt.Printf("\nOverall Status: ")
	if unresolvedThreads == 0 && totalThreads > 0 {
		fmt.Printf("✓ All feedback addressed\n")
	} else if unresolvedThreads > 0 {
		fmt.Printf("⚠️  %d thread(s) need attention\n", unresolvedThreads)
	} else {
		fmt.Printf("No review threads found\n")
	}

	return nil
}

