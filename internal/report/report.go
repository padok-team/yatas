package report

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/fatih/color"
	"github.com/stangirard/yatas/config"
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

func countResultOkOverall(results []config.Result) (int, int) {
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

func IsIgnored(c *config.Config, r config.Result, check config.Check) bool {
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

func RemoveIgnored(c *config.Config, tests []config.Tests) []config.Tests {
	resultsTmp := []config.Tests{}
	for _, test := range tests {
		var testTpm config.Tests
		testTpm.Account = test.Account
		testTpm.Checks = []config.Check{}

		for _, check := range test.Checks {
			checkTmp := check
			checkTmp.Results = []config.Result{}
			checkTmp.InitCheck(check.Name, check.Description, check.Id)
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

func CountChecksPassedOverall(checks []config.Check) (int, int) {
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

func PrettyPrintChecks(checks []config.Tests, c *config.Config) {
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

func ComparePreviousWithNew(previous []config.Tests, new []config.Tests) []config.Tests {
	returnedResults := []config.Tests{}
	for _, tests := range new {

		var checks []config.Check
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

func ReadPreviousResults() []config.Tests {
	d, err := ioutil.ReadFile("results.yaml")
	if err != nil {
		return []config.Tests{}
	}
	var checks []config.Tests
	err = yaml.Unmarshal(d, &checks)
	if err != nil {
		panic(err)
	}
	return checks
}

func WriteChecksToFile(checks []config.Tests, c *config.Config) {
	for _, tests := range checks {
		var checksToWrite []config.Check

		for _, check := range tests.Checks {
			if !c.CheckExclude(check.Id) {
				checksToWrite = append(checksToWrite, check)
			}
		}
		tests.Checks = checksToWrite

	}
	d, err := yaml.Marshal(checks)
	if err != nil {
		panic(err)
	}

	// Write to results.yaml
	err = ioutil.WriteFile("results.yaml", d, 0644)
	if err != nil {
		panic(err)
	}

}

func ExitCode(checks []config.Tests) int {
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
