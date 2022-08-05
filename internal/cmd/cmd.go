package cmd

import (
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/plugins"
	"github.com/stangirard/yatas/internal/report"
)

func Execute() error {

	config, err := config.ParseConfig(".yatas.yml")
	if err != nil {
		return err
	}

	checks, err := plugins.Execute(&config)
	if err != nil {
		return err
	}
	report.PrettyPrintChecks(checks)

	return nil
}
