package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-talk",
	Short: "Manage GitHub PR and Issue conversations",
	Long: `gh-talk is a GitHub CLI extension for managing conversations
on Pull Requests and Issues from the terminal.

Features:
  • Reply to review threads
  • Add emoji reactions
  • Resolve/unresolve threads
  • Filter and list conversations
  • Bulk operations

Never leave the terminal for code review conversations.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global persistent flags
	rootCmd.PersistentFlags().StringP("repo", "R", "", "Repository (OWNER/REPO)")
	rootCmd.PersistentFlags().Int("pr", 0, "PR number")
	rootCmd.PersistentFlags().Int("issue", 0, "Issue number")

	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(replyCmd)
	rootCmd.AddCommand(resolveCmd)
	rootCmd.AddCommand(unresolveCmd)
	rootCmd.AddCommand(reactCmd)
	rootCmd.AddCommand(hideCmd)
	rootCmd.AddCommand(unhideCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(statusCmd)
}
