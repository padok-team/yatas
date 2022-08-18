package ecr

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfImageScanningEnabled(checkConfig yatas.CheckConfig, ecr []types.Repository, testName string) {
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
	checkConfig.Queue <- check
}
