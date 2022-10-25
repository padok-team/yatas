package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/padok-team/yatas/internal/cli"
	"github.com/padok-team/yatas/internal/report"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	generateReadme = flag.Bool("readme", false, "generate README.md checks")
	initFlag       = flag.Bool("init", false, "init yatas")
)

//go:embed .yatas.yml.example
var exampleConfig string

func WriteExampleConfig() error {
	err := os.WriteFile(".yatas.yml", []byte(exampleConfig), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Run YATAS
func run() error {
	flag.Parse()

	if *initFlag {
		err := WriteExampleConfig()
		if err != nil {
			return err
		}
		fmt.Println("Config file created in current directory ✅")
		return nil
	}

	if *generateReadme {
		return report.WriteReadme("README.md", "results.yaml")
	}
	if err := cli.Execute(); err != nil {
		return err
	}
	return nil
}
