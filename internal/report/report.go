package report

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/fatih/color"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
	"gopkg.in/yaml.v3"
)

var status = map[string]string{
	"OK":   "✅",
	"WARN": "⚠️",
	"FAIL": "❌",
}

var (
	details = flag.Bool("details", false, "print detailed results")
	resume  = flag.Bool("resume", false, "print resume results")
)

func countResultOkOverall(results []results.Result) (int, int) {
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

func IsIgnored(c *yatas.Config, r results.Result, check results.Check) bool {
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

func RemoveIgnored(c *yatas.Config, tests []results.Tests) []results.Tests {
	for _, checks := range tests {
		var newChecks []results.Check
		for _, check := range checks.Checks {
			var checktmp results.Check
			checktmp.Id = check.Id
			checktmp.Name = check.Name
			checktmp.Status = "OK"
			checktmp.Results = []results.Result{}
			for _, result := range check.Results {
				if !IsIgnored(c, result, check) {
					if result.Status == "FAIL" {
						checktmp.Status = "FAIL"
					}
					checktmp.Results = append(checktmp.Results, result)
				}
			}
			newChecks = append(newChecks, checktmp)
		}
		checks.Checks = newChecks
	}
	return tests
}

func CountChecksPassedOverall(checks []results.Check) (int, int) {
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

func PrettyPrintChecks(checks []results.Tests, c *yatas.Config) {
	flag.Parse()
	for _, tests := range checks {
		fmt.Println("\n🔥 Account:", tests.Account, "🔥")
		if *resume {
			ok, all := CountChecksPassedOverall(tests.Checks)
			fmt.Printf("\t%d/%d checks passed\n", ok, all)
		} else {
			for _, check := range tests.Checks {
				if c.CheckExclude(check.Id) {
					continue
				}
				ok, all := countResultOkOverall(check.Results)
				count := fmt.Sprintf("%d/%d", ok, all)
				fmt.Println(status[check.Status], check.Id, check.Name, "-", count)
				if *details {
					for _, result := range check.Results {

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

func ComparePreviousWithNew(previous []results.Tests, new []results.Tests) []results.Tests {
	returnedResults := []results.Tests{}
	for _, tests := range new {

		var checks []results.Check
		for _, check := range tests.Checks {
			found := false
			for _, previousTests := range previous {
				for _, previousCheck := range previousTests.Checks {
					fmt.Println("Previous account ", previousTests.Account, " current account ", tests.Account)
					if check.Id == previousCheck.Id && tests.Account == previousTests.Account {
						if check.Status != previousCheck.Status {
							checks = append(checks, check)
							fmt.Println("Found check ", check.Id, " with status ", check.Status, " in previous results")

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

func ReadPreviousResults() []results.Tests {
	d, err := ioutil.ReadFile("results.yaml")
	if err != nil {
		return []results.Tests{}
	}
	var checks []results.Tests
	err = yaml.Unmarshal(d, &checks)
	if err != nil {
		panic(err)
	}
	return checks
}

func WriteChecksToFile(checks []results.Tests, c *yatas.Config) {
	for _, tests := range checks {
		var checksToWrite []results.Check

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

func ExitCode(checks []results.Tests) int {
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
