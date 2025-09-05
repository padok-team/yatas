package cli

import (
	"flag"
	"os"
	"sort"

	"github.com/padok-team/yatas/internal/report"
	"github.com/padok-team/yatas/plugins/commons"
	"github.com/padok-team/yatas/plugins/logger"
	"github.com/padok-team/yatas/plugins/manager"
)

var (
	compare = flag.Bool("compare", false, "compare with previous report")
	ci      = flag.Bool("ci", false, "run in CI with exit code")
	install = flag.Bool("install", false, "install plugins")
	hds     = flag.Bool("hds", false, "only run HDS checks")
)

// initialisePlugins installs plugins if needed and validates their configuration.
func initialisePlugins(configuration commons.Config) error {
	for _, plugin := range configuration.Plugins {
		err := plugin.Validate()
		if err != nil {
			logger.Error("Error validating plugin", "plugin_name", plugin.Name, "error", err)
			return err
		}
		if *install {
			_, err := plugin.Install()
			if err != nil {
				logger.Error("Error installing plugin", "plugin_name", plugin.Name, "error", err)
				return err
			}
		}
	}
	return nil
}

// runChecksPlugins runs checks plugins and appends their test results to the checks slice.
func runChecksPlugins(configuration *commons.Config, checks *[]commons.Tests) {
	for _, plugin := range configuration.Plugins {
		if plugin.Type == "checks" || plugin.Type == "" {
			latestVersion, _ := commons.GetLatestReleaseTag(plugin)
			if plugin.Version != latestVersion {
				logger.Info("New version available for plugin", "plugin_name", plugin.Name, "latest_version", latestVersion)
			}
			checksFromPlugin := manager.RunPlugin(plugin, configuration)
			*checks = append(*checks, checksFromPlugin...)
		}
	}
}

// parseConfig parses the configuration file and returns a Config object.
func parseConfig() (*commons.Config, error) {
	configuration, err := commons.ParseConfig(".yatas.yml")
	if err != nil {
		logger.Error("Error parsing config", "error", err)
		return nil, err
	}
	return configuration, nil
}

// compareResults compares the current test results with the previous ones and writes the results to a file.
func compareResults(configuration *commons.Config, checks *[]commons.Tests) {
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

// ciReporting sets the exit code for the CI based on the test results.
func ciReporting(checks []commons.Tests) {
	if *ci {
		os.Exit(report.ExitCode(checks))
	}
}

// runModPlugins runs mod plugins and appends their test results to the checks slice. Returns true if any mod plugin is executed.
func runModPlugins(configuration *commons.Config, checks *[]commons.Tests) bool {
	mod := false
	for _, plugin := range configuration.Plugins {
		if plugin.Type == "mod" {
			mod = true
			latestVersion, err := commons.GetLatestReleaseTag(plugin)
			if err != nil {
				logger.Error("Error getting latest release tag", "plugin_name", plugin.Name, "error", err)
			}
			if plugin.Version != latestVersion {
				logger.Info("New version available for plugin", "plugin_name", plugin.Name, "latest_version", latestVersion)
			}
			checksFromPlugin := manager.RunPlugin(plugin, configuration)
			*checks = append(*checks, checksFromPlugin...)
		}
	}
	return mod
}

// runReportPlugins runs report plugins.
func runReportPlugins(configuration *commons.Config, checks *[]commons.Tests) {
	for _, plugin := range configuration.Plugins {
		if plugin.Type == "report" {
			logger.Debug("Running report plugin", "plugin_name", plugin.Name)
			latestVersion, _ := commons.GetLatestReleaseTag(plugin)

			if plugin.Version != latestVersion {
				logger.Info("New version available for plugin", "plugin_name", plugin.Name, "latest_version", latestVersion)
			}
			logger.Debug("Running report plugin", "plugin_name", plugin.Name)
			manager.RunPlugin(plugin, configuration)
		}
	}
}

// Execute runs YATAS.
func Execute() error {

	// Parse the config file
	configuration, err := parseConfig()
	if err != nil {
		logger.Error("Error parsing config", "error", err)
		return err
	}

	// Initialise plugins by installing them if needed and checking if the config is valid
	err = initialisePlugins(*configuration)
	if err != nil {
		logger.Error("Error initializing plugins", "error", err)
		return err
	}

	checks := []commons.Tests{}

	// Run Mods plugins
	if runModPlugins(configuration, &checks) {
		logger.Debug("Mod plugins executed, skipping checks")
		return nil
	}

	// Run checks plugins
	logger.Debug("Running checks plugins")
	runChecksPlugins(configuration, &checks)

	if *hds {
		// Filter HDS checks if --hds flag is set
		logger.Debug("Filtering HDS checks only")
		checks = report.FilterHDSChecks(checks)
	} else {
		// Clean results
		logger.Debug("Cleaning results")
		checks = report.RemoveIgnored(configuration, checks)
	}

	// Sort checks by ID
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Account < checks[j].Account
	})
	for _, check := range checks {
		sort.Slice(check.Checks, func(i, j int) bool {
			return check.Checks[i].Id < check.Checks[j].Id
		})
	}

	configuration.Tests = checks

	// Compare with previous report
	logger.Debug("Comparing with previous report")
	compareResults(configuration, &checks)

	// CI reporting
	logger.Debug("CI reporting")
	ciReporting(checks)

	// Run report plugins
	logger.Debug("Running report plugins")
	runReportPlugins(configuration, &configuration.Tests)

	logger.Debug("Done")
	return nil
}
