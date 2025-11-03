package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var reactCmd = &cobra.Command{
	Use:   "react <comment-id> <emoji>",
	Short: "Add emoji reaction to a comment",
	Long: `Add an emoji reaction to a comment.

Arguments:
  comment-id  Comment ID (PRRC_... or IC_...)
  emoji       Emoji or name (👍, THUMBS_UP, +1, etc.)

Supported reactions:
  👍 THUMBS_UP     😄 LAUGH      ❤️ HEART
  👎 THUMBS_DOWN   🎉 HOORAY     🚀 ROCKET
  😕 CONFUSED      👀 EYES

Examples:
  # Add thumbs up
  gh talk react PRRC_kwDOQN97u86UHqK7 👍

  # Add rocket (by name)
  gh talk react PRRC_kwDOQN97u86UHqK7 ROCKET

  # Add heart (slack-style)
  gh talk react PRRC_kwDOQN97u86UHqK7 :heart:

  # Remove reaction
  gh talk react PRRC_kwDOQN97u86UHqK7 👍 --remove`,
	Args: cobra.ExactArgs(2),
	RunE: runReact,
}

func init() {
	reactCmd.Flags().Bool("remove", false, "Remove reaction instead of adding")
}

func runReact(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	commentID := args[0]
	emojiInput := args[1]

	// Validate comment ID
	if !strings.HasPrefix(commentID, "PRRC_") && !strings.HasPrefix(commentID, "IC_") {
		return fmt.Errorf("invalid comment ID: %s\n\nExpected format: PRRC_... or IC_...", commentID)
	}

	// Parse emoji to GraphQL enum
	content, err := parseEmoji(emojiInput)
	if err != nil {
		return err
	}

	// Create client
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Add or remove reaction
	remove, _ := cmd.Flags().GetBool("remove")
	if remove {
		err = client.RemoveReaction(ctx, commentID, content)
		if err != nil {
			return err
		}
		emoji := contentToEmoji(content)
		fmt.Printf("✓ Removed %s reaction from %s\n", emoji, commentID)
	} else {
		err = client.AddReaction(ctx, commentID, content)
		if err != nil {
			return err
		}
		emoji := contentToEmoji(content)
		fmt.Printf("✓ Added %s reaction to %s\n", emoji, commentID)
	}

	return nil
}

// parseEmoji converts user input to GraphQL ReactionContent enum
func parseEmoji(input string) (string, error) {
	// Normalize
	input = strings.TrimSpace(input)
	lower := strings.ToLower(input)

	// Map various formats to GraphQL enum
	emojiMap := map[string]string{
		// Unicode
		"👍": "THUMBS_UP",
		"👎": "THUMBS_DOWN",
		"😄": "LAUGH",
		"🎉": "HOORAY",
		"😕": "CONFUSED",
		"❤️": "HEART",
		"🚀": "ROCKET",
		"👀": "EYES",
		// Lowercase names
		"thumbs_up":   "THUMBS_UP",
		"thumbs_down": "THUMBS_DOWN",
		"laugh":       "LAUGH",
		"hooray":      "HOORAY",
		"confused":    "CONFUSED",
		"heart":       "HEART",
		"rocket":      "ROCKET",
		"eyes":        "EYES",
		// Slack-style
		":thumbs_up:":   "THUMBS_UP",
		":+1:":          "THUMBS_UP",
		":thumbs_down:": "THUMBS_DOWN",
		":-1:":          "THUMBS_DOWN",
		":laugh:":       "LAUGH",
		":smile:":       "LAUGH",
		":hooray:":      "HOORAY",
		":tada:":        "HOORAY",
		":confused:":    "CONFUSED",
		":heart:":       "HEART",
		":rocket:":      "ROCKET",
		":eyes:":        "EYES",
		// Shorthand
		"+1": "THUMBS_UP",
		"-1": "THUMBS_DOWN",
	}

	if mapped, ok := emojiMap[lower]; ok {
		return mapped, nil
	}

	// Try uppercase (THUMBS_UP → THUMBS_UP)
	upper := strings.ToUpper(input)
	validEnums := []string{"THUMBS_UP", "THUMBS_DOWN", "LAUGH", "HOORAY", "CONFUSED", "HEART", "ROCKET", "EYES"}
	for _, valid := range validEnums {
		if upper == valid {
			return upper, nil
		}
	}

	return "", fmt.Errorf("invalid emoji: %s\n\nSupported reactions:\n  👍 THUMBS_UP     😄 LAUGH      ❤️ HEART\n  👎 THUMBS_DOWN   🎉 HOORAY     🚀 ROCKET\n  😕 CONFUSED      👀 EYES", input)
}


