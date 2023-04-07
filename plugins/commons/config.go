package commons

import (
	"strings"

	"github.com/padok-team/yatas/internal/helpers"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Plugins      []Plugin                 `yaml:"plugins"`
	Ignore       []Ignore                 `yaml:"ignore"`
	PluginConfig []map[string]interface{} `yaml:"pluginsConfiguration"`
	Tests        []Tests                  `yaml:"tests"`
}

func ParseConfig(configFile string) (*Config, error) {
	// Read the file .yatas.yml
	// File to array of bytes
	data, err := helpers.ReadFile(configFile)
	if err != nil {
		return &Config{}, err
	}

	// Parse the yaml file
	config := Config{}
	err = unmarshalYAML(data, &config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}

func unmarshalYAML(data []byte, config *Config) error {
	err := yaml.Unmarshal([]byte(data), &config)

	return err
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
