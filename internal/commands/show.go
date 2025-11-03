package commands

import (
	"context"
	"fmt"

	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <thread-id>",
	Short: "Show thread details",
	Long: `Show detailed information about a review thread.

Arguments:
  thread-id   Thread ID (PRRT_...)

Examples:
  # Show thread details
  gh talk show PRRT_kwDOQN97u85gQeTN`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func runShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	threadID, err := parseThreadID(args[0])
	if err != nil {
		return err
	}

	// Get repository context (needed to fetch thread details)
	owner, name, err := getRepository(cmd)
	if err != nil {
		return err
	}

	prNum, err := getCurrentPR(cmd)
	if err != nil {
		return err
	}

	// Create client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Fetch all threads to find the one we want
	threads, err := client.ListThreads(ctx, owner, name, prNum)
	if err != nil {
		return err
	}

	// Find thread
	var thread *api.Thread
	for _, t := range threads {
		if t.ID == threadID {
			thread = &t
			break
		}
	}

	if thread == nil {
		return fmt.Errorf("thread not found: %s\n\nRun 'gh talk list threads' to see available threads", threadID)
	}

	// Display thread details
	fmt.Printf("Thread: %s\n", thread.ID)
	fmt.Printf("File:   %s:%d\n", thread.Path, thread.Line)
	fmt.Printf("Status: ")
	if thread.IsResolved {
		fmt.Printf("✓ RESOLVED")
		if thread.ResolvedBy != nil {
			fmt.Printf(" by @%s", thread.ResolvedBy.Login)
		}
		fmt.Println()
	} else {
		fmt.Println("○ OPEN")
	}

	if thread.IsOutdated {
		fmt.Println("⚠️  Outdated (code has changed since comment)")
	}

	fmt.Printf("\nConversation (%d comments):\n\n", len(thread.Comments))

	for i, comment := range thread.Comments {
		fmt.Printf("─────────────────────────────────────────────────────\n")
		fmt.Printf("@%s", comment.Author.Login)
		if comment.ReplyTo != nil {
			fmt.Printf(" (in reply)")
		}
		fmt.Printf("\n\n")
		fmt.Println(comment.Body)

		// Show reactions
		if len(comment.ReactionGroups) > 0 {
			fmt.Printf("\nReactions: ")
			for j, rg := range comment.ReactionGroups {
				if j > 0 {
					fmt.Printf(" ")
				}
				emoji := contentToEmoji(rg.Content)
				fmt.Printf("%s %d", emoji, rg.Users.TotalCount)
			}
			fmt.Println()
		}

		if i < len(thread.Comments)-1 {
			fmt.Println()
		}
	}

	fmt.Printf("─────────────────────────────────────────────────────\n")

	return nil
}


