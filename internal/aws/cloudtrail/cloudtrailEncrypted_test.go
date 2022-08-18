package cloudtrail

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfCloudtrailsEncrypted(t *testing.T) {
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
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				cloudtrails: []types.Trail{
					{
						Name:     aws.String("test"),
						KmsKeyId: aws.String("test"),
						TrailARN: aws.String("test"),
					},
				},
				testName: "TestCheckIfCloudtrailsEncrypted",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCloudtrailsEncrypted(tt.args.checkConfig, tt.args.cloudtrails, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "OK" {
						t.Errorf("CheckIfCloudtrailsEncrypted() = %v, want %v", check.Results[0].Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfCloudtrailsEncryptedFail(t *testing.T) {
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
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				cloudtrails: []types.Trail{
					{
						Name:     aws.String("test"),
						TrailARN: aws.String("test"),
					},
				},
				testName: "TestCheckIfCloudtrailsEncrypted",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCloudtrailsEncrypted(tt.args.checkConfig, tt.args.cloudtrails, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "FAIL" {
						t.Errorf("CheckIfCloudtrailsEncrypted() = %v, want %v", check.Results[0].Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
