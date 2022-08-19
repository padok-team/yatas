package vpc

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func Test_checkIfVPCFLowLogsEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		VpcFlowLogs []VpcToFlowLogs
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfVPCFLowLogsEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				VpcFlowLogs: []VpcToFlowLogs{
					{
						VpcID: "vpc-12345678",
						FlowLogs: []types.FlowLog{
							{
								FlowLogId: aws.String("fl-12345678"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfVPCFLowLogsEnabled(tt.args.checkConfig, tt.args.VpcFlowLogs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("checkIfVPCFLowLogsEnabled() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func Test_checkIfVPCFLowLogsEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		VpcFlowLogs []VpcToFlowLogs
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfVPCFLowLogsEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				VpcFlowLogs: []VpcToFlowLogs{
					{
						VpcID: "vpc-12345678",
					},
				},
			},
		},
		{
			name: "Test_checkIfVPCFLowLogsEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				VpcFlowLogs: []VpcToFlowLogs{
					{
						VpcID:    "vpc-12345678",
						FlowLogs: []types.FlowLog{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfVPCFLowLogsEnabled(tt.args.checkConfig, tt.args.VpcFlowLogs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("checkIfVPCFLowLogsEnabled() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
