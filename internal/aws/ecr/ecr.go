package ecr

import (
	"context"
	"fmt"
	"sync"

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

func CheckIfImageScanningEnabled(wg *sync.WaitGroup, s aws.Config, ecr []types.Repository, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Image Scanning Enabled", "Check if all ECRs have image scanning enabled", testName)
	for _, ecr := range ecr {
		if !ecr.ImageScanningConfiguration.ScanOnPush {
			Message := "ECR " + *ecr.RepositoryName + " has image scanning disabled"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		} else {
			Message := "ECR " + *ecr.RepositoryName + " has image scanning enabled"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	ecr := GetECRs(s)
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)
	go yatas.CheckTest(&wg, c, "AWS_ECR_001", CheckIfImageScanningEnabled)(&wg, s, ecr, "AWS_ECR_001", queueResults)
	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
