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

var details = flag.Bool("details", false, "print detailed results")

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

func RemoveIgnored(c *yatas.Config, checks []results.Check) []results.Check {
	var newChecks []results.Check
	for _, check := range checks {
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
	return newChecks
}

func PrettyPrintChecks(checks []results.Check, c *yatas.Config) {
	flag.Parse()
	for _, check := range checks {
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

func ComparePreviousWithNew(previous []results.Check, new []results.Check) []results.Check {
	var checks []results.Check
	for _, check := range new {
		found := false
		for _, previousCheck := range previous {
			if check.Id == previousCheck.Id {
				if check.Status != previousCheck.Status {
					checks = append(checks, check)
				}
				found = true
			}
		}
		if !found {
			checks = append(checks, check)
		}
	}
	return checks
}

func ReadPreviousResults() []results.Check {
	d, err := ioutil.ReadFile("results.yaml")
	if err != nil {
		return []results.Check{}
	}
	var checks []results.Check
	err = yaml.Unmarshal(d, checks)
	if err != nil {
		panic(err)
	}
	return checks
}

func WriteChecksToFile(checks []results.Check, c *yatas.Config) {
	var checksToWrite []results.Check
	for _, check := range checks {
		if !c.CheckExclude(check.Id) {
			checksToWrite = append(checksToWrite, check)
		}
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

func ExitCode(checks []results.Check) int {
	var exitCode int
	for _, check := range checks {
		if check.Status == "FAIL" {
			exitCode = 1
		}
	}
	return exitCode
}
