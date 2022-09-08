package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/helpers"
	"github.com/vbauerster/mpb/v7"
	"gopkg.in/yaml.v3"
)

type Plugin struct {
	Name           string   `yaml:"name"`
	Enabled        bool     `yaml:"enabled"`
	Source         string   `yaml:"source"`
	Version        string   `yaml:"version"`
	Description    string   `yaml:"description"`
	Exclude        []string `yaml:"exclude"`
	Include        []string `yaml:"include"`
	Command        string   `yaml:"command"`
	Args           []string `yaml:"args"`
	ExpectedOutput string   `yaml:"expected_output"`
	ExpectedStatus int      `yaml:"expected_status"`

	// Parsed source attributes
	SourceOwner string
	SourceRepo  string
}

// InstallPath returns an installation path from the plugin directory.
func (c *Plugin) InstallPath() string {
	return filepath.Join(c.Source, c.Version, fmt.Sprintf("yatas-%s", c.Name))
}

func (c *Plugin) TagName() string {
	if c.Version == "latest" {
		return "latest"
	}
	return fmt.Sprintf("v%s", c.Version)
}

// AssetName returns a name that the asset contained in the release should meet.
// The name must be in a format similar to `yatas-aws_darwin_amd64.zip`.
func (c *Plugin) AssetName() string {
	return fmt.Sprintf("yatas-%s_%s_%s.zip", c.Name, runtime.GOOS, runtime.GOARCH)
}

func (c *Plugin) Validate() error {
	if c.Version != "" && c.Source == "" {
		return fmt.Errorf("plugin `%s`: `source` attribute cannot be omitted when specifying `version`", c.Name)
	}

	if c.Source != "" {
		if c.Version == "" {
			return fmt.Errorf("plugin `%s`: `version` attribute cannot be omitted when specifying `source`", c.Name)
		}

		parts := strings.Split(c.Source, "/")
		// Expected `github.com/owner/repo` format
		if len(parts) != 3 {
			return fmt.Errorf("plugin `%s`: `source` is invalid. Must be in the format `github.com/owner/repo`", c.Name)
		}
		if parts[0] != "github.com" {
			return fmt.Errorf("plugin `%s`: `source` is invalid. Hostname must be `github.com`", c.Name)
		}
		c.SourceOwner = parts[1]
		c.SourceRepo = parts[2]
	}

	return nil
}

type Ignore struct {
	ID     string   `yaml:"id"`
	Regex  bool     `yaml:"regex"`
	Values []string `yaml:"values"`
}

type AWS_Account struct {
	Name    string `yaml:"name"`
	Profile string `yaml:"profile"`
	SSO     bool   `yaml:"sso"`
	Region  string `yaml:"region"`
}
type Progress struct {
	Bar   *mpb.Bar
	Value int
}

type Config struct {
	Plugins []Plugin      `yaml:"plugins"`
	AWS     []AWS_Account `yaml:"aws"`
	Ignore  []Ignore      `yaml:"ignore"`
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

// Check test if a wrapper around a check that allows to verify if the check is included or excluded and add some custom logic.
// It allows for a simpler integration of new tests without bloating the code.
func CheckTest[A, B, C any](wg *sync.WaitGroup, config *Config, id string, test func(A, B, C)) func(A, B, C) {
	if !config.CheckExclude(id) && config.CheckInclude(id) {
		wg.Add(1)
		// if config.CheckProgress.Bar != nil {

		// 	config.CheckProgress.Value++
		// 	config.CheckProgress.Bar.SetTotal(int64(config.CheckProgress.Value), false)
		// 	time.Sleep(time.Millisecond * 10)

		// }
		return test
	} else {
		return func(A, B, C) {}
	}

}

// Check Macro test is a wrapper around a category that runs all the checks in the category.
// It allows for a simpler integration of new categories without bloating the code.
func CheckMacroTest[A, B, C, D any](wg *sync.WaitGroup, config *Config, test func(A, B, C, D)) func(A, B, C, D) {
	wg.Add(1)
	// TODO check
	// if config.ServiceProgress.Bar != nil {
	// 	config.ServiceProgress.Value++
	// 	config.ServiceProgress.Bar.SetTotal(int64(config.ServiceProgress.Value), false)
	// 	time.Sleep(time.Millisecond * 10)

	// }

	return test
}

// CheckConfig is a struct that contains all the information needed to run a check.
type CheckConfig struct {
	Wg          *sync.WaitGroup
	ConfigAWS   aws.Config
	Queue       chan Check
	ConfigYatas *Config
}

// Init the check config struct. Particularly useful in the categories. It allows to pass the config to the checks and allows
// them to be run in parallel by adding the results to the queue.
func (c *CheckConfig) Init(s aws.Config, config *Config) {
	c.Wg = &sync.WaitGroup{}
	c.ConfigAWS = s
	c.Queue = make(chan Check, 10)
	c.ConfigYatas = config
}
