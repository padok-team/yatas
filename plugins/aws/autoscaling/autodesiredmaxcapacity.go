package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfDesiredCapacityMaxCapacityBelow80percent(checkConfig yatas.CheckConfig, groups []types.AutoScalingGroup, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("Autoscaling maximum capacity is below 80%", "Check if all autoscaling groups have a desired capacity below 80%", testName)
	for _, group := range groups {
		if group.DesiredCapacity != nil && group.MaxSize != nil && float64(*group.DesiredCapacity) > float64(*group.MaxSize)*0.8 {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has a desired capacity above 80%"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		} else {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has a desired capacity below 80%"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
