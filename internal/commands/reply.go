package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var replyCmd = &cobra.Command{
	Use:   "reply [thread-id] [message]",
	Short: "Reply to a review thread",
	Long: `Reply to a review thread with a message.

Arguments:
  thread-id   Thread ID (PRRT_...), or omit for interactive selection
  message     Reply message text, or use --editor

Examples:
  # Interactive mode (prompts for thread and message)
  gh talk reply

  # With thread ID and message
  gh talk reply PRRT_kwDOQN97u85gQeTN "Fixed in latest commit"

  # Reply and resolve
  gh talk reply PRRT_kwDOQN97u85gQeTN "Done!" --resolve

  # Using editor
  gh talk reply PRRT_kwDOQN97u85gQeTN --editor`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runReply,
}

func init() {
	replyCmd.Flags().BoolP("editor", "e", false, "Open editor for message composition")
	replyCmd.Flags().Bool("resolve", false, "Resolve thread after replying")
	replyCmd.Flags().StringP("message", "m", "", "Message text (alternative to positional argument)")

	replyCmd.MarkFlagsMutuallyExclusive("editor", "message")
}

func runReply(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse thread ID
	var threadID string
	var message string

	switch len(args) {
	case 0:
		// Interactive mode for thread selection
		owner, name, err := getRepository(cmd)
		if err != nil {
			return err
		}
		prNum, err := getCurrentPR(cmd)
		if err != nil {
			return err
		}

		id, err := selectThreadInteractive(ctx, owner, name, prNum)
		if err != nil {
			return err
		}
		threadID = id

	case 1:
		// Thread ID provided, message from flag or interactive
		id, err := parseThreadID(args[0])
		if err != nil {
			return err
		}
		threadID = id

	case 2:
		// Both provided
		id, err := parseThreadID(args[0])
		if err != nil {
			return err
		}
		threadID = id
		message = args[1]
	}

	// Get message if not already set
	if message == "" {
		useEditor, _ := cmd.Flags().GetBool("editor")
		msgFlag, _ := cmd.Flags().GetString("message")

		if msgFlag != "" {
			message = msgFlag
		} else if useEditor {
			msg, err := openEditor()
			if err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}
			message = msg
		} else {
			// Prompt for message
			p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
			msg, err := p.Input("Reply message:", "")
			if err != nil {
				return fmt.Errorf("failed to get message: %w", err)
			}
			message = msg
		}
	}

	// Validate message
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Post reply
	err = client.ReplyToThread(ctx, threadID, message)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Replied to thread %s\n", threadID)

	// Resolve if requested
	shouldResolve, _ := cmd.Flags().GetBool("resolve")
	if shouldResolve {
		err = client.ResolveThread(ctx, threadID)
		if err != nil {
			return fmt.Errorf("replied successfully but failed to resolve: %w", err)
		}
		fmt.Printf("✓ Resolved thread\n")
	}

	return nil
}

func selectThreadInteractive(ctx context.Context, owner, name string, pr int) (string, error) {
	// Create client
	client, err := api.NewClient()
	if err != nil {
		return "", err
	}

	// Fetch threads
	threads, err := client.ListThreads(ctx, owner, name, pr)
	if err != nil {
		return "", err
	}

	// Filter to unresolved
	unresolved := make([]api.Thread, 0)
	for _, t := range threads {
		if !t.IsResolved {
			unresolved = append(unresolved, t)
		}
	}

	if len(unresolved) == 0 {
		return "", fmt.Errorf("no unresolved threads found")
	}

	// Build options
	options := make([]string, len(unresolved))
	for i, t := range unresolved {
		preview := ""
		if len(t.Comments) > 0 {
			preview = truncate(t.Comments[0].Body, 50)
		}
		options[i] = fmt.Sprintf("%s:%d - %s", t.Path, t.Line, preview)
	}

	// Prompt
	p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
	idx, err := p.Select("Select thread:", "", options)
	if err != nil {
		return "", err
	}

	return unresolved[idx].ID, nil
}

func openEditor() (string, error) {
	editor := os.Getenv("GH_TALK_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vim"
	}

	// TODO: Implement actual editor integration
	// For now, return error
	return "", fmt.Errorf("editor integration not yet implemented\n\nUse --message flag or provide message as argument")
}
