package loadbalancers

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfAccessLogsEnabled(t *testing.T) {
	type args struct {
		checkConfig   yatas.CheckConfig
		loadBalancers []LoadBalancerAttributes
		testName      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAccessLogsEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerName: "test",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test/1a2b3c4d5e6f",
						Output: &elasticloadbalancingv2.DescribeLoadBalancerAttributesOutput{
							Attributes: []types.LoadBalancerAttribute{
								{
									Key:   aws.String("access_logs.s3.enabled"),
									Value: aws.String("true"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAccessLogsEnabled(tt.args.checkConfig, tt.args.loadBalancers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "OK" {
						t.Errorf("CheckifAccessLogsEnabled() = %v, want %v", check.Results[0].Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfAccessLogsEnabledFail(t *testing.T) {
	type args struct {
		checkConfig   yatas.CheckConfig
		loadBalancers []LoadBalancerAttributes
		testName      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAccessLogsEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerName: "test",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test/1a2b3c4d5e6f",
						Output: &elasticloadbalancingv2.DescribeLoadBalancerAttributesOutput{
							Attributes: []types.LoadBalancerAttribute{
								{
									Key:   aws.String("access_logs.s3.enabled"),
									Value: aws.String("false"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAccessLogsEnabled(tt.args.checkConfig, tt.args.loadBalancers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "FAIL" {
						t.Errorf("CheckifAccessLogsEnabled() = %v, want %v", check.Results[0].Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
