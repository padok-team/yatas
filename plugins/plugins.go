package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/stangirard/yatas/example"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
	"github.com/stangirard/yatas/plugins/aws"
)

// Here is a real implementation of Greeter
type GreeterHello struct {
	logger hclog.Logger
}

func (g *GreeterHello) Run(c *yatas.Config) []yatas.Tests {
	g.logger.Debug("message from GreeterHello.Run")

	var checksAll []yatas.Tests

	checks, err := runPlugins(c, "aws")
	if err != nil {
		g.logger.Error("Error running plugins", "error", err)
	}
	checksAll = append(checksAll, checks...)
	return checksAll
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	greeter := &GreeterHello{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"aws": &example.GreeterPlugin{Impl: greeter},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

// Run the plugins that are enabled in the config with a switch based on the name of the plugin
func runPlugins(c *yatas.Config, plugin string) ([]yatas.Tests, error) {
	var checksAll []yatas.Tests

	logger.Debug(fmt.Sprint("Running plugin: ", plugin))

	checks, err := aws.Run(c)
	checksAll = append(checksAll, checks...)
	if err != nil {
		return nil, err
	}

	return checksAll, nil
}

// Returns a list of plugins that are enabled in the config
func findPlugins(c *yatas.Config) []string {
	var plugins []string
	for _, plugin := range c.Plugins {
		if plugin.Enabled {
			plugins = append(plugins, plugin.Name)
		}
	}
	logger.Debug(fmt.Sprint("Plugins Found in config: ", plugins))

	return plugins
}
