package ecr

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetECRs(s aws.Config) []types.Repository {
	svc := ecr.NewFromConfig(s)
	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int32(100),
	}
	result, err := svc.DescribeRepositories(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Repositories
}

func CheckIfImageScanningEnabled(s aws.Config, ecr []types.Repository, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Image Scanning Enabled"
	check.Id = testName
	check.Description = "Check if all ECRs have image scanning enabled"
	check.Status = "OK"
	for _, ecr := range ecr {
		if !ecr.ImageScanningConfiguration.ScanOnPush {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "ECR " + *ecr.RepositoryName + " has image scanning disabled"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *ecr.RepositoryArn})
		} else {
			status := "OK"
			Message := "ECR " + *ecr.RepositoryName + " has image scanning enabled"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *ecr.RepositoryArn})
		}
	}
	*c = append(*c, check)
}

func RunECRTests(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	ecr := GetECRs(s)
	yatas.CheckTest(c, "AWS_ECR_001", CheckIfImageScanningEnabled)(s, ecr, "AWS_ECR_001", &checks)
	return checks
}
