package main

import (
	"testing"

	_ "github.com/hamishmorgan/gh-talk/internal/cache"
	_ "github.com/hamishmorgan/gh-talk/internal/tui"
	_ "github.com/hamishmorgan/gh-talk/pkg/api"
	_ "github.com/hamishmorgan/gh-talk/pkg/config"
	_ "github.com/hamishmorgan/gh-talk/pkg/filter"
	_ "github.com/hamishmorgan/gh-talk/pkg/format"
)

func TestDirectoryStructure(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
	}{
		{"API package", "github.com/hamishmorgan/gh-talk/pkg/api"},
		{"Filter package", "github.com/hamishmorgan/gh-talk/pkg/filter"},
		{"Format package", "github.com/hamishmorgan/gh-talk/pkg/format"},
		{"Config package", "github.com/hamishmorgan/gh-talk/pkg/config"},
		{"Cache package", "github.com/hamishmorgan/gh-talk/internal/cache"},
		{"TUI package", "github.com/hamishmorgan/gh-talk/internal/tui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("✓ %s imported successfully from %s", tt.name, tt.pkgPath)
		})
	}
}

