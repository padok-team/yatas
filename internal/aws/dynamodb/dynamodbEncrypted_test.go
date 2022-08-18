package dynamodb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfDynamodbEncrypted(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		dynamodbs   []*dynamodb.DescribeTableOutput
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfDynamodbEncrypted",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				dynamodbs: []*dynamodb.DescribeTableOutput{
					{
						Table: &types.TableDescription{
							TableArn: aws.String("arn:aws:dynamodb:us-east-1:123456789012:table/DynamoDB-XXX"),
							SSEDescription: &types.SSEDescription{
								Status: types.SSEStatusEnabled,
							},
							TableName: aws.String("DynamoDB-XXX"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfDynamodbEncrypted(tt.args.checkConfig, tt.args.dynamodbs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "OK" {
						t.Errorf("CheckifDynamodbEncrypted() = %v, want %v", check.Results[0].Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfDynamodbEncryptedFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		dynamodbs   []*dynamodb.DescribeTableOutput
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfDynamodbEncrypted",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				dynamodbs: []*dynamodb.DescribeTableOutput{
					{
						Table: &types.TableDescription{
							TableArn:  aws.String("arn:aws:dynamodb:us-east-1:123456789012:table/DynamoDB-XXX"),
							TableName: aws.String("DynamoDB-XXX"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfDynamodbEncrypted(tt.args.checkConfig, tt.args.dynamodbs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "FAIL" {
						t.Errorf("CheckifDynamodbEncrypted() = %v, want %v", check.Results[0].Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
