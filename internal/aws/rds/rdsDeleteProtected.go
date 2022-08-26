package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfDeleteProtectionEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS have the deletion protection enabled", "Check if RDS delete protection is enabled", testName)
	for _, instance := range instances {
		if instance.DeletionProtection {
			Message := "RDS delete protection is enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS delete protection is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
