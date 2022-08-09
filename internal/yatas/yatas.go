package yatas

import (
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/stangirard/yatas/internal/helpers"
	"gopkg.in/yaml.v3"
)

type Plugin struct {
	Name           string   `yaml:"name"`
	Enabled        bool     `yaml:"enabled"`
	Description    string   `yaml:"description"`
	Exclude        []string `yaml:"exclude"`
	Include        []string `yaml:"include"`
	Command        string   `yaml:"command"`
	Args           []string `yaml:"args"`
	ExpectedOutput string   `yaml:"expected_output"`
	ExpectedStatus int      `yaml:"expected_status"`
}

type Ignore struct {
	ID     string   `yaml:"id"`
	Regex  bool     `yaml:"regex"`
	Values []string `yaml:"values"`
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
	Ignore   []Ignore `yaml:"ignore"`
	Progress *progressbar.ProgressBar
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
	//Split Id at _ and get first value
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
	if config.Progress != nil {
		config.Progress.Add(1)
	}
	if !config.CheckExclude(id) && config.CheckInclude(id) {
		return test
	} else {
		return func(A, B, C, D) {}
	}

}
