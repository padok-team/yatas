package plugins

import (
	"fmt"

	"github.com/stangirard/yatas/internal/aws"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func Execute(c *config.Config) ([]types.Check, error) {

	plugins := findPlugins(c)

	checks, err := runPlugins(c, plugins)
	if err != nil {
		return nil, err
	}

	return checks, nil
}

func runPlugins(c *config.Config, plugins []string) ([]types.Check, error) {
	var checksAll []types.Check
	for _, plugin := range plugins {
		logger.Debug(fmt.Sprint("Running plugin: ", plugin))
		switch plugin {
		case "aws":
			checks, err := aws.Run(c)
			checksAll = append(checksAll, checks...)
			if err != nil {
				return nil, err
			}
		default:
			logger.Error(fmt.Sprint("Plugin not found: ", plugin))
		}
	}
	return checksAll, nil
}

func findPlugins(c *config.Config) []string {
	var plugins []string
	for _, plugin := range c.Plugins {
		if plugin.Enabled {
			plugins = append(plugins, plugin.Name)
		}
	}
	logger.Debug(fmt.Sprint("Plugins Found in config: ", plugins))

	return plugins
}
