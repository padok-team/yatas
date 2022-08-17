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
