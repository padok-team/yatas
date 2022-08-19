package vpc

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfAtLeast2Subnets(t *testing.T) {
	type args struct {
		checkConfig  yatas.CheckConfig
		vpcToSubnets []VPCToSubnet
		testName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAtLeast2Subnets",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnets: []VPCToSubnet{
					{
						VpcID: "test",
						Subnets: []types.Subnet{
							{
								SubnetId: aws.String("test"),
							},
							{
								SubnetId: aws.String("test"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAtLeast2Subnets(tt.args.checkConfig, tt.args.vpcToSubnets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfAtLeast2Subnets() = %v", t)
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfAtLeast2SubnetsFail(t *testing.T) {
	type args struct {
		checkConfig  yatas.CheckConfig
		vpcToSubnets []VPCToSubnet
		testName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAtLeast2Subnets",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnets: []VPCToSubnet{
					{
						VpcID: "test",
						Subnets: []types.Subnet{
							{
								SubnetId: aws.String("test"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAtLeast2Subnets(tt.args.checkConfig, tt.args.vpcToSubnets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfAtLeast2Subnets() = %v", t)
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
