package ec2

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfMonitoringEnabled(t *testing.T) {
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
			name: "TestCheckIfMonitoringEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				instances: []types.Instance{
					{
						InstanceId:      aws.String("i-12345678"),
						PublicIpAddress: aws.String("192828282828"),
						Monitoring: &types.Monitoring{
							State: types.MonitoringStateEnabled,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfMonitoringEnabled(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "OK" {
						t.Errorf("CheckifEC2MonitoringEnabled() = %v, want %v", check.Results[0].Status, "PASS")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})

	}
}

func TestCheckIfMonitoringEnabledFail(t *testing.T) {
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
			name: "TestCheckIfMonitoringEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				instances: []types.Instance{
					{
						InstanceId:      aws.String("i-12345678"),
						PublicIpAddress: aws.String("192828282828"),
						Monitoring: &types.Monitoring{
							State: types.MonitoringStateDisabled,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfMonitoringEnabled(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "FAIL" {
						t.Errorf("CheckifEC2MonitoringEnabled() = %v, want %v", check.Results[0].Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})

	}
}
