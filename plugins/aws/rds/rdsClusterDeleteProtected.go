package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfClusterDeleteProtectionEnabled(checkConfig yatas.CheckConfig, instances []types.DBCluster, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("Aurora RDS have the deletion protection enabled", "Check if Aurora RDS delete protection is enabled", testName)
	for _, instance := range instances {
		if instance.DeletionProtection != nil && *instance.DeletionProtection {
			Message := "RDS delete protection is enabled on " + *instance.DBClusterIdentifier
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS delete protection is not enabled on " + *instance.DBClusterIdentifier
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
