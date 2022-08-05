package cmd

import (
	"fmt"

	"github.com/stangirard/yatas/internal/config"
)

func Execute() error {

	config, err := config.ParseConfig(".yatas.yml")

	for i := range config.Plugins {
		fmt.Printf("%v\n", config.Plugins[i])
		fmt.Printf("%v\n", config.Plugins[i].Name)
	}

	return err
}
