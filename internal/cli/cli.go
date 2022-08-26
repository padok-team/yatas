package cli

import (
	"flag"
	"os"
	"sort"
	"time"

	"github.com/stangirard/yatas/internal/report"
	"github.com/stangirard/yatas/internal/yatas"
	"github.com/stangirard/yatas/plugins"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
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
		p := mpb.New(mpb.WithWidth(64))
		bar := p.AddBar(0, mpb.PrependDecorators(
			decor.Name("Looking at Services: "),
			decor.CountersNoUnit(" %d / %d")),
			mpb.AppendDecorators(

				decor.Percentage(),
			),
		)
		config.Progress = bar

		bar2 := p.AddBar(0,

			mpb.PrependDecorators(
				decor.Name("Running Checks: "),
				decor.CountersNoUnit("%d / %d")),
			mpb.AppendDecorators(
				decor.Percentage(),
			),
		)

		config.ProgressDetailed = bar2

	}
	checks, err := plugins.Execute(config)
	config.Lock()
	config.Progress.SetTotal(config.Progress.Current(), true)
	config.Unlock()
	time.Sleep(time.Millisecond * 100)
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
