package cloudtrail

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfCloudtrailsGlobalServiceEventsEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		cloudtrails []types.Trail
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCloudtrailsEncrypted",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				cloudtrails: []types.Trail{
					{
						Name:                       aws.String("test"),
						TrailARN:                   aws.String("test"),
						IncludeGlobalServiceEvents: aws.Bool(true),
					},
				},
				testName: "TestCheckIfCloudtrailsEncrypted",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCloudtrailsGlobalServiceEventsEnabled(tt.args.checkConfig, tt.args.cloudtrails, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfCloudtrailsEncrypted() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfCloudtrailsGlobalServiceEventsEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		cloudtrails []types.Trail
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCloudtrailsEncrypted",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				cloudtrails: []types.Trail{
					{
						Name:                       aws.String("test"),
						TrailARN:                   aws.String("test"),
						IncludeGlobalServiceEvents: aws.Bool(false),
					},
				},
				testName: "TestCheckIfCloudtrailsEncrypted",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCloudtrailsGlobalServiceEventsEnabled(tt.args.checkConfig, tt.args.cloudtrails, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfCloudtrailsEncrypted() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
