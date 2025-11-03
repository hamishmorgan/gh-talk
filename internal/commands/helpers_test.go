package commands

import (
	"testing"
)

func TestParseThreadID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid thread ID",
			input:   "PRRT_kwDOQN97u85gQeTN",
			want:    "PRRT_kwDOQN97u85gQeTN",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			input:   "INVALID_123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "URL not supported yet",
			input:   "https://github.com/owner/repo/pull/123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "short ID not supported",
			input:   "1",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseThreadID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseThreadID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseThreadID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "shorter than max",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "exactly max",
			input:  "hello",
			maxLen: 5,
			want:   "hello",
		},
		{
			name:   "longer than max",
			input:  "hello world this is a long string",
			maxLen: 20,
			want:   "hello world this ...",
		},
		{
			name:   "very short max",
			input:  "hello",
			maxLen: 3,
			want:   "hel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentToEmoji(t *testing.T) {
	tests := []struct {
		content string
		want    string
	}{
		{"THUMBS_UP", "👍"},
		{"THUMBS_DOWN", "👎"},
		{"LAUGH", "😄"},
		{"HOORAY", "🎉"},
		{"CONFUSED", "😕"},
		{"HEART", "❤️"},
		{"ROCKET", "🚀"},
		{"EYES", "👀"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			got := contentToEmoji(tt.content)
			if got != tt.want {
				t.Errorf("contentToEmoji(%s) = %v, want %v", tt.content, got, tt.want)
			}
		})
	}
}
