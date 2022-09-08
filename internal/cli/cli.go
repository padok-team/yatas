package cli

import (
	"flag"
	"os"
	"sort"

	"github.com/stangirard/yatas/config"
	"github.com/stangirard/yatas/internal/report"
	"github.com/stangirard/yatas/plugins/manager"
)

var (
	compare = flag.Bool("compare", false, "compare with previous report")
	ci      = flag.Bool("ci", false, "run in CI with exit code")
	install = flag.Bool("install", false, "install plugins")
)

func Execute() error {
	configuration, err := config.ParseConfig(".yatas.yml")
	if err != nil {
		return err
	}
	for _, plugins := range configuration.Plugins {
		plugins.Validate()
		if *install {
			plugins.Install()
		}
	}
	var checks []config.Tests
	for _, plugins := range configuration.Plugins {
		checks = manager.RunPlugin(plugins, configuration)
	}

	if err != nil {
		return err
	}
	checks = report.RemoveIgnored(configuration, checks)
	// if !*progressflag {

	// }
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
		report.PrettyPrintChecks(checksCompare, configuration)
		report.WriteChecksToFile(checks, configuration)
		checks = checksCompare
	} else {
		report.PrettyPrintChecks(checks, configuration)
		report.WriteChecksToFile(checks, configuration)

	}
	if *ci {
		os.Exit(report.ExitCode(checks))
	}

	return nil
}
