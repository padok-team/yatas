package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func checkIfClusterEncryptionEnabled(checkConfig yatas.CheckConfig, instances []types.DBCluster, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("Aurora RDS are encrypted", "Check if Aurora RDS encryption is enabled", testName)
	for _, instance := range instances {
		if !instance.StorageEncrypted {
			Message := "RDS encryption is not enabled on " + *instance.DBClusterIdentifier
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS encryption is enabled on " + *instance.DBClusterIdentifier
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
