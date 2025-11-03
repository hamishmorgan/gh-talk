package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve [thread-id...]",
	Short: "Resolve review threads",
	Long: `Mark one or more review threads as resolved.

Arguments:
  thread-id   One or more thread IDs, or omit for interactive selection

Examples:
  # Interactive selection
  gh talk resolve

  # Single thread
  gh talk resolve PRRT_kwDOQN97u85gQeTN

  # Multiple threads
  gh talk resolve PRRT_abc123 PRRT_def456 PRRT_ghi789

  # With message first
  gh talk resolve PRRT_abc123 --message "Fixed in commit abc123"`,
	Args: cobra.MinimumNArgs(0),
	RunE: runResolve,
}

var unresolveCmd = &cobra.Command{
	Use:   "unresolve [thread-id...]",
	Short: "Unresolve review threads",
	Long: `Mark one or more review threads as unresolved (reopen discussion).

Arguments:
  thread-id   One or more thread IDs, or omit for interactive selection

Examples:
  # Interactive selection
  gh talk unresolve

  # Single thread
  gh talk unresolve PRRT_kwDOQN97u85gQeTN

  # Multiple threads
  gh talk unresolve PRRT_abc123 PRRT_def456`,
	Args: cobra.MinimumNArgs(0),
	RunE: runUnresolve,
}

func init() {
	resolveCmd.Flags().StringP("message", "m", "", "Message to post before resolving")
	resolveCmd.Flags().BoolP("yes", "y", false, "Skip confirmation for multiple threads")

	unresolveCmd.Flags().BoolP("yes", "y", false, "Skip confirmation for multiple threads")
}

func runResolve(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var threadIDs []string

	if len(args) == 0 {
		// Interactive selection
		owner, name, err := getRepository(cmd)
		if err != nil {
			return err
		}
		prNum, err := getCurrentPR(cmd)
		if err != nil {
			return err
		}

		ids, err := selectThreadsInteractive(ctx, owner, name, prNum, false)
		if err != nil {
			return err
		}
		threadIDs = ids
	} else {
		// Validate all IDs
		for _, arg := range args {
			id, err := parseThreadID(arg)
			if err != nil {
				return err
			}
			threadIDs = append(threadIDs, id)
		}
	}

	if len(threadIDs) == 0 {
		return fmt.Errorf("no threads selected")
	}

	// Confirm for multiple threads
	if len(threadIDs) > 1 {
		skipConfirm, _ := cmd.Flags().GetBool("yes")
		if !skipConfirm {
			p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
			confirmed, err := p.Confirm(fmt.Sprintf("Resolve %d threads?", len(threadIDs)), false)
			if err != nil || !confirmed {
				return fmt.Errorf("cancelled")
			}
		}
	}

	// Create client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Post message first if provided
	message, _ := cmd.Flags().GetString("message")
	if message != "" {
		for _, id := range threadIDs {
			if err := client.ReplyToThread(ctx, id, message); err != nil {
				return fmt.Errorf("failed to add message to %s: %w", id, err)
			}
		}
	}

	// Resolve each thread
	for _, id := range threadIDs {
		if err := client.ResolveThread(ctx, id); err != nil {
			return fmt.Errorf("failed to resolve %s: %w", id, err)
		}
		fmt.Printf("✓ Resolved %s\n", id)
	}

	return nil
}

func runUnresolve(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var threadIDs []string

	if len(args) == 0 {
		// Interactive selection from resolved threads
		owner, name, err := getRepository(cmd)
		if err != nil {
			return err
		}
		prNum, err := getCurrentPR(cmd)
		if err != nil {
			return err
		}

		ids, err := selectThreadsInteractive(ctx, owner, name, prNum, true)
		if err != nil {
			return err
		}
		threadIDs = ids
	} else {
		for _, arg := range args {
			id, err := parseThreadID(arg)
			if err != nil {
				return err
			}
			threadIDs = append(threadIDs, id)
		}
	}

	if len(threadIDs) == 0 {
		return fmt.Errorf("no threads selected")
	}

	// Confirm for multiple
	if len(threadIDs) > 1 {
		skipConfirm, _ := cmd.Flags().GetBool("yes")
		if !skipConfirm {
			p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
			confirmed, err := p.Confirm(fmt.Sprintf("Unresolve %d threads?", len(threadIDs)), false)
			if err != nil || !confirmed {
				return fmt.Errorf("cancelled")
			}
		}
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	for _, id := range threadIDs {
		if err := client.UnresolveThread(ctx, id); err != nil {
			return fmt.Errorf("failed to unresolve %s: %w", id, err)
		}
		fmt.Printf("✓ Unresolved %s\n", id)
	}

	return nil
}

func selectThreadsInteractive(ctx context.Context, owner, name string, pr int, onlyResolved bool) ([]string, error) {
	client, err := api.NewClient()
	if err != nil {
		return nil, err
	}

	threads, err := client.ListThreads(ctx, owner, name, pr)
	if err != nil {
		return nil, err
	}

	// Filter by resolution status
	filtered := make([]api.Thread, 0)
	for _, t := range threads {
		if onlyResolved && t.IsResolved {
			filtered = append(filtered, t)
		} else if !onlyResolved && !t.IsResolved {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) == 0 {
		if onlyResolved {
			return nil, fmt.Errorf("no resolved threads found")
		}
		return nil, fmt.Errorf("no unresolved threads found")
	}

	// Build options
	options := make([]string, len(filtered))
	for i, t := range filtered {
		preview := ""
		if len(t.Comments) > 0 {
			preview = truncate(t.Comments[0].Body, 50)
		}
		status := "○"
		if t.IsResolved {
			status = "✓"
		}
		options[i] = fmt.Sprintf("%s %s:%d - %s", status, t.Path, t.Line, preview)
	}

	// Multi-select
	p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
	indices, err := p.MultiSelect("Select threads:", nil, options)
	if err != nil {
		return nil, err
	}

	if len(indices) == 0 {
		return nil, fmt.Errorf("no threads selected")
	}

	// Return selected IDs
	ids := make([]string, len(indices))
	for i, idx := range indices {
		ids[i] = filtered[idx].ID
	}

	return ids, nil
}
