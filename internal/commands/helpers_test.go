package commands

import (
	"testing"
)

func TestIsNumericID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"pure numeric", "2486231843", true},
		{"zero", "0", true},
		{"large number", "999999999999", true},
		{"node ID", "PRRT_kwDOQN97u85gQeTN", false},
		{"comment node ID", "PRRC_kwDOQN97u86UHqK7", false},
		{"issue comment ID", "IC_kwDOQN97u87PVA8l", false},
		{"alphanumeric", "abc123", false},
		{"with spaces", "123 456", false},
		{"empty", "", false},
		{"text", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNumericID(tt.input); got != tt.want {
				t.Errorf("isNumericID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseThreadID(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		wantID    string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid thread ID",
			arg:       "PRRT_kwDOQN97u85gQeTN",
			wantID:    "PRRT_kwDOQN97u85gQeTN",
			wantError: false,
		},
		{
			name:      "empty string",
			arg:       "",
			wantID:    "",
			wantError: true,
			errorMsg:  "thread ID required",
		},
		{
			name:      "numeric database ID",
			arg:       "2486231843",
			wantID:    "",
			wantError: true,
			errorMsg:  "numeric database ID",
		},
		{
			name:      "URL format",
			arg:       "https://github.com/owner/repo/pull/123#discussion_r12345",
			wantID:    "",
			wantError: true,
			errorMsg:  "URL parsing not yet supported",
		},
		{
			name:      "invalid format",
			arg:       "invalid_id",
			wantID:    "",
			wantError: true,
			errorMsg:  "invalid thread ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := parseThreadID(tt.arg)
			if (err != nil) != tt.wantError {
				t.Errorf("parseThreadID(%q) error = %v, wantError %v", tt.arg, err, tt.wantError)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("parseThreadID(%q) = %v, want %v", tt.arg, gotID, tt.wantID)
			}
			if err != nil && tt.errorMsg != "" {
				if !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("parseThreadID(%q) error message = %q, should contain %q", tt.arg, err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
