package plugins

import (
	"testing"

	"github.com/stangirard/yatas/internal/yatas"
)

var config = yatas.Config{
	Plugins: []yatas.Plugin{
		{
			Name:        "aws",
			Enabled:     true,
			Description: "AWS Plugin",
			Exclude:     []string{},
			Include:     []string{},
		},
	},
}

func TestFindPlugins(t *testing.T) {
	plugins := findPlugins(&config)
	if len(plugins) != 1 {
		t.Error("Expected 1 plugin, got", len(plugins))
	}
}

// Test the RunPlugins function
func TestRunPlugins(t *testing.T) {
	plugins := findPlugins(&config)
	checks, err := runPlugins(&config, plugins)
	if err != nil {
		t.Error(err)
	}
	if len(checks) != 0 {
		t.Error("Expected 0 check, got", len(checks))
	}
}

// Test the Execute function
func TestExecute(t *testing.T) {
	checks, err := Execute(&config)
	if err != nil {
		t.Error(err)
	}
	if len(checks) != 0 {
		t.Error("Expected 0 check, got", len(checks))
	}
}
