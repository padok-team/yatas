package cloudfront

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfHTTPSOnly(t *testing.T) {
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
			name: "TestCheckIfHTTPSOnly",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						ViewerCertificate: &types.ViewerCertificate{
							MinimumProtocolVersion: types.MinimumProtocolVersionTLSv122021,
						},
						Id: aws.String("test"),
						DefaultCacheBehavior: &types.DefaultCacheBehavior{
							ViewerProtocolPolicy: types.ViewerProtocolPolicyHttpsOnly,
						},
					},
				},
				testName: "AWS_CF_001",
			},
		},
		{
			name: "TestCheckIfHTTPSOnly",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				d: []types.DistributionSummary{
					{
						ViewerCertificate: &types.ViewerCertificate{
							MinimumProtocolVersion: types.MinimumProtocolVersionTLSv122021,
						},
						Id: aws.String("test"),
						DefaultCacheBehavior: &types.DefaultCacheBehavior{
							ViewerProtocolPolicy: types.ViewerProtocolPolicyRedirectToHttps,
						},
					},
				},
				testName: "AWS_CF_001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfHTTPSOnly(tt.args.checkConfig, tt.args.d, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for r := range tt.args.checkConfig.Queue {
					if r.Status != "OK" {
						t.Errorf("CheckIfHTTPSOnly() = %v, want %v", r.Status, results.StatusOK)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
