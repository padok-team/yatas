package cloudfront

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfCloudfrontTLS1_2Minimum(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		d           []types.DistributionSummary
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCloudfrontTLS1_2Minimum",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						ViewerCertificate: &types.ViewerCertificate{
							MinimumProtocolVersion: types.MinimumProtocolVersionTLSv122021,
						},
						Id: aws.String("test"),
					},
				},
				testName: "AWS_CF_001",
			},
		},
		{
			name: "TestCheckIfCloudfrontTLS1_2Minimum",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						ViewerCertificate: &types.ViewerCertificate{
							MinimumProtocolVersion: types.MinimumProtocolVersionTLSv122019,
						},
						Id: aws.String("test"),
					},
				},
				testName: "AWS_CF_001",
			},
		},
		{
			name: "TestCheckIfCloudfrontTLS1_2Minimum",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						ViewerCertificate: &types.ViewerCertificate{
							MinimumProtocolVersion: types.MinimumProtocolVersionTLSv122018,
						},
						Id: aws.String("test"),
					},
				},
				testName: "AWS_CF_001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCloudfrontTLS1_2Minimum(tt.args.checkConfig, tt.args.d, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfCloudfrontTLS1_2Minimum() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfCloudfrontTLS1_2MinimumFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		d           []types.DistributionSummary
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCloudfrontTLS1_2Minimum",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						ViewerCertificate: &types.ViewerCertificate{
							MinimumProtocolVersion: types.MinimumProtocolVersionTLSv1,
						},
						Id: aws.String("test"),
					},
				},
				testName: "AWS_CF_001",
			},
		},
		{
			name: "TestCheckIfCloudfrontTLS1_2Minimum",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						Id: aws.String("test"),
					},
				},
				testName: "AWS_CF_001",
			},
		},
		{
			name: "TestCheckIfCloudfrontTLS1_2Minimum",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						ViewerCertificate: &types.ViewerCertificate{
							MinimumProtocolVersion: types.MinimumProtocolVersionTLSv12016,
						},
						Id: aws.String("test"),
					},
				},
				testName: "AWS_CF_001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCloudfrontTLS1_2Minimum(tt.args.checkConfig, tt.args.d, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfCloudfrontTLS1_2Minimum() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
