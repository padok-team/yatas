package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func checkIfClusterBackupEnabled(checkConfig yatas.CheckConfig, instances []types.DBCluster, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("Aurora RDS are backedup automatically with PITR", "Check if Aurora RDS backup is enabled", testName)
	for _, instance := range instances {
		if instance.BackupRetentionPeriod == nil || *instance.BackupRetentionPeriod == 0 {
			Message := "RDS backup is not enabled on " + *instance.DBClusterIdentifier
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS backup is enabled on " + *instance.DBClusterIdentifier
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
