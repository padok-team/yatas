package report

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/stangirard/yatas/internal/results"
	"gopkg.in/yaml.v3"
)

func parseReportYaml(filename string) ([]results.Check, error) {
	var report []results.Check
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return report, err
	}
	err = yaml.Unmarshal(data, &report)
	return report, err
}

func GenerateReadme(filename string) error {
	report, err := parseReportYaml(filename)
	if err != nil {
		return err
	}
	splitFirst := ""
	splitSecond := ""
	for _, check := range report {
		split := strings.Split(check.Id, "_")
		if split[0] != splitFirst {
			splitFirst = split[0]
			fmt.Printf("\n## %s\n", split[0])
		}
		if split[1] != splitSecond {
			splitSecond = split[1]
			fmt.Printf("\n### %s\n", split[1])
		}
		fmt.Printf("- %s %s\n", check.Id, check.Name)

	}
	fmt.Printf("Number of checks: %d\n", len(report))
	return nil
}
