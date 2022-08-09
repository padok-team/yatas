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

func CheckIfAccessLogsEnabled(wg *sync.WaitGroup, s aws.Config, loadBalancers []types.LoadBalancer, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "ELB Access Logs Enabled"
	check.Id = testName
	check.Description = "Check if all load balancers have access logs enabled"
	check.Status = "OK"
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
					status := "OK"
					Message := "Access logs are enabled on : " + *loadBalancer.LoadBalancerName
					check.Results = append(check.Results, results.Result{Status: status, Message: Message})
				} else if *attribute.Key == "access_logs.s3.enabled" && *attribute.Value == "false" {
					check.Status = "FAIL"
					status := "FAIL"
					Message := "Access logs are not enabled on : " + *loadBalancer.LoadBalancerName
					check.Results = append(check.Results, results.Result{Status: status, Message: Message})
				} else {
					continue
				}
			}
		}
	}

	*c = append(*c, check)
	wg.Done()
}

func RunChecks(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	loadBalancers := GetElasticLoadBalancers(s)
	var wg sync.WaitGroup

	go yatas.CheckTest(&wg, c, "AWS_LB_001", CheckIfAccessLogsEnabled)(&wg, s, loadBalancers, "AWS_ELB_001", &checks)
	wg.Wait()
	return checks
}
