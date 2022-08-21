package autoscaling

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

type mockAutoScaling func()

func (m mockAutoScaling) DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	return &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []types.AutoScalingGroup{
			{
				AutoScalingGroupName:    aws.String("123"),
				DefaultCooldown:         aws.Int32(123),
				DesiredCapacity:         aws.Int32(123),
				HealthCheckType:         aws.String("123"),
				LaunchConfigurationName: aws.String("123"),
				MaxSize:                 aws.Int32(123),
				MinSize:                 aws.Int32(123),
				VPCZoneIdentifier:       aws.String("123"),
			},
		},
	}, nil
}

func TestGetAutoscalingGroups(t *testing.T) {
	type args struct {
		svc AutoscalingGroupApi
	}
	tests := []struct {
		name string
		args args
		want []types.AutoScalingGroup
	}{
		{
			name: "One autoscaling group",
			args: args{
				svc: mockAutoScaling(nil),
			},
			want: []types.AutoScalingGroup{
				{
					AutoScalingGroupName:    aws.String("123"),
					DefaultCooldown:         aws.Int32(123),
					DesiredCapacity:         aws.Int32(123),
					HealthCheckType:         aws.String("123"),
					LaunchConfigurationName: aws.String("123"),
					MaxSize:                 aws.Int32(123),
					MinSize:                 aws.Int32(123),
					VPCZoneIdentifier:       aws.String("123"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAutoscalingGroups(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAutoscalingGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}
