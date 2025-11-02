package main

import (
	"testing"

	_ "github.com/hamishmorgan/gh-talk/internal/api"
	_ "github.com/hamishmorgan/gh-talk/internal/cache"
	_ "github.com/hamishmorgan/gh-talk/internal/commands"
	_ "github.com/hamishmorgan/gh-talk/internal/config"
	_ "github.com/hamishmorgan/gh-talk/internal/filter"
	_ "github.com/hamishmorgan/gh-talk/internal/format"
	_ "github.com/hamishmorgan/gh-talk/internal/tui"
)

func TestDirectoryStructure(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
	}{
		{"API package", "github.com/hamishmorgan/gh-talk/internal/api"},
		{"Commands package", "github.com/hamishmorgan/gh-talk/internal/commands"},
		{"Filter package", "github.com/hamishmorgan/gh-talk/internal/filter"},
		{"Format package", "github.com/hamishmorgan/gh-talk/internal/format"},
		{"Config package", "github.com/hamishmorgan/gh-talk/internal/config"},
		{"Cache package", "github.com/hamishmorgan/gh-talk/internal/cache"},
		{"TUI package", "github.com/hamishmorgan/gh-talk/internal/tui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("✓ %s imported successfully from %s", tt.name, tt.pkgPath)
		})
	}
}

