package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/stangirard/yatas/internal/cmd"
)

var (
	printVersion = flag.Bool("version", false, "print version and exit") // Displays the version of the programm
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//go:embed VERSION
var version string // Version of the application embeded

// Run YATAS
func run() error {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		return nil
	}

	if err := cmd.Execute(); err != nil {
		return err
	}
	fmt.Println("\nYATAS is done 🔥")

	return nil
}
