package cmd

import (
	"flag"
	"fmt"
	"sort"

	"github.com/schollz/progressbar/v3"
	"github.com/stangirard/yatas/internal/plugins"
	"github.com/stangirard/yatas/internal/report"
	"github.com/stangirard/yatas/internal/yatas"
)

var (
	compare  = flag.Bool("compare", false, "compare with previous report")
	progress = flag.Bool("progress", true, "show progress bar")
)

func Execute() error {

	config, err := yatas.ParseConfig(".yatas.yml")
	if err != nil {
		return err
	}

	if *progress {
		config.Progress = progressbar.Default(-1)
	}
	checks, err := plugins.Execute(&config)
	if err != nil {
		return err
	}
	checks = report.RemoveIgnored(&config, checks)

	if *progress {
		config.Progress.Finish()
	}
	// Sort checks by ID
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Id < checks[j].Id
	})
	fmt.Println()
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
