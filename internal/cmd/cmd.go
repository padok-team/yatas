package cmd

import (
	"flag"

	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/plugins"
	"github.com/stangirard/yatas/internal/report"
)

var (
	compare = flag.Bool("compare", false, "compare with previous report")
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

	if *compare {
		previous := report.ReadPreviousResults()
		if err != nil {
			return err
		}
		checksCompare := report.ComparePreviousWithNew(previous, checks)
		report.PrettyPrintChecks(checksCompare, &config)
		report.WriteChecksToFile(checks, &config)
	} else {
		report.PrettyPrintChecks(checks, &config)
		report.WriteChecksToFile(checks, &config)

	}

	return nil
}
