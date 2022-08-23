package autoscaling

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

type AutoscalingGroupApi interface {
	DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
}

func GetAutoscalingGroups(svc AutoscalingGroupApi) []types.AutoScalingGroup {
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	var groups []types.AutoScalingGroup
	result, err := svc.DescribeAutoScalingGroups(context.TODO(), input)
	groups = append(groups, result.AutoScalingGroups...)
	if err != nil {
		return nil
	}
	for {
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
		result, err = svc.DescribeAutoScalingGroups(context.TODO(), input)
		if err != nil {
			return nil
		}
		groups = append(groups, result.AutoScalingGroups...)
	}
	return groups
}
