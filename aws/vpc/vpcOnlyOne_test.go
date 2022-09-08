package vpc

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func Test_checkIfOnlyOneVPC(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		vpcs        []types.Vpc
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfOnlyOneVPC",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				vpcs: []types.Vpc{
					{
						CidrBlock: aws.String("32.32.32.0/20"),
						VpcId:     aws.String("test"),
					},
				},
				testName: "AWS_VPC_001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfOnlyOneVPC(tt.args.checkConfig, tt.args.vpcs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("checkIfOnlyOneVPC() = %v, expected %s", t, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func Test_checkIfOnlyOneVPCFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		vpcs        []types.Vpc
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfOnlyOneVPC",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				vpcs: []types.Vpc{
					{
						CidrBlock: aws.String("32.32.32.0/20"),
						VpcId:     aws.String("test"),
					},
					{
						CidrBlock: aws.String("32.32.32.0/20"),
						VpcId:     aws.String("test"),
					},
				},
				testName: "AWS_VPC_001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfOnlyOneVPC(tt.args.checkConfig, tt.args.vpcs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("checkIfOnlyOneVPC() = %v, expected %s", t, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
