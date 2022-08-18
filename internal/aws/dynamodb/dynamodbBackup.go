package dynamodb

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfDynamodbContinuousBackupsEnabled(checkConfig yatas.CheckConfig, dynamodbs []TableBackups, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Dynamodb Continuous Backups", "Check if DynamoDB continuous backups are enabled", testName)
	for _, d := range dynamodbs {
		if d.Backups.ContinuousBackupsStatus != "ENABLED" {
			Message := "Dynamodb continuous backups are not enabled on " + d.TableName
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: d.TableName}
			check.AddResult(result)
		} else {
			Message := "Dynamodb continuous backups are enabled on " + d.TableName
			result := results.Result{Status: "OK", Message: Message, ResourceID: d.TableName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
