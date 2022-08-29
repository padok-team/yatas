package plugins

import (
	"fmt"
	"regexp"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
	"github.com/stangirard/yatas/plugins/aws"
	"github.com/stangirard/yatas/plugins/custom"
)

// Runs all the plugins that are enabled in the config
func Execute(c *yatas.Config) ([]yatas.Tests, error) {

	plugins := findPlugins(c)

	checks, err := runPlugins(c, plugins)
	if err != nil {
		return nil, err
	}

	return checks, nil
}

// Run the plugins that are enabled in the config with a switch based on the name of the plugin
func runPlugins(c *yatas.Config, plugins []string) ([]yatas.Tests, error) {
	var checksAll []yatas.Tests
	if c.Progress != nil {
		c.AddBar("Plugins : ", "Plugins", len(plugins), 1, c.Progress)
	}

	for _, plugin := range plugins {
		logger.Debug(fmt.Sprint("Running plugin: ", plugin))
		var commandPat = regexp.MustCompile(`custom.*`)
		switch cmd := plugin; {
		case cmd == "aws":
			checks, err := aws.Run(c)
			checksAll = append(checksAll, checks...)
			if err != nil {
				return nil, err
			}
		case commandPat.MatchString(plugin):
			checks, err := custom.Run(c, cmd)
			checksAll = append(checksAll, checks)
			if err != nil {
				return nil, err
			}

		default:
			logger.Error(fmt.Sprint("Plugin not found: ", plugin))
		}
		if c.Progress != nil {
			c.PluginsProgress["Plugins"].Bar.Increment()
		}
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
