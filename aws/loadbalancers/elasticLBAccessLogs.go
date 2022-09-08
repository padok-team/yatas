package loadbalancers

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfAccessLogsEnabled(checkConfig yatas.CheckConfig, loadBalancers []LoadBalancerAttributes, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("ELB have access logs enabled", "Check if all load balancers have access logs enabled", testName)
	for _, loadBalancer := range loadBalancers {
		for _, attributes := range loadBalancer.Output.Attributes {

			if *attributes.Key == "access_logs.s3.enabled" && *attributes.Value == "true" {
				Message := "Access logs are enabled on : " + loadBalancer.LoadBalancerName
				result := yatas.Result{Status: "OK", Message: Message, ResourceID: loadBalancer.LoadBalancerArn}
				check.AddResult(result)
			} else if *attributes.Key == "access_logs.s3.enabled" && *attributes.Value == "false" {
				Message := "Access logs are not enabled on : " + loadBalancer.LoadBalancerName
				result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: loadBalancer.LoadBalancerArn}
				check.AddResult(result)
			} else {
				continue
			}
		}

	}

	checkConfig.Queue <- check
}
