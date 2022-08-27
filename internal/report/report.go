package report

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/fatih/color"
	"github.com/stangirard/yatas/internal/yatas"
	"gopkg.in/yaml.v3"
)

var status = map[string]string{
	"OK":   "✅",
	"WARN": "⚠️",
	"FAIL": "❌",
}

var (
	details   = flag.Bool("details", false, "print detailed results")
	resume    = flag.Bool("resume", false, "print resume results")
	timeTaken = flag.Bool("time", false, "print time taken for each check")
)

func countResultOkOverall(results []yatas.Result) (int, int) {
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

func IsIgnored(c *yatas.Config, r yatas.Result, check yatas.Check) bool {
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

func RemoveIgnored(c *yatas.Config, tests []yatas.Tests) []yatas.Tests {
	resultsTmp := []yatas.Tests{}
	for _, test := range tests {
		var testTpm yatas.Tests
		testTpm.Account = test.Account
		testTpm.Checks = []yatas.Check{}

		for _, check := range test.Checks {
			checkTmp := check
			checkTmp.Results = []yatas.Result{}
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

func CountChecksPassedOverall(checks []yatas.Check) (int, int) {
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

func PrettyPrintChecks(checks []yatas.Tests, c *yatas.Config) {
	flag.Parse()
	for _, tests := range checks {
		fmt.Println("\nName:", tests.Account)
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
				duration := fmt.Sprintf("%.2fs", check.Duration.Seconds())
				if *timeTaken {
					fmt.Println(status[check.Status], check.Id, check.Name, "-", duration, "-", count)
				}
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

func ComparePreviousWithNew(previous []yatas.Tests, new []yatas.Tests) []yatas.Tests {
	returnedResults := []yatas.Tests{}
	for _, tests := range new {

		var checks []yatas.Check
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

func ReadPreviousResults() []yatas.Tests {
	d, err := ioutil.ReadFile("results.yaml")
	if err != nil {
		return []yatas.Tests{}
	}
	var checks []yatas.Tests
	err = yaml.Unmarshal(d, &checks)
	if err != nil {
		panic(err)
	}
	return checks
}

func WriteChecksToFile(checks []yatas.Tests, c *yatas.Config) {
	for _, tests := range checks {
		var checksToWrite []yatas.Check

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

func ExitCode(checks []yatas.Tests) int {
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
