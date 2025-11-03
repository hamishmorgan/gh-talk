package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/hamishmorgan/gh-talk/internal/api"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List conversations",
}

var listThreadsCmd = &cobra.Command{
	Use:   "threads",
	Short: "List review threads",
	Long: `List review threads from a pull request.

By default, shows only unresolved threads. Use --all to see
both resolved and unresolved threads.

Examples:
  # List unresolved threads in current PR
  gh talk list threads

  # List all threads in PR #123
  gh talk list threads --pr 123 --all

  # List threads on specific file
  gh talk list threads --file src/api.go`,
	RunE: runListThreads,
}

func init() {
	listCmd.AddCommand(listThreadsCmd)

	// Filter flags
	listThreadsCmd.Flags().Bool("unresolved", false, "Show only unresolved threads")
	listThreadsCmd.Flags().Bool("resolved", false, "Show only resolved threads")
	listThreadsCmd.Flags().Bool("all", false, "Show all threads")
	listThreadsCmd.Flags().String("author", "", "Filter by author")
	listThreadsCmd.Flags().String("file", "", "Filter by file path")

	// Output flags
	listThreadsCmd.Flags().String("format", "", "Output format (table, json, tsv)")

	// Make resolution flags mutually exclusive
	listThreadsCmd.MarkFlagsMutuallyExclusive("unresolved", "resolved", "all")
}

func runListThreads(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Fetch threads
	threads, err := client.ListThreads(ctx, owner, name, prNum)
	if err != nil {
		return err
	}

	// Apply filters
	threads = filterThreads(cmd, threads)

	if len(threads) == 0 {
		fmt.Printf("No threads found in %s/%s#%d\n", owner, name, prNum)
		return nil
	}

	// Output
	return outputThreads(cmd, threads)
}

func filterThreads(cmd *cobra.Command, threads []api.Thread) []api.Thread {
	unresolved, _ := cmd.Flags().GetBool("unresolved")
	resolved, _ := cmd.Flags().GetBool("resolved")
	all, _ := cmd.Flags().GetBool("all")
	author, _ := cmd.Flags().GetString("author")
	file, _ := cmd.Flags().GetString("file")

	// Default to unresolved if no filter specified
	if !unresolved && !resolved && !all {
		unresolved = true
	}

	filtered := make([]api.Thread, 0, len(threads))
	for _, t := range threads {
		// Resolution filter
		if unresolved && t.IsResolved {
			continue
		}
		if resolved && !t.IsResolved {
			continue
		}

		// Author filter
		if author != "" {
			hasAuthor := false
			for _, c := range t.Comments {
				if c.Author.Login == author {
					hasAuthor = true
					break
				}
			}
			if !hasAuthor {
				continue
			}
		}

		// File filter
		if file != "" && t.Path != file {
			continue
		}

		filtered = append(filtered, t)
	}

	return filtered
}

func outputThreads(cmd *cobra.Command, threads []api.Thread) error {
	format, _ := cmd.Flags().GetString("format")
	terminal := term.FromEnv()

	// Auto-detect format if not specified
	if format == "" {
		if terminal.IsTerminalOutput() {
			format = "table"
		} else {
			format = "tsv"
		}
	}

	switch format {
	case "table":
		return outputThreadsTable(threads, terminal)
	case "tsv":
		return outputThreadsTSV(threads, terminal)
	case "json":
		return outputThreadsJSON(threads, terminal)
	default:
		return fmt.Errorf("unknown format: %s\n\nValid formats: table, json, tsv", format)
	}
}

func outputThreadsTable(threads []api.Thread, terminal term.Term) error {
	width, _, _ := terminal.Size()
	t := tableprinter.New(terminal.Out(), true, width)

	// Header
	t.AddField("ID")
	t.AddField("File:Line")
	t.AddField("Status")
	t.AddField("Comments")
	t.AddField("Reactions")
	t.AddField("Preview")
	t.EndRow()

	// Rows
	for _, thread := range threads {
		// ID
		t.AddField(thread.ID)

		// File:Line
		fileLine := fmt.Sprintf("%s:%d", thread.Path, thread.Line)
		t.AddField(fileLine)

		// Status
		status := "○ OPEN"
		if thread.IsResolved {
			status = "✓ RESOLVED"
			if thread.ResolvedBy != nil {
				status += " by @" + thread.ResolvedBy.Login
			}
		}
		t.AddField(status)

		// Comment count
		t.AddField(fmt.Sprintf("%d", len(thread.Comments)))

		// Reactions (show non-zero only)
		reactions := formatReactions(thread)
		t.AddField(reactions)

		// Preview (first comment body, truncated)
		preview := ""
		if len(thread.Comments) > 0 {
			preview = truncate(thread.Comments[0].Body, 50)
		}
		t.AddField(preview)

		t.EndRow()
	}

	return t.Render()
}

func outputThreadsTSV(threads []api.Thread, terminal term.Term) error {
	t := tableprinter.New(terminal.Out(), false, 0)

	// Header
	t.AddField("ID")
	t.AddField("Path")
	t.AddField("Line")
	t.AddField("IsResolved")
	t.AddField("CommentCount")
	t.AddField("Preview")
	t.EndRow()

	// Rows
	for _, thread := range threads {
		t.AddField(thread.ID)
		t.AddField(thread.Path)
		t.AddField(fmt.Sprintf("%d", thread.Line))
		t.AddField(fmt.Sprintf("%t", thread.IsResolved))
		t.AddField(fmt.Sprintf("%d", len(thread.Comments)))

		preview := ""
		if len(thread.Comments) > 0 {
			preview = thread.Comments[0].Body
		}
		t.AddField(preview)
		t.EndRow()
	}

	return t.Render()
}

func outputThreadsJSON(threads []api.Thread, terminal term.Term) error {
	// For now, simple JSON output
	// TODO: Support --json <fields> like gh CLI
	fmt.Fprintf(terminal.Out(), "[\n")
	for i, thread := range threads {
		fmt.Fprintf(terminal.Out(), `  {"id":"%s","path":"%s","line":%d,"isResolved":%t,"comments":%d}`,
			thread.ID, thread.Path, thread.Line, thread.IsResolved, len(thread.Comments))
		if i < len(threads)-1 {
			fmt.Fprintf(terminal.Out(), ",")
		}
		fmt.Fprintf(terminal.Out(), "\n")
	}
	fmt.Fprintf(terminal.Out(), "]\n")
	return nil
}

func formatReactions(thread api.Thread) string {
	if len(thread.Comments) == 0 {
		return ""
	}

	// Collect all unique reactions from all comments
	reactions := make(map[string]int)
	for _, comment := range thread.Comments {
		for _, rg := range comment.ReactionGroups {
			emoji := contentToEmoji(rg.Content)
			reactions[emoji] += rg.Users.TotalCount
		}
	}

	if len(reactions) == 0 {
		return ""
	}

	// Format as emoji count pairs
	parts := make([]string, 0, len(reactions))
	for emoji, count := range reactions {
		parts = append(parts, fmt.Sprintf("%s %d", emoji, count))
	}

	return strings.Join(parts, " ")
}

func contentToEmoji(content string) string {
	switch content {
	case "THUMBS_UP":
		return "👍"
	case "THUMBS_DOWN":
		return "👎"
	case "LAUGH":
		return "😄"
	case "HOORAY":
		return "🎉"
	case "CONFUSED":
		return "😕"
	case "HEART":
		return "❤️"
	case "ROCKET":
		return "🚀"
	case "EYES":
		return "👀"
	default:
		return content
	}
}
