package loadbalancers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetElasticLoadBalancers(s *session.Session) []*elbv2.LoadBalancer {
	svc := elbv2.New(s)
	input := &elbv2.DescribeLoadBalancersInput{
		PageSize: aws.Int64(100),
	}
	result, err := svc.DescribeLoadBalancers(input)
	if err != nil {
		panic(err)
	}
	return result.LoadBalancers
}

func CheckIfAccessLogsEnabled(s *session.Session, loadBalancers []*elbv2.LoadBalancer, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "ELB Access Logs Enabled"
	check.Id = testName
	check.Description = "Check if all load balancers have access logs enabled"
	check.Status = "OK"
	svc := elbv2.New(s)
	// Get Load Balancers attributes
	for _, loadBalancer := range loadBalancers {
		input := &elbv2.DescribeLoadBalancerAttributesInput{
			LoadBalancerArn: loadBalancer.LoadBalancerArn,
		}
		result, err := svc.DescribeLoadBalancerAttributes(input)
		if err != nil {
			panic(err)
		}
		for _, attribute := range result.Attributes {
			{
				if *attribute.Key == "access_logs.s3.enabled" && *attribute.Value == "true" {
					status := "OK"
					Message := "Access logs are enabled on : " + *loadBalancer.LoadBalancerName
					check.Results = append(check.Results, types.Result{Status: status, Message: Message})
				} else if *attribute.Key == "access_logs.s3.enabled" && *attribute.Value == "false" {
					check.Status = "FAIL"
					status := "FAIL"
					Message := "Access logs are not enabled on : " + *loadBalancer.LoadBalancerName
					check.Results = append(check.Results, types.Result{Status: status, Message: Message})
				} else {
					continue
				}
			}
		}
	}

	*c = append(*c, check)
}

func RunLoadBalancersTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	loadBalancers := GetElasticLoadBalancers(s)
	config.CheckTest(c, "AWS_LB_001", CheckIfAccessLogsEnabled)(s, loadBalancers, "AWS_ELB_001", &checks)
	return checks
}
