package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetDynamodbs(s *session.Session) []*string {
	svc := dynamodb.New(s)
	dynamodbInput := &dynamodb.ListTablesInput{}
	result, err := svc.ListTables(dynamodbInput)
	if err != nil {
		panic(err)
	}
	return result.TableNames
}

func CheckIfDynamodbEncrypted(s *session.Session, dynamodbs []*string, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Dynamodb Encryption"
	check.Id = testName
	check.Description = "Check if DynamoDB encryption is enabled"
	check.Status = "OK"
	svc := dynamodb.New(s)
	for _, d := range dynamodbs {
		params := &dynamodb.DescribeTableInput{
			TableName: d,
		}
		resp, err := svc.DescribeTable(params)
		if err != nil {
			panic(err)
		}
		if *resp.Table.SSEDescription.Status != "ENABLED" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Dynamodb encryption is not enabled on " + *d
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Dynamodb encryption is enabled on " + *d
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfDynamodbContinuousBackupsEnabled(s *session.Session, dynamodbs []*string, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Dynamodb Continuous Backups"
	check.Id = testName
	check.Description = "Check if DynamoDB continuous backups are enabled"
	check.Status = "OK"
	svc := dynamodb.New(s)
	for _, d := range dynamodbs {
		params := &dynamodb.DescribeContinuousBackupsInput{
			TableName: d,
		}
		resp, err := svc.DescribeContinuousBackups(params)
		if err != nil {
			panic(err)
		}
		if *resp.ContinuousBackupsDescription.ContinuousBackupsStatus != "ENABLED" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Dynamodb continuous backups are not enabled on " + *d
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Dynamodb continuous backups are enabled on " + *d
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func RunDynamodbTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	dynamodbs := GetDynamodbs(s)
	config.CheckTest(c, "AWS_DYN_001", CheckIfDynamodbEncrypted)(s, dynamodbs, "AWS_DYN_001", &checks)
	config.CheckTest(c, "AWS_DYN_002", CheckIfDynamodbContinuousBackupsEnabled)(s, dynamodbs, "AWS_DYN_002", &checks)
	return checks
}
