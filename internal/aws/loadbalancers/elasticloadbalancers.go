package loadbalancers

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetElasticLoadBalancers(s aws.Config) []types.LoadBalancer {
	svc := elasticloadbalancingv2.NewFromConfig(s)
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		PageSize: aws.Int32(100),
	}
	result, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.LoadBalancers
}

func CheckIfAccessLogsEnabled(wg *sync.WaitGroup, s aws.Config, loadBalancers []types.LoadBalancer, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("ELB Access Logs Enabled", "Check if all load balancers have access logs enabled", testName)
	svc := elasticloadbalancingv2.NewFromConfig(s)
	// Get Load Balancers attributes
	for _, loadBalancer := range loadBalancers {
		input := &elasticloadbalancingv2.DescribeLoadBalancerAttributesInput{
			LoadBalancerArn: loadBalancer.LoadBalancerArn,
		}
		result, err := svc.DescribeLoadBalancerAttributes(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		for _, attribute := range result.Attributes {
			{
				if *attribute.Key == "access_logs.s3.enabled" && *attribute.Value == "true" {
					Message := "Access logs are enabled on : " + *loadBalancer.LoadBalancerName
					result := results.Result{Status: "OK", Message: Message, ResourceID: *loadBalancer.LoadBalancerArn}
					check.AddResult(result)
				} else if *attribute.Key == "access_logs.s3.enabled" && *attribute.Value == "false" {
					Message := "Access logs are not enabled on : " + *loadBalancer.LoadBalancerName
					result := results.Result{Status: "FAIL", Message: Message, ResourceID: *loadBalancer.LoadBalancerArn}
					check.AddResult(result)
				} else {
					continue
				}
			}
		}
	}

	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	loadBalancers := GetElasticLoadBalancers(s)
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)

	go yatas.CheckTest(&wg, c, "AWS_LB_001", CheckIfAccessLogsEnabled)(&wg, s, loadBalancers, "AWS_ELB_001", queueResults)
	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
