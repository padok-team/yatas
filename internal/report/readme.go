package report

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/padok-team/yatas/internal/helpers"
)

func WriteReadme(readmeFile string, resultFile string) error {
	// Open the readme File
	file, err := helpers.ReadFile(readmeFile)
	if err != nil {
		return err
	}
	readme := GenerateReadme(resultFile)

	re := regexp.MustCompile("(?s)(?:<!-- BEGIN_YATAS -->)(.*)(?:<!-- END_YATAS -->)")
	s := re.ReplaceAllString(string(file), fmt.Sprintf("<!-- BEGIN_YATAS -->\n%s\n<!-- END_YATAS -->", readme))
	err = ioutil.WriteFile(readmeFile, []byte(s), 0644)
	if err != nil {
		return err
	}
	return nil

}

func GenerateReadme(filename string) string {
	report, err := parseReportYaml(filename)
	readme := ""
	if err != nil {
		panic(err)
	}
	splitFirst := ""
	splitSecond := ""
	splitFirstMap := make(map[string]int)
	for _, tests := range report {
		for _, check := range tests.Checks {
			split := strings.Split(check.Id, "_")
			splitFirstMap[split[0]]++
		}
		for _, check := range tests.Checks {
			split := strings.Split(check.Id, "_")
			if split[0] != splitFirst {
				splitFirst = split[0]
				readme += fmt.Sprintf("\n## %s - %d Checks\n", split[0], splitFirstMap[split[0]])
			}
			if split[1] != splitSecond {
				splitSecond = split[1]
				// If split is in fullName map then use fullName as name
				if fullName, ok := fullName[split[1]]; ok {
					readme += fmt.Sprintf("\n### %s\n", fullName)
				} else {
					readme += fmt.Sprintf("\n### %s\n", split[1])
				}
			}
			readme += fmt.Sprintf("- %s %s\n", check.Id, check.Name)

		}
		break
	}
	return readme
}
