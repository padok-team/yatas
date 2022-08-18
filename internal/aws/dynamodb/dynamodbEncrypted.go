package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfDynamodbEncrypted(checkConfig yatas.CheckConfig, dynamodbs []*dynamodb.DescribeTableOutput, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Dynamodb Encryption", "Check if DynamoDB encryption is enabled", testName)
	for _, d := range dynamodbs {
		if d.Table != nil && d.Table.SSEDescription != nil && d.Table.SSEDescription.Status == "ENABLED" {
			Message := "Dynamodb encryption is enabled on " + *d.Table.TableName
			result := results.Result{Status: "OK", Message: Message, ResourceID: *d.Table.TableArn}
			check.AddResult(result)

		} else {
			Message := "Dynamodb encryption is not enabled on " + *d.Table.TableName
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *d.Table.TableArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
