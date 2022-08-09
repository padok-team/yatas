package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetEC2s(s aws.Config) []types.Instance {
	svc := ec2.NewFromConfig(s)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	var instances []types.Instance
	for _, r := range result.Reservations {
		instances = append(instances, r.Instances...)
	}

	return instances
}

func CheckIfEC2PublicIP(s aws.Config, instances []types.Instance, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "EC2 Public IP"
	check.Id = testName
	check.Description = "Check if all instances have a public IP"
	check.Status = "OK"
	for _, instance := range instances {
		if instance.PublicIpAddress != nil {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "EC2 instance " + *instance.InstanceId + " has a public IP" + *instance.PublicIpAddress
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.InstanceId})
		} else {
			status := "OK"
			Message := "EC2 instance " + *instance.InstanceId + " has no public IP "
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.InstanceId})
		}
	}
	*c = append(*c, check)
}

func RunEC2Tests(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	instances := GetEC2s(s)
	yatas.CheckTest(c, "AWS_EC2_001", CheckIfEC2PublicIP)(s, instances, "AWS_EC2_001", &checks)
	return checks
}
