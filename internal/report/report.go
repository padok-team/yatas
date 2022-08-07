package report

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/fatih/color"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/types"
	"gopkg.in/yaml.v3"
)

var status = map[string]string{
	"OK":   "✅",
	"WARN": "⚠️",
	"FAIL": "❌",
}

var details = flag.Bool("details", false, "print detailed results")

func countResultOkOverall(results []types.Result) (int, int) {
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

func IsIgnored(c *config.Config, r types.Result, check types.Check) bool {
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

func RemoveIgnored(c *config.Config, checks []types.Check) []types.Check {
	var newChecks []types.Check
	for _, check := range checks {
		var checktmp types.Check
		checktmp.Id = check.Id
		checktmp.Name = check.Name
		checktmp.Status = "OK"
		checktmp.Results = []types.Result{}
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

func PrettyPrintChecks(checks []types.Check, c *config.Config) {
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

func ComparePreviousWithNew(previous []types.Check, new []types.Check) []types.Check {
	var checks []types.Check
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

func ReadPreviousResults() []types.Check {
	d, err := ioutil.ReadFile("results.yaml")
	if err != nil {
		return []types.Check{}
	}
	var checks []types.Check
	err = yaml.Unmarshal(d, &checks)
	if err != nil {
		panic(err)
	}
	return checks
}

func WriteChecksToFile(checks []types.Check, c *config.Config) {
	var checksToWrite []types.Check
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
