package commands

import (
	"testing"
)

func TestParseEmoji(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Unicode
		{"unicode thumbs up", "👍", "THUMBS_UP", false},
		{"unicode heart", "❤️", "HEART", false},
		{"unicode rocket", "🚀", "ROCKET", false},

		// Uppercase names
		{"uppercase THUMBS_UP", "THUMBS_UP", "THUMBS_UP", false},
		{"uppercase ROCKET", "ROCKET", "ROCKET", false},

		// Lowercase names
		{"lowercase thumbs_up", "thumbs_up", "THUMBS_UP", false},
		{"lowercase rocket", "rocket", "ROCKET", false},

		// Slack-style
		{"slack :thumbs_up:", ":thumbs_up:", "THUMBS_UP", false},
		{"slack :+1:", ":+1:", "THUMBS_UP", false},
		{"slack :tada:", ":tada:", "HOORAY", false},
		{"slack :heart:", ":heart:", "HEART", false},

		// Shorthand
		{"shorthand +1", "+1", "THUMBS_UP", false},
		{"shorthand -1", "-1", "THUMBS_DOWN", false},

		// Invalid
		{"invalid emoji", "invalid", "", true},
		{"empty string", "", "", true},
		{"random text", "xyz", "", true},

		// All valid emoji types
		{"eyes", "👀", "EYES", false},
		{"laugh", "😄", "LAUGH", false},
		{"hooray", "🎉", "HOORAY", false},
		{"confused", "😕", "CONFUSED", false},
		{"thumbs down", "👎", "THUMBS_DOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseEmoji(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEmoji() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseEmoji(%s) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

