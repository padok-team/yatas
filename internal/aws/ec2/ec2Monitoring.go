package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfMonitoringEnabled(checkConfig yatas.CheckConfig, instances []types.Instance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2s have the monitoring option enabled", "Check if all instances have monitoring enabled", testName)
	for _, instance := range instances {
		if instance.Monitoring.State != types.MonitoringStateEnabled {
			Message := "EC2 instance " + *instance.InstanceId + " has no monitoring enabled"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		} else {
			Message := "EC2 instance " + *instance.InstanceId + " has monitoring enabled"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
