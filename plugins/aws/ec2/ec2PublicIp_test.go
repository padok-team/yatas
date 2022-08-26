package ec2

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfEC2PublicIPFAIL(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		instances   []types.Instance
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfEC2PublicIP",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				instances: []types.Instance{
					{
						InstanceId:      aws.String("i-12345678"),
						PublicIpAddress: aws.String("192828282828"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfEC2PublicIP(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckifEC2PublicIP() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfEC2PublicIP(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		instances   []types.Instance
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfEC2PublicIP",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				instances: []types.Instance{
					{
						InstanceId: aws.String("i-12345678"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfEC2PublicIP(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckifEC2PublicIP() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
