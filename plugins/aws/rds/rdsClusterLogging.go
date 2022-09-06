package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfClusterLoggingEnabled(checkConfig yatas.CheckConfig, instances []types.DBCluster, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("Aurora RDS logs are exported to cloudwatch", "Check if Aurora RDS logging is enabled", testName)
	for _, instance := range instances {
		if instance.EnabledCloudwatchLogsExports != nil {
			found := false
			for _, export := range instance.EnabledCloudwatchLogsExports {
				if export == "audit" {
					Message := "RDS logging is enabled on " + *instance.DBClusterIdentifier
					result := yatas.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
					check.AddResult(result)
					found = true

					break

				}
			}
			if !found {
				Message := "RDS logging is not enabled on " + *instance.DBClusterIdentifier
				result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
				check.AddResult(result)
				continue
			}
		} else {
			Message := "RDS logging is not enabled on " + *instance.DBClusterIdentifier
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
