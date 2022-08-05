package cmd

import (
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/plugins"
)

func Execute() error {

	config, err := config.ParseConfig(".yatas.yml")
	if err != nil {
		return err
	}

	plugins.Run(&config)

	return err
}
