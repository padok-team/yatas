package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfEC2PublicIP(checkConfig yatas.CheckConfig, instances []types.Instance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2s don't have a public IP", "Check if all instances have a public IP", testName)
	for _, instance := range instances {
		if instance.PublicIpAddress != nil {
			Message := "EC2 instance " + *instance.InstanceId + " has a public IP" + *instance.PublicIpAddress
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		} else {
			Message := "EC2 instance " + *instance.InstanceId + " has no public IP "
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
