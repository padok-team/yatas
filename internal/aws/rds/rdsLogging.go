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
	check.InitCheck("RDS logs are exported to cloudwatch", "Check if RDS logging is enabled", testName)
	for _, instance := range instances {
		if instance.EnabledCloudwatchLogsExports != nil {
			found := false
			for _, export := range instance.EnabledCloudwatchLogsExports {
				if export == "audit" {
					Message := "RDS logging is enabled on " + *instance.DBInstanceIdentifier
					result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
					check.AddResult(result)
					found = true

					break

				}
			}
			if !found {
				Message := "RDS logging is not enabled on " + *instance.DBInstanceIdentifier
				result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
				check.AddResult(result)
				continue
			}
		} else {
			Message := "RDS logging is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
