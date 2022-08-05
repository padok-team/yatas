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

func PrettyPrintChecks(checks []types.Check, c *config.Config) {
	flag.Parse()
	for _, check := range checks {
		if c.CheckExclude(check.Id) {
			continue
		}
		fmt.Println(status[check.Status], check.Id, check.Name)
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
