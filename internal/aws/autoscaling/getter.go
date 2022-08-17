package autoscaling

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

func GetAutoscalingGroups(s aws.Config) []types.AutoScalingGroup {
	svc := autoscaling.NewFromConfig(s)
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	result, err := svc.DescribeAutoScalingGroups(context.TODO(), input)
	if err != nil {
		return nil
	}
	return result.AutoScalingGroups
}
