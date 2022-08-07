package ecr

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetECRs(s *session.Session) []*ecr.Repository {
	svc := ecr.New(s)
	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int64(100),
	}
	result, err := svc.DescribeRepositories(input)
	if err != nil {
		panic(err)
	}
	return result.Repositories
}

func CheckIfImageScanningEnabled(s *session.Session, ecr []*ecr.Repository, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Image Scanning Enabled"
	check.Id = testName
	check.Description = "Check if all ECRs have image scanning enabled"
	check.Status = "OK"
	for _, ecr := range ecr {
		if *ecr.ImageScanningConfiguration.ScanOnPush != true {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "ECR " + *ecr.RepositoryName + " has image scanning disabled"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "ECR " + *ecr.RepositoryName + " has image scanning enabled"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func RunECRTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	ecr := GetECRs(s)
	config.CheckTest(c, "AWS_ECR_001", CheckIfImageScanningEnabled)(s, ecr, "AWS_CLD_001", &checks)
	return checks
}