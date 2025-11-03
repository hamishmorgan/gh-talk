package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	gh "github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

// getRepository gets repository from flag or auto-detects
func getRepository(cmd *cobra.Command) (owner, name string, err error) {
	repoFlag, _ := cmd.Flags().GetString("repo")

	if repoFlag != "" {
		repo, err := repository.Parse(repoFlag)
		if err != nil {
			return "", "", fmt.Errorf("invalid repository format: %s", repoFlag)
		}
		return repo.Owner, repo.Name, nil
	}

	repo, err := repository.Current()
	if err != nil {
		return "", "", fmt.Errorf("could not determine repository\n\nRun this from a git repository or use --repo OWNER/REPO")
	}

	return repo.Owner, repo.Name, nil
}

// getCurrentPR gets PR number from flag or current branch
func getCurrentPR(cmd *cobra.Command) (int, error) {
	prNum, _ := cmd.Flags().GetInt("pr")
	if prNum > 0 {
		return prNum, nil
	}

	// Try to get PR for current branch using gh
	stdout, _, err := gh.Exec("pr", "view", "--json", "number")
	if err != nil {
		return 0, fmt.Errorf("no PR found for current branch\n\nUse --pr NUMBER to specify a PR")
	}

	var result struct {
		Number int `json:"number"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return 0, fmt.Errorf("failed to parse PR number")
	}

	return result.Number, nil
}

// parseThreadID validates and returns thread ID
func parseThreadID(arg string) (string, error) {
	if arg == "" {
		return "", fmt.Errorf("thread ID required")
	}

	if strings.HasPrefix(arg, "PRRT_") {
		return arg, nil
	}

	// Could add URL parsing here in Phase 2
	if strings.Contains(arg, "github.com") || strings.Contains(arg, "http") {
		return "", fmt.Errorf("URL parsing not yet supported\n\nUse the full thread ID (PRRT_...)\nRun 'gh talk list threads' to see thread IDs")
	}

	return "", fmt.Errorf("invalid thread ID format: %s\n\nExpected format: PRRT_kwDOQN97u85gQeTN", arg)
}

// truncate truncates a string to maxLen with "..." if needed
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}


