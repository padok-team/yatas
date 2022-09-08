package ecr

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfTagImmutable(checkConfig yatas.CheckConfig, ecr []types.Repository, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("ECRs tags are immutable", "Check if all ECRs are tag immutable", testName)
	for _, ecr := range ecr {
		if ecr.ImageTagMutability == types.ImageTagMutabilityMutable {
			Message := "ECR " + *ecr.RepositoryName + " is not tag immutable"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		} else {
			Message := "ECR " + *ecr.RepositoryName + " is tag immutable"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
