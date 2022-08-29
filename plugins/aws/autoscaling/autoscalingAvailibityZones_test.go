package autoscaling

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfInTwoAvailibilityZones(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		groups      []types.AutoScalingGroup
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestCheckIfInTwoAvailibilityZones",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				groups:      []types.AutoScalingGroup{},
				testName:    "AWS_ASG_001",
			},
			want: "OK",
		},
		{
			name: "TestCheckIfInTwoAvailibilityZones",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				groups: []types.AutoScalingGroup{
					{
						AvailabilityZones:    []string{"us-east-1a", "us-east-1b"},
						AutoScalingGroupName: aws.String("test"),
					},
				},
				testName: "AWS_ASG_001",
			},
			want: "OK",
		},
		{
			name: "TestCheckIfInTwoAvailibilityZones",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				groups: []types.AutoScalingGroup{
					{
						AvailabilityZones:    []string{"us-east-1b"},
						AutoScalingGroupName: aws.String("test"),
					},
				},
				testName: "AWS_ASG_001",
			},
			want: "FAIL",
		},
		{
			name: "TestCheckIfInTwoAvailibilityZones",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				groups: []types.AutoScalingGroup{
					{
						AvailabilityZones:    []string{"us-east-1a", "us-east-1b", "us-east-1c"},
						AutoScalingGroupName: aws.String("test"),
					},
				},
				testName: "AWS_ASG_001",
			},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfInTwoAvailibilityZones(tt.args.checkConfig, tt.args.groups, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					t.Logf("%+v", check)
					if check.Status != tt.want {
						t.Errorf("CheckIfInTwoAvailibilityZones() = %v, want %v", check.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
