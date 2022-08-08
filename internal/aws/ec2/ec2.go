package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetEC2s(s *session.Session) []*ec2.Instance {
	svc := ec2.New(s)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		panic(err)
	}
	var instances []*ec2.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances
}

func CheckIfEC2PublicIP(s *session.Session, instances []*ec2.Instance, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "EC2 Public IP"
	check.Id = testName
	check.Description = "Check if all instances have a public IP"
	check.Status = "OK"
	for _, instance := range instances {
		if instance.PublicIpAddress != nil {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "EC2 instance " + *instance.InstanceId + " has a public IP" + *instance.PublicIpAddress
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *instance.InstanceId})
		} else {
			status := "OK"
			Message := "EC2 instance " + *instance.InstanceId + " has no public IP "
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *instance.InstanceId})
		}
	}
	*c = append(*c, check)
}

func RunEC2Tests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	instances := GetEC2s(s)
	config.CheckTest(c, "AWS_EC2_001", CheckIfEC2PublicIP)(s, instances, "AWS_EC2_001", &checks)
	return checks
}
