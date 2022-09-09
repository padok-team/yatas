package commons

import (
	"sync"

	"github.com/stangirard/yatas/internal/helpers"
	"gopkg.in/yaml.v3"
)

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

		return test
	} else {
		return func(A, B, C) {}
	}

}

// Check Macro test is a wrapper around a category that runs all the checks in the category.
// It allows for a simpler integration of new categories without bloating the code.
func CheckMacroTest[A, B, C, D any](wg *sync.WaitGroup, config *Config, test func(A, B, C, D)) func(A, B, C, D) {
	wg.Add(1)

	return test
}
