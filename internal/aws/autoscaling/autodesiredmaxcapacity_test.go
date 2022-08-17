package autoscaling

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfDesiredCapacityMaxCapacityBelow80percent(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		groups      []types.AutoScalingGroup
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfDesiredCapacityMaxCapacityBelow80percent",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				groups: []types.AutoScalingGroup{
					{
						DesiredCapacity:      aws.Int32(1),
						MaxSize:              aws.Int32(2),
						AutoScalingGroupName: aws.String("test"),
					},
				},
				testName: "AWS_ASG_001",
			},
		},
		{
			name: "TestCheckIfDesiredCapacityMaxCapacityBelow80percent",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				groups: []types.AutoScalingGroup{
					{
						DesiredCapacity:      aws.Int32(8),
						MaxSize:              aws.Int32(12),
						AutoScalingGroupName: aws.String("test"),
					},
				},
				testName: "AWS_ASG_001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfDesiredCapacityMaxCapacityBelow80percent(tt.args.checkConfig, tt.args.groups, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					tt.args.checkConfig.Wg.Done()
					if check.Status != "OK" {
						t.Errorf("CheckIfDesiredCapacityMaxCapacityBelow80percent() = %v", t)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}