package report

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/types"
)

var status = map[string]string{
	"OK":   "✅",
	"WARN": "⚠️",
	"FAIL": "❌",
}

var details = flag.Bool("details", false, "print detailed results")

func countResultOkOverall(results []types.Result) (int, int) {
	var ok int
	var all int
	for _, result := range results {
		if result.Status == "OK" {
			ok++
		}
		all++
	}
	return ok, all
}

func PrettyPrintChecks(checks []types.Check, c *config.Config) {
	flag.Parse()
	for _, check := range checks {
		if c.CheckExclude(check.Id) {
			continue
		}
		ok, all := countResultOkOverall(check.Results)
		count := fmt.Sprintf("%d/%d", ok, all)
		fmt.Println(status[check.Status], check.Id, check.Name, "-", count)
		if *details {
			for _, result := range check.Results {
				if result.Status == "FAIL" {
					color.Red("\t" + result.Message)
				} else {
					fmt.Println("\t" + result.Message)
				}

			}
		}

	}
}
