package config

import (
	"github.com/stangirard/yatas/internal/helpers"
	"gopkg.in/yaml.v3"
)

type Plugin struct {
	Name        string   `yaml:"name"`
	Enabled     bool     `yaml:"enabled"`
	Description string   `yaml:"description"`
	Exclude     []string `yaml:"exclude"`
}

type Config struct {
	Plugins []Plugin `yaml:"plugins"`
	AWS     struct {
		Enabled bool `yaml:"enabled"`
		Account struct {
			Profile string `yaml:"profile"`
			SSO     bool   `yaml:"sso"`
			Region  string `yaml:"region"`
		} `yaml:"account"`
	} `yaml:"aws"`
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

func ParseConfig(configFile string) (Config, error) {
	// Read the file .yatas.yml
	// File to array of bytes
	data, err := helpers.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	// Parse the yaml file
	config := Config{}
	err = unmarshalYAML(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func unmarshalYAML(data []byte, config *Config) error {
	err := yaml.Unmarshal([]byte(data), &config)

	return err
}

func CheckTest[A, B, C, D any](config *Config, id string, test func(A, B, C, D)) func(A, B, C, D) {
	if !config.CheckExclude(id) {
		return test
	} else {
		return func(A, B, C, D) {}
	}

}
