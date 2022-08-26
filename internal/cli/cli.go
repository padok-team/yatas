package cli

import (
	"flag"
	"os"
	"sort"

	"github.com/schollz/progressbar/v3"
	"github.com/stangirard/yatas/internal/report"
	"github.com/stangirard/yatas/internal/yatas"
	"github.com/stangirard/yatas/plugins"
)

var (
	compare      = flag.Bool("compare", false, "compare with previous report")
	progressflag = flag.Bool("progress", false, "show progress bar")
	ci           = flag.Bool("ci", false, "run in CI with exit code")
)

func Execute() error {
	config, err := yatas.ParseConfig(".yatas.yml")
	if err != nil {
		return err
	}
	if !*progressflag {
		config.Progress = progressbar.Default(1)
	}
	checks, err := plugins.Execute(&config)
	if err != nil {
		return err
	}
	checks = report.RemoveIgnored(&config, checks)

	if !*progressflag {
		config.Progress.Add(1)
		config.Progress.Finish()
	}
	// Sort checks by ID
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Account < checks[j].Account
	})
	for _, check := range checks {
		sort.Slice(check.Checks, func(i, j int) bool {
			return check.Checks[i].Id < check.Checks[j].Id
		})
	}

	if *compare {
		previous := report.ReadPreviousResults()
		if err != nil {
			return err
		}
		checksCompare := report.ComparePreviousWithNew(previous, checks)
		report.PrettyPrintChecks(checksCompare, &config)
		report.WriteChecksToFile(checks, &config)
		checks = checksCompare
	} else {
		report.PrettyPrintChecks(checks, &config)
		report.WriteChecksToFile(checks, &config)

	}
	if *ci {
		os.Exit(report.ExitCode(checks))
	}

	return nil
}
