package yatas

import (
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/helpers"
	"github.com/vbauerster/mpb/v7"
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
	sync.Mutex

	Plugins         []Plugin      `yaml:"plugins"`
	AWS             []AWS_Account `yaml:"aws"`
	Ignore          []Ignore      `yaml:"ignore"`
	Progress        *mpb.Progress
	ServiceProgress Progress
	CheckProgress   Progress
	PluginsProgress map[string]Progress
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

func CheckTest[A, B, C any](wg *sync.WaitGroup, config *Config, id string, test func(A, B, C)) func(A, B, C) {
	if !config.CheckExclude(id) && config.CheckInclude(id) {
		wg.Add(1)
		if config.CheckProgress.Bar != nil {
			config.Lock()
			config.CheckProgress.Value++
			config.CheckProgress.Bar.SetTotal(int64(config.CheckProgress.Value), false)
			config.Unlock()
			time.Sleep(time.Millisecond * 10)

		}
		return test
	} else {
		return func(A, B, C) {}
	}

}

func CheckMacroTest[A, B, C, D any](wg *sync.WaitGroup, config *Config, test func(A, B, C, D)) func(A, B, C, D) {
	wg.Add(1)
	// TODO check
	if config.ServiceProgress.Bar != nil {
		config.Lock()
		config.ServiceProgress.Value++
		config.ServiceProgress.Bar.SetTotal(int64(config.ServiceProgress.Value), false)
		config.Unlock()
		time.Sleep(time.Millisecond * 10)

	}

	return test
}

type CheckConfig struct {
	Wg          *sync.WaitGroup
	ConfigAWS   aws.Config
	Queue       chan Check
	ConfigYatas *Config
}

func (c *CheckConfig) Init(s aws.Config, config *Config) {
	c.Wg = &sync.WaitGroup{}
	c.ConfigAWS = s
	c.Queue = make(chan Check, 10)
	c.ConfigYatas = config
}
