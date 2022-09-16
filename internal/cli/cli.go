package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/stangirard/yatas/internal/report"
	"github.com/stangirard/yatas/plugins/commons"
	"github.com/stangirard/yatas/plugins/manager"
)

var (
	compare = flag.Bool("compare", false, "compare with previous report")
	ci      = flag.Bool("ci", false, "run in CI with exit code")
	install = flag.Bool("install", false, "install plugins")
)

func initialisePlugins(configuration commons.Config) error {
	for _, plugins := range configuration.Plugins {
		plugins.Validate()
		if *install {
			_, err := plugins.Install()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RunChecksPlugins(configuration *commons.Config, checks *[]commons.Tests) {
	for _, plugins := range configuration.Plugins {
		if plugins.Type == "checks" || plugins.Type == "" {

			latestVersion, _ := commons.GetLatestReleaseTag(plugins)

			if plugins.Version != latestVersion {
				fmt.Println("New version available for plugin " + plugins.Name + " : " + latestVersion)
			}
			checksFromPlugin := manager.RunPlugin(plugins, configuration)
			*checks = append(*checks, checksFromPlugin...)
		}
	}
}

func parseConfig() (*commons.Config, error) {
	configuration, err := commons.ParseConfig(".yatas.yml")
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

func compareResults(configuration *commons.Config, checks *[]commons.Tests) {
	// Compare with previous report
	if *compare {
		previous := report.ReadPreviousResults()
		checksCompare := report.ComparePreviousWithNew(previous, *checks)
		report.PrettyPrintChecks(checksCompare, configuration)
		report.WriteChecksToFile(*checks, configuration)
		checks = &checksCompare
	} else {
		report.PrettyPrintChecks(*checks, configuration)
		report.WriteChecksToFile(*checks, configuration)

	}
}

func ciReporting(checks []commons.Tests) {
	if *ci {
		os.Exit(report.ExitCode(checks))
	}
}

func runModPlugins(configuration *commons.Config, checks *[]commons.Tests) bool {
	mod := false
	for _, plugins := range configuration.Plugins {
		if plugins.Type == "mod" {
			mod = true
			latestVersion, _ := commons.GetLatestReleaseTag(plugins)

			if plugins.Version != latestVersion {
				fmt.Println("New version available for plugin " + plugins.Name + " : " + latestVersion)
			}
			checksFromPlugin := manager.RunPlugin(plugins, configuration)
			*checks = append(*checks, checksFromPlugin...)
		}
	}
	return mod
}

func RunReportPlugins(configuration *commons.Config, checks *[]commons.Tests) {
	for _, plugins := range configuration.Plugins {
		if plugins.Type == "report" {
			latestVersion, _ := commons.GetLatestReleaseTag(plugins)

			if plugins.Version != latestVersion {
				fmt.Println("New version available for plugin " + plugins.Name + " : " + latestVersion)
			}
			manager.RunPlugin(plugins, configuration)
		}
	}
}

// Execute YATAS
func Execute() error {
	// Parse the config file
	configuration, err := parseConfig()
	if err != nil {
		return err
	}

	// Initialise plugins by installing them if needed and checking if the config is valid
	err = initialisePlugins(*configuration)
	if err != nil {
		return err
	}

	checks := []commons.Tests{}

	// Run Mods plugins
	if runModPlugins(configuration, &checks) {
		return nil
	}

	// Run plugins
	RunChecksPlugins(configuration, &checks)

	// Clean results
	checks = report.RemoveIgnored(configuration, checks)

	// Sort checks by ID
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Account < checks[j].Account
	})
	for _, check := range checks {
		sort.Slice(check.Checks, func(i, j int) bool {
			return check.Checks[i].Id < check.Checks[j].Id
		})
	}

	// Compare with previous report
	compareResults(configuration, &checks)

	// CI reporting
	ciReporting(checks)

	// Run report plugins
	RunReportPlugins(configuration, &checks)

	return nil
}
