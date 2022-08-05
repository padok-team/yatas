package config

import (
	"github.com/stangirard/yatas/internal/helpers"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Plugins []struct {
		Name          string `yaml:"name"`
		Enabled       bool   `yaml:"enabled"`
		Description   string `yaml:"description"`
		CloudProvider string `yaml:"cloud_provider"`
	} `yaml:"plugins"`
	AWS struct {
		Enabled bool `yaml:"enabled"`
		Account struct {
			Profile string `yaml:"profile"`
			SSO     bool   `yaml:"sso"`
			Region  string `yaml:"region"`
		} `yaml:"account"`
	} `yaml:"aws"`
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
