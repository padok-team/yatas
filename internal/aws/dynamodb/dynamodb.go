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

func CheckIfDynamodbEncrypted(wg *sync.WaitGroup, s aws.Config, dynamodbs []string, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Dynamodb Encryption", "Check if DynamoDB encryption is enabled", testName)
	svc := dynamodb.NewFromConfig(s)
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
	queueToAdd <- check
}

func CheckIfDynamodbContinuousBackupsEnabled(wg *sync.WaitGroup, s aws.Config, dynamodbs []string, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Dynamodb Continuous Backups", "Check if DynamoDB continuous backups are enabled", testName)
	svc := dynamodb.NewFromConfig(s)
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
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	dynamodbs := GetDynamodbs(s)
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)
	go yatas.CheckTest(&wg, c, "AWS_DYN_001", CheckIfDynamodbEncrypted)(&wg, s, dynamodbs, "AWS_DYN_001", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_DYN_002", CheckIfDynamodbContinuousBackupsEnabled)(&wg, s, dynamodbs, "AWS_DYN_002", queueResults)

	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
