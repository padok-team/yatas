package vpc

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func Test_checkIfOnlyOneGateway(t *testing.T) {
	type args struct {
		checkConfig         yatas.CheckConfig
		vpcInternetGateways []VpcToInternetGateway
		testName            string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfOnlyOneGateway",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				vpcInternetGateways: []VpcToInternetGateway{
					{
						VpcID: "vpc-12345678",
						InternetGateways: []types.InternetGateway{
							{
								InternetGatewayId: aws.String("igw-12345678"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfOnlyOneGateway(tt.args.checkConfig, tt.args.vpcInternetGateways, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("checkIfOnlyOneGateway() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()

		})
	}
}

func Test_checkIfOnlyOneGatewayFail(t *testing.T) {
	type args struct {
		checkConfig         yatas.CheckConfig
		vpcInternetGateways []VpcToInternetGateway
		testName            string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfOnlyOneGateway",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				vpcInternetGateways: []VpcToInternetGateway{
					{
						VpcID: "vpc-12345678",
						InternetGateways: []types.InternetGateway{
							{
								InternetGatewayId: aws.String("igw-12345678"),
							},
							{
								InternetGatewayId: aws.String("igw-12345678"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfOnlyOneGateway(tt.args.checkConfig, tt.args.vpcInternetGateways, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("checkIfOnlyOneGateway() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()

		})
	}
}
