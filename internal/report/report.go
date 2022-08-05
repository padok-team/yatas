package report

import (
	"flag"
	"fmt"

	"github.com/stangirard/yatas/internal/types"
)

var status = map[string]string{
	"OK":   "✅",
	"WARN": "⚠️",
	"FAIL": "❌",
}

var details = flag.Bool("details", false, "print detailed results")

func PrettyPrintChecks(checks []types.Check) {
	flag.Parse()
	for _, check := range checks {
		fmt.Println("✓ Check: ", check.Name, " - ", status[check.Status])
		if *details {
			fmt.Println("\tDescritpion: ", check.Description)
			fmt.Println("\tResults:")
			for _, result := range check.Results {
				fmt.Println("\t\t", status[result.Status], result.Message)
			}
		}

	}
}
