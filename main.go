package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/stangirard/yatas/internal/cmd"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Run YATAS
func run() error {
	flag.Parse()

	if err := cmd.Execute(); err != nil {
		return err
	}
	fmt.Println("\nYATAS is done 🔥")

	return nil
}
