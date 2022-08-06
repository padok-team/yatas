package report

import (
	"flag"
	"fmt"
	"io/ioutil"

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
					fmt.Println("\t" + result.Message)
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
