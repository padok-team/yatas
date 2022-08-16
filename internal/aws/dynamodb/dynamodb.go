package dynamodb

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetDynamodbs(s aws.Config) []string {
	svc := dynamodb.NewFromConfig(s)
	dynamodbInput := &dynamodb.ListTablesInput{}
	result, err := svc.ListTables(context.TODO(), dynamodbInput)
	if err != nil {
		panic(err)
	}
	return result.TableNames
}

func CheckIfDynamodbEncrypted(checkConfig yatas.CheckConfig, dynamodbs []string, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Dynamodb Encryption", "Check if DynamoDB encryption is enabled", testName)
	svc := dynamodb.NewFromConfig(checkConfig.ConfigAWS)
	for _, d := range dynamodbs {
		params := &dynamodb.DescribeTableInput{
			TableName: &d,
		}
		resp, err := svc.DescribeTable(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.Table != nil && resp.Table.SSEDescription != nil && resp.Table.SSEDescription.Status == "ENABLED" {
			Message := "Dynamodb encryption is enabled on " + d
			result := results.Result{Status: "OK", Message: Message, ResourceID: d}
			check.AddResult(result)

		} else {
			Message := "Dynamodb encryption is not enabled on " + d
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: d}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfDynamodbContinuousBackupsEnabled(checkConfig yatas.CheckConfig, dynamodbs []string, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Dynamodb Continuous Backups", "Check if DynamoDB continuous backups are enabled", testName)
	svc := dynamodb.NewFromConfig(checkConfig.ConfigAWS)
	for _, d := range dynamodbs {
		params := &dynamodb.DescribeContinuousBackupsInput{
			TableName: &d,
		}
		resp, err := svc.DescribeContinuousBackups(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.ContinuousBackupsDescription.ContinuousBackupsStatus != "ENABLED" {
			Message := "Dynamodb continuous backups are not enabled on " + d
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: d}
			check.AddResult(result)
		} else {
			Message := "Dynamodb continuous backups are enabled on " + d
			result := results.Result{Status: "OK", Message: Message, ResourceID: d}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	dynamodbs := GetDynamodbs(s)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_DYN_001", CheckIfDynamodbEncrypted)(checkConfig, dynamodbs, "AWS_DYN_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_DYN_002", CheckIfDynamodbContinuousBackupsEnabled)(checkConfig, dynamodbs, "AWS_DYN_002")

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
