package report

import (
	"fmt"

	"github.com/stangirard/yatas/internal/types"
)

var status = map[string]string{
	"OK":   "✅",
	"WARN": "⚠️",
	"FAIL": "❌",
}

func PrettyPrintChecks(checks []types.Check) {
	for _, check := range checks {
		fmt.Println("✓ Check: ", check.Name)
		fmt.Println("\tDescritpion: ", check.Description)
		fmt.Println("\tStatus: ", status[check.Status])
		for _, result := range check.Results {
			fmt.Println("\t\t🧪Result: ", status[result.Status], result.Message)
		}

	}
}
