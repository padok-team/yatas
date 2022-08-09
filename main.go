package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/stangirard/yatas/internal/cmd"
	"github.com/stangirard/yatas/internal/report"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	generateReadme = flag.Bool("readme", false, "generate README.md checks")
)

// Run YATAS
func run() error {
	flag.Parse()

	if *generateReadme {
		return report.GenerateReadme("results.yaml")
	}
	if err := cmd.Execute(); err != nil {
		return err
	}
	fmt.Println("\nYATAS is done 🔥")

	return nil
}
