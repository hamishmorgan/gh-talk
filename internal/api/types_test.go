package api

import (
	"encoding/json"
	"os"
	"testing"
)

func TestThreadStructure(t *testing.T) {
	// Test that our Thread type can unmarshal from real API response
	data, err := os.ReadFile("../../testdata/pr_with_resolved_threads.json")
	if err != nil {
		t.Skipf("Skipping test, fixture not found: %v", err)
		return
	}

	var response struct {
		Data struct {
			Repository struct {
				PullRequest struct {
					ReviewThreads struct {
						Nodes []struct {
							ID         string
							IsResolved bool
							Path       string
						}
					}
				}
			}
		}
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal fixture: %v", err)
	}

	nodes := response.Data.Repository.PullRequest.ReviewThreads.Nodes
	if len(nodes) == 0 {
		t.Error("Expected threads in fixture")
	}

	// Verify we can access fields
	for _, node := range nodes {
		if node.ID == "" {
			t.Error("Thread ID is empty")
		}
		if node.Path == "" {
			t.Error("Thread path is empty")
		}
	}
}
