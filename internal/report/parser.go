package report

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/stangirard/yatas/internal/results"
	"gopkg.in/yaml.v3"
)

func parseReportYaml(filename string) ([]*results.Check, error) {
	var report []*results.Check
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
	splitFirstMap := make(map[string]int)
	for _, check := range report {
		split := strings.Split(check.Id, "_")
		splitFirstMap[split[0]]++
	}
	for _, check := range report {
		split := strings.Split(check.Id, "_")
		if split[0] != splitFirst {
			splitFirst = split[0]
			fmt.Printf("\n## %s - %d Checks\n", split[0], splitFirstMap[split[0]])
		}
		if split[1] != splitSecond {
			splitSecond = split[1]
			// If split is in fullName map then use fullName as name
			if fullName, ok := fullName[split[1]]; ok {
				fmt.Printf("\n### %s\n", fullName)
			} else {
				fmt.Printf("\n### %s\n", split[1])
			}
		}
		fmt.Printf("- %s %s\n", check.Id, check.Name)

	}
	return nil
}

var fullName = map[string]string{
	"S3":  "S3 Bucket",
	"VOL": "Volume",
	"BAK": "Backup",
	"RDS": "RDS",
	"VPC": "VPC",
	"CLD": "CloudTrail",
	"ECR": "ECR",
	"LMD": "Lambda",
	"DYN": "DynamoDB",
	"EC2": "EC2",
	"IAM": "IAM",
	"CFT": "Cloudfront",
	"APG": "APIGateway",
	"ASG": "AutoScaling",
	"ELB": "LoadBalancer",
}
