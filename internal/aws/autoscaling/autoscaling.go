package autoscaling

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
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

func CheckIfDesiredCapacityMaxCapacityBelow80percent(wg *sync.WaitGroup, s aws.Config, groups []types.AutoScalingGroup, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Autoscaling Desired Capacity vs Max Capacity below 80%", "Check if all autoscaling groups have a desired capacity below 80%", testName)
	for _, group := range groups {
		if group.DesiredCapacity != nil && group.MaxSize != nil && float64(*group.DesiredCapacity) > float64(*group.MaxSize)*0.8 {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has a desired capacity above 80%"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		} else {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has a desired capacity below 80%"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {
	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	groups := GetAutoscalingGroups(s)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_ASG_001", CheckIfDesiredCapacityMaxCapacityBelow80percent)(checkConfig.Wg, checkConfig.ConfigAWS, groups, "AWS_ASG_001", checkConfig.Queue)

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
