package report

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/fatih/color"
	"github.com/padok-team/yatas/plugins/commons"
	"github.com/padok-team/yatas/plugins/logger"
	"gopkg.in/yaml.v3"
)

var status = map[string]string{
	"OK":   "✅",
	"WARN": "⚠️",
	"FAIL": "❌",
}

var (
	details     = flag.Bool("details", false, "print detailed results")
	resume      = flag.Bool("resume", false, "print resume results")
	timeTaken   = flag.Bool("time", false, "print time taken for each check")
	onlyFailure = flag.Bool("only-failure", false, "print only failed checks")
)

// countResultOkOverall counts the number of OK and total results.
func countResultOkOverall(results []commons.Result) (int, int) {
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

// IsIgnored checks if a result is ignored based on the configuration.
func IsIgnored(c *commons.Config, r commons.Result, check commons.Check) bool {
	for _, ignored := range c.Ignore {
		if ignored.ID == check.Id {
			for i := range ignored.Values {
				if ignored.Regex && regexp.MustCompile(ignored.Values[i]).MatchString(r.Message) {
					return true
				} else if !ignored.Regex && r.Message == ignored.Values[i] {
					return true
				}
			}
		}
	}
	return false
}

// RemoveIgnored removes ignored checks from the given tests based on the configuration.
func RemoveIgnored(c *commons.Config, tests []commons.Tests) []commons.Tests {
	resultsTmp := []commons.Tests{}
	for _, test := range tests {
		var testTpm commons.Tests
		testTpm.Account = test.Account
		testTpm.Checks = []commons.Check{}

		for _, check := range test.Checks {
			checkTmp := check
			checkTmp.Results = []commons.Result{}
			checkTmp.InitCheck(check.Name, check.Description, check.Id, check.Categories)
			for _, result := range check.Results {
				if !IsIgnored(c, result, check) {
					checkTmp.AddResult(result)
				}
			}
			testTpm.Checks = append(testTpm.Checks, checkTmp)
		}
		resultsTmp = append(resultsTmp, testTpm)
	}
	return resultsTmp
}

// FilterHDSChecks filters checks to only include those with "HDS" in their categories
func FilterHDSChecks(tests []commons.Tests) []commons.Tests {
	resultsTmp := []commons.Tests{}
	for _, test := range tests {
		var testTmp commons.Tests
		testTmp.Account = test.Account
		testTmp.Checks = []commons.Check{}

		for _, check := range test.Checks {
			// Check if "HDS" is in any of the categories
			hasHDS := false
			for _, category := range check.Categories {
				if category == "HDS" {
					hasHDS = true
					break
				}
			}
			if hasHDS {
				testTmp.Checks = append(testTmp.Checks, check)
			}
		}
		if len(testTmp.Checks) > 0 {
			resultsTmp = append(resultsTmp, testTmp)
		}
	}
	return resultsTmp
}

// CountChecksPassedOverall counts the number of passed and total checks.
func CountChecksPassedOverall(checks []commons.Check) (int, int) {
	var ok int
	var all int
	for _, check := range checks {
		if check.Status == "OK" {
			ok++
		}
		all++
	}
	return ok, all
}

// PrettyPrintChecks prints the checks in a human-readable format.
func PrettyPrintChecks(checks []commons.Tests, c *commons.Config) {
	flag.Parse()
	for _, tests := range checks {
		ok, all := CountChecksPassedOverall(tests.Checks)

		fmt.Printf("\nName: %s (%d/%d)\n", tests.Account, ok, all)
		if !*resume {
			for _, check := range tests.Checks {
				if c.CheckExclude(check.Id) {
					continue
				}
				ok, all := countResultOkOverall(check.Results)
				count := fmt.Sprintf("%d/%d", ok, all)
				duration := fmt.Sprintf("%.2fs", check.Duration.Seconds())
				if *onlyFailure && check.Status == "OK" {
					continue
				}
				if *timeTaken {
					fmt.Println(status[check.Status], check.Id, check.Name, "-", duration, "-", count)
				} else {
					fmt.Println(status[check.Status], check.Id, check.Name, "-", count)
				}

				if *details {
					for _, result := range check.Results {
						if *onlyFailure && result.Status == "OK" {
							continue
						}
						if result.Status == "FAIL" {
							color.Red("\t" + result.Message)
						} else {
							color.Green("\t" + result.Message)
						}

					}
				}
			}
		}
	}
}

// ComparePreviousWithNew compares the previous test results with the new ones and returns the difference.
func ComparePreviousWithNew(previous []commons.Tests, new []commons.Tests) []commons.Tests {
	returnedResults := []commons.Tests{}
	for _, tests := range new {

		var checks []commons.Check
		for _, check := range tests.Checks {
			found := false
			for _, previousTests := range previous {
				for _, previousCheck := range previousTests.Checks {
					if check.Id == previousCheck.Id && tests.Account == previousTests.Account {
						if check.Status != previousCheck.Status {
							checks = append(checks, check)

						} else {
							found = true
						}

					}
				}
			}
			if !found {
				checks = append(checks, check)

			}

		}
		test := tests
		test.Checks = checks
		returnedResults = append(returnedResults, test)

	}
	return returnedResults
}

// ReadPreviousResults reads the previous test results from the results.yaml file.
func ReadPreviousResults() []commons.Tests {
	d, err := ioutil.ReadFile("results.yaml")
	if err != nil {
		return []commons.Tests{}
	}
	var checks []commons.Tests
	err = yaml.Unmarshal(d, &checks)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return checks
}

// WriteChecksToFile writes the test results to the results.yaml file.
func WriteChecksToFile(checks []commons.Tests, c *commons.Config) {
	for _, tests := range checks {
		var checksToWrite []commons.Check

		for _, check := range tests.Checks {
			if !c.CheckExclude(check.Id) {
				checksToWrite = append(checksToWrite, check)
			}
		}
		tests.Checks = checksToWrite

	}
	d, err := yaml.Marshal(checks)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Write to results.yaml
	err = ioutil.WriteFile("results.yaml", d, 0644)
	if err != nil {
		logger.Error(err.Error())
		return
	}

}

// ExitCode returns the exit code for the CI based on the test results.
func ExitCode(checks []commons.Tests) int {
	var exitCode int
	for _, tests := range checks {
		for _, check := range tests.Checks {
			if check.Status == "FAIL" {
				exitCode = 1
			}
		}
	}
	return exitCode
}
