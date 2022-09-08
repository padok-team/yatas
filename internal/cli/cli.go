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
)

func Execute() error {
	config, err := config.ParseConfig(".yatas.yml")
	if err != nil {
		return err
	}
	checks := manager.RunPlugin("aws", config)

	if err != nil {
		return err
	}
	checks = report.RemoveIgnored(config, checks)
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
		report.PrettyPrintChecks(checksCompare, config)
		report.WriteChecksToFile(checks, config)
		checks = checksCompare
	} else {
		report.PrettyPrintChecks(checks, config)
		report.WriteChecksToFile(checks, config)

	}
	if *ci {
		os.Exit(report.ExitCode(checks))
	}

	return nil
}
