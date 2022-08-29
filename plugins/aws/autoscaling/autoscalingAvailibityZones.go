package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfInTwoAvailibilityZones(checkConfig yatas.CheckConfig, groups []types.AutoScalingGroup, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("Autoscaling group are in two availability zones", "Check if all autoscaling groups have at least two availability zones", testName)
	for _, group := range groups {
		if len(group.AvailabilityZones) < 2 {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has less than two availability zones"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		} else {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has two availability zones"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
