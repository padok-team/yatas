package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetAutoscalingGroups(s *session.Session) []*autoscaling.Group {
	svc := autoscaling.New(s)
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	result, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil
	}
	return result.AutoScalingGroups
}

func CheckIfDesiredCapacityMaxCapacityBelow80percent(s *session.Session, groups []*autoscaling.Group, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Autoscaling DesiredCapacity MaxCapacity below 80%"
	check.Id = testName
	check.Description = "Check if all autoscaling groups have a desired capacity below 80%"
	check.Status = "OK"
	for _, group := range groups {
		if group.DesiredCapacity != nil && group.MaxSize != nil && float64(*group.DesiredCapacity) > float64(*group.MaxSize)*0.8 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has a desired capacity above 80%"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *group.AutoScalingGroupName})
		} else {
			status := "OK"
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has a desired capacity below 80%"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *group.AutoScalingGroupName})
		}
	}
	*c = append(*c, check)
}

func RunAutoscalingGroupChecks(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	groups := GetAutoscalingGroups(s)
	config.CheckTest(c, "AWS_ASG_001", CheckIfDesiredCapacityMaxCapacityBelow80percent)(s, groups, "AWS_ASG_001", &checks)
	return checks
}
