package commons

import "strings"

type Config struct {
	Plugins []Plugin      `yaml:"plugins"`
	AWS     []AWS_Account `yaml:"aws"`
	Ignore  []Ignore      `yaml:"ignore"`
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
