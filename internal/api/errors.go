package api

import (
	"errors"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
)

// handleError converts go-gh errors into user-friendly errors
func handleError(err error) error {
	if err == nil {
		return nil
	}

	var gqlErr *api.GraphQLError
	if errors.As(err, &gqlErr) {
		return handleGraphQLError(gqlErr)
	}

	var httpErr *api.HTTPError
	if errors.As(err, &httpErr) {
		return handleHTTPError(httpErr)
	}

	return err
}

func handleGraphQLError(gqlErr *api.GraphQLError) error {
	if len(gqlErr.Errors) == 0 {
		return fmt.Errorf("GraphQL error")
	}

	e := gqlErr.Errors[0]

	switch e.Type {
	case "NOT_FOUND":
		return fmt.Errorf("resource not found: %s\n\nThe thread, comment, or resource may have been deleted", e.Message)
	case "FORBIDDEN":
		return fmt.Errorf("permission denied: %s\n\nYou may not have access to this repository or resource", e.Message)
	case "UNPROCESSABLE":
		return fmt.Errorf("invalid input: %s", e.Message)
	case "RATE_LIMIT":
		return fmt.Errorf("rate limit exceeded\n\nGitHub API rate limit reached. Try again later")
	default:
		return fmt.Errorf("GraphQL error: %s", e.Message)
	}
}

func handleHTTPError(httpErr *api.HTTPError) error {
	switch httpErr.StatusCode {
	case 401:
		return fmt.Errorf("authentication failed\n\nRun 'gh auth login' to authenticate")
	case 403:
		return fmt.Errorf("forbidden: %s\n\nCheck repository permissions", httpErr.Message)
	case 404:
		return fmt.Errorf("not found: %s", httpErr.Message)
	case 422:
		return fmt.Errorf("validation failed: %s", httpErr.Message)
	case 502, 503, 504:
		return fmt.Errorf("GitHub API unavailable\n\nTry again in a few moments")
	default:
		return fmt.Errorf("HTTP %d: %s", httpErr.StatusCode, httpErr.Message)
	}
}


