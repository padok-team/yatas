package commons

import (
	"fmt"
	"strings"

	"github.com/padok-team/yatas/internal/helpers"
	"github.com/padok-team/yatas/plugins/logger"
	"gopkg.in/yaml.v3"
)

// ParseConfig reads the specified configuration file, parses its content, and returns a Config object.
// It returns an error if the file cannot be read, the content is not valid YAML, or the configuration is not valid.
func ParseConfig(configFile string) (*Config, error) {
	// Read the file .yatas.yml
	data, err := helpers.ReadFile(configFile)
	if err != nil {
		logger.Error("Error reading config file", "config_file", configFile, "error", err)
		return nil, err
	}

	// Parse the YAML file
	var config Config
	if err := unmarshalYAML(data, &config); err != nil {
		logger.Error("Error parsing YAML config", "config_file", configFile, "error", err)
		return nil, err
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		logger.Error("Error validating config", "config_file", configFile, "error", err)
		return nil, err
	}

	return &config, nil
}

// unmarshalYAML unmarshals the given YAML data into the provided object.
func unmarshalYAML(data []byte, obj interface{}) error {
	return yaml.Unmarshal(data, obj)
}

func validateConfig(config *Config) error {
	if len(config.Plugins) == 0 {
		return fmt.Errorf("no plugins defined in config file %s", ".yatas.yml")
	}

	for i, plugin := range config.Plugins {
		if plugin.Name == "" {
			return fmt.Errorf("Plugin at index %d must have a name", i)
		}
		if plugin.Source == "" {
			return fmt.Errorf("Plugin '%s' must have a source", plugin.Name)
		}
		if plugin.Type != "checks" && plugin.Type != "mod" && plugin.Type != "report" && plugin.Type != "" {
			return fmt.Errorf("Plugin '%s' has invalid type '%s'", plugin.Name, plugin.Type)
		}
		if plugin.Version == "" {
			return fmt.Errorf("Plugin '%s' must have a version", plugin.Name)
		}
	}

	return nil
}

func (c *Config) FindPluginWithName(name string) *Plugin {
	for _, plugin := range c.Plugins {
		if plugin.Name == name {
			return &plugin
		}
	}
	return nil
}

func (c *Config) CheckExclude(id string) bool {
	for _, plugins := range c.Plugins {
		for _, exclude := range plugins.Exclude {
			if exclude == id {
				return true
			}
		}
	}
	return false
}

func (c *Config) CheckInclude(id string) bool {
	// Split Id at _ and get first value
	idSplit := strings.Split(id, "_")[0]

	for _, plugins := range c.Plugins {
		if strings.ToUpper(plugins.Name) == idSplit {
			if len(plugins.Include) == 0 {
				return true
			} else {
				for _, include := range plugins.Include {
					if include == id {
						return true
					}
				}
				return false
			}
		}
	}
	return true
}
