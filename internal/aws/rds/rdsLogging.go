package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfLoggingEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS Logging", "Check if RDS logging is enabled", testName)
	for _, instance := range instances {
		if instance.EnabledCloudwatchLogsExports != nil {
			for _, export := range instance.EnabledCloudwatchLogsExports {
				if export == "audit" {
					Message := "RDS logging is enabled on " + *instance.DBInstanceIdentifier
					result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
					check.AddResult(result)
					break
				}
			}
		} else {
			Message := "RDS logging is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.Results = append(check.Results, result)
		}
	}
	checkConfig.Queue <- check
}
