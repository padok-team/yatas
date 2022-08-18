package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func GetTables(s aws.Config, dynamodbs []string) []*dynamodb.DescribeTableOutput {
	svc := dynamodb.NewFromConfig(s)
	var tables []*dynamodb.DescribeTableOutput
	for _, d := range dynamodbs {
		params := &dynamodb.DescribeTableInput{
			TableName: &d,
		}
		resp, err := svc.DescribeTable(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		tables = append(tables, resp)

	}
	return tables
}

type TableBackups struct {
	TableName string
	Backups   types.ContinuousBackupsDescription
}

func GetContinuousBackups(s aws.Config, tables []string) []TableBackups {
	svc := dynamodb.NewFromConfig(s)
	var continuousBackups []TableBackups
	for _, d := range tables {
		params := &dynamodb.DescribeContinuousBackupsInput{
			TableName: &d,
		}
		resp, err := svc.DescribeContinuousBackups(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		continuousBackups = append(continuousBackups, TableBackups{d, *resp.ContinuousBackupsDescription})
	}
	return continuousBackups
}
