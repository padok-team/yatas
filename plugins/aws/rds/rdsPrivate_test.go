package rds

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func Test_checkIfRDSPrivateEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		instances   []types.DBInstance
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfRDSPrivateEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				instances: []types.DBInstance{
					{
						DBInstanceIdentifier: aws.String("test"),
						DBInstanceArn:        aws.String("arn:aws:rds:us-east-1:123456789012:db:test"),
						StorageEncrypted:     true,
						PubliclyAccessible:   false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfRDSPrivateEnabled(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfRDSPrivate() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()

		})
	}
}

func Test_checkIfRDSPrivateEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		instances   []types.DBInstance
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfRDSPrivateEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				instances: []types.DBInstance{
					{
						DBInstanceIdentifier: aws.String("test"),
						DBInstanceArn:        aws.String("arn:aws:rds:us-east-1:123456789012:db:test"),
						StorageEncrypted:     true,
						PubliclyAccessible:   true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfRDSPrivateEnabled(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfRDSPrivate() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()

		})
	}
}
